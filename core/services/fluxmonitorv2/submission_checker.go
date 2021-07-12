package fluxmonitorv2

import (
	"math/big"

	"github.com/shopspring/decimal"
)

// SubmissionChecker checks whether an answer is inside the allowable range.
type SubmissionChecker struct {
	Min decimal.Decimal
	Max decimal.Decimal
}

// NewSubmissionChecker initializes a new SubmissionChecker
func NewSubmissionChecker(min *big.Int, max *big.Int) *SubmissionChecker {
	return &SubmissionChecker{
		Min: decimal.NewFromBigInt(min, 0),
		Max: decimal.NewFromBigInt(max, 0),
	}
}

// IsValid checks if the submission is between the min and max
func (c *SubmissionChecker) IsValid(answer decimal.Decimal) bool {
	return answer.GreaterThanOrEqual(c.Min) && answer.LessThanOrEqual(c.Max)
}
