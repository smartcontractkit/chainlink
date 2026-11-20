package ccvexecutor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccv/executor"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
)

func TestMemberChains_OnlyChainsWhereExecutorIsInPool(t *testing.T) {
	cfg := executor.Configuration{
		ExecutorID: "exec-1",
		ChainConfiguration: map[string]executor.ChainConfiguration{
			"1": {ExecutorPool: []string{"exec-1", "exec-2"}},
			"2": {ExecutorPool: []string{"exec-2"}},
			"3": {ExecutorPool: []string{"exec-1"}},
		},
	}

	got, err := memberChains(cfg)
	require.NoError(t, err)
	assert.ElementsMatch(t, []protocol.ChainSelector{1, 3}, got)
}

func TestMemberChains_InvalidSelector(t *testing.T) {
	cfg := executor.Configuration{
		ExecutorID: "exec-1",
		ChainConfiguration: map[string]executor.ChainConfiguration{
			"not-a-number": {ExecutorPool: []string{"exec-1"}},
		},
	}

	_, err := memberChains(cfg)
	require.Error(t, err)
}
