package persistent

import (
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	ctftestenv "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type ExistingDONConfig struct {
	testconfig.CLCluster
	MockServerURL *string `toml:",omitempty"`
}

type NewDONHooks interface {
	// PreStartupHook is called before the DON is started. No containers are running yet. For example, you can use this hook to modify configuration of each node.
	PreStartupHook([]*test_env.ClNode) error
	// PostStartupHook is called after the DON is started. All containers are running. For example, you can use this hook to interact with them using the API.
	PostStartupHook([]*test_env.ClNode) error
}

type NewDockerDONConfig struct {
	testconfig.ChainlinkDeployment
	Options          Options
	ChainlinkConfigs []*chainlink.Config
	NewDONHooks
}

type Options struct {
	Networks  []string
	LogStream *logstream.LogStream
}

type DONConfig struct {
	ExistingDON *ExistingDONConfig
	NewDON      *NewDockerDONConfig
}

type DON struct {
	ClClients []*client.ChainlinkK8sClient
	deployment.Mocks
}

func NewNodes(donConfig DONConfig) (DON, error) {
	if donConfig.NewDON == nil && donConfig.ExistingDON == nil {
		return DON{}, fmt.Errorf("no DON config provided, you need to provide either an existing or new DON config")
	}

	if donConfig.NewDON != nil && donConfig.ExistingDON != nil {
		return DON{}, fmt.Errorf("both new and existing DON config provided, you need to provide either an existing or new DON config")
	}

	if donConfig.NewDON != nil {
		return NewDockerDON(donConfig.NewDON)
	}

	return ConnectToExistingNodes(*donConfig.ExistingDON)
}

func ConnectToExistingNodes(config ExistingDONConfig) (DON, error) {
	noOfNodes := pointer.GetInt(config.NoOfNodes)
	namespace := pointer.GetString(config.Name)

	if noOfNodes != len(config.NodeConfigs) {
		return DON{}, fmt.Errorf("number of nodes %d does not match number of node configs %d", noOfNodes, len(config.NodeConfigs))
	}

	don := DON{}

	for i := 0; i < noOfNodes; i++ {
		cfg := config.NodeConfigs[i]
		if cfg == nil {
			return don, fmt.Errorf("node %d config is nil", i+1)
		}
		clClient, err := client.NewChainlinkK8sClient(cfg, cfg.InternalIP, namespace)
		if err != nil {
			return don, errors.Wrapf(err, "failed to create chainlink client: %w for node %d config %v", i+1, cfg)
		}
		don.ClClients = append(don.ClClients, clClient)
	}

	return don, nil
}

