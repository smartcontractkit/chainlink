package bulletprooftxmanager_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	ksmocks "github.com/smartcontractkit/chainlink/core/services/keystore/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func mustInsertUnstartedEthTx(t *testing.T, db *gorm.DB, fromAddress gethCommon.Address) {
	etx := cltest.NewEthTx(t, fromAddress)
	etx.State = bulletprooftxmanager.EthTxUnstarted
	require.NoError(t, db.Save(&etx).Error)
}

func newBroadcastEthTxAttempt(t *testing.T, etxID int64, gasPrice ...int64) bulletprooftxmanager.EthTxAttempt {
	attempt := cltest.NewEthTxAttempt(t, etxID)
	attempt.State = bulletprooftxmanager.EthTxAttemptBroadcast
	if len(gasPrice) > 0 {
		gp := gasPrice[0]
		attempt.GasPrice = *utils.NewBig(big.NewInt(gp))
	}
	return attempt
}

func newInProgressEthTxAttempt(t *testing.T, etxID int64, gasPrice ...int64) bulletprooftxmanager.EthTxAttempt {
	attempt := cltest.NewEthTxAttempt(t, etxID)
	attempt.State = bulletprooftxmanager.EthTxAttemptInProgress
	if len(gasPrice) > 0 {
		gp := gasPrice[0]
		attempt.GasPrice = *utils.NewBig(big.NewInt(gp))
	}
	return attempt
}

func mustInsertInProgressEthTx(t *testing.T, db *gorm.DB, nonce int64, fromAddress gethCommon.Address) bulletprooftxmanager.EthTx {
	etx := cltest.NewEthTx(t, fromAddress)
	etx.State = bulletprooftxmanager.EthTxInProgress
	etx.Nonce = &nonce
	require.NoError(t, db.Save(&etx).Error)

	return etx
}

func mustInsertConfirmedEthTx(t *testing.T, db *gorm.DB, nonce int64, fromAddress gethCommon.Address) bulletprooftxmanager.EthTx {
	etx := cltest.NewEthTx(t, fromAddress)
	etx.State = bulletprooftxmanager.EthTxConfirmed
	etx.Nonce = &nonce
	now := time.Now()
	etx.BroadcastAt = &now
	require.NoError(t, db.Save(&etx).Error)

	return etx
}

func TestEthConfirmer_SetBroadcastBeforeBlockNum(t *testing.T) {
	t.Parallel()
	db := pgtest.NewGormDB(t)

	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	ethClient := new(mocks.Client)

	key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	ec := cltest.NewEthConfirmer(t, db, ethClient, config, ethKeyStore, []ethkey.Key{key})

	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 0, fromAddress)

	headNum := int64(9000)
	var err error

	t.Run("saves block num to unconfirmed eth_tx_attempts without one", func(t *testing.T) {
		// Do the thing
		require.NoError(t, ec.SetBroadcastBeforeBlockNum(headNum))

		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]

		assert.Equal(t, int64(9000), *attempt.BroadcastBeforeBlockNum)
	})

	t.Run("does not change eth_tx_attempts that already have BroadcastBeforeBlockNum set", func(t *testing.T) {
		n := int64(42)
		attempt := newBroadcastEthTxAttempt(t, etx.ID, 2)
		attempt.BroadcastBeforeBlockNum = &n
		require.NoError(t, db.Save(&attempt).Error)

		// Do the thing
		require.NoError(t, ec.SetBroadcastBeforeBlockNum(headNum))

		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.EthTxAttempts, 2)
		attempt = etx.EthTxAttempts[1]

		assert.Equal(t, int64(42), *attempt.BroadcastBeforeBlockNum)
	})
}

