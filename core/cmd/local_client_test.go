package cmd_test

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/chainlink/core/cmd"
	cmdMocks "github.com/smartcontractkit/chainlink/core/cmd/mocks"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/config/envvar"
	"github.com/smartcontractkit/chainlink/core/config/parse"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/logger/audit"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/store/dialects"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"go.uber.org/zap/zapcore"
)

func TestClient_RunNodeShowsEnv(t *testing.T) {
	// We must override the default log level to pass in CI, so check it separately first.
	llStr, ok := envvar.DefaultValue("LogLevel")
	require.False(t, ok)
	require.Empty(t, llStr)
	ll, invalid := parse.LogLevel(llStr)
	require.Empty(t, invalid)
	require.Equal(t, zapcore.InfoLevel, ll)

	lggr, observed := logger.TestLoggerObserved(t, zapcore.DebugLevel)

	cfg := config.NewGeneralConfig(lggr)
	require.NoError(t, cfg.SetLogLevel(zapcore.DebugLevel))

	db := pgtest.NewSqlxDB(t)
	sessionORM := sessions.NewORM(db, time.Minute, lggr, cfg, audit.NoopLogger)
	keyStore := cltest.NewKeyStore(t, db, cfg)
	_, err := keyStore.Eth().Create(&cltest.FixtureChainID)
	require.NoError(t, err)

	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	ethClient.On("Dial", mock.Anything).Return(nil).Maybe()
	ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil).Maybe()

	app := mocks.NewApplication(t)
	app.On("SessionORM").Return(sessionORM)
	app.On("GetKeyStore").Return(keyStore)
	app.On("GetChains").Return(chainlink.Chains{EVM: cltest.NewChainSetMockWithOneChain(t, ethClient, evmtest.NewChainScopedConfig(t, cfg))}).Maybe()
	app.On("Start", mock.Anything).Return(nil)
	app.On("Stop").Return(nil)
	app.On("ID").Return(uuid.NewV4())

	var logFileSize utils.FileSize
	err = logFileSize.UnmarshalText([]byte("100mb"))
	assert.NoError(t, err)

	runner := cltest.BlockedRunner{Done: make(chan struct{})}
	client := &cmd.Client{
		AppFactory:             cltest.InstanceAppFactory{App: app},
		FallbackAPIInitializer: cltest.NewMockAPIInitializer(t),
		Runner:                 runner,
	}

	// Start RunNode in a goroutine, it will block until we resume the runner
	awaiter := cltest.NewAwaiter()
	go func() {
		cliApp := cmd.NewApp(client)
		original := cliApp.Before
		cliApp.Before = func(c *cli.Context) error {
			if err := original(c); err != nil {
				return err
			}
			client.Logger = lggr
			client.Config = cfg
			return nil
		}
		assert.NoError(t, cliApp.
			Run([]string{"", "node", "start", "-debug", "-password", "../internal/fixtures/correct_password.txt"}))
		awaiter.ItHappened()
	}()

	// Unlock the runner to the client can begin shutdown
	select {
	case runner.Done <- struct{}{}:
	case <-time.After(30 * time.Second):
		t.Fatal("Timed out waiting for runner")
	}

	awaiter.AwaitOrFail(t)

	expected := fmt.Sprintf(`Environment variables
ADVISORY_LOCK_CHECK_INTERVAL: 1s
ADVISORY_LOCK_ID: 1027321974924625846
ALLOW_ORIGINS: http://localhost:3000,http://localhost:6688
BLOCK_BACKFILL_DEPTH: 10
BLOCK_HISTORY_ESTIMATOR_BLOCK_DELAY: 0
BLOCK_HISTORY_ESTIMATOR_BLOCK_HISTORY_SIZE: 0
BLOCK_HISTORY_ESTIMATOR_TRANSACTION_PERCENTILE: 0
BRIDGE_RESPONSE_URL: 
BRIDGE_CACHE_TTL: 
CHAIN_TYPE: 
DATABASE_BACKUP_FREQUENCY: 1h0m0s
DATABASE_BACKUP_MODE: none
DATABASE_BACKUP_ON_VERSION_UPGRADE: true
DATABASE_LOCKING_MODE: dual
ETH_CHAIN_ID: <nil>
DEFAULT_HTTP_LIMIT: 32768
DEFAULT_HTTP_TIMEOUT: 15s
CHAINLINK_DEV: false
SHUTDOWN_GRACE_PERIOD: 5s
EVM_RPC_ENABLED: true
ETH_HTTP_URL: 
ETH_SECONDARY_URLS: []
ETH_URL: 
EXPLORER_URL: 
FM_DEFAULT_TRANSACTION_QUEUE_DEPTH: 1
FEATURE_EXTERNAL_INITIATORS: false
FEATURE_OFFCHAIN_REPORTING: false
GAS_ESTIMATOR_MODE: 
INSECURE_FAST_SCRYPT: false
JSON_CONSOLE: false
JOB_PIPELINE_REAPER_INTERVAL: 1h0m0s
JOB_PIPELINE_REAPER_THRESHOLD: 24h0m0s
KEEPER_DEFAULT_TRANSACTION_QUEUE_DEPTH: 1
KEEPER_GAS_PRICE_BUFFER_PERCENT: 20
KEEPER_GAS_TIP_CAP_BUFFER_PERCENT: 20
KEEPER_BASE_FEE_BUFFER_PERCENT: 20
KEEPER_MAXIMUM_GRACE_PERIOD: 100
KEEPER_REGISTRY_CHECK_GAS_OVERHEAD: 200000
KEEPER_REGISTRY_PERFORM_GAS_OVERHEAD: 300000
KEEPER_REGISTRY_MAX_PERFORM_DATA_SIZE: 5000
KEEPER_REGISTRY_SYNC_INTERVAL: 30m0s
KEEPER_REGISTRY_SYNC_UPKEEP_QUEUE_SIZE: 10
KEEPER_TURN_LOOK_BACK: 1000
LEASE_LOCK_DURATION: 10s
LEASE_LOCK_REFRESH_INTERVAL: 1s
FLAGS_CONTRACT_ADDRESS: 
LINK_CONTRACT_ADDRESS: 
LOG_FILE_DIR: %[1]s
LOG_LEVEL: debug
LOG_SQL: false
LOG_FILE_MAX_SIZE: 5.12gb
LOG_FILE_MAX_AGE: 0
LOG_FILE_MAX_BACKUPS: 1
TRIGGER_FALLBACK_DB_POLL_INTERVAL: 30s
AUDIT_LOGGER_ENABLED: false
AUDIT_LOGGER_FORWARD_TO_URL: 
AUDIT_LOGGER_JSON_WRAPPER_KEY: 
AUDIT_LOGGER_HEADERS: 
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
ROOT: %[1]s
SECURE_COOKIES: true
SESSION_TIMEOUT: 15m0s
TELEMETRY_INGRESS_LOGGING: false
TELEMETRY_INGRESS_SERVER_PUB_KEY: 
TELEMETRY_INGRESS_URL: 
CHAINLINK_TLS_HOST: 
CHAINLINK_TLS_PORT: 6689
CHAINLINK_TLS_REDIRECT: false`, cfg.RootDir())

	logs := observed.Filter(func(e observer.LoggedEntry) bool {
		return strings.Contains(e.Message, "Environment variables")
	}).All()
	require.Len(t, logs, 1)
	require.Contains(t, logs[0].Message, expected)
}

