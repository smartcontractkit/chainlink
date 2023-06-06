package denom

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CoinToAtom is a helper for converting to atom.
// TODO(BCI-915): Remove token specific functions
func CoinToAtom(coin sdk.Coin) (sdk.DecCoin, error) {
	return sdk.ConvertDecCoin(sdk.NewDecCoinFromCoin(coin), "atom")
}

// DecCoinToUAtom is a helper for converting to uatom.
// TODO(BCI-915): Remove token specific functions
func DecCoinToUAtom(coin sdk.DecCoin) (sdk.Coin, error) {
	decCoin, err := sdk.ConvertDecCoin(coin, "uatom")
	if err != nil {
		return sdk.Coin{}, err
	}
	truncated, _ := decCoin.TruncateDecimal()
	return truncated, nil
}
