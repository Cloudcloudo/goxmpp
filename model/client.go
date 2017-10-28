package model

import "strings"

type Client struct {
	User  User
	Msgch chan interface{}
}


func NewClient(fullJid string) *Client {
	domainSplit := strings.Split(fullJid, "@")
	resourceSplit := strings.Split(fullJid, "/")

	var jid, domainpart, localpart string
	if len(domainSplit) > 1 {
		jid = resourceSplit[0]
		localpart = domainSplit[0]
		domainpart = strings.Split(domainSplit[1], "/")[0]
	} else {
		domainpart = resourceSplit[0]
	}

	var resourcepart string
	if len(resourceSplit) > 1 {
		resourcepart = resourceSplit[1]
	}

	return &Client{
		User: User{
			Jid:          jid,
			LocalPart:    localpart,
			DomainPart:   domainpart,
			ResourcePart: resourcepart,
		},
	}
}