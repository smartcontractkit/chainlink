package txmgr_test

import (
	"database/sql"
	"fmt"
	"math/big"
	"testing"
	"time"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestORM_TransactionsWithAttempts(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)

	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 1, from)        // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewLegacyEthTxAttempt(t, tx2.ID)
	attempt.State = txmgrtypes.TxAttemptBroadcast
	attempt.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(3)}
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, txStore.InsertTxAttempt(&attempt))

	// tx 3 has no attempts
	mustCreateUnstartedGeneratedTx(t, txStore, from, &cltest.FixtureChainID)

	var count int
	err := db.Get(&count, `SELECT count(*) FROM evm.txes`)
	require.NoError(t, err)
	require.Equal(t, 3, count)

	txs, count, err := txStore.TransactionsWithAttempts(0, 100) // should omit tx3
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 2)
	assert.Equal(t, evmtypes.Nonce(1), *txs[0].Sequence, "transactions should be sorted by nonce")
	assert.Equal(t, evmtypes.Nonce(0), *txs[1].Sequence, "transactions should be sorted by nonce")
	assert.Len(t, txs[0].TxAttempts, 2, "all eth tx attempts are preloaded")
	assert.Len(t, txs[1].TxAttempts, 1)
	assert.Equal(t, int64(3), *txs[0].TxAttempts[0].BroadcastBeforeBlockNum, "attempts should be sorted by created_at")
	assert.Equal(t, int64(2), *txs[0].TxAttempts[1].BroadcastBeforeBlockNum, "attempts should be sorted by created_at")

	txs, count, err = txStore.TransactionsWithAttempts(0, 1)
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 1, "limit should apply to length of results")
	assert.Equal(t, evmtypes.Nonce(1), *txs[0].Sequence, "transactions should be sorted by nonce")
}

func TestORM_Transactions(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)

	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 1, from)        // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewLegacyEthTxAttempt(t, tx2.ID)
	attempt.State = txmgrtypes.TxAttemptBroadcast
	attempt.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(3)}
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, txStore.InsertTxAttempt(&attempt))

	// tx 3 has no attempts
	mustCreateUnstartedGeneratedTx(t, txStore, from, &cltest.FixtureChainID)

	var count int
	err := db.Get(&count, `SELECT count(*) FROM evm.txes`)
	require.NoError(t, err)
	require.Equal(t, 3, count)

	txs, count, err := txStore.Transactions(0, 100)
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 2)
	assert.Equal(t, evmtypes.Nonce(1), *txs[0].Sequence, "transactions should be sorted by nonce")
	assert.Equal(t, evmtypes.Nonce(0), *txs[1].Sequence, "transactions should be sorted by nonce")
	assert.Len(t, txs[0].TxAttempts, 0, "eth tx attempts should not be preloaded")
	assert.Len(t, txs[1].TxAttempts, 0)
}

func TestORM(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	keyStore := cltest.NewKeyStore(t, db, cfg.Database())
	orm := cltest.NewTestTxStore(t, db, cfg.Database())
	_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth())

	var etx txmgr.Tx
	t.Run("InsertTx", func(t *testing.T) {
		etx = cltest.NewEthTx(fromAddress)
		require.NoError(t, orm.InsertTx(&etx))
		assert.Greater(t, int(etx.ID), 0)
		cltest.AssertCount(t, db, "evm.txes", 1)
	})
	var attemptL txmgr.TxAttempt
	var attemptD txmgr.TxAttempt
	t.Run("InsertTxAttempt", func(t *testing.T) {
		attemptD = cltest.NewDynamicFeeEthTxAttempt(t, etx.ID)
		require.NoError(t, orm.InsertTxAttempt(&attemptD))
		assert.Greater(t, int(attemptD.ID), 0)
		cltest.AssertCount(t, db, "evm.tx_attempts", 1)

		attemptL = cltest.NewLegacyEthTxAttempt(t, etx.ID)
		attemptL.State = txmgrtypes.TxAttemptBroadcast
		attemptL.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(42)}
		require.NoError(t, orm.InsertTxAttempt(&attemptL))
		assert.Greater(t, int(attemptL.ID), 0)
		cltest.AssertCount(t, db, "evm.tx_attempts", 2)
	})
	var r txmgr.Receipt
	t.Run("InsertReceipt", func(t *testing.T) {
		r = newEthReceipt(42, utils.NewHash(), attemptD.Hash, 0x1)
		id, err := orm.InsertReceipt(&r.Receipt)
		r.ID = id
		require.NoError(t, err)
		assert.Greater(t, int(r.ID), 0)
		cltest.AssertCount(t, db, "evm.receipts", 1)
	})
	t.Run("FindTxWithAttempts", func(t *testing.T) {
		var err error
		etx, err = orm.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 2)
		assert.Equal(t, etx.TxAttempts[0].ID, attemptD.ID)
		assert.Equal(t, etx.TxAttempts[1].ID, attemptL.ID)
		require.Len(t, etx.TxAttempts[0].Receipts, 1)
		require.Len(t, etx.TxAttempts[1].Receipts, 0)
		assert.Equal(t, r.BlockHash, etx.TxAttempts[0].Receipts[0].GetBlockHash())
	})
	t.Run("FindTxByHash", func(t *testing.T) {
		foundEtx, err := orm.FindTxByHash(attemptD.Hash)
		require.NoError(t, err)
		assert.Equal(t, etx.ID, foundEtx.ID)
		assert.Equal(t, etx.ChainID, foundEtx.ChainID)
	})
	t.Run("FindTxAttemptsByTxIDs", func(t *testing.T) {
		attempts, err := orm.FindTxAttemptsByTxIDs([]int64{etx.ID})
		require.NoError(t, err)
		require.Len(t, attempts, 2)
		assert.Equal(t, etx.TxAttempts[0].ID, attemptD.ID)
		assert.Equal(t, etx.TxAttempts[1].ID, attemptL.ID)
		require.Len(t, etx.TxAttempts[0].Receipts, 1)
		require.Len(t, etx.TxAttempts[1].Receipts, 0)
		assert.Equal(t, r.BlockHash, etx.TxAttempts[0].Receipts[0].GetBlockHash())
	})
}

func TestORM_FindTxAttemptConfirmedByTxIDs(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	orm := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)

	tx1 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, orm, 0, 1, from) // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, orm, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewLegacyEthTxAttempt(t, tx2.ID)
	attempt.State = txmgrtypes.TxAttemptBroadcast
	attempt.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(3)}
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, orm.InsertTxAttempt(&attempt))

	// add receipt for the second attempt
	r := newEthReceipt(4, utils.NewHash(), attempt.Hash, 0x1)
	_, err := orm.InsertReceipt(&r.Receipt)
	require.NoError(t, err)

	// tx 3 has no attempts
	mustCreateUnstartedGeneratedTx(t, orm, from, &cltest.FixtureChainID)

	cltest.MustInsertUnconfirmedEthTx(t, orm, 3, from)                           // tx4
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, orm, 4, from) // tx5

	var count int
	err = db.Get(&count, `SELECT count(*) FROM evm.txes`)
	require.NoError(t, err)
	require.Equal(t, 5, count)

	err = db.Get(&count, `SELECT count(*) FROM evm.tx_attempts`)
	require.NoError(t, err)
	require.Equal(t, 4, count)

	confirmedAttempts, err := orm.FindTxAttemptConfirmedByTxIDs([]int64{tx1.ID, tx2.ID}) // should omit tx3
	require.NoError(t, err)
	assert.Equal(t, 4, count, "only eth txs with attempts are counted")
	require.Len(t, confirmedAttempts, 1)
	assert.Equal(t, confirmedAttempts[0].ID, attempt.ID)
	require.Len(t, confirmedAttempts[0].Receipts, 1, "should have only one EthRecipts for a confirmed transaction")
	assert.Equal(t, confirmedAttempts[0].Receipts[0].GetBlockHash(), r.BlockHash)
	assert.Equal(t, confirmedAttempts[0].Hash, attempt.Hash, "confirmed Recieipt Hash should match the attempt hash")
}

