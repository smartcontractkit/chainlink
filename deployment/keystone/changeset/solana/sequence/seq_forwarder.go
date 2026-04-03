package sequence

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	commonOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/operations"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/solana/sequence/operation"
)

var (
	DeployForwarderSeq = operations.NewSequence(
		"deploy-forwarder-seq",
		operation.Version1_0_0,
		"Deploys forwarder contract and initializes it",
		deployForwarder,
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
