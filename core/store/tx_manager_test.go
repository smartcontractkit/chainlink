package store_test

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

func TestTxManager_CreateTx_Success(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ethClient := new(mocks.Client)

	config := cltest.NewTestConfig(t)
	keyStore := strpkg.NewKeyStore(config.KeysDir())
	account, err := keyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	require.NoError(t, keyStore.Unlock(cltest.Password))
	manager := strpkg.NewEthTxManager(ethClient, config, keyStore, store.ORM)

	from := account.Address
	to := cltest.NewAddress()
	data := hexutil.MustDecode("0x0000abcdef")
	nonce := uint64(256)

	manager.Register(keyStore.Accounts())

	ethClient.On("GetNonce", from).Return(nonce, nil)

	err = manager.Connect(cltest.Head(nonce))
	require.NoError(t, err)

	ethClient.On("SendRawTx", mock.Anything).Return(cltest.NewHash(), nil)

	tx, err := manager.CreateTx(to, data)
	require.NoError(t, err)

	ntx, err := store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Equal(t, nonce, ntx.Nonce)
	assert.Equal(t, data, ntx.Data)
	assert.Equal(t, to, ntx.To)
	assert.Equal(t, from, ntx.From)
	assert.Equal(t, nonce, ntx.Nonce)
	assert.Len(t, ntx.Attempts, 1)

	ethClient.AssertExpectations(t)
}

func TestTxManager_CreateTx_RoundRobinSuccess(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ethClient := new(mocks.Client)

	config := cltest.NewTestConfig(t)
	keyStore := strpkg.NewKeyStore(config.KeysDir())

	// Add two accounts
	_, err := keyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	_, err = keyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	require.NoError(t, keyStore.Unlock(cltest.Password))

	manager := strpkg.NewEthTxManager(ethClient, config, keyStore, store.ORM)

	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	nonce := uint64(256)
	sentAt := uint64(1)
	bumpAt := sentAt + config.EthGasBumpThreshold()

	manager.Register(keyStore.Accounts())
	require.NoError(t, err)

	ethClient.On("GetNonce", mock.Anything).Return(nonce, nil).Times(2)

	err = manager.Connect(cltest.Head(sentAt))
	require.NoError(t, err)

	ethClient.On("SendRawTx", mock.Anything).Return(cltest.NewHash(), nil)

	createdTx1, err := manager.CreateTx(to, data)
	require.NoError(t, err)

	ntx, err := store.FindTx(createdTx1.ID)
	require.NoError(t, err)
	assert.Len(t, ntx.Attempts, 1)

	manager.OnNewLongestChain(*cltest.Head(bumpAt))
	ethClient.On("GetTxReceipt", createdTx1.Attempts[0].Hash).Return(&models.TxReceipt{}, nil)
	ethClient.On("SendRawTx", mock.Anything).Return(cltest.NewHash(), nil)

	// bump gas
	receipt, state, err := manager.BumpGasUntilSafe(createdTx1.Attempts[0].Hash)
	require.NoError(t, err)
	assert.Nil(t, receipt)
	assert.Equal(t, strpkg.Unconfirmed, state)

	createdTx1, err = store.FindTx(createdTx1.ID)
	require.NoError(t, err)
	require.Len(t, createdTx1.Attempts, 2)

	// Ensure that the From address hasn't been updated on the Tx
	assert.Equal(t, createdTx1.From, ntx.From)

	// ensure second tx uses the first account again
	ethClient.On("SendRawTx", mock.Anything).Return(cltest.NewHash(), nil)

	createdTx2, err := manager.CreateTx(to, data)
	require.NoError(t, err)
	require.NotEqual(t, createdTx1.From.Hex(), createdTx2.From.Hex(), "should come from a different account")

	ethClient.AssertExpectations(t)
}

