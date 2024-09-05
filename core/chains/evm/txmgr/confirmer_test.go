package txmgr_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	commonfee "github.com/smartcontractkit/chainlink/v2/common/fee"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	gasmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/keystore"
	ksmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/keystore/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

func newTestChainScopedConfig(t *testing.T) (chainlink.GeneralConfig, evmconfig.ChainScopedConfig) {
	cfg := configtest.NewTestGeneralConfig(t)
	return cfg, evmtest.NewChainScopedConfig(t, cfg)
}

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
	etx, err := txStore.FindTxWithAttempts(tests.Context(t), tx.ID)
	require.NoError(t, err)
	require.Equal(t, expectedState, etx.State)
}

func newTxReceipt(hash gethCommon.Hash, blockNumber int, txIndex uint) evmtypes.Receipt {
	return evmtypes.Receipt{
		TxHash:           hash,
		BlockHash:        testutils.NewHash(),
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
	require.NoError(t, txStore.InsertTx(tests.Context(t), &etx))

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
	require.NoError(t, txStore.InsertTx(tests.Context(t), &etx))

	return etx
}

func TestEthConfirmer_Lifecycle(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	gconfig, config := newTestChainScopedConfig(t)
	txStore := newTxStore(t, db)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	// Add some fromAddresses
	cltest.MustInsertRandomKey(t, ethKeyStore)
	cltest.MustInsertRandomKey(t, ethKeyStore)
	estimator := gasmocks.NewEvmEstimator(t)
	newEst := func(logger.Logger) gas.EvmEstimator { return estimator }
	lggr := logger.Test(t)
	ge := config.EVM().GasEstimator()
	feeEstimator := gas.NewEvmFeeEstimator(lggr, newEst, ge.EIP1559DynamicFees(), ge, ethClient)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, ethKeyStore, feeEstimator)
	stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, "", assets.NewWei(assets.NewEth(100).ToInt()), config.EVM().Transactions().AutoPurge(), feeEstimator, txStore, ethClient)
	ht := headtracker.NewSimulatedHeadTracker(ethClient, true, 0)
	ec := txmgr.NewEvmConfirmer(txStore, txmgr.NewEvmTxmClient(ethClient, nil), txmgr.NewEvmTxmConfig(config.EVM()), txmgr.NewEvmTxmFeeConfig(ge), config.EVM().Transactions(), gconfig.Database(), ethKeyStore, txBuilder, lggr, stuckTxDetector, ht)
	ctx := tests.Context(t)

	// Can't close unstarted instance
	err := ec.Close()
	require.Error(t, err)

	// Can successfully start once
	err = ec.Start(ctx)
	require.NoError(t, err)

	// Can't start an already started instance
	err = ec.Start(ctx)
	require.Error(t, err)

	latestFinalizedHead := &evmtypes.Head{
		Number: 8,
		Hash:   testutils.NewHash(),
	}
	// We are guaranteed to receive a latestFinalizedHead.
	latestFinalizedHead.IsFinalized.Store(true)

	h9 := &evmtypes.Head{
		Hash:   testutils.NewHash(),
		Number: 9,
	}
	h9.Parent.Store(latestFinalizedHead)
	head := &evmtypes.Head{
		Hash:   testutils.NewHash(),
		Number: 10,
	}
	head.Parent.Store(h9)

	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(head, nil).Once()
	ethClient.On("LatestFinalizedBlock", mock.Anything).Return(latestFinalizedHead, nil).Once()

	err = ec.ProcessHead(ctx, head)
	require.NoError(t, err)
	// Can successfully close once
	err = ec.Close()
	require.NoError(t, err)

	// Can't start more than once (Confirmer uses services.StateMachine)
	err = ec.Start(ctx)
	require.Error(t, err)
	// Can't close more than once (Confirmer use services.StateMachine)
	err = ec.Close()
	require.Error(t, err)

	// Can't closeInternal unstarted instance
	require.Error(t, ec.XXXTestCloseInternal())

	// Can successfully startInternal a previously closed instance
	require.NoError(t, ec.XXXTestStartInternal())
	// Can't startInternal already started instance
	require.Error(t, ec.XXXTestStartInternal())
	// Can successfully closeInternal again
	require.NoError(t, ec.XXXTestCloseInternal())
}

func TestEthConfirmer_CheckForReceipts(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	gconfig, config := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	ec := newEthConfirmer(t, txStore, ethClient, gconfig, config, ethKeyStore, nil)

	nonce := int64(0)
	ctx := tests.Context(t)
	blockNum := int64(0)
	latestFinalizedBlockNum := int64(0)

	t.Run("only finds eth_txes in unconfirmed state with at least one broadcast attempt", func(t *testing.T) {
		mustInsertFatalErrorEthTx(t, txStore, fromAddress)
		mustInsertInProgressEthTx(t, txStore, nonce, fromAddress)
		nonce++
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, nonce, 1, fromAddress)
		nonce++
		mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, nonce, fromAddress)
		nonce++
		mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, config.EVM().ChainID())

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum, latestFinalizedBlockNum))
	})

	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	require.Len(t, etx1.TxAttempts, 1)
	attempt1_1 := etx1.TxAttempts[0]
	hashAttempt1_1 := attempt1_1.Hash
	require.Len(t, attempt1_1.Receipts, 0)

	t.Run("fetches receipt for one unconfirmed eth_tx", func(t *testing.T) {
		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
		// Transaction not confirmed yet, receipt is nil
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], hashAttempt1_1, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &evmtypes.Receipt{}
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum, latestFinalizedBlockNum))

		var err error
		etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
		assert.NoError(t, err)
		require.Len(t, etx1.TxAttempts, 1)
		attempt1_1 = etx1.TxAttempts[0]
		require.NoError(t, err)
		require.Len(t, attempt1_1.Receipts, 0)
	})

	t.Run("saves nothing if returned receipt does not match the attempt", func(t *testing.T) {
		txmReceipt := evmtypes.Receipt{
			TxHash:           testutils.NewHash(),
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}

		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
		// First transaction confirmed
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], hashAttempt1_1, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt
		}).Once()

		// No error because it is merely logged
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum, latestFinalizedBlockNum))

		etx, err := txStore.FindTxWithAttempts(ctx, etx1.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 1)

		require.Len(t, etx.TxAttempts[0].Receipts, 0)
	})

	t.Run("saves nothing if query returns error", func(t *testing.T) {
		txmReceipt := evmtypes.Receipt{
			TxHash:           attempt1_1.Hash,
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}

		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
		// First transaction confirmed
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], hashAttempt1_1, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt
			elems[0].Error = errors.New("foo")
		}).Once()

		// No error because it is merely logged
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum, latestFinalizedBlockNum))

		etx, err := txStore.FindTxWithAttempts(ctx, etx1.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 1)
		require.Len(t, etx.TxAttempts[0].Receipts, 0)
	})

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	require.Len(t, etx2.TxAttempts, 1)
	attempt2_1 := etx2.TxAttempts[0]
	require.Len(t, attempt2_1.Receipts, 0)

	t.Run("saves eth_receipt and marks eth_tx as confirmed when geth client returns valid receipt", func(t *testing.T) {
		txmReceipt := evmtypes.Receipt{
			TxHash:           attempt1_1.Hash,
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
			Status:           uint64(1),
		}

		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				cltest.BatchElemMatchesParams(b[0], attempt1_1.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[1], attempt2_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// First transaction confirmed
			*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt
			// Second transaction still unconfirmed
			elems[1].Result = &evmtypes.Receipt{}
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum, latestFinalizedBlockNum))

		// Check that the receipt was saved
		etx, err := txStore.FindTxWithAttempts(ctx, etx1.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxConfirmed, etx.State)
		assert.Len(t, etx.TxAttempts, 1)
		attempt1_1 = etx.TxAttempts[0]
		require.Len(t, attempt1_1.Receipts, 1)

		ethReceipt := attempt1_1.Receipts[0]

		assert.Equal(t, txmReceipt.TxHash, ethReceipt.GetTxHash())
		assert.Equal(t, txmReceipt.BlockHash, ethReceipt.GetBlockHash())
		assert.Equal(t, txmReceipt.BlockNumber.Int64(), ethReceipt.GetBlockNumber().Int64())
		assert.Equal(t, txmReceipt.TransactionIndex, ethReceipt.GetTransactionIndex())

		receiptJSON, err := json.Marshal(txmReceipt)
		require.NoError(t, err)

		j, err := json.Marshal(ethReceipt)
		require.NoError(t, err)
		assert.JSONEq(t, string(receiptJSON), string(j))
	})

	t.Run("fetches and saves receipts for several attempts in gas price order", func(t *testing.T) {
		attempt2_2 := newBroadcastLegacyEthTxAttempt(t, etx2.ID)
		attempt2_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(10)}

		attempt2_3 := newBroadcastLegacyEthTxAttempt(t, etx2.ID)
		attempt2_3.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(20)}

		// Insert order deliberately reversed to test sorting by gas price
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt2_3))
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt2_2))

		txmReceipt := evmtypes.Receipt{
			TxHash:           attempt2_2.Hash,
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
			Status:           uint64(1),
		}

		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 3 &&
				cltest.BatchElemMatchesParams(b[2], attempt2_1.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[1], attempt2_2.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[0], attempt2_3.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// Most expensive attempt still unconfirmed
			elems[2].Result = &evmtypes.Receipt{}
			// Second most expensive attempt is confirmed
			*(elems[1].Result.(*evmtypes.Receipt)) = txmReceipt
			// Cheapest attempt still unconfirmed
			elems[0].Result = &evmtypes.Receipt{}
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum, latestFinalizedBlockNum))

		// Check that the state was updated
		etx, err := txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(t, err)

		require.Equal(t, txmgrcommon.TxConfirmed, etx.State)
		require.Len(t, etx.TxAttempts, 3)
	})

	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	attempt3_1 := etx3.TxAttempts[0]
	nonce++

	t.Run("ignores receipt missing BlockHash that comes from querying parity too early", func(t *testing.T) {
		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
		receipt := evmtypes.Receipt{
			TxHash: attempt3_1.Hash,
			Status: uint64(1),
		}
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt3_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			*(elems[0].Result.(*evmtypes.Receipt)) = receipt
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum, latestFinalizedBlockNum))

		// No receipt, but no error either
		etx, err := txStore.FindTxWithAttempts(ctx, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
		assert.Len(t, etx.TxAttempts, 1)
		attempt3_1 = etx.TxAttempts[0]
		require.Len(t, attempt3_1.Receipts, 0)
	})

	t.Run("does not panic if receipt has BlockHash but is missing some other fields somehow", func(t *testing.T) {
		// NOTE: This should never happen, but we shouldn't panic regardless
		receipt := evmtypes.Receipt{
			TxHash:    attempt3_1.Hash,
			BlockHash: testutils.NewHash(),
			Status:    uint64(1),
		}
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt3_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			*(elems[0].Result.(*evmtypes.Receipt)) = receipt
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum, latestFinalizedBlockNum))

		// No receipt, but no error either
		etx, err := txStore.FindTxWithAttempts(ctx, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
		assert.Len(t, etx.TxAttempts, 1)
		attempt3_1 = etx.TxAttempts[0]
		require.Len(t, attempt3_1.Receipts, 0)
	})
	t.Run("handles case where eth_receipt already exists somehow", func(t *testing.T) {
		ethReceipt := mustInsertEthReceipt(t, txStore, 42, testutils.NewHash(), attempt3_1.Hash)
		txmReceipt := evmtypes.Receipt{
			TxHash:           attempt3_1.Hash,
			BlockHash:        ethReceipt.BlockHash,
			BlockNumber:      big.NewInt(ethReceipt.BlockNumber),
			TransactionIndex: ethReceipt.TransactionIndex,
			Status:           uint64(1),
		}
		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt3_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum, latestFinalizedBlockNum))

		// Check that the receipt was unchanged
		etx, err := txStore.FindTxWithAttempts(ctx, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxConfirmed, etx.State)
		assert.Len(t, etx.TxAttempts, 1)
		attempt3_1 = etx.TxAttempts[0]
		require.Len(t, attempt3_1.Receipts, 1)

		ethReceipt3_1 := attempt3_1.Receipts[0]

		assert.Equal(t, txmReceipt.TxHash, ethReceipt3_1.GetTxHash())
		assert.Equal(t, txmReceipt.BlockHash, ethReceipt3_1.GetBlockHash())
		assert.Equal(t, txmReceipt.BlockNumber.Int64(), ethReceipt3_1.GetBlockNumber().Int64())
		assert.Equal(t, txmReceipt.TransactionIndex, ethReceipt3_1.GetTransactionIndex())
	})

	etx4 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	attempt4_1 := etx4.TxAttempts[0]
	nonce++

	t.Run("on receipt fetch marks in_progress eth_tx_attempt as broadcast", func(t *testing.T) {
		attempt4_2 := newInProgressLegacyEthTxAttempt(t, etx4.ID)
		attempt4_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(10)}

		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt4_2))

		txmReceipt := evmtypes.Receipt{
			TxHash:           attempt4_2.Hash,
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
			Status:           uint64(1),
		}
		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
		// Second attempt is confirmed
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				cltest.BatchElemMatchesParams(b[0], attempt4_2.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[1], attempt4_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// First attempt still unconfirmed
			elems[1].Result = &evmtypes.Receipt{}
			// Second attempt is confirmed
			*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum, latestFinalizedBlockNum))

		// Check that the state was updated
		var err error
		etx4, err = txStore.FindTxWithAttempts(ctx, etx4.ID)
		require.NoError(t, err)

		attempt4_1 = etx4.TxAttempts[1]
		attempt4_2 = etx4.TxAttempts[0]

		// And the attempts
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt4_1.State)
		require.Nil(t, attempt4_1.BroadcastBeforeBlockNum)
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt4_2.State)
		require.Equal(t, int64(42), *attempt4_2.BroadcastBeforeBlockNum)

		// Check receipts
		require.Len(t, attempt4_1.Receipts, 0)
		require.Len(t, attempt4_2.Receipts, 1)
	})

	etx5 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	attempt5_1 := etx5.TxAttempts[0]
	nonce++

	t.Run("simulate on revert", func(t *testing.T) {
		txmReceipt := evmtypes.Receipt{
			TxHash:           attempt5_1.Hash,
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
			Status:           uint64(0),
		}
		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
		// First attempt is confirmed and reverted
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 &&
				cltest.BatchElemMatchesParams(b[0], attempt5_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// First attempt still unconfirmed
			*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt
		}).Once()
		data, err := utils.ABIEncode(`[{"type":"uint256"}]`, big.NewInt(10))
		require.NoError(t, err)
		sig := utils.Keccak256Fixed([]byte(`MyError(uint256)`))
		ethClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(nil, &client.JsonError{
			Code:    1,
			Message: "reverted",
			Data:    utils.ConcatBytes(sig[:4], data),
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum, latestFinalizedBlockNum))

		// Check that the state was updated
		etx5, err = txStore.FindTxWithAttempts(ctx, etx5.ID)
		require.NoError(t, err)

		attempt5_1 = etx5.TxAttempts[0]

		// And the attempts
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt5_1.State)
		require.NotNil(t, attempt5_1.BroadcastBeforeBlockNum)
		// Check receipts
		require.Len(t, attempt5_1.Receipts, 1)
	})
}

