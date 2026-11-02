package pipeline

import (
	"context"
	stderrors "errors"
	"math"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Return types:
//
//	*decimal.Decimal
type MedianTask struct {
	BaseTask      `mapstructure:",squash"`
	Values        string `json:"values"`
	AllowedFaults string `json:"allowedFaults"`
	// Lax when disabled (default) will return an error if there are no values to medianize or if the input includes nil values.
	// Lax when enabled will return nil with no error if there are no valid values to medianize. If the input includes nil values, they will be excluded from the median calculation and do not count as a fault.
	Lax string
	// CountNilsAsFaults when enabled treats nil values as faults (counted toward allowedFaults) but filters them out before median calculation.
	// This is mutually exclusive with Lax.
	CountNilsAsFaults string
	// StalenessGateSeconds opts into the staleness gate: inputs whose source
	// timestamp trails the freshest input by more than this many seconds are
	// excluded from the median and count as faults (with a fallback to all
	// inputs when those exclusions would fail the allowedFaults check). An
	// explicit zero admits only the freshest input. The gate is off when the
	// parameter is absent.
	//
	// The gate can only tell inputs apart by their source timestamps, so it has
	// no effect when values="$(x)" references a single slice (every element
	// gets x's timestamp), or when every input derives from one shared
	// slow-moving dependency (every input inherits its timestamp).
	StalenessGateSeconds string `json:"stalenessGateSeconds"`
}

var _ Task = (*MedianTask)(nil)

func (t *MedianTask) Type() TaskType {
	return TaskTypeMedian
}

func (t *MedianTask) Run(_ context.Context, lggr logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		maybeAllowedFaults    MaybeUint64Param
		valuesAndErrs         SliceParam
		decimalValues         DecimalSliceParam
		allowedFaults         int
		faults                int
		lax                   BoolParam
		countNilsAsFaults     BoolParam
		maybeStalenessGateSec MaybeUint64Param
	)
	err := stderrors.Join(
		errors.Wrap(ResolveParam(&maybeAllowedFaults, From(t.AllowedFaults)), "allowedFaults"),
		errors.Wrap(ResolveParam(&valuesAndErrs, From(VarExpr(t.Values, vars), JSONWithVarExprs(t.Values, vars, true), Inputs(inputs))), "values"),
		errors.Wrap(ResolveParam(&lax, From(NonemptyString(t.Lax), false)), "lax"),
		errors.Wrap(ResolveParam(&countNilsAsFaults, From(NonemptyString(t.CountNilsAsFaults), false)), "countNilsAsFaults"),
		errors.Wrap(ResolveParam(&maybeStalenessGateSec, From(VarExpr(t.StalenessGateSeconds, vars), t.StalenessGateSeconds)), "stalenessGateSeconds"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	if bool(lax) && bool(countNilsAsFaults) {
		return Result{Error: errors.New("lax and countNilsAsFaults cannot both be enabled")}, runInfo
	}

	// Resolve per-element source timestamps while the element order is known,
	// so they can be carried through filtering in lockstep with the values.
	timestamps := resolveTimestamps(t.Values, vars, valuesAndErrs, inputs)

	// if lax is enabled, filter out nil values
	// nil values are not included in the fault calculations
	if bool(lax) {
		valuesAndErrs, timestamps, _ = filterNilsWithTimestamps(valuesAndErrs, timestamps)
	}

	if allowed, isSet := maybeAllowedFaults.Uint64(); isSet {
		if allowed > math.MaxInt {
			allowedFaults = math.MaxInt
		} else {
			allowedFaults = int(allowed) //nolint:gosec // G115: bounded by math.MaxInt above
		}
	} else {
		allowedFaults = max(len(valuesAndErrs)-1, 0)
	}

	values, timestamps, faults := filterErrorsWithTimestamps(valuesAndErrs, timestamps)

	// If countNilsAsFaults is enabled, filter nils AFTER fault counting
	// so that nils are counted toward allowedFaults
	if bool(countNilsAsFaults) {
		var nilCount int
		values, timestamps, nilCount = filterNilsWithTimestamps(values, timestamps)
		faults += nilCount
	}

	// The staleness gate is opt-in: unset or empty stalenessGateSeconds leaves
	// it off; an explicit zero means zero tolerance (only the freshest input
	// survives).
	stalenessGateSec, gateEnabled := maybeStalenessGateSec.Uint64()
	values, timestamps, faults = applyStalenessGate(lggr, values, timestamps, faults, allowedFaults, stalenessGateSec, gateEnabled)

	// The median is only as fresh as the inputs that went into it. Zero
	// (unknown) lets the scheduler fall back to the dependencies' timestamps.
	timestamp := oldestKnownTimestamp(timestamps)

	switch {
	case faults > allowedFaults:
		return Result{Error: errors.Wrapf(ErrTooManyErrors, "Number of faulty inputs %v to median task > number allowed faults %v", faults, allowedFaults)}, runInfo
	case len(values) == 0 && (bool(lax) || bool(countNilsAsFaults)):
		return Result{}, runInfo // if lax is enabled, return nil result with no error
	case len(values) == 0:
		return Result{Error: errors.Wrap(ErrWrongInputCardinality, "no values to medianize")}, runInfo
	}

	err = decimalValues.UnmarshalPipelineParam(values)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	sort.Slice(decimalValues, func(i, j int) bool {
		return decimalValues[i].LessThan(decimalValues[j])
	})
	k := len(decimalValues) / 2
	if len(decimalValues)%2 == 1 {
		return Result{Value: decimalValues[k], Timestamp: timestamp}, runInfo
	}
	median := decimalValues[k].Add(decimalValues[k-1]).Div(decimal.NewFromInt(2))
	return Result{Value: median, Timestamp: timestamp}, runInfo
}

// resolveTimestamps returns, for each element of the resolved values, the
// source timestamp (unix ms) of the task that produced it. It returns nil when
// timestamps cannot be mapped to every element (e.g. values resolved from
// arbitrary JSON), which disables the staleness gate for the run.
//
// The values expression is a comma-separated list of var expressions, each
// resolving to one element (or to a slice that is flattened into consecutive
// elements), so timestamps are looked up per referenced task via the dotID
// registry and assigned in expansion order. When the expression is empty,
// values came from the task inputs and correspond 1:1; a non-empty expression
// without var references is a literal with no source timestamps.
func resolveTimestamps(expr string, vars Vars, values SliceParam, inputs []Result) []uint64 {
	ts := make([]uint64, len(values))
	if strings.TrimSpace(expr) == "" {
		if len(inputs) != len(values) {
			return nil
		}
		for i := range inputs {
			ts[i] = inputs[i].Timestamp
		}
		return ts
	}
	refs := variableRegexp.FindAllStringSubmatch(expr, -1)
	if len(refs) == 0 {
		return nil
	}

	idx := 0
	for _, ref := range refs {
		if idx >= len(values) {
			break
		}
		keypath := strings.TrimSpace(ref[1])
		root, _, _ := strings.Cut(keypath, KeypathSeparator)
		tsRoot := vars.GetTimestamp(root)
		count := 1
		if v, err := vars.Get(keypath); err == nil {
			if s, is := v.([]any); is {
				count = len(s)
			}
		}
		for i := 0; i < count && idx < len(values); i++ {
			ts[idx] = tsRoot
			idx++
		}
	}
	if idx != len(values) {
		// could not map every value to a recorded source timestamp
		return nil
	}
	return ts
}

// filterErrorsWithTimestamps filters errored values, returning the surviving
// values and their timestamps in lockstep along with the error count.
func filterErrorsWithTimestamps(values SliceParam, ts []uint64) (SliceParam, []uint64, int) {
	outValues, outTs := make(SliceParam, 0, len(values)), make([]uint64, 0, len(values))
	errs := 0
	for i, x := range values {
		if _, is := x.(error); is {
			errs++
			continue
		}
		outValues = append(outValues, x)
		if len(ts) > i {
			outTs = append(outTs, ts[i])
		}
	}
	if len(ts) != len(values) {
		outTs = nil
	}
	return outValues, outTs, errs
}

// filterNilsWithTimestamps filters nil values, returning the surviving values
// and their timestamps in lockstep along with the nil count.
func filterNilsWithTimestamps(values SliceParam, ts []uint64) (SliceParam, []uint64, int) {
	outValues, outTs := make(SliceParam, 0, len(values)), make([]uint64, 0, len(values))
	nils := 0
	for i, x := range values {
		if x == nil {
			nils++
			continue
		}
		outValues = append(outValues, x)
		if len(ts) > i {
			outTs = append(outTs, ts[i])
		}
	}
	if len(ts) != len(values) {
		outTs = nil
	}
	return outValues, outTs, nils
}

// applyStalenessGate excludes values whose source timestamp is older than the
// freshest input by more than stalenessGateSec seconds (the maximum
// staleness), counting them as faults, and returns the surviving values and
// their timestamps in lockstep. Inputs with an unknown (zero) timestamp are
// never excluded. It is a no-op when the gate is not enabled.
//
// When excluding every stale value would fail the fault-threshold check, only
// the stalest values are excluded, oldest first, as many as the threshold
// still allows (staleness fallback). Equally stale values that cannot all be
// excluded are excluded in their order in values, earliest first.
func applyStalenessGate(lggr logger.Logger, values SliceParam, ts []uint64, faults, allowedFaults int, stalenessGateSec uint64, enabled bool) (SliceParam, []uint64, int) {
	if !enabled || len(ts) != len(values) {
		// Gate not opted into, or timestamps unknown for this run.
		return values, ts, faults
	}

	// Guard against overflow: gate seconds near the uint64 maximum mean an
	// effectively unbounded tolerance, so nothing is ever excluded.
	stalenessGateMs := uint64(math.MaxUint64)
	if stalenessGateSec <= math.MaxUint64/1_000 {
		stalenessGateMs = stalenessGateSec * 1_000
	}
	var freshest uint64
	for _, ms := range ts {
		if ms > freshest {
			freshest = ms
		}
	}
	if freshest == 0 {
		return values, ts, faults
	}

	var stale []int
	for i, ms := range ts {
		if ms != 0 && freshest-ms > stalenessGateMs {
			stale = append(stale, i)
		}
	}
	if len(stale) == 0 {
		return values, ts, faults
	}

	if budget := max(allowedFaults-faults, 0); len(stale) > budget {
		// Staleness fallback: exclude only the stalest values the threshold
		// still allows. The stable sort keeps ties in values order.
		logger.Sugared(lggr).Debugw("Median task: staleness gate would exceed allowed faults, excluding only the stalest inputs",
			"stale", len(stale),
			"excluded", budget,
			"faults", faults,
			"allowedFaults", allowedFaults,
		)
		sort.SliceStable(stale, func(i, j int) bool { return ts[stale[i]] < ts[stale[j]] })
		stale = stale[:budget]
	}
	if len(stale) == 0 {
		return values, ts, faults
	}

	excluded := make(map[int]struct{}, len(stale))
	for _, i := range stale {
		excluded[i] = struct{}{}
	}
	outValues := make(SliceParam, 0, len(values)-len(stale))
	outTs := make([]uint64, 0, len(values)-len(stale))
	for i, x := range values {
		if _, is := excluded[i]; is {
			continue
		}
		outValues = append(outValues, x)
		outTs = append(outTs, ts[i])
	}
	return outValues, outTs, faults + len(stale)
}

// oldestKnownTimestamp returns the oldest non-zero timestamp, or zero when
// none is known.
func oldestKnownTimestamp(ts []uint64) uint64 {
	var oldest uint64
	for _, ms := range ts {
		if ms != 0 && (oldest == 0 || ms < oldest) {
			oldest = ms
		}
	}
	return oldest
}
