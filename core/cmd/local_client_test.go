package cmd_test

import (
	"flag"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/dialects"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestClient_RunNodeShowsEnv(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)
	require.NoError(t, keyStore.Eth().Unlock(cltest.Password))
	_, err := keyStore.Eth().CreateNewKey()
	require.NoError(t, err)

	store.Config.Set("LINK_CONTRACT_ADDRESS", "0x514910771AF9Ca656af840dff83E8264EcF986CA")
	store.Config.Set("FLAGS_CONTRACT_ADDRESS", "0x4A5b9B4aD08616D11F3A402FF7cBEAcB732a76C6")
	store.Config.Set("CHAINLINK_PORT", 6688)

	ethClient := new(mocks.Client)
	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil)

	app := new(mocks.Application)
	app.On("GetStore").Return(store)
	app.On("GetKeyStore").Return(keyStore)
	app.On("GetEthClient").Return(ethClient).Maybe()
	app.On("Start").Return(nil)
	app.On("Stop").Return(nil)

	auth := cltest.CallbackAuthenticator{Callback: func(*keystore.Eth, string) (string, error) { return "", nil }}
	runner := cltest.BlockedRunner{Done: make(chan struct{})}
	client := cmd.Client{
		Config:                 store.Config,
		AppFactory:             cltest.InstanceAppFactory{App: app},
		KeyStoreAuthenticator:  auth,
		FallbackAPIInitializer: &cltest.MockAPIInitializer{},
		Runner:                 runner,
	}

	set := flag.NewFlagSet("test", 0)
	set.Bool("debug", true, "")
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
	logs, err := cltest.ReadLogs(store.Config)
	require.NoError(t, err)

	assert.Contains(t, logs, "ALLOW_ORIGINS: http://localhost:3000,http://localhost:6688\\n")
	assert.Contains(t, logs, "BRIDGE_RESPONSE_URL: http://localhost:6688\\n")
	assert.Contains(t, logs, "BLOCK_BACKFILL_DEPTH: 10\\n")
	assert.Contains(t, logs, "CHAINLINK_PORT: 6688\\n")
	assert.Contains(t, logs, "CLIENT_NODE_URL: http://")
	assert.Contains(t, logs, "ETH_CHAIN_ID: 0\\n")
	assert.Contains(t, logs, "ETH_GAS_BUMP_THRESHOLD: 3\\n")
	assert.Contains(t, logs, "ETH_GAS_BUMP_WEI: 5000000000\\n")
	assert.Contains(t, logs, "ETH_GAS_PRICE_DEFAULT: 20000000000\\n")
	assert.Contains(t, logs, "ETH_URL: ws://")
	assert.Contains(t, logs, "FLAGS_CONTRACT_ADDRESS: 0x4A5b9B4aD08616D11F3A402FF7cBEAcB732a76C6\\n")
	assert.Contains(t, logs, "JSON_CONSOLE: false")
	assert.Contains(t, logs, "LINK_CONTRACT_ADDRESS: 0x514910771AF9Ca656af840dff83E8264EcF986CA\\n")
	assert.Contains(t, logs, "LOG_LEVEL: debug\\n")
	assert.Contains(t, logs, "LOG_TO_DISK: true")
	assert.Contains(t, logs, "MIN_INCOMING_CONFIRMATIONS: 1\\n")
	assert.Contains(t, logs, "MIN_OUTGOING_CONFIRMATIONS: 6\\n")
	assert.Contains(t, logs, "MINIMUM_CONTRACT_PAYMENT_LINK_JUELS: 100\\n")
	assert.Contains(t, logs, "OPERATOR_CONTRACT_ADDRESS: \\n")
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
			app.On("GetEthClient").Return(new(mocks.Client)).Maybe()
			app.On("Start").Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)

			ethClient := new(mocks.Client)
			ethClient.On("Dial", mock.Anything).Return(nil)
			ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil)

			cltest.MustInsertRandomKey(t, store.DB)

			var unlocked bool
			callback := func(store *keystore.Eth, phrase string) (string, error) {
				err := keyStore.Eth().Unlock(phrase)
				unlocked = err == nil
				return phrase, err
			}

			auth := cltest.CallbackAuthenticator{Callback: callback}
			apiPrompt := &cltest.MockAPIInitializer{}
			client := cmd.Client{
				Config:                 store.Config,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				KeyStoreAuthenticator:  auth,
				FallbackAPIInitializer: apiPrompt,
				Runner:                 cltest.EmptyRunner{},
			}

			set := flag.NewFlagSet("test", 0)
			set.String("password", test.pwdfile, "")
			c := cli.NewContext(nil, set, nil)

			if test.wantUnlocked {
				assert.NoError(t, client.RunNode(c))
				assert.True(t, unlocked)
				assert.Equal(t, 1, apiPrompt.Count)
			} else {
				assert.Error(t, client.RunNode(c))
				assert.False(t, unlocked)
				assert.Equal(t, 0, apiPrompt.Count)
			}
		})
	}
}

