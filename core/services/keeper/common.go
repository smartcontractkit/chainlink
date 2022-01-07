package keeper

import (
	"time"

	eth "github.com/smartcontractkit/chainlink/core/chains/evm/eth"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
)

var RegistryABI = eth.MustGetABI(keeper_registry_wrapper.KeeperRegistryABI)

type Config interface {
	EvmEIP1559DynamicFees() bool
	KeeperDefaultTransactionQueueDepth() uint32
	KeeperGasPriceBufferPercent() uint32
	KeeperGasTipCapBufferPercent() uint32
	KeeperMaximumGracePeriod() int64
	KeeperRegistryCheckGasOverhead() uint64
	KeeperRegistryPerformGasOverhead() uint64
	KeeperRegistrySyncInterval() time.Duration
	KeeperRegistrySyncUpkeepQueueSize() uint32
	LogSQL() bool
}
