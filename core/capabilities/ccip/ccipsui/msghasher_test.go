package ccipsui

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"
	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

func TestParseExtraDataMap(t *testing.T) {
	tokenReceiverExample := [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20}

	solanaSelector := ccipocr3common.ChainSelector(chainsel.SOLANA_DEVNET.Selector)
	nonSolanaSelector := ccipocr3common.ChainSelector(chainsel.SUI_TESTNET.Selector)

	tests := []struct {
		name           string
		input          map[string]any
		sourceSelector ccipocr3common.ChainSelector
		want           *struct {
			gasLimit      *big.Int
			tokenReceiver [32]byte
		}
		expectErr bool
	}{
		{
			name: "valid input with [32]byte for tokenReceiver",
			input: map[string]any{
				"gasLimit":      new(big.Int).SetInt64(500000),
				"tokenReceiver": [32]byte{0x01},
			},
			sourceSelector: nonSolanaSelector,
			want: &struct {
				gasLimit      *big.Int
				tokenReceiver [32]byte
			}{
				gasLimit:      new(big.Int).SetInt64(500000),
				tokenReceiver: [32]byte{0x01},
			},
			expectErr: false,
		},
		{
			name: "valid input with []byte for tokenReceiver",
			input: map[string]any{
				"gasLimit":      new(big.Int).SetInt64(500000),
				"tokenReceiver": tokenReceiverExample[:], // convert to slice for input
			},
			sourceSelector: nonSolanaSelector,
			want: &struct {
				gasLimit      *big.Int
				tokenReceiver [32]byte
			}{
				gasLimit:      new(big.Int).SetInt64(500000),
				tokenReceiver: tokenReceiverExample,
			},
			expectErr: false,
		},
		{
			name: "invalid input with []byte for tokenReceiver",
			input: map[string]any{
				"gasLimit":      new(big.Int).SetInt64(500000),
				"tokenReceiver": tokenReceiverExample[:16], // 16 bytes, we expect an error due to length
			},
			sourceSelector: nonSolanaSelector,
			want:           nil,
			expectErr:      true,
		},
		{
			// Solana source emits GenericExtraArgsV2 (Borsh); ccipsolana.ExtraDataDecoder
			// surfaces the gas limit under the struct field name "GasLimit" (capital) and
			// does not include a tokenReceiver (GenericExtraArgsV2 has no token_receiver).
			// parseExtraDataMap accepts the capital key and, for a Solana source, defaults
			// tokenReceiver to zero.
			name: "Solana GenericExtraArgsV2: capital GasLimit, no tokenReceiver defaults to zero",
			input: map[string]any{
				"GasLimit":                 new(big.Int).SetInt64(1000000),
				"AllowOutOfOrderExecution": true,
			},
			sourceSelector: solanaSelector,
			want: &struct {
				gasLimit      *big.Int
				tokenReceiver [32]byte
			}{
				gasLimit:      new(big.Int).SetInt64(1000000),
				tokenReceiver: [32]byte{},
			},
			expectErr: false,
		},
		{
			// A Solana source that does carry a tokenReceiver uses it as-is; the zero default
			// only applies when the key is absent.
			name: "Solana source with tokenReceiver present uses it",
			input: map[string]any{
				"GasLimit":      new(big.Int).SetInt64(1000000),
				"tokenReceiver": [32]byte{0x0A},
			},
			sourceSelector: solanaSelector,
			want: &struct {
				gasLimit      *big.Int
				tokenReceiver [32]byte
			}{
				gasLimit:      new(big.Int).SetInt64(1000000),
				tokenReceiver: [32]byte{0x0A},
			},
			expectErr: false,
		},
		{
			name: "no gas limit key of either casing",
			input: map[string]any{
				"tokenReceiver": [32]byte{0x01},
			},
			sourceSelector: nonSolanaSelector,
			want:           nil,
			expectErr:      true,
		},
		{
			// SuiExtraArgsV1 (EVM/Sui sources) uses lowercase "gasLimit" and MUST carry
			// tokenReceiver. A missing tokenReceiver on a non-Solana source is a
			// malformed/partial decode and must error — the zero default is scoped to Solana
			// sources only, so a bad EVM/Sui decode is surfaced rather than silently hashed.
			name: "non-Solana source missing tokenReceiver errors (masking guard)",
			input: map[string]any{
				"gasLimit": new(big.Int).SetInt64(500000),
			},
			sourceSelector: nonSolanaSelector,
			want:           nil,
			expectErr:      true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gasLimit, tokenReceiver, err := parseExtraDataMap(tt.input, tt.sourceSelector)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want.gasLimit, gasLimit)
				assert.Equal(t, tt.want.tokenReceiver, tokenReceiver)
			}
		})
	}
}
