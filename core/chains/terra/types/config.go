package types

import (
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	"github.com/smartcontractkit/chainlink/core/logger"
)

// Config extends terra.Config with an Update method.
type Config interface {
	terra.Config

	Update(dbcfg ChainCfg)
}

var _ Config = (*config)(nil)

type config struct {
	defaults terra.ConfigSet
	chain    ChainCfg
	chainMu  sync.RWMutex
	lggr     logger.Logger
}

// NewConfig returns a Config with defaults overridden by dbcfg.
func NewConfig(dbcfg ChainCfg, lggr logger.Logger) *config {
	return &config{
		defaults: terra.DefaultConfigSet,
		chain:    dbcfg,
		lggr:     lggr.Named("Config"),
	}
}

func (c *config) Update(dbcfg ChainCfg) {
	c.chainMu.Lock()
	c.chain = dbcfg
	c.chainMu.Unlock()
}

func (c *config) ConfirmMaxPolls() int64 {
	c.chainMu.RLock()
	ch := c.chain.ConfirmMaxPolls
	c.chainMu.RUnlock()
	if ch.Valid {
		return ch.Int64
	}
	return c.defaults.ConfirmMaxPolls
}

func (c *config) ConfirmPollPeriod() time.Duration {
	c.chainMu.RLock()
	ch := c.chain.ConfirmPollPeriod
	c.chainMu.RUnlock()
	if ch != nil {
		return ch.Duration()
	}
	return c.defaults.ConfirmPollPeriod
}

func (c *config) FallbackGasPriceULuna() sdk.Dec {
	c.chainMu.RLock()
	ch := c.chain.FallbackGasPriceULuna
	c.chainMu.RUnlock()
	if ch.Valid {
		str := ch.String
		dec, err := sdk.NewDecFromStr(str)
		if err == nil {
			return dec
		}
		c.lggr.Warnw("Invalid value provided for FallbackGasPriceULuna, falling back to default",
			"value", str, "default", c.defaults.FallbackGasPriceULuna, "err", err)
	}
	return c.defaults.FallbackGasPriceULuna
}

func (c *config) GasLimitMultiplier() float64 {
	c.chainMu.RLock()
	ch := c.chain.GasLimitMultiplier
	c.chainMu.RUnlock()
	if ch.Valid {
		return ch.Float64
	}
	return c.defaults.GasLimitMultiplier
}
