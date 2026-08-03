package operation

import (
	"context"
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
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ks_forwarder "github.com/smartcontractkit/chainlink-solana/contracts/generated/keystone_forwarder"

	commonOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/operations"
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
	ClearForwarderConfigOp = operations.NewOperation(
		"clear-forwarder-config-op",
		Version1_0_0,
		"Closes a DON oracles config of the forwarder for Solana Chain",
		clearForwarderConfig,
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
		MCMS                *cldfproposalutils.TimelockConfig // if set, assumes current upgrade authority is the timelock
	}

	SetUpgradeAuthorityOutput struct {
		Proposals []mcms.TimelockProposal // will be returned in case if timelock config is passed
	}

	// ForwarderConfigTarget identifies the oracles config account an instruction operates on and
	// how that instruction has to be submitted. It is shared by every operation that touches the
	// oracles config of a forwarder.
	ForwarderConfigTarget struct {
		MCMS           *cldfproposalutils.TimelockConfig // if set, assumes current owner is the timelock
		ProgramID      solana.PublicKey
		ForwarderState solana.PublicKey
		ConfigPDA      solana.PublicKey
		Owner          solana.PublicKey
		DonID          uint32
		ConfigVersion  uint32
		Type           cldf.ContractType
	}

	ConfigureForwarderInput struct {
		ForwarderConfigTarget

		Signers [][20]uint8
		F       uint8
	}

	ConfigureForwarderOutput struct {
		Batch mcmsTypes.BatchOperation
	}

	ClearForwarderConfigInput struct {
		ForwarderConfigTarget
	}

	ClearForwarderConfigOutput struct {
		Batch mcmsTypes.BatchOperation
	}
)

func initForwarder(b operations.Bundle, deps Deps, in InitForwarderInput) (InitForwarderOutput, error) {
	var out InitForwarderOutput
	// anchor-go bakes a default program id; deploy uses the keystone_forwarder keypair, which can differ.
	// NewInitializeInstruction uses this package var as the program id, so it must match deploy output.
	ks_forwarder.ProgramID = in.ProgramID

	stateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		return out, fmt.Errorf("failed to create random keys: %w", err)
	}

	instruction, err := ks_forwarder.NewInitializeInstruction(stateKey.PublicKey(), deps.Chain.DeployerKey.PublicKey(), solana.SystemProgramID)
	if err != nil {
		return out, fmt.Errorf("failed to build and validate initialize instruction %w", err)
	}

	instructions := []solana.Instruction{instruction}
	if err = deps.Chain.Confirm(instructions, solanaUtils.AddSigners(stateKey)); err != nil {
		return out, fmt.Errorf("failed to confirm: %w", err)
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

	// anchor-go bakes a default program id, which can differ from the deployed one.
	ks_forwarder.ProgramID = in.ProgramID

	oracleExists, err := oraclesConfigExists(b.GetContext(), deps, in.ConfigPDA)
	if err != nil {
		return out, err
	}

	var instruction solana.Instruction
	if oracleExists {
		instruction, err = ks_forwarder.NewUpdateOraclesConfigInstruction(
			in.DonID,
			in.ConfigVersion,
			in.F,
			in.Signers,
			in.ForwarderState,
			in.ConfigPDA,
			in.Owner,
		)
		if err != nil {
			return out, fmt.Errorf("cant build update oracles config instruction: %w", err)
		}
	} else {
		instruction, err = ks_forwarder.NewInitOraclesConfigInstruction(
			in.DonID,
			in.ConfigVersion,
			in.F,
			in.Signers,
			in.ForwarderState,
			in.ConfigPDA,
			in.Owner,
			solana.SystemProgramID,
		)
		if err != nil {
			return out, fmt.Errorf("cant build init oracles config instruction: %w", err)
		}
	}

	out.Batch, err = submitOrPropose(b, deps, in.ForwarderConfigTarget, instruction)

	return out, err
}

// clearForwarderConfig closes the oracles config account of a single DON. The rent is refunded to
// the owner, which makes the config version reusable by a later configure.
func clearForwarderConfig(b operations.Bundle, deps Deps, in ClearForwarderConfigInput) (ClearForwarderConfigOutput, error) {
	var out ClearForwarderConfigOutput

	// anchor-go bakes a default program id, which can differ from the deployed one.
	ks_forwarder.ProgramID = in.ProgramID

	oracleExists, err := oraclesConfigExists(b.GetContext(), deps, in.ConfigPDA)
	if err != nil {
		return out, err
	}
	if !oracleExists {
		return out, fmt.Errorf("no oracles config %s found for don %d config version %d", in.ConfigPDA, in.DonID, in.ConfigVersion)
	}

	instruction, err := ks_forwarder.NewCloseOraclesConfigInstruction(
		in.DonID,
		in.ConfigVersion,
		in.ForwarderState,
		in.ConfigPDA,
		in.Owner,
	)
	if err != nil {
		return out, fmt.Errorf("cant build close oracles config instruction: %w", err)
	}

	out.Batch, err = submitOrPropose(b, deps, in.ForwarderConfigTarget, instruction)

	return out, err
}

// oraclesConfigExists reports whether the oracles config account was already initialized. It reads
// at the commitment level the chain confirms transactions with, so a config written earlier in the
// same run is visible.
func oraclesConfigExists(ctx context.Context, deps Deps, configPDA solana.PublicKey) (bool, error) {
	_, err := deps.Chain.Client.GetAccountInfoWithOpts(ctx, configPDA, &rpc.GetAccountInfoOpts{
		Commitment: cldfsol.SolDefaultCommitment,
	})
	if err == nil {
		return true, nil
	}
	if errors.Is(err, rpc.ErrNotFound) {
		return false, nil
	}

	return false, fmt.Errorf("can't confirm oracles config %s existence: %w", configPDA, err)
}

// submitOrPropose sends the instruction with the deployer key, or, when the forwarder is owned by
// the timelock, returns it as an MCMS batch operation for the caller to propose.
func submitOrPropose(b operations.Bundle, deps Deps, target ForwarderConfigTarget, instruction solana.Instruction) (mcmsTypes.BatchOperation, error) {
	if target.MCMS == nil {
		if err := deps.Chain.Confirm([]solana.Instruction{instruction}); err != nil {
			return mcmsTypes.BatchOperation{}, fmt.Errorf("failed to confirm instruction: %w", err)
		}

		return mcmsTypes.BatchOperation{}, nil
	}

	tx, err := helpers.BuildMCMSTxn(instruction, target.ProgramID.String(), target.Type)
	if err != nil {
		return mcmsTypes.BatchOperation{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	b.Logger.Infof("build mcmstxn contract type: %q program_id: %q", target.Type.String(), target.ProgramID.String())

	return mcmsTypes.BatchOperation{
		ChainSelector: mcmsTypes.ChainSelector(deps.Chain.ChainSelector()),
		Transactions:  []mcmsTypes.Transaction{*tx},
	}, nil
}
