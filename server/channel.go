package pogo

import (
	"log"
	"sync"
	"time"
)

type Channel struct {
	name        string
	subscribers map[string]Subscription
	*sync.Mutex
}

type Subscription struct {
	client    *Client
	createdAt time.Time
}

func NewChannel(name string) Channel {
	return Channel{name, make(map[string]Subscription), &sync.Mutex{}}
}

func (ch *Channel) Subscribe(c *Client) {
	ch.Lock()
	defer ch.Unlock()

	if _, present := ch.subscribers[c.id]; present {
		log.Println("[pogo] Client already subscribed to Channel: " + ch.name)
	} else {
		log.Println("[pogo] Client Subscribed to Channel: " + ch.name)
		ch.subscribers[c.id] = Subscription{c, time.Now()}
		c.subscriptions[ch.name] = ch
	}
}

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

func (ch *Channel) SubscriberCount() int {
	return len(ch.subscribers)
}

func (ch *Channel) Broadcast(msg *Message, sender *Client) {
	ch.Lock()
	defer ch.Unlock()

	for _, subscriber := range ch.subscribers {
		if subscriber.client != sender {
			subscriber.client.Send(msg)
		}
	}
}
