package pipeline

import (
	"context"
	stderrors "errors"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Experimental: NormalizeTask converts sample values onto a common unit using FX rates
// provided as their own Sample slice. Samples already in the target unit (or
// with no unit) pass through unchanged.
//
// Input:  samples ([]Sample) — values in various units
// Input:  rates ([]Sample, optional) — Source names the unit this rate converts
// Output: []Sample — all values in targetUnit; unconvertible samples dropped (or error)
// Fails:  if a sample's unit has no rate and onMissingRate != "drop"
type NormalizeTask struct {
	BaseTask      `mapstructure:",squash"`
	Samples       string `json:"samples"`
	Rates         string `json:"rates"`
	TargetUnit    string `json:"targetUnit"`
	Enabled       string `json:"enabled"`
	OnMissingRate string `json:"onMissingRate"`
}

var _ Task = (*NormalizeTask)(nil)

func (t *NormalizeTask) Type() TaskType {
	return TaskTypeNormalize
}

func (t *NormalizeTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		samplesAndErrs SliceParam
		ratesAndErrs   SliceParam
		targetUnit     StringParam
		enabled        BoolParam
		onMissingRate  StringParam
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
		resolveOpt(&ratesAndErrs, VarExpr(t.Rates, vars), JSONWithVarExprs(t.Rates, vars, true)),
		errors.Wrap(ResolveParam(&targetUnit, From(NonemptyString(t.TargetUnit))), "targetUnit"),
		errors.Wrap(ResolveParam(&enabled, From(NonemptyString(t.Enabled), true)), "enabled"),
		errors.Wrap(ResolveParam(&onMissingRate, From(NonemptyString(t.OnMissingRate), "drop")), "onMissingRate"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	samplesRaw, _ := samplesAndErrs.FilterErrors()
	var samples SampleSliceParam
	if err := samples.UnmarshalPipelineParam(samplesRaw); err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "samples: %v", err)}, runInfo
	}

	if !bool(enabled) {
		return Result{Value: []Sample(samples)}, runInfo
	}

	ratesRaw, _ := ratesAndErrs.FilterErrors()
	var rates SampleSliceParam
	if err := rates.UnmarshalPipelineParam(ratesRaw); err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "rates: %v", err)}, runInfo
	}

	// Build a lookup from unit name to rate value. The proposal says the Source
	// field on a rate sample names the unit it converts; fall back to Unit if
	// Source is empty.
	rateByUnit := make(map[string]decimal.Decimal, len(rates))
	for _, r := range rates {
		if r.Value.IsZero() {
			continue
		}
		unit := r.Source
		if unit == "" {
			unit = r.Unit
		}
		if unit == "" {
			continue
		}
		rateByUnit[unit] = r.Value
	}

	out := make([]Sample, 0, len(samples))
	for _, s := range samples {
		if s.Unit == string(targetUnit) || s.Unit == "" {
			out = append(out, s)
			continue
		}
		rate, ok := rateByUnit[s.Unit]
		if !ok {
			if string(onMissingRate) == "drop" {
				continue
			}
			return Result{Error: errors.Errorf("no rate for unit %q", s.Unit)}, runInfo
		}
		s.Value = s.Value.Mul(rate)
		s.Unit = string(targetUnit)
		out = append(out, s)
	}

	return Result{Value: out}, runInfo
}
