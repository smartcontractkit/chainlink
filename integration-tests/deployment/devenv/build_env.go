package devenv

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"github.com/subosito/gotenv"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	ctftestenv "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	clclient "github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"
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

	evmNetworks := networks.MustGetSelectedNetworkConfig(cfg.GetNetworkConfig())

	// find out if the selected networks are provided with PrivateEthereumNetworks configs
	// if yes, PrivateEthereumNetworkConfig will be used to create simulated private ethereum networks in docker environment
	var privateEthereumNetworks []*ctf_config.EthereumNetworkConfig
	for _, name := range cfg.GetNetworkConfig().SelectedNetworks {
		if network, exists := cfg.CCIP.PrivateEthereumNetworks[name]; exists {
			privateEthereumNetworks = append(privateEthereumNetworks, network)
		}
	}

	builder := test_env.NewCLTestEnvBuilder().
		WithTestConfig(&cfg).
		WithTestInstance(t).
		WithJobDistributor(cfg.CCIP.JobDistributorConfig).
		WithStandardCleanup()

	// if private ethereum networks are provided, we will use them to create the test environment
	// otherwise we will use the network URLs provided in the network config
	if len(privateEthereumNetworks) > 0 {
		builder = builder.WithPrivateEthereumNetworks(privateEthereumNetworks)
	}
	env, err := builder.Build()
	require.NoError(t, err, "Error building test environment")

	// we need to update the URLs for the simulated networks to the private chain RPCs in the docker test environment
	// so that the chainlink nodes and rmn nodes can internally connect to the chain
	env.EVMNetworks = []*blockchain.EVMNetwork{}
	for i, net := range evmNetworks {
		// if network is simulated, update the URLs with private chain RPCs in the docker test environment
		// so that nodes can internally connect to the chain
		if net.Simulated {
			rpcProvider, err := env.GetRpcProvider(net.ChainID)
			require.NoError(t, err, "Error getting rpc provider")
			evmNetworks[i].HTTPURLs = rpcProvider.PrivateHttpUrls()
			evmNetworks[i].URLs = rpcProvider.PrivateWsUrsl()
		}
		env.EVMNetworks = append(env.EVMNetworks, &evmNetworks[i])
	}

	chains := CreateChainConfigFromNetworks(t, env, privateEthereumNetworks, cfg.GetNetworkConfig())

	jdConfig := JDConfig{
		GRPC:  cfg.CCIP.JobDistributorConfig.GetJDGRPC(),
		WSRPC: cfg.CCIP.JobDistributorConfig.GetJDWSRPC(),
	}
	// TODO : move this as a part of test_env setup with an input in testconfig
	// if JD is not provided, we will spin up a new JD
	if jdConfig.GRPC == "" || jdConfig.WSRPC == "" {
		jd := env.JobDistributor
		require.NotNil(t, jd, "JD is not found in test environment")
		jdConfig = JDConfig{
			GRPC: jd.Grpc,
			// we will use internal wsrpc for nodes on same docker network to connect to JD
			WSRPC: jd.InternalWSRPC,
			Creds: insecure.NewCredentials(),
		}
	}
	require.NotEmpty(t, jdConfig, "JD config is empty")

	homeChainSelector, err := cfg.CCIP.GetHomeChainSelector(evmNetworks)
	require.NoError(t, err, "Error getting home chain selector")
	feedChainSelector, err := cfg.CCIP.GetFeedChainSelector(evmNetworks)
	require.NoError(t, err, "Error getting feed chain selector")

	return &EnvironmentConfig{
		Chains:            chains,
		JDConfig:          jdConfig,
		HomeChainSelector: homeChainSelector,
		FeedChainSelector: feedChainSelector,
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
	var evmNetworks []blockchain.EVMNetwork
	for i := range env.EVMNetworks {
		evmNetworks = append(evmNetworks, *env.EVMNetworks[i])
	}
	noOfNodes := pointer.GetInt(cfg.CCIP.CLNode.NoOfPluginNodes) + pointer.GetInt(cfg.CCIP.CLNode.NoOfBootstraps)
	if env.ClCluster == nil {
		env.ClCluster = &test_env.ClCluster{}
	}
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
	if envConfig == nil {
		envConfig = &EnvironmentConfig{}
	}
	envConfig.JDConfig.nodeInfo = nodeInfo
	return nil
}

