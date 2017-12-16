package membroker_test

import (
	"testing"
	"time"

	"github.com/dhickie/go-membroker"
)

func TestSubscriberCanReceiveMessages(t *testing.T) {
	c := make(chan []byte)
	membroker.Subscribe("test", func(msg []byte) {
		c <- msg
	})

	membroker.Publish("test", []byte("message"))

	msg := <-c
	if string(msg) != "message" {
		t.Error("Recieved incorrect message")
	}
}

func TestCorrectMessageIsReceivedOnCorrectTopic(t *testing.T) {
	c1 := make(chan []byte)
	membroker.Subscribe("test1", func(msg []byte) {
		c1 <- msg
	})

	c2 := make(chan []byte)
	membroker.Subscribe("test2", func(msg []byte) {
		c2 <- msg
	})

	membroker.Publish("test1", []byte("message1"))
	membroker.Publish("test2", []byte("message2"))

	msg1 := <-c1
	msg2 := <-c2

	if string(msg1) != "message1" {
		t.Error("Recieved incorrect message on topic 1")
	}
	if string(msg2) != "message2" {
		t.Error("Recieved incorrect message on topic 2")
	}
}

func TestCanUnsubscribeFromTopic(t *testing.T) {
	c := make(chan []byte)
	subId := membroker.Subscribe("unsubTest", func(msg []byte) {
		c <- msg
	})
	membroker.Unsubscribe("unsubTest", subId)

	membroker.Publish("unsubTest", []byte("message"))

	ticker := time.NewTicker(time.Second)
	select {
	case <-ticker.C:
	case <-c:
		t.Error("Incorrectly recieved message on unsubscribed topic")
	}
}

func TestUnsubscribingFromNonExistentTopic(t *testing.T) {
	membroker.Unsubscribe("doesntexist", 0)
}

func TestMultipleSubscribers(t *testing.T) {
	c1 := make(chan []byte)
	membroker.Subscribe("test", func(msg []byte) {
		c1 <- msg
	})

	c2 := make(chan []byte)
	membroker.Subscribe("test", func(msg []byte) {
		c2 <- msg
	})

	membroker.Publish("test", []byte("message"))

	msg1 := <-c1
	msg2 := <-c2

	if string(msg1) != "message" {
		t.Error("Didn't receive correct message on first subscriber")
	}
	if string(msg2) != "message" {
		t.Error("Didn't receive correct message on second subscriber")
	}
}
