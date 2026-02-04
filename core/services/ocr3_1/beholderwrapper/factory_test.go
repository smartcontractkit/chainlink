package beholderwrapper

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_WrapperFactory(t *testing.T) {
	validFactory := NewReportingPluginFactory(
		&fakeFactory[uint]{},
		logger.TestLogger(t),
		"plugin",
	)
	failingFactory := NewReportingPluginFactory(
		&fakeFactory[uint]{err: errors.New("error")},
		logger.TestLogger(t),
		"plugin",
	)

	plugin, _, err := validFactory.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{}, nil)
	require.NoError(t, err)

	_, err = plugin.StateTransition(t.Context(), 1, ocrtypes.AttributedQuery{}, nil, nil, nil)
	require.NoError(t, err)

	_, _, err = failingFactory.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{}, nil)
	require.Error(t, err)
}

type fakeFactory[RI any] struct {
	err error
}

func (f *fakeFactory[RI]) NewReportingPlugin(context.Context, ocr3types.ReportingPluginConfig, ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[RI], ocr3_1types.ReportingPluginInfo, error) {
	if f.err != nil {
		return nil, nil, f.err
	}
	return &fakePlugin[RI]{}, ocr3_1types.ReportingPluginInfo1{}, nil
}
