package txmgr_test

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	clnull "github.com/smartcontractkit/chainlink-common/pkg/utils/null"
	commontxmgr "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"

	evmassets "github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmgas "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func TestInMemoryStore_FindTxesPendingCallback(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	head := evmtypes.Head{
		Hash:   evmutils.NewHash(),
		Number: 10,
		Parent: &evmtypes.Head{
			Hash:   evmutils.NewHash(),
			Number: 9,
			Parent: &evmtypes.Head{
				Number: 8,
				Hash:   evmutils.NewHash(),
				Parent: nil,
			},
		},
	}
	minConfirmations := int64(2)

	pgtest.MustExec(t, db, `SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`)
	// insert the transaction into the persistent store
	// Suspended run waiting for callback
	run1 := cltest.MustInsertPipelineRun(t, db)
	tr1 := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run1.ID)
	pgtest.MustExec(t, db, `UPDATE pipeline_runs SET state = 'suspended' WHERE id = $1`, run1.ID)
	inTx_0 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 3, 1, fromAddress)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET meta='{"FailOnRevert": true}'`)
	attempt1 := inTx_0.TxAttempts[0]
	r_0 := mustInsertEthReceipt(t, persistentStore, head.Number-minConfirmations, head.Hash, attempt1.Hash)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr1.ID, minConfirmations, inTx_0.ID)
	failOnRevert := null.BoolFrom(true)
	b, err := json.Marshal(evmtxmgr.TxMeta{FailOnRevert: failOnRevert})
	require.NoError(t, err)
	meta := sqlutil.JSON(b)
	inTx_0.Meta = &meta
	inTx_0.TxAttempts[0].Receipts = append(inTx_0.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&r_0))
	inTx_0.MinConfirmations = clnull.Uint32From(uint32(minConfirmations))
	inTx_0.PipelineTaskRunID = uuid.NullUUID{UUID: tr1.ID, Valid: true}
	inTx_0.SignalCallback = true

	// Callback to pipeline service completed. Should be ignored
	run2 := cltest.MustInsertPipelineRunWithStatus(t, db, 0, pipeline.RunStatusCompleted)
	tr2 := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run2.ID)
	inTx_1 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 4, 1, fromAddress)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET meta='{"FailOnRevert": false}'`)
	attempt2 := inTx_1.TxAttempts[0]
	r_1 := mustInsertEthReceipt(t, persistentStore, head.Number-minConfirmations, head.Hash, attempt2.Hash)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE, callback_completed = TRUE WHERE id = $3`, &tr2.ID, minConfirmations, inTx_1.ID)
	failOnRevert = null.BoolFrom(false)
	b, err = json.Marshal(evmtxmgr.TxMeta{FailOnRevert: failOnRevert})
	require.NoError(t, err)
	meta = sqlutil.JSON(b)
	inTx_1.Meta = &meta
	inTx_1.TxAttempts[0].Receipts = append(inTx_1.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&r_1))
	inTx_1.MinConfirmations = clnull.Uint32From(uint32(minConfirmations))
	inTx_1.PipelineTaskRunID = uuid.NullUUID{UUID: tr2.ID, Valid: true}
	inTx_1.SignalCallback = true
	inTx_1.CallbackCompleted = true

	// Suspended run younger than minConfirmations. Should be ignored
	run3 := cltest.MustInsertPipelineRun(t, db)
	tr3 := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run3.ID)
	pgtest.MustExec(t, db, `UPDATE pipeline_runs SET state = 'suspended' WHERE id = $1`, run3.ID)
	inTx_2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 5, 1, fromAddress)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET meta='{"FailOnRevert": false}'`)
	attempt3 := inTx_2.TxAttempts[0]
	r_2 := mustInsertEthReceipt(t, persistentStore, head.Number, head.Hash, attempt3.Hash)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr3.ID, minConfirmations, inTx_2.ID)
	failOnRevert = null.BoolFrom(false)
	b, err = json.Marshal(evmtxmgr.TxMeta{FailOnRevert: failOnRevert})
	require.NoError(t, err)
	meta = sqlutil.JSON(b)
	inTx_2.Meta = &meta
	inTx_2.TxAttempts[0].Receipts = append(inTx_2.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&r_2))
	inTx_2.MinConfirmations = clnull.Uint32From(uint32(minConfirmations))
	inTx_2.PipelineTaskRunID = uuid.NullUUID{UUID: tr3.ID, Valid: true}
	inTx_2.SignalCallback = true

	// Tx not marked for callback. Should be ignore
	inTx_3 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 6, 1, fromAddress)
	attempt4 := inTx_3.TxAttempts[0]
	r_3 := mustInsertEthReceipt(t, persistentStore, head.Number, head.Hash, attempt4.Hash)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET min_confirmations = $1 WHERE id = $2`, minConfirmations, inTx_3.ID)
	inTx_3.TxAttempts[0].Receipts = append(inTx_3.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&r_3))
	inTx_3.MinConfirmations = clnull.Uint32From(uint32(minConfirmations))

	// Unconfirmed Tx without receipts. Should be ignored
	inTx_4 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 7, 1, fromAddress)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET min_confirmations = $1 WHERE id = $2`, minConfirmations, inTx_4.ID)
	inTx_4.MinConfirmations = clnull.Uint32From(uint32(minConfirmations))

	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_2))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_3))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_4))

	tcs := []struct {
		name         string
		inHeadNumber int64
		inChainID    *big.Int

		hasErr      bool
		hasReceipts bool
	}{
		{"successfully finds receipts", head.Number, chainID, false, true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actReceipts, actErr := inMemoryStore.FindTxesPendingCallback(ctx, tc.inHeadNumber, tc.inChainID)
			expReceipts, expErr := persistentStore.FindTxesPendingCallback(ctx, tc.inHeadNumber, tc.inChainID)
			require.Equal(t, expErr, actErr)
			if tc.hasErr {
				require.NotNil(t, expErr)
				require.NotNil(t, actErr)
			} else {
				require.Nil(t, expErr)
				require.Nil(t, actErr)
			}
			if tc.hasReceipts {
				require.NotEqual(t, 0, len(expReceipts))
				assert.NotEqual(t, 0, len(actReceipts))
				require.Equal(t, len(expReceipts), len(actReceipts))
				for i := 0; i < len(expReceipts); i++ {
					assert.Equal(t, expReceipts[i].ID, actReceipts[i].ID)
					assert.Equal(t, expReceipts[i].FailOnRevert, actReceipts[i].FailOnRevert)
					assertChainReceiptEqual(t, expReceipts[i].Receipt, actReceipts[i].Receipt)
				}
			} else {
				require.Equal(t, 0, len(expReceipts))
				require.Equal(t, 0, len(actReceipts))
			}
		})
	}
}

