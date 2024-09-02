package persistent

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/integration-tests/client"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type EnvironmentConfig struct {
	persistent_types.ChainConfig
	DONConfig
}

func NewEnvironment(lggr logger.Logger, config EnvironmentConfig) (*deployment.Environment, error) {
	chains, err := NewChains(lggr, config.ChainConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create Chains")
	}

	// TODO add logstream
	// TODO add clean ups? although that should be related to test, not to environment
	if config.DONConfig.NewDON != nil {
		clNodesConfigs, err := NewEVMOnlyChainlinkConfigs(config.DONConfig.NewDON.ChainlinkDeployment, chains)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create chainlink configs")
		}
		config.DONConfig.NewDON.ChainlinkConfigs = clNodesConfigs
	}

	don, err := NewNodes(config.DONConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create nodes")
	}

	keys := make(map[uint64][]client.NodeKeysBundle)

	if config.DONConfig.NewDON != nil {
		var clients []*client.ChainlinkClient
		for _, k8sClient := range don.ClClients {
			clients = append(clients, k8sClient.ChainlinkClient)
		}
		for chainId := range chains {
			_, clNodes, err := client.CreateNodeKeysBundle(clients, "evm", fmt.Sprint(chainId))
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create node keys for chain %d", chainId)
			}
			keys[chainId] = func() []client.NodeKeysBundle {
				var keys []client.NodeKeysBundle
				for _, clNode := range clNodes {
					keys = append(keys, clNode.KeysBundle)
				}
				return keys
			}()
		}
	}

	_ = don

	var nodeIDs []string
	for _, node := range keys {
		for _, keys := range node {
			nodeIDs = append(nodeIDs, keys.PeerID)
		}
		// peer ids are the same for all nodes, so we can iterate only once
		break
	}

	mocks, err := NewMocks(config.DONConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create mocks")
	}

	don.MockServer = mocks.MockServer
	don.KillGrave = mocks.KillGrave

	return &deployment.Environment{
		Name: "persistent",
		//Offchain: NewMemoryJobClient(nodes),
		NodeIDs: nodeIDs,
		Chains:  chains,
		Logger:  lggr,
	}, nil
}
