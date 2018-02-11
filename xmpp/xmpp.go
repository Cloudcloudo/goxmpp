package xmpp

import (
	"database/sql"
	"encoding/base64"
	"io"
	"log"
	"net"
	"strings"

	"goxmpp/model"
	"goxmpp/service"
	"goxmpp/settings"

	_ "github.com/lib/pq"
)

const MsgBuffer = 2

var db *sql.DB

type Server struct {
	Port   string
	Domain string
}

var conf = settings.Settings{}

var connPool *ConnectionPool

func Start(database *sql.DB, connectionPool *ConnectionPool) {
	db = database
	connPool = connectionPool

	conf = settings.GetSettings()

	ln, err := net.Listen("tcp", conf.ListenAddress + ":" + conf.XMPPPort)
	if err != nil {
		log.Fatalf("Unable to listen on address %s port %s %s\n", conf.ListenAddress, conf.XMPPPort, err)
		panic(err)
	}

	defer ln.Close()
	log.Print("Server up and running on ", ln.Addr().String())

	for {
		var err error

		connection, err := ln.Accept()
		if err != nil {
			log.Fatalf("Error occured while accepting connection from "+connection.LocalAddr().String()+"\n", err)
			panic(err)
		}

		conn := NewConn(connection, model.StanzaTypes)
		//go handleNewConn(conn)
		go authenticate(conn)
	}

	log.Println("Server closed")
}

func handleNewConn(conn *Conn) {
	// TODO upgrade to true TLS
	conn.WriteStringf("<stream:features><register xmlns='http://jabber.org/features/iq-register'><starttls xmlns=\"urn:ietf:params:xml:ns:xmpp-tls\"><required/></starttls></stream:features>")
	se, err := conn.Next()
	if err != nil {
		return
	}
	_, stanza, err := conn.ReadStanza(se)
	switch s := stanza.(type) {
	case *model.Starttls:
		log.Println("Client accepted TLS")
		conn.WriteStringf("<proceed xmlns=\"urn:ietf:params:xml:ns:xmpp-tls\"/><?xml version=\"1.0\"?>")

	default:
		log.Println("Not allowed stanza ", s)
	}

	authenticate(conn)
}

func authenticate(conn *Conn) {
	defer func() {
		err := conn.Conn.Close()
		if err == nil {
			log.Println("Connection closed ", conn.Conn.RemoteAddr())
		}
	}()

	conn.WriteStringf("<stream:stream id=\"%s\" version=\"1.0\" xmlns=\"jabber:client\" xmlns:stream=\"http://etherx.jabber.org/streams\">", service.RandString(10))
	conn.WriteStringf("<stream:features><mechanisms xmlns=\"urn:ietf:params:xml:ns:xmpp-sasl\"><mechanism>PLAIN</mechanism></mechanisms></stream:features>")

	for {
		se, err := conn.Next()
		if err != nil {
			return
		}
		_, stanza, err := conn.ReadStanza(se)
		switch s := stanza.(type) {
		case *model.SaslAuth:
			if s.Mechanism == "PLAIN" {
				data, err := base64.StdEncoding.DecodeString(s.Body)
				if err != nil {
					return
				}
				info := strings.Split(string(data), "\x00")

				if service.AuthorizeUser(info[1], info[2], db) {
					conn.WriteStringf("<success xmlns=\"urn:ietf:params:xml:ns:xmpp-sasl\"/>")

					// TODO refactor client struct assign

					var client = model.Client{
						User:			model.User{
							Jid:          info[1] + "@" + conf.Domain,
							LocalPart:    info[1],
							DomainPart:   conf.Domain,
							ResourcePart: "/",
						},
						Msgch:        make(chan interface{}, MsgBuffer),
					}
					log.Print("Client: ", client)

					connPool.AddConnection(&client)

					handleMessages(conn, &client)
					break

				} else {
					conn.WriteStringf("<failure xmlns=\"urn:ietf:params:xml:ns:xmpp-sasl\"><not-authorized/></failure>")
					break
				}
			}
		default:
			log.Println("Not allowed stanza ", s)
			continue
		}
	}
}

