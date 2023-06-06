package chainlink

import (
	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
)

type sentryConfig struct {
	c v2.Sentry
}

func (s sentryConfig) DSN() string {
	return *s.c.DSN
}

func (s sentryConfig) Debug() bool {
	return *s.c.Debug
}

func (s sentryConfig) Environment() string {
	return *s.c.Environment
}

func (s sentryConfig) Release() string {
	return *s.c.Release
}
