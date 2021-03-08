package offchainreporting_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/require"
)

func Test_Transmitter_CreateEthTransaction(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	db, _ := store.DB.DB()

	key := cltest.MustInsertRandomKey(t, store.DB, 0)

	gasLimit := uint64(1000)
	fromAddress := key.Address.Address()
	toAddress := cltest.NewAddress()
	payload := []byte{1, 2, 3}

	transmitter := offchainreporting.NewTransmitter(db, fromAddress, gasLimit, 0)

	require.NoError(t, transmitter.CreateEthTransaction(context.Background(), toAddress, payload))

	etx := models.EthTx{}
	require.NoError(t, store.ORM.DB.First(&etx).Error)

	require.Equal(t, gasLimit, etx.GasLimit)
	require.Equal(t, fromAddress, etx.FromAddress)
	require.Equal(t, toAddress, etx.ToAddress)
	require.Equal(t, payload, etx.EncodedPayload)
	require.Equal(t, assets.NewEthValue(0), etx.Value)
}

func Test_Transmitter_CreateEthTransaction_OutOfEth(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	db, _ := store.DB.DB()

	thisKey := cltest.MustInsertRandomKey(t, store.DB, 1)
	otherKey := cltest.MustInsertRandomKey(t, store.DB, 1)

	gasLimit := uint64(1000)
	toAddress := cltest.NewAddress()

	transmitter := offchainreporting.NewTransmitter(db, thisKey.Address.Address(), gasLimit, 0)

	t.Run("if another key has any transactions with insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, store, 0, otherKey.Address.Address())

		require.NoError(t, transmitter.CreateEthTransaction(context.Background(), toAddress, payload))

		etx := models.EthTx{}
		require.NoError(t, store.ORM.DB.First(&etx, "nonce IS NULL AND from_address = ?", thisKey.Address.Address()).Error)
		require.Equal(t, payload, etx.EncodedPayload)
	})

	require.NoError(t, store.DB.Exec(`DELETE FROM eth_txes WHERE from_address = ?`, thisKey.Address.Address()).Error)

	t.Run("if this key has any transactions with insufficient eth errors, skips transmission entirely", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, store, 0, thisKey.Address.Address())

		err := transmitter.CreateEthTransaction(context.Background(), toAddress, payload)
		require.EqualError(t, err, fmt.Sprintf("Skipped OCR transmission because wallet is out of eth: %s", thisKey.Address.Hex()))
	})

	t.Run("if this key has transactions but no insufficient eth errors, transmits as normal", func(t *testing.T) {
		payload := cltest.MustRandomBytes(t, 100)
		require.NoError(t, store.DB.Exec(`UPDATE eth_tx_attempts SET state = 'broadcast'`).Error)
		require.NoError(t, store.DB.Exec(`UPDATE eth_txes SET nonce = 0, state = 'confirmed', broadcast_at = NOW()`).Error)

		require.NoError(t, transmitter.CreateEthTransaction(context.Background(), toAddress, payload))

		etx := models.EthTx{}
		require.NoError(t, store.ORM.DB.First(&etx, "nonce IS NULL AND from_address = ?", thisKey.Address.Address()).Error)
		require.Equal(t, payload, etx.EncodedPayload)
	})
}