func TestEthConfirmer_CheckForReceipts_batching(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].RPCDefaultBatchSize = ptr[uint32](2)
	})
	txStore := cltest.NewTestTxStore(t, db)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, ethKeyStore, nil)
	ctx := tests.Context(t)

	etx := cltest.MustInsertUnconfirmedEthTx(t, txStore, 0, fromAddress)
	var attempts []txmgr.TxAttempt
	latestFinalizedBlockNum := int64(0)

	// Total of 5 attempts should lead to 3 batched fetches (2, 2, 1)
	for i := 0; i < 5; i++ {
		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, int64(i+2))
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
		attempts = append(attempts, attempt)
	}

	ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)

	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 2 &&
			cltest.BatchElemMatchesParams(b[0], attempts[4].Hash, "eth_getTransactionReceipt") &&
			cltest.BatchElemMatchesParams(b[1], attempts[3].Hash, "eth_getTransactionReceipt")
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &evmtypes.Receipt{}
		elems[1].Result = &evmtypes.Receipt{}
	}).Once()
	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 2 &&
			cltest.BatchElemMatchesParams(b[0], attempts[2].Hash, "eth_getTransactionReceipt") &&
			cltest.BatchElemMatchesParams(b[1], attempts[1].Hash, "eth_getTransactionReceipt")
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &evmtypes.Receipt{}
		elems[1].Result = &evmtypes.Receipt{}
	}).Once()
	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 1 &&
			cltest.BatchElemMatchesParams(b[0], attempts[0].Hash, "eth_getTransactionReceipt")
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &evmtypes.Receipt{}
	}).Once()

	require.NoError(t, ec.CheckForReceipts(ctx, 42, latestFinalizedBlockNum))
}

func TestEthConfirmer_CheckForReceipts_HandlesNonFwdTxsWithForwardingEnabled(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].RPCDefaultBatchSize = ptr[uint32](1)
		c.EVM[0].Transactions.ForwardersEnabled = ptr(true)
	})

	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, ethKeyStore, nil)
	ctx := tests.Context(t)
	latestFinalizedBlockNum := int64(0)

	// tx is not forwarded and doesn't have meta set. EthConfirmer should handle nil meta values
	etx := cltest.MustInsertUnconfirmedEthTx(t, txStore, 0, fromAddress)
	attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, 2)
	attempt.Tx.Meta = nil
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	dbtx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
	require.NoError(t, err)
	require.Equal(t, 0, len(dbtx.TxAttempts[0].Receipts))

	txmReceipt := evmtypes.Receipt{
		TxHash:           attempt.Hash,
		BlockHash:        testutils.NewHash(),
		BlockNumber:      big.NewInt(42),
		TransactionIndex: uint(1),
		Status:           uint64(1),
	}

	ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 1 &&
			cltest.BatchElemMatchesParams(b[0], attempt.Hash, "eth_getTransactionReceipt")
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt // confirmed
	}).Once()

	require.NoError(t, ec.CheckForReceipts(ctx, 42, latestFinalizedBlockNum))

	// Check receipt is inserted correctly.
	dbtx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
	require.NoError(t, err)
	require.Equal(t, 1, len(dbtx.TxAttempts[0].Receipts))
}

func TestEthConfirmer_CheckForReceipts_only_likely_confirmed(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].RPCDefaultBatchSize = ptr[uint32](6)
	})
	txStore := cltest.NewTestTxStore(t, db)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, ethKeyStore, nil)
	ctx := tests.Context(t)
	latestFinalizedBlockNum := int64(0)

	var attempts []txmgr.TxAttempt
	// inserting in DESC nonce order to test DB ASC ordering
	etx2 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 1, fromAddress)
	for i := 0; i < 4; i++ {
		attempt := newBroadcastLegacyEthTxAttempt(t, etx2.ID, int64(100-i))
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	}
	etx := cltest.MustInsertUnconfirmedEthTx(t, txStore, 0, fromAddress)
	for i := 0; i < 4; i++ {
		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, int64(100-i))
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))

		// only adding these because a batch for only those attempts should be sent
		attempts = append(attempts, attempt)
	}

	ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(0), nil)

	var captured []rpc.BatchElem
	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 4
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		captured = append(captured, elems...)
		elems[0].Result = &evmtypes.Receipt{}
		elems[1].Result = &evmtypes.Receipt{}
		elems[2].Result = &evmtypes.Receipt{}
		elems[3].Result = &evmtypes.Receipt{}
	}).Once()

	require.NoError(t, ec.CheckForReceipts(ctx, 42, latestFinalizedBlockNum))

	cltest.BatchElemMustMatchParams(t, captured[0], attempts[0].Hash, "eth_getTransactionReceipt")
	cltest.BatchElemMustMatchParams(t, captured[1], attempts[1].Hash, "eth_getTransactionReceipt")
	cltest.BatchElemMustMatchParams(t, captured[2], attempts[2].Hash, "eth_getTransactionReceipt")
	cltest.BatchElemMustMatchParams(t, captured[3], attempts[3].Hash, "eth_getTransactionReceipt")
}

func TestEthConfirmer_CheckForReceipts_should_not_check_for_likely_unconfirmed(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	gconfig, config := newTestChainScopedConfig(t)
	txStore := cltest.NewTestTxStore(t, db)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)

	ec := newEthConfirmer(t, txStore, ethClient, gconfig, config, ethKeyStore, nil)
	ctx := tests.Context(t)
	latestFinalizedBlockNum := int64(0)

	etx := cltest.MustInsertUnconfirmedEthTx(t, txStore, 1, fromAddress)
	for i := 0; i < 4; i++ {
		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, int64(100-i))
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt))
	}

	// latest nonce is lower that all attempts' nonces
	ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(0), nil)

	require.NoError(t, ec.CheckForReceipts(ctx, 42, latestFinalizedBlockNum))
}

func TestEthConfirmer_CheckForReceipts_confirmed_missing_receipt_scoped_to_key(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress1_1 := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	_, fromAddress1_2 := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	_, fromAddress2_1 := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(20), nil)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, ethKeyStore, nil)
	ctx := tests.Context(t)
	latestFinalizedBlockNum := int64(0)

	// STATE
	// key 1, tx with nonce 0 is unconfirmed
	// key 1, tx with nonce 1 is unconfirmed
	// key 2, tx with nonce 9 is unconfirmed and gets a receipt in block 10
	etx1_0 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 0, fromAddress1_1)
	etx1_1 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 1, fromAddress1_1)
	etx2_9 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 3, fromAddress1_2)
	// there also happens to be a confirmed tx with a higher nonce from a different chain in the DB
	etx_other_chain := cltest.MustInsertUnconfirmedEthTx(t, txStore, 8, fromAddress2_1)
	pgtest.MustExec(t, db, `UPDATE evm.txes SET state='confirmed' WHERE id = $1`, etx_other_chain.ID)

	attempt2_9 := newBroadcastLegacyEthTxAttempt(t, etx2_9.ID, int64(1))
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt2_9))
	txmReceipt2_9 := newTxReceipt(attempt2_9.Hash, 10, 1)

	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt2_9.Hash, "eth_getTransactionReceipt")
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt2_9
	}).Once()

	require.NoError(t, ec.CheckForReceipts(ctx, 10, latestFinalizedBlockNum))

	mustTxBeInState(t, txStore, etx1_0, txmgrcommon.TxUnconfirmed)
	mustTxBeInState(t, txStore, etx1_1, txmgrcommon.TxUnconfirmed)
	mustTxBeInState(t, txStore, etx2_9, txmgrcommon.TxConfirmed)

	// Now etx1_1 gets a receipt in block 11, which should mark etx1_0 as confirmed_missing_receipt
	attempt1_1 := newBroadcastLegacyEthTxAttempt(t, etx1_1.ID, int64(2))
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt1_1))
	txmReceipt1_1 := newTxReceipt(attempt1_1.Hash, 11, 1)

	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt1_1.Hash, "eth_getTransactionReceipt")
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt1_1
	}).Once()

	require.NoError(t, ec.CheckForReceipts(ctx, 11, latestFinalizedBlockNum))

	mustTxBeInState(t, txStore, etx1_0, txmgrcommon.TxConfirmedMissingReceipt)
	mustTxBeInState(t, txStore, etx1_1, txmgrcommon.TxConfirmed)
	mustTxBeInState(t, txStore, etx2_9, txmgrcommon.TxConfirmed)
}

