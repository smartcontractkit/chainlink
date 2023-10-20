package node

import (
	"bytes"
	_ "embed"
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	itutils "github.com/smartcontractkit/ccip/integration-tests/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/config"
)

//go:embed tomls/ccip.toml
var CCIPTOML []byte

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

func WithDBConnectionPool(open, idle int64) node.NodeConfigOpt {
	return func(c *chainlink.Config) {
		c.Database.MaxOpenConns = itutils.Ptr(open)
		c.Database.MaxIdleConns = itutils.Ptr(idle)
	}
}

func WithPrivateEVMs(networks []blockchain.EVMNetwork) node.NodeConfigOpt {
	var evmConfigs []*evmcfg.EVMConfig
	for _, network := range networks {
		var evmNodes []*evmcfg.Node
		for i := range network.URLs {
			evmNodes = append(evmNodes, &evmcfg.Node{
				Name:    itutils.Ptr(fmt.Sprintf("%s-%d", network.Name, i)),
				WSURL:   itutils.MustURL(network.URLs[i]),
				HTTPURL: itutils.MustURL(network.HTTPURLs[i]),
			})
		}
		evmConfig := &evmcfg.EVMConfig{
			ChainID: utils.NewBig(big.NewInt(network.ChainID)),
			Chain: evmcfg.Chain{
				AutoCreateKey:      itutils.Ptr(true),
				FinalityDepth:      itutils.Ptr[uint32](50),
				MinContractPayment: assets.NewLinkFromJuels(0),
				LogPollInterval:    models.MustNewDuration(1 * time.Second),
				HeadTracker: evmcfg.HeadTracker{
					HistoryDepth: itutils.Ptr(uint32(100)),
				},
				GasEstimator: WithCCIPGasEstimator(network.ChainID),
			},
			Nodes: evmNodes,
		}

		evmConfigs = append(evmConfigs, evmConfig)
	}
	return func(c *chainlink.Config) {
		c.EVM = evmConfigs
	}
}

func WithCCIPGasEstimator(chainId int64) evmcfg.GasEstimator {
	cfg := evmcfg.GasEstimator{
		LimitDefault:  itutils.Ptr(uint32(6000000)),
		PriceMax:      assets.GWei(200),
		FeeCapDefault: assets.GWei(200),
	}
	switch chainId {
	case 421613:
		cfg.LimitDefault = itutils.Ptr(uint32(100000000))
		cfg.BumpThreshold = itutils.Ptr(uint32(60))
		cfg.BumpPercent = itutils.Ptr(uint16(20))
		cfg.BumpMin = assets.GWei(100)
		cfg.PriceMax = assets.GWei(400)
	case 420:
		cfg.BumpThreshold = itutils.Ptr(uint32(60))
		cfg.BumpPercent = itutils.Ptr(uint16(20))
		cfg.BumpMin = assets.GWei(100)
		cfg.PriceMax = assets.GWei(150)
		cfg.FeeCapDefault = assets.GWei(150)
		cfg.BlockHistory.BlockHistorySize = itutils.Ptr(uint16(200))
		cfg.BlockHistory.EIP1559FeeCapBufferBlocks = itutils.Ptr(uint16(0))
	case 84531:
		cfg.BumpThreshold = itutils.Ptr(uint32(60))
		cfg.BumpPercent = itutils.Ptr(uint16(20))
		cfg.BumpMin = assets.GWei(100)
		cfg.PriceMax = assets.GWei(150)
		cfg.FeeCapDefault = assets.GWei(150)
		cfg.BlockHistory.BlockHistorySize = itutils.Ptr(uint16(200))
		cfg.BlockHistory.EIP1559FeeCapBufferBlocks = itutils.Ptr(uint16(0))
	case 43113:
		cfg.BumpThreshold = itutils.Ptr(uint32(60))
	case 11155111:
		cfg.BlockHistory.BlockHistorySize = itutils.Ptr(uint16(200))
		cfg.BlockHistory.EIP1559FeeCapBufferBlocks = itutils.Ptr(uint16(0))
	}

	return cfg
}
