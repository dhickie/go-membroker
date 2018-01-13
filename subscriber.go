package membroker

type subscriber struct {
	id int
	fn func(Message)
}

func newSubscriber(id int, callback func(Message)) *subscriber {
	return &subscriber{
		id: id,
		fn: callback,
	}
}

func (s *subscriber) process(message Message) {
	go s.fn(message)
}
