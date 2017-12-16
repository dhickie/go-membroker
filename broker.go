package membroker

import "sync"

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

func (b *broker) subscribe(topic string, callback func([]byte)) int {
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

	// Subscribe to the topic, locking it first to ensure it doesn't get deleted before we subscribe
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
	// If the topic doesn't exist, then there can't be any subscribers -
	// just swallow the message
	if top, ok := b.topics[topic]; ok {
		top.publish(msg)
	}
}
