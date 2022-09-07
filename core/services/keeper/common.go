package keeper

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/keeper_registry_wrapper1_3"
)

var Registry1_1ABI = types.MustGetABI(keeper_registry_wrapper1_1.KeeperRegistryABI)
var Registry1_2ABI = types.MustGetABI(keeper_registry_wrapper1_2.KeeperRegistryABI)
var Registry1_3ABI = types.MustGetABI(keeper_registry_wrapper1_3.KeeperRegistryABI)

type Config interface {
	EvmEIP1559DynamicFees() bool
	KeySpecificMaxGasPriceWei(addr common.Address) *big.Int
	KeeperDefaultTransactionQueueDepth() uint32
	KeeperGasPriceBufferPercent() uint32
	KeeperGasTipCapBufferPercent() uint32
	KeeperBaseFeeBufferPercent() uint32
	KeeperMaximumGracePeriod() int64
	KeeperRegistryCheckGasOverhead() uint32
	KeeperRegistryPerformGasOverhead() uint32
	KeeperRegistrySyncInterval() time.Duration
	KeeperRegistrySyncUpkeepQueueSize() uint32
	KeeperCheckUpkeepGasPriceFeatureEnabled() bool
	KeeperTurnLookBack() int64
	KeeperTurnFlagEnabled() bool
	LogSQL() bool
}

type RegistryGasChecker interface {
	KeeperRegistryCheckGasOverhead() uint32
	KeeperRegistryPerformGasOverhead() uint32
}
