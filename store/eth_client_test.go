package store_test

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestGethGetTransactionReceipt(t *testing.T) {
	response := cltest.LoadJSON("../internal/fixtures/eth/getTransactionReceipt.json")
	mockServer := cltest.NewWSServer(string(response))
	config := cltest.NewConfigWithWSServer(mockServer)
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	ec := store.TxManager.EthClient

	hash, _ := utils.StringToHash("0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238")
	receipt, err := ec.GetTxReceipt(hash)
	assert.Nil(t, err)
	assert.Equal(t, hash, receipt.Hash)
	assert.Equal(t, uint64(11), receipt.BlockNumber)
}

func TestEventLogUnmarshalJSON(t *testing.T) {
	notification := cltest.LoadJSON("../internal/fixtures/eth/subscription_logs.json")
	blockHash, _ := utils.StringToHash("0x61cdb2a09ab99abf791d474f20c2ea89bf8de2923a2d42bb49944c8c993cbf04")
	address, _ := utils.StringToAddress("0x8320fe7702b96808f7bbc0d4a888ed1468216cfd")
	txHash, _ := utils.StringToHash("0xe044554a0a55067caafd07f8020ab9f2af60bdfe337e395ecd84b4877a3d1ab4")
	topic1, _ := utils.StringToBytes("0xd78a0cb8bb633d06981248b816e7bd33c2a35a6089241d099fa519e361cab902")
	data, _ := utils.StringToBytes("0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000003")
	blockNumber := hexutil.Uint64(171655)
	logIndex := hexutil.Uint64(66)
	txIndex := hexutil.Uint64(23)

	en := store.EthNotification{}
	assert.Nil(t, json.Unmarshal(notification, &en))

	el, err := en.UnmarshalLog()
	assert.Nil(t, err)

	assert.Equal(t, blockHash, el.BlockHash)
	assert.Equal(t, address, el.Address)
	assert.Equal(t, txHash, el.TxHash)
	assert.Equal(t, 1, len(el.Topics))
	assert.Equal(t, topic1, el.Topics[0])
	assert.Equal(t, data, el.Data)
	assert.Equal(t, blockNumber, el.BlockNumber)
	assert.Equal(t, logIndex, el.LogIndex)
	assert.Equal(t, txIndex, el.TxIndex)
}

func TestEventLogUnmarshalJSONError(t *testing.T) {
	notification := cltest.LoadJSON("../internal/fixtures/eth/subscription_new_heads.json")

	en := store.EthNotification{}
	assert.Nil(t, json.Unmarshal(notification, &en))

	_, err := en.UnmarshalLog()
	assert.NotNil(t, err)
}
