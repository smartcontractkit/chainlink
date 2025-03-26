package chainreader

import (
	"fmt"

	"github.com/pkg/errors"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func GenerateJobSpecs(donTopology *types.DonTopology, chainID int, networkFamily, logEventTriggerBinaryPath, readContractBinaryPath string) (types.DonsToJobSpecs, error) {
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

				jobSpec := libjobs.WorkerStandardCapability(nodeID, fmt.Sprintf("log-event-trigger-capability-%d", chainID), logEventTriggerBinaryPath, fmt.Sprintf(`'{"chainId":"%d","network":"%s","lookbackBlocks":1000,"pollPeriod":1000}'`, chainID, networkFamily))
				jobDesc := types.JobDescription{Flag: types.LogTriggerCapability, NodeType: types.WorkerNode}

				if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
					donToJobSpecs[donWithMetadata.ID] = make(map[types.JobDescription][]*jobv1.ProposeJobRequest)
				}

				if _, ok := donToJobSpecs[donWithMetadata.ID][jobDesc]; !ok {
					donToJobSpecs[donWithMetadata.ID][jobDesc] = []*jobv1.ProposeJobRequest{jobSpec}
				} else {
					donToJobSpecs[donWithMetadata.ID][jobDesc] = append(donToJobSpecs[donWithMetadata.ID][jobDesc], jobSpec)
				}
			}

			if flags.HasFlag(donWithMetadata.Flags, types.ReadContractCapability) {
				if readContractBinaryPath == "" {
					return nil, errors.New("read contract binary path is empty")
				}

				jobSpec := libjobs.WorkerStandardCapability(nodeID, fmt.Sprintf("read-contract-capability-%d", chainID), readContractBinaryPath, fmt.Sprintf(`'{"chainId":%d,"network":"%s"}'`, chainID, networkFamily))
				jobDesc := types.JobDescription{Flag: types.LogTriggerCapability, NodeType: types.WorkerNode}

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
