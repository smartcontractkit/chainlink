package testsetups

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	ctf_config_types "github.com/smartcontractkit/chainlink-testing-framework/lib/config/types"
	ctftestenv "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	k8config "github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/foundry"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/parrot"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/reorg"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/conversions"

	evmcfg "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/types/config/node"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	integrationnodes "github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func SetResourceProfile(cpu, mem string) map[string]interface{} {
	return map[string]interface{}{
		"requests": map[string]interface{}{
			"cpu":    cpu,
			"memory": mem,
		},
		"limits": map[string]interface{}{
			"cpu":    cpu,
			"memory": mem,
		},
	}
}

func SetNodeConfig(nets []blockchain.EVMNetwork, nodeConfig, commonChain string, configByChain map[string]string) (*corechainlink.Config, string, error) {
	var tomlCfg *corechainlink.Config
	var err error
	var commonChainConfig *evmcfg.Chain
	if commonChain != "" {
		err = config.DecodeTOML(bytes.NewReader([]byte(commonChain)), &commonChainConfig)
		if err != nil {
			return nil, "", err
		}
	}
	configByChainMap := make(map[int64]evmcfg.Chain)
	for k, v := range configByChain {
		var chain evmcfg.Chain
		err = config.DecodeTOML(bytes.NewReader([]byte(v)), &chain)
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
			node.WithPrivateEVMs(nets, commonChainConfig, configByChainMap))
	} else {
		tomlCfg, err = node.NewConfigFromToml([]byte(nodeConfig), node.WithPrivateEVMs(nets, commonChainConfig, configByChainMap))
		if err != nil {
			return nil, "", err
		}
	}
	tomlStr, err := tomlCfg.TOMLString()
	return tomlCfg, tomlStr, err
}

func ChainlinkPropsForUpdate(
	t *testing.T,
	testInputs *CCIPTestConfig,
) (map[string]any, int) {
	updateProps := make(map[string]any)
	upgradeImage := pointer.GetString(testInputs.EnvInput.NewCLCluster.Common.ChainlinkUpgradeImage.Image)
	upgradeTag := pointer.GetString(testInputs.EnvInput.NewCLCluster.Common.ChainlinkUpgradeImage.Version)
	noOfNodesToUpgrade := 0
	if len(testInputs.EnvInput.NewCLCluster.Nodes) > 0 {
		var nodesMap []map[string]any
		for _, clNode := range testInputs.EnvInput.NewCLCluster.Nodes {
			if !pointer.GetBool(clNode.NeedsUpgrade) {
				continue
			}
			upgradeImage = pointer.GetString(clNode.ChainlinkUpgradeImage.Image)
			upgradeTag = pointer.GetString(clNode.ChainlinkUpgradeImage.Version)
			if upgradeImage == "" || upgradeTag == "" {
				continue
			}
			nodeConfig := clNode.BaseConfigTOML
			commonChainConfig := clNode.CommonChainConfigTOML
			chainConfigByChain := clNode.ChainConfigTOMLByChain
			if nodeConfig == "" {
				nodeConfig = testInputs.EnvInput.NewCLCluster.Common.BaseConfigTOML
			}
			if commonChainConfig == "" {
				commonChainConfig = testInputs.EnvInput.NewCLCluster.Common.CommonChainConfigTOML
			}
			if chainConfigByChain == nil {
				chainConfigByChain = testInputs.EnvInput.NewCLCluster.Common.ChainConfigTOMLByChain
			}

			_, tomlStr, err := SetNodeConfig(
				testInputs.SelectedNetworks,
				nodeConfig, commonChainConfig, chainConfigByChain,
			)
			require.NoError(t, err)
			nodesMap = append(nodesMap, map[string]any{
				"name": clNode.Name,
				"chainlink": map[string]any{
					"image": map[string]any{
						"image":   upgradeImage,
						"version": upgradeTag,
					},
				},
				"toml": tomlStr,
			})
			noOfNodesToUpgrade++
		}
		updateProps["nodes"] = nodesMap
	} else {
		if upgradeImage == "" || upgradeTag == "" {
			return nil, 0
		}
		updateProps["chainlink"] = map[string]interface{}{
			"image": map[string]interface{}{
				"image":   upgradeImage,
				"version": upgradeTag,
			},
		}
		_, tomlStr, err := SetNodeConfig(
			testInputs.SelectedNetworks,
			testInputs.EnvInput.NewCLCluster.Common.BaseConfigTOML,
			testInputs.EnvInput.NewCLCluster.Common.CommonChainConfigTOML,
			testInputs.EnvInput.NewCLCluster.Common.ChainConfigTOMLByChain,
		)
		require.NoError(t, err)
		updateProps["toml"] = tomlStr
		noOfNodesToUpgrade = pointer.GetInt(testInputs.EnvInput.NewCLCluster.NoOfNodes)
	}
	return updateProps, noOfNodesToUpgrade
}

