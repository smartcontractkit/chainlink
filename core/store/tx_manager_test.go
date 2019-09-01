package store_test

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/store"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

func TestTxManager_CreateTx_Success(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eth := mocks.NewMockEthClient(ctrl)

	config := cltest.NewTestConfig(t)
	keyStore := strpkg.NewKeyStore(config.KeysDir())
	account, err := keyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	require.NoError(t, keyStore.Unlock(cltest.Password))
	manager := strpkg.NewEthTxManager(eth, config, keyStore, store.ORM)

	from := account.Address
	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	nonce := uint64(256)

	manager.Register(keyStore.Accounts())
	require.NoError(t, err)

	eth.EXPECT().GetNonce(from).Return(nonce, nil)

	err = manager.Connect(cltest.Head(nonce))
	require.NoError(t, err)

	eth.EXPECT().SendRawTx(gomock.Any())

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
}

func TestTxManager_CreateTx_RoundRobinSuccess(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eth := mocks.NewMockEthClient(ctrl)

	config := cltest.NewTestConfig(t)
	keyStore := strpkg.NewKeyStore(config.KeysDir())

	// Add two accounts
	_, err := keyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	_, err = keyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	require.NoError(t, keyStore.Unlock(cltest.Password))

	manager := strpkg.NewEthTxManager(eth, config, keyStore, store.ORM)

	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	nonce := uint64(256)
	sentAt := uint64(1)
	bumpAt := sentAt + config.EthGasBumpThreshold()

	manager.Register(keyStore.Accounts())
	require.NoError(t, err)

	eth.EXPECT().GetNonce(gomock.Any()).Return(nonce, nil).Times(2)

	err = manager.Connect(cltest.Head(sentAt))
	require.NoError(t, err)

	eth.EXPECT().SendRawTx(gomock.Any()).Return(cltest.NewHash(), nil)

	createdTx1, err := manager.CreateTx(to, data)
	require.NoError(t, err)

	ntx, err := store.FindTx(createdTx1.ID)
	require.NoError(t, err)
	assert.Len(t, ntx.Attempts, 1)

	manager.OnNewHead(cltest.Head(bumpAt))
	eth.EXPECT().GetTxReceipt(createdTx1.Attempts[0].Hash).Return(&models.TxReceipt{}, nil)
	eth.EXPECT().SendRawTx(gomock.Any()).Return(cltest.NewHash(), nil)

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
	eth.EXPECT().SendRawTx(gomock.Any()).Return(cltest.NewHash(), nil)

	createdTx2, err := manager.CreateTx(to, data)
	assert.NoError(t, err)
	require.NotEqual(t, createdTx1.From.Hex(), createdTx2.From.Hex(), "should come from a different account")
}

