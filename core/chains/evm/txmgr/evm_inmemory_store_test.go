package txmgr_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

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

func TestInMemoryStore_UpdateTxFatalError(t *testing.T) {
	t.Parallel()

	t.Run("successfully update transaction", func(t *testing.T) {
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
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		inTx.Error = null.StringFrom("no more toilet paper")
		err = inMemoryStore.UpdateTxFatalError(ctx, &inTx)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)

		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]

		assertTxEqual(t, expTx, actTx)
		assert.Equal(t, commontxmgr.TxFatalError, actTx.State)
	})
}

func TestInMemoryStore_DeleteInProgressAttempt(t *testing.T) {
	t.Parallel()

	t.Run("successfully replace tx attempt", func(t *testing.T) {
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
		inTx := mustInsertInProgressEthTxWithAttempt(t, persistentStore, 1, fromAddress)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		oldAttempt := inTx.TxAttempts[0]
		err = inMemoryStore.DeleteInProgressAttempt(ctx, oldAttempt)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)
		assert.Equal(t, 0, len(expTx.TxAttempts))

		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]
		assertTxEqual(t, expTx, actTx)
	})

	t.Run("error parity for in-memory vs persistent store", func(t *testing.T) {
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
		inTx := mustInsertInProgressEthTxWithAttempt(t, persistentStore, 124, fromAddress)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		oldAttempt := inTx.TxAttempts[0]
		t.Run("error when attempt is not in progress", func(t *testing.T) {
			oldAttempt.State = txmgrtypes.TxAttemptBroadcast
			expErr := persistentStore.DeleteInProgressAttempt(ctx, oldAttempt)
			actErr := inMemoryStore.DeleteInProgressAttempt(ctx, oldAttempt)
			assert.Equal(t, expErr, actErr)
			oldAttempt.State = txmgrtypes.TxAttemptInProgress
		})

		t.Run("error when attempt has 0 id", func(t *testing.T) {
			originalID := oldAttempt.ID
			oldAttempt.ID = 0
			expErr := persistentStore.DeleteInProgressAttempt(ctx, oldAttempt)
			actErr := inMemoryStore.DeleteInProgressAttempt(ctx, oldAttempt)
			assert.Equal(t, expErr, actErr)
			oldAttempt.ID = originalID
		})
	})
}

