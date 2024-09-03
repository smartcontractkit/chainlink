package hooks

import (
	"github.com/rs/zerolog"
	"testing"

	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

type DefaultDONHooks struct {
	T                               *testing.T
	L                               zerolog.Logger
	LogStream                       *logstream.LogStream
	RunId                           *string
	CollectTestArtifacts            bool
	ChainlinkNodeLogScannerSettings *test_env.ChainlinkNodeLogScannerSettings
	ShowHTMLCoverageReport          bool
}

func (s *DefaultDONHooks) PreStartupHook(nodes []*test_env.ClNode) error {
	for _, node := range nodes {
		node.SetTestLogger(s.T)
		node.LogStream = s.LogStream
	}

	return nil
}

func (s *DefaultDONHooks) PostStartupHook(nodes []*test_env.ClNode) error {
	test_env.AttachLogStreamCleanUp(s.L, s.T, s.LogStream, nodes, s.ChainlinkNodeLogScannerSettings, s.CollectTestArtifacts)
	test_env.AttachDbDumpingCleanup(s.L, s.T, nodes, s.CollectTestArtifacts)
	test_env.AttachDefaultCleanUp(s.L, s.T, nodes, s.ShowHTMLCoverageReport, s.RunId)
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
