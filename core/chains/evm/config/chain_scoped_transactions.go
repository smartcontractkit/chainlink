package config

import (
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

type transactionsConfig struct {
	c toml.Transactions
}

func (t *transactionsConfig) ForwardersEnabled() bool {
	return *t.c.ForwardersEnabled
}

func (t *transactionsConfig) ReaperInterval() time.Duration {
	return t.c.ReaperInterval.Duration()
}

func (t *transactionsConfig) ReaperThreshold() time.Duration {
	return t.c.ReaperThreshold.Duration()
}

func (t *transactionsConfig) ResendAfterThreshold() time.Duration {
	return t.c.ResendAfterThreshold.Duration()
}

func (t *transactionsConfig) MaxInFlight() uint32 {
	return *t.c.MaxInFlight
}

func (t *transactionsConfig) MaxQueued() uint64 {
	return uint64(*t.c.MaxQueued)
}

func (g *transactionsConfig) AutoPurgeConfig() AutoPurgeConfig {
	return &autoPurgeConfig{c: g.c.AutoPurgeConfig}
}

type autoPurgeConfig struct {
	c toml.AutoPurgeConfig
}

func (a *autoPurgeConfig) AutoPurgeStuckTxs() bool {
	return *a.c.AutoPurgeStuckTxs
}

func (a *autoPurgeConfig) AutoPurgeThreshold() uint32 {
	return *a.c.AutoPurgeThreshold
}

func (a *autoPurgeConfig) AutoPurgeMinAttempts() uint32 {
	return *a.c.AutoPurgeMinAttempts
}

func (a *autoPurgeConfig) AutoPurgeDetectionApiUrl() *url.URL {
	return a.c.AutoPurgeDetectionApiUrl
}
