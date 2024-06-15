package node

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/segmentio/ksuid"
	"go.uber.org/zap/zapcore"

	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
	itutils "github.com/smartcontractkit/chainlink/integration-tests/utils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func NewBaseConfig() *chainlink.Config {
	return &chainlink.Config{
		Core: toml.Core{
			RootDir: ptr.Ptr("/home/chainlink"),
			Database: toml.Database{
				MaxIdleConns:     ptr.Ptr(int64(20)),
				MaxOpenConns:     ptr.Ptr(int64(40)),
				MigrateOnStartup: ptr.Ptr(true),
			},
			Log: toml.Log{
				Level:       ptr.Ptr(toml.LogLevel(zapcore.DebugLevel)),
				JSONConsole: ptr.Ptr(true),
				File: toml.LogFile{
					MaxSize: ptr.Ptr(utils.FileSize(0)),
				},
			},
			WebServer: toml.WebServer{
				AllowOrigins:   ptr.Ptr("*"),
				HTTPPort:       ptr.Ptr[uint16](6688),
				SecureCookies:  ptr.Ptr(false),
				SessionTimeout: commonconfig.MustNewDuration(time.Hour * 999),
				TLS: toml.WebServerTLS{
					HTTPSPort: ptr.Ptr[uint16](0),
				},
				RateLimit: toml.WebServerRateLimit{
					Authenticated:   ptr.Ptr(int64(2000)),
					Unauthenticated: ptr.Ptr(int64(100)),
				},
			},
			Feature: toml.Feature{
				LogPoller:    ptr.Ptr(true),
				FeedsManager: ptr.Ptr(true),
				UICSAKeys:    ptr.Ptr(true),
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

func NewConfigFromToml(tomlConfig []byte, opts ...NodeConfigOpt) (*chainlink.Config, error) {
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

func WithPrivateEVMs(networks []blockchain.EVMNetwork, commonChainConfig *evmcfg.Chain, chainSpecificConfig map[int64]evmcfg.Chain) NodeConfigOpt {
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
		if chainSpecificConfig != nil {
			if overriddenChainCfg, ok := chainSpecificConfig[network.ChainID]; ok {
				evmConfig.Chain = overriddenChainCfg
			}
		}
		evmConfigs = append(evmConfigs, evmConfig)
	}
	return func(c *chainlink.Config) {
		c.EVM = evmConfigs
	}
}

func WithKeySpecificMaxGasPrice(addresses []string, maxGasPriceGWei int64) NodeConfigOpt {
	est := assets.GWei(maxGasPriceGWei)
	var keySpecicifArr []evmcfg.KeySpecific
	for _, addr := range addresses {
		keySpecicifArr = append(keySpecicifArr, evmcfg.KeySpecific{
			Key: ptr.Ptr(types.EIP55Address(addr)),
			GasEstimator: evmcfg.KeySpecificGasEstimator{
				PriceMax: est,
			},
		})
	}
	return func(c *chainlink.Config) {
		c.EVM[0].KeySpecific = keySpecicifArr
	}
}

func WithLogPollInterval(interval time.Duration) NodeConfigOpt {
	return func(c *chainlink.Config) {
		c.EVM[0].Chain.LogPollInterval = commonconfig.MustNewDuration(interval)
	}
}

func BuildChainlinkNodeConfig(nets []blockchain.EVMNetwork, nodeConfig, commonChain string, configByChain map[string]string) (*corechainlink.Config, string, error) {
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
		tomlCfg = NewConfig(
			NewBaseConfig(),
			WithPrivateEVMs(nets, commonChainConfig, configByChainMap))
	} else {
		tomlCfg, err = NewConfigFromToml([]byte(nodeConfig), WithPrivateEVMs(nets, commonChainConfig, configByChainMap))
		if err != nil {
			return nil, "", err
		}
	}

	// we need unique id for each node for OTEL tracing
	if tomlCfg.Tracing.Enabled != nil && *tomlCfg.Tracing.Enabled {
		tomlCfg.Tracing.NodeID = ptr.Ptr(ksuid.New().String())
	}

	tomlStr, err := tomlCfg.TOMLString()
	return tomlCfg, tomlStr, err
}
