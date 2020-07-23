package cmd_test

import (
	"flag"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestClient_DisplayAccountBalance(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	app.EthMock.Register("eth_call", "0x0100")

	require.NoError(t, app.Start())

	client, r := app.NewClientAndRenderer()

	assert.Nil(t, client.DisplayAccountBalance(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))
	from := cltest.GetAccountAddress(t, app.GetStore())
	balances := *r.Renders[0].(*[]presenters.AccountBalance)
	assert.Equal(t, from.Hex(), balances[0].Address)
}

func TestClient_IndexJobSpecs(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	j1 := cltest.NewJob()
	app.Store.CreateJob(&j1)
	j2 := cltest.NewJob()
	app.Store.CreateJob(&j2)

	client, r := app.NewClientAndRenderer()

	require.Nil(t, client.IndexJobSpecs(cltest.EmptyCLIContext()))
	jobs := *r.Renders[0].(*[]models.JobSpec)
	assert.Equal(t, 2, len(jobs))
	assert.Equal(t, j1.ID, jobs[0].ID)
}

func TestClient_ShowJobRun_Exists(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.NewJobWithWebInitiator()
	assert.NoError(t, app.Store.CreateJob(&j))

	jr := cltest.CreateJobRunViaWeb(t, app, j, `{"result":"100"}`)

	client, r := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{jr.ID.String()})
	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.ShowJobRun(c))
	assert.Equal(t, 1, len(r.Renders))
	assert.Equal(t, jr.ID, r.Renders[0].(*presenters.JobRun).ID)
}

func TestClient_ShowJobRun_NotFound(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client, r := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"bogus-ID"})
	c := cli.NewContext(nil, set, nil)
	assert.Error(t, client.ShowJobRun(c))
	assert.Empty(t, r.Renders)
}

func TestClient_IndexJobRuns(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	j := cltest.NewJobWithWebInitiator()
	assert.NoError(t, app.Store.CreateJob(&j))

	jr0 := cltest.NewJobRun(j)
	jr0.Result.Data = cltest.JSONFromString(t, `{"a":"b"}`)
	require.NoError(t, app.Store.CreateJobRun(&jr0))
	jr1 := cltest.NewJobRun(j)
	jr1.Result.Data = cltest.JSONFromString(t, `{"x":"y"}`)
	require.NoError(t, app.Store.CreateJobRun(&jr1))

	client, r := app.NewClientAndRenderer()

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

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	job := cltest.NewJob()
	app.Store.CreateJob(&job)

	client, r := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{job.ID.String()})
	c := cli.NewContext(nil, set, nil)
	require.Nil(t, client.ShowJobSpec(c))
	require.Equal(t, 1, len(r.Renders))
	assert.Equal(t, job.ID, r.Renders[0].(*presenters.JobSpec).ID)
}

func TestClient_ShowJobSpec_NotFound(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client, r := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"bogus-ID"})
	c := cli.NewContext(nil, set, nil)
	assert.Error(t, client.ShowJobSpec(c))
	assert.Empty(t, r.Renders)
}

var EndAt = time.Now().AddDate(0, 10, 0).Round(time.Second).UTC()

func TestClient_CreateServiceAgreement(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client, _ := app.NewClientAndRenderer()

	sa := cltest.MustHelloWorldAgreement(t)
	endAtISO8601 := EndAt.Format(time.RFC3339)
	sa = strings.Replace(sa, "2019-10-19T22:17:19Z", endAtISO8601, 1)
	tmpFile, err := ioutil.TempFile("", "sa.*.json")
	require.NoError(t, err, "while opening temp file for modified service agreement")
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString(sa)

	tests := []struct {
		name        string
		input       string
		jobsCreated bool
		errored     bool
	}{
		{"invalid json", "{bad son}", false, true},
		{"bad file path", "bad/filepath/", false, true},
		{"valid service agreement", string(sa), true, false},
		{"service agreement specified as path", tmpFile.Name(), true, false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {

			set := flag.NewFlagSet("create", 0)
			assert.NoError(t, set.Parse([]string{test.input}))
			c := cli.NewContext(nil, set, nil)

			err := client.CreateServiceAgreement(c)

			cltest.AssertError(t, test.errored, err)
			jobs := cltest.AllJobs(t, app.Store)
			if test.jobsCreated {
				assert.True(t, len(jobs) > 0)
			} else {
				assert.Equal(t, 0, len(jobs))
			}
		})
	}
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
			app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
			defer cleanup()
			require.NoError(t, app.Start())

			client, _ := app.NewClientAndRenderer()

			set := flag.NewFlagSet("create", 0)
			assert.NoError(t, set.Parse(test.args))
			c := cli.NewContext(nil, set, nil)

			err := client.CreateExternalInitiator(c)
			assert.NoError(t, err)

			var exi models.ExternalInitiator
			err = app.Store.RawDB(func(db *gorm.DB) error {
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
			app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
			defer cleanup()
			require.NoError(t, app.Start())

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

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	token := auth.NewToken()
	exi, err := models.NewExternalInitiator(token,
		&models.ExternalInitiatorRequest{Name: "name"},
	)
	require.NoError(t, err)
	err = app.Store.CreateExternalInitiator(exi)
	require.NoError(t, err)

	client, r := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{exi.Name})
	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.DeleteExternalInitiator(c))
	assert.Empty(t, r.Renders)
}

