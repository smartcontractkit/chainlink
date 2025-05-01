package txmgr_test

import (
	"database/sql"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	txmgrcommon "github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/configtest"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
)

func TestORM_TransactionsWithAttempts(t *testing.T) {
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ctx := tests.Context(t)

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)

	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 1, from)        // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewLegacyEthTxAttempt(t, tx2.ID)
	attempt.State = txmgrtypes.TxAttemptBroadcast
	attempt.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(3)}
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))

	// tx 3 has no attempts
	mustCreateUnstartedGeneratedTx(t, txStore, from, testutils.FixtureChainID)

	var count int
	err := db.Get(&count, `SELECT count(*) FROM evm.txes`)
	require.NoError(t, err)
	require.Equal(t, 3, count)

	txs, count, err := txStore.TransactionsWithAttempts(ctx, 0, 100) // should omit tx3
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 2)
	assert.Equal(t, types.Nonce(1), *txs[0].Sequence, "transactions should be sorted by nonce")
	assert.Equal(t, types.Nonce(0), *txs[1].Sequence, "transactions should be sorted by nonce")
	assert.Len(t, txs[0].TxAttempts, 2, "all eth tx attempts are preloaded")
	assert.Len(t, txs[1].TxAttempts, 1)
	assert.Equal(t, int64(3), *txs[0].TxAttempts[0].BroadcastBeforeBlockNum, "attempts should be sorted by created_at")
	assert.Equal(t, int64(2), *txs[0].TxAttempts[1].BroadcastBeforeBlockNum, "attempts should be sorted by created_at")

	txs, count, err = txStore.TransactionsWithAttempts(ctx, 0, 1)
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 1, "limit should apply to length of results")
	assert.Equal(t, types.Nonce(1), *txs[0].Sequence, "transactions should be sorted by nonce")
}

func TestORM_Transactions(t *testing.T) {
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ctx := tests.Context(t)

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)

	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 1, from)        // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewLegacyEthTxAttempt(t, tx2.ID)
	attempt.State = txmgrtypes.TxAttemptBroadcast
	attempt.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(3)}
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))

	// tx 3 has no attempts
	mustCreateUnstartedGeneratedTx(t, txStore, from, testutils.FixtureChainID)

	var count int
	err := db.Get(&count, `SELECT count(*) FROM evm.txes`)
	require.NoError(t, err)
	require.Equal(t, 3, count)

	txs, count, err := txStore.Transactions(ctx, 0, 100)
	require.NoError(t, err)
	assert.Equal(t, 2, count, "only eth txs with attempts are counted")
	assert.Len(t, txs, 2)
	assert.Equal(t, types.Nonce(1), *txs[0].Sequence, "transactions should be sorted by nonce")
	assert.Equal(t, types.Nonce(0), *txs[1].Sequence, "transactions should be sorted by nonce")
	assert.Empty(t, txs[0].TxAttempts, "eth tx attempts should not be preloaded")
	assert.Empty(t, txs[1].TxAttempts)
}

func TestORM(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	orm := cltest.NewTestTxStore(t, db)
	_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth())
	ctx := tests.Context(t)

	var etx txmgr.Tx
	t.Run("InsertTx", func(t *testing.T) {
		etx = cltest.NewEthTx(fromAddress)
		require.NoError(t, orm.InsertTx(ctx, &etx))
		assert.Greater(t, int(etx.ID), 0)
		cltest.AssertCount(t, db, "evm.txes", 1)
	})
	var attemptL txmgr.TxAttempt
	var attemptD txmgr.TxAttempt
	t.Run("InsertTxAttempt", func(t *testing.T) {
		attemptD = cltest.NewDynamicFeeEthTxAttempt(t, etx.ID)
		require.NoError(t, orm.InsertTxAttempt(ctx, &attemptD))
		assert.Greater(t, int(attemptD.ID), 0)
		cltest.AssertCount(t, db, "evm.tx_attempts", 1)

		attemptL = cltest.NewLegacyEthTxAttempt(t, etx.ID)
		attemptL.State = txmgrtypes.TxAttemptBroadcast
		attemptL.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(42)}
		require.NoError(t, orm.InsertTxAttempt(ctx, &attemptL))
		assert.Greater(t, int(attemptL.ID), 0)
		cltest.AssertCount(t, db, "evm.tx_attempts", 2)
	})
	var r txmgr.Receipt
	t.Run("InsertReceipt", func(t *testing.T) {
		r = newEthReceipt(42, utils.NewHash(), attemptD.Hash, 0x1)
		id, err := orm.InsertReceipt(ctx, &r.Receipt)
		r.ID = id
		require.NoError(t, err)
		assert.Greater(t, int(r.ID), 0)
		cltest.AssertCount(t, db, "evm.receipts", 1)
	})
	t.Run("FindTxWithAttempts", func(t *testing.T) {
		var err error
		etx, err = orm.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 2)
		assert.Equal(t, etx.TxAttempts[0].ID, attemptD.ID)
		assert.Equal(t, etx.TxAttempts[1].ID, attemptL.ID)
		require.Len(t, etx.TxAttempts[0].Receipts, 1)
		require.Empty(t, etx.TxAttempts[1].Receipts)
		assert.Equal(t, r.BlockHash, etx.TxAttempts[0].Receipts[0].GetBlockHash())
	})
	t.Run("FindTxByHash", func(t *testing.T) {
		foundEtx, err := orm.FindTxByHash(ctx, attemptD.Hash)
		require.NoError(t, err)
		assert.Equal(t, etx.ID, foundEtx.ID)
		assert.Equal(t, etx.ChainID, foundEtx.ChainID)
	})
	t.Run("FindTxAttemptsByTxIDs", func(t *testing.T) {
		attempts, err := orm.FindTxAttemptsByTxIDs(ctx, []int64{etx.ID})
		require.NoError(t, err)
		require.Len(t, attempts, 2)
		assert.Equal(t, etx.TxAttempts[0].ID, attemptD.ID)
		assert.Equal(t, etx.TxAttempts[1].ID, attemptL.ID)
		require.Len(t, etx.TxAttempts[0].Receipts, 1)
		require.Empty(t, etx.TxAttempts[1].Receipts)
		assert.Equal(t, r.BlockHash, etx.TxAttempts[0].Receipts[0].GetBlockHash())
	})
}

func TestORM_FindTxAttemptConfirmedByTxIDs(t *testing.T) {
	db := testutils.NewSqlxDB(t)
	orm := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ctx := tests.Context(t)

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)

	tx1 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, orm, 0, 1, from) // tx1
	tx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, orm, 1, 2, from) // tx2

	// add 2nd attempt to tx2
	blockNum := int64(3)
	attempt := cltest.NewLegacyEthTxAttempt(t, tx2.ID)
	attempt.State = txmgrtypes.TxAttemptBroadcast
	attempt.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(3)}
	attempt.BroadcastBeforeBlockNum = &blockNum
	require.NoError(t, orm.InsertTxAttempt(ctx, &attempt))

	// add receipt for the second attempt
	r := newEthReceipt(4, utils.NewHash(), attempt.Hash, 0x1)
	_, err := orm.InsertReceipt(ctx, &r.Receipt)
	require.NoError(t, err)
	// tx 3 has no attempts
	mustCreateUnstartedGeneratedTx(t, orm, from, testutils.FixtureChainID)

	cltest.MustInsertUnconfirmedEthTx(t, orm, 3, from)                           // tx4
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, orm, 4, from) // tx5

	var count int
	err = db.Get(&count, `SELECT count(*) FROM evm.txes`)
	require.NoError(t, err)
	require.Equal(t, 5, count)

	err = db.Get(&count, `SELECT count(*) FROM evm.tx_attempts`)
	require.NoError(t, err)
	require.Equal(t, 4, count)

	confirmedAttempts, err := orm.FindTxAttemptConfirmedByTxIDs(ctx, []int64{tx1.ID, tx2.ID}) // should omit tx3
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

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ctx := tests.Context(t)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	t.Run("returns nothing if there are no transactions", func(t *testing.T) {
		olderThan := time.Now()
		attempts, err := txStore.FindTxAttemptsRequiringResend(tests.Context(t), olderThan, 10, testutils.FixtureChainID, fromAddress)
		require.NoError(t, err)
		assert.Empty(t, attempts)
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
	attempt1_2.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(10)}
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt1_2))

	attempt3_2 := newInProgressLegacyEthTxAttempt(t, etxs[2].ID)
	attempt3_2.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(10)}
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt3_2))

	attempt4_2 := cltest.NewDynamicFeeEthTxAttempt(t, etxs[3].ID)
	attempt4_2.TxFee.GasTipCap = assets.NewWeiI(10)
	attempt4_2.TxFee.GasFeeCap = assets.NewWeiI(20)
	attempt4_2.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt4_2))
	attempt4_4 := cltest.NewDynamicFeeEthTxAttempt(t, etxs[3].ID)
	attempt4_4.TxFee.GasTipCap = assets.NewWeiI(30)
	attempt4_4.TxFee.GasFeeCap = assets.NewWeiI(40)
	attempt4_4.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt4_4))
	attempt4_3 := cltest.NewDynamicFeeEthTxAttempt(t, etxs[3].ID)
	attempt4_3.TxFee.GasTipCap = assets.NewWeiI(20)
	attempt4_3.TxFee.GasFeeCap = assets.NewWeiI(30)
	attempt4_3.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt4_3))

	t.Run("returns nothing if there are transactions from a different key", func(t *testing.T) {
		olderThan := time.Now()
		attempts, err := txStore.FindTxAttemptsRequiringResend(tests.Context(t), olderThan, 10, testutils.FixtureChainID, utils.RandomAddress())
		require.NoError(t, err)
		assert.Empty(t, attempts)
	})

	t.Run("returns the highest price attempt for each transaction that was last broadcast before or on the given time", func(t *testing.T) {
		olderThan := time.Unix(1616509200, 0)
		attempts, err := txStore.FindTxAttemptsRequiringResend(tests.Context(t), olderThan, 0, testutils.FixtureChainID, fromAddress)
		require.NoError(t, err)
		assert.Len(t, attempts, 2)
		assert.Equal(t, attempt1_2.ID, attempts[0].ID)
		assert.Equal(t, etxs[1].TxAttempts[0].ID, attempts[1].ID)
	})

	t.Run("returns the highest price attempt for EIP-1559 transactions", func(t *testing.T) {
		olderThan := time.Unix(1616509400, 0)
		attempts, err := txStore.FindTxAttemptsRequiringResend(tests.Context(t), olderThan, 0, testutils.FixtureChainID, fromAddress)
		require.NoError(t, err)
		assert.Len(t, attempts, 4)
		assert.Equal(t, attempt4_4.ID, attempts[3].ID)
	})

	t.Run("applies limit", func(t *testing.T) {
		olderThan := time.Unix(1616509200, 0)
		attempts, err := txStore.FindTxAttemptsRequiringResend(tests.Context(t), olderThan, 1, testutils.FixtureChainID, fromAddress)
		require.NoError(t, err)
		assert.Len(t, attempts, 1)
		assert.Equal(t, attempt1_2.ID, attempts[0].ID)
	})
}

