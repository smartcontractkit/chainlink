package workflow

import (
	"context"
	"fmt"
	"os"
	"time"

	"path"
	"path/filepath"

	dc "github.com/moby/moby/client"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/internal/dockerops"
)

var (
	DefaultWorkflowTargetDir   = "/home/chainlink/workflows"
	DefaultWorkflowNodePattern = "workflow-node"
)

func FindDockerContainerNames(ctx context.Context, pattern string) ([]string, error) {
	return dockerops.FindDockerContainerNames(ctx, pattern)
}

func CopyArtifactsToDockerContainers(containerTargetDir string, containerNamePattern string, filesToCopy ...string) error {
	containerNames, containerNamesErr := dockerops.FindDockerContainerNames(context.Background(), containerNamePattern)
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
				dockerops.CopyFilesToContainers(context.Background(), containerNamePattern, containerTargetDir, []string{fileCopy}),
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

// CopyAndExtractTarballToDockerContainers copies one tarball into each matching container, extracts it
// with GNU tar into containerTargetDir, removes the tarball, and chowns containerTargetDir recursively
// when the container runs as a non-root user.
func CopyAndExtractTarballToDockerContainers(containerTargetDir, containerNamePattern, tarballHostPath string) error {
	frameworkDockerClient, err := framework.NewDockerClient()
	if err != nil {
		return errors.Wrap(err, "failed to create framework Docker client")
	}
	dockerClient, err := dc.New(dc.FromEnv)
	if err != nil {
		return errors.Wrap(err, "failed to create Docker client")
	}
	defer dockerClient.Close()

	containerNames, err := FindDockerContainerNames(context.Background(), containerNamePattern)
	if err != nil {
		return errors.Wrap(err, "failed to find Docker containers")
	}
	if len(containerNames) == 0 {
		return fmt.Errorf("no Docker containers found with name pattern %s", containerNamePattern)
	}

	tarBase := filepath.Base(tarballHostPath)
	tarContainerPath := path.Join(containerTargetDir, tarBase)

	eg := errgroup.Group{}
	eg.SetLimit(4)
	for _, containerName := range containerNames {
		eg.Go(func() error {
			execOutput, execErr := frameworkDockerClient.ExecContainer(containerName, []string{"mkdir", "-p", containerTargetDir})
			if execErr != nil {
				fmt.Fprint(os.Stderr, execOutput)
				return errors.Wrap(execErr, "failed to execute mkdir in Docker container")
			}
			if copyErr := frameworkDockerClient.CopyFile(containerName, tarballHostPath, containerTargetDir); copyErr != nil {
				return errors.Wrap(copyErr, "failed to copy tarball to Docker container")
			}

			execOutput, execErr = frameworkDockerClient.ExecContainer(containerName, []string{"tar", "-xf", tarContainerPath, "-C", containerTargetDir})
			if execErr != nil {
				fmt.Fprint(os.Stderr, execOutput)
				return errors.Wrapf(execErr, "failed to extract tarball in %s", containerName)
			}

			rmOut, rmErr := frameworkDockerClient.ExecContainer(containerName, []string{"rm", "-f", tarContainerPath})
			if rmErr != nil {
				fmt.Fprint(os.Stderr, rmOut)
				return errors.Wrap(rmErr, "failed to remove tarball from Docker container")
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			inspectRes, inspectErr := dockerClient.ContainerInspect(ctx, containerName, dc.ContainerInspectOptions{})
			if inspectErr != nil {
				return errors.Wrap(inspectErr, "failed to inspect Docker container")
			}
			user := inspectRes.Container.Config.User
			if user == "" {
				return nil
			}
			execConfig := dc.ExecCreateOptions{
				Cmd:          []string{"chown", "-R", user, containerTargetDir},
				AttachStdout: true,
				AttachStderr: true,
				User:         "root",
			}
			chOut, chErr := frameworkDockerClient.ExecContainerOptions(containerName, execConfig)
			if chErr != nil {
				fmt.Fprint(os.Stderr, chOut)
				return errors.Wrapf(chErr, "chown %s in %s", containerTargetDir, containerName)
			}
			return nil
		})
	}

	return eg.Wait()
}
