package config

// Deprecated: has no effect. Rate limiting is now configured via cresettings
// DispatcherGlobalRate/DispatcherPerSenderRate instead.
type DispatcherRateLimit interface {
	GlobalRPS() float64
	GlobalBurst() int
	PerSenderRPS() float64
	PerSenderBurst() int
}

type Dispatcher interface {
	SupportedVersion() int
	ReceiverBufferSize() int
	// Deprecated: has no effect. Rate limiting is now configured via cresettings
	// DispatcherGlobalRate/DispatcherPerSenderRate instead.
	RateLimit() DispatcherRateLimit
	SendToSharedPeer() bool
}
