package webxmpp

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"goxmpp/model"
	"goxmpp/service"
	"goxmpp/settings"
	"goxmpp/xmpp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)
const MsgBuffer = 2

var conf = settings.Settings{}
var db *sql.DB
var connPool *xmpp.ConnectionPool
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Start(database *sql.DB, connectionPool *xmpp.ConnectionPool)  {
	db = database
	connPool = connectionPool

	conf = settings.GetSettings()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router := mux.NewRouter()
	router.HandleFunc("/login", loginUser).Methods("POST")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{jid}", getUser).Methods("GET")
	router.HandleFunc("/stream", initStream).Methods("GET").Queries("token", "{token}")

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("webxmpp/client/dist/"))))

	err := http.ListenAndServe(conf.ListenAddress + ":" + conf.WebPort, handlers.CORS(originsOk, headersOk, methodsOk)(router))
	if err != nil {
		log.Fatalf("Web server failed to start \n%s\n", err)
	}
}


func initStream(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Origin", "*")

	param := mux.Vars(req)
	token := param["token"]
	user, err := service.AuthorizeUserToken(token, db)
	if err != nil {
		switch e := err.(type) {
		case *model.TokenMissingError:
			resp, _ := json.Marshal(e)
			w.WriteHeader(http.StatusNotFound)
			w.Write(resp)

		case *model.TokenExpiredError:
			resp, _ := json.Marshal(e)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(resp)

		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	user.DomainPart = conf.Domain
	user.ResourcePart = "browserClient"

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	conn := NewConn(ws, model.WSMessageTypes)
	client := model.NewClient(user.GetFullJid())

	readCh := readMessages(conn)
	client.Msgch = make(chan interface{}, MsgBuffer)
	connPool.AddConnection(client)

	defer func() {
		log.Println("Connection closed", client.User.Jid, conn.Conn.RemoteAddr())
		close(readCh)
		service.SetUserOffline(client.User, db)
		conn.Conn.Close()

		connPool.CloseConnection(&client.User)
	}()

	for {
		select {
		case msg := <-client.Msgch: 	// receive messages from others clients
			var msgType string

			switch msg.(type) {
			case model.ClientMessage:
				msgType = "message"

			case model.ClientPresence:
				msgType = "presence"

			case model.Iq:
				msgType = "iq"

			default:
				log.Println("Unknown WEB message ", msg)
			}
			conn.WriteStanza(msgType, msg)


		case stanza := <-readCh:		// process my messages and stanzas
			switch s := stanza.(type) {
			case *model.ClientMessage:
				to := model.NewClient(s.To)

				s.From = client.User.GetFullJid()
				log.Println("msg from: ", s.From, " to: ", s.To)

				receiver, ok := connPool.GetConnection(to)

				if ok {
					receiver.Msgch <- *s
				} else {
					resp, err := service.SaveMessage(client.User, *s, db)
					if err != nil {
						conn.WriteStanza("message", resp)
						continue
					}
				}
			case *model.ClientPresence:
				resp, err := service.ChangeStatus(client.User, *s, db)

				if err != nil {
					conn.WriteStanza("presence", resp)
					continue
				}

				contactPresence := service.GetPreseceStanzaForUserContacts(client.User, *s, db)

				//acMutex.RLock()
				for _, presence := range contactPresence {
					receiver, ok := connPool.GetConnectionViaJid(presence.To)

					if ok {
						presence.To = ""
						receiver.Msgch <- presence
					}
				}
				//acMutex.RUnlock()

				if s.Type == "unavailable" {
					log.Println("bye! ", s.Status)
					break
				}
			case error:
				return

			default:
				log.Println(s)
			}
		}
	}
}
func readMessages(conn *WsConn) chan interface{} {
	msgCh := make(chan interface{})
	go func() {
		for {
			var wsMessage model.WSMessage
			err := conn.Conn.ReadJSON(&wsMessage)

			if err == io.EOF {
				break

			} else if err != nil {
				log.Println("Error reading stanza ", err)
				//msgCh <- err
				break
			}
			_, stanza, err := conn.ReadStanza(wsMessage)
			if err != nil {
				log.Println("Unable to read stanza ", err)
				//TODO add error response
			} else {
				msgCh <- stanza
			}
		}
	}()
	return msgCh
}

func getUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	param := mux.Vars(req)
	jid := param["jid"]

	err := db.QueryRow("SELECT jid FROM users WHERE jid=$1", jid).Scan(&jid)
	switch {
	case err == sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)

	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)

	default:
		w.WriteHeader(http.StatusOK)
	}
}


func createUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a new record.
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Unmarshal
	var user model.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = service.CreateUser(user, db)
	if err != nil {
		switch e := err.(type) {
		case *model.ValidationErrors:
			resp, _ := json.Marshal(e)

			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write(resp)

		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func loginUser(w http.ResponseWriter, req *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Create a new record.
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Unmarshal
	var user model.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.LocalPart = user.Jid
	user.DomainPart = conf.Domain

	if !service.AuthorizeUser(user.Jid, user.Password, db) {
		w.WriteHeader(http.StatusUnauthorized)
		return

	} else {
		userData, err := service.GetUserData(user, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := service.GenerateToken(user.Jid, db)
		userData.Token = &token
		resp, _ := json.Marshal(userData)

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}