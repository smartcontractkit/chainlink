package sequence

import (
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	commonOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/operations"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana/sequence/operation"
)

var (
	DeployCacheSeq = operations.NewSequence(
		"deploy-cache-seq",
		operation.Version1_0_0,
		"Deploys DataFeeds Cache program and initializes its state",
		deployCache,
	)
)

// DeployCacheSeqInput defines the input for deploying the DataFeeds Cache program.
type DeployCacheSeqInput struct {
	ChainSel    uint64
	ProgramName string
	FeedAdmins  []solana.PublicKey // Feed admins to be added to the cache
}

// DeployCacheSeqOutput defines the output of the deployment sequence.
type DeployCacheSeqOutput struct {
	ProgramID solana.PublicKey
	State     solana.PublicKey
}

func deployCache(b operations.Bundle, deps operation.Deps, in DeployCacheSeqInput) (DeployCacheSeqOutput, error) {
	var out DeployCacheSeqOutput

	// 1. Deploy the DataFeeds Cache program
	deployOut, err := operations.ExecuteOperation(b, operation.DeployCacheOp, commonOps.Deps{Chain: deps.Chain}, commonOps.DeployInput{
		ProgramName: in.ProgramName,
		ChainSel:    in.ChainSel,
	})
	if err != nil {
		return DeployCacheSeqOutput{}, err
	}
	out.ProgramID = deployOut.Output.ProgramID
	fmt.Printf("Cache program deployed with ID: %s\n", out.ProgramID)

	// 2. Initialize the DataFeeds Cache state
	initOut, err := operations.ExecuteOperation(b, operation.InitCacheOp, deps, operation.InitCacheInput{
		ProgramID:  out.ProgramID,
		ChainSel:   in.ChainSel,
		FeedAdmins: in.FeedAdmins,
	})
	if err != nil {
		return DeployCacheSeqOutput{}, err
	}
	out.State = initOut.Output.StatePubKey

	return out, nil
}
