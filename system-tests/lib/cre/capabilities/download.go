package capabilities

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	"github.com/smartcontractkit/chainlink/system-tests/lib/types"
)

func DownloadCapabilityFromRelease(ghToken, version, assetFileName string) (string, error) {
	ghClient := client.NewGithubClient(ghToken)
	content, err := ghClient.DownloadAssetFromRelease("smartcontractkit", "capabilities", version, assetFileName)
	if err != nil {
		return "", err
	}

	fileName := assetFileName
	file, err := os.Create(assetFileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err = file.Write(content); err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func DefaultContainerDirectory(infraType types.InfraType) (string, error) {
	switch infraType {
	case types.CRIB:
		// chainlink user will always have access to this directory
		return "/home/chainlink", nil
	case types.Docker:
		// needs to match what CTFv2 uses by default, we should define a constant there and import it here
		return "/home/capabilities", nil
	default:
		return "", fmt.Errorf("unknown infra type: %s", infraType)
	}
}
