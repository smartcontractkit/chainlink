package test

import (
	"context"
	"fmt"
	"math/big"

	"github.com/stretchr/testify/assert"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// PriceGetter is a static test implementation of
// [testtypes.Evaluator] for [cciptypes.PriceGetter].
var PriceGetter = staticPriceGetter{
	config: staticPriceGetterConfig{
		Prices: map[cciptypes.Address]*big.Int{
			"ETH":  big.NewInt(7),
			"LINK": big.NewInt(11),
		},
		Addresses: []cciptypes.Address{"ETH", "LINK"},
	},
}

type PriceGetterEvaluator interface {
	cciptypes.PriceGetter
	testtypes.Evaluator[cciptypes.PriceGetter]
}
type staticPriceGetterConfig struct {
	Prices    map[cciptypes.Address]*big.Int
	Addresses []cciptypes.Address
}
type staticPriceGetter struct {
	config staticPriceGetterConfig
}

var _ PriceGetterEvaluator = staticPriceGetter{}

// Close implements ccip.PriceGetter.
func (s staticPriceGetter) Close() error {
	return nil
}

// FilterConfiguredTokens implements ccip.PriceGetter.
func (s staticPriceGetter) FilterConfiguredTokens(ctx context.Context, tokens []cciptypes.Address) (configured []cciptypes.Address, unconfigured []cciptypes.Address, err error) {
	configured = []cciptypes.Address{}
	unconfigured = []cciptypes.Address{}

	for _, tk := range tokens {
		found := false
		for _, addr := range s.config.Addresses {
			if addr == tk {
				found = true
				configured = append(configured, tk)
				break
			}
		}
		if !found {
			unconfigured = append(unconfigured, tk)
		}
	}

	return configured, unconfigured, nil
}

// TokenPricesUSD implements ccip.PriceGetter.
func (s staticPriceGetter) TokenPricesUSD(ctx context.Context, tokens []cciptypes.Address) (map[cciptypes.Address]*big.Int, error) {
	if ok := assert.ObjectsAreEqual(s.config.Addresses, tokens); !ok {
		return nil, fmt.Errorf("unexpected tokens: expected %v, got %v", s.config.Addresses, tokens)
	}
	return s.config.Prices, nil
}

// Evaluate implements types_test.Evaluator.
func (s staticPriceGetter) Evaluate(ctx context.Context, other cciptypes.PriceGetter) error {
	got, err := other.TokenPricesUSD(ctx, s.config.Addresses)
	if err != nil {
		return fmt.Errorf("failed to get prices: %w", err)
	}
	ok := assert.ObjectsAreEqualValues(s.config.Prices, got)
	if !ok {
		return fmt.Errorf("unexpected prices: expected %v, got %v", s.config.Prices, got)
	}
	return nil
}