func TestEthConfirmer_CheckForReceipts(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB
	ethClient := new(mocks.Client)
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, []ethkey.Key{key})

	nonce := int64(0)
	ctx := context.Background()
	blockNum := int64(0)

	t.Run("only finds eth_txes in unconfirmed state with at least one broadcast attempt", func(t *testing.T) {
		cltest.MustInsertFatalErrorEthTx(t, db, fromAddress)
		mustInsertInProgressEthTx(t, db, nonce, fromAddress)
		nonce++
		cltest.MustInsertConfirmedEthTxWithAttempt(t, db, nonce, 1, fromAddress)
		nonce++
		cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, db, nonce, fromAddress)
		nonce++
		mustInsertUnstartedEthTx(t, db, fromAddress)

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))
		// No calls
		ethClient.AssertExpectations(t)
	})

	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, nonce, fromAddress)
	nonce++
	require.Len(t, etx1.EthTxAttempts, 1)
	attempt1_1 := etx1.EthTxAttempts[0]
	require.Len(t, attempt1_1.EthReceipts, 0)

	t.Run("fetches receipt for one unconfirmed eth_tx", func(t *testing.T) {
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		// Transaction not confirmed yet, receipt is nil
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesHash(b[0], attempt1_1.Hash)
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &bulletprooftxmanager.Receipt{}
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		var err error
		etx1, err = cltest.FindEthTxWithAttempts(db, etx1.ID)
		require.Len(t, etx1.EthTxAttempts, 1)
		attempt1_1 = etx1.EthTxAttempts[0]
		require.NoError(t, err)
		require.Len(t, attempt1_1.EthReceipts, 0)

		ethClient.AssertExpectations(t)
	})

	t.Run("saves nothing if returned receipt does not match the attempt", func(t *testing.T) {
		bptxmReceipt := bulletprooftxmanager.Receipt{
			TxHash:           utils.NewHash(),
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}

		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		// First transaction confirmed
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesHash(b[0], attempt1_1.Hash)
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &bptxmReceipt
		}).Once()

		// No error because it is merely logged
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		etx, err := cltest.FindEthTxWithAttempts(db, etx1.ID)
		require.NoError(t, err)
		require.Len(t, etx.EthTxAttempts, 1)

		require.Len(t, etx.EthTxAttempts[0].EthReceipts, 0)
	})

	t.Run("saves nothing if query returns error", func(t *testing.T) {
		bptxmReceipt := bulletprooftxmanager.Receipt{
			TxHash:           attempt1_1.Hash,
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}

		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		// First transaction confirmed
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesHash(b[0], attempt1_1.Hash)
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &bptxmReceipt
			elems[0].Error = errors.New("foo")
		}).Once()

		// No error because it is merely logged
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		etx, err := cltest.FindEthTxWithAttempts(db, etx1.ID)
		require.NoError(t, err)
		require.Len(t, etx.EthTxAttempts, 1)
		require.Len(t, etx.EthTxAttempts[0].EthReceipts, 0)
	})

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, nonce, fromAddress)
	nonce++
	require.Len(t, etx2.EthTxAttempts, 1)
	attempt2_1 := etx2.EthTxAttempts[0]
	require.Len(t, attempt2_1.EthReceipts, 0)

	t.Run("saves eth_receipt and marks eth_tx as confirmed when geth client returns valid receipt", func(t *testing.T) {
		bptxmReceipt := bulletprooftxmanager.Receipt{
			TxHash:           attempt1_1.Hash,
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}

		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				cltest.BatchElemMatchesHash(b[0], attempt1_1.Hash) &&
				cltest.BatchElemMatchesHash(b[1], attempt2_1.Hash)

		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// First transaction confirmed
			elems[0].Result = &bptxmReceipt
			// Second transaction still unconfirmed
			elems[1].Result = &bulletprooftxmanager.Receipt{}
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		// Check that the receipt was saved
		etx, err := cltest.FindEthTxWithAttempts(db, etx1.ID)
		require.NoError(t, err)

		assert.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt1_1 = etx.EthTxAttempts[0]
		require.Len(t, attempt1_1.EthReceipts, 1)

		ethReceipt := attempt1_1.EthReceipts[0]

		assert.Equal(t, bptxmReceipt.TxHash, ethReceipt.TxHash)
		assert.Equal(t, bptxmReceipt.BlockHash, ethReceipt.BlockHash)
		assert.Equal(t, bptxmReceipt.BlockNumber.Int64(), ethReceipt.BlockNumber)
		assert.Equal(t, bptxmReceipt.TransactionIndex, ethReceipt.TransactionIndex)

		receiptJSON, err := json.Marshal(bptxmReceipt)
		require.NoError(t, err)

		assert.JSONEq(t, string(receiptJSON), string(ethReceipt.Receipt))

		ethClient.AssertExpectations(t)
	})

	t.Run("fetches and saves receipts for several attempts in gas price order", func(t *testing.T) {
		attempt2_2 := newBroadcastEthTxAttempt(t, etx2.ID)
		attempt2_2.GasPrice = *utils.NewBig(big.NewInt(10))

		attempt2_3 := newBroadcastEthTxAttempt(t, etx2.ID)
		attempt2_3.GasPrice = *utils.NewBig(big.NewInt(20))

		// Insert order deliberately reversed to test sorting by gas price
		require.NoError(t, store.DB.Create(&attempt2_3).Error)
		require.NoError(t, store.DB.Create(&attempt2_2).Error)

		bptxmReceipt := bulletprooftxmanager.Receipt{
			TxHash:           attempt2_2.Hash,
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}

		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 3 &&
				cltest.BatchElemMatchesHash(b[2], attempt2_1.Hash) &&
				cltest.BatchElemMatchesHash(b[1], attempt2_2.Hash) &&
				cltest.BatchElemMatchesHash(b[0], attempt2_3.Hash)

		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// Most expensive attempt still unconfirmed
			elems[2].Result = &bulletprooftxmanager.Receipt{}
			// Second most expensive attempt is confirmed
			elems[1].Result = &bptxmReceipt
			// Cheapest attempt still unconfirmed
			elems[0].Result = &bulletprooftxmanager.Receipt{}
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		ethClient.AssertExpectations(t)

		// Check that the state was updated
		etx, err := cltest.FindEthTxWithAttempts(db, etx2.ID)
		require.NoError(t, err)

		require.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 3)
	})

	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, nonce, fromAddress)
	attempt3_1 := etx3.EthTxAttempts[0]
	nonce++

	t.Run("ignores receipt missing BlockHash that comes from querying parity too early", func(t *testing.T) {
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		receipt := bulletprooftxmanager.Receipt{
			TxHash: attempt3_1.Hash,
		}
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesHash(b[0], attempt3_1.Hash)
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &receipt
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		// No receipt, but no error either
		etx, err := cltest.FindEthTxWithAttempts(db, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt3_1 = etx.EthTxAttempts[0]
		require.Len(t, attempt3_1.EthReceipts, 0)
	})

	t.Run("does not panic if receipt has BlockHash but is missing some other fields somehow", func(t *testing.T) {
		// NOTE: This should never happen, but we shouldn't panic regardless
		receipt := bulletprooftxmanager.Receipt{
			TxHash:    attempt3_1.Hash,
			BlockHash: utils.NewHash(),
		}
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesHash(b[0], attempt3_1.Hash)
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &receipt
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		// No receipt, but no error either
		etx, err := cltest.FindEthTxWithAttempts(db, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt3_1 = etx.EthTxAttempts[0]
		require.Len(t, attempt3_1.EthReceipts, 0)
	})

	t.Run("handles case where eth_receipt already exists somehow", func(t *testing.T) {
		ethReceipt := cltest.MustInsertEthReceipt(t, db, 42, utils.NewHash(), attempt3_1.Hash)

		bptxmReceipt := bulletprooftxmanager.Receipt{
			TxHash:           attempt3_1.Hash,
			BlockHash:        ethReceipt.BlockHash,
			BlockNumber:      big.NewInt(ethReceipt.BlockNumber),
			TransactionIndex: ethReceipt.TransactionIndex,
		}
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 1 && cltest.BatchElemMatchesHash(b[0], attempt3_1.Hash)
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			elems[0].Result = &bptxmReceipt
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		// Check that the receipt was unchanged
		etx, err := cltest.FindEthTxWithAttempts(db, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt3_1 = etx.EthTxAttempts[0]
		require.Len(t, attempt3_1.EthReceipts, 1)

		ethReceipt = attempt3_1.EthReceipts[0]

		assert.Equal(t, bptxmReceipt.TxHash, ethReceipt.TxHash)
		assert.Equal(t, bptxmReceipt.BlockHash, ethReceipt.BlockHash)
		assert.Equal(t, bptxmReceipt.BlockNumber.Int64(), ethReceipt.BlockNumber)
		assert.Equal(t, bptxmReceipt.TransactionIndex, ethReceipt.TransactionIndex)

		ethClient.AssertExpectations(t)
	})

	etx4 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, nonce, fromAddress)
	attempt4_1 := etx4.EthTxAttempts[0]
	nonce++

	t.Run("on receipt fetch marks in_progress eth_tx_attempt as broadcast", func(t *testing.T) {
		attempt4_2 := newInProgressEthTxAttempt(t, etx4.ID)
		attempt4_2.GasPrice = *utils.NewBig(big.NewInt(10))

		require.NoError(t, store.DB.Create(&attempt4_2).Error)

		bptxmReceipt := bulletprooftxmanager.Receipt{
			TxHash:           attempt4_2.Hash,
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		// Second attempt is confirmed
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 2 &&
				cltest.BatchElemMatchesHash(b[0], attempt4_2.Hash) &&
				cltest.BatchElemMatchesHash(b[1], attempt4_1.Hash)
		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// First attempt still unconfirmed
			elems[1].Result = &bulletprooftxmanager.Receipt{}
			// Second attempt is confirmed
			elems[0].Result = &bptxmReceipt
		}).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts(ctx, blockNum))

		ethClient.AssertExpectations(t)

		// Check that the state was updated
		var err error
		etx4, err = cltest.FindEthTxWithAttempts(db, etx4.ID)
		require.NoError(t, err)

		attempt4_1 = etx4.EthTxAttempts[0]
		attempt4_2 = etx4.EthTxAttempts[1]

		// And the attempts
		require.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt4_1.State)
		require.Nil(t, attempt4_1.BroadcastBeforeBlockNum)
		require.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt4_2.State)
		require.Equal(t, int64(42), *attempt4_2.BroadcastBeforeBlockNum)

		// Check receipts
		require.Len(t, attempt4_1.EthReceipts, 0)
		require.Len(t, attempt4_2.EthReceipts, 1)
	})
}

