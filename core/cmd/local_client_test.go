package cmd_test

import (
	"database/sql"
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/cmd"
	cmdMocks "github.com/smartcontractkit/chainlink/core/cmd/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/store/dialects"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/kylelemons/godebug/diff"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"go.uber.org/zap/zapcore"
	null "gopkg.in/guregu/null.v4"
)

func TestClient_RunNodeShowsEnv(t *testing.T) {
	cfg := cltest.NewTestGeneralConfig(t)
	debug := zapcore.DebugLevel
	cfg.Overrides.LogLevel = &debug
	cfg.Overrides.LogToDisk = null.BoolFrom(true)
	db := pgtest.NewSqlxDB(t)
	sessionORM := sessions.NewORM(db, time.Minute, logger.TestLogger(t))
	keyStore := cltest.NewKeyStore(t, db, cfg)
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
	app.On("ID").Return(uuid.NewV4())

	lcfg := logger.Config{
		LogLevel: zapcore.DebugLevel,
		ToDisk:   true,
		Dir:      t.TempDir(),
	}

	runner := cltest.BlockedRunner{Done: make(chan struct{})}
	client := cmd.Client{
		Config:                 cfg,
		Logger:                 lcfg.New(),
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

	logs, err := cltest.ReadLogs(lcfg.Dir)
	require.NoError(t, err)
	msg := cltest.FindLogMessage(logs, func(msg string) bool {
		return strings.HasPrefix(msg, "Environment variables")
	})
	require.NotEmpty(t, msg, "No env var message found")

	expected := fmt.Sprintf(`Environment variables
ADVISORY_LOCK_CHECK_INTERVAL: 1s
ADVISORY_LOCK_ID: 1027321974924625846
ALLOW_ORIGINS: http://localhost:3000,http://localhost:6688
BLOCK_BACKFILL_DEPTH: 10
BLOCK_HISTORY_ESTIMATOR_BLOCK_DELAY: 0
BLOCK_HISTORY_ESTIMATOR_BLOCK_HISTORY_SIZE: 0
BLOCK_HISTORY_ESTIMATOR_TRANSACTION_PERCENTILE: 0
BRIDGE_RESPONSE_URL: http://localhost:6688
CHAIN_TYPE: 
CLIENT_NODE_URL: http://localhost:6688
DATABASE_BACKUP_FREQUENCY: 1h0m0s
DATABASE_BACKUP_MODE: none
DATABASE_LOCKING_MODE: none
ETH_CHAIN_ID: 0
DEFAULT_HTTP_LIMIT: 32768
DEFAULT_HTTP_TIMEOUT: 15s
CHAINLINK_DEV: true
ETH_DISABLED: false
ETH_HTTP_URL: 
ETH_SECONDARY_URLS: []
ETH_URL: 
EXPLORER_URL: 
FM_DEFAULT_TRANSACTION_QUEUE_DEPTH: 1
FEATURE_EXTERNAL_INITIATORS: false
FEATURE_OFFCHAIN_REPORTING: false
GAS_ESTIMATOR_MODE: 
INSECURE_FAST_SCRYPT: true
JSON_CONSOLE: false
JOB_PIPELINE_REAPER_INTERVAL: 1h0m0s
JOB_PIPELINE_REAPER_THRESHOLD: 24h0m0s
KEEPER_DEFAULT_TRANSACTION_QUEUE_DEPTH: 1
KEEPER_GAS_PRICE_BUFFER_PERCENT: 20
KEEPER_GAS_TIP_CAP_BUFFER_PERCENT: 20
KEEPER_MAXIMUM_GRACE_PERIOD: 0
KEEPER_REGISTRY_CHECK_GAS_OVERHEAD: 0
KEEPER_REGISTRY_PERFORM_GAS_OVERHEAD: 0
KEEPER_REGISTRY_SYNC_INTERVAL: 
KEEPER_REGISTRY_SYNC_UPKEEP_QUEUE_SIZE: 0
LEASE_LOCK_DURATION: 30s
LEASE_LOCK_REFRESH_INTERVAL: 1s
FLAGS_CONTRACT_ADDRESS: 
LINK_CONTRACT_ADDRESS: 
LOG_FILE_DIR: %[1]s
LOG_LEVEL: debug
LOG_SQL: false
LOG_SQL_MIGRATIONS: false
LOG_TO_DISK: true
TRIGGER_FALLBACK_DB_POLL_INTERVAL: 30s
OCR_CONTRACT_TRANSMITTER_TRANSMIT_TIMEOUT: 
OCR_DATABASE_TIMEOUT: 
OCR_DEFAULT_TRANSACTION_QUEUE_DEPTH: 1
OCR_TRACE_LOGGING: false
P2P_NETWORKING_STACK: V1
P2P_PEER_ID: 
P2P_INCOMING_MESSAGE_BUFFER_SIZE: 10
P2P_OUTGOING_MESSAGE_BUFFER_SIZE: 10
P2P_BOOTSTRAP_PEERS: []
P2P_LISTEN_IP: 0.0.0.0
P2P_LISTEN_PORT: 
P2P_NEW_STREAM_TIMEOUT: 10s
P2P_DHT_LOOKUP_INTERVAL: 10
P2P_BOOTSTRAP_CHECK_INTERVAL: 20s
P2PV2_ANNOUNCE_ADDRESSES: []
P2PV2_BOOTSTRAPPERS: []
P2PV2_DELTA_DIAL: 15s
P2PV2_DELTA_RECONCILE: 1m0s
P2PV2_LISTEN_ADDRESSES: []
CHAINLINK_PORT: 6688
REAPER_EXPIRATION: 240h0m0s
REPLAY_FROM_BLOCK: -1
ROOT: %[1]s
SECURE_COOKIES: true
SESSION_TIMEOUT: 2m0s
TELEMETRY_INGRESS_LOGGING: false
TELEMETRY_INGRESS_SERVER_PUB_KEY: 
TELEMETRY_INGRESS_URL: 
CHAINLINK_TLS_HOST: 
CHAINLINK_TLS_PORT: 6689
CHAINLINK_TLS_REDIRECT: false`, cfg.RootDir())

	if !strings.Contains(msg, expected) {
		t.Errorf("Expected to find:\n\n%s\n\nWithin:\n\n%s\n\nDiff:\n\n%s", expected, msg, diff.Diff(expected, msg))
	}

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
			db := pgtest.NewSqlxDB(t)
			keyStore := cltest.NewKeyStore(t, db, cfg)
			sessionORM := sessions.NewORM(db, time.Minute, logger.TestLogger(t))
			// Clear out fixture
			err := sessionORM.DeleteUser()
			require.NoError(t, err)

			app := new(mocks.Application)
			app.On("SessionORM").Return(sessionORM)
			app.On("GetKeyStore").Return(keyStore)
			app.On("GetChainSet").Return(cltest.NewChainSetMockWithOneChain(t, cltest.NewEthClientMock(t), evmtest.NewChainScopedConfig(t, cfg))).Maybe()
			app.On("Start").Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)
			app.On("ID").Maybe().Return(uuid.NewV4())

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
	db := pgtest.NewSqlxDB(t)
	sessionORM := sessions.NewORM(db, time.Minute, logger.TestLogger(t))
	keyStore := cltest.NewKeyStore(t, db, cfg)
	_, err := keyStore.Eth().Create(&cltest.FixtureChainID)
	require.NoError(t, err)

	app := new(mocks.Application)
	app.On("SessionORM").Return(sessionORM)
	app.On("GetKeyStore").Return(keyStore)
	app.On("GetChainSet").Return(cltest.NewChainSetMockWithOneChain(t, cltest.NewEthClientMock(t), evmtest.NewChainScopedConfig(t, cfg))).Maybe()
	app.On("Start").Maybe().Return(nil)
	app.On("Stop").Maybe().Return(nil)
	app.On("ID").Maybe().Return(uuid.NewV4())

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
	err = db.Get(&keyState, `SELECT * FROM eth_key_states WHERE is_funding = TRUE`)
	assert.EqualError(t, err, sql.ErrNoRows.Error())

	set := flag.NewFlagSet("test", 0)
	set.String("password", "../internal/fixtures/correct_password.txt", "")
	ctx := cli.NewContext(nil, set, nil)

	assert.NoError(t, client.RunNode(ctx))

	assert.NoError(t, db.Get(&keyState, `SELECT * FROM eth_key_states WHERE is_funding = TRUE`))
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
			db := pgtest.NewSqlxDB(t)
			sessionORM := sessions.NewORM(db, time.Minute, logger.TestLogger(t))
			// Clear out fixture
			err := sessionORM.DeleteUser()
			require.NoError(t, err)
			keyStore := cltest.NewKeyStore(t, db, cfg)
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
			app.On("ID").Maybe().Return(uuid.NewV4())

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
				err = client.RunNode(c)
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), "error creating api initializer: open doesntexist.txt: no such file or directory")
				}
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
			cfg := logger.Config{
				ToDisk: tt.logToDiskValue,
				Dir:    t.TempDir(),
			}
			require.NoError(t, os.MkdirAll(cfg.Dir, os.FileMode(0700)))

			cfg.New().Sync()
			filepath := filepath.Join(cfg.Dir, "log.jsonl")
			_, err := os.Stat(filepath)
			assert.Equal(t, os.IsNotExist(err), !tt.fileShouldExist)
		})
	}
}

