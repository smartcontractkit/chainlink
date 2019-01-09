package store_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxManager_CreateTx_Success(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	manager := store.TxManager

	from := cltest.GetAccountAddress(app.GetStore())
	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	assert.NoError(t, err)
	hash := cltest.NewHash()
	sentAt := uint64(23456)
	nonce := uint64(256)
	ethMock := app.MockEthClient()
	ethMock.Context("app.StartAndConnect()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce))
	})
	assert.NoError(t, app.StartAndConnect())

	require.True(t, manager.Connected())
	ethMock.Context("manager.CreateTx#1", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_sendRawTransaction", hash)
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	})

	tx, err := manager.CreateTx(to, data)
	assert.NoError(t, err)
	_, err = store.FindTx(tx.ID)
	assert.NoError(t, err)
	assert.Equal(t, nonce, tx.Nonce)
	assert.Equal(t, data, tx.Data)
	assert.Equal(t, to, tx.To)
	assert.Equal(t, from, tx.From)
	assert.Equal(t, nonce, tx.Nonce)
	attempts, err := store.TxAttemptsFor(tx.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(attempts))

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_CreateTx_RoundRobinSuccess(t *testing.T) {
	t.Parallel()
	config, _ := cltest.NewConfigWithPrivateKey() // second account
	app, cleanup := cltest.NewApplicationWithConfigAndKeyStore(config)
	defer cleanup()
	store := app.Store
	manager := store.TxManager
	accounts := store.KeyStore.Accounts()

	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	assert.NoError(t, err)
	sentAt := uint64(23456)
	ethMock := app.MockEthClient()
	ethMock.Context("app.StartAndConnect()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", "0x00")
		ethMock.Register("eth_getTransactionCount", "0x10")
	})
	assert.NoError(t, app.StartAndConnect())

	ethMock.Context("manager.CreateTx#1", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_sendRawTransaction", cltest.NewHash())
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	})

	createdTx1, err := manager.CreateTx(to, data)
	require.NoError(t, err)

	attempts, err := store.TxAttemptsFor(createdTx1.ID)
	require.NoError(t, err)
	require.Len(t, attempts, 1)

	// bump gas
	ethMock.Context("manager.bumpGas#1", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionReceipt", strpkg.TxReceipt{})
		ethMock.Register("eth_sendRawTransaction", cltest.NewHash())
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt+config.EthGasBumpThreshold()))
	})

	_, err = manager.EnsureConfirmed(attempts[0].Hash)
	require.NoError(t, err)

	// retrieve new gas bumped second attempt
	attempts, err = store.TxAttemptsFor(createdTx1.ID)
	require.NoError(t, err)
	require.Len(t, attempts, 2)
	a2 := attempts[1]

	// ensure gas bumped attempt does not round robin on the From Address
	// best way to ensure the same from address atm is to compare Hashes, since
	// tx attempts don't have From but rely on parent Tx model.
	etx := createdTx1.EthTx(a2.GasPrice)
	etx, err = store.KeyStore.SignTx(accounts[0], etx, config.ChainID())
	assert.Equal(t, etx.Hash().Hex(), a2.Hash.Hex(), "should be same since they have the same input, include From address")

	// ensure second tx round robins
	ethMock.Context("manager.CreateTx#2", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_sendRawTransaction", cltest.NewHash())
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	})

	createdTx2, err := manager.CreateTx(to, data)
	assert.NoError(t, err)

	require.NotEqual(t, createdTx1.From.Hex(), createdTx2.From.Hex(), "should come from a different account")

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_CreateTx_AttemptErrorDeletesTxAndDoesNotIncrementNonce(t *testing.T) {
	t.Parallel()

	config, configCleanup := cltest.NewConfig()
	defer configCleanup()

	app, cleanup := cltest.NewApplicationWithConfigAndKeyStore(config)
	defer cleanup()

	store := app.Store
	manager := store.TxManager

	from := cltest.GetAccountAddress(app.GetStore())
	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	assert.NoError(t, err)
	sentAt := uint64(23456)
	nonce := uint64(256)
	ethMock := app.MockEthClient()
	ethMock.Context("app.StartAndConnect()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce))
	})
	assert.NoError(t, app.StartAndConnect())

	ethMock.Context("manager.CreateTx#1", func(ethMock *cltest.EthMock) {
		ethMock.RegisterError("eth_sendRawTransaction", "invalid transaction")
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	})

	_, err = manager.CreateTx(to, data)
	assert.Error(t, err)

	txs, err := store.Transactions(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(txs))

	txAttempts, _, err := store.TxAttempts(0, 100)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(txAttempts))

	hash := cltest.NewHash()
	ethMock.Context("manager.CreateTx#2", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_sendRawTransaction", hash)
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
	})

	tx, err := manager.CreateTx(to, data)
	assert.NoError(t, err)
	_, err = store.FindTx(tx.ID)
	assert.NoError(t, err)
	assert.NoError(t, err)
	assert.Equal(t, nonce, tx.Nonce)
	assert.Equal(t, data, tx.Data)
	assert.Equal(t, to, tx.To)
	assert.Equal(t, from, tx.From)

	assert.Equal(t, nonce, tx.Nonce)
	attempts, err := store.TxAttemptsFor(tx.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(attempts))

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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			app, cleanup := cltest.NewApplicationWithKeyStore()
			defer cleanup()
			store := app.Store
			manager := store.TxManager

			from := cltest.GetAccountAddress(store)
			to := cltest.NewAddress()
			data, err := hex.DecodeString("0000abcdef")
			assert.NoError(t, err)
			ethMock := app.MockEthClient()

			nonce1 := uint64(256)
			ethMock.Context("app.StartAndConnect()", func(ethMock *cltest.EthMock) {
				ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce1))
			})
			assert.NoError(t, app.StartAndConnect())

			sentAt := uint64(23456)
			ethMock.Context("manager.CreateTx#1", func(ethMock *cltest.EthMock) {
				ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
				ethMock.RegisterError("eth_sendRawTransaction", test.ethClientErrorMsg)
			})

			hash := cltest.NewHash()
			nonce2 := uint64(257)
			ethMock.Context("manager.CreateTx#2", func(ethMock *cltest.EthMock) {
				ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
				ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce2))
				ethMock.Register("eth_sendRawTransaction", hash)
			})

			a, err := manager.CreateTx(to, data)
			assert.NoError(t, err)
			tx, err := store.FindTx(a.TxID)
			require.NoError(t, err)
			assert.NoError(t, err)
			assert.Equal(t, nonce2, tx.Nonce)
			assert.Equal(t, data, tx.Data)
			assert.Equal(t, to, tx.To)

			assert.Equal(t, from, tx.From)
			assert.Equal(t, nonce2, tx.Nonce)
			attempts, err := store.TxAttemptsFor(tx.ID)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(attempts))

			ethMock.EventuallyAllCalled(t)

		})
	}
}

