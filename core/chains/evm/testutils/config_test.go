package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

func TestNewTestChainScopedConfigOverride(t *testing.T) {
	c := NewTestChainScopedConfig(t, func(c *toml.EVMConfig) {
		finalityDepth := uint32(100)
		c.FinalityDepth = &finalityDepth
	})

	// Overrides values
	assert.Equal(t, uint32(100), c.EVM().FinalityDepth())
	// fallback.toml values
	assert.Equal(t, false, c.EVM().GasEstimator().EIP1559DynamicFees())
}
