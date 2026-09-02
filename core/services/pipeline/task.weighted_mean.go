package pipeline

import (
	"context"
	stderrors "errors"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var ErrInsufficientWeightMass = errors.New("insufficient weight mass")

// Experimental: WeightedMeanTask gates on minimum weight mass, clamps each
// sample's value into a band around the cross-sectional median (when method
// is a winsor mode), then returns the weighted mean. By default (no method
// specified) it computes a plain weighted mean with no clamping.
//
// Input:  samples ([]Sample) — staleness-decayed samples, one per source
// Input:  reference ([]Sample, optional) — samples to derive the clamp band from
// Output: decimal.Decimal — the weighted mean (winsorized if method is set)
// Fails:  ErrInsufficientWeightMass if surviving mass < minWeightMass
// Fails:  ErrWrongInputCardinality if no samples survive
//
// Method values:
//
//	(omitted)    — plain weighted mean, no clamping (default)
//	winsor       — winsorize: clamp to [m×(1-band), m×(1+band)]
//	winsor_abs   — winsorize: clamp to [m-band, m+band]
//	winsor_mad   — winsorize: clamp to [m-band×MAD, m+band×MAD]
type WeightedMeanTask struct {
	BaseTask      `mapstructure:",squash"`
	Samples       string `json:"samples"`
	Reference     string `json:"reference"`
	Band          string `json:"band"`
	MinWeightMass string `json:"minWeightMass"`
	Method        string `json:"method"`
	WinsorRef     string `json:"winsorRef"`
	Precision     string `json:"precision"`
}

var _ Task = (*WeightedMeanTask)(nil)

func (t *WeightedMeanTask) Type() TaskType {
	return TaskTypeWeightedMean
}

func (t *WeightedMeanTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		samplesAndErrs   SliceParam
		referenceAndErrs SliceParam
		band             DecimalParam
		minMass          DecimalParam
		method           StringParam
		winsorRef        StringParam
		maybePrecision   MaybeInt32Param
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
		resolveOpt(&referenceAndErrs, VarExpr(t.Reference, vars), JSONWithVarExprs(t.Reference, vars, true)),
		errors.Wrap(ResolveParam(&band, From(VarExpr(t.Band, vars), NonemptyString(t.Band), "0.03")), "band"),
		errors.Wrap(ResolveParam(&minMass, From(VarExpr(t.MinWeightMass, vars), NonemptyString(t.MinWeightMass), "0")), "minWeightMass"),
		errors.Wrap(ResolveParam(&method, From(NonemptyString(t.Method), "")), "method"),
		errors.Wrap(ResolveParam(&winsorRef, From(NonemptyString(t.WinsorRef), "median")), "winsorRef"),
		errors.Wrap(ResolveParam(&maybePrecision, From(VarExpr(t.Precision, vars), t.Precision)), "precision"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	samplesRaw, _ := samplesAndErrs.FilterErrors()
	var samples SampleSliceParam
	if err := samples.UnmarshalPipelineParam(samplesRaw); err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "samples: %v", err)}, runInfo
	}

	// Drop zero-weight samples (staleness already set weight=0 and dropped them,
	// but this is defensive against upstream tasks that don't filter).
	active := make([]Sample, 0, len(samples))
	for _, s := range samples {
		if !s.Weight.IsZero() {
			active = append(active, s)
		}
	}
	if len(active) == 0 {
		return Result{Error: errors.Wrap(ErrWrongInputCardinality, "no surviving samples")}, runInfo
	}

	// Gate: sum the weight of all surviving samples. If total mass is below
	// minWeightMass, the round fails — too many sources were missing or stale.
	totalMass := decimal.Zero
	for _, s := range active {
		totalMass = totalMass.Add(s.Weight)
	}
	mm := minMass.Decimal()
	if mm.GreaterThan(decimal.Zero) && totalMass.LessThan(mm) {
		return Result{Error: errors.Wrapf(ErrInsufficientWeightMass, "mass %s < min %s", totalMass.String(), mm.String())}, runInfo
	}

	methodMode := strings.ToLower(string(method))

	num := decimal.Zero
	den := decimal.Zero

	if methodMode == "" || methodMode == "none" {
		// Plain weighted mean — no clamping. Skip band computation entirely.
		for _, s := range active {
			num = num.Add(s.Weight.Mul(s.Value))
			den = den.Add(s.Weight)
		}
	} else {
		// Compute reference values for the band. Defaults to the task's own samples.
		ref := active
		if len(referenceAndErrs) > 0 {
			refRaw, _ := referenceAndErrs.FilterErrors()
			var refSamples SampleSliceParam
			if err := refSamples.UnmarshalPipelineParam(refRaw); err != nil {
				return Result{Error: errors.Wrapf(ErrBadInput, "reference: %v", err)}, runInfo
			}
			ref = make([]Sample, 0, len(refSamples))
			for _, s := range refSamples {
				if !s.Weight.IsZero() {
					ref = append(ref, s)
				}
			}
			if len(ref) == 0 {
				return Result{Error: errors.Wrap(ErrWrongInputCardinality, "no reference values to compute band")}, runInfo
			}
		}

		m := bandReference(ref, strings.ToLower(string(winsorRef)))

		b := band.Decimal()
		lo, hi, err := computeBand(m, b, methodMode, ref)
		if err != nil {
			return Result{Error: err}, runInfo
		}

		// Clamp each sample's value into [lo, hi] around the reference median,
		// then accumulate the weighted sum. This bounds how far one corrupted
		// source can move the output: at most omega_k * 2*band*median.
		for _, s := range active {
			q := s.Value
			if q.LessThan(lo) {
				q = lo
			}
			if q.GreaterThan(hi) {
				q = hi
			}
			num = num.Add(s.Weight.Mul(q))
			den = den.Add(s.Weight)
		}
	}
	if den.IsZero() {
		return Result{Error: errors.Wrap(ErrWrongInputCardinality, "zero weight mass")}, runInfo
	}

	var out decimal.Decimal
	if prec, isSet := maybePrecision.Int32(); isSet {
		out = num.DivRound(den, prec)
	} else {
		out = num.Div(den)
	}
	return Result{Value: out}, runInfo
}

