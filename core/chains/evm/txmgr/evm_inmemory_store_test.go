package txmgr_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontxmgr "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"

	evmgas "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestInMemoryStore_SaveInProgressAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("saves new in_progress attempt if attempt is new", func(t *testing.T) {
		// Insert a transaction into persistent store
		inTx := cltest.MustInsertUnconfirmedEthTx(t, persistentStore, 1, fromAddress)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		// generate new attempt
		inTxAttempt := cltest.NewLegacyEthTxAttempt(t, inTx.ID)
		require.Equal(t, int64(0), inTxAttempt.ID)

		err := inMemoryStore.SaveInProgressAttempt(ctx, &inTxAttempt)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)

		// Check that the in-memory store has the new attempt
		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.NotNil(t, actTxs)
		actTx := actTxs[0]
		require.Equal(t, len(expTx.TxAttempts), len(actTx.TxAttempts))

		assertTxEqual(t, expTx, actTx)
	})
	t.Run("updates old attempt to in_progress when insufficient_funds", func(t *testing.T) {
		// Insert a transaction into persistent store
		inTx := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, persistentStore, 23, fromAddress)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		// use old attempt
		inTxAttempt := inTx.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptInsufficientFunds, inTxAttempt.State)
		require.NotEqual(t, int64(0), inTxAttempt.ID)

		inTxAttempt.BroadcastBeforeBlockNum = nil
		inTxAttempt.State = txmgrtypes.TxAttemptInProgress
		err := inMemoryStore.SaveInProgressAttempt(ctx, &inTxAttempt)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)

		// Check that the in-memory store has the new attempt
		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.NotNil(t, actTxs)
		actTx := actTxs[0]
		require.Equal(t, len(expTx.TxAttempts), len(actTx.TxAttempts))

		assertTxEqual(t, expTx, actTx)
	})
	t.Run("handles errors the same way as the persistent store", func(t *testing.T) {
		// Insert a transaction into persistent store
		inTx := cltest.MustInsertUnconfirmedEthTx(t, persistentStore, 55, fromAddress)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		// generate new attempt
		inTxAttempt := cltest.NewLegacyEthTxAttempt(t, inTx.ID)
		require.Equal(t, int64(0), inTxAttempt.ID)

		t.Run("wrong tx id", func(t *testing.T) {
			inTxAttempt.TxID = 999
			actErr := inMemoryStore.SaveInProgressAttempt(ctx, &inTxAttempt)
			expErr := persistentStore.SaveInProgressAttempt(ctx, &inTxAttempt)
			assert.Error(t, actErr)
			assert.Error(t, expErr)
			inTxAttempt.TxID = inTx.ID // reset
		})

		t.Run("wrong state", func(t *testing.T) {
			inTxAttempt.State = txmgrtypes.TxAttemptBroadcast
			actErr := inMemoryStore.SaveInProgressAttempt(ctx, &inTxAttempt)
			expErr := persistentStore.SaveInProgressAttempt(ctx, &inTxAttempt)
			assert.Error(t, actErr)
			assert.Error(t, expErr)
			assert.Equal(t, expErr, actErr)
			inTxAttempt.State = txmgrtypes.TxAttemptInProgress // reset
		})
	})
}

func TestInMemoryStore_UpdateTxCallbackCompleted(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("sets tx callback as completed", func(t *testing.T) {
		// Insert a transaction into persistent store
		inTx := cltest.NewEthTx(fromAddress)
		inTx.PipelineTaskRunID = uuid.NullUUID{UUID: uuid.New(), Valid: true}
		require.NoError(t, persistentStore.InsertTx(ctx, &inTx))
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		err := inMemoryStore.UpdateTxCallbackCompleted(
			testutils.Context(t),
			inTx.PipelineTaskRunID.UUID,
			chainID,
		)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)
		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]
		assertTxEqual(t, expTx, actTx)
		assert.True(t, actTx.CallbackCompleted)

		// wrong PipelineTaskRunID
		wrongPipelineTaskRunID := uuid.NullUUID{UUID: uuid.New(), Valid: true}
		actErr := inMemoryStore.UpdateTxCallbackCompleted(ctx, wrongPipelineTaskRunID.UUID, chainID)
		expErr := persistentStore.UpdateTxCallbackCompleted(ctx, wrongPipelineTaskRunID.UUID, chainID)
		assert.NoError(t, actErr)
		assert.NoError(t, expErr)
	})
}

