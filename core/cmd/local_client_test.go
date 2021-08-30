package cmd_test

import (
	"flag"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/cmd"
	cmdMocks "github.com/smartcontractkit/chainlink/core/cmd/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	null "gopkg.in/guregu/null.v4"
)

func TestClient_RunNodeShowsEnv(t *testing.T) {
	cfg := cltest.NewTestEVMConfig(t)
	store, cleanup := cltest.NewStoreWithConfig(t, cfg)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)
	_, err := keyStore.Eth().Create()
	require.NoError(t, err)

	ethClient := cltest.NewEthClientMock(t)
	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil)

	app := new(mocks.Application)
	app.On("GetStore").Return(store)
	app.On("GetKeyStore").Return(keyStore)
	app.On("GetEthClient").Return(ethClient).Maybe()
	app.On("Start").Return(nil)
	app.On("Stop").Return(nil)

	runner := cltest.BlockedRunner{Done: make(chan struct{})}
	client := cmd.Client{
		Config:                 store.Config,
		AppFactory:             cltest.InstanceAppFactory{App: app},
		FallbackAPIInitializer: &cltest.MockAPIInitializer{},
		Runner:                 runner,
	}

	set := flag.NewFlagSet("test", 0)
	set.Bool("debug", true, "")
	set.String("password", "../internal/fixtures/correct_password.txt", "")
	c := cli.NewContext(nil, set, nil)

	// Start RunNode in a goroutine, it will block until we resume the runner
	go func() {
		assert.NoError(t, client.RunNode(c))
	}()

	// Unlock the runner to the client can begin shutdown
	select {
	case runner.Done <- struct{}{}:
	case <-time.After(30 * time.Second):
		t.Fatal("Timed out waiting for runner")
	}

	logger.Sync()
	logs, err := cltest.ReadLogs(cfg)
	require.NoError(t, err)

	assert.Contains(t, logs, "ALLOW_ORIGINS: http://localhost:3000,http://localhost:6688\\n")
	assert.Contains(t, logs, "BRIDGE_RESPONSE_URL: http://localhost:6688\\n")
	assert.Contains(t, logs, "BLOCK_BACKFILL_DEPTH: 10\\n")
	assert.Contains(t, logs, "CHAINLINK_PORT: 6688\\n")
	assert.Contains(t, logs, "CLIENT_NODE_URL: http://")
	assert.Contains(t, logs, "ETH_CHAIN_ID: 0\\n")
	assert.Contains(t, logs, "ETH_URL: ws://")
	assert.Contains(t, logs, "JSON_CONSOLE: false")
	assert.Contains(t, logs, "LOG_LEVEL: debug\\n")
	assert.Contains(t, logs, "LOG_TO_DISK: true")
	assert.Contains(t, logs, "ROOT: /tmp/chainlink_test/")

	app.AssertExpectations(t)
}

func TestClient_RunNodeWithPasswords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		pwdfile      string
		wantUnlocked bool
	}{
		{"correct", "../internal/fixtures/correct_password.txt", true},
		{"incorrect", "../internal/fixtures/incorrect_password.txt", false},
		{"wrongfile", "doesntexist.txt", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()
			keyStore := cltest.NewKeyStore(t, store.DB)
			// Clear out fixture
			err := store.DeleteUser()
			require.NoError(t, err)

			app := new(mocks.Application)
			app.On("GetStore").Return(store)
			app.On("GetKeyStore").Return(keyStore)
			app.On("GetEthClient").Return(cltest.NewEthClientMock(t)).Maybe()
			app.On("Start").Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)

			ethClient := cltest.NewEthClientMock(t)
			ethClient.On("Dial", mock.Anything).Return(nil)
			ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil)

			cltest.MustInsertRandomKey(t, keyStore.Eth())

			apiPrompt := &cltest.MockAPIInitializer{}
			client := cmd.Client{
				Config:                 store.Config,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				FallbackAPIInitializer: apiPrompt,
				Runner:                 cltest.EmptyRunner{},
			}

			set := flag.NewFlagSet("test", 0)
			set.String("password", test.pwdfile, "")
			c := cli.NewContext(nil, set, nil)

			if test.wantUnlocked {
				assert.NoError(t, client.RunNode(c))
				assert.Equal(t, 1, apiPrompt.Count)
			} else {
				assert.Error(t, client.RunNode(c))
				assert.Equal(t, 0, apiPrompt.Count)
			}
		})
	}
}

