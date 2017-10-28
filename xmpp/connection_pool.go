package xmpp

import (
	"sync"

	"goxmpp/model"
)

type ConnectionPool struct {
	acMutex sync.RWMutex
	Connections map[string]*model.Client
}

func NewConnectionPool() *ConnectionPool  {
	return &ConnectionPool{
		acMutex: sync.RWMutex{},
		Connections: make(map[string]*model.Client),
	}
}

func (c *ConnectionPool) CloseConnection(user *model.User)  {
	c.acMutex.Lock()
	delete(c.Connections, user.Jid)
	c.acMutex.Unlock()
}

func (c *ConnectionPool) AddConnection(client *model.Client)  {
	c.acMutex.Lock()
	c.Connections[client.User.Jid] = client
	c.acMutex.Unlock()
}

func (c *ConnectionPool) GetConnection(client *model.Client) (*model.Client, bool) {
	c.acMutex.RLock()
	conn, ok := c.Connections[client.User.Jid]
	c.acMutex.RUnlock()

	return conn, ok
}

func (c *ConnectionPool) GetConnectionViaJid(jid string) (*model.Client, bool) {
	c.acMutex.RLock()
	conn, ok := c.Connections[jid]
	c.acMutex.RUnlock()

	return conn, ok
}