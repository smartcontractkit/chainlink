package keeper

import (
	"time"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper1_2"
)

var Registry1_1ABI = evmtypes.MustGetABI(keeper_registry_wrapper1_1.KeeperRegistryABI)
var Registry1_2ABI = evmtypes.MustGetABI(keeper_registry_wrapper1_2.KeeperRegistryABI)

type Config interface {
	EvmEIP1559DynamicFees() bool
	KeeperDefaultTransactionQueueDepth() uint32
	KeeperGasPriceBufferPercent() uint32
	KeeperGasTipCapBufferPercent() uint32
	KeeperBaseFeeBufferPercent() uint32
	KeeperMaximumGracePeriod() int64
	KeeperRegistryCheckGasOverhead() uint64
	KeeperRegistryPerformGasOverhead() uint64
	KeeperRegistrySyncInterval() time.Duration
	KeeperRegistrySyncUpkeepQueueSize() uint32
	KeeperCheckUpkeepGasPriceFeatureEnabled() bool
	KeeperTurnLookBack() int64
	KeeperTurnFlagEnabled() bool
	LogSQL() bool
}
