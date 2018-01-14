package store_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
)

func TestEthGetTxReceipt(t *testing.T) {
	t.Parallel()

	config := cltest.NewConfig()
	response := cltest.LoadJSON("../internal/fixtures/web/eth_getTransactionReceipt.json")
	mockServer := cltest.NewWSServer(string(response))
	config.SetEthereumServer(mockServer)
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	eth := store.TxManager

	hash, _ := utils.StringToHash("0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238")
	receipt, err := eth.GetTxReceipt(hash)
	assert.Nil(t, err)
	assert.Equal(t, hash, receipt.Hash)
	assert.Equal(t, uint64(11), receipt.BlockNumber)
}
