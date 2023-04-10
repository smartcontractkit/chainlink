package denom

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func init() {
	for _, d := range []struct {
		denom    string
		decimals int64
	}{
		{"atom", 0},
		{"matom", 3},
		{"uatom", 6},
	} {
		dec := sdk.NewDecWithPrec(1, d.decimals)
		if err := sdk.RegisterDenom(d.denom, dec); err != nil {
			panic(fmt.Errorf("failed to register denomination %q: %w", d.denom, err))
		}
	}
}

// CoinToAtom is a helper for converting to atom.
func CoinToAtom(coin sdk.Coin) (sdk.DecCoin, error) {
	return sdk.ConvertDecCoin(sdk.NewDecCoinFromCoin(coin), "atom")
}

// DecCoinToUAtom is a helper for converting to uatom.
func DecCoinToUAtom(coin sdk.DecCoin) (sdk.Coin, error) {
	decCoin, err := sdk.ConvertDecCoin(coin, "uatom")
	if err != nil {
		return sdk.Coin{}, err
	}
	truncated, _ := decCoin.TruncateDecimal()
	return truncated, nil
}