func ChainlinkChart(
	t *testing.T,
	testInputs *CCIPTestConfig,
	nets []blockchain.EVMNetwork,
) environment.ConnectedChart {
	require.NotNil(t, testInputs.EnvInput.NewCLCluster.Common, "Chainlink Common config is not specified")
	clProps := make(map[string]interface{})
	clProps["prometheus"] = true
	var formattedArgs []string
	if len(testInputs.EnvInput.NewCLCluster.DBArgs) > 0 {
		for _, arg := range testInputs.EnvInput.NewCLCluster.DBArgs {
			formattedArgs = append(formattedArgs, "-c")
			formattedArgs = append(formattedArgs, arg)
		}
	}
	clProps["db"] = map[string]interface{}{
		"resources":                        SetResourceProfile(testInputs.EnvInput.NewCLCluster.DBCPU, testInputs.EnvInput.NewCLCluster.DBMemory),
		"additionalArgs":                   formattedArgs,
		"stateful":                         pointer.GetBool(testInputs.EnvInput.NewCLCluster.IsStateful),
		"capacity":                         testInputs.EnvInput.NewCLCluster.DBCapacity,
		"storageClassName":                 pointer.GetString(testInputs.EnvInput.NewCLCluster.DBStorageClass),
		"enablePrometheusPostgresExporter": pointer.GetBool(testInputs.EnvInput.NewCLCluster.PromPgExporter),
		"image": map[string]any{
			"image":   testInputs.EnvInput.NewCLCluster.Common.DBImage,
			"version": testInputs.EnvInput.NewCLCluster.Common.DBTag,
		},
	}
	clProps["chainlink"] = map[string]interface{}{
		"resources": SetResourceProfile(testInputs.EnvInput.NewCLCluster.NodeCPU, testInputs.EnvInput.NewCLCluster.NodeMemory),
		"image": map[string]any{
			"image":   pointer.GetString(testInputs.EnvInput.NewCLCluster.Common.ChainlinkImage.Image),
			"version": pointer.GetString(testInputs.EnvInput.NewCLCluster.Common.ChainlinkImage.Version),
		},
	}

	require.NotNil(t, testInputs.EnvInput, "no env test input specified")

	if len(testInputs.EnvInput.NewCLCluster.Nodes) > 0 {
		var nodesMap []map[string]any
		for _, clNode := range testInputs.EnvInput.NewCLCluster.Nodes {
			nodeConfig := clNode.BaseConfigTOML
			commonChainConfig := clNode.CommonChainConfigTOML
			chainConfigByChain := clNode.ChainConfigTOMLByChain
			if nodeConfig == "" {
				nodeConfig = testInputs.EnvInput.NewCLCluster.Common.BaseConfigTOML
			}
			if commonChainConfig == "" {
				commonChainConfig = testInputs.EnvInput.NewCLCluster.Common.CommonChainConfigTOML
			}
			if chainConfigByChain == nil {
				chainConfigByChain = testInputs.EnvInput.NewCLCluster.Common.ChainConfigTOMLByChain
			}

			_, tomlStr, err := SetNodeConfig(nets, nodeConfig, commonChainConfig, chainConfigByChain)
			require.NoError(t, err)
			nodesMap = append(nodesMap, map[string]any{
				"name": clNode.Name,
				"chainlink": map[string]any{
					"image": map[string]any{
						"image":   pointer.GetString(clNode.ChainlinkImage.Image),
						"version": pointer.GetString(clNode.ChainlinkImage.Version),
					},
				},
				"db": map[string]any{
					"image": map[string]any{
						"image":   clNode.DBImage,
						"version": clNode.DBTag,
					},
					"storageClassName": "gp3",
				},
				"toml": tomlStr,
			})
		}
		clProps["nodes"] = nodesMap
		return chainlink.New(0, clProps)
	}
	clProps["replicas"] = pointer.GetInt(testInputs.EnvInput.NewCLCluster.NoOfNodes)

	_, tomlStr, err := SetNodeConfig(
		nets,
		testInputs.EnvInput.NewCLCluster.Common.BaseConfigTOML,
		testInputs.EnvInput.NewCLCluster.Common.CommonChainConfigTOML,
		testInputs.EnvInput.NewCLCluster.Common.ChainConfigTOMLByChain,
	)
	require.NoError(t, err)
	clProps["toml"] = tomlStr
	return chainlink.New(0, clProps)
}

