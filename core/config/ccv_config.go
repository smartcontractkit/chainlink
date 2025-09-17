package config

type CCV interface {
	Enabled() bool
	ExecutorIndexerAPIKey() string
}
