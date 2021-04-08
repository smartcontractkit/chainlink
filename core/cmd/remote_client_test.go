package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/smartcontractkit/chainlink/core/services/eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
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

	for k, v := range sopts.Config {
		config.Set(k, v)
	}

	var app *cltest.TestApplication
	var cleanup func()
	if sopts.WithKey {
		app, cleanup = cltest.NewApplicationWithConfigAndKey(t, config, sopts.FlagsAndDeps...)
	} else {
		app, cleanup = cltest.NewApplicationWithConfig(t, config, sopts.FlagsAndDeps...)
	}
	t.Cleanup(cleanup)

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

func newEthMocks(t *testing.T) (*mocks.RPCClient, *mocks.GethClient) {
	t.Helper()

	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	t.Cleanup(assertMocksCalled)

	return rpcClient, gethClient
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
			assert.NoError(t, err)

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

			set := flag.NewFlagSet("create", 0)
			assert.NoError(t, set.Parse(test.args))
			c := cli.NewContext(nil, set, nil)

			err := client.CreateExternalInitiator(c)
			assert.Error(t, err)

			exis := cltest.AllExternalInitiators(t, app.Store)
			assert.Len(t, exis, 0)
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
		{"file", "../internal/fixtures/web/end_at_job.json", 3, false},
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

	eim := new(mocks.ExternalInitiatorManager)
	app := startNewApplication(t, withMocks(eim))
	client, _ := app.NewClientAndRenderer()

	job := cltest.NewJob()
	require.NoError(t, app.Store.CreateJob(&job))

	set := flag.NewFlagSet("archive", 0)
	set.Parse([]string{job.ID.String()})
	c := cli.NewContext(nil, set, nil)

	eim.On("DeleteJob", mock.Anything, mock.MatchedBy(func(id models.JobID) bool {
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

func TestClient_CreateBridge(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	tests := []struct {
		name    string
		param   string
		errored bool
	}{
		{"EmptyString", "", true},
		{"ValidString", `{ "name": "TestBridge", "url": "http://localhost:3000/randomNumber" }`, false},
		{"InvalidString", `{ "noname": "", "nourl": "" }`, true},
		{"InvalidChar", `{ "badname": "path/bridge", "nourl": "" }`, true},
		{"ValidPath", "testdata/create_random_number_bridge_type.json", false},
		{"InvalidPath", "bad/filepath/", true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {

			set := flag.NewFlagSet("bridge", 0)
			set.Parse([]string{test.param})
			c := cli.NewContext(nil, set, nil)
			if test.errored {
				assert.Error(t, client.CreateBridge(c))
			} else {
				assert.Nil(t, client.CreateBridge(c))
			}
		})
	}
}

func TestClient_IndexBridges(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	bt1 := &models.BridgeType{
		Name:          models.MustNewTaskType("testingbridges1"),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	err := app.GetStore().CreateBridgeType(bt1)
	require.NoError(t, err)

	bt2 := &models.BridgeType{
		Name:          models.MustNewTaskType("testingbridges2"),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	err = app.GetStore().CreateBridgeType(bt2)
	require.NoError(t, err)

	require.Nil(t, client.IndexBridges(cltest.EmptyCLIContext()))
	bridges := *r.Renders[0].(*[]models.BridgeType)
	require.Equal(t, 2, len(bridges))
	assert.Equal(t, bt1.Name, bridges[0].Name)
}

func TestClient_ShowBridge(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	bt := &models.BridgeType{
		Name:          models.MustNewTaskType("testingbridges1"),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	require.NoError(t, app.GetStore().CreateBridgeType(bt))

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{bt.Name.String()})
	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.ShowBridge(c))
	require.Len(t, r.Renders, 1)
	assert.Equal(t, bt.Name, r.Renders[0].(*models.BridgeType).Name)
}

func TestClient_RemoveBridge(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	bt := &models.BridgeType{
		Name:          models.MustNewTaskType("testingbridges1"),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	err := app.GetStore().CreateBridgeType(bt)
	require.NoError(t, err)

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{bt.Name.String()})
	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.RemoveBridge(c))
	require.Len(t, r.Renders, 1)
	assert.Equal(t, bt.Name, r.Renders[0].(*models.BridgeType).Name)
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

func TestClient_SendEther_From_BPTXM(t *testing.T) {
	t.Parallel()

	rpcClient, gethClient := newEthMocks(t)
	oca := common.HexToAddress("0xDEADB3333333F")
	app := startNewApplication(t,
		withKey(),
		withConfig(map[string]interface{}{
			"OPERATOR_CONTRACT_ADDRESS": &oca,
		}),
		withMocks(eth.NewClientWith(rpcClient, gethClient)),
		startAndConnect(),
	)
	client, _ := app.NewClientAndRenderer()
	s := app.GetStore()

	set := flag.NewFlagSet("sendether", 0)
	amount := "100.5"
	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, s, 0)
	to := "0x342156c8d3bA54Abc67920d35ba1d1e67201aC9C"
	set.Parse([]string{amount, fromAddress.Hex(), to})

	cliapp := cli.NewApp()
	c := cli.NewContext(cliapp, set, nil)

	assert.NoError(t, client.SendEther(c))

	etx := models.EthTx{}
	require.NoError(t, s.DB.First(&etx).Error)
	require.Equal(t, "100.500000000000000000", etx.Value.String())
	require.Equal(t, fromAddress, etx.FromAddress)
	require.Equal(t, to, etx.ToAddress.Hex())
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
		ChangePasswordRequest: models.ChangePasswordRequest{
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

func TestClient_IndexTransactions(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	store := app.GetStore()
	_, from := cltest.MustAddRandomKeyToKeystore(t, store)

	tx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 0, 1, from)
	attempt := tx.EthTxAttempts[0]

	// page 1
	set := flag.NewFlagSet("test transactions", 0)
	set.Int("page", 1, "doc")
	c := cli.NewContext(nil, set, nil)
	require.Equal(t, 1, c.Int("page"))
	assert.NoError(t, client.IndexTransactions(c))

	renderedTxs := *r.Renders[0].(*[]presenters.EthTx)
	assert.Equal(t, 1, len(renderedTxs))
	assert.Equal(t, attempt.Hash.Hex(), renderedTxs[0].Hash.Hex())

	// page 2 which doesn't exist
	set = flag.NewFlagSet("test txattempts", 0)
	set.Int("page", 2, "doc")
	c = cli.NewContext(nil, set, nil)
	require.Equal(t, 2, c.Int("page"))
	assert.NoError(t, client.IndexTransactions(c))

	renderedTxs = *r.Renders[1].(*[]presenters.EthTx)
	assert.Equal(t, 0, len(renderedTxs))
}

func TestClient_ShowTransaction(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	store := app.GetStore()
	_, from := cltest.MustAddRandomKeyToKeystore(t, store)

	tx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 0, 1, from)
	attempt := tx.EthTxAttempts[0]

	set := flag.NewFlagSet("test get tx", 0)
	set.Parse([]string{attempt.Hash.Hex()})
	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.ShowTransaction(c))

	renderedTx := *r.Renders[0].(*presenters.EthTx)
	assert.Equal(t, &tx.FromAddress, renderedTx.From)
}

func TestClient_IndexTxAttempts(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	store := app.GetStore()
	_, from := cltest.MustAddRandomKeyToKeystore(t, store)

	tx := cltest.MustInsertConfirmedEthTxWithAttempt(t, store, 0, 1, from)

	// page 1
	set := flag.NewFlagSet("test txattempts", 0)
	set.Int("page", 1, "doc")
	c := cli.NewContext(nil, set, nil)
	require.Equal(t, 1, c.Int("page"))
	assert.NoError(t, client.IndexTxAttempts(c))

	renderedAttempts := *r.Renders[0].(*[]presenters.EthTx)
	require.Len(t, tx.EthTxAttempts, 1)
	assert.Equal(t, tx.EthTxAttempts[0].Hash.Hex(), renderedAttempts[0].Hash.Hex())

	// page 2 which doesn't exist
	set = flag.NewFlagSet("test transactions", 0)
	set.Int("page", 2, "doc")
	c = cli.NewContext(nil, set, nil)
	require.Equal(t, 2, c.Int("page"))
	assert.NoError(t, client.IndexTxAttempts(c))

	renderedAttempts = *r.Renders[1].(*[]presenters.EthTx)
	assert.Equal(t, 0, len(renderedAttempts))
}

func TestClient_SetMinimumGasPrice(t *testing.T) {
	t.Parallel()

	// Setup Withdrawals application
	rpcClient, gethClient := newEthMocks(t)
	oca := common.HexToAddress("0xDEADB3333333F")
	app := startNewApplication(t,
		withKey(),
		withConfig(map[string]interface{}{
			"OPERATOR_CONTRACT_ADDRESS": &oca,
		}),
		withMocks(eth.NewClientWith(rpcClient, gethClient)),
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
	tree, err := toml.LoadFile("testdata/oracle-spec.toml")
	require.NoError(t, err)
	err = tree.Unmarshal(&ocrJobSpecFromFile)
	require.NoError(t, err)
	var ocrSpec job.OffchainReportingOracleSpec
	err = tree.Unmarshal(&ocrSpec)
	require.NoError(t, err)
	ocrJobSpecFromFile.OffchainreportingOracleSpec = &ocrSpec
	key := cltest.MustInsertRandomKey(t, app.Store.DB)
	ocrJobSpecFromFile.OffchainreportingOracleSpec.TransmitterAddress = &key.Address

	jobID, _ := app.AddJobV2(context.Background(), ocrJobSpecFromFile, null.String{})

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{strconv.FormatInt(int64(jobID), 10)})
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
	assert.EqualError(t, client.TriggerPipelineRun(c), "parseResponse error: Error; no job found with id 1 (most likely it was deleted)")
}

func TestClient_ListJobsV2(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, r := app.NewClientAndRenderer()

	// Create the job
	toml, err := ioutil.ReadFile("./testdata/direct-request-spec.toml")
	assert.NoError(t, err)

	request, err := json.Marshal(models.CreateJobSpecRequest{
		TOML: string(toml),
	})
	assert.NoError(t, err)

	resp, err := client.HTTP.Post("/v2/jobs", bytes.NewReader(request))
	assert.NoError(t, err)

	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	job := cmd.Job{}
	err = web.ParseJSONAPIResponse(responseBodyBytes, &job)
	assert.NoError(t, err)

	require.Nil(t, client.ListJobsV2(cltest.EmptyCLIContext()))
	jobs := *r.Renders[0].(*[]cmd.Job)
	require.Equal(t, 1, len(jobs))
	assert.Equal(t, job.ID, jobs[0].ID)
}

func TestClient_CreateJobV2(t *testing.T) {
	t.Parallel()

	app := startNewApplication(t)
	client, _ := app.NewClientAndRenderer()

	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.Parse([]string{"./testdata/ocr-bootstrap-spec.toml"})
	err := client.CreateJobV2(cli.NewContext(nil, fs, nil))
	require.NoError(t, err)
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