func TestClient_RunNodeWithPasswords(t *testing.T) {
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
			cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				s.Password.Keystore = models.NewSecret("dummy")
				c.EVM[0].Nodes[0].Name = ptr("fake")
				c.EVM[0].Nodes[0].HTTPURL = models.MustParseURL("http://fake.com")
			})
			db := pgtest.NewSqlxDB(t)
			keyStore := cltest.NewKeyStore(t, db, cfg)
			sessionORM := sessions.NewORM(db, time.Minute, logger.TestLogger(t), cfg, audit.NoopLogger)

			// Purge the fixture users to test assumption of single admin
			// initialUser user created above
			pgtest.MustExec(t, db, "DELETE FROM users;")

			app := new(mocks.Application)
			app.On("SessionORM").Return(sessionORM)
			app.On("GetKeyStore").Return(keyStore)
			app.On("GetChains").Return(chainlink.Chains{EVM: cltest.NewChainSetMockWithOneChain(t, evmtest.NewEthClientMock(t), evmtest.NewChainScopedConfig(t, cfg))}).Maybe()
			app.On("Start", mock.Anything).Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)
			app.On("ID").Maybe().Return(uuid.NewV4())

			ethClient := evmtest.NewEthClientMock(t)
			ethClient.On("Dial", mock.Anything).Return(nil).Maybe()
			ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil).Maybe()

			cltest.MustInsertRandomKey(t, keyStore.Eth())
			apiPrompt := cltest.NewMockAPIInitializer(t)
			lggr := logger.TestLogger(t)

			client := cmd.Client{
				Config:                 cfg,
				FallbackAPIInitializer: apiPrompt,
				Runner:                 cltest.EmptyRunner{},
				AppFactory:             cltest.InstanceAppFactory{App: app},
				Logger:                 lggr,
			}

			set := flag.NewFlagSet("test", 0)
			cltest.CopyFlagSetFromAction(client.RunNode, set, "")

			require.NoError(t, set.Set("password", test.pwdfile))

			c := cli.NewContext(nil, set, nil)

			run := func() error {
				cli := cmd.NewApp(&client)
				if err := cli.Before(c); err != nil {
					return err
				}
				if err := client.RunNode(c); err != nil {
					return err
				}
				return nil
			}

			if test.wantUnlocked {
				assert.NoError(t, run())
				assert.Equal(t, 1, apiPrompt.Count)
			} else {
				assert.Error(t, run())
				assert.Equal(t, 0, apiPrompt.Count)
			}
		})
	}
}

