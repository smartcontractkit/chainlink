package persistent

import (
	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func BuildEVMOnlyChainlinkConfigs(donConfig *testconfig.ChainlinkDeployment, rpcProviders map[uint64]persistent_types.RpcProvider) ([]*chainlink.Config, error) {
	evmNetworks := prepareEVMNetworksEndpoints(rpcProviders)
	noOfNodes := getNumberOfNodes(donConfig)

	return buildNodeConfigs(evmNetworks, noOfNodes, donConfig.Common.BaseConfigTOML, donConfig.Common.CommonChainConfigTOML, donConfig.Common.ChainConfigTOMLByChain)
}

func FetchNodeIds(don *persistent_types.DON) ([]string, error) {
	var nodeIDs []string
	for _, node := range don.ChainlinkClients {
		p2pKeys, err := node.MustReadP2PKeys()
		if err != nil {
			return nil, err
		}

		nodeIDs = append(nodeIDs, p2pKeys.Data[0].Attributes.PeerID)
	}
	return nodeIDs, nil
}

func prepareEVMNetworksEndpoints(rpcProviders map[uint64]persistent_types.RpcProvider) []blockchain.EVMNetwork {
	var evmNetworks []blockchain.EVMNetwork
	for _, rpcProvider := range rpcProviders {
		evmNetwork := rpcProvider.EVMNetwork()
		evmNetwork.HTTPURLs = rpcProvider.PrivateHttpUrls()
		evmNetwork.URLs = rpcProvider.PrivateWsUrls()
		evmNetworks = append(evmNetworks, evmNetwork)
	}
	return evmNetworks
}

func getNumberOfNodes(donConfig *testconfig.ChainlinkDeployment) int {
	if donConfig.NoOfNodes != nil {
		return pointer.GetInt(donConfig.NoOfNodes)
	}
	return len(donConfig.Nodes)
}

func buildNodeConfigs(evmNetworks []blockchain.EVMNetwork, noOfNodes int, baseConfigTOML, commonChainConfigTOML string, chainConfigTOMLByChain map[string]string) ([]*chainlink.Config, error) {
	var clNodeConfigs []*chainlink.Config
	for i := 0; i < noOfNodes; i++ {
		toml, _, err := testsetups.SetNodeConfig(
			evmNetworks,
			baseConfigTOML,
			commonChainConfigTOML,
			chainConfigTOMLByChain,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create node config")
		}
		clNodeConfigs = append(clNodeConfigs, toml)
	}

	return clNodeConfigs, nil
}
