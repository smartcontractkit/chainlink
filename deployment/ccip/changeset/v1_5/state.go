package v1_5

import (
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_contract"
)

// CCIPChainStateLegacy holds 1.5 state for CCIP chains
type CCIPChainStateLegacy struct {
	EVM2EVMOnRamp  map[uint64]*evm_2_evm_onramp.EVM2EVMOnRamp
	CommitStore    map[uint64]*commit_store.CommitStore
	EVM2EVMOffRamp map[uint64]*evm_2_evm_offramp.EVM2EVMOffRamp
	MockRMN        *mock_rmn_contract.MockRMNContract
	PriceRegistry  *price_registry_1_2_0.PriceRegistry
	RMN            *rmn_contract.RMNContract
}
