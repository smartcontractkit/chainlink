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
// Fails:  if method is unknown, or halfLife is missing for exp/exp_cooldown
//
// Methods:
//
//	cutoff    — 1 if age <= threshold, else 0 (binary, discontinuous)
//	linear    — 1 - age/threshold, reaches 0 at threshold
//	exp       — 2^(-age/halfLife), truncated to 0 at threshold
//	exp_cooldown  — 1 if age <= threshold, else 2^(-(age-threshold)/halfLife)
//	piecewise — linear interpolation of user-supplied (age:weight) points
//
// Optional safety parameters:
//
//	decayThreshold (K) — when the decay multiplier falls below K, the sample is
//	dropped entirely (assigned "no data" state). Recommended value: 0.03.
//	Applies to exp and exp_cooldown methods (where decay is asymptotic).
//
//	cutoff — hard time limit; when age exceeds cutoff, the sample is dropped
//	regardless of the decay function value. Acts as a safety valve.
type StalenessTask struct {
	BaseTask       `mapstructure:",squash"`
	Samples        string `json:"samples"`
	Method         string `json:"method"`
	Threshold      string `json:"threshold"`
	HalfLife       string `json:"halfLife"`
	Points         string `json:"points"`
	Cutoff         string `json:"cutoff"`
	DecayThreshold string `json:"decayThreshold"`
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
		cutoffMs       DurationMsParam
		decayThreshold DecimalParam
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
		errors.Wrap(ResolveParam(&thresholdMs, From(NonemptyString(t.Threshold))), "threshold"),
		errors.Wrap(ResolveParam(&halfLifeMs, From(NonemptyString(t.HalfLife), "0s")), "halfLife"),
		errors.Wrap(ResolveParam(&points, From(t.Points, "")), "points"),
		resolveOpt(&cutoffMs, VarExpr(t.Cutoff, vars), NonemptyString(t.Cutoff)),
		resolveOpt(&decayThreshold, VarExpr(t.DecayThreshold, vars), NonemptyString(t.DecayThreshold)),
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
	cutoffMsVal := int64(cutoffMs)
	kVal := decayThreshold.Decimal()
	applyK := (m == "exp" || m == "exp_cooldown") && kVal.GreaterThan(decimal.Zero)

	var pw []piecewisePoint
	if m == "piecewise" {
		var err error
		pw, err = parsePiecewisePoints(string(points))
		if err != nil {
			return Result{Error: errors.Wrap(err, "points")}, runInfo
		}
	}
	if (m == "exp" || m == "exp_cooldown") && halfLifeMs <= 0 {
		return Result{Error: errors.New("halfLife required for exp/exp_cooldown staleness")}, runInfo
	}

	nowMs := time.Now().UnixMilli()
	out := make([]Sample, 0, len(samples))
	for _, s := range samples {
		ageMs := max(nowMs-s.TsMs, 0)

		// Hard cutoff: drop samples older than the cutoff time limit.
		if cutoffMsVal > 0 && ageMs > cutoffMsVal {
			continue
		}

		mult, err := decayMultiplier(m, ageMs, int64(thresholdMs), int64(halfLifeMs), pw)
		if err != nil {
			return Result{Error: err}, runInfo
		}

		// Decay threshold K: for asymptotic methods (exp, exp_cooldown), drop
		// samples whose decay multiplier falls below K ("no data" state).
		if applyK && mult.LessThan(kVal) {
			continue
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
	ageMs int64
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
		points = append(points, piecewisePoint{ageMs: dur.Milliseconds(), value: val})
	}
	if len(points) == 0 {
		return nil, errors.New("piecewise points cannot be empty")
	}
	sort.Slice(points, func(i, j int) bool { return points[i].ageMs < points[j].ageMs })
	return points, nil
}

func piecewiseValue(points []piecewisePoint, ageMs int64) decimal.Decimal {
	if ageMs <= points[0].ageMs {
		return points[0].value
	}
	last := len(points) - 1
	if ageMs >= points[last].ageMs {
		return points[last].value
	}
	for i := range last {
		lo, hi := points[i], points[i+1]
		if ageMs >= lo.ageMs && ageMs <= hi.ageMs {
			if hi.ageMs == lo.ageMs {
				return lo.value
			}
			t := float64(ageMs-lo.ageMs) / float64(hi.ageMs-lo.ageMs)
			return lo.value.Add(hi.value.Sub(lo.value).Mul(decimal.NewFromFloat(t)))
		}
	}
	return points[last].value
}

func decayMultiplier(method string, ageMs, thresholdMs, halfLifeMs int64, points []piecewisePoint) (decimal.Decimal, error) {
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
		if ageMs > thresholdMs {
			return decimal.Zero, nil
		}
		return decimal.NewFromFloat(math.Pow(2, -float64(ageMs)/float64(halfLifeMs))), nil
	case "exp_cooldown":
		if ageMs <= thresholdMs {
			return decimal.NewFromInt(1), nil
		}
		return decimal.NewFromFloat(math.Pow(2, -float64(ageMs-thresholdMs)/float64(halfLifeMs))), nil
	case "piecewise":
		return piecewiseValue(points, ageMs), nil
	default:
		return decimal.Zero, errors.Errorf("unknown staleness method %q", method)
	}
}
