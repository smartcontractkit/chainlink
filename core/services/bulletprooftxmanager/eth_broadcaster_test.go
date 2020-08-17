package bulletprooftxmanager_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	gethAccounts "github.com/ethereum/go-ethereum/accounts"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
)

func mustInsertInProgressEthTxWithAttempt(t *testing.T, store *store.Store, nonce int64) models.EthTx {
	etx := cltest.NewEthTx(t, store)

	etx.BroadcastAt = nil
	etx.Nonce = &nonce
	etx.State = models.EthTxInProgress
	require.NoError(t, store.DB.Save(&etx).Error)
	attempt := cltest.NewEthTxAttempt(t, etx.ID)
	tx := gethTypes.NewTransaction(uint64(nonce), cltest.NewAddress(), big.NewInt(142), 242, big.NewInt(342), []byte{1, 2, 3})
	rlp := new(bytes.Buffer)
	require.NoError(t, tx.EncodeRLP(rlp))
	attempt.SignedRawTx = rlp.Bytes()
	attempt.State = models.EthTxAttemptInProgress
	require.NoError(t, store.DB.Save(&attempt).Error)
	etx, err := store.FindEthTxWithAttempts(etx.ID)
	require.NoError(t, err)
	return etx
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_Success(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	// Use the real KeyStore loaded from database fixtures
	store.KeyStore.Unlock(cltest.Password)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	eb := bulletprooftxmanager.NewEthBroadcaster(store, config)

	keys, err := store.Keys()
	require.NoError(t, err)
	key := keys[0]
	defaultFromAddress := key.Address.Address()
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	timeNow := time.Now()

	encodedPayload := []byte{1, 2, 3}
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)

	t.Run("no eth_txes at all", func(t *testing.T) {
		require.NoError(t, eb.ProcessUnstartedEthTxs(key))
	})

	t.Run("eth_txes exist for a different from address", func(t *testing.T) {
		otherAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
		cltest.MustInsertKey(t, store, otherAddress)

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
			FromAddress:    defaultFromAddress,
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
			FromAddress:    defaultFromAddress,
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
			FromAddress:    defaultFromAddress,
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
		earlierEthTx := models.EthTx{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 0},
			Value:          value,
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 1),
			State:          models.EthTxUnstarted,
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
			require.Equal(t, "0x6acac565c14ca984f1fad43e63036ccc777bfa95a3447930696bb3a33dd99653", tx.Hash().Hex())

			// They must be set to something to indicate that the transaction is signed
			v, r, s := tx.RawSignatureValues()
			require.Equal(t, "41", v.String())
			require.Equal(t, "100125404117036954913117369048685056327836806830110180962458866833256017916154", r.String())
			require.Equal(t, "13748121423502857499887005034545879393991949882113605371398299275605825415634", s.String())
			return true
		})).Return(nil).Once()

		// Later
		laterEthTx := models.EthTx{
			FromAddress:    defaultFromAddress,
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
			require.Equal(t, "0xbea12954ffafe9ac9d89abf6f2f9a563cc50fa14a31e5866e0c0063e76556b60", tx.Hash().Hex())

			// They must be set to something to indicate that the transaction is signed
			v, r, s := tx.RawSignatureValues()
			require.Equal(t, "42", v.String())
			require.Equal(t, "39363214223465398755579021511352119350941428292598621771180906281886401763946", r.String())
			require.Equal(t, "13319764116922590262026344403688458878732413946946573942114704806995618455150", s.String())
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
		assert.Equal(t, defaultFromAddress, earlierTransaction.FromAddress)
		require.NotNil(t, earlierTransaction.Nonce)
		assert.Equal(t, int64(0), *earlierTransaction.Nonce)
		assert.NotNil(t, earlierTransaction.BroadcastAt)
		assert.Len(t, earlierTransaction.EthTxAttempts, 1)

		attempt := earlierTransaction.EthTxAttempts[0]

		assert.Equal(t, earlierTransaction.ID, attempt.EthTxID)
		assert.Equal(t, config.EthGasPriceDefault().String(), attempt.GasPrice.String())

		signedTx, err := attempt.GetSignedTx()
		require.NoError(t, err)
		assert.Equal(t, "0x6acac565c14ca984f1fad43e63036ccc777bfa95a3447930696bb3a33dd99653", signedTx.Hash().Hex())
		assert.Equal(t, "0x6acac565c14ca984f1fad43e63036ccc777bfa95a3447930696bb3a33dd99653", attempt.Hash.Hex())
		assert.Equal(t, "0xf867808504a817c80081f2946c03dda95a2aed917eecc6eddd4b9d16e6380411818e832a2a0029a0dd5cf86fe8e6c6c863c5cc4feb2cbfa5a87b289d8f74b8d82a599931629970faa01e65293571cd92fb96398dfd22362e76cacb527ff9472c5aa14439ae3381e9d2", hexutil.Encode(attempt.SignedRawTx))
		assert.Equal(t, models.EthTxAttemptBroadcast, attempt.State)
		require.Len(t, attempt.EthReceipts, 0)

		// Check laterEthTx and it's attempt
		// This was the later one sent so it has the higher nonce
		laterTransaction, err := store.FindEthTxWithAttempts(laterEthTx.ID)
		require.NoError(t, err)
		assert.Nil(t, laterTransaction.Error)
		require.NotNil(t, laterTransaction.FromAddress)
		assert.Equal(t, defaultFromAddress, laterTransaction.FromAddress)
		require.NotNil(t, laterTransaction.Nonce)
		assert.Equal(t, int64(1), *laterTransaction.Nonce)
		assert.NotNil(t, laterTransaction.BroadcastAt)
		assert.Len(t, laterTransaction.EthTxAttempts, 1)

		attempt = laterTransaction.EthTxAttempts[0]

		assert.Equal(t, laterTransaction.ID, attempt.EthTxID)
		assert.Equal(t, config.EthGasPriceDefault().String(), attempt.GasPrice.String())

		signedTx, err = attempt.GetSignedTx()
		require.NoError(t, err)
		assert.Equal(t, "0xbea12954ffafe9ac9d89abf6f2f9a563cc50fa14a31e5866e0c0063e76556b60", signedTx.Hash().Hex())
		assert.Equal(t, "0xbea12954ffafe9ac9d89abf6f2f9a563cc50fa14a31e5866e0c0063e76556b60", attempt.Hash.Hex())
		assert.Equal(t, "0xf867018504a817c80081f2946c03dda95a2aed917eecc6eddd4b9d16e6380411818e832a2a012aa05706ca2b15c5796218fc602be65cca821d28310135407889fa40bf409c891a6aa01d72b825e1c765c8a3368cbef7ce3c249ceceadc36aa17c60294c4c959545e6e", hexutil.Encode(attempt.SignedRawTx))
		assert.Equal(t, models.EthTxAttemptBroadcast, attempt.State)
		require.Len(t, attempt.EthReceipts, 0)

		ethClient.AssertExpectations(t)
	})
}