func TestClient_RunNodeWithAPICredentialsFile(t *testing.T) {
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
			cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				s.Password.Keystore = models.NewSecret("16charlengthp4SsW0rD1!@#_")
				c.EVM[0].Nodes[0].Name = ptr("fake")
				c.EVM[0].Nodes[0].WSURL = models.MustParseURL("WSS://fake.com/ws")
				c.EVM[0].Nodes[0].HTTPURL = models.MustParseURL("http://fake.com")
			})
			db := pgtest.NewSqlxDB(t)
			sessionORM := sessions.NewORM(db, time.Minute, logger.TestLogger(t), cfg, audit.NoopLogger)

			// Clear out fixture users/users created from the other test cases
			// This asserts that on initial run with an empty users table that the credentials file will instantiate and
			// create/run with a new admin user
			pgtest.MustExec(t, db, "DELETE FROM users;")

			keyStore := cltest.NewKeyStore(t, db, cfg)
			_, err := keyStore.Eth().Create(&cltest.FixtureChainID)
			require.NoError(t, err)

			ethClient := evmtest.NewEthClientMock(t)
			ethClient.On("Dial", mock.Anything).Return(nil).Maybe()
			ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil).Maybe()

			app := new(mocks.Application)
			app.On("SessionORM").Return(sessionORM)
			app.On("GetKeyStore").Return(keyStore)
			app.On("GetChains").Return(chainlink.Chains{EVM: cltest.NewChainSetMockWithOneChain(t, ethClient, evmtest.NewChainScopedConfig(t, cfg))}).Maybe()
			app.On("Start", mock.Anything).Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)
			app.On("ID").Maybe().Return(uuid.NewV4())

			prompter := new(cmdMocks.Prompter)
			prompter.On("IsTerminal").Return(false).Once().Maybe()

			apiPrompt := cltest.NewMockAPIInitializer(t)
			lggr := logger.TestLogger(t)

			client := cmd.Client{
				Config:                 cfg,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				KeyStoreAuthenticator:  cmd.TerminalKeyStoreAuthenticator{prompter},
				FallbackAPIInitializer: apiPrompt,
				Runner:                 cltest.EmptyRunner{},
				Logger:                 lggr,
			}

			set := flag.NewFlagSet("test", 0)
			cltest.CopyFlagSetFromAction(client.RunNode, set, "")

			require.NoError(t, set.Set("api", test.apiFile))

			c := cli.NewContext(nil, set, nil)

			if test.wantError {
				err = client.RunNode(c)
				assert.ErrorContains(t, err, "error creating api initializer: open doesntexist.txt: no such file or directory")
			} else {
				assert.NoError(t, client.RunNode(c))
			}

			assert.Equal(t, test.wantPrompt, apiPrompt.Count > 0)
		})
	}
}

func TestClient_DiskMaxSizeBeforeRotateOptionDisablesAsExpected(t *testing.T) {
	tests := []struct {
		name            string
		logFileSize     func(t *testing.T) utils.FileSize
		fileShouldExist bool
	}{
		{"DiskMaxSizeBeforeRotate = 0 => no log on disk", func(t *testing.T) utils.FileSize {
			return 0
		}, false},
		{"DiskMaxSizeBeforeRotate > 0 => log on disk (positive control)", func(t *testing.T) utils.FileSize {
			var logFileSize utils.FileSize
			err := logFileSize.UnmarshalText([]byte("100mb"))
			assert.NoError(t, err)

			return logFileSize
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := logger.Config{
				Dir:           t.TempDir(),
				FileMaxSizeMB: int(tt.logFileSize(t) / utils.MB),
			}
			assert.NoError(t, os.MkdirAll(cfg.Dir, os.FileMode(0700)))

			lggr, close := cfg.New()
			t.Cleanup(func() { assert.NoError(t, close()) })

			// Tries to create a log file by logging. The log file won't be created if there's no logging happening.
			lggr.Debug("Trying to create a log file by logging.")

			_, err := os.Stat(cfg.LogsFile())
			require.Equal(t, os.IsNotExist(err), !tt.fileShouldExist)
		})
	}
}

