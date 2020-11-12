package cmd_test

import (
	"flag"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestClient_RunNodeShowsEnv(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	_, err := store.KeyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	require.NoError(t, store.KeyStore.Unlock(cltest.Password))

	store.Config.Set("LINK_CONTRACT_ADDRESS", "0x514910771AF9Ca656af840dff83E8264EcF986CA")
	store.Config.Set("FLAGS_CONTRACT_ADDRESS", "0x4A5b9B4aD08616D11F3A402FF7cBEAcB732a76C6")
	store.Config.Set("CHAINLINK_PORT", 6688)

	ethClient := new(mocks.Client)
	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil)
	store.EthClient = ethClient

	app := new(mocks.Application)
	app.On("GetStore").Return(store)
	app.On("Start").Return(nil)
	app.On("Stop").Return(nil)

	auth := cltest.CallbackAuthenticator{Callback: func(*strpkg.Store, string) (string, error) { return "", nil }}
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
	assert.Contains(t, logs, "ETH_CHAIN_ID: 3\\n")
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
	assert.Contains(t, logs, "MINIMUM_CONTRACT_PAYMENT: 0.000000000000000100\\n")
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
			// Clear out fixture
			store.DeleteUser()
			defer cleanup()
			_, err := store.KeyStore.NewAccount(cltest.Password)
			require.NoError(t, err)
			require.NoError(t, store.KeyStore.Unlock(cltest.Password))

			app := new(mocks.Application)
			app.On("GetStore").Return(store)
			app.On("Start").Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)

			ethClient := new(mocks.Client)
			ethClient.On("Dial", mock.Anything).Return(nil)
			ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil)
			store.EthClient = ethClient

			_, err = store.KeyStore.NewAccount("password") // matches correct_password.txt
			require.NoError(t, err)

			var unlocked bool
			callback := func(store *strpkg.Store, phrase string) (string, error) {
				err := store.KeyStore.Unlock(phrase)
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
	_, err := store.KeyStore.NewAccount(cltest.Password)
	require.NoError(t, err)
	require.NoError(t, store.KeyStore.Unlock(cltest.Password))

	app := new(mocks.Application)
	app.On("GetStore").Return(store)
	app.On("Start").Maybe().Return(nil)
	app.On("Stop").Maybe().Return(nil)

	ethClient := new(mocks.Client)
	ethClient.On("Dial", mock.Anything).Return(nil)
	store.EthClient = ethClient

	_, err = store.KeyStore.NewAccount("password") // matches correct_password.txt
	require.NoError(t, err)

	callback := func(store *strpkg.Store, phrase string) (string, error) {
		unlockErr := store.KeyStore.Unlock(phrase)
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

	var fundingKey = models.Key{}
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
			config := orm.NewConfig()

			store, cleanup := cltest.NewStore(t)
			// Clear out fixture
			store.DeleteUser()
			defer cleanup()
			_, err := store.KeyStore.NewAccount(cltest.Password)
			require.NoError(t, err)
			require.NoError(t, store.KeyStore.Unlock(cltest.Password))

			app := new(mocks.Application)
			app.On("GetStore").Return(store)
			app.On("Start").Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)

			ethClient := new(mocks.Client)
			ethClient.On("Dial", mock.Anything).Return(nil)
			ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil)
			store.EthClient = ethClient

			callback := func(*strpkg.Store, string) (string, error) { return "", nil }
			noauth := cltest.CallbackAuthenticator{Callback: callback}
			apiPrompt := &cltest.MockAPIInitializer{}
			client := cmd.Client{
				Config:                 config,
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

	app, cleanup := cltest.NewApplication(t,
		cltest.AllowUnstarted,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()

	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("import", 0)
	set.Parse([]string{"../internal/fixtures/keys/7fc66c61f88A61DFB670627cA715Fe808057123e.json"})
	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.ImportKey(c))

	// importing again simply upserts
	require.NoError(t, client.ImportKey(c))

	keys, err := app.GetStore().SendKeys()
	require.NoError(t, err)

	require.Len(t, keys, 2)
	require.Equal(t, int32(1), keys[0].ID)
	require.Greater(t, keys[1].ID, int32(1))

	addresses := []string{}
	for _, k := range keys {
		addresses = append(addresses, k.Address.String())
	}

	sort.Strings(addresses)
	expectation := []string{cltest.DefaultKey, "0x7fc66c61f88A61DFB670627cA715Fe808057123e"}
	require.Equal(t, expectation, addresses)
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
			require.NoError(t, os.MkdirAll(config.KeysDir(), os.FileMode(0700)))
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
	set.String("address", cltest.DefaultKey, "")
	c := cli.NewContext(nil, set, nil)

	// Use the a non-transactional db for this test because we need to
	// test multiple connections to the database, and changes made within
	// the transaction cannot be seen from another connection.
	config, _, cleanup := cltest.BootstrapThrowawayORM(t, "rebroadcasttransactions", true, true)
	defer cleanup()
	config.Config.Dialect = orm.DialectPostgres
	connectedStore, connectedCleanup := cltest.NewStoreWithConfig(config)
	defer connectedCleanup()

	cltest.MustInsertConfirmedEthTxWithAttempt(t, connectedStore, 7, 42)

	// Use the same config as the connectedStore so that the advisory
	// lock ID is the same. We set the config to be Postgres Without
	// Lock, because the db locking strategy is decided when we
	// initialize the store/ORM.
	config.Config.Dialect = orm.DialectPostgresWithoutLock
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()
	store.KeyStore.Unlock(cltest.Password)
	require.NoError(t, connectedStore.Start())

	app := new(mocks.Application)
	app.On("GetStore").Return(store)
	app.On("Stop").Return(nil)
	ethClient := new(mocks.Client)
	ethClient.On("Dial", mock.Anything).Return(nil)
	store.EthClient = ethClient

	auth := cltest.CallbackAuthenticator{Callback: func(*strpkg.Store, string) (string, error) { return "", nil }}
	client := cmd.Client{
		Config:                 config.Config,
		AppFactory:             cltest.InstanceAppFactory{App: app},
		KeyStoreAuthenticator:  auth,
		FallbackAPIInitializer: &cltest.MockAPIInitializer{},
		Runner:                 cltest.EmptyRunner{},
	}

	config.Config.Dialect = orm.DialectTransactionWrappedPostgres

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
	assert.Equal(t, orm.DialectPostgresWithoutLock, config.Config.GetDatabaseDialectConfiguredOrDefault())

	app.AssertExpectations(t)
	ethClient.AssertExpectations(t)
}