func TestTxManager_CreateTx_BreakTxAttemptLimit(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eth := mocks.NewMockEthClient(ctrl)

	config := cltest.NewTestConfig(t)
	config.Set("CHAINLINK_TX_ATTEMPT_LIMIT", 1)
	keyStore := strpkg.NewKeyStore(config.KeysDir())
	account, err := keyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	require.NoError(t, keyStore.Unlock(cltest.Password))
	manager := strpkg.NewEthTxManager(eth, config, keyStore, store.ORM)

	from := account.Address
	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	nonce := uint64(256)
	sentAt := uint64(1)
	bumpAt := sentAt + config.EthGasBumpThreshold()
	bumpAgainAt := sentAt + 2*config.EthGasBumpThreshold()

	manager.Register(keyStore.Accounts())
	require.NoError(t, err)

	eth.EXPECT().GetNonce(from).Return(nonce, nil)

	err = manager.Connect(cltest.Head(sentAt))
	require.NoError(t, err)

	eth.EXPECT().SendRawTx(gomock.Any())

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

	manager.OnNewHead(cltest.Head(bumpAt))
	eth.EXPECT().GetTxReceipt(gomock.Any()).Return(&models.TxReceipt{}, nil)
	eth.EXPECT().SendRawTx(gomock.Any()).Return(tx.Attempts[0].Hash, nil)

	receipt, state, err := manager.BumpGasUntilSafe(tx.Attempts[0].Hash)
	require.NoError(t, err)
	assert.Nil(t, receipt)
	assert.Equal(t, strpkg.Unconfirmed, state)

	manager.OnNewHead(cltest.Head(bumpAgainAt))
	eth.EXPECT().GetTxReceipt(gomock.Any()).Return(&models.TxReceipt{}, nil)
	eth.EXPECT().GetTxReceipt(gomock.Any()).Return(&models.TxReceipt{}, nil)

	receipt, state, err = manager.BumpGasUntilSafe(tx.Attempts[0].Hash)
	require.NoError(t, err)
	assert.Nil(t, receipt)
	assert.Equal(t, strpkg.Unconfirmed, state)
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
	ethMock := app.MockEthCallerSubscriber()
	ethMock.Context("app.StartAndConnect()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce))
		ethMock.Register("eth_chainId", store.Config.ChainID())
	})
	require.NoError(t, app.Store.ORM.CreateHead(cltest.Head(sentAt)))
	assert.NoError(t, app.StartAndConnect())

	ethMock.Context("manager.CreateTx#1", func(ethMock *cltest.EthMock) {
		ethMock.RegisterError("eth_sendRawTransaction", "invalid transaction")
	})

	_, err = manager.CreateTx(to, data)
	assert.Error(t, err)

	txs, _, err := store.Transactions(0, 10)
	assert.NoError(t, err)
	assert.Len(t, txs, 1)

	txAttempts, _, err := store.TxAttempts(0, 100)
	assert.NoError(t, err)
	assert.Len(t, txAttempts, 1)

	assert.Equal(t, txs[0].Hash, txAttempts[0].Hash)

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
		{"parity", "Transaction nonce is too low. Try incrementing the nonce"},
		{"parity", "Transaction with the same hash was already imported"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			app, cleanup := cltest.NewApplicationWithKey(t)
			defer cleanup()
			store := app.Store
			manager := store.TxManager

			from := cltest.GetAccountAddress(t, store)
			to := cltest.NewAddress()
			data, err := hex.DecodeString("0000abcdef")
			assert.NoError(t, err)
			ethMock := app.MockEthCallerSubscriber()

			nonce1 := uint64(256)
			ethMock.Context("app.StartAndConnect()", func(ethMock *cltest.EthMock) {
				ethMock.Register("eth_chainId", store.Config.ChainID())
				ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce1))
			})
			sentAt := uint64(23456)
			require.NoError(t, app.Store.ORM.CreateHead(cltest.Head(sentAt)))
			assert.NoError(t, app.StartAndConnect())

			ethMock.Context("manager.CreateTx#1", func(ethMock *cltest.EthMock) {
				ethMock.RegisterError("eth_sendRawTransaction", test.ethClientErrorMsg)
			})

			hash := cltest.NewHash()
			nonce2 := uint64(257)
			ethMock.Context("manager.CreateTx#2", func(ethMock *cltest.EthMock) {
				ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce2))
				ethMock.Register("eth_sendRawTransaction", hash)
			})

			a, err := manager.CreateTx(to, data)
			require.NoError(t, err)
			tx, err := store.FindTx(a.ID)
			require.NoError(t, err)
			assert.Equal(t, nonce2, tx.Nonce)
			assert.Equal(t, data, tx.Data)
			assert.Equal(t, to, tx.To)
			assert.Equal(t, from, tx.From)
			assert.Len(t, tx.Attempts, 1)

			ethMock.EventuallyAllCalled(t)
		})
	}
}