func TestEthConfirmer_CheckForReceipts_confirmed_missing_receipt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {})
	txStore := cltest.NewTestTxStore(t, db)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, ethKeyStore, nil)
	ctx := tests.Context(t)
	latestFinalizedBlockNum := int64(0)

	// STATE
	// eth_txes with nonce 0 has two attempts (broadcast before block 21 and 41) the first of which will get a receipt
	// eth_txes with nonce 1 has two attempts (broadcast before block 21 and 41) neither of which will ever get a receipt
	// eth_txes with nonce 2 has an attempt (broadcast before block 41) that will not get a receipt on the first try but will get one later
	// eth_txes with nonce 3 has an attempt (broadcast before block 41) that has been confirmed in block 42
	// All other attempts were broadcast before block 41
	b := int64(21)

	etx0 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 0, fromAddress)
	attempt0_1 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(1))
	attempt0_2 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(2))
	attempt0_2.BroadcastBeforeBlockNum = &b
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt0_1))
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt0_2))

	etx1 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 1, fromAddress)
	attempt1_1 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(1))
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(2))
	attempt1_2.BroadcastBeforeBlockNum = &b
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt1_1))
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt1_2))

	etx2 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 2, fromAddress)
	attempt2_1 := newBroadcastLegacyEthTxAttempt(t, etx2.ID, int64(1))
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt2_1))

	etx3 := cltest.MustInsertUnconfirmedEthTx(t, txStore, 3, fromAddress)
	attempt3_1 := newBroadcastLegacyEthTxAttempt(t, etx3.ID, int64(1))
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt3_1))

	pgtest.MustExec(t, db, `UPDATE evm.tx_attempts SET broadcast_before_block_num = 41 WHERE broadcast_before_block_num IS NULL`)

	t.Run("marks buried eth_txes as 'confirmed_missing_receipt'", func(t *testing.T) {
		txmReceipt0 := evmtypes.Receipt{
			TxHash:           attempt0_2.Hash,
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
			Status:           uint64(1),
		}
		txmReceipt3 := evmtypes.Receipt{
			TxHash:           attempt3_1.Hash,
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
			Status:           uint64(1),
		}
		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(4), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 6 &&
				cltest.BatchElemMatchesParams(b[0], attempt0_2.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[1], attempt0_1.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[2], attempt1_2.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[3], attempt1_1.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[4], attempt2_1.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[5], attempt3_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// First transaction confirmed
			*(elems[0].Result.(*evmtypes.Receipt)) = txmReceipt0
			elems[1].Result = &evmtypes.Receipt{}
			// Second transaction stil unconfirmed
			elems[2].Result = &evmtypes.Receipt{}
			elems[3].Result = &evmtypes.Receipt{}
			// Third transaction still unconfirmed
			elems[4].Result = &evmtypes.Receipt{}
			// Fourth transaction is confirmed
			*(elems[5].Result.(*evmtypes.Receipt)) = txmReceipt3
		}).Once()

		// PERFORM
		// Block num of 43 is one higher than the receipt (as would generally be expected)
		require.NoError(t, ec.CheckForReceipts(ctx, 43, latestFinalizedBlockNum))

		// Expected state is that the "top" eth_tx is now confirmed, with the
		// two below it "confirmed_missing_receipt" and the "bottom" eth_tx also confirmed
		var err error
		etx3, err = txStore.FindTxWithAttempts(ctx, etx3.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx3.State)

		ethReceipt := etx3.TxAttempts[0].Receipts[0]
		require.Equal(t, txmReceipt3.BlockHash, ethReceipt.GetBlockHash())

		etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx2.State)
		etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx1.State)

		etx0, err = txStore.FindTxWithAttempts(ctx, etx0.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx0.State)

		require.Len(t, etx0.TxAttempts, 2)
		require.Len(t, etx0.TxAttempts[0].Receipts, 1)
		ethReceipt = etx0.TxAttempts[0].Receipts[0]
		require.Equal(t, txmReceipt0.BlockHash, ethReceipt.GetBlockHash())
	})

	// STATE
	// eth_txes with nonce 0 is confirmed
	// eth_txes with nonce 1 is confirmed_missing_receipt
	// eth_txes with nonce 2 is confirmed_missing_receipt
	// eth_txes with nonce 3 is confirmed

	t.Run("marks eth_txes with state 'confirmed_missing_receipt' as 'confirmed' if a receipt finally shows up", func(t *testing.T) {
		txmReceipt := evmtypes.Receipt{
			TxHash:           attempt2_1.Hash,
			BlockHash:        testutils.NewHash(),
			BlockNumber:      big.NewInt(43),
			TransactionIndex: uint(1),
			Status:           uint64(1),
		}
		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 3 &&
				cltest.BatchElemMatchesParams(b[0], attempt1_2.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[1], attempt1_1.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[2], attempt2_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// First transaction still unconfirmed
			elems[0].Result = &evmtypes.Receipt{}
			elems[1].Result = &evmtypes.Receipt{}
			// Second transaction confirmed
			*(elems[2].Result.(*evmtypes.Receipt)) = txmReceipt
		}).Once()

		// PERFORM
		// Block num of 44 is one higher than the receipt (as would generally be expected)
		require.NoError(t, ec.CheckForReceipts(ctx, 44, latestFinalizedBlockNum))

		// Expected state is that the "top" two eth_txes are now confirmed, with the
		// one below it still "confirmed_missing_receipt" and the bottom one remains confirmed
		var err error
		etx3, err = txStore.FindTxWithAttempts(ctx, etx3.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx3.State)
		etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx2.State)

		ethReceipt := etx2.TxAttempts[0].Receipts[0]
		require.Equal(t, txmReceipt.BlockHash, ethReceipt.GetBlockHash())

		etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx1.State)
		etx0, err = txStore.FindTxWithAttempts(ctx, etx0.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx0.State)
	})

	// STATE
	// eth_txes with nonce 0 is confirmed
	// eth_txes with nonce 1 is confirmed_missing_receipt
	// eth_txes with nonce 2 is confirmed
	// eth_txes with nonce 3 is confirmed

	t.Run("continues to leave eth_txes with state 'confirmed_missing_receipt' unchanged if at least one attempt is above LatestFinalizedBlockNum", func(t *testing.T) {
		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				cltest.BatchElemMatchesParams(b[0], attempt1_2.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[1], attempt1_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// Both attempts still unconfirmed
			elems[0].Result = &evmtypes.Receipt{}
			elems[1].Result = &evmtypes.Receipt{}
		}).Once()

		latestFinalizedBlockNum = 30

		// PERFORM
		// Block num of 80 puts the first attempt (21) below threshold but second attempt (41) still above
		require.NoError(t, ec.CheckForReceipts(ctx, 80, latestFinalizedBlockNum))

		// Expected state is that the "top" two eth_txes are now confirmed, with the
		// one below it still "confirmed_missing_receipt" and the bottom one remains confirmed
		var err error
		etx3, err = txStore.FindTxWithAttempts(ctx, etx3.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx3.State)
		etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx2.State)
		etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx1.State)
		etx0, err = txStore.FindTxWithAttempts(ctx, etx0.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx0.State)
	})

	// STATE
	// eth_txes with nonce 0 is confirmed
	// eth_txes with nonce 1 is confirmed_missing_receipt
	// eth_txes with nonce 2 is confirmed
	// eth_txes with nonce 3 is confirmed

	t.Run("marks eth_Txes with state 'confirmed_missing_receipt' as 'errored' if a receipt fails to show up and all attempts are buried deeper than LatestFinalizedBlockNum", func(t *testing.T) {
		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(10), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				cltest.BatchElemMatchesParams(b[0], attempt1_2.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[1], attempt1_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// Both attempts still unconfirmed
			elems[0].Result = &evmtypes.Receipt{}
			elems[1].Result = &evmtypes.Receipt{}
		}).Once()

		latestFinalizedBlockNum = 50

		// PERFORM
		// Block num of 100 puts the first attempt (21) and second attempt (41) below threshold
		require.NoError(t, ec.CheckForReceipts(ctx, 100, latestFinalizedBlockNum))

		// Expected state is that the "top" two eth_txes are now confirmed, with the
		// one below it marked as "fatal_error" and the bottom one remains confirmed
		var err error
		etx3, err = txStore.FindTxWithAttempts(ctx, etx3.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx3.State)
		etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx2.State)
		etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxFatalError, etx1.State)
		etx0, err = txStore.FindTxWithAttempts(ctx, etx0.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxConfirmed, etx0.State)
	})
}

func TestEthConfirmer_CheckConfirmedMissingReceipt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {})
	txStore := cltest.NewTestTxStore(t, db)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	ethClient.On("IsL2").Return(false).Maybe()

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, ethKeyStore, nil)
	ctx := tests.Context(t)

	// STATE
	// eth_txes with nonce 0 has two attempts, the later attempt with higher gas fees
	// eth_txes with nonce 1 has two attempts, the later attempt with higher gas fees
	// eth_txes with nonce 2 has one attempt
	originalBroadcastAt := time.Unix(1616509100, 0)
	etx0 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, txStore, 0, 1, originalBroadcastAt, fromAddress)
	attempt0_2 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(2))
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt0_2))
	etx1 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, txStore, 1, 1, originalBroadcastAt, fromAddress)
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(2))
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt1_2))
	etx2 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, txStore, 2, 1, originalBroadcastAt, fromAddress)
	attempt2_1 := etx2.TxAttempts[0]
	etx3 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, txStore, 3, 1, originalBroadcastAt, fromAddress)
	attempt3_1 := etx3.TxAttempts[0]

	ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 4 &&
			cltest.BatchElemMatchesParams(b[0], hexutil.Encode(attempt0_2.SignedRawTx), "eth_sendRawTransaction") &&
			cltest.BatchElemMatchesParams(b[1], hexutil.Encode(attempt1_2.SignedRawTx), "eth_sendRawTransaction") &&
			cltest.BatchElemMatchesParams(b[2], hexutil.Encode(attempt2_1.SignedRawTx), "eth_sendRawTransaction") &&
			cltest.BatchElemMatchesParams(b[3], hexutil.Encode(attempt3_1.SignedRawTx), "eth_sendRawTransaction")
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		// First transaction confirmed
		elems[0].Error = errors.New("nonce too low")
		elems[1].Error = errors.New("transaction underpriced")
		elems[2].Error = nil
		elems[3].Error = errors.New("transaction already finalized")
	}).Once()

	// PERFORM
	require.NoError(t, ec.CheckConfirmedMissingReceipt(ctx))

	// Expected state is that the "top" eth_tx is untouched but the other two
	// are marked as unconfirmed
	var err error
	etx0, err = txStore.FindTxWithAttempts(ctx, etx0.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx0.State)
	assert.Greater(t, etx0.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx1.State)
	assert.Greater(t, etx1.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx2.State)
	assert.Greater(t, etx2.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	etx3, err = txStore.FindTxWithAttempts(ctx, etx3.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx3.State)
	assert.Greater(t, etx3.BroadcastAt.Unix(), originalBroadcastAt.Unix())
}

func TestEthConfirmer_CheckConfirmedMissingReceipt_batchSendTransactions_fails(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {})
	txStore := cltest.NewTestTxStore(t, db)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	ethClient.On("IsL2").Return(false).Maybe()

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, ethKeyStore, nil)
	ctx := tests.Context(t)

	// STATE
	// eth_txes with nonce 0 has two attempts, the later attempt with higher gas fees
	// eth_txes with nonce 1 has two attempts, the later attempt with higher gas fees
	// eth_txes with nonce 2 has one attempt
	originalBroadcastAt := time.Unix(1616509100, 0)
	etx0 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, txStore, 0, 1, originalBroadcastAt, fromAddress)
	attempt0_2 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(2))
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt0_2))
	etx1 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, txStore, 1, 1, originalBroadcastAt, fromAddress)
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(2))
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt1_2))
	etx2 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, txStore, 2, 1, originalBroadcastAt, fromAddress)
	attempt2_1 := etx2.TxAttempts[0]

	ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 3 &&
			cltest.BatchElemMatchesParams(b[0], hexutil.Encode(attempt0_2.SignedRawTx), "eth_sendRawTransaction") &&
			cltest.BatchElemMatchesParams(b[1], hexutil.Encode(attempt1_2.SignedRawTx), "eth_sendRawTransaction") &&
			cltest.BatchElemMatchesParams(b[2], hexutil.Encode(attempt2_1.SignedRawTx), "eth_sendRawTransaction")
	})).Return(errors.New("Timed out")).Once()

	// PERFORM
	require.NoError(t, ec.CheckConfirmedMissingReceipt(ctx))

	// Expected state is that all txes are marked as unconfirmed, since the batch call had failed
	var err error
	etx0, err = txStore.FindTxWithAttempts(ctx, etx0.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx0.State)
	assert.Equal(t, etx0.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx1.State)
	assert.Equal(t, etx1.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx2.State)
	assert.Equal(t, etx2.BroadcastAt.Unix(), originalBroadcastAt.Unix())
}

func TestEthConfirmer_CheckConfirmedMissingReceipt_smallEvmRPCBatchSize_middleBatchSendTransactionFails(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].RPCDefaultBatchSize = ptr[uint32](1)
	})
	txStore := cltest.NewTestTxStore(t, db)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	ethClient.On("IsL2").Return(false).Maybe()

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, ethKeyStore, nil)
	ctx := tests.Context(t)

	// STATE
	// eth_txes with nonce 0 has two attempts, the later attempt with higher gas fees
	// eth_txes with nonce 1 has two attempts, the later attempt with higher gas fees
	// eth_txes with nonce 2 has one attempt
	originalBroadcastAt := time.Unix(1616509100, 0)
	etx0 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, txStore, 0, 1, originalBroadcastAt, fromAddress)
	attempt0_2 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(2))
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt0_2))
	etx1 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, txStore, 1, 1, originalBroadcastAt, fromAddress)
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(2))
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt1_2))
	etx2 := mustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, txStore, 2, 1, originalBroadcastAt, fromAddress)

	// Expect eth_sendRawTransaction in 3 batches. First batch will pass, 2nd will fail, 3rd never attempted.
	ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 1 &&
			cltest.BatchElemMatchesParams(b[0], hexutil.Encode(attempt0_2.SignedRawTx), "eth_sendRawTransaction")
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		// First transaction confirmed
		elems[0].Error = errors.New("nonce too low")
	}).Once()
	ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 1 &&
			cltest.BatchElemMatchesParams(b[0], hexutil.Encode(attempt1_2.SignedRawTx), "eth_sendRawTransaction")
	})).Return(errors.New("Timed out")).Once()

	// PERFORM
	require.NoError(t, ec.CheckConfirmedMissingReceipt(ctx))

	// Expected state is that all transactions since failed batch will be unconfirmed
	var err error
	etx0, err = txStore.FindTxWithAttempts(ctx, etx0.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx0.State)
	assert.Greater(t, etx0.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	etx1, err = txStore.FindTxWithAttempts(ctx, etx1.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx1.State)
	assert.Equal(t, etx1.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgrcommon.TxUnconfirmed, etx2.State)
	assert.Equal(t, etx2.BroadcastAt.Unix(), originalBroadcastAt.Unix())
}

