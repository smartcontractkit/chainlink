package persistent

import (
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	ctftestenv "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

// NewNodes either starts a new Docker-based DON or connects to an existing DON based on the provided config.
func NewNodes(donConfig persistent_types.DONConfig) (*persistent_types.DON, error) {
	if err := validateDONConfig(donConfig); err != nil {
		return nil, err
	}

	if donConfig.NewDON != nil {
		return StartNewDockerDON(donConfig.NewDON)
	}

	return ConnectToExistingDON(donConfig.ExistingDON)
}

func StartNewDockerDON(newDonConfig *persistent_types.NewDockerDONConfig) (*persistent_types.DON, error) {
	don := &persistent_types.DON{}

	if err := appendNewDockerNetworkIfNotSet(&newDonConfig.Options); err != nil {
		return nil, errors.Wrap(err, "failed to append new Docker network")
	}

	var err error
	var clCluster *test_env.ClCluster
	if len(newDonConfig.Nodes) > 0 {
		clCluster, err = buildIdenticalChainlinkCluster(newDonConfig)
		if err != nil {
			return don, errors.Wrap(err, "failed to build identical Chainlink nodes")
		}
	} else {
		clCluster, err = buildCustomChainlinkCluster(newDonConfig)
		if err != nil {
			return don, errors.Wrap(err, "failed to build unique Chainlink nodes")
		}
	}

	if newDonConfig.NewDONHooks != nil {
		err := newDonConfig.NewDONHooks.PreStartupHook(clCluster.Nodes)
		if err != nil {
			return don, errors.Wrap(err, "failed to execute pre startup hook")
		}
	}

	err = startChainlinkClusterAndClients(clCluster, don)
	if err != nil {
		return don, errors.Wrap(err, "failed to start chainlink cluster and clients")
	}

	if newDonConfig.NewDONHooks != nil {
		err := newDonConfig.NewDONHooks.PostStartupHook(clCluster.Nodes)
		if err != nil {
			return don, errors.Wrap(err, "failed to execute post setup hook")
		}
	}

	return don, nil
}

func ConnectToExistingDON(config *persistent_types.ExistingDONConfig) (*persistent_types.DON, error) {
	noOfNodes := pointer.GetInt(config.NoOfNodes)
	namespace := pointer.GetString(config.Name)

	if noOfNodes != len(config.NodeConfigs) {
		return nil, fmt.Errorf("number of nodes %d does not match number of node configs %d", noOfNodes, len(config.NodeConfigs))
	}

	don := &persistent_types.DON{}

	for i := 0; i < noOfNodes; i++ {
		cfg := config.NodeConfigs[i]
		if cfg == nil {
			return nil, fmt.Errorf("node %d config is nil", i+1)
		}
		clClient, err := client.NewChainlinkK8sClient(cfg, cfg.InternalIP, namespace)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create chainlink clientfor node %d config %v", i+1, cfg)
		}
		don.ChainlinkClients = append(don.ChainlinkClients, clClient)
	}

	return don, nil
}

func appendNewDockerNetworkIfNotSet(options *persistent_types.DockerOptions) error {
	if len(options.Networks) == 0 {
		dockerNetwork, err := docker.CreateNetwork(logging.GetLogger(nil, "CORE_DOCKER_ENV_LOG_LEVEL"))
		if err != nil {
			return errors.Wrap(err, "failed to create docker network")
		}
		options.Networks = []string{dockerNetwork.Name}
	}

	return nil
}

func buildIdenticalChainlinkCluster(newDonConfig *persistent_types.NewDockerDONConfig) (*test_env.ClCluster, error) {
	clCluster := &test_env.ClCluster{}
	noOfNodes := pointer.GetInt(newDonConfig.NoOfNodes)

	if len(newDonConfig.Nodes) != len(newDonConfig.ChainlinkConfigs) {
		return clCluster, fmt.Errorf("number of nodes %d does not match number of chainlink configs %d", noOfNodes, len(newDonConfig.ChainlinkConfigs))
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
			return clCluster, errors.Wrapf(err, "failed to build new chainlink node")
		}
		clCluster.Nodes = append(clCluster.Nodes, node)

	}

	return clCluster, nil
}

func buildCustomChainlinkCluster(newDonConfig *persistent_types.NewDockerDONConfig) (*test_env.ClCluster, error) {
	clCluster := &test_env.ClCluster{}
	noOfNodes := pointer.GetInt(newDonConfig.NoOfNodes)

	if noOfNodes != len(newDonConfig.ChainlinkConfigs) {
		return clCluster, fmt.Errorf("number of nodes %d does not match number of chainlink configs %d", noOfNodes, len(newDonConfig.ChainlinkConfigs))
	}
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
			return clCluster, errors.Wrapf(err, "failed to build new chainlink node")
		}
		clCluster.Nodes = append(clCluster.Nodes, node)
	}

	return clCluster, nil
}

func startChainlinkClusterAndClients(clCluster *test_env.ClCluster, don *persistent_types.DON) error {
	startErr := clCluster.Start()
	if startErr != nil {
		return errors.Wrap(startErr, "failed to start chainlink cluster")
	}

	for _, node := range clCluster.Nodes {
		don.ChainlinkClients = append(don.ChainlinkClients, &client.ChainlinkK8sClient{
			ChainlinkClient: node.API,
		})
	}

	don.ChainlinkContainers = clCluster.Nodes

	return nil
}

func validateDONConfig(donConfig persistent_types.DONConfig) error {
	if donConfig.NewDON == nil && donConfig.ExistingDON == nil {
		return fmt.Errorf("no DON config provided, you need to provide either an existing or new DON config")
	}

	if donConfig.NewDON != nil && donConfig.ExistingDON != nil {
		return fmt.Errorf("both new and existing DON config provided, you need to provide either an existing or new DON config")
	}

	return nil
}
