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

func Test_InstrumentedBlobBroadcastFetcher(t *testing.T) {
	metrics, err := newPluginMetrics("test", "abc")
	require.NoError(t, err)

	inner := &fakeBlobBroadcastFetcher{
		broadcastPayload: []byte("broadcast-handle"),
		fetchPayload:     []byte("fetched-data"),
	}

	wrapped := &instrumentedBlobBroadcastFetcher{
		inner:   inner,
		metrics: metrics,
		instrumentedBlobFetcher: instrumentedBlobFetcher{
			inner:   inner,
			metrics: metrics,
		},
	}

	// BroadcastBlob delegates and records metrics
	handle, err := wrapped.BroadcastBlob(t.Context(), []byte("payload"), ocr3_1types.BlobExpirationHintSequenceNumber{SeqNr: 1})
	require.NoError(t, err)
	require.Equal(t, 1, inner.broadcastCalls)

	// FetchBlob delegates and records metrics
	data, err := wrapped.FetchBlob(t.Context(), handle)
	require.NoError(t, err)
	require.Equal(t, []byte("fetched-data"), data)
	require.Equal(t, 1, inner.fetchCalls)
}

func Test_InstrumentedBlobBroadcastFetcher_PropagatesErrors(t *testing.T) {
	metrics, err := newPluginMetrics("test", "abc")
	require.NoError(t, err)

	expectedErr := errors.New("blob error")
	inner := &fakeBlobBroadcastFetcher{err: expectedErr}
	wrapped := &instrumentedBlobBroadcastFetcher{
		inner:   inner,
		metrics: metrics,
		instrumentedBlobFetcher: instrumentedBlobFetcher{
			inner:   inner,
			metrics: metrics,
		},
	}

	_, err = wrapped.BroadcastBlob(t.Context(), []byte("payload"), ocr3_1types.BlobExpirationHintSequenceNumber{SeqNr: 1})
	require.ErrorIs(t, err, expectedErr)

	_, err = wrapped.FetchBlob(t.Context(), ocr3_1types.BlobHandle{})
	require.ErrorIs(t, err, expectedErr)
}

func Test_InstrumentedBlobFetcher(t *testing.T) {
	metrics, err := newPluginMetrics("test", "abc")
	require.NoError(t, err)

	inner := &fakeBlobFetcher{fetchPayload: []byte("fetched-data")}
	wrapped := &instrumentedBlobFetcher{inner: inner, metrics: metrics}

	data, err := wrapped.FetchBlob(t.Context(), ocr3_1types.BlobHandle{})
	require.NoError(t, err)
	require.Equal(t, []byte("fetched-data"), data)
	require.Equal(t, 1, inner.fetchCalls)
}

func Test_InstrumentedBlobFetcher_PropagatesErrors(t *testing.T) {
	metrics, err := newPluginMetrics("test", "abc")
	require.NoError(t, err)

	expectedErr := errors.New("fetch error")
	inner := &fakeBlobFetcher{err: expectedErr}
	wrapped := &instrumentedBlobFetcher{inner: inner, metrics: metrics}

	_, err = wrapped.FetchBlob(t.Context(), ocr3_1types.BlobHandle{})
	require.ErrorIs(t, err, expectedErr)
}

func Test_ReportingPlugin_WrapsBlobs(t *testing.T) {
	metrics, err := newPluginMetrics("test", "abc")
	require.NoError(t, err)

	innerBlob := &fakeBlobBroadcastFetcher{
		broadcastPayload: []byte("handle"),
		fetchPayload:     []byte("data"),
	}
	innerFetcher := &fakeBlobFetcher{fetchPayload: []byte("data")}

	capturingPlugin := &blobCapturingPlugin[uint]{}
	plugin := newReportingPlugin[uint](capturingPlugin, metrics)

	// Query wraps BlobBroadcastFetcher
	_, _ = plugin.Query(t.Context(), 1, nil, innerBlob)
	require.IsType(t, &instrumentedBlobBroadcastFetcher{}, capturingPlugin.lastBroadcastFetcher)

	// Observation wraps BlobBroadcastFetcher
	_, _ = plugin.Observation(t.Context(), 1, ocrtypes.AttributedQuery{}, nil, innerBlob)
	require.IsType(t, &instrumentedBlobBroadcastFetcher{}, capturingPlugin.lastBroadcastFetcher)

	// ValidateObservation wraps BlobFetcher
	_ = plugin.ValidateObservation(t.Context(), 1, ocrtypes.AttributedQuery{}, ocrtypes.AttributedObservation{}, nil, innerFetcher)
	require.IsType(t, &instrumentedBlobFetcher{}, capturingPlugin.lastFetcher)

	// ObservationQuorum wraps BlobFetcher
	_, _ = plugin.ObservationQuorum(t.Context(), 1, ocrtypes.AttributedQuery{}, nil, nil, innerFetcher)
	require.IsType(t, &instrumentedBlobFetcher{}, capturingPlugin.lastFetcher)

	// StateTransition wraps BlobFetcher
	_, _ = plugin.StateTransition(t.Context(), 1, ocrtypes.AttributedQuery{}, nil, nil, innerFetcher)
	require.IsType(t, &instrumentedBlobFetcher{}, capturingPlugin.lastFetcher)

	// nil is preserved
	_, _ = plugin.Query(t.Context(), 1, nil, nil)
	require.Nil(t, capturingPlugin.lastBroadcastFetcher)

	_ = plugin.ValidateObservation(t.Context(), 1, ocrtypes.AttributedQuery{}, ocrtypes.AttributedObservation{}, nil, nil)
	require.Nil(t, capturingPlugin.lastFetcher)
}

