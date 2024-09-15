package config

type DispatcherRateLimit interface {
	GlobalRPS() float64
	GlobalBurst() int
	PerSenderRPS() float64
	PerSenderBurst() int
}

type Dispatcher interface {
	SupportedVersion() int
	ReceiverBufferSize() int
	RateLimit() DispatcherRateLimit
}
