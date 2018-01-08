package store_test

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestEthGetTxReceipt(t *testing.T) {
	store := cltest.NewStore()
	defer cltest.CleanUpStore(store)
	eth := store.Eth

	response := cltest.LoadJSON("../internal/fixtures/web/eth_getTransactionReceipt.json")
	gock.New(store.Config.EthereumURL).
		Post("").
		Reply(200).
		JSON(response)

	hash, _ := utils.StringToHash("0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238")
	receipt, err := eth.GetTxReceipt(hash)
	assert.Nil(t, err)
	assert.Equal(t, hash, receipt.Hash)
	assert.Equal(t, uint64(11), receipt.BlockNumber)
}
