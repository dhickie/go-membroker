package membroker

import "sync"

var (
	b    *broker
	lock sync.Mutex
)

// Subscribe will subscribe the provided callback function to messages
// on the provided topic. It returns an ID of the subscriber, in case
// it wants to unsubscribe later.
func Subscribe(topic string, callback func(Message)) int {
	createBrokerIfNeeded()
	return b.subscribe(topic, callback)
}

// Unsubscribe will unsubcribe the subscriber with the provided ID from
// the provided topic.
func Unsubscribe(topic string, id int) {
	createBrokerIfNeeded()
	b.unsubscribe(topic, id)
}

// Publish publishes the provided message to subscribers on the provided topic
func Publish(topic string, msg []byte) {
	createBrokerIfNeeded()
	b.publish(topic, msg)
}

// Request publishes a message and then recieves a response
func Request(topic string, msg []byte, timeout int) (Message, error) {
	createBrokerIfNeeded()
	return b.request(topic, msg, timeout)
}

func createBrokerIfNeeded() {
	if b == nil {
		lock.Lock()
		if b == nil {
			b = newBroker()
		}
		lock.Unlock()
	}
}
