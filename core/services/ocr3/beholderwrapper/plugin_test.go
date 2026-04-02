package beholderwrapper

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/beholderwrapper/metrics"
)

func Test_ReportingPlugin_WrapsAllMethods(t *testing.T) {
	m, err := metrics.NewPluginMetrics(MetricPrefix, "test", "abc")
	require.NoError(t, err)

	plugin := newReportingPlugin[uint](
		&fakePlugin[uint]{reports: make([]ocr3types.ReportPlus[uint], 2), observationSize: 5, outcomeSize: 3},
		m,
	)

	// Test Query
	_, err = plugin.Query(t.Context(), ocr3types.OutcomeContext{})
	require.NoError(t, err)

	// Test Observation
	obs, err := plugin.Observation(t.Context(), ocr3types.OutcomeContext{}, ocrtypes.Query{})
	require.NoError(t, err)
	require.Len(t, obs, 5)

	// Test ValidateObservation
	err = plugin.ValidateObservation(t.Context(), ocr3types.OutcomeContext{}, ocrtypes.Query{}, ocrtypes.AttributedObservation{})
	require.NoError(t, err)

	// Test ObservationQuorum
	quorum, err := plugin.ObservationQuorum(t.Context(), ocr3types.OutcomeContext{}, ocrtypes.Query{}, nil)
	require.NoError(t, err)
	require.True(t, quorum)

	// Test Outcome
	outcome, err := plugin.Outcome(t.Context(), ocr3types.OutcomeContext{}, ocrtypes.Query{}, nil)
	require.NoError(t, err)
	require.Len(t, outcome, 3)

	// Test Reports
	reports, err := plugin.Reports(t.Context(), 1, nil)
	require.NoError(t, err)
	require.Len(t, reports, 2)

	// Test ShouldAcceptAttestedReport
	accept, err := plugin.ShouldAcceptAttestedReport(t.Context(), 1, ocr3types.ReportWithInfo[uint]{})
	require.NoError(t, err)
	require.True(t, accept)

	// Test ShouldTransmitAcceptedReport
	transmit, err := plugin.ShouldTransmitAcceptedReport(t.Context(), 1, ocr3types.ReportWithInfo[uint]{})
	require.NoError(t, err)
	require.True(t, transmit)

	// Test Close
	err = plugin.Close()
	require.NoError(t, err)
}

func Test_ReportingPlugin_PropagatesErrors(t *testing.T) {
	m, err := metrics.NewPluginMetrics(MetricPrefix, "test", "abc")
	require.NoError(t, err)

	expectedErr := errors.New("test error")
	plugin := newReportingPlugin[uint](
		&fakePlugin[uint]{err: expectedErr},
		m,
	)

	_, err = plugin.Query(t.Context(), ocr3types.OutcomeContext{})
	require.ErrorIs(t, err, expectedErr)

	_, err = plugin.Observation(t.Context(), ocr3types.OutcomeContext{}, ocrtypes.Query{})
	require.ErrorIs(t, err, expectedErr)

	err = plugin.ValidateObservation(t.Context(), ocr3types.OutcomeContext{}, ocrtypes.Query{}, ocrtypes.AttributedObservation{})
	require.ErrorIs(t, err, expectedErr)

	_, err = plugin.ObservationQuorum(t.Context(), ocr3types.OutcomeContext{}, ocrtypes.Query{}, nil)
	require.ErrorIs(t, err, expectedErr)

	_, err = plugin.Outcome(t.Context(), ocr3types.OutcomeContext{}, ocrtypes.Query{}, nil)
	require.ErrorIs(t, err, expectedErr)

	_, err = plugin.Reports(t.Context(), 1, nil)
	require.ErrorIs(t, err, expectedErr)

	_, err = plugin.ShouldAcceptAttestedReport(t.Context(), 1, ocr3types.ReportWithInfo[uint]{})
	require.ErrorIs(t, err, expectedErr)

	_, err = plugin.ShouldTransmitAcceptedReport(t.Context(), 1, ocr3types.ReportWithInfo[uint]{})
	require.ErrorIs(t, err, expectedErr)

	err = plugin.Close()
	require.ErrorIs(t, err, expectedErr)
}

func Test_BoolToInt(t *testing.T) {
	require.Equal(t, 1, boolToInt(true))
	require.Equal(t, 0, boolToInt(false))
}

type fakePlugin[RI any] struct {
	reports         []ocr3types.ReportPlus[RI]
	observationSize int
	outcomeSize     int
	err             error
}

func (f *fakePlugin[RI]) Query(context.Context, ocr3types.OutcomeContext) (ocrtypes.Query, error) {
	if f.err != nil {
		return nil, f.err
	}
	return ocrtypes.Query{}, nil
}

func (f *fakePlugin[RI]) Observation(context.Context, ocr3types.OutcomeContext, ocrtypes.Query) (ocrtypes.Observation, error) {
	if f.err != nil {
		return nil, f.err
	}
	return make([]byte, f.observationSize), nil
}

func (f *fakePlugin[RI]) ValidateObservation(context.Context, ocr3types.OutcomeContext, ocrtypes.Query, ocrtypes.AttributedObservation) error {
	return f.err
}

func (f *fakePlugin[RI]) ObservationQuorum(context.Context, ocr3types.OutcomeContext, ocrtypes.Query, []ocrtypes.AttributedObservation) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return true, nil
}

func (f *fakePlugin[RI]) Outcome(context.Context, ocr3types.OutcomeContext, ocrtypes.Query, []ocrtypes.AttributedObservation) (ocr3types.Outcome, error) {
	if f.err != nil {
		return nil, f.err
	}
	return make([]byte, f.outcomeSize), nil
}

func (f *fakePlugin[RI]) Reports(context.Context, uint64, ocr3types.Outcome) ([]ocr3types.ReportPlus[RI], error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.reports, nil
}

func (f *fakePlugin[RI]) ShouldAcceptAttestedReport(context.Context, uint64, ocr3types.ReportWithInfo[RI]) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return true, nil
}

func (f *fakePlugin[RI]) ShouldTransmitAcceptedReport(context.Context, uint64, ocr3types.ReportWithInfo[RI]) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return true, nil
}

func (f *fakePlugin[RI]) Close() error {
	return f.err
}
