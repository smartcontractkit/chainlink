package observation

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-data-streams/llo"

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

type ObservationContext interface {
	Observe(ctx context.Context, streamID streams.StreamID, opts llo.DSOpts) (val llo.StreamValue, err error)
}

type execution struct {
	done <-chan struct{}

	run  *pipeline.Run
	trrs pipeline.TaskRunResults
	err  error
}

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

func (oc *observationContext) Observe(ctx context.Context, streamID streams.StreamID, opts llo.DSOpts) (val llo.StreamValue, err error) {
	run, trrs, err := oc.run(ctx, streamID)
	observationFinishedAt := time.Now()
	if err != nil {
		// FIXME: This is a hack specific for V3 telemetry, future schemas should
		// use a generic stream value telemetry instead
		// https://smartcontract-it.atlassian.net/browse/MERC-6290
		oc.t.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, val, err)
		return nil, err
	}
	// Extract stream value based on streamID attribute
	found := false
	for _, trr := range trrs {
		if trr.Task.TaskStreamID() != nil && *trr.Task.TaskStreamID() == streamID {
			val, err = resultToStreamValue(trr.Result.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to convert result to StreamValue for streamID %d: %w", streamID, err)
			}
			if trr.FinishedAt.Valid {
				observationFinishedAt = trr.FinishedAt.Time
			}
			found = true
			break
		}
	}
	if !found {
		// If no streamID attribute is found in the task results, then assume the
		// final output is the stream ID and return that. This is safe to do since
		// the registry will never return a spec that doesn't match either by tag
		// or by spec streamID.
		val, err = extractFinalResultAsStreamValue(trrs)
		if oc.t.CaptureEATelemetry() {
			// FIXME: This is a hack specific for V3 telemetry, future schemas should
			// use the generic stream value telemetry instead
			// https://smartcontract-it.atlassian.net/browse/MERC-6290
			oc.t.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, val, err)
		}
	}
	if ch := GetObservationTelemetryCh(ctx); ch != nil {
		cd := opts.ConfigDigest()
		ot := &telem.LLOObservationTelemetry{
			StreamId:              streamID,
			ObservationTimestamp:  opts.ObservationTimestamp().UnixNano(),
			ObservationFinishedAt: observationFinishedAt.UnixNano(),
			SeqNr:                 opts.SeqNr(),
			ConfigDigest:          cd[:],
		}
		if err != nil {
			ot.ObservationError = new(string)
			*ot.ObservationError = err.Error()
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
	return
}

func resultToStreamValue(val interface{}) (llo.StreamValue, error) {
	switch v := val.(type) {
	case decimal.Decimal:
		return llo.ToDecimal(v), nil
	case float64:
		return llo.ToDecimal(decimal.NewFromFloat(v)), nil
	case int64:
		return llo.ToDecimal(decimal.NewFromInt(v)), nil
	case pipeline.ObjectParam:
		switch v.Type {
		case pipeline.DecimalType:
			return llo.ToDecimal(decimal.Decimal(v.DecimalValue)), nil
		default:
			return nil, fmt.Errorf("don't know how to convert pipeline.ObjectParam with type %d to llo.StreamValue", v.Type)
		}
	case map[string]interface{}:
		sv, err := resultMapToStreamValue(v)
		if err != nil {
			return nil, fmt.Errorf("don't know how to convert map to StreamValue: %w; got: %v", err, v)
		}
		return sv, nil
	default:
		return nil, fmt.Errorf("don't know how to convert pipeline output result of type %T to llo.StreamValue (got: %v)", val, val)
	}
}

// Converts arbitrary JSON (parsed to map) to a StreamValue
func resultMapToStreamValue(m map[string]interface{}) (llo.StreamValue, error) {
	var streamValueType llo.LLOStreamValue_Type
	{
		raw, exists := m["streamValueType"]
		if !exists {
			return nil, errors.New("expected a key labeled 'streamValueType' in the map")
		}
		rawInt64, ok := raw.(int64)
		if !ok {
			return nil, fmt.Errorf("expected 'streamValueType' to be a int64, got: %T", raw)
		}
		if rawInt64 < 0 || rawInt64 > math.MaxUint32 || rawInt64 >= int64(llo.LLOStreamValue_Type(len(llo.LLOStreamValue_Type_name))) { //nolint:gosec // G115 // won't overflow
			return nil, fmt.Errorf("invalid streamValueType: %v", rawInt64)
		}
		streamValueType = llo.LLOStreamValue_Type(rawInt64) //nolint:gosec // G115 // won't overflow due to check above
	}
	switch streamValueType {
	case llo.LLOStreamValue_TimestampedStreamValue:
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
//	  "timestamps": {
//	  	"providerIndicatedTimeUnixMs": 1234567890,
//	  	"providerDataReceivedUnixMs": 1234567890
//	  },
//	  "result": "123.456"
//	}
func resultMapToTimestampedStreamValue(m map[string]interface{}) (*llo.TimestampedStreamValue, error) {
	ts, ok := m["timestamps"].(map[string]interface{})
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

	return &llo.TimestampedStreamValue{
		ObservedAtNanoseconds: observedAtMillis * 1e6, // convert ms to ns
		StreamValue:           llo.ToDecimal(svd),
	}, nil
}

func toUint64(v interface{}) (uint64, error) {
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
func extractFinalResultAsStreamValue(trrs pipeline.TaskRunResults) (llo.StreamValue, error) {
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
		return llo.ToDecimal(val), nil
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
		return &llo.Quote{
			Benchmark: results[0],
			Bid:       results[1],
			Ask:       results[2],
		}, nil
	default:
		return nil, fmt.Errorf("invalid number of results, expected: 1 or 3, got: %d", len(finaltrrs))
	}
}

func toDecimal(val interface{}) (decimal.Decimal, error) {
	return utils.ToDecimal(val)
}

type MissingStreamError struct {
	StreamID streams.StreamID
}

func (e MissingStreamError) Error() string {
	return fmt.Sprintf("no pipeline for stream: %d", e.StreamID)
}

func (oc *observationContext) run(ctx context.Context, streamID streams.StreamID) (*pipeline.Run, pipeline.TaskRunResults, error) {
	p, exists := oc.r.Get(streamID)
	if !exists {
		return nil, nil, MissingStreamError{StreamID: streamID}
	}

	// In case of multiple streamIDs per pipeline then the
	// first call executes and the others wait for result
	oc.executionsMu.Lock()
	ex, isExecuting := oc.executions[p]
	if isExecuting {
		oc.executionsMu.Unlock()
		// wait for it to finish
		select {
		case <-ex.done:
			return ex.run, ex.trrs, ex.err
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}

	// execute here
	ch := make(chan struct{})
	ex = &execution{done: ch}
	oc.executions[p] = ex
	oc.executionsMu.Unlock()

	run, trrs, err := p.Run(ctx)
	ex.run = run
	ex.trrs = trrs
	ex.err = err
	close(ch)

	return run, trrs, err
}
