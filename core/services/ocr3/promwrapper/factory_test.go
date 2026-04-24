package promwrapper

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_WrapperFactory(t *testing.T) {
	validFactory := NewReportingPluginFactory(
		fakeFactory[uint]{},
		logger.TestLogger(t),
		"solana",
		"5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d",
		"plugin",
	)
	failingFactory := NewReportingPluginFactory(
		fakeFactory[uint]{err: errors.New("error")},
		logger.TestLogger(t),
		"evm",
		"11155111",
		"plugin",
	)

	plugin, _, err := validFactory.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{})
	require.NoError(t, err)

	_, err = plugin.Outcome(t.Context(), ocr3types.OutcomeContext{}, nil, nil)
	require.NoError(t, err)

	require.Equal(t, 1, counterFromHistogramByLabels(t, promOCR3Durations, "solana", "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d", "plugin", "outcome", "true"))
	require.Equal(t, 0, counterFromHistogramByLabels(t, promOCR3Durations, "solana", "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d", "plugin", "outcome", "false"))

	_, _, err = failingFactory.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{})
	require.Error(t, err)
}

type fakeFactory[RI any] struct {
	err error
}

func (f fakeFactory[RI]) NewReportingPlugin(context.Context, ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[RI], ocr3types.ReportingPluginInfo, error) {
	if f.err != nil {
		return nil, ocr3types.ReportingPluginInfo{}, f.err
	}
	return fakePlugin[RI]{}, ocr3types.ReportingPluginInfo{}, nil
}
