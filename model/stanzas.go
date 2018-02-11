package model

import "encoding/xml"
import "reflect"

const (
	// NsStream stream namesapce
	NsStream = "http://etherx.jabber.org/streams"

	// NsTLS xmpp-tls xml namespace
	NsTLS = "urn:ietf:params:xml:ns:xmpp-tls"

	// NsSASL xmpp-sasl xml namespace
	NsSASL = "urn:ietf:params:xml:ns:xmpp-sasl"

	// NsBind xmpp-bind xml namespace
	NsBind = "urn:ietf:params:xml:ns:xmpp-bind"

	// NsSession xmpp-session xml namespace
	NsSession = "urn:ietf:params:xml:ns:xmpp-session"

	// NsClient jabbet client namespace
	NsClient = "jabber:client"
)

// StreamError element
type StreamError struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams error"`
	Any     xml.Name `xml:",any"`
	Text    string   `xml:"text"`
}

type Stream struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams stream"`
	Lang    string   `xml:"lang,attr,omitempty"`
}

// RFC 3920  C.4  SASL name space

// SaslMechanisms element
type SaslMechanisms struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	Mechanism []string `xml:"mechanism"`
}

// SaslAuth element
type SaslAuth struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl auth"`
	Mechanism string   `xml:"mechanism,attr"`
	Body      string   `xml:",chardata"`
}

type Starttls struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls auth"`
}

// RFC 3920  C.5  Resource binding name space

// BindBind element
type BindBind struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string   `xml:"resource,omitempty"`
	Jid      string   `xml:"jid,omitempty"`
}

// XEP-0203: Delayed Delivery of <message/> and <presence/> stanzas.

// Delay element
type Delay struct {
	XMLName xml.Name `xml:"urn:xmpp:delay delay" json:"-"`
	From    string   `xml:"from,attr,omitempty" json:"from,omitempty"`
	Stamp   string   `xml:"stamp,attr,omitempty" json:"stamp,omitempty"`

	Body string `xml:",chardata" json:"body,omitempty"`
}

// RFC 3921  B.1  jabber:client

// ClientMessage element
type ClientMessage struct {
	XMLName xml.Name `xml:"jabber:client message" json:"-"`
	From    string   `xml:"from,attr,omitempty" json:"from"`
	ID      string   `xml:"id,attr" json:"id"`
	To      string   `xml:"to,attr" json:"to"`
	Type    string   `xml:"type,attr" json:"type"` // chat, error, groupchat, headline, or normal

	// These should technically be []clientText,
	// but string is much more convenient.
	Subject string `xml:"subject,omitempty" json:"subject,omitempty"`
	Body    string `xml:"body" json:"body"`
	Html 	*HtmlMsg `xml:"http://jabber.org/protocol/xhtml-im html,omitempty" json:"html,omitempty"`
	Thread  string `xml:"thread,omitempty" json:"thread,omitempty"`
	Delay   *Delay `xml:"delay,omitempty" json:"delay,omitempty"`
	Error   *Error `xml:"error,omitempty" json:"error,omitempty"`
	// STATES
	Composing *Composing `xml:"composing,omitempty" json:"composing,omitempty"`
	Active    *Active    `xml:"active,omitempty" json:"active,omitempty"`
	Starting  *Starting  `xml:"starting,omitempty" json:"starting,omitempty"`
	Paused    *Paused    `xml:"paused,omitempty" json:"paused,omitempty"`
	Inactive  *Inactive  `xml:"inactive,omitempty" json:"inactive,omitempty"`
	Gone      *Gone      `xml:"gone,omitempty" json:"gone,omitempty"`
	// END STATES
	Nick string `xml:"http://jabber.org/protocol/nick nick,omitempty" json:"nick,omitempty"`
}
type HtmlMsg struct {
	XMLName  xml.Name 	`xml:"http://jabber.org/protocol/xhtml-im html,omitempty"`
	Child []byte		`xml:",innerxml"`
}

