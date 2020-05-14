package bulletprooftxmanager_test

import (
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	gethAccounts "github.com/ethereum/go-ethereum/accounts"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
)

func TestTxBroadcaster_NewBulletproofTxManager(t *testing.T) {
	// TODO: write this test
}

func TestBulletproofTxManager_ProcessUnbroadcastEthTransactions_Locking(t *testing.T) {
	store1, cleanup := cltest.NewStore(t)
	defer cleanup()
	// Use the real KeyStore loaded from database fixtures
	store1.KeyStore.Unlock(cltest.Password)

	config, cleanup := cltest.NewConfig(t)
	gethClient := new(mocks.GethClient)
	gethWrapper := cltest.NewSimpleGethWrapper(gethClient)
	eb1 := bulletprooftxmanager.NewEthBroadcaster(store1, gethWrapper, config)

	keys, err := store1.Keys()
	require.NoError(t, err)
	key := keys[0]
	defaultFromAddress := key.Address.Address()
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)

	chSendingTx := make(chan bool)
	chFinish := make(chan struct{})

	etx := models.EthTransaction{
		FromAddress:    defaultFromAddress,
		ToAddress:      toAddress,
		EncodedPayload: []byte{42, 42, 0},
		Value:          value,
		GasLimit:       gasLimit,
		CreatedAt:      time.Unix(0, 0),
	}
	gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
		chSendingTx <- true
		<-chFinish
		return true
	})).Return(nil).Once()

	require.NoError(t, store1.GetRawDB().Save(&etx).Error)

	// First one gets the lock
	go func() {
		require.NoError(t, eb1.ProcessUnbroadcastEthTransactions(*key))
	}()

	// Wait until first one is in the middle of its run
	<-chSendingTx

	// Simulate another node
	store2, cleanup := cltest.NewStore(t)
	defer cleanup()
	eb2 := bulletprooftxmanager.NewEthBroadcaster(store2, gethWrapper, config)

	// Second attempt to get lock fails
	require.EqualError(t, eb2.ProcessUnbroadcastEthTransactions(*key), fmt.Sprintf("could not get advisory lock for key %v", key.ID))

	// Resume original run
	close(chFinish)
}

