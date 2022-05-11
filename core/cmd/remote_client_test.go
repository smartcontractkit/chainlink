package cmd_test

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/sessions"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/core/web"
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
	config.Overrides.SetHTTPServerWriteTimeout(10 * time.Second)

	// Generally speaking, most tests that use startNewApplication don't
	// actually need ChainSets loaded. We can greatly reduce test
	// overhead by disabling EVM here. If you need EVM interactions in
	// your tests, you can manually override and turn it on using
	// withConfigSet.
	config.Overrides.EVMEnabled = null.BoolFrom(false)
	config.Overrides.P2PEnabled = null.BoolFrom(false)

	if sopts.SetConfig != nil {
		sopts.SetConfig(config)
	}

	app := cltest.NewApplicationWithConfigAndKey(t, config, sopts.FlagsAndDeps...)
	require.NoError(t, app.Start(testutils.Context(t)))

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

func newEthMock(t *testing.T) *evmmocks.Client {
	t.Helper()
	return cltest.NewEthMocksWithStartupAssertions(t)
}

func newEthMockWithTransactionsOnBlocksAssertions(t *testing.T) *evmmocks.Client {
	t.Helper()

	return cltest.NewEthMocksWithTransactionsOnBlocksAssertions(t)
}

func keyNameForTest(t *testing.T) string {
	return fmt.Sprintf("%s_test_key.json", t.Name())
}

func deleteKeyExportFile(t *testing.T) {
	keyName := keyNameForTest(t)
	err := os.Remove(keyName)
	if err == nil || os.IsNotExist(err) {
		return
	}
	require.NoError(t, err)
}

