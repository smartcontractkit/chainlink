package config

type Feature interface {
	FeedsManager() bool
	UICSAKeys() bool
	LogPoller() bool
	EVMClientV2() bool
}
