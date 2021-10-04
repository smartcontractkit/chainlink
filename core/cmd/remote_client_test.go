package cmd_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"
)

var (
	nilContext = cli.NewContext(nil, nil, nil)
)

type startOptions struct {
	// Set the config options
	SetConfig func(cfg *configtest.TestGeneralConfig)
	// Use to set up mocks on the app
	FlagsAndDeps []interface{}
	// Add a key on start up
	WithKey bool
}

func startNewApplication(t *testing.T, setup ...func(opts *startOptions)) *cltest.TestApplication {
	t.Helper()

	sopts := &startOptions{
		FlagsAndDeps: []interface{}{},
	}
	for _, fn := range setup {
		fn(sopts)
	}

	// Setup config
	config := cltest.NewTestGeneralConfig(t)
	config.Overrides.SetDefaultHTTPTimeout(30 * time.Millisecond)
	config.Overrides.DefaultMaxHTTPAttempts = null.IntFrom(1)

	// Generally speaking, most tests that use startNewApplication don't
	// actually need ChainSets loaded. We can greatly reduce test
	// overhead by setting EVM_DISABLED here. If you need EVM interactions in
	// your tests, you can manually override and turn it on using
	// withConfigSet.
	config.Overrides.EVMDisabled = null.BoolFrom(true)

	if sopts.SetConfig != nil {
		sopts.SetConfig(config)
	}

	app := cltest.NewApplicationWithConfigAndKey(t, config, sopts.FlagsAndDeps...)
	require.NoError(t, app.Start())

	return app
}

// withConfig is a function option which sets config on the app
func withConfigSet(cfgSet func(*configtest.TestGeneralConfig)) func(opts *startOptions) {
	return func(opts *startOptions) {
		opts.SetConfig = cfgSet
	}
}

func withMocks(mks ...interface{}) func(opts *startOptions) {
	return func(opts *startOptions) {
		opts.FlagsAndDeps = mks
	}
}

func withKey() func(opts *startOptions) {
	return func(opts *startOptions) {
		opts.WithKey = true
	}
}

func newEthMock(t *testing.T) (*mocks.Client, func()) {
	t.Helper()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)

	return ethClient, assertMocksCalled
}

func keyNameForTest(t *testing.T) string {
	return fmt.Sprintf("%s_test_key.json", t.Name())
}

func deleteKeyExportFile(t *testing.T) {
	keyName := keyNameForTest(t)
	err := os.Remove(keyName)
	if err == nil || os.IsNotExist(err) {
		return
	} else {
		require.NoError(t, err)
	}
}

func TestClient_ReplayBlocks(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t,
		withConfigSet(func(c *configtest.TestGeneralConfig) {
			c.Overrides.EVMDisabled = null.BoolFrom(false)
			c.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(false)
			c.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
			c.Overrides.GlobalGasEstimatorMode = null.StringFrom("FixedPrice")
		}))
	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("flagset", 0)
	set.Int64("block-number", 42, "")
	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.ReplayFromBlock(c))
}

func TestClient_CreateExternalInitiator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{"create external initiator", []string{"exi", "http://testing.com/external_initiators"}},
		{"create external initiator w/ query params", []string{"exiqueryparams", "http://testing.com/external_initiators?query=param"}},
		{"create external initiator w/o url", []string{"exi_no_url"}},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			app := startNewApplication(t)
			client, _ := app.NewClientAndRenderer()

			set := flag.NewFlagSet("create", 0)
			assert.NoError(t, set.Parse(test.args))
			c := cli.NewContext(nil, set, nil)

			err := client.CreateExternalInitiator(c)
			require.NoError(t, err)

			var exi bridges.ExternalInitiator
			err = app.GetDB().Where("name = ?", test.args[0]).Find(&exi).Error
			require.NoError(t, err)

			if len(test.args) > 1 {
				assert.Equal(t, test.args[1], exi.URL.String())
			}
		})
	}
}

