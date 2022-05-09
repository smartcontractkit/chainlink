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

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func newTestChainScopedConfig(t *testing.T) evmconfig.ChainScopedConfig {
	cfg := configtest.NewTestGeneralConfig(t)
	return evmtest.NewChainScopedConfig(t, cfg)
}

func mustInsertUnstartedEthTx(t *testing.T, borm txmgr.ORM, fromAddress gethCommon.Address) {
	etx := cltest.NewEthTx(t, fromAddress)
	etx.State = txmgr.EthTxUnstarted
	require.NoError(t, borm.InsertEthTx(&etx))
}

func newBroadcastLegacyEthTxAttempt(t *testing.T, etxID int64, gasPrice ...int64) txmgr.EthTxAttempt {
	attempt := cltest.NewLegacyEthTxAttempt(t, etxID)
	attempt.State = txmgr.EthTxAttemptBroadcast
	if len(gasPrice) > 0 {
		gp := gasPrice[0]
		attempt.GasPrice = utils.NewBig(big.NewInt(gp))
	}
	return attempt
}

func newInProgressLegacyEthTxAttempt(t *testing.T, etxID int64, gasPrice ...int64) txmgr.EthTxAttempt {
	attempt := cltest.NewLegacyEthTxAttempt(t, etxID)
	attempt.State = txmgr.EthTxAttemptInProgress
	if len(gasPrice) > 0 {
		gp := gasPrice[0]
		attempt.GasPrice = utils.NewBig(big.NewInt(gp))
	}
	return attempt
}

func mustInsertInProgressEthTx(t *testing.T, borm txmgr.ORM, nonce int64, fromAddress gethCommon.Address) txmgr.EthTx {
	etx := cltest.NewEthTx(t, fromAddress)
	etx.State = txmgr.EthTxInProgress
	etx.Nonce = &nonce
	require.NoError(t, borm.InsertEthTx(&etx))

	return etx
}

func mustInsertConfirmedEthTx(t *testing.T, borm txmgr.ORM, nonce int64, fromAddress gethCommon.Address) txmgr.EthTx {
	etx := cltest.NewEthTx(t, fromAddress)
	etx.State = txmgr.EthTxConfirmed
	etx.Nonce = &nonce
	now := time.Now()
	etx.BroadcastAt = &now
	etx.InitialBroadcastAt = &now
	require.NoError(t, borm.InsertEthTx(&etx))

	return etx
}

func TestEthConfirmer_SetBroadcastBeforeBlockNum(t *testing.T) {
	t.Parallel()
	db := pgtest.NewSqlxDB(t)
	cfg := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	state, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)
	ec := cltest.NewEthConfirmer(t, db, ethClient, cfg, ethKeyStore, []ethkey.State{state}, nil)
	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, fromAddress)

	headNum := int64(9000)
	var err error

	t.Run("saves block num to unconfirmed eth_tx_attempts without one", func(t *testing.T) {
		// Do the thing
		require.NoError(t, ec.SetBroadcastBeforeBlockNum(headNum))

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]

		assert.Equal(t, int64(9000), *attempt.BroadcastBeforeBlockNum)
	})

	t.Run("does not change eth_tx_attempts that already have BroadcastBeforeBlockNum set", func(t *testing.T) {
		n := int64(42)
		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, 2)
		attempt.BroadcastBeforeBlockNum = &n
		require.NoError(t, borm.InsertEthTxAttempt(&attempt))

		// Do the thing
		require.NoError(t, ec.SetBroadcastBeforeBlockNum(headNum))

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.EthTxAttempts, 2)
		attempt = etx.EthTxAttempts[0]

		assert.Equal(t, int64(42), *attempt.BroadcastBeforeBlockNum)
	})

	t.Run("only updates eth_tx_attempts for the current chain", func(t *testing.T) {
		etxThisChain := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 1, fromAddress, cfg.DefaultChainID())
		etxOtherChain := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, fromAddress, big.NewInt(1337))

		require.NoError(t, ec.SetBroadcastBeforeBlockNum(headNum))

		etxThisChain, err = borm.FindEthTxWithAttempts(etxThisChain.ID)
		require.NoError(t, err)
		require.Len(t, etxThisChain.EthTxAttempts, 1)
		attempt := etxThisChain.EthTxAttempts[0]

		assert.Equal(t, int64(9000), *attempt.BroadcastBeforeBlockNum)

		etxOtherChain, err = borm.FindEthTxWithAttempts(etxOtherChain.ID)
		require.NoError(t, err)
		require.Len(t, etxOtherChain.EthTxAttempts, 1)
		attempt = etxOtherChain.EthTxAttempts[0]

		assert.Nil(t, attempt.BroadcastBeforeBlockNum)
	})
}

func TestEthConfirmer_CheckForReceipts(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	config := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, config)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	state := cltest.MustGetStateForKey(t, ethKeyStore, key)

	ec := cltest.NewEthConfirmer(t, db, ethClient, config, ethKeyStore, []ethkey.State{state}, nil)

	nonce := int64(0)
	ctx := context.Background()
	blockNum := int64(0)

	t.Run("only finds eth_txes in unconfirmed state with at least one broadcast attempt", func(t *testing.T) {
		cltest.MustInsertFatalErrorEthTx(t, borm, fromAddress)
		mustInsertInProgressEthTx(t, borm, nonce, fromAddress)
		nonce++
		cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, nonce, 1, fromAddress)
		nonce++
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, nonce, fromAddress)
		nonce++
		mustInsertUnstartedEthTx(t, borm, fromAddress)

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))
		// No calls
		ethClient.AssertExpectations(t)
	})

	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, nonce, fromAddress)
	nonce++
	require.Len(t, etx1.EthTxAttempts, 1)
	attempt1_1 := etx1.EthTxAttempts[0]
	require.Len(t, attempt1_1.EthReceipts, 0)

	t.Run("fetches receipt for one unconfirmed eth_tx", func(t *testing.T) {
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		// Transaction not confirmed yet, receipt is nil
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt1_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &evmtypes.Receipt{}
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		var err error
		etx1, err = borm.FindEthTxWithAttempts(etx1.ID)
		assert.NoError(t, err)
		require.Len(t, etx1.EthTxAttempts, 1)
		attempt1_1 = etx1.EthTxAttempts[0]
		require.NoError(t, err)
		require.Len(t, attempt1_1.EthReceipts, 0)

		ethClient.AssertExpectations(t)
	})

	t.Run("saves nothing if returned receipt does not match the attempt", func(t *testing.T) {
		txmReceipt := evmtypes.Receipt{
			TxHash:           utils.NewHash(),
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}

		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		// First transaction confirmed
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt1_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &txmReceipt
		}).Once()

		// No error because it is merely logged
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		etx, err := borm.FindEthTxWithAttempts(etx1.ID)
		require.NoError(t, err)
		require.Len(t, etx.EthTxAttempts, 1)

		require.Len(t, etx.EthTxAttempts[0].EthReceipts, 0)
	})

	t.Run("saves nothing if query returns error", func(t *testing.T) {
		txmReceipt := evmtypes.Receipt{
			TxHash:           attempt1_1.Hash,
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}

		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		// First transaction confirmed
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt1_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &txmReceipt
			elems[0].Error = errors.New("foo")
		}).Once()

		// No error because it is merely logged
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		etx, err := borm.FindEthTxWithAttempts(etx1.ID)
		require.NoError(t, err)
		require.Len(t, etx.EthTxAttempts, 1)
		require.Len(t, etx.EthTxAttempts[0].EthReceipts, 0)
	})

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, nonce, fromAddress)
	nonce++
	require.Len(t, etx2.EthTxAttempts, 1)
	attempt2_1 := etx2.EthTxAttempts[0]
	require.Len(t, attempt2_1.EthReceipts, 0)

	t.Run("saves eth_receipt and marks eth_tx as confirmed when geth client returns valid receipt", func(t *testing.T) {
		txmReceipt := evmtypes.Receipt{
			TxHash:           attempt1_1.Hash,
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}

		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				cltest.BatchElemMatchesParams(b[0], attempt1_1.Hash, "eth_getTransactionReceipt") &&
				cltest.BatchElemMatchesParams(b[1], attempt2_1.Hash, "eth_getTransactionReceipt")

		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// First transaction confirmed
			elems[0].Result = &txmReceipt
			// Second transaction still unconfirmed
			elems[1].Result = &evmtypes.Receipt{}
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		// Check that the receipt was saved
		etx, err := borm.FindEthTxWithAttempts(etx1.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxConfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt1_1 = etx.EthTxAttempts[0]
		require.Len(t, attempt1_1.EthReceipts, 1)

		ethReceipt := attempt1_1.EthReceipts[0]

		assert.Equal(t, txmReceipt.TxHash, ethReceipt.TxHash)
		assert.Equal(t, txmReceipt.BlockHash, ethReceipt.BlockHash)
		assert.Equal(t, txmReceipt.BlockNumber.Int64(), ethReceipt.BlockNumber)
		assert.Equal(t, txmReceipt.TransactionIndex, ethReceipt.TransactionIndex)

		receiptJSON, err := json.Marshal(txmReceipt)
		require.NoError(t, err)

		assert.JSONEq(t, string(receiptJSON), string(ethReceipt.Receipt))

		ethClient.AssertExpectations(t)
	})

	t.Run("fetches and saves receipts for several attempts in gas price order", func(t *testing.T) {
		attempt2_2 := newBroadcastLegacyEthTxAttempt(t, etx2.ID)
		attempt2_2.GasPrice = utils.NewBig(big.NewInt(10))

		attempt2_3 := newBroadcastLegacyEthTxAttempt(t, etx2.ID)
		attempt2_3.GasPrice = utils.NewBig(big.NewInt(20))

		// Insert order deliberately reversed to test sorting by gas price
		require.NoError(t, borm.InsertEthTxAttempt(&attempt2_3))
		require.NoError(t, borm.InsertEthTxAttempt(&attempt2_2))

		txmReceipt := evmtypes.Receipt{
			TxHash:           attempt2_2.Hash,
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}

		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
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
			elems[1].Result = &txmReceipt
			// Cheapest attempt still unconfirmed
			elems[0].Result = &evmtypes.Receipt{}
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		ethClient.AssertExpectations(t)

		// Check that the state was updated
		etx, err := borm.FindEthTxWithAttempts(etx2.ID)
		require.NoError(t, err)

		require.Equal(t, txmgr.EthTxConfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 3)
	})

	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, nonce, fromAddress)
	attempt3_1 := etx3.EthTxAttempts[0]
	nonce++

	t.Run("ignores receipt missing BlockHash that comes from querying parity too early", func(t *testing.T) {
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		receipt := evmtypes.Receipt{
			TxHash: attempt3_1.Hash,
		}
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt3_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &receipt
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		// No receipt, but no error either
		etx, err := borm.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxUnconfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt3_1 = etx.EthTxAttempts[0]
		require.Len(t, attempt3_1.EthReceipts, 0)
	})

	t.Run("does not panic if receipt has BlockHash but is missing some other fields somehow", func(t *testing.T) {
		// NOTE: This should never happen, but we shouldn't panic regardless
		receipt := evmtypes.Receipt{
			TxHash:    attempt3_1.Hash,
			BlockHash: utils.NewHash(),
		}
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt3_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &receipt
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		// No receipt, but no error either
		etx, err := borm.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxUnconfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt3_1 = etx.EthTxAttempts[0]
		require.Len(t, attempt3_1.EthReceipts, 0)
	})

	t.Run("handles case where eth_receipt already exists somehow", func(t *testing.T) {
		ethReceipt := cltest.MustInsertEthReceipt(t, borm, 42, utils.NewHash(), attempt3_1.Hash)

		txmReceipt := evmtypes.Receipt{
			TxHash:           attempt3_1.Hash,
			BlockHash:        ethReceipt.BlockHash,
			BlockNumber:      big.NewInt(ethReceipt.BlockNumber),
			TransactionIndex: ethReceipt.TransactionIndex,
		}
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesParams(b[0], attempt3_1.Hash, "eth_getTransactionReceipt")
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &txmReceipt
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		// Check that the receipt was unchanged
		etx, err := borm.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxConfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt3_1 = etx.EthTxAttempts[0]
		require.Len(t, attempt3_1.EthReceipts, 1)

		ethReceipt = attempt3_1.EthReceipts[0]

		assert.Equal(t, txmReceipt.TxHash, ethReceipt.TxHash)
		assert.Equal(t, txmReceipt.BlockHash, ethReceipt.BlockHash)
		assert.Equal(t, txmReceipt.BlockNumber.Int64(), ethReceipt.BlockNumber)
		assert.Equal(t, txmReceipt.TransactionIndex, ethReceipt.TransactionIndex)

		ethClient.AssertExpectations(t)
	})

	etx4 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, nonce, fromAddress)
	attempt4_1 := etx4.EthTxAttempts[0]
	nonce++

	t.Run("on receipt fetch marks in_progress eth_tx_attempt as broadcast", func(t *testing.T) {
		attempt4_2 := newInProgressLegacyEthTxAttempt(t, etx4.ID)
		attempt4_2.GasPrice = utils.NewBig(big.NewInt(10))

		require.NoError(t, borm.InsertEthTxAttempt(&attempt4_2))

		txmReceipt := evmtypes.Receipt{
			TxHash:           attempt4_2.Hash,
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
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
			elems[0].Result = &txmReceipt
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		ethClient.AssertExpectations(t)

		// Check that the state was updated
		var err error
		etx4, err = borm.FindEthTxWithAttempts(etx4.ID)
		require.NoError(t, err)

		attempt4_1 = etx4.EthTxAttempts[1]
		attempt4_2 = etx4.EthTxAttempts[0]

		// And the attempts
		require.Equal(t, txmgr.EthTxAttemptBroadcast, attempt4_1.State)
		require.Nil(t, attempt4_1.BroadcastBeforeBlockNum)
		require.Equal(t, txmgr.EthTxAttemptBroadcast, attempt4_2.State)
		require.Equal(t, int64(42), *attempt4_2.BroadcastBeforeBlockNum)

		// Check receipts
		require.Len(t, attempt4_1.EthReceipts, 0)
		require.Len(t, attempt4_2.EthReceipts, 1)
	})
}

