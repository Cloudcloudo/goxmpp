package service

import (
	"database/sql"
	"log"
	"time"

	"goxmpp/model"
)

import "golang.org/x/crypto/bcrypt"

func AuthorizeUser(jid, password string, db *sql.DB) bool {

	var hashedPwd string

	err := db.QueryRow("SELECT password FROM users WHERE jid=$1", jid).Scan(&hashedPwd)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that jid %s", jid)
		return false

	case err != nil:
		log.Printf("Error occured checking user %s \n %s", jid, err)
		return false
	}

	byteHash := []byte(hashedPwd)
	err = bcrypt.CompareHashAndPassword(byteHash, []byte(password))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func AuthorizeUserToken(token string, db *sql.DB) (*model.User, error) {
	var jid string
	var tokenExpries time.Time

	err := db.QueryRow("SELECT jid, token_expiration  FROM users WHERE token=$1", token).Scan(&jid, &tokenExpries)

	switch {
	case err == sql.ErrNoRows:
		return nil, &model.TokenMissingError{Type: "invalid_token",Description: "Token "+ token + " missing"}

	case err != nil:
		log.Printf("Error occured checking user token %s \n %s", token, err)
		return nil, &model.TokenMissingError{Type: "unknown_error",Description: err.Error()}
	}


	if tokenExpries.Before(time.Now())  {
		return nil, &model.TokenExpiredError{Type: "token_expired",Description: "Token "+ token + " expired"}
	}

	user := model.NewUser(jid)

	return user, nil
}

func GenerateToken(jid string, db *sql.DB) (string, error) {
	token := RandString(48)

	expires := time.Now().UTC().AddDate(0, 0, 1)

	res, err := db.Exec("UPDATE users SET token=$1, token_expiration=$2 WHERE jid=$3", token, expires, jid)
	if err != nil {
		log.Printf("Error occured while update user token for user %s \n %s", jid, err)
		return "", err
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Printf("Error occured while update token for user %s \n %s", jid, err)
		return "", err
	}
	return token, nil
}

func GetUserLastSeen(user model.User, db *sql.DB) (*time.Time, error) {
	var lastSeen *time.Time

	err := db.QueryRow("SELECT last_seen FROM users WHERE jid=$1", user.LocalPart).Scan(&lastSeen)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that jid %s", user.GetFullJid())
		return nil, err

	case err != nil:
		log.Printf("Error occured checking user %s \n %s", user.GetFullJid(), err)
		return nil, err
	}

	return lastSeen, nil
}

func ChangeStatus(user model.User, req model.ClientPresence, db *sql.DB) (*model.ClientPresence, error) {
	var resp *model.ClientPresence

	if req.Type == "unavailable" {
		req.Show = req.Type
	}

	res, err := db.Exec("UPDATE users SET status=$1, presence=$2, last_seen=$4 WHERE jid=$3", req.Status, req.Show, user.LocalPart, time.Now().UTC())
	if err != nil {
		log.Printf("Error occured while update status for user %s \n %s", user.GetFullJid(), err)
		resp.Error = &model.Error{
			Type: "cancel",
			Text: "unable to set presence",
		}
		return resp, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Printf("Error occured while update status for user %s \n %s", user.GetFullJid(), err)
		resp.Error = &model.Error{
			Type: "cancel",
			Text: "unable to set presence",
		}
		return resp, err
	}

	return nil, nil
}

func SetUserOffline(user model.User, db *sql.DB) (*model.ClientPresence, error) {
	var resp *model.ClientPresence

	res, err := db.Exec("UPDATE users SET presence='unavailable', last_seen=$2 WHERE jid=$1", user.LocalPart, time.Now().UTC())
	if err != nil {
		log.Printf("Error occured while update status for user %s \n %s", user.GetFullJid(), err)
		resp.Error = &model.Error{
			Type: "cancel",
			Text: "unable to set presence",
		}
		return resp, err
	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Printf("Error occured while update status for user %s \n %s", user.GetFullJid(), err)
		resp.Error = &model.Error{
			Type: "cancel",
			Text: "unable to set presence",
		}
		return resp, err
	}

	return nil, nil
}

func CreateUser(u model.User, db *sql.DB) error {

	err := u.Validate()
	if err != nil {
		log.Printf("Invalid Jid", u.Jid, err)
		return err
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error occured encoding password for user %s \n %s", u.Jid, err)
		return err
	}

	_, err = db.Exec("INSERT INTO users (jid, password) VALUES ($1, $2)", u.Jid, hashedPwd)
	if err != nil {
		log.Printf("Error saving user %s \n %s", u.Jid, err)
		return err
	}

	return nil
}

func GetUserData(u model.User, db *sql.DB) (*model.UserData, error)  {
	query := `
		SELECT full_name, nick_name, birthdate,phone, www, email, company, 
		departament, "position", role, street, street_2, city, state, 
		zip_code, country, about, last_seen, presence, status, avatar
		FROM users 
		WHERE jid=$1
	`
	var userdata model.UserData
	err := db.QueryRow(query, u.Jid).Scan(
		&userdata.FullName, &userdata.NickName, &userdata.Birthday, &userdata.Phone, &userdata.Www, &userdata.Email, &userdata.Company,
		&userdata.Department, &userdata.Position, &userdata.Role, &userdata.Street, &userdata.Streer2, &userdata.City, &userdata.State,
		&userdata.ZipCode, &userdata.Country, &userdata.About, &userdata.LastSeen, &userdata.Presence, &userdata.Status, &userdata.Avatar)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that jid %s", u.Jid)
		return nil, err

	case err != nil:
		log.Printf("Error occured while get user %s data\n %s", u.Jid, err)
		return nil, err
	}
	fullJid := u.GetFullJid()
	userdata.Jid = &fullJid
	return &userdata, nil
}