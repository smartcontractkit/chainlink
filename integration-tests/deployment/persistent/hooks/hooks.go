package hooks

import (
	"testing"

	ctf_test_env "github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logstream"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
)

type DefaultDONHooks struct {
	T *testing.T
}

func (s *DefaultDONHooks) PreStartupHook(nodes []*test_env.ClNode) error {
	for _, node := range nodes {
		node.SetTestLogger(s.T)
	}

	return nil
}

func (s *DefaultDONHooks) PostStartupHook(_ []*test_env.ClNode) error {
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
