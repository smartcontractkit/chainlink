package pipeline

import (
	"context"
	stderrors "errors"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// DurationMsParam parses a duration string (e.g. "10s") into milliseconds.
// Integer durations keep linear staleness exact and avoid floating-point drift.
type DurationMsParam int64

func (d *DurationMsParam) UnmarshalPipelineParam(val any) error {
	switch v := val.(type) {
	case string:
		if v == "" {
			return ErrParameterEmpty
		}
		dur, err := time.ParseDuration(v)
		if err != nil {
			return errors.Wrap(ErrBadInput, err.Error())
		}
		*d = DurationMsParam(dur.Milliseconds())
		return nil
	case time.Duration:
		*d = DurationMsParam(v.Milliseconds())
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
// Fails:  if method is unknown, or halfLife is missing for exp/cooldown
//
// Methods:
//   cutoff    — 1 if age <= threshold, else 0 (binary, discontinuous)
//   linear    — 1 - age/threshold, reaches 0 at threshold
//   exp       — 2^(-age/halfLife), truncated to 0 at threshold
//   cooldown  — 1 if age <= threshold, else 2^(-(age-threshold)/halfLife)
//   piecewise — linear interpolation of user-supplied (age:weight) points
type StalenessTask struct {
	BaseTask `mapstructure:",squash"`
	Samples  string `json:"samples"`
	Method   string `json:"method"`
	Threshold string `json:"threshold"`
	HalfLife string `json:"halfLife"`
	Points   string `json:"points"`
}

var _ Task = (*StalenessTask)(nil)

func (t *StalenessTask) Type() TaskType {
	return TaskTypeStaleness
}

func (t *StalenessTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		samplesAndErrs SliceParam
		method         StringParam
		thresholdMs    DurationMsParam
		halfLifeMs     DurationMsParam
		points         StringParam
	)

	err := stderrors.Join(
		errors.Wrap(ResolveParam(&samplesAndErrs, From(VarExpr(t.Samples, vars), JSONWithVarExprs(t.Samples, vars, true), Inputs(inputs))), "samples"),
		errors.Wrap(ResolveParam(&method, From(NonemptyString(t.Method))), "method"),
		errors.Wrap(ResolveParam(&thresholdMs, From(NonemptyString(t.Threshold))), "threshold"),
		errors.Wrap(ResolveParam(&halfLifeMs, From(NonemptyString(t.HalfLife), "0s")), "halfLife"),
		errors.Wrap(ResolveParam(&points, From(t.Points, "")), "points"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	samplesRaw, _ := samplesAndErrs.FilterErrors()
	var samples SampleSliceParam
	if err := samples.UnmarshalPipelineParam(samplesRaw); err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "samples: %v", err)}, runInfo
	}

	m := strings.ToLower(string(method))
	thresholdS := float64(thresholdMs) / 1000.0
	halfLifeS := float64(halfLifeMs) / 1000.0

	var pw []piecewisePoint
	if m == "piecewise" {
		var err error
		pw, err = parsePiecewisePoints(string(points))
		if err != nil {
			return Result{Error: errors.Wrap(err, "points")}, runInfo
		}
	}
	if (m == "exp" || m == "cooldown") && halfLifeS <= 0 {
		return Result{Error: errors.New("halfLife required for exp/cooldown staleness")}, runInfo
	}

	nowMs := time.Now().UnixMilli()
	out := make([]Sample, 0, len(samples))
	for _, s := range samples {
		ageMs := nowMs - s.TsMs
		if ageMs < 0 {
			ageMs = 0
		}
		ageS := float64(ageMs) / 1000.0

		mult, err := decayMultiplier(m, ageMs, ageS, int64(thresholdMs), thresholdS, halfLifeS, pw)
		if err != nil {
			return Result{Error: err}, runInfo
		}
		if mult.IsZero() {
			continue
		}
		s.Weight = s.Weight.Mul(mult)
		out = append(out, s)
	}

	return Result{Value: out}, runInfo
}

type piecewisePoint struct {
	ageS  float64
	value decimal.Decimal
}

func parsePiecewisePoints(s string) ([]piecewisePoint, error) {
	parts := strings.Split(s, ";")
	points := make([]piecewisePoint, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		av := strings.SplitN(p, ":", 2)
		if len(av) != 2 {
			return nil, errors.Errorf("invalid piecewise point %q", p)
		}
		dur, err := time.ParseDuration(strings.TrimSpace(av[0]))
		if err != nil {
			return nil, err
		}
		val, err := decimal.NewFromString(strings.TrimSpace(av[1]))
		if err != nil {
			return nil, err
		}
		points = append(points, piecewisePoint{ageS: dur.Seconds(), value: val})
	}
	if len(points) == 0 {
		return nil, errors.New("piecewise points cannot be empty")
	}
	sort.Slice(points, func(i, j int) bool { return points[i].ageS < points[j].ageS })
	return points, nil
}

func piecewiseValue(points []piecewisePoint, ageS float64) decimal.Decimal {
	if ageS <= points[0].ageS {
		return points[0].value
	}
	last := len(points) - 1
	if ageS >= points[last].ageS {
		return points[last].value
	}
	for i := 0; i < last; i++ {
		lo, hi := points[i], points[i+1]
		if ageS >= lo.ageS && ageS <= hi.ageS {
			if hi.ageS == lo.ageS {
				return lo.value
			}
			t := (ageS - lo.ageS) / (hi.ageS - lo.ageS)
			return lo.value.Add(hi.value.Sub(lo.value).Mul(decimal.NewFromFloat(t)))
		}
	}
	return points[last].value
}

func decayMultiplier(method string, ageMs int64, ageS float64, thresholdMs int64, thresholdS, halfLifeS float64, points []piecewisePoint) (decimal.Decimal, error) {
	switch method {
	case "cutoff":
		if ageMs <= thresholdMs {
			return decimal.NewFromInt(1), nil
		}
		return decimal.Zero, nil
	case "linear":
		if ageMs <= 0 {
			return decimal.NewFromInt(1), nil
		}
		if ageMs >= thresholdMs {
			return decimal.Zero, nil
		}
		return decimal.NewFromInt(thresholdMs - ageMs).Div(decimal.NewFromInt(thresholdMs)), nil
	case "exp":
		if ageS > thresholdS {
			return decimal.Zero, nil
		}
		return decimal.NewFromFloat(math.Pow(2, -ageS/halfLifeS)), nil
	case "cooldown":
		if ageS <= thresholdS {
			return decimal.NewFromInt(1), nil
		}
		return decimal.NewFromFloat(math.Pow(2, -(ageS-thresholdS)/halfLifeS)), nil
	case "piecewise":
		return piecewiseValue(points, ageS), nil
	default:
		return decimal.Zero, errors.Errorf("unknown staleness method %q", method)
	}
}
