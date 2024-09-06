package devenv

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AlekSi/pointer"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	ctftestenv "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/stretchr/testify/require"
	"github.com/subosito/gotenv"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

const (
	DevEnv = "devenv"
)

type EnvironmentConfig struct {
	Chains   []ChainConfig
	nodeInfo []NodeInfo
	JDConfig JDConfig
}

func NewEnvironment(ctx context.Context, lggr logger.Logger, config EnvironmentConfig) (*deployment.Environment, error) {
	chains, err := NewChains(lggr, config.Chains)
	if err != nil {
		return nil, fmt.Errorf("failed to create chains: %w", err)
	}
	offChain, err := NewJDClient(config.JDConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create JD client: %w", err)
	}

	jd, ok := offChain.(JobDistributor)
	if !ok {
		return nil, fmt.Errorf("offchain client does not implement JobDistributor")
	}
	don, err := NewRegisteredDON(ctx, config.nodeInfo, jd)
	if err != nil {
		return nil, fmt.Errorf("failed to create registered DON: %w", err)
	}
	nodeIDs := don.NodeIds()

	err = don.CreateSupportedChains(ctx, config.Chains)
	if err != nil {
		return nil, err
	}

	return &deployment.Environment{
		Name:     DevEnv,
		Offchain: offChain,
		NodeIDs:  nodeIDs,
		Chains:   chains,
		Logger:   lggr,
	}, nil
}

// DeployPrivateChains deploys private chains and returns the chain configs and a deploy function which
// can be used to deploy the Chainlink nodes.
func DeployPrivateChains(t *testing.T) ([]ChainConfig, func() error) {
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
	deployCL := func() error {
		evmNetworks := networks.MustGetSelectedNetworkConfig(cfg.GetNetworkConfig())
		noOfNodes := pointer.GetInt(cfg.CCIP.CLNode.NoOfPluginNodes) + pointer.GetInt(cfg.CCIP.CLNode.NoOfBootstraps)
		for i := 1; i <= noOfNodes; i++ {
			toml, _, err := testsetups.SetNodeConfig(
				evmNetworks,
				cfg.NodeConfig.BaseConfigTOML,
				cfg.NodeConfig.CommonChainConfigTOML,
				cfg.NodeConfig.ChainConfigTOMLByChainID,
			)
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
		return env.ClCluster.Start()
	}
	return chains, deployCL
}
