package server

import (
	"log"
	"sync"
	"time"
)

// Channel maintains a list of subscribers and is used
// as a hub to broadcast messages to clients
type Channel struct {
	name        string
	subscribers map[string]*Subscription
	*sync.Mutex
}

type Subscription struct {
	socket    *Socket
	createdAt time.Time
}

// Creates a new channel
func NewChannel(name string) *Channel {
	// TODO: validate channel name
	return &Channel{name, make(map[string]*Subscription), &sync.Mutex{}}
}

// Subscribe appends a client to the list of channel subscribers.
// It also adds a reference to the channel on the clients subscription list.
// This method is thread-safe.
func (ch *Channel) Subscribe(s *Socket) {
	ch.Lock()
	defer ch.Unlock()

	if _, present := ch.subscribers[s.client.id]; present {
		log.Println("[pogo] Client already subscribed to Channel: " + ch.name)
	} else {
		log.Println("[pogo] Client Subscribed to Channel: " + ch.name)
		ch.subscribers[s.client.id] = &Subscription{s, time.Now()}
		s.subscriptions[ch.name] = ch
	}
}

// Unsubscribe removes a client from the list of channel subscribers.
// It removes reference to channel from the clients subscription list.
// This method is thread-safe.
func (ch *Channel) Unsubscribe(s *Socket) {
	ch.Lock()
	defer ch.Unlock()

	if _, present := ch.subscribers[s.client.id]; present {
		log.Println("[pogo] Removing Client from Channel: " + ch.name)
		delete(ch.subscribers, s.client.id)
		delete(s.subscriptions, ch.name)
	} else {
		log.Println("[pogo] Client not found in channel: " + ch.name)
	}
}

// Subscribers returns a threadsafe list of subscribers to the channel.
// This method is thread-safe.
func (ch *Channel) Subscribers() map[string]*Subscription {
	ch.Lock()
	defer ch.Unlock()
	return ch.subscribers
}

// Returns the total number of subscribers to the channel
func (ch *Channel) SubscriberCount() int {
	return len(ch.subscribers)
}

// Broadcast sends a message to all subscribers of the channel.
func (ch *Channel) Broadcast(msg *Message, sender *Socket) {
	for _, subscriber := range ch.Subscribers() {
		if subscriber.socket != sender {
			subscriber.socket.Send(msg)
		}
	}
}
