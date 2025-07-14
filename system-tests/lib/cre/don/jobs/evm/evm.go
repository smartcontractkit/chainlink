package evmJob

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	libjobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

// todo lautaro marker
var EVMJobSpecFactoryFn = func(chainID int, networkFamily string, pollInterval time.Duration, evmBinaryPath string) types.JobSpecFactoryFn {
	return func(input *types.JobSpecFactoryInput) (types.DonsToJobSpecs, error) {
		return GenerateJobSpecs(input.DonTopology, chainID, networkFamily, pollInterval, evmBinaryPath)
	}
}

var jobName = func(chainID int) string {
	return fmt.Sprintf("evm-capability-%d-%s", chainID)
}

func GenerateJobSpecs(donTopology *types.DonTopology, chainID int, networkFamily string, interval time.Duration, evmBinaryPath string) (types.DonsToJobSpecs, error) {
	if donTopology == nil {
		return nil, errors.New("topology is nil")
	}
	donToJobSpecs := make(types.DonsToJobSpecs)

	for _, donWithMetadata := range donTopology.DonsWithMetadata {
		if !flags.HasFlag(donWithMetadata.Flags, types.EVMCapability) {
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

			if evmBinaryPath == "" {
				return nil, errors.New("log event trigger binary path is empty")
			}
			//todo lautaro marker (print the config for evm)
			jobSpec := libjobs.WorkerStandardCapability(nodeID, jobName(chainID), evmBinaryPath,
				fmt.Sprintf(
					`'{"chainId":"%d","network":"%s","logTriggerPollInterval":%s,"CREForwarderAddress":%s,"ReceiverGasMinimum":%d,"NodeAddress":%s}'`,
					//`'{"chainId":"%d","network":"%s","logTriggerPollInterval":%s}'`,
					//	CREForwarderAddress:    "1234567890abcdef1234567890abcdef12345678", //fake address for testing
					//		ReceiverGasMinimum:     1,
					//		NodeAddress:            "fakeAddressForTesting", //fake address for testing
					chainID,
					networkFamily,
					interval,
					"1234567890abcdef1234567890abcdef12345678", // fake address for testing
					1, // ReceiverGasMinimum
					"fakeAddressForTesting",
				),
			)

			if _, ok := donToJobSpecs[donWithMetadata.ID]; !ok {
				donToJobSpecs[donWithMetadata.ID] = make(types.DonJobs, 0)
			}

			donToJobSpecs[donWithMetadata.ID] = append(donToJobSpecs[donWithMetadata.ID], jobSpec)
		}
	}

	return donToJobSpecs, nil
}
