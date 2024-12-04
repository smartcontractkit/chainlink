package ccip

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	chainsel "github.com/smartcontractkit/chain-selectors"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
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
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	integrationnodes "github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/subosito/gotenv"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink/deployment/environment/devenv"
	clclient "github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	ccipactions "github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// DeployedLocalDevEnvironment is a helper struct for setting up a local dev environment with docker
type DeployedLocalDevEnvironment struct {
	changeset.DeployedEnv
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

func NewLocalDevEnvironmentWithDefaultPrice(t *testing.T, lggr logger.Logger, tCfg *changeset.TestConfigs) (changeset.DeployedEnv, *test_env.CLClusterTestEnv, tc.TestConfig) {
	return NewLocalDevEnvironment(t, lggr, changeset.MockLinkPrice, changeset.MockWethPrice, tCfg)
}

func NewLocalDevEnvironment(
	t *testing.T,
	lggr logger.Logger,
	linkPrice, wethPrice *big.Int,
	tCfg *changeset.TestConfigs,
) (changeset.DeployedEnv, *test_env.CLClusterTestEnv, tc.TestConfig) {
	if tCfg == nil {
		// set to the default constructed value
		tCfg = &changeset.TestConfigs{}
	}

	ctx := testcontext.Get(t)
	// create a local docker environment with simulated chains and job-distributor
	// we cannot create the chainlink nodes yet as we need to deploy the capability registry first
	envConfig, testEnv, cfg := CreateDockerEnv(t)
	require.NotNil(t, envConfig)
	require.NotEmpty(t, envConfig.Chains, "chainConfigs should not be empty")
	require.NotEmpty(t, envConfig.JDConfig, "jdUrl should not be empty")
	users := make(map[uint64][]*bind.TransactOpts)
	for _, chain := range envConfig.Chains {
		sel, err := chainsel.SelectorFromChainId(chain.ChainID)
		require.NoError(t, err)
		users[sel] = chain.Users
	}
	chains, err := devenv.NewChains(lggr, envConfig.Chains)
	require.NoError(t, err)
	// locate the home chain
	homeChainSel := cfg.CCIP.GetHomeChainSelector()
	require.NotEmpty(t, homeChainSel, "homeChainSel should not be empty")
	feedSel := cfg.CCIP.GetFeedChainSelector()
	require.NotEmpty(t, feedSel, "feedSel should not be empty")
	replayBlocks, err := changeset.LatestBlocksByChain(ctx, chains)
	require.NoError(t, err)

	ab := deployment.NewMemoryAddressBook()
	crConfig := changeset.DeployTestContracts(t, lggr, ab, homeChainSel, feedSel, chains, linkPrice, wethPrice)

	// start the chainlink nodes with the CR address
	err = StartChainlinkNodes(t, envConfig,
		crConfig,
		testEnv, cfg)
	require.NoError(t, err)

	e, don, err := devenv.NewEnvironment(func() context.Context { return ctx }, lggr, *envConfig)
	require.NoError(t, err)
	require.NotNil(t, e)
	e.ExistingAddresses = ab

	// fund the nodes
	zeroLogLggr := logging.GetTestLogger(t)
	FundNodes(t, zeroLogLggr, testEnv, cfg, don.PluginNodes())

	env := *e
	envNodes, err := deployment.NodeInfo(env.NodeIDs, env.Offchain)
	require.NoError(t, err)
	allChains := env.AllChainSelectors()
	var usdcChains []uint64
	if tCfg.IsUSDC {
		usdcChains = allChains
	}
	mcmsCfg := make(map[uint64]commontypes.MCMSWithTimelockConfig)
	for _, c := range env.AllChainSelectors() {
		mcmsCfg[c] = commontypes.MCMSWithTimelockConfig{
			Canceller:         commonchangeset.SingleGroupMCMS(t),
			Bypasser:          commonchangeset.SingleGroupMCMS(t),
			Proposer:          commonchangeset.SingleGroupMCMS(t),
			TimelockExecutors: env.AllDeployerKeys(),
			TimelockMinDelay:  big.NewInt(0),
		}
	}
	// Need to deploy prerequisites first so that we can form the USDC config
	// no proposals to be made, timelock can be passed as nil here
	env, err = commonchangeset.ApplyChangesets(t, env, nil, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployHomeChain),
			Config: changeset.DeployHomeChainConfig{
				HomeChainSel:     homeChainSel,
				RMNStaticConfig:  changeset.NewTestRMNStaticConfig(),
				RMNDynamicConfig: changeset.NewTestRMNDynamicConfig(),
				NodeOperators:    changeset.NewTestNodeOperator(chains[homeChainSel].DeployerKey.From),
				NodeP2PIDsPerNodeOpAdmin: map[string][][32]byte{
					"NodeOperator": envNodes.NonBootstraps().PeerIDs(),
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployPrerequisites),
			Config: changeset.DeployPrerequisiteConfig{
				ChainSelectors: allChains,
				Opts: []changeset.PrerequisiteOpt{
					changeset.WithUSDCChains(usdcChains),
					changeset.WithMulticall3(tCfg.IsMultiCall3),
				},
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(commonchangeset.DeployMCMSWithTimelock),
			Config:    mcmsCfg,
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.DeployChainContracts),
			Config: changeset.DeployChainContractsConfig{
				ChainSelectors:    allChains,
				HomeChainSelector: homeChainSel,
			},
		},
	})
	require.NoError(t, err)
	state, err := changeset.LoadOnchainState(env)
	require.NoError(t, err)

	var tokenDataProviders []pluginconfig.TokenDataObserverConfig
	if len(usdcChains) > 0 {
		var endpoint string
		err = ccipactions.SetMockServerWithUSDCAttestation(testEnv.MockAdapter, nil, tCfg.IsUSDCAttestationMissing)
		require.NoError(t, err)
		endpoint = testEnv.MockAdapter.InternalEndpoint
		cctpContracts := make(map[cciptypes.ChainSelector]pluginconfig.USDCCCTPTokenConfig)
		for _, usdcChain := range usdcChains {
			cctpContracts[cciptypes.ChainSelector(usdcChain)] = pluginconfig.USDCCCTPTokenConfig{
				SourcePoolAddress:            state.Chains[usdcChain].USDCTokenPool.Address().String(),
				SourceMessageTransmitterAddr: state.Chains[usdcChain].MockUSDCTransmitter.Address().String(),
			}
		}
		tokenDataProviders = append(tokenDataProviders, pluginconfig.TokenDataObserverConfig{
			Type:    pluginconfig.USDCCCTPHandlerType,
			Version: "1.0",
			USDCCCTPObserverConfig: &pluginconfig.USDCCCTPObserverConfig{
				Tokens:                 cctpContracts,
				AttestationAPI:         endpoint,
				AttestationAPITimeout:  commonconfig.MustNewDuration(time.Second),
				AttestationAPIInterval: commonconfig.MustNewDuration(500 * time.Millisecond),
			}})
	}

	// Build the per chain config.
	tokenConfig := changeset.NewTestTokenConfig(state.Chains[feedSel].USDFeeds)
	chainConfigs := make(map[uint64]changeset.CCIPOCRParams)
	timelocksPerChain := make(map[uint64]*gethwrappers.RBACTimelock)
	for _, chain := range allChains {
		timelocksPerChain[chain] = state.Chains[chain].Timelock
		tokenInfo := tokenConfig.GetTokenInfo(e.Logger, state.Chains[chain].LinkToken, state.Chains[chain].Weth9)
		ocrParams := changeset.DefaultOCRParams(feedSel, tokenInfo, tokenDataProviders)
		if tCfg.OCRConfigOverride != nil {
			ocrParams = tCfg.OCRConfigOverride(ocrParams)
		}
		chainConfigs[chain] = ocrParams
	}

	// Deploy second set of changesets to deploy and configure the CCIP contracts.
	env, err = commonchangeset.ApplyChangesets(t, env, timelocksPerChain, []commonchangeset.ChangesetApplication{
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.ConfigureNewChains),
			Config: changeset.NewChainsConfig{
				HomeChainSel:       homeChainSel,
				FeedChainSel:       feedSel,
				OCRSecrets:         deployment.XXXGenerateTestOCRSecrets(),
				ChainConfigByChain: chainConfigs,
			},
		},
		{
			Changeset: commonchangeset.WrapChangeSet(changeset.CCIPCapabilityJobspec),
		},
	})
	require.NoError(t, err)

	// Ensure capreg logs are up to date.
	changeset.ReplayLogs(t, e.Offchain, replayBlocks)

	return changeset.DeployedEnv{
		Env:          env,
		HomeChainSel: homeChainSel,
		FeedChainSel: feedSel,
		ReplayBlocks: replayBlocks,
		Users:        users,
	}, testEnv, cfg
}