func Test_InstrumentedKVStateReader(t *testing.T) {
	metrics, err := newPluginMetrics("test", "abc")
	require.NoError(t, err)

	inner := &fakeKVStateReader{data: map[string][]byte{"key1": []byte("value1")}}
	wrapped := &instrumentedKVStateReader{inner: inner, ctx: t.Context(), metrics: metrics}

	data, err := wrapped.Read([]byte("key1"))
	require.NoError(t, err)
	require.Equal(t, []byte("value1"), data)
	require.Equal(t, 1, inner.readCalls)

	// Missing key returns nil
	data, err = wrapped.Read([]byte("missing"))
	require.NoError(t, err)
	require.Nil(t, data)
	require.Equal(t, 2, inner.readCalls)
}

func Test_InstrumentedKVStateReader_PropagatesErrors(t *testing.T) {
	metrics, err := newPluginMetrics("test", "abc")
	require.NoError(t, err)

	expectedErr := errors.New("read error")
	inner := &fakeKVStateReader{err: expectedErr}
	wrapped := &instrumentedKVStateReader{inner: inner, ctx: t.Context(), metrics: metrics}

	_, err = wrapped.Read([]byte("key"))
	require.ErrorIs(t, err, expectedErr)
}

func Test_InstrumentedKVStateReadWriter(t *testing.T) {
	metrics, err := newPluginMetrics("test", "abc")
	require.NoError(t, err)

	inner := &fakeKVStateReadWriter{fakeKVStateReader: fakeKVStateReader{data: map[string][]byte{}}}
	wrapped := &instrumentedKVStateReadWriter{
		instrumentedKVStateReader: instrumentedKVStateReader{inner: inner, ctx: t.Context(), metrics: metrics},
		writer:                    inner,
	}

	// Write
	err = wrapped.Write([]byte("key1"), []byte("value1"))
	require.NoError(t, err)
	require.Equal(t, 1, inner.writeCalls)

	// Read back through the wrapper
	data, err := wrapped.Read([]byte("key1"))
	require.NoError(t, err)
	require.Equal(t, []byte("value1"), data)

	// Delete
	err = wrapped.Delete([]byte("key1"))
	require.NoError(t, err)
	require.Equal(t, 1, inner.deleteCalls)

	// Read returns nil after delete
	data, err = wrapped.Read([]byte("key1"))
	require.NoError(t, err)
	require.Nil(t, data)
}

func Test_InstrumentedKVStateReadWriter_PropagatesErrors(t *testing.T) {
	metrics, err := newPluginMetrics("test", "abc")
	require.NoError(t, err)

	expectedErr := errors.New("write error")
	inner := &fakeKVStateReadWriter{fakeKVStateReader: fakeKVStateReader{err: expectedErr}}
	wrapped := &instrumentedKVStateReadWriter{
		instrumentedKVStateReader: instrumentedKVStateReader{inner: inner, ctx: t.Context(), metrics: metrics},
		writer:                    inner,
	}

	_, err = wrapped.Read([]byte("key"))
	require.ErrorIs(t, err, expectedErr)

	err = wrapped.Write([]byte("key"), []byte("value"))
	require.ErrorIs(t, err, expectedErr)

	err = wrapped.Delete([]byte("key"))
	require.ErrorIs(t, err, expectedErr)
}