func TestEthConfirmer_CheckForReceipts_batching(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	cfg.Overrides.GlobalEvmRPCDefaultBatchSize = null.IntFrom(2)
	borm := cltest.NewTxmORM(t, db, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	state, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ec := cltest.NewEthConfirmer(t, db, ethClient, evmcfg, ethKeyStore, []ethkey.State{state}, nil)

	ctx := context.Background()

	etx := cltest.MustInsertUnconfirmedEthTx(t, borm, 0, fromAddress)
	var attempts []txmgr.EthTxAttempt

	// Total of 5 attempts should lead to 3 batched fetches (2, 2, 1)
	for i := 0; i < 5; i++ {
		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, int64(i+2))
		require.NoError(t, borm.InsertEthTxAttempt(&attempt))
		attempts = append(attempts, attempt)
	}

	ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)

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

	require.NoError(t, ec.CheckForReceipts(ctx, 42))
	ethClient.AssertExpectations(t)
}

func TestEthConfirmer_CheckForReceipts_only_likely_confirmed(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	cfg.Overrides.GlobalEvmRPCDefaultBatchSize = null.IntFrom(6)
	borm := cltest.NewTxmORM(t, db, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	state, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ec := cltest.NewEthConfirmer(t, db, ethClient, evmcfg, ethKeyStore, []ethkey.State{state}, nil)

	ctx := context.Background()

	var attempts []txmgr.EthTxAttempt
	// inserting in DESC nonce order to test DB ASC ordering
	etx2 := cltest.MustInsertUnconfirmedEthTx(t, borm, 1, fromAddress)
	for i := 0; i < 4; i++ {
		attempt := newBroadcastLegacyEthTxAttempt(t, etx2.ID, int64(100-i))
		require.NoError(t, borm.InsertEthTxAttempt(&attempt))
	}
	etx := cltest.MustInsertUnconfirmedEthTx(t, borm, 0, fromAddress)
	for i := 0; i < 4; i++ {
		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, int64(100-i))
		require.NoError(t, borm.InsertEthTxAttempt(&attempt))

		// only adding these because a batch for only those attempts should be sent
		attempts = append(attempts, attempt)
	}

	ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(0), nil)

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

	require.NoError(t, ec.CheckForReceipts(ctx, 42))

	cltest.BatchElemMustMatchParams(t, captured[0], attempts[0].Hash, "eth_getTransactionReceipt")
	cltest.BatchElemMustMatchParams(t, captured[1], attempts[1].Hash, "eth_getTransactionReceipt")
	cltest.BatchElemMustMatchParams(t, captured[2], attempts[2].Hash, "eth_getTransactionReceipt")
	cltest.BatchElemMustMatchParams(t, captured[3], attempts[3].Hash, "eth_getTransactionReceipt")

	ethClient.AssertExpectations(t)
}

func TestEthConfirmer_CheckForReceipts_should_not_check_for_likely_unconfirmed(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	config := newTestChainScopedConfig(t)
	borm := cltest.NewTxmORM(t, db, config)

	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	state, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	ec := cltest.NewEthConfirmer(t, db, ethClient, config, ethKeyStore, []ethkey.State{state}, nil)

	ctx := context.Background()

	etx := cltest.MustInsertUnconfirmedEthTx(t, borm, 1, fromAddress)
	for i := 0; i < 4; i++ {
		attempt := newBroadcastLegacyEthTxAttempt(t, etx.ID, int64(100-i))
		require.NoError(t, borm.InsertEthTxAttempt(&attempt))
	}

	// latest nonce is lower that all attempts' nonces
	ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(0), nil)

	require.NoError(t, ec.CheckForReceipts(ctx, 42))

	// no BatchCallContext calls
	ethClient.AssertExpectations(t)
}

