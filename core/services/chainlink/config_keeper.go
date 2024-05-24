package chainlink

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.Keeper = (*keeperConfig)(nil)

type registryConfig struct {
	c toml.KeeperRegistry
}

func (r *registryConfig) CheckGasOverhead() uint32 {
	return *r.c.CheckGasOverhead
}

func (r *registryConfig) PerformGasOverhead() uint32 {
	return *r.c.PerformGasOverhead
}

func (r *registryConfig) MaxPerformDataSize() uint32 {
	return *r.c.MaxPerformDataSize
}

func (r *registryConfig) SyncInterval() time.Duration {
	return r.c.SyncInterval.Duration()
}

func (r *registryConfig) SyncUpkeepQueueSize() uint32 {
	return *r.c.SyncUpkeepQueueSize
}

type keeperConfig struct {
	c toml.Keeper
}

func (k *keeperConfig) Registry() config.Registry {
	return &registryConfig{c: k.c.Registry}
}

func (k *keeperConfig) DefaultTransactionQueueDepth() uint32 {
	return *k.c.DefaultTransactionQueueDepth
}

func (k *keeperConfig) GasPriceBufferPercent() uint16 {
	return *k.c.GasPriceBufferPercent
}

func (k *keeperConfig) GasTipCapBufferPercent() uint16 {
	return *k.c.GasTipCapBufferPercent
}

func (k *keeperConfig) BaseFeeBufferPercent() uint16 {
	return *k.c.BaseFeeBufferPercent
}

func (k *keeperConfig) MaxGracePeriod() int64 {
	return *k.c.MaxGracePeriod
}

func (k *keeperConfig) TurnLookBack() int64 {
	return *k.c.TurnLookBack
}
