package keeper

import (
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_contract"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

var RegistryABI = eth.MustGetABI(keeper_registry_contract.KeeperRegistryContractABI)
