package model

import (
	"time"

	"goxmpp/settings"
)

type User struct {
	Jid      string `json:"jid"`
	Password string `json:"password"`
	LocalPart    string
	DomainPart   string
	ResourcePart string
}

func NewUser(localPart string) *User {
	conf := settings.GetSettings()
	return &User{
		Jid: localPart + "@" + conf.Domain,
		LocalPart:localPart,
	}
}

func (c *User) GetFullJid() string {
	return c.LocalPart + "@" + c.DomainPart + "/" + c.ResourcePart
}

func (c *User) Validate() error {
	var errors ValidationErrors
	errors.AddError(JidValidator(c.Jid))
	errors.AddError(PasswdValidator(c.Password))

	return &errors
}


type UserData struct {
	Token		*string		`json:"token"`
	Jid         *string		`json:"jid"`
	FullName	*string		`json:"full_name"`
	NickName	*string		`json:"nick_name"`
	
	Company   *string    `json:"company"`
	Department*string    `json:"department"`
	Position  *string    `json:"position"`
	Role      *string    `json:"role"`
	Street    *string    `json:"street"`
	Streer2   *string    `json:"streer_2"`
	City      *string    `json:"city"`
	State     *string    `json:"state"`
	ZipCode   *string    `json:"zip_code"`
	Country   *string    `json:"country"`
	About     *string    `json:"about"`
	LastSeen   *time.Time `json:"last_seen"`
	Presence  *string    `json:"presence"`
	Status		*string	 `json:"status"`
	Email		*string		`json:"email"`
	Www			*string		`json:"www"`
	Phone		*string 		`json:"phone"`
	Birthday 	*time.Time 		`json:"birthday"`
	Avatar		*string		`json:"avatar"`
}