func TestClient_RebroadcastTransactions_Txm(t *testing.T) {
	// Use a non-transactional db for this test because we need to
	// test multiple connections to the database, and changes made within
	// the transaction cannot be seen from another connection.
	config, sqlxDB := heavyweight.FullTestDBV2(t, "rebroadcasttransactions", nil)
	keyStore := cltest.NewKeyStore(t, sqlxDB, config)
	_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth(), 0)

	borm := cltest.NewTxmORM(t, sqlxDB, config)
	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, 7, 42, fromAddress)

	app := mocks.NewApplication(t)
	app.On("GetSqlxDB").Return(sqlxDB)
	app.On("GetKeyStore").Return(keyStore)
	app.On("Stop").Return(nil)
	app.On("ID").Maybe().Return(uuid.NewV4())
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	app.On("GetChains").Return(chainlink.Chains{EVM: cltest.NewChainSetMockWithOneChain(t, ethClient, evmtest.NewChainScopedConfig(t, config))}).Maybe()
	ethClient.On("Dial", mock.Anything).Return(nil)

	lggr := logger.TestLogger(t)

	client := cmd.Client{
		Config:                 config,
		AppFactory:             cltest.InstanceAppFactory{App: app},
		FallbackAPIInitializer: cltest.NewMockAPIInitializer(t),
		Runner:                 cltest.EmptyRunner{},
		Logger:                 lggr,
	}

	beginningNonce := uint(7)
	endingNonce := uint(10)
	gasPrice := big.NewInt(100000000000)
	gasLimit := uint64(3000000)
	set := flag.NewFlagSet("test", 0)
	cltest.CopyFlagSetFromAction(client.RebroadcastTransactions, set, "")

	require.NoError(t, set.Set("beginningNonce", strconv.FormatUint(uint64(beginningNonce), 10)))
	require.NoError(t, set.Set("endingNonce", strconv.FormatUint(uint64(endingNonce), 10)))
	require.NoError(t, set.Set("gasPriceWei", gasPrice.String()))
	require.NoError(t, set.Set("gasLimit", strconv.FormatUint(gasLimit, 10)))
	require.NoError(t, set.Set("address", fromAddress.Hex()))
	require.NoError(t, set.Set("password", "../internal/fixtures/correct_password.txt"))

	c := cli.NewContext(nil, set, nil)

	for i := beginningNonce; i <= endingNonce; i++ {
		n := i
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return uint(tx.Nonce()) == n
		})).Once().Return(nil)
	}

	assert.NoError(t, client.RebroadcastTransactions(c))
}

func TestClient_RebroadcastTransactions_OutsideRange_Txm(t *testing.T) {
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
			config, sqlxDB := heavyweight.FullTestDBV2(t, "rebroadcasttransactions_outsiderange", func(c *chainlink.Config, s *chainlink.Secrets) {
				c.Database.Dialect = dialects.Postgres
			})

			keyStore := cltest.NewKeyStore(t, sqlxDB, config)

			_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth(), 0)

			borm := cltest.NewTxmORM(t, sqlxDB, config)
			cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, borm, int64(test.nonce), 42, fromAddress)

			app := mocks.NewApplication(t)
			app.On("GetSqlxDB").Return(sqlxDB)
			app.On("GetKeyStore").Return(keyStore)
			app.On("Stop").Return(nil)
			app.On("ID").Maybe().Return(uuid.NewV4())
			ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
			ethClient.On("Dial", mock.Anything).Return(nil)
			app.On("GetChains").Return(chainlink.Chains{EVM: cltest.NewChainSetMockWithOneChain(t, ethClient, evmtest.NewChainScopedConfig(t, config))}).Maybe()

			lggr := logger.TestLogger(t)

			client := cmd.Client{
				Config:                 config,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				FallbackAPIInitializer: cltest.NewMockAPIInitializer(t),
				Runner:                 cltest.EmptyRunner{},
				Logger:                 lggr,
			}

			set := flag.NewFlagSet("test", 0)
			cltest.CopyFlagSetFromAction(client.RebroadcastTransactions, set, "")

			require.NoError(t, set.Set("beginningNonce", strconv.FormatUint(uint64(beginningNonce), 10)))
			require.NoError(t, set.Set("endingNonce", strconv.FormatUint(uint64(endingNonce), 10)))
			require.NoError(t, set.Set("gasPriceWei", gasPrice.String()))
			require.NoError(t, set.Set("gasLimit", strconv.FormatUint(gasLimit, 10)))
			require.NoError(t, set.Set("address", fromAddress.Hex()))
			require.NoError(t, set.Set("password", "../internal/fixtures/correct_password.txt"))
			c := cli.NewContext(nil, set, nil)

			for i := beginningNonce; i <= endingNonce; i++ {
				n := i
				ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
					return uint(tx.Nonce()) == n
				})).Once().Return(nil)
			}

			assert.NoError(t, client.RebroadcastTransactions(c))

			cltest.AssertEthTxAttemptCountStays(t, app.GetSqlxDB(), 1)
		})
	}
}