func DeployLocalCluster(
	t *testing.T,
	testInputs *CCIPTestConfig,
) (*test_env.CLClusterTestEnv, func() error) {
	selectedNetworks := testInputs.SelectedNetworks

	privateEthereumNetworks := []*ctf_config.EthereumNetworkConfig{}
	for _, network := range testInputs.EnvInput.PrivateEthereumNetworks {
		privateEthereumNetworks = append(privateEthereumNetworks, network)

		for _, networkCfg := range networks.MustGetSelectedNetworkConfig(testInputs.EnvInput.Network) {
			for _, key := range networkCfg.PrivateKeys {
				address, err := conversions.PrivateKeyHexToAddress(key)
				require.NoError(t, err, "failed to convert private key to address: %w", err)
				network.EthereumChainConfig.AddressesToFund = append(
					network.EthereumChainConfig.AddressesToFund, address.Hex(),
				)
			}
		}
	}

	if len(selectedNetworks) > len(privateEthereumNetworks) {
		seen := make(map[int64]bool)
		missing := []blockchain.EVMNetwork{}

		for _, network := range privateEthereumNetworks {
			seen[int64(network.EthereumChainConfig.ChainID)] = true
		}

		for _, network := range selectedNetworks {
			if !seen[network.ChainID] {
				missing = append(missing, network)
			}
		}

		for _, network := range missing {
			chainConfig := &ctf_config.EthereumChainConfig{}
			err := chainConfig.Default()
			if err != nil {
				require.NoError(t, err, "failed to get default chain config: %w", err)
			} else {
				chainConfig.ChainID = int(network.ChainID)

				eth1 := ctf_config_types.EthereumVersion_Eth1
				geth := ctf_config_types.ExecutionLayer_Geth

				privateEthereumNetworks = append(privateEthereumNetworks, &ctf_config.EthereumNetworkConfig{
					EthereumVersion:     &eth1,
					ExecutionLayer:      &geth,
					EthereumChainConfig: chainConfig,
				})
			}
		}

		require.Equal(t, len(selectedNetworks), len(privateEthereumNetworks), "failed to create undefined selected networks. Maybe some of them had the same chain ids?")
	}

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestConfig(testInputs.EnvInput).
		WithTestInstance(t).
		WithPrivateEthereumNetworks(privateEthereumNetworks).
		WithMockAdapter().
		WithoutCleanup().
		Build()
	require.NoError(t, err)
	// the builder builds network with a static network config, we don't want that.
	env.EVMNetworks = []*blockchain.EVMNetwork{}
	for i, networkCfg := range selectedNetworks {
		rpcProvider, err := env.GetRpcProvider(networkCfg.ChainID)
		require.NoError(t, err, "Error getting rpc provider")
		selectedNetworks[i].URLs = rpcProvider.PrivateWsUrsl()
		selectedNetworks[i].HTTPURLs = rpcProvider.PrivateHttpUrls()
		newNetwork := networkCfg
		newNetwork.URLs = rpcProvider.PublicWsUrls()
		newNetwork.HTTPURLs = rpcProvider.PublicHttpUrls()
		env.EVMNetworks = append(env.EVMNetworks, &newNetwork)
	}
	testInputs.SelectedNetworks = selectedNetworks

	// a func to start the CL nodes asynchronously
	deployCL := func() error {
		noOfNodes := pointer.GetInt(testInputs.EnvInput.NewCLCluster.NoOfNodes)
		if env.ClCluster == nil {
			env.ClCluster = &test_env.ClCluster{}
		}
		// if individual nodes are specified, then deploy them with specified configs
		if len(testInputs.EnvInput.NewCLCluster.Nodes) > 0 {
			for _, clNode := range testInputs.EnvInput.NewCLCluster.Nodes {
				toml, _, err := SetNodeConfig(
					selectedNetworks,
					clNode.BaseConfigTOML,
					clNode.CommonChainConfigTOML,
					clNode.ChainConfigTOMLByChain,
				)
				if err != nil {
					return err
				}
				ccipNode, err := test_env.NewClNode(
					[]string{env.DockerNetwork.Name},
					pointer.GetString(clNode.ChainlinkImage.Image),
					pointer.GetString(clNode.ChainlinkImage.Version),
					toml,
					test_env.WithPgDBOptions(
						ctftestenv.WithPostgresImageName(clNode.DBImage),
						ctftestenv.WithPostgresImageVersion(clNode.DBTag),
					),
				)
				if err != nil {
					return err
				}
				ccipNode.SetTestLogger(t)
				env.ClCluster.Nodes = append(env.ClCluster.Nodes, ccipNode)
			}
		} else {
			// if no individual nodes are specified, then deploy the number of nodes specified in the env input with common config
			for i := 0; i < noOfNodes; i++ {
				toml, _, err := SetNodeConfig(
					selectedNetworks,
					testInputs.EnvInput.NewCLCluster.Common.BaseConfigTOML,
					testInputs.EnvInput.NewCLCluster.Common.CommonChainConfigTOML,
					testInputs.EnvInput.NewCLCluster.Common.ChainConfigTOMLByChain,
				)
				if err != nil {
					return err
				}
				ccipNode, err := test_env.NewClNode(
					[]string{env.DockerNetwork.Name},
					pointer.GetString(testInputs.EnvInput.NewCLCluster.Common.ChainlinkImage.Image),
					pointer.GetString(testInputs.EnvInput.NewCLCluster.Common.ChainlinkImage.Version),
					toml,
					test_env.WithPgDBOptions(
						ctftestenv.WithPostgresImageName(testInputs.EnvInput.NewCLCluster.Common.DBImage),
						ctftestenv.WithPostgresImageVersion(testInputs.EnvInput.NewCLCluster.Common.DBTag),
					),
				)
				if err != nil {
					return err
				}
				ccipNode.SetTestLogger(t)
				env.ClCluster.Nodes = append(env.ClCluster.Nodes, ccipNode)
			}
		}
		return env.ClCluster.Start()
	}
	return env, deployCL
}

