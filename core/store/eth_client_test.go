package store_test

import (
	"encoding/json"
	"testing"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEthClient_GetTxReceipt(t *testing.T) {
	response := cltest.MustReadFile(t, "../internal/fixtures/eth/getTransactionReceipt.json")
	mockServer, wsCleanup := cltest.NewWSServer(string(response))
	defer wsCleanup()
	config := cltest.NewConfigWithWSServer(t, mockServer)
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	ec := store.TxManager.(*strpkg.EthTxManager).EthClient

	hash := common.HexToHash("0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238")
	receipt, err := ec.GetTxReceipt(hash)
	assert.NoError(t, err)
	assert.Equal(t, hash, receipt.Hash)
	assert.Equal(t, cltest.Int(uint64(11)), receipt.BlockNumber)
}

func TestTxReceipt_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantLogLen int
	}{
		{"basic", "testdata/getTransactionReceipt.json", 0},
		{"runlog request", "testdata/runlogReceipt.json", 4},
		{"runlog response", "testdata/responseReceipt.json", 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonStr := cltest.JSONFromFixture(t, test.path).Get("result").String()
			var receipt models.TxReceipt
			err := json.Unmarshal([]byte(jsonStr), &receipt)
			require.NoError(t, err)

			assert.Equal(t, test.wantLogLen, len(receipt.Logs))
		})
	}
}

func TestTxReceipt_FulfilledRunlog(t *testing.T) {
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
		t.Run(test.name, func(t *testing.T) {
			receipt := cltest.TxReceiptFromFixture(t, test.path)
			assert.Equal(t, test.want, receipt.FulfilledRunLog())
		})
	}
}

func TestEthClient_GetNonce(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	ethMock := app.MockEthClient()
	ethClientObject := app.Store.TxManager.(*strpkg.EthTxManager).EthClient
	ethMock.Register("eth_getTransactionCount", "0x0100")
	result, err := ethClientObject.GetNonce(cltest.NewAddress())
	assert.NoError(t, err)
	var expected uint64 = 256
	assert.Equal(t, result, expected)
}

func TestEthClient_GetBlockNumber(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	ethMock := app.MockEthClient()
	ethClientObject := app.Store.TxManager.(*strpkg.EthTxManager).EthClient
	ethMock.Register("eth_blockNumber", "0x0100")
	result, err := ethClientObject.GetBlockNumber()
	assert.NoError(t, err)
	var expected uint64 = 256
	assert.Equal(t, result, expected)
}

func TestEthClient_SendRawTx(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	ethMock := app.MockEthClient()
	ethClientObject := app.Store.TxManager.(*strpkg.EthTxManager).EthClient
	ethMock.Register("eth_sendRawTransaction", common.Hash{1})
	result, err := ethClientObject.SendRawTx("test")
	assert.NoError(t, err)
	assert.Equal(t, result, common.Hash{1})
}

func TestEthClient_GetEthBalance(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic", "0x0100", "0.000000000000000256"},
		{"larger than signed 64 bit integer", "0x4b3b4ca85a86c47a098a224000000000", "100000000000000000000.000000000000000000"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ethMock := app.MockEthClient()
			ethClientObject := app.Store.TxManager.(*strpkg.EthTxManager).EthClient

			ethMock.Register("eth_getBalance", test.input)
			result, err := ethClientObject.GetEthBalance(cltest.NewAddress())
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result.String())
		})
	}
}

func TestEthClient_GetERC20Balance(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	ethMock := app.MockEthClient()
	ethClientObject := app.Store.TxManager.(*strpkg.EthTxManager).EthClient

	ethMock.Register("eth_call", "0x0100") // 256
	result, err := ethClientObject.GetERC20Balance(cltest.NewAddress(), cltest.NewAddress())
	assert.NoError(t, err)
	expected := big.NewInt(256)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	ethMock.Register("eth_call", "0x4b3b4ca85a86c47a098a224000000000") // 1e38
	result, err = ethClientObject.GetERC20Balance(cltest.NewAddress(), cltest.NewAddress())
	expected = big.NewInt(0)
	expected.SetString("100000000000000000000000000000000000000", 10)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
