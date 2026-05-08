package vault

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

const (
	stageBlobChosen         = "blob_chosen"
	stageBlobBroadcastOK    = "blob_broadcast_ok"
	stagePendingQueueWrite  = "pending_queue_write"
	stageObservationInBatch = "observation_pending_batch"
	stageStateTransitionOut = "state_transition_outcome"
	stageTransmit           = "transmit"
	stageCapabilityResponse = "capability_response"
)

type requestLifecycleMetrics struct {
	stageLatencyMs     metric.Int64Histogram
	roundDelta         metric.Int64Histogram
	outcomeTotal       metric.Int64Counter
	timeoutTotal       metric.Int64Counter
	roundMissingBase   metric.Int64Counter
	responseErrorTotal metric.Int64Counter
}

func newRequestLifecycleMetrics() (*requestLifecycleMetrics, error) {
	lat, err := beholder.GetMeter().Int64Histogram(
		"platform_vault_request_lifecycle_stage_latency_ms",
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, fmt.Errorf("request lifecycle stage latency histogram: %w", err)
	}

	rounds, err := beholder.GetMeter().Int64Histogram(
		"platform_vault_request_lifecycle_stage_rounds_delta",
		metric.WithDescription("OCR seqNr delta from blob_chosen seq for stages 3-7"),
	)
	if err != nil {
		return nil, fmt.Errorf("request lifecycle rounds histogram: %w", err)
	}

	outcome, err := beholder.GetMeter().Int64Counter("platform_vault_capability_request_outcome_total")
	if err != nil {
		return nil, fmt.Errorf("vault capability request outcome counter: %w", err)
	}

	timeout, err := beholder.GetMeter().Int64Counter("platform_vault_request_lifecycle_timeout_total")
	if err != nil {
		return nil, fmt.Errorf("vault request lifecycle timeout counter: %w", err)
	}

	missBase, err := beholder.GetMeter().Int64Counter("platform_vault_request_lifecycle_round_delta_skipped_total")
	if err != nil {
		return nil, fmt.Errorf("vault request lifecycle round delta skipped counter: %w", err)
	}

	respErr, err := beholder.GetMeter().Int64Counter("platform_vault_capability_request_response_error_total")
	if err != nil {
		return nil, fmt.Errorf("vault capability response error counter: %w", err)
	}

	return &requestLifecycleMetrics{
		stageLatencyMs:     lat,
		roundDelta:         rounds,
		outcomeTotal:       outcome,
		timeoutTotal:       timeout,
		roundMissingBase:   missBase,
		responseErrorTotal: respErr,
	}, nil
}

// RequestLifecycleTracker records per-request lifecycle data for Vault capability and
// OCR plugin instrumentation. Methods are no-ops when the receiver is nil.
type RequestLifecycleTracker struct {
	lggr    logger.Logger
	traces  map[string]*requestLifecycleTrace
	mu      sync.Mutex
	metrics *requestLifecycleMetrics
	digest  atomic.Value // string
}

// NewRequestLifecycleTracker builds a tracker with OTLP metrics registered via beholder.
func NewRequestLifecycleTracker(lggr logger.Logger) (*RequestLifecycleTracker, error) {
	m, err := newRequestLifecycleMetrics()
	if err != nil {
		return nil, err
	}
	t := &RequestLifecycleTracker{
		lggr:    logger.Named(lggr, "VaultRequestLifecycle"),
		traces:  make(map[string]*requestLifecycleTrace),
		metrics: m,
	}
	t.digest.Store("")
	return t, nil
}

// SetConfigDigest updates the config digest label used on emitted metrics (OCR config).
func (t *RequestLifecycleTracker) SetConfigDigest(digest string) {
	if t == nil {
		return
	}
	t.digest.Store(digest)
}

func (t *RequestLifecycleTracker) configDigestAttr() attribute.KeyValue {
	d, _ := t.digest.Load().(string)
	return attribute.String("config_digest", d)
}