func TestEthConfirmer_FindTxsRequiringRebroadcast(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	txStore := cltest.NewTestTxStore(t, db)
	ctx := tests.Context(t)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	evmFromAddress := fromAddress
	currentHead := int64(30)
	gasBumpThreshold := int64(10)
	tooNew := int64(21)
	onTheMoney := int64(20)
	oldEnough := int64(19)
	nonce := int64(0)

	mustInsertConfirmedEthTx(t, txStore, nonce, fromAddress)
	nonce++

	_, otherAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	evmOtherAddress := otherAddress

	lggr := logger.Test(t)

	ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, ethKeyStore, nil)

	t.Run("returns nothing when there are no transactions", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	mustInsertInProgressEthTx(t, txStore, nonce, fromAddress)
	nonce++

	t.Run("returns nothing when the transaction is in_progress", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	// This one has BroadcastBeforeBlockNum set as nil... which can happen, but it should be ignored
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++

	t.Run("ignores unconfirmed transactions with nil BroadcastBeforeBlockNum", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt1_1 := etx1.TxAttempts[0]
	var dbAttempt txmgr.DbEthTxAttempt
	dbAttempt.FromTxAttempt(&attempt1_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, tooNew, attempt1_1.ID))
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID)
	attempt1_2.BroadcastBeforeBlockNum = &onTheMoney
	attempt1_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(30000)}
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt1_2))

	t.Run("returns nothing when the transaction is unconfirmed with an attempt that is recent", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt2_1 := etx2.TxAttempts[0]
	dbAttempt = txmgr.DbEthTxAttempt{}
	dbAttempt.FromTxAttempt(&attempt2_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, tooNew, attempt2_1.ID))

	t.Run("returns nothing when the transaction has attempts that are too new", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	etxWithoutAttempts := cltest.NewEthTx(fromAddress)
	{
		n := evmtypes.Nonce(nonce)
		etxWithoutAttempts.Sequence = &n
	}
	now := time.Now()
	etxWithoutAttempts.BroadcastAt = &now
	etxWithoutAttempts.InitialBroadcastAt = &now
	etxWithoutAttempts.State = txmgrcommon.TxUnconfirmed
	require.NoError(t, txStore.InsertTx(ctx, &etxWithoutAttempts))
	nonce++

	t.Run("does nothing if the transaction is from a different address than the one given", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmOtherAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	t.Run("returns the transaction if it is unconfirmed and has no attempts (note that this is an invariant violation, but we handle it anyway)", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 1)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
	})

	t.Run("returns nothing for different chain id", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, big.NewInt(42))
		require.NoError(t, err)

		require.Len(t, etxs, 0)
	})

	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt3_1 := etx3.TxAttempts[0]
	dbAttempt = txmgr.DbEthTxAttempt{}
	dbAttempt.FromTxAttempt(&attempt3_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_1.ID))

	// NOTE: It should ignore qualifying eth_txes from a different address
	etxOther := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, otherAddress)
	attemptOther1 := etxOther.TxAttempts[0]
	dbAttempt = txmgr.DbEthTxAttempt{}
	dbAttempt.FromTxAttempt(&attemptOther1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attemptOther1.ID))

	t.Run("returns the transaction if it is unconfirmed with an attempt that is older than gasBumpThreshold blocks", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
	})

	t.Run("returns nothing if threshold is zero", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, 0, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 0)
	})

	t.Run("does not return more transactions for gas bumping than gasBumpThreshold", func(t *testing.T) {
		// Unconfirmed txes in DB are:
		// (unnamed) (nonce 2)
		// etx1 (nonce 3)
		// etx2 (nonce 4)
		// etxWithoutAttempts (nonce 5)
		// etx3 (nonce 6) - ready for bump
		// etx4 (nonce 7) - ready for bump
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 4, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 1) // returns etxWithoutAttempts only - eligible for gas bumping because it technically doesn't have any attempts within gasBumpThreshold blocks
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)

		etxs, err = ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 5, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2) // includes etxWithoutAttempts, etx3 and etx4
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)

		// Zero limit disables it
		etxs, err = ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 0, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2) // includes etxWithoutAttempts, etx3 and etx4
	})

	etx4 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt4_1 := etx4.TxAttempts[0]
	dbAttempt = txmgr.DbEthTxAttempt{}
	dbAttempt.FromTxAttempt(&attempt4_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt4_1.ID))

	t.Run("ignores pending transactions for another key", func(t *testing.T) {
		// Re-use etx3 nonce for another key, it should not affect the results for this key
		etxOther := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, (*etx3.Sequence).Int64(), otherAddress)
		aOther := etxOther.TxAttempts[0]
		dbAttempt = txmgr.DbEthTxAttempt{}
		dbAttempt.FromTxAttempt(&aOther)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, aOther.ID))

		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 6, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 3) // includes etxWithoutAttempts, etx3 and etx4
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
		assert.Equal(t, etx4.ID, etxs[2].ID)
	})

	attempt3_2 := newBroadcastLegacyEthTxAttempt(t, etx3.ID)
	attempt3_2.BroadcastBeforeBlockNum = &oldEnough
	attempt3_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(30000)}
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt3_2))

	t.Run("returns the transaction if it is unconfirmed with two attempts that are older than gasBumpThreshold blocks", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 3)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
		assert.Equal(t, etx4.ID, etxs[2].ID)
	})

	attempt3_3 := newBroadcastLegacyEthTxAttempt(t, etx3.ID)
	attempt3_3.BroadcastBeforeBlockNum = &tooNew
	attempt3_3.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(40000)}
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt3_3))

	t.Run("does not return the transaction if it has some older but one newer attempt", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, *etxWithoutAttempts.Sequence, *(etxs[0].Sequence))
		require.Equal(t, evmtypes.Nonce(5), *etxWithoutAttempts.Sequence)
		assert.Equal(t, etx4.ID, etxs[1].ID)
		assert.Equal(t, *etx4.Sequence, *(etxs[1].Sequence))
		require.Equal(t, evmtypes.Nonce(7), *etx4.Sequence)
	})

	attempt0_1 := newBroadcastLegacyEthTxAttempt(t, etxWithoutAttempts.ID)
	attempt0_1.State = txmgrtypes.TxAttemptInsufficientFunds
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt0_1))

	// This attempt has insufficient_eth, but there is also another attempt4_1
	// which is old enough, so this will be caught by both queries and should
	// not be duplicated
	attempt4_2 := cltest.NewLegacyEthTxAttempt(t, etx4.ID)
	attempt4_2.State = txmgrtypes.TxAttemptInsufficientFunds
	attempt4_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(40000)}
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt4_2))

	etx5 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, nonce, fromAddress)
	nonce++

	// This etx has one attempt that is too new, which would exclude it from
	// the gas bumping query, but it should still be caught by the insufficient
	// eth query
	etx6 := mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, nonce, fromAddress)
	attempt6_2 := newBroadcastLegacyEthTxAttempt(t, etx3.ID)
	attempt6_2.BroadcastBeforeBlockNum = &tooNew
	attempt6_2.TxFee = gas.EvmFee{Legacy: assets.NewWeiI(30001)}
	require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt6_2))

	t.Run("returns unique attempts requiring resubmission due to insufficient eth, ordered by nonce asc", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 0, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 4)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, *etxWithoutAttempts.Sequence, *(etxs[0].Sequence))
		assert.Equal(t, etx4.ID, etxs[1].ID)
		assert.Equal(t, *etx4.Sequence, *(etxs[1].Sequence))
		assert.Equal(t, etx5.ID, etxs[2].ID)
		assert.Equal(t, *etx5.Sequence, *(etxs[2].Sequence))
		assert.Equal(t, etx6.ID, etxs[3].ID)
		assert.Equal(t, *etx6.Sequence, *(etxs[3].Sequence))
	})

	t.Run("applies limit", func(t *testing.T) {
		etxs, err := ec.FindTxsRequiringRebroadcast(tests.Context(t), lggr, evmFromAddress, currentHead, gasBumpThreshold, 10, 2, &cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, *etxWithoutAttempts.Sequence, *(etxs[0].Sequence))
		assert.Equal(t, etx4.ID, etxs[1].ID)
		assert.Equal(t, *etx4.Sequence, *(etxs[1].Sequence))
	})
}

func TestEthConfirmer_RebroadcastWhereNecessary_WithConnectivityCheck(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	db := pgtest.NewSqlxDB(t)
	ethClient := testutils.NewEthClientMockWithDefaultChain(t)

	t.Run("should retry previous attempt if connectivity check failed for legacy transactions", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(false)
			c.EVM[0].GasEstimator.BlockHistory.BlockHistorySize = ptr[uint16](2)
			c.EVM[0].GasEstimator.BlockHistory.CheckInclusionBlocks = ptr[uint16](4)
		})
		ccfg := evmtest.NewChainScopedConfig(t, cfg)

		ctx := tests.Context(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		kst := ksmocks.NewEth(t)

		estimator := gasmocks.NewEvmEstimator(t)
		newEst := func(logger.Logger) gas.EvmEstimator { return estimator }
		estimator.On("BumpLegacyGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, uint64(0), pkgerrors.Wrapf(commonfee.ErrConnectivity, "transaction..."))
		ge := ccfg.EVM().GasEstimator()
		feeEstimator := gas.NewEvmFeeEstimator(lggr, newEst, ge.EIP1559DynamicFees(), ge, ethClient)
		txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, kst, feeEstimator)
		addresses := []gethCommon.Address{fromAddress}
		kst.On("EnabledAddressesForChain", mock.Anything, &cltest.FixtureChainID).Return(addresses, nil).Maybe()
		stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, "", assets.NewWei(assets.NewEth(100).ToInt()), ccfg.EVM().Transactions().AutoPurge(), feeEstimator, txStore, ethClient)
		ht := headtracker.NewSimulatedHeadTracker(ethClient, true, 0)
		// Create confirmer with necessary state
		ec := txmgr.NewEvmConfirmer(txStore, txmgr.NewEvmTxmClient(ethClient, nil), ccfg.EVM(), txmgr.NewEvmTxmFeeConfig(ccfg.EVM().GasEstimator()), ccfg.EVM().Transactions(), cfg.Database(), kst, txBuilder, lggr, stuckTxDetector, ht)
		servicetest.Run(t, ec)
		currentHead := int64(30)
		oldEnough := int64(15)
		nonce := int64(0)
		originalBroadcastAt := time.Unix(1616509100, 0)

		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress, originalBroadcastAt)
		attempt1 := etx.TxAttempts[0]
		var dbAttempt txmgr.DbEthTxAttempt
		dbAttempt.FromTxAttempt(&attempt1)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1.ID))

		// Send transaction and assume success.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(commonclient.Successful, nil).Once()

		err := ec.RebroadcastWhereNecessary(tests.Context(t), currentHead)
		require.NoError(t, err)

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 1)
	})

	t.Run("should retry previous attempt if connectivity check failed for dynamic transactions", func(t *testing.T) {
		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
			c.EVM[0].GasEstimator.BlockHistory.BlockHistorySize = ptr[uint16](2)
			c.EVM[0].GasEstimator.BlockHistory.CheckInclusionBlocks = ptr[uint16](4)
		})
		ccfg := evmtest.NewChainScopedConfig(t, cfg)

		ctx := tests.Context(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
		kst := ksmocks.NewEth(t)

		estimator := gasmocks.NewEvmEstimator(t)
		estimator.On("BumpDynamicFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(gas.DynamicFee{}, pkgerrors.Wrapf(commonfee.ErrConnectivity, "transaction..."))
		newEst := func(logger.Logger) gas.EvmEstimator { return estimator }
		// Create confirmer with necessary state
		ge := ccfg.EVM().GasEstimator()
		feeEstimator := gas.NewEvmFeeEstimator(lggr, newEst, ge.EIP1559DynamicFees(), ge, ethClient)
		txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, kst, feeEstimator)
		addresses := []gethCommon.Address{fromAddress}
		kst.On("EnabledAddressesForChain", mock.Anything, &cltest.FixtureChainID).Return(addresses, nil).Maybe()
		stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, "", assets.NewWei(assets.NewEth(100).ToInt()), ccfg.EVM().Transactions().AutoPurge(), feeEstimator, txStore, ethClient)
		ht := headtracker.NewSimulatedHeadTracker(ethClient, true, 0)
		ec := txmgr.NewEvmConfirmer(txStore, txmgr.NewEvmTxmClient(ethClient, nil), ccfg.EVM(), txmgr.NewEvmTxmFeeConfig(ccfg.EVM().GasEstimator()), ccfg.EVM().Transactions(), cfg.Database(), kst, txBuilder, lggr, stuckTxDetector, ht)
		servicetest.Run(t, ec)
		currentHead := int64(30)
		oldEnough := int64(15)
		nonce := int64(0)
		originalBroadcastAt := time.Unix(1616509100, 0)

		etx := mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, txStore, nonce, fromAddress, originalBroadcastAt)
		attempt1 := etx.TxAttempts[0]
		var dbAttempt txmgr.DbEthTxAttempt
		dbAttempt.FromTxAttempt(&attempt1)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1.ID))

		// Send transaction and assume success.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(commonclient.Successful, nil).Once()

		err := ec.RebroadcastWhereNecessary(tests.Context(t), currentHead)
		require.NoError(t, err)

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.TxAttempts, 1)
	})
}