func TestClient_RebroadcastTransactions_BPTXM(t *testing.T) {
	// Use the a non-transactional db for this test because we need to
	// test multiple connections to the database, and changes made within
	// the transaction cannot be seen from another connection.
	config, sqlxDB := heavyweight.FullTestDB(t, "rebroadcasttransactions", true, true)
	keyStore := cltest.NewKeyStore(t, sqlxDB, config)
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

	borm := cltest.NewBulletproofTxManagerORM(t, sqlxDB, config)
	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 7, 42, fromAddress)

	app := new(mocks.Application)
	app.Test(t)
	app.On("GetSqlxDB").Return(sqlxDB)
	app.On("GetKeyStore").Return(keyStore)
	app.On("Stop").Return(nil)
	app.On("ID").Maybe().Return(uuid.NewV4())
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

	config.Overrides.Dialect = dialects.TransactionWrappedPostgres

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
			config, sqlxDB := heavyweight.FullTestDB(t, "rebroadcasttransactions_outsiderange", true, true)
			config.Overrides.Dialect = dialects.Postgres

			keyStore := cltest.NewKeyStore(t, sqlxDB, config)

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

			borm := cltest.NewBulletproofTxManagerORM(t, sqlxDB, config)
			cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, int64(test.nonce), 42, fromAddress)

			app := new(mocks.Application)
			app.Test(t)
			app.On("GetSqlxDB").Return(sqlxDB)
			app.On("GetKeyStore").Return(keyStore)
			app.On("Stop").Return(nil)
			app.On("ID").Maybe().Return(uuid.NewV4())
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

			config.Overrides.Dialect = dialects.TransactionWrappedPostgres

			for i := beginningNonce; i <= endingNonce; i++ {
				n := i
				ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
					return uint(tx.Nonce()) == n
				})).Once().Return(nil)
			}

			assert.NoError(t, client.RebroadcastTransactions(c))

			cltest.AssertEthTxAttemptCountStays(t, app.GetSqlxDB(), 1)
			app.AssertExpectations(t)
			ethClient.AssertExpectations(t)
		})
	}
}

func TestClient_SetNextNonce(t *testing.T) {
	// Need to use separate database
	config, sqlxDB := heavyweight.FullTestDB(t, "setnextnonce", true, true)
	ethKeyStore := cltest.NewKeyStore(t, sqlxDB, config).Eth()

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
	require.NoError(t, sqlxDB.Get(&state, `SELECT * FROM eth_key_states`))
	require.NotNil(t, state.NextNonce)
	require.Equal(t, int64(42), state.NextNonce)
}
