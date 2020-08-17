package bulletprooftxmanager_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	gethAccounts "github.com/ethereum/go-ethereum/accounts"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func mustInsertUnstartedEthTx(t *testing.T, s *store.Store) {
	etx := cltest.NewEthTx(t, s)
	etx.State = models.EthTxUnstarted
	require.NoError(t, s.DB.Save(&etx).Error)
}

func newBroadcastEthTxAttempt(t *testing.T, etxID int64, store *store.Store, gasPrice ...int64) models.EthTxAttempt {
	attempt := cltest.NewEthTxAttempt(t, etxID)
	attempt.State = models.EthTxAttemptBroadcast
	if len(gasPrice) > 0 {
		gp := gasPrice[0]
		attempt.GasPrice = *utils.NewBig(big.NewInt(gp))
	}
	return attempt
}

func mustInsertInProgressEthTx(t *testing.T, store *store.Store, nonce int64) models.EthTx {
	etx := cltest.NewEthTx(t, store)
	etx.State = models.EthTxInProgress
	etx.Nonce = &nonce
	require.NoError(t, store.DB.Save(&etx).Error)

	return etx
}

func TestEthConfirmer_SetBroadcastBeforeBlockNum(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	ec := bulletprooftxmanager.NewEthConfirmer(store, config)

	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 0)

	headNum := int64(9000)
	var err error

	t.Run("saves block num to unconfirmed eth_tx_attempts without one", func(t *testing.T) {
		// Do the thing
		require.NoError(t, ec.SetBroadcastBeforeBlockNum(headNum))

		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]

		assert.Equal(t, int64(9000), *attempt.BroadcastBeforeBlockNum)
	})

	t.Run("does not change eth_tx_attempts that already have BroadcastBeforeBlockNum set", func(t *testing.T) {
		n := int64(42)
		attempt := newBroadcastEthTxAttempt(t, etx.ID, store, 2)
		attempt.BroadcastBeforeBlockNum = &n
		require.NoError(t, store.DB.Save(&attempt).Error)

		// Do the thing
		require.NoError(t, ec.SetBroadcastBeforeBlockNum(headNum))

		etx, err = store.FindEthTxWithAttempts(etx.ID)
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
	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	ec := bulletprooftxmanager.NewEthConfirmer(store, config)

	nonce := int64(0)
	var err error

	t.Run("only finds eth_txes in unconfirmed state", func(t *testing.T) {
		cltest.MustInsertFatalErrorEthTx(t, store)
		mustInsertInProgressEthTx(t, store, nonce)
		nonce++
		cltest.MustInsertConfirmedEthTxWithAttempt(t, store, nonce, 1)
		nonce++
		mustInsertUnstartedEthTx(t, store)

		// Do the thing
		require.NoError(t, ec.CheckForReceipts())
		// No calls
		ethClient.AssertExpectations(t)
	})

	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, nonce)
	nonce++
	require.Len(t, etx1.EthTxAttempts, 1)
	attempt1_1 := etx1.EthTxAttempts[0]
	require.Len(t, attempt1_1.EthReceipts, 0)

	t.Run("fetches receipt for an unconfirmed eth_tx", func(t *testing.T) {
		// Transaction not confirmed yet, receipt is nil
		ethClient.On("TransactionReceipt", mock.Anything, mock.MatchedBy(func(txHash gethCommon.Hash) bool {
			return txHash == attempt1_1.Hash
		})).Return(nil, errors.New("not found")).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts())

		etx1, err = store.FindEthTxWithAttempts(etx1.ID)
		require.Len(t, etx1.EthTxAttempts, 1)
		attempt1_1 = etx1.EthTxAttempts[0]
		require.NoError(t, err)
		require.Len(t, attempt1_1.EthReceipts, 0)

		ethClient.AssertExpectations(t)
	})

	t.Run("returns error and does not save anything if TransactionReceipt returns error", func(t *testing.T) {
		// First transaction confirmed
		ethClient.On("TransactionReceipt", mock.Anything, mock.MatchedBy(func(txHash gethCommon.Hash) bool {
			return txHash == attempt1_1.Hash
		})).Return(nil, errors.New("something exploded")).Once()

		// Do the thing
		err := ec.CheckForReceipts()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "something exploded")
	})

	t.Run("returns error and saves nothing if returned receipt does not match the attempt", func(t *testing.T) {
		gethReceipt := gethTypes.Receipt{
			TxHash:           cltest.NewHash(),
			BlockHash:        cltest.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}

		// First transaction confirmed
		ethClient.On("TransactionReceipt", mock.Anything, mock.MatchedBy(func(txHash gethCommon.Hash) bool {
			return txHash == attempt1_1.Hash
		})).Return(&gethReceipt, nil).Once()

		// Do the thing
		err := ec.CheckForReceipts()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invariant violation: expected receipt with hash")
	})

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, nonce)
	nonce++
	require.Len(t, etx2.EthTxAttempts, 1)
	attempt2_1 := etx2.EthTxAttempts[0]
	require.Len(t, attempt2_1.EthReceipts, 0)

	t.Run("saves eth_receipt and marks eth_tx as confirmed when geth client returns valid receipt", func(t *testing.T) {
		gethReceipt := gethTypes.Receipt{
			TxHash:           attempt1_1.Hash,
			BlockHash:        cltest.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}

		// First transaction confirmed
		ethClient.On("TransactionReceipt", mock.Anything, mock.MatchedBy(func(txHash gethCommon.Hash) bool {
			return txHash == attempt1_1.Hash
		})).Return(&gethReceipt, nil).Once()
		// Second transaction still unconfirmed
		ethClient.On("TransactionReceipt", mock.Anything, mock.MatchedBy(func(txHash gethCommon.Hash) bool {
			return txHash == attempt2_1.Hash
		})).Return(nil, errors.New("not found")).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts())

		// Check that the receipt was saved
		etx, err := store.FindEthTxWithAttempts(etx1.ID)
		require.NoError(t, err)

		assert.Equal(t, models.EthTxConfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt1_1 = etx.EthTxAttempts[0]
		require.Len(t, attempt1_1.EthReceipts, 1)

		ethReceipt := attempt1_1.EthReceipts[0]

		assert.Equal(t, gethReceipt.TxHash, ethReceipt.TxHash)
		assert.Equal(t, gethReceipt.BlockHash, ethReceipt.BlockHash)
		assert.Equal(t, gethReceipt.BlockNumber.Int64(), ethReceipt.BlockNumber)
		assert.Equal(t, gethReceipt.TransactionIndex, ethReceipt.TransactionIndex)

		receiptJSON, err := json.Marshal(gethReceipt)
		require.NoError(t, err)

		assert.JSONEq(t, string(receiptJSON), string(ethReceipt.Receipt))

		ethClient.AssertExpectations(t)
	})

	t.Run("fetches and saves receipts for several attempts in gas price order", func(t *testing.T) {
		attempt2_2 := newBroadcastEthTxAttempt(t, etx2.ID, store)
		attempt2_2.GasPrice = *utils.NewBig(big.NewInt(10))

		attempt2_3 := newBroadcastEthTxAttempt(t, etx2.ID, store)
		attempt2_3.GasPrice = *utils.NewBig(big.NewInt(20))

		// Insert order deliberately reversed to test sorting by gas price
		require.NoError(t, store.DB.Create(&attempt2_3).Error)
		require.NoError(t, store.DB.Create(&attempt2_2).Error)

		// Most expensive attempt still unconfirmed
		ethClient.On("TransactionReceipt", mock.Anything, mock.MatchedBy(func(txHash gethCommon.Hash) bool {
			return txHash == attempt2_3.Hash
		})).Return(nil, errors.New("not found")).Once()

		gethReceipt := gethTypes.Receipt{
			TxHash:           attempt2_2.Hash,
			BlockHash:        cltest.NewHash(),
			BlockNumber:      big.NewInt(42),
			TransactionIndex: uint(1),
		}
		// Second most expensive attempt is confirmed
		ethClient.On("TransactionReceipt", mock.Anything, mock.MatchedBy(func(txHash gethCommon.Hash) bool {
			return txHash == attempt2_2.Hash
		})).Return(&gethReceipt, nil).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts())

		ethClient.AssertExpectations(t)

		// Check that the state was updated
		etx, err := store.FindEthTxWithAttempts(etx2.ID)
		require.NoError(t, err)

		require.Equal(t, models.EthTxConfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 3)
	})

	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, nonce)
	attempt3_1 := etx3.EthTxAttempts[0]
	nonce++

	t.Run("ignores error that comes from querying parity too early", func(t *testing.T) {
		ethClient.On("TransactionReceipt", mock.Anything, mock.MatchedBy(func(txHash gethCommon.Hash) bool {
			return txHash == attempt3_1.Hash
		})).Return(nil, errors.New("missing required field 'transactionHash' for Log")).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts())

		// No receipt, but no error either
		etx, err := store.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, models.EthTxUnconfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt3_1 = etx.EthTxAttempts[0]
		require.Len(t, attempt3_1.EthReceipts, 0)
	})

	t.Run("ignores partially hydrated receipt that comes from querying parity too early", func(t *testing.T) {
		receipt := gethTypes.Receipt{
			TxHash: attempt3_1.Hash,
		}
		ethClient.On("TransactionReceipt", mock.Anything, mock.MatchedBy(func(txHash gethCommon.Hash) bool {
			return txHash == attempt3_1.Hash
		})).Return(&receipt, nil).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts())

		// No receipt, but no error either
		etx, err := store.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, models.EthTxUnconfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt3_1 = etx.EthTxAttempts[0]
		require.Len(t, attempt3_1.EthReceipts, 0)
	})

	t.Run("handles case where eth_receipt already exists somehow", func(t *testing.T) {
		ethReceipt := cltest.MustInsertEthReceipt(t, store, 42, cltest.NewHash(), attempt3_1.Hash)

		gethReceipt := gethTypes.Receipt{
			TxHash:           attempt3_1.Hash,
			BlockHash:        ethReceipt.BlockHash,
			BlockNumber:      big.NewInt(ethReceipt.BlockNumber),
			TransactionIndex: ethReceipt.TransactionIndex,
		}
		ethClient.On("TransactionReceipt", mock.Anything, mock.MatchedBy(func(txHash gethCommon.Hash) bool {
			return txHash == attempt3_1.Hash
		})).Return(&gethReceipt, nil).Once()

		// Do the thing
		require.NoError(t, ec.CheckForReceipts())

		// Check that the receipt was unchanged
		etx, err := store.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, models.EthTxConfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt3_1 = etx.EthTxAttempts[0]
		require.Len(t, attempt3_1.EthReceipts, 1)

		ethReceipt = attempt3_1.EthReceipts[0]

		assert.Equal(t, gethReceipt.TxHash, ethReceipt.TxHash)
		assert.Equal(t, gethReceipt.BlockHash, ethReceipt.BlockHash)
		assert.Equal(t, gethReceipt.BlockNumber.Int64(), ethReceipt.BlockNumber)
		assert.Equal(t, gethReceipt.TransactionIndex, ethReceipt.TransactionIndex)

		ethClient.AssertExpectations(t)
	})
}

