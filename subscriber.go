package membroker

type subscriber struct {
	id int
	fn func([]byte)
}

func newSubscriber(id int, callback func([]byte)) *subscriber {
	return &subscriber{
		id: id,
		fn: callback,
	}
}

func (s *subscriber) process(message []byte) {
	go s.fn(message)
}
