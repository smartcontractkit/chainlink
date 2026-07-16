package beholderwrapper

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_WrapperFactory(t *testing.T) {
	validFactory := NewReportingPluginFactory[uint](
		&fakeFactory[uint]{},
		logger.TestLogger(t),
		"plugin",
	)
	failingFactory := NewReportingPluginFactory[uint](
		&fakeFactory[uint]{err: errors.New("error")},
		logger.TestLogger(t),
		"plugin",
	)

	plugin, _, err := validFactory.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{})
	require.NoError(t, err)

	// Verify the wrapped plugin works
	_, err = plugin.Outcome(t.Context(), ocr3types.OutcomeContext{}, nil, nil)
	require.NoError(t, err)

	_, _, err = failingFactory.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{})
	require.Error(t, err)
}

func Test_MetricViews(t *testing.T) {
	views := MetricViews()
	require.Len(t, views, 2)
}

type fakeFactory[RI any] struct {
	err error
}

func (f *fakeFactory[RI]) NewReportingPlugin(context.Context, ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[RI], ocr3types.ReportingPluginInfo, error) {
	if f.err != nil {
		return nil, ocr3types.ReportingPluginInfo{}, f.err
	}
	return &fakePlugin[RI]{}, ocr3types.ReportingPluginInfo{}, nil
}
