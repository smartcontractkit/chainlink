package crecli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/client"
)

func DownloadAndInstallChainlinkCLI(ghToken, version string) (string, error) {
	system := runtime.GOOS
	arch := runtime.GOARCH

	switch system {
	case "darwin", "linux":
		// nothing to do, we have the binaries
	default:
		return "", fmt.Errorf("chainlnk-cli does not support OS: %s", system)
	}

	switch arch {
	case "amd64", "arm64":
		// nothing to do, we have the binaries
	default:
		return "", fmt.Errorf("chainlnk-cli does not support arch: %s", arch)
	}

	creCLIAssetFile := fmt.Sprintf("cre_%s_%s_%s.tar.gz", version, system, arch)

	ghClient := client.NewGithubClient(ghToken)
	content, err := ghClient.DownloadAssetFromRelease("smartcontractkit", "dev-platform", version, creCLIAssetFile)
	if err != nil {
		return "", errors.Wrapf(err, "failed to download CRE CLI asset %s", creCLIAssetFile)
	}

	tmpfile, err := os.CreateTemp("", creCLIAssetFile)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create temp file for CRE CLI asset %s", creCLIAssetFile)
	}
	defer tmpfile.Close()

	if _, err = tmpfile.Write(content); err != nil {
		return "", errors.Wrapf(err, "failed to write content to temp file for CRE CLI asset %s", creCLIAssetFile)
	}

	cmd := exec.Command("tar", "-xvf", tmpfile.Name(), "-C", ".") // #nosec G204
	if cmd.Run() != nil {
		return "", errors.Wrapf(err, "failed to extract CRE CLI asset %s", creCLIAssetFile)
	}

	extractedFileName := fmt.Sprintf("cre_%s_%s_%s", version, system, arch)
	cmd = exec.Command("chmod", "+x", extractedFileName)
	if cmd.Run() != nil {
		return "", errors.Wrapf(err, "failed to make %s executable", extractedFileName)
	}

	// set it to absolute path, because some commands (e.g. compile) need to be executed in the context
	// of the workflow directory
	extractedFile, err := os.Open(extractedFileName)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open %s", extractedFileName)
	}

	creCLICommandPath, err := filepath.Abs(extractedFile.Name())
	if err != nil {
		return "", errors.Wrapf(err, "failed to get absolute path for %s", tmpfile.Name())
	}

	return creCLICommandPath, nil
}