func TestEthBroadcaster_AssignsNonceOnFirstRun(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	// Simulate new key by manually setting nonce to null
	require.NoError(t, store.DB.Exec(`UPDATE keys SET next_nonce = NULL`).Error)

	// Use the real KeyStore loaded from database fixtures
	store.KeyStore.Unlock(cltest.Password)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	eb := bulletprooftxmanager.NewEthBroadcaster(store, config)

	keys, err := store.Keys()
	require.NoError(t, err)
	key := keys[0]
	defaultFromAddress := key.Address.Address()
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	gasLimit := uint64(242)

	ethTx := models.EthTx{
		FromAddress:    defaultFromAddress,
		ToAddress:      toAddress,
		EncodedPayload: []byte{42, 42, 0},
		Value:          assets.NewEthValue(0),
		GasLimit:       gasLimit,
		CreatedAt:      time.Unix(0, 0),
		State:          models.EthTxUnstarted,
	}
	require.NoError(t, store.DB.Create(&ethTx).Error)

	t.Run("when eth node returns error", func(t *testing.T) {
		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(account gethCommon.Address) bool {
			return account.Hex() == defaultFromAddress.Hex()
		})).Return(uint64(0), errors.New("something exploded")).Once()

		// First attempt errored
		err = eb.ProcessUnstartedEthTxs(key)
		require.Error(t, err)
		require.Contains(t, err.Error(), "something exploded")

		// Check ethTx that it has no nonce assigned
		ethTx, err = store.FindEthTxWithAttempts(ethTx.ID)
		require.NoError(t, err)

		require.Nil(t, ethTx.Nonce)

		// Check key to make sure it still doesn't have a nonce assigned
		res := store.DB.Exec(`SELECT * FROM keys WHERE next_nonce IS NULL`)
		require.NoError(t, res.Error)
		require.Equal(t, int64(1), res.RowsAffected)

		ethClient.AssertExpectations(t)
	})

	t.Run("when eth node returns nonce", func(t *testing.T) {
		ethNodeNonce := uint64(42)

		ethClient.On("PendingNonceAt", mock.Anything, mock.MatchedBy(func(account gethCommon.Address) bool {
			return account.Hex() == defaultFromAddress.Hex()
		})).Return(ethNodeNonce, nil).Once()
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == ethNodeNonce
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(key))

		// Check ethTx that it has the correct nonce assigned
		ethTx, err = store.FindEthTxWithAttempts(ethTx.ID)
		require.NoError(t, err)

		require.NotNil(t, ethTx.Nonce)
		require.Equal(t, int64(ethNodeNonce), *ethTx.Nonce)

		// Check key to make sure it has correct nonce assigned
		keys, err := store.Keys()
		require.NoError(t, err)
		key := keys[0]

		require.NotNil(t, key.NextNonce)
		require.Equal(t, int64(43), *key.NextNonce)

		ethClient.AssertExpectations(t)
	})
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_ResumingFromCrash(t *testing.T) {
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)
	encodedPayload := []byte{0, 1}
	nextNonce := int64(916714082576372851)

	t.Run("cannot be more than one transaction per address in an unfinished state", func(t *testing.T) {
		firstNonce := nextNonce + 1
		secondNonce := nextNonce + 2

		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		keys, err := store.Keys()
		require.NoError(t, err)
		key := keys[0]
		defaultFromAddress := key.Address.Address()

		firstInProgress := models.EthTx{
			FromAddress:    defaultFromAddress,
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
			FromAddress:    defaultFromAddress,
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
		err = store.DB.Create(&secondInProgress).Error
		require.Error(t, err)
		assert.EqualError(t, err, "pq: duplicate key value violates unique constraint \"idx_only_one_in_progress_tx_per_account\"")
	})

	t.Run("previous run assigned nonce but never broadcast", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		// Use the real KeyStore loaded from database fixtures
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		defer cleanup()

		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		eb := bulletprooftxmanager.NewEthBroadcaster(store, config)

		keys, err := store.Keys()
		require.NoError(t, err)
		key := keys[0]
		defaultFromAddress := key.Address.Address()

		require.NoError(t, store.DB.Exec(`UPDATE keys SET next_nonce = ? WHERE address = ?`, nextNonce, defaultFromAddress.Bytes()).Error)

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		nonce := nextNonce
		inProgressEthTx := mustInsertInProgressEthTxWithAttempt(t, store, nextNonce)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(nonce)
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
		// Use the real KeyStore loaded from database fixtures
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		defer cleanup()

		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		eb := bulletprooftxmanager.NewEthBroadcaster(store, config)

		keys, err := store.Keys()
		require.NoError(t, err)
		key := keys[0]
		defaultFromAddress := key.Address.Address()

		require.NoError(t, store.DB.Exec(`UPDATE keys SET next_nonce = ? WHERE address = ?`, nextNonce, defaultFromAddress.Bytes()).Error)

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		nonce := nextNonce
		inProgressEthTx := mustInsertInProgressEthTxWithAttempt(t, store, nextNonce)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(nonce)
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
		// Use the real KeyStore loaded from database fixtures
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		defer cleanup()

		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		eb := bulletprooftxmanager.NewEthBroadcaster(store, config)

		keys, err := store.Keys()
		require.NoError(t, err)
		key := keys[0]
		defaultFromAddress := key.Address.Address()

		require.NoError(t, store.DB.Exec(`UPDATE keys SET next_nonce = ? WHERE address = ?`, nextNonce, defaultFromAddress.Bytes()).Error)

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		nonce := nextNonce
		inProgressEthTx := mustInsertInProgressEthTxWithAttempt(t, store, nextNonce)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(nonce)
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
		// Use the real KeyStore loaded from database fixtures
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		defer cleanup()

		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		eb := bulletprooftxmanager.NewEthBroadcaster(store, config)

		keys, err := store.Keys()
		require.NoError(t, err)
		key := keys[0]
		defaultFromAddress := key.Address.Address()

		require.NoError(t, store.DB.Exec(`UPDATE keys SET next_nonce = ? WHERE address = ?`, nextNonce, defaultFromAddress.Bytes()).Error)

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		nonce := nextNonce
		inProgressEthTx := mustInsertInProgressEthTxWithAttempt(t, store, nextNonce)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(nonce)
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
		// Use the real KeyStore loaded from database fixtures
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		defer cleanup()

		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		eb := bulletprooftxmanager.NewEthBroadcaster(store, config)

		keys, err := store.Keys()
		require.NoError(t, err)
		key := keys[0]
		defaultFromAddress := key.Address.Address()

		require.NoError(t, store.DB.Exec(`UPDATE keys SET next_nonce = ? WHERE address = ?`, nextNonce, defaultFromAddress.Bytes()).Error)

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		nonce := nextNonce
		inProgressEthTx := mustInsertInProgressEthTxWithAttempt(t, store, nextNonce)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(nonce)
		})).Return(failedToReachNodeError).Once()

		// Do the thing
		err = eb.ProcessUnstartedEthTxs(key)
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
		// Use the real KeyStore loaded from database fixtures
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		defer cleanup()

		ethClient := new(mocks.Client)
		store.EthClient = ethClient

		// Configured gas price changed
		store.Config.Set("ETH_GAS_PRICE_DEFAULT", 500000000000)
		eb := bulletprooftxmanager.NewEthBroadcaster(store, config)

		keys, err := store.Keys()
		require.NoError(t, err)
		key := keys[0]
		defaultFromAddress := key.Address.Address()

		require.NoError(t, store.DB.Exec(`UPDATE keys SET next_nonce = ? WHERE address = ?`, nextNonce, defaultFromAddress.Bytes()).Error)

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_tx so keys.next_nonce has not been
		// incremented yet
		nonce := nextNonce
		inProgressEthTx := mustInsertInProgressEthTxWithAttempt(t, store, nextNonce)
		require.Len(t, inProgressEthTx.EthTxAttempts, 1)
		attempt := inProgressEthTx.EthTxAttempts[0]

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			// Ensure that the gas price is the same as the original attempt
			s, e := attempt.GetSignedTx()
			require.NoError(t, e)
			return tx.Nonce() == uint64(nonce) && tx.GasPrice().Int64() == s.GasPrice().Int64()
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
	return uint64(*n)
}