func TestBulletproofTxManager_ProcessUnbroadcastEthTransactions_Success(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	// Use the real KeyStore loaded from database fixtures
	store.KeyStore.Unlock(cltest.Password)

	config, cleanup := cltest.NewConfig(t)
	gethClient := new(mocks.GethClient)
	gethWrapper := cltest.NewSimpleGethWrapper(gethClient)
	eb := bulletprooftxmanager.NewEthBroadcaster(store, gethWrapper, config)

	keys, err := store.Keys()
	require.NoError(t, err)
	key := keys[0]
	defaultFromAddress := key.Address.Address()
	toAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
	timeNow := time.Now()

	encodedPayload := []byte{1, 2, 3}
	value := assets.NewEthValue(142)
	gasLimit := uint64(242)

	t.Run("no eth_transactions at all", func(t *testing.T) {
		require.NoError(t, eb.ProcessUnbroadcastEthTransactions(*key))
	})

	t.Run("eth_transactions exist for a different from address", func(t *testing.T) {
		otherAddress := gethCommon.HexToAddress("0x6C03DDA95a2AEd917EeCc6eddD4b9D16E6380411")
		cltest.MustInsertKey(t, store, otherAddress)

		etx := models.EthTransaction{
			FromAddress:    otherAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
		}
		require.NoError(t, store.GetRawDB().Save(&etx).Error)

		require.NoError(t, eb.ProcessUnbroadcastEthTransactions(*key))
	})

	t.Run("existing eth_transactions with broadcast_at or error", func(t *testing.T) {
		nonce := int64(342)
		errStr := "some error"

		etxWithNonce := models.EthTransaction{
			Nonce:          &nonce,
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			BroadcastAt:    &timeNow,
			Error:          nil,
		}
		etxWithError := models.EthTransaction{
			Nonce:          nil,
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			Error:          &errStr,
		}

		require.NoError(t, store.GetRawDB().Save(&etxWithNonce).Error)
		require.NoError(t, store.GetRawDB().Save(&etxWithError).Error)

		require.NoError(t, eb.ProcessUnbroadcastEthTransactions(*key))
	})

	t.Run("sends two EthTransactions in order starting from the earliest", func(t *testing.T) {
		// Earlier
		earlierEthTransaction := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 0},
			Value:          value,
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(0, 0),
		}
		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			if tx.Nonce() != uint64(0) {
				return false
			}
			require.Equal(t, config.ChainID(), tx.ChainId())
			require.Equal(t, gasLimit, tx.Gas())
			require.Equal(t, config.EthGasPriceDefault(), tx.GasPrice())
			require.Equal(t, toAddress, *tx.To())
			require.Equal(t, value.ToInt().String(), tx.Value().String())
			require.Equal(t, earlierEthTransaction.EncodedPayload, tx.Data())

			// They must be set to something to indicate that the transaction is signed
			v, r, s := tx.RawSignatureValues()
			require.Equal(t, "41", v.String())
			require.Equal(t, "100125404117036954913117369048685056327836806830110180962458866833256017916154", r.String())
			require.Equal(t, "13748121423502857499887005034545879393991949882113605371398299275605825415634", s.String())
			return true
		})).Return(nil).Once()

		// Later
		laterEthTransaction := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: []byte{42, 42, 1},
			Value:          value,
			GasLimit:       gasLimit,
			CreatedAt:      time.Unix(1, 0),
		}
		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			if tx.Nonce() != uint64(1) {
				return false
			}
			require.Equal(t, config.ChainID(), tx.ChainId())
			require.Equal(t, gasLimit, tx.Gas())
			require.Equal(t, config.EthGasPriceDefault(), tx.GasPrice())
			require.Equal(t, toAddress, *tx.To())
			require.Equal(t, value.ToInt().String(), tx.Value().String())
			require.Equal(t, laterEthTransaction.EncodedPayload, tx.Data())

			// They must be set to something to indicate that the transaction is signed
			v, r, s := tx.RawSignatureValues()
			require.Equal(t, "42", v.String())
			require.Equal(t, "39363214223465398755579021511352119350941428292598621771180906281886401763946", r.String())
			require.Equal(t, "13319764116922590262026344403688458878732413946946573942114704806995618455150", s.String())
			return true
		})).Return(nil).Once()

		// Insertion order deliberately reversed to test order by created at
		require.NoError(t, store.GetRawDB().Save(&laterEthTransaction).Error)
		require.NoError(t, store.GetRawDB().Save(&earlierEthTransaction).Error)

		// Do the thing
		require.NoError(t, eb.ProcessUnbroadcastEthTransactions(*key))

		// Check earlierEthTransaction and it's attempt
		// This was the earlier one sent so it has the lower nonce
		earlierTransaction, err := store.FindEthTransactionWithAttempts(earlierEthTransaction.ID)
		require.NoError(t, err)
		assert.Nil(t, earlierTransaction.Error)
		require.NotNil(t, earlierTransaction.FromAddress)
		assert.Equal(t, defaultFromAddress, earlierTransaction.FromAddress)
		require.NotNil(t, earlierTransaction.Nonce)
		assert.Equal(t, int64(0), *earlierTransaction.Nonce)
		assert.NotNil(t, earlierTransaction.BroadcastAt)
		assert.Len(t, earlierTransaction.EthTransactionAttempts, 1)

		attempt := earlierTransaction.EthTransactionAttempts[0]

		assert.Equal(t, earlierTransaction.ID, attempt.EthTransactionID)
		assert.Equal(t, config.EthGasPriceDefault().String(), attempt.GasPrice.String())
		assert.Nil(t, attempt.Hash)
		assert.Nil(t, attempt.Error)
		assert.Nil(t, attempt.ConfirmedInBlockNum)
		assert.Nil(t, attempt.ConfirmedInBlockHash)
		assert.Nil(t, attempt.ConfirmedAt)

		assert.Equal(t, "0xf867808504a817c80081f2946c03dda95a2aed917eecc6eddd4b9d16e6380411818e832a2a0029a0dd5cf86fe8e6c6c863c5cc4feb2cbfa5a87b289d8f74b8d82a599931629970faa01e65293571cd92fb96398dfd22362e76cacb527ff9472c5aa14439ae3381e9d2", hexutil.Encode(attempt.SignedRawTx))

		// Check laterEthTransaction and it's attempt
		// This was the later one sent so it has the higher nonce
		laterTransaction, err := store.FindEthTransactionWithAttempts(laterEthTransaction.ID)
		require.NoError(t, err)
		assert.Nil(t, laterTransaction.Error)
		require.NotNil(t, laterTransaction.FromAddress)
		assert.Equal(t, defaultFromAddress, laterTransaction.FromAddress)
		require.NotNil(t, laterTransaction.Nonce)
		assert.Equal(t, int64(1), *laterTransaction.Nonce)
		assert.NotNil(t, laterTransaction.BroadcastAt)
		assert.Len(t, laterTransaction.EthTransactionAttempts, 1)

		attempt = laterTransaction.EthTransactionAttempts[0]

		assert.Equal(t, laterTransaction.ID, attempt.EthTransactionID)
		assert.Equal(t, config.EthGasPriceDefault().String(), attempt.GasPrice.String())
		assert.Nil(t, attempt.Hash)
		assert.Nil(t, attempt.Error)
		assert.Nil(t, attempt.ConfirmedInBlockNum)
		assert.Nil(t, attempt.ConfirmedInBlockHash)
		assert.Nil(t, attempt.ConfirmedAt)

		assert.Equal(t, "0xf867018504a817c80081f2946c03dda95a2aed917eecc6eddd4b9d16e6380411818e832a2a012aa05706ca2b15c5796218fc602be65cca821d28310135407889fa40bf409c891a6aa01d72b825e1c765c8a3368cbef7ce3c249ceceadc36aa17c60294c4c959545e6e", hexutil.Encode(attempt.SignedRawTx))

		gethClient.AssertExpectations(t)
	})
}

