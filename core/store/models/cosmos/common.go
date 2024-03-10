package cosmos

import sdk "github.com/cosmos/cosmos-sdk/types"

// SendRequest represents a request to transfer Cosmos coins.
type SendRequest struct {
	DestinationAddress sdk.AccAddress `json:"address"`
	FromAddress        sdk.AccAddress `json:"from"`
	Amount             sdk.Dec        `json:"amount"`
	CosmosChainID      string         `json:"cosmosChainID"`
	Token              string         `json:"token"`
	AllowHigherAmounts bool           `json:"allowHigherAmounts"`
}
