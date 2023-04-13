package cmd_test

import (
	"flag"
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	cmdMocks "github.com/smartcontractkit/chainlink/v2/core/cmd/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/mocks"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	gethTypes "github.com/ethereum/go-ethereum/core/types"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

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

			app := mocks.NewApplication(t)
			app.On("SessionORM").Return(sessionORM).Maybe()
			app.On("GetKeyStore").Return(keyStore).Maybe()
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
			cltest.FlagSetApplyFromAction(client.RunNode, set, "")

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

			app := mocks.NewApplication(t)
			app.On("SessionORM").Return(sessionORM)
			app.On("GetKeyStore").Return(keyStore)
			app.On("GetChains").Return(chainlink.Chains{EVM: cltest.NewChainSetMockWithOneChain(t, ethClient, evmtest.NewChainScopedConfig(t, cfg))}).Maybe()
			app.On("Start", mock.Anything).Maybe().Return(nil)
			app.On("Stop").Maybe().Return(nil)
			app.On("ID").Maybe().Return(uuid.NewV4())

			prompter := cmdMocks.NewPrompter(t)

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
			cltest.FlagSetApplyFromAction(client.RunNode, set, "")

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
	config, sqlxDB := heavyweight.FullTestDBV2(t, "rebroadcasttransactions", func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Database.Dialect = dialects.Postgres
		// evm config is used in this test. but if set, it must be pass config validation.
		// simplest to make it nil
		c.EVM = nil
	})
	keyStore := cltest.NewKeyStore(t, sqlxDB, config)
	_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth(), 0)

	txStore := cltest.NewTxStore(t, sqlxDB, config)
	cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, 7, 42, fromAddress)

	app := mocks.NewApplication(t)
	app.On("GetSqlxDB").Return(sqlxDB)
	app.On("GetKeyStore").Return(keyStore)
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

	beginningNonce := uint64(7)
	endingNonce := uint64(10)
	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.RebroadcastTransactions, set, "")

	require.NoError(t, set.Set("beginningNonce", strconv.FormatUint(beginningNonce, 10)))
	require.NoError(t, set.Set("endingNonce", strconv.FormatUint(endingNonce, 10)))
	require.NoError(t, set.Set("gasPriceWei", "100000000000"))
	require.NoError(t, set.Set("gasLimit", "3000000"))
	require.NoError(t, set.Set("address", fromAddress.Hex()))
	require.NoError(t, set.Set("password", "../internal/fixtures/correct_password.txt"))

	c := cli.NewContext(nil, set, nil)

	for i := beginningNonce; i <= endingNonce; i++ {
		n := i
		ethClient.On("SendTransaction", mock.Anything, mock.MatchedBy(func(tx *gethTypes.Transaction) bool {
			return tx.Nonce() == n
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
			// Use the non-transactional db for this test because we need to
			// test multiple connections to the database, and changes made within
			// the transaction cannot be seen from another connection.
			config, sqlxDB := heavyweight.FullTestDBV2(t, "rebroadcasttransactions_outsiderange", func(c *chainlink.Config, s *chainlink.Secrets) {
				c.Database.Dialect = dialects.Postgres
				// evm config is used in this test. but if set, it must be pass config validation.
				// simplest to make it nil
				c.EVM = nil
			})

			keyStore := cltest.NewKeyStore(t, sqlxDB, config)

			_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth(), 0)

			txStore := cltest.NewTxStore(t, sqlxDB, config)
			cltest.MustInsertConfirmedEthTxWithLegacyAttempt(t, txStore, int64(test.nonce), 42, fromAddress)

			app := mocks.NewApplication(t)
			app.On("GetSqlxDB").Return(sqlxDB)
			app.On("GetKeyStore").Return(keyStore)
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
			cltest.FlagSetApplyFromAction(client.RebroadcastTransactions, set, "")

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

func TestClient_RebroadcastTransactions_AddressCheck(t *testing.T) {
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

			config, sqlxDB := heavyweight.FullTestDBV2(t, "rebroadcasttransactions_outsiderange", func(c *chainlink.Config, s *chainlink.Secrets) {
				c.Database.Dialect = dialects.Postgres
				c.EVM = nil
			})

			keyStore := cltest.NewKeyStore(t, sqlxDB, config)

			_, fromAddress := cltest.MustInsertRandomKey(t, keyStore.Eth(), 0)

			if !test.enableAddress {
				keyStore.Eth().Disable(fromAddress, big.NewInt(0))
			}

			app := mocks.NewApplication(t)
			app.On("GetSqlxDB").Maybe().Return(sqlxDB)
			app.On("GetKeyStore").Return(keyStore)
			app.On("ID").Maybe().Return(uuid.NewV4())
			ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
			ethClient.On("Dial", mock.Anything).Return(nil)
			app.On("GetChains").Return(chainlink.Chains{EVM: cltest.NewChainSetMockWithOneChain(t, ethClient, evmtest.NewChainScopedConfig(t, config))}).Maybe()

			ethClient.On("SendTransaction", mock.Anything, mock.Anything).Maybe().Return(nil)

			lggr := logger.TestLogger(t)

			client := cmd.Client{
				Config:                 config,
				AppFactory:             cltest.InstanceAppFactory{App: app},
				FallbackAPIInitializer: cltest.NewMockAPIInitializer(t),
				Runner:                 cltest.EmptyRunner{},
				Logger:                 lggr,
			}

			set := flag.NewFlagSet("test", 0)
			cltest.FlagSetApplyFromAction(client.RebroadcastTransactions, set, "")

			require.NoError(t, set.Set("address", fromAddress.Hex()))
			require.NoError(t, set.Set("password", "../internal/fixtures/correct_password.txt"))
			c := cli.NewContext(nil, set, nil)
			if test.shouldError {
				require.ErrorContains(t, client.RebroadcastTransactions(c), test.errorContains)
			} else {
				require.NoError(t, client.RebroadcastTransactions(c))
			}

		})
	}
}