func TestEthConfirmer_RebroadcastWhereNecessary_MaxFeeScenario(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.PriceMax = assets.GWei(500)
	})
	txStore := cltest.NewTestTxStore(t, db)
	ctx := tests.Context(t)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	_, _ = cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	kst := ksmocks.NewEth(t)
	addresses := []gethCommon.Address{fromAddress}
	kst.On("EnabledAddressesForChain", mock.Anything, &cltest.FixtureChainID).Return(addresses, nil).Maybe()
	// Use a mock keystore for this test
	ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, kst, nil)
	currentHead := int64(30)
	oldEnough := int64(19)
	nonce := int64(0)

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress, originalBroadcastAt)
	attempt1_1 := etx.TxAttempts[0]
	var dbAttempt txmgr.DbEthTxAttempt
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_1.ID))

	t.Run("treats an exceeds max fee attempt as a success", func(t *testing.T) {
		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx", mock.Anything,
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if tx.Nonce() != uint64(*etx.Sequence) {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(evmcfg.EVM().ChainID()) == 0
			})).Return(&ethTx, nil).Once()

		// Once for the bumped attempt which exceeds limit
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence) && tx.GasPrice().Int64() == int64(20000000000)
		}), fromAddress).Return(commonclient.ExceedsMaxFee, errors.New("tx fee (1.10 ether) exceeds the configured cap (1.00 ether)")).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		// Check that the attempt is saved
		require.Len(t, etx.TxAttempts, 2)

		// broadcast_at did change
		require.Greater(t, etx.BroadcastAt.Unix(), originalBroadcastAt.Unix())
		require.Equal(t, etx.InitialBroadcastAt.Unix(), originalBroadcastAt.Unix())
	})
}

func TestEthConfirmer_RebroadcastWhereNecessary(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.PriceMax = assets.GWei(500)
	})
	txStore := cltest.NewTestTxStore(t, db)
	ctx := tests.Context(t)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	_, _ = cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	kst := ksmocks.NewEth(t)
	addresses := []gethCommon.Address{fromAddress}
	kst.On("EnabledAddressesForChain", mock.Anything, &cltest.FixtureChainID).Return(addresses, nil).Maybe()
	// Use a mock keystore for this test
	ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, kst, nil)
	currentHead := int64(30)
	oldEnough := int64(19)
	nonce := int64(0)

	t.Run("does nothing if no transactions require bumping", func(t *testing.T) {
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
	})

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress, originalBroadcastAt)
	nonce++
	attempt1_1 := etx.TxAttempts[0]
	var dbAttempt txmgr.DbEthTxAttempt
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_1.ID))

	t.Run("re-sends previous transaction on keystore error", func(t *testing.T) {
		// simulate bumped transaction that is somehow impossible to sign
		kst.On("SignTx", mock.Anything, fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				return tx.Nonce() == uint64(*etx.Sequence)
			}),
			mock.Anything).Return(nil, errors.New("signing error")).Once()

		// Do the thing
		err := ec.RebroadcastWhereNecessary(tests.Context(t), currentHead)
		require.Error(t, err)
		require.Contains(t, err.Error(), "signing error")

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)

		require.Len(t, etx.TxAttempts, 1)
	})

	t.Run("does nothing and continues on fatal error", func(t *testing.T) {
		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx", mock.Anything,
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if tx.Nonce() != uint64(*etx.Sequence) {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(evmcfg.EVM().ChainID()) == 0
			})).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence)
		}), fromAddress).Return(commonclient.Fatal, errors.New("exceeds block gas limit")).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 1)
	})

	var attempt1_2 txmgr.TxAttempt
	ethClient = testutils.NewEthClientMockWithDefaultChain(t)
	ec.XXXTestSetClient(txmgr.NewEvmTxmClient(ethClient, nil))

	t.Run("creates new attempt with higher gas price if transaction has an attempt older than threshold", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx", mock.Anything,
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(evmcfg.EVM().ChainID()) == 0
			})).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 2)
		require.Equal(t, attempt1_1.ID, etx.TxAttempts[1].ID)

		// Got the new attempt
		attempt1_2 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.Legacy.ToInt().Int64())
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt1_2.State)
	})

	t.Run("does nothing if there is an attempt without BroadcastBeforeBlockNum set", func(t *testing.T) {
		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 2)
	})
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_2.ID))
	var attempt1_3 txmgr.TxAttempt

	t.Run("creates new attempt with higher gas price if transaction is already in mempool (e.g. due to previous crash before we could save the new attempt)", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(25000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx", mock.Anything,
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != *etx.Sequence || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, fmt.Errorf("known transaction: %s", ethTx.Hash().Hex())).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 3)
		require.Equal(t, attempt1_1.ID, etx.TxAttempts[2].ID)
		require.Equal(t, attempt1_2.ID, etx.TxAttempts[1].ID)

		// Got the new attempt
		attempt1_3 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_3.TxFee.Legacy.ToInt().Int64())
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt1_3.State)
	})

	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_3.ID))
	var attempt1_4 txmgr.TxAttempt

	t.Run("saves new attempt even for transaction that has already been confirmed (nonce already used)", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(30000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		receipt := evmtypes.Receipt{BlockNumber: big.NewInt(40)}
		kst.On("SignTx", mock.Anything,
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != *etx.Sequence || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				receipt.TxHash = tx.Hash()
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.TransactionAlreadyKnown, errors.New("nonce too low")).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx.State)

		// Got the new attempt
		attempt1_4 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_4.TxFee.Legacy.ToInt().Int64())

		require.Len(t, etx.TxAttempts, 4)
		require.Equal(t, attempt1_1.ID, etx.TxAttempts[3].ID)
		require.Equal(t, attempt1_2.ID, etx.TxAttempts[2].ID)
		require.Equal(t, attempt1_3.ID, etx.TxAttempts[1].ID)
		require.Equal(t, attempt1_4.ID, etx.TxAttempts[0].ID)
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, etx.TxAttempts[0].State)
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, etx.TxAttempts[1].State)
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, etx.TxAttempts[2].State)
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, etx.TxAttempts[3].State)
	})

	// Mark original tx as confirmed, so we won't pick it up anymore
	pgtest.MustExec(t, db, `UPDATE evm.txes SET state = 'confirmed'`)

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt2_1 := etx2.TxAttempts[0]
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt2_1.ID))
	var attempt2_2 txmgr.TxAttempt

	t.Run("saves in-progress attempt on temporary error and returns error", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt2_1.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		n := *etx2.Sequence
		kst.On("SignTx", mock.Anything,
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != n || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == n && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Unknown, errors.New("some network error")).Once()

		// Do the thing
		err := ec.RebroadcastWhereNecessary(tests.Context(t), currentHead)
		require.Error(t, err)
		require.Contains(t, err.Error(), "some network error")

		etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx2.State)

		// Old attempt is untouched
		require.Len(t, etx2.TxAttempts, 2)
		require.Equal(t, attempt2_1.ID, etx2.TxAttempts[1].ID)
		attempt2_1 = etx2.TxAttempts[1]
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt2_1.State)
		assert.Equal(t, oldEnough, *attempt2_1.BroadcastBeforeBlockNum)

		// New in_progress attempt saved
		attempt2_2 = etx2.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptInProgress, attempt2_2.State)
		assert.Nil(t, attempt2_2.BroadcastBeforeBlockNum)

		// Do it again and move the attempt into "broadcast"
		n = *etx2.Sequence
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == n && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))

		// Attempt marked "broadcast"
		etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx2.State)

		// New in_progress attempt saved
		require.Len(t, etx2.TxAttempts, 2)
		require.Equal(t, attempt2_2.ID, etx2.TxAttempts[0].ID)
		attempt2_2 = etx2.TxAttempts[0]
		require.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt2_2.State)
		assert.Nil(t, attempt2_2.BroadcastBeforeBlockNum)
	})

	// Set BroadcastBeforeBlockNum again so the next test will pick it up
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt2_2.ID))

	t.Run("assumes that 'nonce too low' error means confirmed_missing_receipt", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(25000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt2_1.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		n := *etx2.Sequence
		kst.On("SignTx", mock.Anything,
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != n || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == n && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.TransactionAlreadyKnown, errors.New("nonce too low")).Once()

		// Creates new attempt as normal if currentHead is not high enough
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx2, err = txStore.FindTxWithAttempts(ctx, etx2.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxConfirmedMissingReceipt, etx2.State)

		// One new attempt saved
		require.Len(t, etx2.TxAttempts, 3)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, etx2.TxAttempts[0].State)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, etx2.TxAttempts[1].State)
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, etx2.TxAttempts[2].State)
	})

	// Original tx is confirmed, so we won't pick it up anymore
	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt3_1 := etx3.TxAttempts[0]
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1, gas_price=$2 WHERE id=$3 RETURNING *`, oldEnough, assets.NewWeiI(35000000000), attempt3_1.ID))

	var attempt3_2 txmgr.TxAttempt

	t.Run("saves attempt anyway if replacement transaction is underpriced because the bumped gas price is insufficiently higher than the previous one", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(42000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt3_1.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx", mock.Anything,
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != *etx3.Sequence || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx3.Sequence && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, errors.New("replacement transaction underpriced")).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx3, err = txStore.FindTxWithAttempts(ctx, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx3.State)

		require.Len(t, etx3.TxAttempts, 2)
		require.Equal(t, attempt3_1.ID, etx3.TxAttempts[1].ID)
		attempt3_2 = etx3.TxAttempts[0]

		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt3_2.TxFee.Legacy.ToInt().Int64())
	})

	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_2.ID))
	var attempt3_3 txmgr.TxAttempt

	t.Run("handles case where transaction is already known somehow", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(50400000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt3_1.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx", mock.Anything,
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != *etx3.Sequence || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx3.Sequence && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, fmt.Errorf("known transaction: %s", ethTx.Hash().Hex())).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx3, err = txStore.FindTxWithAttempts(ctx, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx3.State)

		require.Len(t, etx3.TxAttempts, 3)
		attempt3_3 = etx3.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt3_3.TxFee.Legacy.ToInt().Int64())
	})

	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_3.ID))
	var attempt3_4 txmgr.TxAttempt

	t.Run("pretends it was accepted and continues the cycle if rejected for being temporarily underpriced", func(t *testing.T) {
		// This happens if parity is rejecting transactions that are not priced high enough to even get into the mempool at all
		// It should pretend it was accepted into the mempool and hand off to the next cycle to continue bumping gas as normal
		temporarilyUnderpricedError := "There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee."

		expectedBumpedGasPrice := big.NewInt(60480000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt3_2.TxFee.Legacy.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx", mock.Anything,
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != *etx3.Sequence || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx3.Sequence && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, errors.New(temporarilyUnderpricedError)).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx3, err = txStore.FindTxWithAttempts(ctx, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx3.State)

		require.Len(t, etx3.TxAttempts, 4)
		attempt3_4 = etx3.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt3_4.TxFee.Legacy.ToInt().Int64())
	})

	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_4.ID))

	t.Run("resubmits at the old price and does not create a new attempt if one of the bumped transactions would exceed EVM.GasEstimator.PriceMax", func(t *testing.T) {
		// Set price such that the next bump will exceed EVM.GasEstimator.PriceMax
		// Existing gas price is: 60480000000
		gasPrice := attempt3_4.TxFee.Legacy.ToInt()
		gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.PriceMax = assets.NewWeiI(60500000000)
		})
		newCfg := evmtest.NewChainScopedConfig(t, gcfg)
		ec2 := newEthConfirmer(t, txStore, ethClient, gcfg, newCfg, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx3.Sequence && gasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, errors.New("already known")).Once() // we already submitted at this price, now it's time to bump and submit again but since we simply resubmitted rather than increasing gas price, geth already knows about this tx

		// Do the thing
		require.NoError(t, ec2.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx3, err = txStore.FindTxWithAttempts(ctx, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx3.State)

		// No new tx attempts
		require.Len(t, etx3.TxAttempts, 4)
		attempt3_4 = etx3.TxAttempts[0]
		assert.Equal(t, gasPrice.Int64(), attempt3_4.TxFee.Legacy.ToInt().Int64())
	})

	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_4.ID))

	t.Run("resubmits at the old price and does not create a new attempt if the current price is exactly EVM.GasEstimator.PriceMax", func(t *testing.T) {
		// Set price such that the current price is already at EVM.GasEstimator.PriceMax
		// Existing gas price is: 60480000000
		gasPrice := attempt3_4.TxFee.Legacy.ToInt()
		gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.PriceMax = assets.NewWeiI(60480000000)
		})
		newCfg := evmtest.NewChainScopedConfig(t, gcfg)
		ec2 := newEthConfirmer(t, txStore, ethClient, gcfg, newCfg, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx3.Sequence && gasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, errors.New("already known")).Once() // we already submitted at this price, now it's time to bump and submit again but since we simply resubmitted rather than increasing gas price, geth already knows about this tx

		// Do the thing
		require.NoError(t, ec2.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx3, err = txStore.FindTxWithAttempts(ctx, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx3.State)

		// No new tx attempts
		require.Len(t, etx3.TxAttempts, 4)
		attempt3_4 = etx3.TxAttempts[0]
		assert.Equal(t, gasPrice.Int64(), attempt3_4.TxFee.Legacy.ToInt().Int64())
	})

	// The EIP-1559 etx and attempt
	etx4 := mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, txStore, nonce, fromAddress)
	attempt4_1 := etx4.TxAttempts[0]
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1, gas_tip_cap=$2, gas_fee_cap=$3 WHERE id=$4 RETURNING *`,
		oldEnough, assets.GWei(35), assets.GWei(100), attempt4_1.ID))
	var attempt4_2 txmgr.TxAttempt

	t.Run("EIP-1559: bumps using EIP-1559 rules when existing attempts are of type 0x2", func(t *testing.T) {
		ethTx := *types.NewTx(&types.DynamicFeeTx{})
		kst.On("SignTx", mock.Anything,
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != *etx4.Sequence {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		// This is the new, EIP-1559 attempt
		gasTipCap := assets.GWei(42)
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx4.Sequence && gasTipCap.ToInt().Cmp(tx.GasTipCap()) == 0
		}), fromAddress).Return(commonclient.Successful, nil).Once()
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx4, err = txStore.FindTxWithAttempts(ctx, etx4.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx4.State)

		// A new, bumped attempt
		require.Len(t, etx4.TxAttempts, 2)
		attempt4_2 = etx4.TxAttempts[0]
		assert.Nil(t, attempt4_2.TxFee.Legacy)
		assert.Equal(t, assets.GWei(42).String(), attempt4_2.TxFee.DynamicTipCap.String())
		assert.Equal(t, assets.GWei(120).String(), attempt4_2.TxFee.DynamicFeeCap.String())
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt1_2.State)
	})

	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1, gas_tip_cap=$2, gas_fee_cap=$3 WHERE id=$4 RETURNING *`,
		oldEnough, assets.GWei(999), assets.GWei(1000), attempt4_2.ID))

	t.Run("EIP-1559: resubmits at the old price and does not create a new attempt if one of the bumped EIP-1559 transactions would have its tip cap exceed EVM.GasEstimator.PriceMax", func(t *testing.T) {
		gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.PriceMax = assets.GWei(1000)
		})
		newCfg := evmtest.NewChainScopedConfig(t, gcfg)
		ec2 := newEthConfirmer(t, txStore, ethClient, gcfg, newCfg, ethKeyStore, nil)

		// Third attempt failed to bump, resubmits old one instead
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx4.Sequence && attempt4_2.Hash.String() == tx.Hash().String()
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		require.NoError(t, ec2.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx4, err = txStore.FindTxWithAttempts(ctx, etx4.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx4.State)

		// No new tx attempts
		require.Len(t, etx4.TxAttempts, 2)
		assert.Equal(t, assets.GWei(999).Int64(), etx4.TxAttempts[0].TxFee.DynamicTipCap.ToInt().Int64())
		assert.Equal(t, assets.GWei(1000).Int64(), etx4.TxAttempts[0].TxFee.DynamicFeeCap.ToInt().Int64())
	})

	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1, gas_tip_cap=$2, gas_fee_cap=$3 WHERE id=$4 RETURNING *`,
		oldEnough, assets.GWei(45), assets.GWei(100), attempt4_2.ID))

	t.Run("EIP-1559: saves attempt anyway if replacement transaction is underpriced because the bumped gas price is insufficiently higher than the previous one", func(t *testing.T) {
		// NOTE: This test case was empirically impossible when I tried it on eth mainnet (any EIP1559 transaction with a higher tip cap is accepted even if it's only 1 wei more) but appears to be possible on Polygon/Matic, probably due to poor design that applies the 10% minimum to the overall value (base fee + tip cap)
		expectedBumpedTipCap := assets.GWei(54)
		require.Greater(t, expectedBumpedTipCap.Int64(), attempt4_2.TxFee.DynamicTipCap.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx", mock.Anything,
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if evmtypes.Nonce(tx.Nonce()) != *etx4.Sequence || expectedBumpedTipCap.ToInt().Cmp(tx.GasTipCap()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return evmtypes.Nonce(tx.Nonce()) == *etx4.Sequence && expectedBumpedTipCap.ToInt().Cmp(tx.GasTipCap()) == 0
		}), fromAddress).Return(commonclient.Successful, errors.New("replacement transaction underpriced")).Once()

		// Do it
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx4, err = txStore.FindTxWithAttempts(ctx, etx4.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx4.State)

		require.Len(t, etx4.TxAttempts, 3)
		require.Equal(t, attempt4_1.ID, etx4.TxAttempts[2].ID)
		require.Equal(t, attempt4_2.ID, etx4.TxAttempts[1].ID)
		attempt4_3 := etx4.TxAttempts[0]

		assert.Equal(t, expectedBumpedTipCap.Int64(), attempt4_3.TxFee.DynamicTipCap.ToInt().Int64())
	})
}

