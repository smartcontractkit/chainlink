package adapters_test

import (
	"testing"

	gock "github.com/h2non/gock"
	"github.com/smartcontractkit/chainlink-go/adapters"
	"github.com/smartcontractkit/chainlink-go/internal/cltest"
	"github.com/smartcontractkit/chainlink-go/store/models"
	"github.com/smartcontractkit/chainlink-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestSendingEthereumTx(t *testing.T) {
	store := cltest.NewStore()
	defer store.Close()
	defer cltest.CloseGock(t)
	config := store.Config

	value := "0000abcdef"
	input := models.RunResultWithValue(value)

	response := `{"result": "0x0100"}`
	gock.New(config.EthereumURL).
		Post("").
		Reply(200).
		JSON(response)

	adapter := adapters.EthSendRawTx{}
	result := adapter.Perform(input, store)
	assert.Equal(t, "0x0100", result.Value())
}

func TestSigningEthereumTx(t *testing.T) {
	defer cltest.CloseGock(t)

	app := cltest.NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	config := app.Store.Config
	sender := store.KeyStore.GetAccount().Address.String()
	password := "password"

	response := `{"result": "0x11"}`
	gock.New(config.EthereumURL).
		Post("").
		Reply(200).
		JSON(response)

	err := store.KeyStore.Unlock(password)
	assert.Nil(t, err)

	data := "0000abcdef"
	recipient := "0xb70a511bac46ec6442ac6d598eac327334e634db"
	fid := "0x12345678"
	input := models.RunResultWithValue(data)

	adapter := adapters.EthSignTx{
		Address:     recipient,
		FunctionID:  fid,
	}
	result := adapter.Perform(input, store)
	assert.Contains(t, result.Value(), data)
	assert.Contains(t, result.Value(), recipient[2:len(recipient)])

	tx, err := utils.DecodeTxFromHex(result.Value(), config.ChainID)
	assert.Nil(t, err)
	assert.Equal(t, uint64(17), tx.Nonce())

	actual, err := utils.SenderFromTxHex(result.Value(), config.ChainID)
	assert.Equal(t, sender, actual.Hex())
}

func TestSigningAndSendingTx(t *testing.T) {
	defer cltest.CloseGock(t)

	app := cltest.NewApplicationWithKeyStore()
	defer app.Stop()
	store := app.Store
	eth := app.MockEthClient()
	eth.RegisterError("eth_getTransactionCount", "Cannot connect to nodes")

	adapter := adapters.EthSignTx{
		Address:     "recipient",
		FunctionID:  "fid",
	}
	input := models.RunResultWithValue("Hello World!")
	output := adapter.Perform(input, store)

	assert.True(t, output.HasError())
	assert.Equal(t, output.Error(), "Cannot connect to nodes")
}
