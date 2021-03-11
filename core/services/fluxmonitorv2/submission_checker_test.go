package fluxmonitorv2_test

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/stretchr/testify/assert"
)

func TestSubmissionChecker_IsValid(t *testing.T) {
	testCases := []struct {
		name   string
		answer decimal.Decimal
		want   bool
	}{
		{
			name:   "equal to min",
			answer: decimal.NewFromFloat(0.01),
			want:   true,
		},
		{
			name:   "in between",
			answer: decimal.NewFromFloat(0.02),
			want:   true,
		},
		{
			name:   "equal to max",
			answer: decimal.NewFromFloat(0.03),
			want:   true,
		},
		{
			name:   "below min",
			answer: decimal.NewFromFloat(0.0),
			want:   false,
		},
		{
			name:   "over max",
			answer: decimal.NewFromFloat(0.04),
			want:   false,
		},
	}

	checker := fluxmonitorv2.NewSubmissionChecker(
		big.NewInt(1),
		big.NewInt(3),
		2,
	)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, checker.IsValid(tc.answer))
		})
	}
}
