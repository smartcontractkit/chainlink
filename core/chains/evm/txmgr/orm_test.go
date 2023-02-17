package txmgr_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestORM_EthTransactionsWithAttempts(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	orm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, orm, 0, 1, from)        // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, orm, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewLegacyEthTxAttempt(t, tx2.ID)
	attempt.State = txmgr.EthTxAttemptBroadcast
	attempt.GasPrice = assets.NewWeiI(3)
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
	cfg := configtest.NewGeneralConfig(t, nil)
	orm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, orm, 0, 1, from)        // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, orm, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewLegacyEthTxAttempt(t, tx2.ID)
	attempt.State = txmgr.EthTxAttemptBroadcast
	attempt.GasPrice = assets.NewWeiI(3)
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
	cfg := configtest.NewGeneralConfig(t, nil)
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
		attemptL.GasPrice = assets.NewWeiI(42)
		err = orm.InsertEthTxAttempt(&attemptL)
		require.NoError(t, err)
		assert.Greater(t, int(attemptL.ID), 0)
		cltest.AssertCount(t, db, "eth_tx_attempts", 2)
	})
	var r txmgr.EthReceipt
	t.Run("InsertEthReceipt", func(t *testing.T) {
		r = cltest.NewEthReceipt(t, 42, utils.NewHash(), attemptD.Hash, 0x1)
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

func TestORM_FindEthTxAttemptConfirmedByEthTxIDs(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	orm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	tx1 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, orm, 0, 1, from) // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, orm, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewLegacyEthTxAttempt(t, tx2.ID)
	attempt.State = txmgr.EthTxAttemptBroadcast
	attempt.GasPrice = assets.NewWeiI(3)
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, orm.InsertEthTxAttempt(&attempt))

	// add receipt for the second attempt
	r := cltest.NewEthReceipt(t, 4, utils.NewHash(), attempt.Hash, 0x1)
	require.NoError(t, orm.InsertEthReceipt(&r))

	// tx 3 has no attempts
	tx3 := cltest.NewEthTx(t, from)
	tx3.State = txmgr.EthTxUnstarted
	tx3.FromAddress = from
	require.NoError(t, orm.InsertEthTx(&tx3))

	cltest.MustInsertUnconfirmedEthTx(t, orm, 3, from)                           // tx4
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, orm, 4, from) // tx5

	var count int
	err := db.Get(&count, `SELECT count(*) FROM eth_txes`)
	require.NoError(t, err)
	require.Equal(t, 5, count)

	err = db.Get(&count, `SELECT count(*) FROM eth_tx_attempts`)
	require.NoError(t, err)
	require.Equal(t, 4, count)

	confirmedAttempts, err := orm.FindEthTxAttemptConfirmedByEthTxIDs([]int64{tx1.ID, tx2.ID}) // should omit tx3
	require.NoError(t, err)
	assert.Equal(t, 4, count, "only eth txs with attempts are counted")
	require.Len(t, confirmedAttempts, 1)
	assert.Equal(t, confirmedAttempts[0].ID, attempt.ID)
	require.Len(t, confirmedAttempts[0].EthReceipts, 1, "should have only one EthRecipts for a confirmed transaction")
	assert.Equal(t, confirmedAttempts[0].EthReceipts[0].BlockHash, r.BlockHash)
	assert.Equal(t, confirmedAttempts[0].Hash, attempt.Hash, "confirmed Recieipt Hash should match the attempt hash")
}