// FundNodes sends funds to the chainlink nodes based on the provided test config
// It also sets up a clean-up function to return the funds back to the deployer account once the test is done
// It assumes that the chainlink nodes are already started and the account addresses for all chains are available
func FundNodes(t *testing.T, lggr zerolog.Logger, env *test_env.CLClusterTestEnv, cfg tc.TestConfig, nodes []Node) {
	evmNetworks := networks.MustGetSelectedNetworkConfig(cfg.GetNetworkConfig())
	for i, net := range evmNetworks {
		// if network is simulated, update the URLs with deployed chain RPCs in the docker test environment
		if net.Simulated {
			rpcProvider, err := env.GetRpcProvider(net.ChainID)
			require.NoError(t, err, "Error getting rpc provider")
			evmNetworks[i].HTTPURLs = rpcProvider.PublicHttpUrls()
			evmNetworks[i].URLs = rpcProvider.PublicWsUrls()
		}
	}
	t.Cleanup(func() {
		for i := range evmNetworks {
			// if simulated no need for balance return
			if evmNetworks[i].Simulated {
				continue
			}
			evmNetwork := evmNetworks[i]
			sethClient, err := utils.TestAwareSethClient(t, cfg, &evmNetwork)
			require.NoError(t, err, "Error getting seth client for network %s", evmNetwork.Name)
			require.Greater(t, len(sethClient.PrivateKeys), 0, seth.ErrNoKeyLoaded)
			var keyExporters []contracts.ChainlinkKeyExporter
			for j := range nodes {
				node := nodes[j]
				keyExporters = append(keyExporters, &node)
			}
			if err := actions.ReturnFundsFromKeyExporterNodes(lggr, sethClient, keyExporters); err != nil {
				lggr.Error().Err(err).Str("Network", evmNetwork.Name).
					Msg("Error attempting to return funds from chainlink nodes to network's default wallet. " +
						"Environment is left running so you can try manually!")
			}
		}
	})
	for i := range evmNetworks {
		evmNetwork := evmNetworks[i]
		sethClient, err := utils.TestAwareSethClient(t, cfg, &evmNetwork)
		require.NoError(t, err, "Error getting seth client for network %s", evmNetwork.Name)
		require.Greater(t, len(sethClient.PrivateKeys), 0, seth.ErrNoKeyLoaded)
		privateKey := sethClient.PrivateKeys[0]
		for _, node := range nodes {
			nodeAddr, ok := node.AccountAddr[uint64(evmNetwork.ChainID)]
			require.True(t, ok, "Account address not found for chain %d", evmNetwork.ChainID)
			fromAddress, err := actions.PrivateKeyToAddress(privateKey)
			require.NoError(t, err, "Error getting address from private key")
			amount := big.NewFloat(pointer.GetFloat64(cfg.Common.ChainlinkNodeFunding))
			toAddr := common.HexToAddress(nodeAddr)
			receipt, err := actions.SendFunds(lggr, sethClient, actions.FundsToSendPayload{
				ToAddress:  toAddr,
				Amount:     conversions.EtherToWei(amount),
				PrivateKey: privateKey,
			})
			require.NoError(t, err, "Error sending funds to node %s", node.Name)
			require.NotNil(t, receipt, "Receipt is nil")
			txHash := "(none)"
			if receipt != nil {
				txHash = receipt.TxHash.String()
			}
			lggr.Info().
				Str("From", fromAddress.Hex()).
				Str("To", toAddr.String()).
				Str("TxHash", txHash).
				Str("Amount", amount.String()).
				Msg("Funded Chainlink node")
		}
	}
}

// CreateChainConfigFromNetworks creates a list of ChainConfig from the network config provided in test config.
// It either creates it from the private ethereum networks created by the test environment or from the
// network URLs provided in the network config ( if the network is a live testnet).
// It uses the private keys from the network config to create the deployer key for each chain.
func CreateChainConfigFromNetworks(
	t *testing.T,
	env *test_env.CLClusterTestEnv,
	privateEthereumNetworks []*ctf_config.EthereumNetworkConfig,
	networkConfig *ctf_config.NetworkConfig,
) []ChainConfig {
	evmNetworks := networks.MustGetSelectedNetworkConfig(networkConfig)
	networkPvtKeys := make(map[int64]string)
	for _, net := range evmNetworks {
		require.Greater(t, len(net.PrivateKeys), 0, "No private keys found for network")
		networkPvtKeys[net.ChainID] = net.PrivateKeys[0]
	}
	var chains []ChainConfig
	// if private ethereum networks are not provided, we will create chains from the network URLs
	if len(privateEthereumNetworks) == 0 {
		for _, net := range evmNetworks {
			chainId := net.ChainID
			chainName, err := chainselectors.NameFromChainId(uint64(chainId))
			require.NoError(t, err, "Error getting chain name")
			pvtKeyStr, exists := networkPvtKeys[chainId]
			require.Truef(t, exists, "Private key not found for chain id %d", chainId)
			pvtKey, err := crypto.HexToECDSA(pvtKeyStr)
			require.NoError(t, err)
			deployer, err := bind.NewKeyedTransactorWithChainID(pvtKey, big.NewInt(chainId))
			require.NoError(t, err)
			deployer.GasLimit = net.DefaultGasLimit
			chains = append(chains, ChainConfig{
				ChainID:     uint64(chainId),
				ChainName:   chainName,
				ChainType:   "EVM",
				WSRPCs:      net.URLs,
				HTTPRPCs:    net.HTTPURLs,
				DeployerKey: deployer,
			})
		}
		return chains
	}
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
			ChainID:     uint64(chainId),
			ChainName:   chainName,
			ChainType:   EVMChainType,
			WSRPCs:      rpcProvider.PublicWsUrls(),
			HTTPRPCs:    rpcProvider.PublicHttpUrls(),
			DeployerKey: deployer,
		})
	}
	return chains
}

// RestartChainlinkNodes restarts the chainlink nodes in the test environment
func RestartChainlinkNodes(t *testing.T, env *test_env.CLClusterTestEnv) error {
	errGrp := errgroup.Group{}
	if env == nil || env.ClCluster == nil {
		return errors.Wrap(errors.New("no testenv or clcluster found "), "error restarting node")
	}
	for _, n := range env.ClCluster.Nodes {
		n := n
		errGrp.Go(func() error {
			if err := n.Container.Terminate(testcontext.Get(t)); err != nil {
				return err
			}
			err := n.RestartContainer()
			if err != nil {
				return err
			}
			return nil
		})

	}
	return errGrp.Wait()
}