func TestClient_ReplayBlocks(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t,
		withConfigSet(func(c *configtest.TestGeneralConfig) {
			c.Overrides.EVMEnabled = null.BoolFrom(true)
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
			err = app.GetSqlxDB().Get(&exi, `SELECT * FROM external_initiators WHERE name = $1`, test.args[0])
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

			initialExis := len(cltest.AllExternalInitiators(t, app.GetSqlxDB()))

			set := flag.NewFlagSet("create", 0)
			assert.NoError(t, set.Parse(test.args))
			c := cli.NewContext(nil, set, nil)

			err := client.CreateExternalInitiator(c)
			assert.Error(t, err)

			exis := cltest.AllExternalInitiators(t, app.GetSqlxDB())
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

	app := startNewApplication(t)

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
			prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}
			client := app.NewAuthenticatingClient(prompter)

			set := flag.NewFlagSet("test", 0)
			set.String("file", test.file, "")
			set.Bool("bypass-version-check", true, "")
			set.String("admin-credentials-file", "", "")
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

func TestClient_RemoteBuildCompatibility(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	enteredStrings := []string{cltest.APIEmail, cltest.Password}
	prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: append(enteredStrings, enteredStrings...)}
	client := app.NewAuthenticatingClient(prompter)

	remoteVersion, remoteSha := "test"+static.Version, "abcd"+static.Sha
	client.HTTP = &mockHTTPClient{client.HTTP, remoteVersion, remoteSha}

	expErr := cmd.ErrIncompatible{
		CLIVersion:    static.Version,
		CLISha:        static.Sha,
		RemoteVersion: remoteVersion,
		RemoteSha:     remoteSha,
	}.Error()

	// Fails without bypass
	set := flag.NewFlagSet("test", 0)
	set.Bool("bypass-version-check", false, "")
	c := cli.NewContext(nil, set, nil)
	err := client.RemoteLogin(c)
	assert.Error(t, err)
	assert.EqualError(t, err, expErr)

	// Defaults to false
	set = flag.NewFlagSet("test", 0)
	c = cli.NewContext(nil, set, nil)
	err = client.RemoteLogin(c)
	assert.Error(t, err)
	assert.EqualError(t, err, expErr)
}

func TestClient_CheckRemoteBuildCompatibility(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	tests := []struct {
		name                         string
		remoteVersion, remoteSha     string
		cliVersion, cliSha           string
		bypassVersionFlag, wantError bool
	}{
		{"success match", "1.1.1", "53120d5", "1.1.1", "53120d5", false, false},
		{"cli unset fails", "1.1.1", "53120d5", "unset", "unset", false, true},
		{"remote unset fails", "unset", "unset", "1.1.1", "53120d5", false, true},
		{"mismatch fail", "1.1.1", "53120d5", "1.6.9", "13230sas", false, true},
		{"mismatch but using bypass_version_flag", "1.1.1", "53120d5", "1.6.9", "13230sas", true, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			enteredStrings := []string{cltest.APIEmail, cltest.Password}
			prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}
			client := app.NewAuthenticatingClient(prompter)

			client.HTTP = &mockHTTPClient{client.HTTP, test.remoteVersion, test.remoteSha}

			err := client.CheckRemoteBuildCompatibility(logger.TestLogger(t), test.bypassVersionFlag, test.cliVersion, test.cliSha)
			if test.wantError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, cmd.ErrIncompatible{
					RemoteVersion: test.remoteVersion,
					RemoteSha:     test.remoteSha,
					CLIVersion:    test.cliVersion,
					CLISha:        test.cliSha,
				})
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type mockHTTPClient struct {
	HTTP        cmd.HTTPClient
	mockVersion string
	mockSha     string
}

func (h *mockHTTPClient) Get(path string, headers ...map[string]string) (*http.Response, error) {
	if path == "/v2/build_info" {
		// Return mocked response here
		json := fmt.Sprintf(`{"version":"%s","commitSHA":"%s"}`, h.mockVersion, h.mockSha)
		r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	return h.HTTP.Get(path, headers...)
}

func (h *mockHTTPClient) Post(path string, body io.Reader) (*http.Response, error) {
	return h.HTTP.Post(path, body)
}

func (h *mockHTTPClient) Put(path string, body io.Reader) (*http.Response, error) {
	return h.HTTP.Put(path, body)
}

func (h *mockHTTPClient) Patch(path string, body io.Reader, headers ...map[string]string) (*http.Response, error) {
	return h.HTTP.Patch(path, body, headers...)
}

func (h *mockHTTPClient) Delete(path string) (*http.Response, error) {
	return h.HTTP.Delete(path)
}

func TestClient_ChangePassword(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)

	enteredStrings := []string{cltest.APIEmail, cltest.Password}
	prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}

	client := app.NewAuthenticatingClient(prompter)
	otherClient := app.NewAuthenticatingClient(prompter)

	set := flag.NewFlagSet("test", 0)
	set.String("file", "../internal/fixtures/apicredentials", "")
	set.Bool("bypass-version-check", true, "")
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

func TestClient_Profile_InvalidSecondsParam(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	enteredStrings := []string{cltest.APIEmail, cltest.Password}
	prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}

	client := app.NewAuthenticatingClient(prompter)

	set := flag.NewFlagSet("test", 0)
	set.String("file", "../internal/fixtures/apicredentials", "")
	set.Bool("bypass-version-check", true, "")
	c := cli.NewContext(nil, set, nil)
	err := client.RemoteLogin(c)
	require.NoError(t, err)

	set.Uint("seconds", 10, "")

	err = client.Profile(cli.NewContext(nil, set, nil))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "profile duration should be less than server write timeout.")

}

