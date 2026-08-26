package observation

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"math"
	"math/big"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	lloprotocol "github.com/smartcontractkit/chainlink-data-streams/llo/protocol"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// ObservationContext ensures that each pipeline is only executed once. It is
// intended to be instantiated and used then discarded as part of one
// Observation cycle. Subsequent calls to Observe will return the same cached
// values.

var _ ObservationContext = (*observationContext)(nil)

type ObservationContext interface { //nolint:revive // ObservationContext is the established interface name in this package
	Observe(ctx context.Context, streamID streams.StreamID, opts telem.DSOpts) (val lloprotocol.StreamValue, err error)
}

type execution struct {
	done <-chan struct{}

	run  *pipeline.Run
	trrs pipeline.TaskRunResults
	err  error

	// vals holds the stream values resolved from trrs by the executing
	// goroutine, keyed by stream ID. Read-only once done is closed.
	vals resolvedStreamValues
}

// resolvedStreamValue is the outcome of converting one pipeline output into a
// StreamValue. err is scoped to this stream only; a group-fatal failure is
// reported through execution.err instead.
type resolvedStreamValue struct {
	val        lloprotocol.StreamValue
	err        error
	finishedAt time.Time // zero if the task did not report a finish time
	// fallback is true if the value came from the terminal outputs of the
	// pipeline rather than from a streamID-tagged task.
	fallback bool
}

type resolvedStreamValues map[streams.StreamID]resolvedStreamValue

type observationContext struct {
	l logger.Logger
	r Registry
	t Telemeter

	executionsMu sync.Mutex
	// only execute each pipeline once
	executions map[streams.Pipeline]*execution
}

func NewObservationContext(l logger.Logger, r Registry, t Telemeter) ObservationContext {
	return newObservationContext(l, r, t)
}

func newObservationContext(l logger.Logger, r Registry, t Telemeter) *observationContext {
	return &observationContext{l, r, t, sync.Mutex{}, make(map[streams.Pipeline]*execution)}
}

func (oc *observationContext) Observe(ctx context.Context, streamID streams.StreamID, opts telem.DSOpts) (val lloprotocol.StreamValue, err error) {
	ex, err := oc.run(ctx, streamID)
	observationFinishedAt := time.Now()
	if err != nil {
		// Either the pipeline itself failed, or one of the stream values it
		// produced violated an invariant. In both cases every stream served by
		// this pipeline fails for this observation round.
		var run *pipeline.Run
		var trrs pipeline.TaskRunResults
		if ex != nil {
			run, trrs = ex.run, ex.trrs
		}
		// FIXME: This is a hack specific for V3 telemetry, future schemas should
		// use a generic stream value telemetry instead
		// https://smartcontract-it.atlassian.net/browse/MERC-6290
		oc.t.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, nil, err)
		oc.enqueueObservationTelemetry(ctx, streamID, opts, nil, observationFinishedAt, err)
		return nil, err
	}

	rsv, found := ex.vals[streamID]
	if !found {
		// The pipeline tags stream IDs, but none of them is the one we were
		// asked for. Assume the final output is this stream, same as for an
		// untagged pipeline.
		rsv = resolveFinalResult(ex.trrs)
		if rsv.err == nil {
			rsv.err = validateStreamValue(streamID, rsv.val)
			if rsv.err != nil {
				rsv.val = nil
			}
		}
	}
	if rsv.fallback && oc.t.CaptureEATelemetry() {
		// FIXME: This is a hack specific for V3 telemetry, future schemas should
		// use the generic stream value telemetry instead
		// https://smartcontract-it.atlassian.net/browse/MERC-6290
		oc.t.EnqueueV3PremiumLegacy(ex.run, ex.trrs, streamID, opts, rsv.val, rsv.err)
	}
	if !rsv.finishedAt.IsZero() {
		observationFinishedAt = rsv.finishedAt
	}
	oc.enqueueObservationTelemetry(ctx, streamID, opts, rsv.val, observationFinishedAt, rsv.err)
	if rsv.err != nil {
		// Scoped to this stream: sibling streams from the same pipeline are
		// still usable.
		return nil, rsv.err
	}
	return rsv.val, nil
}

