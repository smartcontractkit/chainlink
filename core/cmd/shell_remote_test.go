package cmd_test

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/kylelemons/godebug/diff"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/static"
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

	config := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.JobPipeline.HTTPRequest.DefaultTimeout = commonconfig.MustNewDuration(30 * time.Millisecond)
		f := false
		c.EVM[0].Enabled = &f
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

func TestShell_ReplayBlocks(t *testing.T) {
	t.Parallel()
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.EVM[0].NonceAutoSync = ptr(false)
		c.EVM[0].BalanceMonitor.Enabled = ptr(false)
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
	})
	client, _ := app.NewShellAndRenderer()

	set := flag.NewFlagSet("flagset", 0)
	flagSetApplyFromAction(client.ReplayFromBlock, set, "")

	require.NoError(t, set.Set("block-number", "42"))
	require.NoError(t, set.Set("evm-chain-id", "12345678"))
	c := cli.NewContext(nil, set, nil)
	assert.ErrorContains(t, client.ReplayFromBlock(c), "chain id does not match any local chains")

	require.NoError(t, set.Set("evm-chain-id", "0"))
	c = cli.NewContext(nil, set, nil)
	assert.NoError(t, client.ReplayFromBlock(c))
}

func TestShell_CreateExternalInitiator(t *testing.T) {
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
			ctx := testutils.Context(t)
			app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.JobPipeline.ExternalInitiatorsEnabled = ptr(true)
			})
			client, _ := app.NewShellAndRenderer()

			set := flag.NewFlagSet("create", 0)
			flagSetApplyFromAction(client.CreateExternalInitiator, set, "")
			assert.NoError(t, set.Parse(test.args))
			c := cli.NewContext(nil, set, nil)

			err := client.CreateExternalInitiator(c)
			require.NoError(t, err)

			var exi bridges.ExternalInitiator
			err = app.GetDB().GetContext(ctx, &exi, `SELECT * FROM external_initiators WHERE name = $1`, test.args[0])
			require.NoError(t, err)

			if len(test.args) > 1 {
				assert.Equal(t, test.args[1], exi.URL.String())
			}
		})
	}
}

func TestShell_CreateExternalInitiator_Errors(t *testing.T) {
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
			app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.JobPipeline.ExternalInitiatorsEnabled = ptr(true)
			})
			client, _ := app.NewShellAndRenderer()

			initialExis := len(cltest.AllExternalInitiators(t, app.GetDB()))

			set := flag.NewFlagSet("create", 0)
			flagSetApplyFromAction(client.CreateExternalInitiator, set, "")

			assert.NoError(t, set.Parse(test.args))
			c := cli.NewContext(nil, set, nil)

			err := client.CreateExternalInitiator(c)
			assert.Error(t, err)

			exis := cltest.AllExternalInitiators(t, app.GetDB())
			assert.Len(t, exis, initialExis)
		})
	}
}

func TestShell_DestroyExternalInitiator(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.JobPipeline.ExternalInitiatorsEnabled = ptr(true)
	})
	client, r := app.NewShellAndRenderer()

	token := auth.NewToken()
	exi, err := bridges.NewExternalInitiator(token,
		&bridges.ExternalInitiatorRequest{Name: uuid.New().String()},
	)
	require.NoError(t, err)
	err = app.BridgeORM().CreateExternalInitiator(testutils.Context(t), exi)
	require.NoError(t, err)

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.DeleteExternalInitiator, set, "")

	require.NoError(t, set.Parse([]string{exi.Name}))

	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.DeleteExternalInitiator(c))
	assert.Empty(t, r.Renders)
}

func TestShell_DestroyExternalInitiator_NotFound(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.JobPipeline.ExternalInitiatorsEnabled = ptr(true)
	})
	client, r := app.NewShellAndRenderer()

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.DeleteExternalInitiator, set, "")

	require.NoError(t, set.Parse([]string{"bogus-ID"}))

	c := cli.NewContext(nil, set, nil)
	assert.Error(t, client.DeleteExternalInitiator(c))
	assert.Empty(t, r.Renders)
}

func TestShell_RemoteLogin(t *testing.T) {
	app := startNewApplicationV2(t, nil)
	orm := app.AuthenticationProvider()

	u := cltest.NewUserWithSession(t, orm)

	tests := []struct {
		name, file string
		email, pwd string
		wantError  bool
	}{
		{"success prompt", "", u.Email, cltest.Password, false},
		{"success file", "../internal/fixtures/apicredentials", "", "", false},
		{"failure prompt", "", "wrong@email.com", "wrongpwd", true},
		{"failure file", "/tmp/doesntexist", "", "", true},
		{"failure file w correct prompt", "/tmp/doesntexist", u.Email, cltest.Password, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			enteredStrings := []string{test.email, test.pwd}
			prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}
			client := app.NewAuthenticatingShell(prompter)

			set := flag.NewFlagSet("test", 0)
			flagSetApplyFromAction(client.RemoteLogin, set, "")

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

func TestShell_RemoteBuildCompatibility(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	u := cltest.NewUserWithSession(t, app.AuthenticationProvider())
	enteredStrings := []string{u.Email, cltest.Password}
	prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: append(enteredStrings, enteredStrings...)}
	client := app.NewAuthenticatingShell(prompter)

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
	flagSetApplyFromAction(client.RemoteLogin, set, "")

	require.NoError(t, set.Set("bypass-version-check", "false"))

	c := cli.NewContext(nil, set, nil)
	err := client.RemoteLogin(c)
	assert.Error(t, err)
	assert.EqualError(t, err, expErr)

	// Defaults to false
	set = flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.RemoteLogin, set, "")
	c = cli.NewContext(nil, set, nil)
	err = client.RemoteLogin(c)
	assert.Error(t, err)
	assert.EqualError(t, err, expErr)
}

