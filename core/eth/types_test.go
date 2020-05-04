package eth_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/core/eth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"

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

func TestSafeByteSlice_Success(t *testing.T) {
	tests := []struct {
		ary      eth.UntrustedBytes
		start    int
		end      int
		expected []byte
	}{
		{[]byte{1, 2, 3}, 0, 0, []byte{}},
		{[]byte{1, 2, 3}, 0, 1, []byte{1}},
		{[]byte{1, 2, 3}, 1, 3, []byte{2, 3}},
	}

	for i, test := range tests {
		t.Run(string(i), func(t *testing.T) {
			actual, err := test.ary.SafeByteSlice(test.start, test.end)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestSafeByteSlice_Error(t *testing.T) {
	tests := []struct {
		ary   eth.UntrustedBytes
		start int
		end   int
	}{
		{[]byte{1, 2, 3}, 2, -1},
		{[]byte{1, 2, 3}, 0, 4},
		{[]byte{1, 2, 3}, 3, 4},
		{[]byte{1, 2, 3}, 3, 2},
		{[]byte{1, 2, 3}, -1, 2},
	}

	for i, test := range tests {
		t.Run(string(i), func(t *testing.T) {
			actual, err := test.ary.SafeByteSlice(test.start, test.end)
			assert.EqualError(t, err, "out of bounds slice access")
			var expected []byte
			assert.Equal(t, expected, actual)
		})
	}
}

func TestBlock_Unmarshal(t *testing.T) {
	t.Parallel()

	input := `{
		"number": "0x1b4", 
		"hash": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
		"parentHash": "0x9646252be9520f6e71339a8df9c55e4d7619deeb018d2a3f2d21fc165dde5eb5",
		"nonce": "0xe04d296d2460cfb8472af2c5fd05b5a214109c25688d3704aed5484f9a7792f2",
		"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"logsBloom": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
		"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"stateRoot": "0xd5855eb08b3387c0af375e9cdb6acfc05eb8f519e419b874b6ff2ffda7ed1dff",
		"miner": "0x4e65fda2159562a496f9f3522f89122a3088497a",
		"difficulty": "0xFFFFFFFFFFFF9DDB99A168BD2A000001", 
		"totalDifficulty":  "0x027f07", 
		"extraData": "0x0000000000000000000000000000000000000000000000000000000000000000",
		"size":  "0x027f07", 
		"gasLimit": "0x9f759", 
		"gasUsed": "0x9f759", 
		"timestamp": "0x54e34e8e",
		"transactions": [
			{
			"blockHash":"0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
			"blockNumber":"0x5daf3b", 
			"from":"0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
			"gas":"0xc350", 
			"gasPrice":"0x4a817c800", 
			"hash":"0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b",
			"input":"0x68656c6c6f21",
			"nonce":"0x15", 
			"to":"0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
			"transactionIndex":"0x41", 
			"value":"0xf3dbb76162000", 
			"v":"0x25", 
			"r":"0x1b5e176d927f8e9ab405058b2d2457392da3e20f328b16ddabcebc33eaac5fea",
			"s":"0x4ba69724e8f69de52f0125ad8b3c5c2cef33019bac3249e2c0a2192766d1721c"
		  },
		  {
			"blockHash":"0x1d59ff54b1eb26b013ce3cb5fc9dab3705b415a67127a003c3e61eb445bb8df2",
			"blockNumber":"0x5daf3b", 
			"from":"0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
			"gas":"0xc350", 
			"gasPrice":"0x4a817c801", 
			"hash":"0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b",
			"input":"0x68656c6c6f21",
			"nonce":"0x15", 
			"to":"0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
			"transactionIndex":"0x41", 
			"value":"0xf3dbb76162000", 
			"v":"0x25", 
			"r":"0x1b5e176d927f8e9ab405058b2d2457392da3e20f328b16ddabcebc33eaac5fea",
			"s":"0x4ba69724e8f69de52f0125ad8b3c5c2cef33019bac3249e2c0a2192766d1721c"
		  }
		], 
		"uncles": ["0x1606e5", "0xd5145a9"]
	}`

	var block eth.Block
	err := json.Unmarshal([]byte(input), &block)
	require.NoError(t, err)

	assert.Len(t, block.Transactions, 2)
}
