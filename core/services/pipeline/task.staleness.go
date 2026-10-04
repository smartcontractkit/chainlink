package pipeline

import (
	"context"
	stderrors "errors"
	"math"
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
type StalenessTask struct {
	BaseTask       `mapstructure:",squash"`
	Samples        string `json:"samples"`
	Method         string `json:"method"`
	Threshold      string `json:"threshold"`
	HalfLife       string `json:"halfLife"`
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
		threshold      DurationParam
		halfLife       DurationParam
		cutoff         DurationParam
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
		resolveOpt(&threshold, NonemptyString(t.Threshold)),
		errors.Wrap(ResolveParam(&halfLife, From(NonemptyString(t.HalfLife), "0s")), "halfLife"),
		resolveOpt(&cutoff, VarExpr(t.Cutoff, vars), NonemptyString(t.Cutoff)),
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
	cutoffVal := time.Duration(cutoff)
	kVal := decayThreshold.Decimal()
	applyK := (m == "exp" || m == "exp_cooldown") && kVal.GreaterThan(decimal.Zero)

	if time.Duration(threshold) <= 0 {
		return Result{Error: errors.New("threshold required for staleness")}, runInfo
	}
	if (m == "exp" || m == "exp_cooldown") && halfLife <= 0 {
		return Result{Error: errors.New("halfLife required for exp/exp_cooldown staleness")}, runInfo
	}

	nowNs := time.Now().UnixNano()
	out := make([]Sample, 0, len(samples))
	for _, s := range samples {
		ageNs := max(nowNs-s.TsMs*int64(time.Millisecond), 0)

		// Hard cutoff: drop samples older than the cutoff time limit.
		if cutoffVal > 0 && ageNs > int64(cutoffVal) {
			continue
		}

		mult, err := decayMultiplier(m, time.Duration(ageNs), time.Duration(threshold), time.Duration(halfLife))
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