func TestClient_RunNode_CreateFundingKeyIfNotExists(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)
	_, err := keyStore.Eth().Create()
	require.NoError(t, err)

	app := new(mocks.Application)
	app.On("GetStore").Return(store)
	app.On("GetKeyStore").Return(keyStore)
	app.On("GetEthClient").Return(cltest.NewEthClientMock(t)).Maybe()
	app.On("Start").Maybe().Return(nil)
	app.On("Stop").Maybe().Return(nil)

	ethClient := cltest.NewEthClientMock(t)
	ethClient.On("Dial", mock.Anything).Return(nil)

	_, err = keyStore.Eth().Create()
	require.NoError(t, err)

	apiPrompt := &cltest.MockAPIInitializer{}
	client := cmd.Client{
		Config:                 store.Config,
		AppFactory:             cltest.InstanceAppFactory{App: app},
		FallbackAPIInitializer: apiPrompt,
		Runner:                 cltest.EmptyRunner{},
	}

	var keyState = ethkey.State{}
	err = store.DB.Where("is_funding = TRUE").Find(&keyState).Error
	require.NoError(t, err)
	assert.Empty(t, keyState.ID, "expected no funding key")

	set := flag.NewFlagSet("test", 0)
	set.String("password", "../internal/fixtures/correct_password.txt", "")
	ctx := cli.NewContext(nil, set, nil)

	assert.NoError(t, client.RunNode(ctx))

	assert.NoError(t, store.DB.Where("is_funding = TRUE").First(&keyState).Error)
	assert.NotEmpty(t, keyState.ID, "expected a new funding key")
}

func TestClient_RunNodeWithAPICredentialsFile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		apiFile    string
		wantPrompt bool
		wantError  bool
	}{
		{"correct", "../internal/fixtures/apicredentials", false, false},
		{"no file", "", true, false},
		{"wrong file", "doesntexist.txt", false, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := cltest.NewTestEVMConfig(t)

			store, cleanup := cltest.NewStoreWithConfig(t, cfg)
			// Clear out fixture
			store.DeleteUser()
			defer cleanup()
			keyStore := cltest.NewKeyStore(t, store.DB)
			_, err := keyStore.Eth().Create()
			require.NoError(t, err)

			app := new(mocks.Application)
			app.On("GetStore").Return(store)
			app.On("GetKeyStore").Return(keyStore)
			app.On("GetEthClient").Return(cltest.NewEthClientMock(t)).Maybe()
			app.On("Start").Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)

			ethClient := cltest.NewEthClientMock(t)
			ethClient.On("Dial", mock.Anything).Return(nil)
			ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil)

			prompter := new(cmdMocks.Prompter)
			prompter.On("IsTerminal").Return(false).Once().Maybe()

			apiPrompt := &cltest.MockAPIInitializer{}
			client := cmd.Client{
				Config:                 cfg,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				KeyStoreAuthenticator:  cmd.TerminalKeyStoreAuthenticator{prompter},
				FallbackAPIInitializer: apiPrompt,
				Runner:                 cltest.EmptyRunner{},
			}

			set := flag.NewFlagSet("test", 0)
			set.String("api", test.apiFile, "")
			set.String("password", "../internal/fixtures/correct_password.txt", "")
			c := cli.NewContext(nil, set, nil)

			if test.wantError {
				assert.EqualError(t, client.RunNode(c), "error creating api initializer: open doesntexist.txt: no such file or directory")
			} else {
				assert.NoError(t, client.RunNode(c))
			}
			assert.Equal(t, test.wantPrompt, apiPrompt.Count > 0)

			app.AssertExpectations(t)
		})
	}
}