func TestInMemoryStore_FindTxAttemptsRequiringResend(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// insert the transaction into the persistent store
	inTx_1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, 1, fromAddress, time.Unix(1616509200, 0))
	inTx_3 := mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, persistentStore, 3, fromAddress, time.Unix(1616509400, 0))
	inTx_0 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, 0, fromAddress, time.Unix(1616509100, 0))
	inTx_2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, 2, fromAddress, time.Unix(1616509300, 0))
	// modify the attempts
	attempt0_2 := newBroadcastLegacyEthTxAttempt(t, inTx_0.ID)
	attempt0_2.TxFee = evmgas.EvmFee{Legacy: evmassets.NewWeiI(10)}
	require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt0_2))

	attempt2_2 := newInProgressLegacyEthTxAttempt(t, inTx_2.ID)
	attempt2_2.TxFee = evmgas.EvmFee{Legacy: evmassets.NewWeiI(10)}
	require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt2_2))

	attempt3_2 := cltest.NewDynamicFeeEthTxAttempt(t, inTx_3.ID)
	attempt3_2.TxFee.DynamicTipCap = evmassets.NewWeiI(10)
	attempt3_2.TxFee.DynamicFeeCap = evmassets.NewWeiI(20)
	attempt3_2.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt3_2))
	attempt3_4 := cltest.NewDynamicFeeEthTxAttempt(t, inTx_3.ID)
	attempt3_4.TxFee.DynamicTipCap = evmassets.NewWeiI(30)
	attempt3_4.TxFee.DynamicFeeCap = evmassets.NewWeiI(40)
	attempt3_4.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt3_4))
	attempt3_3 := cltest.NewDynamicFeeEthTxAttempt(t, inTx_3.ID)
	attempt3_3.TxFee.DynamicTipCap = evmassets.NewWeiI(20)
	attempt3_3.TxFee.DynamicFeeCap = evmassets.NewWeiI(30)
	attempt3_3.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt3_3))
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_2))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_3))

	tcs := []struct {
		name                      string
		inOlderThan               time.Time
		inMaxInFlightTransactions uint32
		inChainID                 *big.Int
		inFromAddress             common.Address

		hasErr        bool
		hasTxAttempts bool
	}{
		{"finds nothing if transactions from a different key", time.Now(), 10, chainID, evmutils.RandomAddress(), false, false},
		{"returns the highest price attempt for each transaction that was last broadcast before or on the given time", time.Unix(1616509200, 0), 0, chainID, fromAddress, false, true},
		{"returns the highest price attempt for EIP-1559 transactions", time.Unix(1616509400, 0), 0, chainID, fromAddress, false, true},
		{"applies limit", time.Unix(1616509200, 0), 1, chainID, fromAddress, false, true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actTxAttempts, actErr := inMemoryStore.FindTxAttemptsRequiringResend(ctx, tc.inOlderThan, tc.inMaxInFlightTransactions, tc.inChainID, tc.inFromAddress)
			expTxAttempts, expErr := persistentStore.FindTxAttemptsRequiringResend(ctx, tc.inOlderThan, tc.inMaxInFlightTransactions, tc.inChainID, tc.inFromAddress)
			require.Equal(t, expErr, actErr)
			if tc.hasErr {
				require.NotNil(t, expErr)
				require.NotNil(t, actErr)
			} else {
				require.Nil(t, expErr)
				require.Nil(t, actErr)
			}
			if tc.hasTxAttempts {
				require.NotEqual(t, 0, len(expTxAttempts))
				assert.NotEqual(t, 0, len(actTxAttempts))
				require.Equal(t, len(expTxAttempts), len(actTxAttempts))
				for i := 0; i < len(expTxAttempts); i++ {
					assertTxAttemptEqual(t, expTxAttempts[i], actTxAttempts[i])
				}
			} else {
				require.Equal(t, 0, len(expTxAttempts))
				require.Equal(t, 0, len(actTxAttempts))
			}
		})
	}
}

func TestInMemoryStore_FindTxesWithMetaFieldByReceiptBlockNum(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// initialize the Meta field which is sqlutil.JSON
	subID := uint64(123)
	b, err := json.Marshal(evmtxmgr.TxMeta{SubID: &subID})
	require.NoError(t, err)
	meta := sqlutil.JSON(b)
	timeNow := time.Now()
	nonce := evmtypes.Nonce(123)
	blockNum := int64(3)
	broadcastBeforeBlockNum := int64(3)
	// initialize transactions
	inTx_0 := cltest.NewEthTx(fromAddress)
	inTx_0.BroadcastAt = &timeNow
	inTx_0.InitialBroadcastAt = &timeNow
	inTx_0.Sequence = &nonce
	inTx_0.State = commontxmgr.TxConfirmed
	inTx_0.MinConfirmations.SetValid(6)
	inTx_0.Meta = &meta
	// insert the transaction into the persistent store
	require.NoError(t, persistentStore.InsertTx(ctx, &inTx_0))
	attempt := cltest.NewLegacyEthTxAttempt(t, inTx_0.ID)
	attempt.BroadcastBeforeBlockNum = &broadcastBeforeBlockNum
	attempt.State = txmgrtypes.TxAttemptBroadcast
	require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt))
	inTx_0.TxAttempts = append(inTx_0.TxAttempts, attempt)
	// insert the transaction receipt into the persistent store
	rec_0 := mustInsertEthReceipt(t, persistentStore, 3, evmutils.NewHash(), inTx_0.TxAttempts[0].Hash)
	inTx_0.TxAttempts[0].Receipts = append(inTx_0.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&rec_0))
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))

	tcs := []struct {
		name        string
		inMetaField string
		inBlockNum  int64
		inChainID   *big.Int

		hasErr bool
		hasTxs bool
	}{
		{"successfully finds tx", "SubId", blockNum, chainID, false, true},
		{"unknown meta_field: finds no txs", "unknown", blockNum, chainID, false, false},
		{"incorrect meta_field: finds no txs", "MaxLink", blockNum, chainID, false, false},
		{"incorrect blockNum: finds no txs", "SubId", 12, chainID, false, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actTxs, actErr := inMemoryStore.FindTxesWithMetaFieldByReceiptBlockNum(ctx, tc.inMetaField, tc.inBlockNum, tc.inChainID)
			expTxs, expErr := persistentStore.FindTxesWithMetaFieldByReceiptBlockNum(ctx, tc.inMetaField, tc.inBlockNum, tc.inChainID)
			require.Equal(t, expErr, actErr)
			if tc.hasErr {
				require.NotNil(t, expErr)
				require.NotNil(t, actErr)
			} else {
				require.Nil(t, expErr)
				require.Nil(t, actErr)
			}
			if tc.hasTxs {
				require.NotEqual(t, 0, len(expTxs))
				assert.NotEqual(t, 0, len(actTxs))
				require.Equal(t, len(expTxs), len(actTxs))
				for i := 0; i < len(expTxs); i++ {
					assertTxEqual(t, *expTxs[i], *actTxs[i])
				}
			} else {
				require.Equal(t, 0, len(expTxs))
				require.Equal(t, 0, len(actTxs))
			}
		})
	}
}

