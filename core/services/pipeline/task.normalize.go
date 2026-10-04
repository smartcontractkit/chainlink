package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// Experimental: NormalizeTask converts sample values onto a common unit using a
// user-provided unit→factor map. Samples already in the target unit (or with
// no unit) pass through unchanged.
//
// Input:  samples ([]Sample) — values in various units
// Input:  unitMap (map, required when enabled) — explicit {unit: factor} mapping
//
//	e.g. unitMap={"USDC":$(usdc_factor),"USDT":$(usdt_factor)} — var-refs are bare JSON values
//
// Output: []Sample — all values in targetUnit; unconvertible samples dropped (or error)
// Fails:  if a sample's unit has no factor and onMissingRate != "drop"
type NormalizeTask struct {
	BaseTask      `mapstructure:",squash"`
	Samples       string `json:"samples"`
	UnitMap       string `json:"unitMap"`
	TargetUnit    string `json:"targetUnit"`
	Enabled       string `json:"enabled"`
	OnMissingRate string `json:"onMissingRate"`
}

// UnitMapParam parses a JSON object mapping unit names to rate values.
// Values may be Samples (from upstream tasks) or raw numbers.
type UnitMapParam map[string]decimal.Decimal

func (u *UnitMapParam) UnmarshalPipelineParam(val any) error {
	switch v := val.(type) {
	case nil:
		return ErrParameterEmpty
	case map[string]any:
		out := make(map[string]decimal.Decimal, len(v))
		for unit, raw := range v {
			if s, ok := raw.(Sample); ok {
				out[unit] = s.Value
				continue
			}
			n, ok := raw.(json.Number)
			if ok {
				d, err := decimal.NewFromString(n.String())
				if err != nil {
					return errors.Wrapf(ErrBadInput, "unitMap[%q]: %v", unit, err)
				}
				out[unit] = d
				continue
			}
			d, err := utils.ToDecimal(raw)
			if err != nil {
				return errors.Wrapf(ErrBadInput, "unitMap[%q]: %v", unit, err)
			}
			out[unit] = d
		}
		*u = out
		return nil
	case []byte:
		jd := json.NewDecoder(bytes.NewReader(v))
		jd.UseNumber()
		var m map[string]any
		if err := jd.Decode(&m); err != nil {
			return errors.Wrap(ErrBadInput, err.Error())
		}
		return u.UnmarshalPipelineParam(m)
	case string:
		return u.UnmarshalPipelineParam([]byte(v))
	default:
		return errors.Wrapf(ErrBadInput, "expected unit map, got %T", val)
	}
}

var _ Task = (*NormalizeTask)(nil)

func (t *NormalizeTask) Type() TaskType {
	return TaskTypeNormalize
}

func (t *NormalizeTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		samplesAndErrs SliceParam
		unitMapParam   UnitMapParam
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
		resolveOpt(&unitMapParam, VarExpr(t.UnitMap, vars), JSONWithVarExprs(t.UnitMap, vars, true)),
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

	factorByUnit := map[string]decimal.Decimal(unitMapParam)

	out := make([]Sample, 0, len(samples))
	for _, s := range samples {
		if s.Unit == string(targetUnit) || s.Unit == "" {
			out = append(out, s)
			continue
		}
		factor, ok := factorByUnit[s.Unit]
		if !ok {
			if string(onMissingRate) == "drop" {
				continue
			}
			return Result{Error: errors.Errorf("no factor for unit %q", s.Unit)}, runInfo
		}
		s.Value = s.Value.Mul(factor)
		s.Unit = string(targetUnit)
		out = append(out, s)
	}

	return Result{Value: out}, runInfo
}
