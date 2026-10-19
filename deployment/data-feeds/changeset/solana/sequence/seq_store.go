package sequence

import (
	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	commonOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/operations"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana/sequence/operation"
)

var DeployStoreSeq = operations.NewSequence(
	"deploy-data-feeds-store-seq",
	operation.Version1_0_0,
	"Deploys the DataFeeds Store program and initializes its state",
	deployStore,
)

// DeployStoreSeqInput defines the input for deploying the DataFeeds Store program.
type DeployStoreSeqInput struct {
	ChainSel                  uint64
	ProgramName               string
	ForwarderProgramID        solana.PublicKey   // forwarder program allowed to report to the store
	FeedAdmins                []solana.PublicKey // strictly increasing
	DefaultMaxStoredReports   uint64
	MaxObservationForwardSecs uint32
	ContractType              datastore.ContractType
	Version                   *semver.Version
	Qualifier                 string
}

// DeployStoreSeqOutput defines the output of the deployment sequence.
type DeployStoreSeqOutput struct {
	ProgramID solana.PublicKey
	// Deployed is false when the program was already in the datastore
	Deployed bool
	State    solana.PublicKey
}

func deployStore(b operations.Bundle, deps operation.Deps, in DeployStoreSeqInput) (DeployStoreSeqOutput, error) {
	var out DeployStoreSeqOutput

	// 1. Deploy the program, unless it already is
	ref, err := deps.Datastore.Addresses().Get(datastore.NewAddressRefKey(in.ChainSel, in.ContractType, in.Version, in.Qualifier))
	if err != nil {
		deployOut, err := operations.ExecuteOperation(b, operation.DeployStoreOp, commonOps.Deps{Chain: deps.Chain}, commonOps.DeployInput{
			ProgramName: in.ProgramName,
			ChainSel:    in.ChainSel,
		})
		if err != nil {
			return out, err
		}
		out.ProgramID = deployOut.Output.ProgramID
		out.Deployed = true
	} else {
		out.ProgramID, err = solana.PublicKeyFromBase58(ref.Address)
		if err != nil {
			return out, err
		}
	}

	// 2. Initialize a store state
	initOut, err := operations.ExecuteOperation(b, operation.InitStoreOp, deps, operation.InitStoreInput{
		ProgramID:                 out.ProgramID,
		ChainSel:                  in.ChainSel,
		FeedAdmins:                in.FeedAdmins,
		ForwarderProgramID:        in.ForwarderProgramID,
		DefaultMaxStoredReports:   in.DefaultMaxStoredReports,
		MaxObservationForwardSecs: in.MaxObservationForwardSecs,
	})
	if err != nil {
		return out, err
	}
	out.State = initOut.Output.StatePubKey

	return out, nil
}
