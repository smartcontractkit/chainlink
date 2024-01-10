package txmgr_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

func newBroadcastLegacyEthTxAttempt(t *testing.T, etxID int64, gasPrice ...int64) txmgr.TxAttempt {
	attempt := cltest.NewLegacyEthTxAttempt(t, etxID)
	attempt.State = txmgrtypes.TxAttemptBroadcast
	if len(gasPrice) > 0 {
		gp := gasPrice[0]
		attempt.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(gp)}
	}
	return attempt
}

func mustTxBeInState(t *testing.T, txStore txmgr.TestEvmTxStore, tx txmgr.Tx, expectedState txmgrtypes.TxState) {
	etx, err := txStore.FindTxWithAttempts(tx.ID)
	require.NoError(t, err)
	require.Equal(t, expectedState, etx.State)
}

func newTxReceipt(hash gethCommon.Hash, blockNumber int, txIndex uint) evmtypes.Receipt {
	return evmtypes.Receipt{
		TxHash:           hash,
		BlockHash:        utils.NewHash(),
		BlockNumber:      big.NewInt(int64(blockNumber)),
		TransactionIndex: txIndex,
		Status:           uint64(1),
	}
}

func newInProgressLegacyEthTxAttempt(t *testing.T, etxID int64, gasPrice ...int64) txmgr.TxAttempt {
	attempt := cltest.NewLegacyEthTxAttempt(t, etxID)
	attempt.State = txmgrtypes.TxAttemptInProgress
	if len(gasPrice) > 0 {
		gp := gasPrice[0]
		attempt.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(gp)}
	}
	return attempt
}

func mustInsertInProgressEthTx(t *testing.T, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress gethCommon.Address) txmgr.Tx {
	etx := cltest.NewEthTx(fromAddress)
	etx.State = txmgrcommon.TxInProgress
	n := evmtypes.Nonce(nonce)
	etx.Sequence = &n
	require.NoError(t, txStore.InsertTx(&etx))

	return etx
}

func mustInsertConfirmedEthTx(t *testing.T, txStore txmgr.TestEvmTxStore, nonce int64, fromAddress gethCommon.Address) txmgr.Tx {
	etx := cltest.NewEthTx(fromAddress)
	etx.State = txmgrcommon.TxConfirmed
	n := evmtypes.Nonce(nonce)
	etx.Sequence = &n
	now := time.Now()
	etx.BroadcastAt = &now
	etx.InitialBroadcastAt = &now
	require.NoError(t, txStore.InsertTx(&etx))

	return etx
}

