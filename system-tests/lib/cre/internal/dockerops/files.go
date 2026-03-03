package dockerops

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

func FindContainerNames(ctx context.Context, pattern string) ([]string, error) {
	dockerClient, dockerClientErr := dc.NewClientWithOpts(dc.FromEnv, dc.WithAPIVersionNegotiation())
	if dockerClientErr != nil {
		return nil, errors.Wrap(dockerClientErr, "failed to create Docker client")
	}
	defer dockerClient.Close()

	containers, containersErr := dockerClient.ContainerList(ctx, ctypes.ListOptions{})
	if containersErr != nil {
		return nil, errors.Wrap(containersErr, "failed to list Docker containers")
	}

	containerNames := make([]string, 0)
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, pattern) {
				containerNames = append(containerNames, strings.TrimPrefix(name, "/"))
			}
		}
	}

	return containerNames, nil
}

func CopyFilesToContainers(ctx context.Context, containerNames []string, targetDir string, files []string) error {
	frameworkDockerClient, frameworkDockerClientErr := framework.NewDockerClient()
	if frameworkDockerClientErr != nil {
		return errors.Wrap(frameworkDockerClientErr, "failed to create framework Docker client")
	}

	dockerClient, dockerClientErr := dc.NewClientWithOpts(dc.FromEnv, dc.WithAPIVersionNegotiation())
	if dockerClientErr != nil {
		return errors.Wrap(dockerClientErr, "failed to create Docker client")
	}
	defer dockerClient.Close()

	for _, containerName := range containerNames {
		execOutput, execOutputErr := frameworkDockerClient.ExecContainer(containerName, []string{"mkdir", "-p", targetDir})
		if execOutputErr != nil {
			fmt.Fprint(os.Stderr, execOutput)
			return errors.Wrap(execOutputErr, "failed to execute mkdir command in Docker container")
		}

		for _, filePath := range files {
			framework.L.Info().Msgf("Copying file '%s' to Docker container %s", filePath, containerName)
			copyErr := frameworkDockerClient.CopyFile(containerName, filePath, targetDir)
			if copyErr != nil {
				fmt.Fprint(os.Stderr, execOutput)
				return errors.Wrap(copyErr, "failed to copy artifact to Docker container")
			}
		}

		inspectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		containerJSON, inspectErr := dockerClient.ContainerInspect(inspectCtx, containerName)
		cancel()
		if inspectErr != nil {
			return errors.Wrap(inspectErr, "failed to inspect Docker container")
		}

		user := containerJSON.Config.User
		if user == "" {
			continue
		}
		for _, filePath := range files {
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
				return errors.Wrap(execOutputErr, "failed to execute chown command in Docker container")
			}
		}
	}

	return nil
}
