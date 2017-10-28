package webxmpp

import (
	"encoding/json"
	"errors"
	"reflect"

	"goxmpp/model"

	"github.com/gorilla/websocket"
)

type WsConn struct {
	Conn        *websocket.Conn
	StanzaTypes map[string]reflect.Type
}

func NewConn(connection *websocket.Conn, stanzaTypes map[string]reflect.Type) *WsConn {
	conn := &WsConn{
		Conn:        connection,
		StanzaTypes: stanzaTypes,
	}
	return conn
}

func (c *WsConn) ReadStanza(element model.WSMessage) (string, interface{}, error) {
	var stanzaInterface interface{}

	stanzaType, ok := c.StanzaTypes[element.Type]

	if ok {
		stanzaInterface = reflect.New(stanzaType).Interface()

		// skip unclosed Stream stanza
		switch stanzaInterface.(type) {
		case *model.Stream:
			return element.Type, stanzaInterface, nil
		}
	} else {
		return element.Type, nil, errors.New("Stanza not implemented " + element.Type)
	}

	err := json.Unmarshal(element.Payload, &stanzaInterface)
	if err != nil {
		return element.Type, nil, err
	}

	return element.Type, stanzaInterface, nil
}

func (c *WsConn) WriteStanza(msgtype string, stanza interface{}) error {
	data, err := json.Marshal(stanza)
	if err != nil {
		return err
	}

	wsMsg := &model.WSMessage{
		Type:msgtype,
		Payload: data,
	}

	err = c.Conn.WriteJSON(wsMsg)
	return err
}