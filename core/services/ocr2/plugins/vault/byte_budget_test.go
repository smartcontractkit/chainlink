package vault

import (
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

func TestComputePluginByteBudgets_MatchesDirectFormula(t *testing.T) {
	ctx := t.Context()
	cfg := makeReportingPluginConfig(t, 10, nil, nil, 100, 2000, 64, 64, 64, 10, 0)
	pl, err := initializePluginLimits(ctx, limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)
	n, f := 10, 3
	batch, err := cfg.MaxBatchSize.Limit(ctx)
	require.NoError(t, err)
	bh := ocr3_1types.BlobHandleMarshalledBytesUpperBound(n, f)
	wantObs := pl.MaxObservationBytes - 2*batch*bh - sortNonceWireBytes - observationsOuterProtoOverhead
	wantPrec := pl.MaxReportsPlusPrecursorBytes - outcomesOuterProtoOverhead

	gotObs, gotPrec, err := computePluginByteBudgets(ctx, cfg, pl.MaxObservationBytes, pl.MaxReportsPlusPrecursorBytes, n, f)
	require.NoError(t, err)
	require.Equal(t, wantObs, gotObs)
	require.Equal(t, wantPrec, gotPrec)
}

func TestComputePluginByteBudgets_InvalidN(t *testing.T) {
	ctx := t.Context()
	cfg := makeReportingPluginConfig(t, 10, nil, nil, 100, 2000, 64, 64, 64, 10, 0)
	pl, err := initializePluginLimits(ctx, limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)
	_, _, err = computePluginByteBudgets(ctx, cfg, pl.MaxObservationBytes, pl.MaxReportsPlusPrecursorBytes, 0, 0)
	require.Error(t, err)
}

func TestComputePluginByteBudgets_ObservationBudgetDecreasesWithLargerDON(t *testing.T) {
	ctx := t.Context()
	cfg := makeReportingPluginConfig(t, 10, nil, nil, 100, 2000, 64, 64, 64, 10, 0)
	pl, err := initializePluginLimits(ctx, limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)

	obsSmall, precSmall, err := computePluginByteBudgets(ctx, cfg, pl.MaxObservationBytes, pl.MaxReportsPlusPrecursorBytes, 4, 1)
	require.NoError(t, err)
	obsLarge, precLarge, err := computePluginByteBudgets(ctx, cfg, pl.MaxObservationBytes, pl.MaxReportsPlusPrecursorBytes, 31, 10)
	require.NoError(t, err)
	require.Less(t, obsLarge, obsSmall, "larger N/F increases marshalled blob handle size reserved from MaxObservationBytes")
	require.Equal(t, precSmall, precLarge, "precursor budget must not depend on N/F")
}

func TestApplyTestByteBudgetsMatchesComputePluginByteBudgets(t *testing.T) {
	ctx := t.Context()
	cfg := makeReportingPluginConfig(t, 10, nil, nil, 100, 2000, 64, 64, 64, 10, 0)
	pl, err := initializePluginLimits(ctx, limits.Factory{Settings: cresettings.DefaultGetter})
	require.NoError(t, err)
	applyTestByteBudgets(t, cfg, ocr3types.ReportingPluginConfig{N: 10, F: 3})
	wantObs, wantPrec, err := computePluginByteBudgets(ctx, cfg, pl.MaxObservationBytes, pl.MaxReportsPlusPrecursorBytes, 10, 3)
	require.NoError(t, err)
	require.Equal(t, wantObs, cfg.ObsArrayBudgetBytes)
	require.Equal(t, wantPrec, cfg.PrecursorArrayBudgetBytes)
}