func TestORM_FindTxAttemptsRequiringResend(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	logCfg := pgtest.NewQConfig(true)
	txStore := cltest.NewTestTxStore(t, db, logCfg)

	ethKeyStore := cltest.NewKeyStore(t, db, logCfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	t.Run("returns nothing if there are no transactions", func(t *testing.T) {
		olderThan := time.Now()
		attempts, err := txStore.FindTxAttemptsRequiringResend(testutils.Context(t), olderThan, 10, &cltest.FixtureChainID, fromAddress)
		require.NoError(t, err)
		assert.Len(t, attempts, 0)
	})

	// Mix up the insert order to assure that they come out sorted by nonce not implicitly or by ID
	e1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress, time.Unix(1616509200, 0))
	e3 := mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, txStore, 3, fromAddress, time.Unix(1616509400, 0))
	e0 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress, time.Unix(1616509100, 0))
	e2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress, time.Unix(1616509300, 0))

	etxs := []txmgr.Tx{
		e0,
		e1,
		e2,
		e3,
	}
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etxs[0].ID)
	attempt1_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(10)}
	require.NoError(t, txStore.InsertTxAttempt(&attempt1_2))

	attempt3_2 := newInProgressLegacyEthTxAttempt(t, etxs[2].ID)
	attempt3_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(10)}
	require.NoError(t, txStore.InsertTxAttempt(&attempt3_2))

	attempt4_2 := cltest.NewDynamicFeeEthTxAttempt(t, etxs[3].ID)
	attempt4_2.TxFee.DynamicTipCap = assets.NewWeiI(10)
	attempt4_2.TxFee.DynamicFeeCap = assets.NewWeiI(20)
	attempt4_2.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, txStore.InsertTxAttempt(&attempt4_2))
	attempt4_4 := cltest.NewDynamicFeeEthTxAttempt(t, etxs[3].ID)
	attempt4_4.TxFee.DynamicTipCap = assets.NewWeiI(30)
	attempt4_4.TxFee.DynamicFeeCap = assets.NewWeiI(40)
	attempt4_4.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, txStore.InsertTxAttempt(&attempt4_4))
	attempt4_3 := cltest.NewDynamicFeeEthTxAttempt(t, etxs[3].ID)
	attempt4_3.TxFee.DynamicTipCap = assets.NewWeiI(20)
	attempt4_3.TxFee.DynamicFeeCap = assets.NewWeiI(30)
	attempt4_3.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, txStore.InsertTxAttempt(&attempt4_3))

	t.Run("returns nothing if there are transactions from a different key", func(t *testing.T) {
		olderThan := time.Now()
		attempts, err := txStore.FindTxAttemptsRequiringResend(testutils.Context(t), olderThan, 10, &cltest.FixtureChainID, utils.RandomAddress())
		require.NoError(t, err)
		assert.Len(t, attempts, 0)
	})

	t.Run("returns the highest price attempt for each transaction that was last broadcast before or on the given time", func(t *testing.T) {
		olderThan := time.Unix(1616509200, 0)
		attempts, err := txStore.FindTxAttemptsRequiringResend(testutils.Context(t), olderThan, 0, &cltest.FixtureChainID, fromAddress)
		require.NoError(t, err)
		assert.Len(t, attempts, 2)
		assert.Equal(t, attempt1_2.ID, attempts[0].ID)
		assert.Equal(t, etxs[1].TxAttempts[0].ID, attempts[1].ID)
	})

	t.Run("returns the highest price attempt for EIP-1559 transactions", func(t *testing.T) {
		olderThan := time.Unix(1616509400, 0)
		attempts, err := txStore.FindTxAttemptsRequiringResend(testutils.Context(t), olderThan, 0, &cltest.FixtureChainID, fromAddress)
		require.NoError(t, err)
		assert.Len(t, attempts, 4)
		assert.Equal(t, attempt4_4.ID, attempts[3].ID)
	})

	t.Run("applies limit", func(t *testing.T) {
		olderThan := time.Unix(1616509200, 0)
		attempts, err := txStore.FindTxAttemptsRequiringResend(testutils.Context(t), olderThan, 1, &cltest.FixtureChainID, fromAddress)
		require.NoError(t, err)
		assert.Len(t, attempts, 1)
		assert.Equal(t, attempt1_2.ID, attempts[0].ID)
	})
}

func TestORM_UpdateBroadcastAts(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	keyStore := cltest.NewKeyStore(t, db, cfg.Database())
	orm := cltest.NewTestTxStore(t, db, cfg.Database())
	_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth())

	t.Run("does not update when broadcast_at is NULL", func(t *testing.T) {
		t.Parallel()

		etx := mustCreateUnstartedGeneratedTx(t, orm, fromAddress, &cltest.FixtureChainID)

		var nullTime *time.Time
		assert.Equal(t, nullTime, etx.BroadcastAt)

		currTime := time.Now()
		err := orm.UpdateBroadcastAts(testutils.Context(t), currTime, []int64{etx.ID})
		require.NoError(t, err)
		etx, err = orm.FindTxWithAttempts(etx.ID)

		require.NoError(t, err)
		assert.Equal(t, nullTime, etx.BroadcastAt)
	})

	t.Run("updates when broadcast_at is non-NULL", func(t *testing.T) {
		t.Parallel()

		time1 := time.Now()
		etx := cltest.NewEthTx(fromAddress)
		etx.Sequence = new(evmtypes.Nonce)
		etx.State = txmgrcommon.TxUnconfirmed
		etx.BroadcastAt = &time1
		etx.InitialBroadcastAt = &time1
		err := orm.InsertTx(&etx)
		require.NoError(t, err)

		time2 := time.Date(2077, 8, 14, 10, 0, 0, 0, time.UTC)
		err = orm.UpdateBroadcastAts(testutils.Context(t), time2, []int64{etx.ID})
		require.NoError(t, err)
		etx, err = orm.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)
		// assert year due to time rounding after database save
		assert.Equal(t, etx.BroadcastAt.Year(), time2.Year())
	})
}

