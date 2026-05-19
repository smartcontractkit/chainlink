package cmd_test

import (
	"context"
	"errors"
	"flag"
	"math/big"
	"net/url"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"go.uber.org/zap/zapcore"

	commonkeystore "github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	pgcommon "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr/txmgrtest"
	"github.com/smartcontractkit/chainlink-framework/multinode"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	cmdMocks "github.com/smartcontractkit/chainlink/v2/core/cmd/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	chainlinkmocks "github.com/smartcontractkit/chainlink/v2/core/services/chainlink/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/sessions/localauth"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

// resetShellForTest resets the Shell's cleanup guard for test isolation.
// This allows multiple initialization/cleanup cycles in tests.
func resetShellForTest(shell *cmd.Shell) {
	shell.CleanupOnce = sync.Once{}
	shell.LDB = nil
	shell.DS = nil
	shell.KeyStore = nil
	shell.BeholderClient = nil
}

type stubLockedDB struct{}

func (stubLockedDB) Open(context.Context) error { return nil }
func (stubLockedDB) Close() error               { return nil }
func (stubLockedDB) DB() *sqlx.DB               { return nil }

func genTestEVMRelayers(t *testing.T, cfg chainlink.GeneralConfig, ds sqlutil.DataSource, ethKeystore keystore.Eth, csaKeystore core.Keystore) *chainlink.CoreRelayerChainInteroperators {
	lggr := logger.TestLogger(t)
	f := chainlink.RelayerFactory{
		Logger: lggr,
		LoopRegistry: plugins.NewLoopRegistry(lggr, cfg.AppID().String(), cfg.Feature().LogPoller(), cfg.Database(),
			cfg.Mercury(), cfg.Pyroscope(), cfg.AutoPprof(), cfg.Tracing(), cfg.Telemetry(), nil, "", cfg.LOOPP()),
		CapabilitiesRegistry: capabilities.NewRegistry(lggr),
	}

	relayers, err := chainlink.NewCoreRelayerChainInteroperators(chainlink.InitEVM(f, chainlink.EVMFactoryConfig{
		ChainOpts: legacyevm.ChainOpts{
			ChainConfigs:   cfg.EVMConfigs(),
			DatabaseConfig: cfg.Database(),
			ListenerConfig: cfg.Database().Listener(),
			FeatureConfig:  cfg.Feature(),
			MailMon:        &mailbox.Monitor{},
			DS:             ds,
		},
		EthKeystore: ethKeystore,
		CSAKeystore: csaKeystore,
	}))
	if err != nil {
		t.Fatal(err)
	}
	return relayers
}

func TestShell_RunNodeWithAPICredentialsFile(t *testing.T) {
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
				c.EVM[0].Nodes[0].WSURL = commonconfig.MustParseURL("WSS://fake.com/ws")
				c.EVM[0].Nodes[0].HTTPURL = commonconfig.MustParseURL("http://fake.com")
				// seems to be needed for config validate
				c.Insecure.OCRDevelopmentMode = nil
			})
			db := pgtest.NewSqlxDB(t)
			authProviderORM := localauth.NewORM(db, time.Minute, logger.TestLogger(t), audit.NoopLogger)

			// Clear out fixture users/users created from the other test cases
			// This asserts that on initial run with an empty users table that the credentials file will instantiate and
			// create/run with a new admin user
			pgtest.MustExec(t, db, "DELETE FROM users;")

			keyStore := cltest.NewKeyStore(t, db)
			_, err := keyStore.Eth().Create(testutils.Context(t), &cltest.FixtureChainID)
			require.NoError(t, err)

			ethClient := evmtest.NewEthClientMock(t)
			ethClient.On("Dial", mock.Anything).Return(nil).Maybe()
			ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil).Maybe()

			testRelayers := genTestEVMRelayers(t, cfg, db, keyStore.Eth(), &keystore.CSASigner{CSA: keyStore.CSA()})
			app := mocks.NewApplication(t)
			app.On("BasicAdminUsersORM").Return(authProviderORM)
			app.On("GetKeyStore").Return(keyStore)
			app.On("GetRelayers").Return(testRelayers).Maybe()
			app.On("Start", mock.Anything).Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)
			app.On("ID").Maybe().Return(uuid.New())

			prompter := cmdMocks.NewPrompter(t)

			apiPrompt := cltest.NewMockAPIInitializer(t)

			client := cmd.Shell{
				Config:                 cfg,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				KeyStoreAuthenticator:  cmd.TerminalKeyStoreAuthenticator{prompter},
				FallbackAPIInitializer: apiPrompt,
				Runner:                 cltest.EmptyRunner{},
				Logger:                 logger.TestLogger(t),
			}

			set := flag.NewFlagSet("test", 0)
			flagSetApplyFromAction(client.RunNode, set, "")

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

