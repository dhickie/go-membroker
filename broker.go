package membroker

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

// ErrRequestTimeout is sent if a request times out before we recieve a response
var ErrRequestTimeout = errors.New("The request timed out before a reply was recieved")

type broker struct {
	topics     map[string]*topic
	brokerLock sync.Mutex
}

func newBroker() *broker {
	return &broker{
		topics:     make(map[string]*topic),
		brokerLock: sync.Mutex{},
	}
}

func (b *broker) subscribe(topic string, callback func(Message)) int {
	top, ok := b.topics[topic]

	if !ok {
		// Lock the broker lock to make sure another thread hasn't just created this topic
		b.brokerLock.Lock()
		top, ok = b.topics[topic]
		if !ok {
			// If the topic doesn't exist yet, then make it
			top = newTopic(topic)
			b.topics[topic] = top
		}
		b.brokerLock.Unlock()
	}

	// Subscribe to the topic
	id := top.subscribe(callback)
	return id
}

func (b *broker) unsubscribe(topic string, id int) {
	top, ok := b.topics[topic]

	// If the topic doesn't exist, we don't need to do anything
	if !ok {
		b.brokerLock.Lock()
		top, ok = b.topics[topic]
		if !ok {
			b.brokerLock.Unlock()
			return
		}
		b.brokerLock.Unlock()
	}

	top.unsubscribe(id)
}

func (b *broker) publish(topic string, msg []byte) {
	message := Message{
		Data:  msg,
		Topic: topic,
	}

	// If the topic doesn't exist, then there can't be any subscribers -
	// just swallow the message
	if top, ok := b.topics[topic]; ok {
		top.publish(message)
	}
}

func (b *broker) request(topic string, msg []byte, timeout int) (Message, error) {
	// Generate a topic for the reply to be sent to
	// Topic + current time + first byte of message is a good way to ensure uniqueness
	reply := topic + strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(int(msg[0]))
	message := Message{
		Data:  msg,
		Topic: topic,
		Reply: reply,
	}

	// Subcribe the broker to the reply topic to recieve the reply
	c := make(chan Message)
	b.subscribe(reply, func(replyMsg Message) {
		c <- replyMsg
	})

	// Nothing else should be subscribed to this topic, so we can delete it once the reply is recieved
	// (or we timeout)
	defer delete(b.topics, reply)

	// If the topic doesn't exist, then there can't be any subscribers -
	// just swallow the message
	if top, ok := b.topics[topic]; ok {
		top.publish(message)
	}

	// Wait for either the reply, or the timeout
	ticker := time.NewTicker(time.Duration(timeout) * time.Millisecond)
	select {
	case <-ticker.C:
		return Message{}, ErrRequestTimeout
	case replyMsg := <-c:
		return replyMsg, nil
	}
}
