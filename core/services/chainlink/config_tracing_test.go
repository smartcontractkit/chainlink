package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

func TestTracing_Config(t *testing.T) {
	// Test when all fields are non-nil
	enabled := true
	collectorTarget := "http://localhost:9000"
	nodeID := "Node1"
	samplingRatio := 0.5
	attributes := map[string]string{"key": "value"}
	tracing := toml.Tracing{
		Enabled:         &enabled,
		CollectorTarget: &collectorTarget,
		NodeID:          &nodeID,
		SamplingRatio:   &samplingRatio,
		Attributes:      attributes,
	}
	tConfig := tracingConfig{s: tracing}

	assert.True(t, tConfig.Enabled())
	assert.Equal(t, "http://localhost:9000", tConfig.CollectorTarget())
	assert.Equal(t, "Node1", tConfig.NodeID())
	assert.Equal(t, 0.5, tConfig.SamplingRatio())
	assert.Equal(t, map[string]string{"key": "value"}, tConfig.Attributes())

	// Test when all fields are nil
	nilTracing := toml.Tracing{}
	nilConfig := tracingConfig{s: nilTracing}

	assert.Panics(t, func() { nilConfig.Enabled() })
	assert.Panics(t, func() { nilConfig.CollectorTarget() })
	assert.Panics(t, func() { nilConfig.NodeID() })
	assert.Panics(t, func() { nilConfig.SamplingRatio() })
	assert.Nil(t, nilConfig.Attributes())
}
