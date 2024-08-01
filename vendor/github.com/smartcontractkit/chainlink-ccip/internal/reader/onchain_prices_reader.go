package reader

import (
	"context"
	"fmt"
	"math/big"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"golang.org/x/sync/errgroup"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

type TokenPrices interface {
	// GetTokenPricesUSD returns the prices of the provided tokens in USD.
	// The order of the returned prices corresponds to the order of the provided tokens.
	GetTokenPricesUSD(ctx context.Context, tokens []ocr2types.Account) ([]*big.Int, error)
}

type TokenPriceConfig struct {
	// This is mainly used for inputTokens on testnet to give them a price
	StaticPrices map[ocr2types.Account]big.Int `json:"staticPrices"`
}

type OnchainTokenPricesReader struct {
	TokenPriceConfig TokenPriceConfig
	// Reader for the chain that will have the token prices on-chain
	ContractReader commontypes.ContractReader
}

func NewOnchainTokenPricesReader(
	tokenPriceConfig TokenPriceConfig, contractReader commontypes.ContractReader,
) *OnchainTokenPricesReader {
	return &OnchainTokenPricesReader{
		TokenPriceConfig: tokenPriceConfig,
		ContractReader:   contractReader,
	}
}

func (pr *OnchainTokenPricesReader) GetTokenPricesUSD(
	ctx context.Context, tokens []ocr2types.Account,
) ([]*big.Int, error) {
	const (
		contractName = "PriceAggregator"
		functionName = "getTokenPrice"
	)
	prices := make([]*big.Int, len(tokens))
	eg := new(errgroup.Group)
	for idx, token := range tokens {
		idx := idx
		token := token
		eg.Go(func() error {
			price := new(big.Int)
			if staticPrice, exists := pr.TokenPriceConfig.StaticPrices[token]; exists {
				price.Set(&staticPrice)
			} else {
				if err :=
					pr.ContractReader.GetLatestValue(
						ctx,
						contractName,
						functionName,
						primitives.Finalized,
						token,
						price); err != nil {
					return fmt.Errorf("failed to get token price for %s: %w", token, err)
				}
			}
			prices[idx] = price
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get all token prices successfully: %w", err)
	}

	for _, price := range prices {
		if price == nil {
			return nil, fmt.Errorf("failed to get all token prices successfully, some prices are nil")
		}
	}

	return prices, nil
}

// Ensure OnchainTokenPricesReader implements TokenPrices
var _ TokenPrices = (*OnchainTokenPricesReader)(nil)