func TestShell_DiskMaxSizeBeforeRotateOptionDisablesAsExpected(t *testing.T) {
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
			require.NoError(t, os.MkdirAll(cfg.Dir, os.FileMode(0o700)))

			lggr, closeFn := cfg.New()
			t.Cleanup(func() { assert.NoError(t, closeFn()) })

			// Tries to create a log file by logging. The log file won't be created if there's no logging happening.
			lggr.Debug("Trying to create a log file by logging.")

			_, err := os.Stat(cfg.LogsFile())
			require.Equal(t, os.IsNotExist(err), !tt.fileShouldExist)
		})
	}
}

func TestShell_RebroadcastTransactions_Txm(t *testing.T) {
	t.Parallel()
	// Use a non-transactional db for this test because we need to
	// test multiple connections to the database, and changes made within
	// the transaction cannot be seen from another connection.
	config, sqlxDB := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.DriverName = pgcommon.DriverPostgres
		// evm config is used in this test. but if set, it must be pass config validation.
		// simplest to make it nil
		c.EVM = nil
		// seems to be needed for config validate
		c.Insecure.OCRDevelopmentMode = nil
	})
	keyStore := cltest.NewKeyStore(t, sqlxDB)
	_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth())

	txStore := txmgrtest.NewTestTxStore(t, sqlxDB)
	txmgrtest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 7, 42, fromAddress)

	lggr := logger.TestLogger(t)

	app := mocks.NewApplication(t)
	app.On("GetDB").Return(sqlxDB)
	app.On("GetKeyStore").Return(keyStore)
	app.On("ID").Maybe().Return(uuid.New())
	app.On("GetConfig").Return(config)
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	legacy := cltest.NewLegacyChainsWithMockChain(t, ethClient, config)

	mockRelayerChainInteroperators := &chainlinkmocks.FakeRelayerChainInteroperators{EVMChains: legacy}
	app.On("GetRelayers").Return(mockRelayerChainInteroperators).Maybe()
	ethClient.On("Dial", mock.Anything).Return(nil)

	c := cmd.Shell{
		Config:                 config,
		AppFactory:             cltest.InstanceAppFactory{App: app},
		FallbackAPIInitializer: cltest.NewMockAPIInitializer(t),
		Runner:                 cltest.EmptyRunner{},
		Logger:                 lggr,
	}

	beginningNonce := uint64(7)
	endingNonce := uint64(10)
	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(c.RebroadcastTransactions, set, "")

	require.NoError(t, set.Set("evmChainID", testutils.FixtureChainID.String()))
	require.NoError(t, set.Set("beginningNonce", strconv.FormatUint(beginningNonce, 10)))
	require.NoError(t, set.Set("endingNonce", strconv.FormatUint(endingNonce, 10)))
	require.NoError(t, set.Set("gasPriceWei", "100000000000"))
	require.NoError(t, set.Set("gasLimit", "3000000"))
	require.NoError(t, set.Set("address", fromAddress.Hex()))
	require.NoError(t, set.Set("password", "../internal/fixtures/correct_password.txt"))

	ctx := cli.NewContext(nil, set, nil)

	// Run before hook to initialize components with authentication
	err := c.BeforeNode(ctx)
	require.NoError(t, err)

	for i := beginningNonce; i <= endingNonce; i++ {
		n := i
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == n
		}), mock.Anything).Once().Return(multinode.Successful, nil)
	}

	assert.NoError(t, c.RebroadcastTransactions(ctx))
}

