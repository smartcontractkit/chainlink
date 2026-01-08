package chainlink

import (
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

const defaultArbiterPort = 9876
const defaultArbiterPollInterval = time.Second * 12
const defaultArbiterRetryInterval = time.Second * 12
const defaultShardIndex = 0
const defaultShardOrchestratorPort = 50051

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

func (s *shardingConfig) ShardIndex() uint16 {
	if s.s.ShardIndex == nil {
		return defaultShardIndex
	}
	return *s.s.ShardIndex
}

func (s *shardingConfig) ShardOrchestratorPort() uint16 {
	if s.s.ShardOrchestratorPort == nil {
		return defaultShardOrchestratorPort
	}
	return *s.s.ShardOrchestratorPort
}

func (s *shardingConfig) ShardOrchestratorAddress() *url.URL {
	if s.s.ShardOrchestratorAddress == nil {
		return nil
	}
	return s.s.ShardOrchestratorAddress.URL()
}