func TestORM_UpdateBroadcastAts(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)
	orm := cltest.NewTestTxStore(t, db)
	_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth())

	t.Run("does not update when broadcast_at is NULL", func(t *testing.T) {
		t.Parallel()
		ctx := tests.Context(t)
		etx := mustCreateUnstartedGeneratedTx(t, orm, fromAddress, testutils.FixtureChainID)

		var nullTime *time.Time
		assert.Equal(t, nullTime, etx.BroadcastAt)

		currTime := time.Now()
		err := orm.UpdateBroadcastAts(tests.Context(t), currTime, []int64{etx.ID})
		require.NoError(t, err)
		etx, err = orm.FindTxWithAttempts(ctx, etx.ID)

		require.NoError(t, err)
		assert.Equal(t, nullTime, etx.BroadcastAt)
	})

	t.Run("updates when broadcast_at is non-NULL", func(t *testing.T) {
		t.Parallel()

		ctx := tests.Context(t)
		time1 := time.Now()
		etx := cltest.NewEthTx(fromAddress)
		etx.Sequence = new(types.Nonce)
		etx.State = txmgrcommon.TxUnconfirmed
		etx.BroadcastAt = &time1
		etx.InitialBroadcastAt = &time1
		err := orm.InsertTx(ctx, &etx)
		require.NoError(t, err)

		time2 := time.Date(2077, 8, 14, 10, 0, 0, 0, time.UTC)
		err = orm.UpdateBroadcastAts(ctx, time2, []int64{etx.ID})
		require.NoError(t, err)
		etx, err = orm.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		// assert year due to time rounding after database save
		assert.Equal(t, etx.BroadcastAt.Year(), time2.Year())
	})
}

func TestORM_SetBroadcastBeforeBlockNum(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	cfg := configtest.NewChainScopedConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress)
	chainID := ethClient.ConfiguredChainID()
	ctx := tests.Context(t)

	headNum := int64(9000)
	var err error

	t.Run("saves block num to unconfirmed evm.tx_attempts without one", func(t *testing.T) {
		// Do the thing
		require.NoError(t, txStore.SetBroadcastBeforeBlockNum(tests.Context(t), headNum, chainID))

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 1)
		attempt := etx.TxAttempts[0]

		assert.Equal(t, int64(9000), *attempt.BroadcastBeforeBlockNum)
	})

	t.Run("does not change evm.tx_attempts that already have BroadcastBeforeBlockNum set", func(t *testing.T) {
		n := int64(42)
		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, 2)
		attempt.BroadcastBeforeBlockNum = &n
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))

		// Do the thing
		require.NoError(t, txStore.SetBroadcastBeforeBlockNum(tests.Context(t), headNum, chainID))

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 2)
		attempt = etx.TxAttempts[0]

		assert.Equal(t, int64(42), *attempt.BroadcastBeforeBlockNum)
	})

	t.Run("only updates evm.tx_attempts for the current chain", func(t *testing.T) {
		require.NoError(t, ethKeyStore.Add(tests.Context(t), fromAddress, testutils.SimulatedChainID))
		require.NoError(t, ethKeyStore.Enable(tests.Context(t), fromAddress, testutils.SimulatedChainID))
		etxThisChain := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress, cfg.EVM().ChainID())
		etxOtherChain := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress, testutils.SimulatedChainID)

		require.NoError(t, txStore.SetBroadcastBeforeBlockNum(tests.Context(t), headNum, chainID))

		etxThisChain, err = txStore.FindTxWithAttempts(ctx, etxThisChain.ID)
		require.NoError(t, err)
		require.Len(t, etxThisChain.TxAttempts, 1)
		attempt := etxThisChain.TxAttempts[0]

		assert.Equal(t, int64(9000), *attempt.BroadcastBeforeBlockNum)

		etxOtherChain, err = txStore.FindTxWithAttempts(ctx, etxOtherChain.ID)
		require.NoError(t, err)
		require.Len(t, etxOtherChain.TxAttempts, 1)
		attempt = etxOtherChain.TxAttempts[0]

		assert.Nil(t, attempt.BroadcastBeforeBlockNum)
	})
}

func TestORM_UpdateTxConfirmed(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	etx0 := mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 0, fromAddress, txmgrtypes.TxAttemptBroadcast)
	etx1 := mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 1, fromAddress, txmgrtypes.TxAttemptInProgress)
	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx0.State)
	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx1.State)
	require.NoError(t, txStore.UpdateTxConfirmed(tests.Context(t), []int64{etx0.ID, etx1.ID}))

	var err error
	etx0, err = txStore.FindTxWithAttempts(ctx, etx0.ID)
	require.NoError(t, err)
	assert.Equal(t, txmgrcommon.TxConfirmed, etx0.State)
	assert.Len(t, etx0.TxAttempts, 1)
	assert.Equal(t, txmgrtypes.TxAttemptBroadcast, etx0.TxAttempts[0].State)
	etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
	require.NoError(t, err)
	assert.Equal(t, txmgrcommon.TxConfirmed, etx1.State)
	assert.Len(t, etx1.TxAttempts, 1)
	assert.Equal(t, txmgrtypes.TxAttemptBroadcast, etx1.TxAttempts[0].State)
}

func TestORM_SaveFetchedReceipts(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	ctx := tests.Context(t)

	tx1 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 100, fromAddress)
	require.Len(t, tx1.TxAttempts, 1)

	tx2 := mustInsertTerminallyStuckTxWithAttempt(t, txStore, fromAddress, 1, 100)
	require.Len(t, tx2.TxAttempts, 1)

	// create receipts associated with transactions
	txmReceipt1 := types.Receipt{
		TxHash:           tx1.TxAttempts[0].Hash,
		BlockHash:        utils.NewHash(),
		BlockNumber:      big.NewInt(42),
		TransactionIndex: uint(1),
	}
	txmReceipt2 := types.Receipt{
		TxHash:           tx2.TxAttempts[0].Hash,
		BlockHash:        utils.NewHash(),
		BlockNumber:      big.NewInt(42),
		TransactionIndex: uint(1),
	}

	err := txStore.SaveFetchedReceipts(tests.Context(t), []*types.Receipt{&txmReceipt1, &txmReceipt2})
	require.NoError(t, err)

	tx1, err = txStore.FindTxWithAttempts(ctx, tx1.ID)
	require.NoError(t, err)
	require.Len(t, tx1.TxAttempts, 1)
	require.Len(t, tx1.TxAttempts[0].Receipts, 1)
	require.Equal(t, txmReceipt1.BlockHash, tx1.TxAttempts[0].Receipts[0].GetBlockHash())
	require.Equal(t, txmgrcommon.TxConfirmed, tx1.State)

	tx2, err = txStore.FindTxWithAttempts(ctx, tx2.ID)
	require.NoError(t, err)
	require.Len(t, tx2.TxAttempts, 1)
	require.Len(t, tx2.TxAttempts[0].Receipts, 1)
	require.Equal(t, txmReceipt2.BlockHash, tx2.TxAttempts[0].Receipts[0].GetBlockHash())
	require.Equal(t, txmgrcommon.TxFatalError, tx2.State)
}