func TestShell_RebroadcastTransactions_OutsideRange_Txm(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			// Use the non-transactional db for this test because we need to
			// test multiple connections to the database, and changes made within
			// the transaction cannot be seen from another connection.
			config, sqlxDB := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.Database.DriverName = pgcommon.DriverPostgres
				// evm config is used in this test. but if set, it must be pass config validation.
				// simplest to make it nil
				c.EVM = nil
				// seems to be needed for config validate
				c.Insecure.OCRDevelopmentMode = nil
			})

			keyStore := cltest.NewKeyStore(t, sqlxDB)

			_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth())

			txStore := txmgrtest.NewTestTxStore(t, sqlxDB)
			txmgrtest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, int64(test.nonce), 42, fromAddress)

			lggr := logger.TestLogger(t)

			app := mocks.NewApplication(t)
			app.On("GetDB").Return(sqlxDB)
			app.On("GetKeyStore").Return(keyStore)
			app.On("ID").Maybe().Return(uuid.New())
			app.On("GetConfig").Return(config)
			ethClient := clienttest.NewClientWithDefaultChainID(t)
			ethClient.On("Dial", mock.Anything).Return(nil)
			legacy := cltest.NewLegacyChainsWithMockChain(t, ethClient, config)

			mockRelayerChainInteroperators := &chainlinkmocks.FakeRelayerChainInteroperators{EVMChains: legacy}
			app.On("GetRelayers").Return(mockRelayerChainInteroperators).Maybe()

			c := cmd.Shell{
				Config:                 config,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				FallbackAPIInitializer: cltest.NewMockAPIInitializer(t),
				Runner:                 cltest.EmptyRunner{},
				Logger:                 lggr,
			}

			set := flag.NewFlagSet("test", 0)
			flagSetApplyFromAction(c.RebroadcastTransactions, set, "")

			require.NoError(t, set.Set("evmChainID", testutils.FixtureChainID.String()))
			require.NoError(t, set.Set("beginningNonce", strconv.FormatUint(uint64(beginningNonce), 10)))
			require.NoError(t, set.Set("endingNonce", strconv.FormatUint(uint64(endingNonce), 10)))
			require.NoError(t, set.Set("gasPriceWei", gasPrice.String()))
			require.NoError(t, set.Set("gasLimit", strconv.FormatUint(gasLimit, 10)))
			require.NoError(t, set.Set("address", fromAddress.Hex()))

			require.NoError(t, set.Set("password", "../internal/fixtures/correct_password.txt"))
			ctx := cli.NewContext(nil, set, nil)

			// Run before hook to initialize components with authentication
			err := c.BeforeNode(ctx)
			require.NoError(t, err)

			for i := beginningNonce; i <= endingNonce; i++ {
				n := i
				ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
					return uint(tx.Nonce()) == n
				}), mock.Anything).Once().Return(multinode.Successful, nil)
			}

			assert.NoError(t, c.RebroadcastTransactions(ctx))

			cltest.AssertEthTxAttemptCountStays(t, txStore, 1)
		})
	}
}

