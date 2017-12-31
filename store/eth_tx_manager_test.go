package store_test

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestTxManagerCreateTx(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	assert.Nil(t, store.KeyStore.Unlock("password"))
	config := store.Config
	manager := store.Tx

	to := "0xb70a511baC46ec6442aC6D598eaC327334e634dB"
	data := "0000abcdef"
	txid := "0x86300ee06a57eb27fbd8a6d5380783d4f8cb7210747689fe452e40f049d3de08"
	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", "0x0100") // 256
	ethMock.Register("eth_sendRawTransaction", txid)

	tx, err := manager.CreateTx(to, data)
	assert.Nil(t, err)
	assert.Equal(t, uint64(256), tx.Nonce())
	assert.Equal(t, common.FromHex(data), tx.Data())
	assert.Equal(t, to, tx.To().Hex())

	signer := store.KeyStore.GetAccount().Address.String()
	rlp := bytes.NewBuffer([]byte{})
	assert.Nil(t, tx.EncodeRLP(rlp))
	rlpHex := common.ToHex(rlp.Bytes())
	sender, err := utils.SenderFromTxHex(rlpHex, config.ChainID)
	assert.Equal(t, signer, sender.Hex())
	assert.True(t, ethMock.AllCalled())
}

func TestTxManagerNewSignedTx(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	assert.Nil(t, store.KeyStore.Unlock("password"))
	config := store.Config
	manager := store.Tx

	data := "0000abcdef"
	to := "0xb70a511baC46ec6442aC6D598eaC327334e634dB"
	eth := app.MockEthClient()
	eth.Register("eth_getTransactionCount", "0x0100") // 256

	tx, err := manager.NewSignedTx(to, data)
	assert.Nil(t, err)
	assert.Equal(t, uint64(256), tx.Nonce())
	assert.Equal(t, common.FromHex(data), tx.Data())
	assert.Equal(t, to, tx.To().Hex())

	signer := store.KeyStore.GetAccount().Address.String()
	rlp := bytes.NewBuffer([]byte{})
	assert.Nil(t, tx.EncodeRLP(rlp))
	rlpHex := common.ToHex(rlp.Bytes())
	sender, err := utils.SenderFromTxHex(rlpHex, config.ChainID)
	assert.Equal(t, signer, sender.Hex())
}

func TestTxManagerConfirmTxTrue(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	manager := store.Tx

	txid := "0x86300ee06a57eb27fbd8a6d5380783d4f8cb7210747689fe452e40f049d3de08"
	eth := app.MockEthClient()
	eth.Register("eth_getTransactionReceipt", types.Receipt{TxHash: common.StringToHash(txid)})

	confirmed, err := manager.TxConfirmed(txid)
	assert.Nil(t, err)
	assert.True(t, confirmed)
}

func TestTxManagerConfirmTxFalse(t *testing.T) {
	t.Parallel()
	app := cltest.NewApplicationWithKeyStore()
	store := app.Store
	defer app.Stop()
	manager := store.Tx

	txid := "0x86300ee06a57eb27fbd8a6d5380783d4f8cb7210747689fe452e40f049d3de08"
	eth := app.MockEthClient()
	eth.Register("eth_getTransactionReceipt", types.Receipt{})

	confirmed, err := manager.TxConfirmed(txid)
	assert.Nil(t, err)
	assert.False(t, confirmed)
}