func TestEthConfirmer_CheckForReceipts_batching(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

	ethClient := new(mocks.Client)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	config.Set("ETH_RPC_DEFAULT_BATCH_SIZE", 2)

	ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, []ethkey.Key{key})

	ctx := context.Background()

	etx := cltest.MustInsertUnconfirmedEthTx(t, db, 0, fromAddress)
	var attempts []bulletprooftxmanager.EthTxAttempt

	// Total of 5 attempts should lead to 3 batched fetches (2, 2, 1)
	for i := 0; i < 5; i++ {
		attempt := newBroadcastEthTxAttempt(t, etx.ID, int64(i+2))
		require.NoError(t, store.DB.Create(&attempt).Error)
		attempts = append(attempts, attempt)
	}

	ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)

	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 2 &&
			cltest.BatchElemMatchesHash(b[0], attempts[4].Hash) &&
			cltest.BatchElemMatchesHash(b[1], attempts[3].Hash)
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &bulletprooftxmanager.Receipt{}
		elems[1].Result = &bulletprooftxmanager.Receipt{}
	}).Once()
	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 2 &&
			cltest.BatchElemMatchesHash(b[0], attempts[2].Hash) &&
			cltest.BatchElemMatchesHash(b[1], attempts[1].Hash)
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &bulletprooftxmanager.Receipt{}
		elems[1].Result = &bulletprooftxmanager.Receipt{}
	}).Once()
	ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
		return len(b) == 1 &&
			cltest.BatchElemMatchesHash(b[0], attempts[0].Hash)
	})).Return(nil).Run(func(args mock.Arguments) {
		elems := args.Get(1).([]rpc.BatchElem)
		elems[0].Result = &bulletprooftxmanager.Receipt{}
	}).Once()

	require.NoError(t, ec.CheckForReceipts(ctx, 42))
	ethClient.AssertExpectations(t)
}

func TestEthConfirmer_CheckForReceipts_only_likely_confirmed(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

	ethClient := new(mocks.Client)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	config.Set("ETH_RPC_DEFAULT_BATCH_SIZE", 6)

	ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, []ethkey.Key{key})

	ctx := context.Background()

	var attempts []bulletprooftxmanager.EthTxAttempt
	// inserting in DESC nonce order to test DB ASC ordering
	etx2 := cltest.MustInsertUnconfirmedEthTx(t, db, 1, fromAddress)
	for i := 0; i < 4; i++ {
		attempt := newBroadcastEthTxAttempt(t, etx2.ID, int64(100-i))
		require.NoError(t, store.DB.Create(&attempt).Error)
	}
	etx := cltest.MustInsertUnconfirmedEthTx(t, db, 0, fromAddress)
	for i := 0; i < 4; i++ {
		attempt := newBroadcastEthTxAttempt(t, etx.ID, int64(100-i))
		require.NoError(t, store.DB.Create(&attempt).Error)

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
		elems[0].Result = &bulletprooftxmanager.Receipt{}
		elems[1].Result = &bulletprooftxmanager.Receipt{}
		elems[2].Result = &bulletprooftxmanager.Receipt{}
		elems[3].Result = &bulletprooftxmanager.Receipt{}
	}).Once()

	require.NoError(t, ec.CheckForReceipts(ctx, 42))

	cltest.BatchElemMustMatchHash(t, captured[0], attempts[0].Hash)
	cltest.BatchElemMustMatchHash(t, captured[1], attempts[1].Hash)
	cltest.BatchElemMustMatchHash(t, captured[2], attempts[2].Hash)
	cltest.BatchElemMustMatchHash(t, captured[3], attempts[3].Hash)

	ethClient.AssertExpectations(t)
}

func TestEthConfirmer_CheckForReceipts_should_not_check_for_likely_unconfirmed(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB

	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

	ethClient := new(mocks.Client)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, []ethkey.Key{key})

	ctx := context.Background()

	etx := cltest.MustInsertUnconfirmedEthTx(t, db, 1, fromAddress)
	for i := 0; i < 4; i++ {
		attempt := newBroadcastEthTxAttempt(t, etx.ID, int64(100-i))
		require.NoError(t, store.DB.Create(&attempt).Error)
	}

	// latest nonce is lower that all attempts' nonces
	ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(0), nil)

	require.NoError(t, ec.CheckForReceipts(ctx, 42))

	// no BatchCallContext calls
	ethClient.AssertExpectations(t)
}

