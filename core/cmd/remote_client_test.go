package cmd_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/kylelemons/godebug/diff"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	configtest2 "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/web"
)

var (
	nilContext = cli.NewContext(nil, nil, nil)
)

type startOptions struct {
	// Use to set up mocks on the app
	FlagsAndDeps []interface{}
	// Add a key on start up
	WithKey bool
}

func startNewApplicationV2(t *testing.T, overrideFn func(c *chainlink.Config, s *chainlink.Secrets), setup ...func(opts *startOptions)) *cltest.TestApplication {
	t.Helper()

	sopts := &startOptions{
		FlagsAndDeps: []interface{}{},
	}
	for _, fn := range setup {
		fn(sopts)
	}

	config := configtest2.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.JobPipeline.HTTPRequest.DefaultTimeout = models.MustNewDuration(30 * time.Millisecond)
		f := false
		c.EVM[0].Enabled = &f
		c.P2P.V1.Enabled = &f
		c.P2P.V2.Enabled = &f

		if overrideFn != nil {
			overrideFn(c, s)
		}
	})

	app := cltest.NewApplicationWithConfigAndKey(t, config, sopts.FlagsAndDeps...)
	require.NoError(t, app.Start(testutils.Context(t)))

	return app
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

func newEthMock(t *testing.T) *evmclimocks.Client {
	t.Helper()
	return cltest.NewEthMocksWithStartupAssertions(t)
}

func newEthMockWithTransactionsOnBlocksAssertions(t *testing.T) *evmclimocks.Client {
	t.Helper()

	return cltest.NewEthMocksWithTransactionsOnBlocksAssertions(t)
}

func keyNameForTest(t *testing.T) string {
	return fmt.Sprintf("%s/%s_test_key.json", t.TempDir(), t.Name())
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
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
	})
	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("flagset", 0)
	cltest.FlagSetApplyFromAction(client.ReplayFromBlock, set, "")

	require.NoError(t, set.Set("block-number", "42"))

	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.ReplayFromBlock(c))

	require.NoError(t, set.Set("evm-chain-id", "12345678"))
	c = cli.NewContext(nil, set, nil)
	assert.ErrorContains(t, client.ReplayFromBlock(c), "evmChainID does not match any local chains")

	require.NoError(t, set.Set("evm-chain-id", "0"))
	c = cli.NewContext(nil, set, nil)
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
			app := startNewApplicationV2(t, nil)
			client, _ := app.NewClientAndRenderer()

			set := flag.NewFlagSet("create", 0)
			cltest.FlagSetApplyFromAction(client.CreateExternalInitiator, set, "")
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
			app := startNewApplicationV2(t, nil)
			client, _ := app.NewClientAndRenderer()

			initialExis := len(cltest.AllExternalInitiators(t, app.GetSqlxDB()))

			set := flag.NewFlagSet("create", 0)
			cltest.FlagSetApplyFromAction(client.CreateExternalInitiator, set, "")

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

	app := startNewApplicationV2(t, nil)
	client, r := app.NewClientAndRenderer()

	token := auth.NewToken()
	exi, err := bridges.NewExternalInitiator(token,
		&bridges.ExternalInitiatorRequest{Name: "name"},
	)
	require.NoError(t, err)
	err = app.BridgeORM().CreateExternalInitiator(exi)
	require.NoError(t, err)

	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.DeleteExternalInitiator, set, "")

	require.NoError(t, set.Parse([]string{exi.Name}))

	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.DeleteExternalInitiator(c))
	assert.Empty(t, r.Renders)
}

func TestClient_DestroyExternalInitiator_NotFound(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, r := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.DeleteExternalInitiator, set, "")

	require.NoError(t, set.Parse([]string{"bogus-ID"}))

	c := cli.NewContext(nil, set, nil)
	assert.Error(t, client.DeleteExternalInitiator(c))
	assert.Empty(t, r.Renders)
}

func TestClient_RemoteLogin(t *testing.T) {

	app := startNewApplicationV2(t, nil)

	tests := []struct {
		name, file string
		email, pwd string
		wantError  bool
	}{
		{"success prompt", "", cltest.APIEmailAdmin, cltest.Password, false},
		{"success file", "../internal/fixtures/apicredentials", "", "", false},
		{"failure prompt", "", "wrong@email.com", "wrongpwd", true},
		{"failure file", "/tmp/doesntexist", "", "", true},
		{"failure file w correct prompt", "/tmp/doesntexist", cltest.APIEmailAdmin, cltest.Password, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			enteredStrings := []string{test.email, test.pwd}
			prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}
			client := app.NewAuthenticatingClient(prompter)

			set := flag.NewFlagSet("test", 0)
			cltest.FlagSetApplyFromAction(client.RemoteLogin, set, "")

			require.NoError(t, set.Set("file", test.file))
			require.NoError(t, set.Set("bypass-version-check", "true"))

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

	app := startNewApplicationV2(t, nil)
	enteredStrings := []string{cltest.APIEmailAdmin, cltest.Password}
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
	cltest.FlagSetApplyFromAction(client.RemoteLogin, set, "")

	require.NoError(t, set.Set("bypass-version-check", "false"))

	c := cli.NewContext(nil, set, nil)
	err := client.RemoteLogin(c)
	assert.Error(t, err)
	assert.EqualError(t, err, expErr)

	// Defaults to false
	set = flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.RemoteLogin, set, "")
	c = cli.NewContext(nil, set, nil)
	err = client.RemoteLogin(c)
	assert.Error(t, err)
	assert.EqualError(t, err, expErr)
}

