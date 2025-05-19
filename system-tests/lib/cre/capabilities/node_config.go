package capabilities

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
)

func AppendBinariesPathsNodeSpec(nodeSetInput *types.CapabilitiesAwareNodeSet, donMetadata *types.DonMetadata, customBinariesPaths map[types.CapabilityFlag]string) (*types.CapabilitiesAwareNodeSet, error) {
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
		for capabilityFlag, binaryPath := range customBinariesPaths {
			if binaryPath == "" {
				return nil, fmt.Errorf("binary path for capability %s is empty", capabilityFlag)
			}

			if flags.HasFlag(donMetadata.Flags, capabilityFlag) {
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
						return nil, errors.Wrap(nErr, "failed to find index label")
					}

					nodeIndex, nIErr := strconv.Atoi(nodeIndexStr)
					if nIErr != nil {
						return nil, errors.Wrap(nIErr, "failed to convert index label value to int")
					}

					nodeSetInput.NodeSpecs[nodeIndex].Node.CapabilitiesBinaryPaths = append(nodeSetInput.NodeSpecs[nodeIndex].Node.CapabilitiesBinaryPaths, binaryPath)
				}
			}
		}
	}

	return nodeSetInput, nil
}
