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
	for _, file := range filesToCopy {
		if _, err := os.Stat(file); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: File '%s' does not exist. Skipping file copying to docker containers\n", file)
			continue
		}

		workflowCopyErr := copyArtifactToDockerContainers(file, containerNamePattern, containerTargetDir)
		if workflowCopyErr != nil {
			return errors.Wrapf(workflowCopyErr, "failed to copy a file (%s) to docker containers", file)
		}
	}
	return nil
}

func copyArtifactToDockerContainers(filePath string, containerNamePattern string, targetDir string) error {
	framework.L.Info().Msgf("Copying file '%s' to Docker containers", filePath)
	containerNames, containerNamesErr := findAllDockerContainerNames(containerNamePattern)
	if containerNamesErr != nil {
		return errors.Wrap(containerNamesErr, "failed to find Docker containers")
	}

	if len(containerNames) == 0 {
		return fmt.Errorf("no Docker containers found with name pattern %s", containerNamePattern)
	}

	frameworkDockerClient, frameworkDockerClientErr := framework.NewDockerClient()
	if frameworkDockerClientErr != nil {
		return errors.Wrap(frameworkDockerClientErr, "failed to create framework Docker client")
	}

	for _, containerName := range containerNames {
		execOutput, execOutputErr := frameworkDockerClient.ExecContainer(containerName, []string{"mkdir", "-p", targetDir})
		if execOutputErr != nil {
			fmt.Fprint(os.Stderr, execOutput)
			return errors.Wrap(execOutputErr, "failed to execute mkdir command in Docker container")
		}

		copyErr := frameworkDockerClient.CopyFile(containerName, filePath, targetDir)
		if copyErr != nil {
			fmt.Fprint(os.Stderr, execOutput)
			return errors.Wrap(copyErr, "failed to copy artifact to Docker container")
		}

		dockerClient, dockerClientErr := dc.NewClientWithOpts(dc.FromEnv, dc.WithAPIVersionNegotiation())
		if dockerClientErr != nil {
			return errors.Wrap(dockerClientErr, "failed to create Docker client")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		containerJSON, ispectErr := dockerClient.ContainerInspect(ctx, containerName)
		if ispectErr != nil {
			cancel()
			return errors.Wrap(ispectErr, "failed to inspect Docker container")
		}
		cancel()
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
			execOutput, execOutputErr := frameworkDockerClient.ExecContainerOptions(containerName, execConfig)
			if execOutputErr != nil {
				fmt.Fprint(os.Stderr, execOutput)
				return errors.Wrap(execOutputErr, "failed to execute mkdir command in Docker container")
			}
			fmt.Println("output " + execOutput)
		}
	}

	return nil
}
