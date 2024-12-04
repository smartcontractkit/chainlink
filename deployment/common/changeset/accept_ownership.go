package changeset

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type OwnershipAcceptor interface {
	AcceptOwnership(opts *bind.TransactOpts) (*gethtypes.Transaction, error)
	Address() common.Address
}

type AcceptOwnershipConfig struct {
	// OwnersPerChain is a mapping from chain selector to the owner contract address on that chain.
	OwnersPerChain map[uint64]common.Address

	// ProposerMCMSes is a mapping from chain selector to the proposer MCMS contract on that chain.
	ProposerMCMSes map[uint64]*gethwrappers.ManyChainMultiSig

	// Contracts is a mapping from chain selector to the ownership acceptors on that chain.
	// Proposal will be generated for these contracts.
	Contracts map[uint64][]OwnershipAcceptor

	// MinDelay is the minimum amount of time that must pass before the proposal
	// can be executed onchain.
	// This is typically set to 3 hours but can be set to 0 for immediate execution (useful for tests).
	MinDelay time.Duration
}

func (a AcceptOwnershipConfig) Validate() error {
	// check that we have owners and proposer mcmses for the chains
	// in the Contracts field.
	for chainSelector := range a.Contracts {
		if _, ok := a.OwnersPerChain[chainSelector]; !ok {
			return fmt.Errorf("missing owner for chain %d", chainSelector)
		}
		if _, ok := a.ProposerMCMSes[chainSelector]; !ok {
			return fmt.Errorf("missing proposer MCMS for chain %d", chainSelector)
		}
	}

	return nil
}

// type assertion - comply with deployment.ChangeSet interface
var _ deployment.ChangeSet[AcceptOwnershipConfig] = NewAcceptOwnershipChangeset

// NewAcceptOwnershipChangeset creates a changeset that contains a proposal to accept ownership of the contracts
// provided in the configuration.
func NewAcceptOwnershipChangeset(
	e deployment.Environment,
	cfg AcceptOwnershipConfig,
) (deployment.ChangesetOutput, error) {
	if err := cfg.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid accept ownership config: %w", err)
	}

	var batches []timelock.BatchChainOperation
	for chainSelector, ownershipAcceptors := range cfg.Contracts {
		var ops []mcms.Operation
		for _, ownershipAcceptor := range ownershipAcceptors {
			tx, err := ownershipAcceptor.AcceptOwnership(deployment.SimTransactOpts())
			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate accept ownership calldata of %T: %w", ownershipAcceptor, err)
			}

			ops = append(ops, mcms.Operation{
				To:    ownershipAcceptor.Address(),
				Data:  tx.Data(),
				Value: big.NewInt(0),
			})
		}
		batches = append(batches, timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(chainSelector),
			Batch:           ops,
		})
	}

	proposal, err := proposalutils.BuildProposalFromBatches(
		cfg.OwnersPerChain,
		cfg.ProposerMCMSes,
		batches,
		"Accept ownership of contracts",
		cfg.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal from batch: %w, batches: %+v", err, batches)
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*proposal},
	}, nil
}
