package sequence

import (
	"github.com/gagliardetto/solana-go"

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
	deployOut, err := operations.ExecuteOperation(b, operation.DeployForwarderOp, commonOps.Deps{Chain: deps.Chain}, commonOps.DeployInput{
		ProgramName:  in.ProgramName,
		Overallocate: in.Overallocate,
		Size:         KeystoneForwarderProgramSize,
		ChainSel:     in.ChainSel,
	})

	if err != nil {
		return DeployForwarderSeqOutput{}, nil
	}
	out.ProgramID = deployOut.Output.ProgramID

	// 2. Initialize
	initOut, err := operations.ExecuteOperation(b, operation.InitForwarderOp, deps, operation.InitForwarderInput{
		ProgramID: out.ProgramID,
		ChainSel:  in.ChainSel,
	})

	if err != nil {
		return DeployForwarderSeqOutput{}, nil
	}
	out.State = initOut.Output.StatePubKey

	return out, nil
}
