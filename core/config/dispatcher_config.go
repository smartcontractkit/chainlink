package config

type Dispatcher interface {
	ReceiverBufferSize() int
}