func NewDockerDON(newDonConfig *NewDockerDONConfig) (DON, error) {
	don := DON{}

	// maybe we should validate this and return err if not set instead of generating here
	if len(newDonConfig.Options.Networks) == 0 {
		dockerNetwork, err := docker.CreateNetwork(logging.GetLogger(nil, "CORE_DOCKER_ENV_LOG_LEVEL"))
		if err != nil {
			return don, errors.Wrap(err, "failed to create docker network")
		}
		newDonConfig.Options.Networks = []string{dockerNetwork.Name}
	}

	clCluster := test_env.ClCluster{}
	noOfNodes := pointer.GetInt(newDonConfig.NoOfNodes)
	// if individual nodes are specified, then deploy them with specified configs
	// TODO probably best to put it in a reusable method, use it here and also in integration-tests/ccip-tests/testsetups/test_env.go
	if len(newDonConfig.Nodes) > 0 {
		if len(newDonConfig.Nodes) != len(newDonConfig.ChainlinkConfigs) {
			return don, fmt.Errorf("number of nodes %d does not match number of chainlink configs %d", noOfNodes, len(newDonConfig.ChainlinkConfigs))
		}
		for i, clNode := range newDonConfig.Nodes {
			node, err := test_env.NewClNode(
				newDonConfig.Options.Networks,
				pointer.GetString(clNode.ChainlinkImage.Image),
				pointer.GetString(clNode.ChainlinkImage.Version),
				newDonConfig.ChainlinkConfigs[i],
				newDonConfig.Options.LogStream,
				test_env.WithPgDBOptions(
					ctftestenv.WithPostgresImageName(clNode.DBImage),
					ctftestenv.WithPostgresImageVersion(clNode.DBTag),
				),
			)
			if err != nil {
				return don, errors.Wrapf(err, "failed to build new chainlink node")
			}
			// node.SetTestLogger(t)
			clCluster.Nodes = append(clCluster.Nodes, node)
		}
	} else {
		if noOfNodes != len(newDonConfig.ChainlinkConfigs) {
			return don, fmt.Errorf("number of nodes %d does not match number of chainlink configs %d", noOfNodes, len(newDonConfig.ChainlinkConfigs))
		}
		// if no individual nodes are specified, then deploy the number of nodes specified in the env input with common config
		for i := 0; i < noOfNodes; i++ {
			node, err := test_env.NewClNode(
				newDonConfig.Options.Networks,
				pointer.GetString(newDonConfig.Common.ChainlinkImage.Image),
				pointer.GetString(newDonConfig.Common.ChainlinkImage.Version),
				newDonConfig.ChainlinkConfigs[i],
				newDonConfig.Options.LogStream,
				test_env.WithPgDBOptions(
					ctftestenv.WithPostgresImageName(newDonConfig.Common.DBImage),
					ctftestenv.WithPostgresImageVersion(newDonConfig.Common.DBTag),
				),
			)
			if err != nil {
				return don, errors.Wrapf(err, "failed to build new chainlink node")
			}
			//node.SetTestLogger(t)
			clCluster.Nodes = append(clCluster.Nodes, node)
		}
	}

	if newDonConfig.NewDONHooks != nil {
		err := newDonConfig.NewDONHooks.PreStartupHook(clCluster.Nodes)
		if err != nil {
			return don, errors.Wrap(err, "failed to execute post setup hook")
		}
	}

	startErr := clCluster.Start()
	if startErr != nil {
		return don, errors.Wrap(startErr, "failed to start chainlink cluster")
	}

	for _, node := range clCluster.Nodes {
		don.ClClients = append(don.ClClients, &client.ChainlinkK8sClient{
			ChainlinkClient: node.API,
		})
	}

	if newDonConfig.NewDONHooks != nil {
		err := newDonConfig.NewDONHooks.PostStartupHook(clCluster.Nodes)
		if err != nil {
			return don, errors.Wrap(err, "failed to execute post setup hook")
		}
	}

	return don, nil
}

func NewEVMOnlyChainlinkConfigs(donConfig testconfig.ChainlinkDeployment, chains map[uint64]deployment.Chain, rpcProviders map[uint64]persistent_types.RpcProvider) ([]*chainlink.Config, error) {
	var evmNetworks []blockchain.EVMNetwork
	for _, rpcProvider := range rpcProviders {
		evmNetwork := rpcProvider.EVMNetwork()
		evmNetwork.HTTPURLs = rpcProvider.PrivateHttpUrls()
		evmNetwork.URLs = rpcProvider.PrivateWsUrls()
		evmNetworks = append(evmNetworks, evmNetwork)
	}

	var clNodeConfigs []*chainlink.Config

	noOfNodes := pointer.GetInt(donConfig.NoOfNodes)
	// if individual nodes are specified, then deploy them with specified configs
	// TODO probably best to put it in a reusable method, use it here and also in integration-tests/ccip-tests/testsetups/test_env.go
	if len(donConfig.Nodes) > 0 {
		for range donConfig.Nodes {
			toml, _, err := testsetups.SetNodeConfig(
				evmNetworks,
				donConfig.Common.BaseConfigTOML,
				donConfig.Common.CommonChainConfigTOML,
				donConfig.Common.ChainConfigTOMLByChain,
			)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create node config")
			}
			clNodeConfigs = append(clNodeConfigs, toml)
		}
	} else {
		// if no individual nodes are specified, then deploy the number of nodes specified in the env input with common config
		for i := 0; i < noOfNodes; i++ {
			toml, _, err := testsetups.SetNodeConfig(
				evmNetworks,
				donConfig.Common.BaseConfigTOML,
				donConfig.Common.CommonChainConfigTOML,
				donConfig.Common.ChainConfigTOMLByChain,
			)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create node config")
			}
			clNodeConfigs = append(clNodeConfigs, toml)
		}
	}

	return clNodeConfigs, nil
}
