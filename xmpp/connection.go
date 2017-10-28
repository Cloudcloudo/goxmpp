package xmpp

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"reflect"

	"goxmpp/model"
)

// Connection represents server - client connection
type Conn struct {
	Conn        net.Conn
	StanzaTypes map[xml.Name]reflect.Type
	tx          *xml.Encoder
	rx          *xml.Decoder
}

// connection creator
func NewConn(connection net.Conn, stanzaTypes map[xml.Name]reflect.Type) *Conn {
	conn := &Conn{
		Conn:        connection,
		StanzaTypes: stanzaTypes,
		tx:          xml.NewEncoder(connection),
		rx:          xml.NewDecoder(connection),
	}
	return conn
}

// Next scans the stream to find the next xml.StartElement
func (c *Conn) Next() (xml.StartElement, error) {
	// loop until a start element token is found
	for {
		nextToken, err := c.rx.Token()
		if err != nil {
			return xml.StartElement{}, err
		}
		switch nextToken.(type) {
		case xml.StartElement:
			return nextToken.(xml.StartElement), nil
		}
	}
}

func (c *Conn) ReadStanza(element xml.StartElement) (xml.Name, interface{}, error) {
	var stanzaInterface interface{}

	stanzaType, ok := c.StanzaTypes[element.Name]

	if ok {
		stanzaInterface = reflect.New(stanzaType).Interface()

		// skip unclosed Stream stanza
		switch stanzaInterface.(type) {
		case *model.Stream:
			return element.Name, stanzaInterface, nil
		}
	} else {
		return xml.Name{}, nil, errors.New("Stanza not implemented " + element.Name.Space + " <" + element.Name.Local + "/>")
	}

	err := c.rx.DecodeElement(stanzaInterface, &element)

	if err != nil {
		return xml.Name{}, nil, err
	}

	return element.Name, stanzaInterface, nil
}

// WriteStanza decode stanza interface to xml and send to connection
func (c *Conn) WriteStanza(stanza interface{}) error {
	data, err := xml.Marshal(stanza)
	if err != nil {
		return err
	}
	_, err = c.Conn.Write(data)
	return err
}

// WriteStringf sends string via connection
func (c *Conn) WriteStringf(format string, a ...interface{}) error {
	_, err := fmt.Fprintf(c.Conn, format, a...)
	return err
}
