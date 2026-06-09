package dockerops

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	dc "github.com/moby/moby/client"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func FindDockerContainerNames(ctx context.Context, pattern string) ([]string, error) {
	dockerClient, dockerClientErr := dc.New(dc.FromEnv)
	if dockerClientErr != nil {
		return nil, errors.Wrap(dockerClientErr, "failed to create Docker client")
	}
	defer dockerClient.Close()

	listRes, containersErr := dockerClient.ContainerList(context.Background(), dc.ContainerListOptions{})
	if containersErr != nil {
		return nil, errors.Wrap(containersErr, "failed to list Docker containers")
	}

	containerNames := []string{}
	for _, ctr := range listRes.Items {
		for _, name := range ctr.Names {
			if strings.Contains(name, pattern) {
				// Remove leading slash from container name
				cleanName := strings.TrimPrefix(name, "/")
				containerNames = append(containerNames, cleanName)
			}
		}
	}

	return containerNames, nil
}

func CopyFilesToContainers(ctx context.Context, containerNamePattern string, targetDir string, files []string) error {
	containerNames, err := FindDockerContainerNames(context.Background(), containerNamePattern)
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
	dockerClient, err := dc.New(dc.FromEnv)
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

			for _, filePath := range files {
				framework.L.Info().Msgf("Copying file '%s' to Docker container %s", filePath, containerName)
				copyErr := frameworkDockerClient.CopyFile(containerName, filePath, targetDir)
				if copyErr != nil {
					fmt.Fprint(os.Stderr, execOutput)
					return errors.Wrap(copyErr, "failed to copy artifact to Docker container")
				}
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			inspectRes, inspectErr := dockerClient.ContainerInspect(ctx, containerName, dc.ContainerInspectOptions{})
			if inspectErr != nil {
				return errors.Wrap(inspectErr, "failed to inspect Docker container")
			}
			user := inspectRes.Container.Config.User

			// if not running as root, change ownership to user that is running the container to avoid permission issues
			if user == "" {
				return nil
			}

			for _, filePath := range files {
				targetFilePath := filepath.Join(targetDir, filepath.Base(filePath))
				execConfig := dc.ExecCreateOptions{
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
