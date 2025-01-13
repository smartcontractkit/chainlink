package solana

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
)

func TestChainWriterConfigRaw(t *testing.T) {
	tests := []struct {
		name          string
		fromAddress   string
		expectedError string
	}{
		{
			name:          "valid input",
			fromAddress:   "4Nn9dsYBcSTzRbK9hg9kzCUdrCSkMZq1UR6Vw1Tkaf6A",
			expectedError: "",
		},
		{
			name:          "zero fromAddress",
			fromAddress:   "",
			expectedError: "invalid from address : decode: zero length string",
		},
	}

	commonAddressesLookupTable := solana.MustPublicKeyFromBase58("4Nn9dsYBcSTzRbK9hg9kzCUdrCSkMZq1UR6Vw1Tkaf6H")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetSolanaChainWriterConfig("4Nn9dsYBcSTzRbK9hg9kzCUdrCSkMZq1UR6Vw1Tkaf6H", commonAddressesLookupTable, tt.fromAddress)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