func TestTxManager_CreateTx_BreakTxAttemptLimit(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ethClient := new(mocks.Client)

	config := cltest.NewTestConfig(t)
	config.Set("CHAINLINK_TX_ATTEMPT_LIMIT", 1)
	keyStore := strpkg.NewKeyStore(config.KeysDir())
	account, err := keyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	require.NoError(t, keyStore.Unlock(cltest.Password))
	manager := strpkg.NewEthTxManager(ethClient, config, keyStore, store.ORM)

	from := account.Address
	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	nonce := uint64(256)
	sentAt := uint64(1)
	bumpAt := sentAt + config.EthGasBumpThreshold()
	bumpAgainAt := sentAt + 2*config.EthGasBumpThreshold()

	manager.Register(keyStore.Accounts())
	require.NoError(t, err)

	ethClient.On("GetNonce", from).Return(nonce, nil)

	err = manager.Connect(cltest.Head(sentAt))
	require.NoError(t, err)

	ethClient.On("SendRawTx", mock.Anything).Return(cltest.NewHash(), nil)

	tx, err := manager.CreateTx(to, data)
	require.NoError(t, err)

	ntx, err := store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Equal(t, nonce, ntx.Nonce)
	assert.Equal(t, data, ntx.Data)
	assert.Equal(t, to, ntx.To)
	assert.Equal(t, from, ntx.From)
	assert.Equal(t, nonce, ntx.Nonce)
	assert.Len(t, ntx.Attempts, 1)

	manager.OnNewLongestChain(*cltest.Head(bumpAt))
	ethClient.On("GetTxReceipt", mock.Anything).Once().Return(&models.TxReceipt{}, nil)
	ethClient.On("SendRawTx", mock.Anything).Return(tx.Attempts[0].Hash, nil)

	receipt, state, err := manager.BumpGasUntilSafe(tx.Attempts[0].Hash)
	require.NoError(t, err)
	assert.Nil(t, receipt)
	assert.Equal(t, strpkg.Unconfirmed, state)

	manager.OnNewLongestChain(*cltest.Head(bumpAgainAt))
	ethClient.On("GetTxReceipt", mock.Anything).Twice().Return(&models.TxReceipt{}, nil)

	receipt, state, err = manager.BumpGasUntilSafe(tx.Attempts[0].Hash)
	require.NoError(t, err)
	assert.Nil(t, receipt)
	assert.Equal(t, strpkg.Unconfirmed, state)

	ethClient.AssertExpectations(t)
}

func TestTxManager_CreateTx_AttemptErrorDoesNotIncrementNonce(t *testing.T) {
	t.Parallel()

	config, configCleanup := cltest.NewConfig(t)
	defer configCleanup()

	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config)
	defer cleanup()

	store := app.Store
	manager := store.TxManager

	from := cltest.GetAccountAddress(t, app.GetStore())
	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	assert.NoError(t, err)
	sentAt := uint64(23456)
	nonce := uint64(256)
	ethMock := app.EthMock
	ethMock.Context("app.StartAndConnect()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce))
		ethMock.Register("eth_chainId", store.Config.ChainID())
	})
	require.NoError(t, app.Store.ORM.IdempotentInsertHead(*cltest.Head(sentAt)))
	assert.NoError(t, app.StartAndConnect())

	ethMock.Context("manager.CreateTx#1", func(ethMock *cltest.EthMock) {
		ethMock.RegisterError("eth_sendRawTransaction", "invalid transaction")
	})

	_, err = manager.CreateTx(to, data)
	assert.Error(t, err)

	txs, _, err := store.Transactions(0, 10)
	assert.NoError(t, err)
	assert.Len(t, txs, 1)

	hash := cltest.NewHash()
	ethMock.Context("manager.CreateTx#2", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_sendRawTransaction", hash)
	})

	tx, err := manager.CreateTx(to, data)
	require.NoError(t, err)

	ntx, err := store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Equal(t, nonce, ntx.Nonce)
	assert.Equal(t, data, ntx.Data)
	assert.Equal(t, to, ntx.To)
	assert.Equal(t, from, ntx.From)
	assert.Equal(t, nonce, ntx.Nonce)
	assert.Len(t, ntx.Attempts, 1)

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_CreateTx_NonceTooLowReloadSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		ethClientErrorMsg string
	}{
		{"geth", "nonce too low"},
		{"geth", "replacement transaction underpriced"},
		{"parity", "Transaction nonce is too low. Try incrementing the nonce"},
		{"parity", "Transaction with the same hash was already imported"},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			ethClient := new(mocks.Client)
			config := cltest.NewTestConfig(t)
			require.NoError(t, store.KeyStore.Unlock(cltest.Password))
			manager := strpkg.NewEthTxManager(ethClient, config, store.KeyStore, store.ORM)
			manager.Register(store.KeyStore.Accounts())

			from := cltest.GetAccountAddress(t, store)
			to := cltest.NewAddress()
			data := hexutil.MustDecode("0x0000abcdef")

			nonce := uint64(256)
			ethClient.On("GetNonce", from).Once().Return(nonce, nil)
			require.NoError(t, manager.Connect(cltest.Head(nonce)))

			ethClient.On("SendRawTx", mock.Anything).Once().Return(nil, errors.New("nonce is too low"))
			nonce2 := uint64(257)
			ethClient.On("GetNonce", from).Once().Return(nonce2, nil)
			ethClient.On("SendRawTx", mock.Anything).Once().Return(nil, nil)

			a, err := manager.CreateTx(to, data)
			require.NoError(t, err)
			tx, err := store.FindTx(a.ID)
			require.NoError(t, err)
			assert.Equal(t, nonce2, tx.Nonce)
			assert.Equal(t, data, tx.Data)
			assert.Equal(t, to, tx.To)
			assert.Equal(t, from, tx.From)
			assert.Len(t, tx.Attempts, 1)

			ethClient.AssertExpectations(t)
		})
	}
}

