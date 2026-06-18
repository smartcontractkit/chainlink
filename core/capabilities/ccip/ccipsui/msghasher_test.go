package ccipsui

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseExtraDataMap(t *testing.T) {
	tokenReceiverExample := [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20}

	tests := []struct {
		name  string
		input map[string]any
		want  *struct {
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
			want:      nil,
			expectErr: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gasLimit, tokenReceiver, err := parseExtraDataMap(tt.input)
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
