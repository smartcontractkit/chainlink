package s4_test

import (
	"errors"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/s4"
	s4_mocks "github.com/smartcontractkit/chainlink/v2/core/services/s4/mocks"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/stretchr/testify/require"
)

func TestS4ReportingPluginFactory_NewReportingPlugin(t *testing.T) {
	t.Parallel()

	logger := commonlogger.NewOCRWrapper(logger.TestLogger(t), true, func(msg string) {})
	orm := s4_mocks.NewORM(t)

	f := s4.S4ReportingPluginFactory{
		Logger: logger,
		ORM:    orm,
		ConfigDecoder: func([]byte) (*s4.PluginConfig, *types.ReportingPluginLimits, error) {
			return &s4.PluginConfig{
					ProductName:             "test",
					NSnapshotShards:         1,
					MaxObservationEntries:   10,
					MaxReportEntries:        20,
					MaxDeleteExpiredEntries: 30,
				}, &types.ReportingPluginLimits{
					MaxQueryLength:       100,
					MaxObservationLength: 200,
					MaxReportLength:      300,
				}, nil
		},
	}

	rpConfig := types.ReportingPluginConfig{
		OffchainConfig: make([]byte, 100),
	}
	plugin, pluginInfo, err := f.NewReportingPlugin(rpConfig)
	require.NoError(t, err)
	require.NotNil(t, plugin)
	require.Equal(t, types.ReportingPluginInfo{
		Name:          s4.S4ReportingPluginName,
		UniqueReports: false,
		Limits: types.ReportingPluginLimits{
			MaxQueryLength:       100,
			MaxObservationLength: 200,
			MaxReportLength:      300,
		},
	}, pluginInfo)

	t.Run("error while decoding", func(t *testing.T) {
		f := s4.S4ReportingPluginFactory{
			Logger: logger,
			ORM:    orm,
			ConfigDecoder: func([]byte) (*s4.PluginConfig, *types.ReportingPluginLimits, error) {
				return nil, nil, errors.New("some error")
			},
		}

		rpConfig := types.ReportingPluginConfig{
			OffchainConfig: make([]byte, 100),
		}
		plugin, _, err := f.NewReportingPlugin(rpConfig)
		require.ErrorContains(t, err, "some error")
		require.Nil(t, plugin)
	})
}
