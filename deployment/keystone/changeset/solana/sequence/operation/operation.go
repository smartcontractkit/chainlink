package operation

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solanaUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	cldfsol "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ks_forwarder "github.com/smartcontractkit/chainlink-solana/contracts/generated/keystone_forwarder"

	commonOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/operations"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

var Version1_0_0 = semver.MustParse("1.0.0")

var (
	DeployForwarderOp = operations.NewOperation(
		"deploy-forwarder-op",
		Version1_0_0,
		"Deploys deploys forwarder for Solana Chain",
		commonOps.Deploy,
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
	ConfigureForwarderOp = operations.NewOperation(
		"configure-forwarder-op",
		Version1_0_0,
		"Configure forwarder for Solana Chain",
		configureForwarder,
	)
)

type (
	Deps struct {
		Env       cldf.Environment
		Chain     cldfsol.Chain
		Datastore datastore.DataStore
	}

	InitForwarderInput struct {
		ProgramID solana.PublicKey
		ChainSel  uint64
	}

	InitForwarderOutput struct {
		StatePubKey solana.PublicKey
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

	ConfigureForwarderInput struct {
		MCMS           *proposalutils.TimelockConfig // if set, assumes current owner is the timelock
		ConfigPDA      string
		ProgramID      solana.PublicKey
		ForwarderState solana.PublicKey
		Owner          string
		Signers        [][20]uint8
		DonID          uint32
		ConfigVersion  uint32
		F              uint8
		Type           cldf.ContractType
	}

	ConfigureForwarderOutput struct {
		Batch mcmsTypes.BatchOperation
	}
)

func initForwarder(b operations.Bundle, deps Deps, in InitForwarderInput) (InitForwarderOutput, error) {
	var out InitForwarderOutput
	if ks_forwarder.ProgramID.IsZero() {
		ks_forwarder.SetProgramID(in.ProgramID)
	}

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

	out.StatePubKey = stateKey.PublicKey()

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
		timelockSignerPDA, err := helpers.FetchTimelockSigner(deps.Datastore.Addresses().Filter(datastore.AddressRefByChainSelector(in.ChainSel)))
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
	var out ConfigureForwarderOutput

	var instructions *ks_forwarder.Instruction
	if ks_forwarder.ProgramID.IsZero() {
		ks_forwarder.SetProgramID(in.ProgramID)
	}

	configPDA := solana.MustPublicKeyFromBase58(in.ConfigPDA)

	var oracleExists bool

	_, err := deps.Chain.Client.GetAccountInfo(b.GetContext(), configPDA)
	if err != nil {
		if !errors.Is(err, rpc.ErrNotFound) {
			return out, fmt.Errorf("can't confirm oracle existence: %w", err)
		}
		oracleExists = false
	} else {
		oracleExists = true
	}

	owner := solana.MustPublicKeyFromBase58(in.Owner)

	if !oracleExists {
		instructions, err = ks_forwarder.NewInitOraclesConfigInstruction(
			in.DonID,
			in.ConfigVersion,
			in.F,
			in.Signers,
			in.ForwarderState,
			configPDA,
			owner,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return out, fmt.Errorf("cant build init oracle instruction: %w", err)
		}
	} else {
		instructions, err = ks_forwarder.NewUpdateOraclesConfigInstruction(
			in.DonID,
			in.ConfigVersion,
			in.F,
			in.Signers,
			in.ForwarderState,
			configPDA,
			owner,
			solana.SystemProgramID,
		).ValidateAndBuild()
		if err != nil {
			return out, fmt.Errorf("cant build init oracle instruction: %w", err)
		}
	}

	if in.MCMS == nil {
		err := deps.Chain.Confirm([]solana.Instruction{instructions})
		return out, err
	}

	tx, err := helpers.BuildMCMSTxn(
		instructions,
		in.ProgramID.String(),
		in.Type)
	if err != nil {
		return out, fmt.Errorf("failed to create transaction: %w", err)
	}

	b.Logger.Infof("build mcmstxn contract type: %q program_id: %q", in.Type.String(), in.ProgramID.String())
	out.Batch = mcmsTypes.BatchOperation{
		ChainSelector: mcmsTypes.ChainSelector(deps.Chain.ChainSelector()),
		Transactions:  []mcmsTypes.Transaction{*tx},
	}

	return out, nil
}
