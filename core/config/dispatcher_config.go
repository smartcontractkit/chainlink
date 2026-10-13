package config

type DispatcherRateLimit interface {
	GlobalRPS() float64
	GlobalBurst() int
	PerSenderRPS() float64
	PerSenderBurst() int
}

type Dispatcher interface {
	ReceiverBufferSize() int
	RateLimit() DispatcherRateLimit
}