func Test_ReportingPlugin_WrapsKV(t *testing.T) {
	metrics, err := newPluginMetrics("test", "abc")
	require.NoError(t, err)

	innerReader := &fakeKVStateReader{data: map[string][]byte{}}
	innerReadWriter := &fakeKVStateReadWriter{fakeKVStateReader: fakeKVStateReader{data: map[string][]byte{}}}

	capturingPlugin := &kvCapturingPlugin[uint]{}
	plugin := newReportingPlugin[uint](capturingPlugin, metrics)

	// Query wraps KeyValueStateReader
	_, _ = plugin.Query(t.Context(), 1, innerReader, nil)
	require.IsType(t, &instrumentedKVStateReader{}, capturingPlugin.lastReader)

	// Observation wraps KeyValueStateReader
	_, _ = plugin.Observation(t.Context(), 1, ocrtypes.AttributedQuery{}, innerReader, nil)
	require.IsType(t, &instrumentedKVStateReader{}, capturingPlugin.lastReader)

	// ValidateObservation wraps KeyValueStateReader
	_ = plugin.ValidateObservation(t.Context(), 1, ocrtypes.AttributedQuery{}, ocrtypes.AttributedObservation{}, innerReader, nil)
	require.IsType(t, &instrumentedKVStateReader{}, capturingPlugin.lastReader)

	// ObservationQuorum wraps KeyValueStateReader
	_, _ = plugin.ObservationQuorum(t.Context(), 1, ocrtypes.AttributedQuery{}, nil, innerReader, nil)
	require.IsType(t, &instrumentedKVStateReader{}, capturingPlugin.lastReader)

	// StateTransition wraps KeyValueStateReadWriter
	_, _ = plugin.StateTransition(t.Context(), 1, ocrtypes.AttributedQuery{}, nil, innerReadWriter, nil)
	require.IsType(t, &instrumentedKVStateReadWriter{}, capturingPlugin.lastReadWriter)

	// Committed wraps KeyValueStateReader
	_ = plugin.Committed(t.Context(), 1, innerReader)
	require.IsType(t, &instrumentedKVStateReader{}, capturingPlugin.lastReader)

	// nil is preserved
	_, _ = plugin.Query(t.Context(), 1, nil, nil)
	require.Nil(t, capturingPlugin.lastReader)

	_, _ = plugin.StateTransition(t.Context(), 1, ocrtypes.AttributedQuery{}, nil, nil, nil)
	require.Nil(t, capturingPlugin.lastReadWriter)
}

type fakeKVStateReader struct {
	data      map[string][]byte
	err       error
	readCalls int
}

func (f *fakeKVStateReader) Read(key []byte) ([]byte, error) {
	f.readCalls++
	if f.err != nil {
		return nil, f.err
	}
	return f.data[string(key)], nil
}

type fakeKVStateReadWriter struct {
	fakeKVStateReader
	writeCalls  int
	deleteCalls int
}

func (f *fakeKVStateReadWriter) Write(key []byte, value []byte) error {
	f.writeCalls++
	if f.err != nil {
		return f.err
	}
	f.data[string(key)] = value
	return nil
}

func (f *fakeKVStateReadWriter) Delete(key []byte) error {
	f.deleteCalls++
	if f.err != nil {
		return f.err
	}
	delete(f.data, string(key))
	return nil
}

// kvCapturingPlugin captures the KV reader/writer it receives so tests can assert on wrapping.
type kvCapturingPlugin[RI any] struct {
	lastReader     ocr3_1types.KeyValueStateReader
	lastReadWriter ocr3_1types.KeyValueStateReadWriter
}

func (p *kvCapturingPlugin[RI]) Query(_ context.Context, _ uint64, r ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Query, error) {
	p.lastReader = r
	return ocrtypes.Query{}, nil
}

func (p *kvCapturingPlugin[RI]) Observation(_ context.Context, _ uint64, _ ocrtypes.AttributedQuery, r ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Observation, error) {
	p.lastReader = r
	return ocrtypes.Observation{}, nil
}

func (p *kvCapturingPlugin[RI]) ValidateObservation(_ context.Context, _ uint64, _ ocrtypes.AttributedQuery, _ ocrtypes.AttributedObservation, r ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobFetcher) error {
	p.lastReader = r
	return nil
}

func (p *kvCapturingPlugin[RI]) ObservationQuorum(_ context.Context, _ uint64, _ ocrtypes.AttributedQuery, _ []ocrtypes.AttributedObservation, r ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobFetcher) (bool, error) {
	p.lastReader = r
	return true, nil
}

func (p *kvCapturingPlugin[RI]) StateTransition(_ context.Context, _ uint64, _ ocrtypes.AttributedQuery, _ []ocrtypes.AttributedObservation, rw ocr3_1types.KeyValueStateReadWriter, _ ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	p.lastReadWriter = rw
	return nil, nil
}

