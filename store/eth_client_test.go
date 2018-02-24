package store_test

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestEthClient_GetTxReceipt(t *testing.T) {
	response := cltest.LoadJSON("../internal/fixtures/eth/getTransactionReceipt.json")
	mockServer := cltest.NewWSServer(string(response))
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
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	ethMock := app.MockEthClient()
	ethClientObject := app.Store.TxManager.EthClient
	ethMock.Register("eth_sendRawTransaction", common.Hash{1})
	result, err := ethClientObject.SendRawTx("test")
	assert.Nil(t, err)
	assert.Equal(t, result, common.Hash{1})
}

func TestEthGetBalance(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	ethMock := app.MockEthClient()
	ethClientObject := app.Store.TxManager.EthClient

	ethMock.Register("eth_getBalance", "0x0100") // 256
	result, err := ethClientObject.GetEthBalance(cltest.NewAddress())
	assert.Nil(t, err)
	expected := 256e-18
	assert.Nil(t, err)
	assert.Equal(t, expected, result)

	ethMock.Register("eth_getBalance", "0x4b3b4ca85a86c4000000000000000000") // 1e38
	result, err = ethClientObject.GetEthBalance(cltest.NewAddress())
	expected = 1e20
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestBlockHeader_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	var bh models.BlockHeader

	data := cltest.LoadJSON("../internal/fixtures/eth/subscription_new_heads.json")
	value := gjson.Get(string(data), "params.result")
	assert.Nil(t, json.Unmarshal([]byte(value.String()), &bh))

	assert.Equal(t, cltest.BigHexInt(uint64(1263817)), bh.Number)
}