func TestTxManager_CreateTx_NonceTooLowReloadLimit(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ethClient := new(mocks.Client)

	config := cltest.NewTestConfig(t)
	keyStore := strpkg.NewKeyStore(config.KeysDir())
	account, err := keyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	require.NoError(t, keyStore.Unlock(cltest.Password))
	manager := strpkg.NewEthTxManager(ethClient, config, keyStore, store.ORM)

	manager.Register(keyStore.Accounts())

	from := account.Address
	nonce := uint64(256)
	ethClient.On("GetNonce", from).Return(nonce, nil)
	err = manager.Connect(cltest.Head(nonce))
	require.NoError(t, err)

	ethClient.On("SendRawTx", mock.Anything).Times(4).Return(nil, errors.New("nonce is too low"))

	to := cltest.NewAddress()
	data := hexutil.MustDecode("0x0000abcdef")
	_, err = manager.CreateTx(to, data)
	assert.EqualError(t, err, "transaction reattempt limit reached for 'nonce is too low' error. Limit: 3")

	ethClient.AssertExpectations(t)
}

func TestTxManager_CreateTx_ErrPendingConnection(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	manager := store.TxManager

	to := cltest.NewAddress()
	data := hexutil.MustDecode("0x0000abcdef")

	_, err := manager.CreateTx(to, data)
	assert.Contains(t, err.Error(), strpkg.ErrPendingConnection.Error())
}