func TestORM_FindEthTxAttemptsRequiringResend(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logCfg := pgtest.NewQConfig(true)
	borm := cltest.NewTxmORM(t, db, logCfg)

	ethKeyStore := cltest.NewKeyStore(t, db, logCfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	t.Run("returns nothing if there are no transactions", func(t *testing.T) {
		olderThan := time.Now()
		attempts, err := borm.FindEthTxAttemptsRequiringResend(olderThan, 10, cltest.FixtureChainID, fromAddress)
		require.NoError(t, err)
		assert.Len(t, attempts, 0)
	})

	// Mix up the insert order to assure that they come out sorted by nonce not implicitly or by ID
	e1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 1, fromAddress, time.Unix(1616509200, 0))
	e3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, borm, 3, fromAddress, time.Unix(1616509400, 0))
	e0 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, fromAddress, time.Unix(1616509100, 0))
	e2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 2, fromAddress, time.Unix(1616509300, 0))

	etxs := []txmgr.EthTx{
		e0,
		e1,
		e2,
		e3,
	}
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etxs[0].ID)
	attempt1_2.GasPrice = assets.NewWeiI(10)
	require.NoError(t, borm.InsertEthTxAttempt(&attempt1_2))

	attempt3_2 := newInProgressLegacyEthTxAttempt(t, etxs[2].ID)
	attempt3_2.GasPrice = assets.NewWeiI(10)
	require.NoError(t, borm.InsertEthTxAttempt(&attempt3_2))

	attempt4_2 := cltest.NewDynamicFeeEthTxAttempt(t, etxs[3].ID)
	attempt4_2.GasTipCap = assets.NewWeiI(10)
	attempt4_2.GasFeeCap = assets.NewWeiI(20)
	attempt4_2.State = txmgr.EthTxAttemptBroadcast
	require.NoError(t, borm.InsertEthTxAttempt(&attempt4_2))
	attempt4_4 := cltest.NewDynamicFeeEthTxAttempt(t, etxs[3].ID)
	attempt4_4.GasTipCap = assets.NewWeiI(30)
	attempt4_4.GasFeeCap = assets.NewWeiI(40)
	attempt4_4.State = txmgr.EthTxAttemptBroadcast
	require.NoError(t, borm.InsertEthTxAttempt(&attempt4_4))
	attempt4_3 := cltest.NewDynamicFeeEthTxAttempt(t, etxs[3].ID)
	attempt4_3.GasTipCap = assets.NewWeiI(20)
	attempt4_3.GasFeeCap = assets.NewWeiI(30)
	attempt4_3.State = txmgr.EthTxAttemptBroadcast
	require.NoError(t, borm.InsertEthTxAttempt(&attempt4_3))

	t.Run("returns nothing if there are transactions from a different key", func(t *testing.T) {
		olderThan := time.Now()
		attempts, err := borm.FindEthTxAttemptsRequiringResend(olderThan, 10, cltest.FixtureChainID, utils.RandomAddress())
		require.NoError(t, err)
		assert.Len(t, attempts, 0)
	})

	t.Run("returns the highest price attempt for each transaction that was last broadcast before or on the given time", func(t *testing.T) {
		olderThan := time.Unix(1616509200, 0)
		attempts, err := borm.FindEthTxAttemptsRequiringResend(olderThan, 0, cltest.FixtureChainID, fromAddress)
		require.NoError(t, err)
		assert.Len(t, attempts, 2)
		assert.Equal(t, attempt1_2.ID, attempts[0].ID)
		assert.Equal(t, etxs[1].EthTxAttempts[0].ID, attempts[1].ID)
	})

	t.Run("returns the highest price attempt for EIP-1559 transactions", func(t *testing.T) {
		olderThan := time.Unix(1616509400, 0)
		attempts, err := borm.FindEthTxAttemptsRequiringResend(olderThan, 0, cltest.FixtureChainID, fromAddress)
		require.NoError(t, err)
		assert.Len(t, attempts, 4)
		assert.Equal(t, attempt4_4.ID, attempts[3].ID)
	})

	t.Run("applies limit", func(t *testing.T) {
		olderThan := time.Unix(1616509200, 0)
		attempts, err := borm.FindEthTxAttemptsRequiringResend(olderThan, 1, cltest.FixtureChainID, fromAddress)
		require.NoError(t, err)
		assert.Len(t, attempts, 1)
		assert.Equal(t, attempt1_2.ID, attempts[0].ID)
	})
}

func TestORM_UpdateBroadcastAts(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	orm := cltest.NewTxmORM(t, db, cfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth(), 0)

	t.Run("does not update when broadcast_at is NULL", func(t *testing.T) {
		t.Parallel()

		etx := cltest.NewEthTx(t, fromAddress)
		err := orm.InsertEthTx(&etx)
		require.NoError(t, err)

		var nullTime *time.Time
		assert.Equal(t, nullTime, etx.BroadcastAt)

		currTime := time.Now()
		err = orm.UpdateBroadcastAts(currTime, []int64{etx.ID})
		require.NoError(t, err)
		etx, err = orm.FindEthTxWithAttempts(etx.ID)

		require.NoError(t, err)
		assert.Equal(t, nullTime, etx.BroadcastAt)
	})

	t.Run("updates when broadcast_at is non-NULL", func(t *testing.T) {
		t.Parallel()

		time1 := time.Now()
		etx := cltest.NewEthTx(t, fromAddress)
		etx.Nonce = new(int64)
		etx.State = txmgr.EthTxUnconfirmed
		etx.BroadcastAt = &time1
		etx.InitialBroadcastAt = &time1
		err := orm.InsertEthTx(&etx)
		require.NoError(t, err)

		time2 := time.Date(2077, 8, 14, 10, 0, 0, 0, time.UTC)
		err = orm.UpdateBroadcastAts(time2, []int64{etx.ID})
		require.NoError(t, err)
		etx, err = orm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		// assert year due to time rounding after database save
		assert.Equal(t, etx.BroadcastAt.Year(), time2.Year())
	})
}

