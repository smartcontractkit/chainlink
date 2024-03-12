package txmgr_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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

func TestInMemoryStore_UpdateTxUnstartedToInProgress(t *testing.T) {
	t.Parallel()

	t.Run("successfully updates unstarted tx to inprogress", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
		persistentStore := cltest.NewTestTxStore(t, db, dbcfg)
		kst := cltest.NewKeyStore(t, db, dbcfg)
		_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		lggr := logger.TestSugared(t)
		chainID := ethClient.ConfiguredChainID()
		ctx := context.Background()

		inMemoryStore, err := commontxmgr.NewInMemoryStore[
			*big.Int,
			common.Address, common.Hash, common.Hash,
			*evmtypes.Receipt,
			evmtypes.Nonce,
			evmgas.EvmFee,
		](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
		require.NoError(t, err)

		nonce := evmtypes.Nonce(123)
		// Insert a transaction into persistent store
		inTx := mustCreateUnstartedGeneratedTx(t, persistentStore, fromAddress, chainID)
		inTx.Sequence = &nonce
		inTxAttempt := cltest.NewLegacyEthTxAttempt(t, inTx.ID)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		// Update the transaction to in-progress
		require.NoError(t, inMemoryStore.UpdateTxUnstartedToInProgress(testutils.Context(t), &inTx, &inTxAttempt))

		expTx, err := persistentStore.FindTxWithAttempts(inTx.ID)
		require.NoError(t, err)
		assert.Equal(t, commontxmgr.TxInProgress, expTx.State)
		assert.Equal(t, 1, len(expTx.TxAttempts))

		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]
		assertTxEqual(t, expTx, actTx)
		assert.Equal(t, commontxmgr.TxInProgress, actTx.State)
	})

	t.Run("wrong input error scenarios", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
		persistentStore := cltest.NewTestTxStore(t, db, dbcfg)
		kst := cltest.NewKeyStore(t, db, dbcfg)
		_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())

		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
		lggr := logger.TestSugared(t)
		chainID := ethClient.ConfiguredChainID()
		ctx := context.Background()

		inMemoryStore, err := commontxmgr.NewInMemoryStore[
			*big.Int,
			common.Address, common.Hash, common.Hash,
			*evmtypes.Receipt,
			evmtypes.Nonce,
			evmgas.EvmFee,
		](ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
		require.NoError(t, err)

		nonce1 := evmtypes.Nonce(1)
		nonce2 := evmtypes.Nonce(2)
		// Insert a transaction into persistent store
		inTx1 := mustCreateUnstartedGeneratedTx(t, persistentStore, fromAddress, chainID)
		inTx2 := mustCreateUnstartedGeneratedTx(t, persistentStore, fromAddress, chainID)
		inTx1.Sequence = &nonce1
		inTx2.Sequence = &nonce2
		inTxAttempt1 := cltest.NewLegacyEthTxAttempt(t, inTx1.ID)
		inTxAttempt2 := cltest.NewLegacyEthTxAttempt(t, inTx2.ID)
		// Insert the transaction into the in-memory store
		//inTx2 := cltest.NewEthTx(fromAddress)
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx1))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx2))

		// sequence nil
		inTx1.Sequence = nil
		inTx2.Sequence = nil
		expErr := persistentStore.UpdateTxUnstartedToInProgress(testutils.Context(t), &inTx1, &inTxAttempt1)
		actErr := inMemoryStore.UpdateTxUnstartedToInProgress(testutils.Context(t), &inTx2, &inTxAttempt2)
		assert.Equal(t, expErr, actErr)
		assert.Error(t, actErr)
		assert.Error(t, expErr)
		inTx1.Sequence = &nonce1 // reset
		inTx2.Sequence = &nonce2 // reset

		// tx not in unstarted state
		inTx1.State = commontxmgr.TxInProgress
		inTx2.State = commontxmgr.TxInProgress
		expErr = persistentStore.UpdateTxUnstartedToInProgress(testutils.Context(t), &inTx1, &inTxAttempt1)
		actErr = inMemoryStore.UpdateTxUnstartedToInProgress(testutils.Context(t), &inTx2, &inTxAttempt2)
		assert.Error(t, actErr)
		assert.Error(t, expErr)
		inTx1.State = commontxmgr.TxUnstarted // reset
		inTx2.State = commontxmgr.TxUnstarted // reset

		// tx attempt not in in-progress state
		inTxAttempt1.State = txmgrtypes.TxAttemptBroadcast
		inTxAttempt2.State = txmgrtypes.TxAttemptBroadcast
		expErr = persistentStore.UpdateTxUnstartedToInProgress(testutils.Context(t), &inTx1, &inTxAttempt1)
		actErr = inMemoryStore.UpdateTxUnstartedToInProgress(testutils.Context(t), &inTx2, &inTxAttempt2)
		assert.Equal(t, expErr, actErr)
		assert.Error(t, actErr)
		assert.Error(t, expErr)
		inTxAttempt1.State = txmgrtypes.TxAttemptInProgress // reset
		inTxAttempt2.State = txmgrtypes.TxAttemptInProgress // reset

		// wrong from address
		inTx1.FromAddress = cltest.NewEIP55Address().Address()
		inTx2.FromAddress = cltest.NewEIP55Address().Address()
		expErr = persistentStore.UpdateTxUnstartedToInProgress(testutils.Context(t), &inTx1, &inTxAttempt1)
		actErr = inMemoryStore.UpdateTxUnstartedToInProgress(testutils.Context(t), &inTx2, &inTxAttempt2)
		assert.NoError(t, actErr)
		assert.NoError(t, expErr)
		inTx1.FromAddress = fromAddress // reset
		inTx2.FromAddress = fromAddress // reset
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
