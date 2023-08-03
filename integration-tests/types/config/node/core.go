package node

import (
	"math/big"
	"net"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"go.uber.org/zap/zapcore"
)

var BaseConf = chainlink.Config{
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
		JobPipeline: toml.JobPipeline{
			MaxSuccessfulRuns: ptr[uint64](0),
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

type NodeConfigOpt = func(c *chainlink.Config)

func NewConfig(baseConf chainlink.Config, opts ...NodeConfigOpt) chainlink.Config {
	for _, opt := range opts {
		opt(&baseConf)
	}
	return baseConf
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

func WithSimulatedEVM(httpUrl, wsUrl string) NodeConfigOpt {
	return func(c *chainlink.Config) {
		c.EVM = evmcfg.EVMConfigs{
			{
				ChainID: utils.NewBig(big.NewInt(1337)),
				Chain: evmcfg.Chain{
					AutoCreateKey:      ptr(true),
					FinalityDepth:      ptr[uint32](1),
					MinContractPayment: assets.NewLinkFromJuels(0),
					GasEstimator: evmcfg.GasEstimator{
						LimitDefault: ptr[uint32](3500000),
					},
					Transactions: evmcfg.Transactions{
						MaxQueued: ptr[uint32](10000),
					},
				},
				Nodes: []*evmcfg.Node{
					{
						Name:     ptr("1337_primary_local_0"),
						WSURL:    mustURL(wsUrl),
						HTTPURL:  mustURL(httpUrl),
						SendOnly: ptr(false),
					},
				},
			},
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