func TestEthConfirmer_CheckForReceipts_confirmed_missing_receipt(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

	ethClient := new(mocks.Client)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	config.Set("ETH_FINALITY_DEPTH", 50)
	ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, []ethkey.Key{key})

	ctx := context.Background()

	// STATE
	// eth_txes with nonce 0 has two attempts (broadcast before block 21 and 41) the first of which will get a receipt
	// eth_txes with nonce 1 has two attempts (broadcast before block 21 and 41) neither of which will ever get a receipt
	// eth_txes with nonce 2 has an attempt (broadcast before block 41) that will not get a receipt on the first try but will get one later
	// eth_txes with nonce 3 has an attempt (broadcast before block 41) that has been confirmed in block 42
	// All other attempts were broadcast before block 41
	b := int64(21)

	etx0 := cltest.MustInsertUnconfirmedEthTx(t, db, 0, fromAddress)
	attempt0_1 := newBroadcastEthTxAttempt(t, etx0.ID, int64(1))
	attempt0_2 := newBroadcastEthTxAttempt(t, etx0.ID, int64(2))
	attempt0_2.BroadcastBeforeBlockNum = &b
	require.NoError(t, store.DB.Create(&attempt0_1).Error)
	require.NoError(t, store.DB.Create(&attempt0_2).Error)

	etx1 := cltest.MustInsertUnconfirmedEthTx(t, db, 1, fromAddress)
	attempt1_1 := newBroadcastEthTxAttempt(t, etx1.ID, int64(1))
	attempt1_2 := newBroadcastEthTxAttempt(t, etx1.ID, int64(2))
	attempt1_2.BroadcastBeforeBlockNum = &b
	require.NoError(t, store.DB.Create(&attempt1_1).Error)
	require.NoError(t, store.DB.Create(&attempt1_2).Error)

	etx2 := cltest.MustInsertUnconfirmedEthTx(t, db, 2, fromAddress)
	attempt2_1 := newBroadcastEthTxAttempt(t, etx2.ID, int64(1))
	require.NoError(t, store.DB.Create(&attempt2_1).Error)

	etx3 := cltest.MustInsertUnconfirmedEthTx(t, db, 3, fromAddress)
	attempt3_1 := newBroadcastEthTxAttempt(t, etx3.ID, int64(1))
	require.NoError(t, store.DB.Create(&attempt3_1).Error)

	d, _ := store.ORM.DB.DB()
	_, err := d.Exec(`UPDATE eth_tx_attempts SET broadcast_before_block_num = 41 WHERE broadcast_before_block_num IS NULL`)
	require.NoError(t, err)

	t.Run("marks buried eth_txes as 'confirmed_missing_receipt'", func(t *testing.T) {
		bptxmReceipt0 := bulletprooftxmanager.Receipt{
			TxHash:           attempt0_2.Hash,
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}
		bptxmReceipt3 := bulletprooftxmanager.Receipt{
			TxHash:           attempt3_1.Hash,
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(4), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 6 &&
				cltest.BatchElemMatchesHash(b[0], attempt0_2.Hash) &&
				cltest.BatchElemMatchesHash(b[1], attempt0_1.Hash) &&
				cltest.BatchElemMatchesHash(b[2], attempt1_2.Hash) &&
				cltest.BatchElemMatchesHash(b[3], attempt1_1.Hash) &&
				cltest.BatchElemMatchesHash(b[4], attempt2_1.Hash) &&
				cltest.BatchElemMatchesHash(b[5], attempt3_1.Hash)

		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// First transaction confirmed
			elems[0].Result = &bptxmReceipt0
			elems[1].Result = &bulletprooftxmanager.Receipt{}
			// Second transaction stil unconfirmed
			elems[2].Result = &bulletprooftxmanager.Receipt{}
			elems[3].Result = &bulletprooftxmanager.Receipt{}
			// Third transaction still unconfirmed
			elems[4].Result = &bulletprooftxmanager.Receipt{}
			// Fourth transaction is confirmed
			elems[5].Result = &bptxmReceipt3
		}).Once()

		// PERFORM
		// Block num of 43 is one higher than the receipt (as would generally be expected)
		require.NoError(t, ec.CheckForReceipts(ctx, 43))

		ethClient.AssertExpectations(t)

		// Expected state is that the "top" eth_tx is now confirmed, with the
		// two below it "confirmed_missing_receipt" and the "bottom" eth_tx also confirmed
		etx3, err = cltest.FindEthTxWithAttempts(db, etx3.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx3.State)

		ethReceipt := etx3.EthTxAttempts[0].EthReceipts[0]
		require.Equal(t, bptxmReceipt3.BlockHash, ethReceipt.BlockHash)

		etx2, err = cltest.FindEthTxWithAttempts(db, etx2.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmedMissingReceipt, etx2.State)
		etx1, err = cltest.FindEthTxWithAttempts(db, etx1.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmedMissingReceipt, etx1.State)

		etx0, err = cltest.FindEthTxWithAttempts(db, etx0.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx0.State)

		ethReceipt = etx0.EthTxAttempts[1].EthReceipts[0]
		require.Equal(t, bptxmReceipt0.BlockHash, ethReceipt.BlockHash)
	})

	// STATE
	// eth_txes with nonce 0 is confirmed
	// eth_txes with nonce 1 is confirmed_missing_receipt
	// eth_txes with nonce 2 is confirmed_missing_receipt
	// eth_txes with nonce 3 is confirmed

	t.Run("marks eth_txes with state 'confirmed_missing_receipt' as 'confirmed' if a receipt finally shows up", func(t *testing.T) {
		bptxmReceipt := bulletprooftxmanager.Receipt{
			TxHash:           attempt2_1.Hash,
			BlockHash:        utils.NewHash(),
			BlockNumber:      big.NewInt(43),
			TransactionIndex: uint(1),
		}
		ethClient.On("NonceAt", mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
		ethClient.On("BatchCallContext", mock.Anything, mock.MatchedBy(func(b []rpc.BatchElem) bool {
			return len(b) == 3 &&
				cltest.BatchElemMatchesHash(b[0], attempt1_2.Hash) &&
				cltest.BatchElemMatchesHash(b[1], attempt1_1.Hash) &&
				cltest.BatchElemMatchesHash(b[2], attempt2_1.Hash)

		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// First transaction still unconfirmed
			elems[0].Result = &bulletprooftxmanager.Receipt{}
			elems[1].Result = &bulletprooftxmanager.Receipt{}
			// Second transaction confirmed
			elems[2].Result = &bptxmReceipt
		}).Once()

		// PERFORM
		// Block num of 44 is one higher than the receipt (as would generally be expected)
		require.NoError(t, ec.CheckForReceipts(ctx, 44))

		ethClient.AssertExpectations(t)

		// Expected state is that the "top" two eth_txes are now confirmed, with the
		// one below it still "confirmed_missing_receipt" and the bottom one remains confirmed
		etx3, err = cltest.FindEthTxWithAttempts(db, etx3.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx3.State)
		etx2, err = cltest.FindEthTxWithAttempts(db, etx2.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx2.State)

		ethReceipt := etx2.EthTxAttempts[0].EthReceipts[0]
		require.Equal(t, bptxmReceipt.BlockHash, ethReceipt.BlockHash)

		etx1, err = cltest.FindEthTxWithAttempts(db, etx1.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmedMissingReceipt, etx1.State)
		etx0, err = cltest.FindEthTxWithAttempts(db, etx0.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx0.State)
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
				cltest.BatchElemMatchesHash(b[0], attempt1_2.Hash) &&
				cltest.BatchElemMatchesHash(b[1], attempt1_1.Hash)

		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// Both attempts still unconfirmed
			elems[0].Result = &bulletprooftxmanager.Receipt{}
			elems[1].Result = &bulletprooftxmanager.Receipt{}
		}).Once()

		// PERFORM
		// Block num of 80 puts the first attempt (21) below threshold but second attempt (41) still above
		require.NoError(t, ec.CheckForReceipts(ctx, 80))

		ethClient.AssertExpectations(t)

		// Expected state is that the "top" two eth_txes are now confirmed, with the
		// one below it still "confirmed_missing_receipt" and the bottom one remains confirmed
		etx3, err = cltest.FindEthTxWithAttempts(db, etx3.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx3.State)
		etx2, err = cltest.FindEthTxWithAttempts(db, etx2.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx2.State)
		etx1, err = cltest.FindEthTxWithAttempts(db, etx1.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmedMissingReceipt, etx1.State)
		etx0, err = cltest.FindEthTxWithAttempts(db, etx0.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx0.State)
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
				cltest.BatchElemMatchesHash(b[0], attempt1_2.Hash) &&
				cltest.BatchElemMatchesHash(b[1], attempt1_1.Hash)

		})).Return(nil).Run(func(args mock.Arguments) {
			elems := args.Get(1).([]rpc.BatchElem)
			// Both attempts still unconfirmed
			elems[0].Result = &bulletprooftxmanager.Receipt{}
			elems[1].Result = &bulletprooftxmanager.Receipt{}
		}).Once()

		// PERFORM
		// Block num of 100 puts the first attempt (21) and second attempt (41) below threshold
		require.NoError(t, ec.CheckForReceipts(ctx, 100))

		ethClient.AssertExpectations(t)

		// Expected state is that the "top" two eth_txes are now confirmed, with the
		// one below it marked as "fatal_error" and the bottom one remains confirmed
		etx3, err = cltest.FindEthTxWithAttempts(db, etx3.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx3.State)
		etx2, err = cltest.FindEthTxWithAttempts(db, etx2.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx2.State)
		etx1, err = cltest.FindEthTxWithAttempts(db, etx1.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxFatalError, etx1.State)
		etx0, err = cltest.FindEthTxWithAttempts(db, etx0.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx0.State)
	})
}

func TestEthConfirmer_FindEthTxsRequiringResubmissionDueToInsufficientEth(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)
	_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

	// Insert order is mixed up to test sorting
	etx2 := cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, db, 1, fromAddress)
	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 2, fromAddress)
	attempt3_2 := cltest.NewEthTxAttempt(t, etx3.ID)
	attempt3_2.State = bulletprooftxmanager.EthTxAttemptInsufficientEth
	attempt3_2.GasPrice = *utils.NewBig(big.NewInt(100))
	require.NoError(t, store.DB.Save(&attempt3_2).Error)
	etx1 := cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, db, 0, fromAddress)

	// These should never be returned
	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 3, fromAddress)
	cltest.MustInsertConfirmedEthTxWithAttempt(t, db, 4, 100, fromAddress)
	cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, db, 0, otherAddress)

	t.Run("returns all eth_txes with at least one attempt that is in insufficient_eth state", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringResubmissionDueToInsufficientEth(store.DB, fromAddress)
		require.NoError(t, err)

		assert.Len(t, etxs, 3)

		assert.Equal(t, *etx1.Nonce, *etxs[0].Nonce)
		assert.Equal(t, etx1.ID, etxs[0].ID)
		assert.Equal(t, *etx2.Nonce, *etxs[1].Nonce)
		assert.Equal(t, etx2.ID, etxs[1].ID)
		assert.Equal(t, *etx3.Nonce, *etxs[2].Nonce)
		assert.Equal(t, etx3.ID, etxs[2].ID)
	})

	t.Run("does not return confirmed or fatally errored eth_txes", func(t *testing.T) {
		require.NoError(t, store.DB.Exec(`UPDATE eth_txes SET state='confirmed' WHERE id = ?`, etx1.ID).Error)
		require.NoError(t, store.DB.Exec(`UPDATE eth_txes SET state='fatal_error', nonce=NULL, error='foo', broadcast_at=NULL WHERE id = ?`, etx2.ID).Error)

		etxs, err := bulletprooftxmanager.FindEthTxsRequiringResubmissionDueToInsufficientEth(store.DB, fromAddress)
		require.NoError(t, err)

		assert.Len(t, etxs, 1)

		assert.Equal(t, *etx3.Nonce, *etxs[0].Nonce)
		assert.Equal(t, etx3.ID, etxs[0].ID)
	})
}