func TestTxManager_CreateTx_NonceTooLowReloadLimit(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	manager := store.TxManager

	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	assert.NoError(t, err)
	ethMock := app.MockEthClient()

	nonce1 := uint64(256)
	ethMock.Context("app.StartAndConnect()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce1))
	})
	assert.NoError(t, app.StartAndConnect())

	sentAt := uint64(23456)
	ethMock.Context("manager.CreateTx#1", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
		ethMock.RegisterError("eth_sendRawTransaction", "nonce is too low")
	})

	nonce2 := uint64(257)
	ethMock.Context("manager.CreateTx#2", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce2))
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
		ethMock.RegisterError("eth_sendRawTransaction", "nonce is too low")
	})

	_, err = manager.CreateTx(to, data)
	assert.EqualError(
		t,
		err,
		"Transaction reattempt limit reached for 'nonce is too low' error. Limit: 1, Reattempt: 1",
	)

	ethMock.EventuallyAllCalled(t)
}

func TestTxManager_CreateTx_ErrPendingConnection(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	manager := store.TxManager

	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	assert.NoError(t, err)

	assert.NoError(t, app.Start())

	_, err = manager.CreateTx(to, data)
	assert.Equal(t, strpkg.ErrPendingConnection, err)
}

func TestTxManager_EnsureConfirmed(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	ethMock := app.MockEthClient()
	ethMock.Register("eth_getTransactionCount", "0x0")
	require.NoError(t, app.StartAndConnect())

	store := app.Store
	config := store.Config
	txm := store.TxManager
	from := cltest.GetAccountAddress(store)
	sentAt := uint64(23456)
	gasThreshold := sentAt + config.EthGasBumpThreshold()
	minConfs := config.MinOutgoingConfirmations() - 1

	tests := []struct {
		name             string
		currentHeight    uint64
		receipt          strpkg.TxReceipt
		sendsTransaction bool
		wantConfirmed    bool
		wantLength       int
	}{
		{"< gas bump threshold", (gasThreshold - 1), strpkg.TxReceipt{}, false, false, 1},
		{"== gas bump threshold", gasThreshold, strpkg.TxReceipt{}, true, false, 2},
		{"> gas bump threshold", (gasThreshold + 1), strpkg.TxReceipt{}, true, false, 2},
		{"confirmed && < min confs", (gasThreshold + minConfs - 1), strpkg.TxReceipt{Hash: cltest.NewHash(), BlockNumber: cltest.Int(gasThreshold)}, false, false, 1},
		{"confirmed && == min confs", (gasThreshold + minConfs), strpkg.TxReceipt{Hash: cltest.NewHash(), BlockNumber: cltest.Int(gasThreshold)}, false, true, 1},
		{"confirmed && > min confs", (gasThreshold + minConfs + 1), strpkg.TxReceipt{Hash: cltest.NewHash(), BlockNumber: cltest.Int(gasThreshold)}, false, true, 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tx := cltest.CreateTxAndAttempt(store, from, sentAt)
			attempts, err := store.TxAttemptsFor(tx.ID)
			assert.NoError(t, err)
			a := attempts[0]

			ethMock.Register("eth_getTransactionReceipt", test.receipt)
			ethMock.Register("eth_blockNumber", utils.Uint64ToHex(test.currentHeight))
			if test.sendsTransaction {
				ethMock.Register("eth_sendRawTransaction", cltest.NewHash())
			}

			confirmed, err := txm.EnsureConfirmed(a.Hash)
			assert.NoError(t, err)
			assert.Equal(t, test.wantConfirmed, confirmed)
			attempts, err = store.TxAttemptsFor(tx.ID)
			assert.NoError(t, err)
			assert.Equal(t, test.wantLength, len(attempts))

			ethMock.EventuallyAllCalled(t)
		})
	}
}