func TestClient_Profile(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	enteredStrings := []string{cltest.APIEmail, cltest.Password}
	prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}

	client := app.NewAuthenticatingClient(prompter)

	set := flag.NewFlagSet("test", 0)
	set.String("file", "../internal/fixtures/apicredentials", "")
	set.Bool("bypass-version-check", true, "")
	c := cli.NewContext(nil, set, nil)
	err := client.RemoteLogin(c)
	require.NoError(t, err)

	set.Uint("seconds", 8, "")
	set.String("output_dir", t.TempDir(), "")

	err = client.Profile(cli.NewContext(nil, set, nil))
	require.NoError(t, err)
}
func TestClient_SetDefaultGasPrice(t *testing.T) {
	t.Parallel()

	ethMock := newEthMock(t)
	app := startNewApplication(t,
		withKey(),
		withMocks(ethMock),
		withConfigSet(func(c *configtest.TestGeneralConfig) {
			c.Overrides.EVMEnabled = null.BoolFrom(true)
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
		ch, err := app.GetChains().EVM.Default()
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

		ch, err := app.GetChains().EVM.Default()
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
		ch, err := app.GetChains().EVM.Default()
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

	cp := *r.Renders[0].(*config.ConfigPrinter)
	assert.Equal(t, cp.EnvPrinter.BridgeResponseURL, cfg.BridgeResponseURL().String())
	assert.Equal(t, cp.EnvPrinter.DefaultChainID, cfg.DefaultChainID().String())
	assert.Equal(t, cp.EnvPrinter.Dev, cfg.Dev())
	assert.Equal(t, cp.EnvPrinter.LogLevel, cfg.LogLevel())
	assert.Equal(t, cp.EnvPrinter.LogSQL, cfg.LogSQL())
	assert.Equal(t, cp.EnvPrinter.RootDir, cfg.RootDir())
	assert.Equal(t, cp.EnvPrinter.SessionTimeout, cfg.SessionTimeout())
}

func TestClient_RunOCRJob_HappyPath(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t, withConfigSet(func(c *configtest.TestGeneralConfig) {
		c.Overrides.EVMEnabled = null.BoolFrom(true)
		c.Overrides.FeatureOffchainReporting = null.BoolFrom(true)
		c.Overrides.GlobalGasEstimatorMode = null.StringFrom("FixedPrice")
	}))
	client, _ := app.NewClientAndRenderer()

	app.KeyStore.OCR().Add(cltest.DefaultOCRKey)
	app.KeyStore.P2P().Add(cltest.DefaultP2PKey)

	_, bridge := cltest.MustCreateBridge(t, app.GetSqlxDB(), cltest.BridgeOpts{}, app.GetConfig())
	_, bridge2 := cltest.MustCreateBridge(t, app.GetSqlxDB(), cltest.BridgeOpts{}, app.GetConfig())

	var jb job.Job
	ocrspec := testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{DS1BridgeName: bridge.Name.String(), DS2BridgeName: bridge2.Name.String()})
	err := toml.Unmarshal([]byte(ocrspec.Toml()), &jb)
	require.NoError(t, err)
	var ocrSpec job.OCROracleSpec
	err = toml.Unmarshal([]byte(ocrspec.Toml()), &ocrspec)
	require.NoError(t, err)
	jb.OCROracleSpec = &ocrSpec
	key, _ := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())
	jb.OCROracleSpec.TransmitterAddress = &key.Address

	err = app.AddJobV2(context.Background(), &jb)
	require.NoError(t, err)

	set := flag.NewFlagSet("test", 0)
	set.Bool("bypass-version-check", true, "")
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
	set.Bool("bypass-version-check", true, "")
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
	set.Bool("bypass-version-check", true, "")
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.RemoteLogin(c))
	err := client.TriggerPipelineRun(c)
	assert.Contains(t, err.Error(), "parseResponse error: Error; job ID 1")
}

func TestClient_AutoLogin(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)

	user := cltest.MustRandomUser(t)
	require.NoError(t, app.SessionORM().CreateUser(&user))

	sr := sessions.SessionRequest{
		Email:    user.Email,
		Password: cltest.Password,
	}
	client, _ := app.NewClientAndRenderer()
	client.CookieAuthenticator = cmd.NewSessionCookieAuthenticator(app.NewClientOpts(), &cmd.MemoryCookieStore{}, logger.TestLogger(t))
	client.HTTP = cmd.NewAuthenticatedHTTPClient(app.Logger, app.NewClientOpts(), client.CookieAuthenticator, sr)

	fs := flag.NewFlagSet("", flag.ExitOnError)
	err := client.ListJobs(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)

	// Expire the session and then try again
	pgtest.MustExec(t, app.GetSqlxDB(), "TRUNCATE sessions")
	err = client.ListJobs(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)
}

func TestClient_AutoLogin_AuthFails(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)

	user := cltest.MustRandomUser(t)
	require.NoError(t, app.SessionORM().CreateUser(&user))

	sr := sessions.SessionRequest{
		Email:    user.Email,
		Password: cltest.Password,
	}
	client, _ := app.NewClientAndRenderer()
	client.CookieAuthenticator = FailingAuthenticator{}
	client.HTTP = cmd.NewAuthenticatedHTTPClient(app.Logger, app.NewClientOpts(), client.CookieAuthenticator, sr)

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

// Remove a session ID from disk
func (FailingAuthenticator) Logout() error {
	return errors.New("no luck")
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
	assert.Equal(t, sqlEnabled, app.Config.LogSQL())

	sqlEnabled = false
	set = flag.NewFlagSet("logsql", 0)
	set.Bool("disable", true, "")
	c = cli.NewContext(nil, set, nil)

	err = client.SetLogSQL(c)
	assert.NoError(t, err)
	assert.Equal(t, sqlEnabled, app.Config.LogSQL())
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

	level, ok := logger.NewORM(app.GetSqlxDB(), logger.TestLogger(t)).GetServiceLogLevel(logPkg)
	require.True(t, ok)
	assert.Equal(t, logLevel, level)
}