func (oc *observationContext) enqueueObservationTelemetry(ctx context.Context, streamID streams.StreamID, opts telem.DSOpts, val lloprotocol.StreamValue, observationFinishedAt time.Time, obsErr error) {
	ch := GetObservationTelemetryCh(ctx)
	if ch == nil {
		return
	}
	cd := opts.ConfigDigest()
	ot := &telem.LLOObservationTelemetry{
		StreamId:              streamID,
		ObservationTimestamp:  opts.ObservationTimestamp().UnixNano(),
		ObservationFinishedAt: observationFinishedAt.UnixNano(),
		SeqNr:                 opts.SeqNr(),
		ConfigDigest:          cd[:],
	}
	if obsErr != nil {
		ot.ObservationError = new(string)
		*ot.ObservationError = obsErr.Error()
	}
	if val != nil {
		ot.StreamValueType = int32(val.Type())
		b, err := val.MarshalBinary()
		if err != nil {
			oc.l.Errorw("failed to MarshalBinary on stream value", "error", err)
		} else {
			ot.StreamValueBinary = b
		}
		s, err := val.MarshalText()
		if err != nil {
			oc.l.Errorw("failed to MarshalText on stream value", "error", err)
		} else {
			ot.StreamValueText = string(s)
		}
	}
	select {
	case ch <- ot:
	default:
		oc.l.Error("telemetry channel is full, dropping observation telemetry")
	}
}

// resolvePipelineStreamValues converts every output of a single pipeline run
// into StreamValues, keyed by stream ID. All streamID-tagged tasks are
// converted, not just the requested one, so that an invariant violation on a
// sibling stream is visible to the caller. If the requested stream ID has no
// tagged task, the terminal outputs are resolved for it as a fallback.
//
// Conversion failures are recorded per stream rather than returned, so that one
// broken stream does not take down its siblings.
func resolvePipelineStreamValues(p streams.Pipeline, streamID streams.StreamID, trrs pipeline.TaskRunResults) resolvedStreamValues {
	resolved := make(resolvedStreamValues, len(trrs))
	for _, trr := range trrs {
		sid := trr.Task.TaskStreamID()
		if sid == nil {
			continue
		}
		if _, exists := resolved[*sid]; exists {
			// Duplicate tag; first one wins.
			continue
		}
		var rsv resolvedStreamValue
		val, err := resultToStreamValue(trr.Result.Value)
		if err != nil {
			rsv.err = fmt.Errorf("failed to convert result to StreamValue for streamID %d: %w", *sid, err)
		} else {
			rsv.val = val
		}
		if trr.FinishedAt.Valid {
			rsv.finishedAt = trr.FinishedAt.Time
		}
		resolved[*sid] = rsv
	}
	// Any stream served by this pipeline that has no tagged task takes its value
	// from the terminal outputs. Resolve those here too, even if the caller did
	// not ask for them, so that a quote produced this way is validated together
	// with its siblings.
	untagged := append([]streams.StreamID(nil), p.StreamIDs()...)
	untagged = append(untagged, streamID)
	for _, sid := range untagged {
		if _, exists := resolved[sid]; exists {
			continue
		}
		resolved[sid] = resolveFinalResult(trrs)
	}
	return resolved
}

// resolveFinalResult treats the terminal outputs of the pipeline as the value
// for a stream with no streamID-tagged task. This is safe to do since the
// registry will never return a spec that doesn't match either by tag or by spec
// streamID.
func resolveFinalResult(trrs pipeline.TaskRunResults) resolvedStreamValue {
	val, err := extractFinalResultAsStreamValue(trrs)
	if err != nil {
		return resolvedStreamValue{err: err, fallback: true}
	}
	return resolvedStreamValue{val: val, fallback: true}
}

// QuoteInvariantError is returned when a Quote stream value violates the
// Bid <= Benchmark <= Ask invariant. It fails every stream served by the
// pipeline that produced it, not just the offending stream, since a pipeline
// emitting a nonsensical quote cannot be trusted for its sibling streams
// either.
type QuoteInvariantError struct {
	StreamID streams.StreamID
	Quote    *lloprotocol.Quote
}

