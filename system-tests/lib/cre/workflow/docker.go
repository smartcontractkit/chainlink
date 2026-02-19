package workflow

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"os"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/internal/dockerops"
)

var (
	DefaultWorkflowTargetDir   = "/home/chainlink/workflows"
	DefaultWorkflowNodePattern = "workflow-node"
)

func findAllDockerContainerNames(pattern string) ([]string, error) {
	return FindDockerContainerNames(context.Background(), pattern)
}

func FindDockerContainerNames(ctx context.Context, pattern string) ([]string, error) {
	return dockerops.FindContainerNames(ctx, pattern)
}

func CopyArtifactsToDockerContainers(containerTargetDir string, containerNamePattern string, filesToCopy ...string) error {
	containerNames, containerNamesErr := findAllDockerContainerNames(containerNamePattern)
	if containerNamesErr != nil {
		return errors.Wrap(containerNamesErr, "failed to find Docker containers")
	}
	if len(containerNames) == 0 {
		return fmt.Errorf("no Docker containers found with name pattern %s", containerNamePattern)
	}

	existingFiles := make([]string, 0, len(filesToCopy))
	for _, file := range filesToCopy {
		if _, err := os.Stat(file); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: File '%s' does not exist. Skipping file copying to docker containers\n", file)
			continue
		}
		existingFiles = append(existingFiles, file)
	}
	if len(existingFiles) == 0 {
		return nil
	}
	return CopyFilesToDockerContainers(context.Background(), containerNames, containerTargetDir, existingFiles)
}

func CopyFilesToDockerContainers(ctx context.Context, containerNames []string, targetDir string, files []string) error {
	return dockerops.CopyFilesToContainers(ctx, containerNames, targetDir, files)
}