func TestClient_RunNode_CreateFundingKeyIfNotExists(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	// Clear out fixture
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)
	require.NoError(t, keyStore.Eth().Unlock(cltest.Password))
	_, err := keyStore.Eth().CreateNewKey()
	require.NoError(t, err)

	app := new(mocks.Application)
	app.On("GetStore").Return(store)
	app.On("GetKeyStore").Return(keyStore)
	app.On("GetEthClient").Return(new(mocks.Client)).Maybe()
	app.On("Start").Maybe().Return(nil)
	app.On("Stop").Maybe().Return(nil)

	ethClient := new(mocks.Client)
	ethClient.On("Dial", mock.Anything).Return(nil)

	_, err = keyStore.Eth().CreateNewKey()
	require.NoError(t, err)

	callback := func(store *keystore.Eth, phrase string) (string, error) {
		unlockErr := keyStore.Eth().Unlock(phrase)
		return phrase, unlockErr
	}
	auth := cltest.CallbackAuthenticator{Callback: callback}
	apiPrompt := &cltest.MockAPIInitializer{}
	client := cmd.Client{
		Config:                 store.Config,
		AppFactory:             cltest.InstanceAppFactory{App: app},
		KeyStoreAuthenticator:  auth,
		FallbackAPIInitializer: apiPrompt,
		Runner:                 cltest.EmptyRunner{},
	}

	var fundingKey = ethkey.Key{}
	_ = store.DB.Where("is_funding = TRUE").First(&fundingKey).Error
	assert.Empty(t, fundingKey.ID, "expected no funding key")

	set := flag.NewFlagSet("test", 0)
	set.String("password", "../internal/fixtures/correct_password.txt", "")
	ctx := cli.NewContext(nil, set, nil)

	assert.NoError(t, client.RunNode(ctx))

	assert.NoError(t, store.DB.Where("is_funding = TRUE").First(&fundingKey).Error)
	assert.NotEmpty(t, fundingKey.ID, "expected a new funding key")
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
			cfg := config.NewConfig()

			store, cleanup := cltest.NewStore(t)
			// Clear out fixture
			store.DeleteUser()
			defer cleanup()
			keyStore := cltest.NewKeyStore(t, store.DB)
			require.NoError(t, keyStore.Eth().Unlock(cltest.Password))
			_, err := keyStore.Eth().CreateNewKey()
			require.NoError(t, err)

			app := new(mocks.Application)
			app.On("GetStore").Return(store)
			app.On("GetKeyStore").Return(keyStore)
			app.On("GetEthClient").Return(new(mocks.Client)).Maybe()
			app.On("Start").Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)

			ethClient := new(mocks.Client)
			ethClient.On("Dial", mock.Anything).Return(nil)
			ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil)

			callback := func(*keystore.Eth, string) (string, error) { return "", nil }
			noauth := cltest.CallbackAuthenticator{Callback: callback}
			apiPrompt := &cltest.MockAPIInitializer{}
			client := cmd.Client{
				Config:                 cfg,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				KeyStoreAuthenticator:  noauth,
				FallbackAPIInitializer: apiPrompt,
				Runner:                 cltest.EmptyRunner{},
			}

			set := flag.NewFlagSet("test", 0)
			set.String("api", test.apiFile, "")
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