func TestEthConfirmer_FindEthTxsRequiringRebroadcast(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

	currentHead := int64(30)
	gasBumpThreshold := int64(10)
	tooNew := int64(21)
	onTheMoney := int64(20)
	oldEnough := int64(19)
	nonce := int64(0)

	mustInsertConfirmedEthTx(t, db, nonce, fromAddress)
	nonce++

	_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

	t.Run("returns nothing when there are no transactions", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 10, 0)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	mustInsertInProgressEthTx(t, db, nonce, fromAddress)
	nonce++

	t.Run("returns nothing when the transaction is in_progress", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 10, 0)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	// This one has BroadcastBeforeBlockNum set as nil... which can happen, but it should be ignored
	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, nonce, fromAddress)
	nonce++

	t.Run("ignores unconfirmed transactions with nil BroadcastBeforeBlockNum", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 10, 0)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, nonce, fromAddress)
	nonce++
	attempt1_1 := etx1.EthTxAttempts[0]
	attempt1_1.BroadcastBeforeBlockNum = &tooNew
	require.NoError(t, store.DB.Save(&attempt1_1).Error)
	attempt1_2 := newBroadcastEthTxAttempt(t, etx1.ID)
	attempt1_2.BroadcastBeforeBlockNum = &onTheMoney
	attempt1_2.GasPrice = *utils.NewBigI(30000)
	require.NoError(t, store.DB.Save(&attempt1_2).Error)

	t.Run("returns nothing when the transaction is unconfirmed with an attempt that is recent", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 10, 0)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, nonce, fromAddress)
	nonce++
	attempt2_1 := etx2.EthTxAttempts[0]
	attempt2_1.BroadcastBeforeBlockNum = &tooNew
	require.NoError(t, store.DB.Save(&attempt2_1).Error)

	t.Run("returns nothing when the transaction has attempts that are too new", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 10, 0)
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
	etxWithoutAttempts.State = bulletprooftxmanager.EthTxUnconfirmed
	require.NoError(t, store.DB.Save(&etxWithoutAttempts).Error)
	nonce++

	t.Run("does nothing if the transaction is from a different address than the one given", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, otherAddress, currentHead, gasBumpThreshold, 10, 0)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	t.Run("returns the transaction if it is unconfirmed and has no attempts (note that this is an invariant violation, but we handle it anyway)", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 10, 0)
		require.NoError(t, err)

		require.Len(t, etxs, 1)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
	})

	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, nonce, fromAddress)
	nonce++
	attempt3_1 := etx3.EthTxAttempts[0]
	attempt3_1.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt3_1).Error)

	// NOTE: It should ignore qualifying eth_txes from a different address
	etxOther := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 0, otherAddress)
	attemptOther1 := etxOther.EthTxAttempts[0]
	attemptOther1.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attemptOther1).Error)

	t.Run("returns the transaction if it is unconfirmed with an attempt that is older than gasBumpThreshold blocks", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 10, 0)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
	})

	t.Run("returns nothing if threshold is zero", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, 0, 10, 0)
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
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 4, 0)
		require.NoError(t, err)

		require.Len(t, etxs, 1) // returns etxWithoutAttempts only - eligible for gas bumping because it technically doesn't have any attempts withing gasBumpThreshold blocks
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)

		etxs, err = bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 5, 0)
		require.NoError(t, err)

		require.Len(t, etxs, 2) // includes etxWithoutAttempts, etx3 and etx4
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)

		// Zero limit disables it
		etxs, err = bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 0, 0)
		require.NoError(t, err)

		require.Len(t, etxs, 2) // includes etxWithoutAttempts, etx3 and etx4
	})

	etx4 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, nonce, fromAddress)
	nonce++
	attempt4_1 := etx4.EthTxAttempts[0]
	attempt4_1.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt4_1).Error)

	t.Run("ignores pending transactions for another key", func(t *testing.T) {
		// Re-use etx3 nonce for another key, it should not affect the results for this key
		etxOther := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, *etx3.Nonce, otherAddress)
		aOther := etxOther.EthTxAttempts[0]
		aOther.BroadcastBeforeBlockNum = &oldEnough
		require.NoError(t, store.DB.Save(&aOther).Error)

		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 6, 0)
		require.NoError(t, err)

		require.Len(t, etxs, 3) // includes etxWithoutAttempts, etx3 and etx4
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
		assert.Equal(t, etx4.ID, etxs[2].ID)
	})

	attempt3_2 := newBroadcastEthTxAttempt(t, etx3.ID)
	attempt3_2.BroadcastBeforeBlockNum = &oldEnough
	attempt3_2.GasPrice = *utils.NewBigI(30000)
	require.NoError(t, store.DB.Save(&attempt3_2).Error)

	t.Run("returns the transaction if it is unconfirmed with two attempts that are older than gasBumpThreshold blocks", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 10, 0)
		require.NoError(t, err)

		require.Len(t, etxs, 3)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
		assert.Equal(t, etx4.ID, etxs[2].ID)
	})

	attempt3_3 := newBroadcastEthTxAttempt(t, etx3.ID)
	attempt3_3.BroadcastBeforeBlockNum = &tooNew
	attempt3_3.GasPrice = *utils.NewBigI(40000)
	require.NoError(t, store.DB.Save(&attempt3_3).Error)

	t.Run("does not return the transaction if it has some older but one newer attempt", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 10, 0)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, *etxWithoutAttempts.Nonce, *(etxs[0].Nonce))
		require.Equal(t, int64(5), *etxWithoutAttempts.Nonce)
		assert.Equal(t, etx4.ID, etxs[1].ID)
		assert.Equal(t, *etx4.Nonce, *(etxs[1].Nonce))
		require.Equal(t, int64(7), *etx4.Nonce)
	})

	attempt0_1 := newBroadcastEthTxAttempt(t, etxWithoutAttempts.ID)
	attempt0_1.State = bulletprooftxmanager.EthTxAttemptInsufficientEth
	require.NoError(t, store.DB.Save(&attempt0_1).Error)

	// This attempt has insufficient_eth, but there is also another attempt4_1
	// which is old enough, so this will be caught by both queries and should
	// not be duplicated
	attempt4_2 := cltest.NewEthTxAttempt(t, etx4.ID)
	attempt4_2.State = bulletprooftxmanager.EthTxAttemptInsufficientEth
	attempt4_2.GasPrice = *utils.NewBigI(40000)
	require.NoError(t, store.DB.Save(&attempt4_2).Error)

	etx5 := cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, db, nonce, fromAddress)
	nonce++

	// This etx has one attempt that is too new, which would exclude it from
	// the gas bumping query, but it should still be caught by the insufficient
	// eth query
	etx6 := cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, db, nonce, fromAddress)
	attempt6_2 := newBroadcastEthTxAttempt(t, etx3.ID)
	attempt6_2.BroadcastBeforeBlockNum = &tooNew
	attempt6_2.GasPrice = *utils.NewBigI(30001)
	require.NoError(t, store.DB.Save(&attempt6_2).Error)
	nonce++

	t.Run("returns unique attempts requiring resubmission due to insufficient eth, ordered by nonce asc", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 10, 0)
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
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringRebroadcast(store.DB, fromAddress, currentHead, gasBumpThreshold, 10, 2)
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

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB
	ethClient := new(mocks.Client)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	config.Set("ETH_MAX_GAS_PRICE_WEI", 500000000000) // 500GWei

	otherKey := cltest.MustInsertRandomKey(t, store.DB)
	key := cltest.MustInsertRandomKey(t, store.DB)
	fromAddress := key.Address.Address()
	keys := []ethkey.Key{key, otherKey}

	kst := new(ksmocks.EthKeyStoreInterface)
	// Use a mock keystore for this test
	ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, kst, keys)
	currentHead := int64(30)
	oldEnough := int64(19)
	nonce := int64(0)

	t.Run("does nothing if no transactions require bumping", func(t *testing.T) {
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))
	})

	originalBroadcastAt := time.Unix(1616509100, 0)
	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, nonce, fromAddress, originalBroadcastAt)
	nonce++
	attempt1_1 := etx.EthTxAttempts[0]
	attempt1_1.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt1_1).Error)
	var err error

	t.Run("returns on keystore error", func(t *testing.T) {
		// simulate transaction that is somehow impossible to sign
		kst.On("SignTx", fromAddress,
			mock.MatchedBy(func(tx *types.Transaction) bool {
				return tx.Nonce() == uint64(*etx.Nonce)
			}),
			mock.Anything).Return(nil, errors.New("signing error")).Once()

		// Do the thing
		err = ec.RebroadcastWhereNecessary(context.TODO(), currentHead)
		require.Error(t, err)
		require.Contains(t, err.Error(), "signing error")

		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)
		require.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx.State)

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
				return chainID.Cmp(config.ChainID()) == 0
			})).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Nonce)
		})).Return(errors.New("exceeds block gas limit")).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.EthTxAttempts, 1)

		kst.AssertExpectations(t)
		ethClient.AssertExpectations(t)
	})

	ethClient = new(mocks.Client)
	bulletprooftxmanager.SetEthClientOnEthConfirmer(ethClient, ec)

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
				return chainID.Cmp(config.ChainID()) == 0
			})).Return(&ethTx, nil).Once()

		// Once for the bumped attempt which exceeds limit
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Nonce) && tx.GasPrice().Int64() == int64(25000000000)
		})).Return(errors.New("tx fee (1.10 ether) exceeds the configured cap (1.00 ether)")).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		// Did not create an additional attempt
		require.Len(t, etx.EthTxAttempts, 1)

		// broadcast_at did not change
		require.Equal(t, etx.BroadcastAt.Unix(), originalBroadcastAt.Unix())

		kst.AssertExpectations(t)
		ethClient.AssertExpectations(t)
	})

	var attempt1_2 bulletprooftxmanager.EthTxAttempt
	ethClient = new(mocks.Client)
	bulletprooftxmanager.SetEthClientOnEthConfirmer(ethClient, ec)

	t.Run("creates new attempt with higher gas price if transaction has an attempt older than threshold", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(25000000000)
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
				return chainID.Cmp(config.ChainID()) == 0
			})).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.EthTxAttempts, 2)
		require.Equal(t, attempt1_1.ID, etx.EthTxAttempts[0].ID)

		// Got the new attempt
		attempt1_2 = etx.EthTxAttempts[1]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.GasPrice.ToInt().Int64())
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt1_2.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("does nothing if there is an attempt without BroadcastBeforeBlockNum set", func(t *testing.T) {
		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.EthTxAttempts, 2)
	})

	attempt1_2.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt1_2).Error)
	var attempt1_3 bulletprooftxmanager.EthTxAttempt

	t.Run("creates new attempt with higher gas price if transaction is already in mempool (e.g. due to previous crash before we could save the new attempt)", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(30000000000)
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
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.EthTxAttempts, 3)
		require.Equal(t, attempt1_1.ID, etx.EthTxAttempts[0].ID)
		require.Equal(t, attempt1_2.ID, etx.EthTxAttempts[1].ID)

		// Got the new attempt
		attempt1_3 = etx.EthTxAttempts[2]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_3.GasPrice.ToInt().Int64())
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt1_3.State)

		kst.AssertExpectations(t)
		ethClient.AssertExpectations(t)
	})

	attempt1_3.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt1_3).Error)
	var attempt1_4 bulletprooftxmanager.EthTxAttempt

	t.Run("saves new attempt even for transaction that has already been confirmed (nonce already used)", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(36000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_2.GasPrice.ToInt().Int64())

		ethTx := *types.NewTx(&types.LegacyTx{})
		receipt := bulletprooftxmanager.Receipt{BlockNumber: big.NewInt(40)}
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
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx.State)

		// Got the new attempt
		attempt1_4 = etx.EthTxAttempts[3]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_4.GasPrice.ToInt().Int64())

		require.Len(t, etx.EthTxAttempts, 4)
		require.Equal(t, attempt1_1.ID, etx.EthTxAttempts[0].ID)
		require.Equal(t, attempt1_2.ID, etx.EthTxAttempts[1].ID)
		require.Equal(t, attempt1_3.ID, etx.EthTxAttempts[2].ID)
		require.Equal(t, attempt1_4.ID, etx.EthTxAttempts[3].ID)
		require.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, etx.EthTxAttempts[0].State)
		require.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, etx.EthTxAttempts[1].State)
		require.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, etx.EthTxAttempts[2].State)
		require.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, etx.EthTxAttempts[3].State)

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	// Mark original tx as confirmed so we won't pick it up any more
	require.NoError(t, store.DB.Exec(`UPDATE eth_txes SET state = 'confirmed'`).Error)

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, nonce, fromAddress)
	nonce++
	attempt2_1 := etx2.EthTxAttempts[0]
	attempt2_1.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt2_1).Error)
	var attempt2_2 bulletprooftxmanager.EthTxAttempt

	t.Run("saves in_progress attempt on temporary error and returns error", func(t *testing.T) {
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
		})).Return(errors.New("some network error")).Once()

		// Do the thing
		err = ec.RebroadcastWhereNecessary(context.TODO(), currentHead)
		require.Error(t, err)
		require.Contains(t, err.Error(), "some network error")

		etx2, err = cltest.FindEthTxWithAttempts(db, etx2.ID)
		require.NoError(t, err)

		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx2.State)

		// Old attempt is untouched
		require.Len(t, etx2.EthTxAttempts, 2)
		require.Equal(t, attempt2_1.ID, etx2.EthTxAttempts[0].ID)
		attempt2_1 = etx2.EthTxAttempts[0]
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt2_1.State)
		assert.Equal(t, oldEnough, *attempt2_1.BroadcastBeforeBlockNum)

		// New in_progress attempt saved
		attempt2_2 = etx2.EthTxAttempts[1]
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptInProgress, attempt2_2.State)
		assert.Nil(t, attempt2_2.BroadcastBeforeBlockNum)

		// Do it again and move the attempt into "broadcast"
		n = *etx2.Nonce
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == n && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(nil).Once()

		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		// Attempt marked "broadcast"
		etx2, err = cltest.FindEthTxWithAttempts(db, etx2.ID)
		require.NoError(t, err)

		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx2.State)

		// New in_progress attempt saved
		require.Len(t, etx2.EthTxAttempts, 2)
		require.Equal(t, attempt2_2.ID, etx2.EthTxAttempts[1].ID)
		attempt2_2 = etx2.EthTxAttempts[1]
		require.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt2_2.State)
		assert.Nil(t, attempt2_2.BroadcastBeforeBlockNum)

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	// Set BroadcastBeforeBlockNum again so the next test will pick it up
	attempt2_2.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt2_2).Error)

	t.Run("assumes that 'nonce too low' error means success", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(30000000000)
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
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))
		etx2, err = cltest.FindEthTxWithAttempts(db, etx2.ID)
		require.NoError(t, err)
		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx2.State)

		// One new attempt saved
		require.Len(t, etx2.EthTxAttempts, 3)
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, etx2.EthTxAttempts[0].State)
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, etx2.EthTxAttempts[1].State)
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, etx2.EthTxAttempts[2].State)

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	// Original tx is confirmed so we won't pick it up any more
	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, nonce, fromAddress)
	nonce++
	attempt3_1 := etx3.EthTxAttempts[0]
	attempt3_1.BroadcastBeforeBlockNum = &oldEnough
	attempt3_1.GasPrice = *utils.NewBig(big.NewInt(35000000000))
	require.NoError(t, store.DB.Save(&attempt3_1).Error)

	var attempt3_2 bulletprooftxmanager.EthTxAttempt

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
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		etx3, err = cltest.FindEthTxWithAttempts(db, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx3.State)

		require.Len(t, etx3.EthTxAttempts, 2)
		require.Equal(t, attempt3_1.ID, etx3.EthTxAttempts[0].ID)
		attempt3_2 = etx3.EthTxAttempts[1]

		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt3_2.GasPrice.ToInt().Int64())

		kst.AssertExpectations(t)
		ethClient.AssertExpectations(t)
	})

	attempt3_2.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt3_2).Error)
	var attempt3_3 bulletprooftxmanager.EthTxAttempt

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
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		etx3, err = cltest.FindEthTxWithAttempts(db, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx3.State)

		require.Len(t, etx3.EthTxAttempts, 3)
		attempt3_3 = etx3.EthTxAttempts[2]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt3_3.GasPrice.ToInt().Int64())

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	attempt3_3.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt3_3).Error)
	var attempt3_4 bulletprooftxmanager.EthTxAttempt

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
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		etx3, err = cltest.FindEthTxWithAttempts(db, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx3.State)

		require.Len(t, etx3.EthTxAttempts, 4)
		attempt3_4 = etx3.EthTxAttempts[3]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt3_4.GasPrice.ToInt().Int64())

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	attempt3_4.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt3_4).Error)

	t.Run("resubmits at the old price and does not create a new attempt if one of the bumped transactions would exceed ETH_MAX_GAS_PRICE_WEI", func(t *testing.T) {
		// Set price such that the next bump will exceed ETH_MAX_GAS_PRICE_WEI
		// Existing gas price is: 60480000000
		gasPrice := attempt3_4.GasPrice.ToInt()
		config.Set("ETH_MAX_GAS_PRICE_WEI", 60500000000)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == *etx3.Nonce && gasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(errors.New("already known")).Once() // we already submitted at this price, now its time to bump and submit again but since we simply resubmitted rather than increasing gas price, geth already knows about this tx

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		etx3, err = cltest.FindEthTxWithAttempts(db, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx3.State)

		// No new tx attempts
		require.Len(t, etx3.EthTxAttempts, 4)
		attempt3_4 = etx3.EthTxAttempts[3]
		assert.Equal(t, gasPrice.Int64(), attempt3_4.GasPrice.ToInt().Int64())

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	attempt3_4.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt3_4).Error)

	t.Run("resubmits at the old price and does not create a new attempt if the current price is exactly ETH_MAX_GAS_PRICE_WEI", func(t *testing.T) {
		// Set price such that the current price is already at ETH_MAX_GAS_PRICE_WEI
		// Existing gas price is: 60480000000
		gasPrice := attempt3_4.GasPrice.ToInt()
		config.Set("ETH_MAX_GAS_PRICE_WEI", 60480000000)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return int64(tx.Nonce()) == *etx3.Nonce && gasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(errors.New("already known")).Once() // we already submitted at this price, now its time to bump and submit again but since we simply resubmitted rather than increasing gas price, geth already knows about this tx

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		etx3, err = cltest.FindEthTxWithAttempts(db, etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx3.State)

		// No new tx attempts
		require.Len(t, etx3.EthTxAttempts, 4)
		attempt3_4 := etx3.EthTxAttempts[3]
		assert.Equal(t, gasPrice.Int64(), attempt3_4.GasPrice.ToInt().Int64())

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	kst.AssertExpectations(t)
	ethClient.AssertExpectations(t)
}

func TestEthConfirmer_RebroadcastWhereNecessary_WhenOutOfEth(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB
	ethClient := new(mocks.Client)
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)
	ethKeyStore.Unlock(cltest.Password)

	keys, err := ethKeyStore.SendingKeys()
	require.NoError(t, err)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	currentHead := int64(30)
	oldEnough := int64(19)
	nonce := int64(0)

	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, nonce, fromAddress)
	nonce++
	attempt1_1 := etx.EthTxAttempts[0]
	attempt1_1.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt1_1).Error)
	var attempt1_2 bulletprooftxmanager.EthTxAttempt

	insufficientEthError := errors.New("insufficient funds for gas * price + value")

	t.Run("saves attempt with state 'insufficient_eth' if eth node returns this error", func(t *testing.T) {
		ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, keys)

		expectedBumpedGasPrice := big.NewInt(25000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.GasPrice.ToInt().Int64())

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(insufficientEthError).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.EthTxAttempts, 2)
		require.Equal(t, attempt1_1.ID, etx.EthTxAttempts[0].ID)

		// Got the new attempt
		attempt1_2 = etx.EthTxAttempts[1]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.GasPrice.ToInt().Int64())
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptInsufficientEth, attempt1_2.State)
		assert.Nil(t, attempt1_2.BroadcastBeforeBlockNum)

		ethClient.AssertExpectations(t)
	})

	t.Run("does not bump gas when previous error was 'out of eth', instead resubmits existing transaction", func(t *testing.T) {
		ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, keys)

		expectedBumpedGasPrice := big.NewInt(25000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.GasPrice.ToInt().Int64())

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(insufficientEthError).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		// New attempt was NOT created
		require.Len(t, etx.EthTxAttempts, 2)

		// The attempt is still "out of eth"
		attempt1_2 = etx.EthTxAttempts[1]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.GasPrice.ToInt().Int64())
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptInsufficientEth, attempt1_2.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("saves the attempt as broadcast after node wallet has been topped up with sufficient balance", func(t *testing.T) {
		ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, keys)

		expectedBumpedGasPrice := big.NewInt(25000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.GasPrice.ToInt().Int64())

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		etx, err = cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)

		// New attempt was NOT created
		require.Len(t, etx.EthTxAttempts, 2)

		// Attempt is now 'broadcast'
		attempt1_2 = etx.EthTxAttempts[1]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.GasPrice.ToInt().Int64())
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt1_2.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("resubmitting due to insufficient eth is not limited by ETH_GAS_BUMP_TX_DEPTH", func(t *testing.T) {
		depth := 2
		etxCount := 4

		config.Set("ETH_GAS_BUMP_TX_DEPTH", depth)
		ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, keys)

		for i := 0; i < etxCount; i++ {
			n := nonce
			cltest.MustInsertUnconfirmedEthTxWithInsufficientEthAttempt(t, db, nonce, fromAddress)
			ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
				return tx.Nonce() == uint64(n)
			})).Return(nil).Once()

			nonce++
		}

		require.NoError(t, ec.RebroadcastWhereNecessary(context.TODO(), currentHead))

		var attempts []bulletprooftxmanager.EthTxAttempt
		require.NoError(t, store.DB.Where("state = 'insufficient_eth'").Find(&attempts).Error)
		require.Len(t, attempts, 0)

		ethClient.AssertExpectations(t)
	})
}