func TestORM_SetBroadcastBeforeBlockNum(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress)
	chainID := ethClient.ConfiguredChainID()

	headNum := int64(9000)
	var err error

	t.Run("saves block num to unconfirmed evm.tx_attempts without one", func(t *testing.T) {
		// Do the thing
		require.NoError(t, txStore.SetBroadcastBeforeBlockNum(testutils.Context(t), headNum, chainID))

		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 1)
		attempt := etx.TxAttempts[0]

		assert.Equal(t, int64(9000), *attempt.BroadcastBeforeBlockNum)
	})

	t.Run("does not change evm.tx_attempts that already have BroadcastBeforeBlockNum set", func(t *testing.T) {
		n := int64(42)
		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, 2)
		attempt.BroadcastBeforeBlockNum = &n
		require.NoError(t, txStore.InsertTxAttempt(&attempt))

		// Do the thing
		require.NoError(t, txStore.SetBroadcastBeforeBlockNum(testutils.Context(t), headNum, chainID))

		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 2)
		attempt = etx.TxAttempts[0]

		assert.Equal(t, int64(42), *attempt.BroadcastBeforeBlockNum)
	})

	t.Run("only updates evm.tx_attempts for the current chain", func(t *testing.T) {
		require.NoError(t, ethKeyStore.Add(fromAddress, testutils.SimulatedChainID))
		require.NoError(t, ethKeyStore.Enable(fromAddress, testutils.SimulatedChainID))
		etxThisChain := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress, cfg.EVM().ChainID())
		etxOtherChain := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress, testutils.SimulatedChainID)

		require.NoError(t, txStore.SetBroadcastBeforeBlockNum(testutils.Context(t), headNum, chainID))

		etxThisChain, err = txStore.FindTxWithAttempts(etxThisChain.ID)
		require.NoError(t, err)
		require.Len(t, etxThisChain.TxAttempts, 1)
		attempt := etxThisChain.TxAttempts[0]

		assert.Equal(t, int64(9000), *attempt.BroadcastBeforeBlockNum)

		etxOtherChain, err = txStore.FindTxWithAttempts(etxOtherChain.ID)
		require.NoError(t, err)
		require.Len(t, etxOtherChain.TxAttempts, 1)
		attempt = etxOtherChain.TxAttempts[0]

		assert.Nil(t, attempt.BroadcastBeforeBlockNum)
	})
}

func TestORM_FindTxAttemptsConfirmedMissingReceipt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx0 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, txStore, 0, 1, originalBroadcastAt, fromAddress)

	attempts, err := txStore.FindTxAttemptsConfirmedMissingReceipt(testutils.Context(t), ethClient.ConfiguredChainID())

	require.NoError(t, err)

	assert.Len(t, attempts, 1)
	assert.Len(t, etx0.TxAttempts, 1)
	assert.Equal(t, etx0.TxAttempts[0].ID, attempts[0].ID)
}

func TestORM_UpdateTxsUnconfirmed(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx0 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, txStore, 0, 1, originalBroadcastAt, fromAddress)
	assert.Equal(t, etx0.State, txmgrcommon.TxConfirmedMissingReceipt)
	require.NoError(t, txStore.UpdateTxsUnconfirmed(testutils.Context(t), []int64{etx0.ID}))

	etx0, err := txStore.FindTxWithAttempts(etx0.ID)
	require.NoError(t, err)
	assert.Equal(t, etx0.State, txmgrcommon.TxUnconfirmed)
}

func TestORM_FindTxAttemptsRequiringReceiptFetch(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx0 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, txStore, 0, 1, originalBroadcastAt, fromAddress)

	attempts, err := txStore.FindTxAttemptsRequiringReceiptFetch(testutils.Context(t), ethClient.ConfiguredChainID())
	require.NoError(t, err)
	assert.Len(t, attempts, 1)
	assert.Len(t, etx0.TxAttempts, 1)
	assert.Equal(t, etx0.TxAttempts[0].ID, attempts[0].ID)
}

func TestORM_SaveFetchedReceipts(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx0 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, txStore, 0, 1, originalBroadcastAt, fromAddress)
	require.Len(t, etx0.TxAttempts, 1)

	// create receipt associated with transaction
	txmReceipt := evmtypes.Receipt{
		TxHash:           etx0.TxAttempts[0].Hash,
		BlockHash:        utils.NewHash(),
		BlockNumber:      big.NewInt(42),
		TransactionIndex: uint(1),
	}

	err := txStore.SaveFetchedReceipts(testutils.Context(t), []*evmtypes.Receipt{&txmReceipt}, ethClient.ConfiguredChainID())

	require.NoError(t, err)
	etx0, err = txStore.FindTxWithAttempts(etx0.ID)
	require.NoError(t, err)
	require.Len(t, etx0.TxAttempts, 1)
	require.Len(t, etx0.TxAttempts[0].Receipts, 1)
	require.Equal(t, txmReceipt.BlockHash, etx0.TxAttempts[0].Receipts[0].GetBlockHash())
	require.Equal(t, txmgrcommon.TxConfirmed, etx0.State)
}

func TestORM_MarkAllConfirmedMissingReceipt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	// create transaction 0 (nonce 0) that is unconfirmed (block 7)
	etx0_blocknum := int64(7)
	etx0 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 0, fromAddress)
	etx0_attempt := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(1))
	etx0_attempt.BroadcastBeforeBlockNum = &etx0_blocknum
	require.NoError(t, txStore.InsertTxAttempt(&etx0_attempt))
	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx0.State)

	// create transaction 1 (nonce 1) that is confirmed (block 77)
	etx1 := mustInsertConfirmedEthTxBySaveFetchedReceipts(t, txStore, fromAddress, int64(1), int64(77), *ethClient.ConfiguredChainID())
	assert.Equal(t, etx1.State, txmgrcommon.TxConfirmed)

	// mark transaction 0 confirmed_missing_receipt
	err := txStore.MarkAllConfirmedMissingReceipt(testutils.Context(t), ethClient.ConfiguredChainID())
	require.NoError(t, err)
	etx0, err = txStore.FindTxWithAttempts(etx0.ID)
	require.NoError(t, err)
	assert.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx0.State)
}

func TestORM_PreloadTxes(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("loads eth transaction", func(t *testing.T) {
		// insert etx with attempt
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, int64(7), fromAddress)

		// create unloaded attempt
		unloadedAttempt := txmgr.TxAttempt{TxID: etx.ID}

		// uninitialized EthTx
		assert.Equal(t, int64(0), unloadedAttempt.Tx.ID)

		attempts := []txmgr.TxAttempt{unloadedAttempt}

		err := txStore.PreloadTxes(testutils.Context(t), attempts)
		require.NoError(t, err)

		assert.Equal(t, etx.ID, attempts[0].Tx.ID)
	})

	t.Run("returns nil when attempts slice is empty", func(t *testing.T) {
		emptyAttempts := []txmgr.TxAttempt{}
		err := txStore.PreloadTxes(testutils.Context(t), emptyAttempts)
		require.NoError(t, err)
	})
}

func TestORM_GetInProgressTxAttempts(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	// insert etx with attempt
	etx := mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, int64(7), fromAddress, txmgrtypes.TxAttemptInProgress)

	// fetch attempt
	attempts, err := txStore.GetInProgressTxAttempts(testutils.Context(t), fromAddress, ethClient.ConfiguredChainID())
	require.NoError(t, err)

	assert.Len(t, attempts, 1)
	assert.Equal(t, etx.TxAttempts[0].ID, attempts[0].ID)
}

