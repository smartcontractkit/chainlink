package store_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	strpkg "github.com/smartcontractkit/chainlink-go/store"
	"github.com/smartcontractkit/chainlink-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestEthCreateTx(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	assert.Nil(t, store.KeyStore.Unlock(cltest.Password))
	manager := store.Eth

	to := "0xb70a511baC46ec6442aC6D598eaC327334e634dB"
	data := "0000abcdef"
	txid := "0x86300ee06a57eb27fbd8a6d5380783d4f8cb7210747689fe452e40f049d3de08"
	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", "0x0100") // 256
	ethMock.Register("eth_sendRawTransaction", txid)

	tx, err := manager.CreateTx(to, data)
	assert.Nil(t, err)
	assert.Equal(t, uint64(256), tx.Nonce)
	assert.Equal(t, data, tx.Data)
	assert.Equal(t, to, tx.To)

	assert.Nil(t, store.One("From", tx.From, tx))
	assert.Equal(t, uint64(256), tx.Nonce)
	assert.Equal(t, 1, len(tx.Attempts))

	assert.True(t, ethMock.AllCalled())
}

func TestEthConfirmTxUnconfirmed(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	manager := store.Eth

	txid := "0x86300ee06a57eb27fbd8a6d5380783d4f8cb7210747689fe452e40f049d3de08"
	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})

	confirmed, err := manager.TxConfirmed(txid)
	assert.Nil(t, err)
	assert.False(t, confirmed)

	assert.True(t, ethMock.AllCalled())
}

func TestEthConfirmTxNotEnoughConfs(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	config := store.Config
	manager := store.Eth

	txid := "0x86300ee06a57eb27fbd8a6d5380783d4f8cb7210747689fe452e40f049d3de08"
	bNum := uint64(17)
	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{
		TXID:        txid,
		BlockNumber: bNum,
	})
	current := utils.Uint64ToHex(bNum + config.EthConfMin - 1)
	ethMock.Register("eth_blockNumber", current)

	confirmed, err := manager.TxConfirmed(txid)
	assert.Nil(t, err)
	assert.False(t, confirmed)

	assert.True(t, ethMock.AllCalled())
}

func TestEthConfirmTxTrue(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	config := store.Config
	manager := store.Eth

	txid := "0x86300ee06a57eb27fbd8a6d5380783d4f8cb7210747689fe452e40f049d3de08"
	bNum := uint64(17)
	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{
		TXID:        txid,
		BlockNumber: bNum,
	})
	current := utils.Uint64ToHex(bNum + config.EthConfMin)
	ethMock.Register("eth_blockNumber", current)

	confirmed, err := manager.TxConfirmed(txid)
	assert.Nil(t, err)
	assert.True(t, confirmed)

	assert.True(t, ethMock.AllCalled())
}
