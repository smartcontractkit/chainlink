package devenv

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	ctftestenv "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env/job_distributor"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/stretchr/testify/require"
	"github.com/subosito/gotenv"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	clclient "github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

// CreateDockerEnv creates a new test environment with simulated private ethereum networks and job distributor
// It returns the EnvironmentConfig which holds the chain config and JD config
// The test environment is then used to start chainlink nodes
func CreateDockerEnv(t *testing.T) (
	*EnvironmentConfig,
	*test_env.CLClusterTestEnv,
	tc.TestConfig,
) {
	if _, err := os.Stat(".env"); err == nil || !os.IsNotExist(err) {
		require.NoError(t, gotenv.Load(".env"), "Error loading .env file")
	}

	cfg, err := tc.GetChainAndTestTypeSpecificConfig("Smoke", tc.CCIP)
	require.NoError(t, err, "Error getting config")

	var privateEthereumNetworks []*ctf_config.EthereumNetworkConfig
	for _, network := range cfg.CCIP.PrivateEthereumNetworks {
		privateEthereumNetworks = append(privateEthereumNetworks, network)
	}
	env, err := test_env.NewCLTestEnvBuilder().
		WithTestConfig(&cfg).
		WithTestInstance(t).
		WithPrivateEthereumNetworks(privateEthereumNetworks).
		WithStandardCleanup().
		Build()
	require.NoError(t, err, "Error building test environment")
	chains := CreateChainConfigFromPrivateEthereumNetworks(t, env, cfg.CCIP.PrivateEthereumNetworks, cfg.GetNetworkConfig())

	var jdConfig JDConfig
	// TODO : move this as a part of test_env setup with an input in testconfig
	// if JD is not provided, we will spin up a new JD
	if cfg.CCIP.GetJDGRPC() == "" && cfg.CCIP.GetJDWSRPC() == "" {
		jdDB, err := ctftestenv.NewPostgresDb(
			[]string{env.DockerNetwork.Name},
			ctftestenv.WithPostgresDbName(cfg.CCIP.GetJDDBName()),
			ctftestenv.WithPostgresImageVersion(cfg.CCIP.GetJDDBVersion()),
		)
		require.NoError(t, err)
		err = jdDB.StartContainer()
		require.NoError(t, err)

		jd := job_distributor.New([]string{env.DockerNetwork.Name},
			job_distributor.WithImage(cfg.CCIP.GetJDImage()),
			job_distributor.WithVersion(cfg.CCIP.GetJDVersion()),
			job_distributor.WithDBURL(jdDB.InternalURL.String()),
		)
		err = jd.StartContainer()
		require.NoError(t, err)
		jdConfig = JDConfig{
			GRPC: jd.Grpc,
			// we will use internal wsrpc for nodes on same docker network to connect to JD
			WSRPC: jd.InternalWSRPC,
		}
	} else {
		jdConfig = JDConfig{
			GRPC:  cfg.CCIP.GetJDGRPC(),
			WSRPC: cfg.CCIP.GetJDWSRPC(),
		}
	}
	require.NotEmpty(t, jdConfig, "JD config is empty")

	return &EnvironmentConfig{
		Chains:   chains,
		JDConfig: jdConfig,
	}, env, cfg
}

