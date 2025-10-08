package capabilities

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func MakeBinariesExecutable(customBinariesPaths map[cre.CapabilityFlag]string) error {
	for capabilityFlag, binaryPath := range customBinariesPaths {
		if binaryPath == "" {
			return fmt.Errorf("binary path for capability %s is empty. Please set the binary path in the capabilities TOML config", capabilityFlag)
		}

		// Check if file exists
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			absPath, absErr := filepath.Abs(binaryPath)
			if absErr != nil {
				return errors.Wrapf(absErr, "failed to get absolute path for binary %s", binaryPath)
			}

			return fmt.Errorf("no binary file for capability %s not found at '%s'. Please make sure the path is correct, update it in the capabilities TOML config or copy the binary to the expected location", absPath, capabilityFlag)
		}

		// Make the binary executable
		if err := os.Chmod(binaryPath, 0755); err != nil {
			return errors.Wrapf(err, "failed to make binary %s executable for capability %s", binaryPath, capabilityFlag)
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
				return nil, fmt.Errorf("binary path for capability %s is empty. Make sure you have set the binary path in the TOML config", capabilityFlag)
			}

			workerNodes, wErr := donMetadata.Workers()
			if wErr != nil {
				return nil, errors.Wrap(wErr, "failed to find worker nodes")
			}

			for _, workerNode := range workerNodes {
				nodeSetInput.NodeSpecs[workerNode.Index].Node.CapabilitiesBinaryPaths = append(nodeSetInput.NodeSpecs[workerNode.Index].Node.CapabilitiesBinaryPaths, binaryPath)
			}
		}
	}

	return nodeSetInput, nil
}
