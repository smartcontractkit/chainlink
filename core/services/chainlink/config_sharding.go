package chainlink

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

const defaultArbiterPort = 9876
const defaultArbiterPollInterval = time.Second * 12
const defaultArbiterRetryInterval = time.Second * 12

var _ config.Sharding = (*shardingConfig)(nil)

type shardingConfig struct {
	s toml.Sharding
}

func (s *shardingConfig) ArbiterPort() uint16 {
	if s.s.ArbiterPort == nil {
		return defaultArbiterPort
	}
	return *s.s.ArbiterPort
}

func (s *shardingConfig) ArbiterPollInterval() time.Duration {
	if s.s.ArbiterPollInterval == nil || s.s.ArbiterPollInterval.Duration() <= 0 {
		return defaultArbiterPollInterval
	}
	return s.s.ArbiterPollInterval.Duration()
}

func (s *shardingConfig) ArbiterRetryInterval() time.Duration {
	if s.s.ArbiterRetryInterval == nil || s.s.ArbiterRetryInterval.Duration() <= 0 {
		return defaultArbiterRetryInterval
	}
	return s.s.ArbiterRetryInterval.Duration()
}