func TestEthConfirmer_FindEthTxsRequiringNewAttempt(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	currentHead := int64(30)
	gasBumpThreshold := int64(10)
	tooNew := int64(21)
	onTheMoney := int64(20)
	oldEnough := int64(19)
	nonce := int64(0)

	t.Run("returns nothing when there are no transactions", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringNewAttempt(store.DB, currentHead, gasBumpThreshold)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	mustInsertInProgressEthTx(t, store, nonce)
	nonce++

	t.Run("returns nothing when the transaction is in_progress", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringNewAttempt(store.DB, currentHead, gasBumpThreshold)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	// This one has BroadcastBeforeBlockNum set as nil... which can happen, but it should be ignored
	cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, nonce)
	nonce++

	t.Run("ignores unconfirmed transactions with nil BroadcastBeforeBlockNum", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringNewAttempt(store.DB, currentHead, gasBumpThreshold)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, nonce)
	nonce++
	attempt1_1 := etx1.EthTxAttempts[0]
	attempt1_1.BroadcastBeforeBlockNum = &tooNew
	require.NoError(t, store.DB.Save(&attempt1_1).Error)
	attempt1_2 := newBroadcastEthTxAttempt(t, etx1.ID, store)
	attempt1_2.BroadcastBeforeBlockNum = &onTheMoney
	attempt1_2.GasPrice = *utils.NewBigI(30000)
	require.NoError(t, store.DB.Save(&attempt1_2).Error)

	t.Run("returns nothing when the transaction is unconfirmed with an attempt that is recent", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringNewAttempt(store.DB, currentHead, gasBumpThreshold)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, nonce)
	nonce++
	attempt2_1 := etx2.EthTxAttempts[0]
	attempt2_1.BroadcastBeforeBlockNum = &tooNew
	require.NoError(t, store.DB.Save(&attempt2_1).Error)

	t.Run("returns nothing when the transaction has attempts that are too new", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringNewAttempt(store.DB, currentHead, gasBumpThreshold)
		require.NoError(t, err)

		assert.Len(t, etxs, 0)
	})

	etxWithoutAttempts := cltest.NewEthTx(t, store)
	etxWithoutAttempts.Nonce = &nonce
	now := time.Now()
	etxWithoutAttempts.BroadcastAt = &now
	etxWithoutAttempts.State = models.EthTxUnconfirmed
	require.NoError(t, store.DB.Save(&etxWithoutAttempts).Error)
	nonce++

	t.Run("returns the transaction if it is unconfirmed and has no attempts (note that this is an invariant violation, but we handle it anyway)", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringNewAttempt(store.DB, currentHead, gasBumpThreshold)
		require.NoError(t, err)

		require.Len(t, etxs, 1)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
	})

	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, nonce)
	nonce++
	attempt3_1 := etx3.EthTxAttempts[0]
	attempt3_1.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt3_1).Error)

	t.Run("returns the transaction if it is unconfirmed with an attempt that is older than gasBumpThreshold blocks", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringNewAttempt(store.DB, currentHead, gasBumpThreshold)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
	})

	attempt3_2 := newBroadcastEthTxAttempt(t, etx3.ID, store)
	attempt3_2.BroadcastBeforeBlockNum = &oldEnough
	attempt3_2.GasPrice = *utils.NewBigI(30000)
	require.NoError(t, store.DB.Save(&attempt3_2).Error)

	t.Run("returns the transaction if it is unconfirmed with two attempts that are older than gasBumpThreshold blocks", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringNewAttempt(store.DB, currentHead, gasBumpThreshold)
		require.NoError(t, err)

		require.Len(t, etxs, 2)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
		assert.Equal(t, etx3.ID, etxs[1].ID)
	})

	attempt3_3 := newBroadcastEthTxAttempt(t, etx3.ID, store)
	attempt3_3.BroadcastBeforeBlockNum = &tooNew
	attempt3_3.GasPrice = *utils.NewBigI(40000)
	require.NoError(t, store.DB.Save(&attempt3_3).Error)

	t.Run("does not return the transaction if it has some older but one newer attempt", func(t *testing.T) {
		etxs, err := bulletprooftxmanager.FindEthTxsRequiringNewAttempt(store.DB, currentHead, gasBumpThreshold)
		require.NoError(t, err)

		require.Len(t, etxs, 1)
		assert.Equal(t, etxWithoutAttempts.ID, etxs[0].ID)
	})
}

