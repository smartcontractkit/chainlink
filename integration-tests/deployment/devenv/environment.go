package devenv

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/AlekSi/pointer"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	ctftestenv "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/stretchr/testify/require"
	"github.com/subosito/gotenv"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
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
func DeployPrivateChains(t *testing.T) ([]ChainConfig, string, func([]ChainConfig, string, deployment.RegistryConfig) (*EnvironmentConfig, error)) {
	if info, err := os.Stat("./env"); os.IsNotExist(err) || !info.IsDir() {
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
	for _, networkCfg := range cfg.CCIP.PrivateEthereumNetworks {
		chainId := networkCfg.EthereumChainConfig.ChainID
		chainName, err := chain_selectors.NameFromChainId(uint64(chainId))
		require.NoError(t, err, "Error getting chain name")
		rpcProvider, err := env.GetRpcProvider(int64(chainId))
		require.NoError(t, err, "Error getting rpc provider")
		chains = append(chains, ChainConfig{
			ChainId:         uint64(chainId),
			ChainName:       chainName,
			WsRpcs:          rpcProvider.PublicWsUrls(),
			HttpRpcs:        rpcProvider.PublicHttpUrls(),
			PrivateHttpRpcs: rpcProvider.PrivateHttpUrls(),
			PrivateWsRpcs:   rpcProvider.PrivateWsUrsl(),
		})
	}
	var jdUrl string
	if cfg.CCIP.JobDistributorConfig != nil {

	}
	deployCL := func(chains []ChainConfig, jdUrl string, registryConfig deployment.RegistryConfig) (*EnvironmentConfig, error) {
		evmNetworks := networks.MustGetSelectedNetworkConfig(cfg.GetNetworkConfig())
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
					Name:        fmt.Sprintf("node-%d", i),
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
				return nil, err
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
				return nil, err
			}
			ccipNode.SetTestLogger(t)
			env.ClCluster.Nodes = append(env.ClCluster.Nodes, ccipNode)
		}
		err := env.ClCluster.Start()
		if err != nil {
			return nil, err
		}
		return &EnvironmentConfig{
			Chains:   chains,
			JDConfig: JDConfig{URL: jdUrl},
			nodeInfo: nodeInfo,
		}, nil
	}
	return chains, jdUrl, deployCL
}