func TestInMemoryStore_ReapTxHistory(t *testing.T) {
	t.Parallel()

	t.Run("reap all confirmed txs", func(t *testing.T) {
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
		inTx_0 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 7, 1, fromAddress)
		r_0 := mustInsertEthReceipt(t, persistentStore, 1, utils.NewHash(), inTx_0.TxAttempts[0].Hash)
		inTx_0.TxAttempts[0].Receipts = append(inTx_0.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&r_0))
		inTx_1 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 8, 2, fromAddress)
		r_1 := mustInsertEthReceipt(t, persistentStore, 2, utils.NewHash(), inTx_1.TxAttempts[0].Hash)
		inTx_1.TxAttempts[0].Receipts = append(inTx_1.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&r_1))
		inTx_2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 9, 3, fromAddress)
		r_2 := mustInsertEthReceipt(t, persistentStore, 3, utils.NewHash(), inTx_2.TxAttempts[0].Hash)
		inTx_2.TxAttempts[0].Receipts = append(inTx_2.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&r_2))
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_2))

		minBlockNumberToKeep := int64(3)
		timeThreshold := inTx_2.CreatedAt
		expErr := persistentStore.ReapTxHistory(ctx, minBlockNumberToKeep, timeThreshold, chainID)
		actErr := inMemoryStore.ReapTxHistory(ctx, minBlockNumberToKeep, timeThreshold, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)

		fn := func(tx *evmtxmgr.Tx) bool { return true }
		// Check that the transactions were reaped in persistent store
		expTx_0, err := persistentStore.FindTxWithAttempts(ctx, inTx_0.ID)
		require.Error(t, err)
		require.Equal(t, int64(0), expTx_0.ID)
		expTx_1, err := persistentStore.FindTxWithAttempts(ctx, inTx_1.ID)
		require.Error(t, err)
		require.Equal(t, int64(0), expTx_1.ID)
		// Check that the transactions were reaped in in-memory store
		actTxs_0 := inMemoryStore.XXXTestFindTxs(nil, fn, inTx_0.ID)
		require.Equal(t, 0, len(actTxs_0))
		actTxs_1 := inMemoryStore.XXXTestFindTxs(nil, fn, inTx_1.ID)
		require.Equal(t, 0, len(actTxs_1))

		// Check that the transaction was not reaped
		expTx_2, err := persistentStore.FindTxWithAttempts(ctx, inTx_2.ID)
		require.NoError(t, err)
		require.Equal(t, inTx_2.ID, expTx_2.ID)
		actTxs_2 := inMemoryStore.XXXTestFindTxs(nil, fn, inTx_2.ID)
		require.Equal(t, 1, len(actTxs_2))
		assertTxEqual(t, expTx_2, actTxs_2[0])
	})
	t.Run("reap all fatal error txs", func(t *testing.T) {
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
		inTx_0 := cltest.NewEthTx(fromAddress)
		inTx_0.Error = null.StringFrom("something exploded")
		inTx_0.State = commontxmgr.TxFatalError
		inTx_0.CreatedAt = time.Unix(1000, 0)
		require.NoError(t, persistentStore.InsertTx(ctx, &inTx_0))
		inTx_1 := cltest.NewEthTx(fromAddress)
		inTx_1.Error = null.StringFrom("something exploded")
		inTx_1.State = commontxmgr.TxFatalError
		inTx_1.CreatedAt = time.Unix(2000, 0)
		require.NoError(t, persistentStore.InsertTx(ctx, &inTx_1))
		inTx_2 := cltest.NewEthTx(fromAddress)
		inTx_2.Error = null.StringFrom("something exploded")
		inTx_2.State = commontxmgr.TxFatalError
		inTx_2.CreatedAt = time.Unix(3000, 0)
		require.NoError(t, persistentStore.InsertTx(ctx, &inTx_2))
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_2))

		minBlockNumberToKeep := int64(3)
		timeThreshold := time.Unix(2500, 0) // Only reap txs created before this time
		expErr := persistentStore.ReapTxHistory(ctx, minBlockNumberToKeep, timeThreshold, chainID)
		actErr := inMemoryStore.ReapTxHistory(ctx, minBlockNumberToKeep, timeThreshold, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)

		fn := func(tx *evmtxmgr.Tx) bool { return true }
		// Check that the transactions were reaped in persistent store
		expTx_0, err := persistentStore.FindTxWithAttempts(ctx, inTx_0.ID)
		require.Error(t, err)
		require.Equal(t, int64(0), expTx_0.ID)
		expTx_1, err := persistentStore.FindTxWithAttempts(ctx, inTx_1.ID)
		require.Error(t, err)
		require.Equal(t, int64(0), expTx_1.ID)
		// Check that the transactions were reaped in in-memory store
		actTxs_0 := inMemoryStore.XXXTestFindTxs(nil, fn, inTx_0.ID)
		require.Equal(t, 0, len(actTxs_0))
		actTxs_1 := inMemoryStore.XXXTestFindTxs(nil, fn, inTx_1.ID)
		require.Equal(t, 0, len(actTxs_1))

		// Check that the transaction was not reaped
		expTx_2, err := persistentStore.FindTxWithAttempts(ctx, inTx_2.ID)
		require.NoError(t, err)
		require.Equal(t, inTx_2.ID, expTx_2.ID)
		actTxs_2 := inMemoryStore.XXXTestFindTxs(nil, fn, inTx_2.ID)
		require.Equal(t, 1, len(actTxs_2))
		assertTxEqual(t, expTx_2, actTxs_2[0])
	})
}

func TestInMemoryStore_MarkOldTxesMissingReceiptAsErrored(t *testing.T) {
	t.Parallel()
	blockNum := int64(10)
	finalityDepth := uint32(2)

	t.Run("successfully mark errored transaction", func(t *testing.T) {
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
		inTx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, persistentStore, 1, 7, time.Now(), fromAddress)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		err = inMemoryStore.MarkOldTxesMissingReceiptAsErrored(ctx, blockNum, finalityDepth, chainID)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)

		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]

		assertTxEqual(t, expTx, actTx)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, actTx.TxAttempts[0].State)
		assert.Equal(t, commontxmgr.TxFatalError, actTx.State)
	})
}

