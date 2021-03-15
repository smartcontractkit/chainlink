package fluxmonitorv2

import (
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/logger"
)

// DeviationThresholds carries parameters used by the threshold-trigger logic
type DeviationThresholds struct {
	Rel float64 // Relative change required, i.e. |new-old|/|old| >= Rel
	Abs float64 // Absolute change required, i.e. |new-old| >= Abs
}

// DeviationChecker checks the deviation of the next answer against the current
// answer.
type DeviationChecker struct {
	Thresholds DeviationThresholds
}

// NewDeviationChecker constructs a new deviation checker with thresholds.
func NewDeviationChecker(rel, abs float64) *DeviationChecker {
	return &DeviationChecker{
		Thresholds: DeviationThresholds{
			Rel: rel,
			Abs: abs,
		},
	}
}

// NewZeroDeviationChecker constructs a new deviation checker with 0 as thresholds.
func NewZeroDeviationChecker() *DeviationChecker {
	return &DeviationChecker{
		Thresholds: DeviationThresholds{
			Rel: 0,
			Abs: 0,
		},
	}
}

// OutsideDeviation checks whether the next price is outside the threshold.
// If both thresholds are zero (default value), always returns true.
func (c *DeviationChecker) OutsideDeviation(curAnswer, nextAnswer decimal.Decimal) bool {
	loggerFields := []interface{}{
		"threshold", c.Thresholds.Rel,
		"absoluteThreshold", c.Thresholds.Abs,
		"currentAnswer", curAnswer,
		"nextAnswer", nextAnswer,
	}

	if c.Thresholds.Rel == 0 && c.Thresholds.Abs == 0 {
		logger.Debugw(
			"Deviation thresholds both zero; short-circuiting deviation checker to "+
				"true, regardless of feed values", loggerFields...)
		return true
	}
	diff := curAnswer.Sub(nextAnswer).Abs()
	loggerFields = append(loggerFields, "absoluteDeviation", diff)

	if !diff.GreaterThan(decimal.NewFromFloat(c.Thresholds.Abs)) {
		logger.Debugw("Absolute deviation threshold not met", loggerFields...)
		return false
	}

	if curAnswer.IsZero() {
		if nextAnswer.IsZero() {
			logger.Debugw("Relative deviation is undefined; can't satisfy threshold", loggerFields...)
			return false
		}
		logger.Infow("Threshold met: relative deviation is ∞", loggerFields...)
		return true
	}

	// 100*|new-old|/|old|: Deviation (relative to curAnswer) as a percentage
	percentage := diff.Div(curAnswer.Abs()).Mul(decimal.NewFromInt(100))

	loggerFields = append(loggerFields, "percentage", percentage)

	if percentage.LessThan(decimal.NewFromFloat(c.Thresholds.Rel)) {
		logger.Debugw("Relative deviation threshold not met", loggerFields...)
		return false
	}
	logger.Infow("Relative and absolute deviation thresholds both met", loggerFields...)
	return true
}
