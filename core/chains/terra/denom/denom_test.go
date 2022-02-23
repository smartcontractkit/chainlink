package denom

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestConvertToLuna(t *testing.T) {
	tests := []struct {
		coin types.Coin
		exp  string
	}{
		{types.NewInt64Coin("uluna", 1), "0.000001000000000000luna"},
		{types.NewInt64Coin("uluna", 0), "0.000000000000000000luna"},
		{types.NewInt64Coin("luna", 1), "1.000000000000000000luna"},
		{types.NewInt64Coin("uluna", 1000000), "1.000000000000000000luna"},
		{types.NewInt64Coin("mluna", 1000000), "1000.000000000000000000luna"},
		{types.NewInt64Coin("mluna", 123456789), "123456.789000000000000000luna"},
	}
	for _, tt := range tests {
		t.Run(tt.coin.String(), func(t *testing.T) {
			got, err := ConvertToLuna(tt.coin)
			require.NoError(t, err)
			require.Equal(t, tt.exp, got.String())
		})
	}
}
