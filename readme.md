# GO-MEMBROKER

## Introduction

Membroker is a simple, lightweight in-memory message broker for Go applications than run on a single machine. It facilitates decoupling of application code using a pub/sub architecture, keeping components isolated and maintainable.

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
}

func subscriber(message []byte) {
	fmt.Printf("I got a message! - %v", string(message))
	c <- true
}
```