func TestTxManager_BumpGasUntilSafe_lessThanGasBumpThreshold(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	config := store.Config

	txm := store.TxManager
	from := cltest.GetAccountAddress(t, store)
	sentAt := uint64(23456)
	gasThreshold := sentAt + config.EthGasBumpThreshold()
	ethMock := app.EthMock
	ethMock.Register("eth_getTransactionCount", "0x0")
	ethMock.Register("eth_chainId", config.ChainID())
	require.NoError(t, app.Store.ORM.IdempotentInsertHead(*cltest.Head(gasThreshold - 1)))
	require.NoError(t, app.StartAndConnect())

	tx := cltest.CreateTx(t, store, from, sentAt)
	require.Greater(t, len(tx.Attempts), 0)

	ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{})

	receipt, state, err := txm.BumpGasUntilSafe(tx.Attempts[0].Hash)
	assert.NoError(t, err)
	assert.Nil(t, receipt)
	assert.Equal(t, strpkg.Unconfirmed, state)

	tx, err = store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Len(t, tx.Attempts, 1)

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_BumpGasUntilSafe_atGasBumpThreshold(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	config := store.Config

	txm := store.TxManager
	from := cltest.GetAccountAddress(t, store)
	sentAt := uint64(23456)
	gasThreshold := sentAt + config.EthGasBumpThreshold()
	ethMock := app.EthMock
	ethMock.Register("eth_getTransactionCount", "0x0")
	ethMock.Register("eth_chainId", config.ChainID())
	require.NoError(t, app.Store.ORM.IdempotentInsertHead(*cltest.Head(gasThreshold)))
	require.NoError(t, app.StartAndConnect())

	tx := cltest.CreateTx(t, store, from, sentAt)
	require.Greater(t, len(tx.Attempts), 0)

	ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{})
	ethMock.Register("eth_sendRawTransaction", cltest.NewHash())

	receipt, state, err := txm.BumpGasUntilSafe(tx.Attempts[0].Hash)
	assert.NoError(t, err)
	assert.Nil(t, receipt)
	assert.Equal(t, strpkg.Unconfirmed, state)

	tx, err = store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Len(t, tx.Attempts, 2)

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_BumpGasUntilSafe_atGasBumpThreshold_bumpsGasMoreInCaseOfUnderpricedTransaction(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	config := store.Config
	config.Set("ETH_GAS_BUMP_PERCENT", 10)

	txm := store.TxManager
	from := cltest.GetAccountAddress(t, store)
	sentAt := uint64(23456)
	gasThreshold := sentAt + config.EthGasBumpThreshold()
	ethMock := app.EthMock
	ethMock.Register("eth_getTransactionCount", "0x0")
	ethMock.Register("eth_chainId", config.ChainID())
	require.NoError(t, app.Store.ORM.IdempotentInsertHead(*cltest.Head(gasThreshold)))
	require.NoError(t, app.StartAndConnect())

	tx := cltest.CreateTxWithNonceAndGasPrice(t, store, from, sentAt, 0, 48000000000)
	require.Greater(t, len(tx.Attempts), 0)

	ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{})

	// Simulate two bumps that receive `replacement transaction underpriced`
	// and the third and final one successful
	ethMock.RegisterError("eth_sendRawTransaction", "replacement transaction underpriced")
	ethMock.RegisterError("eth_sendRawTransaction", "replacement transaction underpriced")
	ethMock.Register("eth_sendRawTransaction", cltest.NewHash())

	receipt, state, err := txm.BumpGasUntilSafe(tx.Attempts[0].Hash)
	assert.NoError(t, err)
	assert.Nil(t, receipt)
	assert.Equal(t, strpkg.Unconfirmed, state)

	tx, err = store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Len(t, tx.Attempts, 2)

	latestAttempt := tx.Attempts[1]
	gasPrice, err := latestAttempt.GasPrice.Value()
	require.NoError(t, err)

	// Initial gas price is 48Gwei
	// Bump to 53Gwei and fail
	// Bump to 58.3Gwei and fail
	// Bump to 64.13Gwei and pass
	assert.Equal(t, "64130000000", gasPrice)

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_BumpGasUntilSafe_atGasBumpThreshold_CapsAtMaxIfMaxGasPriceIsReached(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	config := store.Config

	txm := store.TxManager
	from := cltest.GetAccountAddress(t, store)
	sentAt := uint64(23456)
	gasThreshold := sentAt + config.EthGasBumpThreshold()
	ethMock := app.EthMock
	ethMock.Register("eth_getTransactionCount", "0x0")
	ethMock.Register("eth_chainId", config.ChainID())
	ethMock.Register("eth_sendRawTransaction", cltest.NewHash())
	require.NoError(t, app.Store.ORM.IdempotentInsertHead(*cltest.Head(gasThreshold)))
	require.NoError(t, app.StartAndConnect())

	tx := cltest.CreateTxWithNonceAndGasPrice(t, store, from, sentAt, 0, 499000000000)
	store.SaveTx(tx)
	require.Greater(t, len(tx.Attempts), 0)

	ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{})

	receipt, state, err := txm.BumpGasUntilSafe(tx.Attempts[0].Hash)
	require.NoError(t, err)
	assert.Nil(t, receipt)
	assert.Equal(t, strpkg.Unconfirmed, state)

	tx, err = store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Len(t, tx.Attempts, 2)

	latestAttempt := tx.Attempts[1]
	gasPrice, err := latestAttempt.GasPrice.Value()
	require.NoError(t, err)
	assert.Equal(t, "500000000000", gasPrice)

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_BumpGasUntilSafe_exceedsGasBumpThreshold(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	config := store.Config

	txm := store.TxManager
	from := cltest.GetAccountAddress(t, store)
	sentAt := uint64(23456)
	gasThreshold := sentAt + config.EthGasBumpThreshold()
	ethMock := app.EthMock
	ethMock.Register("eth_getTransactionCount", "0x0")
	ethMock.Register("eth_chainId", config.ChainID())
	require.NoError(t, app.Store.ORM.IdempotentInsertHead(*cltest.Head(gasThreshold + 1)))
	require.NoError(t, app.StartAndConnect())

	tx := cltest.CreateTx(t, store, from, sentAt)
	require.Greater(t, len(tx.Attempts), 0)

	ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{})
	ethMock.Register("eth_sendRawTransaction", cltest.NewHash())

	receipt, state, err := txm.BumpGasUntilSafe(tx.Attempts[0].Hash)
	assert.NoError(t, err)
	assert.Nil(t, receipt)
	assert.Equal(t, strpkg.Unconfirmed, state)

	tx, err = store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Len(t, tx.Attempts, 2)

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_BumpGasUntilSafe_confirmed(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	app.EthMock.Context("app.Start()", func(meth *cltest.EthMock) {
		meth.Register("eth_getTransactionCount", "0x1")
		meth.Register("eth_chainId", app.Store.Config.ChainID())
	})
	store := app.Store
	config := store.Config

	sentAt := uint64(23456)
	nonce := uint64(234)
	gasThreshold := sentAt + config.EthGasBumpThreshold()
	minConfs := config.MinRequiredOutgoingConfirmations() - 1
	require.NoError(t, app.Store.ORM.IdempotentInsertHead(*cltest.Head(gasThreshold + minConfs - 1)))
	require.NoError(t, app.StartAndConnect())

	txm := store.TxManager
	from := cltest.GetAccountAddress(t, store)

	tx := cltest.CreateTxWithNonceAndGasPrice(t, store, from, sentAt, nonce, 1)
	require.Greater(t, len(tx.Attempts), 0)

	app.EthMock.Register("eth_getTransactionReceipt", models.TxReceipt{Hash: cltest.NewHash(), BlockNumber: cltest.Int(gasThreshold)})

	receipt, state, err := txm.BumpGasUntilSafe(tx.Attempts[0].Hash)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, strpkg.Confirmed, state)

	tx, err = store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Len(t, tx.Attempts, 1)

	etm := txm.(*strpkg.EthTxManager)
	aa := etm.GetAvailableAccount(from)
	assert.NotEqual(t, tx.Nonce, aa.PublicLastSafeNonce())

	app.EthMock.EventuallyAllCalled(t)
}