func TestShell_RebroadcastTransactions_AddressCheck(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		enableAddress bool
		shouldError   bool
		errorContains string
	}{
		{"Rebroadcast: enabled address", true, false, ""},
		{"Rebroadcast: disabled address", false, true, "exists but is disabled for chain"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			config, sqlxDB := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.Database.DriverName = pgcommon.DriverPostgres

				c.EVM = nil
				// seems to be needed for config validate
				c.Insecure.OCRDevelopmentMode = nil
			})

			keyStore := cltest.NewKeyStore(t, sqlxDB)

			_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth())

			if !test.enableAddress {
				err := keyStore.Eth().Disable(testutils.Context(t), fromAddress, testutils.FixtureChainID)
				require.NoError(t, err, "failed to disable test key")
			}

			lggr := logger.TestLogger(t)

			app := mocks.NewApplication(t)
			app.On("GetDB").Maybe().Return(sqlxDB)
			app.On("GetKeyStore").Return(keyStore)
			app.On("ID").Maybe().Return(uuid.New())
			ethClient := clienttest.NewClientWithDefaultChainID(t)
			ethClient.On("Dial", mock.Anything).Return(nil)
			legacy := cltest.NewLegacyChainsWithMockChain(t, ethClient, config)

			mockRelayerChainInteroperators := &chainlinkmocks.FakeRelayerChainInteroperators{EVMChains: legacy}
			app.On("GetRelayers").Return(mockRelayerChainInteroperators).Maybe()
			ethClient.On("SendTransactionReturnCode", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(multinode.Successful, nil)

			client := cmd.Shell{
				Config:                 config,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				FallbackAPIInitializer: cltest.NewMockAPIInitializer(t),
				Runner:                 cltest.EmptyRunner{},
				Logger:                 lggr,
			}

			set := flag.NewFlagSet("test", 0)
			flagSetApplyFromAction(client.RebroadcastTransactions, set, "")

			require.NoError(t, set.Set("evmChainID", testutils.FixtureChainID.String()))
			require.NoError(t, set.Set("address", fromAddress.Hex()))
			require.NoError(t, set.Set("password", "../internal/fixtures/correct_password.txt"))
			c := cli.NewContext(nil, set, nil)

			// Run before hook to initialize components with authentication
			err := client.BeforeNode(c)
			require.NoError(t, err)

			if test.shouldError {
				require.ErrorContains(t, client.RebroadcastTransactions(c), test.errorContains)
			} else {
				app.On("GetConfig").Return(config).Once()
				require.NoError(t, client.RebroadcastTransactions(c))
			}
		})
	}
}

func TestShell_CleanupChainTables(t *testing.T) {
	// Just check if it doesn't error, command itself shouldn't be changed unless major schema changes were made.
	// It would be really hard to write a test that accounts for schema changes, so this should be enough to alarm us that something broke.
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) { c.Database.DriverName = pgcommon.DriverPostgres })
	client := cmd.Shell{
		Config: config,
		Logger: logger.TestLogger(t),
	}

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.CleanupChainTables, set, "")
	require.NoError(t, set.Set("id", testutils.FixtureChainID.String()))
	require.NoError(t, set.Set("type", "EVM"))
	// heavyweight creates test db named chainlink_test_uid, while usual naming is chainlink_test
	// CleanupChainTables handles test db name with chainlink_test, but because of heavyweight test db naming we have to set danger flag
	require.NoError(t, set.Set("danger", "true"))
	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.CleanupChainTables(c))
}

func TestShell_RemoveBlocks(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		s.Password.Keystore = models.NewSecret("dummy")
		c.EVM[0].Nodes[0].Name = ptr("fake")
		c.EVM[0].Nodes[0].HTTPURL = commonconfig.MustParseURL("http://fake.com")
		c.EVM[0].Nodes[0].WSURL = commonconfig.MustParseURL("WSS://fake.com/ws")
		// seems to be needed for config validate
		c.Insecure.OCRDevelopmentMode = nil
	})

	lggr := logger.TestLogger(t)

	app := mocks.NewApplication(t)
	app.On("GetSqlxDB").Maybe().Return(db)
	shell := cmd.Shell{
		Config:                 cfg,
		AppFactory:             cltest.InstanceAppFactory{App: app},
		FallbackAPIInitializer: cltest.NewMockAPIInitializer(t),
		Runner:                 cltest.EmptyRunner{},
		Logger:                 lggr,
	}

	t.Run("Returns error, if --start is not positive", func(t *testing.T) {
		set := flag.NewFlagSet("test", 0)
		flagSetApplyFromAction(shell.RemoveBlocks, set, "")
		require.NoError(t, set.Set("start", "0"))
		require.NoError(t, set.Set("evm-chain-id", "12"))
		c := cli.NewContext(nil, set, nil)
		err := shell.RemoveBlocks(c)
		require.ErrorContains(t, err, "Must pass a positive value in '--start' parameter")
	})
	t.Run("Returns error, if removal fails", func(t *testing.T) {
		set := flag.NewFlagSet("test", 0)
		flagSetApplyFromAction(shell.RemoveBlocks, set, "")
		require.NoError(t, set.Set("start", "10000"))
		require.NoError(t, set.Set("evm-chain-id", "12"))
		expectedError := errors.New("failed to delete log poller's data")
		app.On("DeleteLogPollerDataAfter", mock.Anything, big.NewInt(12), int64(10000)).Return(expectedError).Once()
		c := cli.NewContext(nil, set, nil)
		err := shell.RemoveBlocks(c)
		require.ErrorContains(t, err, expectedError.Error())
	})
	t.Run("Happy path", func(t *testing.T) {
		set := flag.NewFlagSet("test", 0)
		flagSetApplyFromAction(shell.RemoveBlocks, set, "")
		require.NoError(t, set.Set("start", "10000"))
		require.NoError(t, set.Set("evm-chain-id", "12"))
		app.On("DeleteLogPollerDataAfter", mock.Anything, big.NewInt(12), int64(10000)).Return(nil).Once()
		c := cli.NewContext(nil, set, nil)
		err := shell.RemoveBlocks(c)
		require.NoError(t, err)
	})
}

