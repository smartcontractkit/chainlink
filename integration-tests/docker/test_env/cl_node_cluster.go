package test_env

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
)

var ErrGetNodeCSAKeys = "failed get CL node CSA keys"

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

func (c *ClCluster) Stop() error {
	var eg errgroup.Group
	nodes := c.Nodes
	timeout := time.Minute * 1

	for i := 0; i < len(nodes); i++ {
		nodeIndex := i
		eg.Go(func() error {
			if container := nodes[nodeIndex].Container; container != nil {
				return container.Stop(context.Background(), &timeout)
			}
			return nil
		})
	}

	return eg.Wait()
}

func (c *ClCluster) NodeAPIs() []*nodeclient.ChainlinkClient {
	clients := make([]*nodeclient.ChainlinkClient, 0)
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
			finalDestPath := filepath.Join(destPath, fmt.Sprintf("node_%d", id))
			if err := os.MkdirAll(finalDestPath, 0755); err != nil {
				errors <- fmt.Errorf("failed to create directory for node %d: %w", id, err)
				return
			}
			err := copyFolderFromContainerUsingDockerCP(ctx, n.Container.GetContainerID(), srcPath, finalDestPath)
			if err != nil {
				errors <- fmt.Errorf("failed to copy folder for node %d: %w", id, err)
				return
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

func copyFolderFromContainerUsingDockerCP(ctx context.Context, containerID, srcPath, destPath string) error {
	source := fmt.Sprintf("%s:%s", containerID, srcPath)
	cmd := exec.CommandContext(ctx, "docker", "cp", source, destPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrapf(err, "docker cp command failed: %s, output: %s", cmd, string(output))
	}
	return nil
}
