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

func AddUserContacts(user model.User, req model.Iq, db *sql.DB) (*model.Iq, error) {
	var res sql.Result
	var err error

	for _, item := range req.Roster.Item {

		if item.Subscription == "" {

			row := db.QueryRow("SELECT jid  FROM contacts WHERE user_jid=$1 AND jid=$2", user.LocalPart, item.Jid)
			var jid string

			switch err = row.Scan(&jid); err {
			case sql.ErrNoRows:
				item.Subscription = "both"
				res, err = db.Exec(
					"INSERT INTO contacts (user_jid, jid, \"group\", nick, subscrbed) VALUES ($1, $2, $3, $4, $5)",
					user.LocalPart, item.Jid, pq.Array(item.Group), item.Name, item.Subscription)

			case nil:
				res, err = db.Exec(
					"UPDATE contacts SET \"group\"=$3, nick=$4 WHERE user_jid=$1 AND jid=$2",
					user.LocalPart, item.Jid, pq.Array(item.Group), item.Name)
			}

		} else if item.Subscription == "remove" {
			res, err = db.Exec(
				"DELETE FROM contacts WHERE user_jid=$1 AND jid=$2",
				user.LocalPart, item.Jid)
		}

		if err != nil {
			log.Printf("Error saving contact for user %s \n %s", user.GetFullJid(), err)
			resp := &model.Iq{
				ID:      req.ID,
				To:      req.To,
				From:    req.From,
				Type:    "error",
				Error: &model.Error{
					Type: "cancel",
					Text: "unable to save contact",
				},
			}

			return resp, err
		}
		_, err = res.RowsAffected()
		if err != nil {
			log.Printf("Error saving contact for user %s \n %s", user.GetFullJid(), err)
			resp := &model.Iq{
				ID:      req.ID,
				To:      req.To,
				From:    req.From,
				Type:    "error",
				Error: &model.Error{
					Type: "cancel",
					Text: "unable to save contact",
				},
			}

			return resp, err
		}
	}

	return nil, nil
}