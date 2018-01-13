# GO-MEMBROKER

## Introduction

Membroker is a simple, lightweight in-memory message broker for Go applications than run on a single machine. It facilitates decoupling of application code using a pub/sub architecture, keeping components isolated and maintainable. It also supports a Request/Reply mechanism for the decoupled components.

Note that Membroker does not do any message queueing (and therefore no throttling or rate limiting) - any published messages are passed immediately to all subscribers, each in their own goroutine. For this reason, any subscribers to a topic **must** be threadsafe, as multiple messages could be sent at the same time.

## Installation

Get it:

```
go get github.com/dhickie/go-membroker
```

## Examples

```
package main

import (
	"fmt"

	"github.com/dhickie/go-membroker"
)

var c chan bool

func main() {
	c = make(chan bool)

	// Subscribe to the "example" topic and get a subscriber ID
	subID := membroker.Subscribe("example", subscriber)

	// Publish a message to the topic
	membroker.Publish("example", []byte("Hello :)"))

	// Wait for the message to be recieved
	<-c

	// Unsubscribe from the topic using the subscriber ID
	membroker.Unsubscribe("example", subID)

	// Subscribe the replier function
	membroker.Subscribe("request", replySubscriber)

	// Make a request with a timeout of 1000ms, and get a reply
	reply, _ := membroker.Request("request", []byte("I want a reply :("), 1000)
	fmt.Printf("Got a reply :) - %v\r\n", string(reply.Data))
}

func subscriber(message membroker.Message) {
	fmt.Printf("I got a message! - %v\r\n", string(message.Data))
	c <- true
}

func replySubscriber(message membroker.Message) {
	fmt.Printf("Replying to request! - %v\r\n", string(message.Data))
	membroker.Publish(message.Reply, []byte("Here's your reply :)"))
}
```

## Release notes

**1.1**:
- Added support for the request/reply architecture.
- Subscribers will now recieve a `Message` object instead of the raw binary.

**1.0**: 
- Initial release with support for basic publishing and subscribing to messages.