func TestORM_FindTxesPendingCallback(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

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

	// Suspended run waiting for callback
	run1 := cltest.MustInsertPipelineRun(t, db)
	tr1 := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run1.ID)
	pgtest.MustExec(t, db, `UPDATE pipeline_runs SET state = 'suspended' WHERE id = $1`, run1.ID)
	etx1 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 3, 1, fromAddress)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET meta='{"FailOnRevert": true}'`)
	attempt1 := etx1.TxAttempts[0]
	mustInsertEthReceipt(t, txStore, head.Number-minConfirmations, head.Hash, attempt1.Hash)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr1.ID, minConfirmations, etx1.ID)

	// Callback to pipeline service completed. Should be ignored
	run2 := cltest.MustInsertPipelineRunWithStatus(t, db, 0, pipeline.RunStatusCompleted)
	tr2 := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run2.ID)
	etx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 4, 1, fromAddress)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET meta='{"FailOnRevert": false}'`)
	attempt2 := etx2.TxAttempts[0]
	mustInsertEthReceipt(t, txStore, head.Number-minConfirmations, head.Hash, attempt2.Hash)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE, callback_completed = TRUE WHERE id = $3`, &tr2.ID, minConfirmations, etx2.ID)

	// Suspended run younger than minConfirmations. Should be ignored
	run3 := cltest.MustInsertPipelineRun(t, db)
	tr3 := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run3.ID)
	pgtest.MustExec(t, db, `UPDATE pipeline_runs SET state = 'suspended' WHERE id = $1`, run3.ID)
	etx3 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 5, 1, fromAddress)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET meta='{"FailOnRevert": false}'`)
	attempt3 := etx3.TxAttempts[0]
	mustInsertEthReceipt(t, txStore, head.Number, head.Hash, attempt3.Hash)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr3.ID, minConfirmations, etx3.ID)

	// Tx not marked for callback. Should be ignore
	etx4 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 6, 1, fromAddress)
	attempt4 := etx4.TxAttempts[0]
	mustInsertEthReceipt(t, txStore, head.Number, head.Hash, attempt4.Hash)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET min_confirmations = $1 WHERE id = $2`, minConfirmations, etx4.ID)

	// Unconfirmed Tx without receipts. Should be ignored
	etx5 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 7, 1, fromAddress)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET min_confirmations = $1 WHERE id = $2`, minConfirmations, etx5.ID)

	// Search evm.txes table for tx requiring callback
	receiptsPlus, err := txStore.FindTxesPendingCallback(testutils.Context(t), head.Number, ethClient.ConfiguredChainID())
	require.NoError(t, err)
	assert.Len(t, receiptsPlus, 1)
	assert.Equal(t, tr1.ID, receiptsPlus[0].ID)
}

func Test_FindTxWithIdempotencyKey(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("returns nil if no results", func(t *testing.T) {
		idempotencyKey := "777"
		etx, err := txStore.FindTxWithIdempotencyKey(testutils.Context(t), idempotencyKey, big.NewInt(0))
		require.NoError(t, err)
		assert.Nil(t, etx)
	})

	t.Run("returns transaction if it exists", func(t *testing.T) {
		idempotencyKey := "777"
		cfg.EVM().ChainID()
		etx := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, big.NewInt(0),
			txRequestWithIdempotencyKey(idempotencyKey))
		require.Equal(t, idempotencyKey, *etx.IdempotencyKey)

		res, err := txStore.FindTxWithIdempotencyKey(testutils.Context(t), idempotencyKey, big.NewInt(0))
		require.NoError(t, err)
		assert.Equal(t, etx.Sequence, res.Sequence)
		require.Equal(t, idempotencyKey, *res.IdempotencyKey)
	})
}

func TestORM_FindTxWithSequence(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("returns nil if no results", func(t *testing.T) {
		etx, err := txStore.FindTxWithSequence(testutils.Context(t), fromAddress, evmtypes.Nonce(777))
		require.NoError(t, err)
		assert.Nil(t, etx)
	})

	t.Run("returns transaction if it exists", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 777, 1, fromAddress)
		require.Equal(t, evmtypes.Nonce(777), *etx.Sequence)

		res, err := txStore.FindTxWithSequence(testutils.Context(t), fromAddress, evmtypes.Nonce(777))
		require.NoError(t, err)
		assert.Equal(t, etx.Sequence, res.Sequence)
	})
}

func TestORM_UpdateTxForRebroadcast(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("delete all receipts for eth transaction", func(t *testing.T) {
		etx := mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 777, 1)
		etx, err := txStore.FindTxWithAttempts(etx.ID)
		assert.NoError(t, err)
		// assert attempt state
		attempt := etx.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
		// assert tx state
		assert.Equal(t, txmgrcommon.TxConfirmed, etx.State)
		// assert receipt
		assert.Len(t, etx.TxAttempts[0].Receipts, 1)

		// use exported method
		err = txStore.UpdateTxForRebroadcast(testutils.Context(t), etx, attempt)
		require.NoError(t, err)

		resultTx, err := txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Len(t, resultTx.TxAttempts, 1)
		resultTxAttempt := resultTx.TxAttempts[0]

		// assert attempt state
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, resultTxAttempt.State)
		assert.Nil(t, resultTxAttempt.BroadcastBeforeBlockNum)
		// assert tx state
		assert.Equal(t, txmgrcommon.TxUnconfirmed, resultTx.State)
		// assert receipt
		assert.Len(t, resultTxAttempt.Receipts, 0)
	})
}

func TestORM_IsTxFinalized(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	t.Run("confirmed tx not past finality_depth", func(t *testing.T) {
		confirmedAddr := cltest.MustGenerateRandomKey(t).Address
		tx := mustInsertConfirmedEthTxWithReceipt(t, txStore, confirmedAddr, 123, 1)
		finalized, err := txStore.IsTxFinalized(testutils.Context(t), 2, tx.ID, ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.False(t, finalized)
	})

	t.Run("confirmed tx past finality_depth", func(t *testing.T) {
		confirmedAddr := cltest.MustGenerateRandomKey(t).Address
		tx := mustInsertConfirmedEthTxWithReceipt(t, txStore, confirmedAddr, 123, 1)
		finalized, err := txStore.IsTxFinalized(testutils.Context(t), 10, tx.ID, ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.True(t, finalized)
	})
}

func TestORM_FindTransactionsConfirmedInBlockRange(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

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
		etx_8 := mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 700, 8)
		etx_9 := mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 777, 9)

		etxes, err := txStore.FindTransactionsConfirmedInBlockRange(testutils.Context(t), head.Number, 8, ethClient.ConfiguredChainID())
		require.NoError(t, err)
		assert.Len(t, etxes, 2)
		assert.Equal(t, etxes[0].Sequence, etx_8.Sequence)
		assert.Equal(t, etxes[1].Sequence, etx_9.Sequence)
	})
}

