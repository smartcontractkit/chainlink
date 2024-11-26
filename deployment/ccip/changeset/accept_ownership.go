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

type AcceptOwnershipConfig struct {
	State             CCIPOnChainState
	ChainSelectors    []uint64
	HomeChainSelector uint64
}

// type assertion - comply with deployment.ChangeSet interface
var _ deployment.ChangeSet[AcceptOwnershipConfig] = NewAcceptOwnershipChangeset

// NewAcceptOwnershipChangeset creates a changeset that accepts ownership of all the
// ccip contracts deployed on the given chain selectors.
// New chain contracts are:
// * OnRamp
// * OffRamp
// * FeeQuoter
// * NonceManager
// * RMNRemote
// Home chain contracts are:
// * CCIPHome
// * RMNHome
// * CapabilityRegistry
func NewAcceptOwnershipChangeset(
	e deployment.Environment,
	cfg AcceptOwnershipConfig,
) (deployment.ChangesetOutput, error) {
	// basic validation
	if len(cfg.ChainSelectors) == 0 || cfg.HomeChainSelector == 0 {
		return deployment.ChangesetOutput{}, fmt.Errorf("no chain selectors provided")
	}

	if len(cfg.State.Chains) == 0 {
		return deployment.ChangesetOutput{}, fmt.Errorf("no chains in state")
	}

	// gen one batch per chain for the chain contracts
	var batches []timelock.BatchChainOperation
	for _, chainSelector := range cfg.ChainSelectors {
		chainState, ok := cfg.State.Chains[chainSelector]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain %d not found in state", chainSelector)
		}

		ops, err := genAcceptOwnershipOps(chainState)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate accept ownership batch for chain %d: %w",
				chainSelector, err)
		}

		batches = append(batches, timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(chainSelector),
			Batch:           ops,
		})
	}

	// gen separate batch chain operation for home chain contracts
	homeChainState, ok := cfg.State.Chains[cfg.HomeChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("home chain %d not found in state", cfg.HomeChainSelector)
	}
	homeChainOps, err := genHomeChainAcceptOwnershipOps(homeChainState)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to generate accept ownership batch for home chain %d: %w",
			cfg.HomeChainSelector, err)
	}

	batches = append(batches, timelock.BatchChainOperation{
		ChainIdentifier: mcms.ChainIdentifier(cfg.HomeChainSelector),
		Batch:           homeChainOps,
	})

	proposal, err := BuildProposalFromBatches(
		cfg.State,
		batches,
		"Accept ownership of all CCIP contracts",
		time.Duration(0), // minDelay
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal from batch: %w, batches: %+v", err, batches)
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*proposal},
	}, nil
}

func genAcceptOwnershipOps(chainState CCIPChainState) (ops []mcms.Operation, err error) {
	for _, contract := range []ownershipAcceptor{
		chainState.OnRamp,
		chainState.OffRamp,
		chainState.FeeQuoter,
		chainState.NonceManager,
		chainState.RMNRemote,
	} {
		acceptOwnershipTx, err := contract.AcceptOwnership(deployment.SimTransactOpts())
		if err != nil {
			return nil, fmt.Errorf("failed to generate accept ownership calldata of %T: %w", contract, err)
		}

		ops = append(ops, mcms.Operation{
			To:    contract.Address(),
			Data:  acceptOwnershipTx.Data(),
			Value: big.NewInt(0),
		})
	}

	return ops, nil
}

func genHomeChainAcceptOwnershipOps(homeChainState CCIPChainState) (ops []mcms.Operation, err error) {
	for _, contract := range []ownershipAcceptor{
		homeChainState.CapabilityRegistry,
		homeChainState.CCIPHome,
		homeChainState.RMNHome,
	} {
		acceptOwnershipTx, err := contract.AcceptOwnership(deployment.SimTransactOpts())
		if err != nil {
			return nil, fmt.Errorf("failed to generate accept ownership calldata of %T: %w", contract, err)
		}

		ops = append(ops, mcms.Operation{
			To:    contract.Address(),
			Data:  acceptOwnershipTx.Data(),
			Value: big.NewInt(0),
		})
	}

	return ops, nil
}
