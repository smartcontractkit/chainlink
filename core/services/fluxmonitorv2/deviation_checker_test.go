package fluxmonitorv2_test

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
)

type outsideDeviationRow struct {
	name                string
	curPrice, nextPrice decimal.Decimal
	threshold           float64 // in percentage
	absoluteThreshold   float64
	expectation         bool
}

func (o outsideDeviationRow) String() string {
	return fmt.Sprintf(
		`{name: "%s", curPrice: %s, nextPrice: %s, threshold: %.2f, `+
			"absoluteThreshold: %f, expectation: %v}", o.name, o.curPrice, o.nextPrice,
		o.threshold, o.absoluteThreshold, o.expectation)
}

func TestDeviationChecker_OutsideDeviation(t *testing.T) {
	t.Parallel()

	f, i := decimal.NewFromFloat, decimal.NewFromInt
	testCases := []outsideDeviationRow{
		// Start with a huge absoluteThreshold, to test relative threshold behavior
		{"0 current price, outside deviation", i(0), i(100), 2, 0, true},
		{"0 current and next price", i(0), i(0), 2, 0, false},

		{"inside deviation", i(100), i(101), 2, 0, false},
		{"equal to deviation", i(100), i(102), 2, 0, true},
		{"outside deviation", i(100), i(103), 2, 0, true},
		{"outside deviation zero", i(100), i(0), 2, 0, true},

		{"inside deviation, crosses 0 backwards", f(0.1), f(-0.1), 201, 0, false},
		{"equal to deviation, crosses 0 backwards", f(0.1), f(-0.1), 200, 0, true},
		{"outside deviation, crosses 0 backwards", f(0.1), f(-0.1), 199, 0, true},

		{"inside deviation, crosses 0 forwards", f(-0.1), f(0.1), 201, 0, false},
		{"equal to deviation, crosses 0 forwards", f(-0.1), f(0.1), 200, 0, true},
		{"outside deviation, crosses 0 forwards", f(-0.1), f(0.1), 199, 0, true},

		{"thresholds=0, deviation", i(0), i(100), 0, 0, true},
		{"thresholds=0, no deviation", i(100), i(100), 0, 0, true},
		{"thresholds=0, all zeros", i(0), i(0), 0, 0, true},
	}

	c := func(tc outsideDeviationRow) {
		checker := fluxmonitorv2.NewDeviationChecker(tc.threshold, tc.absoluteThreshold, logger.TestLogger(t))

		assert.Equal(t, tc.expectation,
			checker.OutsideDeviation(tc.curPrice, tc.nextPrice),
			"check on OutsideDeviation failed for %s", tc,
		)
	}

	for _, tc := range testCases {
		tc := tc
		// Checks on relative threshold
		t.Run(tc.name, func(t *testing.T) { c(tc) })
		// Check corresponding absolute threshold tests; make relative threshold
		// always pass (as long as curPrice and nextPrice aren't both 0.)
		test2 := tc
		test2.threshold = 0
		// absoluteThreshold is initially zero, so any change will trigger
		test2.expectation = test2.curPrice.Sub(tc.nextPrice).Abs().GreaterThan(i(0)) ||
			test2.absoluteThreshold == 0
		t.Run(tc.name+" threshold zeroed", func(t *testing.T) { c(test2) })
		// Huge absoluteThreshold means trigger always fails
		test3 := tc
		test3.absoluteThreshold = 1e307
		test3.expectation = false
		t.Run(tc.name+" max absolute threshold", func(t *testing.T) { c(test3) })
	}
}