func TestTxManager_BumpGasUntilSafe_safe(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		confsDiff uint64
	}{
		{"at threshold", 0},
		{"above threshold", 1},
	}

	sentAt := uint64(23456)
	nonce := uint64(234)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, cleanup := cltest.NewApplicationWithKey(t)
			defer cleanup()
			app.EthMock.Context("app.Start()", func(meth *cltest.EthMock) {
				meth.Register("eth_getTransactionCount", "0x1")
				meth.Register("eth_chainId", app.Store.Config.ChainID())
			})
			store := app.Store
			config := store.Config

			gasThreshold := sentAt + config.EthGasBumpThreshold()
			minConfs := config.MinRequiredOutgoingConfirmations() - 1
			head := cltest.Head(gasThreshold + minConfs + test.confsDiff)
			require.NoError(t, app.Store.ORM.IdempotentInsertHead(*head))
			require.NoError(t, app.StartAndConnect())

			txm := store.TxManager
			from := cltest.GetAccountAddress(t, store)

			tx := cltest.CreateTxWithNonceAndGasPrice(t, store, from, sentAt, nonce, 1)
			require.Greater(t, len(tx.Attempts), 0)

			app.EthMock.Register("eth_getTransactionReceipt", models.TxReceipt{Hash: cltest.NewHash(), BlockNumber: cltest.Int(gasThreshold)})

			receipt, state, err := txm.BumpGasUntilSafe(tx.Attempts[0].Hash)
			assert.NoError(t, err)
			assert.NotNil(t, receipt)
			assert.Equal(t, strpkg.Safe, state)

			tx, err = store.FindTx(tx.ID)
			require.NoError(t, err)
			assert.Len(t, tx.Attempts, 1)

			etm := txm.(*strpkg.EthTxManager)
			aa := etm.GetAvailableAccount(from)
			assert.Equal(t, tx.Nonce, aa.PublicLastSafeNonce())

			app.EthMock.EventuallyAllCalled(t)
		})
	}
}

func TestTxManager_BumpGasUntilSafe_laterConfirmedTx(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	app.EthMock.Context("app.Start()", func(meth *cltest.EthMock) {
		meth.Register("eth_getTransactionCount", "0x1")
		meth.Register("eth_chainId", app.Store.Config.ChainID())
	})

	require.NoError(t, app.StartAndConnect())

	store := app.Store
	txm := store.TxManager
	from := cltest.GetAccountAddress(t, store)
	sentAt := uint64(12345)

	tx1 := cltest.CreateTxWithNonceAndGasPrice(t, store, from, sentAt, 1, 1)
	require.Len(t, tx1.Attempts, 1)
	tx2 := cltest.CreateTxWithNonceAndGasPrice(t, store, from, sentAt, 2, 1)
	require.Len(t, tx2.Attempts, 1)

	etm := txm.(*strpkg.EthTxManager)
	aa := etm.GetAvailableAccount(from)
	aa.SetLastSafeNonce(tx2.Nonce)

	app.EthMock.Register("eth_getTransactionReceipt", models.TxReceipt{})

	receipt, state, err := txm.BumpGasUntilSafe(tx1.Attempts[0].Hash)
	assert.Nil(t, receipt)
	assert.Equal(t, strpkg.Safe, state)
	assert.Error(t, err)

	tx, err := store.FindTx(tx1.ID)
	require.NoError(t, err)
	assert.True(t, tx.Confirmed)
	assert.Equal(t, tx.Hash.Hex(), "0x0000000000000000000000000000000000000000000000000000000000000000")
	assert.Len(t, tx.Attempts, 1)

	app.EthMock.EventuallyAllCalled(t)
}

