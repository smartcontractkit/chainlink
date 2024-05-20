package test_env

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"slices"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/seth"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/osutil"
	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
)

type CleanUpType string

const (
	CleanUpTypeNone     CleanUpType = "none"
	CleanUpTypeStandard CleanUpType = "standard"
	CleanUpTypeCustom   CleanUpType = "custom"
)

type ChainlinkNodeLogScannerSettings struct {
	FailingLogLevel zapcore.Level
	Threshold       uint
	AllowedMessages []testreporters.AllowedLogMessage
}

type CLTestEnvBuilder struct {
	hasLogStream                    bool
	hasKillgrave                    bool
	hasForwarders                   bool
	hasSeth                         bool
	hasEVMClient                    bool
	clNodeConfig                    *chainlink.Config
	secretsConfig                   string
	clNodesCount                    int
	clNodesOpts                     []func(*ClNode)
	customNodeCsaKeys               []string
	defaultNodeCsaKeys              []string
	l                               zerolog.Logger
	t                               *testing.T
	te                              *CLClusterTestEnv
	isNonEVM                        bool
	cleanUpType                     CleanUpType
	cleanUpCustomFn                 func()
	chainOptionsFn                  []ChainOption
	evmNetworkOption                []EVMNetworkOption
	privateEthereumNetworks         []*ctf_config.EthereumNetworkConfig
	testConfig                      ctf_config.GlobalTestConfig
	chainlinkNodeLogScannerSettings *ChainlinkNodeLogScannerSettings

	/* funding */
	ETHFunds *big.Float
}

var DefaultAllowedMessages = []testreporters.AllowedLogMessage{
	testreporters.NewAllowedLogMessage("Failed to get LINK balance", "Happens only when we deploy LINK token for test purposes. Harmless.", zapcore.ErrorLevel, testreporters.WarnAboutAllowedMsgs_No),
}

var DefaultChainlinkNodeLogScannerSettings = ChainlinkNodeLogScannerSettings{
	FailingLogLevel: zapcore.DPanicLevel,
	Threshold:       1, // we want to fail on the first concerning log
	AllowedMessages: DefaultAllowedMessages,
}

func GetDefaultChainlinkNodeLogScannerSettingsWithExtraAllowedMessages(extraAllowedMessages ...testreporters.AllowedLogMessage) ChainlinkNodeLogScannerSettings {
	allowedMessages := append(DefaultAllowedMessages, extraAllowedMessages...)
	return ChainlinkNodeLogScannerSettings{
		FailingLogLevel: zapcore.DPanicLevel,
		Threshold:       1,
		AllowedMessages: allowedMessages,
	}
}