func TestClient_DestroyExternalInitiator_NotFound(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	client, r := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"bogus-ID"})
	c := cli.NewContext(nil, set, nil)
	assert.Error(t, client.DeleteExternalInitiator(c))
	assert.Empty(t, r.Renders)
}

func TestClient_CreateJobSpec(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())
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

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	job := cltest.NewJob()
	require.NoError(t, app.Store.CreateJob(&job))

	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("archive", 0)
	set.Parse([]string{job.ID.String()})
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.ArchiveJobSpec(c))

	jobs := cltest.AllJobs(t, app.Store)
	require.Len(t, jobs, 0)
}

func TestClient_CreateJobSpec_JSONAPIErrors(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

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

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

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

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

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

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

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

	client, r := app.NewClientAndRenderer()

	require.Nil(t, client.IndexBridges(cltest.EmptyCLIContext()))
	bridges := *r.Renders[0].(*[]models.BridgeType)
	require.Equal(t, 2, len(bridges))
	assert.Equal(t, bt1.Name, bridges[0].Name)
}

func TestClient_ShowBridge(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.StartAndConnect())

	bt := &models.BridgeType{
		Name:          models.MustNewTaskType("testingbridges1"),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	require.NoError(t, app.GetStore().CreateBridgeType(bt))

	client, r := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{bt.Name.String()})
	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.ShowBridge(c))
	require.Len(t, r.Renders, 1)
	assert.Equal(t, bt.Name, r.Renders[0].(*models.BridgeType).Name)
}

func TestClient_RemoveBridge(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	bt := &models.BridgeType{
		Name:          models.MustNewTaskType("testingbridges1"),
		URL:           cltest.WebURL(t, "https://testing.com/bridges"),
		Confirmations: 0,
	}
	err := app.GetStore().CreateBridgeType(bt)
	require.NoError(t, err)

	client, r := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{bt.Name.String()})
	c := cli.NewContext(nil, set, nil)
	require.NoError(t, client.RemoveBridge(c))
	require.Len(t, r.Renders, 1)
	assert.Equal(t, bt.Name, r.Renders[0].(*models.BridgeType).Name)
}

