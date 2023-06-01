package config

type Sentry interface {
	SentryDSN() string
	SentryDebug() bool
	SentryEnvironment() string
	SentryRelease() string
}
