package solana

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	ks_forwarder "github.com/smartcontractkit/chainlink-solana/contracts/generated/keystone_forwarder"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal/solana"
)

var _ cldf.ChangeSet[*DeployRequest] = DeployForwarder

func DeployForwarder(env cldf.Environment, req *DeployRequest) (cldf.ChangesetOutput, error) {
	if req.BuildConfig != nil {
		err := BuildSolana(env, req.BuildConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	req.deployFn = solana.DeployForwarder
	return deploy(env, req)
}

type InitializeForwardContractsRequest struct {
}

var _ cldf.ChangeSet[InitializeForwardContractsRequest] = InitializeForwarderContract

func InitializeForwarderContract(env cldf.Environment, req InitializeForwardContractsRequest) (cldf.ChangesetOutput, error) {
	instruction, err := ks_forwarder.NewInitializeInstruction()
}
