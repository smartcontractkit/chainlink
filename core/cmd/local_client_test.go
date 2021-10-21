package cmd_test

import (
	"flag"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/cmd"
	cmdMocks "github.com/smartcontractkit/chainlink/core/cmd/mocks"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/store/dialects"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"go.uber.org/zap/zapcore"
	null "gopkg.in/guregu/null.v4"
)

func TestClient_RunNodeShowsEnv(t *testing.T) {
	cfg := cltest.NewTestGeneralConfig(t)
	debug := config.LogLevel{Level: zapcore.DebugLevel}
	cfg.Overrides.LogLevel = &debug
	cfg.Overrides.LogToDisk = null.BoolFrom(true)
	db := pgtest.NewGormDB(t)
	sessionORM := sessions.NewORM(postgres.UnwrapGormDB(db), time.Minute)
	keyStore := cltest.NewKeyStore(t, db)
	_, err := keyStore.Eth().Create(&cltest.FixtureChainID)
	require.NoError(t, err)

	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	ethClient.On("Dial", mock.Anything).Return(nil)
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil)

	app := new(mocks.Application)
	app.On("SessionORM").Return(sessionORM)
	app.On("GetKeyStore").Return(keyStore)
	app.On("GetChainSet").Return(cltest.NewChainSetMockWithOneChain(t, ethClient, evmtest.NewChainScopedConfig(t, cfg))).Maybe()
	app.On("Start").Return(nil)
	app.On("Stop").Return(nil)

	runner := cltest.BlockedRunner{Done: make(chan struct{})}
	client := cmd.Client{
		Config:                 cfg,
		Logger:                 logger.TestLogger(t),
		AppFactory:             cltest.InstanceAppFactory{App: app},
		FallbackAPIInitializer: cltest.NewMockAPIInitializer(t),
		Runner:                 runner,
	}

	// Start RunNode in a goroutine, it will block until we resume the runner
	go func() {
		assert.NoError(t, cmd.NewApp(&client).
			Run([]string{"", "node", "start", "-debug", "-password", "../internal/fixtures/correct_password.txt"}))
	}()

	// Unlock the runner to the client can begin shutdown
	select {
	case runner.Done <- struct{}{}:
	case <-time.After(30 * time.Second):
		t.Fatal("Timed out waiting for runner")
	}

	logs, err := cltest.ReadLogs(cfg)
	require.NoError(t, err)

	assert.Contains(t, logs, "ALLOW_ORIGINS: http://localhost:3000,http://localhost:6688\\n")
	assert.Contains(t, logs, "BRIDGE_RESPONSE_URL: http://localhost:6688\\n")
	assert.Contains(t, logs, "BLOCK_BACKFILL_DEPTH: 10\\n")
	assert.Contains(t, logs, "CHAINLINK_PORT: 6688\\n")
	assert.Contains(t, logs, "CLIENT_NODE_URL: http://")
	assert.Contains(t, logs, "ETH_CHAIN_ID: 0\\n")
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
			cfg := cltest.NewTestGeneralConfig(t)
			db := pgtest.NewGormDB(t)
			keyStore := cltest.NewKeyStore(t, db)
			sessionORM := sessions.NewORM(postgres.UnwrapGormDB(db), time.Minute)
			// Clear out fixture
			err := sessionORM.DeleteUser()
			require.NoError(t, err)

			app := new(mocks.Application)
			app.On("SessionORM").Return(sessionORM)
			app.On("GetKeyStore").Return(keyStore)
			app.On("GetChainSet").Return(cltest.NewChainSetMockWithOneChain(t, cltest.NewEthClientMock(t), evmtest.NewChainScopedConfig(t, cfg))).Maybe()
			app.On("Start").Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)

			ethClient := cltest.NewEthClientMock(t)
			ethClient.On("Dial", mock.Anything).Return(nil)
			ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil)

			cltest.MustInsertRandomKey(t, keyStore.Eth())

			apiPrompt := cltest.NewMockAPIInitializer(t)
			client := cmd.Client{
				Config:                 cfg,
				Logger:                 logger.TestLogger(t),
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

	cfg := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewGormDB(t)
	sessionORM := sessions.NewORM(postgres.UnwrapGormDB(db), time.Minute)
	keyStore := cltest.NewKeyStore(t, db)
	_, err := keyStore.Eth().Create(&cltest.FixtureChainID)
	require.NoError(t, err)

	app := new(mocks.Application)
	app.On("SessionORM").Return(sessionORM)
	app.On("GetKeyStore").Return(keyStore)
	app.On("GetChainSet").Return(cltest.NewChainSetMockWithOneChain(t, cltest.NewEthClientMock(t), evmtest.NewChainScopedConfig(t, cfg))).Maybe()
	app.On("Start").Maybe().Return(nil)
	app.On("Stop").Maybe().Return(nil)

	ethClient := cltest.NewEthClientMock(t)
	ethClient.On("Dial", mock.Anything).Return(nil)

	_, err = keyStore.Eth().Create(&cltest.FixtureChainID)
	require.NoError(t, err)

	apiPrompt := cltest.NewMockAPIInitializer(t)
	client := cmd.Client{
		Config:                 cfg,
		Logger:                 logger.TestLogger(t),
		AppFactory:             cltest.InstanceAppFactory{App: app},
		FallbackAPIInitializer: apiPrompt,
		Runner:                 cltest.EmptyRunner{},
	}

	var keyState = ethkey.State{}
	err = db.Where("is_funding = TRUE").Find(&keyState).Error
	require.NoError(t, err)
	assert.Empty(t, keyState.ID, "expected no funding key")

	set := flag.NewFlagSet("test", 0)
	set.String("password", "../internal/fixtures/correct_password.txt", "")
	ctx := cli.NewContext(nil, set, nil)

	assert.NoError(t, client.RunNode(ctx))

	assert.NoError(t, db.Where("is_funding = TRUE").First(&keyState).Error)
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
			cfg := cltest.NewTestGeneralConfig(t)
			db := pgtest.NewGormDB(t)
			sessionORM := sessions.NewORM(postgres.UnwrapGormDB(db), time.Minute)
			// Clear out fixture
			err := sessionORM.DeleteUser()
			require.NoError(t, err)
			keyStore := cltest.NewKeyStore(t, db)
			_, err = keyStore.Eth().Create(&cltest.FixtureChainID)
			require.NoError(t, err)

			ethClient := cltest.NewEthClientMock(t)
			ethClient.On("Dial", mock.Anything).Return(nil)
			ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil)

			app := new(mocks.Application)
			app.On("SessionORM").Return(sessionORM)
			app.On("GetKeyStore").Return(keyStore)
			app.On("GetChainSet").Return(cltest.NewChainSetMockWithOneChain(t, ethClient, evmtest.NewChainScopedConfig(t, cfg))).Maybe()
			app.On("Start").Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)

			prompter := new(cmdMocks.Prompter)
			prompter.On("IsTerminal").Return(false).Once().Maybe()

			apiPrompt := cltest.NewMockAPIInitializer(t)
			client := cmd.Client{
				Config:                 cfg,
				Logger:                 logger.TestLogger(t),
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
			config := cltest.NewTestGeneralConfig(t)
			config.Overrides.Dev = null.BoolFrom(true)
			config.Overrides.LogToDisk = null.BoolFrom(tt.logToDiskValue)
			require.NoError(t, os.MkdirAll(config.RootDir(), os.FileMode(0700)))
			defer os.RemoveAll(config.RootDir())

			logger.ApplicationLogger(config).Sync()
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
	config, _, db := heavyweight.FullTestDB(t, "rebroadcasttransactions", true, true)
	keyStore := cltest.NewKeyStore(t, db)
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

	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, db, 7, 42, fromAddress)

	app := new(mocks.Application)
	app.On("GetDB").Return(db)
	app.On("GetKeyStore").Return(keyStore)
	app.On("Stop").Return(nil)
	ethClient := cltest.NewEthClientMockWithDefaultChain(t)
	app.On("GetChainSet").Return(cltest.NewChainSetMockWithOneChain(t, ethClient, evmtest.NewChainScopedConfig(t, config))).Maybe()
	ethClient.On("Dial", mock.Anything).Return(nil)

	client := cmd.Client{
		Config:                 config,
		Logger:                 logger.TestLogger(t),
		AppFactory:             cltest.InstanceAppFactory{App: app},
		FallbackAPIInitializer: cltest.NewMockAPIInitializer(t),
		Runner:                 cltest.EmptyRunner{},
	}

	config.SetDialect(dialects.TransactionWrappedPostgres)

	for i := beginningNonce; i <= endingNonce; i++ {
		n := i
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return uint(tx.Nonce()) == n
		})).Once().Return(nil)
	}

	assert.NoError(t, client.RebroadcastTransactions(c))

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
			config, _, db := heavyweight.FullTestDB(t, "rebroadcasttransactions_outsiderange", true, true)
			config.SetDialect(dialects.Postgres)

			keyStore := cltest.NewKeyStore(t, db)

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

			cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, db, int64(test.nonce), 42, fromAddress)

			app := new(mocks.Application)
			app.On("GetDB").Return(db)
			app.On("GetKeyStore").Return(keyStore)
			app.On("Stop").Return(nil)
			ethClient := cltest.NewEthClientMockWithDefaultChain(t)
			ethClient.On("Dial", mock.Anything).Return(nil)
			app.On("GetChainSet").Return(cltest.NewChainSetMockWithOneChain(t, ethClient, evmtest.NewChainScopedConfig(t, config))).Maybe()

			client := cmd.Client{
				Config:                 config,
				Logger:                 logger.TestLogger(t),
				AppFactory:             cltest.InstanceAppFactory{App: app},
				FallbackAPIInitializer: cltest.NewMockAPIInitializer(t),
				Runner:                 cltest.EmptyRunner{},
			}

			config.SetDialect(dialects.TransactionWrappedPostgres)

			for i := beginningNonce; i <= endingNonce; i++ {
				n := i
				ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
					return uint(tx.Nonce()) == n
				})).Once().Return(nil)
			}

			assert.NoError(t, client.RebroadcastTransactions(c))

			cltest.AssertEthTxAttemptCountStays(t, app.GetDB(), 1)
			app.AssertExpectations(t)
			ethClient.AssertExpectations(t)
		})
	}
}

func TestClient_SetNextNonce(t *testing.T) {
	// Need to use separate database
	config, _, db := heavyweight.FullTestDB(t, "setnextnonce", true, true)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()

	client := cmd.Client{
		Config: config,
		Logger: logger.TestLogger(t),
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
	require.NoError(t, db.First(&state).Error)
	require.NotNil(t, state.NextNonce)
	require.Equal(t, int64(42), state.NextNonce)
}
