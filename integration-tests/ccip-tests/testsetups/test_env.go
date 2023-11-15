package testsetups

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/client"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/cdk8s/blockscout"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/k8s/pkg/helm/reorg"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/types/config/node"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	integrationnodes "github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/utils/config"
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

func setNodeConfig(nets []blockchain.EVMNetwork, nodeConfig, commonChain string, configByChain map[string]string) (*corechainlink.Config, string, error) {
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
	if configByChain != nil {
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

func ChainlinkChart(
	t *testing.T,
	testInputs *CCIPTestConfig,
	nets []blockchain.EVMNetwork,
) environment.ConnectedChart {
	require.NotNil(t, testInputs.EnvInput.Chainlink.Common, "Chainlink Common config is not specified")
	clProps := make(map[string]interface{})
	clProps["prometheus"] = true
	var formattedArgs []string
	if len(testInputs.EnvInput.Chainlink.DBArgs) > 0 {
		for _, arg := range testInputs.EnvInput.Chainlink.DBArgs {
			formattedArgs = append(formattedArgs, "-c")
			formattedArgs = append(formattedArgs, arg)
		}
	}
	clProps["db"] = map[string]interface{}{
		"resources":      SetResourceProfile(testInputs.EnvInput.Chainlink.DBCPU, testInputs.EnvInput.Chainlink.DBMemory),
		"additionalArgs": formattedArgs,
		"stateful":       pointer.GetBool(testInputs.EnvInput.Chainlink.IsStateful),
		"capacity":       testInputs.EnvInput.Chainlink.DBCapacity,
		"image": map[string]any{
			"image":   testInputs.EnvInput.Chainlink.Common.DBImage,
			"version": testInputs.EnvInput.Chainlink.Common.DBTag,
		},
	}
	clProps["chainlink"] = map[string]interface{}{
		"resources": SetResourceProfile(testInputs.EnvInput.Chainlink.NodeCPU, testInputs.EnvInput.Chainlink.NodeMemory),
		"image": map[string]any{
			"image":   testInputs.EnvInput.Chainlink.Common.Image,
			"version": testInputs.EnvInput.Chainlink.Common.Tag,
		},
	}

	require.NotNil(t, testInputs.EnvInput, "no env test input specified")

	if len(testInputs.EnvInput.Chainlink.Nodes) > 0 {
		var nodesMap []map[string]any
		for _, clNode := range testInputs.EnvInput.Chainlink.Nodes {
			nodeConfig := clNode.BaseConfigTOML
			commonChainConfig := clNode.CommonChainConfigTOML
			chainConfigByChain := clNode.ChainConfigTOMLByChain
			if nodeConfig == "" {
				nodeConfig = testInputs.EnvInput.Chainlink.Common.BaseConfigTOML
			}
			if commonChainConfig == "" {
				commonChainConfig = testInputs.EnvInput.Chainlink.Common.CommonChainConfigTOML
			}
			if chainConfigByChain == nil {
				chainConfigByChain = testInputs.EnvInput.Chainlink.Common.ChainConfigTOMLByChain
			}

			_, tomlStr, err := setNodeConfig(nets, nodeConfig, commonChainConfig, chainConfigByChain)
			require.NoError(t, err)
			nodesMap = append(nodesMap, map[string]any{
				"name": clNode.Name,
				"chainlink": map[string]any{
					"image": map[string]any{
						"image":   clNode.Image,
						"version": clNode.Tag,
					},
				},
				"db": map[string]any{
					"image": map[string]any{
						"image":   clNode.DBImage,
						"version": clNode.DBTag,
					},
				},
				"toml": tomlStr,
			})
		}
		clProps["nodes"] = nodesMap
		return chainlink.New(0, clProps)
	}
	clProps["replicas"] = pointer.GetInt(testInputs.EnvInput.Chainlink.NoOfNodes)
	_, tomlStr, err := setNodeConfig(
		nets,
		testInputs.EnvInput.Chainlink.Common.BaseConfigTOML,
		testInputs.EnvInput.Chainlink.Common.CommonChainConfigTOML,
		testInputs.EnvInput.Chainlink.Common.ChainConfigTOMLByChain,
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
	env, err := test_env.NewCLTestEnvBuilder().
		WithTestLogger(t).
		WithPrivateGethChains(selectedNetworks).
		WithoutCleanup().
		Build()
	require.NoError(t, err)
	for _, n := range env.PrivateChain {
		primaryNode := n.GetPrimaryNode()
		require.NotNil(t, primaryNode, "Primary node is nil in PrivateChain interface")
		for i, networkCfg := range selectedNetworks {
			if networkCfg.ChainID == n.GetNetworkConfig().ChainID {
				selectedNetworks[i].URLs = []string{primaryNode.GetInternalWsUrl()}
				selectedNetworks[i].HTTPURLs = []string{primaryNode.GetInternalHttpUrl()}
			}
		}
	}

	// a func to start the CL nodes asynchronously
	deployCL := func() error {
		toml, _, err := setNodeConfig(
			selectedNetworks,
			testInputs.EnvInput.Chainlink.Common.BaseConfigTOML,
			testInputs.EnvInput.Chainlink.Common.CommonChainConfigTOML,
			testInputs.EnvInput.Chainlink.Common.ChainConfigTOMLByChain,
		)
		if err != nil {
			return err
		}

		noOfNodes := pointer.GetInt(testInputs.EnvInput.Chainlink.NoOfNodes)
		return env.StartClCluster(toml, noOfNodes, "")
	}
	return env, deployCL
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
	useBlockscout := testInputs.TestGroupInput.Blockscout
	selectedNetworks := testInputs.SelectedNetworks
	testEnvironment := environment.New(envconfig)
	numOfTxNodes := 1
	for _, network := range selectedNetworks {
		if !network.Simulated {
			continue
		}
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
	err := testEnvironment.Run()
	require.NoError(t, err)

	if testEnvironment.WillUseRemoteRunner() {
		return testEnvironment
	}
	urlFinder := func(network blockchain.EVMNetwork) ([]string, []string) {
		if !network.Simulated {
			return network.URLs, network.HTTPURLs
		}
		networkName := network.Name
		var internalWsURLs, internalHttpURLs []string
		for i := 0; i < numOfTxNodes; i++ {
			podName := fmt.Sprintf("%s-ethereum-geth:%d", networkName, i)
			txNodeInternalWs, err := testEnvironment.Fwd.FindPort(podName, "geth", "ws-rpc").As(client.RemoteConnection, client.WS)
			require.NoError(t, err, "Error finding WS ports")
			internalWsURLs = append(internalWsURLs, txNodeInternalWs)
			txNodeInternalHttp, err := testEnvironment.Fwd.FindPort(podName, "geth", "http-rpc").As(client.RemoteConnection, client.HTTP)
			require.NoError(t, err, "Error finding HTTP ports")
			internalHttpURLs = append(internalHttpURLs, txNodeInternalHttp)
		}
		return internalWsURLs, internalHttpURLs
	}
	var nets []blockchain.EVMNetwork
	for i := range selectedNetworks {
		nets = append(nets, selectedNetworks[i])
		nets[i].URLs, nets[i].HTTPURLs = urlFinder(selectedNetworks[i])
		if useBlockscout {
			testEnvironment.AddChart(blockscout.New(&blockscout.Props{
				Name:    fmt.Sprintf("%s-blockscout", selectedNetworks[i].Name),
				WsURL:   selectedNetworks[i].URLs[0],
				HttpURL: selectedNetworks[i].HTTPURLs[0],
			}))
		}
	}

	err = testEnvironment.
		AddHelm(ChainlinkChart(t, testInputs, nets)).
		Run()
	require.NoError(t, err)
	return testEnvironment
}
