package test_env

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"

	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
)

var (
	ErrGetNodeCSAKeys = "failed get CL node CSA keys"
)

type ClCluster struct {
	Nodes []*ClNode `json:"nodes"`
}

// Start all nodes in the cluster./docker/tests/functional/api
func (c *ClCluster) Start() error {
	eg := &errgroup.Group{}
	nodes := c.Nodes

	for i := 0; i < len(nodes); i++ {
		nodeIndex := i
		eg.Go(func() error {
			err := nodes[nodeIndex].StartContainer()
			if err != nil {
				return err
			}
			return nil
		})
	}

	return eg.Wait()
}

func (c *ClCluster) NodeAPIs() []*client.ChainlinkClient {
	clients := make([]*client.ChainlinkClient, 0)
	for _, c := range c.Nodes {
		clients = append(clients, c.API)
	}
	return clients
}

// Return all the on-chain wallet addresses for a set of Chainlink nodes
func (c *ClCluster) NodeAddresses() ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	for _, n := range c.Nodes {
		primaryAddress, err := n.ChainlinkNodeAddress()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, primaryAddress)
	}
	return addresses, nil
}

func (c *ClCluster) NodeCSAKeys() ([]string, error) {
	var keys []string
	for _, n := range c.Nodes {
		csaKeys, err := n.GetNodeCSAKeys()
		if err != nil {
			return nil, fmt.Errorf("%s, err: %w", ErrGetNodeCSAKeys, err)
		}
		keys = append(keys, csaKeys.Data[0].ID)
	}
	return keys, nil
}

func (c *ClCluster) CopyFolderFromNodes(ctx context.Context, srcPath, destPath string) error {
	var wg sync.WaitGroup
	errors := make(chan error, len(c.Nodes))

	for i, node := range c.Nodes {
		wg.Add(1)
		go func(n *ClNode, id int) {
			defer wg.Done()
			// Create a unique subdirectory for each node based on an identifier
			finalDestPath := filepath.Join(destPath, fmt.Sprintf("node_%d", id)) // Use node ID or another unique identifier
			if err := os.MkdirAll(finalDestPath, 0755); err != nil {
				errors <- fmt.Errorf("failed to create directory for node %d: %w", id, err)
				return
			}
			err := copyFolderFromContainerUsingDockerCP(ctx, n.Container.GetContainerID(), srcPath, finalDestPath)
			if err != nil {
				errors <- fmt.Errorf("failed to copy folder for node %d: %w", id, err)
			}
		}(node, i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFolderFromContainer(ctx context.Context, container tc.Container, srcPath, destPath string) error {
	// List all files and directories recursively inside the container
	lsCmd := []string{"find", srcPath, "-type", "f"} // Lists only files, omitting directories
	outputCode, outputReader, err := container.Exec(ctx, lsCmd)
	if err != nil {
		return fmt.Errorf("failed to list files in container: %w", err)
	}
	if outputCode != 0 {
		return fmt.Errorf("could not list files in the container. Command exited with code: %d", outputCode)
	}

	// Read the output into a slice of file paths
	output, err := io.ReadAll(outputReader)
	if err != nil {
		return fmt.Errorf("failed to read command output: %w", err)
	}
	outStr := string(output)
	files := strings.Split(outStr, "\n")

	// Ensure destination path exists or create it
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Iterate over each file path
	for _, file := range files {
		if file == "" {
			continue
		}

		// Define the path for the file on the host
		relPath, err := filepath.Rel(srcPath, file)
		if err != nil {
			return fmt.Errorf("failed to compute relative path: %w", err)
		}
		hostPath := filepath.Join(destPath, relPath)

		// Ensure the subdirectory exists
		if err := os.MkdirAll(filepath.Dir(hostPath), 0755); err != nil {
			return fmt.Errorf("failed to create subdirectory: %w", err)
		}

		// Copy the file from the container
		reader, err := container.CopyFileFromContainer(ctx, file)
		if err != nil {
			return fmt.Errorf("failed to copy file %s from container: %w", file, err)
		}
		defer reader.Close()

		// Create the file on the host
		localFile, err := os.Create(hostPath)
		if err != nil {
			return fmt.Errorf("failed to create file on host: %w", err)
		}
		defer localFile.Close()

		// Copy data from reader to local file
		if _, err := io.Copy(localFile, reader); err != nil {
			return fmt.Errorf("failed to copy file content: %w", err)
		}
	}

	return nil
}

func copyFolderFromContainerUsingDockerCP(ctx context.Context, containerID, srcPath, destPath string) error {
	source := fmt.Sprintf("%s:%s", containerID, srcPath)

	// Prepare the docker cp command
	cmd := exec.CommandContext(ctx, "docker", "cp", source, destPath)

	// Execute the docker cp command
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("docker cp command failed: %s, output: %s", err, string(output))
	}

	return nil
}
