package txmgr_test

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"testing"
	"time"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
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
	attempt.State = txmgrtypes.TxAttemptBroadcast
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
	attempt.State = txmgrtypes.TxAttemptBroadcast
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
		attemptL.State = txmgrtypes.TxAttemptBroadcast
		attemptL.GasPrice = assets.NewWeiI(42)
		err = orm.InsertEthTxAttempt(&attemptL)
		require.NoError(t, err)
		assert.Greater(t, int(attemptL.ID), 0)
		cltest.AssertCount(t, db, "eth_tx_attempts", 2)
	})
	var r txmgr.EvmReceipt
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
	attempt.State = txmgrtypes.TxAttemptBroadcast
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
	attempt4_2.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, borm.InsertEthTxAttempt(&attempt4_2))
	attempt4_4 := cltest.NewDynamicFeeEthTxAttempt(t, etxs[3].ID)
	attempt4_4.GasTipCap = assets.NewWeiI(30)
	attempt4_4.GasFeeCap = assets.NewWeiI(40)
	attempt4_4.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, borm.InsertEthTxAttempt(&attempt4_4))
	attempt4_3 := cltest.NewDynamicFeeEthTxAttempt(t, etxs[3].ID)
	attempt4_3.GasTipCap = assets.NewWeiI(20)
	attempt4_3.GasFeeCap = assets.NewWeiI(30)
	attempt4_3.State = txmgrtypes.TxAttemptBroadcast
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

func TestORM_MarkAllConfirmedMissingReceipt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	// create transaction 0 (nonce 0) that is unconfirmed (block 7)
	etx0_blocknum := int64(7)
	etx0 := cltest.MustInsertUnconfirmedEthTx(t, borm, 0, fromAddress)
	etx0_attempt := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(1))
	etx0_attempt.BroadcastBeforeBlockNum = &etx0_blocknum
	require.NoError(t, borm.InsertEthTxAttempt(&etx0_attempt))
	assert.Equal(t, txmgr.EthTxUnconfirmed, etx0.State)

	// create transaction 1 (nonce 1) that is confirmed (block 77)
	etx1 := cltest.MustInsertConfirmedEthTxBySaveFetchedReceipts(t, borm, fromAddress, int64(1), int64(77), *ethClient.ChainID())
	assert.Equal(t, etx1.State, txmgr.EthTxConfirmed)

	// mark transaction 0 confirmed_missing_receipt
	err := borm.MarkAllConfirmedMissingReceipt(*ethClient.ChainID())
	require.NoError(t, err)
	etx0, err = borm.FindEthTxWithAttempts(etx0.ID)
	require.NoError(t, err)
	assert.Equal(t, txmgr.EthTxConfirmedMissingReceipt, etx0.State)
}

func TestORM_PreloadEthTxes(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	t.Run("loads eth transaction", func(t *testing.T) {
		// insert etx with attempt
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, int64(7), fromAddress)

		// create unloaded attempt
		unloadedAttempt := txmgr.EthTxAttempt{EthTxID: etx.ID}

		// uninitialized EthTx
		assert.Equal(t, int64(0), unloadedAttempt.EthTx.ID)

		attempts := []txmgr.EthTxAttempt{unloadedAttempt}

		err := borm.PreloadEthTxes(attempts)
		require.NoError(t, err)

		assert.Equal(t, etx.ID, attempts[0].EthTx.ID)
	})

	t.Run("returns nil when attempts slice is empty", func(t *testing.T) {
		emptyAttempts := []txmgr.EthTxAttempt{}
		err := borm.PreloadEthTxes(emptyAttempts)
		require.NoError(t, err)
	})
}

func TestORM_GetInProgressEthTxAttempts(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	// insert etx with attempt
	etx := cltest.MustInsertUnconfirmedEthTxWithAttemptState(t, borm, int64(7), fromAddress, txmgrtypes.TxAttemptInProgress)

	// fetch attempt
	attempts, err := borm.GetInProgressEthTxAttempts(context.Background(), fromAddress, *ethClient.ChainID())
	require.NoError(t, err)

	assert.Len(t, attempts, 1)
	assert.Equal(t, etx.EthTxAttempts[0].ID, attempts[0].ID)
}

