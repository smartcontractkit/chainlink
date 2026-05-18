package test_env

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/testsummary"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/osutil"

	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
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
	clNodeConfig                    *chainlink.Config
	secretsConfig                   string
	clNodesCount                    int
	clNodesOpts                     []func(*ClNode)
	customNodeCsaKeys               []string
	defaultNodeCsaKeys              []string
	l                               zerolog.Logger
	t                               *testing.T
	te                              *CLClusterTestEnv
	isEVM                           bool
	cleanUpType                     CleanUpType
	cleanUpCustomFn                 func()
	evmNetworkOption                []EVMNetworkOption
	privateEthereumNetworks         []*ctf_config.EthereumNetworkConfig
	testConfig                      ctf_config.GlobalTestConfig
	chainlinkNodeLogScannerSettings *ChainlinkNodeLogScannerSettings
}

var DefaultAllowedMessages = []testreporters.AllowedLogMessage{
	testreporters.NewAllowedLogMessage("Failed to get LINK balance", "Happens only when we deploy LINK token for test purposes. Harmless.", zapcore.ErrorLevel, testreporters.WarnAboutAllowedMsgs_No),
	testreporters.NewAllowedLogMessage("Error stopping job service", "It's a known issue with lifecycle. There's ongoing work that will fix it.", zapcore.DPanicLevel, testreporters.WarnAboutAllowedMsgs_No),
}

var DefaultChainlinkNodeLogScannerSettings = ChainlinkNodeLogScannerSettings{
	FailingLogLevel: zapcore.DPanicLevel,
	Threshold:       1, // we want to fail on the first concerning log
	AllowedMessages: DefaultAllowedMessages,
}

func GetDefaultChainlinkNodeLogScannerSettingsWithExtraAllowedMessages(extraAllowedMessages ...testreporters.AllowedLogMessage) ChainlinkNodeLogScannerSettings {
	allowedMessages := append(DefaultAllowedMessages, extraAllowedMessages...)
	return ChainlinkNodeLogScannerSettings{
		FailingLogLevel: DefaultChainlinkNodeLogScannerSettings.FailingLogLevel,
		Threshold:       DefaultChainlinkNodeLogScannerSettings.Threshold,
		AllowedMessages: allowedMessages,
	}
}