// Note that all of these tests share the same database, and ordering matters.
// This in order to more deeply test ProcessUnstartedEthTxs over
// multiple runs with previous errors in the database.
func TestEthBroadcaster_ProcessUnstartedEthTxs_Errors(t *testing.T) {
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)
	encodedPayload := []byte{0, 1}

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	// Use the real KeyStore loaded from database fixtures
	store.KeyStore.Unlock(cltest.Password)
	keys, err := store.Keys()
	require.NoError(t, err)
	key := keys[0]
	defaultFromAddress := key.Address.Address()

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	eb := bulletprooftxmanager.NewEthBroadcaster(store, config)

	t.Run("external wallet sent a transaction from the account and now the nonce is one higher than it should be", func(t *testing.T) {
		localNextNonce := getLocalNextNonce(t, store, defaultFromAddress)
		require.Equal(t, 0, int(localNextNonce))
		remoteNextNonce := uint64(1)

		etx := models.EthTx{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          models.EthTxUnstarted,
		}
		require.NoError(t, store.DB.Save(&etx).Error)
		taskRunID := cltest.MustInsertTaskRun(t, store)
		ethTaskRunTx := models.EthTaskRunTx{
			EthTxID:   etx.ID,
			TaskRunID: taskRunID.UUID(),
		}
		require.NoError(t, store.DB.Save(&ethTaskRunTx).Error)

		// First send, nonce too low
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New("nonce too low")).Once()

		// Second send with higher nonce
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == remoteNextNonce
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnstartedEthTxs(key))

		// Check that two transactions were sent
		ethClient.AssertExpectations(t)

		// Check the 'nonce too low' transaction was saved correctly with its attempt
		var etx1 models.EthTx
		etx1, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)
		require.NotNil(t, etx1.BroadcastAt)
		require.NotNil(t, etx1.Nonce)
		assert.Equal(t, int64(localNextNonce), *etx1.Nonce)
		assert.Nil(t, etx1.Error)
		assert.Len(t, etx1.EthTxAttempts, 1)

		// Check that the second transaction was saved correctly with its attempt
		var latestID int64
		var etx2 models.EthTx
		require.NoError(t, store.DB.Raw("SELECT max(id) FROM eth_txes").Row().Scan(&latestID))
		etx2, err = store.FindEthTxWithAttempts(latestID)
		require.NoError(t, err)
		require.NotNil(t, etx2.BroadcastAt)
		assert.NotEqual(t, etx2.CreatedAt, *etx2.BroadcastAt)
		require.NotNil(t, etx2.Nonce)
		assert.Equal(t, int64(localNextNonce+1), *etx2.Nonce)
		assert.Nil(t, etx2.Error)
		assert.Len(t, etx2.EthTxAttempts, 1)

		// Check that the second transaction is later than the first but otherwise identical
		assert.Greater(t, (*etx2.BroadcastAt).UnixNano(), (*etx1.BroadcastAt).UnixNano())
		assert.Equal(t, etx1.FromAddress, etx2.FromAddress)
		assert.Equal(t, etx1.ToAddress, etx2.ToAddress)
		assert.Equal(t, etx1.EncodedPayload, etx2.EncodedPayload)
		assert.Equal(t, etx1.Value, etx2.Value)
		assert.Equal(t, etx1.GasLimit, etx2.GasLimit)
		assert.Greater(t, etx2.CreatedAt.UnixNano(), etx1.CreatedAt.UnixNano())

		// Check that the eth_task_run_tx was re-assigned
		res := store.DB.Exec(`SELECT * FROM eth_task_run_txes WHERE eth_tx_id = ?`, etx.ID)
		require.NoError(t, res.Error)
		require.Equal(t, int64(0), res.RowsAffected)
		res = store.DB.Exec(`SELECT * FROM eth_task_run_txes WHERE eth_tx_id = ?`, etx2.ID)
		require.NoError(t, res.Error)
		require.Equal(t, int64(1), res.RowsAffected)

		// Check that the local nonce was incremented by two
		var finalNextNonce *int64
		finalNextNonce, err = bulletprooftxmanager.GetNextNonce(store.DB, defaultFromAddress)
		require.NoError(t, err)
		require.NotNil(t, finalNextNonce)
		require.Equal(t, int64(2), *finalNextNonce)
	})

	t.Run("geth client returns an error in the fatal errors category", func(t *testing.T) {
		fatalErrorExample := "exceeds block gas limit"
		localNextNonce := getLocalNextNonce(t, store, defaultFromAddress)

		etx := models.EthTx{
			FromAddress:    defaultFromAddress,
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
		require.Equal(t, int64(localNextNonce), *key.NextNonce)

		ethClient.AssertExpectations(t)
	})

	t.Run("eth client call fails with an unexpected random error (e.g. insufficient funds)", func(t *testing.T) {
		retryableErrorExample := "insufficient funds for transfer"
		localNextNonce := getLocalNextNonce(t, store, defaultFromAddress)

		etx := models.EthTx{
			FromAddress:    defaultFromAddress,
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
		require.Contains(t, err.Error(), fmt.Sprintf("error while sending transaction %v: insufficient funds for transfer", etx.ID))

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

		etx := models.EthTx{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          models.EthTxUnstarted,
		}
		require.NoError(t, store.DB.Save(&etx).Error)

		// First was underpriced
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.GasPrice().Cmp(store.Config.EthGasPriceDefault()) == 0
		})).Return(errors.New(underpricedError)).Once()

		// Second with gas bump was still underpriced
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.GasPrice().Cmp(big.NewInt(25000000000)) == 0
		})).Return(errors.New(underpricedError)).Once()

		// Third succeeded
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.GasPrice().Cmp(big.NewInt(30000000000)) == 0
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

	t.Run("failed to reach node for some reason", func(t *testing.T) {
		failedToReachNodeError := context.DeadlineExceeded
		localNextNonce := getLocalNextNonce(t, store, defaultFromAddress)

		etx := models.EthTx{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          models.EthTxUnstarted,
		}
		require.NoError(t, store.DB.Save(&etx).Error)

		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(failedToReachNodeError).Once()

		// Do the thing
		err = eb.ProcessUnstartedEthTxs(key)
		require.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("error while sending transaction %v: context deadline exceeded", etx.ID))

		// Check it was left in the unfinished state
		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.Nonce)
		assert.Nil(t, etx.Error)
		assert.Equal(t, models.EthTxInProgress, etx.State)
		assert.Len(t, etx.EthTxAttempts, 1)
		assert.Equal(t, models.EthTxAttemptInProgress, etx.EthTxAttempts[0].State)

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
	// Use a mock keystore for this test
	store.KeyStore = kst
	keys, err := store.Keys()
	require.NoError(t, err)
	key := keys[0]
	defaultFromAddress := key.Address.Address()

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	ethClient := new(mocks.Client)
	store.EthClient = ethClient

	eb := bulletprooftxmanager.NewEthBroadcaster(store, config)

	t.Run("keystore does not have the unlocked key", func(t *testing.T) {
		etx := models.EthTx{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          models.EthTxUnstarted,
		}
		require.NoError(t, store.DB.Save(&etx).Error)

		kst.On("GetAccountByAddress", defaultFromAddress).Return(gethAccounts.Account{}, errors.New("authentication needed: password or unlock")).Once()

		// Do the thing
		err := eb.ProcessUnstartedEthTxs(key)
		require.Error(t, err)
		require.Contains(t, err.Error(), "authentication needed: password or unlock")

		// Check that the transaction is left in unstarted state
		etx, err = store.FindEthTxWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Equal(t, models.EthTxUnstarted, etx.State)
		assert.Len(t, etx.EthTxAttempts, 0)

		// Check that the key did not have its nonce incremented
		require.NoError(t, store.DB.First(&key).Error)
		require.NotNil(t, key.NextNonce)
		require.Equal(t, int64(localNonce), *key.NextNonce)

		kst.AssertExpectations(t)
	})

	t.Run("tx signing fails", func(t *testing.T) {
		etx := models.EthTx{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			State:          models.EthTxUnstarted,
		}
		require.NoError(t, store.DB.Save(&etx).Error)

		signingAccount := gethAccounts.Account{Address: defaultFromAddress}
		kst.On("GetAccountByAddress", defaultFromAddress).Return(signingAccount, nil).Once()

		tx := gethTypes.Transaction{}
		kst.On("SignTx",
			mock.AnythingOfType("accounts.Account"),
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
		require.Equal(t, int64(localNonce), *key.NextNonce)

		kst.AssertExpectations(t)
	})

	// Should have done nothing
	ethClient.AssertExpectations(t)
}

