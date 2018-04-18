package store_test

import (
	"testing"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestEthClient_GetTxReceipt(t *testing.T) {
	t.Parallel()
	response := cltest.LoadJSON("../internal/fixtures/eth/getTransactionReceipt.json")
	mockServer := cltest.NewWSServer(string(response))
	defer mockServer.Close()
	config := cltest.NewConfigWithWSServer(mockServer)
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	ec := store.TxManager.EthClient

	hash := common.HexToHash("0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238")
	receipt, err := ec.GetTxReceipt(hash)
	assert.Nil(t, err)
	assert.Equal(t, hash, receipt.Hash)
	assert.Equal(t, cltest.BigHexInt(uint64(11)), receipt.BlockNumber)
}

func TestEthClient_GetNonce(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	ethMock := app.MockEthClient()
	ethClientObject := app.Store.TxManager.EthClient
	ethMock.Register("eth_getTransactionCount", "0x0100")
	result, err := ethClientObject.GetNonce(cltest.NewAddress())
	assert.Nil(t, err)
	var expected uint64 = 256
	assert.Equal(t, result, expected)
}

func TestEthClient_GetBlockNumber(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	ethMock := app.MockEthClient()
	ethClientObject := app.Store.TxManager.EthClient
	ethMock.Register("eth_blockNumber", "0x0100")
	result, err := ethClientObject.GetBlockNumber()
	assert.Nil(t, err)
	var expected uint64 = 256
	assert.Equal(t, result, expected)
}

func TestEthClient_SendRawTx(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	ethMock := app.MockEthClient()
	ethClientObject := app.Store.TxManager.EthClient
	ethMock.Register("eth_sendRawTransaction", common.Hash{1})
	result, err := ethClientObject.SendRawTx("test")
	assert.Nil(t, err)
	assert.Equal(t, result, common.Hash{1})
}

func bigRat(s string) *big.Rat {
	n, ok := new(big.Rat).SetString(s)
	if !ok {
		panic("big rational number could not be parsed")
	}
	return n
}

func TestEthClient_GetEthBalance(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	tests := []struct {
		name     string
		input    string
		expected *big.Rat
	}{
		{"basic", "0x0100", bigRat("256/1000000000000000000")},
		{"larger than signed 64 bit integer", "0x4b3b4ca85a86c47a098a224000000000", bigRat("100000000000000000000/1")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ethMock := app.MockEthClient()
			ethClientObject := app.Store.TxManager.EthClient

			ethMock.Register("eth_getBalance", test.input)
			result, err := ethClientObject.GetEthBalance(cltest.NewAddress())
			assert.Nil(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestEthClient_GetERC20Balance(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	ethMock := app.MockEthClient()
	ethClientObject := app.Store.TxManager.EthClient

	ethMock.Register("eth_call", "0x0100") // 256
	result, err := ethClientObject.GetERC20Balance(cltest.NewAddress(), cltest.NewAddress())
	assert.Nil(t, err)
	expected := big.NewInt(256)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)

	ethMock.Register("eth_call", "0x4b3b4ca85a86c47a098a224000000000") // 1e38
	result, err = ethClientObject.GetERC20Balance(cltest.NewAddress(), cltest.NewAddress())
	expected = big.NewInt(0)
	expected.SetString("100000000000000000000000000000000000000", 10)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
