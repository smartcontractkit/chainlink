package node

import (
	"bytes"
	"embed"
	"fmt"
	"math/big"
	"net"
	"path/filepath"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/vrfv2_actions/vrfv2_constants"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils/config"

	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	BaseConf = &chainlink.Config{
		Core: toml.Core{
			RootDir: ptr("/home/chainlink"),
			Database: toml.Database{
				MaxIdleConns:     ptr(int64(20)),
				MaxOpenConns:     ptr(int64(40)),
				MigrateOnStartup: ptr(true),
			},
			Log: toml.Log{
				Level:       ptr(toml.LogLevel(zapcore.DebugLevel)),
				JSONConsole: ptr(true),
			},
			WebServer: toml.WebServer{
				AllowOrigins:   ptr("*"),
				HTTPPort:       ptr[uint16](6688),
				SecureCookies:  ptr(false),
				SessionTimeout: models.MustNewDuration(time.Hour * 999),
				TLS: toml.WebServerTLS{
					HTTPSPort: ptr[uint16](0),
				},
				RateLimit: toml.WebServerRateLimit{
					Authenticated:   ptr(int64(2000)),
					Unauthenticated: ptr(int64(100)),
				},
			},
			Feature: toml.Feature{
				LogPoller:    ptr(true),
				FeedsManager: ptr(true),
				UICSAKeys:    ptr(true),
			},
			P2P: toml.P2P{},
		},
	}
	//go:embed defaults/*.toml
	defaultsFS embed.FS
)

type NodeConfigOpt = func(c *chainlink.Config)

func NewConfig(baseConf *chainlink.Config, opts ...NodeConfigOpt) *chainlink.Config {
	for _, opt := range opts {
		opt(baseConf)
	}
	return baseConf
}

func NewConfigFromToml(tomlFile string, opts ...NodeConfigOpt) (*chainlink.Config, error) {
	path := filepath.Join("defaults", tomlFile+".toml")
	b, err := defaultsFS.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg chainlink.Config
	err = config.DecodeTOML(bytes.NewReader(b), &cfg)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &cfg, err
}

func WithOCR1() NodeConfigOpt {
	return func(c *chainlink.Config) {
		c.OCR = toml.OCR{
			Enabled: ptr(true),
		}
	}
}

func WithOCR2() NodeConfigOpt {
	return func(c *chainlink.Config) {
		c.OCR2 = toml.OCR2{
			Enabled: ptr(true),
		}
	}
}

func WithP2Pv1() NodeConfigOpt {
	return func(c *chainlink.Config) {
		c.P2P.V1 = toml.P2PV1{
			Enabled:    ptr(true),
			ListenIP:   mustIP("0.0.0.0"),
			ListenPort: ptr[uint16](6690),
		}
	}
}

func WithP2Pv2() NodeConfigOpt {
	return func(c *chainlink.Config) {
		c.P2P.V2 = toml.P2PV2{
			Enabled:         ptr(true),
			ListenAddresses: &[]string{"0.0.0.0:6690"},
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
		for i, _ := range wsUrls {
			node := evmcfg.Node{
				Name:     ptr(fmt.Sprintf("node_%d_%s", i, chain.Name)),
				WSURL:    mustURL(wsUrls[i]),
				HTTPURL:  mustURL(httpUrls[i]),
				SendOnly: ptr(false),
			}

			nodes = append(nodes, &node)
		}
		var chainConfig evmcfg.Chain
		if chain.Simulated {
			chainConfig = evmcfg.Chain{
				AutoCreateKey:      ptr(true),
				FinalityDepth:      ptr[uint32](1),
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
				ForwardersEnabled: ptr(true),
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
				AutoCreateKey:      ptr(true),
				FinalityDepth:      ptr[uint32](50),
				MinContractPayment: assets.NewLinkFromJuels(0),
				LogPollInterval:    models.MustNewDuration(1 * time.Second),
				HeadTracker: evmcfg.HeadTracker{
					HistoryDepth: ptr(uint32(100)),
				},
				GasEstimator: evmcfg.GasEstimator{
					LimitDefault:  ptr(uint32(6000000)),
					PriceMax:      assets.GWei(200),
					FeeCapDefault: assets.GWei(200),
				},
			},
			Nodes: []*evmcfg.Node{
				{
					Name:     ptr(network.Name),
					WSURL:    mustURL(network.URLs[0]),
					HTTPURL:  mustURL(network.HTTPURLs[0]),
					SendOnly: ptr(false),
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
				Key: ptr(ethkey.EIP55Address(addr)),
				GasEstimator: evmcfg.KeySpecificGasEstimator{
					PriceMax: est,
				},
			},
		}
		c.EVM[0].Chain.GasEstimator = evmcfg.GasEstimator{
			LimitDefault: ptr[uint32](3500000),
		}
		c.EVM[0].Chain.Transactions = evmcfg.Transactions{
			MaxQueued: ptr[uint32](10000),
		}
	}
}

func ptr[T any](t T) *T { return &t }

func mustURL(s string) *models.URL {
	var u models.URL
	if err := u.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return &u
}

func mustIP(s string) *net.IP {
	var ip net.IP
	if err := ip.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return &ip
}