func TestInMemoryStore_SaveInsufficientFundsAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	defaultDuration := time.Second * 5
	t.Run("updates attempt state and checks error returns", func(t *testing.T) {
		// Insert a transaction into persistent store
		inTx := mustInsertInProgressEthTxWithAttempt(t, persistentStore, 1, fromAddress)
		now := time.Now()
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		err := inMemoryStore.SaveInsufficientFundsAttempt(
			ctx,
			defaultDuration,
			&inTx.TxAttempts[0],
			now,
		)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)
		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]
		assertTxEqual(t, expTx, actTx)
		assert.Equal(t, txmgrtypes.TxAttemptInsufficientFunds, actTx.TxAttempts[0].State)

		// wrong tx id
		inTx.TxAttempts[0].TxID = 123
		actErr := inMemoryStore.SaveInsufficientFundsAttempt(ctx, defaultDuration, &inTx.TxAttempts[0], now)
		expErr := persistentStore.SaveInsufficientFundsAttempt(ctx, defaultDuration, &inTx.TxAttempts[0], now)
		assert.NoError(t, actErr)
		assert.NoError(t, expErr)
		inTx.TxAttempts[0].TxID = inTx.ID // reset

		// wrong attempt state
		inTx.TxAttempts[0].State = txmgrtypes.TxAttemptBroadcast
		actErr = inMemoryStore.SaveInsufficientFundsAttempt(ctx, defaultDuration, &inTx.TxAttempts[0], now)
		expErr = persistentStore.SaveInsufficientFundsAttempt(ctx, defaultDuration, &inTx.TxAttempts[0], now)
		assert.Error(t, actErr)
		assert.Error(t, expErr)
		inTx.TxAttempts[0].State = txmgrtypes.TxAttemptInsufficientFunds // reset
	})
}

func TestInMemoryStore_SaveSentAttempt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	defaultDuration := time.Second * 5
	t.Run("updates attempt state to broadcast and checks error returns", func(t *testing.T) {
		// Insert a transaction into persistent store
		inTx := mustInsertInProgressEthTxWithAttempt(t, persistentStore, 1, fromAddress)
		require.Nil(t, inTx.BroadcastAt)
		now := time.Now()
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		err := inMemoryStore.SaveSentAttempt(
			ctx,
			defaultDuration,
			&inTx.TxAttempts[0],
			now,
		)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)
		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]
		assertTxEqual(t, expTx, actTx)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, actTx.TxAttempts[0].State)

		// wrong tx id
		inTx.TxAttempts[0].TxID = 123
		actErr := inMemoryStore.SaveSentAttempt(ctx, defaultDuration, &inTx.TxAttempts[0], now)
		expErr := persistentStore.SaveSentAttempt(ctx, defaultDuration, &inTx.TxAttempts[0], now)
		assert.Error(t, actErr)
		assert.Error(t, expErr)
		inTx.TxAttempts[0].TxID = inTx.ID // reset

		// wrong attempt state
		inTx.TxAttempts[0].State = txmgrtypes.TxAttemptBroadcast
		actErr = inMemoryStore.SaveSentAttempt(ctx, defaultDuration, &inTx.TxAttempts[0], now)
		expErr = persistentStore.SaveSentAttempt(ctx, defaultDuration, &inTx.TxAttempts[0], now)
		assert.Error(t, actErr)
		assert.Error(t, expErr)
		inTx.TxAttempts[0].State = txmgrtypes.TxAttemptInProgress // reset
	})
}

func TestInMemoryStore_Abandon(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore[
		*big.Int,
		common.Address, common.Hash, common.Hash,
		*evmtypes.Receipt,
		evmtypes.Nonce,
		evmgas.EvmFee,
	](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("Abandon transactions successfully", func(t *testing.T) {
		nTxs := 3
		for i := 0; i < nTxs; i++ {
			inTx := cltest.NewEthTx(fromAddress)
			// insert the transaction into the persistent store
			require.NoError(t, persistentStore.InsertTx(ctx, &inTx))
			// insert the transaction into the in-memory store
			require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))
		}

		actErr := inMemoryStore.Abandon(ctx, chainID, fromAddress)
		expErr := persistentStore.Abandon(ctx, chainID, fromAddress)
		require.NoError(t, actErr)
		require.NoError(t, expErr)

		expTxs, err := persistentStore.FindTxesByFromAddressAndState(ctx, fromAddress, "fatal_error")
		require.NoError(t, err)
		require.NotNil(t, expTxs)
		require.Equal(t, nTxs, len(expTxs))

		// Check the in-memory store
		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn)
		require.NotNil(t, actTxs)
		require.Equal(t, nTxs, len(actTxs))

		for i := 0; i < nTxs; i++ {
			assertTxEqual(t, *expTxs[i], actTxs[i])
		}
	})
}