func TestEthBroadcaster_ProcessUnstartedEthTxs_Locking(t *testing.T) {
	store1, cleanup := cltest.NewStore(t)
	defer cleanup()
	// Use the real KeyStore loaded from database fixtures
	store1.KeyStore.Unlock(cltest.Password)

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	ethClient := new(mocks.Client)
	store1.EthClient = ethClient

	eb1 := bulletprooftxmanager.NewEthBroadcaster(store1, config)

	// Simulate another node
	store2, cleanup := cltest.NewStore(t)
	defer cleanup()
	eb2 := bulletprooftxmanager.NewEthBroadcaster(store2, config)

	keys, err := store1.Keys()
	require.NoError(t, err)
	key := keys[0]
	defaultFromAddress := key.Address.Address()
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)

	chSendingTx := make(chan struct{})
	chMidway := make(chan struct{})
	chFinish := make(chan struct{})

	etx := models.EthTx{
		FromAddress:    defaultFromAddress,
		ToAddress:      toAddress,
		EncodedPayload: []byte{42, 42, 0},
		Value:          value,
		GasLimit:       gasLimit,
		CreatedAt:      time.Unix(0, 0),
		State:          models.EthTxUnstarted,
	}
	ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
		close(chSendingTx)
		<-chMidway
		return true
	})).Return(nil).Once()

	require.NoError(t, store1.DB.Save(&etx).Error)

	// First one gets the lock
	go func() {
		err2 := eb1.ProcessUnstartedEthTxs(key)
		assert.NoError(t, err2)
		close(chFinish)
	}()

	g := gomega.NewGomegaWithT(t)

	// Wait until first one is in the middle of its run
	g.Eventually(chSendingTx).Should(gomega.BeClosed())

	// Second node's attempt to get lock fails
	err = eb2.ProcessUnstartedEthTxs(key)
	require.Error(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("could not get advisory lock for classID, objectID %v, %v", 0, key.ID))

	// Resume original run
	close(chMidway)

	// Ensure all go routines exited
	g.Eventually(chFinish).Should(gomega.BeClosed())
}