func TestEthConfirmer_BumpGasWhereNecessary(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	kst := new(mocks.KeyStoreInterface)
	// Use a mock keystore for this test
	store.KeyStore = kst
	ec := bulletprooftxmanager.NewEthConfirmer(store, config)
	currentHead := int64(30)
	oldEnough := int64(19)
	nonce := int64(0)

	defaultFromAddress := cltest.GetDefaultFromAddress(t, store)
	kst.On("GetAccountByAddress", defaultFromAddress).
		Return(gethAccounts.Account{Address: defaultFromAddress}, nil)

	t.Run("does nothing if no transactions require bumping", func(t *testing.T) {
		require.NoError(t, ec.BumpGasWhereNecessary(currentHead))
	})

	etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, nonce)
	nonce++
	attempt1_1 := etx.EthTxAttempts[0]
	attempt1_1.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt1_1).Error)
	var err error

	t.Run("returns on keystore error", func(t *testing.T) {
		// simulate transaction that is somehow impossible to sign
		kst.On("SignTx", mock.Anything,
			mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				return tx.Nonce() == uint64(*etx.Nonce)
			}),
			mock.Anything).Return(nil, errors.New("signing error")).Once()

		// Do the thing
		err = ec.BumpGasWhereNecessary(currentHead)
		require.Error(t, err)
		require.Contains(t, err.Error(), "signing error")

		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.Equal(t, models.EthTxUnconfirmed, etx.State)

		require.Len(t, etx.EthTxAttempts, 1)

		kst.AssertExpectations(t)
	})

	kst = new(mocks.KeyStoreInterface)
	store.KeyStore = kst
	kst.On("GetAccountByAddress", defaultFromAddress).
		Return(gethAccounts.Account{Address: defaultFromAddress}, nil)

	t.Run("does nothing and continues on fatal error", func(t *testing.T) {
		ethTx := gethTypes.Transaction{}
		kst.On("SignTx",
			mock.AnythingOfType("accounts.Account"),
			mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				if tx.Nonce() != uint64(*etx.Nonce) {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(store.Config.ChainID()) == 0
			})).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(*etx.Nonce)
		})).Return(errors.New("exceeds block gas limit")).Once()

		// Do the thing
		require.NoError(t, ec.BumpGasWhereNecessary(currentHead))

		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.EthTxAttempts, 1)

		kst.AssertExpectations(t)
	})

	kst = new(mocks.KeyStoreInterface)
	store.KeyStore = kst
	kst.On("GetAccountByAddress", defaultFromAddress).
		Return(gethAccounts.Account{Address: defaultFromAddress}, nil)
	var attempt1_2 models.EthTxAttempt

	t.Run("creates new attempt with higher gas price if transaction has an attempt older than threshold", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(25000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_1.GasPrice.ToInt().Int64())

		ethTx := gethTypes.Transaction{}
		kst.On("SignTx",
			mock.AnythingOfType("accounts.Account"),
			mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				if expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(store.Config.ChainID()) == 0
			})).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, ec.BumpGasWhereNecessary(currentHead))

		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.EthTxAttempts, 2)
		require.Equal(t, attempt1_1.ID, etx.EthTxAttempts[0].ID)

		// Got the new attempt
		attempt1_2 = etx.EthTxAttempts[1]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_2.GasPrice.ToInt().Int64())
		assert.Equal(t, models.EthTxAttemptBroadcast, attempt1_2.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("does nothing if there is an attempt without BroadcastBeforeBlockNum set", func(t *testing.T) {
		// Do the thing
		require.NoError(t, ec.BumpGasWhereNecessary(currentHead))

		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.EthTxAttempts, 2)
	})

	attempt1_2.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt1_2).Error)
	var attempt1_3 models.EthTxAttempt

	t.Run("creates new attempt with higher gas price if transaction is already in mempool (e.g. due to previous crash before we could save the new attempt)", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(30000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_2.GasPrice.ToInt().Int64())

		ethTx := gethTypes.Transaction{}
		kst.On("SignTx",
			mock.AnythingOfType("accounts.Account"),
			mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				if int64(tx.Nonce()) != *etx.Nonce || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(fmt.Errorf("known transaction: %s", ethTx.Hash().Hex())).Once()

		// Do the thing
		require.NoError(t, ec.BumpGasWhereNecessary(currentHead))

		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		require.Len(t, etx.EthTxAttempts, 3)
		require.Equal(t, attempt1_1.ID, etx.EthTxAttempts[0].ID)
		require.Equal(t, attempt1_2.ID, etx.EthTxAttempts[1].ID)

		// Got the new attempt
		attempt1_3 = etx.EthTxAttempts[2]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt1_3.GasPrice.ToInt().Int64())
		assert.Equal(t, models.EthTxAttemptBroadcast, attempt1_3.State)

		kst.AssertExpectations(t)
		ethClient.AssertExpectations(t)
	})

	attempt1_3.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt1_3).Error)

	t.Run("does not save new attempt for transaction that has already been confirmed (nonce already used)", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(36000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt1_2.GasPrice.ToInt().Int64())

		ethTx := gethTypes.Transaction{}
		receipt := gethTypes.Receipt{BlockNumber: big.NewInt(40)}
		kst.On("SignTx",
			mock.AnythingOfType("accounts.Account"),
			mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				if int64(tx.Nonce()) != *etx.Nonce || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				receipt.TxHash = tx.Hash()
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(errors.New("nonce too low")).Once()

		// Do the thing
		require.NoError(t, ec.BumpGasWhereNecessary(currentHead))

		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Equal(t, models.EthTxUnconfirmed, etx.State)

		require.Len(t, etx.EthTxAttempts, 3)
		require.Equal(t, attempt1_1.ID, etx.EthTxAttempts[0].ID)
		require.Equal(t, attempt1_2.ID, etx.EthTxAttempts[1].ID)
		require.Equal(t, attempt1_3.ID, etx.EthTxAttempts[2].ID)
		require.Equal(t, models.EthTxAttemptBroadcast, etx.EthTxAttempts[0].State)
		require.Equal(t, models.EthTxAttemptBroadcast, etx.EthTxAttempts[1].State)
		require.Equal(t, models.EthTxAttemptBroadcast, etx.EthTxAttempts[2].State)

		kst.AssertExpectations(t)
		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	// Mark original tx as confirmed so we won't pick it up any more
	require.NoError(t, store.DB.Exec(`UPDATE eth_txes SET state = 'confirmed'`).Error)

	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, nonce)
	nonce++
	attempt2_1 := etx2.EthTxAttempts[0]
	attempt2_1.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt2_1).Error)
	var attempt2_2 models.EthTxAttempt

	t.Run("saves in_progress attempt on temporary error and returns error", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(25000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt2_1.GasPrice.ToInt().Int64())

		ethTx := gethTypes.Transaction{}
		n := *etx2.Nonce
		kst.On("SignTx",
			mock.AnythingOfType("accounts.Account"),
			mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				if int64(tx.Nonce()) != n || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return int64(tx.Nonce()) == n && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(errors.New("some network error")).Once()

		// Do the thing
		err = ec.BumpGasWhereNecessary(currentHead)
		require.Error(t, err)
		require.Contains(t, err.Error(), "some network error")

		etx2, err = store.FindEthTxWithAttempts(etx2.ID)
		require.NoError(t, err)

		assert.Equal(t, models.EthTxUnconfirmed, etx2.State)

		// Old attempt is untouched
		require.Len(t, etx2.EthTxAttempts, 2)
		require.Equal(t, attempt2_1.ID, etx2.EthTxAttempts[0].ID)
		attempt2_1 = etx2.EthTxAttempts[0]
		assert.Equal(t, models.EthTxAttemptBroadcast, attempt2_1.State)
		assert.Equal(t, oldEnough, *attempt2_1.BroadcastBeforeBlockNum)

		// New in_progress attempt saved
		attempt2_2 = etx2.EthTxAttempts[1]
		assert.Equal(t, models.EthTxAttemptInProgress, attempt2_2.State)
		assert.Nil(t, attempt2_2.BroadcastBeforeBlockNum)

		// Do it again and move the attempt into "broadcast"
		n = *etx2.Nonce
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return int64(tx.Nonce()) == n && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(nil).Once()

		require.NoError(t, ec.BumpGasWhereNecessary(currentHead))

		// Attempt marked "broadcast"
		etx2, err = store.FindEthTxWithAttempts(etx2.ID)
		require.NoError(t, err)

		assert.Equal(t, models.EthTxUnconfirmed, etx2.State)

		// New in_progress attempt saved
		require.Len(t, etx2.EthTxAttempts, 2)
		require.Equal(t, attempt2_2.ID, etx2.EthTxAttempts[1].ID)
		attempt2_2 = etx2.EthTxAttempts[1]
		require.Equal(t, models.EthTxAttemptBroadcast, attempt2_2.State)
		assert.Nil(t, attempt2_2.BroadcastBeforeBlockNum)

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	// Set BroadcastBeforeBlockNum again so the next test will pick it up
	attempt2_2.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt2_2).Error)

	t.Run("handles case where nonce is too low but receipt is nil indicating that an external wallet used the nonce (until finalized)", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(30000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt2_1.GasPrice.ToInt().Int64())

		ethTx := gethTypes.Transaction{}
		n := *etx2.Nonce
		kst.On("SignTx",
			mock.AnythingOfType("accounts.Account"),
			mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				if int64(tx.Nonce()) != n || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Twice()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return int64(tx.Nonce()) == n && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(errors.New("nonce too low")).Twice()

		// Does nothing if currentHead is not high enough
		require.NoError(t, ec.BumpGasWhereNecessary(currentHead))
		etx2, err = store.FindEthTxWithAttempts(etx2.ID)
		require.NoError(t, err)
		assert.Equal(t, models.EthTxUnconfirmed, etx2.State)

		// No new attempts saved
		require.Len(t, etx2.EthTxAttempts, 2)
		assert.Equal(t, models.EthTxAttemptBroadcast, etx2.EthTxAttempts[0].State)
		assert.Equal(t, models.EthTxAttemptBroadcast, etx2.EthTxAttempts[1].State)

		// When currentHead reaches the threshold, we save it as failed
		require.NoError(t, ec.BumpGasWhereNecessary(currentHead+100))

		etx2, err = store.FindEthTxWithAttempts(etx2.ID)
		require.NoError(t, err)
		assert.Equal(t, models.EthTxFatalError, etx2.State)

		// No new attempts saved
		require.Len(t, etx2.EthTxAttempts, 2)
		assert.Equal(t, models.EthTxAttemptBroadcast, etx2.EthTxAttempts[0].State)
		assert.Equal(t, models.EthTxAttemptBroadcast, etx2.EthTxAttempts[1].State)

		ethClient.AssertExpectations(t)
		kst.AssertExpectations(t)
	})

	// Original tx is confirmed so we won't pick it up any more
	etx3 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, nonce)
	nonce++
	attempt3_1 := etx3.EthTxAttempts[0]
	attempt3_1.BroadcastBeforeBlockNum = &oldEnough
	attempt3_1.GasPrice = *utils.NewBig(big.NewInt(35000000000))
	require.NoError(t, store.DB.Save(&attempt3_1).Error)

	var attempt3_2 models.EthTxAttempt

	t.Run("saves attempt anyway if replacement transaction is underpriced because the bumped gas price is insufficiently higher than the previous one", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(42000000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt3_1.GasPrice.ToInt().Int64())

		ethTx := gethTypes.Transaction{}
		kst.On("SignTx",
			mock.AnythingOfType("accounts.Account"),
			mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				if int64(tx.Nonce()) != *etx3.Nonce || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return int64(tx.Nonce()) == *etx3.Nonce && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(errors.New("replacement transaction underpriced")).Once()

		// Do the thing
		require.NoError(t, ec.BumpGasWhereNecessary(currentHead))

		etx3, err = store.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, models.EthTxUnconfirmed, etx3.State)

		require.Len(t, etx3.EthTxAttempts, 2)
		require.Equal(t, attempt3_1.ID, etx3.EthTxAttempts[0].ID)
		attempt3_2 = etx3.EthTxAttempts[1]

		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt3_2.GasPrice.ToInt().Int64())

		kst.AssertExpectations(t)
		ethClient.AssertExpectations(t)
	})

	attempt3_2.BroadcastBeforeBlockNum = &oldEnough
	require.NoError(t, store.DB.Save(&attempt3_2).Error)

	t.Run("handles case where transaction is already known somehow", func(t *testing.T) {
		expectedBumpedGasPrice := big.NewInt(50400000000)
		require.Greater(t, expectedBumpedGasPrice.Int64(), attempt3_1.GasPrice.ToInt().Int64())

		ethTx := gethTypes.Transaction{}
		kst.On("SignTx",
			mock.AnythingOfType("accounts.Account"),
			mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				if int64(tx.Nonce()) != *etx3.Nonce || expectedBumpedGasPrice.Cmp(tx.GasPrice()) != 0 {
					return false
				}
				ethTx = *tx
				return true
			}),
			mock.Anything).Return(&ethTx, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return int64(tx.Nonce()) == *etx3.Nonce && expectedBumpedGasPrice.Cmp(tx.GasPrice()) == 0
		})).Return(fmt.Errorf("known transaction: %s", ethTx.Hash().Hex())).Once()

		// Do the thing
		require.NoError(t, ec.BumpGasWhereNecessary(currentHead))

		etx3, err = store.FindEthTxWithAttempts(etx3.ID)
		require.NoError(t, err)

		assert.Equal(t, models.EthTxUnconfirmed, etx3.State)

		require.Len(t, etx3.EthTxAttempts, 3)
		attempt3_3 := etx3.EthTxAttempts[2]
		assert.Equal(t, expectedBumpedGasPrice.Int64(), attempt3_3.GasPrice.ToInt().Int64())
	})
	kst.AssertExpectations(t)
	ethClient.AssertExpectations(t)
}

