package server

import (
	"sync"
)

// Namespace maintains the connected clients, channels and handlers
type Namespace struct {
	name        string
	sockets     []*Socket
	channels    map[string]*Channel
	handlers    map[string]Handler
	chanLock    *sync.Mutex
	socketsLock *sync.Mutex
}

// Creates new Namespace for specified pattern.
func NewNamespace(name string) *Namespace {
	sockets := make([]*Socket, 0)
	channels := make(map[string]*Channel)
	handlers := make(map[string]Handler)
	return &Namespace{name, sockets, channels, handlers, &sync.Mutex{}, &sync.Mutex{}}
}

func (ns *Namespace) Connect(c *Client) *Socket {
	ns.socketsLock.Lock()
	defer ns.socketsLock.Unlock()

	socket := NewSocket(c, ns)
	ns.sockets = append(ns.sockets, socket)
	c.sockets[ns.name] = socket
	return socket
}

func (ns *Namespace) On(event string, handler Handler) {
	ns.handlers[event] = handler
}

func (ns *Namespace) CreateChannel(name string) *Channel {
	ch := NewChannel(name)
	ns.channels[name] = ch
	return ch
}

// returns a registered channel or an error if it isn't registered
func (ns *Namespace) Channel(name string) (ch *Channel, ok bool) {
	ns.chanLock.Lock()
	defer ns.chanLock.Unlock()

	ch, ok = ns.channels[name]

	return
}
