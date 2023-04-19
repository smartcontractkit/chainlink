package assets_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
)

func TestAssets_Units(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func(int64) *assets.Wei
		factor *big.Int
	}{
		{name: "Wei", fn: assets.NewWeiI[int64], factor: big.NewInt(params.Wei)},
		{name: "GWei", fn: assets.GWei[int64], factor: big.NewInt(params.GWei)},
		{name: "UEther", fn: assets.UEther[int64], factor: big.NewInt(params.GWei * 1000)},
		{name: "Ether", fn: assets.Ether[int64], factor: big.NewInt(params.Ether)},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			expected := assets.NewWeiI(0)
			assert.Equal(t, expected, test.fn(0))

			expected = assets.NewWeiI(100)
			expected = expected.Mul(test.factor)
			assert.Equal(t, expected, test.fn(100))
		})
	}
}