func TestTxManager_EnsureConfirmed_erroring(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig()
	defer cleanup()

	sentAt1 := uint64(23456)
	sentAt2 := sentAt1 + config.EthGasBumpThreshold()
	confirmedAt := sentAt2 + 1
	safeAt := confirmedAt + config.MinOutgoingConfirmations()

	nonConfedReceipt := strpkg.TxReceipt{}
	confedReceipt := strpkg.TxReceipt{Hash: cltest.NewHash(), BlockNumber: cltest.Int(confirmedAt)}

	tests := []struct {
		name          string
		blockHeight   uint64
		mockSetup     func(*cltest.EthMock)
		wantConfirmed bool
		wantErrored   bool
	}{
		{"no conf, no error", (sentAt2 + 1), func(ethMock *cltest.EthMock) {
			ethMock.Register("eth_getTransactionReceipt", nonConfedReceipt)
			ethMock.Register("eth_getTransactionReceipt", nonConfedReceipt)
		}, false, false},
		{"no conf, early error", (sentAt2 + 1), func(ethMock *cltest.EthMock) {
			ethMock.RegisterError("eth_getTransactionReceipt", "FUBAR")
			ethMock.Register("eth_getTransactionReceipt", nonConfedReceipt)
		}, false, true},
		{"no conf, later error", (sentAt2 + 1), func(ethMock *cltest.EthMock) {
			ethMock.Register("eth_getTransactionReceipt", nonConfedReceipt)
			ethMock.RegisterError("eth_getTransactionReceipt", "FUBAR")
		}, false, true},
		{"early conf, no error", (safeAt + 1), func(ethMock *cltest.EthMock) {
			ethMock.Register("eth_getTransactionReceipt", confedReceipt)
			ethMock.Register("eth_getTransactionReceipt", nonConfedReceipt)
		}, true, false},
		{"early conf, later error", (safeAt + 1), func(ethMock *cltest.EthMock) {
			ethMock.Register("eth_getTransactionReceipt", confedReceipt)
			ethMock.RegisterError("eth_getTransactionReceipt", "FUBAR")
		}, true, false},
		{"later conf, no error", (safeAt + 1), func(ethMock *cltest.EthMock) {
			ethMock.Register("eth_getTransactionReceipt", nonConfedReceipt)
			ethMock.Register("eth_getTransactionReceipt", confedReceipt)
		}, true, false},
		{"later conf, early error", (safeAt + 1), func(ethMock *cltest.EthMock) {
			ethMock.RegisterError("eth_getTransactionReceipt", "FUBAR")
			ethMock.Register("eth_getTransactionReceipt", confedReceipt)
		}, true, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, cleanup := cltest.NewApplicationWithKeyStore()
			defer cleanup()

			store := app.Store
			txm := store.TxManager
			from := cltest.GetAccountAddress(store)
			tx := cltest.CreateTxAndAttempt(store, from, sentAt1)
			a, err := store.AddTxAttempt(tx, tx.EthTx(big.NewInt(2)), sentAt2)
			assert.NoError(t, err)

			ethMock := app.MockEthClient()
			ethMock.Context("txm.EnsureConfirmed()", test.mockSetup)
			ethMock.Register("eth_blockNumber", utils.Uint64ToHex(test.blockHeight))

			require.NoError(t, app.StartAndConnect())
			confirmed, err := txm.EnsureConfirmed(a.Hash)
			assert.Equal(t, test.wantConfirmed, confirmed)
			cltest.AssertError(t, test.wantErrored, err)
		})
	}
}