func TestTxManager_BumpGasUntilSafe_erroring(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	sentAt1 := uint64(23456)
	sentAt2 := sentAt1 + config.EthGasBumpThreshold()
	confirmedAt := sentAt2 + 1
	safeAt := confirmedAt + config.MinRequiredOutgoingConfirmations()

	nonConfedReceipt := models.TxReceipt{}
	confedReceipt := models.TxReceipt{Hash: cltest.NewHash(), BlockNumber: cltest.Int(confirmedAt)}

	tests := []struct {
		name        string
		blockHeight uint64
		mockSetup   func(*cltest.EthMock)
		wantReceipt bool
		wantErrored bool
	}{
		{"no conf, no error", (sentAt2 + 1), func(ethMock *cltest.EthMock) {
			ethMock.Register("eth_getTransactionCount", "0x0")
			ethMock.Register("eth_getTransactionReceipt", nonConfedReceipt)
			ethMock.Register("eth_getTransactionReceipt", nonConfedReceipt)
		}, false, false},
		{"no conf, early error", (sentAt2 + 1), func(ethMock *cltest.EthMock) {
			ethMock.Register("eth_getTransactionCount", "0x0")
			ethMock.RegisterError("eth_getTransactionReceipt", "FUBAR")
			ethMock.Register("eth_getTransactionReceipt", nonConfedReceipt)
		}, false, true},
		{"no conf, later error", (sentAt2 + 1), func(ethMock *cltest.EthMock) {
			ethMock.Register("eth_getTransactionCount", "0x0")
			ethMock.Register("eth_getTransactionReceipt", nonConfedReceipt)
			ethMock.RegisterError("eth_getTransactionReceipt", "FUBAR")
		}, false, true},
		{"early conf", (safeAt + 1), func(ethMock *cltest.EthMock) {
			ethMock.Register("eth_getTransactionCount", "0x0")
			ethMock.Register("eth_getTransactionReceipt", confedReceipt)
		}, true, false},
		{"later conf, no error", (safeAt + 1), func(ethMock *cltest.EthMock) {
			ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(0))
			ethMock.Register("eth_getTransactionReceipt", confedReceipt)
		}, true, false},
		{"later conf, early error", (safeAt + 1), func(ethMock *cltest.EthMock) {
			ethMock.Register("eth_getTransactionCount", "0x0")
			ethMock.RegisterError("eth_getTransactionReceipt", "FUBAR")
			ethMock.Register("eth_getTransactionReceipt", confedReceipt)
		}, true, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config)
			defer cleanup()

			store := app.Store
			txm := store.TxManager
			from := cltest.GetAccountAddress(t, store)

			tx := cltest.CreateTx(t, store, from, sentAt1)
			a := cltest.AddTxAttempt(t, store, tx, tx.EthTx(big.NewInt(2)), sentAt2)

			ethMock := app.EthMock
			ethMock.ShouldCall(test.mockSetup).During(func() {
				require.NoError(t, app.Store.ORM.IdempotentInsertHead(*cltest.Head(test.blockHeight)))
				ethMock.Register("eth_chainId", store.Config.ChainID())
				ethMock.Register("eth_sendRawTransaction", cltest.NewHash())

				require.NoError(t, app.StartAndConnect())
				receipt, _, err := txm.BumpGasUntilSafe(a.Hash)

				receiptPresent := receipt != nil
				require.Equal(t, test.wantReceipt, receiptPresent)
				cltest.AssertError(t, test.wantErrored, err)
			})
		})
	}
}

