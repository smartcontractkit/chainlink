package v1_5_1

import (
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	hybrid_external "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/hybrid_with_external_minter_fast_transfer_token_pool"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

// GroupUpdate represents a group update for a specific remote chain
type GroupUpdate struct {
	RemoteChainSelector uint64
	Group               uint8 // 0 = LOCK_AND_RELEASE, 1 = BURN_AND_MINT
	RemoteChainSupply   *big.Int
}

// UpdateGroupsInput defines the input for updating groups on hybrid token pools
type UpdateGroupsInput struct {
	GroupUpdates []GroupUpdate
}

var (
	// HybridWithExternalMinterTokenPoolUpdateGroupsOp updates groups on hybrid token pool contracts
	HybridWithExternalMinterTokenPoolUpdateGroupsOp = opsutil.NewEVMCallOperation(
		"HybridWithExternalMinterTokenPoolUpdateGroupsOp",
		semver.MustParse("1.0.0"),
		"Update groups on HybridWithExternalMinter token pool contract",
		hybrid_external.HybridWithExternalMinterFastTransferTokenPoolABI,
		shared.HybridWithExternalMinterFastTransferTokenPool,
		func(address common.Address, backend bind.ContractBackend) (interface{}, error) {
			return hybrid_external.NewHybridWithExternalMinterFastTransferTokenPool(address, backend)
		},
		func(pool interface{}, opts *bind.TransactOpts, input UpdateGroupsInput) (*types.Transaction, error) {
			hybridPool := pool.(*hybrid_external.HybridWithExternalMinterFastTransferTokenPool)

			// Convert our GroupUpdate struct to the contract's expected format
			var groupUpdates []hybrid_external.HybridTokenPoolAbstractGroupUpdate
			for _, update := range input.GroupUpdates {
				groupUpdates = append(groupUpdates, hybrid_external.HybridTokenPoolAbstractGroupUpdate{
					RemoteChainSelector: update.RemoteChainSelector,
					Group:               update.Group,
					RemoteChainSupply:   update.RemoteChainSupply,
				})
			}

			return hybridPool.UpdateGroups(opts, groupUpdates)
		},
	)
)
