package assets_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/stretchr/testify/assert"
)

func TestAssets_Units(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func(int64) *big.Int
		factor *big.Int
	}{
		{name: "Wei", fn: assets.Wei, factor: big.NewInt(params.Wei)},
		{name: "GWei", fn: assets.GWei, factor: big.NewInt(params.GWei)},
		{name: "UEther", fn: assets.UEther, factor: big.NewInt(params.GWei * 1000)},
		{name: "Ether", fn: assets.Ether, factor: big.NewInt(params.Ether)},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			expected := big.NewInt(0)
			assert.Equal(t, expected, test.fn(0))

			expected = big.NewInt(100)
			expected = new(big.Int).Mul(expected, test.factor)
			assert.Equal(t, expected, test.fn(100))

			expected = big.NewInt(-100)
			expected = new(big.Int).Mul(expected, test.factor)
			assert.Equal(t, expected, test.fn(-100))
		})
	}
}
