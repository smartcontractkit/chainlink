package workflow

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	ctypes "github.com/docker/docker/api/types/container"
	dc "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

var (
	DefaultWorkflowTargetDir   = "/home/chainlink/workflows"
	DefaultWorkflowNodePattern = "workflow-node"
)

func findAllDockerContainerNames(pattern string) ([]string, error) {
	dockerClient, dockerClientErr := dc.NewClientWithOpts(dc.FromEnv, dc.WithAPIVersionNegotiation())
	if dockerClientErr != nil {
		return nil, errors.Wrap(dockerClientErr, "failed to create Docker client")
	}
	defer dockerClient.Close()

	containers, containersErr := dockerClient.ContainerList(context.Background(), ctypes.ListOptions{})
	if containersErr != nil {
		return nil, errors.Wrap(containersErr, "failed to list Docker containers")
	}

	containerNames := []string{}
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, pattern) {
				// Remove leading slash from container name
				cleanName := strings.TrimPrefix(name, "/")
				containerNames = append(containerNames, cleanName)
			}
		}
	}

	return containerNames, nil
}

func CopyArtifactsToDockerContainers(containerTargetDir string, containerNamePattern string, filesToCopy ...string) error {
	start := time.Now()
	framework.L.Info().
		Int("file_count", len(filesToCopy)).
		Str("container_pattern", containerNamePattern).
		Msg("Copying workflow artifacts to Docker containers (parallel)")

	eg := errgroup.Group{}
	eg.SetLimit(4)
	for _, file := range filesToCopy {
		if _, err := os.Stat(file); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: File '%s' does not exist. Skipping file copying to docker containers\n", file)
			continue
		}
		eg.Go(func() error {
			return errors.Wrapf(
				copyArtifactToDockerContainers(file, containerNamePattern, containerTargetDir),
				"failed to copy a file (%s) to docker containers", file,
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

func copyArtifactToDockerContainers(filePath string, containerNamePattern string, targetDir string) error {
	framework.L.Info().Msgf("Copying file '%s' to Docker containers", filePath)
	containerNames, err := findAllDockerContainerNames(containerNamePattern)
	if err != nil {
		return errors.Wrap(err, "failed to find Docker containers")
	}
	if len(containerNames) == 0 {
		return fmt.Errorf("no Docker containers found with name pattern %s", containerNamePattern)
	}

	frameworkDockerClient, err := framework.NewDockerClient()
	if err != nil {
		return errors.Wrap(err, "failed to create framework Docker client")
	}
	dockerClient, err := dc.NewClientWithOpts(dc.FromEnv, dc.WithAPIVersionNegotiation())
	if err != nil {
		return errors.Wrap(err, "failed to create Docker client")
	}
	defer dockerClient.Close()

	eg := errgroup.Group{}
	eg.SetLimit(4)
	for _, containerName := range containerNames {
		eg.Go(func() error {
			execOutput, execErr := frameworkDockerClient.ExecContainer(containerName, []string{"mkdir", "-p", targetDir})
			if execErr != nil {
				fmt.Fprint(os.Stderr, execOutput)
				return errors.Wrap(execErr, "failed to execute mkdir command in Docker container")
			}
			if copyErr := frameworkDockerClient.CopyFile(containerName, filePath, targetDir); copyErr != nil {
				return errors.Wrap(copyErr, "failed to copy artifact to Docker container")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			containerJSON, inspectErr := dockerClient.ContainerInspect(ctx, containerName)
			if inspectErr != nil {
				return errors.Wrap(inspectErr, "failed to inspect Docker container")
			}
			user := containerJSON.Config.User
			// if not running as root, change ownership to user that is running the container to avoid permission issues
			if user != "" {
				targetFilePath := filepath.Join(targetDir, filepath.Base(filePath))
				execConfig := ctypes.ExecOptions{
					Cmd:          []string{"chown", user, targetFilePath},
					AttachStdout: true,
					AttachStderr: true,
					User:         "root",
				}
				execOutput, execErr = frameworkDockerClient.ExecContainerOptions(containerName, execConfig)
				if execErr != nil {
					fmt.Fprint(os.Stderr, execOutput)
					return errors.Wrap(execErr, "failed to execute chown command in Docker container")
				}
				framework.L.Debug().Str("container", containerName).Msgf("chown output: %s", execOutput)
			}
			return nil
		})
	}
	return eg.Wait()
}
