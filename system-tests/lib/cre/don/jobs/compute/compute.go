package compute

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	crenode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

var ComputeJobSpecFactoryFn = func(input *types.JobSpecFactoryInput) (types.DonsToJobSpecs, error) {
	return GenerateJobSpecs(input.DonTopology)
}

func GenerateJobSpecs(donTopology *types.DonTopology) (types.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(types.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		workflowNodeSet, err := crenode.FindManyWithLabel(donWithMetadata.NodesMetadata, &types.Label{Key: crenode.NodeTypeKey, Value: types.WorkerNode}, crenode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := crenode.FindLabelValue(workerNode, crenode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			if flags.HasFlag(donWithMetadata.Flags, types.CustomComputeCapability) {
				config := `"""
					NumWorkers = 3
					[rateLimiter]
					globalRPS = 20.0
					globalBurst = 30
					perSenderRPS = 1.0
					perSenderBurst = 5
					"""`
				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerStandardCapability(nodeID, types.CustomComputeCapability, "__builtin_custom-compute-action", config))
			}
		}
	}

	return donToJobSpecs, nil
}
