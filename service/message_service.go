package service

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"goxmpp/model"
)

func SaveMessage(user model.User, msg model.ClientMessage, db *sql.DB) (*model.ClientMessage, error) {
	var resp *model.ClientMessage

	toJid := strings.Split(string(msg.To), "@")

	res, err := db.Exec(
		"INSERT INTO messages (user_jid, type, subject, nick, body, created_at, \"from\") VALUES ($1, $2, $3, $4, $5, $6, $7)",
		toJid[0], msg.Type, msg.Subject, msg.Nick, msg.Body, time.Now().UTC(), msg.From)

	if err != nil {
		log.Printf("Error saving for user %s \n %s", user.GetFullJid(), err)
		resp = &model.ClientMessage{
			ID:      msg.ID,
			To:      msg.To,
			From:    msg.From,
			Subject: msg.Subject,
			Body:    msg.Body,
			Type:    "error",
			Error: &model.Error{
				Type: "cancel",
				Text: "unable to send message",
			},
		}

		return resp, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Printf("Error occured while update status for user %s \n %s", user.GetFullJid(), err)
		resp = &model.ClientMessage{
			ID:      msg.ID,
			To:      msg.To,
			From:    msg.From,
			Subject: msg.Subject,
			Body:    msg.Body,
			Type:    "error",
			Error: &model.Error{
				Type: "cancel",
				Text: "unable to send message",
			},
		}
		return resp, err
	}
	return nil, nil
}

func GetMyMessages(user model.User, fromTime time.Time, db *sql.DB) (*[]model.ClientMessage, error) {
	var entries []model.ClientMessage

	rows, err := db.Query("SELECT id, \"type\", subject, nick, body, created_at, \"from\" FROM messages WHERE user_jid=$1 AND created_at>=$2", user.LocalPart, fromTime)
	switch {
	case err == sql.ErrNoRows:
		return &entries, nil

	case err != nil:
		log.Printf("Error occured while reading user %s messages\n %s", user.GetFullJid(), err)
		return &entries, err
	}

	for rows.Next() {
		var msg model.ClientMessage
		var msgTime time.Time

		err = rows.Scan(&msg.ID, &msg.Type, &msg.Subject, &msg.Nick, &msg.Body, &msgTime, &msg.From)
		if err != nil {
			// handle this error
			panic(err)
		}
		msg.Delay = &model.Delay{
			From:  msg.From,
			Stamp: msgTime.UTC().Format("2006-01-02T15:04:05Z"),
		}
		msg.To = user.GetFullJid()

		entries = append(entries, msg)
	}
	return &entries, nil
}