func TestClient_LogToDiskOptionDisablesAsExpected(t *testing.T) {
	tests := []struct {
		name            string
		logToDiskValue  bool
		fileShouldExist bool
	}{
		{"LogToDisk = false => no log on disk", false, false},
		{"LogToDisk = true => log on disk (positive control)", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := cltest.NewTestEVMConfig(t)
			config.GeneralConfig.Overrides.Dev = null.BoolFrom(true)
			config.GeneralConfig.Overrides.LogToDisk = null.BoolFrom(tt.logToDiskValue)
			require.NoError(t, os.MkdirAll(config.RootDir(), os.FileMode(0700)))
			defer os.RemoveAll(config.RootDir())

			previousLogger := logger.Default
			logger.SetLogger(config.CreateProductionLogger())
			defer logger.SetLogger(previousLogger)
			filepath := filepath.Join(config.RootDir(), "log.jsonl")
			_, err := os.Stat(filepath)
			assert.Equal(t, os.IsNotExist(err), !tt.fileShouldExist)
		})
	}
}

func TestClient_RebroadcastTransactions_BPTXM(t *testing.T) {
	// Use the a non-transactional db for this test because we need to
	// test multiple connections to the database, and changes made within
	// the transaction cannot be seen from another connection.
	config, _, cleanup := heavyweight.FullTestORM(t, "rebroadcasttransactions", true, true)
	defer cleanup()
	connectedStore, connectedCleanup := cltest.NewStoreWithConfig(t, config)
	defer connectedCleanup()
	keyStore := cltest.NewKeyStore(t, connectedStore.DB)
	_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth(), 0)

	beginningNonce := uint(7)
	endingNonce := uint(10)
	gasPrice := big.NewInt(100000000000)
	gasLimit := uint64(3000000)
	set := flag.NewFlagSet("test", 0)
	set.Bool("debug", true, "")
	set.Uint("beginningNonce", beginningNonce, "")
	set.Uint("endingNonce", endingNonce, "")
	set.Uint64("gasPriceWei", gasPrice.Uint64(), "")
	set.Uint64("gasLimit", gasLimit, "")
	set.String("address", fromAddress.Hex(), "")
	set.String("password", "../internal/fixtures/correct_password.txt", "")
	c := cli.NewContext(nil, set, nil)

	cltest.MustInsertConfirmedEthTxWithAttempt(t, connectedStore.DB, 7, 42, fromAddress)

	// Use the same config as the connectedStore so that the advisory
	// lock ID is the same. We set the config to be Postgres Without
	// Lock, because the db locking strategy is decided when we
	// initialize the store/ORM.

	config.SetDialect(dialects.PostgresWithoutLock)

	store, cleanup := cltest.NewStoreWithConfig(t, config)
	defer cleanup()
	require.NoError(t, connectedStore.Start())

	app := new(mocks.Application)
	app.On("GetStore").Return(store)
	app.On("GetKeyStore").Return(keyStore)
	app.On("Stop").Return(nil)
	ethClient := cltest.NewEthClientMock(t)
	app.On("GetEthClient").Return(ethClient).Maybe()
	ethClient.On("Dial", mock.Anything).Return(nil)

	client := cmd.Client{
		Config:                 config,
		AppFactory:             cltest.InstanceAppFactory{App: app},
		FallbackAPIInitializer: &cltest.MockAPIInitializer{},
		Runner:                 cltest.EmptyRunner{},
	}

	config.SetDialect(dialects.TransactionWrappedPostgres)

	for i := beginningNonce; i <= endingNonce; i++ {
		n := i
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return uint(tx.Nonce()) == n
		})).Once().Return(nil)
	}

	// We set the dialect back after initialization so that we can check
	// that it was set back to WithoutLock at the end of the test.
	assert.NoError(t, client.RebroadcastTransactions(c))
	// Check that the Dialect was set back when the command was run.
	assert.Equal(t, dialects.PostgresWithoutLock, config.GetDatabaseDialectConfiguredOrDefault())

	app.AssertExpectations(t)
	ethClient.AssertExpectations(t)
}