func TestORM_FindEarliestUnconfirmedBroadcastTime(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("no unconfirmed eth txes", func(t *testing.T) {
		broadcastAt, err := txStore.FindEarliestUnconfirmedBroadcastTime(testutils.Context(t), ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.False(t, broadcastAt.Valid)
	})

	t.Run("verify broadcast time", func(t *testing.T) {
		tx := cltest.MustInsertUnconfirmedEthTx(t, txStore, 123, fromAddress)
		broadcastAt, err := txStore.FindEarliestUnconfirmedBroadcastTime(testutils.Context(t), ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.True(t, broadcastAt.Ptr().Equal(*tx.BroadcastAt))
	})
}

func TestORM_FindEarliestUnconfirmedTxAttemptBlock(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	_, fromAddress2 := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("no earliest unconfirmed tx block", func(t *testing.T) {
		earliestBlock, err := txStore.FindEarliestUnconfirmedTxAttemptBlock(testutils.Context(t), ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.False(t, earliestBlock.Valid)
	})

	t.Run("verify earliest unconfirmed tx block", func(t *testing.T) {
		var blockNum int64 = 2
		tx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, txStore, 123, blockNum, time.Now(), fromAddress)
		_ = mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, txStore, 123, blockNum, time.Now().Add(time.Minute), fromAddress2)
		err := txStore.UpdateTxsUnconfirmed(testutils.Context(t), []int64{tx.ID})
		require.NoError(t, err)

		earliestBlock, err := txStore.FindEarliestUnconfirmedTxAttemptBlock(testutils.Context(t), ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.True(t, earliestBlock.Valid)
		require.Equal(t, blockNum, earliestBlock.Int64)
	})
}

func TestORM_SaveInsufficientEthAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	defaultDuration, err := time.ParseDuration("5s")
	require.NoError(t, err)

	t.Run("updates attempt state", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, 1, fromAddress)
		now := time.Now()

		err = txStore.SaveInsufficientFundsAttempt(testutils.Context(t), defaultDuration, &etx.TxAttempts[0], now)
		require.NoError(t, err)

		attempt, err := txStore.FindTxAttempt(etx.TxAttempts[0].Hash)
		require.NoError(t, err)
		assert.Equal(t, txmgrtypes.TxAttemptInsufficientFunds, attempt.State)
	})
}

func TestORM_SaveSentAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	defaultDuration, err := time.ParseDuration("5s")
	require.NoError(t, err)

	t.Run("updates attempt state to 'broadcast'", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, 1, fromAddress)
		require.Nil(t, etx.BroadcastAt)
		now := time.Now()

		err = txStore.SaveSentAttempt(testutils.Context(t), defaultDuration, &etx.TxAttempts[0], now)
		require.NoError(t, err)

		attempt, err := txStore.FindTxAttempt(etx.TxAttempts[0].Hash)
		require.NoError(t, err)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
	})
}

func TestORM_SaveConfirmedMissingReceiptAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	defaultDuration, err := time.ParseDuration("5s")
	require.NoError(t, err)

	t.Run("updates attempt to 'broadcast' and transaction to 'confirm_missing_receipt'", func(t *testing.T) {
		etx := mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 1, fromAddress, txmgrtypes.TxAttemptInProgress)
		now := time.Now()

		err = txStore.SaveConfirmedMissingReceiptAttempt(testutils.Context(t), defaultDuration, &etx.TxAttempts[0], now)
		require.NoError(t, err)

		etx, err := txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx.State)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, etx.TxAttempts[0].State)
	})
}

func TestORM_DeleteInProgressAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("deletes in_progress attempt", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, 1, fromAddress)
		attempt := etx.TxAttempts[0]

		err := txStore.DeleteInProgressAttempt(testutils.Context(t), etx.TxAttempts[0])
		require.NoError(t, err)

		nilResult, err := txStore.FindTxAttempt(attempt.Hash)
		assert.Nil(t, nilResult)
		require.Error(t, err)
	})
}

func TestORM_SaveInProgressAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("saves new in_progress attempt if attempt is new", func(t *testing.T) {
		etx := cltest.MustInsertUnconfirmedEthTx(t, txStore, 1, fromAddress)

		attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)
		require.Equal(t, int64(0), attempt.ID)

		err := txStore.SaveInProgressAttempt(testutils.Context(t), &attempt)
		require.NoError(t, err)

		attemptResult, err := txStore.FindTxAttempt(attempt.Hash)
		require.NoError(t, err)
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, attemptResult.State)
	})

	t.Run("updates old attempt to in_progress when insufficient_eth", func(t *testing.T) {
		etx := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 23, fromAddress)
		attempt := etx.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptInsufficientFunds, attempt.State)
		require.NotEqual(t, 0, attempt.ID)

		attempt.BroadcastBeforeBlockNum = nil
		attempt.State = txmgrtypes.TxAttemptInProgress
		err := txStore.SaveInProgressAttempt(testutils.Context(t), &attempt)

		require.NoError(t, err)
		attemptResult, err := txStore.FindTxAttempt(attempt.Hash)
		require.NoError(t, err)
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, attemptResult.State)

	})
}

func TestORM_FindTxsRequiringGasBump(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	currentBlockNum := int64(10)

	t.Run("gets txs requiring gas bump", func(t *testing.T) {
		etx := mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 1, fromAddress, txmgrtypes.TxAttemptBroadcast)
		err := txStore.SetBroadcastBeforeBlockNum(testutils.Context(t), currentBlockNum, ethClient.ConfiguredChainID())
		require.NoError(t, err)

		// this tx will require gas bump
		etx, err = txStore.FindTxWithAttempts(etx.ID)
		attempts := etx.TxAttempts
		require.NoError(t, err)
		assert.Len(t, attempts, 1)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempts[0].State)
		assert.Equal(t, currentBlockNum, *attempts[0].BroadcastBeforeBlockNum)

		// this tx will not require gas bump
		mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 2, fromAddress, txmgrtypes.TxAttemptBroadcast)
		err = txStore.SetBroadcastBeforeBlockNum(testutils.Context(t), currentBlockNum+1, ethClient.ConfiguredChainID())
		require.NoError(t, err)

		// any tx broadcast <= 10 will require gas bump
		newBlock := int64(12)
		gasBumpThreshold := int64(2)
		etxs, err := txStore.FindTxsRequiringGasBump(testutils.Context(t), fromAddress, newBlock, gasBumpThreshold, int64(0), ethClient.ConfiguredChainID())
		require.NoError(t, err)
		assert.Len(t, etxs, 1)
		assert.Equal(t, etx.ID, etxs[0].ID)
	})
}

func TestEthConfirmer_FindTxsRequiringResubmissionDueToInsufficientEth(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())

	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	_, otherAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	// Insert order is mixed up to test sorting
	etx2 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 1, fromAddress)
	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress)
	attempt3_2 := cltest.NewLegacyEthTxAttempt(t, etx3.ID)
	attempt3_2.State = txmgrtypes.TxAttemptInsufficientFunds
	attempt3_2.TxFee.Legacy = assets.NewWeiI(100)
	require.NoError(t, txStore.InsertTxAttempt(&attempt3_2))
	etx1 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, fromAddress)

	// These should never be returned
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 3, fromAddress)
	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 4, 100, fromAddress)
	mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, otherAddress)

	t.Run("returns all eth_txes with at least one attempt that is in insufficient_eth state", func(t *testing.T) {
		etxs, err := txStore.FindTxsRequiringResubmissionDueToInsufficientFunds(testutils.Context(t), fromAddress, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 3)

		assert.Equal(t, *etx1.Sequence, *etxs[0].Sequence)
		assert.Equal(t, etx1.ID, etxs[0].ID)
		assert.Equal(t, *etx2.Sequence, *etxs[1].Sequence)
		assert.Equal(t, etx2.ID, etxs[1].ID)
		assert.Equal(t, *etx3.Sequence, *etxs[2].Sequence)
		assert.Equal(t, etx3.ID, etxs[2].ID)
	})

	t.Run("does not return eth_txes with different chain ID", func(t *testing.T) {
		etxs, err := txStore.FindTxsRequiringResubmissionDueToInsufficientFunds(testutils.Context(t), fromAddress, big.NewInt(42))
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	t.Run("does not return confirmed or fatally errored eth_txes", func(t *testing.T) {
		pgtest.MustExec(t, db, `UPDATE evm.txes SET state='confirmed' WHERE id = $1`, etx1.ID)
		pgtest.MustExec(t, db, `UPDATE evm.txes SET state='fatal_error', nonce=NULL, error='foo', broadcast_at=NULL, initial_broadcast_at=NULL WHERE id = $1`, etx2.ID)

		etxs, err := txStore.FindTxsRequiringResubmissionDueToInsufficientFunds(testutils.Context(t), fromAddress, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 1)

		assert.Equal(t, *etx3.Sequence, *etxs[0].Sequence)
		assert.Equal(t, etx3.ID, etxs[0].ID)
	})
}

func TestORM_MarkOldTxesMissingReceiptAsErrored(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	// tx state should be confirmed missing receipt
	// attempt should be broadcast before cutoff time
	t.Run("successfully mark errored transactions", func(t *testing.T) {
		etx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, txStore, 1, 7, time.Now(), fromAddress)

		err := txStore.MarkOldTxesMissingReceiptAsErrored(testutils.Context(t), 10, 2, ethClient.ConfiguredChainID())
		require.NoError(t, err)

		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxFatalError, etx.State)
	})

	t.Run("successfully mark errored transactions w/ qopt passing in sql.Tx", func(t *testing.T) {
		etx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, txStore, 1, 7, time.Now(), fromAddress)
		err := txStore.MarkOldTxesMissingReceiptAsErrored(testutils.Context(t), 10, 2, ethClient.ConfiguredChainID())
		require.NoError(t, err)

		// must run other query outside of postgres transaction so changes are committed
		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxFatalError, etx.State)
	})
}

