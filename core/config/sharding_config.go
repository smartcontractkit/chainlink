package config

import (
	"net/url"
	"time"
)

type Sharding interface {
	ArbiterPort() uint16
	ArbiterPollInterval() time.Duration
	ArbiterRetryInterval() time.Duration
	ShardIndex() uint16
	ShardOrchestratorPort() uint16
	ShardOrchestratorAddress() *url.URL
}
