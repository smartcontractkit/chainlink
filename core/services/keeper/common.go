package keeper

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

var RegistryABI = eth.MustGetABI(keeper_registry_wrapper.KeeperRegistryABI)

type Config interface {
	KeeperDefaultTransactionQueueDepth() uint32
	KeeperGasPriceBufferPercent() uint32
	KeeperMaximumGracePeriod() int64
	KeeperMinimumRequiredConfirmations() uint64
	KeeperRegistryCheckGasOverhead() uint64
	KeeperRegistryPerformGasOverhead() uint64
	KeeperRegistrySyncInterval() time.Duration
}
