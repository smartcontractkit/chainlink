package webapi

import (
	"github.com/pkg/errors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func GenerateJobSpecs(donTopology *types.DonTopology) (types.DonsToJobSpecs, error) {
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
				jobSpec := libjobs.WorkerStandardCapability(nodeID, "web-api-trigger-capability", "__builtin_web-api-trigger", libjobs.EmptyStdCapConfig)
				jobDesc := types.JobDescription{Flag: types.WebAPITargetCapability, NodeType: types.WorkerNode}

				if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
					donToJobSpecs[donWithMetadata.ID] = make(map[types.JobDescription][]*jobv1.ProposeJobRequest)
				}

				if _, ok := donToJobSpecs[donWithMetadata.ID][jobDesc]; !ok {
					donToJobSpecs[donWithMetadata.ID][jobDesc] = []*jobv1.ProposeJobRequest{jobSpec}
				} else {
					donToJobSpecs[donWithMetadata.ID][jobDesc] = append(donToJobSpecs[donWithMetadata.ID][jobDesc], jobSpec)
				}
			}

			if flags.HasFlag(donWithMetadata.Flags, types.WebAPITargetCapability) {
				config := `"""
						[rateLimiter]
						GlobalRPS = 1000.0
						GlobalBurst = 1000
						PerSenderRPS = 1000.0
						PerSenderBurst = 1000
						"""`

				jobSpec := libjobs.WorkerStandardCapability(nodeID, "web-api-target-capability", "__builtin_web-api-target", config)
				jobDesc := types.JobDescription{Flag: types.WebAPITargetCapability, NodeType: types.WorkerNode}

				if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
					donToJobSpecs[donWithMetadata.ID] = make(map[types.JobDescription][]*jobv1.ProposeJobRequest)
				}

				if _, ok := donToJobSpecs[donWithMetadata.ID][jobDesc]; !ok {
					donToJobSpecs[donWithMetadata.ID][jobDesc] = []*jobv1.ProposeJobRequest{jobSpec}
				} else {
					donToJobSpecs[donWithMetadata.ID][jobDesc] = append(donToJobSpecs[donWithMetadata.ID][jobDesc], jobSpec)
				}
			}
		}
	}

	return donToJobSpecs, nil
}