func TestClient_CreateExternalInitiator_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{"no arguments", []string{}},
		{"too many arguments", []string{"bitcoin", "https://valid.url", "extra arg"}},
		{"invalid url", []string{"bitcoin", "not a url"}},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			app := startNewApplication(t)
			client, _ := app.NewClientAndRenderer()

			initialExis := len(cltest.AllExternalInitiators(t, app.GetDB()))

			set := flag.NewFlagSet("create", 0)
			assert.NoError(t, set.Parse(test.args))
			c := cli.NewContext(nil, set, nil)

			err := client.CreateExternalInitiator(c)
			assert.Error(t, err)

			exis := cltest.AllExternalInitiators(t, app.GetDB())
			assert.Len(t, exis, initialExis)
		})
	}
}

func TestClient_DestroyExternalInitiator(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	token := auth.NewToken()
	exi, err := bridges.NewExternalInitiator(token,
		&bridges.ExternalInitiatorRequest{Name: "name"},
	)
	require.NoError(t, err)
	err = app.BridgeORM().CreateExternalInitiator(exi)
	require.NoError(t, err)

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{exi.Name})
	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.DeleteExternalInitiator(c))
	assert.Empty(t, r.Renders)
}

func TestClient_DestroyExternalInitiator_NotFound(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"bogus-ID"})
	c := cli.NewContext(nil, set, nil)
	assert.Error(t, client.DeleteExternalInitiator(c))
	assert.Empty(t, r.Renders)
}

func TestClient_RemoteLogin(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t, withConfigSet(func(c *configtest.TestGeneralConfig) {
		c.Overrides.AdminCredentialsFile = null.StringFrom("")
	}))

	tests := []struct {
		name, file string
		email, pwd string
		wantError  bool
	}{
		{"success prompt", "", cltest.APIEmail, cltest.Password, false},
		{"success file", "../internal/fixtures/apicredentials", "", "", false},
		{"failure prompt", "", "wrong@email.com", "wrongpwd", true},
		{"failure file", "/tmp/doesntexist", "", "", true},
		{"failure file w correct prompt", "/tmp/doesntexist", cltest.APIEmail, cltest.Password, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			enteredStrings := []string{test.email, test.pwd}
			prompter := &cltest.MockCountingPrompter{EnteredStrings: enteredStrings}
			client := app.NewAuthenticatingClient(prompter)

			set := flag.NewFlagSet("test", 0)
			set.String("file", test.file, "")
			c := cli.NewContext(nil, set, nil)

			err := client.RemoteLogin(c)
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_ChangePassword(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)

	enteredStrings := []string{cltest.APIEmail, cltest.Password}
	prompter := &cltest.MockCountingPrompter{EnteredStrings: enteredStrings}

	client := app.NewAuthenticatingClient(prompter)
	otherClient := app.NewAuthenticatingClient(prompter)

	set := flag.NewFlagSet("test", 0)
	set.String("file", "../internal/fixtures/apicredentials", "")
	c := cli.NewContext(nil, set, nil)
	err := client.RemoteLogin(c)
	require.NoError(t, err)

	err = otherClient.RemoteLogin(c)
	require.NoError(t, err)

	client.ChangePasswordPrompter = cltest.MockChangePasswordPrompter{
		UpdatePasswordRequest: web.UpdatePasswordRequest{
			OldPassword: cltest.Password,
			NewPassword: "_p4SsW0rD1!@#",
		},
	}
	err = client.ChangePassword(cli.NewContext(nil, nil, nil))
	assert.NoError(t, err)

	// otherClient should now be logged out
	err = otherClient.IndexBridges(c)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Unauthorized")
}

func TestClient_SetDefaultGasPrice(t *testing.T) {
	t.Parallel()

	ethMock, assertMocksCalled := newEthMock(t)
	defer assertMocksCalled()
	app := startNewApplication(t,
		withKey(),
		withMocks(ethMock),
		withConfigSet(func(c *configtest.TestGeneralConfig) {
			c.Overrides.EVMDisabled = null.BoolFrom(false)
			c.Overrides.GlobalEvmNonceAutoSync = null.BoolFrom(false)
			c.Overrides.GlobalBalanceMonitorEnabled = null.BoolFrom(false)
		}),
	)
	client, _ := app.NewClientAndRenderer()

	t.Run("without specifying chain id setting value", func(t *testing.T) {
		set := flag.NewFlagSet("setgasprice", 0)
		set.Parse([]string{"8616460799"})

		c := cli.NewContext(nil, set, nil)

		assert.NoError(t, client.SetEvmGasPriceDefault(c))
		ch, err := app.GetChainSet().Default()
		require.NoError(t, err)
		cfg := ch.Config()
		assert.Equal(t, big.NewInt(8616460799), cfg.EvmGasPriceDefault())

		client, _ = app.NewClientAndRenderer()
		set = flag.NewFlagSet("setgasprice", 0)
		set.String("amount", "", "")
		set.Bool("gwei", true, "")
		set.Parse([]string{"-gwei", "861.6460799"})

		c = cli.NewContext(nil, set, nil)
		assert.NoError(t, client.SetEvmGasPriceDefault(c))
		assert.Equal(t, big.NewInt(861646079900), cfg.EvmGasPriceDefault())
	})

	t.Run("specifying wrong chain id", func(t *testing.T) {
		set := flag.NewFlagSet("setgasprice", 0)
		set.String("evmChainID", "", "")
		set.Parse([]string{"-evmChainID", "985435435435", "8616460799"})

		c := cli.NewContext(nil, set, nil)

		err := client.SetEvmGasPriceDefault(c)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "evmChainID does not match any local chains")

		ch, err := app.GetChainSet().Default()
		require.NoError(t, err)
		cfg := ch.Config()
		assert.Equal(t, big.NewInt(861646079900), cfg.EvmGasPriceDefault())
	})

	t.Run("specifying correct chain id", func(t *testing.T) {
		set := flag.NewFlagSet("setgasprice", 0)
		set.String("evmChainID", "", "")
		set.Parse([]string{"-evmChainID", "0", "12345678900"})

		c := cli.NewContext(nil, set, nil)

		assert.NoError(t, client.SetEvmGasPriceDefault(c))
		ch, err := app.GetChainSet().Default()
		require.NoError(t, err)
		cfg := ch.Config()

		assert.Equal(t, big.NewInt(12345678900), cfg.EvmGasPriceDefault())
	})
}

