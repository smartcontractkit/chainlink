package shared

import (
	"encoding/base64"
	"testing"

	bin "github.com/gagliardetto/binary"
	solToken "github.com/gagliardetto/solana-go/programs/token"
	"github.com/stretchr/testify/require"
)

// mainnetSPLMintAccount is the live account data of a 6-decimal SPL mint, captured from
// Solana mainnet. It pins the layout the decimals read depends on.
const mainnetSPLMintAccount = "AQAAAJqmjmC3XUgeV/3BTNcDsJeT3fWVS1LxIRu9+UDUwYbiKgCizQheBAAGAQEAAABpABO3JyzZq1lqrNeYIqP0v/TOjpRY5+Ltw3RArhK21A=="

// TestSolanaMintDecode confirms a real mainnet mint account decodes to the expected
// decimals, covering the layout assumption behind solanaMintDecimals.
func TestSolanaMintDecode(t *testing.T) {
	t.Parallel()

	data, err := base64.StdEncoding.DecodeString(mainnetSPLMintAccount)
	require.NoError(t, err)
	require.Len(t, data, 82, "SPL mint accounts are 82 bytes")

	var mint solToken.Mint
	require.NoError(t, bin.NewBinDecoder(data).Decode(&mint))

	require.Equal(t, uint8(6), mint.Decimals)
	require.True(t, mint.IsInitialized)
}

func TestParseAptosDecimals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		values     []any
		want       uint8
		wantErrMsg string
	}{
		{
			// Shape returned by 0x1::fungible_asset::decimals on mainnet.
			name:   "Success - single number",
			values: []any{float64(6)},
			want:   6,
		},
		{
			name:   "Success - zero decimals",
			values: []any{float64(0)},
			want:   0,
		},
		{
			name:   "Success - max uint8",
			values: []any{float64(255)},
			want:   255,
		},
		{
			name:       "Failure - no values",
			values:     nil,
			wantErrMsg: "returned no value",
		},
		{
			name:       "Failure - string instead of number",
			values:     []any{"6"},
			wantErrMsg: "want a number",
		},
		{
			name:       "Failure - above uint8",
			values:     []any{float64(256)},
			wantErrMsg: "not a uint8",
		},
		{
			name:       "Failure - negative",
			values:     []any{float64(-1)},
			wantErrMsg: "not a uint8",
		},
		{
			name:       "Failure - fractional",
			values:     []any{float64(6.5)},
			wantErrMsg: "not a uint8",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseAptosDecimals(tc.values)
			if tc.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