func TestBulletproofTxManager_ProcessUnbroadcastEthTransactions_ResumingFromCrash(t *testing.T) {
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

		firstUnfinished := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			Nonce:          &firstNonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			BroadcastAt:    nil,
			Error:          nil,
		}

		secondUnfinished := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			Nonce:          &secondNonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			BroadcastAt:    nil,
			Error:          nil,
		}

		require.NoError(t, store.GetRawDB().Create(&firstUnfinished).Error)
		err = store.GetRawDB().Create(&secondUnfinished).Error
		require.Error(t, err)
		assert.EqualError(t, err, "pq: duplicate key value violates unique constraint \"idx_only_one_in_progress_tx_per_account\"")
	})

	t.Run("previous run assigned nonce but never broadcast", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		// Use the real KeyStore loaded from database fixtures
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		gethClient := new(mocks.GethClient)
		gethWrapper := cltest.NewSimpleGethWrapper(gethClient)
		eb := bulletprooftxmanager.NewEthBroadcaster(store, gethWrapper, config)

		keys, err := store.Keys()
		require.NoError(t, err)
		key := keys[0]
		defaultFromAddress := key.Address.Address()

		require.NoError(t, store.GetRawDB().Exec(`UPDATE keys SET next_nonce = ? WHERE address = ?`, nextNonce, defaultFromAddress.Bytes()).Error)

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_transaction so keys.next_nonce has not been
		// incremented yet
		nonce := nextNonce
		unbroadcastEthTransactionWithNonce := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			Nonce:          &nonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			BroadcastAt:    nil,
		}

		require.NoError(t, store.GetRawDB().Create(&unbroadcastEthTransactionWithNonce).Error)

		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(nonce)
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnbroadcastEthTransactions(*key))

		// Check it was saved correctly with its attempt
		etx, err := store.FindEthTransactionWithAttempts(unbroadcastEthTransactionWithNonce.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTransactionAttempts, 1)

		gethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast but it unretryably errored before we could save", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		// Use the real KeyStore loaded from database fixtures
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		gethClient := new(mocks.GethClient)
		gethWrapper := cltest.NewSimpleGethWrapper(gethClient)
		eb := bulletprooftxmanager.NewEthBroadcaster(store, gethWrapper, config)

		keys, err := store.Keys()
		require.NoError(t, err)
		key := keys[0]
		defaultFromAddress := key.Address.Address()

		require.NoError(t, store.GetRawDB().Exec(`UPDATE keys SET next_nonce = ? WHERE address = ?`, nextNonce, defaultFromAddress.Bytes()).Error)

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_transaction so keys.next_nonce has not been
		// incremented yet
		nonce := nextNonce
		unbroadcastEthTransactionWithNonce := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			Nonce:          &nonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			BroadcastAt:    nil,
		}

		require.NoError(t, store.GetRawDB().Create(&unbroadcastEthTransactionWithNonce).Error)

		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(nonce)
		})).Return(errors.New("exceeds block gas limit")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnbroadcastEthTransactions(*key))

		// Check it was saved correctly with its attempt
		etx, err := store.FindEthTransactionWithAttempts(unbroadcastEthTransactionWithNonce.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.Error)
		assert.Equal(t, "exceeds block gas limit", *etx.Error)
		assert.Len(t, etx.EthTransactionAttempts, 0)

		gethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast and is now in mempool", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		// Use the real KeyStore loaded from database fixtures
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		gethClient := new(mocks.GethClient)
		gethWrapper := cltest.NewSimpleGethWrapper(gethClient)
		eb := bulletprooftxmanager.NewEthBroadcaster(store, gethWrapper, config)

		keys, err := store.Keys()
		require.NoError(t, err)
		key := keys[0]
		defaultFromAddress := key.Address.Address()

		require.NoError(t, store.GetRawDB().Exec(`UPDATE keys SET next_nonce = ? WHERE address = ?`, nextNonce, defaultFromAddress.Bytes()).Error)

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_transaction so keys.next_nonce has not been
		// incremented yet
		nonce := nextNonce
		unbroadcastEthTransactionWithNonce := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			Nonce:          &nonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			BroadcastAt:    nil,
		}

		require.NoError(t, store.GetRawDB().Create(&unbroadcastEthTransactionWithNonce).Error)

		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(nonce)
		})).Return(errors.New("known transaction: a1313bd99a81fb4d8ad1d2e90b67c6b3fa77545c990d6251444b83b70b6f8980")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnbroadcastEthTransactions(*key))

		// Check it was saved correctly with its attempt
		etx, err := store.FindEthTransactionWithAttempts(unbroadcastEthTransactionWithNonce.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTransactionAttempts, 1)

		gethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and broadcast and now the transaction has been confirmed", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		// Use the real KeyStore loaded from database fixtures
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		gethClient := new(mocks.GethClient)
		gethWrapper := cltest.NewSimpleGethWrapper(gethClient)
		eb := bulletprooftxmanager.NewEthBroadcaster(store, gethWrapper, config)

		keys, err := store.Keys()
		require.NoError(t, err)
		key := keys[0]
		defaultFromAddress := key.Address.Address()

		require.NoError(t, store.GetRawDB().Exec(`UPDATE keys SET next_nonce = ? WHERE address = ?`, nextNonce, defaultFromAddress.Bytes()).Error)

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_transaction so keys.next_nonce has not been
		// incremented yet
		nonce := nextNonce
		unbroadcastEthTransactionWithNonce := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			Nonce:          &nonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			BroadcastAt:    nil,
		}

		require.NoError(t, store.GetRawDB().Create(&unbroadcastEthTransactionWithNonce).Error)

		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(nonce)
		})).Return(errors.New("nonce too low")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnbroadcastEthTransactions(*key))

		// Check it was saved correctly with its attempt
		etx, err := store.FindEthTransactionWithAttempts(unbroadcastEthTransactionWithNonce.ID)
		require.NoError(t, err)

		require.NotNil(t, etx.BroadcastAt)
		assert.Equal(t, *etx.BroadcastAt, etx.CreatedAt)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTransactionAttempts, 1)

		gethClient.AssertExpectations(t)
	})

	t.Run("previous run assigned nonce and then failed to reach node for some reason and node is still down", func(t *testing.T) {
		failedToReachNodeError := context.DeadlineExceeded
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		// Use the real KeyStore loaded from database fixtures
		store.KeyStore.Unlock(cltest.Password)

		config, cleanup := cltest.NewConfig(t)
		gethClient := new(mocks.GethClient)
		gethWrapper := cltest.NewSimpleGethWrapper(gethClient)
		eb := bulletprooftxmanager.NewEthBroadcaster(store, gethWrapper, config)

		keys, err := store.Keys()
		require.NoError(t, err)
		key := keys[0]
		defaultFromAddress := key.Address.Address()

		require.NoError(t, store.GetRawDB().Exec(`UPDATE keys SET next_nonce = ? WHERE address = ?`, nextNonce, defaultFromAddress.Bytes()).Error)

		// Crashed right after we commit the database transaction that saved
		// the nonce to the eth_transaction so keys.next_nonce has not been
		// incremented yet
		nonce := nextNonce
		unbroadcastEthTransactionWithNonce := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			Nonce:          &nonce,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
			BroadcastAt:    nil,
		}

		require.NoError(t, store.GetRawDB().Create(&unbroadcastEthTransactionWithNonce).Error)

		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == uint64(nonce)
		})).Return(failedToReachNodeError).Once()

		// Do the thing
		err = eb.ProcessUnbroadcastEthTransactions(*key)
		require.Error(t, err)
		assert.EqualError(t, failedToReachNodeError, err.Error())

		// Check it was left in the unfinished state
		etx, err := store.FindEthTransactionWithAttempts(unbroadcastEthTransactionWithNonce.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.Equal(t, nextNonce, *etx.Nonce)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTransactionAttempts, 0)

		gethClient.AssertExpectations(t)
	})
}

