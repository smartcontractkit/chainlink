package store_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestGethGetTransactionReceipt(t *testing.T) {
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
	assert.Equal(t, uint64(11), receipt.BlockNumber)
}

//func TestEventLogUnmarshalJSONError(t *testing.T) {
//notification := cltest.LoadJSON("../internal/fixtures/eth/subscription_new_heads.json")

//en := store.EthNotification{}
//assert.Nil(t, json.Unmarshal(notification, &en))

//_, err := en.UnmarshalLog()
//assert.NotNil(t, err)
//}