// assertTxEqual asserts that two transactions are equal
func assertTxEqual(t *testing.T, exp, act evmtxmgr.Tx) {
	assert.Equal(t, exp.ID, act.ID)
	assert.Equal(t, exp.IdempotencyKey, act.IdempotencyKey)
	assert.Equal(t, exp.Sequence, act.Sequence)
	assert.Equal(t, exp.FromAddress, act.FromAddress)
	assert.Equal(t, exp.ToAddress, act.ToAddress)
	assert.Equal(t, exp.EncodedPayload, act.EncodedPayload)
	assert.Equal(t, exp.Value, act.Value)
	assert.Equal(t, exp.FeeLimit, act.FeeLimit)
	assert.Equal(t, exp.Error, act.Error)
	assert.Equal(t, exp.BroadcastAt, act.BroadcastAt)
	assert.Equal(t, exp.InitialBroadcastAt, act.InitialBroadcastAt)
	assert.Equal(t, exp.CreatedAt, act.CreatedAt)
	assert.Equal(t, exp.State, act.State)
	assert.Equal(t, exp.Meta, act.Meta)
	assert.Equal(t, exp.Subject, act.Subject)
	assert.Equal(t, exp.ChainID, act.ChainID)
	assert.Equal(t, exp.PipelineTaskRunID, act.PipelineTaskRunID)
	assert.Equal(t, exp.MinConfirmations, act.MinConfirmations)
	assert.Equal(t, exp.TransmitChecker, act.TransmitChecker)
	assert.Equal(t, exp.SignalCallback, act.SignalCallback)
	assert.Equal(t, exp.CallbackCompleted, act.CallbackCompleted)

	require.Len(t, exp.TxAttempts, len(act.TxAttempts))
	for i := 0; i < len(exp.TxAttempts); i++ {
		assertTxAttemptEqual(t, exp.TxAttempts[i], act.TxAttempts[i])
	}
}

func assertTxAttemptEqual(t *testing.T, exp, act evmtxmgr.TxAttempt) {
	assert.Equal(t, exp.ID, act.ID)
	assert.Equal(t, exp.TxID, act.TxID)
	assert.Equal(t, exp.TxFee, act.TxFee)
	assert.Equal(t, exp.ChainSpecificFeeLimit, act.ChainSpecificFeeLimit)
	assert.Equal(t, exp.SignedRawTx, act.SignedRawTx)
	assert.Equal(t, exp.Hash, act.Hash)
	assert.Equal(t, exp.CreatedAt, act.CreatedAt)
	assert.Equal(t, exp.BroadcastBeforeBlockNum, act.BroadcastBeforeBlockNum)
	assert.Equal(t, exp.State, act.State)
	assert.Equal(t, exp.TxType, act.TxType)

	require.Equal(t, len(exp.Receipts), len(act.Receipts))
	for i := 0; i < len(exp.Receipts); i++ {
		assertChainReceiptEqual(t, exp.Receipts[i], act.Receipts[i])
	}
}

func assertChainReceiptEqual(t *testing.T, exp, act evmtxmgr.ChainReceipt) {
	assert.Equal(t, exp.GetStatus(), act.GetStatus())
	assert.Equal(t, exp.GetTxHash(), act.GetTxHash())
	assert.Equal(t, exp.GetBlockNumber(), act.GetBlockNumber())
	assert.Equal(t, exp.IsZero(), act.IsZero())
	assert.Equal(t, exp.IsUnmined(), act.IsUnmined())
	assert.Equal(t, exp.GetFeeUsed(), act.GetFeeUsed())
	assert.Equal(t, exp.GetTransactionIndex(), act.GetTransactionIndex())
	assert.Equal(t, exp.GetBlockHash(), act.GetBlockHash())
}