func TestEthConfirmer_EnsureConfirmedTransactionsInLongestChain(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	ec := bulletprooftxmanager.NewEthConfirmer(store, config)

	head := models.Head{
		Hash:   cltest.NewHash(),
		Number: 10,
		Parent: &models.Head{
			Hash:   cltest.NewHash(),
			Number: 9,
			Parent: &models.Head{
				Number: 8,
				Hash:   cltest.NewHash(),
				Parent: nil,
			},
		},
	}

	t.Run("does nothing if there aren't any transactions", func(t *testing.T) {
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(head))
	})

	t.Run("does nothing to unconfirmed transactions", func(t *testing.T) {
		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 0)

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(head))

		etx, err := store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, models.EthTxUnconfirmed, etx.State)
	})

	t.Run("does nothing to confirmed transactions with receipts within head height of the chain and included in the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 2, 1)
		attempt := etx.EthTxAttempts[0]
		cltest.MustInsertEthReceipt(t, store, head.Number, head.Hash, attempt.Hash)

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(head))

		etx, err := store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, models.EthTxConfirmed, etx.State)
	})

	t.Run("does nothing to confirmed transactions that only have receipts older than the start of the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 3, 1)
		attempt := etx.EthTxAttempts[0]
		// Add receipt that is older than the lowest block of the chain
		cltest.MustInsertEthReceipt(t, store, head.Parent.Parent.Number-1, cltest.NewHash(), attempt.Hash)

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(head))

		etx, err := store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, models.EthTxConfirmed, etx.State)
	})

	t.Run("unconfirms and rebroadcasts transactions that have receipts within head height of the chain but not included in the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 4, 1)
		attempt := etx.EthTxAttempts[0]
		// Include one within head height but a different block hash
		cltest.MustInsertEthReceipt(t, store, head.Parent.Number, cltest.NewHash(), attempt.Hash)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			atx, err := attempt.GetSignedTx()
			require.NoError(t, err)
			// Keeps gas price and nonce the same
			return atx.GasPrice().Cmp(tx.GasPrice()) == 0 && atx.Nonce() == tx.Nonce()
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(head))

		etx, err := store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, models.EthTxUnconfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt = etx.EthTxAttempts[0]
		assert.Equal(t, models.EthTxAttemptBroadcast, attempt.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("unconfirms and rebroadcasts transactions that have receipts within head height of chain but not included in the chain even if a receipt exists older than the start of the chain", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 5, 1)
		attempt := etx.EthTxAttempts[0]
		// Add receipt that is older than the lowest block of the chain
		cltest.MustInsertEthReceipt(t, store, head.Parent.Parent.Number-1, cltest.NewHash(), attempt.Hash)
		// Include one within head height but a different block hash
		cltest.MustInsertEthReceipt(t, store, head.Parent.Number, cltest.NewHash(), attempt.Hash)

		ethClient.On("SendTransaction", mock.Anything, mock.Anything).Return(nil).Once()

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(head))

		etx, err := store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, models.EthTxUnconfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt = etx.EthTxAttempts[0]
		assert.Equal(t, models.EthTxAttemptBroadcast, attempt.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("if more than one attempt has a receipt (unlikely but allowed within constraints of system, and possible in the event of forks) unconfirms and rebroadcasts only the attempt with the highest gas price", func(t *testing.T) {
		etx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 6, 1)
		require.Len(t, etx.EthTxAttempts, 1)
		// Sanity check to assert the included attempt has the lowest gas price
		require.Less(t, etx.EthTxAttempts[0].GasPrice.ToInt().Int64(), int64(30000))

		attempt2 := newBroadcastEthTxAttempt(t, etx.ID, store, 30000)
		attempt2.SignedRawTx = hexutil.MustDecode("0xf88c8301f3a98503b9aca000832ab98094f5fff180082d6017036b771ba883025c654bc93580a4daa6d556000000000000000000000000000000000000000000000000000000000000000026a0f25601065ee369b6470c0399a2334afcfbeb0b5c8f3d9a9042e448ed29b5bcbda05b676e00248b85faf4dd889f0e2dcf91eb867e23ac9eeb14a73f9e4c14972cdf")
		attempt3 := newBroadcastEthTxAttempt(t, etx.ID, store, 40000)
		attempt3.SignedRawTx = hexutil.MustDecode("0xf88c8301f3a88503b9aca0008316e36094151445852b0cfdf6a4cc81440f2af99176e8ad0880a4daa6d556000000000000000000000000000000000000000000000000000000000000000026a0dcb5a7ad52b96a866257134429f944c505820716567f070e64abb74899803855a04c13eff2a22c218e68da80111e1bb6dc665d3dea7104ab40ff8a0275a99f630d")
		require.NoError(t, store.DB.Create(&attempt2).Error)
		require.NoError(t, store.DB.Create(&attempt3).Error)

		// Receipt is within head height but a different block hash
		cltest.MustInsertEthReceipt(t, store, head.Parent.Number, cltest.NewHash(), attempt2.Hash)
		// Receipt is within head height but a different block hash
		cltest.MustInsertEthReceipt(t, store, head.Parent.Number, cltest.NewHash(), attempt3.Hash)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			s, err := attempt3.GetSignedTx()
			require.NoError(t, err)
			return tx.Hash() == s.Hash()
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, ec.EnsureConfirmedTransactionsInLongestChain(head))

		etx, err := store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		assert.Equal(t, models.EthTxUnconfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 3)
		attempt1 := etx.EthTxAttempts[0]
		assert.Equal(t, models.EthTxAttemptBroadcast, attempt1.State)
		attempt2 = etx.EthTxAttempts[1]
		assert.Equal(t, models.EthTxAttemptBroadcast, attempt2.State)
		attempt3 = etx.EthTxAttempts[2]
		assert.Equal(t, models.EthTxAttemptBroadcast, attempt3.State)

		ethClient.AssertExpectations(t)
	})
}

