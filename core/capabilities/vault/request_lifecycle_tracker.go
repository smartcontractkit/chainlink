package vault

import (
	"context"
	"errors"
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
	stageReceived                   = "received"
	stageBlobBroadcasting           = "blob_broadcasting"
	stageBlobBroadcasted            = "blob_broadcasted"
	stageWrittenToPendingQueue      = "written_to_pending_queue"
	stageObservedOutcome            = "observed_outcome"
	stageStateTransitionOut         = "state_transition_outcome"
	stageTransmitted                = "transmitted"
	stageCapabilityResponseReceived = "capability_response_received"
)

// requestLifecycleTrace holds timestamps and OCR seqNrs for each Vault request as it
// progresses through the capability handler and OCR plugin. Values are keyed in
// RequestLifecycleTracker by the same requestID used in handleRequest and ReportInfo.
type requestLifecycleTrace struct {
	receivedAt time.Time

	blobBroadcastingAt  time.Time
	blobBroadcastingSeq uint64
	hasBlobBroadcasting bool

	blobBroadcastedAt  time.Time
	blobBroadcastedSeq uint64
	hasBlobBroadcasted bool

	pendingQueueAt  time.Time
	pendingQueueSeq uint64
	hasPendingQueue bool

	obsBatchAt  time.Time
	obsBatchSeq uint64
	hasObsBatch bool

	stateTransitionAt  time.Time
	stateTransitionSeq uint64
	hasStateTransition bool

	transmittedAt  time.Time
	transmittedSeq uint64
	hasTransmitted bool

	capabilityResponseAt          time.Time
	hasCapabilityResponseReceived bool
}

type requestLifecycleMetrics struct {
	stageLatencyMs              metric.Int64Histogram
	roundDelta                  metric.Int64Histogram
	outcomeTotal                metric.Int64Counter
	timeoutTotal                metric.Int64Counter
	roundMissingBase            metric.Int64Counter
	responseErrorTotal          metric.Int64Counter
	requestsReceivedTotal       metric.Int64Counter
	pendingQueueNotInLocalQueue metric.Int64Counter
	transmitNotInLocalQueue     metric.Int64Counter
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
		metric.WithDescription("OCR seqNr delta from the blob_broadcasting stage to subsequent lifecycle stages."),
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

	received, err := beholder.GetMeter().Int64Counter(
		"platform_vault_capability_requests_received_total",
		metric.WithDescription("Total Vault capability requests received by this node (for deriving request rate)."),
	)
	if err != nil {
		return nil, fmt.Errorf("vault capability requests received counter: %w", err)
	}

	pqNoLocal, err := beholder.GetMeter().Int64Counter(
		"platform_vault_request_lifecycle_pending_queue_not_in_local_queue_total",
		metric.WithDescription("Pending-queue write observed for a request ID that was not present in this node's local Queue (request appeared in the OCR round but was not received locally first)."),
	)
	if err != nil {
		return nil, fmt.Errorf("vault pending queue not in local Queue counter: %w", err)
	}

	txNoLocal, err := beholder.GetMeter().Int64Counter(
		"platform_vault_request_lifecycle_transmit_not_in_local_queue_total",
		metric.WithDescription("OCR transmit for a request ID not present in this node's local Queue (DON responded but this node had not recorded the request locally)."),
	)
	if err != nil {
		return nil, fmt.Errorf("vault transmit not in local Queue counter: %w", err)
	}

	return &requestLifecycleMetrics{
		stageLatencyMs:              lat,
		roundDelta:                  rounds,
		outcomeTotal:                outcome,
		timeoutTotal:                timeout,
		roundMissingBase:            missBase,
		responseErrorTotal:          respErr,
		requestsReceivedTotal:       received,
		pendingQueueNotInLocalQueue: pqNoLocal,
		transmitNotInLocalQueue:     txNoLocal,
	}, nil
}

// RequestLifecycleTracker records per-request lifecycle data for Vault capability and
// OCR plugin instrumentation. Methods are no-ops when the receiver is nil.
type RequestLifecycleTracker struct {
	lggr    logger.Logger
	traces  map[string]*requestLifecycleTrace
	mu      sync.Mutex
	metrics requestLifecycleMetrics
	digest  atomic.Value // string
}