func TestInMemoryStore_UpdateTxForRebroadcast(t *testing.T) {
	t.Parallel()

	t.Run("delete all receipts for transaction", func(t *testing.T) {
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
		inTx := mustInsertConfirmedEthTxWithReceipt(t, persistentStore, fromAddress, 777, 1)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		txAttempt := inTx.TxAttempts[0]
		err = inMemoryStore.UpdateTxForRebroadcast(ctx, inTx, txAttempt)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)
		require.Len(t, expTx.TxAttempts, 1)

		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]

		assertTxEqual(t, expTx, actTx)
		assert.Equal(t, commontxmgr.TxUnconfirmed, actTx.State)
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, actTx.TxAttempts[0].State)
		assert.Nil(t, actTx.TxAttempts[0].BroadcastBeforeBlockNum)
		assert.Equal(t, 0, len(actTx.TxAttempts[0].Receipts))
	})

	t.Run("error parity for in-memory vs persistent store", func(t *testing.T) {
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
		inTx := mustInsertConfirmedEthTxWithReceipt(t, persistentStore, fromAddress, 777, 1)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		txAttempt := inTx.TxAttempts[0]

		t.Run("error when attempt is not in Broadcast state", func(t *testing.T) {
			txAttempt.State = txmgrtypes.TxAttemptInProgress
			expErr := persistentStore.UpdateTxForRebroadcast(ctx, inTx, txAttempt)
			actErr := inMemoryStore.UpdateTxForRebroadcast(ctx, inTx, txAttempt)
			assert.Error(t, expErr)
			assert.Error(t, actErr)
			txAttempt.State = txmgrtypes.TxAttemptBroadcast
		})
		t.Run("error when transaction is not in confirmed state", func(t *testing.T) {
			inTx.State = commontxmgr.TxUnconfirmed
			expErr := persistentStore.UpdateTxForRebroadcast(ctx, inTx, txAttempt)
			actErr := inMemoryStore.UpdateTxForRebroadcast(ctx, inTx, txAttempt)
			assert.Error(t, expErr)
			assert.Error(t, actErr)
			inTx.State = commontxmgr.TxConfirmed
		})
		t.Run("wrong fromAddress has no error", func(t *testing.T) {
			inTx.FromAddress = common.Address{}
			expErr := persistentStore.UpdateTxForRebroadcast(ctx, inTx, txAttempt)
			actErr := inMemoryStore.UpdateTxForRebroadcast(ctx, inTx, txAttempt)
			assert.Equal(t, expErr, actErr)
			assert.Nil(t, actErr)
			inTx.FromAddress = fromAddress
		})
	})
}

func TestInMemoryStore_MarkAllConfirmedMissingReceipt(t *testing.T) {
	t.Parallel()

	t.Run("successfully mark all confirmed missing receipt", func(t *testing.T) {
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

		// create transaction 0 that is unconfirmed (block 7)
		// Insert a transaction into persistent store
		blocknum := int64(7)
		inTx_0 := cltest.MustInsertUnconfirmedEthTx(t, persistentStore, 0, fromAddress)
		inTxAttempt_0 := newBroadcastLegacyEthTxAttempt(t, inTx_0.ID, int64(1))
		inTxAttempt_0.BroadcastBeforeBlockNum = &blocknum
		require.NoError(t, persistentStore.InsertTxAttempt(ctx, &inTxAttempt_0))
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
		err = inMemoryStore.MarkAllConfirmedMissingReceipt(ctx, chainID)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx_0.ID)
		require.NoError(t, err)

		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx_0.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]
		assert.Equal(t, commontxmgr.TxConfirmedMissingReceipt, actTx.State)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, actTx.TxAttempts[0].State)
		assertTxEqual(t, expTx, actTx)
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
	if exp.BroadcastAt != nil {
		require.NotNil(t, act.BroadcastAt)
		assert.Equal(t, exp.BroadcastAt.Unix(), act.BroadcastAt.Unix())
	} else {
		assert.Equal(t, exp.BroadcastAt, act.BroadcastAt)
	}
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
