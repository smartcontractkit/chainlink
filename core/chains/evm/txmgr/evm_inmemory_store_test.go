package txmgr_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontxmgr "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"

	evmgas "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
)

func TestInMemoryStore_UpdateTxAttemptInProgressToBroadcast(t *testing.T) {
	t.Parallel()

	t.Run("successfully updates transactions from in progress to broadcast", func(t *testing.T) {
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

		// Insert a transaction into persistent store
		inTx := mustInsertInProgressEthTxWithAttempt(t, persistentStore, 13, fromAddress)
		inTxAttempt := inTx.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptInProgress, inTxAttempt.State)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		time1 := time.Now()
		inTx.BroadcastAt = &time1
		inTx.InitialBroadcastAt = &time1

		// Update the transaction attempt
		require.NoError(t, inMemoryStore.UpdateTxAttemptInProgressToBroadcast(
			ctx,
			&inTx,
			inTxAttempt,
			txmgrtypes.TxAttemptBroadcast,
		))

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)
		assert.Equal(t, commontxmgr.TxUnconfirmed, expTx.State)
		assert.Equal(t, 1, len(expTx.TxAttempts))
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, expTx.TxAttempts[0].State)

		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]
		assert.Equal(t, commontxmgr.TxUnconfirmed, actTx.State)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, actTx.TxAttempts[0].State)

		// verify that both the in-memory and persistent store have the same transaction state
		assertTxEqual(t, expTx, actTx)
	})

	t.Run("verify in-memory error handling has parity with persistent store", func(t *testing.T) {
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

		// Insert a transaction into persistent store
		inTx := mustInsertInProgressEthTxWithAttempt(t, persistentStore, 13, fromAddress)
		inTxAttempt := inTx.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptInProgress, inTxAttempt.State)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		time1 := time.Now()
		inTx.BroadcastAt = &time1
		inTx.InitialBroadcastAt = &time1

		t.Run("nil broadcast at", func(t *testing.T) {
			inTx.BroadcastAt = nil
			expErr := persistentStore.UpdateTxAttemptInProgressToBroadcast(ctx, &inTx, inTxAttempt, txmgrtypes.TxAttemptBroadcast)
			actErr := inMemoryStore.UpdateTxAttemptInProgressToBroadcast(ctx, &inTx, inTxAttempt, txmgrtypes.TxAttemptBroadcast)
			assert.Equal(t, expErr, actErr)
			inTx.BroadcastAt = &time1 // reset
		})

		t.Run("nil initial broadcast at", func(t *testing.T) {
			inTx.InitialBroadcastAt = nil
			expErr := persistentStore.UpdateTxAttemptInProgressToBroadcast(ctx, &inTx, inTxAttempt, txmgrtypes.TxAttemptBroadcast)
			actErr := inMemoryStore.UpdateTxAttemptInProgressToBroadcast(ctx, &inTx, inTxAttempt, txmgrtypes.TxAttemptBroadcast)
			assert.Equal(t, expErr, actErr)
			inTx.InitialBroadcastAt = &time1 // reset
		})

		t.Run("transaction not in progress", func(t *testing.T) {
			inTx.State = commontxmgr.TxConfirmed
			expErr := persistentStore.UpdateTxAttemptInProgressToBroadcast(ctx, &inTx, inTxAttempt, txmgrtypes.TxAttemptBroadcast)
			actErr := inMemoryStore.UpdateTxAttemptInProgressToBroadcast(ctx, &inTx, inTxAttempt, txmgrtypes.TxAttemptBroadcast)
			assert.ErrorContains(t, expErr, "can only transition to unconfirmed from in_progress")
			assert.ErrorContains(t, actErr, "can only transition to unconfirmed from in_progress")
			inTx.State = commontxmgr.TxInProgress // reset
		})

		t.Run("transaction attempt not in progress", func(t *testing.T) {
			inTxAttempt.State = txmgrtypes.TxAttemptBroadcast
			expErr := persistentStore.UpdateTxAttemptInProgressToBroadcast(ctx, &inTx, inTxAttempt, txmgrtypes.TxAttemptBroadcast)
			actErr := inMemoryStore.UpdateTxAttemptInProgressToBroadcast(ctx, &inTx, inTxAttempt, txmgrtypes.TxAttemptBroadcast)
			assert.Equal(t, expErr, actErr)
			inTxAttempt.State = txmgrtypes.TxAttemptInProgress // reset
		})

		t.Run("new attempt state not broadcast", func(t *testing.T) {
			expErr := persistentStore.UpdateTxAttemptInProgressToBroadcast(ctx, &inTx, inTxAttempt, txmgrtypes.TxAttemptInsufficientFunds)
			actErr := inMemoryStore.UpdateTxAttemptInProgressToBroadcast(ctx, &inTx, inTxAttempt, txmgrtypes.TxAttemptInsufficientFunds)
			assert.ErrorContains(t, expErr, "new attempt state must be broadcast")
			assert.ErrorContains(t, actErr, "new attempt state must be broadcast")
		})

		t.Run("incorrect from address", func(t *testing.T) {
			inTx.FromAddress = utils.RandomAddress()
			inTx.State = commontxmgr.TxInProgress
			inTxAttempt.State = txmgrtypes.TxAttemptInProgress
			expErr := persistentStore.UpdateTxAttemptInProgressToBroadcast(ctx, &inTx, inTxAttempt, txmgrtypes.TxAttemptBroadcast)
			inTx.State = commontxmgr.TxInProgress
			inTxAttempt.State = txmgrtypes.TxAttemptInProgress
			actErr := inMemoryStore.UpdateTxAttemptInProgressToBroadcast(ctx, &inTx, inTxAttempt, txmgrtypes.TxAttemptBroadcast)
			assert.NoError(t, expErr)
			assert.NoError(t, actErr)
			inTx.FromAddress = fromAddress // reset
		})

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
	if exp.BroadcastAt != nil && act.BroadcastAt != nil {
		assert.Equal(t, exp.BroadcastAt.Unix(), act.BroadcastAt.Unix())
	} else {
		assert.Equal(t, exp.BroadcastAt, act.BroadcastAt)
	}
	if exp.InitialBroadcastAt != nil && act.InitialBroadcastAt != nil {
		assert.Equal(t, exp.InitialBroadcastAt.Unix(), act.InitialBroadcastAt.Unix())
	} else {
		assert.Equal(t, exp.InitialBroadcastAt, act.InitialBroadcastAt)
	}
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
	assert.Equal(t, exp.Tx, act.Tx)
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
