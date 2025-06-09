package sequence

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
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
		ProgramID   string
		StatePubKey string
		// will see what return here
	}
	ConfigureForwarderSeqInput struct {
	}
	ConfigureForwarderSeqOutput struct {
	}
)

func deployForwarder(b operations.Bundle, deps operation.Deps, in DeployForwarderSeqInput) (DeployForwarderSeqOutput, error) {
	var out DeployForwarderSeqOutput

	// 1. Deploy
	// IsUpgrade false because it's initial deploy sequence
	deployOut, err := operations.ExecuteOperation(b, operation.DeployForwarderOp, deps, operation.DeployInput{
		ProgramName:  in.ProgramName,
		Overallocate: in.Overallocate,
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
	out.StatePubKey = initOut.Output.StatePubKey

	return out, nil
}

func configureForwarder(b operations.Bundle, deps operation.Deps, in ConfigureForwarderSeqInput) (ConfigureForwarderSeqOutput, error)
