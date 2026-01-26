package beholderwrapper

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func Test_ReportingPlugin_WrapsAllMethods(t *testing.T) {
	metrics, err := newPluginMetrics("test", "abc")
	require.NoError(t, err)

	plugin := newReportingPlugin(
		&fakePlugin[uint]{reports: make([]ocr3types.ReportPlus[uint], 2), observationSize: 5, stateTransitionSize: 3},
		metrics,
	)

	// Test Query
	_, err = plugin.Query(t.Context(), 1, nil, nil)
	require.NoError(t, err)

	// Test Observation
	obs, err := plugin.Observation(t.Context(), 1, ocrtypes.AttributedQuery{}, nil, nil)
	require.NoError(t, err)
	require.Len(t, obs, 5)

	// Test ValidateObservation
	err = plugin.ValidateObservation(t.Context(), 1, ocrtypes.AttributedQuery{}, ocrtypes.AttributedObservation{}, nil, nil)
	require.NoError(t, err)

	// Test ObservationQuorum
	quorum, err := plugin.ObservationQuorum(t.Context(), 1, ocrtypes.AttributedQuery{}, nil, nil, nil)
	require.NoError(t, err)
	require.True(t, quorum)

	// Test StateTransition
	st, err := plugin.StateTransition(t.Context(), 1, ocrtypes.AttributedQuery{}, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, st, 3)

	// Test Committed
	err = plugin.Committed(t.Context(), 1, nil)
	require.NoError(t, err)

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
	metrics, err := newPluginMetrics("test", "abc")
	require.NoError(t, err)

	expectedErr := errors.New("test error")
	plugin := newReportingPlugin(
		&fakePlugin[uint]{err: expectedErr},
		metrics,
	)

	_, err = plugin.Query(t.Context(), 1, nil, nil)
	require.ErrorIs(t, err, expectedErr)

	_, err = plugin.Observation(t.Context(), 1, ocrtypes.AttributedQuery{}, nil, nil)
	require.ErrorIs(t, err, expectedErr)

	err = plugin.ValidateObservation(t.Context(), 1, ocrtypes.AttributedQuery{}, ocrtypes.AttributedObservation{}, nil, nil)
	require.ErrorIs(t, err, expectedErr)

	_, err = plugin.ObservationQuorum(t.Context(), 1, ocrtypes.AttributedQuery{}, nil, nil, nil)
	require.ErrorIs(t, err, expectedErr)

	_, err = plugin.StateTransition(t.Context(), 1, ocrtypes.AttributedQuery{}, nil, nil, nil)
	require.ErrorIs(t, err, expectedErr)

	err = plugin.Committed(t.Context(), 1, nil)
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

type fakePlugin[RI any] struct {
	reports             []ocr3types.ReportPlus[RI]
	observationSize     int
	stateTransitionSize int
	err                 error
}

func (f *fakePlugin[RI]) Query(context.Context, uint64, ocr3_1types.KeyValueStateReader, ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Query, error) {
	if f.err != nil {
		return nil, f.err
	}
	return ocrtypes.Query{}, nil
}

func (f *fakePlugin[RI]) Observation(context.Context, uint64, ocrtypes.AttributedQuery, ocr3_1types.KeyValueStateReader, ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Observation, error) {
	if f.err != nil {
		return nil, f.err
	}
	return make([]byte, f.observationSize), nil
}

func (f *fakePlugin[RI]) ValidateObservation(context.Context, uint64, ocrtypes.AttributedQuery, ocrtypes.AttributedObservation, ocr3_1types.KeyValueStateReader, ocr3_1types.BlobFetcher) error {
	return f.err
}

func (f *fakePlugin[RI]) ObservationQuorum(context.Context, uint64, ocrtypes.AttributedQuery, []ocrtypes.AttributedObservation, ocr3_1types.KeyValueStateReader, ocr3_1types.BlobFetcher) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return true, nil
}

func (f *fakePlugin[RI]) StateTransition(context.Context, uint64, ocrtypes.AttributedQuery, []ocrtypes.AttributedObservation, ocr3_1types.KeyValueStateReadWriter, ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	if f.err != nil {
		return nil, f.err
	}
	return make([]byte, f.stateTransitionSize), nil
}

func (f *fakePlugin[RI]) Committed(context.Context, uint64, ocr3_1types.KeyValueStateReader) error {
	return f.err
}

func (f *fakePlugin[RI]) Reports(context.Context, uint64, ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[RI], error) {
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
