package node

import (
	"bytes"
	_ "embed"
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

func WithPrivateEVMs(networks []blockchain.EVMNetwork) node.NodeConfigOpt {
	var evmConfigs []*evmcfg.EVMConfig
	for _, network := range networks {
		evmConfigs = append(evmConfigs, &evmcfg.EVMConfig{
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
			Nodes: []*evmcfg.Node{
				{
					Name:     itutils.Ptr(network.Name),
					WSURL:    itutils.MustURL(network.URLs[0]),
					HTTPURL:  itutils.MustURL(network.HTTPURLs[0]),
					SendOnly: itutils.Ptr(false),
				},
			},
		})
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
	case 420:
		cfg.BumpThreshold = itutils.Ptr(uint32(60))
		cfg.BumpPercent = itutils.Ptr(uint16(20))
		cfg.BumpMin = assets.GWei(100)
	case 5:
		cfg.PriceMax = assets.GWei(500)
		cfg.FeeCapDefault = assets.GWei(500)
	}

	return cfg
}