func TestClient_CheckRemoteBuildCompatibility(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
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
			enteredStrings := []string{cltest.APIEmailAdmin, cltest.Password}
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
		r := io.NopCloser(bytes.NewReader([]byte(json)))
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

	app := startNewApplicationV2(t, nil)

	enteredStrings := []string{cltest.APIEmailAdmin, cltest.Password}
	prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}

	client := app.NewAuthenticatingClient(prompter)
	otherClient := app.NewAuthenticatingClient(prompter)

	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.RemoteLogin, set, "")

	require.NoError(t, set.Set("file", "../internal/fixtures/apicredentials"))
	require.NoError(t, set.Set("bypass-version-check", "true"))

	c := cli.NewContext(nil, set, nil)
	err := client.RemoteLogin(c)
	require.NoError(t, err)

	err = otherClient.RemoteLogin(c)
	require.NoError(t, err)

	client.ChangePasswordPrompter = cltest.MockChangePasswordPrompter{
		UpdatePasswordRequest: web.UpdatePasswordRequest{
			OldPassword: testutils.Password,
			NewPassword: "12345",
		},
	}
	err = client.ChangePassword(cli.NewContext(nil, nil, nil))
	require.Error(t, err)
	assert.ErrorContains(t, err, "Expected password complexity")

	client.ChangePasswordPrompter = cltest.MockChangePasswordPrompter{
		UpdatePasswordRequest: web.UpdatePasswordRequest{
			OldPassword: testutils.Password,
			NewPassword: testutils.Password + "foo",
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

	app := startNewApplicationV2(t, nil)
	enteredStrings := []string{cltest.APIEmailAdmin, cltest.Password}
	prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}

	client := app.NewAuthenticatingClient(prompter)

	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.RemoteLogin, set, "")

	require.NoError(t, set.Set("file", "../internal/fixtures/apicredentials"))
	require.NoError(t, set.Set("bypass-version-check", "true"))

	c := cli.NewContext(nil, set, nil)
	err := client.RemoteLogin(c)
	require.NoError(t, err)

	// pick a value larger than the default http service write timeout
	d := app.Config.HTTPServerWriteTimeout() + 2*time.Second
	set.Uint("seconds", uint(d.Seconds()), "")
	tDir := t.TempDir()
	set.String("output_dir", tDir, "")
	err = client.Profile(cli.NewContext(nil, set, nil))
	wantErr := cmd.ErrProfileTooLong
	require.ErrorAs(t, err, &wantErr)

}

func TestClient_Profile(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	enteredStrings := []string{cltest.APIEmailAdmin, cltest.Password}
	prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}

	client := app.NewAuthenticatingClient(prompter)

	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.RemoteLogin, set, "")

	require.NoError(t, set.Set("file", "../internal/fixtures/apicredentials"))
	require.NoError(t, set.Set("bypass-version-check", "true"))

	c := cli.NewContext(nil, set, nil)
	err := client.RemoteLogin(c)
	require.NoError(t, err)

	set.Uint("seconds", 1, "")
	tDir := t.TempDir()
	set.String("output_dir", tDir, "")

	// we don't care about the cli behavior, i.e. the before func,
	// so call the client func directly
	err = client.Profile(cli.NewContext(nil, set, nil))
	require.NoError(t, err)

	ents, err := os.ReadDir(tDir)
	require.NoError(t, err)
	require.Greater(t, len(ents), 0, "ents %+v", ents)
}

func TestClient_Profile_Unauthenticated(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)

	client := app.NewAuthenticatingClient(&cltest.MockCountingPrompter{T: t, EnteredStrings: []string{}})

	set := flag.NewFlagSet("test", 0)
	set.Uint("seconds", 1, "")
	set.String("output_dir", t.TempDir(), "")

	err := client.Profile(cli.NewContext(nil, set, nil))
	require.ErrorContains(t, err, "profile collection failed:")
	require.ErrorContains(t, err, "Unauthorized")
}

