package types

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func mustBig(t *testing.T, s string) *big.Int {
	t.Helper()
	v, ok := new(big.Int).SetString(s, 10)
	require.True(t, ok, "bad big.Int literal %q", s)
	return v
}

func TestCCIPBridgeTransfer_UnmarshalAmount(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		json string
		want *big.Int
	}{
		{
			name: "unquoted small",
			json: `{"to_chain_selector":5009297550715157269,"amount":3000000000000000,"destination_address":"0xabc"}`,
			want: big.NewInt(3000000000000000),
		},
		{
			name: "quoted small",
			json: `{"to_chain_selector":5009297550715157269,"amount":"3000000000000000","destination_address":"0xabc"}`,
			want: big.NewInt(3000000000000000),
		},
		{
			name: "quoted large (50 ETH, exceeds int64)",
			json: `{"to_chain_selector":5009297550715157269,"amount":"50000000000000000000","destination_address":"0xabc"}`,
			want: mustBig(t, "50000000000000000000"),
		},
		{
			name: "quoted very large (100 ETH)",
			json: `{"to_chain_selector":5009297550715157269,"amount":"100000000000000000000","destination_address":"0xabc"}`,
			want: mustBig(t, "100000000000000000000"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			var tr CCIPBridgeTransfer
			require.NoError(t, json.Unmarshal([]byte(c.json), &tr))
			require.NotNil(t, tr.Amount)
			require.Zero(t, tr.Amount.Cmp(c.want), "amount: got %s want %s", tr.Amount, c.want)
			require.Equal(t, uint64(5009297550715157269), tr.ToChainSelector)
			require.Equal(t, "0xabc", tr.DestinationAddress)
		})
	}
}

func TestCCIPBridgeTransfer_UnmarshalMaxCCIPFeeAndOmitted(t *testing.T) {
	t.Parallel()

	// max_ccip_fee as a quoted string is parsed.
	var withFee CCIPBridgeTransfer
	require.NoError(t, json.Unmarshal(
		[]byte(`{"to_chain_selector":1,"amount":"1","max_ccip_fee":"20000000000000000"}`), &withFee))
	require.NotNil(t, withFee.MaxCCIPFee)
	require.Zero(t, withFee.MaxCCIPFee.Cmp(mustBig(t, "20000000000000000")))

	// Omitted max_ccip_fee stays nil (optional field).
	var noFee CCIPBridgeTransfer
	require.NoError(t, json.Unmarshal([]byte(`{"to_chain_selector":1,"amount":"1"}`), &noFee))
	require.Nil(t, noFee.MaxCCIPFee)

	// Invalid amount is rejected.
	var bad CCIPBridgeTransfer
	require.Error(t, json.Unmarshal([]byte(`{"to_chain_selector":1,"amount":"not-a-number"}`), &bad))
}

// TestCCIPBridgeConfig_UnmarshalTransfers ensures the per-transfer custom unmarshaler is invoked when
// decoding the full config payload.
func TestCCIPBridgeConfig_UnmarshalTransfers(t *testing.T) {
	t.Parallel()

	payload := `{
		"from_chain_selector": 15971525489660198786,
		"token_type": "native",
		"test_needed": true,
		"transfers": [
			{"to_chain_selector": 5009297550715157269, "amount": "50000000000000000000", "destination_address": "0xA9B195EaB94c3877B67FbA9b6Cdc68558D1fc83b"}
		]
	}`

	var cfg CCIPBridgeConfig
	require.NoError(t, json.Unmarshal([]byte(payload), &cfg))
	require.Len(t, cfg.Transfers, 1)
	require.Equal(t, TokenTypeNative, cfg.TokenType)
	require.True(t, cfg.TestNeeded)
	require.Zero(t, cfg.Transfers[0].Amount.Cmp(mustBig(t, "50000000000000000000")))
}