func TestShell_CheckRemoteBuildCompatibility(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	u := cltest.NewUserWithSession(t, app.AuthenticationProvider())
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
			enteredStrings := []string{u.Email, cltest.Password}
			prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}
			client := app.NewAuthenticatingShell(prompter)

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

func (h *mockHTTPClient) Get(ctx context.Context, path string, headers ...map[string]string) (*http.Response, error) {
	if path == "/v2/build_info" {
		// Return mocked response here
		json := fmt.Sprintf(`{"version":"%s","commitSHA":"%s"}`, h.mockVersion, h.mockSha)
		r := io.NopCloser(bytes.NewReader([]byte(json)))
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	return h.HTTP.Get(ctx, path, headers...)
}

func (h *mockHTTPClient) Post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	return h.HTTP.Post(ctx, path, body)
}

func (h *mockHTTPClient) Put(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	return h.HTTP.Put(ctx, path, body)
}

func (h *mockHTTPClient) Patch(ctx context.Context, path string, body io.Reader, headers ...map[string]string) (*http.Response, error) {
	return h.HTTP.Patch(ctx, path, body, headers...)
}

func (h *mockHTTPClient) Delete(ctx context.Context, path string) (*http.Response, error) {
	return h.HTTP.Delete(ctx, path)
}

func TestShell_ChangePassword(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	u := cltest.NewUserWithSession(t, app.AuthenticationProvider())

	enteredStrings := []string{u.Email, cltest.Password}
	prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}

	client := app.NewAuthenticatingShell(prompter)
	otherClient := app.NewAuthenticatingShell(prompter)

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.RemoteLogin, set, "")

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

func TestShell_Profile_InvalidSecondsParam(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	u := cltest.NewUserWithSession(t, app.AuthenticationProvider())
	enteredStrings := []string{u.Email, cltest.Password}
	prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}

	client := app.NewAuthenticatingShell(prompter)

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.RemoteLogin, set, "")

	require.NoError(t, set.Set("file", "../internal/fixtures/apicredentials"))
	require.NoError(t, set.Set("bypass-version-check", "true"))

	c := cli.NewContext(nil, set, nil)
	err := client.RemoteLogin(c)
	require.NoError(t, err)

	// pick a value larger than the default http service write timeout
	d := app.Config.WebServer().HTTPWriteTimeout() + 2*time.Second
	set.Uint("seconds", uint(d.Seconds()), "")
	tDir := t.TempDir()
	set.String("output_dir", tDir, "")
	err = client.Profile(cli.NewContext(nil, set, nil))
	wantErr := cmd.ErrProfileTooLong
	require.ErrorAs(t, err, &wantErr)
}

func TestShell_Profile(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	u := cltest.NewUserWithSession(t, app.AuthenticationProvider())
	enteredStrings := []string{u.Email, cltest.Password}
	prompter := &cltest.MockCountingPrompter{T: t, EnteredStrings: enteredStrings}

	client := app.NewAuthenticatingShell(prompter)

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.RemoteLogin, set, "")

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

func TestShell_Profile_Unauthenticated(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)

	client := app.NewAuthenticatingShell(&cltest.MockCountingPrompter{T: t, EnteredStrings: []string{}})

	set := flag.NewFlagSet("test", 0)
	set.Uint("seconds", 1, "")
	set.String("output_dir", t.TempDir(), "")

	err := client.Profile(cli.NewContext(nil, set, nil))
	require.ErrorContains(t, err, "profile collection failed:")
	require.ErrorContains(t, err, "Unauthorized")
}

func TestShell_ConfigV2(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, _ := app.NewShellAndRenderer()
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

func TestShell_RunOCRJob_HappyPath(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	app := startNewApplicationV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].Enabled = ptr(true)
		c.OCR.Enabled = ptr(true)
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", freeport.GetOne(t))}
		c.P2P.PeerID = &cltest.DefaultP2PPeerID
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
	}, func(opts *startOptions) {
		opts.FlagsAndDeps = append(opts.FlagsAndDeps, cltest.DefaultP2PKey)
	})
	client, _ := app.NewShellAndRenderer()

	require.NoError(t, app.KeyStore.OCR().Add(ctx, cltest.DefaultOCRKey))

	_, bridge := cltest.MustCreateBridge(t, app.GetDB(), cltest.BridgeOpts{})
	_, bridge2 := cltest.MustCreateBridge(t, app.GetDB(), cltest.BridgeOpts{})

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
	flagSetApplyFromAction(client.RemoteLogin, set, "")

	require.NoError(t, set.Set("bypass-version-check", "true"))
	require.NoError(t, set.Parse([]string{strconv.FormatInt(int64(jb.ID), 10)}))

	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.RemoteLogin(c))
	require.NoError(t, client.TriggerPipelineRun(c))
}

