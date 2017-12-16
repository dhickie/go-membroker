package membroker

import (
	"sync"
)

type topic struct {
	name   string
	subs   map[int]*subscriber
	idLock sync.Mutex
	lastID int
}

func newTopic(name string) *topic {
	return &topic{
		name:   name,
		subs:   make(map[int]*subscriber),
		idLock: sync.Mutex{},
		lastID: 0,
	}
}

func (t *topic) subscribe(callback func([]byte)) int {
	id := t.getSubID()
	t.subs[id] = newSubscriber(id, callback)

	return id
}

func (t *topic) unsubscribe(id int) {
	delete(t.subs, id)
}

func (t *topic) publish(msg []byte) {
	// Publish the message to each subscriber
	for _, val := range t.subs {
		val.process(msg)
	}
}

func (t *topic) getSubID() int {
	t.idLock.Lock()
	defer t.idLock.Unlock()
	t.lastID++
	return t.lastID
}