type Starting struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates starting,omitempty" json:"-"`
	ID      string   `xml:"id,attr,omitempty" json:"id,omitempty"`
}
type Composing struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates composing,omitempty" json:"-"`
	ID      string   `xml:"id,attr,omitempty"`
}
type Active struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates active,omitempty" json:"-"`
	ID      string   `xml:"id,attr,omitempty" json:"id,omitempty"`
}
type Paused struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates paused,omitempty" json:"-"`
	ID      string   `xml:"id,attr,omitempty" json:"id,omitempty"`
}
type Inactive struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates inactive,omitempty" json:"-"`
	ID      string   `xml:"id,attr,omitempty" json:"id,omitempty"`
}
type Gone struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/chatstates gone,omitempty" json:"-"`
	ID      string   `xml:"id,attr,omitempty" json:"id,omitempty"`
}

// ClientPresence element
type ClientPresence struct {
	XMLName xml.Name `xml:"jabber:client presence" json:"-"`
	From    string   `xml:"from,attr,omitempty" json:"from,omitempty"`
	ID      string   `xml:"id,attr,omitempty" json:"id,omitempty"`
	To      string   `xml:"to,attr,omitempty" json:"to,omitempty"`
	Type    string   `xml:"type,attr,omitempty" json:"type,omitempty"` // error, probe, subscribe, subscribed, unavailable, unsubscribe, unsubscribed
	Lang    string   `xml:"lang,attr,omitempty" json:"lang,omitempty"`

	Show     string `xml:"show,omitempty" json:"show,omitempty"`   // av, away, chat, dnd, xa
	Status   string `xml:"status,omitempty" json:"status,omitempty"` // sb []clientText
	Priority string `xml:"priority,omitempty" json:"priority,omitempty"`
	Caps     *Caps  `xml:"c" json:"caps,omitempty"`
	Error    *Error `xml:"error" json:"error,omitempty"`
	Delay    *Delay `xml:"delay" json:"delay,omitempty"`
}

// ClientCaps element
type Caps struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/caps c" json:"-"`
	Ext     string   `xml:"ext,attr" json:"ext,omitempty"`
	Hash    string   `xml:"hash,attr" json:"hash,omitempty"`
	Node    string   `xml:"node,attr" json:"node,omitempty"`
	Ver     string   `xml:"ver,attr" json:"ver,omitempty"`
}

// Error element
type Error struct {
	XMLName xml.Name `xml:"jabber:client error" json:"-"`
	Code    string   `xml:"code,attr,omitempty" json:"code,omitempty"`
	Type    string   `xml:"type,attr,omitempty" json:"type,omitempty"`
	Text  string `xml:"text,omitempty" json:"text,omitempty"`
	Child []byte `xml:",innerxml" json:"-"`
}

// Iq element
type Iq struct {
	// info/query
	XMLName         xml.Name         `xml:"jabber:client iq" json:"-"`
	From            string           `xml:"from,attr,omitempty" json:"from,omitempty"`
	ID              string           `xml:"id,attr" json:"id"`
	To              string           `xml:"to,attr,omitempty" json:"to,omitempty"`
	Type            string           `xml:"type,attr" json:"type"` // error, get, result, set
	Error           *Error           `xml:"error,omitempty" json:"error,omitempty"`
	Bind            *BindBind        `xml:"bind,omitempty" json:"bind,omitempty"`
	Roster          *Roster          `xml:"jabber:iq:roster query,omitempty" json:"query-roster,omitempty"`
	Session         *struct{}        `xml:"urn:ietf:params:xml:ns:xmpp-session session,omitempty" json:"session,omitempty"`
	Stream          *Si              `xml:"si,omitempty" json:"si,omitempty"`
	Ping            *struct{}        `xml:"urn:xmpp:ping ping,omitempty" json:"ping,omitempty"`
	QueryByteStream *QueryByteStream `xml:"http://jabber.org/protocol/bytestreams query,omitempty" json:"query-byte-stream,omitempty"`
	QueryVersion    *QueryVersion    `xml:"jabber:iq:version query,omitempty" json:"query-version,omitempty"`
	DiscItms        *DiscItms        `xml:"http://jabber.org/protocol/disco#items query" json:"disc-itms,omitempty"`
	DiscInf         *DiscInf         `xml:"http://jabber.org/protocol/disco#info query" json:"disc-inf,omitempty"`
}

