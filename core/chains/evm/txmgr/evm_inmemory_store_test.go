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

func TestInMemoryStore_MarkAllConfirmedMissingReceipt(t *testing.T) {
	t.Parallel()

	t.Run("successfully mark all confirmed missing receipt", func(t *testing.T) {
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

		// create transaction 0 that is unconfirmed (block 7)
		// Insert a transaction into persistent store
		blocknum := int64(7)
		inTx_0 := cltest.MustInsertUnconfirmedEthTx(t, persistentStore, 0, fromAddress)
		inTxAttempt_0 := newBroadcastLegacyEthTxAttempt(t, inTx_0.ID, int64(1))
		inTxAttempt_0.BroadcastBeforeBlockNum = &blocknum
		require.NoError(t, persistentStore.InsertTxAttempt(&inTxAttempt_0))
		assert.Equal(t, commontxmgr.TxUnconfirmed, inTx_0.State)
		// Insert the transaction into the in-memory store
		inTx_0.TxAttempts = []evmtxmgr.TxAttempt{inTxAttempt_0}
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))

		// create transaction 1 that is confirmed (block 77)
		inTx_1 := mustInsertConfirmedEthTxBySaveFetchedReceipts(t, persistentStore, fromAddress, 1, 77, *chainID)
		assert.Equal(t, commontxmgr.TxConfirmed, inTx_1.State)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))

		// mark transaction 0 as confirmed missing receipt
		err = inMemoryStore.MarkAllConfirmedMissingReceipt(testutils.Context(t), chainID)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(inTx_0.ID)
		require.NoError(t, err)

		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx_0.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]
		assert.Equal(t, commontxmgr.TxConfirmedMissingReceipt, actTx.State)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, actTx.TxAttempts[0].State)
		assertTxEqual(t, expTx, actTx)
	})

	t.Run("error parity for in-memory vs persistent store", func(t *testing.T) {
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

		// create transaction 0 that is unconfirmed (block 7)
		// Insert a transaction into persistent store
		blocknum := int64(7)
		inTx_0 := cltest.MustInsertUnconfirmedEthTx(t, persistentStore, 0, fromAddress)
		inTxAttempt_0 := newBroadcastLegacyEthTxAttempt(t, inTx_0.ID, int64(1))
		inTxAttempt_0.BroadcastBeforeBlockNum = &blocknum
		require.NoError(t, persistentStore.InsertTxAttempt(&inTxAttempt_0))
		assert.Equal(t, commontxmgr.TxUnconfirmed, inTx_0.State)
		// Insert the transaction into the in-memory store
		inTx_0.TxAttempts = []evmtxmgr.TxAttempt{inTxAttempt_0}
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))

		// create transaction 1 that is confirmed (block 77)
		inTx_1 := mustInsertConfirmedEthTxBySaveFetchedReceipts(t, persistentStore, fromAddress, 1, 77, *chainID)
		assert.Equal(t, commontxmgr.TxConfirmed, inTx_1.State)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))

		t.Run("wrong chain ID", func(t *testing.T) {
			wrongChainID := big.NewInt(1)
			expErr := persistentStore.MarkAllConfirmedMissingReceipt(testutils.Context(t), wrongChainID)
			actErr := inMemoryStore.MarkAllConfirmedMissingReceipt(testutils.Context(t), wrongChainID)
			assert.Equal(t, expErr, actErr)
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

	require.Equal(t, len(exp.TxAttempts), len(act.TxAttempts))
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
