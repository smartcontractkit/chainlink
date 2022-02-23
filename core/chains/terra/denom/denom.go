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
		{"luna", 0},
		{"mluna", 3},
		{"uluna", 6},
	} {
		dec := sdk.NewDecWithPrec(1, d.decimals)
		if err := sdk.RegisterDenom(d.denom, dec); err != nil {
			panic(fmt.Errorf("failed to register denomination %q: %w", d.denom, err))
		}
	}
}

// ConvertToLuna is just a helper for converting to luna from other denominations.
func ConvertToLuna(coin sdk.Coin) (sdk.DecCoin, error) {
	return sdk.ConvertDecCoin(sdk.NewDecCoinFromCoin(coin), "luna")
}
