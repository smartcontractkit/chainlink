package pipeline

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// DurationParam parses a duration string (e.g. "10s") into a time.Duration.
type DurationParam time.Duration

func (d *DurationParam) UnmarshalPipelineParam(val any) error {
	switch v := val.(type) {
	case string:
		if v == "" {
			return ErrParameterEmpty
		}
		dur, err := time.ParseDuration(v)
		if err != nil {
			return errors.Wrap(ErrBadInput, err.Error())
		}
		*d = DurationParam(dur)
		return nil
	case time.Duration:
		*d = DurationParam(v)
		return nil
	default:
		return errors.Wrapf(ErrBadInput, "expected duration, got %T", val)
	}
}

// Experimental: StalenessTask scales each sample's weight by a decay function of its age
// (now - ts_ms). Samples whose weight decays to zero are dropped.
//
// Input:  samples ([]Sample) — fresh from sample or normalize
// Output: []Sample — same samples with Weight replaced by W * f(age); w==0 dropped
// Fails:  if method is unknown, or halfLife is missing for exp/exp_cooldown
//
// Methods:
//
//	cutoff    — 1 if age <= threshold, else 0 (binary, discontinuous)
//	linear    — 1 - age/threshold, reaches 0 at threshold
//	exp       — 2^(-age/halfLife), truncated to 0 at threshold
//	exp_cooldown  — 1 if age <= threshold, else 2^(-(age-threshold)/halfLife)
//
// Optional safety parameters:
//
//	decayThreshold (K) — when the decay multiplier falls below K, the sample is
//	dropped entirely (assigned "no data" state). Recommended value: 0.03.
//	Applies to exp and exp_cooldown methods (where decay is asymptotic).
//
//	cutoff — hard time limit; when age exceeds cutoff, the sample is dropped
//	regardless of the decay function value. Acts as a safety valve.
//
// ErrStaleSample marks a sample excluded by the staleness task.
var ErrStaleSample = errors.New("stale sample")

// StaleSampleError is emitted in place of an excluded sample when onStale=fault,
// so that downstream aggregators count it as a fault. It serializes with the
// sample's source and timestamps so persisted runs show what was cut.
type StaleSampleError struct {
	Source         string
	TsMs           int64
	ReferenceTsMs  int64
	ReferenceIsNow bool
}

func (e StaleSampleError) Error() string {
	ref := "freshest sample"
	if e.ReferenceIsNow {
		ref = "now"
	}
	return fmt.Sprintf("%s %q: ts %d is %dms behind %s (%d)", ErrStaleSample, e.Source, e.TsMs, e.ReferenceTsMs-e.TsMs, ref, e.ReferenceTsMs)
}

func (e StaleSampleError) Unwrap() error { return ErrStaleSample }

func (e StaleSampleError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Stale         bool   `json:"stale"`
		Source        string `json:"source"`
		TsMs          int64  `json:"ts_ms"`
		ReferenceTsMs int64  `json:"reference_ts_ms"`
	}{true, e.Source, e.TsMs, e.ReferenceTsMs})
}

// Experimental: StalenessTask decays sample weights by age and excludes stale samples.
//
// Input:  samples ([]Sample) — may contain upstream errors
// Output: []Sample (onStale=drop) or []any of Samples and faults (onStale=fault)
//
// relativeTo selects the age reference: "now" (default) is the node's wall
// clock; "freshest" is the newest known sample timestamp, in which case samples
// with an unknown (zero) timestamp are never excluded.
//
// onStale selects what happens to excluded samples: "drop" (default) removes
// them and discards upstream errors; "fault" emits a StaleSampleError in their
// place and passes upstream errors through, so a downstream median counts both
// against its allowedFaults.
//
// allowedFaults caps exclusions: upstream errors spend the budget first, then
// samples are excluded stalest first (ties in input order) until it runs out.
// Samples that cannot be excluded pass through with their decayed weight.
// Absent means no cap. Set it to the downstream aggregator's allowedFaults.
type StalenessTask struct {
	BaseTask       `mapstructure:",squash"`
	Samples        string `json:"samples"`
	Method         string `json:"method"`
	Threshold      string `json:"threshold"`
	HalfLife       string `json:"halfLife"`
	Cutoff         string `json:"cutoff"`
	DecayThreshold string `json:"decayThreshold"`
	RelativeTo     string `json:"relativeTo"`
	OnStale        string `json:"onStale"`
	AllowedFaults  string `json:"allowedFaults"`
}