func handleMessages(conn *Conn, client *model.Client) {

	readCh := readMessages(conn)
	defer func() {
		log.Println("Connection closed", client.User.Jid, conn.Conn.RemoteAddr())
		close(readCh)
		service.SetUserOffline(client.User, db)
		conn.Conn.Close()

		connPool.CloseConnection(&client.User)
	}()

	for {
		select {
		case msg := <-client.Msgch: // receive messages
			conn.WriteStanza(msg)

		case stanza := <-readCh:
			switch s := stanza.(type) {
			case *model.Stream:
				conn.WriteStringf("<stream:stream id=\"%s\" version=\"1.0\" xmlns=\"jabber:client\" xmlns:stream=\"http://etherx.jabber.org/streams\">", service.RandString(10))
				conn.WriteStringf("<stream:features><bind xmlns=\"urn:ietf:params:xml:ns:xmpp-bind\"><optional/></bind><session xmlns=\"urn:ietf:params:xml:ns:xmpp-session\"><optional/></session></stream:features>")

			case *model.Iq:
				if s.Bind != nil && s.Type == "set" {
					client.User.ResourcePart = s.Bind.Resource
					respBind := model.BindBind{
						Jid: client.User.GetFullJid(),
					}
					resp := model.Iq{
						ID:   s.ID,
						Type: "result",
						Bind: &respBind,
					}
					conn.WriteStanza(resp)
				} else if s.Session != nil && s.Type == "set" {
					resp := model.Iq{
						ID:      s.ID,
						Type:    "result",
						Session: s.Session,
					}
					conn.WriteStanza(resp)
				} else if s.Roster != nil && s.Type == "get" {
					// Fetch contacts list request
					resp := model.Iq{
						ID:   s.ID,
						Type: "result",
						Roster: &model.Roster{
							Item: service.GetUserContacts(client.User.Jid, db),
						},
					}
					// send contacts
					conn.WriteStanza(resp)

					// Send all contacts presence
					presences := service.GetLocalContactsPresence(client.User, db)
					conn.WriteStanza(presences)

					lastSeen, err := service.GetUserLastSeen(client.User, db)
					if err != nil {
						continue
					}

					messages, err := service.GetMyMessages(client.User, *lastSeen, db)
					if err == nil {
						conn.WriteStanza(messages)
					}

				} else if s.Roster != nil && s.Type == "set" {
					//add contact
					resp, err := service.AddUserContacts(client.User, *s, db)

					if err != nil {
						conn.WriteStanza(resp)
					}

				} else if s.To != "" && s.Stream != nil {
					// File stream negotiation?
					to := model.NewClient(s.To)
					resp := model.Iq{
						From:   client.User.Jid,
						ID:     s.ID,
						Type:   s.Type,
						Error:  s.Error,
						Stream: s.Stream,
					}

					receiver, ok := connPool.GetConnection(to)

					if ok {
						receiver.Msgch <- resp
					} else {
						// TODO send error reply, save in db
						log.Println("Not sent")
					}
				} else if s.To != "" && s.QueryByteStream != nil {
					// File accept negotiation?
					to := model.NewClient(s.To)
					resp := model.Iq{
						From:            client.User.Jid,
						ID:              s.ID,
						Type:            s.Type,
						Error:           s.Error,
						QueryByteStream: s.QueryByteStream,
					}

					receiver, ok := connPool.GetConnection(to)

					if ok {
						receiver.Msgch <- resp
					} else {
						// TODO send error reply, save in db
						log.Println("Not sent")
					}
				} else if to := model.NewClient(s.To); s.To != "" && to.User.Jid != "" {
					// other iq requests to sb

					resp := model.Iq{
						From:            client.User.Jid,
						ID:              s.ID,
						Type:            s.Type,
						Error:           s.Error,
						QueryVersion:    s.QueryVersion,
						QueryByteStream: s.QueryByteStream,
					}

					receiver, ok := connPool.GetConnection(to)

					if ok {
						receiver.Msgch <- resp
					} else {
						log.Println("Not sent ")
					}
				} else if s.DiscItms != nil {
					// service discovery request
					to := model.NewClient(s.To)

					// client to our server
					if to.User.LocalPart == "" && to.User.DomainPart == conf.Domain {
						resp := model.Iq{
							ID:       s.ID,
							From:     conf.Domain,
							To:       client.User.GetFullJid(),
							Type:     "result",
							DiscItms: &model.DiscItms{Item: []model.RosterEntry{}},
						}
						conn.WriteStanza(resp)
					}
				} else if s.DiscInf != nil {
					// service discovery info request
					to := model.NewClient(s.To)

					// client to our server
					if to.User.LocalPart == "" && to.User.DomainPart == conf.Domain {
						resp := model.Iq{
							ID:      s.ID,
							From:    conf.Domain,
							To:      client.User.GetFullJid(),
							Type:    "result",
							DiscInf: &model.DiscInf{Item: ""},
						}
						conn.WriteStanza(resp)
					}
				} else if s.Ping != nil {
					// ping request
					to := model.NewClient(s.To)

					// client to our server
					if to.User.LocalPart == "" && to.User.DomainPart == conf.Domain {
						resp := model.Iq{
							ID:   s.ID,
							From: conf.Domain,
							To:   client.User.GetFullJid(),
							Type: "result",
						}
						conn.WriteStanza(resp)
					}
				} else {
					log.Println("Not implemented yet ", s)
				}

			case *model.ClientPresence:
				resp, err := service.ChangeStatus(client.User, *s, db)

				if err != nil {
					conn.WriteStanza(resp)
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
						conn.WriteStanza(resp)
						continue
					}
				}
			case error:
				return

			default:
				log.Println(s)
			}
		}
	}
}

func readMessages(conn *Conn) chan interface{} {
	stanCh := make(chan interface{})
	go func() {
		for {
			se, err := conn.Next()
			if err == io.EOF {
				break

			} else if err != nil {
				log.Println("Error reading stanza ", err)
				stanCh <- err
				break
			}
			_, stanza, err := conn.ReadStanza(se)
			if err != nil {
				log.Println("Unable to read stanza ", err)
				//TODO add error response
			} else {
				stanCh <- stanza
			}
		}
	}()
	return stanCh
}