func getLocalNextNonce(t *testing.T, str *store.Store, fromAddress gethCommon.Address) uint64 {
	n, err := bulletprooftxmanager.GetNextNonce(str.GetRawDB(), fromAddress)
	require.NoError(t, err)
	return uint64(n)
}

// Note that all of these tests share the same database, and ordering matters.
// This in order to more deeply test ProcessUnbroadcastEthTransactions over
// multiple runs with previous errors in the database.
func TestBulletproofTxManager_ProcessUnbroadcastEthTransactions_Errors(t *testing.T) {
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
	gethClient := new(mocks.GethClient)
	gethWrapper := cltest.NewSimpleGethWrapper(gethClient)
	eb := bulletprooftxmanager.NewEthBroadcaster(store, gethWrapper, config)

	t.Run("external wallet sent a transction from the account and now the nonce is one higher than it should be", func(t *testing.T) {
		// TODO: Describe this behaviour
		localNextNonce := getLocalNextNonce(t, store, defaultFromAddress)
		require.Equal(t, uint64(0), localNextNonce)
		remoteNextNonce := uint64(1)

		etx := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
		}
		require.NoError(t, store.GetRawDB().Save(&etx).Error)

		// First send, nonce too low
		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New("nonce too low")).Once()

		// Second send with higher nonce
		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == remoteNextNonce
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnbroadcastEthTransactions(*key))

		// Check that two transactions were sent
		gethClient.AssertExpectations(t)

		// Check the 'nonce too low' transaction was saved correctly with its attempt
		etx1, err := store.FindEthTransactionWithAttempts(etx.ID)
		require.NoError(t, err)
		require.NotNil(t, etx1.BroadcastAt)
		assert.Equal(t, etx1.CreatedAt, *etx1.BroadcastAt)
		require.NotNil(t, etx1.Nonce)
		assert.Equal(t, int64(localNextNonce), *etx1.Nonce)
		assert.Nil(t, etx1.Error)
		assert.Len(t, etx1.EthTransactionAttempts, 1)

		// Check that the second transaction was saved correctly with its attempt
		var latestID int64
		require.NoError(t, store.GetRawDB().Raw("SELECT max(id) FROM eth_transactions").Row().Scan(&latestID))
		etx2, err := store.FindEthTransactionWithAttempts(latestID)
		require.NoError(t, err)
		require.NotNil(t, etx2.BroadcastAt)
		assert.NotEqual(t, etx2.CreatedAt, *etx2.BroadcastAt)
		require.NotNil(t, etx2.Nonce)
		assert.Equal(t, int64(localNextNonce+1), *etx2.Nonce)
		assert.Nil(t, etx2.Error)
		assert.Len(t, etx2.EthTransactionAttempts, 1)

		// Check that the second transaction is later than the first but otherwise identical
		assert.Greater(t, (*etx2.BroadcastAt).UnixNano(), (*etx1.BroadcastAt).UnixNano())
		assert.Equal(t, etx1.FromAddress, etx2.FromAddress)
		assert.Equal(t, etx1.ToAddress, etx2.ToAddress)
		assert.Equal(t, etx1.EncodedPayload, etx2.EncodedPayload)
		assert.Equal(t, etx1.Value, etx2.Value)
		assert.Equal(t, etx1.GasLimit, etx2.GasLimit)
		assert.Greater(t, etx2.CreatedAt.UnixNano(), etx1.CreatedAt.UnixNano())

		// Check that the local nonce was incremented by two
		finalNextNonce, err := bulletprooftxmanager.GetNextNonce(store.GetRawDB(), defaultFromAddress)
		require.NoError(t, err)
		require.Equal(t, int64(2), finalNextNonce)
	})

	t.Run("geth client returns an error in the unretryable errors category", func(t *testing.T) {
		unretryableErrorExample := "exceeds block gas limit"
		localNextNonce := getLocalNextNonce(t, store, defaultFromAddress)

		etx := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
		}
		require.NoError(t, store.GetRawDB().Save(&etx).Error)

		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(unretryableErrorExample)).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnbroadcastEthTransactions(*key))

		// Check it was saved correctly with its attempt
		etx, err = store.FindEthTransactionWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		require.Nil(t, etx.Nonce)
		assert.NotNil(t, etx.Error)
		assert.Contains(t, *etx.Error, "exceeds block gas limit")
		assert.Len(t, etx.EthTransactionAttempts, 0)

		// Check that the key had its nonce reset
		var key models.Key
		require.NoError(t, store.GetRawDB().First(&key).Error)
		// Saved NextNonce must be the same as before because this transaction
		// was not accepted by the eth node and never can be
		require.Equal(t, int64(localNextNonce), key.NextNonce)

		gethClient.AssertExpectations(t)
	})

	t.Run("gethclient call fails with an unexpected random error", func(t *testing.T) {
		retryableErrorExample := "some old bollocks"
		localNextNonce := getLocalNextNonce(t, store, defaultFromAddress)

		etx := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
		}
		require.NoError(t, store.GetRawDB().Save(&etx).Error)

		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(errors.New(retryableErrorExample)).Once()

		// Do the thing
		require.EqualError(t, eb.ProcessUnbroadcastEthTransactions(*key), "some old bollocks")

		// Check it was saved correctly with its attempt
		etx, err = store.FindEthTransactionWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTransactionAttempts, 0)

		gethClient.AssertExpectations(t)

		// Now on the second run, it is successful
		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(nil).Once()

		require.NoError(t, eb.ProcessUnbroadcastEthTransactions(*key))

		// Check it was saved correctly with its attempt
		etx, err = store.FindEthTransactionWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTransactionAttempts, 1)

		gethClient.AssertExpectations(t)
	})

	t.Run("eth node returns 'underpriced transaction'", func(t *testing.T) {
		// This happens if a transaction's gas price is below the minimum
		// configured for the transaction pool.
		// This is a configuration error by the node operator, since it means they set the base gas level too low.
		// We should enter a gas bumping loop until this error no longer occurs
		underpricedError := "transaction underpriced"

		etx := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
		}
		require.NoError(t, store.GetRawDB().Save(&etx).Error)

		// First was underpriced
		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.GasPrice().Cmp(store.Config.EthGasPriceDefault()) == 0
		})).Return(errors.New(underpricedError)).Once()

		// Second with gas bump was still underpriced
		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.GasPrice().Cmp(big.NewInt(25000000000)) == 0
		})).Return(errors.New(underpricedError)).Once()

		// Third succeeded
		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.GasPrice().Cmp(big.NewInt(30000000000)) == 0
		})).Return(nil).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnbroadcastEthTransactions(*key))

		gethClient.AssertExpectations(t)

		// Check it was saved correctly with its attempt
		etx, err = store.FindEthTransactionWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.NotNil(t, etx.BroadcastAt)
		require.NotNil(t, etx.Nonce)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTransactionAttempts, 1)
		attempt := etx.EthTransactionAttempts[0]
		assert.Equal(t, big.NewInt(30000000000).String(), attempt.GasPrice.String())
	})

	t.Run("failed to reach node for some reason", func(t *testing.T) {
		failedToReachNodeError := context.DeadlineExceeded
		localNextNonce := getLocalNextNonce(t, store, defaultFromAddress)

		etx := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
		}
		require.NoError(t, store.GetRawDB().Save(&etx).Error)

		gethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == localNextNonce
		})).Return(failedToReachNodeError).Once()

		// Do the thing
		err = eb.ProcessUnbroadcastEthTransactions(*key)
		require.Error(t, err)
		assert.Equal(t, failedToReachNodeError, err)

		// Check it was left in the unfinished state
		etx, err = store.FindEthTransactionWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		assert.NotNil(t, etx.Nonce)
		assert.Nil(t, etx.Error)
		assert.Len(t, etx.EthTransactionAttempts, 0)

		gethClient.AssertExpectations(t)
	})

	// TODO:
	// - "replacement transaction underpriced"
	// - 'nonce too high'
}