func TestInMemoryStore_FindTxesWithMetaFieldByStates(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// initialize the Meta field which is sqlutil.JSON
	subID := uint64(123)
	b, err := json.Marshal(evmtxmgr.TxMeta{SubID: &subID})
	require.NoError(t, err)
	meta := sqlutil.JSON(b)
	// initialize transactions
	inTx_0 := cltest.NewEthTx(fromAddress)
	inTx_0.Meta = &meta
	// insert the transaction into the persistent store
	require.NoError(t, persistentStore.InsertTx(ctx, &inTx_0))
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))

	tcs := []struct {
		name        string
		inMetaField string
		inStates    []txmgrtypes.TxState
		inChainID   *big.Int

		hasErr bool
		hasTxs bool
	}{
		{"successfully finds tx", "SubId", []txmgrtypes.TxState{commontxmgr.TxUnstarted}, chainID, false, true},
		{"incorrect state: finds no txs", "SubId", []txmgrtypes.TxState{commontxmgr.TxConfirmed}, chainID, false, false},
		{"unknown meta_field: finds no txs", "unknown", []txmgrtypes.TxState{commontxmgr.TxUnstarted}, chainID, false, false},
		{"incorrect meta_field: finds no txs", "MaxLink", []txmgrtypes.TxState{commontxmgr.TxUnstarted}, chainID, false, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actTxs, actErr := inMemoryStore.FindTxesWithMetaFieldByStates(ctx, tc.inMetaField, tc.inStates, tc.inChainID)
			expTxs, expErr := persistentStore.FindTxesWithMetaFieldByStates(ctx, tc.inMetaField, tc.inStates, tc.inChainID)
			require.Equal(t, expErr, actErr)
			if !tc.hasErr {
				require.Nil(t, expErr)
				require.Nil(t, actErr)
			}
			if tc.hasTxs {
				require.NotEqual(t, 0, len(expTxs))
				assert.NotEqual(t, 0, len(actTxs))
				require.Equal(t, len(expTxs), len(actTxs))
				for i := 0; i < len(expTxs); i++ {
					assertTxEqual(t, *expTxs[i], *actTxs[i])
				}
			} else {
				require.Equal(t, 0, len(expTxs))
				require.Equal(t, 0, len(actTxs))
			}
		})
	}
}

func TestInMemoryStore_FindTxesByMetaFieldAndStates(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// initialize the Meta field which is sqlutil.JSON
	subID := uint64(123)
	b, err := json.Marshal(evmtxmgr.TxMeta{SubID: &subID})
	require.NoError(t, err)
	meta := sqlutil.JSON(b)
	// initialize transactions
	inTx_0 := cltest.NewEthTx(fromAddress)
	inTx_0.Meta = &meta
	// insert the transaction into the persistent store
	require.NoError(t, persistentStore.InsertTx(ctx, &inTx_0))
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))

	tcs := []struct {
		name        string
		inMetaField string
		inMetaValue string
		inStates    []txmgrtypes.TxState
		inChainID   *big.Int

		hasErr bool
		hasTxs bool
	}{
		{"successfully finds tx", "SubId", "123", []txmgrtypes.TxState{commontxmgr.TxUnstarted}, chainID, false, true},
		{"incorrect state: finds no txs", "SubId", "123", []txmgrtypes.TxState{commontxmgr.TxConfirmed}, chainID, false, false},
		{"incorrect meta_value: finds no txs", "SubId", "incorrect", []txmgrtypes.TxState{commontxmgr.TxUnstarted}, chainID, false, false},
		{"unknown meta_field: finds no txs", "unknown", "123", []txmgrtypes.TxState{commontxmgr.TxUnstarted}, chainID, false, false},
		{"incorrect meta_field: finds no txs", "JobID", "123", []txmgrtypes.TxState{commontxmgr.TxUnstarted}, chainID, false, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actTxs, actErr := inMemoryStore.FindTxesByMetaFieldAndStates(ctx, tc.inMetaField, tc.inMetaValue, tc.inStates, tc.inChainID)
			expTxs, expErr := persistentStore.FindTxesByMetaFieldAndStates(ctx, tc.inMetaField, tc.inMetaValue, tc.inStates, tc.inChainID)
			require.Equal(t, expErr, actErr)
			if !tc.hasErr {
				require.Nil(t, expErr)
				require.Nil(t, actErr)
			}
			if tc.hasTxs {
				require.NotEqual(t, 0, len(expTxs))
				assert.NotEqual(t, 0, len(actTxs))
				require.Equal(t, len(expTxs), len(actTxs))
				for i := 0; i < len(expTxs); i++ {
					assertTxEqual(t, *expTxs[i], *actTxs[i])
				}
			} else {
				require.Equal(t, 0, len(expTxs))
				require.Equal(t, 0, len(actTxs))
			}
		})
	}
}