func TestClient_ImportKey(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	kst := cltest.NewKeyStore(t, store.DB).Eth()

	ethClient, _, assertMocksCalled := cltest.NewEthMocks(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplication(t, ethClient, kst)
	defer cleanup()

	client, _ := app.NewClientAndRenderer()

	path := "../internal/fixtures/keys/7fc66c61f88A61DFB670627cA715Fe808057123e.json"

	set := flag.NewFlagSet("import", 0)
	set.Parse([]string{path})
	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.ImportKey(c))
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
			config, configCleanup := cltest.NewConfig(t)
			defer configCleanup()
			config.Set("CHAINLINK_DEV", true)
			config.Set("LOG_TO_DISK", tt.logToDiskValue)
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
	config.Config.Dialect = dialects.PostgresWithoutLock
	connectedStore, connectedCleanup := cltest.NewStoreWithConfig(t, config)
	defer connectedCleanup()
	keyStore := cltest.NewKeyStore(t, connectedStore.DB)
	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, keyStore.Eth(), 0)

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
	c := cli.NewContext(nil, set, nil)

	cltest.MustInsertConfirmedEthTxWithAttempt(t, connectedStore.DB, 7, 42, fromAddress)

	// Use the same config as the connectedStore so that the advisory
	// lock ID is the same. We set the config to be Postgres Without
	// Lock, because the db locking strategy is decided when we
	// initialize the store/ORM.
	config.Config.Dialect = dialects.PostgresWithoutLock
	store, cleanup := cltest.NewStoreWithConfig(t, config)
	defer cleanup()
	keyStore.Eth().Unlock(cltest.Password)
	require.NoError(t, connectedStore.Start())

	app := new(mocks.Application)
	app.On("GetStore").Return(store)
	app.On("GetKeyStore").Return(keyStore)
	app.On("Stop").Return(nil)
	ethClient := new(mocks.Client)
	app.On("GetEthClient").Return(ethClient).Maybe()
	ethClient.On("Dial", mock.Anything).Return(nil)

	auth := cltest.CallbackAuthenticator{Callback: func(*keystore.Eth, string) (string, error) { return "", nil }}
	client := cmd.Client{
		Config:                 config.Config,
		AppFactory:             cltest.InstanceAppFactory{App: app},
		KeyStoreAuthenticator:  auth,
		FallbackAPIInitializer: &cltest.MockAPIInitializer{},
		Runner:                 cltest.EmptyRunner{},
	}

	config.Config.Dialect = dialects.TransactionWrappedPostgres

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
	assert.Equal(t, dialects.PostgresWithoutLock, config.Config.GetDatabaseDialectConfiguredOrDefault())

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
			config.Config.Dialect = dialects.Postgres
			connectedStore, connectedCleanup := cltest.NewStoreWithConfig(t, config)
			defer connectedCleanup()
			keyStore := cltest.NewKeyStore(t, connectedStore.DB)

			_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, keyStore.Eth(), 0)

			set := flag.NewFlagSet("test", 0)
			set.Bool("debug", true, "")
			set.Uint("beginningNonce", beginningNonce, "")
			set.Uint("endingNonce", endingNonce, "")
			set.Uint64("gasPriceWei", gasPrice.Uint64(), "")
			set.Uint64("gasLimit", gasLimit, "")
			set.String("address", fromAddress.Hex(), "")
			c := cli.NewContext(nil, set, nil)

			cltest.MustInsertConfirmedEthTxWithAttempt(t, connectedStore.DB, int64(test.nonce), 42, fromAddress)

			// Use the same config as the connectedStore so that the advisory
			// lock ID is the same. We set the config to be Postgres Without
			// Lock, because the db locking strategy is decided when we
			// initialize the store/ORM.
			config.Config.Dialect = dialects.PostgresWithoutLock
			store, cleanup := cltest.NewStoreWithConfig(t, config)
			defer cleanup()
			keyStore.Eth().Unlock(cltest.Password)
			require.NoError(t, connectedStore.Start())

			app := new(mocks.Application)
			app.On("GetStore").Return(store)
			app.On("GetKeyStore").Return(keyStore)
			app.On("Stop").Return(nil)
			ethClient := new(mocks.Client)
			ethClient.On("Dial", mock.Anything).Return(nil)
			app.On("GetEthClient").Return(ethClient).Maybe()

			auth := cltest.CallbackAuthenticator{Callback: func(*keystore.Eth, string) (string, error) { return "", nil }}
			client := cmd.Client{
				Config:                 config.Config,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				KeyStoreAuthenticator:  auth,
				FallbackAPIInitializer: &cltest.MockAPIInitializer{},
				Runner:                 cltest.EmptyRunner{},
			}

			config.Config.Dialect = dialects.TransactionWrappedPostgres

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
			assert.Equal(t, dialects.PostgresWithoutLock, config.Config.GetDatabaseDialectConfiguredOrDefault())

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
	config.Config.Dialect = dialects.Postgres
	store, cleanup := cltest.NewStoreWithConfig(t, config)
	defer cleanup()
	ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth()

	client := cmd.Client{
		Config: config.Config,
		Runner: cltest.EmptyRunner{},
	}

	set := flag.NewFlagSet("test", 0)
	set.Bool("debug", true, "")
	set.Uint("nextNonce", 42, "")
	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore, 0)
	set.String("address", fromAddress.Hex(), "")
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.SetNextNonce(c))

	var key ethkey.Key
	require.NoError(t, store.DB.First(&key).Error)
	require.NotNil(t, key.NextNonce)
	require.Equal(t, int64(42), key.NextNonce)
}
