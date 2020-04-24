package fluxmonitor

// This implements trigger functions for reporting of fluxmonitor observations.

import (
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/store/models/triggerfns"

	"github.com/shopspring/decimal"
)

// Register factories with the factory serializer/deserializer, so that they
// will be recognized in job specs
func init() {
	triggerfns.RegisterTriggerFunctionFactory("relativeThreshold", floatThresholdFactory)
	triggerfns.RegisterTriggerFunctionFactory("absoluteThreshold", floatThresholdFactory)
}

type floatTriggerFn struct {
	triggering func(onchain, recent decimal.Decimal, extraData ...interface{}) (bool, error)
	factory    string
	parameters float64
}

var _ triggerfns.TriggerFn = floatTriggerFn{}

func (t floatTriggerFn) Triggering(onchain, recent decimal.Decimal,
	extraData ...interface{}) (bool, error) {
	return t.triggering(onchain, recent, extraData...)
}

func (t floatTriggerFn) Parameters() interface{} { return t.parameters }
func (t floatTriggerFn) Factory() string         { return t.factory }

type thresholdTrigger func(onchain, recent decimal.Decimal,
	extraData ...interface{}) (bool, error)

func absoluteTrigger(dthreshold decimal.Decimal) thresholdTrigger {
	return func(onchain, recent decimal.Decimal, extraData ...interface{}) (bool, error) {
		if onchain.Sign() != 0 { // current != 0, so |current-new|/|current| is well-defined and finite
			// Trigger if |current-new|/|current| >= threshold
			triggered := !onchain.Sub(recent).Div(onchain).Abs().
				// Treat threshold value as a percentage (multiplication by 100 here
				// is the same as division by 100 on the other side of the
				// inequality.)
				Mul(decimal.NewFromFloat(100)).
				LessThan(dthreshold)
			return triggered, nil
		} else { // current == 0
			// If new != 0, |current-new|/|current| = ∞ > threshold, so trigger
			// If new == 0, new == current, so do not trigger (no deviation)
			return recent.Sign() != 0, nil
		}
	}
}

func relativeTrigger(dthreshold decimal.Decimal) thresholdTrigger {
	return func(onchain, recent decimal.Decimal, extraData ...interface{}) (bool, error) {
		// Trigger if |current-new| >= threshold
		return !onchain.Sub(recent).Abs().LessThan(dthreshold), nil
	}
}

var thresholdFactories = map[string]func(decimal.Decimal) thresholdTrigger{
	"absoluteThreshold": absoluteTrigger,
	"relativeThreshold": relativeTrigger,
}

// Memoize the trigger functions based on their factory/parameters. This makes
// it possible to meaningfully compare functions between deserializations of
// ValueTriggers, which is useful for testing.
var memo map[string]map[float64]triggerfns.TriggerFn

func floatThresholdFactory(name string, params interface{}) (triggerfns.TriggerFn, error) {
	threshold, ok := params.(float64)
	if !ok {
		return nil, errors.Errorf(
			"%s parameter should be number, got %+v", name, params)
	}
	if threshold <= 0 {
		return nil, errors.Errorf("threshold must be positive for %s trigger, got %f",
			name, threshold)
	}
	if memo == nil {
		memo = make(map[string]map[float64]triggerfns.TriggerFn)
	}
	fns, ok := memo[name]
	if !ok {
		memo[name] = make(map[float64]triggerfns.TriggerFn)
		fns = memo[name]
	}
	fn, ok := fns[threshold]
	if ok {
		return fn, nil
	}
	tf := thresholdFactories[name](decimal.NewFromFloat(threshold))
	rv := floatTriggerFn{factory: name, parameters: threshold, triggering: tf}
	fns[threshold] = rv
	return rv, nil
}