// StartChainlinkNodes starts docker containers for chainlink nodes on the existing test environment based on provided test config
// Once the nodes starts, it updates the devenv EnvironmentConfig with the node info
// which includes chainlink API URL, email, password and internal IP
func StartChainlinkNodes(
	t *testing.T,
	envConfig *EnvironmentConfig,
	registryConfig deployment.CapabilityRegistryConfig,
	env *test_env.CLClusterTestEnv,
	cfg tc.TestConfig,
) error {
	evmNetworks := networks.MustGetSelectedNetworkConfig(cfg.GetNetworkConfig())
	for i, net := range evmNetworks {
		rpcProvider, err := env.GetRpcProvider(net.ChainID)
		require.NoError(t, err, "Error getting rpc provider")
		evmNetworks[i].HTTPURLs = rpcProvider.PrivateHttpUrls()
		evmNetworks[i].URLs = rpcProvider.PrivateWsUrsl()
	}
	noOfNodes := pointer.GetInt(cfg.CCIP.CLNode.NoOfPluginNodes) + pointer.GetInt(cfg.CCIP.CLNode.NoOfBootstraps)
	var nodeInfo []NodeInfo
	for i := 1; i <= noOfNodes; i++ {
		if i <= pointer.GetInt(cfg.CCIP.CLNode.NoOfBootstraps) {
			nodeInfo = append(nodeInfo, NodeInfo{
				IsBootstrap: true,
				Name:        fmt.Sprintf("bootstrap-%d", i),
				// TODO : make this configurable
				P2PPort: "6690",
			})
		} else {
			nodeInfo = append(nodeInfo, NodeInfo{
				IsBootstrap: false,
				Name:        fmt.Sprintf("node-%d", i-1),
				// TODO : make this configurable
				P2PPort: "6690",
			})
		}
		toml, _, err := testsetups.SetNodeConfig(
			evmNetworks,
			cfg.NodeConfig.BaseConfigTOML,
			cfg.NodeConfig.CommonChainConfigTOML,
			cfg.NodeConfig.ChainConfigTOMLByChainID,
		)

		toml.Capabilities.ExternalRegistry.NetworkID = ptr.Ptr(relay.NetworkEVM)
		toml.Capabilities.ExternalRegistry.ChainID = ptr.Ptr(strconv.FormatUint(registryConfig.EVMChainID, 10))
		toml.Capabilities.ExternalRegistry.Address = ptr.Ptr(registryConfig.Contract.String())

		if err != nil {
			return err
		}
		ccipNode, err := test_env.NewClNode(
			[]string{env.DockerNetwork.Name},
			pointer.GetString(cfg.GetChainlinkImageConfig().Image),
			pointer.GetString(cfg.GetChainlinkImageConfig().Version),
			toml,
			env.LogStream,
			test_env.WithPgDBOptions(
				ctftestenv.WithPostgresImageVersion(pointer.GetString(cfg.GetChainlinkImageConfig().PostgresVersion)),
			),
		)
		if err != nil {
			return err
		}
		ccipNode.SetTestLogger(t)
		env.ClCluster.Nodes = append(env.ClCluster.Nodes, ccipNode)
	}
	err := env.ClCluster.Start()
	if err != nil {
		return err
	}
	for i, n := range env.ClCluster.Nodes {
		nodeInfo[i].CLConfig = clclient.ChainlinkConfig{
			URL:        n.API.URL(),
			Email:      n.UserEmail,
			Password:   n.UserPassword,
			InternalIP: n.API.InternalIP(),
		}
	}

	envConfig.nodeInfo = nodeInfo
	return nil
}

// CreateChainConfigFromPrivateEthereumNetworks creates a list of ChainConfig from the private ethereum networks created by the test environment.
// It uses the private keys from the network config to create the deployer key for each chain.
func CreateChainConfigFromPrivateEthereumNetworks(
	t *testing.T,
	env *test_env.CLClusterTestEnv,
	privateEthereumNetworks map[string]*ctf_config.EthereumNetworkConfig,
	networkConfig *ctf_config.NetworkConfig,
) []ChainConfig {
	evmNetworks := networks.MustGetSelectedNetworkConfig(networkConfig)
	networkPvtKeys := make(map[int64]string)
	for _, net := range evmNetworks {
		require.Greater(t, len(net.PrivateKeys), 0, "No private keys found for network")
		networkPvtKeys[net.ChainID] = net.PrivateKeys[0]
	}
	var chains []ChainConfig
	for _, networkCfg := range privateEthereumNetworks {
		chainId := networkCfg.EthereumChainConfig.ChainID
		chainName, err := chainselectors.NameFromChainId(uint64(chainId))
		require.NoError(t, err, "Error getting chain name")
		rpcProvider, err := env.GetRpcProvider(int64(chainId))
		require.NoError(t, err, "Error getting rpc provider")
		pvtKeyStr, exists := networkPvtKeys[int64(chainId)]
		require.Truef(t, exists, "Private key not found for chain id %d", chainId)
		pvtKey, err := crypto.HexToECDSA(pvtKeyStr)
		require.NoError(t, err)
		deployer, err := bind.NewKeyedTransactorWithChainID(pvtKey, big.NewInt(int64(chainId)))
		require.NoError(t, err)
		chains = append(chains, ChainConfig{
			ChainID:         uint64(chainId),
			ChainName:       chainName,
			ChainType:       "EVM",
			WSRPCs:          rpcProvider.PublicWsUrls(),
			HTTPRPCs:        rpcProvider.PublicHttpUrls(),
			PrivateHTTPRPCs: rpcProvider.PrivateHttpUrls(),
			PrivateWSRPCs:   rpcProvider.PrivateWsUrsl(),
			DeployerKey:     deployer,
		})
	}
	return chains
}
