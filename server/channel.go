package pogo

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
	client    *Client
	createdAt time.Time
}

// Creates a new channel
func NewChannel(name string) Channel {
	// TODO: validate channel name
	return Channel{name, make(map[string]*Subscription), &sync.Mutex{}}
}

// Subscribe appends a client to the list of channel subscribers.
// It also adds a reference to the channel on the clients subscription list.
// This method is thread-safe.
func (ch *Channel) Subscribe(c *Client) {
	ch.Lock()
	defer ch.Unlock()

	if _, present := ch.subscribers[c.id]; present {
		log.Println("[pogo] Client already subscribed to Channel: " + ch.name)
	} else {
		log.Println("[pogo] Client Subscribed to Channel: " + ch.name)
		ch.subscribers[c.id] = &Subscription{c, time.Now()}
		c.subscriptions[ch.name] = ch
	}
}

// Unsubscribe removes a client from the list of channel subscribers.
// It removes reference to channel from the clients subscription list.
// This method is thread-safe.
func (ch *Channel) Unsubscribe(c *Client) {
	ch.Lock()
	defer ch.Unlock()
	if _, present := ch.subscribers[c.id]; present {
		log.Println("[pogo] Removing Client from Channel: " + ch.name)
		delete(ch.subscribers, c.id)
		delete(c.subscriptions, ch.name)
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
// Sending is done in a separate goroutine
func (ch *Channel) Broadcast(msg *Message, sender *Client) {
	go func(subscribers map[string]*Subscription) {
		for _, subscriber := range subscribers {
			if subscriber.client != sender {
				subscriber.client.Send(msg)
			}
		}
		return
	}(ch.Subscribers())
}
