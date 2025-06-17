package webapi

import (
	"github.com/pkg/errors"

	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

var WebAPITargetJobSpecFactoryFn = func(input *types.JobSpecFactoryInput) (types.DonsToJobSpecs, error) {
	return GenerateWebAPITargetJobSpecs(input.DonTopology)
}

func GenerateWebAPITargetJobSpecs(donTopology *types.DonTopology) (types.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(types.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		workflowNodeSet, err := libnode.FindManyWithLabel(donWithMetadata.NodesMetadata, &types.Label{Key: libnode.NodeTypeKey, Value: types.WorkerNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := libnode.FindLabelValue(workerNode, libnode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			if flags.HasFlag(donWithMetadata.Flags, types.WebAPITargetCapability) {
				config := `"""
						[rateLimiter]
						GlobalRPS = 1000.0
						GlobalBurst = 1000
						PerSenderRPS = 1000.0
						PerSenderBurst = 1000
						"""`

				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], libjobs.WorkerStandardCapability(nodeID, types.WebAPITargetCapability, "__builtin_web-api-target", config))
			}
		}
	}

	return donToJobSpecs, nil
}
