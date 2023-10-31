package node

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/segmentio/ksuid"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/config"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions/vrfv2_constants"
	utils2 "github.com/smartcontractkit/chainlink/integration-tests/utils"
)

func NewBaseConfig() *chainlink.Config {
	return &chainlink.Config{
		Core: toml.Core{
			RootDir: utils2.Ptr("/home/chainlink"),
			Database: toml.Database{
				MaxIdleConns:     utils2.Ptr(int64(20)),
				MaxOpenConns:     utils2.Ptr(int64(40)),
				MigrateOnStartup: utils2.Ptr(true),
			},
			Log: toml.Log{
				Level:       utils2.Ptr(toml.LogLevel(zapcore.DebugLevel)),
				JSONConsole: utils2.Ptr(true),
				File: toml.LogFile{
					MaxSize: utils2.Ptr(utils.FileSize(0)),
				},
			},
			WebServer: toml.WebServer{
				AllowOrigins:   utils2.Ptr("*"),
				HTTPPort:       utils2.Ptr[uint16](6688),
				SecureCookies:  utils2.Ptr(false),
				SessionTimeout: models.MustNewDuration(time.Hour * 999),
				TLS: toml.WebServerTLS{
					HTTPSPort: utils2.Ptr[uint16](0),
				},
				RateLimit: toml.WebServerRateLimit{
					Authenticated:   utils2.Ptr(int64(2000)),
					Unauthenticated: utils2.Ptr(int64(100)),
				},
			},
			Feature: toml.Feature{
				LogPoller:    utils2.Ptr(true),
				FeedsManager: utils2.Ptr(true),
				UICSAKeys:    utils2.Ptr(true),
			},
			P2P: toml.P2P{},
		},
	}
}

type NodeConfigOpt = func(c *chainlink.Config)

func NewConfig(baseConf *chainlink.Config, opts ...NodeConfigOpt) *chainlink.Config {
	for _, opt := range opts {
		opt(baseConf)
	}
	return baseConf
}

func NewConfigFromToml(tomlFile string, opts ...NodeConfigOpt) (*chainlink.Config, error) {
	readFile, err := os.ReadFile(tomlFile)
	if err != nil {
		return nil, err
	}
	var cfg chainlink.Config
	if err != nil {
		return nil, err
	}
	err = config.DecodeTOML(bytes.NewReader(readFile), &cfg)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &cfg, nil
}

func WithOCR1() NodeConfigOpt {
	return func(c *chainlink.Config) {
		c.OCR = toml.OCR{
			Enabled: utils2.Ptr(true),
		}
	}
}

func WithOCR2() NodeConfigOpt {
	return func(c *chainlink.Config) {
		c.OCR2 = toml.OCR2{
			Enabled: utils2.Ptr(true),
		}
	}
}

func WithP2Pv1() NodeConfigOpt {
	return func(c *chainlink.Config) {
		c.P2P.V1 = toml.P2PV1{
			Enabled:    utils2.Ptr(true),
			ListenIP:   utils2.MustIP("0.0.0.0"),
			ListenPort: utils2.Ptr[uint16](6690),
		}
		// disabled default
		c.P2P.V2 = toml.P2PV2{Enabled: utils2.Ptr(false)}
	}
}

func WithP2Pv2() NodeConfigOpt {
	return func(c *chainlink.Config) {
		c.P2P.V2 = toml.P2PV2{
			ListenAddresses: &[]string{"0.0.0.0:6690"},
		}
	}
}

func WithTracing() NodeConfigOpt {
	return func(c *chainlink.Config) {
		c.Tracing = toml.Tracing{
			Enabled:         utils2.Ptr(true),
			CollectorTarget: utils2.Ptr("otel-collector:4317"),
			// ksortable unique id
			NodeID: utils2.Ptr(ksuid.New().String()),
			Attributes: map[string]string{
				"env": "smoke",
			},
			SamplingRatio: utils2.Ptr(1.0),
		}
	}
}

func SetChainConfig(
	cfg *chainlink.Config,
	wsUrls,
	httpUrls []string,
	chain blockchain.EVMNetwork,
	forwarders bool,
) {
	if cfg.EVM == nil {
		var nodes []*evmcfg.Node
		for i := range wsUrls {
			node := evmcfg.Node{
				Name:     utils2.Ptr(fmt.Sprintf("node_%d_%s", i, chain.Name)),
				WSURL:    utils2.MustURL(wsUrls[i]),
				HTTPURL:  utils2.MustURL(httpUrls[i]),
				SendOnly: utils2.Ptr(false),
			}

			nodes = append(nodes, &node)
		}
		var chainConfig evmcfg.Chain
		if chain.Simulated {
			chainConfig = evmcfg.Chain{
				AutoCreateKey:      utils2.Ptr(true),
				FinalityDepth:      utils2.Ptr[uint32](1),
				MinContractPayment: assets.NewLinkFromJuels(0),
			}
		}
		cfg.EVM = evmcfg.EVMConfigs{
			{
				ChainID: utils.NewBig(big.NewInt(chain.ChainID)),
				Chain:   chainConfig,
				Nodes:   nodes,
			},
		}
		if forwarders {
			cfg.EVM[0].Transactions = evmcfg.Transactions{
				ForwardersEnabled: utils2.Ptr(true),
			}
		}
	}
}

func WithPrivateEVMs(networks []blockchain.EVMNetwork) NodeConfigOpt {
	var evmConfigs []*evmcfg.EVMConfig
	for _, network := range networks {
		evmConfigs = append(evmConfigs, &evmcfg.EVMConfig{
			ChainID: utils.NewBig(big.NewInt(network.ChainID)),
			Chain: evmcfg.Chain{
				AutoCreateKey:      utils2.Ptr(true),
				FinalityDepth:      utils2.Ptr[uint32](50),
				MinContractPayment: assets.NewLinkFromJuels(0),
				LogPollInterval:    models.MustNewDuration(1 * time.Second),
				HeadTracker: evmcfg.HeadTracker{
					HistoryDepth: utils2.Ptr(uint32(100)),
				},
				GasEstimator: evmcfg.GasEstimator{
					LimitDefault:  utils2.Ptr(uint32(6000000)),
					PriceMax:      assets.GWei(200),
					FeeCapDefault: assets.GWei(200),
				},
			},
			Nodes: []*evmcfg.Node{
				{
					Name:     utils2.Ptr(network.Name),
					WSURL:    utils2.MustURL(network.URLs[0]),
					HTTPURL:  utils2.MustURL(network.HTTPURLs[0]),
					SendOnly: utils2.Ptr(false),
				},
			},
		})
	}
	return func(c *chainlink.Config) {
		c.EVM = evmConfigs
	}
}

func WithVRFv2EVMEstimator(addr string) NodeConfigOpt {
	est := assets.GWei(vrfv2_constants.MaxGasPriceGWei)
	return func(c *chainlink.Config) {
		c.EVM[0].KeySpecific = evmcfg.KeySpecificConfig{
			{
				Key: utils2.Ptr(ethkey.EIP55Address(addr)),
				GasEstimator: evmcfg.KeySpecificGasEstimator{
					PriceMax: est,
				},
			},
		}
		c.EVM[0].Chain.GasEstimator = evmcfg.GasEstimator{
			LimitDefault: utils2.Ptr[uint32](3500000),
		}
		c.EVM[0].Chain.Transactions = evmcfg.Transactions{
			MaxQueued: utils2.Ptr[uint32](10000),
		}
	}
}
