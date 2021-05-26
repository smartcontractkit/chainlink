package bulletprooftxmanager_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	gormpostgrestypes "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/store/dialects"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	gethCommon "github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
)

func TestEthBroadcaster_ProcessUnstartedEthTxs_Success(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, store, 0)
	store.KeyStore.Unlock(cltest.Password)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	eb, cleanup := cltest.NewEthBroadcaster(t, store, config, key)
	defer cleanup()

	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	timeNow := time.Now()

	encodedPayload := []byte{1, 2, 3}
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)

	t.Run("no eth_txes at all", func(t *testing.T) {
		require.NoError(t, eb.ProcessUnstartedEthTxs(key))
	})

	t.Run("eth_txes exist for a different from address", func(t *testing.T) {
		_, otherAddress := cltest.MustAddRandomKeyToKeystore(t, store)

		etx := models.EthTx{
			FromAddress:    otherAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          models.EthTxUnstarted,
		}
		require.NoError(t, store.DB.Save(&etx).Error)

		require.NoError(t, eb.ProcessUnstartedEthTxs(key))
	})

	t.Run("existing eth_txes with broadcast_at or error", func(t *testing.T) {
		nonce := int64(342)
		errStr := "some error"

		etxUnconfirmed := models.EthTx{
			Nonce:          &nonce,
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			BroadcastAt:    &timeNow,
			Error:          nil,
			State:          models.EthTxUnconfirmed,
		}
		etxWithError := models.EthTx{
			Nonce:          nil,
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			Error:          &errStr,
			State:          models.EthTxFatalError,
		}

		require.NoError(t, store.DB.Save(&etxUnconfirmed).Error)
		require.NoError(t, store.DB.Save(&etxWithError).Error)

		require.NoError(t, eb.ProcessUnstartedEthTxs(key))
	})

	t.Run("sends 3 EthTxs in order with higher value last, and lower values starting from the earliest", func(t *testing.T) {
		// Higher value
		expensiveEthTx := models.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 0},
			Value:          assets.NewEthValue(242),
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 0),
			State:          models.EthTxUnstarted,
		}
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(2) && tx.Value().Cmp(big.NewInt(242)) == 0
		})).Return(nil).Once()

		// Earlier
		h := gethCommon.HexToHash("0x4ea4bb19d0847a0465d003ea11bcaef62935cec3c673238d057f6cacb3e7a405")
		tr := uuid.NewV4()
		b, err := json.Marshal(models.EthTxMeta{TaskRunID: tr, RunRequestID: &h, RunRequestTxHash: &h})
		require.NoError(t, err)
		earlierEthTx := models.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 0},
			Value:          value,
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 1),
			State:          models.EthTxUnstarted,
			Meta:           gormpostgrestypes.Jsonb{RawMessage: b},
		}
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			if tx.Nonce() != uint64(0) {
				return false
			}
			require.Equal(t, config.ChainID(), tx.ChainId())
			require.Equal(t, gasLimit, tx.Gas())
			require.Equal(t, config.EthGasPriceDefault(), tx.GasPrice())
			require.Equal(t, toAddress, *tx.To())
			require.Equal(t, value.ToInt().String(), tx.Value().String())
			require.Equal(t, earlierEthTx.EncodedPayload, tx.Data())
			return true
		})).Return(nil).Once()

		// Later
		laterEthTx := models.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 1},
			Value:          value,
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(1, 0),
			State:          models.EthTxUnstarted,
		}
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			if tx.Nonce() != uint64(1) {
				return false
			}
			require.Equal(t, config.ChainID(), tx.ChainId())
			require.Equal(t, gasLimit, tx.Gas())
			require.Equal(t, config.EthGasPriceDefault(), tx.GasPrice())
			require.Equal(t, toAddress, *tx.To())
			require.Equal(t, value.ToInt().String(), tx.Value().String())
			require.Equal(t, laterEthTx.EncodedPayload, tx.Data())
			return true
		})).Return(nil).Once()

		// Insertion order deliberately reversed to test ordering
		require.NoError(t, store.DB.Save(&expensiveEthTx).Error)
		require.NoError(t, store.DB.Save(&laterEthTx).Error)
		require.NoError(t, store.DB.Save(&earlierEthTx).Error)

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(key))

		// Check earlierEthTx and it's attempt
		// This was the earlier one sent so it has the lower nonce
		earlierTransaction, err := store.FindEthTxWithAttempts(earlierEthTx.ID)
		require.NoError(t, err)
		assert.Nil(t, earlierTransaction.Error)
		require.NotNil(t, earlierTransaction.FromAddress)
		assert.Equal(t, fromAddress, earlierTransaction.FromAddress)
		require.NotNil(t, earlierTransaction.Nonce)
		assert.Equal(t, int64(0), *earlierTransaction.Nonce)
		assert.NotNil(t, earlierTransaction.BroadcastAt)
		assert.Len(t, earlierTransaction.EthTxAttempts, 1)
		var m models.EthTxMeta
		err = json.Unmarshal(earlierEthTx.Meta.RawMessage, &m)
		require.NoError(t, err)
		assert.Equal(t, tr, m.TaskRunID)
		assert.Equal(t, h.String(), m.RunRequestTxHash.String())
		assert.Equal(t, h.String(), m.RunRequestID.String())

		attempt := earlierTransaction.EthTxAttempts[0]

		assert.Equal(t, earlierTransaction.ID, attempt.EthTxID)
		assert.Equal(t, config.EthGasPriceDefault().String(), attempt.GasPrice.String())

		_, err = attempt.GetSignedTx()
		require.NoError(t, err)
		assert.Equal(t, models.EthTxAttemptBroadcast, attempt.State)
		require.Len(t, attempt.EthReceipts, 0)

		// Check laterEthTx and it's attempt
		// This was the later one sent so it has the higher nonce
		laterTransaction, err := store.FindEthTxWithAttempts(laterEthTx.ID)
		require.NoError(t, err)
		assert.Nil(t, laterTransaction.Error)
		require.NotNil(t, laterTransaction.FromAddress)
		assert.Equal(t, fromAddress, laterTransaction.FromAddress)
		require.NotNil(t, laterTransaction.Nonce)
		assert.Equal(t, int64(1), *laterTransaction.Nonce)
		assert.NotNil(t, laterTransaction.BroadcastAt)
		assert.Len(t, laterTransaction.EthTxAttempts, 1)

		attempt = laterTransaction.EthTxAttempts[0]

		assert.Equal(t, laterTransaction.ID, attempt.EthTxID)
		assert.Equal(t, config.EthGasPriceDefault().String(), attempt.GasPrice.String())

		_, err = attempt.GetSignedTx()
		require.NoError(t, err)
		assert.Equal(t, models.EthTxAttemptBroadcast, attempt.State)
		require.Len(t, attempt.EthReceipts, 0)

		ethClient.AssertExpectations(t)
	})
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_Success_OnOptimism(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, store, 0)
	store.KeyStore.Unlock(cltest.Password)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	config.Set("OPTIMISM_GAS_FEES", "true")

	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	eb, cleanup := cltest.NewEthBroadcaster(t, store, config, key)
	defer cleanup()

	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")

	estimatedGas := uint64(9007199254740993)

	tx := models.EthTx{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: []byte{42, 42, 0},
		Value:          assets.NewEthValue(242),
		GasLimit:       estimatedGas,
		CreatedAt:      time.Unix(0, 0),
		State:          models.EthTxUnstarted,
	}
	ethClient.On("EstimateGas", mock.Anything, mock.Anything).Return(estimatedGas, nil).Once()
	ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
		assert.Equal(t, big.NewInt(1000000000), tx.GasPrice())
		assert.Equal(t, estimatedGas, tx.Gas())
		return true
	})).Return(nil).Once()

	require.NoError(t, store.DB.Save(&tx).Error)

	// Do the thing
	require.NoError(t, eb.ProcessUnstartedEthTxs(key))
	ethClient.AssertExpectations(t)
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_Success_WithMultiplier(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, store, 0)
	store.KeyStore.Unlock(cltest.Password)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	config.Set("ETH_GAS_LIMIT_MULTIPLIER", "1.3")

	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	eb, cleanup := cltest.NewEthBroadcaster(t, store, config, key)
	defer cleanup()

	ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
		assert.Equal(t, uint64(1600), tx.Gas())
		return true
	})).Return(nil).Once()

	tx := models.EthTx{
		FromAddress:    fromAddress,
		ToAddress:      gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411"),
		EncodedPayload: []byte{42, 42, 0},
		Value:          assets.NewEthValue(242),
		GasLimit:       1231,
		CreatedAt:      time.Unix(0, 0),
		State:          models.EthTxUnstarted,
	}
	require.NoError(t, store.DB.Save(&tx).Error)

	// Do the thing
	require.NoError(t, eb.ProcessUnstartedEthTxs(key))
	ethClient.AssertExpectations(t)
}

