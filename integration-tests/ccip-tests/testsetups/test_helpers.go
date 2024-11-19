package testsetups

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	ctfconfig "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	ctftestenv "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipdeployment "github.com/smartcontractkit/chainlink/deployment/ccip"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	clclient "github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/subosito/gotenv"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials/insecure"
)

// DeployedLocalDevEnvironment is a helper struct for setting up a local dev environment with docker
type DeployedLocalDevEnvironment struct {
	ccipdeployment.DeployedEnv
	testEnv *test_env.CLClusterTestEnv
	DON     *devenv.DON
}

func (d DeployedLocalDevEnvironment) RestartChainlinkNodes(t *testing.T) error {
	errGrp := errgroup.Group{}
	for _, n := range d.testEnv.ClCluster.Nodes {
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

func NewLocalDevEnvironmentWithDefaultPrice(
	t *testing.T,
	lggr logger.Logger) (ccipdeployment.DeployedEnv, *test_env.CLClusterTestEnv, testconfig.TestConfig) {
	return NewLocalDevEnvironment(t, lggr, ccipdeployment.MockLinkPrice, ccipdeployment.MockWethPrice)
}

func NewLocalDevEnvironment(
	t *testing.T,
	lggr logger.Logger,
	linkPrice, wethPrice *big.Int) (ccipdeployment.DeployedEnv, *test_env.CLClusterTestEnv, testconfig.TestConfig) {
	ctx := testcontext.Get(t)
	// create a local docker environment with simulated chains and job-distributor
	// we cannot create the chainlink nodes yet as we need to deploy the capability registry first
	envConfig, testEnv, cfg := CreateDockerEnv(t)
	require.NotNil(t, envConfig)
	require.NotEmpty(t, envConfig.Chains, "chainConfigs should not be empty")
	require.NotEmpty(t, envConfig.JDConfig, "jdUrl should not be empty")
	chains, err := devenv.NewChains(lggr, envConfig.Chains)
	require.NoError(t, err)
	// locate the home chain
	homeChainSel := envConfig.HomeChainSelector
	require.NotEmpty(t, homeChainSel, "homeChainSel should not be empty")
	feedSel := envConfig.FeedChainSelector
	require.NotEmpty(t, feedSel, "feedSel should not be empty")
	replayBlocks, err := ccipdeployment.LatestBlocksByChain(ctx, chains)
	require.NoError(t, err)

	ab := deployment.NewMemoryAddressBook()
	crConfig := ccipdeployment.DeployTestContracts(t, lggr, ab, homeChainSel, feedSel, chains, linkPrice, wethPrice)

	// start the chainlink nodes with the CR address
	err = StartChainlinkNodes(t, envConfig,
		crConfig,
		testEnv, cfg)
	require.NoError(t, err)
	e, don, err := devenv.NewEnvironment(ctx, lggr, *envConfig)
	require.NoError(t, err)
	require.NotNil(t, e)
	e.ExistingAddresses = ab
	e.MockAdapter = testEnv.MockAdapter

	envNodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	require.NoError(t, err)
	_, err = ccipdeployment.DeployHomeChain(lggr, *e, e.ExistingAddresses, chains[homeChainSel],
		ccipdeployment.NewTestRMNStaticConfig(),
		ccipdeployment.NewTestRMNDynamicConfig(),
		ccipdeployment.NewTestNodeOperator(chains[homeChainSel].DeployerKey.From),
		map[string][][32]byte{
			"NodeOperator": envNodes.NonBootstraps().PeerIDs(),
		},
	)
	require.NoError(t, err)
	zeroLogLggr := logging.GetTestLogger(t)
	// fund the nodes
	FundNodes(t, zeroLogLggr, testEnv, cfg, don.PluginNodes())

	return ccipdeployment.DeployedEnv{
		Env:          *e,
		HomeChainSel: homeChainSel,
		FeedChainSel: feedSel,
		ReplayBlocks: replayBlocks,
	}, testEnv, cfg
}

func NewLocalDevEnvironmentWithRMN(
	t *testing.T,
	lggr logger.Logger,
	numRmnNodes int,
) (ccipdeployment.DeployedEnv, devenv.RMNCluster) {
	tenv, dockerenv, testCfg := NewLocalDevEnvironmentWithDefaultPrice(t, lggr)
	state, err := ccipdeployment.LoadOnchainState(tenv.Env)
	require.NoError(t, err)

	output, err := changeset.DeployPrerequisites(tenv.Env, changeset.DeployPrerequisiteConfig{
		ChainSelectors: tenv.Env.AllChainSelectors(),
	})
	require.NoError(t, err)
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(output.AddressBook))

	// Deploy CCIP contracts.
	newAddresses := deployment.NewMemoryAddressBook()
	err = ccipdeployment.DeployCCIPContracts(tenv.Env, newAddresses, ccipdeployment.DeployCCIPContractConfig{
		HomeChainSel:   tenv.HomeChainSel,
		FeedChainSel:   tenv.FeedChainSel,
		ChainsToDeploy: tenv.Env.AllChainSelectors(),
		TokenConfig:    ccipdeployment.NewTestTokenConfig(state.Chains[tenv.FeedChainSel].USDFeeds),
		MCMSConfig:     ccipdeployment.NewTestMCMSConfig(t, tenv.Env),
		OCRSecrets:     deployment.XXXGenerateTestOCRSecrets(),
	})
	require.NoError(t, err)
	require.NoError(t, tenv.Env.ExistingAddresses.Merge(newAddresses))

	l := logging.GetTestLogger(t)
	config := GenerateTestRMNConfig(t, numRmnNodes, tenv, MustNetworksToRPCMap(dockerenv.EVMNetworks))
	require.NotNil(t, testCfg.CCIP)
	rmnCluster, err := devenv.NewRMNCluster(
		t, l,
		[]string{dockerenv.DockerNetwork.ID},
		config,
		testCfg.CCIP.RMNConfig.GetProxyImage(),
		testCfg.CCIP.RMNConfig.GetProxyVersion(),
		testCfg.CCIP.RMNConfig.GetAFN2ProxyImage(),
		testCfg.CCIP.RMNConfig.GetAFN2ProxyVersion(),
		dockerenv.LogStream,
	)
	require.NoError(t, err)
	return tenv, *rmnCluster
}

func MustNetworksToRPCMap(evmNetworks []*blockchain.EVMNetwork) map[uint64]string {
	rpcs := make(map[uint64]string)
	for _, network := range evmNetworks {
		sel, err := chainsel.SelectorFromChainId(uint64(network.ChainID))
		if err != nil {
			panic(err)
		}
		rpcs[sel] = network.HTTPURLs[0]
	}
	return rpcs
}

func MustCCIPNameToRMNName(a string) string {
	m := map[string]string{
		chainsel.GETH_TESTNET.Name:  "DevnetAlpha",
		chainsel.GETH_DEVNET_2.Name: "DevnetBeta",
		// TODO: Add more as needed.
	}
	v, ok := m[a]
	if !ok {
		panic(fmt.Sprintf("no mapping for %s", a))
	}
	return v
}

func GenerateTestRMNConfig(t *testing.T, nRMNNodes int, tenv ccipdeployment.DeployedEnv, rpcMap map[uint64]string) map[string]devenv.RMNConfig {
	// Find the bootstrappers.
	nodes, err := deployment.NodeInfo(tenv.Env.NodeIDs, tenv.Env.Offchain)
	require.NoError(t, err)
	bootstrappers := nodes.BootstrapLocators()

	// Just set all RMN nodes to support all chains.
	state, err := ccipdeployment.LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	var chainParams []devenv.ChainParam
	var remoteChains []devenv.RemoteChains

	var rpcs []devenv.Chain
	for chainSel, chain := range state.Chains {
		c, _ := chainsel.ChainBySelector(chainSel)
		rmnName := MustCCIPNameToRMNName(c.Name)
		chainParams = append(chainParams, devenv.ChainParam{
			Name: rmnName,
			Stability: devenv.Stability{
				Type:              "ConfirmationDepth",
				SoftConfirmations: 0,
				HardConfirmations: 0,
			},
		})
		remoteChains = append(remoteChains, devenv.RemoteChains{
			Name:                   rmnName,
			OnRampStartBlockNumber: 0,
			OffRamp:                chain.OffRamp.Address().String(),
			OnRamp:                 chain.OnRamp.Address().String(),
			RMNRemote:              chain.RMNRemote.Address().String(),
		})
		rpcs = append(rpcs, devenv.Chain{
			Name: rmnName,
			RPC:  rpcMap[chainSel],
		})
	}
	hc, _ := chainsel.ChainBySelector(tenv.HomeChainSel)
	shared := devenv.SharedConfig{
		Networking: devenv.SharedConfigNetworking{
			Bootstrappers: bootstrappers,
		},
		HomeChain: devenv.HomeChain{
			Name:                 MustCCIPNameToRMNName(hc.Name),
			CapabilitiesRegistry: state.Chains[tenv.HomeChainSel].CapabilityRegistry.Address().String(),
			CCIPConfig:           state.Chains[tenv.HomeChainSel].CCIPHome.Address().String(),
			RMNHome:              state.Chains[tenv.HomeChainSel].RMNHome.Address().String(),
		},
		RemoteChains: remoteChains,
		ChainParams:  chainParams,
	}

	rmnConfig := make(map[string]devenv.RMNConfig)
	for i := 0; i < nRMNNodes; i++ {
		// Listen addresses _should_ be able to operator on the same port since
		// they are inside the docker network.
		proxyLocal := devenv.ProxyLocalConfig{
			ListenAddresses:   []string{devenv.DefaultProxyListenAddress},
			AnnounceAddresses: []string{},
			ProxyAddress:      devenv.DefaultRageProxy,
			DiscovererDbPath:  devenv.DefaultDiscovererDbPath,
		}
		rmnConfig[fmt.Sprintf("rmn_%d", i)] = devenv.RMNConfig{
			Shared: shared,
			Local: devenv.LocalConfig{
				Networking: devenv.LocalConfigNetworking{
					RageProxy: devenv.DefaultRageProxy,
				},
				Chains: rpcs,
			},
			ProxyShared: devenv.DefaultRageProxySharedConfig,
			ProxyLocal:  proxyLocal,
		}
	}
	return rmnConfig
}

// CreateDockerEnv creates a new test environment with simulated private ethereum networks and job distributor
// It returns the EnvironmentConfig which holds the chain config and JD config
// The test environment is then used to start chainlink nodes
func CreateDockerEnv(t *testing.T) (
	*devenv.EnvironmentConfig,
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
	var privateEthereumNetworks []*ctfconfig.EthereumNetworkConfig
	for _, name := range cfg.GetNetworkConfig().SelectedNetworks {
		if network, exists := cfg.CCIP.PrivateEthereumNetworks[name]; exists {
			privateEthereumNetworks = append(privateEthereumNetworks, network)
		}
	}

	builder := test_env.NewCLTestEnvBuilder().
		WithTestConfig(&cfg).
		WithTestInstance(t).
		WithMockAdapter().
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

	jdConfig := devenv.JDConfig{
		GRPC:  cfg.CCIP.JobDistributorConfig.GetJDGRPC(),
		WSRPC: cfg.CCIP.JobDistributorConfig.GetJDWSRPC(),
	}
	// TODO : move this as a part of test_env setup with an input in testconfig
	// if JD is not provided, we will spin up a new JD
	if jdConfig.GRPC == "" || jdConfig.WSRPC == "" {
		jd := env.JobDistributor
		require.NotNil(t, jd, "JD is not found in test environment")
		jdConfig = devenv.JDConfig{
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

	return &devenv.EnvironmentConfig{
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
	envConfig *devenv.EnvironmentConfig,
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
	var nodeInfo []devenv.NodeInfo
	for i := 1; i <= noOfNodes; i++ {
		if i <= pointer.GetInt(cfg.CCIP.CLNode.NoOfBootstraps) {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: true,
				Name:        fmt.Sprintf("bootstrap-%d", i),
				// TODO : make this configurable
				P2PPort: "6690",
			})
		} else {
			nodeInfo = append(nodeInfo, devenv.NodeInfo{
				IsBootstrap: false,
				Name:        fmt.Sprintf("node-%d", i-1),
				// TODO : make this configurable
				P2PPort: "6690",
			})
		}
		toml, _, err := SetNodeConfig(
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
		envConfig = &devenv.EnvironmentConfig{}
	}
	envConfig.JDConfig.NodeInfo = nodeInfo
	return nil
}

// FundNodes sends funds to the chainlink nodes based on the provided test config
// It also sets up a clean-up function to return the funds back to the deployer account once the test is done
// It assumes that the chainlink nodes are already started and the account addresses for all chains are available
func FundNodes(t *testing.T, lggr zerolog.Logger, env *test_env.CLClusterTestEnv, cfg tc.TestConfig, nodes []devenv.Node) {
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
	privateEthereumNetworks []*ctfconfig.EthereumNetworkConfig,
	networkConfig *ctfconfig.NetworkConfig,
) []devenv.ChainConfig {
	evmNetworks := networks.MustGetSelectedNetworkConfig(networkConfig)
	networkPvtKeys := make(map[int64]string)
	for _, net := range evmNetworks {
		require.Greater(t, len(net.PrivateKeys), 0, "No private keys found for network")
		networkPvtKeys[net.ChainID] = net.PrivateKeys[0]
	}
	var chains []devenv.ChainConfig
	// if private ethereum networks are not provided, we will create chains from the network URLs
	if len(privateEthereumNetworks) == 0 {
		for _, net := range evmNetworks {
			chainId := net.ChainID
			chainName, err := chainsel.NameFromChainId(uint64(chainId))
			require.NoError(t, err, "Error getting chain name")
			pvtKeyStr, exists := networkPvtKeys[chainId]
			require.Truef(t, exists, "Private key not found for chain id %d", chainId)
			pvtKey, err := crypto.HexToECDSA(pvtKeyStr)
			require.NoError(t, err)
			deployer, err := bind.NewKeyedTransactorWithChainID(pvtKey, big.NewInt(chainId))
			require.NoError(t, err)
			deployer.GasLimit = net.DefaultGasLimit
			chains = append(chains, devenv.ChainConfig{
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
		chainName, err := chainsel.NameFromChainId(uint64(chainId))
		require.NoError(t, err, "Error getting chain name")
		rpcProvider, err := env.GetRpcProvider(int64(chainId))
		require.NoError(t, err, "Error getting rpc provider")
		pvtKeyStr, exists := networkPvtKeys[int64(chainId)]
		require.Truef(t, exists, "Private key not found for chain id %d", chainId)
		pvtKey, err := crypto.HexToECDSA(pvtKeyStr)
		require.NoError(t, err)
		deployer, err := bind.NewKeyedTransactorWithChainID(pvtKey, big.NewInt(int64(chainId)))
		require.NoError(t, err)
		chains = append(chains, devenv.ChainConfig{
			ChainID:     uint64(chainId),
			ChainName:   chainName,
			ChainType:   devenv.EVMChainType,
			WSRPCs:      rpcProvider.PublicWsUrls(),
			HTTPRPCs:    rpcProvider.PublicHttpUrls(),
			DeployerKey: deployer,
		})
	}
	return chains
}
