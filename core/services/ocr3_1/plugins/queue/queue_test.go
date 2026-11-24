package queue

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3_1/plugins/queue/queuepb"
)

// TestConsensusQueuePluginFactory_NewReportingPlugin asserts the plugin takes its limits from the OCR offchain config
// and refuses to start without a valid one. libocr calls NewReportingPlugin with no recover, so a panic here would
// take down the whole node.
func TestConsensusQueuePluginFactory_NewReportingPlugin(t *testing.T) {
	t.Parallel()

	validCfg, err := proto.Marshal(&queuepb.ReportingPluginConfig{
		Limits: &queuepb.ReportingPluginLimits{
			MaxObservationBytes:                    1024,
			MaxKeyValueModifiedKeys:                10,
			MaxKeyValueModifiedKeysPlusValuesBytes: 2048,
		},
	})
	require.NoError(t, err)

	t.Run("limits come from the OCR config", func(t *testing.T) {
		t.Parallel()
		factory, err := NewConsensusQueuePluginFactory(logger.Test(t), limits.Factory{})
		require.NoError(t, err)

		plugin, info, err := factory.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{N: 4, F: 1, OffchainConfig: validCfg}, nil)
		require.NoError(t, err)
		require.NotNil(t, plugin)
		t.Cleanup(func() { require.NoError(t, plugin.Close()) })

		info1, ok := info.(ocr3_1types.ReportingPluginInfo1)
		require.True(t, ok, "unexpected ReportingPluginInfo type %T", info)
		require.Equal(t, "ConsensusQueueReportingPlugin", info1.Name)
		want := ocr3_1types.ReportingPluginLimits{
			MaxObservationBytes:                    1024,
			MaxKeyValueModifiedKeys:                10,
			MaxKeyValueModifiedKeysPlusValuesBytes: 2048,
		}
		require.Equal(t, want, info1.Limits)
	})

	tests := []struct {
		name string
		got  []byte
	}{
		{name: "empty config", got: nil},
		{name: "malformed config", got: []byte{0xff}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			factory, err := NewConsensusQueuePluginFactory(logger.Test(t), limits.Factory{})
			require.NoError(t, err)

			plugin, _, err := factory.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{N: 4, F: 1, OffchainConfig: tt.got}, nil)
			require.Error(t, err)
			require.Nil(t, plugin)
		})
	}
}