// NewRequestLifecycleTracker builds a tracker with OTLP metrics registered via beholder.
func NewRequestLifecycleTracker(lggr logger.Logger) (*RequestLifecycleTracker, error) {
	m, err := newRequestLifecycleMetrics()
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("request lifecycle metrics cannot be nil")
	}
	t := &RequestLifecycleTracker{
		lggr:    logger.Named(lggr, "VaultRequestLifecycle"),
		traces:  make(map[string]*requestLifecycleTrace),
		metrics: *m,
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

// furthestStageReached returns the latest lifecycle stage reached for this trace (used when a request ends without a full success path).
func (t *RequestLifecycleTracker) furthestStageReached(tr *requestLifecycleTrace) string {
	if tr == nil {
		return stageReceived
	}
	switch {
	case tr.hasCapabilityResponseReceived:
		return stageCapabilityResponseReceived
	case tr.hasTransmitted:
		return stageTransmitted
	case tr.hasStateTransition:
		return stageStateTransitionOut
	case tr.hasObsBatch:
		return stageObservedOutcome
	case tr.hasPendingQueue:
		return stageWrittenToPendingQueue
	case tr.hasBlobBroadcasted:
		return stageBlobBroadcasted
	case tr.hasBlobBroadcasting:
		return stageBlobBroadcasting
	default:
		return stageReceived
	}
}

// RecordReceived starts tracking at capability.handleRequest — stage: received.
func (t *RequestLifecycleTracker) RecordReceived(ctx context.Context, requestID string, at time.Time) {
	if t == nil {
		return
	}
	t.metrics.requestsReceivedTotal.Add(ctx, 1, t.attrs())
	t.mu.Lock()
	defer t.mu.Unlock()
	t.traces[requestID] = &requestLifecycleTrace{receivedAt: at}
}

func (t *RequestLifecycleTracker) markBlobBroadcasting(tr *requestLifecycleTrace, seq uint64, at time.Time) {
	if tr.hasBlobBroadcasting {
		return
	}
	tr.blobBroadcastingAt = at
	tr.blobBroadcastingSeq = seq
	tr.hasBlobBroadcasting = true
}

// RecordBlobBroadcasting records when a request is first chosen for blob broadcast — stage: blob_broadcasting.
func (t *RequestLifecycleTracker) RecordBlobBroadcasting(requestID string, seq uint64, at time.Time) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tr, ok := t.traces[requestID]
	if !ok {
		return
	}
	t.markBlobBroadcasting(tr, seq, at)
}

func (t *RequestLifecycleTracker) markBlobBroadcasted(tr *requestLifecycleTrace, seq uint64, at time.Time) {
	if tr.hasBlobBroadcasted {
		return
	}
	tr.blobBroadcastedAt = at
	tr.blobBroadcastedSeq = seq
	tr.hasBlobBroadcasted = true
}

// RecordBlobBroadcasted records a successful blob broadcast for the request — stage: blob_broadcasted.
func (t *RequestLifecycleTracker) RecordBlobBroadcasted(requestID string, seq uint64, at time.Time) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tr, ok := t.traces[requestID]
	if !ok {
		return
	}
	t.markBlobBroadcasted(tr, seq, at)
}

// RecordWrittenToPendingQueue records consensus pending-queue persistence — stage: written_to_pending_queue.
func (t *RequestLifecycleTracker) RecordWrittenToPendingQueue(ctx context.Context, requestID string, seq uint64, at time.Time) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tr, ok := t.traces[requestID]
	if !ok {
		t.metrics.pendingQueueNotInLocalQueue.Add(ctx, 1, t.attrs())
		return
	}
	if tr.hasPendingQueue {
		return
	}
	tr.pendingQueueAt = at
	tr.pendingQueueSeq = seq
	tr.hasPendingQueue = true
}

// RecordObservedOutcome records when the request appears in the observation proto built from the KV pending queue batch — stage: observed_outcome.
func (t *RequestLifecycleTracker) RecordObservedOutcome(requestID string, seq uint64, at time.Time) {
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

// RecordStateTransitionOutcome records when an outcome for the request is included in the state transition result — stage: state_transition_outcome.
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

// RecordTransmitted records OCR Transmit for this request id — stage: transmitted.
func (t *RequestLifecycleTracker) RecordTransmitted(ctx context.Context, requestID string, seq uint64, at time.Time) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tr, ok := t.traces[requestID]
	if !ok {
		t.metrics.transmitNotInLocalQueue.Add(ctx, 1, t.attrs())
		return
	}
	if tr.hasTransmitted {
		return
	}
	tr.transmittedAt = at
	tr.transmittedSeq = seq
	tr.hasTransmitted = true
}

func (t *RequestLifecycleTracker) remove(requestID string) *requestLifecycleTrace {
	t.mu.Lock()
	defer t.mu.Unlock()
	tr := t.traces[requestID]
	delete(t.traces, requestID)
	return tr
}

// FinalizeSuccess emits latency / round metrics and removes the trace — stage: capability_response_received on success path.
func (t *RequestLifecycleTracker) FinalizeSuccess(ctx context.Context, requestID string, respondedAt time.Time) {
	if t == nil {
		return
	}
	tr := t.remove(requestID)
	if tr == nil {
		return
	}
	tr.capabilityResponseAt = respondedAt
	tr.hasCapabilityResponseReceived = true
	t.emitLatenciesAndRounds(ctx, tr)
	t.metrics.outcomeTotal.Add(ctx, 1, t.attrs(attribute.String("outcome", "success")))
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
	furthest := t.furthestStageReached(tr)
	t.emitLatenciesAndRounds(ctx, tr)
	t.metrics.timeoutTotal.Add(ctx, 1, t.attrs())
	t.metrics.outcomeTotal.Add(ctx, 1, t.attrs(
		attribute.String("outcome", "timeout"),
		attribute.String("furthest_stage", furthest),
	))
	t.lggr.Warnw("vault request timed out in capability.handleRequest before a response was delivered",
		append([]any{"requestID", requestID, "furthest_stage", furthest}, traceLogFields(tr)...)...)
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
	tr.hasCapabilityResponseReceived = true
	t.lggr.Warnw("vault request closed with OCR error response", "requestID", requestID, "err", errMsg, "lifecycle", traceLogFields(tr))
	t.emitLatenciesAndRounds(ctx, tr)
	t.metrics.responseErrorTotal.Add(ctx, 1, t.attrs())
	t.metrics.outcomeTotal.Add(ctx, 1, t.attrs(attribute.String("outcome", "response_error")))
}

