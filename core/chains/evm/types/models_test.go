package types_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/utils"
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
			num := evmtypes.NewHead(test.input, utils.NewHash(), utils.NewHash(), 0, nil)
			assert.Equal(t, test.want, fmt.Sprintf("%x", num.ToInt()))
		})
	}
}

func TestHead_GreaterThan(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		left    *evmtypes.Head
		right   *evmtypes.Head
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
		bn   *evmtypes.Head
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

func TestEthTx_GetID(t *testing.T) {
	tx := txmgr.EthTx{ID: math.MinInt64}
	assert.Equal(t, "-9223372036854775808", tx.GetID())
}

func TestEthTxAttempt_GetSignedTx(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)
	tx := gethTypes.NewTransaction(uint64(42), testutils.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})

	chainID := big.NewInt(3)

	signedTx, err := ethKeyStore.SignTx(fromAddress, tx, chainID)
	require.NoError(t, err)
	signedTx.Size() // Needed to write the size for equality checking
	rlp := new(bytes.Buffer)
	require.NoError(t, signedTx.EncodeRLP(rlp))

	attempt := txmgr.EthTxAttempt{SignedRawTx: rlp.Bytes()}

	gotSignedTx, err := attempt.GetSignedTx()
	require.NoError(t, err)
	decodedEncoded := new(bytes.Buffer)
	require.NoError(t, gotSignedTx.EncodeRLP(decodedEncoded))

	require.Equal(t, signedTx.Hash(), gotSignedTx.Hash())
	require.Equal(t, attempt.SignedRawTx, decodedEncoded.Bytes())
}

func TestHead_ChainLength(t *testing.T) {
	head := evmtypes.Head{
		Parent: &evmtypes.Head{
			Parent: &evmtypes.Head{},
		},
	}

	assert.Equal(t, uint32(3), head.ChainLength())

	var head2 *evmtypes.Head
	assert.Equal(t, uint32(0), head2.ChainLength())
}

func TestModels_HexToFunctionSelector(t *testing.T) {
	t.Parallel()
	fid := evmtypes.HexToFunctionSelector("0xb3f98adc")
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_HexToFunctionSelectorOverflow(t *testing.T) {
	t.Parallel()
	fid := evmtypes.HexToFunctionSelector("0xb3f98adc123456")
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_FunctionSelectorUnmarshalJSON(t *testing.T) {
	t.Parallel()
	bytes := []byte(`"0xb3f98adc"`)
	var fid evmtypes.FunctionSelector
	err := json.Unmarshal(bytes, &fid)
	assert.NoError(t, err)
	assert.Equal(t, "0xb3f98adc", fid.String())
}

func TestModels_FunctionSelectorUnmarshalJSONLiteral(t *testing.T) {
	t.Parallel()
	literalSelectorBytes := []byte(`"setBytes(bytes)"`)
	var fid evmtypes.FunctionSelector
	err := json.Unmarshal(literalSelectorBytes, &fid)
	assert.NoError(t, err)
	assert.Equal(t, "0xda359dc8", fid.String())
}

func TestModels_FunctionSelectorUnmarshalJSONError(t *testing.T) {
	t.Parallel()
	bytes := []byte(`"0xb3f98adc123456"`)
	var fid evmtypes.FunctionSelector
	err := json.Unmarshal(bytes, &fid)
	assert.Error(t, err)
}

func TestSafeByteSlice_Success(t *testing.T) {
	tests := []struct {
		ary      evmtypes.UntrustedBytes
		start    int
		end      int
		expected []byte
	}{
		{[]byte{1, 2, 3}, 0, 0, []byte{}},
		{[]byte{1, 2, 3}, 0, 1, []byte{1}},
		{[]byte{1, 2, 3}, 1, 3, []byte{2, 3}},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual, err := test.ary.SafeByteSlice(test.start, test.end)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestSafeByteSlice_Error(t *testing.T) {
	tests := []struct {
		ary   evmtypes.UntrustedBytes
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
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual, err := test.ary.SafeByteSlice(test.start, test.end)
			assert.EqualError(t, err, "out of bounds slice access")
			var expected []byte
			assert.Equal(t, expected, actual)
		})
	}
}

func TestHead_EarliestInChain(t *testing.T) {
	head := evmtypes.Head{
		Number: 3,
		Parent: &evmtypes.Head{
			Number: 2,
			Parent: &evmtypes.Head{
				Number: 1,
			},
		},
	}

	assert.Equal(t, int64(1), head.EarliestInChain().Number)
}

func TestHead_IsInChain(t *testing.T) {
	hash1 := utils.NewHash()
	hash2 := utils.NewHash()
	hash3 := utils.NewHash()

	head := evmtypes.Head{
		Number: 3,
		Hash:   hash3,
		Parent: &evmtypes.Head{
			Hash:   hash2,
			Number: 2,
			Parent: &evmtypes.Head{
				Hash:   hash1,
				Number: 1,
			},
		},
	}

	assert.True(t, head.IsInChain(hash1))
	assert.True(t, head.IsInChain(hash2))
	assert.True(t, head.IsInChain(hash3))
	assert.False(t, head.IsInChain(utils.NewHash()))
	assert.False(t, head.IsInChain(common.Hash{}))
}

func TestTxReceipt_ReceiptIndicatesRunLogFulfillment(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"basic", "../../../testdata/jsonrpc/getTransactionReceipt.json", false},
		{"runlog request", "../../../testdata/jsonrpc/runlogReceipt.json", false},
		{"runlog response", "../../../testdata/jsonrpc/responseReceipt.json", true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			receipt := cltest.TxReceiptFromFixture(t, test.path)
			require.Equal(t, test.want, evmtypes.ReceiptIndicatesRunLogFulfillment(*receipt))
		})
	}
}

