package cron

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

var CronJobSpecFactoryFn = func(cronBinaryPath string) types.JobSpecFactoryFn {
	return func(input *types.JobSpecFactoryInput) (types.DonsToJobSpecs, error) {
		return GenerateJobSpecs(
			input.DonTopology,
			cronBinaryPath,
		)
	}
}

func GenerateJobSpecs(donTopology *types.DonTopology, cronBinaryPath string) (types.DonsToJobSpecs, error) {
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
			nodeID, nodeIDErr := node.FindLabelValue(workerNode, node.NodeIDKey)
			if nodeIDErr != nil {
				return nil, errors.Wrap(nodeIDErr, "failed to get node id from labels")
			}

			if flags.HasFlag(donWithMetadata.Flags, types.CronCapability) {
				donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobs.WorkerStandardCapability(nodeID, types.CronCapability, cronBinaryPath, jobs.EmptyStdCapConfig))
			}
		}
	}

	return donToJobSpecs, nil
}
