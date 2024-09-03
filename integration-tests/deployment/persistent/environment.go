package persistent

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/solclient"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type EnvironmentConfig struct {
	persistent_types.ChainConfig
	DONConfig
}

type ChainProducer[ChainType any] func(persistent_types.ChainConfig) (map[uint64]ChainType, error)
type NodeConfigProducer[ChainType any] func(map[uint64]ChainType) ([]*chainlink.Config, error)

func NewEnvironment[ChainType any](lggr logger.Logger, config EnvironmentConfig, chainsFn ChainProducer[ChainType], configFn NodeConfigProducer[ChainType]) (*deployment.Environment[ChainType], error) {
	chains, err := chainsFn(config.ChainConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create chains")
	}

	if config.DONConfig.NewDON != nil {
		clNodesConfigs, err := configFn(chains)
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

	return &deployment.Environment[ChainType]{
		Name: "persistent",
		//Offchain: NewMemoryJobClient(nodes),
		NodeIDs: nodeIDs,
		Chains:  chains,
		Logger:  lggr,
	}, nil
}

func NewGenericEVMEnvironment(lggr logger.Logger, config EnvironmentConfig) (*deployment.Environment[deployment.Chain], error) {
	var chainFn ChainProducer[deployment.Chain] = func(persistent_types.ChainConfig) (map[uint64]deployment.Chain, error) {
		return NewChains(lggr, config.ChainConfig)
	}

	configFn := func(chains map[uint64]deployment.Chain) ([]*chainlink.Config, error) {
		return NewEVMOnlyChainlinkConfigs(config.DONConfig.NewDON.ChainlinkDeployment, chains)
	}

	return NewEnvironment[deployment.Chain](lggr, config, chainFn, configFn)
}

func NewGenericSolanaEnvironment(lggr logger.Logger, config EnvironmentConfig) (*deployment.Environment[deployment.SolanaChain], error) {
	solNetworks := make(map[uint64]*solclient.SolNetwork)
	chainFn := func(persistent_types.ChainConfig) (map[uint64]deployment.SolanaChain, error) {
		chains, networks, err := NewSolanaChains(lggr, config.ChainConfig)
		if err != nil {
			return nil, err
		}

		solNetworks = networks
		return chains, nil
	}

	configFn := func(chains map[uint64]deployment.SolanaChain) ([]*chainlink.Config, error) {
		return NewSolanaChainlinkConfigs(config.DONConfig.NewDON.ChainlinkDeployment, solNetworks, chains)
	}

	return NewEnvironment[deployment.SolanaChain](lggr, config, chainFn, configFn)
}

func NewEVMEnvironment(lggr logger.Logger, config EnvironmentConfig) (*deployment.Environment[deployment.Chain], error) {
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

	return &deployment.Environment[deployment.Chain]{
		Name: "persistent",
		//Offchain: NewMemoryJobClient(nodes),
		NodeIDs: nodeIDs,
		Chains:  chains,
		Logger:  lggr,
	}, nil
}