func TestClient_RemoteLogin(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

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

func setupWithdrawalsApplication(t *testing.T, config *cltest.TestConfig) (*cltest.TestApplication, func()) {
	oca := common.HexToAddress("0xDEADB3333333F")
	config.Set("ORACLE_CONTRACT_ADDRESS", &oca)
	app, cleanup := cltest.NewApplicationWithConfigAndKey(t, config,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	return app, cleanup
}

func TestClient_SendEther_From_LegacyTxManager(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig(t)
	config.Set("ENABLE_BULLETPROOF_TX_MANAGER", "false")
	defer cleanup()
	app, cleanup := setupWithdrawalsApplication(t, config)
	defer cleanup()
	app.EthMock.Register("eth_getTransactionCount", "0x100")

	require.NoError(t, app.StartAndConnect())

	client, _ := app.NewClientAndRenderer()
	set := flag.NewFlagSet("sendether", 0)
	set.Parse([]string{"100", app.Store.TxManager.NextActiveAccount().Address.String(), "0x342156c8d3bA54Abc67920d35ba1d1e67201aC9C"})

	app.EthMock.Context("manager.CreateTx#1", func(ethMock *cltest.EthMock) {
		ethMock.Register("eth_sendRawTransaction", cltest.NewHash())
	})

	cliapp := cli.NewApp()
	c := cli.NewContext(cliapp, set, nil)

	assert.NoError(t, client.SendEther(c))
}

func TestClient_SendEther_From_BPTXM(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig(t)
	config.Set("ENABLE_BULLETPROOF_TX_MANAGER", "true")
	defer cleanup()
	app, cleanup := setupWithdrawalsApplication(t, config)
	defer cleanup()
	s := app.GetStore()

	require.NoError(t, app.StartAndConnect())

	client, _ := app.NewClientAndRenderer()
	set := flag.NewFlagSet("sendether", 0)
	amount := "100.5"
	from := cltest.GetDefaultFromAddress(t, s)
	to := "0x342156c8d3bA54Abc67920d35ba1d1e67201aC9C"
	set.Parse([]string{amount, from.Hex(), to})

	cliapp := cli.NewApp()
	c := cli.NewContext(cliapp, set, nil)

	assert.NoError(t, client.SendEther(c))

	etx := models.EthTx{}
	require.NoError(t, s.DB.First(&etx).Error)
	require.Equal(t, "100.500000000000000000", etx.Value.String())
	require.Equal(t, from, etx.FromAddress)
	require.Equal(t, to, etx.ToAddress.Hex())
}

func TestClient_ChangePassword(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

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
			NewPassword: "password",
		},
	}
	err = client.ChangePassword(cli.NewContext(nil, nil, nil))
	assert.NoError(t, err)

	// otherClient should now be logged out
	err = otherClient.IndexBridges(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "401 Unauthorized")
}

func TestClient_IndexTransactions(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	store := app.GetStore()
	from := cltest.GetAccountAddress(t, store)
	tx := cltest.CreateTx(t, store, from, 1)

	client, r := app.NewClientAndRenderer()

	// page 1
	set := flag.NewFlagSet("test transactions", 0)
	set.Int("page", 1, "doc")
	c := cli.NewContext(nil, set, nil)
	require.Equal(t, 1, c.Int("page"))
	assert.NoError(t, client.IndexTransactions(c))

	renderedTxs := *r.Renders[0].(*[]presenters.Tx)
	assert.Equal(t, 1, len(renderedTxs))
	assert.Equal(t, tx.Hash.Hex(), renderedTxs[0].Hash.Hex())

	// page 2 which doesn't exist
	set = flag.NewFlagSet("test txattempts", 0)
	set.Int("page", 2, "doc")
	c = cli.NewContext(nil, set, nil)
	require.Equal(t, 2, c.Int("page"))
	assert.NoError(t, client.IndexTransactions(c))

	renderedTxs = *r.Renders[1].(*[]presenters.Tx)
	assert.Equal(t, 0, len(renderedTxs))
}

func TestClient_ShowTransaction(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	store := app.GetStore()
	from := cltest.GetAccountAddress(t, store)
	tx := cltest.CreateTx(t, store, from, 1)

	client, r := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test get tx", 0)
	set.Parse([]string{tx.Hash.Hex()})
	c := cli.NewContext(nil, set, nil)
	assert.NoError(t, client.ShowTransaction(c))

	renderedTx := *r.Renders[0].(*presenters.Tx)
	assert.Equal(t, &tx.From, renderedTx.From)
}

func TestClient_IndexTxAttempts(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	store := app.GetStore()
	from := cltest.GetAccountAddress(t, store)
	tx := cltest.CreateTx(t, store, from, 1)

	client, r := app.NewClientAndRenderer()

	// page 1
	set := flag.NewFlagSet("test txattempts", 0)
	set.Int("page", 1, "doc")
	c := cli.NewContext(nil, set, nil)
	require.Equal(t, 1, c.Int("page"))
	assert.NoError(t, client.IndexTxAttempts(c))

	renderedAttempts := *r.Renders[0].(*[]models.TxAttempt)
	require.Len(t, tx.Attempts, 1)
	assert.Equal(t, tx.Attempts[0].Hash.Hex(), renderedAttempts[0].Hash.Hex())

	// page 2 which doesn't exist
	set = flag.NewFlagSet("test transactions", 0)
	set.Int("page", 2, "doc")
	c = cli.NewContext(nil, set, nil)
	require.Equal(t, 2, c.Int("page"))
	assert.NoError(t, client.IndexTxAttempts(c))

	renderedAttempts = *r.Renders[1].(*[]models.TxAttempt)
	assert.Equal(t, 0, len(renderedAttempts))
}