func TestTxManager_CheckAttempt(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	config := store.Config

	ethMock := app.EthMock
	ethMock.Register("eth_getTransactionCount", "0x0")
	ethMock.Register("eth_chainId", config.ChainID())
	require.NoError(t, app.StartAndConnect())

	txm := store.TxManager

	from := cltest.GetAccountAddress(t, store)
	sentAt := uint64(14770)
	hash := cltest.NewHash()
	gasBumpThreshold := sentAt + config.EthGasBumpThreshold()

	tx := cltest.CreateTx(t, store, from, sentAt)
	require.Len(t, tx.Attempts, 1)

	// Initial check, no receipt, no change of the block height
	retrievedReceipt := models.TxReceipt{}
	ethMock.Register("eth_getTransactionReceipt", retrievedReceipt)

	receipt, state, err := txm.CheckAttempt(tx.Attempts[0], sentAt)
	require.NoError(t, err)
	assert.Equal(t, strpkg.Unconfirmed, state)
	assert.Equal(t, receipt, &retrievedReceipt)

	ethMock.EventuallyAllCalled(t)

	// A receipt is found, but is not yet safe
	retrievedReceipt = models.TxReceipt{Hash: hash, BlockNumber: cltest.Int(sentAt)}
	ethMock.Register("eth_getTransactionReceipt", retrievedReceipt)

	receipt, state, err = txm.CheckAttempt(tx.Attempts[0], sentAt)
	require.NoError(t, err)
	assert.Equal(t, strpkg.Confirmed, state)
	assert.Equal(t, receipt, &retrievedReceipt)

	ethMock.EventuallyAllCalled(t)

	// A receipt is found, and now is safe
	ethMock.Register("eth_getTransactionReceipt", retrievedReceipt)

	receipt, state, err = txm.CheckAttempt(tx.Attempts[0], sentAt+gasBumpThreshold)
	require.NoError(t, err)
	assert.Equal(t, strpkg.Safe, state)
	assert.Equal(t, receipt, &retrievedReceipt)

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_CheckAttempt_error(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	ethMock := app.EthMock
	ethMock.Register("eth_getTransactionCount", "0x0")
	ethMock.Register("eth_chainId", store.Config.ChainID())
	require.NoError(t, app.StartAndConnect())

	txm := store.TxManager

	sentAt := uint64(14770)

	// Initial check, no receipt, no change of the block height
	ethMock.RegisterError("eth_getTransactionReceipt", "that aint gonna work chief")

	txAttempt := &models.TxAttempt{}
	receipt, state, err := txm.CheckAttempt(txAttempt, sentAt)
	require.Error(t, err)
	assert.Equal(t, strpkg.Unknown, state)
	assert.Nil(t, receipt)

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_Register(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ethMock := &cltest.EthMock{}
	txm := strpkg.NewEthTxManager(
		&eth.CallerSubscriberClient{CallerSubscriber: ethMock},
		orm.NewConfig(),
		nil,
		store.ORM,
	)

	ethMock.Register("eth_getTransactionCount", `0x2D0`)
	account := accounts.Account{Address: common.HexToAddress("0xbf4ed7b27f1d666546e30d74d50d173d20bca754")}
	txm.Register([]accounts.Account{account})
	txm.Connect(cltest.Head(1))
	ethMock.EventuallyAllCalled(t)

	aa := txm.NextActiveAccount()
	assert.Equal(t, account.Address, aa.Address)
	assert.Equal(t, uint64(0x2d0), aa.Nonce())
}

func TestTxManager_NextActiveAccount_RoundRobin(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ethMock := &cltest.EthMock{}
	txm := strpkg.NewEthTxManager(
		&eth.CallerSubscriberClient{CallerSubscriber: ethMock},
		orm.NewConfig(),
		nil,
		store.ORM,
	)

	accounts := []accounts.Account{
		accounts.Account{Address: common.HexToAddress("0xbf4ed7b27f1d666546e30d74d50d173d20bca001")},
		accounts.Account{Address: common.HexToAddress("0xbf4ed7b27f1d666546e30d74d50d173d20bca002")},
	}

	ethMock.Register("eth_getTransactionCount", `0x1D0`)
	ethMock.Register("eth_getTransactionCount", `0x2D0`)

	txm.Register(accounts)
	txm.Connect(cltest.Head(1))
	ethMock.EventuallyAllCalled(t)

	a0 := txm.NextActiveAccount()
	assert.Equal(t, accounts[0].Address, a0.Address)
	assert.Equal(t, uint64(0x1d0), a0.Nonce())

	a1 := txm.NextActiveAccount()
	assert.Equal(t, accounts[1].Address, a1.Address)
	assert.Equal(t, uint64(0x2d0), a1.Nonce())

	a2 := txm.NextActiveAccount()
	assert.Equal(t, a0, a2)
}

func TestTxManager_ReloadNonce(t *testing.T) {
	t.Parallel()

	ethClient := new(mocks.Client)
	txm := strpkg.NewEthTxManager(
		ethClient,
		orm.NewConfig(),
		nil,
		nil,
	)

	from := common.HexToAddress("0xbf4ed7b27f1d666546e30d74d50d173d20bca754")
	account := accounts.Account{Address: from}
	ma := strpkg.NewManagedAccount(account, 0)

	nonce := uint64(234)
	ethClient.On("GetNonce", from).Return(nonce, nil)

	err := ma.ReloadNonce(txm)
	assert.NoError(t, err)

	assert.Equal(t, account.Address, ma.Address)
	assert.Equal(t, nonce, ma.Nonce())

	ethClient.AssertExpectations(t)
}

func TestManagedAccount_GetAndIncrementNonce_YieldsCurrentNonceAndIncrements(t *testing.T) {
	account := accounts.Account{Address: common.HexToAddress("0xbf4ed7b27f1d666546e30d74d50d173d20bca754")}
	managedAccount := strpkg.NewManagedAccount(account, 0)

	managedAccount.GetAndIncrementNonce(func(y uint64) error {
		assert.Equal(t, uint64(0), y)
		return nil
	})
	assert.Equal(t, uint64(1), managedAccount.Nonce())

	managedAccount.GetAndIncrementNonce(func(y uint64) error {
		assert.Equal(t, uint64(1), y)
		return nil
	})
	assert.Equal(t, uint64(2), managedAccount.Nonce())
}

func TestManagedAccount_GetAndIncrementNonce_DoesNotIncrementWhenCallbackThrowsException(t *testing.T) {
	account := accounts.Account{Address: common.HexToAddress("0xbf4ed7b27f1d666546e30d74d50d173d20bca754")}
	managedAccount := strpkg.NewManagedAccount(account, 0)

	err := managedAccount.GetAndIncrementNonce(func(y uint64) error {
		assert.Equal(t, uint64(0), y)
		return errors.New("Should not increment")
	})
	assert.Error(t, err)
	err = managedAccount.GetAndIncrementNonce(func(y uint64) error {
		assert.Equal(t, uint64(0), y)
		return errors.New("Should not increment again")
	})
	assert.Error(t, err)
	assert.Equal(t, uint64(0), managedAccount.Nonce())
}

func TestTxManager_LogsETHAndLINKBalancesAfterSuccessfulTx(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ethClient := new(mocks.Client)

	config := cltest.NewTestConfig(t)
	keyStore := strpkg.NewKeyStore(config.KeysDir())
	account, err := keyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	require.NoError(t, keyStore.Unlock(cltest.Password))
	manager := strpkg.NewEthTxManager(ethClient, config, keyStore, store.ORM)

	from := account.Address
	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	nonce := uint64(256)
	sentAt := uint64(1)
	confirmedAt := sentAt + config.MinRequiredOutgoingConfirmations()

	manager.Register(keyStore.Accounts())
	require.NoError(t, err)

	ethClient.On("GetNonce", from).Return(nonce, nil)

	err = manager.Connect(cltest.Head(sentAt))
	require.NoError(t, err)

	ethClient.On("SendRawTx", mock.Anything).Return(cltest.NewHash(), nil)

	tx, err := manager.CreateTx(to, data)
	require.NoError(t, err)

	ntx, err := store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Equal(t, nonce, ntx.Nonce)
	assert.Equal(t, data, ntx.Data)
	assert.Equal(t, to, ntx.To)
	assert.Equal(t, from, ntx.From)
	assert.Equal(t, nonce, ntx.Nonce)
	assert.Len(t, ntx.Attempts, 1)

	confirmedReceipt := models.TxReceipt{
		Hash:        tx.Attempts[0].Hash,
		BlockNumber: cltest.Int(sentAt),
	}
	manager.OnNewLongestChain(*cltest.Head(confirmedAt))
	ethClient.On("GetTxReceipt", tx.Attempts[0].Hash).Return(&confirmedReceipt, nil)

	receipt, state, err := manager.BumpGasUntilSafe(tx.Attempts[0].Hash)
	require.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, strpkg.Safe, state)

	ethClient.AssertExpectations(t)
}

func TestTxManager_CreateTxWithGas(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store
	config := store.Config
	manager := store.TxManager

	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	assert.NoError(t, err)
	nonce := uint64(256)
	ethMock := app.EthMock
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce))
		ethMock.Register("eth_chainId", config.ChainID())
	})
	require.NoError(t, app.Store.ORM.IdempotentInsertHead(*cltest.Head(1)))
	assert.NoError(t, app.StartAndConnect())

	customGasPrice := utils.NewBig(big.NewInt(1337))
	customGasLimit := uint64(10009)

	defaultGasPrice := utils.NewBig(config.EthGasPriceDefault())

	tests := []struct {
		name             string
		gasPrice         *utils.Big
		gasLimit         uint64
		expectedGasPrice *utils.Big
		expectedGasLimit uint64
	}{
		{"not dev", customGasPrice, customGasLimit, defaultGasPrice, config.EthGasLimitDefault()},
		{"not dev not set", nil, 0, defaultGasPrice, config.EthGasLimitDefault()},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ethMock.Context("manager.CreateTx", func(ethMock *cltest.EthMock) {
				ethMock.Register("eth_sendRawTransaction", cltest.NewHash())
			})

			tx, err := manager.CreateTxWithGas(null.String{}, to, data, test.gasPrice.ToInt(), test.gasLimit)
			require.NoError(t, err)
			assert.Equal(t, test.expectedGasLimit, tx.GasLimit)

			require.Len(t, tx.Attempts, 1)
			assert.Equal(t, test.expectedGasPrice, tx.Attempts[0].GasPrice)

			ethMock.EventuallyAllCalled(t)
		})
	}
}

