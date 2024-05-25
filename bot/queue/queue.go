package queue

type Queue interface {
	SendMessage(interface{}) error
	Close() error
}