func TestHead_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected evmtypes.Head
	}{
		{"geth",
			`{"difficulty":"0xf3a00","extraData":"0xd883010503846765746887676f312e372e318664617277696e","gasLimit":"0xffc001","gasUsed":"0x0","hash":"0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0xd1aeb42885a43b72b518182ef893125814811048","mixHash":"0x0f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","nonce":"0x0ece08ea8c49dfd9","number":"0x100","parentHash":"0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x218","stateRoot":"0xc7b01007a10da045eacb90385887dd0c38fcb5db7393006bdde24b93873c334b","timestamp":"0x58318da2","totalDifficulty":"0x1f3a00","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]}`,
			evmtypes.Head{
				Hash:       common.HexToHash("0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a"),
				Number:     0x100,
				ParentHash: common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d"),
				Timestamp:  time.Unix(0x58318da2, 0).UTC(),
			},
		},
		{"parity",
			`{"author":"0xd1aeb42885a43b72b518182ef893125814811048","difficulty":"0xf3a00","extraData":"0xd883010503846765746887676f312e372e318664617277696e","gasLimit":"0xffc001","gasUsed":"0x0","hash":"0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0xd1aeb42885a43b72b518182ef893125814811048","mixHash":"0x0f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","nonce":"0x0ece08ea8c49dfd9","number":"0x100","parentHash":"0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sealFields":["0xa00f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","0x880ece08ea8c49dfd9"],"sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x218","stateRoot":"0xc7b01007a10da045eacb90385887dd0c38fcb5db7393006bdde24b93873c334b","timestamp":"0x58318da2","totalDifficulty":"0x1f3a00","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]}`,
			evmtypes.Head{
				Hash:       common.HexToHash("0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a"),
				Number:     0x100,
				ParentHash: common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d"),
				Timestamp:  time.Unix(0x58318da2, 0).UTC(),
			},
		},
		{"arbitrum",
			`{"number":"0x15156","hash":"0x752dab43f7a2482db39227d46cd307623b26167841e2207e93e7566ab7ab7871","parentHash":"0x923ad1e27c1d43cb2d2fb09e26d2502ca4b4914a2e0599161d279c6c06117d34","mixHash":"0x0000000000000000000000000000000000000000000000000000000000000000","nonce":"0x0000000000000000","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","transactionsRoot":"0x71448077f5ce420a8e24db62d4d58e8d8e6ad2c7e76318868e089d41f7e0faf3","stateRoot":"0x0000000000000000000000000000000000000000000000000000000000000000","receiptsRoot":"0x2c292672b8fc9d223647a2569e19721f0757c96a1421753a93e141f8e56cf504","miner":"0x0000000000000000000000000000000000000000","difficulty":"0x0","totalDifficulty":"0x0","extraData":"0x","size":"0x0","gasLimit":"0x11278208","gasUsed":"0x3d1fe9","timestamp":"0x60d0952d","transactions":["0xa1ea93556b93ed3b45cb24f21c8deb584e6a9049c35209242651bf3533c23b98","0xfc6593c45ba92351d17173aa1381e84734d252ab0169887783039212c4a41024","0x85ee9d04fd0ebb5f62191eeb53cb45d9c0945d43eba444c3548de2ac8421682f","0x50d120936473e5b75f6e04829ad4eeca7a1df7d3c5026ebb5d34af936a39b29c"],"uncles":[],"l1BlockNumber":"0x8652f9"}`,
			evmtypes.Head{
				Hash:          common.HexToHash("0x752dab43f7a2482db39227d46cd307623b26167841e2207e93e7566ab7ab7871"),
				Number:        0x15156,
				ParentHash:    common.HexToHash("0x923ad1e27c1d43cb2d2fb09e26d2502ca4b4914a2e0599161d279c6c06117d34"),
				Timestamp:     time.Unix(0x60d0952d, 0).UTC(),
				L1BlockNumber: null.Int64From(0x8652f9),
			},
		},
		{"not found",
			`null`,
			evmtypes.Head{},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			var head evmtypes.Head
			err := head.UnmarshalJSON([]byte(test.json))
			require.NoError(t, err)
			assert.Equal(t, test.expected.Hash, head.Hash)
			assert.Equal(t, test.expected.Number, head.Number)
			assert.Equal(t, test.expected.ParentHash, head.ParentHash)
			assert.Equal(t, test.expected.Timestamp.UTC().Unix(), head.Timestamp.UTC().Unix())
			assert.Equal(t, test.expected.L1BlockNumber, head.L1BlockNumber)
		})
	}
}

