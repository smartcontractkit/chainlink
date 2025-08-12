package capabilities

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	libnode "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/node"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

func MakeBinariesExecutable(customBinariesPaths map[cre.CapabilityFlag]string) error {
	for capabilityFlag, binaryPath := range customBinariesPaths {
		if binaryPath == "" {
			return fmt.Errorf("binary path for capability %s is empty", capabilityFlag)
		}

		// Check if file exists
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			return fmt.Errorf("binary path %s for capability %s does not exist", binaryPath, capabilityFlag)
		}

		// Make the binary executable
		if err := os.Chmod(binaryPath, 0755); err != nil {
			return fmt.Errorf("failed to make binary %s executable for capability %s: %w", binaryPath, capabilityFlag, err)
		}
	}

	return nil
}

func AppendBinariesPathsNodeSpec(nodeSetInput *cre.CapabilitiesAwareNodeSet, donMetadata *cre.DonMetadata, customBinariesPaths map[cre.CapabilityFlag]string) (*cre.CapabilitiesAwareNodeSet, error) {
	if len(customBinariesPaths) == 0 {
		return nodeSetInput, nil
	}

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
				workerNodes, wErr := libnode.FindManyWithLabel(donMetadata.NodesMetadata, &cre.Label{
					Key:   libnode.NodeTypeKey,
					Value: cre.WorkerNode,
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
