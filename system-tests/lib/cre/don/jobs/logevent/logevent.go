package logevent

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

var LogEventTriggerJobSpecFactoryFn = func(chainID int, networkFamily, logEventTriggerBinaryPath string) cre.JobSpecFactoryFn {
	return func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
		return GenerateJobSpecs(input.DonTopology, chainID, networkFamily, logEventTriggerBinaryPath)
	}
}

var LogEventTriggerJobName = func(chainID int) string {
	return fmt.Sprintf("log-event-trigger-%d", chainID)
}

func GenerateJobSpecs(donTopology *cre.DonTopology, chainID int, networkFamily, logEventTriggerBinaryPath string) (cre.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(cre.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, cre.LogTriggerCapability) {
			continue
		}

		workflowNodeSet, err := libnode.FindManyWithLabel(donWithMetadata.NodesMetadata, &cre.Label{Key: libnode.NodeTypeKey, Value: cre.WorkerNode}, libnode.EqualLabels)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find worker nodes")
		}

		for _, workerNode := range workflowNodeSet {
			nodeID, nodeIDErr := libnode.FindLabelValue(workerNode, libnode.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			if logEventTriggerBinaryPath == "" {
				return nil, errors.New("log event trigger binary path is empty")
			}

			jobSpec := libjobs.WorkerStandardCapability(nodeID, LogEventTriggerJobName(chainID), logEventTriggerBinaryPath, fmt.Sprintf(`'{"chainId":"%d","network":"%s","lookbackBlocks":1000,"pollPeriod":1000}'`, chainID, networkFamily))

			if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
				donToJobSpecs[donWithMetadata.ID] = make(cre.DonJobs, 0)
			}

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
		}
	}

	return donToJobSpecs, nil
}