func TestORM_PreloadTxes(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("loads eth transaction", func(t *testing.T) {
		// insert etx with attempt
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, int64(7), fromAddress)

		// create unloaded attempt
		unloadedAttempt := txmgr.TxAttempt{TxID: etx.ID}

		// uninitialized EthTx
		assert.Equal(t, int64(0), unloadedAttempt.Tx.ID)

		attempts := []txmgr.TxAttempt{unloadedAttempt}

		err := txStore.PreloadTxes(tests.Context(t), attempts)
		require.NoError(t, err)

		assert.Equal(t, etx.ID, attempts[0].Tx.ID)
	})

	t.Run("returns nil when attempts slice is empty", func(t *testing.T) {
		emptyAttempts := []txmgr.TxAttempt{}
		err := txStore.PreloadTxes(tests.Context(t), emptyAttempts)
		require.NoError(t, err)
	})
}

func TestORM_GetInProgressTxAttempts(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	// insert etx with attempt
	etx := mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, int64(7), fromAddress, txmgrtypes.TxAttemptInProgress)

	// fetch attempt
	attempts, err := txStore.GetInProgressTxAttempts(tests.Context(t), fromAddress, ethClient.ConfiguredChainID())
	require.NoError(t, err)

	assert.Len(t, attempts, 1)
	assert.Equal(t, etx.TxAttempts[0].ID, attempts[0].ID)
}

