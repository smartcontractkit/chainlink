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

	"github.com/smartcontractkit/chainlink/core/logger"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/cron"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	webhookmocks "github.com/smartcontractkit/chainlink/core/services/webhook/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

var (
	nilContext = cli.NewContext(nil, nil, nil)
)

type startOptions struct {
	// Set the config options
	Config map[string]interface{}
	// Use to set up mocks on the app
	FlagsAndDeps []interface{}
	// Add a key on start up
	WithKey bool
	// Use app.StartAndConnect instead of app.Start
	StartAndConnect bool
}

func startNewApplication(t *testing.T, setup ...func(opts *startOptions)) *cltest.TestApplication {
	t.Helper()

	sopts := &startOptions{
		Config:       map[string]interface{}{},
		FlagsAndDeps: []interface{}{},
	}
	for _, fn := range setup {
		fn(sopts)
	}

	// Setup config
	config, cfgCleanup := cltest.NewConfig(t)
	t.Cleanup(cfgCleanup)
	config.Set("DEFAULT_HTTP_TIMEOUT", "30ms")
	config.Set("MAX_HTTP_ATTEMPTS", "1")
	config.Set("ENABLE_LEGACY_JOB_PIPELINE", true)

	for k, v := range sopts.Config {
		config.Set(k, v)
	}

	var app *cltest.TestApplication
	var cleanup func()
	app, cleanup = cltest.NewApplicationWithConfigAndKey(t, config, sopts.FlagsAndDeps...)
	t.Cleanup(cleanup)
	app.Logger = app.Config.CreateProductionLogger()
	app.Logger.SetDB(app.GetStore().DB)

	if sopts.StartAndConnect {
		require.NoError(t, app.StartAndConnect())
	} else {
		require.NoError(t, app.Start())
	}

	return app
}