func NewCLTestEnvBuilder() *CLTestEnvBuilder {
	return &CLTestEnvBuilder{
		l:                               log.Logger,
		hasLogStream:                    true,
		hasEVMClient:                    true,
		chainlinkNodeLogScannerSettings: &DefaultChainlinkNodeLogScannerSettings,
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

func (b *CLTestEnvBuilder) WithoutChainlinkNodeLogScanner() *CLTestEnvBuilder {
	b.chainlinkNodeLogScannerSettings = &ChainlinkNodeLogScannerSettings{}
	return b
}

func (b *CLTestEnvBuilder) WithChainlinkNodeLogScanner(settings ChainlinkNodeLogScannerSettings) *CLTestEnvBuilder {
	b.chainlinkNodeLogScannerSettings = &settings
	return b
}

func (b *CLTestEnvBuilder) WithCLNodes(clNodesCount int) *CLTestEnvBuilder {
	b.clNodesCount = clNodesCount
	return b
}

func (b *CLTestEnvBuilder) WithTestConfig(cfg ctf_config.GlobalTestConfig) *CLTestEnvBuilder {
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

func (b *CLTestEnvBuilder) WithSeth() *CLTestEnvBuilder {
	b.hasSeth = true
	b.hasEVMClient = false
	return b
}

func (b *CLTestEnvBuilder) WithPrivateEthereumNetwork(en ctf_config.EthereumNetworkConfig) *CLTestEnvBuilder {
	b.privateEthereumNetworks = append(b.privateEthereumNetworks, &en)
	return b
}

func (b *CLTestEnvBuilder) WithPrivateEthereumNetworks(ens []*ctf_config.EthereumNetworkConfig) *CLTestEnvBuilder {
	b.privateEthereumNetworks = ens
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

type EVMNetworkOption = func(*blockchain.EVMNetwork) *blockchain.EVMNetwork

func (b *CLTestEnvBuilder) EVMNetworkOptions(opts ...EVMNetworkOption) *CLTestEnvBuilder {
	b.evmNetworkOption = make([]EVMNetworkOption, 0)
	b.evmNetworkOption = append(b.evmNetworkOption, opts...)

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

	b.te.TestConfig = b.testConfig

	var err error
	if b.t != nil {
		b.te.WithTestInstance(b.t)
	}

	if b.hasLogStream {
		loggingConfig := b.testConfig.GetLoggingConfig()
		// we need to enable logging to file if we want to scan logs
		if b.chainlinkNodeLogScannerSettings != nil && !slices.Contains(loggingConfig.LogStream.LogTargets, string(logstream.File)) {
			b.l.Debug().Msg("Enabling logging to file in order to support Chainlink node log scanning")
			loggingConfig.LogStream.LogTargets = append(loggingConfig.LogStream.LogTargets, string(logstream.File))
		}
		b.te.LogStream, err = logstream.NewLogStream(b.te.t, b.testConfig.GetLoggingConfig())
		if err != nil {
			return nil, err
		}

		// this clean up has to be added as the FIRST one, because cleanup functions are executed in reverse order (LIFO)
		if b.t != nil {
			b.t.Cleanup(func() {
				b.l.Info().Msg("Shutting down LogStream")
				logPath, err := osutil.GetAbsoluteFolderPath("logs")
				if err != nil {
					b.l.Info().Str("Absolute path", logPath).Msg("LogStream logs folder location")
				}

				var scanClNodeLogs = func() {
					//filter out non-cl logs
					logLocation := b.te.LogStream.GetLogLocation()
					logFiles, err := testreporters.FindAllLogFilesToScan(logLocation, "cl-node")
					if err != nil {
						b.l.Warn().Err(err).Msg("Error looking for Chainlink Node log files to scan")
					} else {
						// we ignore the context returned by errgroup here, since we have no way of interrupting ongoing scanning of logs
						verifyLogsGroup, _ := errgroup.WithContext(context.Background())
						for _, f := range logFiles {
							file := f
							verifyLogsGroup.Go(func() error {
								logErr := testreporters.VerifyLogFile(file, b.chainlinkNodeLogScannerSettings.FailingLogLevel, b.chainlinkNodeLogScannerSettings.Threshold, b.chainlinkNodeLogScannerSettings.AllowedMessages...)
								if logErr != nil {
									return errors.Wrapf(logErr, "Found a concerning log in %s", file.Name())
								}
								return nil
							})
						}
						if err := verifyLogsGroup.Wait(); err != nil {
							b.l.Error().Err(err).Msg("Found a concerning log. Failing test.")
							b.t.Fatalf("Found a concerning log in Chainklink Node logs: %v", err)
						}
					}
					b.l.Info().Msg("Finished scanning Chainlink Node logs for concerning errors")
				}

				if b.t.Failed() || *b.testConfig.GetLoggingConfig().TestLogCollect {
					// we can't do much if this fails, so we just log the error in logstream
					flushErr := b.te.LogStream.FlushAndShutdown()
					if flushErr != nil {
						b.l.Error().Err(flushErr).Msg("Error flushing and shutting down LogStream")
						return
					}
					b.te.LogStream.PrintLogTargetsLocations()
					b.te.LogStream.SaveLogLocationInTestSummary()

					// if test hasn't failed, but we have chainlinkNodeLogScannerSettings, we should check the logs
					if !b.t.Failed() && b.chainlinkNodeLogScannerSettings != nil {
						scanClNodeLogs()
					}
				} else if b.chainlinkNodeLogScannerSettings != nil {
					flushErr := b.te.LogStream.FlushAndShutdown()
					if flushErr != nil {
						b.l.Error().Err(flushErr).Msg("Error flushing and shutting down LogStream")
						return
					}

					scanClNodeLogs()
				}
			})
		}
	}

	if b.hasKillgrave {
		if b.te.DockerNetwork == nil {
			return nil, fmt.Errorf("test environment builder failed: %w", fmt.Errorf("cannot start mock adapter without a network"))
		}

		b.te.MockAdapter = test_env.NewKillgrave([]string{b.te.DockerNetwork.Name}, "", test_env.WithLogStream(b.te.LogStream))

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
			// Cleanup test environment
			if err := b.te.Cleanup(CleanupOpts{TestName: b.t.Name()}); err != nil {
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
		// this is not the cleanest way to do this, but when we originally build ethereum networks, we don't have the logstream reference
		// so we need to rebuild them here and pass logstream to them
		for i := range b.privateEthereumNetworks {
			builder := test_env.NewEthereumNetworkBuilder()
			netWithLs, err := builder.
				WithExistingConfig(*b.privateEthereumNetworks[i]).
				WithLogStream(b.te.LogStream).
				Build()
			if err != nil {
				return nil, err
			}
			b.privateEthereumNetworks[i] = &netWithLs.EthereumNetworkConfig
		}
	}

	if b.te.LogStream == nil && b.chainlinkNodeLogScannerSettings != nil {
		log.Warn().Msg("Chainlink node log scanner settings provided, but LogStream is not enabled. Ignoring Chainlink node log scanner settings, as no logs will be available.")
	}

	// in this case we will use the builder only to start chains, not the cluster, because currently we support only 1 network config per cluster
	if len(b.privateEthereumNetworks) > 1 {
		b.te.rpcProviders = make(map[int64]*test_env.RpcProvider)
		b.te.EVMNetworks = make([]*blockchain.EVMNetwork, 0)
		b.te.evmClients = make(map[int64]blockchain.EVMClient)
		for _, en := range b.privateEthereumNetworks {
			en.DockerNetworkNames = []string{b.te.DockerNetwork.Name}
			networkConfig, rpcProvider, err := b.te.StartEthereumNetwork(en)
			if err != nil {
				return nil, err
			}

			if b.hasEVMClient {
				evmClient, err := blockchain.NewEVMClientFromNetwork(networkConfig, b.l)
				if err != nil {
					return nil, err
				}
				b.te.evmClients[networkConfig.ChainID] = evmClient
			}

			if b.hasSeth {
				seth, err := actions_seth.GetChainClient(b.testConfig, networkConfig)
				if err != nil {
					return nil, err
				}

				b.te.sethClients[networkConfig.ChainID] = seth
			}

			b.te.rpcProviders[networkConfig.ChainID] = &rpcProvider
			b.te.EVMNetworks = append(b.te.EVMNetworks, &networkConfig)

		}
		err = b.te.StartClCluster(b.clNodeConfig, b.clNodesCount, b.secretsConfig, b.testConfig, b.clNodesOpts...)
		if err != nil {
			return nil, err
		}

		b.te.isSimulatedNetwork = true

		return b.te, nil
	}

	b.te.rpcProviders = make(map[int64]*test_env.RpcProvider)
	networkConfig := networks.MustGetSelectedNetworkConfig(b.testConfig.GetNetworkConfig())[0]
	// This has some hidden behavior so I'm not the biggest fan, but it matches expected behavior.
	// That is, when we specify we want to run on a live network in our config, we will run on the live network and not bother with a private network.
	// Even if we explicitly declare that we want to run on a private network in the test.
	// Keeping this a Kludge for now as SETH transition should change all of this anyway.
	if len(b.privateEthereumNetworks) == 1 {
		if networkConfig.Simulated {
			// TODO here we should save the ethereum network config to te.Cfg, but it doesn't exist at this point
			// in general it seems we have no methods for saving config to file and we only load it from file
			// but I don't know how that config file is to be created or whether anyone ever done that
			var rpcProvider test_env.RpcProvider
			b.privateEthereumNetworks[0].DockerNetworkNames = []string{b.te.DockerNetwork.Name}
			networkConfig, rpcProvider, err = b.te.StartEthereumNetwork(b.privateEthereumNetworks[0])
			if err != nil {
				return nil, err
			}
			b.te.rpcProviders[networkConfig.ChainID] = &rpcProvider
			b.te.PrivateEthereumConfigs = b.privateEthereumNetworks

			b.te.isSimulatedNetwork = true
		} else { // Only start and connect to a private network if we are using a private simulated network
			b.te.l.Warn().
				Str("Network", networkConfig.Name).
				Int64("Chain ID", networkConfig.ChainID).
				Msg("Private network config provided, but we are running on a live network. Ignoring private network config.")
			rpcProvider := test_env.NewRPCProvider(networkConfig.HTTPURLs, networkConfig.URLs, networkConfig.HTTPURLs, networkConfig.URLs)
			b.te.rpcProviders[networkConfig.ChainID] = &rpcProvider
			b.te.isSimulatedNetwork = false
		}

	}

	if !b.hasSeth && !b.hasEVMClient {
		return nil, errors.New("you need to specify, which evm client to use: Seth or EVMClient")
	}

	if b.hasSeth && b.hasEVMClient {
		return nil, errors.New("you can't use both Seth and EMVClient at the same time")
	}

	if !b.isNonEVM {
		if b.evmNetworkOption != nil && len(b.evmNetworkOption) > 0 {
			for _, fn := range b.evmNetworkOption {
				fn(&networkConfig)
			}
		}
		if b.hasEVMClient {
			bc, err := blockchain.NewEVMClientFromNetwork(networkConfig, b.l)
			if err != nil {
				return nil, err
			}

			b.te.evmClients = make(map[int64]blockchain.EVMClient)
			b.te.evmClients[networkConfig.ChainID] = bc

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

		if b.hasSeth {
			b.te.sethClients = make(map[int64]*seth.Client)
			seth, err := actions_seth.GetChainClient(b.testConfig, networkConfig)
			if err != nil {
				return nil, err
			}

			b.te.sethClients[networkConfig.ChainID] = seth
		}
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
			rpcProvider, ok := b.te.rpcProviders[networkConfig.ChainID]
			if !ok {
				return nil, fmt.Errorf("rpc provider for chain %d not found", networkConfig.ChainID)
			}
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
						chainCfg := evmCfg.Chain
						fn(&chainCfg)
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

	if len(b.privateEthereumNetworks) > 0 && b.clNodesCount > 0 && b.ETHFunds != nil {
		if b.hasEVMClient {
			b.te.ParallelTransactions(true)
			defer b.te.ParallelTransactions(false)
			if err := b.te.FundChainlinkNodes(b.ETHFunds); err != nil {
				return nil, err
			}
		}
		if b.hasSeth {
			for _, sethClient := range b.te.sethClients {
				if err := actions_seth.FundChainlinkNodesFromRootAddress(b.l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(b.te.ClCluster.NodeAPIs()), b.ETHFunds); err != nil {
					return nil, err
				}
			}
		}
	}

	var enDesc string
	if len(b.te.PrivateEthereumConfigs) > 0 {
		for _, en := range b.te.PrivateEthereumConfigs {
			enDesc += en.Describe()
		}
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