func (e QuoteInvariantError) Error() string {
	desc := "<nil>"
	if e.Quote != nil {
		if b, err := e.Quote.MarshalText(); err == nil {
			desc = string(b)
		} else {
			desc = fmt.Sprintf("Q{Bid: %s, Benchmark: %s, Ask: %s}", e.Quote.Bid.String(), e.Quote.Benchmark.String(), e.Quote.Ask.String())
		}
	}
	return fmt.Sprintf("quote invariant violation for stream %d, expected Bid <= Benchmark <= Ask, got: %s", e.StreamID, desc)
}

// validateStreamValues returns the first invariant violation found across all
// streams produced by one pipeline run. Streams are visited in ascending stream
// ID order so the reported violation is deterministic.
func validateStreamValues(resolved resolvedStreamValues) error {
	for _, sid := range slices.Sorted(maps.Keys(resolved)) {
		rsv := resolved[sid]
		if rsv.err != nil {
			// Conversion already failed for this stream; it is reported to that
			// stream alone.
			continue
		}
		if err := validateStreamValue(sid, rsv.val); err != nil {
			return err
		}
	}
	return nil
}

func validateStreamValue(streamID streams.StreamID, val lloprotocol.StreamValue) error {
	q, ok := val.(*lloprotocol.Quote)
	if !ok || q.IsValid() {
		return nil
	}
	return QuoteInvariantError{StreamID: streamID, Quote: q}
}

func resultToStreamValue(val any) (lloprotocol.StreamValue, error) {
	switch v := val.(type) {
	case decimal.Decimal:
		return lloprotocol.ToDecimal(v), nil
	case float64:
		return lloprotocol.ToDecimal(decimal.NewFromFloat(v)), nil
	case int64:
		return lloprotocol.ToDecimal(decimal.NewFromInt(v)), nil
	case pipeline.ObjectParam:
		switch v.Type {
		case pipeline.DecimalType:
			return lloprotocol.ToDecimal(decimal.Decimal(v.DecimalValue)), nil
		default:
			return nil, fmt.Errorf("don't know how to convert pipeline.ObjectParam with type %d to lloprotocol.StreamValue", v.Type)
		}
	case map[string]any:
		sv, err := resultMapToStreamValue(v)
		if err != nil {
			return nil, fmt.Errorf("don't know how to convert map to StreamValue: %w; got: %v", err, v)
		}
		return sv, nil
	default:
		return nil, fmt.Errorf("don't know how to convert pipeline output result of type %T to lloprotocol.StreamValue (got: %v)", val, val)
	}
}

// Converts arbitrary JSON (parsed to map) to a StreamValue
func resultMapToStreamValue(m map[string]any) (lloprotocol.StreamValue, error) {
	var streamValueType lloprotocol.LLOStreamValue_Type
	{
		raw, exists := m["streamValueType"]
		if !exists {
			return nil, errors.New("expected a key labeled 'streamValueType' in the map")
		}
		rawInt64, ok := raw.(int64)
		if !ok {
			return nil, fmt.Errorf("expected 'streamValueType' to be a int64, got: %T", raw)
		}
		if rawInt64 < 0 || rawInt64 > math.MaxUint32 || rawInt64 >= int64(lloprotocol.LLOStreamValue_Type(len(lloprotocol.LLOStreamValue_Type_name))) { //nolint:gosec // G115 // won't overflow
			return nil, fmt.Errorf("invalid streamValueType: %v", rawInt64)
		}
		streamValueType = lloprotocol.LLOStreamValue_Type(rawInt64) //nolint:gosec // G115 // won't overflow due to check above
	}
	switch streamValueType {
	case lloprotocol.LLOStreamValue_Quote:
		r, err := resultMapToQuote(m)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Quote: %w", err)
		}
		return r, nil
	case lloprotocol.LLOStreamValue_TimestampedStreamValue:
		r, err := resultMapToTimestampedStreamValue(m)
		if err != nil {
			return nil, fmt.Errorf("failed to parse TimestampedStreamValue: %w", err)
		}
		return r, nil
	default:
		return nil, fmt.Errorf("unknown streamValueType: %v", m["streamValueType"])
	}
}

