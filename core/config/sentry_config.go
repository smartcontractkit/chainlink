package config

type Sentry interface {
	DSN() string
	Debug() bool
	Environment() string
	Release() string
}
