package changeset

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"github.com/smartcontractkit/chainlink/deployment"
)

type ownershipAcceptor interface {
	AcceptOwnership(opts *bind.TransactOpts) (*gethtypes.Transaction, error)
	Address() common.Address
}

// NewAcceptOwnershipChangeset creates a changeset that accepts ownership of all the
// chain contracts on the given chainSelector.
// New chain contracts are:
// * OnRamp
// * OffRamp
// * FeeQuoter
// * NonceManager
// * RMNRemote
func NewAcceptOwnershipChangeset(
	e deployment.Environment,
	state CCIPOnChainState,
	chainSelector uint64,
) (deployment.ChangesetOutput, error) {
	chainState, ok := state.Chains[chainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("desired chain selector %d not found in onchain state", chainSelector)
	}

	var batch timelock.BatchChainOperation
	for _, contract := range []ownershipAcceptor{
		chainState.OnRamp,
		chainState.OffRamp,
		chainState.FeeQuoter,
		chainState.NonceManager,
		chainState.RMNRemote,
	} {
		acceptOwnershipTx, err := contract.AcceptOwnership(deployment.SimTransactOpts())
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate accept ownership calldata of %T: %w", contract, err)
		}

		batch.Batch = append(batch.Batch, mcms.Operation{
			To:    contract.Address(),
			Data:  acceptOwnershipTx.Data(),
			Value: big.NewInt(0),
		})
	}

	proposal, err := BuildProposalFromBatches(
		state,
		[]timelock.BatchChainOperation{batch},
		"Accept ownership of all CCIP chain contracts",
		time.Duration(0), // minDelay
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal from batch: %w, batch: %+v", err, batch)
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*proposal},
	}, nil
}
