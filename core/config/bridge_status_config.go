package config

import "time"

const MinimumPollingInterval = time.Minute

type BridgeStatusReporter interface {
	Enabled() bool
	StatusPath() string
	PollingInterval() time.Duration
	IgnoreInvalidBridges() bool
	IgnoreJoblessBridges() bool
}
