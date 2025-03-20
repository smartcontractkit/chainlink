package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	solToken "github.com/gagliardetto/solana-go/programs/token"

	"github.com/smartcontractkit/chainlink/deployment"
)

type TokenView struct {
	TokenProgramName string `json:"tokenProgramName,omitempty"`
	MintAuthority    string `json:"mintAuthority,omitempty"`
	Supply           uint64 `json:"supply,omitempty"`
	Decimals         uint8  `json:"decimals,omitempty"`
	IsInitialized    bool   `json:"isInitialized,omitempty"`
	FreezeAuthority  string `json:"freezeAuthority,omitempty"`
}

func GenerateTokenView(chain deployment.SolChain, tokenAddress solana.PublicKey, tokenProgram string) (TokenView, error) {
	view := TokenView{}
	view.TokenProgramName = tokenProgram
	var tokenMint solToken.Mint
	err := chain.GetAccountDataBorshInto(context.Background(), tokenAddress, &tokenMint)
	if err != nil {
		return view, fmt.Errorf("token not found in existing state %d", chain.Selector)
	}
	if tokenMint.MintAuthority == nil {
		view.MintAuthority = "None"
	} else {
		view.MintAuthority = tokenMint.MintAuthority.String()
	}
	view.Supply = tokenMint.Supply
	view.Decimals = tokenMint.Decimals
	view.IsInitialized = tokenMint.IsInitialized
	if tokenMint.FreezeAuthority == nil {
		view.FreezeAuthority = "None"
	} else {
		view.FreezeAuthority = tokenMint.FreezeAuthority.String()
	}
	return view, nil
}