func TestBulletproofTxManager_ProcessUnbroadcastEthTransactions_KeystoreErrors(t *testing.T) {
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
	gethClient := new(mocks.GethClient)
	gethWrapper := cltest.NewSimpleGethWrapper(gethClient)
	eb := bulletprooftxmanager.NewEthBroadcaster(store, gethWrapper, config)

	t.Run("keystore does not have the unlocked key", func(t *testing.T) {
		etx := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
		}
		require.NoError(t, store.GetRawDB().Save(&etx).Error)

		kst.On("GetAccountByAddress", defaultFromAddress).Return(gethAccounts.Account{}, errors.New("authentication needed: password or unlock")).Once()

		// Do the thing
		require.NoError(t, eb.ProcessUnbroadcastEthTransactions(*key))

		// Check that the transaction shows the error
		// Check it was saved correctly with its attempt
		etx, err = store.FindEthTransactionWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		require.Nil(t, etx.Nonce)
		assert.NotNil(t, etx.Error)
		assert.Contains(t, *etx.Error, "authentication needed: password or unlock")
		assert.Len(t, etx.EthTransactionAttempts, 0)

		// Check that the key had its nonce reset
		var key models.Key
		require.NoError(t, store.GetRawDB().First(&key).Error)
		require.Equal(t, int64(localNonce), key.NextNonce)

		kst.AssertExpectations(t)
	})

	t.Run("tx signing fails", func(t *testing.T) {
		etx := models.EthTransaction{
			FromAddress:    defaultFromAddress,
			ToAddress:      toAddress,
			EncodedPayload: encodedPayload,
			Value:          value,
			GasLimit:       gasLimit,
		}
		require.NoError(t, store.GetRawDB().Save(&etx).Error)

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
		require.NoError(t, eb.ProcessUnbroadcastEthTransactions(*key))

		// Check that the transaction shows the error
		// Check it was saved correctly with its attempt
		etx, err = store.FindEthTransactionWithAttempts(etx.ID)
		require.NoError(t, err)

		assert.Nil(t, etx.BroadcastAt)
		require.Nil(t, etx.Nonce)
		assert.NotNil(t, etx.Error)
		assert.Contains(t, *etx.Error, "could not sign transaction")
		assert.Len(t, etx.EthTransactionAttempts, 0)

		// Check that the key had its nonce reset
		var key models.Key
		require.NoError(t, store.GetRawDB().First(&key).Error)
		require.Equal(t, int64(localNonce), key.NextNonce)

		kst.AssertExpectations(t)
	})

	// Should have done nothing
	gethClient.AssertExpectations(t)
}

