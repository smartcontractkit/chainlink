package capabilities

import (
	"strconv"

	"github.com/pkg/errors"

	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

var DefaultBinariesPathsFactory = func(cronBinaryPath string) types.CapabilitiesBinaryPathFactoryFn {
	return func(donMetadata *types.DonMetadata) ([]string, error) {
		binaries := []string{}
		if flags.HasFlag(donMetadata.Flags, types.CronCapability) {
			binaries = append(binaries, cronBinaryPath)
		}

		return binaries, nil
	}
}

func AppendBinariesPathsNodeSpec(nodeSetInput *types.CapabilitiesAwareNodeSet, donMetadata *types.DonMetadata, pathFactoryFns []types.CapabilitiesBinaryPathFactoryFn) (*types.CapabilitiesAwareNodeSet, error) {
	// if no capabilities are defined in TOML, but DON has ones that we know require custom binaries
	// append them to the node specification
	hasCapabilitiesBinaries := false
	for _, nodeInput := range nodeSetInput.NodeSpecs {
		if len(nodeInput.Node.CapabilitiesBinaryPaths) > 0 {
			hasCapabilitiesBinaries = true
			break
		}
	}

	if !hasCapabilitiesBinaries {
		binariesToAppend := []string{}
		for _, pathFactoryFn := range pathFactoryFns {
			binaries, err := pathFactoryFn(donMetadata)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get capabilities binaries' paths")
			}

			binariesToAppend = append(binariesToAppend, binaries...)
		}

		workerNodes, wErr := libnode.FindManyWithLabel(donMetadata.NodesMetadata, &types.Label{
			Key:   libnode.NodeTypeKey,
			Value: types.WorkerNode,
		}, libnode.EqualLabels)

		if wErr != nil {
			return nil, errors.Wrap(wErr, "failed to find worker nodes")
		}

		for _, node := range workerNodes {
			nodeIndexStr, nErr := libnode.FindLabelValue(node, libnode.IndexKey)
			if nErr != nil {
				return nil, errors.Wrap(nErr, "failed to find node index in labels")
			}

			nodeIndex, nIErr := strconv.Atoi(nodeIndexStr)
			if nIErr != nil {
				return nil, errors.Wrap(nIErr, "failed to convert index to int")
			}

			nodeSetInput.NodeSpecs[nodeIndex].Node.CapabilitiesBinaryPaths = append(nodeSetInput.NodeSpecs[nodeIndex].Node.CapabilitiesBinaryPaths, binariesToAppend...)
		}
	}

	return nodeSetInput, nil
}