func TestClient_GetConfiguration(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()
	cfg := app.GetConfig()

	assert.NoError(t, client.GetConfiguration(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))

	cp := *r.Renders[0].(*presenters.ConfigPrinter)
	assert.Equal(t, cp.EnvPrinter.BridgeResponseURL, cfg.BridgeResponseURL().String())
	assert.Equal(t, cp.EnvPrinter.DefaultChainID, cfg.DefaultChainID().String())
	assert.Equal(t, cp.EnvPrinter.Dev, cfg.Dev())
	assert.Equal(t, cp.EnvPrinter.LogLevel, cfg.LogLevel())
	assert.Equal(t, cp.EnvPrinter.LogSQLStatements, cfg.LogSQLStatements())
	assert.Equal(t, cp.EnvPrinter.RootDir, cfg.RootDir())
	assert.Equal(t, cp.EnvPrinter.SessionTimeout, cfg.SessionTimeout())
}

func TestClient_RunOCRJob_HappyPath(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t, withConfigSet(func(c *configtest.TestGeneralConfig) {
		c.Overrides.EVMDisabled = null.BoolFrom(false)
		c.Overrides.GlobalGasEstimatorMode = null.StringFrom("FixedPrice")
	}))
	client, _ := app.NewClientAndRenderer()

	app.KeyStore.OCR().Add(cltest.DefaultOCRKey)
	app.KeyStore.P2P().Add(cltest.DefaultP2PKey)

	_, bridge := cltest.NewBridgeType(t, "voter_turnout", "http://blah.com")
	require.NoError(t, app.GetDB().Create(bridge).Error)
	_, bridge2 := cltest.NewBridgeType(t, "election_winner", "http://blah.com")
	require.NoError(t, app.GetDB().Create(bridge2).Error)

	var ocrJobSpecFromFile job.Job
	tree, err := toml.LoadFile("../testdata/tomlspecs/oracle-spec.toml")
	require.NoError(t, err)
	err = tree.Unmarshal(&ocrJobSpecFromFile)
	require.NoError(t, err)
	var ocrSpec job.OffchainReportingOracleSpec
	err = tree.Unmarshal(&ocrSpec)
	require.NoError(t, err)
	ocrJobSpecFromFile.OffchainreportingOracleSpec = &ocrSpec
	key, _ := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())
	ocrJobSpecFromFile.OffchainreportingOracleSpec.TransmitterAddress = &key.Address

	jb, _ := app.AddJobV2(context.Background(), ocrJobSpecFromFile, null.String{})

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{strconv.FormatInt(int64(jb.ID), 10)})
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.RemoteLogin(c))
	require.NoError(t, client.TriggerPipelineRun(c))
}

