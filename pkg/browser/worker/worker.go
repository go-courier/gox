package worker

type Worker interface {
	Receiver() <-chan interface{}
	PostMessage(v interface{})
	Close() error
}
