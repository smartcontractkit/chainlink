package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMergedRegistry(t *testing.T) {
	tests := []struct {
		name            string
		newTokens       map[string][]TokenSymbol
		testDescription string
		expectedFound   bool
	}{
		{
			name:            "Empty input - should only contain defaults",
			newTokens:       map[string][]TokenSymbol{},
			testDescription: USDCUSD,
			expectedFound:   true,
		},
		{
			name: "Add new token for existing description - should merge without duplicates",
			newTokens: map[string][]TokenSymbol{
				USDCUSD: {"TESTUSDC"},
			},
			testDescription: USDCUSD,
			expectedFound:   true,
		},
		{
			name: "Add duplicate token for existing description - should not duplicate",
			newTokens: map[string][]TokenSymbol{
				EthUSD: {WethSymbol}, // Adding WethSymbol to EthUSD which already has WethSymbol
			},
			testDescription: EthUSD,
			expectedFound:   true,
		},
		{
			name: "Add new description with tokens",
			newTokens: map[string][]TokenSymbol{
				"NEW / USD": {CCIPBnMSymbol, CCIPLnMSymbol},
			},
			testDescription: "NEW / USD",
			expectedFound:   true,
		},
		{
			name: "Add both new token and duplicated token for existing description - should merge without duplicates",
			newTokens: map[string][]TokenSymbol{
				USDCUSD: {"TESTUSDC"},             // Adding new TESTUSDC to USDC/USD
				EthUSD:  {WethSymbol, "TESTWETH"}, // Adding WethSymbol & TestWETH to EthUSD which already has WethSymbol
			},
			testDescription: USDCUSD,
			expectedFound:   true,
		},
		{
			name: "Query non-existent description",
			newTokens: map[string][]TokenSymbol{
				"REAL / USD": {APTSymbol},
			},
			testDescription: "NON_EXISTENT / USD",
			expectedFound:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get symbols before merge (from default registry)
			symbolsBeforeMerge, _ := registry.GetSymbols(tt.testDescription)

			// Create merged registry
			registry := NewMergedRegistry(tt.newTokens)

			// Get symbols after merge
			symbolsAfterMerge, found := registry.GetSymbols(tt.testDescription)
			assert.Equal(t, tt.expectedFound, found, "Expected found status mismatch")

			if tt.expectedFound {
				// Check if all the symbols before merge is present in after merge
				for _, beforeMergeSymbol := range symbolsBeforeMerge {
					assert.Contains(t, symbolsAfterMerge, beforeMergeSymbol, "Expected symbol %s to be present in merged registry", beforeMergeSymbol)
				}

				// Check that new symbols from the input are present
				if newSymbols, hasNewSymbols := tt.newTokens[tt.testDescription]; hasNewSymbols {
					for _, newSymbol := range newSymbols {
						assert.Contains(t, symbolsAfterMerge, newSymbol, "New symbol %s should be present after merge", newSymbol)
					}
				}
			} else {
				assert.Nil(t, symbolsAfterMerge, "Expected symbols to be nil when not found")
			}
		})
	}
}
