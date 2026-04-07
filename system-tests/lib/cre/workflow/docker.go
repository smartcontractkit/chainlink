package workflow

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
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

	start := time.Now()
	framework.L.Info().
		Int("file_count", len(existingFiles)).
		Str("container_pattern", containerNamePattern).
		Msg("Copying workflow artifacts to Docker containers (parallel)")

	eg := errgroup.Group{}
	eg.SetLimit(4)
	for _, file := range existingFiles {
		fileCopy := file
		eg.Go(func() error {
			return errors.Wrapf(
				CopyFilesToDockerContainers(context.Background(), containerNames, containerTargetDir, []string{fileCopy}),
				"failed to copy a file (%s) to docker containers", fileCopy,
			)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	framework.L.Info().
		Dur("duration", time.Since(start)).
		Msg("Workflow artifacts copied to Docker containers")
	return nil
}

func CopyFilesToDockerContainers(ctx context.Context, containerNames []string, targetDir string, files []string) error {
	return dockerops.CopyFilesToContainers(ctx, containerNames, targetDir, files)
}