func TestInMemoryStore_FindTxWithIdempotencyKey(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	idempotencyKey := "777"
	inTx := cltest.NewEthTx(fromAddress)
	inTx.IdempotencyKey = &idempotencyKey
	// insert the transaction into the persistent store
	require.NoError(t, persistentStore.InsertTx(ctx, &inTx))
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

	tcs := []struct {
		name             string
		inIdempotencyKey string
		inChainID        *big.Int

		hasErr bool
		hasTx  bool
	}{
		{"no idempotency key", "", chainID, false, false},
		{"wrong idempotency key", "wrong", chainID, false, false},
		{"finds tx with idempotency key", idempotencyKey, chainID, false, true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actTx, actErr := inMemoryStore.FindTxWithIdempotencyKey(ctx, tc.inIdempotencyKey, tc.inChainID)
			expTx, expErr := persistentStore.FindTxWithIdempotencyKey(ctx, tc.inIdempotencyKey, tc.inChainID)
			require.Equal(t, expErr, actErr)
			if !tc.hasErr {
				require.Nil(t, actErr)
				require.Nil(t, expErr)
			}
			if tc.hasTx {
				require.NotNil(t, actTx)
				require.NotNil(t, expTx)
				assertTxEqual(t, *expTx, *actTx)
			} else {
				require.Nil(t, actTx)
				require.Nil(t, expTx)
			}
		})
	}
}

func TestInMemoryStore_CheckTxQueueCapacity(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// insert the transaction into the persistent store
	// insert the transaction into the in-memory store
	tx1 := cltest.NewEthTx(fromAddress)
	require.NoError(t, persistentStore.InsertTx(ctx, &tx1))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &tx1))
	tx2 := cltest.NewEthTx(fromAddress)
	require.NoError(t, persistentStore.InsertTx(ctx, &tx2))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &tx2))

	tcs := []struct {
		name           string
		inFromAddress  common.Address
		inMaxQueuedTxs uint64
		inChainID      *big.Int

		hasErr bool
	}{
		{"capacity reached", fromAddress, 2, chainID, true},
		{"above capacity", fromAddress, 1, chainID, true},
		{"below capacity", fromAddress, 3, chainID, false},
		{"wrong address", common.Address{}, 2, chainID, false},
		{"max queued txs is 0", fromAddress, 0, chainID, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actErr := inMemoryStore.CheckTxQueueCapacity(ctx, tc.inFromAddress, tc.inMaxQueuedTxs, tc.inChainID)
			expErr := persistentStore.CheckTxQueueCapacity(ctx, tc.inFromAddress, tc.inMaxQueuedTxs, tc.inChainID)
			if tc.hasErr {
				require.NotNil(t, expErr)
				require.NotNil(t, actErr)
			} else {
				require.NoError(t, expErr)
				require.NoError(t, actErr)
			}
		})
	}
}

func TestInMemoryStore_CountUnstartedTransactions(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// insert the transaction into the persistent store
	// insert the transaction into the in-memory store
	tx1 := cltest.NewEthTx(fromAddress)
	require.NoError(t, persistentStore.InsertTx(ctx, &tx1))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &tx1))
	tx2 := cltest.NewEthTx(fromAddress)
	require.NoError(t, persistentStore.InsertTx(ctx, &tx2))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &tx2))

	tcs := []struct {
		name          string
		inFromAddress common.Address
		inChainID     *big.Int

		expUnstartedCount uint32
		hasErr            bool
	}{
		{"return correct total transactions", fromAddress, chainID, 2, false},
		{"invalid address", common.Address{}, chainID, 0, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actMemoryCount, actErr := inMemoryStore.CountUnstartedTransactions(ctx, tc.inFromAddress, tc.inChainID)
			actPersistentCount, expErr := persistentStore.CountUnstartedTransactions(ctx, tc.inFromAddress, tc.inChainID)
			if tc.hasErr {
				require.NotNil(t, expErr)
				require.NotNil(t, actErr)
			} else {
				require.NoError(t, expErr)
				require.NoError(t, actErr)
			}
			assert.Equal(t, tc.expUnstartedCount, actMemoryCount)
			assert.Equal(t, tc.expUnstartedCount, actPersistentCount)
		})
	}
}

func TestInMemoryStore_CountUnconfirmedTransactions(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// initialize unconfirmed transactions
	inNonces := []int64{1, 2, 3}
	for _, inNonce := range inNonces {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertUnconfirmedEthTx(t, persistentStore, inNonce, fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))
	}

	tcs := []struct {
		name          string
		inFromAddress common.Address
		inChainID     *big.Int

		expUnconfirmedCount uint32
		hasErr              bool
	}{
		{"return correct total transactions", fromAddress, chainID, 3, false},
		{"invalid address", common.Address{}, chainID, 0, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actMemoryCount, actErr := inMemoryStore.CountUnconfirmedTransactions(ctx, tc.inFromAddress, tc.inChainID)
			actPersistentCount, expErr := persistentStore.CountUnconfirmedTransactions(ctx, tc.inFromAddress, tc.inChainID)
			if tc.hasErr {
				require.NotNil(t, expErr)
				require.NotNil(t, actErr)
			} else {
				require.NoError(t, expErr)
				require.NoError(t, actErr)
			}
			assert.Equal(t, tc.expUnconfirmedCount, actMemoryCount)
			assert.Equal(t, tc.expUnconfirmedCount, actPersistentCount)
		})
	}
}