// expects something in the form of:
//
//	{
//	  "streamValueType": 1,
//	  "bid": "123.456",
//	  "benchmark": "123.789",
//	  "ask": "124.012"
//	}
//
// "mid" is accepted as an alias for "benchmark"; supplying both is an error.
func resultMapToQuote(m map[string]any) (*lloprotocol.Quote, error) {
	const benchmarkKey, midKey = "benchmark", "mid"
	_, hasBenchmark := m[benchmarkKey]
	_, hasMid := m[midKey]
	if hasBenchmark && hasMid {
		return nil, fmt.Errorf("expected exactly one of '%s' or '%s', got both", benchmarkKey, midKey)
	}
	benchmarkAlias := benchmarkKey
	if hasMid {
		benchmarkAlias = midKey
	}
	q := new(lloprotocol.Quote)
	for _, f := range []struct {
		key string
		dst *decimal.Decimal
	}{
		{"bid", &q.Bid},
		{benchmarkAlias, &q.Benchmark},
		{"ask", &q.Ask},
	} {
		raw, exists := m[f.key]
		if !exists || raw == nil {
			return nil, fmt.Errorf("expected a key labeled '%s'", f.key)
		}
		d, err := utils.ToDecimal(raw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse '%s' as a decimal: %w", f.key, err)
		}
		*f.dst = d
	}
	return q, nil
}

// expects something in the form of:
//
//	{
//	  "timestamps": {
//	  	"providerIndicatedTimeUnixMs": 1234567890,
//	  	"providerDataReceivedUnixMs": 1234567890
//	  },
//	  "result": "123.456"
//	}
func resultMapToTimestampedStreamValue(m map[string]any) (*lloprotocol.TimestampedStreamValue, error) {
	ts, ok := m["timestamps"].(map[string]any)
	if !ok {
		return nil, errors.New("expected a key labeled 'timestamps' as map[string]interface{}")
	}
	// providerIndicatedTimeUnixMs is the best option, with providerDataReceivedUnixMs as a fallback
	k := "providerIndicatedTimeUnixMs"
	rawObservedAtMillis, exists := ts[k]
	if !exists || rawObservedAtMillis == nil {
		k = "providerDataReceivedUnixMs"
		rawObservedAtMillis, exists = ts[k]
		if !exists || rawObservedAtMillis == nil {
			return nil, errors.New("expected a key labeled 'providerIndicatedTimeUnixMs' or 'providerDataReceivedUnixMs'")
		}
	}
	observedAtMillis, err := toUint64(rawObservedAtMillis)
	if err != nil {
		return nil, fmt.Errorf("failed to parse '%s' as a uint64: %w", k, err)
	}
	rStreamValue, exists := m["result"]
	if !exists {
		return nil, errors.New("expected a key labeled 'result'")
	}
	// Assume it's always a decimal for now
	svd, err := utils.ToDecimal(rStreamValue)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 'result' as a decimal: %w", err)
	}

	if observedAtMillis > (math.MaxUint64 / uint64(1e6)) {
		return nil, fmt.Errorf("observedAtMillis too large, got: %d", observedAtMillis)
	}

	return &lloprotocol.TimestampedStreamValue{
		ObservedAtNanoseconds: observedAtMillis * 1e6, // convert ms to ns
		StreamValue:           lloprotocol.ToDecimal(svd),
	}, nil
}

func toUint64(v any) (uint64, error) {
	switch v := v.(type) {
	case float64:
		if v < 0 {
			return 0, fmt.Errorf("expected positive float64, got: %f", v)
		}
		if v > math.MaxUint64 {
			return 0, fmt.Errorf("float64 too large, got: %f", v)
		}
		return uint64(v), nil
	case string:
		return strconv.ParseUint(v, 10, 64)
	case int64:
		if v < 0 {
			return 0, fmt.Errorf("expected positive int64, got: %d", v)
		}
		return uint64(v), nil
	case decimal.Decimal:
		if v.IsNegative() {
			return 0, fmt.Errorf("expected positive decimal, got: %s", v.String())
		}
		// Convert the integer part of the decimal to a big.Int, discard fractional part
		bigIntVal := v.BigInt()

		// Check if the value is greater than maxUint64
		if bigIntVal.Cmp(new(big.Int).SetUint64(math.MaxUint64)) > 0 {
			return 0, fmt.Errorf("decimal too large, got: %s", v.String())
		}
		return bigIntVal.Uint64(), nil
	default:
		return 0, fmt.Errorf("expected float64, int64, string, or decimal, got: %T", v)
	}
}

