package denom

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ConvertDecCoinToDenom is a helper for converting a DecCoin to a given denomination, rounded
// down with the remainder discarded. Requires InitCosmosSdk to be called first to register
// both the source and destinations token denominations, otherwise will return an error.
func ConvertDecCoinToDenom(coin sdk.DecCoin, denom string) (sdk.Coin, error) {
	decCoin, err := sdk.ConvertDecCoin(coin, denom)
	if err != nil {
		return sdk.Coin{}, err
	}
	truncated, _ := decCoin.TruncateDecimal()
	return truncated, nil
}
