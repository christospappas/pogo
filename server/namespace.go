package server

import (
	"sync"
)

// Namespace maintains the connected clients, channels and handlers
type Namespace struct {
	name        string
	sockets     map[string]*Socket
	channels    map[string]*Channel
	handlers    map[string]Handler
	toChannels  []string
	chanLock    *sync.Mutex
	socketsLock *sync.Mutex
}

// Creates new Namespace for specified pattern.
func NewNamespace(name string) *Namespace {
	sockets := make(map[string]*Socket)
	channels := make(map[string]*Channel)
	handlers := make(map[string]Handler)
	toChannels := make([]string, 0)
	return &Namespace{name, sockets, channels, handlers, toChannels, &sync.Mutex{}, &sync.Mutex{}}
}

func (ns *Namespace) AddClient(c *Client) *Socket {
	ns.socketsLock.Lock()
	defer ns.socketsLock.Unlock()

	socket := NewSocket(c, ns)
	ns.sockets[socket.id] = socket
	return socket
}

func (ns *Namespace) RemoveClient(s *Socket) {
	ns.socketsLock.Lock()
	defer ns.socketsLock.Unlock()

	delete(ns.sockets, s.id)
	s.client.Disconnect(ns.name)
}

func (ns *Namespace) On(event string, handler Handler) {
	ns.handlers[event] = handler
}

func (ns *Namespace) AddChannel(ch *Channel) {
	ns.chanLock.Lock()
	defer ns.chanLock.Unlock()
	ns.channels[ch.name] = ch
}

func (ns *Namespace) RemoveChannel(ch *Channel) {
	ns.chanLock.Lock()
	defer ns.chanLock.Unlock()
	delete(ns.channels, ch.name)
}

// returns a registered channel or an error if it isn't registered
func (ns *Namespace) Channel(name string) (ch *Channel, ok bool) {
	ns.chanLock.Lock()
	defer ns.chanLock.Unlock()

	ch, ok = ns.channels[name]
	return
}

// Broadcast sends a message to all sockets of the namespace or channel
func (ns *Namespace) Send(msg *Message, sender *Socket) {
	if len(ns.channels) > 0 {
		for ch := range ns.channels {
			ns.channels[ch].Broadcast(msg, sender)
		}
	}
}
