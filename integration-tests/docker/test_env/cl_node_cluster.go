package test_env

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"golang.org/x/sync/errgroup"
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
			return nil, errors.Wrap(err, ErrGetNodeCSAKeys)
		}
		keys = append(keys, csaKeys.Data[0].ID)
	}
	return keys, nil
}
