package cosmos

import (
	"math/big"
)

// SendRequest represents a request to transfer Cosmos coins.
type SendRequest struct {
	DestinationAddress string   `json:"address"`
	FromAddress        string   `json:"from"`
	Amount             *big.Int `json:"amount"`
	CosmosChainID      string   `json:"cosmosChainID"`
	Token              string   `json:"token"`
	AllowHigherAmounts bool     `json:"allowHigherAmounts"`
}
