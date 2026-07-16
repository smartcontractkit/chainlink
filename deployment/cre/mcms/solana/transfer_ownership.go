package solana

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmssolanasdk "github.com/smartcontractkit/mcms/sdk/solana"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	proposeutils "github.com/smartcontractkit/cld-changesets/legacy/mcms/proposeutils"
	solstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/solana"
	solchangesets "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/solana/changesets"
	pdasol "github.com/smartcontractkit/cld-changesets/pkg/family/solana"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	mcmscontracts "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/contracts/mcms"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

var transferMCMSOwnershipVersion = semver.MustParse("1.0.0")

// TransferSolanaMCMSOwnershipConfig is the input for transferring the Solana MCMS stack
// to the timelock signer PDA.
type TransferSolanaMCMSOwnershipConfig struct {
	Chains  []uint64
	MCMSCfg cldfproposalutils.TimelockConfig
}

// TransferSolanaMCMSOwnership transfers MCMS program ownership to the timelock signer PDA.
// MCMS addresses are resolved via solstate.GetState, which reads the legacy address book
// first and falls back to the CLD datastore.
type TransferSolanaMCMSOwnership struct{}

var _ cldf.ChangeSetV2[TransferSolanaMCMSOwnershipConfig] = TransferSolanaMCMSOwnership{}

func (TransferSolanaMCMSOwnership) VerifyPreconditions(
	env cldf.Environment, cfg TransferSolanaMCMSOwnershipConfig,
) error {
	for _, chainSelector := range cfg.Chains {
		if err := verifySelector(env, chainSelector); err != nil {
			return err
		}

		state, err := solstate.GetState(env, chainSelector)
		if err != nil {
			return fmt.Errorf("load MCMS timelock state for chain %d: %w", chainSelector, err)
		}
		if err := state.Validate(); err != nil {
			return fmt.Errorf("MCMS timelock state for chain %d incomplete: %w", chainSelector, err)
		}
	}

	return nil
}

func (TransferSolanaMCMSOwnership) Apply(
	env cldf.Environment, cfg TransferSolanaMCMSOwnershipConfig,
) (cldf.ChangesetOutput, error) {
	contractsByChain := make(map[uint64][]solchangesets.OwnableContract, len(cfg.Chains))
	for _, chainSelector := range cfg.Chains {
		chainState, err := solstate.GetState(env, chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("load MCMS timelock state for chain %d: %w", chainSelector, err)
		}

		contractsByChain[chainSelector] = ownableMCMSContracts(chainState)
	}

	return applyTransferMCMSToTimelock(env, contractsByChain, cfg.MCMSCfg)
}

func ownableMCMSContracts(chainState *solstate.MCMSWithTimelockState) []solchangesets.OwnableContract {
	return []solchangesets.OwnableContract{
		{
			ProgramID: chainState.McmProgram,
			Seed:      chainState.ProposerMcmSeed,
			OwnerPDA:  pdasol.GetMCMConfigPDA(chainState.McmProgram, chainState.ProposerMcmSeed),
			Type:      mcmscontracts.ProposerManyChainMultisig,
		},
		{
			ProgramID: chainState.McmProgram,
			Seed:      chainState.CancellerMcmSeed,
			OwnerPDA:  pdasol.GetMCMConfigPDA(chainState.McmProgram, chainState.CancellerMcmSeed),
			Type:      mcmscontracts.CancellerManyChainMultisig,
		},
		{
			ProgramID: chainState.McmProgram,
			Seed:      chainState.BypasserMcmSeed,
			OwnerPDA:  pdasol.GetMCMConfigPDA(chainState.McmProgram, chainState.BypasserMcmSeed),
			Type:      mcmscontracts.BypasserManyChainMultisig,
		},
		{
			ProgramID: chainState.TimelockProgram,
			Seed:      chainState.TimelockSeed,
			OwnerPDA:  pdasol.GetTimelockConfigPDA(chainState.TimelockProgram, chainState.TimelockSeed),
			Type:      mcmscontracts.RBACTimelock,
		},
		{
			ProgramID: chainState.AccessControllerProgram,
			OwnerPDA:  chainState.ProposerAccessControllerAccount,
			Type:      mcmscontracts.ProposerAccessControllerAccount,
		},
		{
			ProgramID: chainState.AccessControllerProgram,
			OwnerPDA:  chainState.ExecutorAccessControllerAccount,
			Type:      mcmscontracts.ExecutorAccessControllerAccount,
		},
		{
			ProgramID: chainState.AccessControllerProgram,
			OwnerPDA:  chainState.CancellerAccessControllerAccount,
			Type:      mcmscontracts.CancellerAccessControllerAccount,
		},
		{
			ProgramID: chainState.AccessControllerProgram,
			OwnerPDA:  chainState.BypasserAccessControllerAccount,
			Type:      mcmscontracts.BypasserAccessControllerAccount,
		},
	}
}

func applyTransferMCMSToTimelock(
	env cldf.Environment,
	contractsByChain map[uint64][]solchangesets.OwnableContract,
	mcmsCfg cldfproposalutils.TimelockConfig,
) (cldf.ChangesetOutput, error) {
	solChains := env.BlockChains.SolanaChains()

	batches := []mcmstypes.BatchOperation{}
	timelocks := map[uint64]string{}
	proposers := map[uint64]string{}
	inspectors := map[uint64]mcmssdk.Inspector{}

	var out cldf.ChangesetOutput
	for chainSelector, contractsToTransfer := range contractsByChain {
		solChain, ok := solChains[chainSelector]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("solana chain not found in environment (selector: %v)", chainSelector)
		}

		chainState, err := solstate.GetState(env, chainSelector)
		if err != nil {
			return out, fmt.Errorf("load MCMS timelock state for chain %d: %w", chainSelector, err)
		}

		inspectors[chainSelector] = mcmssolanasdk.NewInspector(solChain.Client)
		timelocks[chainSelector] = mcmssolanasdk.ContractAddress(chainState.TimelockProgram, mcmssolanasdk.PDASeed(chainState.TimelockSeed))
		proposers[chainSelector] = mcmssolanasdk.ContractAddress(chainState.McmProgram, mcmssolanasdk.PDASeed(chainState.ProposerMcmSeed))

		for _, contract := range contractsToTransfer {
			execOut, execErr := operations.ExecuteOperation(env.OperationsBundle,
				operations.NewOperation(
					"transfer-ownership",
					transferMCMSOwnershipVersion,
					"transfers ownership of contracts to mcms",
					solchangesets.TransferToTimelockSolanaOp,
				),
				solchangesets.Deps{
					Env:   env,
					State: chainState,
					Chain: solChain,
				},
				solchangesets.TransferToTimelockInput{
					Contract: contract,
					MCMSCfg:  mcmsCfg,
				},
			)
			if execErr != nil {
				return out, execErr
			}

			batches = append(batches, execOut.Output.Batches...)
		}
	}

	proposal, err := proposeutils.BuildProposalFromBatchesV2(
		env, timelocks, proposers, inspectors,
		batches, "proposal to transfer ownership of contracts to timelock", mcmsCfg,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}
	env.Logger.Debugw("created timelock proposal", "# batches", len(batches))

	out.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}

	return out, nil
}