func TestEthConfirmer_RebroadcastWhereNecessary_TerminallyUnderpriced_ThenGoesThrough(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.PriceMax = assets.GWei(500)
	})
	txStore := cltest.NewTestTxStore(t, db)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	_, _ = cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	// Use a mock keystore for this test
	kst := ksmocks.NewEth(t)
	addresses := []gethCommon.Address{fromAddress}
	kst.On("EnabledAddressesForChain", mock.Anything, &cltest.FixtureChainID).Return(addresses, nil).Maybe()
	currentHead := int64(30)
	oldEnough := 5
	nonce := int64(0)

	t.Run("terminally underpriced transaction with in_progress attempt is retried with more gas", func(t *testing.T) {
		ethClient := testutils.NewEthClientMockWithDefaultChain(t)
		ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, kst, nil)

		originalBroadcastAt := time.Unix(1616509100, 0)
		etx := mustInsertUnconfirmedEthTxWithAttemptState(t, txStore, nonce, fromAddress, txmgrtypes.TxAttemptInProgress, originalBroadcastAt)
		require.Equal(t, originalBroadcastAt, *etx.BroadcastAt)
		nonce++
		attempt := etx.TxAttempts[0]
		signedTx, err := txmgr.GetGethSignedTx(attempt.SignedRawTx)
		require.NoError(t, err)

		// Fail the first time with terminally underpriced.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			commonclient.Underpriced, errors.New("Transaction gas price is too low. It does not satisfy your node's minimal gas price")).Once()
		// Succeed the second time after bumping gas.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			commonclient.Successful, nil).Once()
		kst.On("SignTx", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
			signedTx, nil,
		).Once()
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
	})

	t.Run("multiple gas bumps with existing broadcast attempts are retried with more gas until success in legacy mode", func(t *testing.T) {
		ethClient := testutils.NewEthClientMockWithDefaultChain(t)
		ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, kst, nil)

		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
		nonce++
		legacyAttempt := etx.TxAttempts[0]
		var dbAttempt txmgr.DbEthTxAttempt
		dbAttempt.FromTxAttempt(&legacyAttempt)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, legacyAttempt.ID))

		// Fail a few times with terminally underpriced
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			commonclient.Underpriced, errors.New("Transaction gas price is too low. It does not satisfy your node's minimal gas price")).Times(3)
		// Succeed the second time after bumping gas.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			commonclient.Successful, nil).Once()
		signedLegacyTx := new(types.Transaction)
		kst.On("SignTx", mock.Anything, mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Type() == 0x0 && tx.Nonce() == uint64(*etx.Sequence)
		}), mock.Anything).Return(
			signedLegacyTx, nil,
		).Run(func(args mock.Arguments) {
			unsignedLegacyTx := args.Get(2).(*types.Transaction)
			// Use the real keystore to do the actual signing
			thisSignedLegacyTx, err := ethKeyStore.SignTx(tests.Context(t), fromAddress, unsignedLegacyTx, testutils.FixtureChainID)
			require.NoError(t, err)
			*signedLegacyTx = *thisSignedLegacyTx
		}).Times(4) // 3 failures 1 success
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
	})

	t.Run("multiple gas bumps with existing broadcast attempts are retried with more gas until success in EIP-1559 mode", func(t *testing.T) {
		ethClient := testutils.NewEthClientMockWithDefaultChain(t)
		ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, kst, nil)

		etx := mustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, txStore, nonce, fromAddress)
		nonce++
		dxFeeAttempt := etx.TxAttempts[0]
		var dbAttempt txmgr.DbEthTxAttempt
		dbAttempt.FromTxAttempt(&dxFeeAttempt)
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, dxFeeAttempt.ID))

		// Fail a few times with terminally underpriced
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			commonclient.Underpriced, errors.New("transaction underpriced")).Times(3)
		// Succeed the second time after bumping gas.
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			commonclient.Successful, nil).Once()
		signedDxFeeTx := new(types.Transaction)
		kst.On("SignTx", mock.Anything, mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Type() == 0x2 && tx.Nonce() == uint64(*etx.Sequence)
		}), mock.Anything).Return(
			signedDxFeeTx, nil,
		).Run(func(args mock.Arguments) {
			unsignedDxFeeTx := args.Get(2).(*types.Transaction)
			// Use the real keystore to do the actual signing
			thisSignedDxFeeTx, err := ethKeyStore.SignTx(tests.Context(t), fromAddress, unsignedDxFeeTx, testutils.FixtureChainID)
			require.NoError(t, err)
			*signedDxFeeTx = *thisSignedDxFeeTx
		}).Times(4) // 3 failures 1 success
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
	})
}

func TestEthConfirmer_RebroadcastWhereNecessary_WhenOutOfEth(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ctx := tests.Context(t)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	_, err := ethKeyStore.EnabledKeysForChain(tests.Context(t), testutils.FixtureChainID)
	require.NoError(t, err)
	require.NoError(t, err)
	// keyStates, err := ethKeyStore.GetStatesForKeys(keys)
	// require.NoError(t, err)

	gconfig, config := newTestChainScopedConfig(t)
	currentHead := int64(30)
	oldEnough := int64(19)
	nonce := int64(0)

	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress)
	nonce++
	attempt1_1 := etx.TxAttempts[0]
	var dbAttempt txmgr.DbEthTxAttempt
	dbAttempt.FromTxAttempt(&attempt1_1)
	require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_1.ID))
	var attempt1_2 txmgr.TxAttempt

	insufficientEthError := errors.New("insufficient funds for gas * price + value")

	t.Run("saves attempt with state 'insufficient_eth' if eth node returns this error", func(t *testing.T) {
		ec := newEthConfirmer(t, txStore, ethClient, gconfig, config, ethKeyStore, nil)

		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.TxFee.Legacy.ToInt().Int64())

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.InsufficientFunds, insufficientEthError).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 2)
		require.Equal(t, attempt1_1.ID, etx.TxAttempts[1].ID)

		// Got the new attempt
		attempt1_2 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.Legacy.ToInt().Int64())
		assert.Equal(t, txmgrtypes.TxAttemptInsufficientFunds, attempt1_2.State)
		assert.Nil(t, attempt1_2.BroadcastBeforeBlockNum)
	})

	t.Run("does not bump gas when previous error was 'out of eth', instead resubmits existing transaction", func(t *testing.T) {
		ec := newEthConfirmer(t, txStore, ethClient, gconfig, config, ethKeyStore, nil)

		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.TxFee.Legacy.ToInt().Int64())

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.InsufficientFunds, insufficientEthError).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		// New attempt was NOT created
		require.Len(t, etx.TxAttempts, 2)

		// The attempt is still "out of eth"
		attempt1_2 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.Legacy.ToInt().Int64())
		assert.Equal(t, txmgrtypes.TxAttemptInsufficientFunds, attempt1_2.State)
	})

	t.Run("saves the attempt as broadcast after node wallet has been topped up with sufficient balance", func(t *testing.T) {
		ec := newEthConfirmer(t, txStore, ethClient, gconfig, config, ethKeyStore, nil)

		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.TxFee.Legacy.ToInt().Int64())

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))

		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		// New attempt was NOT created
		require.Len(t, etx.TxAttempts, 2)

		// Attempt is now 'broadcast'
		attempt1_2 = etx.TxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.TxFee.Legacy.ToInt().Int64())
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt1_2.State)
	})

	t.Run("resubmitting due to insufficient eth is not limited by EVM.GasEstimator.BumpTxDepth", func(t *testing.T) {
		depth := 2
		etxCount := 4

		cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.EVM[0].GasEstimator.BumpTxDepth = ptr(uint32(depth))
		})
		evmcfg := evmtest.NewChainScopedConfig(t, cfg)
		ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, ethKeyStore, nil)

		for i := 0; i < etxCount; i++ {
			n := nonce
			mustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, txStore, nonce, fromAddress)
			ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
				return tx.Nonce() == uint64(n)
			}), fromAddress).Return(commonclient.Successful, nil).Once()

			nonce++
		}

		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))

		var dbAttempts []txmgr.DbEthTxAttempt

		require.NoError(t, db.Select(&dbAttempts, "SELECT * FROM evm.tx_attempts WHERE state = 'insufficient_eth'"))
		require.Len(t, dbAttempts, 0)
	})
}

