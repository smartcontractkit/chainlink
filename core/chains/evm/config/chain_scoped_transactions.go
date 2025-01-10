package config

import (
	"net/url"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

type transactionsConfig struct {
	c toml.Transactions
}

func (t *transactionsConfig) Enabled() bool {
	return *t.c.Enabled
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

func (t *transactionsConfig) TransactionManagerV2() TransactionManagerV2 {
	return &transactionManagerV2Config{c: t.c.TransactionManagerV2}
}

type transactionManagerV2Config struct {
	c toml.TransactionManagerV2Config
}

func (t *transactionManagerV2Config) Enabled() bool {
	return *t.c.Enabled
}

func (t *transactionManagerV2Config) BlockTime() *time.Duration {
	d := t.c.BlockTime.Duration()
	return &d
}

func (t *transactionManagerV2Config) CustomURL() *url.URL {
	return t.c.CustomURL.URL()
}

func (t *transactionManagerV2Config) DualBroadcast() *bool {
	return t.c.DualBroadcast
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