func TestEthConfirmer_EnsureConfirmedTransactionsInLongestChain(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)

	ethClient := new(mocks.Client)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, []ethkey.Key{key})

	head := models.Head{
		Hash:   utils.NewHash(),
		Number: 10,
		Parent: &models.Head{
			Hash:   utils.NewHash(),
			Number: 9,
			Parent: &models.Head{
				Number: 8,
				Hash:   utils.NewHash(),
				Parent: nil,
			},
		},
	}

	t.Run("does nothing if there aren't any transactions", func(t *testing.T) {
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(context.TODO(), head))
	})

	t.Run("does nothing to unconfirmed transactions", func(t *testing.T) {
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 0, fromAddress)

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(context.TODO(), head))

		etx, err := cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx.State)
	})

	t.Run("does nothing to confirmed transactions with receipts within head height of the chain and included in the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, db, 2, 1, fromAddress)
		attempt := etx.EthTxAttempts[0]
		cltest.MustInsertEthReceipt(t, db, head.Number, head.Hash, attempt.Hash)

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(context.TODO(), head))

		etx, err := cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx.State)
	})

	t.Run("does nothing to confirmed transactions that only have receipts older than the start of the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, db, 3, 1, fromAddress)
		attempt := etx.EthTxAttempts[0]
		// Add receipt that is older than the lowest block of the chain
		cltest.MustInsertEthReceipt(t, db, head.Parent.Parent.Number-1, utils.NewHash(), attempt.Hash)

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(context.TODO(), head))

		etx, err := cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx.State)
	})

	t.Run("unconfirms and rebroadcasts transactions that have receipts within head height of the chain but not included in the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, db, 4, 1, fromAddress)
		attempt := etx.EthTxAttempts[0]
		// Include one within head height but a different block hash
		cltest.MustInsertEthReceipt(t, db, head.Parent.Number, utils.NewHash(), attempt.Hash)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			atx, err := attempt.GetSignedTx()
			require.NoError(t, err)
			// Keeps gas price and nonce the same
			return atx.GasPrice().Cmp(tx.GasPrice()) == 0 && atx.Nonce() == tx.Nonce()
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(context.TODO(), head))

		etx, err := cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt = etx.EthTxAttempts[0]
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("unconfirms and rebroadcasts transactions that have receipts within head height of chain but not included in the chain even if a receipt exists older than the start of the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, db, 5, 1, fromAddress)
		attempt := etx.EthTxAttempts[0]
		// Add receipt that is older than the lowest block of the chain
		cltest.MustInsertEthReceipt(t, db, head.Parent.Parent.Number-1, utils.NewHash(), attempt.Hash)
		// Include one within head height but a different block hash
		cltest.MustInsertEthReceipt(t, db, head.Parent.Number, utils.NewHash(), attempt.Hash)

		ethClient.On("SendTransaction", mock.Anything, mock.Anything).Return(nil).Once()

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(context.TODO(), head))

		etx, err := cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt = etx.EthTxAttempts[0]
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("if more than one attempt has a receipt (should not be possible but isn't prevented by database constraints) unconfirms and rebroadcasts only the attempt with the highest gas price", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, db, 6, 1, fromAddress)
		require.Len(t, etx.EthTxAttempts, 1)
		// Sanity check to assert the included attempt has the lowest gas price
		require.Less(t, etx.EthTxAttempts[0].GasPrice.ToInt().Int64(), int64(30000))

		attempt2 := newBroadcastEthTxAttempt(t, etx.ID, 30000)
		attempt2.SignedRawTx = hexutil.MustDecode("0xf88c8301f3a98503b9aca000832ab98094f5fff180082d6017036b771ba883025c654bc93580a4daa6d556000000000000000000000000000000000000000000000000000000000000000026a0f25601065ee369b6470c0399a2334afcfbeb0b5c8f3d9a9042e448ed29b5bcbda05b676e00248b85faf4dd889f0e2dcf91eb867e23ac9eeb14a73f9e4c14972cdf")
		attempt3 := newBroadcastEthTxAttempt(t, etx.ID, 40000)
		attempt3.SignedRawTx = hexutil.MustDecode("0xf88c8301f3a88503b9aca0008316e36094151445852b0cfdf6a4cc81440f2af99176e8ad0880a4daa6d556000000000000000000000000000000000000000000000000000000000000000026a0dcb5a7ad52b96a866257134429f944c505820716567f070e64abb74899803855a04c13eff2a22c218e68da80111e1bb6dc665d3dea7104ab40ff8a0275a99f630d")
		require.NoError(t, store.DB.Create(&attempt2).Error)
		require.NoError(t, store.DB.Create(&attempt3).Error)

		// Receipt is within head height but a different block hash
		cltest.MustInsertEthReceipt(t, db, head.Parent.Number, utils.NewHash(), attempt2.Hash)
		// Receipt is within head height but a different block hash
		cltest.MustInsertEthReceipt(t, db, head.Parent.Number, utils.NewHash(), attempt3.Hash)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			s, err := attempt3.GetSignedTx()
			require.NoError(t, err)
			return tx.Hash() == s.Hash()
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(context.TODO(), head))

		etx, err := cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, bulletprooftxmanager.EthTxUnconfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 3)
		attempt1 := etx.EthTxAttempts[0]
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt1.State)
		attempt2 = etx.EthTxAttempts[1]
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt2.State)
		attempt3 = etx.EthTxAttempts[2]
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt3.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("if receipt has a block number that is in the future, does not mark for rebroadcast (the safe thing to do is simply wait until heads catches up)", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, db, 7, 1, fromAddress)
		attempt := etx.EthTxAttempts[0]
		// Add receipt that is higher than head
		cltest.MustInsertEthReceipt(t, db, head.Number+1, utils.NewHash(), attempt.Hash)

		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(context.TODO(), head))

		etx, err := cltest.FindEthTxWithAttempts(db, etx.ID)
		require.NoError(t, err)
		assert.Equal(t, bulletprooftxmanager.EthTxConfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt = etx.EthTxAttempts[0]
		assert.Equal(t, bulletprooftxmanager.EthTxAttemptBroadcast, attempt.State)
		assert.Len(t, attempt.EthReceipts, 1)

		ethClient.AssertExpectations(t)
	})
}