func TestEthBroadcaster_GetNextNonce(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// Fixture key has nonce 0
	var key models.Key
	require.NoError(t, store.DB.First(&key).Error)
	require.NotNil(t, key.NextNonce)
	require.Equal(t, int64(0), *key.NextNonce)

	nonce, err := bulletprooftxmanager.GetNextNonce(store.DB, key.Address.Address())
	assert.NoError(t, err)
	require.NotNil(t, nonce)
	assert.Equal(t, int64(0), *nonce)
}

func TestEthBroadcaster_IncrementNextNonce(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// Fixture key had nonce 0
	var key models.Key
	require.NoError(t, store.DB.First(&key).Error)
	require.NotNil(t, key.NextNonce)
	require.Equal(t, int64(0), *key.NextNonce)

	previouslyUpdatedAt := key.UpdatedAt

	// Cannot increment if supplied nonce doesn't match existing
	require.Error(t, bulletprooftxmanager.IncrementNextNonce(store.DB, key.Address.Address(), int64(42)))

	require.NoError(t, bulletprooftxmanager.IncrementNextNonce(store.DB, key.Address.Address(), int64(0)))

	// Nonce bumped to 1
	require.NoError(t, store.DB.First(&key).Error)
	require.NotNil(t, key.NextNonce)
	require.Equal(t, int64(1), *key.NextNonce)
	// Updated at
	require.Greater(t, key.UpdatedAt.Unix(), previouslyUpdatedAt.Unix())
}

func TestEthBroadcaster_Trigger(t *testing.T) {
	// Simple sanity check to make sure it doesn't block
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	eb := bulletprooftxmanager.NewEthBroadcaster(store, config)

	eb.Trigger()
	eb.Trigger()
	eb.Trigger()
}