func TestClient_ConfigV2(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, _ := app.NewClientAndRenderer()
	user, effective := app.Config.ConfigTOML()

	t.Run("user", func(t *testing.T) {
		got, err := client.ConfigV2Str(true)
		require.NoError(t, err)
		assert.Equal(t, user, got, diff.Diff(user, got))
	})
	t.Run("effective", func(t *testing.T) {
		got, err := client.ConfigV2Str(false)
		require.NoError(t, err)
		assert.Equal(t, effective, got, diff.Diff(effective, got))
	})
}

func TestClient_RunOCRJob_HappyPath(t *testing.T) {
	t.Parallel()
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.OCR.Enabled = ptr(true)
		c.P2P.V1.Enabled = ptr(true)
		c.P2P.PeerID = &cltest.DefaultP2PPeerID
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
	}, func(opts *startOptions) {
		opts.FlagsAndDeps = append(opts.FlagsAndDeps, cltest.DefaultP2PKey)
	})
	client, _ := app.NewClientAndRenderer()

	require.NoError(t, app.KeyStore.OCR().Add(cltest.DefaultOCRKey))

	_, bridge := cltest.MustCreateBridge(t, app.GetSqlxDB(), cltest.BridgeOpts{}, app.GetConfig())
	_, bridge2 := cltest.MustCreateBridge(t, app.GetSqlxDB(), cltest.BridgeOpts{}, app.GetConfig())

	var jb job.Job
	ocrspec := testspecs.GenerateOCRSpec(testspecs.OCRSpecParams{DS1BridgeName: bridge.Name.String(), DS2BridgeName: bridge2.Name.String()})
	err := toml.Unmarshal([]byte(ocrspec.Toml()), &jb)
	require.NoError(t, err)
	var ocrSpec job.OCROracleSpec
	err = toml.Unmarshal([]byte(ocrspec.Toml()), &ocrSpec)
	require.NoError(t, err)
	jb.OCROracleSpec = &ocrSpec
	key, _ := cltest.MustInsertRandomKey(t, app.KeyStore.Eth())
	jb.OCROracleSpec.TransmitterAddress = &key.EIP55Address

	err = app.AddJobV2(testutils.Context(t), &jb)
	require.NoError(t, err)

	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.RemoteLogin, set, "")

	require.NoError(t, set.Set("bypass-version-check", "true"))
	require.NoError(t, set.Parse([]string{strconv.FormatInt(int64(jb.ID), 10)}))

	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.RemoteLogin(c))
	require.NoError(t, client.TriggerPipelineRun(c))
}

func TestClient_RunOCRJob_MissingJobID(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.RemoteLogin, set, "")

	require.NoError(t, set.Set("bypass-version-check", "true"))

	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.RemoteLogin(c))
	assert.EqualError(t, client.TriggerPipelineRun(c), "Must pass the job id to trigger a run")
}

func TestClient_RunOCRJob_JobNotFound(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	cltest.FlagSetApplyFromAction(client.RemoteLogin, set, "")

	require.NoError(t, set.Parse([]string{"1"}))
	require.NoError(t, set.Set("bypass-version-check", "true"))

	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.RemoteLogin(c))
	err := client.TriggerPipelineRun(c)
	assert.Contains(t, err.Error(), "findJob failed: failed to load job")
}

func TestClient_AutoLogin(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)

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
	cltest.FlagSetApplyFromAction(client.ListJobs, fs, "")

	err := client.ListJobs(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)

	// Expire the session and then try again
	pgtest.MustExec(t, app.GetSqlxDB(), "TRUNCATE sessions")
	err = client.ListJobs(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)
}

func TestClient_AutoLogin_AuthFails(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)

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
	cltest.FlagSetApplyFromAction(client.ListJobs, fs, "")
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

	app := startNewApplicationV2(t, nil)
	client, _ := app.NewClientAndRenderer()

	logLevel := "warn"
	set := flag.NewFlagSet("loglevel", 0)
	cltest.FlagSetApplyFromAction(client.SetLogLevel, set, "")

	require.NoError(t, set.Set("level", logLevel))

	c := cli.NewContext(nil, set, nil)

	err := client.SetLogLevel(c)
	require.NoError(t, err)
	assert.Equal(t, logLevel, app.Config.LogLevel().String())

	sqlEnabled := true
	set = flag.NewFlagSet("logsql", 0)
	cltest.FlagSetApplyFromAction(client.SetLogSQL, set, "")

	require.NoError(t, set.Set("enable", strconv.FormatBool(sqlEnabled)))
	c = cli.NewContext(nil, set, nil)

	err = client.SetLogSQL(c)
	assert.NoError(t, err)
	assert.Equal(t, sqlEnabled, app.Config.LogSQL())

	sqlEnabled = false
	set = flag.NewFlagSet("logsql", 0)
	cltest.FlagSetApplyFromAction(client.SetLogSQL, set, "")

	require.NoError(t, set.Set("disable", "true"))
	c = cli.NewContext(nil, set, nil)

	err = client.SetLogSQL(c)
	assert.NoError(t, err)
	assert.Equal(t, sqlEnabled, app.Config.LogSQL())
}
