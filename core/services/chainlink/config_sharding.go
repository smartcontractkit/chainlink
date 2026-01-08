package chainlink

import (
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.Sharding = (*shardingConfig)(nil)

type shardingConfig struct {
	s toml.Sharding
}

func (s *shardingConfig) ArbiterPort() uint16 {
	return *s.s.ArbiterPort
}

func (s *shardingConfig) ArbiterPollInterval() time.Duration {
	return s.s.ArbiterPollInterval.Duration()
}

func (s *shardingConfig) ArbiterRetryInterval() time.Duration {
	return s.s.ArbiterRetryInterval.Duration()
}

func (s *shardingConfig) ShardIndex() uint16 {
	return *s.s.ShardIndex
}

func (s *shardingConfig) ShardOrchestratorPort() uint16 {
	return *s.s.ShardOrchestratorPort
}

func (s *shardingConfig) ShardOrchestratorAddress() *url.URL {
	if s.s.ShardOrchestratorAddress == nil {
		return nil
	}
	return s.s.ShardOrchestratorAddress.URL()
}