func (p *kvCapturingPlugin[RI]) Committed(_ context.Context, _ uint64, r ocr3_1types.KeyValueStateReader) error {
	p.lastReader = r
	return nil
}

func (p *kvCapturingPlugin[RI]) Reports(context.Context, uint64, ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[RI], error) {
	return nil, nil
}

func (p *kvCapturingPlugin[RI]) ShouldAcceptAttestedReport(context.Context, uint64, ocr3types.ReportWithInfo[RI]) (bool, error) {
	return true, nil
}

func (p *kvCapturingPlugin[RI]) ShouldTransmitAcceptedReport(context.Context, uint64, ocr3types.ReportWithInfo[RI]) (bool, error) {
	return true, nil
}

func (p *kvCapturingPlugin[RI]) Close() error {
	return nil
}

type fakeBlobBroadcastFetcher struct {
	broadcastPayload []byte
	fetchPayload     []byte
	err              error
	broadcastCalls   int
	fetchCalls       int
}

func (f *fakeBlobBroadcastFetcher) BroadcastBlob(_ context.Context, _ []byte, _ ocr3_1types.BlobExpirationHint) (ocr3_1types.BlobHandle, error) {
	f.broadcastCalls++
	return ocr3_1types.BlobHandle{}, f.err
}

func (f *fakeBlobBroadcastFetcher) FetchBlob(_ context.Context, _ ocr3_1types.BlobHandle) ([]byte, error) {
	f.fetchCalls++
	if f.err != nil {
		return nil, f.err
	}
	return f.fetchPayload, nil
}

type fakeBlobFetcher struct {
	fetchPayload []byte
	err          error
	fetchCalls   int
}

func (f *fakeBlobFetcher) FetchBlob(_ context.Context, _ ocr3_1types.BlobHandle) ([]byte, error) {
	f.fetchCalls++
	if f.err != nil {
		return nil, f.err
	}
	return f.fetchPayload, nil
}

// blobCapturingPlugin captures the blob fetcher/broadcaster it receives so tests can assert on wrapping.
type blobCapturingPlugin[RI any] struct {
	lastBroadcastFetcher ocr3_1types.BlobBroadcastFetcher
	lastFetcher          ocr3_1types.BlobFetcher
}

func (p *blobCapturingPlugin[RI]) Query(_ context.Context, _ uint64, _ ocr3_1types.KeyValueStateReader, bbf ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Query, error) {
	p.lastBroadcastFetcher = bbf
	return ocrtypes.Query{}, nil
}

func (p *blobCapturingPlugin[RI]) Observation(_ context.Context, _ uint64, _ ocrtypes.AttributedQuery, _ ocr3_1types.KeyValueStateReader, bbf ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Observation, error) {
	p.lastBroadcastFetcher = bbf
	return ocrtypes.Observation{}, nil
}

func (p *blobCapturingPlugin[RI]) ValidateObservation(_ context.Context, _ uint64, _ ocrtypes.AttributedQuery, _ ocrtypes.AttributedObservation, _ ocr3_1types.KeyValueStateReader, bf ocr3_1types.BlobFetcher) error {
	p.lastFetcher = bf
	return nil
}

func (p *blobCapturingPlugin[RI]) ObservationQuorum(_ context.Context, _ uint64, _ ocrtypes.AttributedQuery, _ []ocrtypes.AttributedObservation, _ ocr3_1types.KeyValueStateReader, bf ocr3_1types.BlobFetcher) (bool, error) {
	p.lastFetcher = bf
	return true, nil
}

func (p *blobCapturingPlugin[RI]) StateTransition(_ context.Context, _ uint64, _ ocrtypes.AttributedQuery, _ []ocrtypes.AttributedObservation, _ ocr3_1types.KeyValueStateReadWriter, bf ocr3_1types.BlobFetcher) (ocr3_1types.ReportsPlusPrecursor, error) {
	p.lastFetcher = bf
	return nil, nil
}

func (p *blobCapturingPlugin[RI]) Committed(context.Context, uint64, ocr3_1types.KeyValueStateReader) error {
	return nil
}

func (p *blobCapturingPlugin[RI]) Reports(context.Context, uint64, ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[RI], error) {
	return nil, nil
}

func (p *blobCapturingPlugin[RI]) ShouldAcceptAttestedReport(context.Context, uint64, ocr3types.ReportWithInfo[RI]) (bool, error) {
	return true, nil
}

func (p *blobCapturingPlugin[RI]) ShouldTransmitAcceptedReport(context.Context, uint64, ocr3types.ReportWithInfo[RI]) (bool, error) {
	return true, nil
}

func (p *blobCapturingPlugin[RI]) Close() error {
	return nil
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
