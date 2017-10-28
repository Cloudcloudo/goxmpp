package service

import (
	"database/sql"
	"log"
	"strings"

	"goxmpp/model"

	"github.com/lib/pq"
)

func GetUserContacts(jid string, db *sql.DB) []model.RosterEntry {
	var entries []model.RosterEntry
	jid_sep := strings.Split(jid, "@")

	rows, err := db.Query("SELECT jid, subscrbed, nick, \"group\"  FROM contacts WHERE user_jid=$1", jid_sep[0])
	switch {
	case err == sql.ErrNoRows:
		log.Printf("User %s has empty contact list", jid)

	case err != nil:
		log.Printf("Error occured while reading  user %s contacts\n %s", jid, err)
	}

	for rows.Next() {
		var contact model.RosterEntry

		err = rows.Scan(&contact.Jid, &contact.Subscription, &contact.Name, pq.Array(&contact.Group))
		if err != nil {
			// handle this error
			panic(err)
		}
		entries = append(entries, contact)
	}

	defer rows.Close()

	return entries
}

func GetPreseceStanzaForUserContacts(user model.User, req model.ClientPresence, db *sql.DB) []model.ClientPresence {
	var entries []model.ClientPresence

	rows, err := db.Query("SELECT jid  FROM contacts WHERE user_jid=$1 AND ( subscrbed ='from' OR subscrbed='both')", user.LocalPart)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("User %s has empty contact list", user.LocalPart)

	case err != nil:
		log.Printf("Error occured while reading  user %s contacts\n %s", user.LocalPart, err)
	}

	for rows.Next() {
		presence := model.ClientPresence{
			From:   user.Jid,
			Status: req.Status,
			Show:   req.Show,
			Type:   req.Type,
		}

		err = rows.Scan(&presence.To)
		if err != nil {
			// handle this error
			panic(err)
		}

		entries = append(entries, presence)
	}

	defer rows.Close()

	return entries
}

func GetLocalContactsPresence(user model.User, db *sql.DB) []model.ClientPresence {
	var entries []model.ClientPresence

	query := `SELECT contacts.jid, users.presence, users.status
		FROM contacts
		INNER JOIN users ON contacts.jid = users.jid || '@' || $2
		WHERE contacts.user_jid=$1 AND ( contacts.subscrbed ='to' OR contacts.subscrbed='both')
	`
	rows, err := db.Query(query, user.LocalPart, user.DomainPart)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("User %s has empty contact list", user.GetFullJid())

	case err != nil:
		log.Printf("Error occured while reading  user %s contacts\n %s", user.GetFullJid(), err)
	}

	for rows.Next() {
		var contact model.ClientPresence

		err = rows.Scan(&contact.From, &contact.Show, &contact.Status)
		if err != nil {
			// handle this error
			panic(err)
		}
		if contact.Show == "unavailable" {
			contact.Show = ""
			contact.Type = "unavailable"
		}
		contact.To = user.GetFullJid()

		entries = append(entries, contact)
	}
	return entries
}