// UpgradeNodes restarts chainlink nodes in the given range with upgrade image
// startIndex and endIndex are inclusive
func UpgradeNodes(
	t *testing.T,
	lggr *zerolog.Logger,
	testInputs *CCIPTestConfig,
	ccipEnv *actions.CCIPTestEnv,
) error {
	lggr.Info().
		Msg("Upgrading node version")
	// if the test is running on local docker
	if pointer.GetBool(testInputs.TestGroupInput.LocalCluster) {
		env := ccipEnv.LocalCluster
		for i, clNode := range env.ClCluster.Nodes {
			upgradeImage := pointer.GetString(testInputs.EnvInput.NewCLCluster.Common.ChainlinkUpgradeImage.Image)
			upgradeTag := pointer.GetString(testInputs.EnvInput.NewCLCluster.Common.ChainlinkUpgradeImage.Version)
			// if individual node upgrade image is provided, use that
			if len(testInputs.EnvInput.NewCLCluster.Nodes) > 0 {
				if i < len(testInputs.EnvInput.NewCLCluster.Nodes) {
					upgradeImage = pointer.GetString(testInputs.EnvInput.NewCLCluster.Nodes[i].ChainlinkUpgradeImage.Image)
					upgradeTag = pointer.GetString(testInputs.EnvInput.NewCLCluster.Nodes[i].ChainlinkUpgradeImage.Version)
				}
			}
			if upgradeImage == "" || upgradeTag == "" {
				continue
			}
			err := clNode.UpgradeVersion(upgradeImage, upgradeTag)
			if err != nil {
				return err
			}
		}
	} else {
		// if the test is running on k8s
		k8Env := ccipEnv.K8Env
		if k8Env == nil {
			return errors.New("k8s environment is nil, cannot restart nodes")
		}
		props, noOfNodesToUpgrade := ChainlinkPropsForUpdate(t, testInputs)
		chartName := ccipEnv.CLNodes[0].ChartName
		// explicitly set the env var into false to allow manifest update
		// if tests are run in remote runner, it might be set to true to disable manifest update
		err := os.Setenv(k8config.EnvVarSkipManifestUpdate, "false")
		if err != nil {
			return err
		}
		k8Env.Cfg.SkipManifestUpdate = false
		lggr.Info().
			Str("Chart Name", chartName).
			Interface("Upgrade Details", props).
			Msg("Upgrading Chainlink Node")
		k8Env, err = k8Env.UpdateHelm(chartName, props)
		if err != nil {
			return err
		}
		err = k8Env.RunUpdated(noOfNodesToUpgrade)
		// Run the new environment and wait for changes to show
		if err != nil {
			return err
		}
	}
	return nil
}