func TestInMemoryStore_FindTxAttemptsConfirmedMissingReceipt(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// initialize transactions
	inTxDatas := []struct {
		nonce                   int64
		broadcastBeforeBlockNum int64
		broadcastAt             time.Time
	}{
		{0, 1, time.Unix(1616509300, 0)},
		{1, 1, time.Unix(1616509400, 0)},
		{2, 1, time.Unix(1616509500, 0)},
	}
	for _, inTxData := range inTxDatas {
		// insert the transaction into the persistent store
		inTx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
			t, persistentStore, inTxData.nonce, inTxData.broadcastBeforeBlockNum,
			inTxData.broadcastAt, fromAddress,
		)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))
	}

	tcs := []struct {
		name      string
		inChainID *big.Int

		expTxAttemptsCount int
		hasError           bool
	}{
		{"finds tx attempts confirmed missing receipt", chainID, 3, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actTxAttempts, actErr := inMemoryStore.FindTxAttemptsConfirmedMissingReceipt(ctx, tc.inChainID)
			expTxAttempts, expErr := persistentStore.FindTxAttemptsConfirmedMissingReceipt(ctx, tc.inChainID)
			if tc.hasError {
				require.NotNil(t, actErr)
				require.NotNil(t, expErr)
			} else {
				require.NoError(t, actErr)
				require.NoError(t, expErr)
				require.Equal(t, tc.expTxAttemptsCount, len(expTxAttempts))
				require.Equal(t, tc.expTxAttemptsCount, len(actTxAttempts))
				for i := 0; i < len(expTxAttempts); i++ {
					assertTxAttemptEqual(t, expTxAttempts[i], actTxAttempts[i])
				}
			}
		})
	}
}

func TestInMemoryStore_FindTxAttemptsRequiringReceiptFetch(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	// initialize transactions
	inTxDatas := []struct {
		nonce                   int64
		broadcastBeforeBlockNum int64
		broadcastAt             time.Time
	}{
		{0, 1, time.Unix(1616509300, 0)},
		{1, 1, time.Unix(1616509400, 0)},
		{2, 1, time.Unix(1616509500, 0)},
	}
	for _, inTxData := range inTxDatas {
		// insert the transaction into the persistent store
		inTx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
			t, persistentStore, inTxData.nonce, inTxData.broadcastBeforeBlockNum,
			inTxData.broadcastAt, fromAddress,
		)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))
	}

	tcs := []struct {
		name      string
		inChainID *big.Int

		expTxAttemptsCount int
		hasError           bool
	}{
		{"finds tx attempts requiring receipt fetch", chainID, 3, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			expTxAttempts, expErr := persistentStore.FindTxAttemptsRequiringReceiptFetch(ctx, tc.inChainID)
			actTxAttempts, actErr := inMemoryStore.FindTxAttemptsRequiringReceiptFetch(ctx, tc.inChainID)
			if tc.hasError {
				require.NotNil(t, actErr)
				require.NotNil(t, expErr)
			} else {
				require.NoError(t, actErr)
				require.NoError(t, expErr)
				require.Equal(t, tc.expTxAttemptsCount, len(expTxAttempts))
				require.Equal(t, tc.expTxAttemptsCount, len(actTxAttempts))
				for i := 0; i < len(expTxAttempts); i++ {
					assertTxAttemptEqual(t, expTxAttempts[i], actTxAttempts[i])
				}
			}
		})
	}
}

func TestInMemoryStore_GetInProgressTxAttempts(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("gets 0 in progress transaction", func(t *testing.T) {
		expTxAttempts, expErr := persistentStore.GetInProgressTxAttempts(ctx, fromAddress, chainID)
		actTxAttempts, actErr := inMemoryStore.GetInProgressTxAttempts(ctx, fromAddress, chainID)
		require.NoError(t, actErr)
		require.NoError(t, expErr)
		assert.Equal(t, len(expTxAttempts), len(actTxAttempts))
	})

	t.Run("gets 1 in progress transaction", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := mustInsertUnconfirmedEthTxWithAttemptState(t, persistentStore, int64(7), fromAddress, txmgrtypes.TxAttemptInProgress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expTxAttempts, expErr := persistentStore.GetInProgressTxAttempts(ctx, fromAddress, chainID)
		actTxAttempts, actErr := inMemoryStore.GetInProgressTxAttempts(ctx, fromAddress, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.Equal(t, len(expTxAttempts), len(actTxAttempts))
		for i := 0; i < len(expTxAttempts); i++ {
			assertTxAttemptEqual(t, expTxAttempts[i], actTxAttempts[i])
		}
	})
}

func TestInMemoryStore_HasInProgressTransaction(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no in progress transaction", func(t *testing.T) {
		expExists, expErr := persistentStore.HasInProgressTransaction(ctx, fromAddress, chainID)
		actExists, actErr := inMemoryStore.HasInProgressTransaction(ctx, fromAddress, chainID)
		require.NoError(t, actErr)
		require.NoError(t, expErr)
		assert.Equal(t, expExists, actExists)
	})

	t.Run("has an in progress transaction", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := mustInsertInProgressEthTxWithAttempt(t, persistentStore, 7, fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expExists, expErr := persistentStore.HasInProgressTransaction(ctx, fromAddress, chainID)
		actExists, actErr := inMemoryStore.HasInProgressTransaction(ctx, fromAddress, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.Equal(t, expExists, actExists)
	})
}

func TestInMemoryStore_GetTxByID(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no transaction", func(t *testing.T) {
		expTx, expErr := persistentStore.GetTxByID(ctx, 0)
		actTx, actErr := inMemoryStore.GetTxByID(ctx, 0)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Nil(t, expTx)
		assert.Nil(t, actTx)
	})

	t.Run("successfully get transaction by ID", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := mustInsertInProgressEthTxWithAttempt(t, persistentStore, 7, fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expTx, expErr := persistentStore.GetTxByID(ctx, inTx.ID)
		actTx, actErr := inMemoryStore.GetTxByID(ctx, inTx.ID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.NotNil(t, expTx)
		require.NotNil(t, actTx)
		assertTxEqual(t, *expTx, *actTx)
	})
}

func TestInMemoryStore_FindTxWithSequence(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no results", func(t *testing.T) {
		expTx, expErr := persistentStore.FindTxWithSequence(ctx, fromAddress, evmtypes.Nonce(666))
		actTx, actErr := inMemoryStore.FindTxWithSequence(ctx, fromAddress, evmtypes.Nonce(666))
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Nil(t, expTx)
		assert.Nil(t, actTx)
	})

	t.Run("successfully get transaction by ID", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 666, 1, fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expTx, expErr := persistentStore.FindTxWithSequence(ctx, fromAddress, evmtypes.Nonce(666))
		actTx, actErr := inMemoryStore.FindTxWithSequence(ctx, fromAddress, evmtypes.Nonce(666))
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.NotNil(t, expTx)
		require.NotNil(t, actTx)
		assertTxEqual(t, *expTx, *actTx)
	})

	t.Run("incorrect from address", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 777, 7, fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		wrongFromAddress := common.Address{}
		expTx, expErr := persistentStore.FindTxWithSequence(ctx, wrongFromAddress, evmtypes.Nonce(777))
		actTx, actErr := inMemoryStore.FindTxWithSequence(ctx, wrongFromAddress, evmtypes.Nonce(777))
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.Nil(t, expTx)
		require.Nil(t, actTx)
	})
}

