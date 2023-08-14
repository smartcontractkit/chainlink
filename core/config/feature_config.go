package config

type Feature interface {
	FeedsManager() bool
	UICSAKeys() bool
	LogPoller() bool
	CCIP() bool
	LegacyGasStation() bool
}
