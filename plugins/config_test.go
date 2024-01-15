package plugins

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestRegistrarConfig(t *testing.T) {
	mockCfgTracing := &MockCfgTracing{}
	registry := make(map[string]*RegisteredLoop)

	// Create a LoopRegistry instance with mockCfgTracing
	loopRegistry := &LoopRegistry{
		lggr:       logger.TestLogger(t),
		registry:   registry,
		cfgTracing: mockCfgTracing,
	}

	opts := loop.GRPCOpts{}
	rConf := NewRegistrarConfig(opts, loopRegistry.Register)

	assert.Equal(t, opts, rConf.GRPCOpts())

	id := "command-id"
	rl, err := rConf.RegisterLOOP(id, "./COMMAND")
	require.NoError(t, err)
	assert.Equal(t, registry[id], rl)
}
