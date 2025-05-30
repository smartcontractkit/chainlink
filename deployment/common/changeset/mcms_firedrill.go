package changeset

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

type FireDrillConfig struct {
	TimelockCfg proposalutils.TimelockConfig
	Selectors   []uint64
}

// buildNoOPEVM builds a dummy tx that transfers 0 to the RBACTimelock
func buildNoOPEVM(e cldf.Environment, selector uint64) (mcmstypes.Transaction, error) {
	chain, ok := e.BlockChains.EVMChains()[selector]
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
func MCMSSignFireDrillChangeset(e cldf.Environment, cfg FireDrillConfig) (cldf.ChangesetOutput, error) {
	allSelectors := cfg.Selectors
	if len(allSelectors) == 0 {
		solSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilySolana))
		evmSelectors := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chainsel.FamilyEVM))
		allSelectors = append(allSelectors, solSelectors...)
		allSelectors = append(allSelectors, evmSelectors...)
	}
	operations := make([]mcmstypes.BatchOperation, 0, len(allSelectors))
	timelocks := map[uint64]string{}
	mcmAddresses := map[uint64]string{}
	inspectors, err := proposalutils.McmsInspectors(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	for _, selector := range allSelectors {
		family, err := chainsel.GetSelectorFamily(selector)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		switch family {
		case chainsel.FamilyEVM:
			//nolint:staticcheck  // need to migrate CCIP changesets so we can do it alongside this too
			addresses, err := e.ExistingAddresses.AddressesForChain(selector)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			state, err := state.MaybeLoadMCMSWithTimelockChainState(e.BlockChains.EVMChains()[selector], addresses)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			timelocks[selector] = state.Timelock.Address().String()
			mcmAddress, err := cfg.TimelockCfg.MCMBasedOnAction(*state)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			mcmAddresses[selector] = mcmAddress.Address().String()
			tx, err := buildNoOPEVM(e, selector)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			operations = append(operations, mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(selector),
				Transactions:  []mcmstypes.Transaction{tx},
			})
		case chainsel.FamilySolana:
			//nolint:staticcheck // need to migrate CCIP changesets so we can do it alongside this too
			addresses, err := e.ExistingAddresses.AddressesForChain(selector)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			state, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(e.BlockChains.SolanaChains()[selector], addresses)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			timelocks[selector] = mcmssolanasdk.ContractAddress(state.TimelockProgram, mcmssolanasdk.PDASeed(state.TimelockSeed))
			mcmAddress, err := cfg.TimelockCfg.MCMBasedOnActionSolana(*state)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			mcmAddresses[selector] = mcmAddress
			tx, err := buildNoOPSolana()
			if err != nil {
				return cldf.ChangesetOutput{}, err
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
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
	}, nil
}