func TestEthConfirmer_CheckForReceipts_confirmed_missing_receipt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	state, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	cfg.Overrides.GlobalEvmFinalityDepth = null.IntFrom(50)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ec := cltest.NewEthConfirmer(t, db, ethClient, evmcfg, ethKeyStore, []ethkey.State{state}, nil)

	ctx := context.Background()

	// STATE
	// eth_txes with nonce 0 has two attempts (broadcast before block 21 and 41) the first of which will get a receipt
	// eth_txes with nonce 1 has two attempts (broadcast before block 21 and 41) neither of which will ever get a receipt
	// eth_txes with nonce 2 has an attempt (broadcast before block 41) that will not get a receipt on the first try but will get one later
	// eth_txes with nonce 3 has an attempt (broadcast before block 41) that has been confirmed in block 42
	// All other attempts were broadcast before block 41
	b := int64(21)

	etx0 := cltest.MustInsertUnconfirmedEthTx(t, borm, 0, fromAddress)
	attempt0_1 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(1))
	attempt0_2 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(2))
	attempt0_2.BroadcastBeforeBlockNum = &b
	require.NoError(t, borm.InsertEthTxAttempt(&attempt0_1))
	require.NoError(t, borm.InsertEthTxAttempt(&attempt0_2))

	etx1 := cltest.MustInsertUnconfirmedEthTx(t, borm, 1, fromAddress)
	attempt1_1 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(1))
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(2))
	attempt1_2.BroadcastBeforeBlockNum = &b
	require.NoError(t, borm.InsertEthTxAttempt(&attempt1_1))
	require.NoError(t, borm.InsertEthTxAttempt(&attempt1_2))

	etx2 := cltest.MustInsertUnconfirmedEthTx(t, borm, 2, fromAddress)
	attempt2_1 := newBroadcastLegacyEthTxAttempt(t, etx2.ID, int64(1))
	require.NoError(t, borm.InsertEthTxAttempt(&attempt2_1))

	etx3 := cltest.MustInsertUnconfirmedEthTx(t, borm, 3, fromAddress)
	attempt3_1 := newBroadcastLegacyEthTxAttempt(t, etx3.ID, int64(1))
	require.NoError(t, borm.InsertEthTxAttempt(&attempt3_1))

	pgtest.MustExec(t, db, `UPDATE eth_tx_attempts SET broadcast_before_block_num = 41 WHERE broadcast_before_block_num IS NULL`)

	t.Run("marks buried eth_txes as 'confirmed_missing_receipt'", func(t *testing.T) {
		txmReceipt0 := evmtypes.Receipt{
			TxHash:           attempt0_2.Hash,
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}
		txmReceipt3 := evmtypes.Receipt{
			TxHash:           attempt3_1.Hash,
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(4), nil)
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
			elems[0].Result = &txmReceipt0
			elems[1].Result = &evmtypes.Receipt{}
			// Second transaction stil unconfirmed
			elems[2].Result = &evmtypes.Receipt{}
			elems[3].Result = &evmtypes.Receipt{}
			// Third transaction still unconfirmed
			elems[4].Result = &evmtypes.Receipt{}
			// Fourth transaction is confirmed
			elems[5].Result = &txmReceipt3
		}).Once()

		// PERFORM
		// Block num of 43 is one higher than the receipt (as would generally be expected)
		require.NoError(t, ec.CheckForReceipts(ctx, 43))

		ethClient.AssertExpectations(t)

		// Expected state is that the "top" eth_tx is now confirmed, with the
		// two below it "confirmed_missing_receipt" and the "bottom" eth_tx also confirmed
		etx3, err := borm.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmed, etx3.State)

		ethReceipt := etx3.EthTxAttempts[0].EthReceipts[0]
		require.Equal(t, txmReceipt3.BlockHash, ethReceipt.BlockHash)

		etx2, err = borm.FindEthTxWithAttempts(etx2.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmedMissingReceipt, etx2.State)
		etx1, err = borm.FindEthTxWithAttempts(etx1.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmedMissingReceipt, etx1.State)

		etx0, err = borm.FindEthTxWithAttempts(etx0.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmed, etx0.State)

		require.Len(t, etx0.EthTxAttempts, 2)
		require.Len(t, etx0.EthTxAttempts[0].EthReceipts, 1)
		ethReceipt = etx0.EthTxAttempts[0].EthReceipts[0]
		require.Equal(t, txmReceipt0.BlockHash, ethReceipt.BlockHash)
	})

	// STATE
	// eth_txes with nonce 0 is confirmed
	// eth_txes with nonce 1 is confirmed_missing_receipt
	// eth_txes with nonce 2 is confirmed_missing_receipt
	// eth_txes with nonce 3 is confirmed

	t.Run("marks eth_txes with state 'confirmed_missing_receipt' as 'confirmed' if a receipt finally shows up", func(t *testing.T) {
		txmReceipt := evmtypes.Receipt{
			TxHash:           attempt2_1.Hash,
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(43),
			TransactionIndex: uint(1),
		}
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
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
			elems[2].Result = &txmReceipt
		}).Once()

		// PERFORM
		// Block num of 44 is one higher than the receipt (as would generally be expected)
		require.NoError(t, ec.CheckForReceipts(ctx, 44))

		ethClient.AssertExpectations(t)

		// Expected state is that the "top" two eth_txes are now confirmed, with the
		// one below it still "confirmed_missing_receipt" and the bottom one remains confirmed
		etx3, err := borm.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmed, etx3.State)
		etx2, err = borm.FindEthTxWithAttempts(etx2.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmed, etx2.State)

		ethReceipt := etx2.EthTxAttempts[0].EthReceipts[0]
		require.Equal(t, txmReceipt.BlockHash, ethReceipt.BlockHash)

		etx1, err = borm.FindEthTxWithAttempts(etx1.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmedMissingReceipt, etx1.State)
		etx0, err = borm.FindEthTxWithAttempts(etx0.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmed, etx0.State)
	})

	// STATE
	// eth_txes with nonce 0 is confirmed
	// eth_txes with nonce 1 is confirmed_missing_receipt
	// eth_txes with nonce 2 is confirmed
	// eth_txes with nonce 3 is confirmed

	t.Run("continues to leave eth_txes with state 'confirmed_missing_receipt' unchanged if at least one attempt is above ETH_FINALITY_DEPTH", func(t *testing.T) {
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
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

		// PERFORM
		// Block num of 80 puts the first attempt (21) below threshold but second attempt (41) still above
		require.NoError(t, ec.CheckForReceipts(ctx, 80))

		ethClient.AssertExpectations(t)

		// Expected state is that the "top" two eth_txes are now confirmed, with the
		// one below it still "confirmed_missing_receipt" and the bottom one remains confirmed
		etx3, err := borm.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmed, etx3.State)
		etx2, err = borm.FindEthTxWithAttempts(etx2.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmed, etx2.State)
		etx1, err = borm.FindEthTxWithAttempts(etx1.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmedMissingReceipt, etx1.State)
		etx0, err = borm.FindEthTxWithAttempts(etx0.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmed, etx0.State)
	})

	// STATE
	// eth_txes with nonce 0 is confirmed
	// eth_txes with nonce 1 is confirmed_missing_receipt
	// eth_txes with nonce 2 is confirmed
	// eth_txes with nonce 3 is confirmed

	t.Run("marks eth_Txes with state 'confirmed_missing_receipt' as 'errored' if a receipt fails to show up and all attempts are buried deeper than ETH_FINALITY_DEPTH", func(t *testing.T) {
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
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

		// PERFORM
		// Block num of 100 puts the first attempt (21) and second attempt (41) below threshold
		require.NoError(t, ec.CheckForReceipts(ctx, 100))

		ethClient.AssertExpectations(t)

		// Expected state is that the "top" two eth_txes are now confirmed, with the
		// one below it marked as "fatal_error" and the bottom one remains confirmed
		etx3, err := borm.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmed, etx3.State)
		etx2, err = borm.FindEthTxWithAttempts(etx2.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmed, etx2.State)
		etx1, err = borm.FindEthTxWithAttempts(etx1.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxFatalError, etx1.State)
		etx0, err = borm.FindEthTxWithAttempts(etx0.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxConfirmed, etx0.State)
	})
}

func TestEthConfirmer_CheckConfirmedMissingReceipt(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	state, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	cfg.Overrides.GlobalEvmFinalityDepth = null.IntFrom(50)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ec := cltest.NewEthConfirmer(t, db, ethClient, evmcfg, ethKeyStore, []ethkey.State{state}, nil)

	ctx := context.Background()

	// STATE
	// eth_txes with nonce 0 has two attempts, the later attempt with higher gas fees
	// eth_txes with nonce 1 has two attempts, the later attempt with higher gas fees
	// eth_txes with nonce 2 has one attempt
	originalBroadcastAt := time.Unix(1616509100, 0)
	etx0 := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, borm, 0, 1, originalBroadcastAt, fromAddress)
	attempt0_2 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(2))
	require.NoError(t, borm.InsertEthTxAttempt(&attempt0_2))
	etx1 := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, borm, 1, 1, originalBroadcastAt, fromAddress)
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(2))
	require.NoError(t, borm.InsertEthTxAttempt(&attempt1_2))
	etx2 := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, borm, 2, 1, originalBroadcastAt, fromAddress)
	attempt2_1 := etx2.EthTxAttempts[0]
	etx3 := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, borm, 3, 1, originalBroadcastAt, fromAddress)
	attempt3_1 := etx3.EthTxAttempts[0]

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

	ethClient.AssertExpectations(t)

	// Expected state is that the "top" eth_tx is untouched but the other two
	// are marked as unconfirmed
	etx0, err := borm.FindEthTxWithAttempts(etx0.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgr.EthTxConfirmedMissingReceipt, etx0.State)
	assert.Greater(t, etx0.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	etx1, err = borm.FindEthTxWithAttempts(etx1.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgr.EthTxUnconfirmed, etx1.State)
	assert.Greater(t, etx1.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	etx2, err = borm.FindEthTxWithAttempts(etx2.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgr.EthTxUnconfirmed, etx2.State)
	assert.Greater(t, etx2.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	etx3, err = borm.FindEthTxWithAttempts(etx3.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgr.EthTxConfirmedMissingReceipt, etx3.State)
	assert.Greater(t, etx3.BroadcastAt.Unix(), originalBroadcastAt.Unix())
}

func TestEthConfirmer_CheckConfirmedMissingReceipt_batchSendTransactions_fails(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	state, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	cfg.Overrides.GlobalEvmFinalityDepth = null.IntFrom(50)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ec := cltest.NewEthConfirmer(t, db, ethClient, evmcfg, ethKeyStore, []ethkey.State{state}, nil)

	ctx := context.Background()

	// STATE
	// eth_txes with nonce 0 has two attempts, the later attempt with higher gas fees
	// eth_txes with nonce 1 has two attempts, the later attempt with higher gas fees
	// eth_txes with nonce 2 has one attempt
	originalBroadcastAt := time.Unix(1616509100, 0)
	etx0 := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, borm, 0, 1, originalBroadcastAt, fromAddress)
	attempt0_2 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(2))
	require.NoError(t, borm.InsertEthTxAttempt(&attempt0_2))
	etx1 := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, borm, 1, 1, originalBroadcastAt, fromAddress)
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(2))
	require.NoError(t, borm.InsertEthTxAttempt(&attempt1_2))
	etx2 := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, borm, 2, 1, originalBroadcastAt, fromAddress)
	attempt2_1 := etx2.EthTxAttempts[0]

	ethClient.On("BatchCallContextAll", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 3 &&
			cltest.BatchElemMatchesParams(b[0], hexutil.Encode(attempt0_2.SignedRawTx), "eth_sendRawTransaction") &&
			cltest.BatchElemMatchesParams(b[1], hexutil.Encode(attempt1_2.SignedRawTx), "eth_sendRawTransaction") &&
			cltest.BatchElemMatchesParams(b[2], hexutil.Encode(attempt2_1.SignedRawTx), "eth_sendRawTransaction")
	})).Return(errors.New("Timed out")).Once()

	// PERFORM
	require.NoError(t, ec.CheckConfirmedMissingReceipt(ctx))

	ethClient.AssertExpectations(t)

	// Expected state is that all txes are marked as unconfirmed, since the batch call had failed
	etx0, err := borm.FindEthTxWithAttempts(etx0.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgr.EthTxUnconfirmed, etx0.State)
	assert.Equal(t, etx0.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	etx1, err = borm.FindEthTxWithAttempts(etx1.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgr.EthTxUnconfirmed, etx1.State)
	assert.Equal(t, etx1.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	etx2, err = borm.FindEthTxWithAttempts(etx2.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgr.EthTxUnconfirmed, etx2.State)
	assert.Equal(t, etx2.BroadcastAt.Unix(), originalBroadcastAt.Unix())
}

func TestEthConfirmer_CheckConfirmedMissingReceipt_smallEvmRPCBatchSize_middleBatchSendTransactionFails(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	state, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	cfg.Overrides.GlobalEvmFinalityDepth = null.IntFrom(50)
	cfg.Overrides.GlobalEvmRPCDefaultBatchSize = null.IntFrom(1) // Set default batch size to 1
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	ec := cltest.NewEthConfirmer(t, db, ethClient, evmcfg, ethKeyStore, []ethkey.State{state}, nil)

	ctx := context.Background()

	// STATE
	// eth_txes with nonce 0 has two attempts, the later attempt with higher gas fees
	// eth_txes with nonce 1 has two attempts, the later attempt with higher gas fees
	// eth_txes with nonce 2 has one attempt
	originalBroadcastAt := time.Unix(1616509100, 0)
	etx0 := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, borm, 0, 1, originalBroadcastAt, fromAddress)
	attempt0_2 := newBroadcastLegacyEthTxAttempt(t, etx0.ID, int64(2))
	require.NoError(t, borm.InsertEthTxAttempt(&attempt0_2))
	etx1 := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, borm, 1, 1, originalBroadcastAt, fromAddress)
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID, int64(2))
	require.NoError(t, borm.InsertEthTxAttempt(&attempt1_2))
	etx2 := cltest.MustInsertConfirmedMissingReceiptEthTxWithLegacyAttempt(
		t, borm, 2, 1, originalBroadcastAt, fromAddress)

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

	ethClient.AssertExpectations(t)

	// Expected state is that all transactions since failed batch will be unconfirmed
	etx0, err := borm.FindEthTxWithAttempts(etx0.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgr.EthTxConfirmedMissingReceipt, etx0.State)
	assert.Greater(t, etx0.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	etx1, err = borm.FindEthTxWithAttempts(etx1.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgr.EthTxUnconfirmed, etx1.State)
	assert.Equal(t, etx1.BroadcastAt.Unix(), originalBroadcastAt.Unix())
	etx2, err = borm.FindEthTxWithAttempts(etx2.ID)
	assert.NoError(t, err)
	assert.Equal(t, txmgr.EthTxUnconfirmed, etx2.State)
	assert.Equal(t, etx2.BroadcastAt.Unix(), originalBroadcastAt.Unix())
}