func TestEthConfirmer_RebroadcastWhereNecessary_TerminallyStuckError(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.PriceMax = assets.GWei(500)
	})
	txStore := cltest.NewTestTxStore(t, db)
	ctx := tests.Context(t)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	// Use a mock keystore for this test
	ec := newEthConfirmer(t, txStore, ethClient, cfg, evmcfg, ethKeyStore, nil)
	currentHead := int64(30)
	oldEnough := int64(19)
	nonce := int64(0)
	terminallyStuckError := "failed to add tx to the pool: not enough step counters to continue the execution"

	t.Run("terminally stuck transaction replaced with purge attempt", func(t *testing.T) {
		originalBroadcastAt := time.Unix(1616509100, 0)
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, nonce, fromAddress, originalBroadcastAt)
		nonce++
		attempt1_1 := etx.TxAttempts[0]
		var dbAttempt txmgr.DbEthTxAttempt
		require.NoError(t, db.Get(&dbAttempt, `UPDATE evm.tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_1.ID))

		// Return terminally stuck error on first rebroadcast
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence)
		}), fromAddress).Return(commonclient.TerminallyStuck, errors.New(terminallyStuckError)).Once()
		// Return successful for purge attempt
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Sequence)
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		// Start processing transactions for rebroadcast
		require.NoError(t, ec.RebroadcastWhereNecessary(tests.Context(t), currentHead))
		var err error
		etx, err = txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.TxAttempts, 2)
		purgeAttempt := etx.TxAttempts[0]
		require.True(t, purgeAttempt.IsPurgeAttempt)
	})
}

func TestEthConfirmer_EnsureConfirmedTransactionsInLongestChain(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ctx := tests.Context(t)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)

	gconfig, config := newTestChainScopedConfig(t)
	ec := newEthConfirmer(t, txStore, ethClient, gconfig, config, ethKeyStore, nil)

	latestFinalizedHead := evmtypes.Head{
		Number: 8,
		Hash:   testutils.NewHash(),
	}

	h8 := &evmtypes.Head{
		Number: 8,
		Hash:   testutils.NewHash(),
	}
	h9 := &evmtypes.Head{
		Hash:   testutils.NewHash(),
		Number: 9,
	}
	h9.Parent.Store(h8)
	head := &evmtypes.Head{
		Hash:   testutils.NewHash(),
		Number: 10,
	}
	head.Parent.Store(h9)
	t.Run("does nothing if there aren't any transactions", func(t *testing.T) {
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(tests.Context(t), head, latestFinalizedHead.BlockNumber()))
	})

	t.Run("does nothing to unconfirmed transactions", func(t *testing.T) {
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress)

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(tests.Context(t), head, latestFinalizedHead.BlockNumber()))

		etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
	})

	t.Run("does nothing to confirmed transactions with receipts within head height of the chain and included in the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 2, 1, fromAddress)
		mustInsertEthReceipt(t, txStore, head.Number, head.Hash, etx.TxAttempts[0].Hash)

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(tests.Context(t), head, latestFinalizedHead.BlockNumber()))

		etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxConfirmed, etx.State)
	})

	t.Run("does nothing to confirmed transactions that only have receipts older than the start of the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 3, 1, fromAddress)
		// Add receipt that is older than the lowest block of the chain
		mustInsertEthReceipt(t, txStore, h8.Number-1, testutils.NewHash(), etx.TxAttempts[0].Hash)

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(tests.Context(t), head, latestFinalizedHead.BlockNumber()))

		etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxConfirmed, etx.State)
	})

	t.Run("unconfirms and rebroadcasts transactions that have receipts within head height of the chain but not included in the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 4, 1, fromAddress)
		attempt := etx.TxAttempts[0]
		// Include one within head height but a different block hash
		mustInsertEthReceipt(t, txStore, head.Parent.Load().Number, testutils.NewHash(), attempt.Hash)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			atx, err := txmgr.GetGethSignedTx(attempt.SignedRawTx)
			require.NoError(t, err)
			// Keeps gas price and nonce the same
			return atx.GasPrice().Cmp(tx.GasPrice()) == 0 && atx.Nonce() == tx.Nonce()
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(tests.Context(t), head, latestFinalizedHead.BlockNumber()))

		etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
		require.Len(t, etx.TxAttempts, 1)
		attempt = etx.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
	})

	t.Run("unconfirms and rebroadcasts transactions that have receipts within head height of chain but not included in the chain even if a receipt exists older than the start of the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 5, 1, fromAddress)
		attempt := etx.TxAttempts[0]
		attemptHash := attempt.Hash
		// Add receipt that is older than the lowest block of the chain
		mustInsertEthReceipt(t, txStore, h8.Number-1, testutils.NewHash(), attemptHash)
		// Include one within head height but a different block hash
		mustInsertEthReceipt(t, txStore, head.Parent.Load().Number, testutils.NewHash(), attemptHash)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(
			commonclient.Successful, nil).Once()

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(tests.Context(t), head, latestFinalizedHead.BlockNumber()))

		etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
		require.Len(t, etx.TxAttempts, 1)
		attempt = etx.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
	})

	t.Run("if more than one attempt has a receipt (should not be possible but isn't prevented by database constraints) unconfirms and rebroadcasts only the attempt with the highest gas price", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 6, 1, fromAddress)
		require.Len(t, etx.TxAttempts, 1)
		// Sanity check to assert the included attempt has the lowest gas price
		require.Less(t, etx.TxAttempts[0].TxFee.Legacy.ToInt().Int64(), int64(30000))

		attempt2 := newBroadcastLegacyEthTxAttempt(t, etx.ID, 30000)
		attempt2.SignedRawTx = hexutil.MustDecode("0xf88c8301f3a98503b9aca000832ab98094f5fff180082d6017036b771ba883025c654bc93580a4daa6d556000000000000000000000000000000000000000000000000000000000000000026a0f25601065ee369b6470c0399a2334afcfbeb0b5c8f3d9a9042e448ed29b5bcbda05b676e00248b85faf4dd889f0e2dcf91eb867e23ac9eeb14a73f9e4c14972cdf")
		attempt3 := newBroadcastLegacyEthTxAttempt(t, etx.ID, 40000)
		attempt3.SignedRawTx = hexutil.MustDecode("0xf88c8301f3a88503b9aca0008316e36094151445852b0cfdf6a4cc81440f2af99176e8ad0880a4daa6d556000000000000000000000000000000000000000000000000000000000000000026a0dcb5a7ad52b96a866257134429f944c505820716567f070e64abb74899803855a04c13eff2a22c218e68da80111e1bb6dc665d3dea7104ab40ff8a0275a99f630d")
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt2))
		require.NoError(t, txStore.InsertTxAttempt(ctx, &attempt3))

		// Receipt is within head height but a different block hash
		mustInsertEthReceipt(t, txStore, head.Parent.Load().Number, testutils.NewHash(), attempt2.Hash)
		// Receipt is within head height but a different block hash
		mustInsertEthReceipt(t, txStore, head.Parent.Load().Number, testutils.NewHash(), attempt3.Hash)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			s, err := txmgr.GetGethSignedTx(attempt3.SignedRawTx)
			require.NoError(t, err)
			return tx.Hash() == s.Hash()
		}), fromAddress).Return(commonclient.Successful, nil).Once()

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(tests.Context(t), head, latestFinalizedHead.BlockNumber()))

		etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxUnconfirmed, etx.State)
		require.Len(t, etx.TxAttempts, 3)
		attempt1 := etx.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt1.State)
		attempt2 = etx.TxAttempts[1]
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt2.State)
		attempt3 = etx.TxAttempts[2]
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt3.State)
	})

	t.Run("if receipt has a block number that is in the future, does not mark for rebroadcast (the safe thing to do is simply wait until heads catches up)", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 7, 1, fromAddress)
		attempt := etx.TxAttempts[0]
		// Add receipt that is higher than head
		mustInsertEthReceipt(t, txStore, head.Number+1, testutils.NewHash(), attempt.Hash)

		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(tests.Context(t), head, latestFinalizedHead.BlockNumber()))

		etx, err := txStore.FindTxWithAttempts(ctx, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgrcommon.TxConfirmed, etx.State)
		require.Len(t, etx.TxAttempts, 1)
		attempt = etx.TxAttempts[0]
		assert.Equal(t, txmgrtypes.TxAttemptBroadcast, attempt.State)
		assert.Len(t, attempt.Receipts, 1)
	})
}

func TestEthConfirmer_ForceRebroadcast(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)

	gconfig, config := newTestChainScopedConfig(t)
	mustCreateUnstartedGeneratedTx(t, txStore, fromAddress, config.EVM().ChainID())
	mustInsertInProgressEthTx(t, txStore, 0, fromAddress)
	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress)
	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress)

	gasPriceWei := gas.EvmFee{Legacy: assets.GWei(52)}
	overrideGasLimit := uint64(20000)

	t.Run("rebroadcasts one eth_tx if it falls within in nonce range", func(t *testing.T) {
		ethClient := testutils.NewEthClientMockWithDefaultChain(t)
		ec := newEthConfirmer(t, txStore, ethClient, gconfig, config, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Sequence) &&
				tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() &&
				tx.Gas() == overrideGasLimit &&
				reflect.DeepEqual(tx.Data(), etx1.EncodedPayload) &&
				tx.To().String() == etx1.ToAddress.String()
		}), mock.Anything).Return(commonclient.Successful, nil).Once()

		require.NoError(t, ec.ForceRebroadcast(tests.Context(t), []evmtypes.Nonce{1}, gasPriceWei, fromAddress, overrideGasLimit))
	})

	t.Run("uses default gas limit if overrideGasLimit is 0", func(t *testing.T) {
		ethClient := testutils.NewEthClientMockWithDefaultChain(t)
		ec := newEthConfirmer(t, txStore, ethClient, gconfig, config, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Sequence) &&
				tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() &&
				tx.Gas() == etx1.FeeLimit &&
				reflect.DeepEqual(tx.Data(), etx1.EncodedPayload) &&
				tx.To().String() == etx1.ToAddress.String()
		}), mock.Anything).Return(commonclient.Successful, nil).Once()

		require.NoError(t, ec.ForceRebroadcast(tests.Context(t), []evmtypes.Nonce{(1)}, gasPriceWei, fromAddress, 0))
	})

	t.Run("rebroadcasts several eth_txes in nonce range", func(t *testing.T) {
		ethClient := testutils.NewEthClientMockWithDefaultChain(t)
		ec := newEthConfirmer(t, txStore, ethClient, gconfig, config, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Sequence) && tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() && tx.Gas() == overrideGasLimit
		}), mock.Anything).Return(commonclient.Successful, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx2.Sequence) && tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() && tx.Gas() == overrideGasLimit
		}), mock.Anything).Return(commonclient.Successful, nil).Once()

		require.NoError(t, ec.ForceRebroadcast(tests.Context(t), []evmtypes.Nonce{(1), (2)}, gasPriceWei, fromAddress, overrideGasLimit))
	})

	t.Run("broadcasts zero transactions if eth_tx doesn't exist for that nonce", func(t *testing.T) {
		ethClient := testutils.NewEthClientMockWithDefaultChain(t)
		ec := newEthConfirmer(t, txStore, ethClient, gconfig, config, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(1)
		}), mock.Anything).Return(commonclient.Successful, nil).Once()
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(2)
		}), mock.Anything).Return(commonclient.Successful, nil).Once()
		for i := 3; i <= 5; i++ {
			nonce := i
			ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
				return tx.Nonce() == uint64(nonce) &&
					tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() &&
					tx.Gas() == overrideGasLimit &&
					*tx.To() == fromAddress &&
					tx.Value().Cmp(big.NewInt(0)) == 0 &&
					len(tx.Data()) == 0
			}), mock.Anything).Return(commonclient.Successful, nil).Once()
		}
		nonces := []evmtypes.Nonce{(1), (2), (3), (4), (5)}

		require.NoError(t, ec.ForceRebroadcast(tests.Context(t), nonces, gasPriceWei, fromAddress, overrideGasLimit))
	})

	t.Run("zero transactions use default gas limit if override wasn't specified", func(t *testing.T) {
		ethClient := testutils.NewEthClientMockWithDefaultChain(t)
		ec := newEthConfirmer(t, txStore, ethClient, gconfig, config, ethKeyStore, nil)

		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(0) && tx.GasPrice().Int64() == gasPriceWei.Legacy.Int64() && tx.Gas() == config.EVM().GasEstimator().LimitDefault()
		}), mock.Anything).Return(commonclient.Successful, nil).Once()

		require.NoError(t, ec.ForceRebroadcast(tests.Context(t), []evmtypes.Nonce{(0)}, gasPriceWei, fromAddress, 0))
	})
}

func TestEthConfirmer_ResumePendingRuns(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	config := configtest.NewTestGeneralConfig(t)
	txStore := cltest.NewTestTxStore(t, db)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

	ethClient := testutils.NewEthClientMockWithDefaultChain(t)

	evmcfg := evmtest.NewChainScopedConfig(t, config)

	h8 := &evmtypes.Head{
		Number: 8,
		Hash:   testutils.NewHash(),
	}
	h9 := &evmtypes.Head{
		Hash:   testutils.NewHash(),
		Number: 9,
	}
	h9.Parent.Store(h8)
	head := evmtypes.Head{
		Hash:   testutils.NewHash(),
		Number: 10,
	}
	head.Parent.Store(h9)

	minConfirmations := int64(2)

	pgtest.MustExec(t, db, `SET CONSTRAINTS fk_pipeline_runs_pruning_key DEFERRED`)
	pgtest.MustExec(t, db, `SET CONSTRAINTS pipeline_runs_pipeline_spec_id_fkey DEFERRED`)

	t.Run("doesn't process task runs that are not suspended (possibly already previously resumed)", func(t *testing.T) {
		ec := newEthConfirmer(t, txStore, ethClient, config, evmcfg, ethKeyStore, func(context.Context, uuid.UUID, interface{}, error) error {
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

		err := ec.ResumePendingTaskRuns(tests.Context(t), &head)
		require.NoError(t, err)
	})

	t.Run("doesn't process task runs where the receipt is younger than minConfirmations", func(t *testing.T) {
		ec := newEthConfirmer(t, txStore, ethClient, config, evmcfg, ethKeyStore, func(context.Context, uuid.UUID, interface{}, error) error {
			t.Fatal("No value expected")
			return nil
		})

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 2, 1, fromAddress)
		mustInsertEthReceipt(t, txStore, head.Number, head.Hash, etx.TxAttempts[0].Hash)

		pgtest.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		err := ec.ResumePendingTaskRuns(tests.Context(t), &head)
		require.NoError(t, err)
	})

	t.Run("processes eth_txes with receipts older than minConfirmations", func(t *testing.T) {
		ch := make(chan interface{})
		nonce := evmtypes.Nonce(3)
		var err error
		ec := newEthConfirmer(t, txStore, ethClient, config, evmcfg, ethKeyStore, func(ctx context.Context, id uuid.UUID, value interface{}, thisErr error) error {
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
			err2 := ec.ResumePendingTaskRuns(tests.Context(t), &head)
			if !assert.NoError(t, err2) {
				return
			}
			// Retrieve Tx to check if callback completed flag was set to true
			updateTx, err3 := txStore.FindTxWithSequence(tests.Context(t), fromAddress, nonce)
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
		ec := newEthConfirmer(t, txStore, ethClient, config, evmcfg, ethKeyStore, func(ctx context.Context, id uuid.UUID, value interface{}, err error) error {
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
			err2 := ec.ResumePendingTaskRuns(tests.Context(t), &head)
			if !assert.NoError(t, err2) {
				return
			}
			// Retrieve Tx to check if callback completed flag was set to true
			updateTx, err3 := txStore.FindTxWithSequence(tests.Context(t), fromAddress, nonce)
			if assert.NoError(t, err3) {
				assert.Equal(t, true, updateTx.CallbackCompleted)
			}
		}()

		select {
		case data := <-ch:
			assert.Error(t, data.error)

			assert.EqualError(t, data.error, fmt.Sprintf("transaction %s reverted on-chain", etx.TxAttempts[0].Hash.String()))

			assert.Nil(t, data.value)

		case <-time.After(tests.WaitTimeout(t)):
			t.Fatal("no value received")
		}
	})

	t.Run("does not mark callback complete if callback fails", func(t *testing.T) {
		nonce := evmtypes.Nonce(5)
		ec := newEthConfirmer(t, txStore, ethClient, config, evmcfg, ethKeyStore, func(context.Context, uuid.UUID, interface{}, error) error {
			return errors.New("error")
		})

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, int64(nonce), 1, fromAddress)
		mustInsertEthReceipt(t, txStore, head.Number-minConfirmations, head.Hash, etx.TxAttempts[0].Hash)
		pgtest.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, min_confirmations = $2, signal_callback = TRUE WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		err := ec.ResumePendingTaskRuns(tests.Context(t), &head)
		require.Error(t, err)

		// Retrieve Tx to check if callback completed flag was left unchanged
		updateTx, err := txStore.FindTxWithSequence(tests.Context(t), fromAddress, nonce)
		require.NoError(t, err)
		require.Equal(t, false, updateTx.CallbackCompleted)
	})
}

func TestEthConfirmer_ProcessStuckTransactions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	txStore := cltest.NewTestTxStore(t, db)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)
	ethClient := testutils.NewEthClientMockWithDefaultChain(t)
	ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, fromAddress).Return(commonclient.Successful, nil).Once()
	lggr := logger.Test(t)
	feeEstimator := gasmocks.NewEvmFeeEstimator(t)

	// Return 10 gwei as market gas price
	marketGasPrice := tenGwei
	fee := gas.EvmFee{Legacy: marketGasPrice}
	bumpedLegacy := assets.GWei(30)
	bumpedFee := gas.EvmFee{Legacy: bumpedLegacy}
	feeEstimator.On("GetFee", mock.Anything, []byte{}, uint64(0), mock.Anything, mock.Anything, mock.Anything).Return(fee, uint64(0), nil)
	feeEstimator.On("BumpFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(bumpedFee, uint64(10_000), nil)
	autoPurgeThreshold := uint32(5)
	autoPurgeMinAttempts := uint32(3)
	limitDefault := uint64(100)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.LimitDefault = ptr(limitDefault)
		c.EVM[0].Transactions.AutoPurge.Enabled = ptr(true)
		c.EVM[0].Transactions.AutoPurge.Threshold = ptr(autoPurgeThreshold)
		c.EVM[0].Transactions.AutoPurge.MinAttempts = ptr(autoPurgeMinAttempts)
	})
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)
	ge := evmcfg.EVM().GasEstimator()
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, ethKeyStore, feeEstimator)
	stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, "", assets.NewWei(assets.NewEth(100).ToInt()), evmcfg.EVM().Transactions().AutoPurge(), feeEstimator, txStore, ethClient)
	ht := headtracker.NewSimulatedHeadTracker(ethClient, true, 0)
	ec := txmgr.NewEvmConfirmer(txStore, txmgr.NewEvmTxmClient(ethClient, nil), txmgr.NewEvmTxmConfig(evmcfg.EVM()), txmgr.NewEvmTxmFeeConfig(ge), evmcfg.EVM().Transactions(), cfg.Database(), ethKeyStore, txBuilder, lggr, stuckTxDetector, ht)
	fn := func(ctx context.Context, id uuid.UUID, result interface{}, err error) error {
		require.ErrorContains(t, err, client.TerminallyStuckMsg)
		return nil
	}
	ec.SetResumeCallback(fn)
	servicetest.Run(t, ec)

	ctx := tests.Context(t)
	blockNum := int64(100)

	t.Run("detects and processes stuck transactions", func(t *testing.T) {
		nonce := int64(0)
		// Create attempts so that the oldest broadcast attempt's block num is what meets the threshold check
		// Create autoPurgeMinAttempts number of attempts to ensure the broadcast attempt count check is not being triggered
		// Create attempts broadcasted autoPurgeThreshold block ago to ensure broadcast block num check is not being triggered
		tx := mustInsertUnconfirmedTxWithBroadcastAttempts(t, txStore, nonce, fromAddress, autoPurgeMinAttempts, blockNum-int64(autoPurgeThreshold), marketGasPrice.Add(oneGwei))
		// Update tx to signal callback once it is identified as terminally stuck
		pgtest.MustExec(t, db, `UPDATE evm.txes SET pipeline_task_run_id = $1, signal_callback = TRUE WHERE id = $2`, uuid.New(), tx.ID)
		head := evmtypes.Head{
			Hash:   testutils.NewHash(),
			Number: blockNum,
		}
		head.IsFinalized.Store(true)

		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(&head, nil).Once()
		ethClient.On("LatestFinalizedBlock", mock.Anything).Return(&head, nil).Once()
		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(0), nil).Once()
		ethClient.On("BatchCallContext", mock.Anything, mock.Anything).Return(nil).Once()

		// First call to ProcessHead should:
		// 1. Detect a stuck transaction
		// 2. Create a purge attempt for it
		// 3. Save the purge attempt to the DB
		// 4. Send the purge attempt
		err := ec.ProcessHead(ctx, &head)
		require.NoError(t, err)

		// Check if the purge attempt was saved to the DB properly
		dbTx, err := txStore.FindTxWithAttempts(ctx, tx.ID)
		require.NoError(t, err)
		require.NotNil(t, dbTx)
		latestAttempt := dbTx.TxAttempts[0]
		require.Equal(t, true, latestAttempt.IsPurgeAttempt)
		require.Equal(t, limitDefault, latestAttempt.ChainSpecificFeeLimit)
		require.Equal(t, bumpedFee.Legacy, latestAttempt.TxFee.Legacy)

		head = evmtypes.Head{
			Hash:   testutils.NewHash(),
			Number: blockNum + 1,
		}
		head.IsFinalized.Store(true)
		ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(&head, nil).Once()
		ethClient.On("LatestFinalizedBlock", mock.Anything).Return(&head, nil).Once()
		ethClient.On("SequenceAt", mock.Anything, mock.Anything, mock.Anything).Return(evmtypes.Nonce(1), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 4 && cltest.BatchElemMatchesParams(b[0], latestAttempt.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// First transaction confirmed
			*(elems[0].Result.(*evmtypes.Receipt)) = evmtypes.Receipt{
				TxHash:           latestAttempt.Hash,
				BlockHash:        testutils.NewHash(),
				BlockNumber:      big.NewInt(blockNum + 1),
				TransactionIndex: uint(1),
				Status:           uint64(1),
			}
		}).Once()

		// Second call to ProcessHead on next head should:
		// 1. Check for receipts for purged transaction
		// 2. When receipts are found for a purge attempt, the transaction is marked in the DB as fatal error with error message
		err = ec.ProcessHead(ctx, &head)
		require.NoError(t, err)
		dbTx, err = txStore.FindTxWithAttempts(ctx, tx.ID)
		require.NoError(t, err)
		require.NotNil(t, dbTx)
		require.Equal(t, txmgrcommon.TxFatalError, dbTx.State)
		require.Equal(t, client.TerminallyStuckMsg, dbTx.Error.String)
		require.Equal(t, true, dbTx.CallbackCompleted)
	})
}

func ptr[T any](t T) *T { return &t }

func newEthConfirmer(t testing.TB, txStore txmgr.EvmTxStore, ethClient client.Client, gconfig chainlink.GeneralConfig, config evmconfig.ChainScopedConfig, ks keystore.Eth, fn txmgrcommon.ResumeCallback) *txmgr.Confirmer {
	lggr := logger.Test(t)
	ge := config.EVM().GasEstimator()
	estimator := gas.NewEvmFeeEstimator(lggr, func(lggr logger.Logger) gas.EvmEstimator {
		return gas.NewFixedPriceEstimator(ge, nil, ge.BlockHistory(), lggr, nil)
	}, ge.EIP1559DynamicFees(), ge, ethClient)
	txBuilder := txmgr.NewEvmTxAttemptBuilder(*ethClient.ConfiguredChainID(), ge, ks, estimator)
	stuckTxDetector := txmgr.NewStuckTxDetector(lggr, testutils.FixtureChainID, "", assets.NewWei(assets.NewEth(100).ToInt()), config.EVM().Transactions().AutoPurge(), estimator, txStore, ethClient)
	ht := headtracker.NewSimulatedHeadTracker(ethClient, true, 0)
	ec := txmgr.NewEvmConfirmer(txStore, txmgr.NewEvmTxmClient(ethClient, nil), txmgr.NewEvmTxmConfig(config.EVM()), txmgr.NewEvmTxmFeeConfig(ge), config.EVM().Transactions(), gconfig.Database(), ks, txBuilder, lggr, stuckTxDetector, ht)
	ec.SetResumeCallback(fn)
	servicetest.Run(t, ec)
	return ec
}
