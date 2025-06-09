package operation

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	solanaUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	cldfsol "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ks_forwarder "github.com/smartcontractkit/chainlink-solana/contracts/generated/keystone_forwarder"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
)

var Version1_0_0 = semver.MustParse("1.0.0")

var (
	DeployForwarderOp = operations.NewOperation(
		"deploy-forwarder-op",
		Version1_0_0,
		"Deploys deploys forwarder for Solana Chain",
		deploy,
	)
	InitForwarderOp = operations.NewOperation(
		"init-forwarder-op",
		Version1_0_0,
		"Initialize forwarder for Solana Chain",
		initForwarder,
	)
	SetUpgradeAuthorityOp = operations.NewOperation(
		"set-upgrade-authority-op",
		Version1_0_0,
		"Sets upgrade forwarder's upgrade authority for Solana Chain",
		setUpgradeAuthority,
	)
)

type (
	Deps struct {
		Env       cldf.Environment
		Chain     cldfsol.Chain
		Datastore datastore.DataStore
	}

	DeployInput struct {
		ChainSel     uint64
		ProgramName  string
		Overallocate bool
	}

	DeployOutput struct {
		ProgramID string
	}

	InitForwarderInput struct {
		ProgramID string
		ChainSel  uint64
	}

	InitForwarderOutput struct {
		StatePubKey string
	}

	SetUpgradeAuthorityInput struct {
		ChainSel            uint64
		ProgramID           string
		NewUpgradeAuthority string
		MCMS                *proposalutils.TimelockConfig // if set, assumes current upgrade authority is the timelock
	}

	SetUpgradeAuthorityOutput struct {
		Proposals []mcms.TimelockProposal // will be returned in case if timelock config is passed
	}

	ConfigureForwarderInput struct{}

	ConfigureForwarderOutput struct {
		Batch mcmsTypes.BatchOperation
	}
)

func deploy(b operations.Bundle, deps Deps, in DeployInput) (DeployOutput, error) {
	var out DeployOutput

	programID, err := deps.Chain.DeployProgram(deps.Env.Logger, cldfsol.ProgramInfo{
		Name:  in.ProgramName,
		Bytes: deployment.SolanaProgramBytes[in.ProgramName],
	}, false, in.Overallocate)
	if err != nil {
		return out, err
	}

	out.ProgramID = programID

	return out, nil
}

func initForwarder(b operations.Bundle, deps Deps, in InitForwarderInput) (InitForwarderOutput, error) {
	var out InitForwarderOutput
	address, err := solana.PublicKeyFromBase58(in.ProgramID)
	if err != nil {
		return out, err
	}
	ks_forwarder.SetProgramID(address)

	stateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		return out, fmt.Errorf("failed to create random keys: %w", err)
	}

	instruction, err := ks_forwarder.NewInitializeInstruction(stateKey.PublicKey(), deps.Chain.DeployerKey.PublicKey(), solana.SystemProgramID).ValidateAndBuild()
	if err != nil {
		return out, fmt.Errorf("failed to build and validate initialize instruction %w", err)
	}

	instructions := []solana.Instruction{instruction}
	if err = deps.Chain.Confirm(instructions, solanaUtils.AddSigners(stateKey)); err != nil {
		return out, errors.New("failed to confirm ")
	}

	out.StatePubKey = stateKey.PublicKey().String()

	return out, nil
}

func setUpgradeAuthority(b operations.Bundle, deps Deps, in SetUpgradeAuthorityInput) (SetUpgradeAuthorityOutput, error) {
	var out SetUpgradeAuthorityOutput

	programID, err := solana.PublicKeyFromBase58(in.ProgramID)
	if err != nil {
		return out, fmt.Errorf("failed parse programID: %w", err)
	}

	newAuthority, err := solana.PublicKeyFromBase58(in.NewUpgradeAuthority)
	if err != nil {
		return out, fmt.Errorf("failed parse upgrade authority: %w", err)
	}

	currentAuthority := deps.Chain.DeployerKey.PublicKey()
	if in.MCMS != nil {
		timelockSignerPDA, err := helpers.FetchTimelockSigner(deps.Env, deps.Chain.Selector)
		if err != nil {
			return out, fmt.Errorf("failed to get timelock signer: %w", err)
		}
		currentAuthority = timelockSignerPDA
	}

	mcmsTxns := make([]mcmsTypes.Transaction, 0)

	ixn := helpers.SetUpgradeAuthority(&deps.Env, programID, currentAuthority, newAuthority, false)

	if in.MCMS == nil {
		if err := deps.Chain.Confirm([]solana.Instruction{ixn}); err != nil {
			return out, fmt.Errorf("failed to confirm instructions: %w", err)
		}

		return out, nil
	}

	// build MCMS proposal
	tx, err := helpers.BuildMCMSTxn(
		ixn,
		solana.BPFLoaderUpgradeableProgramID.String(),
		cldf.ContractType(solana.BPFLoaderUpgradeableProgramID.String()))
	if err != nil {
		return out, fmt.Errorf("failed to create transaction: %w", err)
	}
	mcmsTxns = append(mcmsTxns, *tx)

	proposal, err := helpers.BuildProposalsForTxns(
		deps.Env, in.ChainSel, "proposal to SetUpgradeAuthority in Solana", in.MCMS.MinDelay, mcmsTxns)
	if err != nil {
		return out, fmt.Errorf("failed to build proposal: %w", err)
	}
	out.Proposals = []mcms.TimelockProposal{*proposal}

	return out, nil
}

func configureForwarder(b operations.Bundle, deps Deps, in ConfigureForwarderInput) (ConfigureForwarderOutput, error) {
}
