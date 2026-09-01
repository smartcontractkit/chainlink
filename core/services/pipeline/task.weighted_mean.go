package pipeline

import (
	"context"
	stderrors "errors"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Experimental: WeightedMeanTask collapses samples by source, gates on minimum weight mass,
// then returns the weighted mean. No clamping — use this for bid/ask where
// winsorization would destroy spread information.
//
// Input:  samples ([]Sample) — staleness-decayed samples from one or more sources
// Output: decimal.Decimal — the weighted mean (Sum w*q / Sum w)
// Fails:  ErrInsufficientWeightMass if surviving mass < minWeightMass
// Fails:  ErrWrongInputCardinality if no sources survive collapse
type WeightedMeanTask struct {
	BaseTask      `mapstructure:",squash"`
	Samples       string `json:"samples"`
	MinWeightMass string `json:"minWeightMass"`
	Precision     string `json:"precision"`
}

var _ Task = (*WeightedMeanTask)(nil)

func (t *WeightedMeanTask) Type() TaskType {
	return TaskTypeWeightedMean
}

func (t *WeightedMeanTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		samplesAndErrs SliceParam
		minMass        DecimalParam
		maybePrecision MaybeInt32Param
	)

	err := stderrors.Join(
		errors.Wrap(ResolveParam(&samplesAndErrs, From(VarExpr(t.Samples, vars), JSONWithVarExprs(t.Samples, vars, true), Inputs(inputs))), "samples"),
		errors.Wrap(ResolveParam(&minMass, From(VarExpr(t.MinWeightMass, vars), NonemptyString(t.MinWeightMass), "0")), "minWeightMass"),
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

	totalMass := decimal.Zero
	for _, c := range collapsed {
		totalMass = totalMass.Add(c.weight)
	}
	mm := minMass.Decimal()
	if mm.GreaterThan(decimal.Zero) && totalMass.LessThan(mm) {
		return Result{Error: errors.Wrapf(ErrInsufficientWeightMass, "mass %s < min %s", totalMass.String(), mm.String())}, runInfo
	}

	num := decimal.Zero
	den := decimal.Zero
	for _, c := range collapsed {
		num = num.Add(c.weight.Mul(c.value))
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
