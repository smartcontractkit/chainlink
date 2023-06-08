package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEncodeThresholdReportingPluginConfig(t *testing.T) {
	// Functions plugin config
	raw := []byte{8, 144, 78, 16, 144, 78, 24, 144, 78, 32, 10, 48, 1, 58, 17, 8, 144, 78, 16, 144, 78, 24, 144, 78, 32, 100, 40, 160, 141, 6, 48, 1}

	decoded, err := GetThresholdReportingPluginConfig(raw)
	require.NoError(t, err)
	require.NotNil(t, decoded)

	encoded, err := EncodeThresholdPluginConfig(decoded)
	require.NoError(t, err)

	// Threshold plugin config
	expected := []byte{8, 144, 78, 16, 144, 78, 24, 144, 78, 32, 100, 40, 160, 141, 6, 48, 1}

	assert.Equal(t, expected, encoded)
}