func TestEthBroadcaster_AssignsNonceOnStart(t *testing.T) {
	var err error
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	k1, fromAddress := cltest.MustAddRandomKeyToKeystore(t, store, true)
	k2, dummyAddress := cltest.MustAddRandomKeyToKeystore(t, store, false)
	keys := []models.Key{k1, k2}

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	config.Set("ETH_NONCE_AUTO_SYNC", "true")

	ethNodeNonce := uint64(22)

	t.Run("when eth node returns error", func(t *testing.T) {
		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		eb, cleanup := cltest.NewEthBroadcaster(t, store, config, keys...)
		defer cleanup()

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(account gethCommon.Address) bool {
			return account.Hex() == dummyAddress.Hex()
		})).Return(uint64(0), nil).Once()
		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(account gethCommon.Address) bool {
			return account.Hex() == fromAddress.Hex()
		})).Return(ethNodeNonce, errors.New("something exploded")).Once()

		err = eb.Start()
		require.Error(t, err)
		defer eb.Close()
		require.Contains(t, err.Error(), "something exploded")

		// dummy address got updated
		var n int
		err := store.DB.Raw(`SELECT next_nonce FROM keys WHERE address = ?`, dummyAddress).Scan(&n).Error
		require.NoError(t, err)
		require.Equal(t, 0, n)

		// real address did not update (it errored)
		err = store.DB.Raw(`SELECT next_nonce FROM keys WHERE address = ?`, fromAddress).Scan(&n).Error
		require.NoError(t, err)
		require.Equal(t, 0, n)

		ethClient.AssertExpectations(t)
	})

	t.Run("when eth node returns nonce", func(t *testing.T) {
		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		eb, cleanup := cltest.NewEthBroadcaster(t, store, config, keys...)
		defer cleanup()

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(account gethCommon.Address) bool {
			return account.Hex() == dummyAddress.Hex()
		})).Return(uint64(0), nil).Once()
		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(account gethCommon.Address) bool {
			return account.Hex() == fromAddress.Hex()
		})).Return(ethNodeNonce, nil).Once()

		require.NoError(t, eb.Start())
		defer eb.Close()

		// Check key to make sure it has correct nonce assigned
		var keys []models.Key
		err := store.DB.Order("created_at asc").Find(&keys).Error
		require.NoError(t, err)
		key := keys[0]

		assert.NotNil(t, key.NextNonce)
		assert.Equal(t, int64(ethNodeNonce), key.NextNonce)

		// The dummy key did not get updated
		key2 := keys[1]
		assert.Equal(t, dummyAddress.Hex(), key2.Address.Hex())
		assert.Equal(t, 0, int(key2.NextNonce))

		ethClient.AssertExpectations(t)
	})
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_ResumingFromCrash(t *testing.T) {
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)
	encodedPayload := []byte{0, 1}
	nextNonce := int64(916714082576372851)
	firstNonce := nextNonce
	secondNonce := nextNonce + 1

	t.Run("cannot be more than one transaction per address in an unfinished state", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, store, nextNonce)

		firstInProgress := models.EthTx{
			FromAddress:    fromAddress,
			Nonce:          &firstNonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			BroadcastAt:    nil,
			Error:          nil,
			State:          models.EthTxInProgress,
		}

		secondInProgress := models.EthTx{
			FromAddress:    fromAddress,
			Nonce:          &secondNonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			BroadcastAt:    nil,
			Error:          nil,
			State:          models.EthTxInProgress,
		}

		require.NoError(t, store.DB.Create(&firstInProgress).Error)
		err := store.DB.Create(&secondInProgress).Error
		require.Error(t, err)
		assert.EqualError(t, err, "ERROR: duplicate key value violates unique constraint \"idx_only_one_in_progress_tx_per_account\" (SQLSTATE 23505)")
	})

	t.Run("previous run assigned nonce but never broadcast", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, store, nextNonce)

		config, cleanup := cltest.NewConfig(t)
		defer cleanup()

		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		eb, cleanup := cltest.NewEthBroadcaster(t, store, config, key)
		defer cleanup()

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, store, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(key))

		// Check it was saved correctly with its attempt
		etx, err := store.FindEthTxWithAttempts(inProgressEthTx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTxAttempts, 1)
		assert.Equal(t, models.EthTxAttemptBroadcast, etx.EthTxAttempts[0].State)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast but it fatally errored before we could save", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, store, nextNonce)
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		defer cleanup()

		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		eb, cleanup := cltest.NewEthBroadcaster(t, store, config, key)
		defer cleanup()

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, store, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(errors.New("exceeds block gas limit")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(key))

		// Check it was saved correctly with its attempt
		etx, err := store.FindEthTxWithAttempts(inProgressEthTx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.Error)
		assert.Equal(t, "exceeds block gas limit", *etx.Error)
		assert.Len(t, etx.EthTxAttempts, 0)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast and is now in mempool", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, store, nextNonce)

		config, cleanup := cltest.NewConfig(t)
		defer cleanup()

		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		eb, cleanup := cltest.NewEthBroadcaster(t, store, config, key)
		defer cleanup()

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, store, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(errors.New("known transaction: a1313bd99a81fb4d8ad1d2e90b67c6b3fa77545c990d6251444b83b70b6f8980")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(key))

		// Check it was saved correctly with its attempt
		etx, err := store.FindEthTxWithAttempts(inProgressEthTx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTxAttempts, 1)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast and now the transaction has been confirmed", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, store, nextNonce)
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		defer cleanup()

		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		eb, cleanup := cltest.NewEthBroadcaster(t, store, config, key)
		defer cleanup()

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, store, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(errors.New("nonce too low")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(key))

		// Check it was saved correctly with its attempt
		etx, err := store.FindEthTxWithAttempts(inProgressEthTx.ID)
		require.NoError(t, err)

		require.NotNil(t, etx.BroadcastAt)
		assert.Equal(t, *etx.BroadcastAt, etx.CreatedAt)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTxAttempts, 1)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and then failed to reach node for some reason and node is still down", func(t *testing.T) {
		failedToReachNodeError := context.DeadlineExceeded
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, store, nextNonce)
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		defer cleanup()

		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		eb, cleanup := cltest.NewEthBroadcaster(t, store, config, key)
		defer cleanup()

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, store, firstNonce, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(firstNonce)
		})).Return(failedToReachNodeError).Once()

		// Do the thing
		err := eb.ProcessUnstartedEthTxs(key)
		require.Error(t, err)
		assert.Contains(t, err.Error(), failedToReachNodeError.Error())

		// Check it was left in the unfinished state
		etx, err := store.FindEthTxWithAttempts(inProgressEthTx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Equal(t, nextNonce, *etx.Nonce)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTxAttempts, 1)

		ethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast transaction then crashed and rebooted with a different configured gas price", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, store, nextNonce)
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		defer cleanup()

		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		// Configured gas price changed
		store.Config.Set("ETH_GAS_PRICE_DEFAULT", 500000000000)

		eb, cleanup := cltest.NewEthBroadcaster(t, store, config, key)
		defer cleanup()

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		inProgressEthTx := cltest.MustInsertInProgressEthTxWithAttempt(t, store, firstNonce, fromAddress)
		require.Len(t, inProgressEthTx.EthTxAttempts, 1)
		attempt := inProgressEthTx.EthTxAttempts[0]

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			// Ensure that the gas price is the same as the original attempt
			s, e := attempt.GetSignedTx()
			require.NoError(t, e)
			return tx.Nonce() == uint64(firstNonce) && tx.GasPrice().Int64() == s.GasPrice().Int64()
		})).Return(errors.New("known transaction: a1313bd99a81fb4d8ad1d2e90b67c6b3fa77545c990d6251444b83b70b6f8980")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(key))

		// Check it was saved correctly with its attempt
		etx, err := store.FindEthTxWithAttempts(inProgressEthTx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt = etx.EthTxAttempts[0]
		s, err := attempt.GetSignedTx()
		require.NoError(t, err)
		assert.Equal(t, int64(342), s.GasPrice().Int64())
		assert.Equal(t, models.EthTxAttemptBroadcast, attempt.State)

		ethClient.AssertExpectations(t)
	})
}

