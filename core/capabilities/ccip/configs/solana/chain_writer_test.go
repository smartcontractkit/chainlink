package solana

import (
	"encoding/json"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainwriter"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
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
			config, err := GetSolanaChainWriterConfig("4Nn9dsYBcSTzRbK9hg9kzCUdrCSkMZq1UR6Vw1Tkaf6H", commonAddressesLookupTable, tt.fromAddress)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)

				raw, err := json.Marshal(config)
				require.NoError(t, err)
				var result chainwriter.ChainWriterConfig
				err = json.Unmarshal(raw, &result)
				require.NoError(t, err)
				require.EqualValues(t, config, result)
			}
		})
	}
}
