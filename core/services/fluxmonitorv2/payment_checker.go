package fluxmonitorv2

import (
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
)

// MinFundedRounds defines the minimum number of rounds that needs to be paid
// to oracles on a contract
const MinFundedRounds int64 = 3

// PaymentChecker provides helper functions to check whether payments are valid
type PaymentChecker struct {
	// The minimum amount for a payment set in the ENV Var
	MinContractPayment *assets.Link
	// The minimum amount for a payment set in the job
	MinJobPayment *assets.Link
}

// NewPaymentChecker constructs a new payment checker
func NewPaymentChecker(minContractPayment, minJobPayment *assets.Link) *PaymentChecker {
	return &PaymentChecker{
		MinContractPayment: minContractPayment,
		MinJobPayment:      minJobPayment,
	}
}

// SufficientFunds checks if the contract has sufficient funding to pay all the
// oracles on a contract for a minimum number of rounds, based on the payment
// amount in the contract
func (c *PaymentChecker) SufficientFunds(availableFunds *big.Int, paymentAmount *big.Int, oracleCount uint8) bool {
	min := big.NewInt(int64(oracleCount))
	min = min.Mul(min, big.NewInt(MinFundedRounds))
	min = min.Mul(min, paymentAmount)

	return availableFunds.Cmp(min) >= 0
}

// SufficientPayment checks if the available payment is enough to submit an
// answer. It compares the payment amount on chain with the min payment amount
// listed in the job / ENV var.
func (c *PaymentChecker) SufficientPayment(payment *big.Int) bool {
	aboveOrEqMinGlobalPayment := payment.Cmp(c.MinContractPayment.ToInt()) >= 0
	aboveOrEqMinJobPayment := true
	if c.MinJobPayment != nil {
		aboveOrEqMinJobPayment = payment.Cmp(c.MinJobPayment.ToInt()) >= 0
	}
	return aboveOrEqMinGlobalPayment && aboveOrEqMinJobPayment
}
