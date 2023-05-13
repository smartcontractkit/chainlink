package denom

import (
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/params"
)

func TestMain(m *testing.M) {
	params.InitCosmosSdk(
		/* bech32Prefix= */ "wasm",
		/* token= */ "atom",
	)
	code := m.Run()
	os.Exit(code)
}

func TestCoinToAtom(t *testing.T) {
	tests := []struct {
		coin types.Coin
		exp  string
	}{
		{types.NewInt64Coin("uatom", 1), "0.000001000000000000atom"},
		{types.NewInt64Coin("uatom", 0), "0.000000000000000000atom"},
		{types.NewInt64Coin("atom", 1), "1.000000000000000000atom"},
		{types.NewInt64Coin("uatom", 1000000), "1.000000000000000000atom"},
		{types.NewInt64Coin("matom", 1000000), "1000.000000000000000000atom"},
		{types.NewInt64Coin("matom", 123456789), "123456.789000000000000000atom"},
	}
	for _, tt := range tests {
		t.Run(tt.coin.String(), func(t *testing.T) {
			got, err := CoinToAtom(tt.coin)
			require.NoError(t, err)
			require.Equal(t, tt.exp, got.String())
		})
	}
}
