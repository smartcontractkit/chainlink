package models_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHead_NewHead(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input *big.Int
		want  string
	}{
		{big.NewInt(0), "0"},
		{big.NewInt(0xf), "f"},
		{big.NewInt(0x10), "10"},
	}
	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			num := models.NewHead(test.input, cltest.NewHash(), cltest.NewHash(), 0)
			assert.Equal(t, test.want, fmt.Sprintf("%x", num.ToInt()))
		})
	}
}

func TestHead_GreaterThan(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		left    *models.Head
		right   *models.Head
		greater bool
	}{
		{"nil nil", nil, nil, false},
		{"present nil", cltest.Head(1), nil, true},
		{"nil present", nil, cltest.Head(1), false},
		{"less", cltest.Head(1), cltest.Head(2), false},
		{"equal", cltest.Head(2), cltest.Head(2), false},
		{"greater", cltest.Head(2), cltest.Head(1), true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.greater, test.left.GreaterThan(test.right))
		})
	}
}

func TestHead_NextInt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		bn   *models.Head
		want *big.Int
	}{
		{"nil", nil, nil},
		{"one", cltest.Head(1), big.NewInt(2)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.bn.NextInt())
		})
	}
}

func TestTx_PresenterMatchesHex(t *testing.T) {
	t.Parallel()

	createdTx := models.Tx{
		From:        common.HexToAddress("0xf208"),
		To:          common.HexToAddress("0x70"),
		Data:        []byte(`{"data": "is wilding out"}`),
		Nonce:       0x8008,
		Value:       utils.NewBig(big.NewInt(777)),
		GasLimit:    1999,
		Hash:        common.HexToHash("0x0"),
		GasPrice:    utils.NewBig(big.NewInt(333)),
		Confirmed:   true,
		SentAt:      1745,
		SignedRawTx: hexutil.MustDecode("0xcafe"),
	}

	ptx := presenters.NewTx(&createdTx)
	bytes, err := json.Marshal(ptx)
	require.NoError(t, err)
	assert.JSONEq(t, `{`+
		`"confirmed":true,`+
		`"data":"0x7b2264617461223a202269732077696c64696e67206f7574227d",`+
		`"from":"0x000000000000000000000000000000000000f208",`+
		`"gasLimit":"1999",`+
		`"gasPrice":"333",`+
		`"hash":"0x0000000000000000000000000000000000000000000000000000000000000000",`+
		`"rawHex":"0xcafe",`+
		`"nonce":"32776",`+
		`"sentAt":"1745",`+
		`"to":"0x0000000000000000000000000000000000000070",`+
		`"value":"777"`+
		`}`, string(bytes))
}

func TestHighestPricedTxAttemptPerTx(t *testing.T) {
	items := []models.TxAttempt{
		{TxID: 1, GasPrice: utils.NewBig(big.NewInt(5555))},
		{TxID: 1, GasPrice: utils.NewBig(big.NewInt(444))},
		{TxID: 1, GasPrice: utils.NewBig(big.NewInt(2))},
		{TxID: 1, GasPrice: utils.NewBig(big.NewInt(33333))},
		{TxID: 2, GasPrice: utils.NewBig(big.NewInt(4444))},
		{TxID: 2, GasPrice: utils.NewBig(big.NewInt(999))},
		{TxID: 2, GasPrice: utils.NewBig(big.NewInt(12211))},
	}

	items = models.HighestPricedTxAttemptPerTx(items)

	sort.Slice(items, func(i, j int) bool { return items[i].TxID < items[j].TxID })

	assert.Len(t, items, 2)
	assert.True(t, items[0].GasPrice.ToInt().Cmp(big.NewInt(33333)) == 0)
	assert.True(t, items[1].GasPrice.ToInt().Cmp(big.NewInt(12211)) == 0)
}