func TestORM_LoadEthTxesAttempts(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("load eth tx attempt", func(t *testing.T) {
		etx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, txStore, 1, 7, time.Now(), fromAddress)
		etx.TxAttempts = []txmgr.TxAttempt{}

		err := txStore.LoadTxesAttempts([]*txmgr.Tx{&etx})
		require.NoError(t, err)
		assert.Len(t, etx.TxAttempts, 1)
	})

	t.Run("load new attempt inserted in current postgres transaction", func(t *testing.T) {
		etx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, txStore, 3, 9, time.Now(), fromAddress)
		etx.TxAttempts = []txmgr.TxAttempt{}

		q := pg.NewQ(db, logger.Test(t), cfg.Database())

		newAttempt := cltest.NewDynamicFeeEthTxAttempt(t, etx.ID)
		var dbAttempt txmgr.DbEthTxAttempt
		dbAttempt.FromTxAttempt(&newAttempt)
		err := q.Transaction(func(tx pg.Queryer) error {
			const insertEthTxAttemptSQL = `INSERT INTO evm.tx_attempts (eth_tx_id, gas_price, signed_raw_tx, hash, broadcast_before_block_num, state, created_at, chain_specific_gas_limit, tx_type, gas_tip_cap, gas_fee_cap) VALUES (
				:eth_tx_id, :gas_price, :signed_raw_tx, :hash, :broadcast_before_block_num, :state, NOW(), :chain_specific_gas_limit, :tx_type, :gas_tip_cap, :gas_fee_cap
				) RETURNING *`
			_, err := tx.NamedExec(insertEthTxAttemptSQL, dbAttempt)
			require.NoError(t, err)

			err = txStore.LoadTxesAttempts([]*txmgr.Tx{&etx}, pg.WithQueryer(tx))
			require.NoError(t, err)
			assert.Len(t, etx.TxAttempts, 2)

			return nil
		})
		require.NoError(t, err)
		// also check after postgres transaction is committed
		etx.TxAttempts = []txmgr.TxAttempt{}
		err = txStore.LoadTxesAttempts([]*txmgr.Tx{&etx})
		require.NoError(t, err)
		assert.Len(t, etx.TxAttempts, 2)
	})
}

func TestORM_SaveReplacementInProgressAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("replace eth tx attempt", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, 123, fromAddress)
		oldAttempt := etx.TxAttempts[0]

		newAttempt := cltest.NewDynamicFeeEthTxAttempt(t, etx.ID)
		err := txStore.SaveReplacementInProgressAttempt(testutils.Context(t), oldAttempt, &newAttempt)
		require.NoError(t, err)

		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Len(t, etx.TxAttempts, 1)
		require.Equal(t, etx.TxAttempts[0].Hash, newAttempt.Hash)
	})
}

func TestORM_FindNextUnstartedTransactionFromAddress(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("cannot find unstarted tx", func(t *testing.T) {
		mustInsertInProgressEthTxWithAttempt(t, txStore, 13, fromAddress)

		resultEtx := new(txmgr.Tx)
		err := txStore.FindNextUnstartedTransactionFromAddress(testutils.Context(t), resultEtx, fromAddress, ethClient.ConfiguredChainID())
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("finds unstarted tx", func(t *testing.T) {
		mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, &cltest.FixtureChainID)
		resultEtx := new(txmgr.Tx)
		err := txStore.FindNextUnstartedTransactionFromAddress(testutils.Context(t), resultEtx, fromAddress, ethClient.ConfiguredChainID())
		require.NoError(t, err)
	})
}

func TestORM_UpdateTxFatalError(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("update successful", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, 13, fromAddress)
		etxPretendError := null.StringFrom("no more toilet paper")
		etx.Error = etxPretendError

		err := txStore.UpdateTxFatalError(testutils.Context(t), &etx)
		require.NoError(t, err)
		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Len(t, etx.TxAttempts, 0)
		assert.Equal(t, txmgrcommon.TxFatalError, etx.State)
	})
}

func TestORM_UpdateTxAttemptInProgressToBroadcast(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("update successful", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, 13, fromAddress)
		attempt := etx.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptInProgress, attempt.State)

		time1 := time.Now()
		i := int16(0)
		etx.BroadcastAt = &time1
		etx.InitialBroadcastAt = &time1
		err := txStore.UpdateTxAttemptInProgressToBroadcast(testutils.Context(t), &etx, attempt, txmgrtypes.TxAttemptBroadcast)
		require.NoError(t, err)
		// Increment sequence
		i++

		attemptResult, err := txStore.FindTxAttempt(attempt.Hash)
		require.NoError(t, err)
		require.Equal(t, attempt.Hash, attemptResult.Hash)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attemptResult.State)
		assert.Equal(t, int16(1), i)
	})
}