func TestShell_BeforeNode(t *testing.T) {
	testutils.SkipShortDB(t)
	tests := []struct {
		name            string
		pwdfile         string
		wantUnlocked    bool
		prePopulateKeys bool
	}{
		{"correct password", "../internal/fixtures/correct_password.txt", true, false},
		{"incorrect password", "../internal/fixtures/incorrect_password.txt", false, true},
		{"wrong file", "doesntexist.txt", false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, db := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.Database.DriverName = pgcommon.DriverPostgres
				c.EVM = nil
				c.Insecure.OCRDevelopmentMode = nil
			})

			// Seed key material so the wrong password actually fails decryption.
			// An empty keystore accepts any password.
			if test.prePopulateKeys {
				correctPwd, err := utils.PasswordFromFile("../internal/fixtures/correct_password.txt")
				require.NoError(t, err)
				ks := keystore.New(db, commonkeystore.FastScryptParams, logger.TestLogger(t).Infof)
				require.NoError(t, ks.Unlock(testutils.Context(t), correctPwd))
				_, err = ks.CSA().Create(testutils.Context(t))
				require.NoError(t, err)
			}

			shell := cmd.Shell{
				Config: cfg,
				KeyStoreAuthenticator: cmd.TerminalKeyStoreAuthenticator{
					Prompter: &cltest.MockCountingPrompter{T: t, NotTerminal: true},
				},
				Logger: logger.TestLogger(t),
			}
			// Reset for test isolation
			defer resetShellForTest(&shell)

			set := flag.NewFlagSet("test", 0)
			flagSetApplyFromAction(shell.RunNode, set, "")
			require.NoError(t, set.Set("password", test.pwdfile))

			c := cli.NewContext(nil, set, nil)

			// Create full CLI app and run the Before hook first
			app := cmd.NewApp(&shell)
			err := app.Before(c)
			if err != nil && test.wantUnlocked {
				t.Fatalf("CLI Before hook failed: %v", err)
			}

			// Run before hook to initialize components with authentication
			err = shell.BeforeNode(c)

			if test.wantUnlocked {
				require.NoError(t, err)
				// Verify that shell components were initialized
				assert.NotNil(t, shell.KeyStore)
				assert.NotNil(t, shell.DS)
				assert.NotNil(t, shell.LDB)

				// Verify keystore is unlocked by checking if we can access keys
				keys, keysErr := shell.KeyStore.CSA().GetAll()
				require.NoError(t, keysErr)
				assert.NotEmpty(t, keys)
			} else {
				require.Error(t, err)
			}
			// Clean up database if it was opened
			cleanupShell(t, &shell, c)
		})
	}
}