func TestORM_FindEthReceiptsPendingConfirmation(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	pgtest.MustExec(t, db, `SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`)

	head := evmtypes.Head{
		Hash:   utils.NewHash(),
		Number: 10,
		Parent: &evmtypes.Head{
			Hash:   utils.NewHash(),
			Number: 9,
			Parent: &evmtypes.Head{
				Number: 8,
				Hash:   utils.NewHash(),
				Parent: nil,
			},
		},
	}

	minConfirmations := int64(2)

	run := cltest.MustInsertPipelineRun(t, db)
	tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)
	pgtest.MustExec(t, db, `UPDATE pipeline_runs SET state = 'suspended' WHERE id = $1`, run.ID)

	etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 3, 1, fromAddress)
	pgtest.MustExec(t, db, `UPDATE eth_txes SET meta='{"FailOnRevert": true}'`)
	attempt := etx.EthTxAttempts[0]
	cltest.MustInsertEthReceipt(t, borm, head.Number-minConfirmations, head.Hash, attempt.Hash)

	pgtest.MustExec(t, db, `UPDATE eth_txes SET pipeline_task_run_id = $1, min_confirmations = $2 WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

	receiptsPlus, err := borm.FindEthReceiptsPendingConfirmation(testutils.Context(t), head.Number, *ethClient.ChainID())
	require.NoError(t, err)
	assert.Len(t, receiptsPlus, 1)
	assert.Equal(t, tr.ID, receiptsPlus[0].ID)
}

func TestORM_FindEthTxWithNonce(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	t.Run("returns nil if no results", func(t *testing.T) {
		etx, err := borm.FindEthTxWithNonce(fromAddress, 777)
		require.NoError(t, err)
		assert.Nil(t, etx)
	})

	t.Run("returns transaction if it exists", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 777, 1, fromAddress)
		require.Equal(t, int64(777), *etx.Nonce)

		res, err := borm.FindEthTxWithNonce(fromAddress, 777)
		require.NoError(t, err)
		assert.Equal(t, etx.Nonce, res.Nonce)
	})
}

func TestORM_UpdateEthTxForRebroadcast(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	t.Run("delete all receipts for eth transaction", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithReceipt(t, borm, fromAddress, 777, 1)
		etx, err := borm.FindEthTxWithAttempts(etx.ID)
		assert.NoError(t, err)
		// assert attempt state
		attempt := etx.EthTxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
		// assert tx state
		assert.Equal(t, txmgr.EthTxConfirmed, etx.State)
		// assert receipt
		assert.Len(t, etx.EthTxAttempts[0].EthReceipts, 1)

		// use exported method
		borm.UpdateEthTxForRebroadcast(etx, attempt)

		resultTx, err := borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Len(t, resultTx.EthTxAttempts, 1)
		resultTxAttempt := resultTx.EthTxAttempts[0]

		// assert attempt state
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, resultTxAttempt.State)
		assert.Nil(t, resultTxAttempt.BroadcastBeforeBlockNum)
		// assert tx state
		assert.Equal(t, txmgr.EthTxUnconfirmed, resultTx.State)
		// assert receipt
		assert.Len(t, resultTxAttempt.EthReceipts, 0)
	})
}

func TestORM_FindTransactionsConfirmedInBlockRange(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	head := evmtypes.Head{
		Hash:   utils.NewHash(),
		Number: 10,
		Parent: &evmtypes.Head{
			Hash:   utils.NewHash(),
			Number: 9,
			Parent: &evmtypes.Head{
				Number: 8,
				Hash:   utils.NewHash(),
				Parent: nil,
			},
		},
	}

	t.Run("find all transactions confirmed in range", func(t *testing.T) {
		etx_8 := cltest.MustInsertConfirmedEthTxWithReceipt(t, borm, fromAddress, 700, 8)
		etx_9 := cltest.MustInsertConfirmedEthTxWithReceipt(t, borm, fromAddress, 777, 9)

		etxes, err := borm.FindTransactionsConfirmedInBlockRange(head.Number, 8, *ethClient.ChainID())
		require.NoError(t, err)
		assert.Len(t, etxes, 2)
		assert.Equal(t, etxes[0].Nonce, etx_8.Nonce)
		assert.Equal(t, etxes[1].Nonce, etx_9.Nonce)
	})
}

func TestORM_SaveInsufficientEthAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)
	defaultDuration, err := time.ParseDuration("5s")
	require.NoError(t, err)

	t.Run("updates attempt state", func(t *testing.T) {
		etx := cltest.MustInsertInProgressEthTxWithAttempt(t, borm, 1, fromAddress)
		now := time.Now()

		err = borm.SaveInsufficientEthAttempt(defaultDuration, &etx.EthTxAttempts[0], now)
		require.NoError(t, err)

		attempt, err := borm.FindEthTxAttempt(etx.EthTxAttempts[0].Hash)
		require.NoError(t, err)
		assert.Equal(t, txmgrtypes.TxAttemptInsufficientEth, attempt.State)
	})
}

func TestORM_SaveSentAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)
	defaultDuration, err := time.ParseDuration("5s")
	require.NoError(t, err)

	t.Run("updates attempt state to 'broadcast'", func(t *testing.T) {
		etx := cltest.MustInsertInProgressEthTxWithAttempt(t, borm, 1, fromAddress)
		require.Nil(t, etx.BroadcastAt)
		now := time.Now()

		err = borm.SaveSentAttempt(defaultDuration, &etx.EthTxAttempts[0], now)
		require.NoError(t, err)

		attempt, err := borm.FindEthTxAttempt(etx.EthTxAttempts[0].Hash)
		require.NoError(t, err)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
	})
}

func TestORM_SaveConfirmedMissingReceiptAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)
	defaultDuration, err := time.ParseDuration("5s")
	require.NoError(t, err)

	t.Run("updates attempt to 'broadcast' and transaction to 'confirm_missing_receipt'", func(t *testing.T) {
		etx := cltest.MustInsertUnconfirmedEthTxWithAttemptState(t, borm, 1, fromAddress, txmgrtypes.TxAttemptInProgress)
		now := time.Now()

		err = borm.SaveConfirmedMissingReceiptAttempt(context.Background(), defaultDuration, &etx.EthTxAttempts[0], now)
		require.NoError(t, err)

		etx, err := borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxConfirmedMissingReceipt, etx.State)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, etx.EthTxAttempts[0].State)
	})
}

func TestORM_DeleteInProgressAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	t.Run("deletes in_progress attempt", func(t *testing.T) {
		etx := cltest.MustInsertInProgressEthTxWithAttempt(t, borm, 1, fromAddress)
		attempt := etx.EthTxAttempts[0]

		err := borm.DeleteInProgressAttempt(testutils.Context(t), etx.EthTxAttempts[0])
		require.NoError(t, err)

		nilResult, err := borm.FindEthTxAttempt(attempt.Hash)
		assert.Nil(t, nilResult)
		require.Error(t, err)
	})
}

func TestORM_SaveInProgressAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	t.Run("saves new in_progress attempt if attempt is new", func(t *testing.T) {
		etx := cltest.MustInsertUnconfirmedEthTx(t, borm, 1, fromAddress)

		attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)
		require.Equal(t, int64(0), attempt.ID)

		err := borm.SaveInProgressAttempt(&attempt)
		require.NoError(t, err)

		attemptResult, err := borm.FindEthTxAttempt(attempt.Hash)
		require.NoError(t, err)
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, attemptResult.State)
	})

	t.Run("updates old attempt to in_progress when insufficient_eth", func(t *testing.T) {
		etx := cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, 23, fromAddress)
		attempt := etx.EthTxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptInsufficientEth, attempt.State)
		require.NotEqual(t, 0, attempt.ID)

		attempt.BroadcastBeforeBlockNum = nil
		attempt.State = txmgrtypes.TxAttemptInProgress
		err := borm.SaveInProgressAttempt(&attempt)

		require.NoError(t, err)
		attemptResult, err := borm.FindEthTxAttempt(attempt.Hash)
		require.NoError(t, err)
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, attemptResult.State)

	})
}

func TestORM_FindEthTxsRequiringGasBump(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	currentBlockNum := int64(10)

	t.Run("gets txs requiring gas bump", func(t *testing.T) {
		etx := cltest.MustInsertUnconfirmedEthTxWithAttemptState(t, borm, 1, fromAddress, txmgrtypes.TxAttemptBroadcast)
		borm.SetBroadcastBeforeBlockNum(currentBlockNum, *ethClient.ChainID())

		// this tx will require gas bump
		etx, err := borm.FindEthTxWithAttempts(etx.ID)
		attempts := etx.EthTxAttempts
		require.NoError(t, err)
		assert.Len(t, attempts, 1)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempts[0].State)
		assert.Equal(t, currentBlockNum, *attempts[0].BroadcastBeforeBlockNum)

		// this tx will not require gas bump
		cltest.MustInsertUnconfirmedEthTxWithAttemptState(t, borm, 2, fromAddress, txmgrtypes.TxAttemptBroadcast)
		borm.SetBroadcastBeforeBlockNum(currentBlockNum+1, *ethClient.ChainID())

		// any tx broadcast <= 10 will require gas bump
		newBlock := int64(12)
		gasBumpThreshold := int64(2)
		etxs, err := borm.FindEthTxsRequiringGasBump(context.Background(), fromAddress, newBlock, gasBumpThreshold, int64(0), *ethClient.ChainID())
		require.NoError(t, err)
		assert.Len(t, etxs, 1)
		assert.Equal(t, etx.ID, etxs[0].ID)
	})
}

func TestEthConfirmer_FindEthTxsRequiringResubmissionDueToInsufficientEth(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	borm := cltest.NewTxmORM(t, db, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)
	_, otherAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	// Insert order is mixed up to test sorting
	etx2 := cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, 1, fromAddress)
	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 2, fromAddress)
	attempt3_2 := cltest.NewLegacyEthTxAttempt(t, etx3.ID)
	attempt3_2.State = txmgrtypes.TxAttemptInsufficientEth
	attempt3_2.GasPrice = assets.NewWeiI(100)
	require.NoError(t, borm.InsertEthTxAttempt(&attempt3_2))
	etx1 := cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, 0, fromAddress)

	// These should never be returned
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 3, fromAddress)
	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 4, 100, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, 0, otherAddress)

	t.Run("returns all eth_txes with at least one attempt that is in insufficient_eth state", func(t *testing.T) {
		etxs, err := borm.FindEthTxsRequiringResubmissionDueToInsufficientEth(fromAddress, cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 3)

		assert.Equal(t, *etx1.Nonce, *etxs[0].Nonce)
		assert.Equal(t, etx1.ID, etxs[0].ID)
		assert.Equal(t, *etx2.Nonce, *etxs[1].Nonce)
		assert.Equal(t, etx2.ID, etxs[1].ID)
		assert.Equal(t, *etx3.Nonce, *etxs[2].Nonce)
		assert.Equal(t, etx3.ID, etxs[2].ID)
	})

	t.Run("does not return eth_txes with different chain ID", func(t *testing.T) {
		etxs, err := borm.FindEthTxsRequiringResubmissionDueToInsufficientEth(fromAddress, *big.NewInt(42))
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	t.Run("does not return confirmed or fatally errored eth_txes", func(t *testing.T) {
		pgtest.MustExec(t, db, `UPDATE eth_txes SET state='confirmed' WHERE id = $1`, etx1.ID)
		pgtest.MustExec(t, db, `UPDATE eth_txes SET state='fatal_error', nonce=NULL, error='foo', broadcast_at=NULL, initial_broadcast_at=NULL WHERE id = $1`, etx2.ID)

		etxs, err := borm.FindEthTxsRequiringResubmissionDueToInsufficientEth(fromAddress, cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 1)

		assert.Equal(t, *etx3.Nonce, *etxs[0].Nonce)
		assert.Equal(t, etx3.ID, etxs[0].ID)
	})
}

func TestORM_MarkOldTxesMissingReceiptAsErrored(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	// tx state should be confirmed missing receipt
	// attempt should be broadcast before cutoff time
	t.Run("succesfully mark errored transactions", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, borm, 1, 7, time.Now(), fromAddress)

		err := borm.MarkOldTxesMissingReceiptAsErrored(10, 2, *ethClient.ChainID())
		require.NoError(t, err)

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxFatalError, etx.State)
	})

	t.Run("succesfully mark errored transactions w/ qopt passing in sql.Tx", func(t *testing.T) {
		q := pg.NewQ(db, logger.TestLogger(t), cfg)

		etx := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, borm, 1, 7, time.Now(), fromAddress)
		q.Transaction(func(q pg.Queryer) error {
			err := borm.MarkOldTxesMissingReceiptAsErrored(10, 2, *ethClient.ChainID(), pg.WithQueryer(q))
			require.NoError(t, err)
			return nil
		})
		// must run other query outside of postgres transaction so changes are committed
		etx, err := borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxFatalError, etx.State)
	})
}

func TestORM_LoadEthTxesAttempts(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	t.Run("load eth tx attempt", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, borm, 1, 7, time.Now(), fromAddress)
		etx.EthTxAttempts = []txmgr.EthTxAttempt{}

		err := borm.LoadEthTxesAttempts([]*txmgr.EthTx{&etx})
		require.NoError(t, err)
		assert.Len(t, etx.EthTxAttempts, 1)
	})

	t.Run("load new attempt inserted in current postgres transaction", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, borm, 3, 9, time.Now(), fromAddress)
		etx.EthTxAttempts = []txmgr.EthTxAttempt{}

		q := pg.NewQ(db, logger.TestLogger(t), cfg)

		newAttempt := cltest.NewDynamicFeeEthTxAttempt(t, etx.ID)
		q.Transaction(func(tx pg.Queryer) error {
			const insertEthTxAttemptSQL = `INSERT INTO eth_tx_attempts (eth_tx_id, gas_price, signed_raw_tx, hash, broadcast_before_block_num, state, created_at, chain_specific_gas_limit, tx_type, gas_tip_cap, gas_fee_cap) VALUES (
				:eth_tx_id, :gas_price, :signed_raw_tx, :hash, :broadcast_before_block_num, :state, NOW(), :chain_specific_gas_limit, :tx_type, :gas_tip_cap, :gas_fee_cap
				) RETURNING *`
			_, err := tx.NamedExec(insertEthTxAttemptSQL, newAttempt)
			require.NoError(t, err)

			err = borm.LoadEthTxesAttempts([]*txmgr.EthTx{&etx}, pg.WithQueryer(tx))
			require.NoError(t, err)
			assert.Len(t, etx.EthTxAttempts, 2)

			return nil
		})
		// also check after postgres transaction is committed
		etx.EthTxAttempts = []txmgr.EthTxAttempt{}
		err := borm.LoadEthTxesAttempts([]*txmgr.EthTx{&etx})
		require.NoError(t, err)
		assert.Len(t, etx.EthTxAttempts, 2)
	})
}

func TestORM_SaveReplacementInProgressAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	t.Run("replace eth tx attempt", func(t *testing.T) {
		etx := cltest.MustInsertInProgressEthTxWithAttempt(t, borm, 123, fromAddress)
		oldAttempt := etx.EthTxAttempts[0]

		newAttempt := cltest.NewDynamicFeeEthTxAttempt(t, etx.ID)
		err := borm.SaveReplacementInProgressAttempt(oldAttempt, &newAttempt)
		require.NoError(t, err)

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Len(t, etx.EthTxAttempts, 1)
		require.Equal(t, etx.EthTxAttempts[0].Hash, newAttempt.Hash)
	})
}

func TestORM_FindNextUnstartedTransactionFromAddress(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	t.Run("cannot find unstarted tx", func(t *testing.T) {
		cltest.MustInsertInProgressEthTxWithAttempt(t, borm, 13, fromAddress)

		resultEtx := new(txmgr.EthTx)
		err := borm.FindNextUnstartedTransactionFromAddress(resultEtx, fromAddress, *ethClient.ChainID())
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("finds unstarted tx", func(t *testing.T) {
		cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)

		resultEtx := new(txmgr.EthTx)
		err := borm.FindNextUnstartedTransactionFromAddress(resultEtx, fromAddress, *ethClient.ChainID())
		require.NoError(t, err)
	})
}

func TestORM_UpdateEthTxFatalError(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	t.Run("update successful", func(t *testing.T) {
		etx := cltest.MustInsertInProgressEthTxWithAttempt(t, borm, 13, fromAddress)
		etxPretendError := null.StringFrom("no more toilet paper")
		etx.Error = etxPretendError

		err := borm.UpdateEthTxFatalError(&etx)
		require.NoError(t, err)
		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Len(t, etx.EthTxAttempts, 0)
		assert.Equal(t, txmgr.EthTxFatalError, etx.State)
	})
}

func TestORM_UpdateEthTxAttemptInProgressToBroadcast(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	t.Run("update successful", func(t *testing.T) {
		etx := cltest.MustInsertInProgressEthTxWithAttempt(t, borm, 13, fromAddress)
		attempt := etx.EthTxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptInProgress, attempt.State)

		time1 := time.Now()
		i := int16(0)
		etx.BroadcastAt = &time1
		etx.InitialBroadcastAt = &time1
		err := borm.UpdateEthTxAttemptInProgressToBroadcast(&etx, attempt, txmgrtypes.TxAttemptBroadcast, func(_ pg.Queryer) error {
			// dummy function because tests do not use keystore as source of truth for next nonce number
			i++
			return nil
		})
		require.NoError(t, err)

		attemptResult, err := borm.FindEthTxAttempt(attempt.Hash)
		require.NoError(t, err)
		require.Equal(t, attempt.Hash, attemptResult.Hash)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attemptResult.State)
		assert.Equal(t, int16(1), i)
	})
}

func TestORM_UpdateEthTxUnstartedToInProgress(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)
	q := pg.NewQ(db, logger.TestLogger(t), cfg)

	t.Run("update successful", func(t *testing.T) {
		nonce := int64(123)
		etx := cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)
		etx.Nonce = &nonce

		attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)

		err := borm.UpdateEthTxUnstartedToInProgress(&etx, &attempt)
		require.NoError(t, err)

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxInProgress, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
	})

	t.Run("update fails because tx is removed", func(t *testing.T) {
		nonce := int64(123)
		etx := cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)
		etx.Nonce = &nonce

		attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)

		err := q.ExecQ("DELETE FROM eth_txes WHERE id = $1", etx.ID)
		require.NoError(t, err)

		err = borm.UpdateEthTxUnstartedToInProgress(&etx, &attempt)
		require.ErrorContains(t, err, "eth_tx removed")
	})
}

func TestORM_GetEthTxInProgress(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	t.Run("gets 0 in progress eth transaction", func(t *testing.T) {
		etxResult, err := borm.GetEthTxInProgress(fromAddress)
		require.NoError(t, err)
		require.Nil(t, etxResult)
	})

	t.Run("get 1 in progress eth transaction", func(t *testing.T) {
		etx := cltest.MustInsertInProgressEthTxWithAttempt(t, borm, 123, fromAddress)

		etxResult, err := borm.GetEthTxInProgress(fromAddress)
		require.NoError(t, err)
		assert.Equal(t, etxResult.ID, etx.ID)
	})
}

func TestORM_HasInProgressTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	t.Run("no in progress eth transaction", func(t *testing.T) {
		exists, err := borm.HasInProgressTransaction(fromAddress, *ethClient.ChainID())
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("has in progress eth transaction", func(t *testing.T) {
		cltest.MustInsertInProgressEthTxWithAttempt(t, borm, 123, fromAddress)

		exists, err := borm.HasInProgressTransaction(fromAddress, *ethClient.ChainID())
		require.NoError(t, err)
		require.True(t, exists)
	})
}

func TestORM_UpdateEthKeyNextNonce(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	ethKeyState, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	t.Run("update next nonce", func(t *testing.T) {
		assert.Equal(t, int64(0), ethKeyState.NextNonce)
		err := borm.UpdateEthKeyNextNonce(int64(24), int64(0), fromAddress, *ethClient.ChainID())
		require.NoError(t, err)

		newNextNonce, err := ethKeyStore.GetNextNonce(fromAddress, ethClient.ChainID())
		require.NoError(t, err)
		assert.Equal(t, int64(24), newNextNonce)
	})

	t.Run("no rows found", func(t *testing.T) {
		err := borm.UpdateEthKeyNextNonce(int64(100), int64(123), fromAddress, *ethClient.ChainID())
		require.Error(t, err)
	})
}

func TestORM_CountUnconfirmedTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)
	_, otherAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, otherAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 1, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 2, fromAddress)

	count, err := borm.CountUnconfirmedTransactions(fromAddress, cltest.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, int(count), 3)
}

func TestORM_CountUnstartedTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)
	_, otherAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)
	cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)
	cltest.MustInsertUnstartedEthTx(t, borm, otherAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 2, fromAddress)

	count, err := borm.CountUnstartedTransactions(fromAddress, cltest.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, int(count), 2)
}

func TestORM_CheckEthTxQueueCapacity(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	var maxUnconfirmedTransactions uint64 = 2

	t.Run("with no eth_txes returns nil", func(t *testing.T) {
		err := borm.CheckEthTxQueueCapacity(fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertUnstartedEthTx(t, borm, otherAddress)
	}

	t.Run("with eth_txes from another address returns nil", func(t *testing.T) {
		err := borm.CheckEthTxQueueCapacity(fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertFatalErrorEthTx(t, borm, otherAddress)
	}

	t.Run("ignores fatally_errored transactions", func(t *testing.T) {
		err := borm.CheckEthTxQueueCapacity(fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	var n int64
	cltest.MustInsertInProgressEthTxWithAttempt(t, borm, n, fromAddress)
	n++
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, n, fromAddress)
	n++

	t.Run("unconfirmed and in_progress transactions do not count", func(t *testing.T) {
		err := borm.CheckEthTxQueueCapacity(fromAddress, 1, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, n, 42, fromAddress)
		n++
	}

	t.Run("with many confirmed eth_txes from the same address returns nil", func(t *testing.T) {
		err := borm.CheckEthTxQueueCapacity(fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	for i := 0; i < int(maxUnconfirmedTransactions)-1; i++ {
		cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)
	}

	t.Run("with fewer unstarted eth_txes than limit returns nil", func(t *testing.T) {
		err := borm.CheckEthTxQueueCapacity(fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.NoError(t, err)
	})

	cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)

	t.Run("with equal or more unstarted eth_txes than limit returns error", func(t *testing.T) {
		err := borm.CheckEthTxQueueCapacity(fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("cannot create transaction; too many unstarted transactions in the queue (2/%d). WARNING: Hitting EVM.Transactions.MaxQueued", maxUnconfirmedTransactions))

		cltest.MustInsertUnstartedEthTx(t, borm, fromAddress)
		err = borm.CheckEthTxQueueCapacity(fromAddress, maxUnconfirmedTransactions, cltest.FixtureChainID)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("cannot create transaction; too many unstarted transactions in the queue (3/%d). WARNING: Hitting EVM.Transactions.MaxQueued", maxUnconfirmedTransactions))
	})

	t.Run("with different chain ID ignores txes", func(t *testing.T) {
		err := borm.CheckEthTxQueueCapacity(fromAddress, maxUnconfirmedTransactions, *big.NewInt(42))
		require.NoError(t, err)
	})

	t.Run("disables check with 0 limit", func(t *testing.T) {
		err := borm.CheckEthTxQueueCapacity(fromAddress, 0, cltest.FixtureChainID)
		require.NoError(t, err)
	})
}

func TestORM_CreateEthTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	borm := cltest.NewTxmORM(t, db, cfg)
	kst := cltest.NewKeyStore(t, db, cfg)

	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth(), 0)
	toAddress := testutils.NewAddress()
	gasLimit := uint32(1000)
	payload := []byte{1, 2, 3}

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	t.Run("with queue under capacity inserts eth_tx", func(t *testing.T) {
		subject := uuid.NewV4()
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		strategy.On("PruneQueue", mock.AnythingOfType("*txmgr.orm"), mock.AnythingOfType("pg.QOpt")).Return(int64(0), nil)
		etx, err := borm.CreateEthTransaction(txmgr.NewTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			GasLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		}, *ethClient.ChainID())
		assert.NoError(t, err)

		assert.Greater(t, etx.ID, int64(0))
		assert.Equal(t, etx.State, txmgr.EthTxUnstarted)
		assert.Equal(t, gasLimit, etx.GasLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, assets.NewEthValue(0), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)

		cltest.AssertCount(t, db, "eth_txes", 1)

		require.NoError(t, db.Get(&etx, `SELECT * FROM eth_txes ORDER BY id ASC LIMIT 1`))

		assert.Equal(t, etx.State, txmgr.EthTxUnstarted)
		assert.Equal(t, gasLimit, etx.GasLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, assets.NewEthValue(0), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)
	})

	t.Run("doesn't insert eth_tx if a matching tx already exists for that pipeline_task_run_id", func(t *testing.T) {
		id := uuid.NewV4()
		newTx := txmgr.NewTx{
			FromAddress:       fromAddress,
			ToAddress:         testutils.NewAddress(),
			EncodedPayload:    []byte{1, 2, 3},
			GasLimit:          21000,
			PipelineTaskRunID: &id,
			Strategy:          txmgr.SendEveryStrategy{},
		}
		tx1, err := borm.CreateEthTransaction(newTx, *ethClient.ChainID())
		assert.NoError(t, err)

		tx2, err := borm.CreateEthTransaction(newTx, *ethClient.ChainID())
		assert.NoError(t, err)

		assert.Equal(t, tx1.ID, tx2.ID)
	})
}

func TestORM_PruneUnstartedTxQueue(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	subject1 := uuid.NewV4()
	for i := 0; i < 5; i++ {
		cltest.MustInsertUnstartedEthTx(t, borm, fromAddress, subject1)
	}
	subject2 := uuid.NewV4()
	for i := 0; i < 5; i++ {
		cltest.MustInsertUnstartedEthTx(t, borm, fromAddress, subject2)
	}

	t.Run("does not prune if queue has not exceeded capacity", func(t *testing.T) {
		n, err := borm.PruneUnstartedTxQueue(uint32(5), subject1)
		require.NoError(t, err)
		assert.Equal(t, int64(0), n)
	})

	t.Run("prunes if queue has exceeded capacity", func(t *testing.T) {
		n, err := borm.PruneUnstartedTxQueue(uint32(3), subject1)
		require.NoError(t, err)
		assert.Equal(t, int64(2), n)
	})
}