func TestEthConfirmer_FindEthTxsRequiringResubmissionDueToInsufficientEth(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := cltest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	q := pg.NewQ(db, logger.TestLogger(t), cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)
	_, otherAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	// Insert order is mixed up to test sorting
	etx2 := cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, 1, fromAddress)
	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 2, fromAddress)
	attempt3_2 := cltest.NewLegacyEthTxAttempt(t, etx3.ID)
	attempt3_2.State = txmgr.EthTxAttemptInsufficientEth
	attempt3_2.GasPrice = utils.NewBig(big.NewInt(100))
	require.NoError(t, borm.InsertEthTxAttempt(&attempt3_2))
	etx1 := cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, 0, fromAddress)

	// These should never be returned
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 3, fromAddress)
	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 4, 100, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, 0, otherAddress)

	t.Run("returns all eth_txes with at least one attempt that is in insufficient_eth state", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringResubmissionDueToInsufficientEth(context.Background(), q, logger.TestLogger(t), fromAddress, cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 3)

		assert.Equal(t, *etx1.Nonce, *etxs[0].Nonce)
		assert.Equal(t, etx1.ID, etxs[0].ID)
		assert.Equal(t, *etx2.Nonce, *etxs[1].Nonce)
		assert.Equal(t, etx2.ID, etxs[1].ID)
		assert.Equal(t, *etx3.Nonce, *etxs[2].Nonce)
		assert.Equal(t, etx3.ID, etxs[2].ID)
	})

	t.Run("does not return eth_txes with different chain ID", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringResubmissionDueToInsufficientEth(context.Background(), q, logger.TestLogger(t), fromAddress, *big.NewInt(42))
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	t.Run("does not return confirmed or fatally errored eth_txes", func(t *testing.T) {
		pgtest.MustExec(t, db, `UPDATE eth_txes SET state='confirmed' WHERE id = $1`, etx1.ID)
		pgtest.MustExec(t, db, `UPDATE eth_txes SET state='fatal_error', nonce=NULL, error='foo', broadcast_at=NULL, initial_broadcast_at=NULL WHERE id = $1`, etx2.ID)

		etxs, err := txmgr.FindEthTxsRequiringResubmissionDueToInsufficientEth(context.Background(), q, logger.TestLogger(t), fromAddress, cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 1)

		assert.Equal(t, *etx3.Nonce, *etxs[0].Nonce)
		assert.Equal(t, etx3.ID, etxs[0].ID)
	})
}

func TestEthConfirmer_FindEthTxsRequiringRebroadcast(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)
	q := pg.NewQ(db, logger.TestLogger(t), cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	currentHead := int64(30)
	gasBumpThreshold := int64(10)
	tooNew := int64(21)
	onTheMoney := int64(20)
	oldEnough := int64(19)
	nonce := int64(0)

	mustInsertConfirmedEthTx(t, borm, nonce, fromAddress)
	nonce++

	_, otherAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	lggr := logger.TestLogger(t)

	t.Run("returns nothing when there are no transactions", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 10, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	mustInsertInProgressEthTx(t, borm, nonce, fromAddress)
	nonce++

	t.Run("returns nothing when the transaction is in_progress", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 10, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	// This one has BroadcastBeforeBlockNum set as nil... which can happen, but it should be ignored
	cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, nonce, fromAddress)
	nonce++

	t.Run("ignores unconfirmed transactions with nil BroadcastBeforeBlockNum", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 10, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, nonce, fromAddress)
	nonce++
	attempt1_1 := etx1.EthTxAttempts[0]
	require.NoError(t, db.Get(&attempt1_1, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, tooNew, attempt1_1.ID))
	attempt1_2 := newBroadcastLegacyEthTxAttempt(t, etx1.ID)
	attempt1_2.BroadcastBeforeBlockNum = &onTheMoney
	attempt1_2.GasPrice = utils.NewBigI(30000)
	require.NoError(t, borm.InsertEthTxAttempt(&attempt1_2))

	t.Run("returns nothing when the transaction is unconfirmed with an attempt that is recent", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 10, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, nonce, fromAddress)
	nonce++
	attempt2_1 := etx2.EthTxAttempts[0]
	require.NoError(t, db.Get(&attempt2_1, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, tooNew, attempt2_1.ID))

	t.Run("returns nothing when the transaction has attempts that are too new", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 10, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	etxWithoutAttempts := cltest.NewEthTx(t, fromAddress)
	{
		n := nonce
		etxWithoutAttempts.Nonce = &n
	}
	now := time.Now()
	etxWithoutAttempts.BroadcastAt = &now
	etxWithoutAttempts.InitialBroadcastAt = &now
	etxWithoutAttempts.State = txmgr.EthTxUnconfirmed
	require.NoError(t, borm.InsertEthTx(&etxWithoutAttempts))
	nonce++

	t.Run("does nothing if the transaction is from a different address than the one given", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, otherAddress, currentHead, gasBumpThreshold, 10, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	t.Run("returns the transaction if it is unconfirmed and has no attempts (note that this is an invariant violation, but we handle it anyway)", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 10, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 1)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
	})

	t.Run("returns nothing for different chain id", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 10, 0, *big.NewInt(42))
		require.NoError(t, err)

		require.Len(t, etxs, 0)
	})

	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, nonce, fromAddress)
	nonce++
	attempt3_1 := etx3.EthTxAttempts[0]
	require.NoError(t, db.Get(&attempt3_1, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_1.ID))

	// NOTE: It should ignore qualifying eth_txes from a different address
	etxOther := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, otherAddress)
	attemptOther1 := etxOther.EthTxAttempts[0]
	require.NoError(t, db.Get(&attemptOther1, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attemptOther1.ID))

	t.Run("returns the transaction if it is unconfirmed with an attempt that is older than gasBumpThreshold blocks", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 10, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
	})

	t.Run("returns nothing if threshold is zero", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, 0, 10, 0, cltest.FixtureChainID)
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
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 4, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 1) // returns etxWithoutAttempts only - eligible for gas bumping because it technically doesn't have any attempts withing gasBumpThreshold blocks
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)

		etxs, err = txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 5, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2) // includes etxWithoutAttempts, etx3 and etx4
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)

		// Zero limit disables it
		etxs, err = txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 0, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2) // includes etxWithoutAttempts, etx3 and etx4
	})

	etx4 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, nonce, fromAddress)
	nonce++
	attempt4_1 := etx4.EthTxAttempts[0]
	require.NoError(t, db.Get(&attempt4_1, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt4_1.ID))

	t.Run("ignores pending transactions for another key", func(t *testing.T) {
		// Re-use etx3 nonce for another key, it should not affect the results for this key
		etxOther := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, *etx3.Nonce, otherAddress)
		aOther := etxOther.EthTxAttempts[0]
		require.NoError(t, db.Get(&aOther, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, aOther.ID))

		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 6, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 3) // includes etxWithoutAttempts, etx3 and etx4
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
		assert.Equal(t, etx4.ID, etxs[2].ID)
	})

	attempt3_2 := newBroadcastLegacyEthTxAttempt(t, etx3.ID)
	attempt3_2.BroadcastBeforeBlockNum = &oldEnough
	attempt3_2.GasPrice = utils.NewBigI(30000)
	require.NoError(t, borm.InsertEthTxAttempt(&attempt3_2))

	t.Run("returns the transaction if it is unconfirmed with two attempts that are older than gasBumpThreshold blocks", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 10, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 3)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
		assert.Equal(t, etx4.ID, etxs[2].ID)
	})

	attempt3_3 := newBroadcastLegacyEthTxAttempt(t, etx3.ID)
	attempt3_3.BroadcastBeforeBlockNum = &tooNew
	attempt3_3.GasPrice = utils.NewBigI(40000)
	require.NoError(t, borm.InsertEthTxAttempt(&attempt3_3))

	t.Run("does not return the transaction if it has some older but one newer attempt", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 10, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, *etxWithoutAttempts.Nonce, *(etxs[0].Nonce))
		require.Equal(t, int64(5), *etxWithoutAttempts.Nonce)
		assert.Equal(t, etx4.ID, etxs[1].ID)
		assert.Equal(t, *etx4.Nonce, *(etxs[1].Nonce))
		require.Equal(t, int64(7), *etx4.Nonce)
	})

	attempt0_1 := newBroadcastLegacyEthTxAttempt(t, etxWithoutAttempts.ID)
	attempt0_1.State = txmgr.EthTxAttemptInsufficientEth
	require.NoError(t, borm.InsertEthTxAttempt(&attempt0_1))

	// This attempt has insufficient_eth, but there is also another attempt4_1
	// which is old enough, so this will be caught by both queries and should
	// not be duplicated
	attempt4_2 := cltest.NewLegacyEthTxAttempt(t, etx4.ID)
	attempt4_2.State = txmgr.EthTxAttemptInsufficientEth
	attempt4_2.GasPrice = utils.NewBigI(40000)
	require.NoError(t, borm.InsertEthTxAttempt(&attempt4_2))

	etx5 := cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, nonce, fromAddress)
	nonce++

	// This etx has one attempt that is too new, which would exclude it from
	// the gas bumping query, but it should still be caught by the insufficient
	// eth query
	etx6 := cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, nonce, fromAddress)
	attempt6_2 := newBroadcastLegacyEthTxAttempt(t, etx3.ID)
	attempt6_2.BroadcastBeforeBlockNum = &tooNew
	attempt6_2.GasPrice = utils.NewBigI(30001)
	require.NoError(t, borm.InsertEthTxAttempt(&attempt6_2))
	nonce++

	t.Run("returns unique attempts requiring resubmission due to insufficient eth, ordered by nonce asc", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 10, 0, cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 4)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, *etxWithoutAttempts.Nonce, *(etxs[0].Nonce))
		assert.Equal(t, etx4.ID, etxs[1].ID)
		assert.Equal(t, *etx4.Nonce, *(etxs[1].Nonce))
		assert.Equal(t, etx5.ID, etxs[2].ID)
		assert.Equal(t, *etx5.Nonce, *(etxs[2].Nonce))
		assert.Equal(t, etx6.ID, etxs[3].ID)
		assert.Equal(t, *etx6.Nonce, *(etxs[3].Nonce))
	})

	t.Run("applies limit", func(t *testing.T) {
		etxs, err := txmgr.FindEthTxsRequiringRebroadcast(context.Background(), q, lggr, fromAddress, currentHead, gasBumpThreshold, 10, 2, cltest.FixtureChainID)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, *etxWithoutAttempts.Nonce, *(etxs[0].Nonce))
		assert.Equal(t, etx4.ID, etxs[1].ID)
		assert.Equal(t, *etx4.Nonce, *(etxs[1].Nonce))
	})
}

