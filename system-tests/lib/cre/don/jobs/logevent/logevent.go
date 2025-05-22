package logevent

import (
	"fmt"

	"github.com/pkg/errors"

	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

var LogEventTriggerJobSpecFactoryFn = func(chainID int, networkFamily, logEventTriggerBinaryPath string) types.JobSpecFactoryFn {
	return func(input *types.JobSpecFactoryInput) (types.DonsToJobSpecs, error) {
		return GenerateJobSpecs(input.DonTopology, chainID, networkFamily, logEventTriggerBinaryPath)
	}
}

var LogEventTriggerJobName = func(chainID int) string {
	return fmt.Sprintf("log-event-trigger-%d", chainID)
}

func GenerateJobSpecs(donTopology *types.DonTopology, chainID int, networkFamily, logEventTriggerBinaryPath string) (types.DonsToJobSpecs, error) {
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

			if flags.HasFlag(donWithMetadata.Flags, types.LogTriggerCapability) {
				if logEventTriggerBinaryPath == "" {
					return nil, errors.New("log event trigger binary path is empty")
				}

				jobSpec := libjobs.WorkerStandardCapability(nodeID, LogEventTriggerJobName(chainID), logEventTriggerBinaryPath, fmt.Sprintf(`'{"chainId":"%d","network":"%s","lookbackBlocks":1000,"pollPeriod":1000}'`, chainID, networkFamily))

				if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
					donToJobSpecs[donWithMetadata.ID] = make(types.DonJobs, 0)
				}

				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
			}
		}
	}

	return donToJobSpecs, nil
}
