package mock

import (
	"fmt"

	"github.com/pkg/errors"

	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

var MockJobSpecFactoryFn = func(port int) types.JobSpecFactoryFn {
	return func(input *types.JobSpecFactoryInput) (types.DonsToJobSpecs, error) {
		return GenerateJobSpecs(input.DonTopology, port)
	}
}

var MockJobName = func(chainID int) string {
	return fmt.Sprintf("mock-%d", chainID)
}

func GenerateJobSpecs(donTopology *types.DonTopology, port int) (types.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(types.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, types.MockCapability) {
			continue
		}

		workflowNodeSet, err := libnode.FindManyWithLabel(donWithMetadata.NodesMetadata, &types.Label{Key: libnode.NodeTypeKey, Value: types.WorkerNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := libnode.FindLabelValue(workerNode, libnode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			jobSpec := libjobs.WorkerStandardCapability(nodeID, "mock-cap", "mock", fmt.Sprintf(`"""port=%d"""`, port))

			if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
				donToJobSpecs[donWithMetadata.ID] = make(types.DonJobs, 0)
			}

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
		}
	}

	return donToJobSpecs, nil
}
