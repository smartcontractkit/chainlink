package hooks

import (
	"github.com/rs/zerolog"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/config"
	"testing"

	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

type DefaultDONHooks struct {
	T                               *testing.T
	L                               zerolog.Logger
	LogStream                       *logstream.LogStream
	ChainlinkNodeLogScannerSettings *test_env.ChainlinkNodeLogScannerSettings
	// logically these 3 are not related to DON as such, but rather test specific
	RunId                  *string
	CollectTestArtifacts   bool
	ShowHTMLCoverageReport bool
}

func NewDefaultDONHooks(t *testing.T, logger zerolog.Logger, logStream *logstream.LogStream, chainlinkNodeLogScannerSettings *test_env.ChainlinkNodeLogScannerSettings, runId *string, collectTestArtifacts, showHTMLCoverageReport bool) *DefaultDONHooks {
	if chainlinkNodeLogScannerSettings == nil {
		chainlinkNodeLogScannerSettings = &test_env.DefaultChainlinkNodeLogScannerSettings
	} else {
		chainlinkNodeLogScannerSettings.AllowedMessages = append(chainlinkNodeLogScannerSettings.AllowedMessages, test_env.DefaultChainlinkNodeLogScannerSettings.AllowedMessages...)
	}
	return &DefaultDONHooks{
		T:                               t,
		L:                               logger,
		LogStream:                       logStream,
		ChainlinkNodeLogScannerSettings: chainlinkNodeLogScannerSettings,
		RunId:                           runId,
		CollectTestArtifacts:            collectTestArtifacts,
		ShowHTMLCoverageReport:          showHTMLCoverageReport,
	}
}

func NewDefaultDONHooksFromTestConfig(t *testing.T, logger zerolog.Logger, logStream *logstream.LogStream, chainlinkNodeLogScannerSettings *test_env.ChainlinkNodeLogScannerSettings, loggingConfig *ctf_config.LoggingConfig) *DefaultDONHooks {
	return NewDefaultDONHooks(t, logger, logStream, chainlinkNodeLogScannerSettings, nil, loggingConfig.TestLogCollect != nil && *loggingConfig.TestLogCollect, loggingConfig.ShowHTMLCoverageReport != nil && *loggingConfig.ShowHTMLCoverageReport)
}

func (s *DefaultDONHooks) PreStartupHook(nodes []*test_env.ClNode) error {
	for _, node := range nodes {
		node.SetTestLogger(s.T)
		node.LogStream = s.LogStream
	}

	return nil
}

func (s *DefaultDONHooks) PostStartupHook(nodes []*test_env.ClNode) error {
	test_env.AttachLogStreamCleanUp(s.L, s.T, s.LogStream, &test_env.ClCluster{Nodes: nodes}, s.ChainlinkNodeLogScannerSettings, s.CollectTestArtifacts)
	test_env.AttachDbDumpingCleanup(s.L, s.T, &test_env.ClCluster{Nodes: nodes}, s.CollectTestArtifacts)
	test_env.AttachDefaultCleanUp(s.L, s.T, &test_env.ClCluster{Nodes: nodes}, s.ShowHTMLCoverageReport, s.RunId)
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
