package solana

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal/solana"
)

var _ cldf.ChangeSet[*DeployRequest] = DeployForwarder

func DeployForwarder(env deployment.Environment, req *DeployRequest) (cldf.ChangesetOutput, error) {
	req.deployFn = solana.DeployForwarder

	return deploy(env, req)
}