func TestEthConfirmer_RebroadcastWhereNecessary(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	cfg.Overrides.GlobalEvmMaxGasPriceWei = assets.GWei(500)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	otherKey, _ := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	state, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	keys := []ethkey.State{state, otherKey}

	kst := new(ksmocks.Eth)
	kst.Test(t)
	// Use a mock keystore for this test
	ec := cltest.NewEthConfirmer(t, db, ethClient, evmcfg, kst, keys, nil)
	currentHead := int64(30)
	oldEnough := int64(19)
	nonce := int64(0)

	t.Run("does nothing if no transactions require bumping", func(t *testing.T) {
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
	})

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, nonce, fromAddress, originalBroadcastAt)
	nonce++
	attempt1_1 := etx.EthTxAttempts[0]
	require.NoError(t, db.Get(&attempt1_1, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_1.ID))
	var err error

	t.Run("re-sends previous transaction on keystore error", func(t *testing.T) {
		// simulate bumped transaction that is somehow impossible to sign
		kst.On("SignTx", fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				return tx.Nonce() == uint64(*etx.Nonce)
			}),
			mock.Anything).Return(nil, errors.New("signing error")).Once()

		// Do the thing
		err = ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead)
		require.Error(t, err)
		require.Contains(t, err.Error(), "signing error")

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Equal(t, txmgr.EthTxUnconfirmed, etx.State)

		require.Len(t, etx.EthTxAttempts, 1)

		kst.AssertExpectations(t)
	})

	t.Run("does nothing and continues on fatal error", func(t *testing.T) {
		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if tx.Nonce() != uint64(*etx.Nonce) {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(evmcfg.ChainID()) == 0
			})).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Nonce)
		})).Return(errors.New("exceeds block gas limit")).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.EthTxAttempts, 1)

		kst.AssertExpectations(t)
		ethClient.AssertExpectations(t)
	})

	ethClient = cltest.NewEthClientMockWithDefaultChain(t)
	txmgr.SetEthClientOnEthConfirmer(ethClient, ec)

	t.Run("does nothing and continues if bumped attempt transaction was too expensive", func(t *testing.T) {
		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if tx.Nonce() != uint64(*etx.Nonce) {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(evmcfg.ChainID()) == 0
			})).Return(&ethTx, nil).Once()

		// Once for the bumped attempt which exceeds limit
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Nonce) && tx.GasPrice().Int64() == int64(20000000000)
		})).Return(errors.New("tx fee (1.10 ether) exceeds the configured cap (1.00 ether)")).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		// Did not create an additional attempt
		require.Len(t, etx.EthTxAttempts, 1)

		// broadcast_at did not change
		require.Equal(t, etx.BroadcastAt.Unix(), originalBroadcastAt.Unix())
		require.Equal(t, etx.InitialBroadcastAt.Unix(), originalBroadcastAt.Unix())

		kst.AssertExpectations(t)
		ethClient.AssertExpectations(t)
	})

	var attempt1_2 txmgr.EthTxAttempt
	ethClient = cltest.NewEthClientMockWithDefaultChain(t)
	txmgr.SetEthClientOnEthConfirmer(ethClient, ec)

	t.Run("creates new attempt with higher gas price if transaction has an attempt older than threshold", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.GasPrice.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(evmcfg.ChainID()) == 0
			})).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.EthTxAttempts, 2)
		require.Equal(t, attempt1_1.ID, etx.EthTxAttempts[1].ID)

		// Got the new attempt
		attempt1_2 = etx.EthTxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.GasPrice.ToInt().Int64())
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt1_2.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("does nothing if there is an attempt without BroadcastBeforeBlockNum set", func(t *testing.T) {
		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.EthTxAttempts, 2)
	})

	require.NoError(t, db.Get(&attempt1_2, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_2.ID))
	var attempt1_3 txmgr.EthTxAttempt

	t.Run("creates new attempt with higher gas price if transaction is already in mempool (e.g. due to previous crash before we could save the new attempt)", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(25000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_2.GasPrice.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if int64(tx.Nonce()) != *etx.Nonce || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(fmt.Errorf("known transaction: %s", ethTx.Hash().Hex())).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.EthTxAttempts, 3)
		require.Equal(t, attempt1_1.ID, etx.EthTxAttempts[2].ID)
		require.Equal(t, attempt1_2.ID, etx.EthTxAttempts[1].ID)

		// Got the new attempt
		attempt1_3 = etx.EthTxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_3.GasPrice.ToInt().Int64())
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt1_3.State)

		kst.AssertExpectations(t)
		ethClient.AssertExpectations(t)
	})

	require.NoError(t, db.Get(&attempt1_3, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_3.ID))
	var attempt1_4 txmgr.EthTxAttempt

	t.Run("saves new attempt even for transaction that has already been confirmed (nonce already used)", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(30000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_2.GasPrice.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		receipt := evmtypes.Receipt{BlockNumber: big.NewInt(40)}
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if int64(tx.Nonce()) != *etx.Nonce || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				receipt.TxHash = tx.Hash()
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(errors.New("nonce too low")).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxConfirmedMissingReceipt, etx.State)

		// Got the new attempt
		attempt1_4 = etx.EthTxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_4.GasPrice.ToInt().Int64())

		require.Len(t, etx.EthTxAttempts, 4)
		require.Equal(t, attempt1_1.ID, etx.EthTxAttempts[3].ID)
		require.Equal(t, attempt1_2.ID, etx.EthTxAttempts[2].ID)
		require.Equal(t, attempt1_3.ID, etx.EthTxAttempts[1].ID)
		require.Equal(t, attempt1_4.ID, etx.EthTxAttempts[0].ID)
		require.Equal(t, txmgr.EthTxAttemptBroadcast, etx.EthTxAttempts[0].State)
		require.Equal(t, txmgr.EthTxAttemptBroadcast, etx.EthTxAttempts[1].State)
		require.Equal(t, txmgr.EthTxAttemptBroadcast, etx.EthTxAttempts[2].State)
		require.Equal(t, txmgr.EthTxAttemptBroadcast, etx.EthTxAttempts[3].State)

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	// Mark original tx as confirmed so we won't pick it up any more
	pgtest.MustExec(t, db, `UPDATE eth_txes SET state = 'confirmed'`)

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, nonce, fromAddress)
	nonce++
	attempt2_1 := etx2.EthTxAttempts[0]
	require.NoError(t, db.Get(&attempt2_1, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt2_1.ID))
	var attempt2_2 txmgr.EthTxAttempt

	t.Run("saves in_progress attempt on temporary error and returns error", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt2_1.GasPrice.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		n := *etx2.Nonce
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if int64(tx.Nonce()) != n || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == n && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(errors.New("some network error")).Once()

		// Do the thing
		err = ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead)
		require.Error(t, err)
		require.Contains(t, err.Error(), "some network error")

		etx2, err = borm.FindEthTxWithAttempts(etx2.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxUnconfirmed, etx2.State)

		// Old attempt is untouched
		require.Len(t, etx2.EthTxAttempts, 2)
		require.Equal(t, attempt2_1.ID, etx2.EthTxAttempts[1].ID)
		attempt2_1 = etx2.EthTxAttempts[1]
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt2_1.State)
		assert.Equal(t, oldEnough, *attempt2_1.BroadcastBeforeBlockNum)

		// New in_progress attempt saved
		attempt2_2 = etx2.EthTxAttempts[0]
		assert.Equal(t, txmgr.EthTxAttemptInProgress, attempt2_2.State)
		assert.Nil(t, attempt2_2.BroadcastBeforeBlockNum)

		// Do it again and move the attempt into "broadcast"
		n = *etx2.Nonce
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == n && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(nil).Once()

		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		// Attempt marked "broadcast"
		etx2, err = borm.FindEthTxWithAttempts(etx2.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxUnconfirmed, etx2.State)

		// New in_progress attempt saved
		require.Len(t, etx2.EthTxAttempts, 2)
		require.Equal(t, attempt2_2.ID, etx2.EthTxAttempts[0].ID)
		attempt2_2 = etx2.EthTxAttempts[0]
		require.Equal(t, txmgr.EthTxAttemptBroadcast, attempt2_2.State)
		assert.Nil(t, attempt2_2.BroadcastBeforeBlockNum)

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	// Set BroadcastBeforeBlockNum again so the next test will pick it up
	require.NoError(t, db.Get(&attempt2_2, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt2_2.ID))

	t.Run("assumes that 'nonce too low' error means confirmed_missing_receipt", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(25000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt2_1.GasPrice.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		n := *etx2.Nonce
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if int64(tx.Nonce()) != n || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == n && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(errors.New("nonce too low")).Once()

		// Creates new attempt as normal if currentHead is not high enough
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
		etx2, err = borm.FindEthTxWithAttempts(etx2.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxConfirmedMissingReceipt, etx2.State)

		// One new attempt saved
		require.Len(t, etx2.EthTxAttempts, 3)
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, etx2.EthTxAttempts[0].State)
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, etx2.EthTxAttempts[1].State)
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, etx2.EthTxAttempts[2].State)

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	// Original tx is confirmed so we won't pick it up any more
	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, nonce, fromAddress)
	nonce++
	attempt3_1 := etx3.EthTxAttempts[0]
	require.NoError(t, db.Get(&attempt3_1, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1, gas_price=$2 WHERE id=$3 RETURNING *`, oldEnough, utils.NewBig(big.NewInt(35000000000)), attempt3_1.ID))

	var attempt3_2 txmgr.EthTxAttempt

	t.Run("saves attempt anyway if replacement transaction is underpriced because the bumped gas price is insufficiently higher than the previous one", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(42000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt3_1.GasPrice.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if int64(tx.Nonce()) != *etx3.Nonce || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == *etx3.Nonce && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(errors.New("replacement transaction underpriced")).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx3, err = borm.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxUnconfirmed, etx3.State)

		require.Len(t, etx3.EthTxAttempts, 2)
		require.Equal(t, attempt3_1.ID, etx3.EthTxAttempts[1].ID)
		attempt3_2 = etx3.EthTxAttempts[0]

		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt3_2.GasPrice.ToInt().Int64())

		kst.AssertExpectations(t)
		ethClient.AssertExpectations(t)
	})

	require.NoError(t, db.Get(&attempt3_2, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_2.ID))
	var attempt3_3 txmgr.EthTxAttempt

	t.Run("handles case where transaction is already known somehow", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(50400000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt3_1.GasPrice.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if int64(tx.Nonce()) != *etx3.Nonce || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == *etx3.Nonce && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(fmt.Errorf("known transaction: %s", ethTx.Hash().Hex())).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx3, err = borm.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxUnconfirmed, etx3.State)

		require.Len(t, etx3.EthTxAttempts, 3)
		attempt3_3 = etx3.EthTxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt3_3.GasPrice.ToInt().Int64())

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	require.NoError(t, db.Get(&attempt3_3, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_3.ID))
	var attempt3_4 txmgr.EthTxAttempt

	t.Run("pretends it was accepted and continues the cycle if rejected for being temporarily underpriced", func(t *testing.T) {
		// This happens if parity is rejecting transactions that are not priced high enough to even get into the mempool at all
		// It should pretend it was accepted into the mempool and hand off to the next cycle to continue bumping gas as normal
		temporarilyUnderpricedError := "There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee."

		expectedBumpedGasPrice := big.NewInt(60480000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt3_2.GasPrice.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if int64(tx.Nonce()) != *etx3.Nonce || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == *etx3.Nonce && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(errors.New(temporarilyUnderpricedError)).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx3, err = borm.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxUnconfirmed, etx3.State)

		require.Len(t, etx3.EthTxAttempts, 4)
		attempt3_4 = etx3.EthTxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt3_4.GasPrice.ToInt().Int64())

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	require.NoError(t, db.Get(&attempt3_4, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_4.ID))

	t.Run("resubmits at the old price and does not create a new attempt if one of the bumped transactions would exceed ETH_MAX_GAS_PRICE_WEI", func(t *testing.T) {
		// Set price such that the next bump will exceed ETH_MAX_GAS_PRICE_WEI
		// Existing gas price is: 60480000000
		gasPrice := attempt3_4.GasPrice.ToInt()
		cfg.Overrides.GlobalEvmMaxGasPriceWei = assets.Wei(60500000000)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == *etx3.Nonce && gasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(errors.New("already known")).Once() // we already submitted at this price, now its time to bump and submit again but since we simply resubmitted rather than increasing gas price, geth already knows about this tx

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx3, err = borm.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxUnconfirmed, etx3.State)

		// No new tx attempts
		require.Len(t, etx3.EthTxAttempts, 4)
		attempt3_4 = etx3.EthTxAttempts[0]
		assert.Equal(t, gasPrice.Int64(), attempt3_4.GasPrice.ToInt().Int64())

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	require.NoError(t, db.Get(&attempt3_4, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt3_4.ID))

	t.Run("resubmits at the old price and does not create a new attempt if the current price is exactly ETH_MAX_GAS_PRICE_WEI", func(t *testing.T) {
		// Set price such that the current price is already at ETH_MAX_GAS_PRICE_WEI
		// Existing gas price is: 60480000000
		gasPrice := attempt3_4.GasPrice.ToInt()
		cfg.Overrides.GlobalEvmMaxGasPriceWei = assets.Wei(60480000000)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == *etx3.Nonce && gasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(errors.New("already known")).Once() // we already submitted at this price, now its time to bump and submit again but since we simply resubmitted rather than increasing gas price, geth already knows about this tx

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx3, err = borm.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxUnconfirmed, etx3.State)

		// No new tx attempts
		require.Len(t, etx3.EthTxAttempts, 4)
		attempt3_4 := etx3.EthTxAttempts[0]
		assert.Equal(t, gasPrice.Int64(), attempt3_4.GasPrice.ToInt().Int64())

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	// The EIP-1559 etx and attempt
	etx4 := cltest.MustInsertUnconfirmedEthTxWithBroadcastDynamicFeeAttempt(t, borm, nonce, fromAddress)
	nonce++
	attempt4_1 := etx4.EthTxAttempts[0]
	require.NoError(t, db.Get(&attempt4_1, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1, gas_tip_cap=$2, gas_fee_cap=$3 WHERE id=$4 RETURNING *`,
		oldEnough, utils.NewBig(assets.GWei(35)), utils.NewBig(assets.GWei(100)), attempt4_1.ID))
	var attempt4_2 txmgr.EthTxAttempt

	t.Run("EIP-1559: bumps using EIP-1559 rules when existing attempts are of type 0x2", func(t *testing.T) {
		cfg.Overrides.GlobalEvmMaxGasPriceWei = assets.GWei(1000)
		ethTx := *types.NewTx(&types.DynamicFeeTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if int64(tx.Nonce()) != *etx4.Nonce {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		// This is the new, EIP-1559 attempt
		gasTipCap := assets.GWei(42)
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == *etx4.Nonce && gasTipCap.Cmp(tx.GasTipCap()) == 0
		})).Return(nil).Once()

		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx4, err = borm.FindEthTxWithAttempts(etx4.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxUnconfirmed, etx4.State)

		// A new, bumped attempt
		require.Len(t, etx4.EthTxAttempts, 2)
		attempt4_2 = etx4.EthTxAttempts[0]
		assert.Nil(t, attempt4_2.GasPrice)
		assert.Equal(t, assets.GWei(42).String(), attempt4_2.GasTipCap.String())
		assert.Equal(t, assets.GWei(120).String(), attempt4_2.GasFeeCap.String())
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt1_2.State)

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	require.NoError(t, db.Get(&attempt4_2, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1, gas_tip_cap=$2, gas_fee_cap=$3 WHERE id=$4 RETURNING *`,
		oldEnough, utils.NewBig(assets.GWei(999)), utils.NewBig(assets.GWei(1000)), attempt4_2.ID))

	t.Run("EIP-1559: resubmits at the old price and does not create a new attempt if one of the bumped EIP-1559 transactions would have its tip cap exceed ETH_MAX_GAS_PRICE_WEI", func(t *testing.T) {
		cfg.Overrides.GlobalEvmMaxGasPriceWei = assets.GWei(1000)

		// Third attempt failed to bump, resubmits old one instead
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == *etx4.Nonce && attempt4_2.Hash == tx.Hash()
		})).Return(nil).Once()

		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx4, err = borm.FindEthTxWithAttempts(etx4.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxUnconfirmed, etx4.State)

		// No new tx attempts
		require.Len(t, etx4.EthTxAttempts, 2)
		attempt4_2 := etx4.EthTxAttempts[0]
		assert.Equal(t, assets.GWei(999).Int64(), attempt4_2.GasTipCap.ToInt().Int64())
		assert.Equal(t, assets.GWei(1000).Int64(), attempt4_2.GasFeeCap.ToInt().Int64())

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	require.NoError(t, db.Get(&attempt4_2, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1, gas_tip_cap=$2, gas_fee_cap=$3 WHERE id=$4 RETURNING *`,
		oldEnough, utils.NewBig(assets.GWei(45)), utils.NewBig(assets.GWei(100)), attempt4_2.ID))

	t.Run("EIP-1559: saves attempt anyway if replacement transaction is underpriced because the bumped gas price is insufficiently higher than the previous one", func(t *testing.T) {
		// NOTE: This test case was empirically impossible when I tried it on eth mainnet (any EIP1559 transaction with a higher tip cap is accepted even if it's only 1 wei more) but appears to be possible on Polygon/Matic, probably due to poor design that applies the 10% minumum to the overall value (base fee + tip cap)
		expectedBumpedTipCap := assets.GWei(54)
		require.Greater(t, expectedBumpedTipCap.Int64(), attempt4_2.GasTipCap.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				if int64(tx.Nonce()) != *etx4.Nonce || expectedBumpedTipCap.Cmp(tx.GasTipCap()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == *etx4.Nonce && expectedBumpedTipCap.Cmp(tx.GasTipCap()) == 0
		})).Return(errors.New("replacement transaction underpriced")).Once()

		// Do it
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx4, err = borm.FindEthTxWithAttempts(etx4.ID)
		require.NoError(t, err)

		assert.Equal(t, txmgr.EthTxUnconfirmed, etx4.State)

		require.Len(t, etx4.EthTxAttempts, 3)
		require.Equal(t, attempt4_1.ID, etx4.EthTxAttempts[2].ID)
		require.Equal(t, attempt4_2.ID, etx4.EthTxAttempts[1].ID)
		attempt4_3 := etx4.EthTxAttempts[0]

		assert.Equal(t, expectedBumpedTipCap.Int64(), attempt4_3.GasTipCap.ToInt().Int64())

		kst.AssertExpectations(t)
		ethClient.AssertExpectations(t)

	})

	kst.AssertExpectations(t)
	ethClient.AssertExpectations(t)
}

