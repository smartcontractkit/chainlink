package solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	solanaUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	ks_forwarder "github.com/smartcontractkit/chainlink-solana/contracts/generated/keystone_forwarder"
	"github.com/smartcontractkit/chainlink/deployment"
	cdeployment "github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/shared"
	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
)

type DeployRequest = struct {
	ChainSel    uint64
	BuildConfig *helpers.BuildSolanaConfig
}

var _ cldf.ChangeSet[*DeployRequest] = DeployForwarder

func DeployForwarder(env cldf.Environment, req *DeployRequest) (cldf.ChangesetOutput, error) {
	if req.BuildConfig != nil {
		err := helpers.BuildSolana(env, *req.BuildConfig, keystoneBuildParams)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	chain := env.SolChains[req.ChainSel]
	ab := cldf.NewMemoryAddressBook()

	address, err := helpers.DeployAndMaybeSaveToAddressBook(env, chain, ab, shared.Forwarder, cdeployment.Version1_0_0, false, "")
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// initialize
	ks_forwarder.SetProgramID(address)

	stateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create random keys: %w", err)
	}

	instruction, err := ks_forwarder.NewInitializeInstruction(stateKey.PublicKey(), chain.DeployerKey.PublicKey(), solana.SystemProgramID).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build and validate initialize instruction %w", err)
	}

	instructions := []solana.Instruction{instruction}
	if err = chain.Confirm(instructions, solanaUtils.AddSigners(stateKey)); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm ")
	}

	tv := cldf.NewTypeAndVersion(shared.ForwarderState, deployment.Version1_0_0)
	err = ab.Save(req.ChainSel, stateKey.PublicKey().String(), tv)

	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save forwarder state address: %w", err)
	}
	return cldf.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

type SetForwarderUpgradeAuthorityRequest = struct {
	ChainSel            uint64
	NewUpgradeAuthority solana.PublicKey
	MCMS                *proposalutils.TimelockConfig // if set, assumes current upgrade authority is the timelock
}

var _ cldf.ChangeSet[*SetForwarderUpgradeAuthorityRequest] = SetForwarderUpgradeAuthority

func SetForwarderUpgradeAuthority(env cldf.Environment, req *SetForwarderUpgradeAuthorityRequest) (cldf.ChangesetOutput, error) {
	chain, ok := env.SolChains[req.ChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("can't get chain for chain selector %d", req.ChainSel)
	}

	state, err := loadOnchainState(env, req.ChainSel)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if state.forwarderProgramID.IsZero() {
		return cldf.ChangesetOutput{}, fmt.Errorf("forwarder not found for chain selector %d", req.ChainSel)
	}

	currentAuthority := chain.DeployerKey.PublicKey()
	if req.MCMS != nil {
		timelockSignerPDA, err := helpers.FetchTimelockSigner(env, chain.Selector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get timelock signer: %w", err)
		}
		currentAuthority = timelockSignerPDA
	}

	env.Logger.Infow("Setting upgrade authority for", state.forwarderProgramID.String(), "newUpgradeAuthority", req.NewUpgradeAuthority.String())
	mcmsTxns := make([]mcmsTypes.Transaction, 0)
	ixn := helpers.SetUpgradeAuthority(&env, state.forwarderProgramID, currentAuthority, req.NewUpgradeAuthority, false)
	if req.MCMS == nil {
		if err := chain.Confirm([]solana.Instruction{ixn}); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm instructions: %w", err)
		}
	} else {
		tx, err := helpers.BuildMCMSTxn(
			ixn,
			solana.BPFLoaderUpgradeableProgramID.String(),
			cldf.ContractType(solana.BPFLoaderUpgradeableProgramID.String()))
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
		}
		mcmsTxns = append(mcmsTxns, *tx)
	}
	if len(mcmsTxns) > 0 {
		proposal, err := helpers.BuildProposalsForTxns(
			env, req.ChainSel, "proposal to SetUpgradeAuthority in Solana", req.MCMS.MinDelay, mcmsTxns)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{}, nil
}
