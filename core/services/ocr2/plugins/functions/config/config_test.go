package config_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestS4ConfigDecoder(t *testing.T) {
	t.Parallel()

	configProto := &config.ReportingPluginConfig{
		S4PluginConfig: &config.S4ReportingPluginConfig{
			MaxQueryLengthBytes:       100,
			MaxObservationLengthBytes: 200,
			MaxReportLengthBytes:      300,
			NSnapshotShards:           1,
			MaxObservationEntries:     111,
			MaxReportEntries:          222,
			MaxDeleteExpiredEntries:   333,
		},
	}

	configBytes, err := proto.Marshal(configProto)
	require.NoError(t, err)

	config, limits, err := config.S4ConfigDecoder(configBytes)
	require.NoError(t, err)
	assert.Equal(t, "functions", config.ProductName)
	assert.Equal(t, uint(1), config.NSnapshotShards)
	assert.Equal(t, uint(111), config.MaxObservationEntries)
	assert.Equal(t, uint(222), config.MaxReportEntries)
	assert.Equal(t, uint(333), config.MaxDeleteExpiredEntries)
	assert.Equal(t, 100, limits.MaxQueryLength)
	assert.Equal(t, 200, limits.MaxObservationLength)
	assert.Equal(t, 300, limits.MaxReportLength)
}
