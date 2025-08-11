package httpaction

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

const serviceConfigTemplate = `"""
{
	"proxyMode": "gateway",
	"incomingRateLimiter": {
		"globalBurst": 10,
		"globalRPS": 50,
		"perSenderBurst": 10,
		"perSenderRPS": 10
	},
	"outgoingRateLimiter": {
		"globalBurst": 10,
		"globalRPS": 50,
		"perSenderBurst": 10,
		"perSenderRPS": 10
	}
}
"""
`

var HTTPActionJobSpecFactoryFn = func(httpActionBinaryPath string) cre.JobSpecFactoryFn {
	return func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
		return GenerateJobSpecs(
			input.DonTopology,
			httpActionBinaryPath,
		)
	}
}

func GenerateJobSpecs(donTopology *cre.DonTopology, httpActionBinaryPath string) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.HTTPActionCapability) {
			continue
		}
		workflowNodeSet, err := crenode.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: crenode.NodeTypeKey, Value: cre.WorkerNode}, crenode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := crenode.FindLabelValue(workerNode, crenode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerStandardCapability(nodeID, cre.HTTPActionCapability, httpActionBinaryPath, serviceConfigTemplate, ""))
		}
	}

	return donToJobSpecs, nil
}
