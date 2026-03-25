package chainlink

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountEnabledEVMConfigs(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 0, countEnabledEVMConfigs(nil))
}

func TestIsUnknownChainSelectorErr(t *testing.T) {
	t.Parallel()
	assert.False(t, isUnknownChainSelectorErr(nil))
	assert.True(t, isUnknownChainSelectorErr(errors.New("cannot create evm relayer: chain-selectors missing chain id 123: chain selector not found for chain 123")))
	assert.True(t, isUnknownChainSelectorErr(errors.New("wrapped: chain selector not found for chain 1")))
	assert.False(t, isUnknownChainSelectorErr(errors.New("some other failure")))
}

func TestEVMChainConfigHealth_Ready(t *testing.T) {
	t.Parallel()
	h := NewEVMChainConfigHealth([]string{"1", "2"})
	require.NoError(t, h.Ready())

	h.RecordSkipped("2", evmSkipReasonRelayerInit, errors.New("boom"))
	err := h.Ready()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "degraded")
}

func TestEnabledEVMChainIDStrings(t *testing.T) {
	t.Parallel()
	assert.Empty(t, enabledEVMChainIDStrings(nil))
}
