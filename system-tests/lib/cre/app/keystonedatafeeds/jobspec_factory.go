package keystonedatafeeds

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	mockcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	mock_llo "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/triggers/llo"
	cretypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func JobSpecFactoryGenerator(feedsAddresses [][]mock_llo.FeedWithStreamID, mockCapabilities []*mockcapability.MockCapabilities) func(input *cretypes.JobSpecFactoryInput) (cretypes.DonsToJobSpecs, error) {
	return func(input *cretypes.JobSpecFactoryInput) (cretypes.DonsToJobSpecs, error) {
		donTojobSpecs := make(cretypes.DonsToJobSpecs, 0)

		for _, donWithMetadata := range input.DonTopology.DonsWithMetadata {
			jobSpecs := make(cretypes.DonJobs, 0)
			workflowNodeSet, err2 := node.FindManyWithLabel(donWithMetadata.NodesMetadata, &cretypes.Label{Key: node.NodeTypeKey, Value: cretypes.WorkerNode}, node.EqualLabels)
			if err2 != nil {
				// there should be no DON without worker nodes, even gateway DON is composed of a single worker node
				return nil, fmt.Errorf("failed to find worker nodes: %w", err2)
			}
			for _, workerNode := range workflowNodeSet {
				nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
				if nodeIDErr != nil {
					return nil, fmt.Errorf("failed to get node id from labels: %w", nodeIDErr)
				}
				if flags.HasFlag(donWithMetadata.Flags, cretypes.WorkflowDON) {
					for i := range feedsAddresses {
						feedConfig := make([]mock_llo.FeedConfig, 0)
						for _, feed := range feedsAddresses[i] {
							feedConfig = append(feedConfig, feed.MustFeedConfig())
						}
						jobSpecs = append(jobSpecs, mock_llo.WorkflowsJob(nodeID, fmt.Sprintf("load_%d", i), feedConfig))
					}
				}

				if flags.HasFlag(donWithMetadata.Flags, cretypes.MockCapability) && mockCapabilities != nil {
					jobSpecs = append(jobSpecs, mockcapability.MockCapabilitiesJob(nodeID, "mock", mockCapabilities))
				}
			}

			donTojobSpecs[donWithMetadata.ID] = jobSpecs
		}

		return donTojobSpecs, nil
	}
}