func getLocalNextNonce(t *testing.T, str *store.Store, fromAddress gethCommon.Address) uint64 {
	n, err := bulletprooftxmanager.GetNextNonce(str.DB, fromAddress)
	require.NoError(t, err)
	require.NotNil(t, n)
	return uint64(n)
}

// Note that all of these tests share the same database, and ordering matters.
// This in order to more deeply test ProcessUnstartedEthTxs over
// multiple runs with previous errors in the database.
func TestEthBroadcaster_ProcessUnstartedEthTxs_Errors(t *testing.T) {
	var err error
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)
	encodedPayload := []byte{0, 1}

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	key, fromAddress := cltest.MustAddRandomKeyToKeystore(t, store, 0)
	store.KeyStore.Unlock(cltest.Password)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	eb, cleanup := cltest.NewEthBroadcaster(t, store, config, key)
	defer cleanup()

	t.Run("if external wallet sent a transaction from the account and now the nonce is one higher than it should be and we got replacement underpriced then we assume a previous transaction of ours was the one that succeeded, and hand off to EthConfirmer", func(t *testing.T) {
		etx := models.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          models.EthTxUnstarted,
		}
		require.NoError(t, store.DB.Save(&etx).Error)
		taskRunID, _ := cltest.MustInsertTaskRun(t, store)
		_, err = store.MustSQLDB().Exec(`INSERT INTO eth_task_run_txes (task_run_id, eth_tx_id) VALUES ($1, $2)`, taskRunID, etx.ID)
		require.NoError(t, err)

		// First send, replacement underpriced
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(0)
		})).Return(errors.New("replacement transaction underpriced")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(key))

		ethClient.AssertExpectations(t)

		// Check that the transaction was saved correctly with its attempt
		// We assume success and hand off to eth confirmer to eventually mark it as failed
		var latestID int64
		var etx1 models.EthTx
		require.NoError(t, store.DB.Raw("SELECT max(id) FROM eth_txes").Row().Scan(&latestID))
		etx1, err = store.FindEthTxWithAttempts(latestID)
		require.NoError(t, err)
		require.NotNil(t, etx1.BroadcastAt)
		assert.NotEqual(t, etx1.CreatedAt, *etx1.BroadcastAt)
		require.NotNil(t, etx1.Nonce)
		assert.Equal(t, int64(0), *etx1.Nonce)
		assert.Nil(t, etx1.Error)
		assert.Len(t, etx1.EthTxAttempts, 1)

		// Check that the local nonce was incremented by one
		var finalNextNonce int64
		finalNextNonce, err = bulletprooftxmanager.GetNextNonce(store.DB, fromAddress)
		require.NoError(t, err)
		require.NotNil(t, finalNextNonce)
		require.Equal(t, int64(1), finalNextNonce)
	})

	t.Run("geth client returns an error in the fatal errors category", func(t *testing.T) {
		fatalErrorExample := "exceeds block gas limit"
		localNextNonce := getLocalNextNonce(t, store, fromAddress)

		etx := models.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          models.EthTxUnstarted,
		}
		require.NoError(t, store.DB.Save(&etx).Error)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(fatalErrorExample)).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(key))

		// Check it was saved correctly with its attempt
		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		require.Nil(t, etx.Nonce)
		assert.NotNil(t, etx.Error)
		assert.Contains(t, *etx.Error, "exceeds block gas limit")
		assert.Len(t, etx.EthTxAttempts, 0)

		// Check that the key had its nonce reset
		require.NoError(t, store.DB.First(&key).Error)
		// Saved NextNonce must be the same as before because this transaction
		// was not accepted by the eth node and never can be
		require.NotNil(t, key.NextNonce)
		require.Equal(t, int64(localNextNonce), key.NextNonce)

		ethClient.AssertExpectations(t)
	})

	t.Run("geth client fails with error indicating that the transaction was too expensive", func(t *testing.T) {
		tooExpensiveError := "tx fee (1.10 ether) exceeds the configured cap (1.00 ether)"
		localNextNonce := getLocalNextNonce(t, store, fromAddress)

		etx := models.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          models.EthTxUnstarted,
		}
		require.NoError(t, store.DB.Save(&etx).Error)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(tooExpensiveError)).Once()

		require.NoError(t, eb.ProcessUnstartedEthTxs(key))

		// Check it was saved with no attempt and a fatal error
		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		require.Nil(t, etx.Nonce)
		assert.NotNil(t, etx.Error)
		assert.Contains(t, *etx.Error, "tx fee (1.10 ether) exceeds the configured cap (1.00 ether)")
		assert.Len(t, etx.EthTxAttempts, 0)

		// Check that the key had its nonce reset
		require.NoError(t, store.DB.First(&key).Error)
		// Saved NextNonce must be the same as before because this transaction
		// was not accepted by the eth node and never can be
		require.NotNil(t, key.NextNonce)
		require.Equal(t, int64(localNextNonce), key.NextNonce)

		ethClient.AssertExpectations(t)
	})

	t.Run("eth client call fails with an unexpected random error", func(t *testing.T) {
		retryableErrorExample := "geth shit the bed again"
		localNextNonce := getLocalNextNonce(t, store, fromAddress)

		etx := models.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          models.EthTxUnstarted,
		}
		require.NoError(t, store.DB.Save(&etx).Error)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(retryableErrorExample)).Once()

		// Do the thing
		err = eb.ProcessUnstartedEthTxs(key)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("error while sending transaction %v: %s", etx.ID, retryableErrorExample))

		// Check it was saved correctly with its attempt
		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.Nil(t, etx.Error)
		assert.Equal(t, models.EthTxInProgress, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]
		assert.Equal(t, models.EthTxAttemptInProgress, attempt.State)

		ethClient.AssertExpectations(t)

		// Now on the second run, it is successful
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(nil).Once()

		require.NoError(t, eb.ProcessUnstartedEthTxs(key))

		// Check it was saved correctly with its attempt
		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.Nil(t, etx.Error)
		assert.Equal(t, models.EthTxUnconfirmed, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt = etx.EthTxAttempts[0]
		assert.Equal(t, models.EthTxAttemptBroadcast, attempt.State)

		ethClient.AssertExpectations(t)
	})

	t.Run("eth node returns underpriced transaction", func(t *testing.T) {
		// This happens if a transaction's gas price is below the minimum
		// configured for the transaction pool.
		// This is a configuration error by the node operator, since it means they set the base gas level too low.
		underpricedError := "transaction underpriced"
		localNextNonce := getLocalNextNonce(t, store, fromAddress)

		etx := models.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          models.EthTxUnstarted,
		}
		require.NoError(t, store.DB.Save(&etx).Error)

		// First was underpriced
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasPrice().Cmp(store.Config.EthGasPriceDefault()) == 0
		})).Return(errors.New(underpricedError)).Once()

		// Second with gas bump was still underpriced
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasPrice().Cmp(big.NewInt(25000000000)) == 0
		})).Return(errors.New(underpricedError)).Once()

		// Third succeeded
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasPrice().Cmp(big.NewInt(30000000000)) == 0
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(key))

		ethClient.AssertExpectations(t)

		// Check it was saved correctly with its attempt
		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]
		assert.Equal(t, big.NewInt(30000000000).String(), attempt.GasPrice.String())
	})

	etxUnfinished := models.EthTx{
		FromAddress:    fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: encodedPayload,
		Value:          value,
		GasLimit:       gasLimit,
		State:          models.EthTxUnstarted,
	}
	require.NoError(t, store.DB.Save(&etxUnfinished).Error)

	t.Run("failed to reach node for some reason", func(t *testing.T) {
		failedToReachNodeError := context.DeadlineExceeded
		localNextNonce := getLocalNextNonce(t, store, fromAddress)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(failedToReachNodeError).Once()

		// Do the thing
		err = eb.ProcessUnstartedEthTxs(key)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("error while sending transaction %v: context deadline exceeded", etxUnfinished.ID))

		// Check it was left in the unfinished state
		etx, err := store.FindEthTxWithAttempts(etxUnfinished.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.Nonce)
		assert.Nil(t, etx.Error)
		assert.Equal(t, models.EthTxInProgress, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		assert.Equal(t, models.EthTxAttemptInProgress, etx.EthTxAttempts[0].State)

		ethClient.AssertExpectations(t)
	})

	t.Run("eth node returns temporarily underpriced transaction", func(t *testing.T) {
		// This happens if parity is rejecting transactions that are not priced high enough to even get into the mempool at all
		// It should pretend it was accepted into the mempool and hand off to ethConfirmer to bump gas as normal
		temporarilyUnderpricedError := "There are too many transactions in the queue. Your transaction was dropped due to limit. Try increasing the fee."
		localNextNonce := getLocalNextNonce(t, store, fromAddress)

		// Re-use the previously unfinished transaction, no need to insert new

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(temporarilyUnderpricedError)).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(key))

		// Check it was saved correctly with its attempt
		etx, err := store.FindEthTxWithAttempts(etxUnfinished.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]
		assert.Equal(t, big.NewInt(20000000000).String(), attempt.GasPrice.String())

		ethClient.AssertExpectations(t)
	})

	t.Run("eth node returns underpriced transaction and bumping gas doesn't increase it", func(t *testing.T) {
		// This happens if a transaction's gas price is below the minimum
		// configured for the transaction pool.
		// This is a configuration error by the node operator, since it means they set the base gas level too low.
		underpricedError := "transaction underpriced"
		localNextNonce := getLocalNextNonce(t, store, fromAddress)
		// In this scenario the node operator REALLY fucked up and set the bump
		// to zero (even though that should not be possible due to config
		// validation)
		config.Set("ETH_GAS_BUMP_WEI", "0")
		config.Set("ETH_GAS_BUMP_PERCENT", "0")

		etx := models.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          models.EthTxUnstarted,
		}
		require.NoError(t, store.DB.Save(&etx).Error)

		// First was underpriced
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce && tx.GasPrice().Cmp(store.Config.EthGasPriceDefault()) == 0
		})).Return(errors.New(underpricedError)).Once()

		// Do the thing
		err := eb.ProcessUnstartedEthTxs(key)
		require.Error(t, err)
		require.Contains(t, err.Error(), "bumped gas price of 20000000000 is equal to original gas price of 20000000000. ACTION REQUIRED: This is a configuration error, you must increase either ETH_GAS_BUMP_PERCENT or ETH_GAS_BUMP_WEI")

		// TEARDOWN: Clear out the unsent tx before the next test
		require.NoError(t, store.DB.Exec(`DELETE FROM eth_txes WHERE nonce = ?`, localNextNonce).Error)

		ethClient.AssertExpectations(t)
	})

	t.Run("eth node returns insufficient eth", func(t *testing.T) {
		insufficientEthError := "insufficient funds for transfer"
		localNextNonce := getLocalNextNonce(t, store, fromAddress)
		etx := models.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          models.EthTxUnstarted,
		}
		require.NoError(t, store.DB.Save(&etx).Error)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(insufficientEthError)).Once()

		err := eb.ProcessUnstartedEthTxs(key)
		require.NoError(t, err)

		// Check it was saved correctly with its attempt
		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.Nil(t, etx.Error)
		assert.Equal(t, models.EthTxUnconfirmed, etx.State)
		require.Len(t, etx.EthTxAttempts, 1)
		attempt := etx.EthTxAttempts[0]
		assert.Equal(t, models.EthTxAttemptInsufficientEth, attempt.State)
		assert.Nil(t, attempt.BroadcastBeforeBlockNum)

		ethClient.AssertExpectations(t)
	})
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_KeystoreErrors(t *testing.T) {
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)
	encodedPayload := []byte{0, 1}
	localNonce := 0

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	kst := new(mocks.KeyStoreInterface)
	store.KeyStore = kst

	key := cltest.MustInsertRandomKey(t, store.DB, 0)
	fromAddress := key.Address.Address()

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	eb, cleanup := cltest.NewEthBroadcaster(t, store, config, key)
	defer cleanup()

	t.Run("tx signing fails", func(t *testing.T) {
		etx := models.EthTx{
			FromAddress:    fromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          models.EthTxUnstarted,
		}
		require.NoError(t, store.DB.Save(&etx).Error)

		tx := *gethTypes.NewTx(&gethTypes.LegacyTx{})
		kst.On("SignTx",
			fromAddress,
			mock.AnythingOfType("*types.Transaction"),
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Cmp(store.Config.ChainID()) == 0
			})).Return(&tx, errors.New("could not sign transaction")).Once()

		// Do the thing
		err := eb.ProcessUnstartedEthTxs(key)
		require.Error(t, err)
		require.Contains(t, err.Error(), "could not sign transaction")

		// Check that the transaction is left in unstarted state
		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Equal(t, models.EthTxUnstarted, etx.State)
		assert.Len(t, etx.EthTxAttempts, 0)

		// Check that the key did not have its nonce incremented
		var key models.Key
		require.NoError(t, store.DB.First(&key).Error)
		require.NotNil(t, key.NextNonce)
		require.Equal(t, int64(localNonce), key.NextNonce)

		kst.AssertExpectations(t)
	})

	// Should have done nothing
	ethClient.AssertExpectations(t)
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_Locking(t *testing.T) {
	advisoryLocker1 := new(mocks.AdvisoryLocker)
	store, cleanup := cltest.NewStore(t, advisoryLocker1)
	defer cleanup()
	key, _ := cltest.MustAddRandomKeyToKeystore(t, store, 0)

	advisoryLocker1.On("WithAdvisoryLock", mock.Anything, mock.AnythingOfType("int32"), key.ID, mock.AnythingOfType("func() error")).Return(nil)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	eb := bulletprooftxmanager.NewEthBroadcaster(store.DB, store.EthClient, config, store.KeyStore, advisoryLocker1, &postgres.NullEventBroadcaster{}, []models.Key{key})

	require.NoError(t, eb.ProcessUnstartedEthTxs(key))

	advisoryLocker1.AssertExpectations(t)
	advisoryLocker1.On("Close").Return(nil)
}

