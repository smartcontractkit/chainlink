package rhea

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenPrice(t *testing.T) {
	var tt = []struct {
		token    Token
		expected *big.Int
	}{
		{
			LINK,
			new(big.Int).Mul(big.NewInt(65), big.NewInt(1e17)),
		},
		{
			SUPER,
			new(big.Int).Mul(big.NewInt(1), big.NewInt(1e18)),
		},
		{
			CACHEGOLD,
			new(big.Int).Mul(big.NewInt(60), big.NewInt(1e18)),
		},
	}
	for _, tc := range tt {
		tc := tc
		a := tc.token.Price()
		assert.Equal(t, tc.expected, a)
	}
}

func TestGetTokenPricePer1e18Units(t *testing.T) {
	var tt = []struct {
		price    *big.Int
		decimals uint8
		expected *big.Int
	}{
		{
			new(big.Int).Mul(big.NewInt(10), big.NewInt(1e18)),
			18,
			new(big.Int).Mul(big.NewInt(10), big.NewInt(1e18)),
		},
		{
			new(big.Int).Mul(big.NewInt(1), big.NewInt(1e18)),
			6,
			new(big.Int).Mul(big.NewInt(1e12), big.NewInt(1e18)),
		},
		{
			new(big.Int).Mul(big.NewInt(60), big.NewInt(1e18)),
			8,
			new(big.Int).Mul(big.NewInt(60e10), big.NewInt(1e18)),
		},
	}
	for _, tc := range tt {
		tc := tc
		a := GetPricePer1e18Units(tc.price, tc.decimals)
		assert.Equal(t, tc.expected, a)
	}
}