//	func TestEthConfirmer_Lifecycle(t *testing.T) {
//		t.Parallel()
//
//		db := pgtest.NewSqlxDB(t)
//		config := newTestChainScopedConfig(t)
//		txStore := newTxStore(t, db, config.Database())
//
//		ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
//		ethKeyStore := cltest.NewKeyStore(t, db, config.Database()).Eth()
//
//		// Add some fromAddresses
//		cltest.MustInsertRandomKey(t, ethKeyStore)
//		cltest.MustInsertRandomKey(t, ethKeyStore)
//		estimator := gasmocks.NewEvmEstimator(t)
//		newEst := func(logger.Logger) gas.EvmEstimator { return estimator }
//		lggr := logger.Test(t)
//		ge := config.EVM().GasEstimator()
//		feeEstimator := gas.NewWrappedEvmEstimator(lggr, newEst, ge.EIP1559DynamicFees(), nil)
//		txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, ethKeyStore, feeEstimator)
//		ec := txmgr.NewEvmConfirmer(txStore, txmgr.NewEvmTxmClient(ethClient), txmgr.NewEvmTxmConfig(config.EVM()), txmgr.NewEvmTxmFeeConfig(ge), config.EVM().Transactions(), config.Database(), ethKeyStore, txBuilder, lggr)
//		ctx := testutils.Context(t)
//
//		// Can't close unstarted instance
//		err := ec.Close()
//		require.Error(t, err)
//
//		// Can successfully start once
//		err = ec.Start(ctx)
//		require.NoError(t, err)
//
//		// Can't start an already started instance
//		err = ec.Start(ctx)
//		require.Error(t, err)
//		head := evmtypes.Head{
//			Hash:   utils.NewHash(),
//			Number: 10,
//			Parent: &evmtypes.Head{
//				Hash:   utils.NewHash(),
//				Number: 9,
//				Parent: &evmtypes.Head{
//					Number: 8,
//					Hash:   utils.NewHash(),
//					Parent: nil,
//				},
//			},
//		}
//		err = ec.ProcessHead(ctx, &head)
//		require.NoError(t, err)
//		// Can successfully close once
//		err = ec.Close()
//		require.NoError(t, err)
//
//		// Can't start more than once (Confirmer uses services.StateMachine)
//		err = ec.Start(ctx)
//		require.Error(t, err)
//		// Can't close more than once (Confirmer use services.StateMachine)
//		err = ec.Close()
//		require.Error(t, err)
//
//		// Can't closeInternal unstarted instance
//		require.Error(t, ec.XXXTestCloseInternal())
//
//		// Can successfully startInternal a previously closed instance
//		require.NoError(t, ec.XXXTestStartInternal())
//		// Can't startInternal already started instance
//		require.Error(t, ec.XXXTestStartInternal())
//		// Can successfully closeInternal again
//		require.NoError(t, ec.XXXTestCloseInternal())
//	}
//func TestEthConfirmer_CheckForReceipts(t *testing.T) {
//	t.Parallel()
//
//	db := pgtest.NewSqlxDB(t)
//	config := newTestChainScopedConfig(t)
//	txStore := cltest.NewTestTxStore(t, db, config.Database())
//
//	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
//	ethKeyStore := cltest.NewKeyStore(t, db, config.Database()).Eth()
//
//	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
//
//	ec := newEthConfirmer(t, txStore, ethClient, config, ethKeyStore, nil)
//
//	nonce := int64(0)
//	ctx := testutils.Context(t)
//	blockNum := int64(0)
//
//	t.Run("only finds eth_txes in unconfirmed state with at least one broadcast attempt", func(t *testing.T) {
//		mustInsertFatalErrorEthTx(t, txStore, fromAddress)
//		mustInsertInProgressEthTx(t, txStore, nonce, fromAddress)
//		nonce++
//		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, nonce, 1, fromAddress)
//		nonce++
//		mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, nonce, fromAddress)
//		nonce++
//		mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, config.EVM().ChainID())
//
//		// Do the thing
//		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))
//	})
//
//	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
//	nonce++
//	require.Len(t, etx1.TxAttempts, 1)
//	attempt1_1 := etx1.TxAttempts[0]
//	hashAttempt1_1 := attempt1_1.Hash
//	require.Len(t, attempt1_1.Receipts, 0)
//
//	t.Run("fetches receipt for one unconfirmed eth_tx", func(t *testing.T) {
//		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
//		// Transaction not confirmed yet, receipt is nil
//		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], hashAttempt1_1, "eth_getTransactionReceipt")
//		})).Return(nil).Run(func(args mock.Arguments) {
//			elems := args.Get(1).([]rpc.BatchElem)
//			elems[0].Result = &evmtypes.Receipt{}
//		}).Once()
//
//		// Do the thing
//		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))
//
//		var err error
//		etx1, err = txStore.FindTxWithAttempts(etx1.ID)
//		assert.NoError(t, err)
//		require.Len(t, etx1.TxAttempts, 1)
//		attempt1_1 = etx1.TxAttempts[0]
//		require.NoError(t, err)
//		require.Len(t, attempt1_1.Receipts, 0)
//	})
//
//	t.Run("saves nothing if returned receipt does not match the attempt", func(t *testing.T) {
//		txmReceipt := evmtypes.Receipt{
//			TxHash:           utils.NewHash(),
//			BlockHash:        utils.NewHash(),
//			BlockNumber:      big.NewInt(42),
//			TransactionIndex: uint(1),
//		}
//
//		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
//		// First transaction confirmed
//		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], hashAttempt1_1, "eth_getTransactionReceipt")
//		})).Return(nil).Run(func(args mock.Arguments) {
//			elems := args.Get(1).([]rpc.BatchElem)
//			*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt
//		}).Once()
//
//		// No error because it is merely logged
//		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))
//
//		etx, err := txStore.FindTxWithAttempts(etx1.ID)
//		require.NoError(t, err)
//		require.Len(t, etx.TxAttempts, 1)
//
//		require.Len(t, etx.TxAttempts[0].Receipts, 0)
//	})
//
//	t.Run("saves nothing if query returns error", func(t *testing.T) {
//		txmReceipt := evmtypes.Receipt{
//			TxHash:           attempt1_1.Hash,
//			BlockHash:        utils.NewHash(),
//			BlockNumber:      big.NewInt(42),
//			TransactionIndex: uint(1),
//		}
//
//		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
//		// First transaction confirmed
//		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], hashAttempt1_1, "eth_getTransactionReceipt")
//		})).Return(nil).Run(func(args mock.Arguments) {
//			elems := args.Get(1).([]rpc.BatchElem)
//			*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt
//			elems[0].Error = errors.New("foo")
//		}).Once()
//
//		// No error because it is merely logged
//		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))
//
//		etx, err := txStore.FindTxWithAttempts(etx1.ID)
//		require.NoError(t, err)
//		require.Len(t, etx.TxAttempts, 1)
//		require.Len(t, etx.TxAttempts[0].Receipts, 0)
//	})
//
//	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
//	nonce++
//	require.Len(t, etx2.TxAttempts, 1)
//	attempt2_1 := etx2.TxAttempts[0]
//	require.Len(t, attempt2_1.Receipts, 0)
//
//	t.Run("saves eth_receipt and marks eth_tx as confirmed when geth client returns valid receipt", func(t *testing.T) {
//		txmReceipt := evmtypes.Receipt{
//			TxHash:           attempt1_1.Hash,
//			BlockHash:        utils.NewHash(),
//			BlockNumber:      big.NewInt(42),
//			TransactionIndex: uint(1),
//			Status:           uint64(1),
//		}
//
//		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
//		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//			return len(b) == 2 &&
//				cltest.BatchElemMatchesParams(b[0], attempt1_1.Hash, "eth_getTransactionReceipt") &&
//				cltest.BatchElemMatchesParams(b[1], attempt2_1.Hash, "eth_getTransactionReceipt")
//
//		})).Return(nil).Run(func(args mock.Arguments) {
//			elems := args.Get(1).([]rpc.BatchElem)
//			// First transaction confirmed
//			*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt
//			// Second transaction still unconfirmed
//			elems[1].Result = &evmtypes.Receipt{}
//		}).Once()
//
//		// Do the thing
//		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))
//
//		// Check that the receipt was saved
//		etx, err := txStore.FindTxWithAttempts(etx1.ID)
//		require.NoError(t, err)
//
//		assert.Equal(t, txmgrcommon.TxConfirmed, etx.State)
//		assert.Len(t, etx.TxAttempts, 1)
//		attempt1_1 = etx.TxAttempts[0]
//		require.Len(t, attempt1_1.Receipts, 1)
//
//		ethReceipt := attempt1_1.Receipts[0]
//
//		assert.Equal(t, txmReceipt.TxHash, ethReceipt.GetTxHash())
//		assert.Equal(t, txmReceipt.BlockHash, ethReceipt.GetBlockHash())
//		assert.Equal(t, txmReceipt.BlockNumber.Int64(), ethReceipt.GetBlockNumber().Int64())
//		assert.Equal(t, txmReceipt.TransactionIndex, ethReceipt.GetTransactionIndex())
//
//		receiptJSON, err := json.Marshal(txmReceipt)
//		require.NoError(t, err)
//
//		j, err := json.Marshal(ethReceipt)
//		require.NoError(t, err)
//		assert.JSONEq(t, string(receiptJSON), string(j))
//	})
//
//	t.Run("fetches and saves receipts for several attempts in gas price order", func(t *testing.T) {
//		attempt2_2 := newBroadcastLegacyEthTxAttempt(t, etx2.ID)
//		attempt2_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(10)}
//
//		attempt2_3 := newBroadcastLegacyEthTxAttempt(t, etx2.ID)
//		attempt2_3.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(20)}
//
//		// Insert order deliberately reversed to test sorting by gas price
//		require.NoError(t, txStore.InsertTxAttempt(&attempt2_3))
//		require.NoError(t, txStore.InsertTxAttempt(&attempt2_2))
//
//		txmReceipt := evmtypes.Receipt{
//			TxHash:           attempt2_2.Hash,
//			BlockHash:        utils.NewHash(),
//			BlockNumber:      big.NewInt(42),
//			TransactionIndex: uint(1),
//			Status:           uint64(1),
//		}
//
//		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
//		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//			return len(b) == 3 &&
//				cltest.BatchElemMatchesParams(b[2], attempt2_1.Hash, "eth_getTransactionReceipt") &&
//				cltest.BatchElemMatchesParams(b[1], attempt2_2.Hash, "eth_getTransactionReceipt") &&
//				cltest.BatchElemMatchesParams(b[0], attempt2_3.Hash, "eth_getTransactionReceipt")
//
//		})).Return(nil).Run(func(args mock.Arguments) {
//			elems := args.Get(1).([]rpc.BatchElem)
//			// Most expensive attempt still unconfirmed
//			elems[2].Result = &evmtypes.Receipt{}
//			// Second most expensive attempt is confirmed
//			*(elems[1].Result.(*evmtypes.Receipt)) = txmReceipt
//			// Cheapest attempt still unconfirmed
//			elems[0].Result = &evmtypes.Receipt{}
//		}).Once()
//
//		// Do the thing
//		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))
//
//		// Check that the state was updated
//		etx, err := txStore.FindTxWithAttempts(etx2.ID)
//		require.NoError(t, err)
//
//		require.Equal(t, txmgrcommon.TxConfirmed, etx.State)
//		require.Len(t, etx.TxAttempts, 3)
//	})
//
//	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
//	attempt3_1 := etx3.TxAttempts[0]
//	nonce++
//
//	t.Run("ignores receipt missing BlockHash that comes from querying parity too early", func(t *testing.T) {
//		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
//		receipt := evmtypes.Receipt{
//			TxHash: attempt3_1.Hash,
//			Status: uint64(1),
//		}
//		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt3_1.Hash, "eth_getTransactionReceipt")
//		})).Return(nil).Run(func(args mock.Arguments) {
//			elems := args.Get(1).([]rpc.BatchElem)
//			*(elems[0].Result.(*evmtypes.Receipt)) = receipt
//		}).Once()
//
//		// Do the thing
//		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))
//
//		// No receipt, but no error either
//		etx, err := txStore.FindTxWithAttempts(etx3.ID)
//		require.NoError(t, err)
//
//		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
//		assert.Len(t, etx.TxAttempts, 1)
//		attempt3_1 = etx.TxAttempts[0]
//		require.Len(t, attempt3_1.Receipts, 0)
//	})
//
//	t.Run("does not panic if receipt has BlockHash but is missing some other fields somehow", func(t *testing.T) {
//		// NOTE: This should never happen, but we shouldn't panic regardless
//		receipt := evmtypes.Receipt{
//			TxHash:    attempt3_1.Hash,
//			BlockHash: utils.NewHash(),
//			Status:    uint64(1),
//		}
//		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt3_1.Hash, "eth_getTransactionReceipt")
//		})).Return(nil).Run(func(args mock.Arguments) {
//			elems := args.Get(1).([]rpc.BatchElem)
//			*(elems[0].Result.(*evmtypes.Receipt)) = receipt
//		}).Once()
//
//		// Do the thing
//		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))
//
//		// No receipt, but no error either
//		etx, err := txStore.FindTxWithAttempts(etx3.ID)
//		require.NoError(t, err)
//
//		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
//		assert.Len(t, etx.TxAttempts, 1)
//		attempt3_1 = etx.TxAttempts[0]
//		require.Len(t, attempt3_1.Receipts, 0)
//	})
//	t.Run("handles case where eth_receipt already exists somehow", func(t *testing.T) {
//		ethReceipt := mustInsertEthReceipt(t, txStore, 42, utils.NewHash(), attempt3_1.Hash)
//		txmReceipt := evmtypes.Receipt{
//			TxHash:           attempt3_1.Hash,
//			BlockHash:        ethReceipt.BlockHash,
//			BlockNumber:      big.NewInt(ethReceipt.BlockNumber),
//			TransactionIndex: ethReceipt.TransactionIndex,
//			Status:           uint64(1),
//		}
//		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
//		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt3_1.Hash, "eth_getTransactionReceipt")
//		})).Return(nil).Run(func(args mock.Arguments) {
//			elems := args.Get(1).([]rpc.BatchElem)
//			*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt
//		}).Once()
//
//		// Do the thing
//		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))
//
//		// Check that the receipt was unchanged
//		etx, err := txStore.FindTxWithAttempts(etx3.ID)
//		require.NoError(t, err)
//
//		assert.Equal(t, txmgrcommon.TxConfirmed, etx.State)
//		assert.Len(t, etx.TxAttempts, 1)
//		attempt3_1 = etx.TxAttempts[0]
//		require.Len(t, attempt3_1.Receipts, 1)
//
//		ethReceipt3_1 := attempt3_1.Receipts[0]
//
//		assert.Equal(t, txmReceipt.TxHash, ethReceipt3_1.GetTxHash())
//		assert.Equal(t, txmReceipt.BlockHash, ethReceipt3_1.GetBlockHash())
//		assert.Equal(t, txmReceipt.BlockNumber.Int64(), ethReceipt3_1.GetBlockNumber().Int64())
//		assert.Equal(t, txmReceipt.TransactionIndex, ethReceipt3_1.GetTransactionIndex())
//	})
//
//	etx4 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
//	attempt4_1 := etx4.TxAttempts[0]
//	nonce++
//
//	t.Run("on receipt fetch marks in_progress eth_tx_attempt as broadcast", func(t *testing.T) {
//		attempt4_2 := newInProgressLegacyEthTxAttempt(t, etx4.ID)
//		attempt4_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(10)}
//
//		require.NoError(t, txStore.InsertTxAttempt(&attempt4_2))
//
//		txmReceipt := evmtypes.Receipt{
//			TxHash:           attempt4_2.Hash,
//			BlockHash:        utils.NewHash(),
//			BlockNumber:      big.NewInt(42),
//			TransactionIndex: uint(1),
//			Status:           uint64(1),
//		}
//		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
//		// Second attempt is confirmed
//		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//			return len(b) == 2 &&
//				cltest.BatchElemMatchesParams(b[0], attempt4_2.Hash, "eth_getTransactionReceipt") &&
//				cltest.BatchElemMatchesParams(b[1], attempt4_1.Hash, "eth_getTransactionReceipt")
//		})).Return(nil).Run(func(args mock.Arguments) {
//			elems := args.Get(1).([]rpc.BatchElem)
//			// First attempt still unconfirmed
//			elems[1].Result = &evmtypes.Receipt{}
//			// Second attempt is confirmed
//			*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt
//		}).Once()
//
//		// Do the thing
//		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))
//
//		// Check that the state was updated
//		var err error
//		etx4, err = txStore.FindTxWithAttempts(etx4.ID)
//		require.NoError(t, err)
//
//		attempt4_1 = etx4.TxAttempts[1]
//		attempt4_2 = etx4.TxAttempts[0]
//
//		// And the attempts
//		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt4_1.State)
//		require.Nil(t, attempt4_1.BroadcastBeforeBlockNum)
//		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt4_2.State)
//		require.Equal(t, int64(42), *attempt4_2.BroadcastBeforeBlockNum)
//
//		// Check receipts
//		require.Len(t, attempt4_1.Receipts, 0)
//		require.Len(t, attempt4_2.Receipts, 1)
//	})
//
//	etx5 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
//	attempt5_1 := etx5.TxAttempts[0]
//	nonce++
//
//	t.Run("simulate on revert", func(t *testing.T) {
//		txmReceipt := evmtypes.Receipt{
//			TxHash:           attempt5_1.Hash,
//			BlockHash:        utils.NewHash(),
//			BlockNumber:      big.NewInt(42),
//			TransactionIndex: uint(1),
//			Status:           uint64(0),
//		}
//		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
//		// First attempt is confirmed and reverted
//		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//			return len(b) == 1 &&
//				cltest.BatchElemMatchesParams(b[0], attempt5_1.Hash, "eth_getTransactionReceipt")
//		})).Return(nil).Run(func(args mock.Arguments) {
//			elems := args.Get(1).([]rpc.BatchElem)
//			// First attempt still unconfirmed
//			*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt
//		}).Once()
//		data, err := utils.ABIEncode(`[{"type":"uint256"}]`, big.NewInt(10))
//		require.NoError(t, err)
//		sig := utils.Keccak256Fixed([]byte(`MyError(uint256)`))
//		ethClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(nil, &client.JsonError{
//			Code:    1,
//			Message: "reverted",
//			Data:    utils.ConcatBytes(sig[:4], data),
//		}).Once()
//
//		// Do the thing
//		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))
//
//		// Check that the state was updated
//		etx5, err = txStore.FindTxWithAttempts(etx5.ID)
//		require.NoError(t, err)
//
//		attempt5_1 = etx5.TxAttempts[0]
//
//		// And the attempts
//		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt5_1.State)
//		require.NotNil(t, attempt5_1.BroadcastBeforeBlockNum)
//		// Check receipts
//		require.Len(t, attempt5_1.Receipts, 1)
//	})
//}