func (t *RequestLifecycleTracker) emitLatenciesAndRounds(ctx context.Context, tr *requestLifecycleTrace) {
	if tr.receivedAt.IsZero() {
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

	emitLatency(stageBlobBroadcasting, tr.blobBroadcastingAt, tr.hasBlobBroadcasting)
	emitLatency(stageBlobBroadcasted, tr.blobBroadcastedAt, tr.hasBlobBroadcasted)
	emitLatency(stageWrittenToPendingQueue, tr.pendingQueueAt, tr.hasPendingQueue)
	emitLatency(stageObservedOutcome, tr.obsBatchAt, tr.hasObsBatch)
	emitLatency(stageStateTransitionOut, tr.stateTransitionAt, tr.hasStateTransition)
	emitLatency(stageTransmitted, tr.transmittedAt, tr.hasTransmitted)
	emitLatency(stageCapabilityResponseReceived, tr.capabilityResponseAt, tr.hasCapabilityResponseReceived)

	if !tr.hasBlobBroadcasting {
		t.emitRoundSkips(ctx, stageBlobBroadcasted, stageWrittenToPendingQueue, stageObservedOutcome, stageStateTransitionOut, stageTransmitted)
		return
	}

	emitRound := func(stage string, seq uint64, ok bool) {
		if !ok {
			return
		}
		delta := uint64SeqDeltaToInt64(seq, tr.blobBroadcastingSeq)
		if delta < 0 {
			delta = 0
		}
		t.metrics.roundDelta.Record(ctx, delta, t.attrs(attribute.String("stage", stage)))
	}

	emitRound(stageBlobBroadcasted, tr.blobBroadcastedSeq, tr.hasBlobBroadcasted)
	emitRound(stageWrittenToPendingQueue, tr.pendingQueueSeq, tr.hasPendingQueue)
	emitRound(stageObservedOutcome, tr.obsBatchSeq, tr.hasObsBatch)
	emitRound(stageStateTransitionOut, tr.stateTransitionSeq, tr.hasStateTransition)
	emitRound(stageTransmitted, tr.transmittedSeq, tr.hasTransmitted)
}

func (t *RequestLifecycleTracker) emitRoundSkips(ctx context.Context, stages ...string) {
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
	baseSeq, baseOK := tr.blobBroadcastingSeq, tr.hasBlobBroadcasting
	return []interface{}{
		"receivedAt", tr.receivedAt,
		"blob_broadcasting", tr.hasBlobBroadcasting, "blob_broadcasting_at", tr.blobBroadcastingAt, "blob_broadcasting_seq", tr.blobBroadcastingSeq,
		"blob_broadcasted", tr.hasBlobBroadcasted, "blob_broadcasted_at", tr.blobBroadcastedAt, "blob_broadcasted_seq", tr.blobBroadcastedSeq,
		"rounds_blob_broadcasted_after_blob_broadcasting", roundDeltaOrNeg(tr.blobBroadcastedSeq, tr.hasBlobBroadcasted, baseSeq, baseOK),
		"written_to_pending_queue", tr.hasPendingQueue, "written_to_pending_queue_at", tr.pendingQueueAt, "written_to_pending_queue_seq", tr.pendingQueueSeq,
		"rounds_written_to_pending_queue_after_blob_broadcasting", roundDeltaOrNeg(tr.pendingQueueSeq, tr.hasPendingQueue, baseSeq, baseOK),
		"observed_outcome", tr.hasObsBatch, "observed_outcome_at", tr.obsBatchAt, "observed_outcome_seq", tr.obsBatchSeq,
		"rounds_observed_outcome_after_blob_broadcasting", roundDeltaOrNeg(tr.obsBatchSeq, tr.hasObsBatch, baseSeq, baseOK),
		"state_transition_outcome", tr.hasStateTransition, "state_transition_outcome_at", tr.stateTransitionAt, "state_transition_outcome_seq", tr.stateTransitionSeq,
		"rounds_state_transition_after_blob_broadcasting", roundDeltaOrNeg(tr.stateTransitionSeq, tr.hasStateTransition, baseSeq, baseOK),
		"transmitted", tr.hasTransmitted, "transmitted_at", tr.transmittedAt, "transmitted_seq", tr.transmittedSeq,
		"rounds_transmitted_after_blob_broadcasting", roundDeltaOrNeg(tr.transmittedSeq, tr.hasTransmitted, baseSeq, baseOK),
		"capability_response_received", tr.hasCapabilityResponseReceived, "capability_response_received_at", tr.capabilityResponseAt,
	}
}
