package membroker_test

import (
	"testing"
	"time"

	"github.com/dhickie/go-membroker"
)

func TestSubscriberCanReceiveMessages(t *testing.T) {
	c := make(chan membroker.Message)
	membroker.Subscribe("test", func(msg membroker.Message) {
		c <- msg
	})

	membroker.Publish("test", []byte("message"))

	msg := <-c
	if string(msg.Data) != "message" {
		t.Error("Recieved incorrect message")
	}
}

func TestCorrectMessageIsReceivedOnCorrectTopic(t *testing.T) {
	c1 := make(chan membroker.Message)
	membroker.Subscribe("test1", func(msg membroker.Message) {
		c1 <- msg
	})

	c2 := make(chan membroker.Message)
	membroker.Subscribe("test2", func(msg membroker.Message) {
		c2 <- msg
	})

	membroker.Publish("test1", []byte("message1"))
	membroker.Publish("test2", []byte("message2"))

	msg1 := <-c1
	msg2 := <-c2

	if string(msg1.Data) != "message1" {
		t.Error("Recieved incorrect message on topic 1")
	}
	if string(msg2.Data) != "message2" {
		t.Error("Recieved incorrect message on topic 2")
	}
}

func TestCanUnsubscribeFromTopic(t *testing.T) {
	c := make(chan membroker.Message)
	subId := membroker.Subscribe("unsubTest", func(msg membroker.Message) {
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
	c1 := make(chan membroker.Message)
	membroker.Subscribe("test", func(msg membroker.Message) {
		c1 <- msg
	})

	c2 := make(chan membroker.Message)
	membroker.Subscribe("test", func(msg membroker.Message) {
		c2 <- msg
	})

	membroker.Publish("test", []byte("message"))

	msg1 := <-c1
	msg2 := <-c2

	if string(msg1.Data) != "message" {
		t.Error("Didn't receive correct message on first subscriber")
	}
	if string(msg2.Data) != "message" {
		t.Error("Didn't receive correct message on second subscriber")
	}
}

func TestRequestReply(t *testing.T) {
	membroker.Subscribe("test", func(msg membroker.Message) {
		membroker.Publish(msg.Reply, []byte("Reply to "+string(msg.Data)))
	})

	reply, err := membroker.Request("test", []byte("request"), 1000)
	if err != nil {
		t.Error("Received error when trying to get reply")
	} else if string(reply.Data) != "Reply to request" {
		t.Error("Received incorrect reply to request")
	}
}