func TestShell_BeholderLifecycle(t *testing.T) {
	testutils.SkipShortDB(t)
	rawDBURL, ok := os.LookupEnv("CL_DATABASE_URL")
	if !ok {
		t.Skip("CL_DATABASE_URL is required for this test")
	}
	parsedDBURL, err := url.Parse(rawDBURL)
	if err != nil || parsedDBURL.Path == "" {
		t.Skip("CL_DATABASE_URL must include a database path for this test")
	}

	newShell := func(t *testing.T, overrideFn func(c *chainlink.Config, s *chainlink.Secrets), setOtelCore func(zapcore.Core)) (*cmd.Shell, *cli.Context) {
		t.Helper()

		cfg, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			c.Database.DriverName = pgcommon.DriverPostgres
			c.EVM = nil
			c.Insecure.OCRDevelopmentMode = nil
			if overrideFn != nil {
				overrideFn(c, s)
			}
		})

		shell := &cmd.Shell{
			Config: cfg,
			KeyStoreAuthenticator: cmd.TerminalKeyStoreAuthenticator{
				Prompter: &cltest.MockCountingPrompter{T: t, NotTerminal: true},
			},
			Logger:      logger.TestLogger(t),
			SetOtelCore: setOtelCore,
		}
		t.Cleanup(func() { resetShellForTest(shell) })

		set := flag.NewFlagSet("test", 0)
		flagSetApplyFromAction(shell.RunNode, set, "")
		require.NoError(t, set.Set("password", "../internal/fixtures/correct_password.txt"))

		ctx := cli.NewContext(nil, set, nil)
		cliApp := cmd.NewApp(shell)
		require.NoError(t, cliApp.Before(ctx))

		return shell, ctx
	}

	t.Run("telemetry disabled assigns noop beholder client", func(t *testing.T) {
		shell, c := newShell(t, nil, nil)
		require.NoError(t, shell.BeforeNode(c))
		require.NotNil(t, shell.BeholderClient, "BeholderClient should be a no-op client when telemetry is disabled")
		require.NoError(t, shell.AfterNode(c))
	})

	t.Run("telemetry enabled starts and closes beholder", func(t *testing.T) {
		shell, c := newShell(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			trueVal := true
			c.Telemetry.Enabled = &trueVal
			endpoint := "localhost:4317"
			c.Telemetry.Endpoint = &endpoint
			c.Telemetry.InsecureConnection = &trueVal
		}, nil)
		require.NoError(t, shell.BeforeNode(c))
		require.NotNil(t, shell.BeholderClient, "BeholderClient should be set when telemetry is enabled")
		assert.NoError(t, shell.BeholderClient.Ready())
		require.NoError(t, shell.AfterNode(c))
	})

	t.Run("after node is idempotent", func(t *testing.T) {
		shell, c := newShell(t, nil, nil)
		require.NoError(t, shell.BeforeNode(c))
		require.NoError(t, shell.AfterNode(c))
		require.NoError(t, shell.AfterNode(c))
	})

	t.Run("log streaming sets otel core", func(t *testing.T) {
		var setOtelCoreCalls int
		shell, c := newShell(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			trueVal := true
			c.Telemetry.Enabled = &trueVal
			endpoint := "localhost:4317"
			c.Telemetry.Endpoint = &endpoint
			c.Telemetry.InsecureConnection = &trueVal
			c.Telemetry.LogStreamingEnabled = &trueVal
		}, func(core zapcore.Core) {
			require.NotNil(t, core)
			setOtelCoreCalls++
		})

		require.NoError(t, shell.BeforeNode(c))
		assert.Equal(t, 1, setOtelCoreCalls)
		require.NoError(t, shell.AfterNode(c))
	})

	t.Run("log streaming fails when SetOtelCore is nil", func(t *testing.T) {
		shell, c := newShell(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			trueVal := true
			c.Telemetry.Enabled = &trueVal
			endpoint := "localhost:4317"
			c.Telemetry.Endpoint = &endpoint
			c.Telemetry.InsecureConnection = &trueVal
			c.Telemetry.LogStreamingEnabled = &trueVal
		}, nil) // SetOtelCore intentionally nil
		err := shell.BeforeNode(c)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "SetOtelCore is nil")
	})
}

func TestShell_AfterNode_NilBeholderClient(t *testing.T) {
	shell := cmd.Shell{
		LDB:            stubLockedDB{},
		Logger:         logger.TestLogger(t),
		BeholderClient: beholder.NewNoopClient(),
	}
	assert.NotPanics(t, func() {
		_ = shell.AfterNode(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))
	})
}

