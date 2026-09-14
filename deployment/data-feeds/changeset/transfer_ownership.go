package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	mcmschangesets "github.com/smartcontractkit/cld-changesets/legacy/mcms/changesets"
	evmstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/evm"
	mcmslib "github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

// TransferOwnershipChangeset is a changeset that will create an MCMS proposal to transfer
// the ownership of contracts currently owned by the timelock to NewOwnerAddress.
// Returns an MCMS proposal to transfer the ownership of contracts. Doesn't return a new addressbook.
// Ownership transfer is two-step: once the proposal is executed, the new owner must call acceptOwnership on each
// contract (e.g. via AcceptOwnershipChangeset if the new owner is another timelock) before the transfer takes effect.
var TransferOwnershipChangeset = cldf.CreateChangeSet(transferOwnershipLogic, transferOwnershipPrecondition)

func transferOwnershipLogic(env cldf.Environment, c types.TransferOwnershipConfig) (cldf.ChangesetOutput, error) {
	chain := env.BlockChains.EVMChains()[c.ChainSelector]

	var mcmsProposals []ProposalData
	for _, contractAddress := range c.ContractAddresses {
		_, contract, err := mcmschangesets.LoadOwnableContract(contractAddress, chain.Client)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to load the contract %w", err)
		}

		tx, err := contract.TransferOwnership(cldf.SimTransactOpts(), c.NewOwnerAddress)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transfer ownership tx %w", err)
		}
		mcmsProposals = append(mcmsProposals, ProposalData{
			contract:          contract.Address().Hex(),
			tx:                tx,
			timeLockQualifier: c.McmsConfig.TimeLockQualifier,
		})
	}

	proposalConfig := MultiChainProposalConfig{c.ChainSelector: mcmsProposals}
	proposal, err := BuildMultiChainProposals(env, "transfer ownership from timelock to "+c.NewOwnerAddress.Hex(), proposalConfig, c.McmsConfig.MinDelay)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

func transferOwnershipPrecondition(env cldf.Environment, c types.TransferOwnershipConfig) error {
	chain, ok := env.BlockChains.EVMChains()[c.ChainSelector]
	if !ok {
		return fmt.Errorf("chain not found in env %d", c.ChainSelector)
	}

	if c.McmsConfig == nil {
		return errors.New("mcms config is required")
	}

	if len(c.ContractAddresses) == 0 {
		return errors.New("at least one contract address is required")
	}

	if c.NewOwnerAddress == (common.Address{}) {
		return errors.New("new owner address cannot be the zero address")
	}

	if err := ValidateMCMSAddresses(env.DataStore.Addresses(), c.ChainSelector); err != nil {
		return err
	}

	mcmsState, err := evmstate.MaybeLoadMCMSWithTimelockStateWithQualifier(env, []uint64{c.ChainSelector}, c.McmsConfig.TimeLockQualifier)
	if err != nil {
		return fmt.Errorf("failed to load MCMS contracts for chain %d: %w", c.ChainSelector, err)
	}
	chainState, ok := mcmsState[c.ChainSelector]
	if !ok || chainState == nil || chainState.Timelock == nil {
		return fmt.Errorf("timelock with qualifier %q not found on chain %d", c.McmsConfig.TimeLockQualifier, c.ChainSelector)
	}
	timelockAddress := chainState.Timelock.Address()

	for _, contractAddress := range c.ContractAddresses {
		owner, _, err := mcmschangesets.LoadOwnableContract(contractAddress, chain.Client)
		if err != nil {
			return fmt.Errorf("failed to load the contract %s: %w", contractAddress.Hex(), err)
		}
		if owner == c.NewOwnerAddress {
			return fmt.Errorf("contract %s is already owned by %s", contractAddress.Hex(), owner.Hex())
		}
		if owner != timelockAddress {
			return fmt.Errorf("contract %s is owned by %s, not by the timelock %s", contractAddress.Hex(), owner.Hex(), timelockAddress.Hex())
		}
	}

	return nil
}
