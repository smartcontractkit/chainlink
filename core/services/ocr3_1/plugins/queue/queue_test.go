package queue

import (
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

// libocr calls NewReportingPlugin with no recover, so a panic here takes down
// the whole node. This guards the constructor and the plugin lifecycle.
func TestConsensusQueuePluginFactory_NewReportingPlugin(t *testing.T) {
	t.Parallel()

	cfg := &PluginConfig{Limits: ocr3_1types.ReportingPluginLimits{MaxKeyValueModifiedKeys: 1}}
	factory, err := NewConsensusQueuePluginFactory(logger.Test(t), cfg, limits.Factory{})
	require.NoError(t, err)

	plugin, info, err := factory.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{N: 4, F: 1}, nil)
	require.NoError(t, err)
	require.NotNil(t, plugin)

	info1, ok := info.(ocr3_1types.ReportingPluginInfo1)
	require.True(t, ok, "unexpected ReportingPluginInfo type %T", info)
	require.Equal(t, "ConsensusQueueReportingPlugin", info1.Name)
	require.Equal(t, cfg.Limits, info1.Limits)

	require.NoError(t, plugin.Close())
}