func TestEthConfirmer_RebroadcastWhereNecessary_TerminallyUnderpriced_ThenGoesThrough(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	cfg.Overrides.GlobalEvmMaxGasPriceWei = assets.GWei(500)
	evmcfg := evmtest.NewChainScopedConfig(t, cfg)

	otherKey, _ := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	state, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore)
	keys := []ethkey.State{state, otherKey}

	kst := new(ksmocks.Eth)
	kst.Test(t)
	// Use a mock keystore for this test
	ec := cltest.NewEthConfirmer(t, db, ethClient, evmcfg, kst, keys, nil)
	currentHead := int64(30)
	nonce := int64(0)

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx := cltest.MustInsertUnconfrimedEthTxWithAttemptState(t, borm, nonce, fromAddress, txmgr.EthTxAttemptInProgress, originalBroadcastAt)
	require.Equal(t, originalBroadcastAt, *etx.BroadcastAt)
	nonce++
	attempt := etx.EthTxAttempts[0]
	signedTx, err := attempt.GetSignedTx()
	require.NoError(t, err)

	t.Run("terminally underpriced transactions are retried with more gas", func(t *testing.T) {
		// Fail the first time with terminally underpriced.
		ethClient.On("SendTransaction", mock.Anything, mock.Anything).Return(
			errors.New("Transaction gas price is too low. It does not satisfy your node's minimal gas price"),
		).Once()
		// Succeed the second time after bumping gas.
		ethClient.On("SendTransaction", mock.Anything, mock.Anything).Return(nil)
		kst.On("SignTx", mock.Anything, mock.Anything, mock.Anything).Return(
			signedTx, nil,
		)
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))
	})
}