func (t *RequestLifecycleTracker) attrs(extra ...attribute.KeyValue) metric.MeasurementOption {
	out := make([]attribute.KeyValue, 0, len(extra)+1)
	out = append(out, t.configDigestAttr())
	out = append(out, extra...)
	return metric.WithAttributes(out...)
}

// RecordReceived starts tracking at capability.handleRequest (event 1).
func (t *RequestLifecycleTracker) RecordReceived(requestID string, at time.Time) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.traces[requestID] = &requestLifecycleTrace{receivedAt: at}
}

func (t *RequestLifecycleTracker) markBlobChosen(tr *requestLifecycleTrace, seq uint64, at time.Time) {
	if tr.hasBlobChosen {
		return
	}
	tr.blobChosenAt = at
	tr.blobChosenSeq = seq
	tr.hasBlobChosen = true
}

// RecordBlobChosen records when a request is first chosen for blob broadcast (event 2).
func (t *RequestLifecycleTracker) RecordBlobChosen(requestID string, seq uint64, at time.Time) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tr, ok := t.traces[requestID]
	if !ok {
		return
	}
	t.markBlobChosen(tr, seq, at)
}

func (t *RequestLifecycleTracker) markBlobBroadcast(tr *requestLifecycleTrace, seq uint64, at time.Time) {
	if tr.hasBlobBroadcast {
		return
	}
	tr.blobBroadcastAt = at
	tr.blobBroadcastSeq = seq
	tr.hasBlobBroadcast = true
}

// RecordBlobBroadcastOK records a successful blob broadcast for the request (event 3).
func (t *RequestLifecycleTracker) RecordBlobBroadcastOK(requestID string, seq uint64, at time.Time) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tr, ok := t.traces[requestID]
	if !ok {
		return
	}
	t.markBlobBroadcast(tr, seq, at)
}

// RecordPendingQueueWrite records consensus pending-queue persistence (event 4).
func (t *RequestLifecycleTracker) RecordPendingQueueWrite(requestID string, seq uint64, at time.Time) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tr, ok := t.traces[requestID]
	if !ok {
		return
	}
	if tr.hasPendingQueue {
		return
	}
	tr.pendingQueueAt = at
	tr.pendingQueueSeq = seq
	tr.hasPendingQueue = true
}

// RecordObservationPendingBatch records when the request appears in the observation
// proto built from the KV pending queue batch (event 5).
func (t *RequestLifecycleTracker) RecordObservationPendingBatch(requestID string, seq uint64, at time.Time) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tr, ok := t.traces[requestID]
	if !ok {
		return
	}
	if tr.hasObsBatch {
		return
	}
	tr.obsBatchAt = at
	tr.obsBatchSeq = seq
	tr.hasObsBatch = true
}

// RecordStateTransitionOutcome records when an outcome for the request is included in
// the state transition result (event 6).
func (t *RequestLifecycleTracker) RecordStateTransitionOutcome(requestID string, seq uint64, at time.Time) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tr, ok := t.traces[requestID]
	if !ok {
		return
	}
	if tr.hasStateTransition {
		return
	}
	tr.stateTransitionAt = at
	tr.stateTransitionSeq = seq
	tr.hasStateTransition = true
}

// RecordTransmit records OCR Transmit for this request id (event 7).
func (t *RequestLifecycleTracker) RecordTransmit(requestID string, seq uint64, at time.Time) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tr, ok := t.traces[requestID]
	if !ok {
		return
	}
	if tr.hasTransmit {
		return
	}
	tr.transmitAt = at
	tr.transmitSeq = seq
	tr.hasTransmit = true
}

func (t *RequestLifecycleTracker) remove(requestID string) *requestLifecycleTrace {
	t.mu.Lock()
	defer t.mu.Unlock()
	tr := t.traces[requestID]
	delete(t.traces, requestID)
	return tr
}

