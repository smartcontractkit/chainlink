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

func (t *transactionsConfig) AutoPurge() AutoPurgeConfig {
	return &autoPurgeConfig{c: t.c.AutoPurge}
}

type autoPurgeConfig struct {
	c toml.AutoPurgeConfig
}

func (a *autoPurgeConfig) Enabled() bool {
	return *a.c.Enabled
}

func (a *autoPurgeConfig) Threshold() *uint32 {
	return a.c.Threshold
}

func (a *autoPurgeConfig) MinAttempts() *uint32 {
	return a.c.MinAttempts
}

func (a *autoPurgeConfig) DetectionApiUrl() *url.URL {
	return a.c.DetectionApiUrl.URL()
}
