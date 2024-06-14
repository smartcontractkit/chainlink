package fluxmonitorv2_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2"
)

func TestPaymentChecker_SufficientFunds(t *testing.T) {
	var (
		checker     = fluxmonitorv2.NewPaymentChecker(nil, nil)
		payment     = 100
		rounds      = 3
		oracleCount = 21
		min         = payment * rounds * oracleCount
	)

	testCases := []struct {
		name  string
		funds int
		want  bool
	}{
		{"above minimum", min + 1, true},
		{"equal to minimum", min, true},
		{"below minimum", min - 1, false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, checker.SufficientFunds(
				big.NewInt(int64(tc.funds)),
				big.NewInt(int64(payment)),
				uint8(oracleCount),
			))
		})
	}
}

func TestPaymentChecker_SufficientPayment(t *testing.T) {
	var (
		payment int64 = 10
		eq            = payment
		gt            = payment + 1
		lt            = payment - 1
	)

	testCases := []struct {
		name               string
		minContractPayment int64
		minJobPayment      interface{} // nil or int64
		want               bool
	}{
		{"payment above min contract payment, no min job payment", lt, nil, true},
		{"payment equal to min contract payment, no min job payment", eq, nil, true},
		{"payment below min contract payment, no min job payment", gt, nil, false},

		{"payment above min contract payment, above min job payment", lt, lt, true},
		{"payment equal to min contract payment, above min job payment", eq, lt, true},
		{"payment below min contract payment, above min job payment", gt, lt, false},

		{"payment above min contract payment, equal to min job payment", lt, eq, true},
		{"payment equal to min contract payment, equal to min job payment", eq, eq, true},
		{"payment below min contract payment, equal to min job payment", gt, eq, false},

		{"payment above minimum contract payment, below min job payment", lt, gt, false},
		{"payment equal to minimum contract payment, below min job payment", eq, gt, false},
		{"payment below minimum contract payment, below min job payment", gt, gt, false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var minJobPayment *assets.Link
			if tc.minJobPayment != nil {
				mjb := assets.Link(*big.NewInt(tc.minJobPayment.(int64)))
				minJobPayment = &mjb
			}

			checker := fluxmonitorv2.NewPaymentChecker(assets.NewLinkFromJuels(tc.minContractPayment), minJobPayment)

			assert.Equal(t, tc.want, checker.SufficientPayment(big.NewInt(payment)))
		})
	}
}