func NewLocalDevEnvironmentWithRMN(
	t *testing.T,
	lggr logger.Logger,
	numRmnNodes int,
) (changeset.DeployedEnv, devenv.RMNCluster) {
	tenv, dockerenv, testCfg := NewLocalDevEnvironmentWithDefaultPrice(t, lggr, nil)
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
		if network.ChainID < 0 {
			panic(fmt.Errorf("negative chain ID: %d", network.ChainID))
		}
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

func GenerateTestRMNConfig(t *testing.T, nRMNNodes int, tenv changeset.DeployedEnv, rpcMap map[uint64]string) map[string]devenv.RMNConfig {
	// Find the bootstrappers.
	nodes, err := deployment.NodeInfo(tenv.Env.NodeIDs, tenv.Env.Offchain)
	require.NoError(t, err)
	bootstrappers := nodes.BootstrapLocators()

	// Just set all RMN nodes to support all chains.
	state, err := changeset.LoadOnchainState(tenv.Env)
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
			CCIPHome:             state.Chains[tenv.HomeChainSel].CCIPHome.Address().String(),
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

	return &devenv.EnvironmentConfig{
		Chains:   chains,
		JDConfig: jdConfig,
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

		toml.Capabilities.ExternalRegistry.NetworkID = ptr.Ptr(registryConfig.NetworkType)
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
	fundGrp := errgroup.Group{}
	for i := range evmNetworks {
		fundGrp.Go(func() error {
			evmNetwork := evmNetworks[i]
			sethClient, err := utils.TestAwareSethClient(t, cfg, &evmNetwork)
			if err != nil {
				return fmt.Errorf("error getting seth client for network %s: %w", evmNetwork.Name, err)
			}
			if len(sethClient.PrivateKeys) == 0 {
				return fmt.Errorf(seth.ErrNoKeyLoaded)
			}
			privateKey := sethClient.PrivateKeys[0]
			if evmNetwork.ChainID < 0 {
				return fmt.Errorf("negative chain ID: %d", evmNetwork.ChainID)
			}
			for _, node := range nodes {
				nodeAddr, ok := node.AccountAddr[uint64(evmNetwork.ChainID)]
				if !ok {
					return fmt.Errorf("account address not found for chain %d", evmNetwork.ChainID)
				}
				fromAddress, err := actions.PrivateKeyToAddress(privateKey)
				if err != nil {
					return fmt.Errorf("error getting address from private key: %w", err)
				}
				amount := big.NewFloat(pointer.GetFloat64(cfg.Common.ChainlinkNodeFunding))
				toAddr := common.HexToAddress(nodeAddr)
				receipt, err := actions.SendFunds(lggr, sethClient, actions.FundsToSendPayload{
					ToAddress:  toAddr,
					Amount:     conversions.EtherToWei(amount),
					PrivateKey: privateKey,
				})
				if err != nil {
					return fmt.Errorf("error sending funds to node %s: %w", node.Name, err)
				}
				if receipt == nil {
					return fmt.Errorf("receipt is nil")
				}
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
			return nil
		})
	}
	require.NoError(t, fundGrp.Wait(), "Error funding chainlink nodes")
}

// CreateChainConfigFromNetworks creates a list of CCIPOCRParams from the network config provided in test config.
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
	networkPvtKeys := make(map[uint64][]string)
	for _, net := range evmNetworks {
		require.Greater(t, len(net.PrivateKeys), 0, "No private keys found for network")
		if net.ChainID < 0 {
			t.Fatalf("negative chain ID: %d", net.ChainID)
		}
		networkPvtKeys[uint64(net.ChainID)] = net.PrivateKeys
	}
	type chainDetails struct {
		chainId  uint64
		wsRPCs   []string
		httpRPCs []string
	}
	var chains []devenv.ChainConfig
	var chaindetails []chainDetails
	if len(privateEthereumNetworks) == 0 {
		for _, net := range evmNetworks {
			chainId := net.ChainID
			if chainId < 0 {
				t.Fatalf("negative chain ID: %d", chainId)
			}
			chaindetails = append(chaindetails, chainDetails{
				chainId:  uint64(chainId),
				wsRPCs:   net.URLs,
				httpRPCs: net.HTTPURLs,
			})
		}
	} else {
		for _, net := range privateEthereumNetworks {
			chainId := net.EthereumChainConfig.ChainID
			if chainId < 0 {
				t.Fatalf("negative chain ID: %d", chainId)
			}
			rpcProvider, err := env.GetRpcProvider(int64(chainId))
			require.NoError(t, err, "Error getting rpc provider")
			chaindetails = append(chaindetails, chainDetails{
				chainId:  uint64(chainId),
				wsRPCs:   rpcProvider.PublicWsUrls(),
				httpRPCs: rpcProvider.PublicHttpUrls(),
			})
		}
	}
	for _, cd := range chaindetails {
		chainId := cd.chainId
		chainName, err := chainsel.NameFromChainId(chainId)
		require.NoError(t, err, "Error getting chain name")
		chainCfg := devenv.ChainConfig{
			ChainID:   chainId,
			ChainName: chainName,
			ChainType: "EVM",
			WSRPCs:    cd.wsRPCs,
			HTTPRPCs:  cd.httpRPCs,
		}
		var pvtKey *string
		// if private keys are provided, use the first private key as deployer key
		// otherwise it will try to load the private key from KMS
		if len(networkPvtKeys[chainId]) > 0 {
			pvtKey = ptr.Ptr(networkPvtKeys[chainId][0])
		}
		require.NoError(t, chainCfg.SetDeployerKey(pvtKey), "Error setting deployer key")
		var additionalPvtKeys []string
		if len(networkPvtKeys[chainId]) > 1 {
			additionalPvtKeys = networkPvtKeys[chainId][1:]
		}
		// if no additional private keys are provided, this will set the users to default deployer key
		require.NoError(t, chainCfg.SetUsers(additionalPvtKeys), "Error setting users")
		chains = append(chains, chainCfg)
	}
	return chains
}

func SetNodeConfig(nets []blockchain.EVMNetwork, nodeConfig, commonChain string, configByChain map[string]string) (*corechainlink.Config, string, error) {
	var tomlCfg *corechainlink.Config
	var err error
	var commonChainConfig *evmcfg.Chain
	if commonChain != "" {
		err = commonconfig.DecodeTOML(bytes.NewReader([]byte(commonChain)), &commonChainConfig)
		if err != nil {
			return nil, "", err
		}
	}
	configByChainMap := make(map[int64]evmcfg.Chain)
	for k, v := range configByChain {
		var chain evmcfg.Chain
		err = commonconfig.DecodeTOML(bytes.NewReader([]byte(v)), &chain)
		if err != nil {
			return nil, "", err
		}
		chainId, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return nil, "", err
		}
		configByChainMap[chainId] = chain
	}
	if nodeConfig == "" {
		tomlCfg = integrationnodes.NewConfig(
			integrationnodes.NewBaseConfig(),
			integrationnodes.WithPrivateEVMs(nets, commonChainConfig, configByChainMap))
	} else {
		tomlCfg, err = integrationnodes.NewConfigFromToml([]byte(nodeConfig), integrationnodes.WithPrivateEVMs(nets, commonChainConfig, configByChainMap))
		if err != nil {
			return nil, "", err
		}
	}
	tomlStr, err := tomlCfg.TOMLString()
	return tomlCfg, tomlStr, err
}