// extractFinalResultAsStreamValue extracts a final StreamValue from a TaskRunResults
func extractFinalResultAsStreamValue(trrs pipeline.TaskRunResults) (lloprotocol.StreamValue, error) {
	// pipeline.TaskRunResults comes ordered asc by index, this is guaranteed
	// by the pipeline executor
	finaltrrs := trrs.Terminals()

	// HACK: Right now we rely on the number of outputs to determine whether
	// its a Decimal or a Quote.
	// This is a hack to support the legacy "Quote" case.
	// Future stream specs should use streamID tags instead.
	switch len(finaltrrs) {
	case 1:
		res := finaltrrs[0].Result
		if res.Error != nil {
			return nil, fmt.Errorf("terminal task error: %w; all task errors: %w", res.Error, trrs.AllErrors())
		}
		val, err := toDecimal(res.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse BenchmarkPrice: %w", err)
		}
		return lloprotocol.ToDecimal(val), nil
	case 3:
		// Expect ordering of Benchmark, Bid, Ask
		results := make([]decimal.Decimal, 3)
		for i, trr := range finaltrrs {
			res := trr.Result
			if res.Error != nil {
				return nil, fmt.Errorf("failed to parse stream output into Quote (task index: %d): %w", i, res.Error)
			}
			val, err := toDecimal(res.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse decimal: %w", err)
			}
			results[i] = val
		}
		return &lloprotocol.Quote{
			Benchmark: results[0],
			Bid:       results[1],
			Ask:       results[2],
		}, nil
	default:
		return nil, fmt.Errorf("invalid number of results, expected: 1 or 3, got: %d", len(finaltrrs))
	}
}

func toDecimal(val any) (decimal.Decimal, error) {
	return utils.ToDecimal(val)
}

type MissingStreamError struct {
	StreamID streams.StreamID
}

func (e MissingStreamError) Error() string {
	return fmt.Sprintf("no pipeline for stream: %d", e.StreamID)
}

func (oc *observationContext) run(ctx context.Context, streamID streams.StreamID) (*execution, error) {
	p, exists := oc.r.Get(streamID)
	if !exists {
		return nil, MissingStreamError{StreamID: streamID}
	}

	// In case of multiple streamIDs per pipeline then the
	// first call executes and the others wait for result
	oc.executionsMu.Lock()
	ex, isExecuting := oc.executions[p]
	if isExecuting {
		oc.executionsMu.Unlock()
		// We intentionally do NOT select on ctx.Done() here.
		// BridgeTask uses overtimeContext (context.WithoutCancel) which
		// detaches from the caller's deadline, so p.Run can still be
		// in-flight after ctx expires. If waiters bail early via
		// ctx.Done(), they return an error while the executor may later
		// succeed. This results in some streams from the pipeline having values
		// while others do not. Blocking on ex.done ensures all goroutines for
		// the same pipeline receive the identical (run, trrs, err) tuple.
		<-ex.done
		return ex, ex.err
	}

	// execute here
	ch := make(chan struct{})
	ex = &execution{done: ch}
	oc.executions[p] = ex
	oc.executionsMu.Unlock()

	run, trrs, err := p.Run(ctx)
	ex.run = run
	ex.trrs = trrs
	if err == nil {
		// Resolve and validate once, here, so that every stream served by this
		// pipeline observes the same outcome. A quote invariant violation on any
		// stream fails the whole set.
		ex.vals = resolvePipelineStreamValues(p, streamID, trrs)
		err = validateStreamValues(ex.vals)
	}
	ex.err = err
	close(ch)

	return ex, err
}