func TestTxManager_Register(t *testing.T) {
	t.Parallel()

	ethMock := &cltest.EthMock{}
	txm := store.NewEthTxManager(
		&strpkg.EthClient{CallerSubscriber: ethMock},
		store.NewConfig(),
		nil,
		nil,
	)

	ethMock.Register("eth_getTransactionCount", `0x2D0`)
	account := accounts.Account{Address: common.HexToAddress("0xbf4ed7b27f1d666546e30d74d50d173d20bca754")}
	txm.Register([]accounts.Account{account})
	txm.Connect(cltest.IndexableBlockNumber(1))
	ethMock.EventuallyAllCalled(t)

	aa := txm.NextActiveAccount()
	assert.Equal(t, account.Address, aa.Address)
	assert.Equal(t, uint64(0x2d0), aa.GetNonce())
}

func TestTxManager_NextActiveAccount_RoundRobin(t *testing.T) {
	t.Parallel()

	ethMock := &cltest.EthMock{}
	txm := store.NewEthTxManager(
		&strpkg.EthClient{CallerSubscriber: ethMock},
		store.NewConfig(),
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
	txm.Connect(cltest.IndexableBlockNumber(1))
	ethMock.EventuallyAllCalled(t)

	a0 := txm.NextActiveAccount()
	assert.Equal(t, accounts[0].Address, a0.Address)
	assert.Equal(t, uint64(0x1d0), a0.GetNonce())

	a1 := txm.NextActiveAccount()
	assert.Equal(t, accounts[1].Address, a1.Address)
	assert.Equal(t, uint64(0x2d0), a1.GetNonce())

	a2 := txm.NextActiveAccount()
	assert.Equal(t, a0, a2)
}

func TestTxManager_ReloadNonce(t *testing.T) {
	t.Parallel()

	ethMock := &cltest.EthMock{}
	txm := store.NewEthTxManager(
		&strpkg.EthClient{CallerSubscriber: ethMock},
		store.NewConfig(),
		nil,
		nil,
	)

	account := accounts.Account{Address: common.HexToAddress("0xbf4ed7b27f1d666546e30d74d50d173d20bca754")}
	ma := strpkg.NewManagedAccount(account, 0)

	ethMock.Register("eth_getTransactionCount", `0x2D1`)
	assert.NoError(t, ma.ReloadNonce(txm))
	ethMock.EventuallyAllCalled(t)

	assert.Equal(t, account.Address, ma.Address)
	assert.Equal(t, uint64(0x2d1), ma.GetNonce())
}

func TestTxManager_WithdrawLink(t *testing.T) {
	t.Parallel()
	config, configCleanup := cltest.NewConfig()
	defer configCleanup()
	oca := common.HexToAddress("0xDEADB3333333F")
	config.Set("ORACLE_CONTRACT_ADDRESS", &oca)
	app, cleanup := cltest.NewApplicationWithConfigAndKeyStore(config)
	defer cleanup()

	txm := app.Store.TxManager

	from := cltest.GetAccountAddress(app.GetStore())
	to := cltest.NewAddress()
	hash := cltest.NewHash()
	sentAt := uint64(23456)
	nonce := uint64(256)
	ethMock := app.MockEthClient()
	ethMock.Context("app.StartAndConnect()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce))
	})
	assert.NoError(t, app.StartAndConnect())

	ethMock.Context("txm.CreateTx#1", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_sendRawTransaction", hash)
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(sentAt))
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
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()

	nonce := uint64(256)
	ethMock := app.MockEthClient()
	ethMock.Context("app.StartAndConnect()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce))
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
	assert.Equal(t, uint64(1), managedAccount.GetNonce())

	managedAccount.GetAndIncrementNonce(func(y uint64) error {
		assert.Equal(t, uint64(1), y)
		return nil
	})
	assert.Equal(t, uint64(2), managedAccount.GetNonce())
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
	assert.Equal(t, uint64(0), managedAccount.GetNonce())
}

func TestTxManager_LogsETHAndLINKBalancesAfterSuccessfulTx(t *testing.T) {
	t.Parallel()

	logsToCheckForBalance := cltest.ObserveLogs()

	config, configCleanup := cltest.NewConfig()
	defer configCleanup()
	oracleAddress := "0xDEADB3333333F"
	oca := common.HexToAddress(oracleAddress)
	config.Set("ORACLE_CONTRACT_ADDRESS", &oca)
	app, cleanup := cltest.NewApplicationWithConfigAndKeyStore(config)
	defer cleanup()

	manager := app.Store.TxManager
	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	assert.NoError(t, err)
	hash := cltest.NewHash()
	nonce := uint64(256)
	sentAt := uint64(23456)
	ethMock := app.MockEthClient()
	mockedEthBalance := "0x100"
	mockedLinkBalance := "256000000000000000000"
	confirmedHeight := sentAt + config.MinOutgoingConfirmations()
	confirmedReceipt := strpkg.TxReceipt{
		Hash:        hash,
		BlockNumber: cltest.Int(sentAt),
	}
	ethMock.Context("app.StartAndConnect()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_blockNumber", "0x1")
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce))
		ethMock.Register("eth_sendRawTransaction", hash)
		ethMock.Register("eth_getTransactionReceipt", confirmedReceipt)
		ethMock.Register("eth_blockNumber", utils.Uint64ToHex(
			confirmedHeight))
		ethMock.Register("eth_getBalance", mockedEthBalance)
		ethMock.Register("eth_call", mockedLinkBalance)
	})
	assert.NoError(t, app.StartAndConnect())

	confirmedTx, err := manager.CreateTx(to, data)
	assert.NoError(t, err)
	txTransmissionAttempts, err := app.Store.TxAttemptsFor(confirmedTx.ID)
	assert.NoError(t, err)
	initialSuccessfulAttempt := txTransmissionAttempts[0]

	txWasConfirmed, err := manager.EnsureConfirmed(
		initialSuccessfulAttempt.Hash)
	assert.NoError(t, err)
	assert.True(t, txWasConfirmed)

	ethMock.EventuallyAllCalled(t)

	targetLog := fmt.Sprintf("Confirmed tx %v", initialSuccessfulAttempt.Hash.String())
	targetLogSeen := false
	for _, log := range logsToCheckForBalance.All() {
		if strings.Contains(log.Entry.Message, targetLog) {
			targetLogSeen = true
			break
		}
	}
	assert.True(t, targetLogSeen)
}

