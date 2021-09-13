package orm_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestORM_CreateExternalInitiator(t *testing.T) {
	store := cltest.NewStore(t)

	token := auth.NewToken()
	req := models.ExternalInitiatorRequest{
		Name: "externalinitiator",
	}
	exi, err := models.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.NoError(t, store.CreateExternalInitiator(exi))

	exi2, err := models.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.Equal(t, `ERROR: duplicate key value violates unique constraint "external_initiators_name_key" (SQLSTATE 23505)`, store.CreateExternalInitiator(exi2).Error())
}

func TestORM_DeleteExternalInitiator(t *testing.T) {
	store := cltest.NewStore(t)

	token := auth.NewToken()
	req := models.ExternalInitiatorRequest{
		Name: "externalinitiator",
	}
	exi, err := models.NewExternalInitiator(token, &req)
	require.NoError(t, err)
	require.NoError(t, store.CreateExternalInitiator(exi))

	_, err = store.FindExternalInitiator(token)
	require.NoError(t, err)

	err = store.DeleteExternalInitiator(exi.Name)
	require.NoError(t, err)

	_, err = store.FindExternalInitiator(token)
	require.Error(t, err)

	require.NoError(t, store.CreateExternalInitiator(exi))
}

func TestORM_EthTransactionsWithAttempts(t *testing.T) {
	store := cltest.NewStore(t)
	db := store.DB
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	cltest.MustInsertConfirmedEthTxWithAttempt(t, db, 0, 1, from)        // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithAttempt(t, db, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewEthTxAttempt(t, tx2.ID)
	attempt.State = bulletprooftxmanager.EthTxAttemptBroadcast
	attempt.GasPrice = *utils.NewBig(big.NewInt(3))
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, store.DB.Create(&attempt).Error)

	// tx 3 has no attempts
	tx3 := cltest.NewEthTx(t, from)
	tx3.State = bulletprooftxmanager.EthTxUnstarted
	tx3.FromAddress = from
	require.NoError(t, store.DB.Save(&tx3).Error)

	count, err := store.CountOf(bulletprooftxmanager.EthTx{})
	require.NoError(t, err)
	require.Equal(t, 3, count)

	txs, count, err := store.EthTransactionsWithAttempts(0, 100) // should omit tx3
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 2)
	assert.Equal(t, int64(1), *txs[0].Nonce, "transactions should be sorted by nonce")
	assert.Equal(t, int64(0), *txs[1].Nonce, "transactions should be sorted by nonce")
	assert.Len(t, txs[0].EthTxAttempts, 2, "all eth tx attempts are preloaded")
	assert.Len(t, txs[1].EthTxAttempts, 1)
	assert.Equal(t, int64(3), *txs[0].EthTxAttempts[0].BroadcastBeforeBlockNum, "attempts shoud be sorted by created_at")
	assert.Equal(t, int64(2), *txs[0].EthTxAttempts[1].BroadcastBeforeBlockNum, "attempts shoud be sorted by created_at")

	txs, count, err = store.EthTransactionsWithAttempts(0, 1)
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 1, "limit should apply to length of results")
	assert.Equal(t, int64(1), *txs[0].Nonce, "transactions should be sorted by nonce")
}
