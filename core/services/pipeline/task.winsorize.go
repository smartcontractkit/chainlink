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

// Experimental: WinsorizedMeanTask collapses samples by source, gates on minimum weight mass,
// clamps each source's value into a band around the cross-sectional median, then
// returns the weighted mean.
//
// Input:  samples ([]Sample) — staleness-decayed samples from one or more sources
// Input:  reference ([]Sample, optional) — samples to derive the clamp band from
// Output: decimal.Decimal — the winsorized weighted mean
// Fails:  ErrInsufficientWeightMass if surviving mass < minWeightMass
// Fails:  ErrWrongInputCardinality if no sources survive collapse
type WinsorizedMeanTask struct {
	BaseTask      `mapstructure:",squash"`
	Samples       string `json:"samples"`
	Reference     string `json:"reference"`
	Band          string `json:"band"`
	MinWeightMass string `json:"minWeightMass"`
	Winsor        string `json:"winsor"`
	WinsorRef     string `json:"winsorRef"`
	Precision     string `json:"precision"`
}

var _ Task = (*WinsorizedMeanTask)(nil)

func (t *WinsorizedMeanTask) Type() TaskType {
	return TaskTypeWinsorizedMean
}

func (t *WinsorizedMeanTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		samplesAndErrs    SliceParam
		referenceAndErrs  SliceParam
		band              DecimalParam
		minMass           DecimalParam
		winsor            StringParam
		winsorRef         StringParam
		maybePrecision    MaybeInt32Param
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
		errors.Wrap(ResolveParam(&winsor, From(NonemptyString(t.Winsor), "rel")), "winsor"),
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

	collapsed := collapseSamples(samples)
	if len(collapsed) == 0 {
		return Result{Error: errors.Wrap(ErrWrongInputCardinality, "no surviving sources")}, runInfo
	}

	// Gate: sum the weight of all surviving sources. If total mass is below
	// minWeightMass, the round fails — too many sources were missing or stale.
	totalMass := decimal.Zero
	for _, c := range collapsed {
		totalMass = totalMass.Add(c.weight)
	}
	mm := minMass.Decimal()
	if mm.GreaterThan(decimal.Zero) && totalMass.LessThan(mm) {
		return Result{Error: errors.Wrapf(ErrInsufficientWeightMass, "mass %s < min %s", totalMass.String(), mm.String())}, runInfo
	}

	// Compute reference values for the band. Defaults to the task's own samples.
	ref := collapsed
	if len(referenceAndErrs) > 0 {
		refRaw, _ := referenceAndErrs.FilterErrors()
		var refSamples SampleSliceParam
		if err := refSamples.UnmarshalPipelineParam(refRaw); err != nil {
			return Result{Error: errors.Wrapf(ErrBadInput, "reference: %v", err)}, runInfo
		}
		ref = collapseSamples(refSamples)
		if len(ref) == 0 {
			return Result{Error: errors.Wrap(ErrWrongInputCardinality, "no reference values to compute band")}, runInfo
		}
	}

	m := bandReference(ref, strings.ToLower(string(winsorRef)))
	if m.IsZero() && len(ref) > 0 {
		// A zero median is legal only if all values are zero. Keep it.
	}

	b := band.Decimal()
	lo, hi, err := computeBand(m, b, strings.ToLower(string(winsor)), ref)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	// Clamp each source's value into [lo, hi] around the reference median,
	// then accumulate the weighted sum. This bounds how far one corrupted
	// source can move the output: at most omega_k * 2*band*median.
	num := decimal.Zero
	den := decimal.Zero
	for _, c := range collapsed {
		q := c.value
		if q.LessThan(lo) {
			q = lo
		}
		if q.GreaterThan(hi) {
			q = hi
		}
		num = num.Add(c.weight.Mul(q))
		den = den.Add(c.weight)
	}
	if den.IsZero() {
		return Result{Error: errors.Wrap(ErrWrongInputCardinality, "zero weight mass")}, runInfo
	}

	var out decimal.Decimal
	if prec, isSet := maybePrecision.Int32(); isSet {
		out = num.DivRound(den, int32(prec))
	} else {
		out = num.Div(den)
	}
	return Result{Value: out}, runInfo
}

type collapsedSample struct {
	source string
	value  decimal.Decimal
	weight decimal.Decimal
}

func collapseSamples(samples []Sample) []collapsedSample {
	type group struct {
		values []decimal.Decimal
		weight decimal.Decimal
	}
	groups := make(map[string]*group)
	for _, s := range samples {
		if s.Weight.IsZero() {
			continue
		}
		g, ok := groups[s.Source]
		if !ok {
			g = &group{}
			groups[s.Source] = g
		}
		g.values = append(g.values, s.Value)
		if s.Weight.GreaterThan(g.weight) {
			g.weight = s.Weight
		}
	}
	out := make([]collapsedSample, 0, len(groups))
	for src, g := range groups {
		out = append(out, collapsedSample{
			source: src,
			value:  medianDecimal(g.values),
			weight: g.weight,
		})
	}
	return out
}

func bandReference(ref []collapsedSample, method string) decimal.Decimal {
	switch method {
	case "weighted_median":
		return weightedMedian(ref)
	default:
		vals := make([]decimal.Decimal, len(ref))
		for i, r := range ref {
			vals[i] = r.value
		}
		return medianDecimal(vals)
	}
}

func weightedMedian(samples []collapsedSample) decimal.Decimal {
	sorted := make([]collapsedSample, len(samples))
	copy(sorted, samples)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].value.LessThan(sorted[j].value)
	})
	total := decimal.Zero
	for _, s := range sorted {
		total = total.Add(s.weight)
	}
	half := total.Div(decimal.NewFromInt(2))
	cum := decimal.Zero
	for _, s := range sorted {
		cum = cum.Add(s.weight)
		if cum.GreaterThanOrEqual(half) {
			return s.value
		}
	}
	if len(sorted) > 0 {
		return sorted[len(sorted)-1].value
	}
	return decimal.Zero
}

func computeBand(m, b decimal.Decimal, winsor string, ref []collapsedSample) (lo, hi decimal.Decimal, err error) {
	one := decimal.NewFromInt(1)
	switch winsor {
	case "rel":
		lo = m.Mul(one.Sub(b))
		hi = m.Mul(one.Add(b))
		if lo.IsNegative() {
			lo = decimal.Zero
		}
	case "abs":
		lo = m.Sub(b)
		hi = m.Add(b)
	case "mad":
		vals := make([]decimal.Decimal, len(ref))
		for i, r := range ref {
			vals[i] = r.value.Sub(m).Abs()
		}
		mad := medianDecimal(vals)
		lo = m.Sub(b.Mul(mad))
		hi = m.Add(b.Mul(mad))
	default:
		return decimal.Zero, decimal.Zero, errors.Errorf("unknown winsor band %q", winsor)
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