func TestBulletproofTxManager_ProcessUnbroadcastEthTransactions_ComplexTest(t *testing.T) {
	// Multiple tx's, some of which fail with different errors on multiple occasions over multiple calls
}

func TestBulletproofTxManager_GetDefaultAddress(t *testing.T) {
	// Test cases:
	// -
}

func TestBulletproofTxManager_GetNextNonce(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// Fixture key has nonce 0
	var key models.Key
	require.NoError(t, store.GetRawDB().First(&key).Error)
	require.Equal(t, int64(0), key.NextNonce)

	nonce, err := bulletprooftxmanager.GetNextNonce(store.GetRawDB(), key.Address.Address())
	assert.NoError(t, err)
	assert.Equal(t, int64(0), nonce)
}

func TestBulletproofTxManager_IncrementNextNonce(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// Fixture key had nonce 0
	var key models.Key
	require.NoError(t, store.GetRawDB().First(&key).Error)
	require.Equal(t, int64(0), key.NextNonce)

	// Cannot increment if supplied nonce doesn't match existing
	require.Error(t, bulletprooftxmanager.IncrementNextNonce(store.GetRawDB(), key.Address.Address(), int64(42)))

	require.NoError(t, bulletprooftxmanager.IncrementNextNonce(store.GetRawDB(), key.Address.Address(), int64(0)))

	// Nonce bumped to 1
	require.NoError(t, store.GetRawDB().First(&key).Error)
	require.Equal(t, int64(1), key.NextNonce)
}

// TODO: Probably better track down the 'The account is being used by another wallet' warning and make sure it also checks the local key nonce
