package devenv

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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

const (
	DevEnv = "devenv"
)

type EnvironmentConfig struct {
	Chains   []ChainConfig
	nodeInfo []NodeInfo
	JDConfig JDConfig
}

func NewEnvironment(ctx context.Context, lggr logger.Logger, config EnvironmentConfig) (*deployment.Environment, *DON, error) {
	chains, err := NewChains(lggr, config.Chains)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create chains: %w", err)
	}
	offChain, err := NewJDClient(config.JDConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create JD client: %w", err)
	}

	jd, ok := offChain.(JobDistributor)
	if !ok {
		return nil, nil, fmt.Errorf("offchain client does not implement JobDistributor")
	}
	don, err := NewRegisteredDON(ctx, config.nodeInfo, jd)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create registered DON: %w", err)
	}
	nodeIDs := don.NodeIds()

	err = don.CreateSupportedChains(ctx, config.Chains)
	if err != nil {
		return nil, nil, err
	}

	return &deployment.Environment{
		Name:     DevEnv,
		Offchain: offChain,
		NodeIDs:  nodeIDs,
		Chains:   chains,
		Logger:   lggr,
	}, don, nil
}

// DeployPrivateChains deploys private chains and returns the chain configs and a function to deploy the Chainlink nodes
func DeployPrivateChains(t *testing.T) (
	*EnvironmentConfig,
	*test_env.CLClusterTestEnv,
	tc.TestConfig,
	func(
		*EnvironmentConfig,
		deployment.RegistryConfig,
		*test_env.CLClusterTestEnv,
		tc.TestConfig,
	) error) {
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
		WithoutCleanup().
		Build()
	require.NoError(t, err, "Error building test environment")
	var chains []ChainConfig
	evmNetworks := networks.MustGetSelectedNetworkConfig(cfg.Network)
	networkPvtKeys := make(map[int64]string)
	for _, net := range evmNetworks {
		require.Greater(t, len(net.PrivateKeys), 0, "No private keys found for network")
		networkPvtKeys[net.ChainID] = net.PrivateKeys[0]
	}
	for _, networkCfg := range cfg.CCIP.PrivateEthereumNetworks {
		chainId := networkCfg.EthereumChainConfig.ChainID
		chainName, err := chain_selectors.NameFromChainId(uint64(chainId))
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
			ChainId:         uint64(chainId),
			ChainName:       chainName,
			ChainType:       "EVM",
			WsRpcs:          rpcProvider.PublicWsUrls(),
			HttpRpcs:        rpcProvider.PublicHttpUrls(),
			PrivateHttpRpcs: rpcProvider.PrivateHttpUrls(),
			PrivateWsRpcs:   rpcProvider.PrivateWsUrsl(),
			DeployerKey:     deployer,
		})
	}
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
	deployCL := func(
		envConfig *EnvironmentConfig,
		registryConfig deployment.RegistryConfig,
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
				})
			} else {
				nodeInfo = append(nodeInfo, NodeInfo{
					IsBootstrap: false,
					Name:        fmt.Sprintf("node-%d", i-1),
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

	return &EnvironmentConfig{
		Chains:   chains,
		JDConfig: jdConfig,
	}, env, cfg, deployCL
}
