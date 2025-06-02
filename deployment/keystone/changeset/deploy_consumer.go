package changeset

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	kslib "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

type DeployFeedsConsumerRequest struct {
	ChainSelector uint64
}

var _ cldf.ChangeSet[*DeployFeedsConsumerRequest] = DeployFeedsConsumer

// DeployFeedsConsumer deploys the FeedsConsumer contract to the chain with the given chainSelector.
func DeployFeedsConsumer(env cldf.Environment, req *DeployFeedsConsumerRequest) (cldf.ChangesetOutput, error) {
	return DeployFeedsConsumerV2(env, &DeployRequestV2{
		ChainSel: req.ChainSelector,
	})
}

func DeployFeedsConsumerV2(env cldf.Environment, req *DeployRequestV2) (cldf.ChangesetOutput, error) {
	req.deployFn = kslib.DeployFeedsConsumer
	return deploy(env, req)
}
