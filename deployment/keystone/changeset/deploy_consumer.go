package changeset

import (
	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

type DeployFeedsConsumerRequest struct {
	ChainSelector uint64
}

var _ deployment.ChangeSet[*DeployFeedsConsumerRequest] = DeployFeedsConsumer

// DeployFeedsConsumer deploys the FeedsConsumer contract to the chain with the given chainSelector.
func DeployFeedsConsumer(env deployment.Environment, req *DeployFeedsConsumerRequest) (deployment.ChangesetOutput, error) {
	return DeployFeedsConsumerV2(env, &DeployRequestV2{
		ChainSel: req.ChainSelector,
	})
}

func DeployFeedsConsumerV2(env deployment.Environment, req *DeployRequestV2) (deployment.ChangesetOutput, error) {
	req.deployFn = kslib.DeployFeedsConsumer
	return deploy(env, req)
}