// DeployEnvironments deploys K8 env for CCIP tests. For tests running on simulated geth it deploys -
// 1. two simulated geth network in non-dev mode
// 2. mockserver ( to set mock price feed details)
// 3. chainlink nodes
func DeployEnvironments(
	t *testing.T,
	envconfig *environment.Config,
	testInputs *CCIPTestConfig,
) *environment.Environment {
	selectedNetworks := testInputs.SelectedNetworks
	testEnvironment := environment.New(envconfig)
	numOfTxNodes := 1
	var charts []string
	for i, network := range selectedNetworks {
		if testInputs.EnvInput.Network.AnvilConfigs != nil {
			// if anvilconfig is specified for a network addhelm for anvil
			if anvilConfig, exists := testInputs.EnvInput.Network.AnvilConfigs[strings.ToUpper(network.Name)]; exists {
				charts = append(charts, foundry.ChartName)
				if anvilConfig.BaseFee == nil {
					anvilConfig.BaseFee = pointer.ToInt64(1000000)
				}
				if anvilConfig.BlockGaslimit == nil {
					anvilConfig.BlockGaslimit = pointer.ToInt64(100000000)
				}
				testEnvironment.
					AddHelm(foundry.NewVersioned("0.2.1", &foundry.Props{
						NetworkName: network.Name,
						Values: map[string]interface{}{
							"fullnameOverride": actions.NetworkName(network.Name),
							"image": map[string]interface{}{
								"repository": "ghcr.io/foundry-rs/foundry",
								"tag":        "v0.3.0",
							},
							"anvil": map[string]interface{}{
								"chainId":                   fmt.Sprintf("%d", network.ChainID),
								"blockTime":                 anvilConfig.BlockTime,
								"forkURL":                   anvilConfig.URL,
								"forkBlockNumber":           anvilConfig.BlockNumber,
								"forkRetries":               anvilConfig.Retries,
								"forkTimeout":               anvilConfig.Timeout,
								"forkComputeUnitsPerSecond": anvilConfig.ComputePerSecond,
								"forkNoRateLimit":           anvilConfig.RateLimitDisabled,
								"blocksToKeepInMemory":      anvilConfig.BlocksToKeepInMem,
								"blockGasLimit":             fmt.Sprintf("%d", pointer.GetInt64(anvilConfig.BlockGaslimit)),
								"baseFee":                   fmt.Sprintf("%d", pointer.GetInt64(anvilConfig.BaseFee)),
							},
							"resources": GethResourceProfile,
							"cache": map[string]interface{}{
								"capacity": "150Gi",
							},
						},
					}))
				selectedNetworks[i].Simulated = true
				actions.NetworkChart = foundry.ChartName
				continue
			}
		}

		if !network.Simulated {
			charts = append(charts, "")
			continue
		}
		charts = append(charts, strings.ReplaceAll(strings.ToLower(network.Name), " ", "-"))
		testEnvironment.
			AddHelm(reorg.New(&reorg.Props{
				NetworkName: network.Name,
				NetworkType: "simulated-geth-non-dev",
				Values: map[string]interface{}{
					"geth": map[string]interface{}{
						"genesis": map[string]interface{}{
							"networkId": fmt.Sprint(network.ChainID),
						},
						"tx": map[string]interface{}{
							"replicas":  strconv.Itoa(numOfTxNodes),
							"resources": testInputs.GethResourceProfile,
						},
						"miner": map[string]interface{}{
							"replicas":  "0",
							"resources": testInputs.GethResourceProfile,
						},
					},
					"bootnode": map[string]interface{}{
						"replicas": "1",
					},
				},
			}))
	}
	if pointer.GetBool(testInputs.TestGroupInput.USDCMockDeployment) ||
		pointer.GetBool(testInputs.TestGroupInput.TokenConfig.WithPipeline) {
		testEnvironment.AddHelm(parrot.New(nil))
	}
	err := testEnvironment.Run()
	require.NoError(t, err)

	if testEnvironment.WillUseRemoteRunner() {
		return testEnvironment
	}
	urlFinder := func(network blockchain.EVMNetwork, chart string) ([]string, []string) {
		if !network.Simulated {
			return network.URLs, network.HTTPURLs
		}
		networkName := actions.NetworkName(network.Name)
		var internalWsURLs, internalHttpURLs []string
		switch chart {
		case foundry.ChartName:
			internalWsURLs = append(internalWsURLs, fmt.Sprintf("ws://%s:8545", networkName))
			internalHttpURLs = append(internalHttpURLs, fmt.Sprintf("http://%s:8545", networkName))
		case networkName:
			for i := 0; i < numOfTxNodes; i++ {
				internalWsURLs = append(internalWsURLs, fmt.Sprintf("ws://%s-ethereum-geth:8546", networkName))
				internalHttpURLs = append(internalHttpURLs, fmt.Sprintf("http://%s-ethereum-geth:8544", networkName))
			}
		default:
			return network.URLs, network.HTTPURLs
		}

		return internalWsURLs, internalHttpURLs
	}
	var nets []blockchain.EVMNetwork
	for i := range selectedNetworks {
		nets = append(nets, selectedNetworks[i])
		nets[i].URLs, nets[i].HTTPURLs = urlFinder(selectedNetworks[i], charts[i])
	}

	err = testEnvironment.
		AddHelm(ChainlinkChart(t, testInputs, nets)).
		Run()
	require.NoError(t, err)
	return testEnvironment
}