func TestORM_FindTxesPendingCallback(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	testutils.MustExec(t, db, `SET CONSTRAINTS fk_pipeline_runs_pruning_key DEFERRED`)
	testutils.MustExec(t, db, `SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`)

	h8 := &types.Head{
		Number: 8,
		Hash:   testutils.NewHash(),
	}
	h9 := &types.Head{
		Hash:   testutils.NewHash(),
		Number: 9,
	}
	h9.Parent.Store(h8)
	head := types.Head{
		Hash:   testutils.NewHash(),
		Number: 10,
	}

	minConfirmations := int64(2)

	// Suspended run waiting for callback
	run1 := cltest.MustInsertPipelineRun(t, db)
	tr1 := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run1.ID)
	testutils.MustExec(t, db, `UPDATE pipeline_runs SET state = 'suspended' WHERE id = $1`, run1.ID)
	etx1 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 3, 1, fromAddress)
	testutils.MustExec(t, db, `UPDATE evm.txes SET meta='{"FailOnRevert": true}'`)
	attempt1 := etx1.TxAttempts[0]
	etxBlockNum := mustInsertEthReceipt(t, txStore, head.Number-minConfirmations, head.Hash, attempt1.Hash).BlockNumber
	testutils.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr1.ID, minConfirmations, etx1.ID)

	// Callback to pipeline service completed. Should be ignored
	run2 := cltest.MustInsertPipelineRunWithStatus(t, db, 0, "completed", 0)
	tr2 := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run2.ID)
	etx2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 4, 1, fromAddress)
	testutils.MustExec(t, db, `UPDATE evm.txes SET meta='{"FailOnRevert": false}'`)
	attempt2 := etx2.TxAttempts[0]
	mustInsertEthReceipt(t, txStore, head.Number-minConfirmations, head.Hash, attempt2.Hash)
	testutils.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE, callback_completed = TRUE WHERE id = $3`, &tr2.ID, minConfirmations, etx2.ID)

	// Suspended run younger than minConfirmations. Should be ignored
	run3 := cltest.MustInsertPipelineRun(t, db)
	tr3 := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run3.ID)
	testutils.MustExec(t, db, `UPDATE pipeline_runs SET state = 'suspended' WHERE id = $1`, run3.ID)
	etx3 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 5, 1, fromAddress)
	testutils.MustExec(t, db, `UPDATE evm.txes SET meta='{"FailOnRevert": false}'`)
	attempt3 := etx3.TxAttempts[0]
	mustInsertEthReceipt(t, txStore, head.Number, head.Hash, attempt3.Hash)
	testutils.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr3.ID, minConfirmations, etx3.ID)

	// Tx not marked for callback. Should be ignore
	etx4 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 6, 1, fromAddress)
	attempt4 := etx4.TxAttempts[0]
	mustInsertEthReceipt(t, txStore, head.Number, head.Hash, attempt4.Hash)
	testutils.MustExec(t, db, `UPDATE evm.txes SET min_confirmations = $1 WHERE id = $2`, minConfirmations, etx4.ID)

	// Unconfirmed Tx without receipts. Should be ignored
	etx5 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 7, 1, fromAddress)
	testutils.MustExec(t, db, `UPDATE evm.txes SET min_confirmations = $1 WHERE id = $2`, minConfirmations, etx5.ID)

	// Search evm.txes table for tx requiring callback
	receiptsPlus, err := txStore.FindTxesPendingCallback(tests.Context(t), head.Number, 0, ethClient.ConfiguredChainID())
	require.NoError(t, err)
	if assert.Len(t, receiptsPlus, 1) {
		assert.Equal(t, tr1.ID, receiptsPlus[0].ID)
	}

	// Clear min_confirmations
	testutils.MustExec(t, db, `UPDATE evm.txes SET min_confirmations = NULL WHERE id = $1`, etx1.ID)

	// Search evm.txes table for tx requiring callback
	receiptsPlus, err = txStore.FindTxesPendingCallback(tests.Context(t), head.Number, 0, ethClient.ConfiguredChainID())
	require.NoError(t, err)
	assert.Empty(t, receiptsPlus)

	// Search evm.txes table for tx requiring callback, with block 1 finalized
	receiptsPlus, err = txStore.FindTxesPendingCallback(tests.Context(t), head.Number, etxBlockNum, ethClient.ConfiguredChainID())
	require.NoError(t, err)
	if assert.Len(t, receiptsPlus, 1) {
		assert.Equal(t, tr1.ID, receiptsPlus[0].ID)
	}
}

func Test_FindTxWithIdempotencyKey(t *testing.T) {
	t.Parallel()
	db := testutils.NewSqlxDB(t)
	cfg := configtest.NewChainScopedConfig(t, nil)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("returns nil error if no results", func(t *testing.T) {
		idempotencyKey := "777"
		etx, err := txStore.FindTxWithIdempotencyKey(tests.Context(t), idempotencyKey, big.NewInt(0))
		require.NoError(t, err)
		assert.Nil(t, etx)
	})

	t.Run("returns transaction if it exists", func(t *testing.T) {
		idempotencyKey := "777"
		cfg.EVM().ChainID()
		etx := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, big.NewInt(0),
			txRequestWithIdempotencyKey(idempotencyKey))
		require.Equal(t, idempotencyKey, *etx.IdempotencyKey)

		res, err := txStore.FindTxWithIdempotencyKey(tests.Context(t), idempotencyKey, big.NewInt(0))
		require.NoError(t, err)
		assert.Equal(t, etx.Sequence, res.Sequence)
		require.Equal(t, idempotencyKey, *res.IdempotencyKey)
	})
}

func Test_FindReceiptWithIdempotencyKey(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	ctx := t.Context()

	idempotencyKey := "654"
	t.Run("returns nil error if no results", func(t *testing.T) {
		r, err := txStore.FindReceiptWithIdempotencyKey(tests.Context(t), idempotencyKey, big.NewInt(0))
		require.NoError(t, err)
		assert.Nil(t, r)
	})

	t.Run("returns receipt if it exists", func(t *testing.T) {
		var etx txmgr.Tx
		// insert tx
		etx = cltest.NewEthTx(fromAddress)
		etx.IdempotencyKey = &idempotencyKey
		require.NoError(t, txStore.InsertTx(ctx, &etx))
		assert.Greater(t, int(etx.ID), 0)
		cltest.AssertCount(t, db, "evm.txes", 1)

		// insert attempt
		var attemptL txmgr.TxAttempt
		var attemptD txmgr.TxAttempt
		attemptD = cltest.NewDynamicFeeEthTxAttempt(t, etx.ID)
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attemptD))
		assert.Greater(t, int(attemptD.ID), 0)
		cltest.AssertCount(t, db, "evm.tx_attempts", 1)

		attemptL = cltest.NewLegacyEthTxAttempt(t, etx.ID)
		attemptL.State = txmgrtypes.TxAttemptBroadcast
		attemptL.TxFee = gas.EvmFee{GasPrice: assets.NewWeiI(42)}
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attemptL))
		assert.Greater(t, int(attemptL.ID), 0)
		cltest.AssertCount(t, db, "evm.tx_attempts", 2)

		// insert receipt
		var r txmgr.Receipt
		r = newEthReceipt(42, utils.NewHash(), attemptD.Hash, 0x1)
		id, err := txStore.InsertReceipt(ctx, &r.Receipt)
		r.ID = id
		require.NoError(t, err)
		assert.Greater(t, int(r.ID), 0)
		cltest.AssertCount(t, db, "evm.receipts", 1)

		res, err := txStore.FindReceiptWithIdempotencyKey(ctx, idempotencyKey, etx.ChainID)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, r.Receipt.GasUsed, res.GetFeeUsed())
		require.Equal(t, r.Receipt.EffectiveGasPrice, res.GetEffectiveGasPrice())
	})
}

func TestORM_FindTxWithSequence(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("returns nil if no results", func(t *testing.T) {
		etx, err := txStore.FindTxWithSequence(tests.Context(t), fromAddress, types.Nonce(777))
		require.NoError(t, err)
		assert.Nil(t, etx)
	})

	t.Run("returns transaction if it exists", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 777, 1, fromAddress)
		require.Equal(t, types.Nonce(777), *etx.Sequence)

		res, err := txStore.FindTxWithSequence(tests.Context(t), fromAddress, types.Nonce(777))
		require.NoError(t, err)
		assert.Equal(t, etx.Sequence, res.Sequence)
	})
}

func TestORM_UpdateTxForRebroadcast(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	ctx := tests.Context(t)

	t.Run("marks confirmed tx as unconfirmed, marks latest attempt as in-progress, deletes receipt", func(t *testing.T) {
		etx := mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 777, 1)
		etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
		assert.NoError(t, err)
		// assert attempt state
		attempt := etx.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
		// assert tx state
		assert.Equal(t, txmgrcommon.TxConfirmed, etx.State)
		// assert receipt
		assert.Len(t, etx.TxAttempts[0].Receipts, 1)

		// use exported method
		err = txStore.UpdateTxsForRebroadcast(tests.Context(t), []int64{etx.ID}, []int64{attempt.ID})
		require.NoError(t, err)

		resultTx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Len(t, resultTx.TxAttempts, 1)
		resultTxAttempt := resultTx.TxAttempts[0]

		// assert attempt state
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, resultTxAttempt.State)
		assert.Nil(t, resultTxAttempt.BroadcastBeforeBlockNum)
		// assert tx state
		assert.Equal(t, txmgrcommon.TxUnconfirmed, resultTx.State)
		// assert receipt
		assert.Empty(t, resultTxAttempt.Receipts)
	})

	t.Run("marks confirmed tx as unconfirmed, clears error, marks latest attempt as in-progress, deletes receipt", func(t *testing.T) {
		blockNum := int64(100)
		etx := mustInsertTerminallyStuckTxWithAttempt(t, txStore, fromAddress, 1, blockNum)
		mustInsertEthReceipt(t, txStore, blockNum, utils.NewHash(), etx.TxAttempts[0].Hash)
		etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		// assert attempt state
		attempt := etx.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
		// assert tx state
		assert.Equal(t, txmgrcommon.TxFatalError, etx.State)
		// assert receipt
		assert.Len(t, etx.TxAttempts[0].Receipts, 1)

		// use exported method
		err = txStore.UpdateTxsForRebroadcast(tests.Context(t), []int64{etx.ID}, []int64{attempt.ID})
		require.NoError(t, err)

		resultTx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Len(t, resultTx.TxAttempts, 1)
		resultTxAttempt := resultTx.TxAttempts[0]

		// assert attempt state
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, resultTxAttempt.State)
		assert.Nil(t, resultTxAttempt.BroadcastBeforeBlockNum)
		// assert tx state
		assert.Equal(t, txmgrcommon.TxUnconfirmed, resultTx.State)
		// assert receipt
		assert.Empty(t, resultTxAttempt.Receipts)
	})
}

func TestORM_FindEarliestUnconfirmedBroadcastTime(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("no unconfirmed eth txes", func(t *testing.T) {
		broadcastAt, err := txStore.FindEarliestUnconfirmedBroadcastTime(tests.Context(t), ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.False(t, broadcastAt.Valid)
	})

	t.Run("verify broadcast time", func(t *testing.T) {
		tx := cltest.MustInsertUnconfirmedEthTx(t, txStore, 123, fromAddress)
		broadcastAt, err := txStore.FindEarliestUnconfirmedBroadcastTime(tests.Context(t), ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.True(t, broadcastAt.Ptr().Equal(*tx.BroadcastAt))
	})
}

func TestORM_FindEarliestUnconfirmedTxAttemptBlock(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	_, fromAddress2 := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("no earliest unconfirmed tx block", func(t *testing.T) {
		earliestBlock, err := txStore.FindEarliestUnconfirmedTxAttemptBlock(tests.Context(t), ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.False(t, earliestBlock.Valid)
	})

	t.Run("verify earliest unconfirmed tx block", func(t *testing.T) {
		var blockNum int64 = 2
		tx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, txStore, 123, blockNum, time.Now(), fromAddress)
		_ = mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, txStore, 123, blockNum, time.Now().Add(time.Minute), fromAddress2)
		err := txStore.UpdateTxsUnconfirmed(tests.Context(t), []int64{tx.ID})
		require.NoError(t, err)

		earliestBlock, err := txStore.FindEarliestUnconfirmedTxAttemptBlock(tests.Context(t), ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.True(t, earliestBlock.Valid)
		require.Equal(t, blockNum, earliestBlock.Int64)
	})
}

func TestORM_SaveInsufficientEthAttempt(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	defaultDuration, err := time.ParseDuration("5s")
	require.NoError(t, err)

	t.Run("updates attempt state", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, 1, fromAddress)
		now := time.Now()

		err = txStore.SaveInsufficientFundsAttempt(tests.Context(t), defaultDuration, &etx.TxAttempts[0], now)
		require.NoError(t, err)

		attempt, err := txStore.FindTxAttempt(ctx, etx.TxAttempts[0].Hash)
		require.NoError(t, err)
		assert.Equal(t, txmgrtypes.TxAttemptInsufficientFunds, attempt.State)
	})
}

func TestORM_SaveSentAttempt(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	defaultDuration, err := time.ParseDuration("5s")
	require.NoError(t, err)

	t.Run("updates attempt state to 'broadcast'", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, 1, fromAddress)
		require.Nil(t, etx.BroadcastAt)
		now := time.Now()

		err = txStore.SaveSentAttempt(tests.Context(t), defaultDuration, &etx.TxAttempts[0], now)
		require.NoError(t, err)

		attempt, err := txStore.FindTxAttempt(ctx, etx.TxAttempts[0].Hash)
		require.NoError(t, err)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
	})
}

func TestORM_SaveConfirmedAttempt(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	defaultDuration, err := time.ParseDuration("5s")
	require.NoError(t, err)

	t.Run("updates attempt to 'broadcast' and transaction to 'confirm_missing_receipt'", func(t *testing.T) {
		etx := mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 1, fromAddress, txmgrtypes.TxAttemptInProgress)
		now := time.Now()

		err = txStore.SaveConfirmedAttempt(tests.Context(t), defaultDuration, &etx.TxAttempts[0], now)
		require.NoError(t, err)

		etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxConfirmed, etx.State)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, etx.TxAttempts[0].State)
	})
}

func TestORM_DeleteInProgressAttempt(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("deletes in_progress attempt", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, 1, fromAddress)
		attempt := etx.TxAttempts[0]

		err := txStore.DeleteInProgressAttempt(tests.Context(t), etx.TxAttempts[0])
		require.NoError(t, err)

		nilResult, err := txStore.FindTxAttempt(ctx, attempt.Hash)
		assert.Nil(t, nilResult)
		require.Error(t, err)
	})
}

func TestORM_SaveInProgressAttempt(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("saves new in_progress attempt if attempt is new", func(t *testing.T) {
		etx := cltest.MustInsertUnconfirmedEthTx(t, txStore, 1, fromAddress)

		attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)
		require.Equal(t, int64(0), attempt.ID)

		err := txStore.SaveInProgressAttempt(tests.Context(t), &attempt)
		require.NoError(t, err)

		attemptResult, err := txStore.FindTxAttempt(ctx, attempt.Hash)
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
		err := txStore.SaveInProgressAttempt(tests.Context(t), &attempt)

		require.NoError(t, err)
		attemptResult, err := txStore.FindTxAttempt(ctx, attempt.Hash)
		require.NoError(t, err)
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, attemptResult.State)
	})
}

func TestORM_FindTxsRequiringGasBump(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	currentBlockNum := int64(10)

	t.Run("gets txs requiring gas bump", func(t *testing.T) {
		etx := mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 1, fromAddress, txmgrtypes.TxAttemptBroadcast)
		err := txStore.SetBroadcastBeforeBlockNum(tests.Context(t), currentBlockNum, ethClient.ConfiguredChainID())
		require.NoError(t, err)

		// this tx will require gas bump
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		attempts := etx.TxAttempts
		require.NoError(t, err)
		assert.Len(t, attempts, 1)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempts[0].State)
		assert.Equal(t, currentBlockNum, *attempts[0].BroadcastBeforeBlockNum)

		// this tx will not require gas bump
		mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 2, fromAddress, txmgrtypes.TxAttemptBroadcast)
		err = txStore.SetBroadcastBeforeBlockNum(tests.Context(t), currentBlockNum+1, ethClient.ConfiguredChainID())
		require.NoError(t, err)

		// any tx broadcast <= 10 will require gas bump
		newBlock := int64(12)
		gasBumpThreshold := int64(2)
		etxs, err := txStore.FindTxsRequiringGasBump(tests.Context(t), fromAddress, newBlock, gasBumpThreshold, int64(0), ethClient.ConfiguredChainID())
		require.NoError(t, err)
		assert.Len(t, etxs, 1)
		assert.Equal(t, etx.ID, etxs[0].ID)
	})
}

func TestEthConfirmer_FindTxsRequiringResubmissionDueToInsufficientEth(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ctx := tests.Context(t)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	_, otherAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	// Insert order is mixed up to test sorting
	etx2 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 1, fromAddress)
	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress)
	attempt3_2 := cltest.NewLegacyEthTxAttempt(t, etx3.ID)
	attempt3_2.State = txmgrtypes.TxAttemptInsufficientFunds
	attempt3_2.TxFee.GasPrice = assets.NewWeiI(100)
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt3_2))
	etx1 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, fromAddress)

	// These should never be returned
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 3, fromAddress)
	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 4, 100, fromAddress)
	mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, 0, otherAddress)

	t.Run("returns all eth_txes with at least one attempt that is in insufficient_eth state", func(t *testing.T) {
		etxs, err := txStore.FindTxsRequiringResubmissionDueToInsufficientFunds(tests.Context(t), fromAddress, testutils.FixtureChainID)
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
		etxs, err := txStore.FindTxsRequiringResubmissionDueToInsufficientFunds(tests.Context(t), fromAddress, big.NewInt(42))
		require.NoError(t, err)

		assert.Empty(t, etxs)
	})

	t.Run("does not return confirmed or fatally errored eth_txes", func(t *testing.T) {
		testutils.MustExec(t, db, `UPDATE evm.txes SET state='confirmed' WHERE id = $1`, etx1.ID)
		testutils.MustExec(t, db, `UPDATE evm.txes SET state='fatal_error', nonce=NULL, error='foo', broadcast_at=NULL, initial_broadcast_at=NULL WHERE id = $1`, etx2.ID)

		etxs, err := txStore.FindTxsRequiringResubmissionDueToInsufficientFunds(tests.Context(t), fromAddress, testutils.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 1)

		assert.Equal(t, *etx3.Sequence, *etxs[0].Sequence)
		assert.Equal(t, etx3.ID, etxs[0].ID)
	})
}

func TestORM_LoadEthTxesAttempts(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("load eth tx attempt", func(t *testing.T) {
		etx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, txStore, 1, 7, time.Now(), fromAddress)
		etx.TxAttempts = []txmgr.TxAttempt{}

		err := txStore.LoadTxesAttempts(ctx, []*txmgr.Tx{&etx})
		require.NoError(t, err)
		assert.Len(t, etx.TxAttempts, 1)
	})

	t.Run("load new attempt inserted in current postgres transaction", func(t *testing.T) {
		etx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, txStore, 3, 9, time.Now(), fromAddress)
		newAttempt := cltest.NewDynamicFeeEthTxAttempt(t, etx.ID)
		var dbAttempt txmgr.DbEthTxAttempt
		dbAttempt.FromTxAttempt(&newAttempt)

		func() {
			tx, err := db.BeginTx(ctx, nil)
			require.NoError(t, err)

			const insertEthTxAttemptSQL = `INSERT INTO evm.tx_attempts (eth_tx_id, gas_price, signed_raw_tx, hash, broadcast_before_block_num, state, created_at, chain_specific_gas_limit, tx_type, gas_tip_cap, gas_fee_cap, is_purge_attempt) VALUES (
					:eth_tx_id, :gas_price, :signed_raw_tx, :hash, :broadcast_before_block_num, :state, NOW(), :chain_specific_gas_limit, :tx_type, :gas_tip_cap, :gas_fee_cap, :is_purge_attempt
					) RETURNING *`
			query, args, err := sqlutil.DataSource(db).BindNamed(insertEthTxAttemptSQL, dbAttempt)
			require.NoError(t, err)
			_, err = tx.ExecContext(ctx, query, args...)
			require.NoError(t, err)

			etx.TxAttempts = []txmgr.TxAttempt{}
			err = txStore.LoadTxesAttempts(ctx, []*txmgr.Tx{&etx})
			require.NoError(t, err)
			assert.Len(t, etx.TxAttempts, 2)

			err = tx.Commit()
			require.NoError(t, err)
		}()

		// also check after postgres transaction is committed
		etx.TxAttempts = []txmgr.TxAttempt{}
		err := txStore.LoadTxesAttempts(ctx, []*txmgr.Tx{&etx})
		require.NoError(t, err)
		assert.Len(t, etx.TxAttempts, 2)
	})
}

func TestORM_SaveReplacementInProgressAttempt(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("replace eth tx attempt", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, 123, fromAddress)
		oldAttempt := etx.TxAttempts[0]

		newAttempt := cltest.NewDynamicFeeEthTxAttempt(t, etx.ID)
		err := txStore.SaveReplacementInProgressAttempt(tests.Context(t), oldAttempt, &newAttempt)
		require.NoError(t, err)

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		assert.Len(t, etx.TxAttempts, 1)
		require.Equal(t, etx.TxAttempts[0].Hash, newAttempt.Hash)
	})
}

func TestORM_FindNextUnstartedTransactionFromAddress(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ethClient := clienttest.NewClientWithDefaultChainID(t)

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("cannot find unstarted tx", func(t *testing.T) {
		mustInsertInProgressEthTxWithAttempt(t, txStore, 13, fromAddress)

		resultEtx, err := txStore.FindNextUnstartedTransactionFromAddress(tests.Context(t), fromAddress, ethClient.ConfiguredChainID())
		assert.ErrorIs(t, err, sql.ErrNoRows)
		assert.Nil(t, resultEtx)
	})

	t.Run("finds unstarted tx", func(t *testing.T) {
		mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, testutils.FixtureChainID)
		resultEtx, err := txStore.FindNextUnstartedTransactionFromAddress(tests.Context(t), fromAddress, ethClient.ConfiguredChainID())
		require.NoError(t, err)
		assert.NotNil(t, resultEtx)
	})
}

func TestORM_UpdateTxFatalErrorAndDeleteAttempts(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("update successful", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, 13, fromAddress)
		etxPretendError := null.StringFrom("no more toilet paper")
		etx.Error = etxPretendError

		err := txStore.UpdateTxFatalErrorAndDeleteAttempts(tests.Context(t), &etx)
		require.NoError(t, err)
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		assert.Empty(t, etx.TxAttempts)
		assert.Equal(t, txmgrcommon.TxFatalError, etx.State)
	})
}

func TestORM_UpdateTxAttemptInProgressToBroadcast(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("update successful", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, 13, fromAddress)
		attempt := etx.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptInProgress, attempt.State)

		time1 := time.Now()
		i := int16(0)
		etx.BroadcastAt = &time1
		etx.InitialBroadcastAt = &time1
		err := txStore.UpdateTxAttemptInProgressToBroadcast(tests.Context(t), &etx, attempt, txmgrtypes.TxAttemptBroadcast)
		require.NoError(t, err)
		// Increment sequence
		i++

		attemptResult, err := txStore.FindTxAttempt(ctx, attempt.Hash)
		require.NoError(t, err)
		require.Equal(t, attempt.Hash, attemptResult.Hash)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attemptResult.State)
		assert.Equal(t, int16(1), i)
	})
}

func TestORM_UpdateTxUnstartedToInProgress(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	nonce := types.Nonce(123)

	t.Run("update successful", func(t *testing.T) {
		etx := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, testutils.FixtureChainID)
		etx.Sequence = &nonce
		attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)

		err := txStore.UpdateTxUnstartedToInProgress(tests.Context(t), &etx, &attempt)
		require.NoError(t, err)

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxInProgress, etx.State)
		assert.Len(t, etx.TxAttempts, 1)
	})

	t.Run("update fails because tx is removed", func(t *testing.T) {
		etx := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, testutils.FixtureChainID)
		etx.Sequence = &nonce

		attempt := cltest.NewLegacyEthTxAttempt(t, etx.ID)

		_, err := db.ExecContext(ctx, "DELETE FROM evm.txes WHERE id = $1", etx.ID)
		require.NoError(t, err)

		err = txStore.UpdateTxUnstartedToInProgress(tests.Context(t), &etx, &attempt)
		require.ErrorContains(t, err, "tx removed")
	})

	db = testutils.NewSqlxDB(t)
	txStore = cltest.NewTestTxStore(t, db)
	ethKeyStore = cltest.NewKeyStore(t, db).Eth()
	_, fromAddress = cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("update replaces abandoned tx with same hash", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, nonce, fromAddress)
		require.Len(t, etx.TxAttempts, 1)

		zero := commonconfig.MustNewDuration(time.Duration(0))
		ccfg := configtest.NewChainScopedConfig(t, func(c *toml.EVMConfig) {
			c.Transactions.ReaperInterval = zero
			c.Transactions.ReaperThreshold = zero
			c.Transactions.ResendAfterThreshold = zero
		})
		evmTxmCfg := txmgr.NewEvmTxmConfig(ccfg.EVM())
		ec := clienttest.NewClientWithDefaultChainID(t)
		txMgr := txmgr.NewEvmTxm(ec.ConfiguredChainID(), evmTxmCfg, ccfg.EVM().Transactions(), nil, logger.Test(t), nil, nil,
			nil, txStore, nil, nil, nil, nil, nil, nil)
		err := txMgr.XXXTestAbandon(fromAddress) // mark transaction as abandoned
		require.NoError(t, err)

		etx2 := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, testutils.FixtureChainID)
		etx2.Sequence = &nonce
		attempt2 := cltest.NewLegacyEthTxAttempt(t, etx2.ID)
		attempt2.Hash = etx.TxAttempts[0].Hash

		// Even though this will initially fail due to idx_eth_tx_attempts_hash constraint, because the conflicting tx has been abandoned
		// it should succeed after removing the abandoned attempt and retrying the insert
		err = txStore.UpdateTxUnstartedToInProgress(tests.Context(t), &etx2, &attempt2)
		require.NoError(t, err)
	})

	_, fromAddress = cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	// Same flow as previous test, but without calling txMgr.Abandon()
	t.Run("duplicate tx hash disallowed in tx_eth_attempts", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, nonce, fromAddress)
		require.Len(t, etx.TxAttempts, 1)

		etx.State = txmgrcommon.TxUnstarted

		// Should fail due to idx_eth_tx_attempt_hash constraint
		err := txStore.UpdateTxUnstartedToInProgress(tests.Context(t), &etx, &etx.TxAttempts[0])
		assert.ErrorContains(t, err, "idx_eth_tx_attempts_hash")
		txStore = cltest.NewTestTxStore(t, db) // current txStore is poisened now, next test will need fresh one
	})
}

func TestORM_GetTxInProgress(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("gets 0 in progress eth transaction", func(t *testing.T) {
		etxResult, err := txStore.GetTxInProgress(tests.Context(t), fromAddress)
		require.NoError(t, err)
		require.Nil(t, etxResult)
	})

	t.Run("get 1 in progress eth transaction", func(t *testing.T) {
		etx := mustInsertInProgressEthTxWithAttempt(t, txStore, 123, fromAddress)

		etxResult, err := txStore.GetTxInProgress(tests.Context(t), fromAddress)
		require.NoError(t, err)
		assert.Equal(t, etxResult.ID, etx.ID)
	})
}

func TestORM_GetAbandonedTransactionsByBatch(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	_, enabled := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	enabledAddrs := []common.Address{enabled}

	t.Run("get 0 abandoned transactions", func(t *testing.T) {
		txes, err := txStore.GetAbandonedTransactionsByBatch(tests.Context(t), ethClient.ConfiguredChainID(), enabledAddrs, 0, 10)
		require.NoError(t, err)
		require.Empty(t, txes)
	})

	t.Run("do not return enabled addresses", func(t *testing.T) {
		_ = mustInsertInProgressEthTxWithAttempt(t, txStore, 123, enabled)
		_ = mustCreateUnstartedGeneratedTx(t, txStore, enabled, ethClient.ConfiguredChainID())
		txes, err := txStore.GetAbandonedTransactionsByBatch(tests.Context(t), ethClient.ConfiguredChainID(), enabledAddrs, 0, 10)
		require.NoError(t, err)
		require.Empty(t, txes)
	})

	t.Run("get in progress, unstarted, and unconfirmed eth transactions", func(t *testing.T) {
		inProgressTx := mustInsertInProgressEthTxWithAttempt(t, txStore, 123, fromAddress)
		unstartedTx := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, ethClient.ConfiguredChainID())

		txes, err := txStore.GetAbandonedTransactionsByBatch(tests.Context(t), ethClient.ConfiguredChainID(), enabledAddrs, 0, 10)
		require.NoError(t, err)
		require.Len(t, txes, 2)

		for _, tx := range txes {
			require.True(t, tx.ID == inProgressTx.ID || tx.ID == unstartedTx.ID)
		}
	})

	t.Run("get batches of transactions", func(t *testing.T) {
		var batchSize uint = 10
		numTxes := 55
		for i := 0; i < numTxes; i++ {
			_ = mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, ethClient.ConfiguredChainID())
		}

		allTxes := make([]*txmgr.Tx, 0)
		err := sqlutil.Batch(func(offset, limit uint) (count uint, err error) {
			batchTxes, err := txStore.GetAbandonedTransactionsByBatch(tests.Context(t), ethClient.ConfiguredChainID(), enabledAddrs, offset, limit)
			require.NoError(t, err)
			allTxes = append(allTxes, batchTxes...)
			return uint(len(batchTxes)), nil
		}, batchSize)
		require.NoError(t, err)
		require.Len(t, allTxes, numTxes+2)
	})
}

func TestORM_GetTxByID(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("no transaction", func(t *testing.T) {
		tx, err := txStore.GetTxByID(tests.Context(t), int64(0))
		require.NoError(t, err)
		require.Nil(t, tx)
	})

	t.Run("get transaction by ID", func(t *testing.T) {
		insertedTx := mustInsertInProgressEthTxWithAttempt(t, txStore, 123, fromAddress)
		tx, err := txStore.GetTxByID(tests.Context(t), insertedTx.ID)
		require.NoError(t, err)
		require.NotNil(t, tx)
	})
}

func TestORM_GetFatalTransactions(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("gets 0 fatal eth transactions", func(t *testing.T) {
		txes, err := txStore.GetFatalTransactions(tests.Context(t))
		require.NoError(t, err)
		require.Empty(t, txes)
	})

	t.Run("get fatal transactions", func(t *testing.T) {
		fatalTx := mustInsertFatalErrorEthTx(t, txStore, fromAddress)
		txes, err := txStore.GetFatalTransactions(tests.Context(t))
		require.NoError(t, err)
		require.Equal(t, txes[0].ID, fatalTx.ID)
	})
}

func TestORM_HasInProgressTransaction(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("no in progress eth transaction", func(t *testing.T) {
		exists, err := txStore.HasInProgressTransaction(tests.Context(t), fromAddress, ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("has in progress eth transaction", func(t *testing.T) {
		mustInsertInProgressEthTxWithAttempt(t, txStore, 123, fromAddress)

		exists, err := txStore.HasInProgressTransaction(tests.Context(t), fromAddress, ethClient.ConfiguredChainID())
		require.NoError(t, err)
		require.True(t, exists)
	})
}

func TestORM_CountUnconfirmedTransactions(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, otherAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, otherAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress)

	count, err := txStore.CountUnconfirmedTransactions(tests.Context(t), fromAddress, testutils.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, int(count), 3)
}

func TestORM_CountTransactionsByState(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress1 := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, fromAddress2 := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, fromAddress3 := cltest.MustInsertRandomKey(t, ethKeyStore)

	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress1)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress2)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress3)

	count, err := txStore.CountTransactionsByState(tests.Context(t), txmgrcommon.TxUnconfirmed, testutils.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, int(count), 3)
}

func TestORM_CountUnstartedTransactions(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, otherAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, testutils.FixtureChainID)
	mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, testutils.FixtureChainID)
	mustCreateUnstartedGeneratedTx(t, txStore, otherAddress, testutils.FixtureChainID)
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress)

	count, err := txStore.CountUnstartedTransactions(tests.Context(t), fromAddress, testutils.FixtureChainID)
	require.NoError(t, err)
	assert.Equal(t, int(count), 2)
}

func TestORM_CheckTxQueueCapacity(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, otherAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	toAddress := testutils.NewAddress()
	encodedPayload := []byte{1, 2, 3}
	feeLimit := uint64(1000000000)
	value := big.Int(assets.NewEthValue(142))
	var maxUnconfirmedTransactions uint64 = 2

	t.Run("with no eth_txes returns nil", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(tests.Context(t), fromAddress, maxUnconfirmedTransactions, testutils.FixtureChainID)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		mustCreateUnstartedTx(t, txStore, otherAddress, toAddress, encodedPayload, feeLimit, value, testutils.FixtureChainID)
	}

	t.Run("with eth_txes from another address returns nil", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(tests.Context(t), fromAddress, maxUnconfirmedTransactions, testutils.FixtureChainID)
		require.NoError(t, err)
	})

	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		mustInsertFatalErrorEthTx(t, txStore, otherAddress)
	}

	t.Run("ignores fatally_errored transactions", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(tests.Context(t), fromAddress, maxUnconfirmedTransactions, testutils.FixtureChainID)
		require.NoError(t, err)
	})

	var n int64
	mustInsertInProgressEthTxWithAttempt(t, txStore, types.Nonce(n), fromAddress)
	n++
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, n, fromAddress)
	n++

	t.Run("unconfirmed and in_progress transactions do not count", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(tests.Context(t), fromAddress, 1, testutils.FixtureChainID)
		require.NoError(t, err)
	})

	// deliberately one extra to exceed limit
	for i := 0; i <= int(maxUnconfirmedTransactions); i++ {
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, n, 42, fromAddress)
		n++
	}

	t.Run("with many confirmed eth_txes from the same address returns nil", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(tests.Context(t), fromAddress, maxUnconfirmedTransactions, testutils.FixtureChainID)
		require.NoError(t, err)
	})

	for i := 0; i < int(maxUnconfirmedTransactions)-1; i++ {
		mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, feeLimit, value, testutils.FixtureChainID)
	}

	t.Run("with fewer unstarted eth_txes than limit returns nil", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(tests.Context(t), fromAddress, maxUnconfirmedTransactions, testutils.FixtureChainID)
		require.NoError(t, err)
	})

	mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, feeLimit, value, testutils.FixtureChainID)

	t.Run("with equal or more unstarted eth_txes than limit returns error", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(tests.Context(t), fromAddress, maxUnconfirmedTransactions, testutils.FixtureChainID)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("cannot create transaction; too many unstarted transactions in the queue (2/%d). WARNING: Hitting EVM.Transactions.MaxQueued", maxUnconfirmedTransactions))

		mustCreateUnstartedTx(t, txStore, fromAddress, toAddress, encodedPayload, feeLimit, value, testutils.FixtureChainID)
		err = txStore.CheckTxQueueCapacity(tests.Context(t), fromAddress, maxUnconfirmedTransactions, testutils.FixtureChainID)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("cannot create transaction; too many unstarted transactions in the queue (3/%d). WARNING: Hitting EVM.Transactions.MaxQueued", maxUnconfirmedTransactions))
	})

	t.Run("with different chain ID ignores txes", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(tests.Context(t), fromAddress, maxUnconfirmedTransactions, big.NewInt(42))
		require.NoError(t, err)
	})

	t.Run("disables check with 0 limit", func(t *testing.T) {
		err := txStore.CheckTxQueueCapacity(tests.Context(t), fromAddress, 0, testutils.FixtureChainID)
		require.NoError(t, err)
	})
}

func TestORM_CreateTransaction(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := newTxStore(t, db)
	kst := cltest.NewKeyStore(t, db)

	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())
	toAddress := testutils.NewAddress()
	gasLimit := uint64(1000)
	payload := []byte{1, 2, 3}

	ethClient := clienttest.NewClientWithDefaultChainID(t)

	t.Run("with queue under capacity inserts eth_tx", func(t *testing.T) {
		subject := uuid.New()
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		etx, err := txStore.CreateTransaction(tests.Context(t), txmgr.TxRequest{
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
		tx1, err := txStore.CreateTransaction(tests.Context(t), txRequest, ethClient.ConfiguredChainID())
		assert.NoError(t, err)

		tx2, err := txStore.CreateTransaction(tests.Context(t), txRequest, ethClient.ConfiguredChainID())
		assert.NoError(t, err)

		assert.Equal(t, tx1.GetID(), tx2.GetID())
	})

	t.Run("sets signal callback flag", func(t *testing.T) {
		subject := uuid.New()
		strategy := newMockTxStrategy(t)
		strategy.On("Subject").Return(uuid.NullUUID{UUID: subject, Valid: true})
		etx, err := txStore.CreateTransaction(tests.Context(t), txmgr.TxRequest{
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
		assert.True(t, etx.SignalCallback)

		cltest.AssertCount(t, db, "evm.txes", 3)

		var dbEthTx txmgr.DbEthTx
		require.NoError(t, db.Get(&dbEthTx, `SELECT * FROM evm.txes ORDER BY id DESC LIMIT 1`))

		assert.Equal(t, fromAddress, dbEthTx.FromAddress)
		assert.True(t, dbEthTx.SignalCallback)
	})
}

func TestORM_PruneUnstartedTxQueue(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := txmgr.NewTxStore(db, logger.Test(t))
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	t.Run("does not prune if queue has not exceeded capacity-1", func(t *testing.T) {
		subject1 := uuid.New()
		strategy1 := txmgrcommon.NewDropOldestStrategy(subject1, uint32(5))
		for i := 0; i < 5; i++ {
			mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, testutils.FixtureChainID, txRequestWithStrategy(strategy1))
		}
		AssertCountPerSubject(t, txStore, int64(4), subject1)
	})

	t.Run("prunes if queue has exceeded capacity-1", func(t *testing.T) {
		subject2 := uuid.New()
		strategy2 := txmgrcommon.NewDropOldestStrategy(subject2, uint32(3))
		for i := 0; i < 5; i++ {
			mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, testutils.FixtureChainID, txRequestWithStrategy(strategy2))
		}
		AssertCountPerSubject(t, txStore, int64(2), subject2)
	})
}

func TestORM_FindTxesWithAttemptsAndReceiptsByIdsAndState(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ctx := tests.Context(t)

	_, from := cltest.MustInsertRandomKey(t, ethKeyStore)

	tx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 0, 1, from)
	r := newEthReceipt(4, utils.NewHash(), tx.TxAttempts[0].Hash, 0x1)
	_, err := txStore.InsertReceipt(ctx, &r.Receipt)
	require.NoError(t, err)

	txes, err := txStore.FindTxesWithAttemptsAndReceiptsByIdsAndState(ctx, []int64{tx.ID}, []txmgrtypes.TxState{txmgrcommon.TxConfirmed}, testutils.FixtureChainID)
	require.NoError(t, err)
	require.Len(t, txes, 1)
	require.Len(t, txes[0].TxAttempts, 1)
	require.Len(t, txes[0].TxAttempts[0].Receipts, 1)
}

func AssertCountPerSubject(t *testing.T, txStore txmgr.TestEvmTxStore, expected int64, subject uuid.UUID) {
	t.Helper()
	count, err := txStore.CountTxesByStateAndSubject(tests.Context(t), "unstarted", subject)
	require.NoError(t, err)
	require.Equal(t, int(expected), count)
}

func TestORM_FindAttemptsRequiringReceiptFetch(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	blockNum := int64(100)

	t.Run("finds confirmed transaction requiring receipt fetch", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		kst := cltest.NewKeyStore(t, db)
		_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())
		// Transactions whose attempts should not be picked up for receipt fetch
		mustInsertFatalErrorEthTx(t, txStore, fromAddress)
		mustInsertUnstartedTx(t, txStore, fromAddress)
		mustInsertInProgressEthTxWithAttempt(t, txStore, 4, fromAddress)
		mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 3, fromAddress, txmgrtypes.TxAttemptBroadcast)
		mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 2, blockNum)
		// Terminally stuck transaction with receipt should NOT be picked up for receipt fetch
		stuckTx := mustInsertTerminallyStuckTxWithAttempt(t, txStore, fromAddress, 1, blockNum)
		mustInsertEthReceipt(t, txStore, blockNum, utils.NewHash(), stuckTx.TxAttempts[0].Hash)
		// Fatal transactions with nil nonce and stored attempts should NOT be picked up for receipt fetch
		fatalTxWithAttempt := mustInsertFatalErrorEthTx(t, txStore, fromAddress)
		attempt := newBroadcastLegacyEthTxAttempt(t, fatalTxWithAttempt.ID)
		err := txStore.InsertTxAttempt(ctx, &attempt)
		require.NoError(t, err)

		// Confirmed transaction without receipt should be picked up for receipt fetch
		confirmedTx := mustInsertConfirmedEthTx(t, txStore, 0, fromAddress)
		attempt = newBroadcastLegacyEthTxAttempt(t, confirmedTx.ID)
		err = txStore.InsertTxAttempt(ctx, &attempt)
		require.NoError(t, err)

		attempts, err := txStore.FindAttemptsRequiringReceiptFetch(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		require.Len(t, attempts, 1)
		require.Equal(t, attempt.Hash.String(), attempts[0].Hash.String())
	})

	t.Run("finds terminally stuck transaction requiring receipt fetch", func(t *testing.T) {
		db := testutils.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		kst := cltest.NewKeyStore(t, db)
		_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())
		// Transactions whose attempts should not be picked up for receipt fetch
		mustInsertFatalErrorEthTx(t, txStore, fromAddress)
		mustInsertUnstartedTx(t, txStore, fromAddress)
		mustInsertInProgressEthTxWithAttempt(t, txStore, 4, fromAddress)
		mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 3, fromAddress, txmgrtypes.TxAttemptBroadcast)
		mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 2, blockNum)
		// Terminally stuck transaction with receipt should NOT be picked up for receipt fetch
		stuckTxWithReceipt := mustInsertTerminallyStuckTxWithAttempt(t, txStore, fromAddress, 1, blockNum)
		mustInsertEthReceipt(t, txStore, blockNum, utils.NewHash(), stuckTxWithReceipt.TxAttempts[0].Hash)
		// Fatal transactions with nil nonce and stored attempts should NOT be picked up for receipt fetch
		fatalTxWithAttempt := mustInsertFatalErrorEthTx(t, txStore, fromAddress)
		attempt := newBroadcastLegacyEthTxAttempt(t, fatalTxWithAttempt.ID)
		err := txStore.InsertTxAttempt(ctx, &attempt)
		require.NoError(t, err)

		// Terminally stuck transaction without receipt should be picked up for receipt fetch
		stuckTxWoutReceipt := mustInsertTerminallyStuckTxWithAttempt(t, txStore, fromAddress, 0, blockNum)

		attempts, err := txStore.FindAttemptsRequiringReceiptFetch(ctx, testutils.FixtureChainID)
		require.NoError(t, err)
		require.Len(t, attempts, 1)
		require.Equal(t, stuckTxWoutReceipt.TxAttempts[0].Hash.String(), attempts[0].Hash.String())
	})
}

func TestORM_UpdateTxStatesToFinalizedUsingTxHashes(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db)
	broadcast := time.Now()
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	t.Run("successfully finalizes a confirmed transaction", func(t *testing.T) {
		nonce := types.Nonce(0)
		tx := &txmgr.Tx{
			Sequence:           &nonce,
			FromAddress:        fromAddress,
			EncodedPayload:     []byte{1, 2, 3},
			State:              txmgrcommon.TxConfirmed,
			BroadcastAt:        &broadcast,
			InitialBroadcastAt: &broadcast,
		}
		err := txStore.InsertTx(ctx, tx)
		require.NoError(t, err)
		attempt := newBroadcastLegacyEthTxAttempt(t, tx.ID)
		err = txStore.InsertTxAttempt(ctx, &attempt)
		require.NoError(t, err)
		mustInsertEthReceipt(t, txStore, 100, testutils.NewHash(), attempt.Hash)
		err = txStore.UpdateTxStatesToFinalizedUsingTxHashes(ctx, []common.Hash{attempt.Hash}, testutils.FixtureChainID)
		require.NoError(t, err)
		etx, err := txStore.FindTxWithAttempts(ctx, tx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxFinalized, etx.State)
	})
}

func TestORM_FindReorgOrIncludedTxs(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db)
	blockNum := int64(100)
	t.Run("finds re-org'd transactions using the mined tx count", func(t *testing.T) {
		_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())
		_, otherAddress := cltest.MustInsertRandomKey(t, kst.Eth())
		// Unstarted can't be re-org'd
		mustInsertUnstartedTx(t, txStore, fromAddress)
		// In-Progress can't be re-org'd
		mustInsertInProgressEthTxWithAttempt(t, txStore, 4, fromAddress)
		// Unconfirmed can't be re-org'd
		mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 3, fromAddress, txmgrtypes.TxAttemptBroadcast)
		// Confirmed and nonce greater than mined tx count so has been re-org'd
		mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 2, blockNum)
		// Fatal error and nonce equal to mined tx count so has been re-org'd
		mustInsertTerminallyStuckTxWithAttempt(t, txStore, fromAddress, 1, blockNum)
		// Nonce lower than mined tx count so has not been re-org
		mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 0, blockNum)

		// Tx for another from address should not be returned
		mustInsertConfirmedEthTxWithReceipt(t, txStore, otherAddress, 1, blockNum)
		mustInsertConfirmedEthTxWithReceipt(t, txStore, otherAddress, 0, blockNum)

		reorgTxs, includedTxs, err := txStore.FindReorgOrIncludedTxs(ctx, fromAddress, types.Nonce(1), testutils.FixtureChainID)
		require.NoError(t, err)
		require.Len(t, reorgTxs, 2)
		require.Empty(t, includedTxs)
	})

	t.Run("finds transactions included on-chain using the mined tx count", func(t *testing.T) {
		_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())
		_, otherAddress := cltest.MustInsertRandomKey(t, kst.Eth())
		// Unstarted can't be included
		mustInsertUnstartedTx(t, txStore, fromAddress)
		// In-Progress can't be included
		mustInsertInProgressEthTxWithAttempt(t, txStore, 5, fromAddress)
		// Unconfirmed with higher nonce can't be included
		mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 4, fromAddress, txmgrtypes.TxAttemptBroadcast)
		// Unconfirmed with nonce less than mined tx count is newly included
		mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 3, fromAddress, txmgrtypes.TxAttemptBroadcast)
		// Unconfirmed with purge attempt with nonce less than mined tx cound is newly included
		mustInsertUnconfirmedEthTxWithBroadcastPurgeAttempt(t, txStore, 2, fromAddress)
		// Fatal error so already included
		mustInsertTerminallyStuckTxWithAttempt(t, txStore, fromAddress, 1, blockNum)
		// Confirmed so already included
		mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 0, blockNum)

		// Tx for another from address should not be returned
		mustInsertConfirmedEthTxWithReceipt(t, txStore, otherAddress, 1, blockNum)
		mustInsertConfirmedEthTxWithReceipt(t, txStore, otherAddress, 0, blockNum)

		reorgTxs, includedTxs, err := txStore.FindReorgOrIncludedTxs(ctx, fromAddress, types.Nonce(4), testutils.FixtureChainID)
		require.NoError(t, err)
		require.Len(t, includedTxs, 2)
		require.Empty(t, reorgTxs)
	})
}

func TestORM_UpdateTxFatalError(t *testing.T) {
	t.Parallel()

	ctx := tests.Context(t)
	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())
	t.Run("successfully marks transaction as fatal with error message", func(t *testing.T) {
		// Unconfirmed with purge attempt with nonce less than mined tx cound is newly included
		tx1 := mustInsertUnconfirmedEthTxWithBroadcastPurgeAttempt(t, txStore, 0, fromAddress)
		tx2 := mustInsertUnconfirmedEthTxWithBroadcastPurgeAttempt(t, txStore, 1, fromAddress)

		err := txStore.UpdateTxFatalError(ctx, []int64{tx1.ID, tx2.ID}, client.TerminallyStuckMsg)
		require.NoError(t, err)

		tx1, err = txStore.FindTxWithAttempts(ctx, tx1.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxFatalError, tx1.State)
		require.Equal(t, client.TerminallyStuckMsg, tx1.Error.String)
		tx2, err = txStore.FindTxWithAttempts(ctx, tx2.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxFatalError, tx2.State)
		require.Equal(t, client.TerminallyStuckMsg, tx2.Error.String)
	})
}

func TestORM_FindTxesByIDs(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ctx := tests.Context(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	// tx state should be confirmed missing receipt
	// attempt should be before latestFinalizedBlockNum
	t.Run("successfully finds transactions with IDs", func(t *testing.T) {
		etx1 := mustInsertInProgressEthTxWithAttempt(t, txStore, 3, fromAddress)
		etx2 := mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 2, fromAddress, txmgrtypes.TxAttemptBroadcast)
		etx3 := mustInsertTerminallyStuckTxWithAttempt(t, txStore, fromAddress, 1, 100)
		etx4 := mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 0, 100)

		etxIDs := []int64{etx1.ID, etx2.ID, etx3.ID, etx4.ID}
		oldTxs, err := txStore.FindTxesByIDs(ctx, etxIDs, testutils.FixtureChainID)
		require.NoError(t, err)
		require.Len(t, oldTxs, 4)
	})
}

func TestORM_DeleteReceiptsByTxHash(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ctx := tests.Context(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	etx1 := mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 0, 100)
	etx2 := mustInsertConfirmedEthTxWithReceipt(t, txStore, fromAddress, 2, 100)

	// Delete one transaction's receipt
	err := txStore.DeleteReceiptByTxHash(ctx, etx1.TxAttempts[0].Hash)
	require.NoError(t, err)

	// receipt has been deleted
	etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
	require.NoError(t, err)
	require.Empty(t, etx1.TxAttempts[0].Receipts)

	// receipt still exists for other tx
	etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
	require.NoError(t, err)
	require.Len(t, etx2.TxAttempts[0].Receipts, 1)
}

func TestORM_Abandon(t *testing.T) {
	t.Parallel()

	db := testutils.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ctx := tests.Context(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	etx1 := mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, testutils.FixtureChainID)
	etx2 := mustInsertInProgressEthTxWithAttempt(t, txStore, 1, fromAddress)
	etx3 := mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, 0, fromAddress, txmgrtypes.TxAttemptBroadcast)

	err := txStore.Abandon(ctx, testutils.FixtureChainID, fromAddress)
	require.NoError(t, err)

	// transactions marked as fatal error with abandon reason, nil sequence, and no attempts
	etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
	require.NoError(t, err)
	require.Equal(t, txmgrcommon.TxFatalError, etx1.State)
	require.Nil(t, etx1.Sequence)
	require.Empty(t, etx1.TxAttempts)

	etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
	require.NoError(t, err)
	require.Equal(t, txmgrcommon.TxFatalError, etx2.State)
	require.Nil(t, etx2.Sequence)
	require.Empty(t, etx2.TxAttempts)

	etx3, err = txStore.FindTxWithAttempts(ctx, etx3.ID)
	require.NoError(t, err)
	require.Equal(t, txmgrcommon.TxFatalError, etx3.State)
	require.Nil(t, etx3.Sequence)
	require.Empty(t, etx3.TxAttempts)
}

func mustInsertTerminallyStuckTxWithAttempt(t testing.TB, txStore txmgr.TestEvmTxStore, fromAddress common.Address, nonceInt int64, broadcastBeforeBlockNum int64) txmgr.Tx {
	ctx := tests.Context(t)
	broadcast := time.Now()
	nonce := types.Nonce(nonceInt)
	tx := txmgr.Tx{
		Sequence:           &nonce,
		FromAddress:        fromAddress,
		EncodedPayload:     []byte{1, 2, 3},
		State:              txmgrcommon.TxFatalError,
		BroadcastAt:        &broadcast,
		InitialBroadcastAt: &broadcast,
		Error:              null.StringFrom(client.TerminallyStuckMsg),
	}
	require.NoError(t, txStore.InsertTx(ctx, &tx))
	attempt := cltest.NewLegacyEthTxAttempt(t, tx.ID)
	attempt.BroadcastBeforeBlockNum = &broadcastBeforeBlockNum
	attempt.State = txmgrtypes.TxAttemptBroadcast
	attempt.IsPurgeAttempt = true
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	tx, err := txStore.FindTxWithAttempts(ctx, tx.ID)
	require.NoError(t, err)
	return tx
}