type DiscItms struct {
	XMLName xml.Name      `xml:"http://jabber.org/protocol/disco#items query"`
	Item    []RosterEntry `xml:"item,omitempty"`
}
type DiscInf struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#info query"`
	Item    string   `xml:"item,omitempty"`
}

type QueryVersion struct {
	XMLName xml.Name `xml:"jabber:iq:version query" json:"-"`
	Name    string   `xml:"name,omitempty" json:"name"`
	Version string   `xml:"version,omitempty" json:"version"`
	Os      string   `xml:"os,omitempty" json:"os"`
}

// Si protocol element
type Si struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/si si"`
	ID      string   `xml:"id,attr"`
	Profile string   `xml:"profile,attr,omitempty"`
	Child   []byte   `xml:",innerxml"`
}

type QueryByteStream struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/bytestreams query"`
	SID     string   `xml:",sid,attr,omitempty"`
	Mode    string   `xml:",mode,attr,omitempty"`
	Child   []byte   `xml:",innerxml"`
}

// Roster element
type Roster struct {
	XMLName xml.Name      `xml:"jabber:iq:roster query" json:"-"`
	Item    []RosterEntry `xml:"item" json:"item"`
}

// RosterEntry element
type RosterEntry struct {
	Jid          string   `xml:"jid,attr" json:"jid"`
	Subscription string   `xml:"subscription,attr" json:"subscription"` // 'none', 'to', 'both', 'from', 'remove'
	Name         string   `xml:"name,attr" json:"name"`
	Group        []string `xml:"group" json:"group"`
}

// RosterRequest is used to request that the server update the user's roster.
// See RFC 6121, section 2.3.
type RosterRequest struct {
	XMLName xml.Name          `xml:"jabber:iq:roster query"`
	Item    RosterRequestItem `xml:"item"`
}

// starttls
type starttls struct{}

// RosterRequestItem element
type RosterRequestItem struct {
	Jid          string   `xml:"jid,attr" json:"jid"`
	Subscription string   `xml:"subscription,attr" json:"subscription"`
	Name         string   `xml:"name,attr" json:"name"`
	Group        []string `xml:"group" json:"group"`
}

// StanzaTypes map of known stanza types
var StanzaTypes = map[xml.Name]reflect.Type{
	xml.Name{Space: NsStream, Local: "error"}:    reflect.TypeOf(StreamError{}),
	xml.Name{Space: NsStream, Local: "stream"}:   reflect.TypeOf(Stream{}),
	xml.Name{Space: NsSASL, Local: "auth"}:       reflect.TypeOf(SaslAuth{}),
	xml.Name{Space: NsSASL, Local: "mechanisms"}: reflect.TypeOf(SaslMechanisms{}),
	xml.Name{Space: NsSASL, Local: "challenge"}:  reflect.TypeOf(""),
	xml.Name{Space: NsTLS, Local: "starttls"}:    reflect.TypeOf(starttls{}),
	xml.Name{Space: NsSASL, Local: "response"}:   reflect.TypeOf(""),
	xml.Name{Space: NsBind, Local: "bind"}:       reflect.TypeOf(BindBind{}),
	xml.Name{Space: NsClient, Local: "message"}:  reflect.TypeOf(ClientMessage{}),
	xml.Name{Space: NsClient, Local: "presence"}: reflect.TypeOf(ClientPresence{}),
	xml.Name{Space: NsClient, Local: "iq"}:       reflect.TypeOf(Iq{}),
	xml.Name{Space: NsClient, Local: "error"}:    reflect.TypeOf(Error{}),
}