func TestEthConfirmer_RebroadcastWhereNecessary_WhenOutOfEth(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	_, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	keys, err := ethKeyStore.SendingKeys(nil)
	require.NoError(t, err)
	keyStates, err := ethKeyStore.GetStatesForKeys(keys)
	require.NoError(t, err)

	config := newTestChainScopedConfig(t)
	currentHead := int64(30)
	oldEnough := int64(19)
	nonce := int64(0)

	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, nonce, fromAddress)
	nonce++
	attempt1_1 := etx.EthTxAttempts[0]
	require.NoError(t, db.Get(&attempt1_1, `UPDATE eth_tx_attempts SET broadcast_before_block_num=$1 WHERE id=$2 RETURNING *`, oldEnough, attempt1_1.ID))
	var attempt1_2 txmgr.EthTxAttempt

	insufficientEthError := errors.New("insufficient funds for gas * price + value")

	t.Run("saves attempt with state 'insufficient_eth' if eth node returns this error", func(t *testing.T) {
		ec := cltest.NewEthConfirmer(t, db, ethClient, config, ethKeyStore, keyStates, nil)

		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.GasPrice.ToInt().Int64())

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(insufficientEthError).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.EthTxAttempts, 2)
		require.Equal(t, attempt1_1.ID, etx.EthTxAttempts[1].ID)

		// Got the new attempt
		attempt1_2 = etx.EthTxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.GasPrice.ToInt().Int64())
		assert.Equal(t, txmgr.EthTxAttemptInsufficientEth, attempt1_2.State)
		assert.Nil(t, attempt1_2.BroadcastBeforeBlockNum)

		ethClient.AssertExpectations(t)
	})

	t.Run("does not bump gas when previous error was 'out of eth', instead resubmits existing transaction", func(t *testing.T) {
		ec := cltest.NewEthConfirmer(t, db, ethClient, config, ethKeyStore, keyStates, nil)

		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.GasPrice.ToInt().Int64())

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(insufficientEthError).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		// New attempt was NOT created
		require.Len(t, etx.EthTxAttempts, 2)

		// The attempt is still "out of eth"
		attempt1_2 = etx.EthTxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.GasPrice.ToInt().Int64())
		assert.Equal(t, txmgr.EthTxAttemptInsufficientEth, attempt1_2.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("saves the attempt as broadcast after node wallet has been topped up with sufficient balance", func(t *testing.T) {
		ec := cltest.NewEthConfirmer(t, db, ethClient, config, ethKeyStore, keyStates, nil)

		expectedBumpedGasPrice := big.NewInt(20000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.GasPrice.ToInt().Int64())

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		etx, err = borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		// New attempt was NOT created
		require.Len(t, etx.EthTxAttempts, 2)

		// Attempt is now 'broadcast'
		attempt1_2 = etx.EthTxAttempts[0]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.GasPrice.ToInt().Int64())
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt1_2.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("resubmitting due to insufficient eth is not limited by ETH_GAS_BUMP_TX_DEPTH", func(t *testing.T) {
		depth := 2
		etxCount := 4

		cfg := configtest.NewTestGeneralConfig(t)
		cfg.Overrides.GlobalEvmGasBumpTxDepth = null.IntFrom(int64(depth))
		evmcfg := evmtest.NewChainScopedConfig(t, cfg)
		ec := cltest.NewEthConfirmer(t, db, ethClient, evmcfg, ethKeyStore, keyStates, nil)

		for i := 0; i < etxCount; i++ {
			n := nonce
			cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, borm, nonce, fromAddress)
			ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
				return tx.Nonce() == uint64(n)
			})).Return(nil).Once()

			nonce++
		}

		require.NoError(t, ec.RebroadcastWhereNecessary(testutils.Context(t), currentHead))

		var attempts []txmgr.EthTxAttempt
		require.NoError(t, db.Select(&attempts, "SELECT * FROM eth_tx_attempts WHERE state = 'insufficient_eth'"))
		require.Len(t, attempts, 0)

		ethClient.AssertExpectations(t)
	})
}

func TestEthConfirmer_EnsureConfirmedTransactionsInLongestChain(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()

	state, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

	config := newTestChainScopedConfig(t)
	ec := cltest.NewEthConfirmer(t, db, ethClient, config, ethKeyStore, []ethkey.State{state}, nil)

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

	t.Run("does nothing if there aren't any transactions", func(t *testing.T) {
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))
	})

	t.Run("does nothing to unconfirmed transactions", func(t *testing.T) {
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 0, fromAddress)

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))

		etx, err := borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxUnconfirmed, etx.State)
	})

	t.Run("does nothing to confirmed transactions with receipts within head height of the chain and included in the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 2, 1, fromAddress)
		attempt := etx.EthTxAttempts[0]
		cltest.MustInsertEthReceipt(t, borm, head.Number, head.Hash, attempt.Hash)

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))

		etx, err := borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxConfirmed, etx.State)
	})

	t.Run("does nothing to confirmed transactions that only have receipts older than the start of the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 3, 1, fromAddress)
		attempt := etx.EthTxAttempts[0]
		// Add receipt that is older than the lowest block of the chain
		cltest.MustInsertEthReceipt(t, borm, head.Parent.Parent.Number-1, utils.NewHash(), attempt.Hash)

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))

		etx, err := borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxConfirmed, etx.State)
	})

	t.Run("unconfirms and rebroadcasts transactions that have receipts within head height of the chain but not included in the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 4, 1, fromAddress)
		attempt := etx.EthTxAttempts[0]
		// Include one within head height but a different block hash
		cltest.MustInsertEthReceipt(t, borm, head.Parent.Number, utils.NewHash(), attempt.Hash)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			atx, err := attempt.GetSignedTx()
			require.NoError(t, err)
			// Keeps gas price and nonce the same
			return atx.GasPrice().Cmp(tx.GasPrice()) == 0 && atx.Nonce() == tx.Nonce()
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))

		etx, err := borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxUnconfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt = etx.EthTxAttempts[0]
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("unconfirms and rebroadcasts transactions that have receipts within head height of chain but not included in the chain even if a receipt exists older than the start of the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 5, 1, fromAddress)
		attempt := etx.EthTxAttempts[0]
		// Add receipt that is older than the lowest block of the chain
		cltest.MustInsertEthReceipt(t, borm, head.Parent.Parent.Number-1, utils.NewHash(), attempt.Hash)
		// Include one within head height but a different block hash
		cltest.MustInsertEthReceipt(t, borm, head.Parent.Number, utils.NewHash(), attempt.Hash)

		ethClient.On("SendTransaction", mock.Anything, mock.Anything).Return(nil).Once()

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))

		etx, err := borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxUnconfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt = etx.EthTxAttempts[0]
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("if more than one attempt has a receipt (should not be possible but isn't prevented by database constraints) unconfirms and rebroadcasts only the attempt with the highest gas price", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 6, 1, fromAddress)
		require.Len(t, etx.EthTxAttempts, 1)
		// Sanity check to assert the included attempt has the lowest gas price
		require.Less(t, etx.EthTxAttempts[0].GasPrice.ToInt().Int64(), int64(30000))

		attempt2 := newBroadcastLegacyEthTxAttempt(t, etx.ID, 30000)
		attempt2.SignedRawTx = hexutil.MustDecode("0xf88c8301f3a98503b9aca000832ab98094f5fff180082d6017036b771ba883025c654bc93580a4daa6d556000000000000000000000000000000000000000000000000000000000000000026a0f25601065ee369b6470c0399a2334afcfbeb0b5c8f3d9a9042e448ed29b5bcbda05b676e00248b85faf4dd889f0e2dcf91eb867e23ac9eeb14a73f9e4c14972cdf")
		attempt3 := newBroadcastLegacyEthTxAttempt(t, etx.ID, 40000)
		attempt3.SignedRawTx = hexutil.MustDecode("0xf88c8301f3a88503b9aca0008316e36094151445852b0cfdf6a4cc81440f2af99176e8ad0880a4daa6d556000000000000000000000000000000000000000000000000000000000000000026a0dcb5a7ad52b96a866257134429f944c505820716567f070e64abb74899803855a04c13eff2a22c218e68da80111e1bb6dc665d3dea7104ab40ff8a0275a99f630d")
		require.NoError(t, borm.InsertEthTxAttempt(&attempt2))
		require.NoError(t, borm.InsertEthTxAttempt(&attempt3))

		// Receipt is within head height but a different block hash
		cltest.MustInsertEthReceipt(t, borm, head.Parent.Number, utils.NewHash(), attempt2.Hash)
		// Receipt is within head height but a different block hash
		cltest.MustInsertEthReceipt(t, borm, head.Parent.Number, utils.NewHash(), attempt3.Hash)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			s, err := attempt3.GetSignedTx()
			require.NoError(t, err)
			return tx.Hash() == s.Hash()
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))

		etx, err := borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxUnconfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 3)
		attempt1 := etx.EthTxAttempts[0]
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt1.State)
		attempt2 = etx.EthTxAttempts[1]
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt2.State)
		attempt3 = etx.EthTxAttempts[2]
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt3.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("if receipt has a block number that is in the future, does not mark for rebroadcast (the safe thing to do is simply wait until heads catches up)", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 7, 1, fromAddress)
		attempt := etx.EthTxAttempts[0]
		// Add receipt that is higher than head
		cltest.MustInsertEthReceipt(t, borm, head.Number+1, utils.NewHash(), attempt.Hash)

		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(testutils.Context(t), &head))

		etx, err := borm.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, txmgr.EthTxConfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt = etx.EthTxAttempts[0]
		assert.Equal(t, txmgr.EthTxAttemptBroadcast, attempt.State)
		assert.Len(t, attempt.EthReceipts, 1)

		ethClient.AssertExpectations(t)
	})
}

