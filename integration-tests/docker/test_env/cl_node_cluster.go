package test_env

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
			err := copyFolderFromContainer(ctx, n.Container, srcPath, finalDestPath)
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
	// Tar the source directory inside the container
	tarCmd := []string{"tar", "-czf", "/tmp/archive.tar.gz", "-C", srcPath, "."}
	_, _, err := container.Exec(ctx, tarCmd)
	if err != nil {
		return fmt.Errorf("failed to tar folder in container: %w", err)
	}

	reader, err := container.CopyFileFromContainer(ctx, "/tmp/archive.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to copy from container: %w", err)
	}
	defer reader.Close()

	// Ensure destination path exists
	if info, err := os.Stat(destPath); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("destination path %s is not a directory", destPath)
		}
	} else if os.IsNotExist(err) {
		return fmt.Errorf("destination path %s does not exist", destPath)
	} else {
		return fmt.Errorf("error checking destination directory: %w", err)
	}

	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create the tar file on the host
	destTarPath := filepath.Join(destPath, "archive.tar.gz")
	localTarFile, err := os.Create(destTarPath)
	if err != nil {
		return fmt.Errorf("failed to create tar file on host: %w", err)
	}
	defer localTarFile.Close()

	// Copy tar data from the container to the host file
	if _, err := io.Copy(localTarFile, reader); err != nil {
		return fmt.Errorf("failed to copy tar file content: %w", err)
	}

	// Extract the tar file
	cmd := exec.Command("tar", "-xzf", destTarPath, "-C", destPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract tar file: %w", err)
	}

	return nil
}
