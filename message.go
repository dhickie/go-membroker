package membroker

// Message represents a message sent over membroker
type Message struct {
	Data  []byte
	Reply string
	Topic string
}