func TestInMemoryStore_CountTransactionsByState(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no results", func(t *testing.T) {
		expCount, expErr := persistentStore.CountTransactionsByState(ctx, commontxmgr.TxUnconfirmed, chainID)
		actCount, actErr := inMemoryStore.CountTransactionsByState(ctx, commontxmgr.TxUnconfirmed, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, expCount, actCount)
	})
	t.Run("3 unconfirmed transactions", func(t *testing.T) {
		for i := int64(0); i < 3; i++ {
			// insert the transaction into the persistent store
			inTx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, i, fromAddress)
			// insert the transaction into the in-memory store
			require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))
		}

		expCount, expErr := persistentStore.CountTransactionsByState(ctx, commontxmgr.TxUnconfirmed, chainID)
		actCount, actErr := inMemoryStore.CountTransactionsByState(ctx, commontxmgr.TxUnconfirmed, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, expCount, actCount)
	})
}

func TestInMemoryStore_FindTxsRequiringResubmissionDueToInsufficientEth(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	_, dbcfg, evmcfg := evmtxmgr.MakeTestConfigs(t)
	persistentStore := cltest.NewTestTxStore(t, db)
	kst := cltest.NewKeyStore(t, db, dbcfg)
	_, fromAddress := cltest.MustInsertRandomKey(t, kst.Eth())
	_, otherAddress := cltest.MustInsertRandomKey(t, kst.Eth())

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	lggr := logger.TestSugared(t)
	chainID := ethClient.ConfiguredChainID()
	ctx := testutils.Context(t)

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no results", func(t *testing.T) {
		expTxs, expErr := persistentStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, fromAddress, chainID)
		actTxs, actErr := inMemoryStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, fromAddress, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, len(expTxs), len(actTxs))
	})

	// Insert order is mixed up to test sorting
	// insert the transaction into the persistent store
	inTx_2 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, persistentStore, 1, fromAddress)
	inTx_3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, 2, fromAddress)
	attempt3_2 := cltest.NewLegacyEthTxAttempt(t, inTx_3.ID)
	attempt3_2.State = txmgrtypes.TxAttemptInsufficientFunds
	attempt3_2.TxFee.Legacy = evmassets.NewWeiI(100)
	require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt3_2))
	inTx_1 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, persistentStore, 0, fromAddress)
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_2))
	inTx_3.TxAttempts = append([]evmtxmgr.TxAttempt{attempt3_2}, inTx_3.TxAttempts...)
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_3))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))

	// These should never be returned
	// insert the transaction into the persistent store
	otx_1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, 3, fromAddress)
	otx_2 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 4, 100, fromAddress)
	otx_3 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, persistentStore, 0, otherAddress)
	// insert the transaction into the in-memory store
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &otx_1))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &otx_2))
	require.NoError(t, inMemoryStore.XXXTestInsertTx(otherAddress, &otx_3))

	t.Run("return all eth_txes with at least one attempt that is in insufficient_eth state", func(t *testing.T) {
		expTxs, expErr := persistentStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, fromAddress, chainID)
		actTxs, actErr := inMemoryStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, fromAddress, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)

		assert.Equal(t, len(expTxs), len(actTxs))
		for i := 0; i < len(expTxs); i++ {
			assertTxEqual(t, *expTxs[i], *actTxs[i])
		}
	})

	t.Run("does not return txes with different fromAddress", func(t *testing.T) {
		anotherFromAddress := common.Address{}
		expTxs, expErr := persistentStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, anotherFromAddress, chainID)
		actTxs, actErr := inMemoryStore.FindTxsRequiringResubmissionDueToInsufficientFunds(ctx, anotherFromAddress, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, len(expTxs), len(actTxs))
	})
}

func TestInMemoryStore_GetNonFatalTransactions(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no results", func(t *testing.T) {
		expTxs, expErr := persistentStore.GetNonFatalTransactions(ctx, chainID)
		actTxs, actErr := inMemoryStore.GetNonFatalTransactions(ctx, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, len(expTxs), len(actTxs))
	})

	t.Run("get in progress, unstarted, and unconfirmed transactions", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx_0 := mustInsertInProgressEthTxWithAttempt(t, persistentStore, 123, fromAddress)
		inTx_1 := mustCreateUnstartedGeneratedTx(t, persistentStore, fromAddress, chainID)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))

		expTxs, expErr := persistentStore.GetNonFatalTransactions(ctx, chainID)
		actTxs, actErr := inMemoryStore.GetNonFatalTransactions(ctx, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.Equal(t, len(expTxs), len(actTxs))

		for i := 0; i < len(expTxs); i++ {
			assertTxEqual(t, *expTxs[i], *actTxs[i])
		}
	})
}

func TestInMemoryStore_FindTransactionsConfirmedInBlockRange(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no results", func(t *testing.T) {
		expTxs, expErr := persistentStore.FindTransactionsConfirmedInBlockRange(ctx, 10, 8, chainID)
		actTxs, actErr := inMemoryStore.FindTransactionsConfirmedInBlockRange(ctx, 10, 8, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, len(expTxs), len(actTxs))
	})

	t.Run("find all transactions confirmed in range", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx_0 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 700, 8, fromAddress)
		rec_0 := mustInsertEthReceipt(t, persistentStore, 8, evmutils.NewHash(), inTx_0.TxAttempts[0].Hash)
		inTx_1 := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 777, 9, fromAddress)
		rec_1 := mustInsertEthReceipt(t, persistentStore, 9, evmutils.NewHash(), inTx_1.TxAttempts[0].Hash)
		// insert the transaction into the in-memory store
		inTx_0.TxAttempts[0].Receipts = append(inTx_0.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&rec_0))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))
		inTx_1.TxAttempts[0].Receipts = append(inTx_1.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&rec_1))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))

		expTxs, expErr := persistentStore.FindTransactionsConfirmedInBlockRange(ctx, 10, 8, chainID)
		actTxs, actErr := inMemoryStore.FindTransactionsConfirmedInBlockRange(ctx, 10, 8, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.Equal(t, len(expTxs), len(actTxs))
		for i := 0; i < len(expTxs); i++ {
			assertTxEqual(t, *expTxs[i], *actTxs[i])
		}
	})
}