// FinalizeSuccess emits latency / round metrics and removes the trace (event 8: successful response).
func (t *RequestLifecycleTracker) FinalizeSuccess(ctx context.Context, requestID string, respondedAt time.Time) {
	if t == nil {
		return
	}
	tr := t.remove(requestID)
	if tr == nil {
		return
	}
	tr.capabilityResponseAt = respondedAt
	tr.hasCapabilityResponse = true
	t.emitLatenciesAndRounds(ctx, tr)
	if t.metrics != nil {
		t.metrics.outcomeTotal.Add(ctx, 1, t.attrs(attribute.String("outcome", "success")))
	}
}

// FinalizeTimeout logs pipeline state, emits timeout/failure telemetry, and removes the trace.
func (t *RequestLifecycleTracker) FinalizeTimeout(ctx context.Context, requestID string) {
	if t == nil {
		return
	}
	tr := t.remove(requestID)
	if tr == nil {
		return
	}
	t.lggr.Errorw("vault request timed out in capability.handleRequest", append([]any{"requestID", requestID}, traceLogFields(tr)...)...)
	if t.metrics != nil {
		t.metrics.timeoutTotal.Add(ctx, 1, t.attrs())
		t.metrics.outcomeTotal.Add(ctx, 1, t.attrs(attribute.String("outcome", "timeout")))
	}
}

// FinalizeResponseError records a capability-layer response error (non-timeout) and removes the trace.
func (t *RequestLifecycleTracker) FinalizeResponseError(ctx context.Context, requestID string, respondedAt time.Time, errMsg string) {
	if t == nil {
		return
	}
	tr := t.remove(requestID)
	if tr == nil {
		return
	}
	tr.capabilityResponseAt = respondedAt
	tr.hasCapabilityResponse = true
	t.lggr.Warnw("vault request closed with OCR error response", "requestID", requestID, "err", errMsg, "lifecycle", traceLogFields(tr))
	t.emitLatenciesAndRounds(ctx, tr)
	if t.metrics != nil {
		t.metrics.responseErrorTotal.Add(ctx, 1, t.attrs())
		t.metrics.outcomeTotal.Add(ctx, 1, t.attrs(attribute.String("outcome", "response_error")))
	}
}

func (t *RequestLifecycleTracker) emitLatenciesAndRounds(ctx context.Context, tr *requestLifecycleTrace) {
	if t.metrics == nil || tr.receivedAt.IsZero() {
		return
	}
	base := tr.receivedAt
	emitLatency := func(stage string, at time.Time, ok bool) {
		if !ok || at.IsZero() {
			return
		}
		ms := at.Sub(base).Milliseconds()
		if ms < 0 {
			ms = 0
		}
		t.metrics.stageLatencyMs.Record(ctx, ms, t.attrs(attribute.String("stage", stage)))
	}

	emitLatency(stageBlobChosen, tr.blobChosenAt, tr.hasBlobChosen)
	emitLatency(stageBlobBroadcastOK, tr.blobBroadcastAt, tr.hasBlobBroadcast)
	emitLatency(stagePendingQueueWrite, tr.pendingQueueAt, tr.hasPendingQueue)
	emitLatency(stageObservationInBatch, tr.obsBatchAt, tr.hasObsBatch)
	emitLatency(stageStateTransitionOut, tr.stateTransitionAt, tr.hasStateTransition)
	emitLatency(stageTransmit, tr.transmitAt, tr.hasTransmit)
	emitLatency(stageCapabilityResponse, tr.capabilityResponseAt, tr.hasCapabilityResponse)

	if !tr.hasBlobChosen {
		t.emitRoundSkips(ctx, stageBlobBroadcastOK, stagePendingQueueWrite, stageObservationInBatch, stageStateTransitionOut, stageTransmit)
		return
	}

	emitRound := func(stage string, seq uint64, ok bool) {
		if !ok {
			return
		}
		delta := uint64SeqDeltaToInt64(seq, tr.blobChosenSeq)
		if delta < 0 {
			delta = 0
		}
		t.metrics.roundDelta.Record(ctx, delta, t.attrs(attribute.String("stage", stage)))
	}

	emitRound(stageBlobBroadcastOK, tr.blobBroadcastSeq, tr.hasBlobBroadcast)
	emitRound(stagePendingQueueWrite, tr.pendingQueueSeq, tr.hasPendingQueue)
	emitRound(stageObservationInBatch, tr.obsBatchSeq, tr.hasObsBatch)
	emitRound(stageStateTransitionOut, tr.stateTransitionSeq, tr.hasStateTransition)
	emitRound(stageTransmit, tr.transmitSeq, tr.hasTransmit)
}