// withConfig is a function option which sets config on the app
func withConfig(cfgs map[string]interface{}) func(opts *startOptions) {
	return func(opts *startOptions) {
		for k, v := range cfgs {
			opts.Config[k] = v
		}
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

func startAndConnect() func(opts *startOptions) {
	return func(opts *startOptions) {
		opts.StartAndConnect = true
	}
}

func newEthMock(t *testing.T) *mocks.Client {
	t.Helper()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	t.Cleanup(assertMocksCalled)

	return ethClient
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

func TestClient_IndexJobSpecs(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	j1 := cltest.NewJob()
	app.Store.CreateJob(&j1)
	j2 := cltest.NewJob()
	app.Store.CreateJob(&j2)

	require.Nil(t, client.IndexJobSpecs(cltest.EmptyCLIContext()))
	jobs := *r.Renders[0].(*[]models.JobSpec)
	require.Equal(t, 2, len(jobs))
	assert.Equal(t, j1.ID, jobs[0].ID)
}

func TestClient_ShowJobRun_Exists(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	j := cltest.NewJobWithWebInitiator()
	assert.NoError(t, app.Store.CreateJob(&j))

	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"result":"100"}`)

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{jr.ID.String()})
	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.ShowJobRun(c))
	assert.Equal(t, 1, len(r.Renders))
	assert.Equal(t, jr.ID, r.Renders[0].(*presenters.JobRun).ID)
}

func TestClient_ShowJobRun_NotFound(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"bogus-ID"})
	c := cli.NewContext(nil, set, nil)
	assert.Error(t, client.ShowJobRun(c))
	assert.Empty(t, r.Renders)
}

func TestClient_ReplayBlocks(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("flagset", 0)
	set.Int64("block-number", 42, "")
	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.ReplayFromBlock(c))
}

func TestClient_IndexJobRuns(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	j := cltest.NewJobWithWebInitiator()
	assert.NoError(t, app.Store.CreateJob(&j))

	jr0 := cltest.NewJobRun(j)
	jr0.Result.Data = cltest.JSONFromString(t, `{"a":"b"}`)
	require.NoError(t, app.Store.CreateJobRun(&jr0))
	jr1 := cltest.NewJobRun(j)
	jr1.Result.Data = cltest.JSONFromString(t, `{"x":"y"}`)
	require.NoError(t, app.Store.CreateJobRun(&jr1))

	require.Nil(t, client.IndexJobRuns(cltest.EmptyCLIContext()))
	runs := *r.Renders[0].(*[]presenters.JobRun)
	require.Len(t, runs, 2)
	assert.Equal(t, jr0.ID, runs[0].ID)
	assert.JSONEq(t, `{"a":"b"}`, runs[0].Result.Data.String())
	assert.Equal(t, jr1.ID, runs[1].ID)
	assert.JSONEq(t, `{"x":"y"}`, runs[1].Result.Data.String())
}

func TestClient_ShowJobSpec_Exists(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	job := cltest.NewJob()
	app.Store.CreateJob(&job)

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{job.ID.String()})
	c := cli.NewContext(nil, set, nil)
	require.Nil(t, client.ShowJobSpec(c))
	require.Equal(t, 1, len(r.Renders))
	assert.Equal(t, job.ID, r.Renders[0].(*presenters.JobSpec).ID)
}

func TestClient_ShowJobSpec_NotFound(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"bogus-ID"})
	c := cli.NewContext(nil, set, nil)
	assert.Error(t, client.ShowJobSpec(c))
	assert.Empty(t, r.Renders)
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

			var exi models.ExternalInitiator
			err = app.Store.RawDBWithAdvisoryLock(func(db *gorm.DB) error {
				return db.Where("name = ?", test.args[0]).Find(&exi).Error
			})
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

			initialExis := len(cltest.AllExternalInitiators(t, app.Store))

			set := flag.NewFlagSet("create", 0)
			assert.NoError(t, set.Parse(test.args))
			c := cli.NewContext(nil, set, nil)

			err := client.CreateExternalInitiator(c)
			assert.Error(t, err)

			exis := cltest.AllExternalInitiators(t, app.Store)
			assert.Len(t, exis, initialExis)
		})
	}
}

func TestClient_DestroyExternalInitiator(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	token := auth.NewToken()
	exi, err := models.NewExternalInitiator(token,
		&models.ExternalInitiatorRequest{Name: "name"},
	)
	require.NoError(t, err)
	err = app.Store.CreateExternalInitiator(exi)
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

func TestClient_CreateJobSpec(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	tests := []struct {
		name, input string
		nJobs       int
		errored     bool
	}{
		{"bad json", "{bad son}", 0, true},
		{"bad filepath", "bad/filepath/", 0, true},
		{"web", `{"initiators":[{"type":"web"}],"tasks":[{"type":"NoOp"}]}`, 1, false},
		{"runAt", `{"initiators":[{"type":"runAt","params":{"time":"3000-01-08T18:12:01.103Z"}}],"tasks":[{"type":"NoOp"}]}`, 2, false},
		{"file", "../testdata/jsonspecs/end_at_job.json", 3, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			set := flag.NewFlagSet("create", 0)
			set.Parse([]string{test.input})
			c := cli.NewContext(nil, set, nil)

			err := client.CreateJobSpec(c)
			cltest.AssertError(t, test.errored, err)

			numberOfJobs := cltest.AllJobs(t, app.Store)
			assert.Equal(t, test.nJobs, len(numberOfJobs))
		})
	}
}

func TestClient_ArchiveJobSpec(t *testing.T) {
	t.Parallel()

	eim := new(webhookmocks.ExternalInitiatorManager)
	app := startNewApplication(t, withMocks(eim))
	client, _ := app.NewClientAndRenderer()

	job := cltest.NewJob()
	require.NoError(t, app.Store.CreateJob(&job))

	set := flag.NewFlagSet("archive", 0)
	set.Parse([]string{job.ID.String()})
	c := cli.NewContext(nil, set, nil)

	eim.On("DeleteJob", mock.MatchedBy(func(id models.JobID) bool {
		return id.String() == job.ID.String()
	})).Once().Return(nil)

	require.NoError(t, client.ArchiveJobSpec(c))

	jobs := cltest.AllJobs(t, app.Store)
	require.Len(t, jobs, 0)
}

func TestClient_CreateJobSpec_JSONAPIErrors(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("create", 0)
	set.Parse([]string{`{"initiators":[{"type":"runAt"}],"tasks":[{"type":"NoOp"}]}`})
	c := cli.NewContext(nil, set, nil)

	err := client.CreateJobSpec(c)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must have a time")
}

func TestClient_CreateJobRun(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	tests := []struct {
		name    string
		json    string
		jobSpec models.JobSpec
		errored bool
	}{
		{"CreateSuccess", `{"result": 100}`, cltest.NewJobWithWebInitiator(), false},
		{"EmptyBody", ``, cltest.NewJobWithWebInitiator(), false},
		{"InvalidBody", `{`, cltest.NewJobWithWebInitiator(), true},
		{"WithoutWebInitiator", ``, cltest.NewJobWithLogInitiator(), true},
		{"NotFound", ``, cltest.NewJobWithWebInitiator(), true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			assert.Nil(t, app.Store.CreateJob(&test.jobSpec))

			args := make([]string, 1)
			args[0] = test.jobSpec.ID.String()
			if test.name == "NotFound" {
				args[0] = "badID"
			}

			if len(test.json) > 0 {
				args = append(args, test.json)
			}

			set := flag.NewFlagSet("run", 0)
			set.Parse(args)
			c := cli.NewContext(nil, set, nil)
			if test.errored {
				assert.Error(t, client.CreateJobRun(c))
			} else {
				assert.Nil(t, client.CreateJobRun(c))
			}
		})
	}
}

func TestClient_RemoteLogin(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t, withConfig(map[string]interface{}{
		"ADMIN_CREDENTIALS_FILE": "",
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

func TestClient_SetMinimumGasPrice(t *testing.T) {
	t.Parallel()

	// Setup Withdrawals application
	oca := common.HexToAddress("0xDEADB3333333F")
	app := startNewApplication(t,
		withKey(),
		withConfig(map[string]interface{}{
			"OPERATOR_CONTRACT_ADDRESS": &oca,
		}),
		withMocks(newEthMock(t)),
		startAndConnect(),
	)
	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("setgasprice", 0)
	set.Parse([]string{"8616460799"})

	c := cli.NewContext(nil, set, nil)

	assert.NoError(t, client.SetMinimumGasPrice(c))
	assert.Equal(t, big.NewInt(8616460799), app.Store.Config.EthGasPriceDefault())

	client, _ = app.NewClientAndRenderer()
	set = flag.NewFlagSet("setgasprice", 0)
	set.String("amount", "861.6460799", "")
	set.Bool("gwei", true, "")
	set.Parse([]string{"-gwei", "861.6460799"})

	c = cli.NewContext(nil, set, nil)
	assert.NoError(t, client.SetMinimumGasPrice(c))
	assert.Equal(t, big.NewInt(861646079900), app.Store.Config.EthGasPriceDefault())
}

func TestClient_GetConfiguration(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	assert.NoError(t, client.GetConfiguration(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))

	cp := *r.Renders[0].(*presenters.ConfigPrinter)
	assert.Equal(t, cp.EnvPrinter.BridgeResponseURL, app.Config.BridgeResponseURL().String())
	assert.Equal(t, cp.EnvPrinter.ChainID, app.Config.ChainID())
	assert.Equal(t, cp.EnvPrinter.Dev, app.Config.Dev())
	assert.Equal(t, cp.EnvPrinter.EthGasBumpThreshold, app.Config.EthGasBumpThreshold())
	assert.Equal(t, cp.EnvPrinter.LogLevel, app.Config.LogLevel())
	assert.Equal(t, cp.EnvPrinter.LogSQLStatements, app.Config.LogSQLStatements())
	assert.Equal(t, cp.EnvPrinter.MinIncomingConfirmations, app.Config.MinIncomingConfirmations())
	assert.Equal(t, cp.EnvPrinter.MinRequiredOutgoingConfirmations, app.Config.MinRequiredOutgoingConfirmations())
	assert.Equal(t, cp.EnvPrinter.MinimumContractPayment, app.Config.MinimumContractPayment())
	assert.Equal(t, cp.EnvPrinter.RootDir, app.Config.RootDir())
	assert.Equal(t, cp.EnvPrinter.SessionTimeout, app.Config.SessionTimeout())
}

func TestClient_CancelJobRun(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, app.Store.CreateJob(&job))
	run := cltest.NewJobRun(job)
	require.NoError(t, app.Store.CreateJobRun(&run))

	set := flag.NewFlagSet("cancel", 0)
	set.Parse([]string{run.ID.String()})
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.CancelJobRun(c))

	runs := cltest.MustAllJobsWithStatus(t, app.Store, models.RunStatusCancelled)
	require.Len(t, runs, 1)
	assert.Equal(t, models.RunStatusCancelled, runs[0].GetStatus())
	assert.NotNil(t, runs[0].FinishedAt)
}

func TestClient_RunOCRJob_HappyPath(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	_, bridge := cltest.NewBridgeType(t, "voter_turnout", "http://blah.com")
	require.NoError(t, app.Store.DB.Create(bridge).Error)
	_, bridge2 := cltest.NewBridgeType(t, "election_winner", "http://blah.com")
	require.NoError(t, app.Store.DB.Create(bridge2).Error)

	var ocrJobSpecFromFile job.Job
	tree, err := toml.LoadFile("../testdata/tomlspecs/oracle-spec.toml")
	require.NoError(t, err)
	err = tree.Unmarshal(&ocrJobSpecFromFile)
	require.NoError(t, err)
	var ocrSpec job.OffchainReportingOracleSpec
	err = tree.Unmarshal(&ocrSpec)
	require.NoError(t, err)
	ocrJobSpecFromFile.OffchainreportingOracleSpec = &ocrSpec
	key := cltest.MustInsertRandomKey(t, app.Store.DB)
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
	require.NoError(t, app.Store.SaveUser(&user))

	sr := models.SessionRequest{
		Email:    user.Email,
		Password: cltest.Password,
	}
	client, _ := app.NewClientAndRenderer()
	client.CookieAuthenticator = cmd.NewSessionCookieAuthenticator(app.Config.Config, &cmd.MemoryCookieStore{})
	client.HTTP = cmd.NewAuthenticatedHTTPClient(app.Config, client.CookieAuthenticator, sr)

	fs := flag.NewFlagSet("", flag.ExitOnError)
	err := client.ListJobsV2(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)

	// Expire the session and then try again
	require.NoError(t, app.GetStore().ORM.DB.Exec("delete from sessions;").Error)
	err = client.ListJobsV2(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)
}

func TestClient_AutoLogin_AuthFails(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)

	user := cltest.MustRandomUser()
	require.NoError(t, app.Store.SaveUser(&user))

	sr := models.SessionRequest{
		Email:    user.Email,
		Password: cltest.Password,
	}
	client, _ := app.NewClientAndRenderer()
	client.CookieAuthenticator = FailingAuthenticator{}
	client.HTTP = cmd.NewAuthenticatedHTTPClient(app.Config, client.CookieAuthenticator, sr)

	fs := flag.NewFlagSet("", flag.ExitOnError)
	err := client.ListJobsV2(cli.NewContext(nil, fs, nil))
	require.Error(t, err)
}

type FailingAuthenticator struct{}

func (FailingAuthenticator) Cookie() (*http.Cookie, error) {
	return &http.Cookie{}, nil
}

// Authenticate retrieves a session ID via a cookie and saves it to disk.
func (FailingAuthenticator) Authenticate(sessionRequest models.SessionRequest) (*http.Cookie, error) {
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

	level, err := app.Logger.ServiceLogLevel(logPkg)
	require.NoError(t, err)
	assert.Equal(t, logLevel, level)
}

func TestClient_MigrateCron(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("migrate", 0)
	set.Parse([]string{"../testdata/jsonspecs/example-cron2.json"})
	c := cli.NewContext(nil, set, nil)

	toml, _, err := client.MigrateJobSpecForResult(c)
	require.NoError(t, err)

	fmt.Println(toml)

	_, err = job.ValidateSpec(toml)
	require.NoError(t, err)

	var jb job.Job
	jb, err = cron.ValidatedCronSpec(toml)
	require.NoError(t, err)

	jb, err = app.AddJobV2(context.Background(), jb, jb.Name)
	require.Error(t, err, "augur-sportsdataio: no such bridge exists")
}

func TestClient_MigrateRunLog(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("migrate", 0)
	set.Parse([]string{"../testdata/jsonspecs/example-runlog.json", "../testdata/jsonspecs/requesters.txt"})
	c := cli.NewContext(nil, set, nil)

	toml, j, err := client.MigrateJobSpecForResult(c)
	require.NoError(t, err)

	fmt.Println(toml)

	_, err = job.ValidateSpec(toml)
	require.NoError(t, err)

	var jb job.Job
	jb, err = directrequest.ValidatedDirectRequestSpec(toml)
	require.NoError(t, err)

	require.Equal(t, "0xfe8F390fFD3c74870367121cE251C744d3DC01Ed", jb.DirectRequestSpec.ContractAddress.String())
	require.Equal(t, clnull.Uint32From(10), jb.DirectRequestSpec.MinIncomingConfirmations)
	require.Equal(t, fmt.Sprintf(
		`Name = "QDT Price Prediction"
SchemaVersion = 1
Type = "directrequest"
contractAddress = "0xfe8F390fFD3c74870367121cE251C744d3DC01Ed"
externalJobID = "%v"
minIncomingConfirmations = "10"
observationSource = """
decode_log [
	abi="OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes32 data)"
	data="$(jobRun.logData)"
	topics="$(jobRun.logTopics)"
	type=ethabidecodelog
	];
	merge_1 [
	right=<{"endpoint":"price"}>
	type=merge
	];
	send_to_bridge_1 [
	name=qdt
	requestData=<{ "data": $(merge_1) }>
	type=bridge
	];
	multiply_1 [
	times=100000000
	type=multiply
	];
	encode_data_4 [
	abi="(uint256 value)"
	data=<{ "value": $(multiply_1) }>
	type=ethabiencode
	];
	encode_tx_4 [
	abi="fulfillOracleRequest(bytes32 requestId, uint256 payment, address callbackAddress, bytes4 callbackFunctionId, uint256 expiration, bytes32 calldata data)"
	data=<{
"requestId":          $(decode_log.requestId),
"payment":            $(decode_log.payment),
"callbackAddress":    $(decode_log.callbackAddr),
"callbackFunctionId": $(decode_log.callbackFunctionId),
"expiration":         $(decode_log.cancelExpiration),
"data":               $(encode_data_4)
}
>
	type=ethabiencode
	];
	send_tx_4 [
	data="$(encode_tx_4)"
	to="0xfe8F390fFD3c74870367121cE251C744d3DC01Ed"
	type=ethtx
	];
	
	// Edge definitions.
	decode_log -> merge_1;
	merge_1 -> send_to_bridge_1;
	send_to_bridge_1 -> multiply_1;
	multiply_1 -> encode_data_4;
	encode_data_4 -> encode_tx_4;
	encode_tx_4 -> send_tx_4;
	"""
requesters = ["0x22DE85B0cD5B3684865ECfEedfBAF12777cd0Ff8", "0x22DE85B0cD5B3684865ECfEedfBAF12777cd0Ff8"]
`, j.ExternalJobID), toml)

	jb, err = app.AddJobV2(context.Background(), jb, jb.Name)
	require.Error(t, err, "augur-sportsdataio: no such bridge exists")
}
