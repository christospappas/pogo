package server

import (
	"log"
)

type Socket struct {
	id            string
	client        *Client
	namespace     *Namespace
	subscriptions map[string]*Channel
}

func NewSocket(c *Client, ns *Namespace) *Socket {
	return &Socket{c.id, c, ns, make(map[string]*Channel)}
}

func (s *Socket) Send(msg *Message) {

}

func (s *Socket) Disconnect() {
	for _, ch := range s.subscriptions {
		ch.Unsubscribe(s)
	}
	s.namespace = nil
}

func (s *Socket) subscribe(chanName string) {
	if ch, ok := s.namespace.Channel(chanName); ok {
		log.Println("[pogo] Adding Client to Channel: " + chanName)
		ch.Subscribe(s)
	} else {
		log.Println("[pogo] Creating new Channel: " + chanName)
		ch := NewChannel(chanName)
		s.namespace.AddChannel(ch)
		ch.Subscribe(s)
	}
}

func (s *Socket) unsubscribe(chanName string) {

	if ch, ok := s.namespace.Channel(chanName); ok {
		ch.Unsubscribe(s)
	} else {
		log.Println("[pogo] Channel not found: " + chanName)
	}

}

func (s *Socket) OnData(msg *Message) {
	log.Println("[pogo] Command: " + msg.Event)

	switch msg.Event {
	case "pogo:subscribe":
		s.subscribe(msg.Channel)
	case "pogo:unsubscribe":
		s.unsubscribe(msg.Channel)
	default:
		if fn, ok := s.namespace.handlers[msg.Event]; ok {
			fn(msg, s.client)
		} else {
			log.Println("[pogo] No Event found: " + msg.Event)
		}
	}
}
