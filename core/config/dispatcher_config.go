package config

type DispatcherRateLimit interface {
	GlobalRPS() float64
	GlobalBurst() int
	RPS() float64
	Burst() int
}

type Dispatcher interface {
	SupportedVersion() int
	ReceiverBufferSize() int
	RateLimit() DispatcherRateLimit
}
