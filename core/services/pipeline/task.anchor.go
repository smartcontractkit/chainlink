package pipeline

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Experimental: AnchorTask rebases a lower or upper sample stream onto a
// reference stream using a spread-derived per-source displacement, then
// returns the weighted mean.
//
// For each surviving source i:
//
//	disp_i = (high_i - low_i) / (2 * ref_i)
//	q_i = clamp(ref_i, m*(1-band), m*(1+band))
//
// When select="low":  adj_i = q_i * (1 - disp_i)
// When select="high": adj_i = q_i * (1 + disp_i)
//
// Weights come from the reference stream. Sources missing from low or high
// are dropped. A reference value of zero drops that source.
//
// Input:  reference ([]Sample) — the anchor stream; drives weights + band
// Input:  low ([]Sample) — the lower stream
// Input:  high ([]Sample) — the upper stream
// Output: decimal.Decimal — the weighted mean of rebased values
// Fails:  ErrInsufficientWeightMass if surviving mass < minWeightMass
type AnchorTask struct {
	BaseTask      `mapstructure:",squash"`
	Reference     string `json:"reference"`
	Low           string `json:"low"`
	High          string `json:"high"`
	Select        string `json:"select"`
	Band          string `json:"band"`
	MinWeightMass string `json:"minWeightMass"`
	Precision     string `json:"precision"`
}

var _ Task = (*AnchorTask)(nil)

func (t *AnchorTask) Type() TaskType {
	return TaskTypeAnchor
}

func (t *AnchorTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		referenceAndErrs, lowAndErrs, highAndErrs SliceParam
		band, minMass                            DecimalParam
		selectParam                              StringParam
		maybePrecision                           MaybeInt32Param
	)

	if err := stderrors.Join(
		errors.Wrap(ResolveParam(&referenceAndErrs, From(VarExpr(t.Reference, vars), JSONWithVarExprs(t.Reference, vars, true))), "reference"),
		errors.Wrap(ResolveParam(&lowAndErrs, From(VarExpr(t.Low, vars), JSONWithVarExprs(t.Low, vars, true), Inputs(inputs))), "low"),
		errors.Wrap(ResolveParam(&highAndErrs, From(VarExpr(t.High, vars), JSONWithVarExprs(t.High, vars, true))), "high"),
		errors.Wrap(ResolveParam(&selectParam, From(NonemptyString(t.Select), "low")), "select"),
		errors.Wrap(ResolveParam(&band, From(VarExpr(t.Band, vars), NonemptyString(t.Band), "0.03")), "band"),
		errors.Wrap(ResolveParam(&minMass, From(VarExpr(t.MinWeightMass, vars), NonemptyString(t.MinWeightMass), "0")), "minWeightMass"),
		errors.Wrap(ResolveParam(&maybePrecision, From(VarExpr(t.Precision, vars), t.Precision)), "precision"),
	); err != nil {
		return Result{Error: err}, runInfo
	}

	toMap := func(andErrs SliceParam, label string) (map[string]Sample, error) {
		raw, _ := andErrs.FilterErrors()
		var s SampleSliceParam
		if err := s.UnmarshalPipelineParam(raw); err != nil {
			return nil, errors.Wrapf(ErrBadInput, "%s: %v", label, err)
		}
		m := make(map[string]Sample, len(s))
		for _, v := range s {
			if v.Weight.IsNegative() {
				return nil, errors.Errorf("negative weight %s on %s sample %q", v.Weight, label, v.Source)
			}
			m[v.Source] = v
		}
		return m, nil
	}

	lowBySource, err := toMap(lowAndErrs, "low")
	if err != nil {
		return Result{Error: err}, runInfo
	}
	highBySource, err := toMap(highAndErrs, "high")
	if err != nil {
		return Result{Error: err}, runInfo
	}

	refRaw, _ := referenceAndErrs.FilterErrors()
	var refSamples SampleSliceParam
	if err := refSamples.UnmarshalPipelineParam(refRaw); err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "reference: %v", err)}, runInfo
	}
	ref := make([]Sample, 0, len(refSamples))
	for _, s := range refSamples {
		if s.Weight.IsNegative() {
			return Result{Error: errors.Errorf("negative weight %s on reference sample %q", s.Weight, s.Source)}, runInfo
		}
		if !s.Weight.IsZero() {
			ref = append(ref, s)
		}
	}
	if len(ref) == 0 {
		return Result{Error: errors.Wrap(ErrWrongInputCardinality, "no surviving reference samples")}, runInfo
	}

	b := band.Decimal()
	if b.IsNegative() {
		return Result{Error: errors.Errorf("band must be non-negative, got %s", b)}, runInfo
	}
	loBand, hiBand, _ := computeBand(bandReference(ref, "median"), b, "winsor", ref)

	sign := decimal.NewFromInt(-1)
	if strings.ToLower(string(selectParam)) == "high" {
		sign = decimal.NewFromInt(1)
	}

	two := decimal.NewFromInt(2)
	num, den := decimal.Zero, decimal.Zero
	for _, r := range ref {
		loS, ok := lowBySource[r.Source]
		if !ok || r.Value.IsZero() {
			continue
		}
		hiS, ok := highBySource[r.Source]
		if !ok {
			continue
		}
		q := r.Value
		if q.LessThan(loBand) {
			q = loBand
		}
		if q.GreaterThan(hiBand) {
			q = hiBand
		}
		disp := hiS.Value.Sub(loS.Value).Div(two.Mul(r.Value))
		adj := q.Mul(decimal.NewFromInt(1).Add(sign.Mul(disp)))
		num = num.Add(r.Weight.Mul(adj))
		den = den.Add(r.Weight)
	}

	if den.IsZero() {
		return Result{Error: errors.Wrap(ErrWrongInputCardinality, "no surviving samples after anchor join")}, runInfo
	}
	if mm := minMass.Decimal(); mm.GreaterThan(decimal.Zero) && den.LessThan(mm) {
		return Result{Error: errors.Wrapf(ErrInsufficientWeightMass, "mass %s < min %s", den, mm)}, runInfo
	}
	if prec, isSet := maybePrecision.Int32(); isSet {
		return Result{Value: num.DivRound(den, prec)}, runInfo
	}
	return Result{Value: num.Div(den)}, runInfo
}
