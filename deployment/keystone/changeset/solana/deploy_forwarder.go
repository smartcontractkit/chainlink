package solana

import (
	ks_forwarder "github.com/smartcontractki/chainlink-solana/contracts/generated/keystone_forwarder"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
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
	cl := env.SolChains[0].Client
	ks_forwarder.Initialize()
	return cldf.ChangesetOutput{}, nil
}