//
//func TestEthConfirmer_CheckForReceipts_batching(t *testing.T) {
//	t.Parallel()
//
//	db := pgtest.NewSqlxDB(t)
//	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
//		c.EVM[0].RPCDefaultBatchSize = ptr[uint32](2)
//	})
//	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
//
//	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
//
//	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
//
//	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
//
//	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
//
//	ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)
//	ctx := testutils.Context(t)
//
//	etx := cltest.MustInsertUnconfirmedEthTx(t, txStore, 0, fromAddress)
//	var attempts []txmgr.TxAttempt
//
//	// Total of 5 attempts should lead to 3 batched fetches (2, 2, 1)
//	for i := 0; i < 5; i++ {
//		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, int64(i+2))
//		require.NoError(t, txStore.InsertTxAttempt(&attempt))
//		attempts = append(attempts, attempt)
//	}
//
//	ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
//
//	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//		return len(b) == 2 &&
//			cltest.BatchElemMatchesParams(b[0], attempts[4].Hash, "eth_getTransactionReceipt") &&
//			cltest.BatchElemMatchesParams(b[1], attempts[3].Hash, "eth_getTransactionReceipt")
//	})).Return(nil).Run(func(args mock.Arguments) {
//		elems := args.Get(1).([]rpc.BatchElem)
//		elems[0].Result = &evmtypes.Receipt{}
//		elems[1].Result = &evmtypes.Receipt{}
//	}).Once()
//	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//		return len(b) == 2 &&
//			cltest.BatchElemMatchesParams(b[0], attempts[2].Hash, "eth_getTransactionReceipt") &&
//			cltest.BatchElemMatchesParams(b[1], attempts[1].Hash, "eth_getTransactionReceipt")
//	})).Return(nil).Run(func(args mock.Arguments) {
//		elems := args.Get(1).([]rpc.BatchElem)
//		elems[0].Result = &evmtypes.Receipt{}
//		elems[1].Result = &evmtypes.Receipt{}
//	}).Once()
//	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//		return len(b) == 1 &&
//			cltest.BatchElemMatchesParams(b[0], attempts[0].Hash, "eth_getTransactionReceipt")
//	})).Return(nil).Run(func(args mock.Arguments) {
//		elems := args.Get(1).([]rpc.BatchElem)
//		elems[0].Result = &evmtypes.Receipt{}
//	}).Once()
//
//	require.NoError(t, ec.CheckForReceipts(ctx, 42))
//}
//
//func TestEthConfirmer_CheckForReceipts_HandlesNonFwdTxsWithForwardingEnabled(t *testing.T) {
//	t.Parallel()
//
//	db := pgtest.NewSqlxDB(t)
//
//	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
//		c.EVM[0].RPCDefaultBatchSize = ptr[uint32](1)
//		c.EVM[0].Transactions.ForwardersEnabled = ptr(true)
//	})
//
//	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
//	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
//	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
//	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
//
//	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
//	ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)
//	ctx := testutils.Context(t)
//	// tx is not forwarded and doesn't have meta set. EthConfirmer should handle nil meta values
//	etx := cltest.MustInsertUnconfirmedEthTx(t, txStore, 0, fromAddress)
//	attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, 2)
//	attempt.Tx.Meta = nil
//	require.NoError(t, txStore.InsertTxAttempt(&attempt))
//	dbtx, err := txStore.FindTxWithAttempts(etx.ID)
//	require.NoError(t, err)
//	require.Equal(t, 0, len(dbtx.TxAttempts[0].Receipts))
//
//	txmReceipt := evmtypes.Receipt{
//		TxHash:           attempt.Hash,
//		BlockHash:        utils.NewHash(),
//		BlockNumber:      big.NewInt(42),
//		TransactionIndex: uint(1),
//		Status:           uint64(1),
//	}
//
//	ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
//	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//		return len(b) == 1 &&
//			cltest.BatchElemMatchesParams(b[0], attempt.Hash, "eth_getTransactionReceipt")
//	})).Return(nil).Run(func(args mock.Arguments) {
//		elems := args.Get(1).([]rpc.BatchElem)
//		*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt // confirmed
//	}).Once()
//
//	require.NoError(t, ec.CheckForReceipts(ctx, 42))
//
//	// Check receipt is inserted correctly.
//	dbtx, err = txStore.FindTxWithAttempts(etx.ID)
//	require.NoError(t, err)
//	require.Equal(t, 1, len(dbtx.TxAttempts[0].Receipts))
//}
//
//func TestEthConfirmer_CheckForReceipts_only_likely_confirmed(t *testing.T) {
//	t.Parallel()
//
//	db := pgtest.NewSqlxDB(t)
//	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
//		c.EVM[0].RPCDefaultBatchSize = ptr[uint32](6)
//	})
//	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
//
//	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
//
//	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
//
//	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
//
//	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
//
//	ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)
//	ctx := testutils.Context(t)
//
//	var attempts []txmgr.TxAttempt
//	// inserting in DESC nonce order to test DB ASC ordering
//	etx2 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 1, fromAddress)
//	for i := 0; i < 4; i++ {
//		attempt := newBroadcastLegacyEthTxAttempt(t, etx2.ID, int64(100-i))
//		require.NoError(t, txStore.InsertTxAttempt(&attempt))
//	}
//	etx := cltest.MustInsertUnconfirmedEthTx(t, txStore, 0, fromAddress)
//	for i := 0; i < 4; i++ {
//		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, int64(100-i))
//		require.NoError(t, txStore.InsertTxAttempt(&attempt))
//
//		// only adding these because a batch for only those attempts should be sent
//		attempts = append(attempts, attempt)
//	}
//
//	ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(0), nil)
//
//	var captured []rpc.BatchElem
//	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//		return len(b) == 4
//	})).Return(nil).Run(func(args mock.Arguments) {
//		elems := args.Get(1).([]rpc.BatchElem)
//		captured = append(captured, elems...)
//		elems[0].Result = &evmtypes.Receipt{}
//		elems[1].Result = &evmtypes.Receipt{}
//		elems[2].Result = &evmtypes.Receipt{}
//		elems[3].Result = &evmtypes.Receipt{}
//	}).Once()
//
//	require.NoError(t, ec.CheckForReceipts(ctx, 42))
//
//	cltest.BatchElemMustMatchParams(t, captured[0], attempts[0].Hash, "eth_getTransactionReceipt")
//	cltest.BatchElemMustMatchParams(t, captured[1], attempts[1].Hash, "eth_getTransactionReceipt")
//	cltest.BatchElemMustMatchParams(t, captured[2], attempts[2].Hash, "eth_getTransactionReceipt")
//	cltest.BatchElemMustMatchParams(t, captured[3], attempts[3].Hash, "eth_getTransactionReceipt")
//}
//
//func TestEthConfirmer_CheckForReceipts_should_not_check_for_likely_unconfirmed(t *testing.T) {
//	t.Parallel()
//
//	db := pgtest.NewSqlxDB(t)
//	config := newTestChainScopedConfig(t)
//	txStore := cltest.NewTestTxStore(t, db, config.Database())
//
//	ethKeyStore := cltest.NewKeyStore(t, db, config.Database()).Eth()
//
//	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
//
//	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
//
//	ec := newEthConfirmer(t, txStore, ethClient, config, ethKeyStore, nil)
//	ctx := testutils.Context(t)
//
//	etx := cltest.MustInsertUnconfirmedEthTx(t, txStore, 1, fromAddress)
//	for i := 0; i < 4; i++ {
//		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, int64(100-i))
//		require.NoError(t, txStore.InsertTxAttempt(&attempt))
//	}
//
//	// latest nonce is lower that all attempts' nonces
//	ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(0), nil)
//
//	require.NoError(t, ec.CheckForReceipts(ctx, 42))
//}
//
//func TestEthConfirmer_CheckForReceipts_confirmed_missing_receipt_scoped_to_key(t *testing.T) {
//	t.Parallel()
//
//	db := pgtest.NewSqlxDB(t)
//	cfg := configtest.NewTestGeneralConfig(t)
//	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
//	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
//
//	_, fromAddress1_1 := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
//	_, fromAddress1_2 := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
//	_, fromAddress2_1 := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
//
//	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
//	ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(20), nil)
//	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
//
//	ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)
//	ctx := testutils.Context(t)
//
//	// STATE
//	// key 1, tx with nonce 0 is unconfirmed
//	// key 1, tx with nonce 1 is unconfirmed
//	// key 2, tx with nonce 9 is unconfirmed and gets a receipt in block 10
//	etx1_0 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 0, fromAddress1_1)
//	etx1_1 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 1, fromAddress1_1)
//	etx2_9 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 3, fromAddress1_2)
//	// there also happens to be a confirmed tx with a higher nonce from a different chain in the DB
//	etx_other_chain := cltest.MustInsertUnconfirmedEthTx(t, txStore, 8, fromAddress2_1)
//	pgtest.MustExec(t, db, `UPDATE evm.txes SET state='confirmed' WHERE id = $1`, etx_other_chain.ID)
//
//	attempt2_9 := newBroadcastLegacyEthTxAttempt(t, etx2_9.ID, int64(1))
//	require.NoError(t, txStore.InsertTxAttempt(&attempt2_9))
//	txmReceipt2_9 := newTxReceipt(attempt2_9.Hash, 10, 1)
//
//	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//		return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt2_9.Hash, "eth_getTransactionReceipt")
//	})).Return(nil).Run(func(args mock.Arguments) {
//		elems := args.Get(1).([]rpc.BatchElem)
//		*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt2_9
//	}).Once()
//
//	require.NoError(t, ec.CheckForReceipts(ctx, 10))
//
//	mustTxBeInState(t, txStore, etx1_0, txmgrcommon.TxUnconfirmed)
//	mustTxBeInState(t, txStore, etx1_1, txmgrcommon.TxUnconfirmed)
//	mustTxBeInState(t, txStore, etx2_9, txmgrcommon.TxConfirmed)
//
//	// Now etx1_1 gets a receipt in block 11, which should mark etx1_0 as confirmed_missing_receipt
//	attempt1_1 := newBroadcastLegacyEthTxAttempt(t, etx1_1.ID, int64(2))
//	require.NoError(t, txStore.InsertTxAttempt(&attempt1_1))
//	txmReceipt1_1 := newTxReceipt(attempt1_1.Hash, 11, 1)
//
//	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//		return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt1_1.Hash, "eth_getTransactionReceipt")
//	})).Return(nil).Run(func(args mock.Arguments) {
//		elems := args.Get(1).([]rpc.BatchElem)
//		*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt1_1
//	}).Once()
//
//	require.NoError(t, ec.CheckForReceipts(ctx, 11))
//
//	mustTxBeInState(t, txStore, etx1_0, txmgrcommon.TxConfirmedMissingReceipt)
//	mustTxBeInState(t, txStore, etx1_1, txmgrcommon.TxConfirmed)
//	mustTxBeInState(t, txStore, etx2_9, txmgrcommon.TxConfirmed)
//}
//
//func TestEthConfirmer_CheckForReceipts_confirmed_missing_receipt(t *testing.T) {
//	t.Parallel()
//
//	db := pgtest.NewSqlxDB(t)
//	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
//		c.EVM[0].FinalityDepth = ptr[uint32](50)
//	})
//	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
//
//	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
//
//	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
//
//	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
//
//	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
//
//	ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)
//	ctx := testutils.Context(t)
//
//	// STATE
//	// eth_txes with nonce 0 has two attempts (broadcast before block 21 and 41) the first of which will get a receipt
//	// eth_txes with nonce 1 has two attempts (broadcast before block 21 and 41) neither of which will ever get a receipt
//	// eth_txes with nonce 2 has an attempt (broadcast before block 41) that will not get a receipt on the first try but will get one later
//	// eth_txes with nonce 3 has an attempt (broadcast before block 41) that has been confirmed in block 42
//	// All other attempts were broadcast before block 41
//	b := int64(21)
//
//	etx0 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 0, fromAddress)
//	attempt0_1 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(1))
//	attempt0_2 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(2))
//	attempt0_2.BroadcastBeforeBlockNum = &b
//	require.NoError(t, txStore.InsertTxAttempt(&attempt0_1))
//	require.NoError(t, txStore.InsertTxAttempt(&attempt0_2))
//
//	etx1 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 1, fromAddress)
//	attempt1_1 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(1))
//	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(2))
//	attempt1_2.BroadcastBeforeBlockNum = &b
//	require.NoError(t, txStore.InsertTxAttempt(&attempt1_1))
//	require.NoError(t, txStore.InsertTxAttempt(&attempt1_2))
//
//	etx2 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 2, fromAddress)
//	attempt2_1 := newBroadcastLegacyEthTxAttempt(t, etx2.ID, int64(1))
//	require.NoError(t, txStore.InsertTxAttempt(&attempt2_1))
//
//	etx3 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 3, fromAddress)
//	attempt3_1 := newBroadcastLegacyEthTxAttempt(t, etx3.ID, int64(1))
//	require.NoError(t, txStore.InsertTxAttempt(&attempt3_1))
//
//	pgtest.MustExec(t, db, `UPDATE evm.tx_attempts SET broadcast_before_block_num = 41 WHERE broadcast_before_block_num IS NULL`)
//
//	t.Run("marks buried eth_txes as 'confirmed_missing_receipt'", func(t *testing.T) {
//		txmReceipt0 := evmtypes.Receipt{
//			TxHash:           attempt0_2.Hash,
//			BlockHash:        utils.NewHash(),
//			BlockNumber:      big.NewInt(42),
//			TransactionIndex: uint(1),
//			Status:           uint64(1),
//		}
//		txmReceipt3 := evmtypes.Receipt{
//			TxHash:           attempt3_1.Hash,
//			BlockHash:        utils.NewHash(),
//			BlockNumber:      big.NewInt(42),
//			TransactionIndex: uint(1),
//			Status:           uint64(1),
//		}
//		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(4), nil)
//		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//			return len(b) == 6 &&
//				cltest.BatchElemMatchesParams(b[0], attempt0_2.Hash, "eth_getTransactionReceipt") &&
//				cltest.BatchElemMatchesParams(b[1], attempt0_1.Hash, "eth_getTransactionReceipt") &&
//				cltest.BatchElemMatchesParams(b[2], attempt1_2.Hash, "eth_getTransactionReceipt") &&
//				cltest.BatchElemMatchesParams(b[3], attempt1_1.Hash, "eth_getTransactionReceipt") &&
//				cltest.BatchElemMatchesParams(b[4], attempt2_1.Hash, "eth_getTransactionReceipt") &&
//				cltest.BatchElemMatchesParams(b[5], attempt3_1.Hash, "eth_getTransactionReceipt")
//
//		})).Return(nil).Run(func(args mock.Arguments) {
//			elems := args.Get(1).([]rpc.BatchElem)
//			// First transaction confirmed
//			*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt0
//			elems[1].Result = &evmtypes.Receipt{}
//			// Second transaction stil unconfirmed
//			elems[2].Result = &evmtypes.Receipt{}
//			elems[3].Result = &evmtypes.Receipt{}
//			// Third transaction still unconfirmed
//			elems[4].Result = &evmtypes.Receipt{}
//			// Fourth transaction is confirmed
//			*(elems[5].Result.(*evmtypes.Receipt)) = txmReceipt3
//		}).Once()
//
//		// PERFORM
//		// Block num of 43 is one higher than the receipt (as would generally be expected)
//		require.NoError(t, ec.CheckForReceipts(ctx, 43))
//
//		// Expected state is that the "top" eth_tx is now confirmed, with the
//		// two below it "confirmed_missing_receipt" and the "bottom" eth_tx also confirmed
//		var err error
//		etx3, err = txStore.FindTxWithAttempts(etx3.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmed, etx3.State)
//
//		ethReceipt := etx3.TxAttempts[0].Receipts[0]
//		require.Equal(t, txmReceipt3.BlockHash, ethReceipt.GetBlockHash())
//
//		etx2, err = txStore.FindTxWithAttempts(etx2.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx2.State)
//		etx1, err = txStore.FindTxWithAttempts(etx1.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx1.State)
//
//		etx0, err = txStore.FindTxWithAttempts(etx0.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmed, etx0.State)
//
//		require.Len(t, etx0.TxAttempts, 2)
//		require.Len(t, etx0.TxAttempts[0].Receipts, 1)
//		ethReceipt = etx0.TxAttempts[0].Receipts[0]
//		require.Equal(t, txmReceipt0.BlockHash, ethReceipt.GetBlockHash())
//	})
//
//	// STATE
//	// eth_txes with nonce 0 is confirmed
//	// eth_txes with nonce 1 is confirmed_missing_receipt
//	// eth_txes with nonce 2 is confirmed_missing_receipt
//	// eth_txes with nonce 3 is confirmed
//
//	t.Run("marks eth_txes with state 'confirmed_missing_receipt' as 'confirmed' if a receipt finally shows up", func(t *testing.T) {
//		txmReceipt := evmtypes.Receipt{
//			TxHash:           attempt2_1.Hash,
//			BlockHash:        utils.NewHash(),
//			BlockNumber:      big.NewInt(43),
//			TransactionIndex: uint(1),
//			Status:           uint64(1),
//		}
//		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
//		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//			return len(b) == 3 &&
//				cltest.BatchElemMatchesParams(b[0], attempt1_2.Hash, "eth_getTransactionReceipt") &&
//				cltest.BatchElemMatchesParams(b[1], attempt1_1.Hash, "eth_getTransactionReceipt") &&
//				cltest.BatchElemMatchesParams(b[2], attempt2_1.Hash, "eth_getTransactionReceipt")
//
//		})).Return(nil).Run(func(args mock.Arguments) {
//			elems := args.Get(1).([]rpc.BatchElem)
//			// First transaction still unconfirmed
//			elems[0].Result = &evmtypes.Receipt{}
//			elems[1].Result = &evmtypes.Receipt{}
//			// Second transaction confirmed
//			*(elems[2].Result.(*evmtypes.Receipt)) = txmReceipt
//		}).Once()
//
//		// PERFORM
//		// Block num of 44 is one higher than the receipt (as would generally be expected)
//		require.NoError(t, ec.CheckForReceipts(ctx, 44))
//
//		// Expected state is that the "top" two eth_txes are now confirmed, with the
//		// one below it still "confirmed_missing_receipt" and the bottom one remains confirmed
//		var err error
//		etx3, err = txStore.FindTxWithAttempts(etx3.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmed, etx3.State)
//		etx2, err = txStore.FindTxWithAttempts(etx2.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmed, etx2.State)
//
//		ethReceipt := etx2.TxAttempts[0].Receipts[0]
//		require.Equal(t, txmReceipt.BlockHash, ethReceipt.GetBlockHash())
//
//		etx1, err = txStore.FindTxWithAttempts(etx1.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx1.State)
//		etx0, err = txStore.FindTxWithAttempts(etx0.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmed, etx0.State)
//	})
//
//	// STATE
//	// eth_txes with nonce 0 is confirmed
//	// eth_txes with nonce 1 is confirmed_missing_receipt
//	// eth_txes with nonce 2 is confirmed
//	// eth_txes with nonce 3 is confirmed
//
//	t.Run("continues to leave eth_txes with state 'confirmed_missing_receipt' unchanged if at least one attempt is above EVM.FinalityDepth", func(t *testing.T) {
//		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
//		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//			return len(b) == 2 &&
//				cltest.BatchElemMatchesParams(b[0], attempt1_2.Hash, "eth_getTransactionReceipt") &&
//				cltest.BatchElemMatchesParams(b[1], attempt1_1.Hash, "eth_getTransactionReceipt")
//
//		})).Return(nil).Run(func(args mock.Arguments) {
//			elems := args.Get(1).([]rpc.BatchElem)
//			// Both attempts still unconfirmed
//			elems[0].Result = &evmtypes.Receipt{}
//			elems[1].Result = &evmtypes.Receipt{}
//		}).Once()
//
//		// PERFORM
//		// Block num of 80 puts the first attempt (21) below threshold but second attempt (41) still above
//		require.NoError(t, ec.CheckForReceipts(ctx, 80))
//
//		// Expected state is that the "top" two eth_txes are now confirmed, with the
//		// one below it still "confirmed_missing_receipt" and the bottom one remains confirmed
//		var err error
//		etx3, err = txStore.FindTxWithAttempts(etx3.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmed, etx3.State)
//		etx2, err = txStore.FindTxWithAttempts(etx2.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmed, etx2.State)
//		etx1, err = txStore.FindTxWithAttempts(etx1.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx1.State)
//		etx0, err = txStore.FindTxWithAttempts(etx0.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmed, etx0.State)
//	})
//
//	// STATE
//	// eth_txes with nonce 0 is confirmed
//	// eth_txes with nonce 1 is confirmed_missing_receipt
//	// eth_txes with nonce 2 is confirmed
//	// eth_txes with nonce 3 is confirmed
//
//	t.Run("marks eth_Txes with state 'confirmed_missing_receipt' as 'errored' if a receipt fails to show up and all attempts are buried deeper than EVM.FinalityDepth", func(t *testing.T) {
//		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
//		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//			return len(b) == 2 &&
//				cltest.BatchElemMatchesParams(b[0], attempt1_2.Hash, "eth_getTransactionReceipt") &&
//				cltest.BatchElemMatchesParams(b[1], attempt1_1.Hash, "eth_getTransactionReceipt")
//
//		})).Return(nil).Run(func(args mock.Arguments) {
//			elems := args.Get(1).([]rpc.BatchElem)
//			// Both attempts still unconfirmed
//			elems[0].Result = &evmtypes.Receipt{}
//			elems[1].Result = &evmtypes.Receipt{}
//		}).Once()
//
//		// PERFORM
//		// Block num of 100 puts the first attempt (21) and second attempt (41) below threshold
//		require.NoError(t, ec.CheckForReceipts(ctx, 100))
//
//		// Expected state is that the "top" two eth_txes are now confirmed, with the
//		// one below it marked as "fatal_error" and the bottom one remains confirmed
//		var err error
//		etx3, err = txStore.FindTxWithAttempts(etx3.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmed, etx3.State)
//		etx2, err = txStore.FindTxWithAttempts(etx2.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmed, etx2.State)
//		etx1, err = txStore.FindTxWithAttempts(etx1.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxFatalError, etx1.State)
//		etx0, err = txStore.FindTxWithAttempts(etx0.ID)
//		require.NoError(t, err)
//		require.Equal(t, txmgrcommon.TxConfirmed, etx0.State)
//	})
//}
//
//func TestEthConfirmer_CheckConfirmedMissingReceipt(t *testing.T) {
//	t.Parallel()
//
//	db := pgtest.NewSqlxDB(t)
//	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
//		c.EVM[0].FinalityDepth = ptr[uint32](50)
//	})
//	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
//
//	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
//
//	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
//
//	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
//
//	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
//
//	ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)
//	ctx := testutils.Context(t)
//
//	// STATE
//	// eth_txes with nonce 0 has two attempts, the later attempt with higher gas fees
//	// eth_txes with nonce 1 has two attempts, the later attempt with higher gas fees
//	// eth_txes with nonce 2 has one attempt
//	originalBroadcastAt := time.Unix(1616509100, 0)
//	etx0 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
//		t, txStore, 0, 1, originalBroadcastAt, fromAddress)
//	attempt0_2 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(2))
//	require.NoError(t, txStore.InsertTxAttempt(&attempt0_2))
//	etx1 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
//		t, txStore, 1, 1, originalBroadcastAt, fromAddress)
//	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(2))
//	require.NoError(t, txStore.InsertTxAttempt(&attempt1_2))
//	etx2 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
//		t, txStore, 2, 1, originalBroadcastAt, fromAddress)
//	attempt2_1 := etx2.TxAttempts[0]
//	etx3 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
//		t, txStore, 3, 1, originalBroadcastAt, fromAddress)
//	attempt3_1 := etx3.TxAttempts[0]
//
//	ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//		return len(b) == 4 &&
//			cltest.BatchElemMatchesParams(b[0], hexutil.Encode(attempt0_2.SignedRawTx), "eth_sendRawTransaction") &&
//			cltest.BatchElemMatchesParams(b[1], hexutil.Encode(attempt1_2.SignedRawTx), "eth_sendRawTransaction") &&
//			cltest.BatchElemMatchesParams(b[2], hexutil.Encode(attempt2_1.SignedRawTx), "eth_sendRawTransaction") &&
//			cltest.BatchElemMatchesParams(b[3], hexutil.Encode(attempt3_1.SignedRawTx), "eth_sendRawTransaction")
//	})).Return(nil).Run(func(args mock.Arguments) {
//		elems := args.Get(1).([]rpc.BatchElem)
//		// First transaction confirmed
//		elems[0].Error = errors.New("nonce too low")
//		elems[1].Error = errors.New("transaction underpriced")
//		elems[2].Error = nil
//		elems[3].Error = errors.New("transaction already finalized")
//	}).Once()
//
//	// PERFORM
//	require.NoError(t, ec.CheckConfirmedMissingReceipt(ctx))
//
//	// Expected state is that the "top" eth_tx is untouched but the other two
//	// are marked as unconfirmed
//	var err error
//	etx0, err = txStore.FindTxWithAttempts(etx0.ID)
//	assert.NoError(t, err)
//	assert.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx0.State)
//	assert.Greater(t, etx0.BroadcastAt.Unix(), originalBroadcastAt.Unix())
//	etx1, err = txStore.FindTxWithAttempts(etx1.ID)
//	assert.NoError(t, err)
//	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx1.State)
//	assert.Greater(t, etx1.BroadcastAt.Unix(), originalBroadcastAt.Unix())
//	etx2, err = txStore.FindTxWithAttempts(etx2.ID)
//	assert.NoError(t, err)
//	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx2.State)
//	assert.Greater(t, etx2.BroadcastAt.Unix(), originalBroadcastAt.Unix())
//	etx3, err = txStore.FindTxWithAttempts(etx3.ID)
//	assert.NoError(t, err)
//	assert.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx3.State)
//	assert.Greater(t, etx3.BroadcastAt.Unix(), originalBroadcastAt.Unix())
//}
//
//func TestEthConfirmer_CheckConfirmedMissingReceipt_batchSendTransactions_fails(t *testing.T) {
//	t.Parallel()
//
//	db := pgtest.NewSqlxDB(t)
//	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
//		c.EVM[0].FinalityDepth = ptr[uint32](50)
//	})
//	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
//
//	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
//
//	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
//
//	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
//
//	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
//
//	ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)
//	ctx := testutils.Context(t)
//
//	// STATE
//	// eth_txes with nonce 0 has two attempts, the later attempt with higher gas fees
//	// eth_txes with nonce 1 has two attempts, the later attempt with higher gas fees
//	// eth_txes with nonce 2 has one attempt
//	originalBroadcastAt := time.Unix(1616509100, 0)
//	etx0 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
//		t, txStore, 0, 1, originalBroadcastAt, fromAddress)
//	attempt0_2 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(2))
//	require.NoError(t, txStore.InsertTxAttempt(&attempt0_2))
//	etx1 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
//		t, txStore, 1, 1, originalBroadcastAt, fromAddress)
//	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(2))
//	require.NoError(t, txStore.InsertTxAttempt(&attempt1_2))
//	etx2 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
//		t, txStore, 2, 1, originalBroadcastAt, fromAddress)
//	attempt2_1 := etx2.TxAttempts[0]
//
//	ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//		return len(b) == 3 &&
//			cltest.BatchElemMatchesParams(b[0], hexutil.Encode(attempt0_2.SignedRawTx), "eth_sendRawTransaction") &&
//			cltest.BatchElemMatchesParams(b[1], hexutil.Encode(attempt1_2.SignedRawTx), "eth_sendRawTransaction") &&
//			cltest.BatchElemMatchesParams(b[2], hexutil.Encode(attempt2_1.SignedRawTx), "eth_sendRawTransaction")
//	})).Return(errors.New("Timed out")).Once()
//
//	// PERFORM
//	require.NoError(t, ec.CheckConfirmedMissingReceipt(ctx))
//
//	// Expected state is that all txes are marked as unconfirmed, since the batch call had failed
//	var err error
//	etx0, err = txStore.FindTxWithAttempts(etx0.ID)
//	assert.NoError(t, err)
//	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx0.State)
//	assert.Equal(t, etx0.BroadcastAt.Unix(), originalBroadcastAt.Unix())
//	etx1, err = txStore.FindTxWithAttempts(etx1.ID)
//	assert.NoError(t, err)
//	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx1.State)
//	assert.Equal(t, etx1.BroadcastAt.Unix(), originalBroadcastAt.Unix())
//	etx2, err = txStore.FindTxWithAttempts(etx2.ID)
//	assert.NoError(t, err)
//	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx2.State)
//	assert.Equal(t, etx2.BroadcastAt.Unix(), originalBroadcastAt.Unix())
//}
//
//func TestEthConfirmer_CheckConfirmedMissingReceipt_smallEvmRPCBatchSize_middleBatchSendTransactionFails(t *testing.T) {
//	t.Parallel()
//
//	db := pgtest.NewSqlxDB(t)
//	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
//		c.EVM[0].FinalityDepth = ptr[uint32](50)
//		c.EVM[0].RPCDefaultBatchSize = ptr[uint32](1)
//	})
//	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
//
//	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
//
//	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
//
//	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
//
//	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
//
//	ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, nil)
//	ctx := testutils.Context(t)
//
//	// STATE
//	// eth_txes with nonce 0 has two attempts, the later attempt with higher gas fees
//	// eth_txes with nonce 1 has two attempts, the later attempt with higher gas fees
//	// eth_txes with nonce 2 has one attempt
//	originalBroadcastAt := time.Unix(1616509100, 0)
//	etx0 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
//		t, txStore, 0, 1, originalBroadcastAt, fromAddress)
//	attempt0_2 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(2))
//	require.NoError(t, txStore.InsertTxAttempt(&attempt0_2))
//	etx1 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
//		t, txStore, 1, 1, originalBroadcastAt, fromAddress)
//	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(2))
//	require.NoError(t, txStore.InsertTxAttempt(&attempt1_2))
//	etx2 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
//		t, txStore, 2, 1, originalBroadcastAt, fromAddress)
//
//	// Expect eth_sendRawTransaction in 3 batches. First batch will pass, 2nd will fail, 3rd never attempted.
//	ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//		return len(b) == 1 &&
//			cltest.BatchElemMatchesParams(b[0], hexutil.Encode(attempt0_2.SignedRawTx), "eth_sendRawTransaction")
//	})).Return(nil).Run(func(args mock.Arguments) {
//		elems := args.Get(1).([]rpc.BatchElem)
//		// First transaction confirmed
//		elems[0].Error = errors.New("nonce too low")
//	}).Once()
//	ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
//		return len(b) == 1 &&
//			cltest.BatchElemMatchesParams(b[0], hexutil.Encode(attempt1_2.SignedRawTx), "eth_sendRawTransaction")
//	})).Return(errors.New("Timed out")).Once()
//
//	// PERFORM
//	require.NoError(t, ec.CheckConfirmedMissingReceipt(ctx))
//
//	// Expected state is that all transactions since failed batch will be unconfirmed
//	var err error
//	etx0, err = txStore.FindTxWithAttempts(etx0.ID)
//	assert.NoError(t, err)
//	assert.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx0.State)
//	assert.Greater(t, etx0.BroadcastAt.Unix(), originalBroadcastAt.Unix())
//	etx1, err = txStore.FindTxWithAttempts(etx1.ID)
//	assert.NoError(t, err)
//	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx1.State)
//	assert.Equal(t, etx1.BroadcastAt.Unix(), originalBroadcastAt.Unix())
//	etx2, err = txStore.FindTxWithAttempts(etx2.ID)
//	assert.NoError(t, err)
//	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx2.State)
//	assert.Equal(t, etx2.BroadcastAt.Unix(), originalBroadcastAt.Unix())
//}
//
//
//func TestEthConfirmer_EnsureConfirmedTransactionsInLongestChain(t *testing.T) {
//	t.Parallel()
//
//	db := pgtest.NewSqlxDB(t)
//	cfg := configtest.NewTestGeneralConfig(t)
//	txStore := cltest.NewTestTxStore(t, db, cfg.Database())
//
//	ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
//
//	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
//
//	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
//
//	config := newTestChainScopedConfig(t)
//	ec := newEthConfirmer(t, txStore, ethClient, config, ethKeyStore, nil)
//
//	head := evmtypes.Head{
//		Hash:   utils.NewHash(),
//		Number: 10,
//		Parent: &evmtypes.Head{
//			Hash:   utils.NewHash(),
//			Number: 9,
//			Parent: &evmtypes.Head{
//				Number: 8,
//				Hash:   utils.NewHash(),
//				Parent: nil,
//			},
//		},
//	}
//
//	t.Run("does nothing if there aren't any transactions", func(t *testing.T) {
//		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))
//	})
//
//	t.Run("does nothing to unconfirmed transactions", func(t *testing.T) {
//		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress)
//
//		// Do the thing
//		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))
//
//		etx, err := txStore.FindTxWithAttempts(etx.ID)
//		require.NoError(t, err)
//		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
//	})
//
//	t.Run("does nothing to confirmed transactions with receipts within head height of the chain and included in the chain", func(t *testing.T) {
//		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 2, 1, fromAddress)
//		mustInsertEthReceipt(t, txStore, head.Number, head.Hash, etx.TxAttempts[0].Hash)
//
//		// Do the thing
//		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))
//
//		etx, err := txStore.FindTxWithAttempts(etx.ID)
//		require.NoError(t, err)
//		assert.Equal(t, txmgrcommon.TxConfirmed, etx.State)
//	})
//
//	t.Run("does nothing to confirmed transactions that only have receipts older than the start of the chain", func(t *testing.T) {
//		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 3, 1, fromAddress)
//		// Add receipt that is older than the lowest block of the chain
//		mustInsertEthReceipt(t, txStore, head.Parent.Parent.Number-1, utils.NewHash(), etx.TxAttempts[0].Hash)
//
//		// Do the thing
//		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))
//
//		etx, err := txStore.FindTxWithAttempts(etx.ID)
//		require.NoError(t, err)
//		assert.Equal(t, txmgrcommon.TxConfirmed, etx.State)
//	})
//
//	t.Run("unconfirms and rebroadcasts transactions that have receipts within head height of the chain but not included in the chain", func(t *testing.T) {
//		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 4, 1, fromAddress)
//		attempt := etx.TxAttempts[0]
//		// Include one within head height but a different block hash
//		mustInsertEthReceipt(t, txStore, head.Parent.Number, utils.NewHash(), attempt.Hash)
//
//		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
//			atx, err := txmgr.GetGethSignedTx(attempt.SignedRawTx)
//			require.NoError(t, err)
//			// Keeps gas price and nonce the same
//			return atx.GasPrice().Cmp(tx.GasPrice()) == 0 && atx.Nonce() == tx.Nonce()
//		}), fromAddress).Return(commonclient.Successful, nil).Once()
//
//		// Do the thing
//		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))
//
//		etx, err := txStore.FindTxWithAttempts(etx.ID)
//		require.NoError(t, err)
//		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
//		require.Len(t, etx.TxAttempts, 1)
//		attempt = etx.TxAttempts[0]
//		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
//	})
//
//	t.Run("unconfirms and rebroadcasts transactions that have receipts within head height of chain but not included in the chain even if a receipt exists older than the start of the chain", func(t *testing.T) {
//		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 5, 1, fromAddress)
//		attempt := etx.TxAttempts[0]
//		attemptHash := attempt.Hash
//		// Add receipt that is older than the lowest block of the chain
//		mustInsertEthReceipt(t, txStore, head.Parent.Parent.Number-1, utils.NewHash(), attemptHash)
//		// Include one within head height but a different block hash
//		mustInsertEthReceipt(t, txStore, head.Parent.Number, utils.NewHash(), attemptHash)
//
//		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
//			commonclient.Successful, nil).Once()
//
//		// Do the thing
//		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))
//
//		etx, err := txStore.FindTxWithAttempts(etx.ID)
//		require.NoError(t, err)
//		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
//		require.Len(t, etx.TxAttempts, 1)
//		attempt = etx.TxAttempts[0]
//		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
//	})
//
//	t.Run("if more than one attempt has a receipt (should not be possible but isn't prevented by database constraints) unconfirms and rebroadcasts only the attempt with the highest gas price", func(t *testing.T) {
//		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 6, 1, fromAddress)
//		require.Len(t, etx.TxAttempts, 1)
//		// Sanity check to assert the included attempt has the lowest gas price
//		require.Less(t, etx.TxAttempts[0].TxFee.Legacy.ToInt().Int64(), int64(30000))
//
//		attempt2 := newBroadcastLegacyEthTxAttempt(t, etx.ID, 30000)
//		attempt2.SignedRawTx = hexutil.MustDecode("0xf88c8301f3a98503b9aca000832ab98094f5fff180082d6017036b771ba883025c654bc93580a4daa6d556000000000000000000000000000000000000000000000000000000000000000026a0f25601065ee369b6470c0399a2334afcfbeb0b5c8f3d9a9042e448ed29b5bcbda05b676e00248b85faf4dd889f0e2dcf91eb867e23ac9eeb14a73f9e4c14972cdf")
//		attempt3 := newBroadcastLegacyEthTxAttempt(t, etx.ID, 40000)
//		attempt3.SignedRawTx = hexutil.MustDecode("0xf88c8301f3a88503b9aca0008316e36094151445852b0cfdf6a4cc81440f2af99176e8ad0880a4daa6d556000000000000000000000000000000000000000000000000000000000000000026a0dcb5a7ad52b96a866257134429f944c505820716567f070e64abb74899803855a04c13eff2a22c218e68da80111e1bb6dc665d3dea7104ab40ff8a0275a99f630d")
//		require.NoError(t, txStore.InsertTxAttempt(&attempt2))
//		require.NoError(t, txStore.InsertTxAttempt(&attempt3))
//
//		// Receipt is within head height but a different block hash
//		mustInsertEthReceipt(t, txStore, head.Parent.Number, utils.NewHash(), attempt2.Hash)
//		// Receipt is within head height but a different block hash
//		mustInsertEthReceipt(t, txStore, head.Parent.Number, utils.NewHash(), attempt3.Hash)
//
//		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
//			s, err := txmgr.GetGethSignedTx(attempt3.SignedRawTx)
//			require.NoError(t, err)
//			return tx.Hash() == s.Hash()
//		}), fromAddress).Return(commonclient.Successful, nil).Once()
//
//		// Do the thing
//		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))
//
//		etx, err := txStore.FindTxWithAttempts(etx.ID)
//		require.NoError(t, err)
//		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
//		require.Len(t, etx.TxAttempts, 3)
//		attempt1 := etx.TxAttempts[0]
//		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt1.State)
//		attempt2 = etx.TxAttempts[1]
//		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt2.State)
//		attempt3 = etx.TxAttempts[2]
//		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt3.State)
//	})
//
//	t.Run("if receipt has a block number that is in the future, does not mark for rebroadcast (the safe thing to do is simply wait until heads catches up)", func(t *testing.T) {
//		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 7, 1, fromAddress)
//		attempt := etx.TxAttempts[0]
//		// Add receipt that is higher than head
//		mustInsertEthReceipt(t, txStore, head.Number+1, utils.NewHash(), attempt.Hash)
//
//		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))
//
//		etx, err := txStore.FindTxWithAttempts(etx.ID)
//		require.NoError(t, err)
//		assert.Equal(t, txmgrcommon.TxConfirmed, etx.State)
//		require.Len(t, etx.TxAttempts, 1)
//		attempt = etx.TxAttempts[0]
//		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
//		assert.Len(t, attempt.Receipts, 1)
//	})
//}
//