func bandReference(ref []Sample, method string) decimal.Decimal {
	switch method {
	case "weighted_median":
		return weightedMedian(ref)
	default:
		vals := make([]decimal.Decimal, len(ref))
		for i, r := range ref {
			vals[i] = r.Value
		}
		return medianDecimal(vals)
	}
}

func weightedMedian(samples []Sample) decimal.Decimal {
	sorted := make([]Sample, len(samples))
	copy(sorted, samples)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value.LessThan(sorted[j].Value)
	})
	total := decimal.Zero
	for _, s := range sorted {
		total = total.Add(s.Weight)
	}
	half := total.Div(decimal.NewFromInt(2))
	cum := decimal.Zero
	for _, s := range sorted {
		cum = cum.Add(s.Weight)
		if cum.GreaterThanOrEqual(half) {
			return s.Value
		}
	}
	if len(sorted) > 0 {
		return sorted[len(sorted)-1].Value
	}
	return decimal.Zero
}

func computeBand(m, b decimal.Decimal, method string, ref []Sample) (lo, hi decimal.Decimal, err error) {
	one := decimal.NewFromInt(1)
	switch method {
	case "winsor":
		lo = m.Mul(one.Sub(b))
		hi = m.Mul(one.Add(b))
		if lo.IsNegative() {
			lo = decimal.Zero
		}
	case "winsor_abs":
		lo = m.Sub(b)
		hi = m.Add(b)
	case "winsor_mad":
		vals := make([]decimal.Decimal, len(ref))
		for i, r := range ref {
			vals[i] = r.Value.Sub(m).Abs()
		}
		mad := medianDecimal(vals)
		lo = m.Sub(b.Mul(mad))
		hi = m.Add(b.Mul(mad))
	default:
		return decimal.Zero, decimal.Zero, errors.Errorf("unknown method %q", method)
	}
	return lo, hi, nil
}

func medianDecimal(vals []decimal.Decimal) decimal.Decimal {
	n := len(vals)
	if n == 0 {
		return decimal.Zero
	}
	sorted := make([]decimal.Decimal, n)
	copy(sorted, vals)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].LessThan(sorted[j]) })
	if n%2 == 1 {
		return sorted[n/2]
	}
	return sorted[n/2].Add(sorted[n/2-1]).Div(decimal.NewFromInt(2))
}
