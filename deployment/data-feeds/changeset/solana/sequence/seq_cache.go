package sequence

import (
	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
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
	ChainSel           uint64
	ProgramName        string
	ForwarderProgramID solana.PublicKey   // ForwarderProgram that is allowed to write to this cache
	FeedAdmins         []solana.PublicKey // Feed admins to be added to the cache
	ContractType       datastore.ContractType
	Version            *semver.Version
	Qualifier          string
}

// DeployCacheSeqOutput defines the output of the deployment sequence.
type DeployCacheSeqOutput struct {
	ProgramID solana.PublicKey
	State     solana.PublicKey
}

func deployCache(b operations.Bundle, deps operation.Deps, in DeployCacheSeqInput) (DeployCacheSeqOutput, error) {
	var out DeployCacheSeqOutput

	// 1. Deploy the DataFeeds Cache program
	var cacheID solana.PublicKey
	programID, err := deps.Datastore.Addresses().Get(datastore.NewAddressRefKey(
		in.ChainSel,
		in.ContractType,
		in.Version,
		in.Qualifier,
	))
	if err != nil {
		deployOut, err2 := operations.ExecuteOperation(b, operation.DeployCacheOp, commonOps.Deps{Chain: deps.Chain}, commonOps.DeployInput{
			ProgramName: in.ProgramName,
			ChainSel:    in.ChainSel,
		})
		if err2 != nil {
			return DeployCacheSeqOutput{}, err2
		}
		cacheID = deployOut.Output.ProgramID
	} else {
		cacheID = solana.MustPublicKeyFromBase58(programID.Address)
	}

	out.ProgramID = cacheID

	// 2. Initialize the DataFeeds Cache state
	initOut, err := operations.ExecuteOperation(b, operation.InitCacheOp, deps, operation.InitCacheInput{
		ProgramID:          out.ProgramID,
		ChainSel:           in.ChainSel,
		FeedAdmins:         in.FeedAdmins,
		ForwarderProgramID: in.ForwarderProgramID,
	})
	if err != nil {
		return DeployCacheSeqOutput{}, err
	}
	out.State = initOut.Output.StatePubKey

	return out, nil
}
