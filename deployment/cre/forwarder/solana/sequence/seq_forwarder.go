package sequence

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	commonOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder/solana/sequence/operation"
)

var (
	DeployForwarderSeq = operations.NewSequence(
		"deploy-forwarder-seq",
		operation.Version1_0_0,
		"Deploys forwarder contract and initializes it",
		deployForwarder,
	)

	UpgradeForwarderSeq = operations.NewSequence(
		"upgrade-forwarder-seq",
		operation.Version1_0_0,
		"Writes the forwarder contract to a buffer and upgrades the deployed program in place",
		upgradeForwarder,
	)
)

type (
	DeployForwarderSeqInput struct {
		ChainSel     uint64
		ProgramName  string
		Qualifier    string
		Version      *semver.Version
		ContractType datastore.ContractType
		Overallocate bool
	}

	DeployForwarderSeqOutput struct {
		ProgramID solana.PublicKey
		State     solana.PublicKey
	}

	UpgradeForwarderSeqInput struct {
		ChainSel    uint64
		ProgramName string
		// ProgramID is the deployed program that gets upgraded in place.
		ProgramID solana.PublicKey
		MCMS      *cldfproposalutils.TimelockConfig // if set, assumes current upgrade authority is the timelock
	}

	UpgradeForwarderSeqOutput struct {
		BufferID  solana.PublicKey
		Proposals []mcms.TimelockProposal // will be returned in case if timelock config is passed
	}
)

const KeystoneForwarderProgramSize = 5 * 1024 * 1024

func deployForwarder(b operations.Bundle, deps operation.Deps, in DeployForwarderSeqInput) (DeployForwarderSeqOutput, error) {
	var out DeployForwarderSeqOutput

	// 1. Deploy
	var forwarderID solana.PublicKey
	programID, err := deps.Datastore.Addresses().Get(datastore.NewAddressRefKey(
		in.ChainSel,
		in.ContractType,
		in.Version,
		in.Qualifier,
	))

	if err != nil {
		deployOut, err2 := operations.ExecuteOperation(b, operation.DeployForwarderOp, commonOps.Deps{Chain: deps.Chain}, commonOps.DeployInput{
			ProgramName:  in.ProgramName,
			Overallocate: in.Overallocate,
			Size:         KeystoneForwarderProgramSize,
			ChainSel:     in.ChainSel,
		})
		if err2 != nil {
			return DeployForwarderSeqOutput{}, fmt.Errorf("deploy forwarder op failed: %w", err2)
		}
		forwarderID = deployOut.Output.ProgramID
	} else {
		deps.Env.Logger.Info("Forwarder program ID is already present in datastore for given version and qualifier. Proceed sequence without deploying")
		forwarderID = solana.MustPublicKeyFromBase58(programID.Address)
	}

	out.ProgramID = forwarderID

	// 2. Initialize
	initOut, err := operations.ExecuteOperation(b, operation.InitForwarderOp, deps, operation.InitForwarderInput{
		ProgramID: out.ProgramID,
		ChainSel:  in.ChainSel,
	})

	if err != nil {
		return DeployForwarderSeqOutput{}, fmt.Errorf("initialize forwarder op failed: %w", err)
	}
	out.State = initOut.Output.StatePubKey

	return out, nil
}

func upgradeForwarder(b operations.Bundle, deps operation.Deps, in UpgradeForwarderSeqInput) (UpgradeForwarderSeqOutput, error) {
	var out UpgradeForwarderSeqOutput

	// 1. Write the new binary to a buffer account. The buffer is owned by the deployer key until the
	// upgrade operation hands it over to the program's upgrade authority.
	deployOut, err := operations.ExecuteOperation(b, operation.DeployForwarderOp, commonOps.Deps{Chain: deps.Chain}, commonOps.DeployInput{
		ProgramName: in.ProgramName,
		ChainSel:    in.ChainSel,
		IsUpgrade:   true,
	})
	if err != nil {
		return out, fmt.Errorf("write forwarder buffer op failed: %w", err)
	}
	out.BufferID = deployOut.Output.ProgramID

	// 2. Upgrade in place
	upgradeOut, err := operations.ExecuteOperation(b, operation.UpgradeForwarderOp, deps, operation.UpgradeForwarderInput{
		ChainSel:  in.ChainSel,
		ProgramID: in.ProgramID,
		BufferID:  out.BufferID,
		MCMS:      in.MCMS,
	})
	if err != nil {
		return out, fmt.Errorf("upgrade forwarder op failed: %w", err)
	}
	out.Proposals = upgradeOut.Output.Proposals

	return out, nil
}
