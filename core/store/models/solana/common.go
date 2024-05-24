package solana

import "github.com/gagliardetto/solana-go"

type SendRequest struct {
	From               solana.PublicKey `json:"from"`
	To                 solana.PublicKey `json:"to"`
	Amount             uint64           `json:"amount"`
	SolanaChainID      string           `json:"solanaChainID"`
	AllowHigherAmounts bool             `json:"allowHigherAmounts"`
}
