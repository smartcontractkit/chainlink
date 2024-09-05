package hooks

import (
	"github.com/rs/zerolog"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	persistent_types "github.com/smartcontractkit/chainlink/integration-tests/deployment/persistent/types"
	"testing"

	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

type DefaultEnvironmentHooks struct {
	L                      zerolog.Logger
	T                      *testing.T
	RunId                  *string
	ShowHTMLCoverageReport bool
	SethConfig             *seth.Config
}

func NewDefaultEnvironmentHooks(t *testing.T, logger zerolog.Logger, sethConfig *seth.Config, runId *string, showHTMLCoverageReport bool) *DefaultEnvironmentHooks {
	return &DefaultEnvironmentHooks{
		T:                      t,
		L:                      logger,
		RunId:                  runId,
		ShowHTMLCoverageReport: showHTMLCoverageReport,
		SethConfig:             sethConfig,
	}
}

func NewDefaultEnvironmentHooksFromTestConfig(t *testing.T, logger zerolog.Logger, sethConfig *seth.Config, loggingConfig *ctf_config.LoggingConfig) *DefaultEnvironmentHooks {
	return NewDefaultEnvironmentHooks(t, logger, sethConfig, loggingConfig.RunId, loggingConfig.ShowHTMLCoverageReport != nil && *loggingConfig.ShowHTMLCoverageReport)
}

func (d *DefaultEnvironmentHooks) PostChainStartupHooks(_ map[uint64]deployment.Chain, _ map[uint64]persistent_types.RpcProvider, _ *persistent_types.EnvironmentConfig) error {
	return nil
}

func (d *DefaultEnvironmentHooks) PostNodeStartupHooks(don *persistent_types.DON, _ *persistent_types.EnvironmentConfig) error {
	test_env.AttachDefaultCleanUp(d.L, d.T, &test_env.ClCluster{Nodes: don.ChainlinkContainers}, d.ShowHTMLCoverageReport, d.RunId)
	test_env.AttachSethCleanup(d.T, d.SethConfig)

	return nil
}

func (d *DefaultEnvironmentHooks) PostMocksStartupHooks(_ *deployment.Mocks, _ *persistent_types.EnvironmentConfig) error {
	return nil
}

type DefaultDockerDONHooks struct {
	T                               *testing.T
	L                               zerolog.Logger
	LogStream                       *logstream.LogStream
	ChainlinkNodeLogScannerSettings *test_env.ChainlinkNodeLogScannerSettings
	// this one also affect other containers, no only DON, but I didn't find a way to properly separate them yet
	SaveContainerArtifacts bool
}

func NewDefaultDONHooks(t *testing.T, logger zerolog.Logger, logStream *logstream.LogStream, chainlinkNodeLogScannerSettings *test_env.ChainlinkNodeLogScannerSettings, saveContainerArtifacts bool) *DefaultDockerDONHooks {
	if chainlinkNodeLogScannerSettings == nil {
		chainlinkNodeLogScannerSettings = &test_env.DefaultChainlinkNodeLogScannerSettings
	} else {
		chainlinkNodeLogScannerSettings.AllowedMessages = append(chainlinkNodeLogScannerSettings.AllowedMessages, test_env.DefaultChainlinkNodeLogScannerSettings.AllowedMessages...)
	}
	return &DefaultDockerDONHooks{
		T:                               t,
		L:                               logger,
		LogStream:                       logStream,
		ChainlinkNodeLogScannerSettings: chainlinkNodeLogScannerSettings,
		//RunId:                           runId,
		SaveContainerArtifacts: saveContainerArtifacts,
		//ShowHTMLCoverageReport:          showHTMLCoverageReport,
	}
}

func NewDefaultDONHooksFromTestConfig(t *testing.T, logger zerolog.Logger, logStream *logstream.LogStream, chainlinkNodeLogScannerSettings *test_env.ChainlinkNodeLogScannerSettings, loggingConfig *ctf_config.LoggingConfig) *DefaultDockerDONHooks {
	return NewDefaultDONHooks(t, logger, logStream, chainlinkNodeLogScannerSettings, loggingConfig.TestLogCollect != nil && *loggingConfig.TestLogCollect)
}

func (s *DefaultDockerDONHooks) PreStartupHook(nodes []*test_env.ClNode) error {
	for _, node := range nodes {
		node.SetTestLogger(s.T)
		node.LogStream = s.LogStream
	}

	return nil
}

func (s *DefaultDockerDONHooks) PostStartupHook(nodes []*test_env.ClNode) error {
	test_env.AttachLogStreamCleanUp(s.L, s.T, s.LogStream, &test_env.ClCluster{Nodes: nodes}, s.ChainlinkNodeLogScannerSettings, s.SaveContainerArtifacts)
	test_env.AttachDbDumpingCleanup(s.L, s.T, &test_env.ClCluster{Nodes: nodes}, s.SaveContainerArtifacts)
	//test_env.AttachDefaultCleanUp(s.L, s.T, &test_env.ClCluster{Nodes: nodes}, s.ShowHTMLCoverageReport, s.RunId)
	//test_env.AttachSethCleanup(s.T, s.SethConfig)

	return nil
}

type DefaultPrivateEVMHooks struct {
	T         *testing.T
	LogStream *logstream.LogStream
}

func (s *DefaultPrivateEVMHooks) PreStartEnvComponentHooks() []ctf_test_env.EnvComponentOption {
	var opts []ctf_test_env.EnvComponentOption
	if s.LogStream != nil {
		opts = append(opts, ctf_test_env.WithLogStream(s.LogStream))
	}

	if s.T != nil {
		opts = append(opts, ctf_test_env.WithTestInstance(s.T))
	}

	return opts
}

func (s *DefaultPrivateEVMHooks) PostStartEnvComponentHooks() []ctf_test_env.EnvComponentOption {
	return nil
}
