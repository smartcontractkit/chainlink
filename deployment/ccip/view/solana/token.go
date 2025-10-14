package solana

import (
	"context"
	"fmt"
	"strings"

	solTokenMetadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"github.com/gagliardetto/solana-go"
	solToken "github.com/gagliardetto/solana-go/programs/token"

	solanashared "github.com/smartcontractkit/chainlink/deployment"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
)

type TokenView struct {
	TokenProgramName string        `json:"tokenProgramName,omitempty"`
	MintAuthority    string        `json:"mintAuthority,omitempty"`
	Supply           uint64        `json:"supply,omitempty"`
	Decimals         uint8         `json:"decimals,omitempty"`
	IsInitialized    bool          `json:"isInitialized,omitempty"`
	FreezeAuthority  string        `json:"freezeAuthority,omitempty"`
	TokenMetadata    TokenMetadata `json:"tokenMetadata"`
}

type TokenMetadata struct {
	UpdateAuthority      string `json:"updateAuthority,omitempty"`
	Name                 string `json:"name,omitempty"`
	Symbol               string `json:"symbol,omitempty"`
	URI                  string `json:"uri,omitempty"`
	SellerFeeBasisPoints uint16 `json:"sellerFeeBasisPoints,omitempty"`
	PrimarySaleHappened  bool   `json:"primarySaleHappened,omitempty"`
	IsMutable            bool   `json:"isMutable,omitempty"`
}

func GenerateTokenView(chain cldf_solana.Chain, tokenAddress solana.PublicKey, tokenProgram string) (TokenView, error) {
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
	var tokenMetadata solTokenMetadata.Metadata
	metadataPDA, err := solanashared.FindMplTokenMetadataPDA(tokenAddress)
	if err != nil {
		return view, fmt.Errorf("failed to find metadata PDA: %w", err)
	}
	// if no metadata, don't return an error
	if err = chain.GetAccountDataBorshInto(context.Background(), metadataPDA, &tokenMetadata); err == nil {
		view.TokenMetadata = TokenMetadata{}
		view.TokenMetadata.UpdateAuthority = tokenMetadata.UpdateAuthority.String()
		view.TokenMetadata.Name = strings.ReplaceAll(tokenMetadata.Data.Name, "\x00", "")
		view.TokenMetadata.Symbol = strings.ReplaceAll(tokenMetadata.Data.Symbol, "\x00", "")
		view.TokenMetadata.URI = strings.ReplaceAll(tokenMetadata.Data.Uri, "\x00", "")
		view.TokenMetadata.SellerFeeBasisPoints = tokenMetadata.Data.SellerFeeBasisPoints
		view.TokenMetadata.PrimarySaleHappened = tokenMetadata.PrimarySaleHappened
		view.TokenMetadata.IsMutable = tokenMetadata.IsMutable
	}
	return view, nil
}