func TestORM_UpdateTxUnstartedToInProgress(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	q := pg.NewQ(db, logger.Test(t), cfg.Database())
	nonce := evmtypes.Nonce(123)

	t.Run("update successful", func(t *testing.T) {
		etx := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, &cltest.FixtureChainID)
		etx.Sequence = &nonce
		attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)

		err := txStore.UpdateTxUnstartedToInProgress(testutils.Context(t), &etx, &attempt)
		require.NoError(t, err)

		etx, err = txStore.FindTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxInProgress, etx.State)
		assert.Len(t, etx.TxAttempts, 1)
	})

	t.Run("update fails because tx is removed", func(t *testing.T) {
		etx := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, &cltest.FixtureChainID)
		etx.Sequence = &nonce

		attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)

		err := q.ExecQ("DELETE FROM evm.txes WHERE id = $1", etx.ID)
		require.NoError(t, err)

		err = txStore.UpdateTxUnstartedToInProgress(testutils.Context(t), &etx, &attempt)
		require.ErrorContains(t, err, "tx removed")
	})

	db = pgtest.NewSqlxDB(t)
	cfg = newTestChainScopedConfig(t)
	txStore = cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore = cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress = cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	q = pg.NewQ(db, logger.Test(t), cfg.Database())

	t.Run("update replaces abandoned tx with same hash", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, nonce, fromAddress)
		require.Len(t, etx.TxAttempts, 1)

		zero := commonconfig.MustNewDuration(time.Duration(0))
		evmCfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].Chain.Transactions.ReaperInterval = zero
			c.EVM[0].Chain.Transactions.ReaperThreshold = zero
			c.EVM[0].Chain.Transactions.ResendAfterThreshold = zero
		})

		ccfg := evmtest.NewChainScopedConfig(t, evmCfg)
		evmTxmCfg := txmgr.NewEvmTxmConfig(ccfg.EVM())
		ec := evmtest.NewEthClientMockWithDefaultChain(t)
		txMgr := txmgr.NewEvmTxm(ec.ConfiguredChainID(), evmTxmCfg, ccfg.EVM().Transactions(), nil, logger.Test(t), nil, nil,
			nil, txStore, nil, nil, nil, nil, nil)
		err := txMgr.XXXTestAbandon(fromAddress) // mark transaction as abandoned
		require.NoError(t, err)

		etx2 := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, &cltest.FixtureChainID)
		etx2.Sequence = &nonce
		attempt2 := cltest.NewLegacyEthTxAttempt(t, etx2.ID)
		attempt2.Hash = etx.TxAttempts[0].Hash

		// Even though this will initially fail due to idx_eth_tx_attempts_hash constraint, because the conflicting tx has been abandoned
		// it should succeed after removing the abandoned attempt and retrying the insert
		err = txStore.UpdateTxUnstartedToInProgress(testutils.Context(t), &etx2, &attempt2)
		require.NoError(t, err)
	})

	_, fromAddress = cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	// Same flow as previous test, but without calling txMgr.Abandon()
	t.Run("duplicate tx hash disallowed in tx_eth_attempts", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, nonce, fromAddress)
		require.Len(t, etx.TxAttempts, 1)

		etx.State = txmgrcommon.TxUnstarted

		// Should fail due to idx_eth_tx_attempt_hash constraint
		err := txStore.UpdateTxUnstartedToInProgress(testutils.Context(t), &etx, &etx.TxAttempts[0])
		assert.ErrorContains(t, err, "idx_eth_tx_attempts_hash")
		txStore = cltest.NewTestTxStore(t, db, cfg.Database()) // current txStore is poisened now, next test will need fresh one
	})
}

func TestORM_GetTxInProgress(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("gets 0 in progress eth transaction", func(t *testing.T) {
		etxResult, err := txStore.GetTxInProgress(testutils.Context(t), fromAddress)
		require.NoError(t, err)
		require.Nil(t, etxResult)
	})

	t.Run("get 1 in progress eth transaction", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, 123, fromAddress)

		etxResult, err := txStore.GetTxInProgress(testutils.Context(t), fromAddress)
		require.NoError(t, err)
		assert.Equal(t, etxResult.ID, etx.ID)
	})
}

func TestORM_GetNonFatalTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("gets 0 non finalized eth transaction", func(t *testing.T) {
		txes, err := txStore.GetNonFatalTransactions(testutils.Context(t), ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.Empty(t, txes)
	})

	t.Run("get in progress, unstarted, and unconfirmed eth transactions", func(t *testing.T) {
		inProgressTx := mustInsertInProgressEthTxWithAttempt(t, txStore, 123, fromAddress)
		unstartedTx := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, ethClient.ConfiguredChainID())

		txes, err := txStore.GetNonFatalTransactions(testutils.Context(t), ethClient.ConfiguredChainID())
		require.NoError(t, err)

		for _, tx := range txes {
			require.True(t, tx.ID == inProgressTx.ID || tx.ID == unstartedTx.ID)
		}
	})
}

func TestORM_GetTxByID(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("no transaction", func(t *testing.T) {
		tx, err := txStore.GetTxByID(testutils.Context(t), int64(0))
		require.NoError(t, err)
		require.Nil(t, tx)
	})

	t.Run("get transaction by ID", func(t *testing.T) {
		insertedTx := mustInsertInProgressEthTxWithAttempt(t, txStore, 123, fromAddress)
		tx, err := txStore.GetTxByID(testutils.Context(t), insertedTx.ID)
		require.NoError(t, err)
		require.NotNil(t, tx)
	})
}

func TestORM_GetFatalTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("gets 0 fatal eth transactions", func(t *testing.T) {
		txes, err := txStore.GetFatalTransactions(testutils.Context(t))
		require.NoError(t, err)
		require.Empty(t, txes)
	})

	t.Run("get fatal transactions", func(t *testing.T) {
		fatalTx := mustInsertFatalErrorEthTx(t, txStore, fromAddress)
		txes, err := txStore.GetFatalTransactions(testutils.Context(t))
		require.NoError(t, err)
		require.Equal(t, txes[0].ID, fatalTx.ID)
	})
}

func TestORM_HasInProgressTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("no in progress eth transaction", func(t *testing.T) {
		exists, err := txStore.HasInProgressTransaction(testutils.Context(t), fromAddress, ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("has in progress eth transaction", func(t *testing.T) {
		mustInsertInProgressEthTxWithAttempt(t, txStore, 123, fromAddress)

		exists, err := txStore.HasInProgressTransaction(testutils.Context(t), fromAddress, ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.True(t, exists)
	})
}

func TestORM_CountUnconfirmedTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, otherAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, otherAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress)

	count, err := txStore.CountUnconfirmedTransactions(testutils.Context(t), fromAddress, &cltest.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, int(count), 3)
}

func TestORM_CountTransactionsByState(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	_, fromAddress1 := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, fromAddress2 := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, fromAddress3 := cltest.MustInsertRandomKey(t, ethKeyStore)

	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress1)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress2)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress3)

	count, err := txStore.CountTransactionsByState(testutils.Context(t), txmgrcommon.TxUnconfirmed, &cltest.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, int(count), 3)
}

func TestORM_CountUnstartedTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, otherAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, &cltest.FixtureChainID)
	mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, &cltest.FixtureChainID)
	mustCreateUnstartedGeneratedTx(t, txStore, otherAddress, &cltest.FixtureChainID)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress)

	count, err := txStore.CountUnstartedTransactions(testutils.Context(t), fromAddress, &cltest.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, int(count), 2)
}