func TestEthBroadcaster_GetNextNonce(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	key, _ := cltest.MustAddRandomKeyToKeystore(t, store, 0)

	nonce, err := bulletprooftxmanager.GetNextNonce(store.DB, key.Address.Address())
	assert.NoError(t, err)
	require.NotNil(t, nonce)
	assert.Equal(t, int64(0), nonce)
}

func TestEthBroadcaster_IncrementNextNonce(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	key, _ := cltest.MustAddRandomKeyToKeystore(t, store, 0)

	// Cannot increment if supplied nonce doesn't match existing
	require.Error(t, bulletprooftxmanager.IncrementNextNonce(store.DB, key.Address.Address(), int64(42)))

	require.NoError(t, bulletprooftxmanager.IncrementNextNonce(store.DB, key.Address.Address(), int64(0)))

	// Nonce bumped to 1
	require.NoError(t, store.DB.First(&key).Error)
	require.NotNil(t, key.NextNonce)
	require.Equal(t, int64(1), key.NextNonce)
}

func TestEthBroadcaster_Trigger(t *testing.T) {
	t.Parallel()

	// Simple sanity check to make sure it doesn't block
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	eb, cleanup := cltest.NewEthBroadcaster(t, store, config)
	defer cleanup()

	eb.Trigger(cltest.NewAddress())
	eb.Trigger(cltest.NewAddress())
	eb.Trigger(cltest.NewAddress())
}

func TestEthBroadcaster_EthTxInsertEventCausesTriggerToFire(t *testing.T) {
	// NOTE: Testing triggers requires committing transactions and does not work with transactional tests
	config, _, cleanup := cltest.BootstrapThrowawayORM(t, "eth_tx_triggers", true, true)
	defer cleanup()
	config.Config.Dialect = dialects.PostgresWithoutLock
	store, cleanup := cltest.NewStoreWithConfig(t, config)
	defer cleanup()
	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, store, 0)
	eventBroadcaster := postgres.NewEventBroadcaster(config.DatabaseURL(), 0, 0)
	eventBroadcaster.Start()
	defer eventBroadcaster.Close()

	ethTxInsertListener, err := eventBroadcaster.Subscribe(postgres.ChannelInsertOnEthTx, "")
	require.NoError(t, err)

	// Give it some time to start listening
	time.Sleep(100 * time.Millisecond)

	mustInsertUnstartedEthTx(t, store, fromAddress)
	gomega.NewGomegaWithT(t).Eventually(ethTxInsertListener.Events()).Should(gomega.Receive())
}