func TestClient_RunOCRJob_MissingJobID(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.RemoteLogin(c))
	assert.EqualError(t, client.TriggerPipelineRun(c), "Must pass the job id to trigger a run")
}

func TestClient_RunOCRJob_JobNotFound(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"1"})
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.RemoteLogin(c))
	assert.EqualError(t, client.TriggerPipelineRun(c), "parseResponse error: Error; job ID 1: record not found")
}

func TestClient_AutoLogin(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)

	user := cltest.MustRandomUser()
	require.NoError(t, app.SessionORM().CreateUser(&user))

	sr := sessions.SessionRequest{
		Email:    user.Email,
		Password: cltest.Password,
	}
	client, _ := app.NewClientAndRenderer()
	client.CookieAuthenticator = cmd.NewSessionCookieAuthenticator(app.GetConfig(), &cmd.MemoryCookieStore{})
	client.HTTP = cmd.NewAuthenticatedHTTPClient(app.Config, client.CookieAuthenticator, sr)

	fs := flag.NewFlagSet("", flag.ExitOnError)
	err := client.ListJobs(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)

	// Expire the session and then try again
	require.NoError(t, app.GetDB().Exec("delete from sessions;").Error)
	err = client.ListJobs(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)
}

func TestClient_AutoLogin_AuthFails(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)

	user := cltest.MustRandomUser()
	require.NoError(t, app.SessionORM().CreateUser(&user))

	sr := sessions.SessionRequest{
		Email:    user.Email,
		Password: cltest.Password,
	}
	client, _ := app.NewClientAndRenderer()
	client.CookieAuthenticator = FailingAuthenticator{}
	client.HTTP = cmd.NewAuthenticatedHTTPClient(app.Config, client.CookieAuthenticator, sr)

	fs := flag.NewFlagSet("", flag.ExitOnError)
	err := client.ListJobs(cli.NewContext(nil, fs, nil))
	require.Error(t, err)
}

type FailingAuthenticator struct{}

func (FailingAuthenticator) Cookie() (*http.Cookie, error) {
	return &http.Cookie{}, nil
}

// Authenticate retrieves a session ID via a cookie and saves it to disk.
func (FailingAuthenticator) Authenticate(sessionRequest sessions.SessionRequest) (*http.Cookie, error) {
	return nil, errors.New("no luck")
}

func TestClient_SetLogConfig(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	logLevel := "warn"
	set := flag.NewFlagSet("loglevel", 0)
	set.String("level", logLevel, "")
	c := cli.NewContext(nil, set, nil)

	err := client.SetLogLevel(c)
	require.NoError(t, err)
	assert.Equal(t, logLevel, app.Config.LogLevel().String())

	sqlEnabled := true
	set = flag.NewFlagSet("logsql", 0)
	set.Bool("enable", sqlEnabled, "")
	c = cli.NewContext(nil, set, nil)

	err = client.SetLogSQL(c)
	assert.NoError(t, err)
	assert.Equal(t, sqlEnabled, app.Config.LogSQLStatements())

	sqlEnabled = false
	set = flag.NewFlagSet("logsql", 0)
	set.Bool("disable", true, "")
	c = cli.NewContext(nil, set, nil)

	err = client.SetLogSQL(c)
	assert.NoError(t, err)
	assert.Equal(t, sqlEnabled, app.Config.LogSQLStatements())
}

func TestClient_SetPkgLogLevel(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	logPkg := logger.HeadTracker
	logLevel := "warn"
	set := flag.NewFlagSet("logpkg", 0)
	set.String("pkg", logPkg, "")
	set.String("level", logLevel, "")
	c := cli.NewContext(nil, set, nil)

	err := client.SetLogPkg(c)
	require.NoError(t, err)

	level, ok := logger.NewORM(app.GetDB()).GetServiceLogLevel(logPkg)
	require.True(t, ok)
	assert.Equal(t, logLevel, level)
}
