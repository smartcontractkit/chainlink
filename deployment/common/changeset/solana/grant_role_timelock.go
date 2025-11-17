package solana

import (
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	timelockbindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/timelock"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/evm/mcms/seqs"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// GrantRoleTimelockSolana grants the given accounts access to the given role on the timelock
type GrantRoleTimelockSolana struct{}

type GrantRoleTimelockSolanaConfig struct {
	Accounts map[uint64][]solana.PublicKey // chain selector to accounts mapping
	Role     timelockbindings.Role
	MCMS     *proposalutils.TimelockConfig
}

func (t GrantRoleTimelockSolana) VerifyPreconditions(
	env cldf.Environment, config GrantRoleTimelockSolanaConfig,
) error {
	if !validTimelockActions(config.MCMS) {
		return fmt.Errorf("invalid mcms action: %v", config.MCMS.MCMSAction)
	}

	solanaChains := env.BlockChains.SolanaChains()
	if len(solanaChains) == 0 {
		return errors.New("no solana chains provided")
	}

	// check the timelock program is deployed
	for chainSelector := range config.Accounts {
		solChain, ok := solanaChains[chainSelector]
		if !ok {
			return fmt.Errorf("solana chain not found for selector %d", chainSelector)
		}

		chainAddresses, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return fmt.Errorf("failed to get existing addresses: %w", err)
		}

		if env.DataStore != nil && env.DataStore.Addresses() != nil {
			datastoreAddresses := env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chainSelector))
			for _, address := range datastoreAddresses {
				if address.Version == nil {
					return fmt.Errorf("address without Version found in data store: %#v", address)
				}

				chainAddresses[address.Address] = cldf.TypeAndVersion{
					Type:    cldf.ContractType(address.Type),
					Version: *address.Version,
					Labels:  cldf.NewLabelSet(address.Labels.List()...),
				}
			}
		}

		mcmState, err := state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, chainAddresses)
		if err != nil {
			return fmt.Errorf("failed to load MCMS state: %w", err)
		}
		if mcmState.TimelockProgram.IsZero() {
			return fmt.Errorf("timelock program not deployed for chain %d", chainSelector)
		}
		if (mcmState.TimelockSeed == state.PDASeed{}) {
			return fmt.Errorf("timelock seed not found for chain %d", chainSelector)
		}
	}

	return nil
}

func (t GrantRoleTimelockSolana) Apply(
	env cldf.Environment, cfg GrantRoleTimelockSolanaConfig,
) (cldf.ChangesetOutput, error) {
	batchOps := []mcmstypes.BatchOperation{}
	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]mcmssdk.Inspector{}

	for chainSelector, accountsList := range cfg.Accounts {
		solChain := env.BlockChains.SolanaChains()[chainSelector]

		addresses, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get existing addresses: %w", err)
		}
		mcmsChainState, _ := state.MaybeLoadMCMSWithTimelockChainStateSolana(solChain, addresses)

		deps := seqs.SeqSolanaGrantRoleTimelockDeps{Chain: solChain}
		input := seqs.SeqSolanaGrantRoleTimelockInput{
			ChainState:         mcmsChainState,
			Role:               cfg.Role,
			Accounts:           accountsList,
			IsDeployerKeyAdmin: cfg.MCMS == nil,
		}
		report, err := operations.ExecuteSequence(env.OperationsBundle, seqs.SeqSolanaGrantRoleTimelock, deps, input)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute operation %q: %w", seqs.SeqSolanaGrantRoleTimelock.ID(), err)
		}

		if cfg.MCMS != nil {
			batchOps = append(batchOps, mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(chainSelector),
				Transactions:  report.Output.McmsTransactions,
			})
			proposers[chainSelector], _ = proposalMCM(mcmsChainState, cfg.MCMS.MCMSAction)
			timelocks[chainSelector] = state.EncodeAddressWithSeed(mcmsChainState.TimelockProgram, mcmsChainState.TimelockSeed)
			inspectors[chainSelector] = mcmssolanasdk.NewInspector(solChain.Client)
		}
	}

	if cfg.MCMS == nil {
		// direct on chain execution; return early
		return cldf.ChangesetOutput{}, nil
	}

	// create MCMS proposal
	proposal, err := proposalutils.BuildProposalFromBatchesV2(env, timelocks, proposers, inspectors,
		batchOps, "proposal to grant role in timelock", *cfg.MCMS)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcms.TimelockProposal{*proposal}}, nil
}

func validTimelockActions(timelockConfig *proposalutils.TimelockConfig) bool {
	if timelockConfig == nil {
		return true
	}

	switch timelockConfig.MCMSAction {
	case "", mcmstypes.TimelockActionSchedule, mcmstypes.TimelockActionBypass:
		return true
	default:
		return false
	}
}

func proposalMCM(mcmsState *state.MCMSWithTimelockStateSolana, action mcmstypes.TimelockAction) (string, error) {
	switch action {
	case "", mcmstypes.TimelockActionSchedule:
		return state.EncodeAddressWithSeed(mcmsState.McmProgram, mcmsState.ProposerMcmSeed), nil
	case mcmstypes.TimelockActionBypass:
		return state.EncodeAddressWithSeed(mcmsState.McmProgram, mcmsState.BypasserMcmSeed), nil
	default:
		return "", fmt.Errorf("invalid mcms action: %v", action)
	}
}