func TestORM_SetBroadcastBeforeBlockNum(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)
	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, fromAddress)
	chainID := *ethClient.ChainID()

	headNum := int64(9000)
	var err error

	t.Run("saves block num to unconfirmed eth_tx_attempts without one", func(t *testing.T) {
		// Do the thing
		require.NoError(t, borm.SetBroadcastBeforeBlockNum(headNum, chainID))

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]

		assert.Equal(t, int64(9000), *attempt.BroadcastBeforeBlockNum)
	})

	t.Run("does not change eth_tx_attempts that already have BroadcastBeforeBlockNum set", func(t *testing.T) {
		n := int64(42)
		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, 2)
		attempt.BroadcastBeforeBlockNum = &n
		require.NoError(t, borm.InsertEthTxAttempt(&attempt))

		// Do the thing
		require.NoError(t, borm.SetBroadcastBeforeBlockNum(headNum, chainID))

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.EthTxAttempts, 2)
		attempt = etx.EthTxAttempts[0]

		assert.Equal(t, int64(42), *attempt.BroadcastBeforeBlockNum)
	})

	t.Run("only updates eth_tx_attempts for the current chain", func(t *testing.T) {
		require.NoError(t, ethKeyStore.Enable(fromAddress, testutils.SimulatedChainID))
		etxThisChain := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 1, fromAddress, cfg.DefaultChainID())
		etxOtherChain := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, fromAddress, testutils.SimulatedChainID)

		require.NoError(t, borm.SetBroadcastBeforeBlockNum(headNum, chainID))

		etxThisChain, err = borm.FindEthTxWithAttempts(etxThisChain.ID)
		require.NoError(t, err)
		require.Len(t, etxThisChain.EthTxAttempts, 1)
		attempt := etxThisChain.EthTxAttempts[0]

		assert.Equal(t, int64(9000), *attempt.BroadcastBeforeBlockNum)

		etxOtherChain, err = borm.FindEthTxWithAttempts(etxOtherChain.ID)
		require.NoError(t, err)
		require.Len(t, etxOtherChain.EthTxAttempts, 1)
		attempt = etxOtherChain.EthTxAttempts[0]

		assert.Nil(t, attempt.BroadcastBeforeBlockNum)
	})
}

func TestORM_FindEtxAttemptsConfirmedMissingReceipt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx0 := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, borm, 0, 1, originalBroadcastAt, fromAddress)

	attempts, err := borm.FindEtxAttemptsConfirmedMissingReceipt(*ethClient.ChainID())

	require.NoError(t, err)

	assert.Len(t, attempts, 1)
	assert.Len(t, etx0.EthTxAttempts, 1)
	assert.Equal(t, etx0.EthTxAttempts[0].ID, attempts[0].ID)
}

func TestORM_UpdateEthTxsUnconfirmed(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx0 := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, borm, 0, 1, originalBroadcastAt, fromAddress)
	assert.Equal(t, etx0.State, txmgr.EthTxConfirmedMissingReceipt)
	require.NoError(t, borm.UpdateEthTxsUnconfirmed([]int64{etx0.ID}))

	etx0, err := borm.FindEthTxWithAttempts(etx0.ID)
	require.NoError(t, err)
	assert.Equal(t, etx0.State, txmgr.EthTxUnconfirmed)
}

func TestORM_FindEthTxAttemptsRequiringReceiptFetch(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx0 := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, borm, 0, 1, originalBroadcastAt, fromAddress)

	attempts, err := borm.FindEthTxAttemptsRequiringReceiptFetch(*ethClient.ChainID())
	require.NoError(t, err)
	assert.Len(t, attempts, 1)
	assert.Len(t, etx0.EthTxAttempts, 1)
	assert.Equal(t, etx0.EthTxAttempts[0].ID, attempts[0].ID)
}

func TestORM_SaveFetchedReceipts(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx0 := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, borm, 0, 1, originalBroadcastAt, fromAddress)
	require.Len(t, etx0.EthTxAttempts, 1)

	// create receipt associated with transaction
	txmReceipt := evmtypes.Receipt{
		TxHash:           etx0.EthTxAttempts[0].Hash,
		BlockHash:        utils.NewHash(),
		BlockNumber:      big.NewInt(42),
		TransactionIndex: uint(1),
	}

	err := borm.SaveFetchedReceipts([]evmtypes.Receipt{txmReceipt}, *ethClient.ChainID())

	require.NoError(t, err)
	etx0, err = borm.FindEthTxWithAttempts(etx0.ID)
	require.NoError(t, err)
	require.Len(t, etx0.EthTxAttempts, 1)
	require.Len(t, etx0.EthTxAttempts[0].EthReceipts, 1)
	require.Equal(t, txmReceipt.BlockHash, etx0.EthTxAttempts[0].EthReceipts[0].BlockHash)
	require.Equal(t, txmgr.EthTxConfirmed, etx0.State)
}