func TestEthConfirmer_ForceRebroadcast(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, cfg)

	ethKeyStore := cltest.NewKeyStore(t, db, cfg).Eth()
	state, fromAddress := cltest.MustInsertRandomKeyReturningState(t, ethKeyStore, 0)

	config := newTestChainScopedConfig(t)
	mustInsertUnstartedEthTx(t, borm, fromAddress)
	mustInsertInProgressEthTx(t, borm, 0, fromAddress)
	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 1, fromAddress)
	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, borm, 2, fromAddress)

	gasPriceWei := uint64(assets.GWei(52).Int64())
	overrideGasLimit := uint64(20000)

	t.Run("rebroadcasts one eth_tx if it falls within in nonce range", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		ec := cltest.NewEthConfirmer(t, db, ethClient, config, ethKeyStore, []ethkey.State{state}, nil)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Nonce) &&
				uint64(tx.GasPrice().Int64()) == gasPriceWei &&
				tx.Gas() == overrideGasLimit &&
				reflect.DeepEqual(tx.Data(), etx1.EncodedPayload) &&
				*tx.To() == etx1.ToAddress
		})).Return(nil).Once()

		require.NoError(t, ec.ForceRebroadcast(1, 1, gasPriceWei, fromAddress, overrideGasLimit))

		ethClient.AssertExpectations(t)
	})

	t.Run("uses default gas limit if overrideGasLimit is 0", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		ec := cltest.NewEthConfirmer(t, db, ethClient, config, ethKeyStore, []ethkey.State{state}, nil)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Nonce) &&
				uint64(tx.GasPrice().Int64()) == gasPriceWei &&
				tx.Gas() == etx1.GasLimit &&
				reflect.DeepEqual(tx.Data(), etx1.EncodedPayload) &&
				*tx.To() == etx1.ToAddress
		})).Return(nil).Once()

		require.NoError(t, ec.ForceRebroadcast(1, 1, gasPriceWei, fromAddress, 0))

		ethClient.AssertExpectations(t)
	})

	t.Run("rebroadcasts several eth_txes in nonce range", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		ec := cltest.NewEthConfirmer(t, db, ethClient, config, ethKeyStore, []ethkey.State{state}, nil)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Nonce) && uint64(tx.GasPrice().Int64()) == gasPriceWei && tx.Gas() == overrideGasLimit
		})).Return(nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx2.Nonce) && uint64(tx.GasPrice().Int64()) == gasPriceWei && tx.Gas() == overrideGasLimit
		})).Return(nil).Once()

		require.NoError(t, ec.ForceRebroadcast(1, 2, gasPriceWei, fromAddress, overrideGasLimit))

		ethClient.AssertExpectations(t)
	})

	t.Run("broadcasts zero transactions if eth_tx doesn't exist for that nonce", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		ec := cltest.NewEthConfirmer(t, db, ethClient, config, ethKeyStore, []ethkey.State{state}, nil)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(1)
		})).Return(nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(2)
		})).Return(nil).Once()
		for i := 3; i <= 5; i++ {
			nonce := i
			ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
				return tx.Nonce() == uint64(nonce) &&
					uint64(tx.GasPrice().Int64()) == gasPriceWei &&
					tx.Gas() == overrideGasLimit &&
					*tx.To() == fromAddress &&
					tx.Value().Cmp(big.NewInt(0)) == 0 &&
					len(tx.Data()) == 0
			})).Return(nil).Once()
		}

		require.NoError(t, ec.ForceRebroadcast(1, 5, gasPriceWei, fromAddress, overrideGasLimit))

		ethClient.AssertExpectations(t)
	})

	t.Run("zero transactions use default gas limit if override wasn't specified", func(t *testing.T) {
		ethClient := cltest.NewEthClientMockWithDefaultChain(t)
		ec := cltest.NewEthConfirmer(t, db, ethClient, config, ethKeyStore, []ethkey.State{state}, nil)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(0) && uint64(tx.GasPrice().Int64()) == gasPriceWei && uint64(tx.Gas()) == config.EvmGasLimitDefault()
		})).Return(nil).Once()

		require.NoError(t, ec.ForceRebroadcast(0, 0, gasPriceWei, fromAddress, 0))

		ethClient.AssertExpectations(t)
	})
}

func TestEthConfirmer_ResumePendingRuns(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	config := configtest.NewTestGeneralConfig(t)
	borm := cltest.NewTxmORM(t, db, config)

	ethKeyStore := cltest.NewKeyStore(t, db, config).Eth()

	key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)
	state := cltest.MustGetStateForKey(t, ethKeyStore, key)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)

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
		ec := cltest.NewEthConfirmer(t, db, ethClient, evmcfg, ethKeyStore, []ethkey.State{state}, func(uuid.UUID, interface{}, error) error {
			t.Fatal("No value expected")
			return nil
		})

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 1, 1, fromAddress)
		attempt := etx.EthTxAttempts[0]
		cltest.MustInsertEthReceipt(t, borm, head.Number-minConfirmations, head.Hash, attempt.Hash)
		pgtest.MustExec(t, db, `UPDATE eth_txes SET pipeline_task_run_id = $1, min_confirmations = $2 WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		err := ec.ResumePendingTaskRuns(context.Background(), &head)
		require.NoError(t, err)

	})

	t.Run("doesn't process task runs where the receipt is younger than minConfirmations", func(t *testing.T) {
		ec := cltest.NewEthConfirmer(t, db, ethClient, evmcfg, ethKeyStore, []ethkey.State{state}, func(uuid.UUID, interface{}, error) error {
			t.Fatal("No value expected")
			return nil
		})

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 2, 1, fromAddress)
		attempt := etx.EthTxAttempts[0]
		cltest.MustInsertEthReceipt(t, borm, head.Number, head.Hash, attempt.Hash)

		pgtest.MustExec(t, db, `UPDATE eth_txes SET pipeline_task_run_id = $1, min_confirmations = $2 WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		err := ec.ResumePendingTaskRuns(context.Background(), &head)
		require.NoError(t, err)

	})

	t.Run("processes eth_txes with receipts older than minConfirmations", func(t *testing.T) {
		ch := make(chan interface{})
		ec := cltest.NewEthConfirmer(t, db, ethClient, evmcfg, ethKeyStore, []ethkey.State{state}, func(id uuid.UUID, value interface{}, err error) error {
			require.Nil(t, err)
			ch <- value
			return nil
		})

		run := cltest.MustInsertPipelineRun(t, db)
		tr := cltest.MustInsertUnfinishedPipelineTaskRun(t, db, run.ID)
		pgtest.MustExec(t, db, `UPDATE pipeline_runs SET state = 'suspended' WHERE id = $1`, run.ID)

		etx := cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 3, 1, fromAddress)
		attempt := etx.EthTxAttempts[0]
		receipt := cltest.MustInsertEthReceipt(t, borm, head.Number-minConfirmations, head.Hash, attempt.Hash)

		pgtest.MustExec(t, db, `UPDATE eth_txes SET pipeline_task_run_id = $1, min_confirmations = $2 WHERE id = $3`, &tr.ID, minConfirmations, etx.ID)

		go func() {
			err := ec.ResumePendingTaskRuns(context.Background(), &head)
			require.NoError(t, err)
		}()

		select {
		case data := <-ch:
			require.IsType(t, []byte{}, data)

			var r evmtypes.Receipt
			err := json.Unmarshal(data.([]byte), &r)
			require.NoError(t, err)
			require.Equal(t, receipt.TxHash, r.TxHash)

		case <-time.After(time.Second):
			t.Fatal("no value received")
		}
	})

}
