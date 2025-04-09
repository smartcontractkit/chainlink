package cmd_test

import (
	"errors"
	"flag"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	pgcommon "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink-framework/multinode"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"

	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
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
	"github.com/smartcontractkit/chainlink/v2/core/sessions/localauth"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func genTestEVMRelayers(t *testing.T, cfg chainlink.GeneralConfig, ds sqlutil.DataSource, ethKeystore keystore.Eth, csaKeystore core.Keystore) *chainlink.CoreRelayerChainInteroperators {
	lggr := logger.TestLogger(t)
	f := chainlink.RelayerFactory{
		Logger:               lggr,
		LoopRegistry:         plugins.NewLoopRegistry(lggr, cfg.Database(), cfg.Tracing(), cfg.Telemetry(), nil, ""),
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

func TestShell_RunNodeWithPasswords(t *testing.T) {
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
				c.EVM[0].Nodes[0].HTTPURL = commonconfig.MustParseURL("http://fake.com")
				c.EVM[0].Nodes[0].WSURL = commonconfig.MustParseURL("WSS://fake.com/ws")
				// seems to be needed for config validate
				c.Insecure.OCRDevelopmentMode = nil
			})
			db := pgtest.NewSqlxDB(t)
			keyStore := cltest.NewKeyStore(t, db)
			authProviderORM := localauth.NewORM(db, time.Minute, logger.TestLogger(t), audit.NoopLogger)

			testRelayers := genTestEVMRelayers(t, cfg, db, keyStore.Eth(), &keystore.CSASigner{CSA: keyStore.CSA()})

			// Purge the fixture users to test assumption of single admin
			// initialUser user created above
			pgtest.MustExec(t, db, "DELETE FROM users;")

			app := mocks.NewApplication(t)
			app.On("AuthenticationProvider").Return(authProviderORM).Maybe()
			app.On("BasicAdminUsersORM").Return(authProviderORM).Maybe()
			app.On("GetKeyStore").Return(keyStore).Maybe()
			app.On("GetRelayers").Return(testRelayers).Maybe()
			app.On("Start", mock.Anything).Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)
			app.On("ID").Maybe().Return(uuid.New())

			ethClient := evmtest.NewEthClientMock(t)
			ethClient.On("Dial", mock.Anything).Return(nil).Maybe()
			ethClient.On("BalanceAt", mock.Anything, mock.Anything, mock.Anything).Return(big.NewInt(10), nil).Maybe()

			cltest.MustInsertRandomKey(t, keyStore.Eth())
			apiPrompt := cltest.NewMockAPIInitializer(t)

			client := cmd.Shell{
				Config:                 cfg,
				FallbackAPIInitializer: apiPrompt,
				Runner:                 cltest.EmptyRunner{},
				AppFactory:             cltest.InstanceAppFactoryWithKeystoreMock{App: app},
				Logger:                 logger.TestLogger(t),
			}

			set := flag.NewFlagSet("test", 0)
			flagSetApplyFromAction(client.RunNode, set, "")

			require.NoError(t, set.Set("password", test.pwdfile))

			c := cli.NewContext(nil, set, nil)

			run := func() error {
				cli := cmd.NewApp(&client)
				if err := cli.Before(c); err != nil {
					return err
				}
				return client.RunNode(c)
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
			assert.NoError(t, os.MkdirAll(cfg.Dir, os.FileMode(0700)))

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

	txStore := cltest.NewTestTxStore(t, sqlxDB)
	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 7, 42, fromAddress)

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

	for i := beginningNonce; i <= endingNonce; i++ {
		n := i
		ethClient.On("SendTransactionReturnCode", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == n
		}), mock.Anything).Once().Return(multinode.Successful, nil)
	}

	assert.NoError(t, c.RebroadcastTransactions(ctx))
}

func TestShell_RebroadcastTransactions_OutsideRange_Txm(t *testing.T) {
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

			txStore := cltest.NewTestTxStore(t, sqlxDB)
			cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, int64(test.nonce), 42, fromAddress)

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