func TestHead_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		head     evmtypes.Head
		expected string
	}{
		{"happy",
			evmtypes.Head{
				Hash:       common.HexToHash("0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a"),
				Number:     0x100,
				ParentHash: common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d"),
				Timestamp:  time.Unix(0x58318da2, 0).UTC(),
			},
			`{"hash":"0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a","number":"0x100","parentHash":"0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d","timestamp":"0x58318da2"}`,
		},
		{"empty",
			evmtypes.Head{},
			`{"number":"0x0"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			bs, err := test.head.MarshalJSON()
			require.NoError(t, err)
			require.Equal(t, test.expected, string(bs))
		})
	}
}

func Test_NullableEIP2930AccessList(t *testing.T) {
	addr := testutils.NewAddress()
	storageKey := utils.NewHash()
	al := gethTypes.AccessList{{Address: addr, StorageKeys: []common.Hash{storageKey}}}
	alb, err := json.Marshal(al)
	require.NoError(t, err)
	jsonStr := fmt.Sprintf(`[{"address":"0x%s","storageKeys":["%s"]}]`, hex.EncodeToString(addr.Bytes()), storageKey.Hex())
	require.Equal(t, jsonStr, string(alb))

	nNull := txmgr.NullableEIP2930AccessList{}
	nValid := txmgr.NullableEIP2930AccessListFrom(al)

	t.Run("MarshalJSON", func(t *testing.T) {
		_, err := json.Marshal(nNull)
		require.NoError(t, err)
		assert.Nil(t, nil)

		b, err := json.Marshal(nValid)
		require.NoError(t, err)
		assert.Equal(t, alb, b)
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		var n txmgr.NullableEIP2930AccessList
		err := json.Unmarshal(nil, &n)
		require.EqualError(t, err, "unexpected end of JSON input")

		err = json.Unmarshal([]byte("null"), &n)
		require.NoError(t, err)
		assert.False(t, n.Valid)

		err = json.Unmarshal([]byte(jsonStr), &n)
		require.NoError(t, err)
		assert.True(t, n.Valid)
		assert.Equal(t, al, n.AccessList)
	})

	t.Run("Value", func(t *testing.T) {
		value, err := nNull.Value()
		require.NoError(t, err)
		assert.Nil(t, value)

		value, err = nValid.Value()
		require.NoError(t, err)
		assert.NotNil(t, value)
		assert.Equal(t, alb, value)
	})

	t.Run("Scan", func(t *testing.T) {
		n := new(txmgr.NullableEIP2930AccessList)
		err := n.Scan(nil)
		require.NoError(t, err)
		assert.False(t, n.Valid)

		err = n.Scan([]byte("null"))
		require.NoError(t, err)
		assert.False(t, n.Valid)

		err = n.Scan([]byte(jsonStr))
		require.NoError(t, err)
		assert.True(t, n.Valid)
		assert.Equal(t, al, n.AccessList)
	})
}