func TestEthConfirmer_ForceRebroadcast(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	db := store.DB
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()
	key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)
	ethKeyStore.Unlock(cltest.Password)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	mustInsertUnstartedEthTx(t, db, fromAddress)
	mustInsertInProgressEthTx(t, db, 0, fromAddress)
	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 1, fromAddress)
	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, db, 2, fromAddress)

	gasPriceWei := uint64(52)
	overrideGasLimit := uint64(20000)

	t.Run("rebroadcasts one eth_tx if it falls within in nonce range", func(t *testing.T) {
		ethClient := new(mocks.Client)
		ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, []ethkey.Key{key})

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
		ethClient := new(mocks.Client)
		ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, []ethkey.Key{key})

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
		ethClient := new(mocks.Client)
		ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, []ethkey.Key{key})

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
		ethClient := new(mocks.Client)
		ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, []ethkey.Key{key})

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
		ethClient := new(mocks.Client)
		ec := cltest.NewEthConfirmer(t, store.DB, ethClient, config, ethKeyStore, []ethkey.Key{key})

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *types.Transaction) bool {
			return tx.Nonce() == uint64(0) && uint64(tx.GasPrice().Int64()) == gasPriceWei && uint64(tx.Gas()) == config.EthGasLimitDefault()
		})).Return(nil).Once()

		require.NoError(t, ec.ForceRebroadcast(0, 0, gasPriceWei, fromAddress, 0))

		ethClient.AssertExpectations(t)
	})
}
