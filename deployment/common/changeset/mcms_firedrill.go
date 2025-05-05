package changeset

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type FireDrillConfig struct {
	TimelockCfg proposalutils.TimelockConfig
	Selectors   []uint64
}

// buildNoOPEVM builds a dummy tx that transfers 0 to the RBACTimelock
func buildNoOPEVM(e deployment.Environment, selector uint64) (mcmstypes.Transaction, error) {
	chain, ok := e.Chains[selector]
	if !ok {
		return mcmstypes.Transaction{}, nil
	}
	//nolint:staticcheck  // need to migrate CCIP changesets so we can do it alongside this too
	addresses, err := e.ExistingAddresses.AddressesForChain(selector)
	if err != nil {
		return mcmstypes.Transaction{}, err
	}
	state, err := state.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	if err != nil {
		return mcmstypes.Transaction{}, err
	}

	// No-op: empty call to timelock (will hit the receive() function)
	tx := mcmsevmsdk.NewTransaction(
		state.Timelock.Address(),
		[]byte{},      // empty calldata
		big.NewInt(0), // no value
		"FireDrillNoop",
		nil,
	)
	return tx, nil
}

// buildNoOPSolana builds a dummy tx that calls the memo program
func buildNoOPSolana() (mcmstypes.Transaction, error) {
	contractID := solana.MemoProgramID
	memo := []byte("noop")

	// Create transaction
	tx, err := mcmssolanasdk.NewTransaction(
		contractID.String(),
		memo,
		big.NewInt(0),           // No lamports
		[]*solana.AccountMeta{}, // No account metas at the transaction level either
		"Memo",
		[]string{}, // Attach the no-op instruction
	)
	if err != nil {
		return mcmstypes.Transaction{}, err
	}

	return tx, nil
}

// MCMSSignFireDrillChangeset creates a changeset for a MCMS signing Fire Drill.
// It is used to make sure team member can effectively sign proposal and that the execution pipelines are healthy.
// The changeset will create a NO-OP transaction for each chain selector in the environment and create a proposal for it.
func MCMSSignFireDrillChangeset(e deployment.Environment, cfg FireDrillConfig) (deployment.ChangesetOutput, error) {
	allSelectors := cfg.Selectors
	if len(allSelectors) == 0 {
		solSelectors := e.AllChainSelectorsSolana()
		evmSelectors := e.AllChainSelectors()
		allSelectors = append(allSelectors, solSelectors...)
		allSelectors = append(allSelectors, evmSelectors...)
	}
	operations := make([]mcmstypes.BatchOperation, 0, len(allSelectors))
	timelocks := map[uint64]string{}
	mcmAddresses := map[uint64]string{}
	inspectors, err := proposalutils.McmsInspectors(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	for _, selector := range allSelectors {
		family, err := chainsel.GetSelectorFamily(selector)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		switch family {
		case chainsel.FamilyEVM:
			//nolint:staticcheck  // need to migrate CCIP changesets so we can do it alongside this too
			addresses, err := e.ExistingAddresses.AddressesForChain(selector)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			state, err := state.MaybeLoadMCMSWithTimelockChainState(e.Chains[selector], addresses)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			timelocks[selector] = state.Timelock.Address().String()
			mcmAddress, err := cfg.TimelockCfg.MCMBasedOnAction(*state)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			mcmAddresses[selector] = mcmAddress.Address().String()
			tx, err := buildNoOPEVM(e, selector)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			operations = append(operations, mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(selector),
				Transactions:  []mcmstypes.Transaction{tx},
			})
		case chainsel.FamilySolana:
			//nolint:staticcheck // need to migrate CCIP changesets so we can do it alongside this too
			addresses, err := e.ExistingAddresses.AddressesForChain(selector)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			state, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(e.SolChains[selector], addresses)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			timelocks[selector] = mcmssolanasdk.ContractAddress(state.TimelockProgram, mcmssolanasdk.PDASeed(state.TimelockSeed))
			mcmAddress, err := cfg.TimelockCfg.MCMBasedOnActionSolana(*state)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			mcmAddresses[selector] = mcmAddress
			tx, err := buildNoOPSolana()
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			operations = append(operations, mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(selector),
				Transactions:  []mcmstypes.Transaction{tx},
			})
		}
	}
	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmAddresses,
		inspectors,
		operations,
		"firedrill proposal",
		cfg.TimelockCfg)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	return deployment.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
	}, nil
}