func TestInMemoryStore_FindEarliestUnconfirmedBroadcastTime(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no results", func(t *testing.T) {
		expBroadcastAt, expErr := persistentStore.FindEarliestUnconfirmedBroadcastTime(ctx, chainID)
		actBroadcastAt, actErr := inMemoryStore.FindEarliestUnconfirmedBroadcastTime(ctx, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.False(t, expBroadcastAt.Valid)
		assert.False(t, actBroadcastAt.Valid)
	})
	t.Run("find broadcast at time", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertUnconfirmedEthTx(t, persistentStore, 123, fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expBroadcastAt, expErr := persistentStore.FindEarliestUnconfirmedBroadcastTime(ctx, chainID)
		actBroadcastAt, actErr := inMemoryStore.FindEarliestUnconfirmedBroadcastTime(ctx, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.True(t, expBroadcastAt.Valid)
		require.True(t, actBroadcastAt.Valid)
		assert.Equal(t, expBroadcastAt.Time.Unix(), actBroadcastAt.Time.Unix())
	})
}

func TestInMemoryStore_FindEarliestUnconfirmedTxAttemptBlock(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("no results", func(t *testing.T) {
		expBlock, expErr := persistentStore.FindEarliestUnconfirmedTxAttemptBlock(ctx, chainID)
		actBlock, actErr := inMemoryStore.FindEarliestUnconfirmedTxAttemptBlock(ctx, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.False(t, expBlock.Valid)
		assert.False(t, actBlock.Valid)
	})

	t.Run("find earliest unconfirmed tx block", func(t *testing.T) {
		broadcastBeforeBlockNum := int64(2)
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertUnconfirmedEthTx(t, persistentStore, 123, fromAddress)
		attempt := cltest.NewLegacyEthTxAttempt(t, inTx.ID)
		attempt.BroadcastBeforeBlockNum = &broadcastBeforeBlockNum
		attempt.State = txmgrtypes.TxAttemptBroadcast
		require.NoError(t, persistentStore.InsertTxAttempt(ctx, &attempt))
		inTx.TxAttempts = append(inTx.TxAttempts, attempt)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expBlock, expErr := persistentStore.FindEarliestUnconfirmedTxAttemptBlock(ctx, chainID)
		actBlock, actErr := inMemoryStore.FindEarliestUnconfirmedTxAttemptBlock(ctx, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.True(t, expBlock.Valid)
		assert.True(t, actBlock.Valid)
		assert.Equal(t, expBlock.Int64, actBlock.Int64)
	})
}

func TestInMemoryStore_LoadTxAttempts(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("load tx attempt", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(t, persistentStore, 1, 7, time.Now(), fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expTx := evmtxmgr.Tx{ID: inTx.ID, TxAttempts: []evmtxmgr.TxAttempt{}, FromAddress: fromAddress} // empty tx attempts for test
		expErr := persistentStore.LoadTxAttempts(ctx, &expTx)
		require.Equal(t, 1, len(expTx.TxAttempts))
		expAttempt := expTx.TxAttempts[0]

		actTx := evmtxmgr.Tx{ID: inTx.ID, TxAttempts: []evmtxmgr.TxAttempt{}, FromAddress: fromAddress} // empty tx attempts for test
		actErr := inMemoryStore.LoadTxAttempts(ctx, &actTx)
		require.Equal(t, 1, len(actTx.TxAttempts))
		actAttempt := actTx.TxAttempts[0]

		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assertTxAttemptEqual(t, expAttempt, actAttempt)
	})
}

func TestInMemoryStore_PreloadTxes(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("load transaction", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, int64(7), fromAddress)
		// insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		expAttempts := []evmtxmgr.TxAttempt{{ID: 0, TxID: inTx.ID}}
		expErr := persistentStore.PreloadTxes(ctx, expAttempts)
		require.Equal(t, 1, len(expAttempts))
		expAttempt := expAttempts[0]

		actAttempts := []evmtxmgr.TxAttempt{{ID: 0, TxID: inTx.ID}}
		actErr := inMemoryStore.PreloadTxes(ctx, actAttempts)
		require.Equal(t, 1, len(actAttempts))
		actAttempt := actAttempts[0]

		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assertTxAttemptEqual(t, expAttempt, actAttempt)
	})
}

func TestInMemoryStore_IsTxFinalized(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("tx not past finality depth", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 111, 1, fromAddress)
		rec := mustInsertEthReceipt(t, persistentStore, 1, evmutils.NewHash(), inTx.TxAttempts[0].Hash)
		// insert the transaction into the in-memory store
		inTx.TxAttempts[0].Receipts = append(inTx.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&rec))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		blockHeight := int64(2)
		expIsFinalized, expErr := persistentStore.IsTxFinalized(ctx, blockHeight, inTx.ID, chainID)
		actIsFinalized, actErr := inMemoryStore.IsTxFinalized(ctx, blockHeight, inTx.ID, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, expIsFinalized, actIsFinalized)
	})

	t.Run("tx is past finality depth", func(t *testing.T) {
		// insert the transaction into the persistent store
		inTx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, persistentStore, 122, 2, fromAddress)
		rec := mustInsertEthReceipt(t, persistentStore, 2, evmutils.NewHash(), inTx.TxAttempts[0].Hash)
		// insert the transaction into the in-memory store
		inTx.TxAttempts[0].Receipts = append(inTx.TxAttempts[0].Receipts, evmtxmgr.DbReceiptToEvmReceipt(&rec))
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		blockHeight := int64(10)
		expIsFinalized, expErr := persistentStore.IsTxFinalized(ctx, blockHeight, inTx.ID, chainID)
		actIsFinalized, actErr := inMemoryStore.IsTxFinalized(ctx, blockHeight, inTx.ID, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		assert.Equal(t, expIsFinalized, actIsFinalized)
	})
}