func TestEthConfirmer_ForceRebroadcast(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.KeyStore.Unlock(cltest.Password)
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	mustInsertUnstartedEthTx(t, store)
	mustInsertInProgressEthTx(t, store, 0)
	etx1 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 1)
	etx2 := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 2)

	gasPriceWei := uint64(52)
	address := cltest.GetDefaultFromAddress(t, store)
	overrideGasLimit := uint64(20000)

	t.Run("rebroadcasts one eth_tx if it falls within in nonce range", func(t *testing.T) {
		ethClient := new(mocks.Client)
		store.EthClient = ethClient
		ec := bulletprooftxmanager.NewEthConfirmer(store, config)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Nonce) &&
				uint64(tx.GasPrice().Int64()) == gasPriceWei &&
				tx.Gas() == overrideGasLimit &&
				reflect.DeepEqual(tx.Data(), etx1.EncodedPayload) &&
				*tx.To() == etx1.ToAddress
		})).Return(nil).Once()

		require.NoError(t, ec.ForceRebroadcast(1, 1, gasPriceWei, address, overrideGasLimit))

		ethClient.AssertExpectations(t)
	})

	t.Run("uses default gas limit if overrideGasLimit is 0", func(t *testing.T) {
		ethClient := new(mocks.Client)
		store.EthClient = ethClient
		ec := bulletprooftxmanager.NewEthConfirmer(store, config)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Nonce) &&
				uint64(tx.GasPrice().Int64()) == gasPriceWei &&
				tx.Gas() == etx1.GasLimit &&
				reflect.DeepEqual(tx.Data(), etx1.EncodedPayload) &&
				*tx.To() == etx1.ToAddress
		})).Return(nil).Once()

		require.NoError(t, ec.ForceRebroadcast(1, 1, gasPriceWei, address, 0))

		ethClient.AssertExpectations(t)
	})

	t.Run("rebroadcasts several eth_txes in nonce range", func(t *testing.T) {
		ethClient := new(mocks.Client)
		store.EthClient = ethClient
		ec := bulletprooftxmanager.NewEthConfirmer(store, config)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(*etx1.Nonce) && uint64(tx.GasPrice().Int64()) == gasPriceWei && tx.Gas() == overrideGasLimit
		})).Return(nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(*etx2.Nonce) && uint64(tx.GasPrice().Int64()) == gasPriceWei && tx.Gas() == overrideGasLimit
		})).Return(nil).Once()

		require.NoError(t, ec.ForceRebroadcast(1, 2, gasPriceWei, address, overrideGasLimit))

		ethClient.AssertExpectations(t)
	})

	t.Run("broadcasts zero transactions if eth_tx doesn't exist for that nonce", func(t *testing.T) {
		ethClient := new(mocks.Client)
		store.EthClient = ethClient
		ec := bulletprooftxmanager.NewEthConfirmer(store, config)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(1)
		})).Return(nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(2)
		})).Return(nil).Once()
		for i := 3; i <= 5; i++ {
			nonce := i
			ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
				return tx.Nonce() == uint64(nonce) &&
					uint64(tx.GasPrice().Int64()) == gasPriceWei &&
					tx.Gas() == overrideGasLimit &&
					*tx.To() == utils.ZeroAddress &&
					tx.Value().Cmp(big.NewInt(0)) == 0 &&
					len(tx.Data()) == 0
			})).Return(nil).Once()
		}

		require.NoError(t, ec.ForceRebroadcast(1, 5, gasPriceWei, address, overrideGasLimit))

		ethClient.AssertExpectations(t)
	})

	t.Run("zero transactions use default gas limit if override wasn't specified", func(t *testing.T) {
		ethClient := new(mocks.Client)
		store.EthClient = ethClient
		ec := bulletprooftxmanager.NewEthConfirmer(store, config)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(0) && uint64(tx.GasPrice().Int64()) == gasPriceWei && uint64(tx.Gas()) == config.EthGasLimitDefault()
		})).Return(nil).Once()

		require.NoError(t, ec.ForceRebroadcast(0, 0, gasPriceWei, address, 0))

		ethClient.AssertExpectations(t)
	})
}
