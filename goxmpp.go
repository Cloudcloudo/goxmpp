package main

import (
	"database/sql"
	"log"

	"goxmpp/settings"
	"goxmpp/webxmpp"
	"goxmpp/xmpp"
)
var db *sql.DB

func main() {
	var err error

	db, err = sql.Open("postgres", settings.GetDBConnectionString())
	if err != nil {
		log.Fatalf("Unable to connect database %s\n", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Connected to database but %s\n", err)
	}

	connPool := xmpp.NewConnectionPool()

	go webxmpp.Start(db, connPool)
	xmpp.Start(db, connPool)
}