func TestInMemoryStore_FindTxsRequiringGasBump(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("gets transactions requiring gas bumping", func(t *testing.T) {
		currentBlockNum := int64(10)

		// insert the transaction into the persistent store
		inTx_0 := mustInsertUnconfirmedEthTxWithAttemptState(t, persistentStore, 1, fromAddress, txmgrtypes.TxAttemptBroadcast)
		require.NoError(t, persistentStore.SetBroadcastBeforeBlockNum(ctx, currentBlockNum, chainID))
		inTx_1 := mustInsertUnconfirmedEthTxWithAttemptState(t, persistentStore, 2, fromAddress, txmgrtypes.TxAttemptBroadcast)
		require.NoError(t, persistentStore.SetBroadcastBeforeBlockNum(ctx, currentBlockNum+1, chainID))
		// insert the transaction into the in-memory store
		inTx_0.TxAttempts[0].BroadcastBeforeBlockNum = &currentBlockNum
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_0))
		tempCurrentBlockNum := currentBlockNum + 1
		inTx_1.TxAttempts[0].BroadcastBeforeBlockNum = &tempCurrentBlockNum
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx_1))

		newBlock := int64(12)
		gasBumpThreshold := int64(2)
		expTxs, expErr := persistentStore.FindTxsRequiringGasBump(ctx, fromAddress, newBlock, gasBumpThreshold, 0, chainID)
		actTxs, actErr := inMemoryStore.FindTxsRequiringGasBump(ctx, fromAddress, newBlock, gasBumpThreshold, 0, chainID)
		require.NoError(t, expErr)
		require.NoError(t, actErr)
		require.Equal(t, len(expTxs), len(actTxs))
		for i := 0; i < len(expTxs); i++ {
			assertTxEqual(t, *expTxs[i], *actTxs[i])
		}
	})
}

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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
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

func TestInMemoryStore_UpdateBroadcastAts(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("does not update when broadcast_at is Null", func(t *testing.T) {
		// Insert a transaction into persistent store
		inTx := mustInsertInProgressEthTxWithAttempt(t, persistentStore, 1, fromAddress)
		require.Nil(t, inTx.BroadcastAt)
		now := time.Now()
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		err := inMemoryStore.UpdateBroadcastAts(
			ctx,
			now,
			[]int64{inTx.ID},
		)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)
		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]
		assertTxEqual(t, expTx, actTx)
		assert.Nil(t, actTx.BroadcastAt)
	})

	t.Run("updates broadcast_at when not null", func(t *testing.T) {
		// Insert a transaction into persistent store
		time1 := time.Now()
		inTx := cltest.NewEthTx(fromAddress)
		inTx.Sequence = new(evmtypes.Nonce)
		inTx.State = commontxmgr.TxUnconfirmed
		inTx.BroadcastAt = &time1
		inTx.InitialBroadcastAt = &time1
		require.NoError(t, persistentStore.InsertTx(ctx, &inTx))
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		time2 := time1.Add(1 * time.Hour)
		err := inMemoryStore.UpdateBroadcastAts(
			ctx,
			time2,
			[]int64{inTx.ID},
		)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)
		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]
		assertTxEqual(t, expTx, actTx)
		assert.NotNil(t, actTx.BroadcastAt)
	})
}

func TestInMemoryStore_SetBroadcastBeforeBlockNum(t *testing.T) {
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
	require.NoError(t, err)

	t.Run("saves block num to unconfirmed evm.tx_attempts without one", func(t *testing.T) {
		// Insert a transaction into persistent store
		inTx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, 1, fromAddress)
		// Insert the transaction into the in-memory store
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		headNum := int64(9000)
		err := inMemoryStore.SetBroadcastBeforeBlockNum(ctx, headNum, chainID)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)
		require.Equal(t, 1, len(expTx.TxAttempts))
		assert.Equal(t, headNum, *expTx.TxAttempts[0].BroadcastBeforeBlockNum)
		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]
		assertTxEqual(t, expTx, actTx)
	})

	t.Run("does not change evm.tx_attempts that already have BroadcastBeforeBlockNum set", func(t *testing.T) {
		n := int64(42)
		// Insert a transaction into persistent store
		inTx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, persistentStore, 11, fromAddress)
		inTxAttempt := newBroadcastLegacyEthTxAttempt(t, inTx.ID, 2)
		inTxAttempt.BroadcastBeforeBlockNum = &n
		require.NoError(t, persistentStore.InsertTxAttempt(ctx, &inTxAttempt))
		// Insert the transaction into the in-memory store
		inTx.TxAttempts = append([]evmtxmgr.TxAttempt{inTxAttempt}, inTx.TxAttempts...)
		require.NoError(t, inMemoryStore.XXXTestInsertTx(fromAddress, &inTx))

		headNum := int64(9000)
		err := inMemoryStore.SetBroadcastBeforeBlockNum(ctx, headNum, chainID)
		require.NoError(t, err)

		expTx, err := persistentStore.FindTxWithAttempts(ctx, inTx.ID)
		require.NoError(t, err)
		require.Equal(t, 2, len(expTx.TxAttempts))
		assert.Equal(t, n, *expTx.TxAttempts[0].BroadcastBeforeBlockNum)
		fn := func(tx *evmtxmgr.Tx) bool { return true }
		actTxs := inMemoryStore.XXXTestFindTxs(nil, fn, inTx.ID)
		require.Equal(t, 1, len(actTxs))
		actTx := actTxs[0]
		assertTxEqual(t, expTx, actTx)
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
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

	inMemoryStore, err := commontxmgr.NewInMemoryStore(ctx, lggr, chainID, kst.Eth(), persistentStore, evmcfg.Transactions())
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
	if exp.BroadcastAt != nil {
		require.NotNil(t, act.BroadcastAt)
		assert.Equal(t, exp.BroadcastAt.Unix(), act.BroadcastAt.Unix())
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

	if len(exp.TxAttempts) == 0 {
		return
	}
	require.Equal(t, len(exp.TxAttempts), len(act.TxAttempts))
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

	if len(exp.Receipts) == 0 {
		return
	}
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