func NewCLTestEnvBuilder() *CLTestEnvBuilder {
	return &CLTestEnvBuilder{
		l:                               log.Logger,
		hasLogStream:                    true,
		isEVM:                           true,
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

func (b *CLTestEnvBuilder) WithPrivateEthereumNetwork(en ctf_config.EthereumNetworkConfig) *CLTestEnvBuilder {
	b.privateEthereumNetworks = append(b.privateEthereumNetworks, &en)
	return b
}

func (b *CLTestEnvBuilder) WithPrivateEthereumNetworks(ens []*ctf_config.EthereumNetworkConfig) *CLTestEnvBuilder {
	b.privateEthereumNetworks = ens
	return b
}

// Deprecated: Use TOML instead
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
	b.isEVM = false
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

type EVMNetworkOption = func(*blockchain.EVMNetwork) *blockchain.EVMNetwork

// WithEVMNetworkOptions sets the options for the EVM network. This is especially useful for simulated networks, which
// by usually use default options, so if we want to change any of them before the configuration is passed to evm client
// or Chainlnik node, we can do it here.
func (b *CLTestEnvBuilder) WithEVMNetworkOptions(opts ...EVMNetworkOption) *CLTestEnvBuilder {
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
		if b.t != nil && b.cleanUpType != CleanUpTypeNone {
			b.t.Cleanup(func() {
				b.l.Info().Msg("Shutting down LogStream")
				logPath, err := osutil.GetAbsoluteFolderPath("logs")
				if err == nil {
					b.l.Info().Str("Absolute path", logPath).Msg("LogStream logs folder location")
				}

				// flush logs when test failed or when we are explicitly told to collect logs
				flushLogStream := b.t.Failed() || *b.testConfig.GetLoggingConfig().TestLogCollect

				// run even if test has failed, as we might be able to catch additional problems without running the test again
				if b.chainlinkNodeLogScannerSettings != nil {
					logProcessor := logstream.NewLogProcessor[int](b.te.LogStream)

					processFn := func(log logstream.LogContent, count *int) error {
						countSoFar := count
						newCount, err := testreporters.ScanLogLine(b.l, string(log.Content), b.chainlinkNodeLogScannerSettings.FailingLogLevel, uint(*countSoFar), b.chainlinkNodeLogScannerSettings.Threshold, b.chainlinkNodeLogScannerSettings.AllowedMessages)
						if err != nil {
							return err
						}
						*count = int(newCount)
						return nil
					}

					// we cannot do parallel processing here, because ProcessContainerLogs() locks a mutex that controls whether
					// new logs can be added to the log stream, so parallel processing would get stuck on waiting for it to be unlocked
				LogScanningLoop:
					for i := 0; i < b.clNodesCount; i++ {
						// if something went wrong during environment setup we might not have all nodes, and we don't want an NPE
						if b == nil || b.te == nil || b.te.ClCluster == nil || b.te.ClCluster.Nodes == nil || len(b.te.ClCluster.Nodes)-1 < i || b.te.ClCluster.Nodes[i] == nil {
							continue
						}
						// ignore count return, because we are only interested in the error
						_, err := logProcessor.ProcessContainerLogs(b.te.ClCluster.Nodes[i].ContainerName, processFn)
						if err != nil && !strings.Contains(err.Error(), testreporters.MultipleLogsAtLogLevelErr) && !strings.Contains(err.Error(), testreporters.OneLogAtLogLevelErr) {
							b.l.Error().Err(err).Msg("Error processing CL node logs")
							continue
						} else if err != nil && (strings.Contains(err.Error(), testreporters.MultipleLogsAtLogLevelErr) || strings.Contains(err.Error(), testreporters.OneLogAtLogLevelErr)) {
							flushLogStream = true
							b.t.Errorf("Found a concerning log in Chainklink Node logs: %v", err)
							break LogScanningLoop
						}
					}
					b.l.Info().Msg("Finished scanning Chainlink Node logs for concerning errors")
				}

				if flushLogStream {
					b.l.Info().Msg("Flushing LogStream logs")
					// we can't do much if this fails, so we just log the error in LogStream
					if err := b.te.LogStream.FlushAndShutdown(); err != nil {
						b.l.Error().Err(err).Msg("Error flushing and shutting down LogStream")
					}
					b.te.LogStream.PrintLogTargetsLocations()
					b.te.LogStream.SaveLogLocationInTestSummary()
				}
				b.l.Info().Msg("Finished shutting down LogStream")

				if b.t.Failed() || *b.testConfig.GetLoggingConfig().TestLogCollect {
					b.l.Info().Msg("Dump state of all Postgres DBs used by Chainlink Nodes")

					dbDumpFolder := "db_dumps"
					dbDumpPath := fmt.Sprintf("%s/%s-%s", dbDumpFolder, b.t.Name(), time.Now().Format("2006-01-02T15-04-05"))
					if err := os.MkdirAll(dbDumpPath, os.ModePerm); err != nil {
						b.l.Error().Err(err).Msg("Error creating folder for Postgres DB dump")
						return
					}

					absDbDumpPath, err := osutil.GetAbsoluteFolderPath(dbDumpFolder)
					if err == nil {
						b.l.Info().Str("Absolute path", absDbDumpPath).Msg("PostgresDB dump folder location")
					}

					for i := 0; i < b.clNodesCount; i++ {
						// if something went wrong during environment setup we might not have all nodes, and we don't want an NPE
						if b == nil || b.te == nil || b.te.ClCluster == nil || b.te.ClCluster.Nodes == nil || len(b.te.ClCluster.Nodes)-1 < i || b.te.ClCluster.Nodes[i] == nil || b.te.ClCluster.Nodes[i].PostgresDb == nil {
							continue
						}

						filePath := filepath.Join(dbDumpPath, fmt.Sprintf("postgres_db_dump_%s.sql", b.te.ClCluster.Nodes[i].ContainerName))
						localDbDumpFile, err := os.Create(filePath)
						if err != nil {
							b.l.Error().Err(err).Msg("Error creating localDbDumpFile for Postgres DB dump")
							_ = localDbDumpFile.Close()
							continue
						}

						if err := b.te.ClCluster.Nodes[i].PostgresDb.ExecPgDumpFromContainer(localDbDumpFile); err != nil {
							b.l.Error().Err(err).Msg("Error dumping Postgres DB")
						}
						_ = localDbDumpFile.Close()
					}
					b.l.Info().Msg("Finished dumping state of all Postgres DBs used by Chainlink Nodes")
				}

				if b.testConfig.GetSethConfig() != nil && ((b.t.Failed() && slices.Contains(b.testConfig.GetSethConfig().TraceOutputs, seth.TraceOutput_DOT) && b.testConfig.GetSethConfig().TracingLevel != seth.TracingLevel_None) || (!b.t.Failed() && slices.Contains(b.testConfig.GetSethConfig().TraceOutputs, seth.TraceOutput_DOT) && b.testConfig.GetSethConfig().TracingLevel == seth.TracingLevel_All)) {
					_ = testsummary.AddEntry(b.t.Name(), "dot_graphs", "true")
				}
			})
		} else {
			b.l.Warn().Msg("LogStream won't be cleaned up, because either test instance is not set or cleanup type is set to none")
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

	if b.te.LogStream == nil && b.chainlinkNodeLogScannerSettings != nil {
		log.Warn().Msg("Chainlink node log scanner settings provided, but LogStream is not enabled. Ignoring Chainlink node log scanner settings, as no logs will be available.")
	}

	// in this case we will use the builder only to start chains, not the cluster, because currently we support only 1 network config per cluster
	if len(b.privateEthereumNetworks) > 1 {
		b.te.rpcProviders = make(map[int64]*test_env.RpcProvider)
		b.te.EVMNetworks = make([]*blockchain.EVMNetwork, 0)
		for _, en := range b.privateEthereumNetworks {
			en.DockerNetworkNames = []string{b.te.DockerNetwork.Name}
			networkConfig, rpcProvider, err := b.te.StartEthereumNetwork(en)
			if err != nil {
				return nil, err
			}

			b.te.rpcProviders[networkConfig.ChainID] = &rpcProvider
			b.te.EVMNetworks = append(b.te.EVMNetworks, &networkConfig)
		}
		if b.clNodesCount > 0 {
			dereferrencedEvms := make([]blockchain.EVMNetwork, 0)
			for _, en := range b.te.EVMNetworks {
				dereferrencedEvms = append(dereferrencedEvms, *en)
			}

			nodeConfigInToml := b.testConfig.GetNodeConfig()

			nodeConfig, _, err := node.BuildChainlinkNodeConfig(
				dereferrencedEvms,
				nodeConfigInToml.BaseConfigTOML,
				nodeConfigInToml.CommonChainConfigTOML,
				nodeConfigInToml.ChainConfigTOMLByChainID,
			)
			if err != nil {
				return nil, err
			}

			err = b.te.StartClCluster(nodeConfig, b.clNodesCount, b.secretsConfig, b.testConfig, b.clNodesOpts...)
			if err != nil {
				return nil, err
			}
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
	b.te.EVMNetworks = make([]*blockchain.EVMNetwork, 0)
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
	} else if len(b.privateEthereumNetworks) == 0 && !networkConfig.Simulated {
		b.te.l.Warn().
			Str("Network", networkConfig.Name).
			Int64("Chain ID", networkConfig.ChainID).
			Msg("Private network config provided, but we are running on a live network. Ignoring private network config.")
		rpcProvider := test_env.NewRPCProvider(networkConfig.HTTPURLs, networkConfig.URLs, networkConfig.HTTPURLs, networkConfig.URLs)
		b.te.rpcProviders[networkConfig.ChainID] = &rpcProvider
		b.te.isSimulatedNetwork = false
	}
	b.te.EVMNetworks = append(b.te.EVMNetworks, &networkConfig)

	if b.isEVM {
		if b.evmNetworkOption != nil && len(b.evmNetworkOption) > 0 {
			for _, fn := range b.evmNetworkOption {
				fn(&networkConfig)
			}
		}
	}

	var nodeCsaKeys []string

	// Start Chainlink Nodes
	if b.clNodesCount > 0 {
		// needed for live networks
		if len(b.te.EVMNetworks) == 0 {
			b.te.EVMNetworks = append(b.te.EVMNetworks, &networkConfig)
		}

		// only add EVM networks to node config if running EVM tests
		dereferrencedEvms := make([]blockchain.EVMNetwork, 0)
		if b.isEVM {
			for _, en := range b.te.EVMNetworks {
				network := *en
				if en.Simulated {
					if rpcs, ok := b.te.rpcProviders[network.ChainID]; ok {
						network.HTTPURLs = rpcs.PrivateHttpUrls()
						network.URLs = rpcs.PrivateWsUrsl()
					} else {
						return nil, fmt.Errorf("rpc provider for chain %d not found", network.ChainID)
					}
				}
				dereferrencedEvms = append(dereferrencedEvms, network)
			}
		}

		nodeConfigInToml := b.testConfig.GetNodeConfig()

		nodeConfig, _, err := node.BuildChainlinkNodeConfig(
			dereferrencedEvms,
			nodeConfigInToml.BaseConfigTOML,
			nodeConfigInToml.CommonChainConfigTOML,
			nodeConfigInToml.ChainConfigTOMLByChainID,
		)
		if err != nil {
			return nil, err
		}

		err = b.te.StartClCluster(nodeConfig, b.clNodesCount, b.secretsConfig, b.testConfig, b.clNodesOpts...)
		if err != nil {
			return nil, err
		}

		nodeCsaKeys, err = b.te.ClCluster.NodeCSAKeys()
		if err != nil {
			return nil, err
		}
		b.defaultNodeCsaKeys = nodeCsaKeys
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
