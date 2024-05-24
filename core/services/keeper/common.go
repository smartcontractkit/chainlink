package keeper

import (
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper1_3"
)

var Registry1_1ABI = types.MustGetABI(keeper_registry_wrapper1_1.KeeperRegistryABI)
var Registry1_2ABI = types.MustGetABI(keeper_registry_wrapper1_2.KeeperRegistryABI)
var Registry1_3ABI = types.MustGetABI(keeper_registry_wrapper1_3.KeeperRegistryABI)

type RegistryGasChecker interface {
	CheckGasOverhead() uint32
	PerformGasOverhead() uint32
	MaxPerformDataSize() uint32
}