func TestTxManager_CreateTx_NonceTooLowReloadLimit(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store
	manager := store.TxManager

	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	assert.NoError(t, err)
	ethMock := app.MockEthCallerSubscriber()

	nonce1 := uint64(256)
	ethMock.Context("app.StartAndConnect()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce1))
		ethMock.Register("eth_chainId", store.Config.ChainID())
	})
	sentAt := uint64(23456)
	require.NoError(t, app.Store.ORM.CreateHead(cltest.Head(sentAt)))
	assert.NoError(t, app.StartAndConnect())

	ethMock.Context("manager.CreateTx#1", func(ethMock *cltest.EthMock) {
		ethMock.RegisterError("eth_sendRawTransaction", "nonce is too low")
	})

	nonce2 := uint64(257)
	ethMock.Context("manager.CreateTx#2", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce2))
		ethMock.RegisterError("eth_sendRawTransaction", "nonce is too low")
	})

	_, err = manager.CreateTx(to, data)
	assert.EqualError(
		t,
		err,
		"Transaction reattempt limit reached for 'nonce is too low' error. Limit: 1",
	)

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_CreateTx_ErrPendingConnection(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store
	manager := store.TxManager

	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	assert.NoError(t, err)

	_, err = manager.CreateTx(to, data)
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
	ethMock := app.MockEthCallerSubscriber(cltest.Strict)
	ethMock.Register("eth_getTransactionCount", "0x0")
	ethMock.Register("eth_chainId", config.ChainID())
	require.NoError(t, app.Store.ORM.CreateHead(cltest.Head(gasThreshold-1)))
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
	ethMock := app.MockEthCallerSubscriber(cltest.Strict)
	ethMock.Register("eth_getTransactionCount", "0x0")
	ethMock.Register("eth_chainId", config.ChainID())
	require.NoError(t, app.Store.ORM.CreateHead(cltest.Head(gasThreshold)))
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
	ethMock := app.MockEthCallerSubscriber(cltest.Strict)
	ethMock.Register("eth_getTransactionCount", "0x0")
	ethMock.Register("eth_chainId", config.ChainID())
	require.NoError(t, app.Store.ORM.CreateHead(cltest.Head(gasThreshold+1)))
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

func TestTxManager_BumpGasUntilSafe_confirmed_lessThanGasThreshold(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	config := store.Config

	minConfs := config.MinOutgoingConfirmations() - 1
	sentAt := uint64(23456)
	gasThreshold := sentAt + config.EthGasBumpThreshold()
	ethMock := app.MockEthCallerSubscriber(cltest.Strict)
	ethMock.Register("eth_getTransactionCount", "0x0")
	ethMock.Register("eth_chainId", config.ChainID())
	require.NoError(t, app.Store.ORM.CreateHead(cltest.Head(gasThreshold+minConfs-1)))
	require.NoError(t, app.StartAndConnect())

	txm := store.TxManager
	from := cltest.GetAccountAddress(t, store)

	tx := cltest.CreateTx(t, store, from, sentAt)
	require.Greater(t, len(tx.Attempts), 0)

	ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{Hash: cltest.NewHash(), BlockNumber: cltest.Int(gasThreshold)})

	receipt, state, err := txm.BumpGasUntilSafe(tx.Attempts[0].Hash)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, strpkg.Confirmed, state)

	tx, err = store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Len(t, tx.Attempts, 1)

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_BumpGasUntilSafe_confirmed_atGasBumpThreshold(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	store := app.Store
	config := store.Config

	from := cltest.GetAccountAddress(t, store)
	sentAt := uint64(23456)
	gasThreshold := sentAt + config.EthGasBumpThreshold()
	minConfs := config.MinOutgoingConfirmations() - 1
	ethMock := app.MockEthCallerSubscriber(cltest.Strict)
	ethMock.Register("eth_getTransactionCount", "0x0")
	ethMock.Register("eth_chainId", config.ChainID())
	require.NoError(t, app.Store.ORM.CreateHead(cltest.Head(gasThreshold+minConfs+1)))
	require.NoError(t, app.StartAndConnect())

	txm := store.TxManager

	tx := cltest.CreateTx(t, store, from, sentAt)
	require.Greater(t, len(tx.Attempts), 0)

	ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{Hash: cltest.NewHash(), BlockNumber: cltest.Int(gasThreshold)})
	ethMock.Register("eth_getBalance", "0x100")
	ethMock.Register("eth_call", "0x100")

	receipt, state, err := txm.BumpGasUntilSafe(tx.Attempts[0].Hash)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, strpkg.Safe, state)

	tx, err = store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Len(t, tx.Attempts, 1)

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_BumpGasUntilSafe_confirmed_exceedsGasBumpThreshold(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()
	store := app.Store
	config := store.Config

	sentAt := uint64(23456)
	gasThreshold := sentAt + config.EthGasBumpThreshold()
	minConfs := config.MinOutgoingConfirmations() - 1
	require.NoError(t, app.Store.ORM.CreateHead(cltest.Head(gasThreshold+minConfs+1)))
	ethMock, err := app.MockStartAndConnect()
	require.NoError(t, err)

	txm := store.TxManager
	from := cltest.GetAccountAddress(t, store)

	tx := cltest.CreateTx(t, store, from, sentAt)
	require.Greater(t, len(tx.Attempts), 0)

	ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{Hash: cltest.NewHash(), BlockNumber: cltest.Int(gasThreshold)})
	ethMock.Register("eth_getBalance", "0x100")
	ethMock.Register("eth_call", "0x100")

	receipt, state, err := txm.BumpGasUntilSafe(tx.Attempts[0].Hash)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, strpkg.Safe, state)

	tx, err = store.FindTx(tx.ID)
	require.NoError(t, err)
	assert.Len(t, tx.Attempts, 1)

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_BumpGasUntilSafe_laterConfirmedTx(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	ethMock, err := app.MockStartAndConnect()
	require.NoError(t, err)

	store := app.Store
	txm := store.TxManager
	from := cltest.GetAccountAddress(t, store)
	sentAt := uint64(12345)

	tx1 := cltest.CreateTxWithNonce(t, store, from, sentAt, 1)
	tx2 := cltest.CreateTxWithNonce(t, store, from, sentAt, 2)
	tx2a := tx2.Attempts[0]
	tx2a.Confirmed = true
	assert.NoError(t, store.MarkTxSafe(tx2, tx2a))

	ethMock.Register("eth_getTransactionReceipt", models.TxReceipt{})

	receipt, state, err := txm.BumpGasUntilSafe(tx1.Attempts[0].Hash)
	assert.Nil(t, receipt)
	assert.Equal(t, strpkg.Confirmed, state)
	assert.Error(t, err)

	tx, err := store.FindTx(tx1.ID)
	require.NoError(t, err)
	assert.True(t, tx.Confirmed)
	assert.Equal(t, tx.Hash.Hex(), "0x0000000000000000000000000000000000000000000000000000000000000000")
	assert.Len(t, tx.Attempts, 1)

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_BumpGasUntilSafe_erroring(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	sentAt1 := uint64(23456)
	sentAt2 := sentAt1 + config.EthGasBumpThreshold()
	confirmedAt := sentAt2 + 1
	safeAt := confirmedAt + config.MinOutgoingConfirmations()

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
			ethMock.Register("eth_call", "0x0100")
			ethMock.Register("eth_getBalance", "0x0100")
			ethMock.Register("eth_getTransactionReceipt", confedReceipt)
		}, true, false},
		{"later conf, no error", (safeAt + 1), func(ethMock *cltest.EthMock) {
			ethMock.Register("eth_call", "0x0100")
			ethMock.Register("eth_getBalance", "0x0100")
			ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(0))
			ethMock.Register("eth_getTransactionReceipt", confedReceipt)
		}, true, false},
		{"later conf, early error", (safeAt + 1), func(ethMock *cltest.EthMock) {
			ethMock.Register("eth_getTransactionCount", "0x0")
			ethMock.Register("eth_call", "0x0100")
			ethMock.Register("eth_getBalance", "0x0100")
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
			a, err := store.AddTxAttempt(tx, tx.EthTx(big.NewInt(2)), sentAt2)
			assert.NoError(t, err)

			ethMock := app.MockEthCallerSubscriber(cltest.Strict)
			ethMock.ShouldCall(test.mockSetup).During(func() {
				require.NoError(t, app.Store.ORM.CreateHead(cltest.Head(test.blockHeight)))
				ethMock.Register("eth_chainId", store.Config.ChainID())
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

	ethMock := app.MockEthCallerSubscriber(cltest.Strict)
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
	ethMock := app.MockEthCallerSubscriber(cltest.Strict)
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

	ethMock := &cltest.EthMock{}
	txm := store.NewEthTxManager(
		&strpkg.EthCallerSubscriber{CallerSubscriber: ethMock},
		orm.NewConfig(),
		nil,
		nil,
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

	ethMock := &cltest.EthMock{}
	txm := store.NewEthTxManager(
		&strpkg.EthCallerSubscriber{CallerSubscriber: ethMock},
		orm.NewConfig(),
		nil,
		nil,
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

	ethMock := &cltest.EthMock{}
	txm := store.NewEthTxManager(
		&strpkg.EthCallerSubscriber{CallerSubscriber: ethMock},
		orm.NewConfig(),
		nil,
		nil,
	)

	account := accounts.Account{Address: common.HexToAddress("0xbf4ed7b27f1d666546e30d74d50d173d20bca754")}
	ma := strpkg.NewManagedAccount(account, 0)

	ethMock.Register("eth_getTransactionCount", `0x2D1`)
	assert.NoError(t, ma.ReloadNonce(txm))
	ethMock.EventuallyAllCalled(t)

	assert.Equal(t, account.Address, ma.Address)
	assert.Equal(t, uint64(0x2d1), ma.Nonce())
}

func TestTxManager_WithdrawLink_HappyPath(t *testing.T) {
	t.Parallel()
	config, configCleanup := cltest.NewConfig(t)
	defer configCleanup()
	oca := common.HexToAddress("0xDEADB3333333F")
	config.Set("ORACLE_CONTRACT_ADDRESS", &oca)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config)
	defer cleanup()

	txm := app.Store.TxManager

	from := cltest.GetAccountAddress(t, app.GetStore())
	to := cltest.NewAddress()
	hash := cltest.NewHash()
	sentAt := uint64(23456)
	nonce := uint64(256)
	ethMock := app.MockEthCallerSubscriber()
	ethMock.Context("app.StartAndConnect()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce))
		ethMock.Register("eth_chainId", config.ChainID())
	})
	require.NoError(t, app.Store.ORM.CreateHead(cltest.Head(sentAt)))
	assert.NoError(t, app.StartAndConnect())

	ethMock.Context("txm.CreateTx#1", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_sendRawTransaction", hash)
	})

	wr := models.WithdrawalRequest{
		DestinationAddress: to,
		Amount:             assets.NewLink(10),
	}

	hash, err := txm.WithdrawLINK(wr)
	assert.NoError(t, err)
	assert.True(t, ethMock.AllCalled(), "Not Called")

	transactions, err := app.Store.TxFrom(from)
	require.NoError(t, err)
	tx := transactions[0]
	assert.Equal(t, hash, tx.Hash)
	assert.Equal(t, nonce, tx.Nonce)
}

func TestTxManager_WithdrawLink_Unconfigured_Oracle(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKey(t)
	defer cleanup()

	nonce := uint64(256)
	ethMock := app.MockEthCallerSubscriber()
	ethMock.Context("app.StartAndConnect()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce))
		ethMock.Register("eth_chainId", app.Store.Config.ChainID())
	})
	assert.NoError(t, app.StartAndConnect())

	wr := models.WithdrawalRequest{
		DestinationAddress: cltest.NewAddress(),
		Amount:             assets.NewLink(10),
	}

	_, err := app.Store.TxManager.WithdrawLINK(wr)
	assert.EqualError(t, err, "OracleContractAddress not set; cannot withdraw")
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eth := mocks.NewMockEthClient(ctrl)

	config := cltest.NewTestConfig(t)
	keyStore := strpkg.NewKeyStore(config.KeysDir())
	account, err := keyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	require.NoError(t, keyStore.Unlock(cltest.Password))
	manager := strpkg.NewEthTxManager(eth, config, keyStore, store.ORM)

	from := account.Address
	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	nonce := uint64(256)
	sentAt := uint64(1)
	confirmedAt := sentAt + config.MinOutgoingConfirmations()

	manager.Register(keyStore.Accounts())
	require.NoError(t, err)

	eth.EXPECT().GetNonce(from).Return(nonce, nil)

	err = manager.Connect(cltest.Head(sentAt))
	require.NoError(t, err)

	eth.EXPECT().SendRawTx(gomock.Any())

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
	manager.OnNewHead(cltest.Head(confirmedAt))
	eth.EXPECT().GetTxReceipt(tx.Attempts[0].Hash).Return(&confirmedReceipt, nil)
	eth.EXPECT().GetERC20Balance(from, gomock.Any())
	eth.EXPECT().GetEthBalance(from)

	receipt, state, err := manager.BumpGasUntilSafe(tx.Attempts[0].Hash)
	require.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, strpkg.Safe, state)
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
	ethMock := app.MockEthCallerSubscriber()
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce))
		ethMock.Register("eth_chainId", config.ChainID())
	})
	require.NoError(t, app.Store.ORM.CreateHead(cltest.Head(1)))
	assert.NoError(t, app.StartAndConnect())

	customGasPrice := models.NewBig(big.NewInt(1337))
	customGasLimit := uint64(10009)

	defaultGasPrice := models.NewBig(config.EthGasPriceDefault())

	tests := []struct {
		name             string
		dev              bool
		gasPrice         *models.Big
		gasLimit         uint64
		expectedGasPrice *models.Big
		expectedGasLimit uint64
	}{
		{"dev", true, customGasPrice, customGasLimit, customGasPrice, customGasLimit},
		{"dev but not set", true, nil, 0, defaultGasPrice, strpkg.DefaultGasLimit},
		{"not dev", false, customGasPrice, customGasLimit, defaultGasPrice, strpkg.DefaultGasLimit},
		{"not dev not set", false, nil, 0, defaultGasPrice, strpkg.DefaultGasLimit},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config.Set("CHAINLINK_DEV", test.dev)

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

func TestGetContract(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		contract  string
		expectErr bool
	}{
		{"Get Oracle contract", "Oracle", false},
		{"Get non-existent contract", "not-a-contract", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			contract, err := strpkg.GetContract(test.contract)
			if test.expectErr {
				assert.Error(t, err)
				assert.Nil(t, contract)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, contract)
			}
		})
	}
}

func TestContract_EncodeMessageCall(t *testing.T) {
	t.Parallel()

	// Test with the Oracle contract
	oracle, err := strpkg.GetContract("Oracle")
	assert.NoError(t, err)

	tests := []struct {
		name      string
		method    string
		args      []interface{}
		expectErr bool
	}{
		{"Withdraw LINK", "withdraw", []interface{}{cltest.NewAddress(), (*big.Int)(assets.NewLink(10))}, false},
		{"Non-existent method", "not-a-method", []interface{}{cltest.NewAddress(), (*big.Int)(assets.NewLink(10))}, true},
		{"Too few arguments", "withdraw", []interface{}{cltest.NewAddress()}, true},
		{"Too many arguments", "withdraw", []interface{}{cltest.NewAddress(), (*big.Int)(assets.NewLink(10)), (*big.Int)(assets.NewLink(10))}, true},
		{"Incorrect argument types", "withdraw", []interface{}{(*big.Int)(assets.NewLink(10)), cltest.NewAddress()}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := oracle.EncodeMessageCall(test.method, test.args...)
			if test.expectErr {
				assert.Error(t, err)
				assert.Nil(t, data)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, data)
			}
		})
	}
}
