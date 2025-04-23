package changeset

import (
	"errors"
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
}

func buildNoOPEVM(e deployment.Environment, selector uint64) (mcmstypes.Transaction, error) {
	chain, ok := e.Chains[selector]
	if !ok {
		return mcmstypes.Transaction{}, nil
	}
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

// buildNoOPSolana builds a dummy tx that safely calls getMinDelay on the RBACTimelock
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
func getMcmAddressFromActionEVM(action mcmstypes.TimelockAction, state *state.MCMSWithTimelockState) (string, error) {
	switch action {
	case mcmstypes.TimelockActionBypass:
		return state.BypasserMcm.Address().String(), nil
	case mcmstypes.TimelockActionSchedule:
		return state.ProposerMcm.Address().String(), nil
	case mcmstypes.TimelockActionCancel:
		return state.CancellerMcm.Address().String(), nil
	}
	return "", errors.New("invalid action")
}
func getMcmAddressFromActionSol(action mcmstypes.TimelockAction, state *state.MCMSWithTimelockStateSolana) (string, error) {
	switch action {
	case mcmstypes.TimelockActionBypass:
		contractID := mcmssolanasdk.ContractAddress(state.McmProgram, mcmssolanasdk.PDASeed(state.BypasserMcmSeed))
		return contractID, nil
	case mcmstypes.TimelockActionSchedule:
		contractID := mcmssolanasdk.ContractAddress(state.McmProgram, mcmssolanasdk.PDASeed(state.ProposerMcmSeed))
		return contractID, nil
	case mcmstypes.TimelockActionCancel:
		contractID := mcmssolanasdk.ContractAddress(state.McmProgram, mcmssolanasdk.PDASeed(state.CancellerMcmSeed))
		return contractID, nil
	}
	return "", errors.New("invalid action")
}

// MCMSSignFireDrillChangeset creates a changeset for a MCMS signing Fire Drill.
// It is used to make sure team member can effectively sign proposal and that the execution pipelines are healthy.
// The changeset will create a NO-OP transaction for each chain selector in the environment and create a proposal for it.
func MCMSSignFireDrillChangeset(e deployment.Environment, cfg FireDrillConfig) (deployment.ChangesetOutput, error) {
	allSelectors := e.AllChainSelectors()
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
			addresses, err := e.ExistingAddresses.AddressesForChain(selector)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			state, err := state.MaybeLoadMCMSWithTimelockChainState(e.Chains[selector], addresses)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			timelocks[selector] = state.Timelock.Address().String()
			mcmAddress, err := getMcmAddressFromActionEVM(cfg.TimelockCfg.MCMSAction, state)
			mcmAddresses[selector] = mcmAddress
			tx, err := buildNoOPEVM(e, selector)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			operations = append(operations, mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(selector),
				Transactions:  []mcmstypes.Transaction{tx},
			})
		case chainsel.FamilySolana:
			addresses, err := e.ExistingAddresses.AddressesForChain(selector)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			state, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(e.SolChains[selector], addresses)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			timelocks[selector] = mcmssolanasdk.ContractAddress(state.TimelockProgram, mcmssolanasdk.PDASeed(state.TimelockSeed))
			mcmAddress, err := getMcmAddressFromActionSol(cfg.TimelockCfg.MCMSAction, state)
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
