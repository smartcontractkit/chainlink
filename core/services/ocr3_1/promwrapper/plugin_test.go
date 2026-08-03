package promwrapper

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Test_Plugin_FunctionLabels guards the OCR3.1 metric labelling: the three
// OCR3.1-only phases must be recorded under their own distinct function labels
// (observationQuorum, stateTransition, committed) and NOT conflated with one
// another. Each method is invoked once and the per-label duration histogram
// sample count is asserted to increment by exactly one.
func Test_Plugin_FunctionLabels(t *testing.T) {
	t.Parallel()
	const (
		fam  = "evm"
		id   = "1"
		plug = "llo"
	)
	funcs := []functionType{query, observation, validateObservation, observationQuorum, stateTransition, committed, reports, shouldAccept, shouldTransmit}

	init := map[functionType]int{}
	for _, f := range funcs {
		init[f] = counterFromHistogramByLabels(t, promOCR3Durations, fam, id, plug, string(f), "true")
	}

	p := newReportingPlugin(
		fakePlugin[uint]{reports: make([]ocr3types.ReportPlus[uint], 2), stateTransitionSize: 4},
		fam, id, plug, "abc",
		promOCR3ReportsGenerated, promOCR3Durations, promOCR3Sizes, promOCR3PluginStatus,
	)

	ctx := t.Context()
	_, err := p.Query(ctx, 1, nil, nil)
	require.NoError(t, err)
	_, err = p.Observation(ctx, 1, ocrtypes.AttributedQuery{}, nil, nil)
	require.NoError(t, err)
	require.NoError(t, p.ValidateObservation(ctx, 1, ocrtypes.AttributedQuery{}, ocrtypes.AttributedObservation{}, nil, nil))
	_, err = p.ObservationQuorum(ctx, 1, ocrtypes.AttributedQuery{}, nil, nil, nil)
	require.NoError(t, err)
	_, err = p.StateTransition(ctx, 1, ocrtypes.AttributedQuery{}, nil, nil, nil)
	require.NoError(t, err)
	require.NoError(t, p.Committed(ctx, 1, nil))
	_, err = p.Reports(ctx, 1, nil)
	require.NoError(t, err)
	_, err = p.ShouldAcceptAttestedReport(ctx, 1, ocr3types.ReportWithInfo[uint]{})
	require.NoError(t, err)
	_, err = p.ShouldTransmitAcceptedReport(ctx, 1, ocr3types.ReportWithInfo[uint]{})
	require.NoError(t, err)

	// Every phase recorded exactly one duration sample under its own label.
	for _, f := range funcs {
		got := counterFromHistogramByLabels(t, promOCR3Durations, fam, id, plug, string(f), "true") - init[f]
		require.Equalf(t, 1, got, "duration label %q should increment once", f)
	}

	// StateTransition precursor size tracked under the stateTransition label.
	require.Equal(t, 4, int(testutil.ToFloat64(promOCR3Sizes.WithLabelValues(fam, id, plug, string(stateTransition)))))
	// Reports counter under the reports label.
	require.Equal(t, 2, int(testutil.ToFloat64(promOCR3ReportsGenerated.WithLabelValues(fam, id, plug, string(reports)))))
}

// Test_Factory covers NewReportingPluginFactory + NewReportingPlugin: the
// factory wraps the origin plugin and the wrapper reports metrics.
func Test_Factory(t *testing.T) {
	t.Parallel()
	factory := NewReportingPluginFactory(
		fakeFactory[uint]{plugin: fakePlugin[uint]{}},
		logger.TestLogger(t), "aptos", "1", "llo",
	)

	cd := ocrtypes.ConfigDigest{9}
	p, info, err := factory.NewReportingPlugin(t.Context(), ocr3types.ReportingPluginConfig{ConfigDigest: cd}, nil)
	require.NoError(t, err)
	require.NotNil(t, p)
	require.NotNil(t, info)

	_, err = p.Query(t.Context(), 1, nil, nil)
	require.NoError(t, err)
	require.Equal(t, 1, int(testutil.ToFloat64(promOCR3PluginStatus.WithLabelValues("aptos", "1", "llo", cd.String()))))

	require.NoError(t, p.Close())
	require.Equal(t, 0, int(testutil.ToFloat64(promOCR3PluginStatus.WithLabelValues("aptos", "1", "llo", cd.String()))))
}

func counterFromHistogramByLabels(t *testing.T, histogramVec *prometheus.HistogramVec, labels ...string) int {
	observer, err := histogramVec.GetMetricWithLabelValues(labels...)
	require.NoError(t, err)

	metricCh := make(chan prometheus.Metric, 1)
	observer.(prometheus.Histogram).Collect(metricCh)
	close(metricCh)

	metric := <-metricCh
	pb := &io_prometheus_client.Metric{}
	err = metric.Write(pb)
	require.NoError(t, err)

	//nolint:gosec // we don't care about that in tests
	return int(pb.GetHistogram().GetSampleCount())
}

type fakeFactory[RI any] struct{ plugin fakePlugin[RI] }

func (f fakeFactory[RI]) NewReportingPlugin(context.Context, ocr3types.ReportingPluginConfig, ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[RI], ocr3_1types.ReportingPluginInfo, error) {
	return f.plugin, ocr3_1types.ReportingPluginInfo1{}, nil
}

type fakePlugin[RI any] struct {
	reports             []ocr3types.ReportPlus[RI]
	observationSize     int
	stateTransitionSize int
	err                 error
}

func (f fakePlugin[RI]) Query(context.Context, uint64, ocr3_1types.KeyValueStateReader, ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Query, error) {
	return ocrtypes.Query{}, f.err
}

func (f fakePlugin[RI]) Observation(context.Context, uint64, ocrtypes.AttributedQuery, ocr3_1types.KeyValueStateReader, ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Observation, error) {
	if f.err != nil {
		return nil, f.err
	}
	return make([]byte, f.observationSize), nil
}

func (f fakePlugin[RI]) ValidateObservation(context.Context, uint64, ocrtypes.AttributedQuery, ocrtypes.AttributedObservation, ocr3_1types.KeyValueStateReader, ocr3_1types.BlobFetcher) error {
	return f.err
}

func (f fakePlugin[RI]) ObservationQuorum(context.Context, uint64, ocrtypes.AttributedQuery, []ocrtypes.AttributedObservation, ocr3_1types.KeyValueStateReader, ocr3_1types.BlobFetcher) (bool, error) {
	return false, f.err
}

func (f fakePlugin[RI]) StateTransition(context.Context, uint64, ocrtypes.AttributedQuery, []ocrtypes.AttributedObservation, ocr3_1types.KeyValueStateReadWriter, ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	if f.err != nil {
		return nil, f.err
	}
	return make([]byte, f.stateTransitionSize), nil
}

func (f fakePlugin[RI]) Committed(context.Context, uint64, ocr3_1types.KeyValueStateReader) error {
	return f.err
}

func (f fakePlugin[RI]) Reports(context.Context, uint64, ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[RI], error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.reports, nil
}

func (f fakePlugin[RI]) ShouldAcceptAttestedReport(context.Context, uint64, ocr3types.ReportWithInfo[RI]) (bool, error) {
	return true, f.err
}

func (f fakePlugin[RI]) ShouldTransmitAcceptedReport(context.Context, uint64, ocr3types.ReportWithInfo[RI]) (bool, error) {
	return true, f.err
}

func (f fakePlugin[RI]) Close() error { return f.err }
