package model

import (
	"encoding/json"
	"reflect"
)

type WSMessage struct {
	Type string				`json:"type"`
	Payload  json.RawMessage `json:"payload"`
}

// StanzaTypes map of known stanza types
var WSMessageTypes = map[string]reflect.Type{
	"stream_error":    	reflect.TypeOf(StreamError{}),
	"stream":   		reflect.TypeOf(Stream{}),
	"auth":       		reflect.TypeOf(SaslAuth{}),
	"bind": 	        reflect.TypeOf(BindBind{}),
	"message":  		reflect.TypeOf(ClientMessage{}),
	"presence":			reflect.TypeOf(ClientPresence{}),
	"iq":		     	reflect.TypeOf(Iq{}),
	"error":    		reflect.TypeOf(Error{}),
}