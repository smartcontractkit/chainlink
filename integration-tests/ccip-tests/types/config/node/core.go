package node

import (
	"bytes"
	"fmt"
	"math"
	"math/big"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	evmcfg "github.com/smartcontractkit/chainlink-integrations/evm/config/toml"
	ubig "github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	itutils "github.com/smartcontractkit/chainlink/integration-tests/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func NewConfigFromToml(tomlConfig []byte, opts ...node.NodeConfigOpt) (*chainlink.Config, error) {
	var cfg chainlink.Config
	err := config.DecodeTOML(bytes.NewReader(tomlConfig), &cfg)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &cfg, nil
}

func WithPrivateEVMs(networks []blockchain.EVMNetwork, commonChainConfig *evmcfg.Chain, chainSpecificConfig map[int64]evmcfg.Chain) node.NodeConfigOpt {
	var evmConfigs []*evmcfg.EVMConfig
	for _, network := range networks {
		var evmNodes []*evmcfg.Node
		for i := range network.URLs {
			evmNodes = append(evmNodes, &evmcfg.Node{
				Name:    ptr.Ptr(fmt.Sprintf("%s-%d", network.Name, i)),
				WSURL:   itutils.MustURL(network.URLs[i]),
				HTTPURL: itutils.MustURL(network.HTTPURLs[i]),
			})
		}
		evmConfig := &evmcfg.EVMConfig{
			ChainID: ubig.New(big.NewInt(network.ChainID)),
			Nodes:   evmNodes,
			Chain:   evmcfg.Chain{},
		}
		if commonChainConfig != nil {
			evmConfig.Chain = *commonChainConfig
		}
		if chainSpecificConfig == nil {
			if overriddenChainCfg, ok := chainSpecificConfig[network.ChainID]; ok {
				evmConfig.Chain = overriddenChainCfg
			}
		}
		if evmConfig.Chain.FinalityDepth == nil && network.FinalityDepth > 0 {
			if network.FinalityDepth > math.MaxUint32 {
				panic(fmt.Errorf("finality depth overflows uint32: %d", network.FinalityDepth))
			}
			evmConfig.Chain.FinalityDepth = ptr.Ptr(uint32(network.FinalityDepth))
		}
		if evmConfig.Chain.FinalityTagEnabled == nil && network.FinalityTag {
			evmConfig.Chain.FinalityTagEnabled = ptr.Ptr(network.FinalityTag)
		}
		evmConfigs = append(evmConfigs, evmConfig)
	}
	return func(c *chainlink.Config) {
		c.EVM = evmConfigs
	}
}