func TestEthConfirmer_ResumePendingRuns(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	config := configtest.NewTestGeneralConfig(t)
	txStore := cltest.NewTestTxStore(t, db, config.Database())

	ethKeyStore := cltest.NewKeyStore(t, db, config.Database()).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)

	evmcfg := evmtest.NewChainScopedConfig(t, config)

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

	pgtest.MustExec(t, db, `SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`)

	t.Run("doesn't process task runs that are not suspended (possibly already previously resumed)", func(t *testing.T) {
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, func(uuid.UUID, interface{}, error) error {
			t.Fatal("No value expected")
			return nil
		})

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 1, 1, fromAddress)
		mustInsertEthReceipt(t, txStore, head.Number-minConfirmations, head.Hash, etx.TxAttempts[0].Hash)
		// Setting both signal_callback and callback_completed to TRUE to simulate a completed pipeline task
		// It would only be in a state past suspended if the resume callback was called and callback_completed was set to TRUE
		pgtest.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE, callback_completed = TRUE WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		err := ec.ResumePendingTaskRuns(testutils.Context(t), &head)
		require.NoError(t, err)
	})

	t.Run("doesn't process task runs where the receipt is younger than minConfirmations", func(t *testing.T) {
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, func(uuid.UUID, interface{}, error) error {
			t.Fatal("No value expected")
			return nil
		})

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 2, 1, fromAddress)
		mustInsertEthReceipt(t, txStore, head.Number, head.Hash, etx.TxAttempts[0].Hash)

		pgtest.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		err := ec.ResumePendingTaskRuns(testutils.Context(t), &head)
		require.NoError(t, err)
	})

	t.Run("processes eth_txes with receipts older than minConfirmations", func(t *testing.T) {
		ch := make(chan interface{})
		nonce := evmtypes.Nonce(3)
		var err error
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, func(id uuid.UUID, value interface{}, thisErr error) error {
			err = thisErr
			ch <- value
			return nil
		})

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)
		pgtest.MustExec(t, db, `UPDATE pipeline_runs SET state = 'suspended' WHERE id = $1`, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, int64(nonce), 1, fromAddress)
		pgtest.MustExec(t, db, `UPDATE evm.txes SET meta='{"FailOnRevert": true}'`)
		receipt := mustInsertEthReceipt(t, txStore, head.Number-minConfirmations, head.Hash, etx.TxAttempts[0].Hash)

		pgtest.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		done := make(chan struct{})
		t.Cleanup(func() { <-done })
		go func() {
			defer close(done)
			err2 := ec.ResumePendingTaskRuns(testutils.Context(t), &head)
			if !assert.NoError(t, err2) {
				return
			}
			// Retrieve Tx to check if callback completed flag was set to true
			updateTx, err3 := txStore.FindTxWithSequence(testutils.Context(t), fromAddress, nonce)
			if assert.NoError(t, err3) {
				assert.Equal(t, true, updateTx.CallbackCompleted)
			}
		}()

		select {
		case data := <-ch:
			assert.NoError(t, err)

			require.IsType(t, &evmtypes.Receipt{}, data)
			r := data.(*evmtypes.Receipt)
			require.Equal(t, receipt.TxHash, r.TxHash)

		case <-time.After(time.Second):
			t.Fatal("no value received")
		}
	})

	pgtest.MustExec(t, db, `DELETE FROM pipeline_runs`)

	t.Run("processes eth_txes with receipt older than minConfirmations that reverted", func(t *testing.T) {
		type data struct {
			value any
			error
		}
		ch := make(chan data)
		nonce := evmtypes.Nonce(4)
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, func(id uuid.UUID, value interface{}, err error) error {
			ch <- data{value, err}
			return nil
		})

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)
		pgtest.MustExec(t, db, `UPDATE pipeline_runs SET state = 'suspended' WHERE id = $1`, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, int64(nonce), 1, fromAddress)
		pgtest.MustExec(t, db, `UPDATE evm.txes SET meta='{"FailOnRevert": true}'`)

		// receipt is not passed through as a value since it reverted and caused an error
		mustInsertRevertedEthReceipt(t, txStore, head.Number-minConfirmations, head.Hash, etx.TxAttempts[0].Hash)

		pgtest.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		done := make(chan struct{})
		t.Cleanup(func() { <-done })
		go func() {
			defer close(done)
			err2 := ec.ResumePendingTaskRuns(testutils.Context(t), &head)
			if !assert.NoError(t, err2) {
				return
			}
			// Retrieve Tx to check if callback completed flag was set to true
			updateTx, err3 := txStore.FindTxWithSequence(testutils.Context(t), fromAddress, nonce)
			if assert.NoError(t, err3) {
				assert.Equal(t, true, updateTx.CallbackCompleted)
			}
		}()

		select {
		case data := <-ch:
			assert.Error(t, data.error)

			assert.EqualError(t, data.error, fmt.Sprintf("transaction %s reverted on-chain", etx.TxAttempts[0].Hash.String()))

			assert.Nil(t, data.value)

		case <-testutils.AfterWaitTimeout(t):
			t.Fatal("no value received")
		}
	})

	t.Run("does not mark callback complete if callback fails", func(t *testing.T) {
		nonce := evmtypes.Nonce(5)
		ec := newEthConfirmer(t, txStore, ethClient, evmcfg, ethKeyStore, func(uuid.UUID, interface{}, error) error {
			return errors.New("error")
		})

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, int64(nonce), 1, fromAddress)
		mustInsertEthReceipt(t, txStore, head.Number-minConfirmations, head.Hash, etx.TxAttempts[0].Hash)
		pgtest.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		err := ec.ResumePendingTaskRuns(testutils.Context(t), &head)
		require.Error(t, err)

		// Retrieve Tx to check if callback completed flag was left unchanged
		updateTx, err := txStore.FindTxWithSequence(testutils.Context(t), fromAddress, nonce)
		require.NoError(t, err)
		require.Equal(t, false, updateTx.CallbackCompleted)
	})
}

func newEthConfirmer(t testing.TB, txStore txmgr.EvmTxStore, ethClient client.Client, config evmconfig.ChainScopedConfig, ks keystore.Eth, fn txmgrcommon.ResumeCallback) *txmgr.Confirmer {
	lggr := logger.Test(t)
	ec := txmgr.NewEvmConfirmer(txStore, txmgr.NewEvmTxmClient(ethClient), txmgr.NewEvmTxmConfig(config.EVM()), config.EVM().Transactions(), ks, lggr)
	ec.SetResumeCallback(fn)
	servicetest.Run(t, ec)
	return ec
}