func TestClient_RebroadcastTransactions_OutsideRange_BPTXM(t *testing.T) {
	beginningNonce := uint(7)
	endingNonce := uint(10)
	gasPrice := big.NewInt(100000000000)
	gasLimit := uint64(3000000)

	tests := []struct {
		name  string
		nonce uint
	}{
		{"below beginning", beginningNonce - 1},
		{"above ending", endingNonce + 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use the a non-transactional db for this test because we need to
			// test multiple connections to the database, and changes made within
			// the transaction cannot be seen from another connection.
			config, _, cleanup := heavyweight.FullTestORM(t, "rebroadcasttransactions_outsiderange", true, true)
			defer cleanup()
			config.SetDialect(dialects.Postgres)
			connectedStore, connectedCleanup := cltest.NewStoreWithConfig(t, config)
			defer connectedCleanup()
			keyStore := cltest.NewKeyStore(t, connectedStore.DB)

			_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth(), 0)

			set := flag.NewFlagSet("test", 0)
			set.Bool("debug", true, "")
			set.Uint("beginningNonce", beginningNonce, "")
			set.Uint("endingNonce", endingNonce, "")
			set.Uint64("gasPriceWei", gasPrice.Uint64(), "")
			set.Uint64("gasLimit", gasLimit, "")
			set.String("address", fromAddress.Hex(), "")
			set.String("password", "../internal/fixtures/correct_password.txt", "")
			c := cli.NewContext(nil, set, nil)

			cltest.MustInsertConfirmedEthTxWithAttempt(t, connectedStore.DB, int64(test.nonce), 42, fromAddress)

			// Use the same config as the connectedStore so that the advisory
			// lock ID is the same. We set the config to be Postgres Without
			// Lock, because the db locking strategy is decided when we
			// initialize the store/ORM.
			config.SetDialect(dialects.PostgresWithoutLock)
			store, cleanup := cltest.NewStoreWithConfig(t, config)
			defer cleanup()
			require.NoError(t, connectedStore.Start())

			app := new(mocks.Application)
			app.On("GetStore").Return(store)
			app.On("GetKeyStore").Return(keyStore)
			app.On("Stop").Return(nil)
			ethClient := cltest.NewEthClientMock(t)
			ethClient.On("Dial", mock.Anything).Return(nil)
			app.On("GetEthClient").Return(ethClient).Maybe()

			client := cmd.Client{
				Config:                 config,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				FallbackAPIInitializer: &cltest.MockAPIInitializer{},
				Runner:                 cltest.EmptyRunner{},
			}

			config.SetDialect(dialects.TransactionWrappedPostgres)

			for i := beginningNonce; i <= endingNonce; i++ {
				n := i
				ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
					return uint(tx.Nonce()) == n
				})).Once().Return(nil)
			}

			// We set the dialect back after initialization so that we can check
			// that it was set back to WithoutLock at the end of the test.
			assert.NoError(t, client.RebroadcastTransactions(c))
			// Check that the Dialect was set back when the command was run.
			assert.Equal(t, dialects.PostgresWithoutLock, config.GetDatabaseDialectConfiguredOrDefault())

			cltest.AssertEthTxAttemptCountStays(t, store, 1)
			app.AssertExpectations(t)
			ethClient.AssertExpectations(t)
		})
	}
}

func TestClient_SetNextNonce(t *testing.T) {
	// Need to use separate database
	config, _, cleanup := heavyweight.FullTestORM(t, "setnextnonce", true, true)
	defer cleanup()
	config.SetDialect(dialects.Postgres)
	store, cleanup := cltest.NewStoreWithConfig(t, config)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	client := cmd.Client{
		Config: config,
		Runner: cltest.EmptyRunner{},
	}

	_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)

	set := flag.NewFlagSet("test", 0)
	set.Bool("debug", true, "")
	set.Uint("nextNonce", 42, "")
	set.String("address", fromAddress.Hex(), "")
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.SetNextNonce(c))

	var state ethkey.State
	require.NoError(t, store.DB.First(&state).Error)
	require.NotNil(t, state.NextNonce)
	require.Equal(t, int64(42), state.NextNonce)
}