func (t *RequestLifecycleTracker) emitRoundSkips(ctx context.Context, stages ...string) {
	if t.metrics == nil {
		return
	}
	for _, s := range stages {
		t.metrics.roundMissingBase.Add(ctx, 1, t.attrs(attribute.String("stage", s)))
	}
}

func roundDeltaOrNeg(seq uint64, ok bool, baseSeq uint64, baseOK bool) int64 {
	if !ok || !baseOK {
		return -1
	}
	return uint64SeqDeltaToInt64(seq, baseSeq)
}

// uint64SeqDeltaToInt64 returns a-b in int64, clamping so uint64→int64 conversions never overflow (gosec G115).
func uint64SeqDeltaToInt64(a, b uint64) int64 {
	if a >= b {
		d := a - b
		if d > uint64(math.MaxInt64) {
			return math.MaxInt64
		}
		return int64(d)
	}
	d := b - a
	if d > uint64(math.MaxInt64) {
		return math.MinInt64
	}
	return -int64(d)
}

func traceLogFields(tr *requestLifecycleTrace) []interface{} {
	baseSeq, baseOK := tr.blobChosenSeq, tr.hasBlobChosen
	return []interface{}{
		"receivedAt", tr.receivedAt,
		"blobChosen", tr.hasBlobChosen, "blobChosenAt", tr.blobChosenAt, "blobChosenSeq", tr.blobChosenSeq,
		"blobBroadcastOk", tr.hasBlobBroadcast, "blobBroadcastAt", tr.blobBroadcastAt, "blobBroadcastSeq", tr.blobBroadcastSeq,
		"rounds_blobBroadcast_after_blobChosen", roundDeltaOrNeg(tr.blobBroadcastSeq, tr.hasBlobBroadcast, baseSeq, baseOK),
		"pendingQueue", tr.hasPendingQueue, "pendingQueueAt", tr.pendingQueueAt, "pendingQueueSeq", tr.pendingQueueSeq,
		"rounds_pendingQueue_after_blobChosen", roundDeltaOrNeg(tr.pendingQueueSeq, tr.hasPendingQueue, baseSeq, baseOK),
		"obsBatch", tr.hasObsBatch, "obsBatchAt", tr.obsBatchAt, "obsBatchSeq", tr.obsBatchSeq,
		"rounds_obsBatch_after_blobChosen", roundDeltaOrNeg(tr.obsBatchSeq, tr.hasObsBatch, baseSeq, baseOK),
		"stateTransition", tr.hasStateTransition, "stateTransitionAt", tr.stateTransitionAt, "stateTransitionSeq", tr.stateTransitionSeq,
		"rounds_stateTransition_after_blobChosen", roundDeltaOrNeg(tr.stateTransitionSeq, tr.hasStateTransition, baseSeq, baseOK),
		"transmit", tr.hasTransmit, "transmitAt", tr.transmitAt, "transmitSeq", tr.transmitSeq,
		"rounds_transmit_after_blobChosen", roundDeltaOrNeg(tr.transmitSeq, tr.hasTransmit, baseSeq, baseOK),
		"capabilityResponse", tr.hasCapabilityResponse, "capabilityResponseAt", tr.capabilityResponseAt,
	}
}
