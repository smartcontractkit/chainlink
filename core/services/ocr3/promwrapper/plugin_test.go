package promwrapper

import (
	"context"
	"errors"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func Test_ReportsGeneratedGauge(t *testing.T) {
	pluginObservationSize := 5
	pluginOutcomeSize := 3

	plugin1 := newReportingPlugin(
		fakePlugin[uint]{reports: make([]ocr3types.ReportPlus[uint], 2)},
		"123", "empty", "abc", promOCR3ReportsGenerated, promOCR3Durations, promOCR3Sizes, promOCR3PluginStatus,
	)
	plugin2 := newReportingPlugin(
		fakePlugin[bool]{reports: make([]ocr3types.ReportPlus[bool], 10), observationSize: pluginObservationSize, outcomeSize: pluginOutcomeSize},
		"solana", "different_plugin", "abc", promOCR3ReportsGenerated, promOCR3Durations, promOCR3Sizes, promOCR3PluginStatus,
	)
	plugin3 := newReportingPlugin(
		fakePlugin[string]{err: errors.New("error")},
		"1234", "empty", "abc", promOCR3ReportsGenerated, promOCR3Durations, promOCR3Sizes, promOCR3PluginStatus,
	)

	r1, err := plugin1.Reports(t.Context(), 1, nil)
	require.NoError(t, err)
	require.Len(t, r1, 2)

	for i := 0; i < 10; i++ {
		r2, err1 := plugin2.Reports(t.Context(), 1, nil)
		require.NoError(t, err1)
		require.Len(t, r2, 10)
	}

	_, err = plugin2.ShouldAcceptAttestedReport(t.Context(), 1, ocr3types.ReportWithInfo[bool]{})
	require.NoError(t, err)

	_, err = plugin3.Reports(t.Context(), 1, nil)
	require.Error(t, err)

	g1 := testutil.ToFloat64(promOCR3ReportsGenerated.WithLabelValues("123", "empty", "reports"))
	require.Equal(t, 2, int(g1))

	g2 := testutil.ToFloat64(promOCR3ReportsGenerated.WithLabelValues("solana", "different_plugin", "reports"))
	require.Equal(t, 100, int(g2))

	g3 := testutil.ToFloat64(promOCR3ReportsGenerated.WithLabelValues("solana", "different_plugin", "shouldAccept"))
	require.Equal(t, 1, int(g3))

	g4 := testutil.ToFloat64(promOCR3ReportsGenerated.WithLabelValues("1234", "empty", "reports"))
	require.Equal(t, 0, int(g4))

	pluginHealth := testutil.ToFloat64(promOCR3PluginStatus.WithLabelValues("123", "empty", "abc"))
	require.Equal(t, 1, int(pluginHealth))

	require.NoError(t, plugin1.Close())
	pluginHealth = testutil.ToFloat64(promOCR3PluginStatus.WithLabelValues("123", "empty", "abc"))
	require.Equal(t, 0, int(pluginHealth))

	iterations := 10
	for i := 0; i < iterations; i++ {
		_, err1 := plugin2.Outcome(t.Context(), ocr3types.OutcomeContext{}, nil, nil)
		require.NoError(t, err1)
	}
	_, err1 := plugin2.Observation(t.Context(), ocr3types.OutcomeContext{}, nil)
	require.NoError(t, err1)

	outcomesLen := testutil.ToFloat64(promOCR3Sizes.WithLabelValues("solana", "different_plugin", "outcome"))
	require.Equal(t, pluginOutcomeSize*iterations, int(outcomesLen))
	observationLen := testutil.ToFloat64(promOCR3Sizes.WithLabelValues("solana", "different_plugin", "observation"))
	require.Equal(t, pluginObservationSize, int(observationLen))
}

func Test_DurationHistograms(t *testing.T) {
	plugin1 := newReportingPlugin(
		fakePlugin[uint]{},
		"123", "empty", "abc", promOCR3ReportsGenerated, promOCR3Durations, promOCR3Sizes, promOCR3PluginStatus,
	)
	plugin2 := newReportingPlugin(
		fakePlugin[uint]{err: errors.New("error")},
		"123", "empty", "abc", promOCR3ReportsGenerated, promOCR3Durations, promOCR3Sizes, promOCR3PluginStatus,
	)
	plugin3 := newReportingPlugin(
		fakePlugin[uint]{},
		"solana", "commit", "abc", promOCR3ReportsGenerated, promOCR3Durations, promOCR3Sizes, promOCR3PluginStatus,
	)

	for _, p := range []*reportingPlugin[uint]{plugin1, plugin2, plugin3} {
		_, _ = p.Query(t.Context(), ocr3types.OutcomeContext{})
		for i := 0; i < 2; i++ {
			_, _ = p.Observation(t.Context(), ocr3types.OutcomeContext{}, nil)
		}
		_ = p.ValidateObservation(t.Context(), ocr3types.OutcomeContext{}, nil, ocrtypes.AttributedObservation{})
		_, _ = p.Outcome(t.Context(), ocr3types.OutcomeContext{}, nil, nil)
		_, _ = p.Reports(t.Context(), 0, nil)
		_, _ = p.ShouldAcceptAttestedReport(t.Context(), 0, ocr3types.ReportWithInfo[uint]{})
		_, _ = p.ShouldTransmitAcceptedReport(t.Context(), 0, ocr3types.ReportWithInfo[uint]{})
	}

	require.Equal(t, 1, counterFromHistogramByLabels(t, promOCR3Durations, "123", "empty", "query", "true"))
	require.Equal(t, 1, counterFromHistogramByLabels(t, promOCR3Durations, "123", "empty", "query", "false"))
	require.Equal(t, 1, counterFromHistogramByLabels(t, promOCR3Durations, "solana", "commit", "query", "true"))

	require.Equal(t, 2, counterFromHistogramByLabels(t, promOCR3Durations, "123", "empty", "observation", "true"))
	require.Equal(t, 2, counterFromHistogramByLabels(t, promOCR3Durations, "123", "empty", "observation", "false"))
	require.Equal(t, 2, counterFromHistogramByLabels(t, promOCR3Durations, "solana", "commit", "observation", "true"))
}

type fakePlugin[RI any] struct {
	reports         []ocr3types.ReportPlus[RI]
	observationSize int
	outcomeSize     int
	err             error
}

func (f fakePlugin[RI]) Query(context.Context, ocr3types.OutcomeContext) (ocrtypes.Query, error) {
	if f.err != nil {
		return nil, f.err
	}
	return ocrtypes.Query{}, nil
}

func (f fakePlugin[RI]) Observation(context.Context, ocr3types.OutcomeContext, ocrtypes.Query) (ocrtypes.Observation, error) {
	if f.err != nil {
		return nil, f.err
	}
	return make([]byte, f.observationSize), nil
}

func (f fakePlugin[RI]) ValidateObservation(context.Context, ocr3types.OutcomeContext, ocrtypes.Query, ocrtypes.AttributedObservation) error {
	return f.err
}

func (f fakePlugin[RI]) ObservationQuorum(context.Context, ocr3types.OutcomeContext, ocrtypes.Query, []ocrtypes.AttributedObservation) (quorumReached bool, err error) {
	return false, f.err
}

func (f fakePlugin[RI]) Outcome(context.Context, ocr3types.OutcomeContext, ocrtypes.Query, []ocrtypes.AttributedObservation) (ocr3types.Outcome, error) {
	if f.err != nil {
		return nil, f.err
	}
	return make([]byte, f.outcomeSize), nil
}

func (f fakePlugin[RI]) Reports(context.Context, uint64, ocr3types.Outcome) ([]ocr3types.ReportPlus[RI], error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.reports, nil
}

func (f fakePlugin[RI]) ShouldAcceptAttestedReport(context.Context, uint64, ocr3types.ReportWithInfo[RI]) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return true, nil
}

func (f fakePlugin[RI]) ShouldTransmitAcceptedReport(context.Context, uint64, ocr3types.ReportWithInfo[RI]) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return true, nil
}

func (f fakePlugin[RI]) Close() error {
	return f.err
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
