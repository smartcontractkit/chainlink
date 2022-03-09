package txmgr_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestORM_EthTransactionsWithAttempts(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	orm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, orm, 0, 1, from)        // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, orm, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewLegacyEthTxAttempt(t, tx2.ID)
	attempt.State = txmgr.EthTxAttemptBroadcast
	attempt.GasPrice = utils.NewBig(big.NewInt(3))
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, orm.InsertEthTxAttempt(&attempt))

	// tx 3 has no attempts
	tx3 := cltest.NewEthTx(t, from)
	tx3.State = txmgr.EthTxUnstarted
	tx3.FromAddress = from
	require.NoError(t, orm.InsertEthTx(&tx3))

	var count int
	err := db.Get(&count, `SELECT count(*) FROM eth_txes`)
	require.NoError(t, err)
	require.Equal(t, 3, count)

	txs, count, err := orm.EthTransactionsWithAttempts(0, 100) // should omit tx3
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 2)
	assert.Equal(t, int64(1), *txs[0].Nonce, "transactions should be sorted by nonce")
	assert.Equal(t, int64(0), *txs[1].Nonce, "transactions should be sorted by nonce")
	assert.Len(t, txs[0].EthTxAttempts, 2, "all eth tx attempts are preloaded")
	assert.Len(t, txs[1].EthTxAttempts, 1)
	assert.Equal(t, int64(3), *txs[0].EthTxAttempts[0].BroadcastBeforeBlockNum, "attempts should be sorted by created_at")
	assert.Equal(t, int64(2), *txs[0].EthTxAttempts[1].BroadcastBeforeBlockNum, "attempts should be sorted by created_at")

	txs, count, err = orm.EthTransactionsWithAttempts(0, 1)
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 1, "limit should apply to length of results")
	assert.Equal(t, int64(1), *txs[0].Nonce, "transactions should be sorted by nonce")
}

func TestORM_EthTransactions(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	orm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, orm, 0, 1, from)        // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, orm, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewLegacyEthTxAttempt(t, tx2.ID)
	attempt.State = txmgr.EthTxAttemptBroadcast
	attempt.GasPrice = utils.NewBig(big.NewInt(3))
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, orm.InsertEthTxAttempt(&attempt))

	// tx 3 has no attempts
	tx3 := cltest.NewEthTx(t, from)
	tx3.State = txmgr.EthTxUnstarted
	tx3.FromAddress = from
	require.NoError(t, orm.InsertEthTx(&tx3))

	var count int
	err := db.Get(&count, `SELECT count(*) FROM eth_txes`)
	require.NoError(t, err)
	require.Equal(t, 3, count)

	txs, count, err := orm.EthTransactions(0, 100)
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 2)
	assert.Equal(t, int64(1), *txs[0].Nonce, "transactions should be sorted by nonce")
	assert.Equal(t, int64(0), *txs[1].Nonce, "transactions should be sorted by nonce")
	assert.Len(t, txs[0].EthTxAttempts, 0, "eth tx attempts should not be preloaded")
	assert.Len(t, txs[1].EthTxAttempts, 0)
}

func TestORM(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	orm := cltest.NewTxmORM(t, db, cfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth(), 0)

	var err error
	var etx txmgr.EthTx
	t.Run("InsertEthTx", func(t *testing.T) {
		etx = cltest.NewEthTx(t, fromAddress)
		err = orm.InsertEthTx(&etx)
		require.NoError(t, err)
		assert.Greater(t, int(etx.ID), 0)
		cltest.AssertCount(t, db, "eth_txes", 1)
	})
	var attemptL txmgr.EthTxAttempt
	var attemptD txmgr.EthTxAttempt
	t.Run("InsertEthTxAttempt", func(t *testing.T) {
		attemptD = cltest.NewDynamicFeeEthTxAttempt(t, etx.ID)
		err = orm.InsertEthTxAttempt(&attemptD)
		require.NoError(t, err)
		assert.Greater(t, int(attemptD.ID), 0)
		cltest.AssertCount(t, db, "eth_tx_attempts", 1)

		attemptL = cltest.NewLegacyEthTxAttempt(t, etx.ID)
		attemptL.State = txmgr.EthTxAttemptBroadcast
		attemptL.GasPrice = utils.NewBigI(42)
		err = orm.InsertEthTxAttempt(&attemptL)
		require.NoError(t, err)
		assert.Greater(t, int(attemptL.ID), 0)
		cltest.AssertCount(t, db, "eth_tx_attempts", 2)
	})
	var r txmgr.EthReceipt
	t.Run("InsertEthReceipt", func(t *testing.T) {
		r = cltest.NewEthReceipt(t, 42, utils.NewHash(), attemptD.Hash)
		err = orm.InsertEthReceipt(&r)
		require.NoError(t, err)
		assert.Greater(t, int(r.ID), 0)
		cltest.AssertCount(t, db, "eth_receipts", 1)
	})
	t.Run("FindEthTxWithAttempts", func(t *testing.T) {
		etx, err = orm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.EthTxAttempts, 2)
		assert.Equal(t, etx.EthTxAttempts[0].ID, attemptD.ID)
		assert.Equal(t, etx.EthTxAttempts[1].ID, attemptL.ID)
		require.Len(t, etx.EthTxAttempts[0].EthReceipts, 1)
		require.Len(t, etx.EthTxAttempts[1].EthReceipts, 0)
		assert.Equal(t, r.BlockHash, etx.EthTxAttempts[0].EthReceipts[0].BlockHash)
	})
	t.Run("FindEthTxByHash", func(t *testing.T) {
		foundEtx, err := orm.FindEthTxByHash(attemptD.Hash)
		require.NoError(t, err)
		assert.Equal(t, etx.ID, foundEtx.ID)
		assert.Equal(t, etx.EVMChainID, foundEtx.EVMChainID)
	})
	t.Run("FindEthTxAttemptsByEthTxIDs", func(t *testing.T) {
		attempts, err := orm.FindEthTxAttemptsByEthTxIDs([]int64{etx.ID})
		require.NoError(t, err)
		require.Len(t, attempts, 2)
		assert.Equal(t, etx.EthTxAttempts[0].ID, attemptD.ID)
		assert.Equal(t, etx.EthTxAttempts[1].ID, attemptL.ID)
		require.Len(t, etx.EthTxAttempts[0].EthReceipts, 1)
		require.Len(t, etx.EthTxAttempts[1].EthReceipts, 0)
		assert.Equal(t, r.BlockHash, etx.EthTxAttempts[0].EthReceipts[0].BlockHash)
	})
}