func TestShell_RunNode_WithBeforeNode(t *testing.T) {
	tests := []struct {
		name        string
		pwdfile     string
		expectStart bool
	}{
		{"correct password", "../internal/fixtures/correct_password.txt", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				s.Password.Keystore = models.NewSecret("dummy")
				c.EVM[0].Nodes[0].Name = ptr("fake")
				c.EVM[0].Nodes[0].HTTPURL = commonconfig.MustParseURL("http://fake.com")
				c.EVM[0].Nodes[0].WSURL = commonconfig.MustParseURL("WSS://fake.com/ws")
				// seems to be needed for config validate
				c.Insecure.OCRDevelopmentMode = nil
			})

			db := pgtest.NewSqlxDB(t)
			keyStore := cltest.NewKeyStore(t, db)
			authProviderORM := localauth.NewORM(db, time.Minute, logger.TestLogger(t), audit.NoopLogger)

			testRelayers := genTestEVMRelayers(t, cfg, db, keyStore.Eth(), &keystore.CSASigner{CSA: keyStore.CSA()})

			// Purge fixture users to test assumption of single admin
			pgtest.MustExec(t, db, "DELETE FROM users;")

			app := mocks.NewApplication(t)
			app.On("AuthenticationProvider").Return(authProviderORM).Maybe()
			app.On("BasicAdminUsersORM").Return(authProviderORM).Maybe()
			app.On("GetKeyStore").Return(keyStore).Maybe()
			app.On("GetRelayers").Return(testRelayers).Maybe()
			app.On("Start", mock.Anything).Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)
			app.On("ID").Maybe().Return(uuid.New())

			ethClient := clienttest.NewClient(t)
			ethClient.On("Dial", mock.Anything).Return(nil).Maybe()
			ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil).Maybe()

			cltest.MustInsertRandomKey(t, keyStore.Eth())
			apiPrompt := cltest.NewMockAPIInitializer(t)

			shell := cmd.Shell{
				Config:                 cfg,
				FallbackAPIInitializer: apiPrompt,
				Runner:                 cltest.EmptyRunner{},
				AppFactory:             cltest.InstanceAppFactory{App: app},
				KeyStoreAuthenticator: cmd.TerminalKeyStoreAuthenticator{
					Prompter: &cltest.MockCountingPrompter{T: t, NotTerminal: true},
				},
				Logger: logger.TestLogger(t),
			}
			// Reset for test isolation
			defer resetShellForTest(&shell)

			set := flag.NewFlagSet("test", 0)
			flagSetApplyFromAction(shell.RunNode, set, "")
			require.NoError(t, set.Set("password", test.pwdfile))

			c := cli.NewContext(nil, set, nil)

			// First initialize components (this includes authentication)
			cliApp := cmd.NewApp(&shell)
			err := cliApp.Before(c)
			require.NoError(t, err)

			err = shell.BeforeNode(c)

			if test.expectStart {
				require.NoError(t, err, "BeforeNode should succeed")
				// Verify components are initialized
				assert.NotNil(t, shell.KeyStore)
				assert.NotNil(t, shell.DS)
				assert.NotNil(t, shell.LDB)

				// Now test RunNode with pre-authenticated keystore
				// Note: RunNode will start the app but we expect it to work since keystore is authenticated
				err = shell.RunNode(c)
				require.NoError(t, err, "RunNode should succeed with authenticated keystore")
				assert.Equal(t, 1, apiPrompt.Count, "API should be initialized")
			} else {
				require.Error(t, err, "BeforeNode should fail with incorrect password")
				// Don't test RunNode if BeforeNode failed
			}
			// Clean up database if it was opened
			cleanupShell(t, &shell, c)
		})
	}
}

func cleanupShell(t *testing.T, shell *cmd.Shell, c *cli.Context) {
	t.Helper()
	if shell.LDB == nil {
		return
	}
	if shell.BeholderClient == nil {
		shell.BeholderClient = beholder.NewNoopClient()
	}
	require.NoError(t, shell.AfterNode(c))
}