func TestTxManager_RebroadcastUnconfirmedTxsOnReconnect(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ethClient := new(mocks.Client)

	config := cltest.NewTestConfig(t)
	config.Set("CHAINLINK_TX_ATTEMPT_LIMIT", 1)
	keyStore := strpkg.NewKeyStore(config.KeysDir())
	_, err := keyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	require.NoError(t, keyStore.Unlock(cltest.Password))
	manager := strpkg.NewEthTxManager(ethClient, config, keyStore, store.ORM)

	to := cltest.NewAddress()
	data := hexutil.MustDecode("0x0000abcdef")
	sentAt := uint64(1)

	manager.Register(keyStore.Accounts())

	ethClient.On("GetNonce", mock.Anything).Times(2).Return(uint64(0), nil)

	err = manager.Connect(cltest.Head(sentAt))
	require.NoError(t, err)

	hash := cltest.NewHash()
	ethClient.On("SendRawTx", mock.Anything).Return(hash, nil)

	_, err = manager.CreateTx(to, data)
	require.NoError(t, err)

	manager.Disconnect()

	ethClient.On("SendRawTx", mock.Anything).Return(hash, nil)
	err = manager.Connect(cltest.Head(sentAt))
	require.NoError(t, err)

	ethClient.AssertExpectations(t)
}

func TestTxManager_BumpGasByIncrement(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ethClient := new(mocks.Client)

	config := cltest.NewTestConfig(t)
	config.Set("CHAINLINK_TX_ATTEMPT_LIMIT", 1)
	config.Set("ETH_GAS_PRICE_DEFAULT", 1)
	config.Set("ETH_GAS_BUMP_PERCENT", 10)
	keyStore := strpkg.NewKeyStore(config.KeysDir())
	txm := strpkg.NewEthTxManager(ethClient, config, keyStore, store.ORM)

	tests := []struct {
		name                   string
		originalGasPrice       *big.Int
		expectedBumpedGasPrice *big.Int
	}{
		{"bumping gas from 5Gwei to 10Gwei", big.NewInt(5000000000), big.NewInt(10000000000)},
		{"bumping gas from 100Gwei to 110Gwei", big.NewInt(100000000000), big.NewInt(110000000000)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := txm.BumpGasByIncrement(test.originalGasPrice)
			assert.Equal(t, test.expectedBumpedGasPrice.String(), actual.String())
		})
	}
}
