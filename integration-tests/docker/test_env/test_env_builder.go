package test_env

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	ctf_docker "github.com/smartcontractkit/chainlink-testing-framework/lib/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/osutil"

	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"
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
	hasParrot                       bool
	jdConfig                        *ccip.JDConfig
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
	testreporters.NewAllowedLogMessage(
		"No live RPC nodes available",
		"Networking or infra issues can cause brief disconnections from the node to RPC nodes, especially at startup. This isn't a concern as long as the test passes otherwise",
		zapcore.DPanicLevel,
		testreporters.WarnAboutAllowedMsgs_Yes,
	),
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
	b.hasParrot = true
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

func (b *CLTestEnvBuilder) WithJobDistributor(cfg ccip.JDConfig) *CLTestEnvBuilder {
	b.jdConfig = &cfg
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
		return nil, errors.New("test config must be set")
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

	// this clean up has to be added as the FIRST one, because cleanup functions are executed in reverse order (LIFO)
	if b.t != nil && b.cleanUpType != CleanUpTypeNone {
		b.t.Cleanup(func() {
			logsDir := fmt.Sprintf("logs/%s-%s", b.t.Name(), time.Now().Format("2006-01-02T15-04-05"))
			loggingErr := ctf_docker.WriteAllContainersLogs(b.l, logsDir)
			if loggingErr != nil {
				b.l.Error().Err(loggingErr).Msg("Error writing all Docker containers logs")
			}

			if b == nil || b.te == nil || b.te.ClCluster == nil || b.te.ClCluster.Nodes == nil {
				log.Warn().Msg("Won't dump container and postgres logs, because test environment doesn't have any nodes")
				return
			}

			if b.chainlinkNodeLogScannerSettings != nil {
				var logFiles []*os.File

				// when tests run in parallel, we need to make sure that we only process logs that belong to nodes created by the current test
				// that is required, because some tests might have custom log messages that are allowed, but only for that test (e.g. because they restart the CL node)
				var belongsToCurrentEnv = func(filePath string) bool {
					for _, clNode := range b.te.ClCluster.Nodes {
						if clNode == nil {
							continue
						}
						if strings.EqualFold(filePath, clNode.ContainerName+".log") {
							return true
						}
					}
					return false
				}

				fileWalkErr := filepath.Walk(logsDir, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() && belongsToCurrentEnv(info.Name()) {
						file, fileErr := os.Open(path)
						if fileErr != nil {
							return fmt.Errorf("failed to open file %s: %w", path, fileErr)
						}
						logFiles = append(logFiles, file)
					}
					return nil
				})

				if len(logFiles) != len(b.te.ClCluster.Nodes) {
					b.l.Warn().Int("Expected", len(b.te.ClCluster.Nodes)).Int("Got", len(logFiles)).Msg("Number of log files does not match number of nodes. Some logs might be missing.")
				}

				if fileWalkErr != nil {
					b.l.Error().Err(fileWalkErr).Msg("Error walking through log files. Skipping log verification.")
				} else {
					verifyLogsGroup := &errgroup.Group{}
					for _, f := range logFiles {
						file := f
						verifyLogsGroup.Go(func() error {
							verifyErr := testreporters.VerifyLogFile(
								file,
								b.chainlinkNodeLogScannerSettings.FailingLogLevel,
								b.chainlinkNodeLogScannerSettings.Threshold,
								b.chainlinkNodeLogScannerSettings.AllowedMessages...,
							)
							_ = file.Close()
							// ignore processing errors
							if verifyErr != nil && !strings.Contains(verifyErr.Error(), testreporters.MultipleLogsAtLogLevelErr) &&
								!strings.Contains(verifyErr.Error(), testreporters.OneLogAtLogLevelErr) {
								b.l.Error().Err(verifyErr).Msg("Error processing CL node logs")

								return nil

								// if it's not a processing error, we want to fail the test; we also can stop processing logs all together at this point
							} else if verifyErr != nil &&
								(strings.Contains(verifyErr.Error(), testreporters.MultipleLogsAtLogLevelErr) ||
									strings.Contains(verifyErr.Error(), testreporters.OneLogAtLogLevelErr)) {
								return verifyErr
							}
							return nil
						})
					}

					if logVerificationErr := verifyLogsGroup.Wait(); logVerificationErr != nil {
						b.t.Errorf("Found a concerning log in Chainlink Node logs: %v", logVerificationErr)
					}
				}
			}

			b.l.Info().Msg("Starting to dump state of all Postgres DBs used by Chainlink Nodes")

			dbDumpFolder := "db_dumps"
			dbDumpPath := fmt.Sprintf("%s/%s-%s", dbDumpFolder, b.t.Name(), time.Now().Format("2006-01-02T15-04-05"))
			if err := os.MkdirAll(dbDumpPath, os.ModePerm); err != nil {
				b.l.Error().Err(err).Msg("Error creating folder for Postgres DB dump")
			} else {
				absDbDumpPath, err := osutil.GetAbsoluteFolderPath(dbDumpFolder)
				if err == nil {
					b.l.Info().Str("Absolute path", absDbDumpPath).Msg("PostgresDB dump folder location")
				}

				dbDumpGroup := sync.WaitGroup{}
				for i := 0; i < b.clNodesCount; i++ {
					dbDumpGroup.Add(1)
					go func() {
						defer dbDumpGroup.Done()
						// if something went wrong during environment setup we might not have all nodes, and we don't want an NPE
						if b == nil || b.te == nil || b.te.ClCluster == nil || b.te.ClCluster.Nodes == nil || len(b.te.ClCluster.Nodes)-1 < i || b.te.ClCluster.Nodes[i] == nil || b.te.ClCluster.Nodes[i].PostgresDb == nil {
							return
						}

						filePath := filepath.Join(dbDumpPath, fmt.Sprintf("postgres_db_dump_%s.sql", b.te.ClCluster.Nodes[i].ContainerName))
						localDbDumpFile, err := os.Create(filePath)
						if err != nil {
							b.l.Error().Err(err).Msg("Error creating localDbDumpFile for Postgres DB dump")
							_ = localDbDumpFile.Close()
							return
						}

						if err := b.te.ClCluster.Nodes[i].PostgresDb.ExecPgDumpFromContainer(localDbDumpFile); err != nil {
							b.l.Error().Err(err).Msg("Error dumping Postgres DB")
						}
						_ = localDbDumpFile.Close()
					}()
				}

				dbDumpGroup.Wait()

				b.l.Info().Msg("Finished dumping state of all Postgres DBs used by Chainlink Nodes")
			}
		})
	} else {
		b.l.Warn().Msg("Won't dump container and postgres logs, because either test instance is not set or cleanup type is set to none")
	}

	if b.hasParrot {
		if b.te.DockerNetwork == nil {
			return nil, fmt.Errorf("test environment builder failed: %w", errors.New("cannot start mock adapter without a network"))
		}

		b.te.MockAdapter = test_env.NewParrot(
			[]string{b.te.DockerNetwork.Name},
			test_env.WithStartupTimeout(2*time.Minute),
		)

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
		return b.te, fmt.Errorf("test environment builder failed: %w", errors.New("explicit cleanup type must be set when building test environment"))
	}

	if b.jdConfig != nil {
		err := b.te.StartJobDistributor(b.jdConfig)
		if err != nil {
			return nil, err
		}
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
		if len(b.evmNetworkOption) > 0 {
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
		Bool("hasParrot", b.hasParrot).
		Bool("hasJobDistributor", b.jdConfig != nil).
		Int("clNodesCount", b.clNodesCount).
		Strs("customNodeCsaKeys", b.customNodeCsaKeys).
		Strs("defaultNodeCsaKeys", b.defaultNodeCsaKeys).
		Msg("Building CL cluster test environment..")

	return b.te, nil
}
