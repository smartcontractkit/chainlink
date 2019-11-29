package eth_test

import (
	"chainlink/core/eth"
	"chainlink/core/internal/cltest"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestLog_UnmarshalEmptyTxHash(t *testing.T) {
	t.Parallel()

	input := `{
		"transactionHash": null,
		"transactionIndex": "0x3",
		"address": "0x1aee7c03606fca5035d204c3818d0660bb230e44",
		"blockNumber": "0x8bf99b",
		"topics": ["0xdeadbeefdeadbeedeadbeedeadbeefffdeadbeefdeadbeedeadbeedeadbeefff"],
		"blockHash": "0xdb777676330c067e3c3a6dbfc2d51282cac5bcc1b7a884dd8d85ba72ca1f147e",
		"data": "0xdeadbeef",
		"logIndex": "0x5",
		"transactionLogIndex": "0x3"
	}`

	var log eth.Log
	err := json.Unmarshal([]byte(input), &log)
	assert.NoError(t, err)
}

func TestReceipt_UnmarshalEmptyBlockHash(t *testing.T) {
	t.Parallel()

	input := `{
		"transactionHash": "0x444172bef57ad978655171a8af2cfd89baa02a97fcb773067aef7794d6913374",
		"blockNumber": "0x8bf99b",
		"blockHash": null
	}`

	var receipt eth.TxReceipt
	err := json.Unmarshal([]byte(input), &receipt)
	require.NoError(t, err)
}

func TestModels_HexToFunctionSelector(t *testing.T) {
	t.Parallel()
	fid := eth.HexToFunctionSelector("0xb3f98adc")
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_HexToFunctionSelectorOverflow(t *testing.T) {
	t.Parallel()
	fid := eth.HexToFunctionSelector("0xb3f98adc123456")
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_FunctionSelectorUnmarshalJSON(t *testing.T) {
	t.Parallel()
	bytes := []byte(`"0xb3f98adc"`)
	var fid eth.FunctionSelector
	err := json.Unmarshal(bytes, &fid)
	assert.NoError(t, err)
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_FunctionSelectorUnmarshalJSONLiteral(t *testing.T) {
	t.Parallel()
	literalSelectorBytes := []byte(`"setBytes(bytes)"`)
	var fid eth.FunctionSelector
	err := json.Unmarshal(literalSelectorBytes, &fid)
	assert.NoError(t, err)
	assert.Equal(t, "0xda359dc8", fid.String())
}

func TestModels_FunctionSelectorUnmarshalJSONError(t *testing.T) {
	t.Parallel()
	bytes := []byte(`"0xb3f98adc123456"`)
	var fid eth.FunctionSelector
	err := json.Unmarshal(bytes, &fid)
	assert.Error(t, err)
}

func TestModels_Header_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		path       string
		wantNumber hexutil.Big
		wantHash   string
	}{
		{
			"parity",
			"testdata/subscription_new_heads_parity.json",
			cltest.BigHexInt(1263817),
			"0xf8e4691ceab8052d1cb478c6c5e0d9b122e747ad838023633f63bd5e81ec5114",
		},
		{
			"geth",
			"testdata/subscription_new_heads_geth.json",
			cltest.BigHexInt(1263817),
			"0xf8e4691ceab8052d1cb478c6c5e0d9b122e747ad838023633f63bd5e81ec5fff",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var header eth.BlockHeader

			data := cltest.MustReadFile(t, test.path)
			value := gjson.Get(string(data), "params.result")
			assert.NoError(t, json.Unmarshal([]byte(value.String()), &header))

			assert.Equal(t, test.wantNumber, header.Number)
			assert.Equal(t, test.wantHash, header.Hash().String())
		})
	}
}