func TestShell_RunOCRJob_MissingJobID(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, _ := app.NewShellAndRenderer()

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.RemoteLogin, set, "")

	require.NoError(t, set.Set("bypass-version-check", "true"))

	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.RemoteLogin(c))
	assert.EqualError(t, client.TriggerPipelineRun(c), "Must pass the job id to trigger a run")
}

func TestShell_RunOCRJob_JobNotFound(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, _ := app.NewShellAndRenderer()

	set := flag.NewFlagSet("test", 0)
	flagSetApplyFromAction(client.RemoteLogin, set, "")

	require.NoError(t, set.Parse([]string{"1"}))
	require.NoError(t, set.Set("bypass-version-check", "true"))

	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.RemoteLogin(c))
	err := client.TriggerPipelineRun(c)
	assert.Contains(t, err.Error(), "findJob failed: failed to load job")
}

func TestShell_AutoLogin(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	app := startNewApplicationV2(t, nil)

	user := cltest.MustRandomUser(t)
	require.NoError(t, app.BasicAdminUsersORM().CreateUser(ctx, &user))

	sr := sessions.SessionRequest{
		Email:    user.Email,
		Password: cltest.Password,
	}
	client, _ := app.NewShellAndRenderer()
	client.CookieAuthenticator = cmd.NewSessionCookieAuthenticator(app.NewClientOpts(), &cmd.MemoryCookieStore{}, logger.TestLogger(t))
	client.HTTP = cmd.NewAuthenticatedHTTPClient(app.Logger, app.NewClientOpts(), client.CookieAuthenticator, sr)

	fs := flag.NewFlagSet("", flag.ExitOnError)
	flagSetApplyFromAction(client.ListJobs, fs, "")

	err := client.ListJobs(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)

	// Expire the session and then try again
	pgtest.MustExec(t, app.GetDB(), "delete from sessions where email = $1", user.Email)
	err = client.ListJobs(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)
}

func TestShell_AutoLogin_AuthFails(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	app := startNewApplicationV2(t, nil)

	user := cltest.MustRandomUser(t)
	require.NoError(t, app.BasicAdminUsersORM().CreateUser(ctx, &user))

	sr := sessions.SessionRequest{
		Email:    user.Email,
		Password: cltest.Password,
	}
	client, _ := app.NewShellAndRenderer()
	client.CookieAuthenticator = FailingAuthenticator{}
	client.HTTP = cmd.NewAuthenticatedHTTPClient(app.Logger, app.NewClientOpts(), client.CookieAuthenticator, sr)

	fs := flag.NewFlagSet("", flag.ExitOnError)
	flagSetApplyFromAction(client.ListJobs, fs, "")
	err := client.ListJobs(cli.NewContext(nil, fs, nil))
	require.Error(t, err)
}

type FailingAuthenticator struct{}

func (FailingAuthenticator) Cookie() (*http.Cookie, error) {
	return &http.Cookie{}, nil
}

// Authenticate retrieves a session ID via a cookie and saves it to disk.
func (FailingAuthenticator) Authenticate(context.Context, sessions.SessionRequest) (*http.Cookie, error) {
	return nil, errors.New("no luck")
}

// Remove a session ID from disk
func (FailingAuthenticator) Logout() error {
	return errors.New("no luck")
}

func TestShell_SetLogConfig(t *testing.T) {
	t.Parallel()

	app := startNewApplicationV2(t, nil)
	client, _ := app.NewShellAndRenderer()

	logLevel := "warn"
	set := flag.NewFlagSet("loglevel", 0)
	flagSetApplyFromAction(client.SetLogLevel, set, "")

	require.NoError(t, set.Set("level", logLevel))

	c := cli.NewContext(nil, set, nil)

	err := client.SetLogLevel(c)
	require.NoError(t, err)
	assert.Equal(t, logLevel, app.Config.Log().Level().String())

	sqlEnabled := true
	set = flag.NewFlagSet("logsql", 0)
	flagSetApplyFromAction(client.SetLogSQL, set, "")

	require.NoError(t, set.Set("enable", strconv.FormatBool(sqlEnabled)))
	c = cli.NewContext(nil, set, nil)

	err = client.SetLogSQL(c)
	assert.NoError(t, err)
	assert.Equal(t, sqlEnabled, app.Config.Database().LogSQL())

	sqlEnabled = false
	set = flag.NewFlagSet("logsql", 0)
	flagSetApplyFromAction(client.SetLogSQL, set, "")

	require.NoError(t, set.Set("disable", "true"))
	c = cli.NewContext(nil, set, nil)

	err = client.SetLogSQL(c)
	assert.NoError(t, err)
	assert.Equal(t, sqlEnabled, app.Config.Database().LogSQL())
}