func TestTxManager_CreateTxWithGas(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	store := app.Store
	manager := store.TxManager

	to := cltest.NewAddress()
	data, err := hex.DecodeString("0000abcdef")
	assert.NoError(t, err)
	nonce := uint64(256)
	ethMock := app.MockEthClient()
	ethMock.Context("app.Start()", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_getTransactionCount", utils.Uint64ToHex(nonce))
	})
	assert.NoError(t, app.StartAndConnect())

	customGasPrice := big.NewInt(1337)
	customGasLimit := uint64(10009)

	tests := []struct {
		name             string
		dev              bool
		gasPrice         *big.Int
		gasLimit         uint64
		expectedGasPrice *big.Int
		expectedGasLimit uint64
	}{
		{"dev", true, customGasPrice, customGasLimit, customGasPrice, customGasLimit},
		{"dev but not set", true, nil, 0, store.Config.EthGasPriceDefault(), strpkg.DefaultGasLimit},
		{"not dev", false, customGasPrice, customGasLimit, store.Config.EthGasPriceDefault(), strpkg.DefaultGasLimit},
		{"not dev not set", false, nil, 0, store.Config.EthGasPriceDefault(), strpkg.DefaultGasLimit},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			strpkg.ExportedSetTxManagerDev(manager, test.dev)
			ethMock.Context("manager.CreateTx", func(ethMock *cltest.EthMock) {
				ethMock.Register("eth_sendRawTransaction", cltest.NewHash())
				ethMock.Register("eth_blockNumber", utils.Uint64ToHex(1))
			})

			tx, err := manager.CreateTxWithGas(to, data, test.gasPrice, test.gasLimit)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedGasLimit, tx.GasLimit)

			attempts, err := store.TxAttemptsFor(tx.ID)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(attempts))
			assert.Equal(t, test.expectedGasPrice, attempts[0].GasPrice)

			ethMock.EventuallyAllCalled(t)
		})
	}
}
