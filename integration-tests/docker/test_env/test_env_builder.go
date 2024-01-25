package test_env

import (
	"fmt"
	"math/big"
	"os"
	"runtime/debug"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/osutil"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

type CleanUpType string

const (
	CleanUpTypeNone     CleanUpType = "none"
	CleanUpTypeStandard CleanUpType = "standard"
	CleanUpTypeCustom   CleanUpType = "custom"
)

type CLTestEnvBuilder struct {
	hasLogStream           bool
	hasKillgrave           bool
	hasForwarders          bool
	clNodeConfig           *chainlink.Config
	secretsConfig          string
	nonDevGethNetworks     []blockchain.EVMNetwork
	clNodesCount           int
	clNodesOpts            []func(*ClNode)
	customNodeCsaKeys      []string
	defaultNodeCsaKeys     []string
	l                      zerolog.Logger
	t                      *testing.T
	te                     *CLClusterTestEnv
	isNonEVM               bool
	cleanUpType            CleanUpType
	cleanUpCustomFn        func()
	chainOptionsFn         []ChainOption
	evmClientNetworkOption []EVMClientNetworkOption
	ethereumNetwork        *test_env.EthereumNetwork
	testConfig             tc.GlobalTestConfig

	/* funding */
	ETHFunds *big.Float
}

func NewCLTestEnvBuilder() *CLTestEnvBuilder {
	return &CLTestEnvBuilder{
		l:            log.Logger,
		hasLogStream: true,
	}
}

// WithTestEnv sets the test environment to use for the test.
// If nil, a new test environment is created.
// If not nil, the test environment is used as-is.
// If TEST_ENV_CONFIG_PATH is set, the test environment is created with the config at that path.
func (b *CLTestEnvBuilder) WithTestEnv(te *CLClusterTestEnv) (*CLTestEnvBuilder, error) {
	envConfigPath, isSet := os.LookupEnv("TEST_ENV_CONFIG_PATH")
	var cfg *TestEnvConfig
	var err error
	if isSet {
		cfg, err = NewTestEnvConfigFromFile(envConfigPath)
		if err != nil {
			return nil, err
		}
	}

	if te != nil {
		b.te = te
	} else {
		b.te, err = NewTestEnv()
		if err != nil {
			return nil, err
		}
	}

	if cfg != nil {
		b.te = b.te.WithTestEnvConfig(cfg)
	}
	return b, nil
}

// WithTestLogger sets the test logger to use for the test.
// Useful for parallel tests so the logging will be separated correctly in the results views.
func (b *CLTestEnvBuilder) WithTestInstance(t *testing.T) *CLTestEnvBuilder {
	b.t = t
	b.l = logging.GetTestLogger(t)
	return b
}

// WithoutLogStream disables LogStream logging component
func (b *CLTestEnvBuilder) WithoutLogStream() *CLTestEnvBuilder {
	b.hasLogStream = false
	return b
}

func (b *CLTestEnvBuilder) WithCLNodes(clNodesCount int) *CLTestEnvBuilder {
	b.clNodesCount = clNodesCount
	return b
}

func (b *CLTestEnvBuilder) WithTestConfig(cfg tc.GlobalTestConfig) *CLTestEnvBuilder {
	b.testConfig = cfg
	return b
}

func (b *CLTestEnvBuilder) WithCLNodeOptions(opt ...ClNodeOption) *CLTestEnvBuilder {
	b.clNodesOpts = append(b.clNodesOpts, opt...)
	return b
}

func (b *CLTestEnvBuilder) WithForwarders() *CLTestEnvBuilder {
	b.hasForwarders = true
	return b
}

func (b *CLTestEnvBuilder) WithFunding(eth *big.Float) *CLTestEnvBuilder {
	b.ETHFunds = eth
	return b
}

// deprecated
// left only for backward compatibility
func (b *CLTestEnvBuilder) WithGeth() *CLTestEnvBuilder {
	ethBuilder := test_env.NewEthereumNetworkBuilder()
	cfg, err := ethBuilder.
		WithConsensusType(test_env.ConsensusType_PoW).
		WithExecutionLayer(test_env.ExecutionLayer_Geth).
		WithTest(b.t).
		Build()

	if err != nil {
		panic(err)
	}

	b.ethereumNetwork = &cfg

	return b
}

func (b *CLTestEnvBuilder) WithPrivateEthereumNetwork(en test_env.EthereumNetwork) *CLTestEnvBuilder {
	b.ethereumNetwork = &en
	return b
}

func (b *CLTestEnvBuilder) WithPrivateGethChains(evmNetworks []blockchain.EVMNetwork) *CLTestEnvBuilder {
	b.nonDevGethNetworks = evmNetworks
	return b
}

func (b *CLTestEnvBuilder) WithCLNodeConfig(cfg *chainlink.Config) *CLTestEnvBuilder {
	b.clNodeConfig = cfg
	return b
}

func (b *CLTestEnvBuilder) WithSecretsConfig(secrets string) *CLTestEnvBuilder {
	b.secretsConfig = secrets
	return b
}

func (b *CLTestEnvBuilder) WithMockAdapter() *CLTestEnvBuilder {
	b.hasKillgrave = true
	return b
}

// WithNonEVM sets the test environment to not use EVM when built.
func (b *CLTestEnvBuilder) WithNonEVM() *CLTestEnvBuilder {
	b.isNonEVM = true
	return b
}

func (b *CLTestEnvBuilder) WithStandardCleanup() *CLTestEnvBuilder {
	b.cleanUpType = CleanUpTypeStandard
	return b
}

func (b *CLTestEnvBuilder) WithoutCleanup() *CLTestEnvBuilder {
	b.cleanUpType = CleanUpTypeNone
	return b
}

func (b *CLTestEnvBuilder) WithCustomCleanup(customFn func()) *CLTestEnvBuilder {
	b.cleanUpType = CleanUpTypeCustom
	b.cleanUpCustomFn = customFn
	return b
}

type ChainOption = func(*evmcfg.Chain) *evmcfg.Chain

func (b *CLTestEnvBuilder) WithChainOptions(opts ...ChainOption) *CLTestEnvBuilder {
	b.chainOptionsFn = make([]ChainOption, 0)
	b.chainOptionsFn = append(b.chainOptionsFn, opts...)

	return b
}

type EVMClientNetworkOption = func(*blockchain.EVMNetwork) *blockchain.EVMNetwork

func (b *CLTestEnvBuilder) EVMClientNetworkOptions(opts ...EVMClientNetworkOption) *CLTestEnvBuilder {
	b.evmClientNetworkOption = make([]EVMClientNetworkOption, 0)
	b.evmClientNetworkOption = append(b.evmClientNetworkOption, opts...)

	return b
}

func (b *CLTestEnvBuilder) Build() (*CLClusterTestEnv, error) {
	if b.testConfig == nil {
		return nil, fmt.Errorf("test config must be set")
	}

	if b.te == nil {
		var err error
		b, err = b.WithTestEnv(nil)
		if err != nil {
			return nil, err
		}
	}

	var err error
	if b.t != nil {
		b.te.WithTestInstance(b.t)
	}

	if b.hasLogStream {
		b.te.LogStream, err = logstream.NewLogStream(b.te.t, b.testConfig.GetLoggingConfig())
		if err != nil {
			return nil, err
		}
	}

	if b.hasKillgrave {
		if b.te.Network == nil {
			return nil, fmt.Errorf("test environment builder failed: %w", fmt.Errorf("cannot start mock adapter without a network"))
		}

		b.te.MockAdapter = test_env.NewKillgrave([]string{b.te.Network.Name}, "", test_env.WithLogStream(b.te.LogStream))

		err = b.te.StartMockAdapter()
		if err != nil {
			return nil, err
		}
	}

	if b.t != nil {
		b.te.WithTestInstance(b.t)
	}

	switch b.cleanUpType {
	case CleanUpTypeStandard:
		b.t.Cleanup(func() {
			if err := b.te.Cleanup(); err != nil {
				b.l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		})
	case CleanUpTypeCustom:
		b.t.Cleanup(b.cleanUpCustomFn)
	case CleanUpTypeNone:
		b.l.Warn().Msg("test environment won't be cleaned up")
	case "":
		return b.te, fmt.Errorf("test environment builder failed: %w", fmt.Errorf("explicit cleanup type must be set when building test environment"))
	}

	if b.te.LogStream != nil {
		b.t.Cleanup(func() {
			b.l.Info().Msg("Shutting down LogStream")
			logPath, err := osutil.GetAbsoluteFolderPath("logs")
			if err != nil {
				b.l.Info().Str("Absolute path", logPath).Msg("LogStream logs folder location")
			}

			if b.t.Failed() || *b.testConfig.GetLoggingConfig().TestLogCollect {
				// we can't do much if this fails, so we just log the error in logstream
				_ = b.te.LogStream.FlushAndShutdown()
				b.te.LogStream.PrintLogTargetsLocations()
				b.te.LogStream.SaveLogLocationInTestSummary()
			}

		})
	}

	if b.nonDevGethNetworks != nil {
		b.te.WithPrivateChain(b.nonDevGethNetworks)
		err := b.te.StartPrivateChain()
		if err != nil {
			return b.te, err
		}
		var nonDevNetworks []blockchain.EVMNetwork
		for i, n := range b.te.PrivateChain {
			primaryNode := n.GetPrimaryNode()
			if primaryNode == nil {
				return b.te, fmt.Errorf("primary node is nil in PrivateChain interface, stack: %s", string(debug.Stack()))
			}
			nonDevNetworks = append(nonDevNetworks, *n.GetNetworkConfig())
			nonDevNetworks[i].URLs = []string{primaryNode.GetInternalWsUrl()}
			nonDevNetworks[i].HTTPURLs = []string{primaryNode.GetInternalHttpUrl()}
		}
		if nonDevNetworks == nil {
			return nil, fmt.Errorf("cannot create nodes with custom config without nonDevNetworks")
		}

		err = b.te.StartClCluster(b.clNodeConfig, b.clNodesCount, b.secretsConfig, b.testConfig, b.clNodesOpts...)
		if err != nil {
			return nil, err
		}
		return b.te, nil
	}

	networkConfig := networks.MustGetSelectedNetworkConfig(b.testConfig.GetNetworkConfig())[0]
	var rpcProvider test_env.RpcProvider
	if b.ethereumNetwork != nil && networkConfig.Simulated {
		// TODO here we should save the ethereum network config to te.Cfg, but it doesn't exist at this point
		// in general it seems we have no methods for saving config to file and we only load it from file
		// but I don't know how that config file is to be created or whether anyone ever done that
		b.ethereumNetwork.DockerNetworkNames = []string{b.te.Network.Name}
		networkConfig, rpcProvider, err = b.te.StartEthereumNetwork(b.ethereumNetwork)
		if err != nil {
			return nil, err
		}
		b.te.RpcProvider = rpcProvider
		b.te.PrivateEthereumConfig = b.ethereumNetwork
	}

	if !b.isNonEVM {
		if b.evmClientNetworkOption != nil && len(b.evmClientNetworkOption) > 0 {
			for _, fn := range b.evmClientNetworkOption {
				fn(&networkConfig)
			}
		}
		bc, err := blockchain.NewEVMClientFromNetwork(networkConfig, b.l)
		if err != nil {
			return nil, err
		}

		b.te.EVMClient = bc
		cd, err := contracts.NewContractDeployer(bc, b.l)
		if err != nil {
			return nil, err
		}
		b.te.ContractDeployer = cd

		cl, err := contracts.NewContractLoader(bc, b.l)
		if err != nil {
			return nil, err
		}
		b.te.ContractLoader = cl
	}

	var nodeCsaKeys []string

	// Start Chainlink Nodes
	if b.clNodesCount > 0 {
		var cfg *chainlink.Config
		if b.clNodeConfig != nil {
			cfg = b.clNodeConfig
		} else {
			cfg = node.NewConfig(node.NewBaseConfig(),
				node.WithOCR1(),
				node.WithP2Pv2(),
			)
		}

		if !b.isNonEVM {
			var httpUrls []string
			var wsUrls []string
			if networkConfig.Simulated {
				httpUrls = rpcProvider.PrivateHttpUrls()
				wsUrls = rpcProvider.PrivateWsUrsl()
			} else {
				httpUrls = networkConfig.HTTPURLs
				wsUrls = networkConfig.URLs
			}

			node.SetChainConfig(cfg, wsUrls, httpUrls, networkConfig, b.hasForwarders)

			if b.chainOptionsFn != nil && len(b.chainOptionsFn) > 0 {
				for _, fn := range b.chainOptionsFn {
					for _, evmCfg := range cfg.EVM {
						fn(&evmCfg.Chain)
					}
				}
			}
		}

		err := b.te.StartClCluster(cfg, b.clNodesCount, b.secretsConfig, b.testConfig, b.clNodesOpts...)
		if err != nil {
			return nil, err
		}

		nodeCsaKeys, err = b.te.ClCluster.NodeCSAKeys()
		if err != nil {
			return nil, err
		}
		b.defaultNodeCsaKeys = nodeCsaKeys
	}

	if b.ethereumNetwork != nil && b.clNodesCount > 0 && b.ETHFunds != nil {
		b.te.ParallelTransactions(true)
		defer b.te.ParallelTransactions(false)
		if err := b.te.FundChainlinkNodes(b.ETHFunds); err != nil {
			return nil, err
		}
	}

	var enDesc string
	if b.te.PrivateEthereumConfig != nil {
		enDesc = b.te.PrivateEthereumConfig.Describe()
	} else {
		enDesc = "none"
	}

	b.l.Info().
		Str("privateEthereumNetwork", enDesc).
		Bool("hasKillgrave", b.hasKillgrave).
		Int("clNodesCount", b.clNodesCount).
		Strs("customNodeCsaKeys", b.customNodeCsaKeys).
		Strs("defaultNodeCsaKeys", b.defaultNodeCsaKeys).
		Msg("Building CL cluster test environment..")

	return b.te, nil
}