var _ Task = (*StalenessTask)(nil)

func (t *StalenessTask) Type() TaskType {
	return TaskTypeStaleness
}

func (t *StalenessTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		samplesAndErrs     SliceParam
		method             StringParam
		threshold          DurationParam
		halfLife           DurationParam
		cutoff             DurationParam
		decayThreshold     DecimalParam
		relativeTo         StringParam
		onStale            StringParam
		maybeAllowedFaults MaybeUint64Param
	)

	resolveOpt := func(out PipelineParamUnmarshaler, getters ...GetterFunc) error {
		err := ResolveParam(out, getters)
		if err == nil || errors.Is(err, ErrParameterEmpty) {
			return nil
		}
		return err
	}

	err := stderrors.Join(
		errors.Wrap(ResolveParam(&samplesAndErrs, From(VarExpr(t.Samples, vars), JSONWithVarExprs(t.Samples, vars, true), Inputs(inputs))), "samples"),
		errors.Wrap(ResolveParam(&method, From(NonemptyString(t.Method))), "method"),
		resolveOpt(&threshold, NonemptyString(t.Threshold)),
		errors.Wrap(ResolveParam(&halfLife, From(NonemptyString(t.HalfLife), "0s")), "halfLife"),
		resolveOpt(&cutoff, VarExpr(t.Cutoff, vars), NonemptyString(t.Cutoff)),
		resolveOpt(&decayThreshold, VarExpr(t.DecayThreshold, vars), NonemptyString(t.DecayThreshold)),
		errors.Wrap(ResolveParam(&relativeTo, From(NonemptyString(t.RelativeTo), "now")), "relativeTo"),
		errors.Wrap(ResolveParam(&onStale, From(NonemptyString(t.OnStale), "drop")), "onStale"),
		errors.Wrap(ResolveParam(&maybeAllowedFaults, From(t.AllowedFaults)), "allowedFaults"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	m := strings.ToLower(string(method))
	cutoffVal := time.Duration(cutoff)
	kVal := decayThreshold.Decimal()
	applyK := (m == "exp" || m == "exp_cooldown") && kVal.GreaterThan(decimal.Zero)

	if time.Duration(threshold) <= 0 {
		return Result{Error: errors.New("threshold required for staleness")}, runInfo
	}
	if (m == "exp" || m == "exp_cooldown") && halfLife <= 0 {
		return Result{Error: errors.New("halfLife required for exp/exp_cooldown staleness")}, runInfo
	}

	var relativeToFreshest bool
	switch strings.ToLower(string(relativeTo)) {
	case "now":
	case "freshest":
		relativeToFreshest = true
	default:
		return Result{Error: errors.Errorf("relativeTo: unknown value %q, expected now or freshest", relativeTo)}, runInfo
	}
	var faultOnStale bool
	switch strings.ToLower(string(onStale)) {
	case "drop":
	case "fault":
		faultOnStale = true
	default:
		return Result{Error: errors.Errorf("onStale: unknown value %q, expected drop or fault", onStale)}, runInfo
	}

	// Parse each element separately so upstream errors keep their position.
	type entry struct {
		sample   Sample
		err      error
		excluded bool
	}
	var entries []entry
	var upstreamFaults int
	for _, x := range samplesAndErrs {
		if e, is := x.(error); is {
			upstreamFaults++
			entries = append(entries, entry{err: e})
			continue
		}
		var parsed SampleSliceParam
		if err := parsed.UnmarshalPipelineParam([]any{x}); err != nil {
			return Result{Error: errors.Wrapf(ErrBadInput, "samples: %v", err)}, runInfo
		}
		for _, sample := range parsed {
			entries = append(entries, entry{sample: sample})
		}
	}

	now := time.Now()
	referenceMs := now.UnixMilli()
	referenceNs := now.UnixNano()
	if relativeToFreshest {
		referenceMs = 0
		for _, e := range entries {
			if e.err == nil && e.sample.TsMs != 0 {
				referenceMs = max(referenceMs, e.sample.TsMs)
			}
		}
		referenceNs = referenceMs * int64(time.Millisecond)
	}

	var excluded []int
	for i := range entries {
		e := &entries[i]
		if e.err != nil || (relativeToFreshest && e.sample.TsMs == 0) {
			continue
		}
		ageNs := max(referenceNs-e.sample.TsMs*int64(time.Millisecond), 0)

		mult, err := decayMultiplier(m, time.Duration(ageNs), time.Duration(threshold), time.Duration(halfLife))
		if err != nil {
			return Result{Error: err}, runInfo
		}
		e.sample.Weight = e.sample.Weight.Mul(mult)

		if (cutoffVal > 0 && ageNs > int64(cutoffVal)) || (applyK && mult.LessThan(kVal)) || mult.IsZero() {
			e.excluded = true
			excluded = append(excluded, i)
		}
	}

	if allowed, isSet := maybeAllowedFaults.Uint64(); isSet {
		allowedFaults := math.MaxInt
		if allowed <= math.MaxInt {
			allowedFaults = int(allowed)
		}
		budget := min(max(allowedFaults-upstreamFaults, 0), len(excluded))
		// Keep the budget's worth of stalest exclusions; the rest pass through.
		sort.SliceStable(excluded, func(a, b int) bool {
			return entries[excluded[a]].sample.TsMs < entries[excluded[b]].sample.TsMs
		})
		for _, i := range excluded[budget:] {
			entries[i].excluded = false
		}
	}

	if !faultOnStale {
		out := make([]Sample, 0, len(entries))
		for _, e := range entries {
			if e.err == nil && !e.excluded {
				out = append(out, e.sample)
			}
		}
		return Result{Value: out}, runInfo
	}

	out := make([]any, 0, len(entries))
	for _, e := range entries {
		switch {
		case e.err != nil:
			out = append(out, e.err)
		case e.excluded:
			out = append(out, StaleSampleError{
				Source:         e.sample.Source,
				TsMs:           e.sample.TsMs,
				ReferenceTsMs:  referenceMs,
				ReferenceIsNow: !relativeToFreshest,
			})
		default:
			out = append(out, e.sample)
		}
	}
	return Result{Value: out}, runInfo
}

func decayMultiplier(method string, age, threshold, halfLife time.Duration) (decimal.Decimal, error) {
	switch method {
	case "cutoff":
		if age <= threshold {
			return decimal.NewFromInt(1), nil
		}
		return decimal.Zero, nil
	case "linear":
		if age <= 0 {
			return decimal.NewFromInt(1), nil
		}
		if age >= threshold {
			return decimal.Zero, nil
		}
		return decimal.NewFromInt(int64(threshold - age)).Div(decimal.NewFromInt(int64(threshold))), nil
	case "exp":
		return decimal.NewFromFloat(math.Pow(2, -float64(age)/float64(halfLife))), nil
	case "exp_cooldown":
		if age <= threshold {
			return decimal.NewFromInt(1), nil
		}
		return decimal.NewFromFloat(math.Pow(2, -float64(age-threshold)/float64(halfLife))), nil
	default:
		return decimal.Zero, errors.Errorf("unknown staleness method %q", method)
	}
}