func TestClient_CreateExtraKey(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	kst := app.Store.KeyStore.(*mocks.KeyStoreInterface)
	defer cleanup()
	require.NoError(t, app.Start())

	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("test", 0)
	set.String("file", "internal/fixtures/apicredentials", "")
	c := cli.NewContext(nil, set, nil)
	err := client.RemoteLogin(c)
	assert.NoError(t, err)

	client.PasswordPrompter = cltest.MockPasswordPrompter{Password: "password"}

	kst.On("Unlock", cltest.Password).Return(nil)
	kst.On("NewAccount", cltest.Password).Return(accounts.Account{}, nil)
	assert.NoError(t, client.CreateExtraKey(c))

	kst.AssertExpectations(t)
}

func TestClient_SetMinimumGasPrice(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	app, cleanup := setupWithdrawalsApplication(t, config)
	defer cleanup()
	app.EthMock.Register("eth_getTransactionCount", "0x100")
	require.NoError(t, app.StartAndConnect())

	client, _ := app.NewClientAndRenderer()
	set := flag.NewFlagSet("setgasprice", 0)
	set.Parse([]string{"8616460799"})

	c := cli.NewContext(nil, set, nil)

	// app.EthMock.Register("eth_call", "0xDE0B6B3A7640000")
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

	app, cleanup := cltest.NewApplicationWithKey(t, cltest.LenientEthMock)
	defer cleanup()
	require.NoError(t, app.Start())

	client, r := app.NewClientAndRenderer()
	assert.NoError(t, client.GetConfiguration(cltest.EmptyCLIContext()))
	require.Equal(t, 1, len(r.Renders))

	cwl := *r.Renders[0].(*presenters.ConfigWhitelist)
	assert.Equal(t, cwl.Whitelist.BridgeResponseURL, app.Config.BridgeResponseURL().String())
	assert.Equal(t, cwl.Whitelist.ChainID, app.Config.ChainID())
	assert.Equal(t, cwl.Whitelist.Dev, app.Config.Dev())
	assert.Equal(t, cwl.Whitelist.EthGasBumpThreshold, app.Config.EthGasBumpThreshold())
	assert.Equal(t, cwl.Whitelist.LogLevel, app.Config.LogLevel())
	assert.Equal(t, cwl.Whitelist.LogSQLStatements, app.Config.LogSQLStatements())
	assert.Equal(t, cwl.Whitelist.MinIncomingConfirmations, app.Config.MinIncomingConfirmations())
	assert.Equal(t, cwl.Whitelist.MinRequiredOutgoingConfirmations, app.Config.MinRequiredOutgoingConfirmations())
	assert.Equal(t, cwl.Whitelist.MinimumContractPayment, app.Config.MinimumContractPayment())
	assert.Equal(t, cwl.Whitelist.RootDir, app.Config.RootDir())
	assert.Equal(t, cwl.Whitelist.SessionTimeout, app.Config.SessionTimeout())
}

func TestClient_CancelJobRun(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication(t,
		cltest.LenientEthMock,
		cltest.EthMockRegisterChainID,
		cltest.EthMockRegisterGetBalance,
	)
	defer cleanup()
	require.NoError(t, app.Start())

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, app.Store.CreateJob(&job))
	run := cltest.NewJobRun(job)
	require.NoError(t, app.Store.CreateJobRun(&run))

	client, _ := app.NewClientAndRenderer()

	set := flag.NewFlagSet("cancel", 0)
	set.Parse([]string{run.ID.String()})
	c := cli.NewContext(nil, set, nil)

	require.NoError(t, client.CancelJobRun(c))

	runs := cltest.MustAllJobsWithStatus(t, app.Store, models.RunStatusCancelled)
	require.Len(t, runs, 1)
	assert.Equal(t, models.RunStatusCancelled, runs[0].GetStatus())
	assert.NotNil(t, runs[0].FinishedAt)
}