func TestEthTxAttempt_GetSignedTx(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	// Use the real KeyStore loaded from database fixtures
	store.KeyStore.Unlock(cltest.Password)
	tx := gethTypes.NewTransaction(uint64(42), cltest.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})

	keys, err := store.Keys()
	require.NoError(t, err)
	key := keys[0]
	fromAddress := key.Address.Address()
	account, err := store.KeyStore.GetAccountByAddress(fromAddress)
	require.NoError(t, err)

	chainID := big.NewInt(3)

	signedTx, err := store.KeyStore.SignTx(account, tx, chainID)
	require.NoError(t, err)
	signedTx.Size() // Needed to write the size for equality checking
	rlp := new(bytes.Buffer)
	require.NoError(t, signedTx.EncodeRLP(rlp))

	attempt := models.EthTxAttempt{SignedRawTx: rlp.Bytes()}

	gotSignedTx, err := attempt.GetSignedTx()
	require.NoError(t, err)
	decodedEncoded := new(bytes.Buffer)
	require.NoError(t, gotSignedTx.EncodeRLP(decodedEncoded))

	require.Equal(t, signedTx, gotSignedTx)
	require.Equal(t, attempt.SignedRawTx, decodedEncoded.Bytes())
}

func TestHead_ChainLength(t *testing.T) {
	head := models.Head{
		Parent: &models.Head{
			Parent: &models.Head{},
		},
	}

	assert.Equal(t, uint32(3), head.ChainLength())
}

func TestReceipt_UnmarshalEmptyBlockHash(t *testing.T) {
	t.Parallel()

	input := `{
        "transactionHash": "0x444172bef57ad978655171a8af2cfd89baa02a97fcb773067aef7794d6913374",
        "blockNumber": "0x8bf99b",
        "blockHash": null
    }`

	var receipt models.TxReceipt
	err := json.Unmarshal([]byte(input), &receipt)
	require.NoError(t, err)
}

func TestModels_HexToFunctionSelector(t *testing.T) {
	t.Parallel()
	fid := models.HexToFunctionSelector("0xb3f98adc")
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_HexToFunctionSelectorOverflow(t *testing.T) {
	t.Parallel()
	fid := models.HexToFunctionSelector("0xb3f98adc123456")
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_FunctionSelectorUnmarshalJSON(t *testing.T) {
	t.Parallel()
	bytes := []byte(`"0xb3f98adc"`)
	var fid models.FunctionSelector
	err := json.Unmarshal(bytes, &fid)
	assert.NoError(t, err)
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_FunctionSelectorUnmarshalJSONLiteral(t *testing.T) {
	t.Parallel()
	literalSelectorBytes := []byte(`"setBytes(bytes)"`)
	var fid models.FunctionSelector
	err := json.Unmarshal(literalSelectorBytes, &fid)
	assert.NoError(t, err)
	assert.Equal(t, "0xda359dc8", fid.String())
}

func TestModels_FunctionSelectorUnmarshalJSONError(t *testing.T) {
	t.Parallel()
	bytes := []byte(`"0xb3f98adc123456"`)
	var fid models.FunctionSelector
	err := json.Unmarshal(bytes, &fid)
	assert.Error(t, err)
}

func TestSafeByteSlice_Success(t *testing.T) {
	tests := []struct {
		ary      models.UntrustedBytes
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
		ary   models.UntrustedBytes
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

	var block models.Block
	err := json.Unmarshal([]byte(input), &block)
	require.NoError(t, err)

	assert.Len(t, block.Transactions, 2)
}

func TestHead_EarliestInChain(t *testing.T) {
	head := models.Head{
		Number: 3,
		Parent: &models.Head{
			Number: 2,
			Parent: &models.Head{
				Number: 1,
			},
		},
	}

	assert.Equal(t, int64(1), head.EarliestInChain().Number)
}

func TestTxReceipt_ReceiptIndicatesRunLogFulfillment(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"basic", "testdata/getTransactionReceipt.json", false},
		{"runlog request", "testdata/runlogReceipt.json", false},
		{"runlog response", "testdata/responseReceipt.json", true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			receipt := cltest.TxReceiptFromFixture(t, test.path)
			assert.Equal(t, test.want, models.ReceiptIndicatesRunLogFulfillment(receipt))
		})
	}
}
