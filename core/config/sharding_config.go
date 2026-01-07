package config

import "time"

type Sharding interface {
	ArbiterPort() uint16
	ArbiterPollInterval() time.Duration
	ArbiterRetryInterval() time.Duration
}