func TestORM_CheckTxQueueCapacity(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, otherAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	toAddress := testutils.NewAddress()
	encodedPayload := []byte{1, 2, 3}
	feeLimit := uint32(1000000000)
	value := big.Int(assets.NewEthValue(142))
	var maxUnconfirmedTransactions uint64 = 2

	t.Run("with no eth_txes returns nil", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(testutils.Context(t), fromAddress, maxUnconfirmedTransactions, &cltest.FixtureChainID)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		mustCreateUnstartedTx(t, txStore, otherAddress, toAddress, encodedPayload, feeLimit, value, &cltest.FixtureChainID)
	}

	t.Run("with eth_txes from another address returns nil", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(testutils.Context(t), fromAddress, maxUnconfirmedTransactions, &cltest.FixtureChainID)
		require.NoError(t, err)
	})

	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		mustInsertFatalErrorEthTx(t, txStore, otherAddress)
	}

	t.Run("ignores fatally_errored transactions", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(testutils.Context(t), fromAddress, maxUnconfirmedTransactions, &cltest.FixtureChainID)
		require.NoError(t, err)
	})

	var n int64
	mustInsertInProgressEthTxWithAttempt(t, txStore, evmtypes.Nonce(n), fromAddress)
	n++
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, n, fromAddress)
	n++

	t.Run("unconfirmed and in_progress transactions do not count", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(testutils.Context(t), fromAddress, 1, &cltest.FixtureChainID)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, n, 42, fromAddress)
		n++
	}

	t.Run("with many confirmed eth_txes from the same address returns nil", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(testutils.Context(t), fromAddress, maxUnconfirmedTransactions, &cltest.FixtureChainID)
		require.NoError(t, err)
	})

	for i := 0; i < int(maxUnconfirmedTransactions)-1; i++ {
		mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, feeLimit, value, &cltest.FixtureChainID)
	}

	t.Run("with fewer unstarted eth_txes than limit returns nil", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(testutils.Context(t), fromAddress, maxUnconfirmedTransactions, &cltest.FixtureChainID)
		require.NoError(t, err)
	})

	mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, feeLimit, value, &cltest.FixtureChainID)

	t.Run("with equal or more unstarted eth_txes than limit returns error", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(testutils.Context(t), fromAddress, maxUnconfirmedTransactions, &cltest.FixtureChainID)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("cannot create transaction; too many unstarted transactions in the queue (2/%d). WARNING: Hitting EVM.Transactions.MaxQueued", maxUnconfirmedTransactions))

		mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, feeLimit, value, &cltest.FixtureChainID)
		err = txStore.CheckTxQueueCapacity(testutils.Context(t), fromAddress, maxUnconfirmedTransactions, &cltest.FixtureChainID)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("cannot create transaction; too many unstarted transactions in the queue (3/%d). WARNING: Hitting EVM.Transactions.MaxQueued", maxUnconfirmedTransactions))
	})

	t.Run("with different chain ID ignores txes", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(testutils.Context(t), fromAddress, maxUnconfirmedTransactions, big.NewInt(42))
		require.NoError(t, err)
	})

	t.Run("disables check with 0 limit", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(testutils.Context(t), fromAddress, 0, &cltest.FixtureChainID)
		require.NoError(t, err)
	})
}

func TestORM_CreateTransaction(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, nil)
	txStore := newTxStore(t, db, cfg.Database())
	kst := cltest.NewKeyStore(t, db, cfg.Database())

	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())
	toAddress := testutils.NewAddress()
	gasLimit := uint32(1000)
	payload := []byte{1, 2, 3}

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	t.Run("with queue under capacity inserts eth_tx", func(t *testing.T) {
		subject := uuid.New()
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		strategy.On("PruneQueue", mock.Anything, mock.AnythingOfType("*txmgr.evmTxStore")).Return(int64(0), nil)
		etx, err := txStore.CreateTransaction(testutils.Context(t), txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
		}, ethClient.ConfiguredChainID())
		assert.NoError(t, err)

		assert.Greater(t, etx.ID, int64(0))
		assert.Equal(t, etx.State, txmgrcommon.TxUnstarted)
		assert.Equal(t, gasLimit, etx.FeeLimit)
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, toAddress, etx.ToAddress)
		assert.Equal(t, payload, etx.EncodedPayload)
		assert.Equal(t, big.Int(assets.NewEthValue(0)), etx.Value)
		assert.Equal(t, subject, etx.Subject.UUID)

		cltest.AssertCount(t, db, "evm.txes", 1)

		var dbEthTx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEthTx, `SELECT * FROM evm.txes ORDER BY id ASC LIMIT 1`))

		assert.Equal(t, dbEthTx.State, txmgrcommon.TxUnstarted)
		assert.Equal(t, gasLimit, dbEthTx.GasLimit)
		assert.Equal(t, fromAddress, dbEthTx.FromAddress)
		assert.Equal(t, toAddress, dbEthTx.ToAddress)
		assert.Equal(t, payload, dbEthTx.EncodedPayload)
		assert.Equal(t, assets.NewEthValue(0), dbEthTx.Value)
		assert.Equal(t, subject, dbEthTx.Subject.UUID)
	})

	t.Run("doesn't insert eth_tx if a matching tx already exists for that pipeline_task_run_id", func(t *testing.T) {
		id := uuid.New()
		txRequest := txmgr.TxRequest{
			FromAddress:       fromAddress,
			ToAddress:         testutils.NewAddress(),
			EncodedPayload:    []byte{1, 2, 3},
			FeeLimit:          21000,
			PipelineTaskRunID: &id,
			Strategy:          txmgrcommon.NewSendEveryStrategy(),
		}
		tx1, err := txStore.CreateTransaction(testutils.Context(t), txRequest, ethClient.ConfiguredChainID())
		assert.NoError(t, err)

		tx2, err := txStore.CreateTransaction(testutils.Context(t), txRequest, ethClient.ConfiguredChainID())
		assert.NoError(t, err)

		assert.Equal(t, tx1.GetID(), tx2.GetID())
	})

	t.Run("sets signal callback flag", func(t *testing.T) {
		subject := uuid.New()
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		strategy.On("PruneQueue", mock.Anything, mock.AnythingOfType("*txmgr.evmTxStore")).Return(int64(0), nil)
		etx, err := txStore.CreateTransaction(testutils.Context(t), txmgr.TxRequest{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: payload,
			FeeLimit:       gasLimit,
			Meta:           nil,
			Strategy:       strategy,
			SignalCallback: true,
		}, ethClient.ConfiguredChainID())
		assert.NoError(t, err)

		assert.Greater(t, etx.ID, int64(0))
		assert.Equal(t, fromAddress, etx.FromAddress)
		assert.Equal(t, true, etx.SignalCallback)

		cltest.AssertCount(t, db, "evm.txes", 3)

		var dbEthTx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEthTx, `SELECT * FROM evm.txes ORDER BY id DESC LIMIT 1`))

		assert.Equal(t, fromAddress, dbEthTx.FromAddress)
		assert.Equal(t, true, dbEthTx.SignalCallback)
	})
}

func TestORM_PruneUnstartedTxQueue(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	txStore := txmgr.NewTxStore(db, logger.Test(t), cfg.Database())
	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
	evmtest.NewEthClientMockWithDefaultChain(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("does not prune if queue has not exceeded capacity", func(t *testing.T) {
		subject1 := uuid.New()
		strategy1 := txmgrcommon.NewDropOldestStrategy(subject1, uint32(5), cfg.Database().DefaultQueryTimeout())
		for i := 0; i < 5; i++ {
			mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, &cltest.FixtureChainID, txRequestWithStrategy(strategy1))
		}
		AssertCountPerSubject(t, txStore, int64(5), subject1)
	})

	t.Run("prunes if queue has exceeded capacity", func(t *testing.T) {
		subject2 := uuid.New()
		strategy2 := txmgrcommon.NewDropOldestStrategy(subject2, uint32(3), cfg.Database().DefaultQueryTimeout())
		for i := 0; i < 5; i++ {
			mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, &cltest.FixtureChainID, txRequestWithStrategy(strategy2))
		}
		AssertCountPerSubject(t, txStore, int64(3), subject2)
	})
}

func AssertCountPerSubject(t *testing.T, txStore txmgr.TestEvmTxStore, expected int64, subject uuid.UUID) {
	t.Helper()
	count, err := txStore.CountTxesByStateAndSubject(testutils.Context(t), "unstarted", subject)
	require.NoError(t, err)
	require.Equal(t, int(expected), count)
}