func TestClient_RebroadcastTransactions_OutsideRange_BPTXM(t *testing.T) {
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
	set.String("address", cltest.DefaultKey, "")
	c := cli.NewContext(nil, set, nil)

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
			config, _, cleanup := cltest.BootstrapThrowawayORM(t, "rebroadcasttransactions_outsiderange", true, true)
			defer cleanup()
			config.Config.Dialect = orm.DialectPostgres
			connectedStore, connectedCleanup := cltest.NewStoreWithConfig(config)
			defer connectedCleanup()

			cltest.MustInsertConfirmedEthTxWithAttempt(t, connectedStore, int64(test.nonce), 42)

			// Use the same config as the connectedStore so that the advisory
			// lock ID is the same. We set the config to be Postgres Without
			// Lock, because the db locking strategy is decided when we
			// initialize the store/ORM.
			config.Config.Dialect = orm.DialectPostgresWithoutLock
			store, cleanup := cltest.NewStoreWithConfig(config)
			defer cleanup()
			store.KeyStore.Unlock(cltest.Password)
			require.NoError(t, connectedStore.Start())

			app := new(mocks.Application)
			app.On("GetStore").Return(store)
			app.On("Stop").Return(nil)
			ethClient := new(mocks.Client)
			ethClient.On("Dial", mock.Anything).Return(nil)
			store.EthClient = ethClient

			auth := cltest.CallbackAuthenticator{Callback: func(*strpkg.Store, string) (string, error) { return "", nil }}
			client := cmd.Client{
				Config:                 config.Config,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				KeyStoreAuthenticator:  auth,
				FallbackAPIInitializer: &cltest.MockAPIInitializer{},
				Runner:                 cltest.EmptyRunner{},
			}

			config.Config.Dialect = orm.DialectTransactionWrappedPostgres

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
			assert.Equal(t, orm.DialectPostgresWithoutLock, config.Config.GetDatabaseDialectConfiguredOrDefault())

			cltest.AssertEthTxAttemptCountStays(t, store, 1)
			app.AssertExpectations(t)
			ethClient.AssertExpectations(t)
		})
	}
}

func TestClient_SetNextNonce(t *testing.T) {
	// Need to use separate database
	config, _, cleanup := cltest.BootstrapThrowawayORM(t, "setnextnonce", true, true)
	defer cleanup()
	config.Config.Dialect = orm.DialectPostgres
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	client := cmd.Client{
		Config: config.Config,
		Runner: cltest.EmptyRunner{},
	}

	set := flag.NewFlagSet("test", 0)
	set.Bool("debug", true, "")
	set.Uint("nextNonce", 42, "")
	defaultFromAddress := cltest.GetDefaultFromAddress(t, store)
	set.String("address", defaultFromAddress.Hex(), "")
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.SetNextNonce(c))

	var key models.Key
	require.NoError(t, store.DB.First(&key).Error)
	require.NotNil(t, key.NextNonce)
	require.Equal(t, int64(42), *key.NextNonce)
}
