package models_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/common"
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

func TestEthTx_GetID(t *testing.T) {
	tx := models.EthTx{ID: math.MinInt64}
	assert.Equal(t, "-9223372036854775808", tx.GetID())
}

func TestEthTxAttempt_GetSignedTx(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth
	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)
	ethKeyStore.Unlock(cltest.Password)
	tx := gethTypes.NewTransaction(uint64(42), cltest.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})

	chainID := big.NewInt(3)

	signedTx, err := ethKeyStore.SignTx(fromAddress, tx, chainID)
	require.NoError(t, err)
	signedTx.Size() // Needed to write the size for equality checking
	rlp := new(bytes.Buffer)
	require.NoError(t, signedTx.EncodeRLP(rlp))

	attempt := models.EthTxAttempt{SignedRawTx: rlp.Bytes()}

	gotSignedTx, err := attempt.GetSignedTx()
	require.NoError(t, err)
	decodedEncoded := new(bytes.Buffer)
	require.NoError(t, gotSignedTx.EncodeRLP(decodedEncoded))

	require.Equal(t, signedTx.Hash(), gotSignedTx.Hash())
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
		t.Run(strconv.Itoa(i), func(t *testing.T) {
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
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual, err := test.ary.SafeByteSlice(test.start, test.end)
			assert.EqualError(t, err, "out of bounds slice access")
			var expected []byte
			assert.Equal(t, expected, actual)
		})
	}
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

func TestHead_IsInChain(t *testing.T) {
	hash1 := cltest.NewHash()
	hash2 := cltest.NewHash()
	hash3 := cltest.NewHash()

	head := models.Head{
		Number: 3,
		Hash:   hash3,
		Parent: &models.Head{
			Hash:   hash2,
			Number: 2,
			Parent: &models.Head{
				Hash:   hash1,
				Number: 1,
			},
		},
	}

	assert.True(t, head.IsInChain(hash1))
	assert.True(t, head.IsInChain(hash2))
	assert.True(t, head.IsInChain(hash3))
	assert.False(t, head.IsInChain(cltest.NewHash()))
	assert.False(t, head.IsInChain(common.Hash{}))
}

func TestTxReceipt_ReceiptIndicatesRunLogFulfillment(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"basic", "../../testdata/jsonrpc/getTransactionReceipt.json", false},
		{"runlog request", "../../testdata/jsonrpc/runlogReceipt.json", false},
		{"runlog response", "../../testdata/jsonrpc/responseReceipt.json", true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			receipt := cltest.TxReceiptFromFixture(t, test.path)
			require.Equal(t, test.want, models.ReceiptIndicatesRunLogFulfillment(*receipt))
		})
	}
}

func TestHead_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected models.Head
	}{
		{"geth",
			`{"difficulty":"0xf3a00","extraData":"0xd883010503846765746887676f312e372e318664617277696e","gasLimit":"0xffc001","gasUsed":"0x0","hash":"0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0xd1aeb42885a43b72b518182ef893125814811048","mixHash":"0x0f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","nonce":"0x0ece08ea8c49dfd9","number":"0x100","parentHash":"0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x218","stateRoot":"0xc7b01007a10da045eacb90385887dd0c38fcb5db7393006bdde24b93873c334b","timestamp":"0x58318da2","totalDifficulty":"0x1f3a00","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]}`,
			models.Head{
				Hash:       common.HexToHash("0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a"),
				Number:     0x100,
				ParentHash: common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d"),
				Timestamp:  time.Unix(0x58318da2, 0).UTC(),
			},
		},
		{"parity",
			`{"author":"0xd1aeb42885a43b72b518182ef893125814811048","difficulty":"0xf3a00","extraData":"0xd883010503846765746887676f312e372e318664617277696e","gasLimit":"0xffc001","gasUsed":"0x0","hash":"0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0xd1aeb42885a43b72b518182ef893125814811048","mixHash":"0x0f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","nonce":"0x0ece08ea8c49dfd9","number":"0x100","parentHash":"0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sealFields":["0xa00f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","0x880ece08ea8c49dfd9"],"sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x218","stateRoot":"0xc7b01007a10da045eacb90385887dd0c38fcb5db7393006bdde24b93873c334b","timestamp":"0x58318da2","totalDifficulty":"0x1f3a00","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]}`,
			models.Head{
				Hash:       common.HexToHash("0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a"),
				Number:     0x100,
				ParentHash: common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d"),
				Timestamp:  time.Unix(0x58318da2, 0).UTC(),
			},
		},
		{"not found",
			`null`,
			models.Head{},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			var head models.Head
			err := head.UnmarshalJSON([]byte(test.json))
			require.NoError(t, err)
			require.Equal(t, test.expected.Hash, head.Hash)
			require.Equal(t, test.expected.Number, head.Number)
			require.Equal(t, test.expected.ParentHash, head.ParentHash)
			require.Equal(t, test.expected.Timestamp.UTC().Unix(), head.Timestamp.UTC().Unix())
		})
	}
}

func TestHead_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		head     models.Head
		expected string
	}{
		{"happy",
			models.Head{
				Hash:       common.HexToHash("0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a"),
				Number:     0x100,
				ParentHash: common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d"),
				Timestamp:  time.Unix(0x58318da2, 0).UTC(),
			},
			`{"hash":"0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a","number":"0x100","parentHash":"0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d","timestamp":"0x58318da2"}`,
		},
		{"empty",
			models.Head{},
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
