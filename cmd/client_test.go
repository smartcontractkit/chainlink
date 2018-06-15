package cmd_test

import (
	"flag"
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	null "gopkg.in/guregu/null.v3"
)

func TestClient_RunNode(t *testing.T) {
	app, _ := cltest.NewApplicationWithKeyStore() // cleanup invoked in client.RunNode
	eth := app.MockEthClient()
	eth.Register("eth_getTransactionCount", `0x1`)

	r := &cltest.RendererMock{}
	var called bool
	auth := cltest.CallbackAuthenticator{Callback: func(*store.Store, string) { called = true }}
	client := cmd.Client{
		Renderer:   r,
		Config:     app.Store.Config,
		AppFactory: cltest.InstanceAppFactory{App: app},
		Auth:       auth,
		Runner:     cltest.EmptyRunner{}}

	set := flag.NewFlagSet("test", 0)
	set.Bool("debug", true, "")
	c := cli.NewContext(nil, set, nil)

	assert.Nil(t, client.RunNode(c))
	assert.True(t, called)
}

func TestClient_DisplayAccountBalance(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKeyStore()
	defer cleanup()
	account, err := app.Store.KeyStore.GetAccount()
	assert.NoError(t, err)

	ethMock := app.MockEthClient()
	ethMock.Register("eth_getBalance", "0x0100")
	ethMock.Register("eth_call", "0x0100")

	client, r := cltest.NewClientAndRenderer(app.Store.Config)

	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, client.DisplayAccountBalance(c))
	assert.Equal(t, 1, len(r.Renders))
	assert.Equal(t, account.Address.Hex(), r.Renders[0].(*presenters.AccountBalance).Address)
}

func TestClient_GetJobSpecs(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j1 := cltest.NewJob()
	app.Store.SaveJob(&j1)
	j2 := cltest.NewJob()
	app.Store.SaveJob(&j2)

	client, r := cltest.NewClientAndRenderer(app.Store.Config)

	assert.Nil(t, client.GetJobSpecs(nil))
	jobs := *r.Renders[0].(*[]models.JobSpec)
	assert.Equal(t, 2, len(jobs))
	assert.Equal(t, j1.ID, jobs[0].ID)
}

func TestClient_ShowJobSpec_Exists(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	job := cltest.NewJob()
	app.Store.SaveJob(&job)

	client, r := cltest.NewClientAndRenderer(app.Store.Config)

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{job.ID})
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, client.ShowJobSpec(c))
	assert.Equal(t, 1, len(r.Renders))
	assert.Equal(t, job.ID, r.Renders[0].(*presenters.JobSpec).ID)
}

func TestClient_ShowJobSpec_NotFound(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	client, r := cltest.NewClientAndRenderer(app.Store.Config)

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"bogus-ID"})
	c := cli.NewContext(nil, set, nil)
	assert.Error(t, client.ShowJobSpec(c))
	assert.Empty(t, r.Renders)
}

func TestClient_CreateJobSpec(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client, _ := cltest.NewClientAndRenderer(app.Store.Config)

	tests := []struct {
		input   string
		nJobs   int
		errored bool
	}{
		{"{bad son}", 0, true},
		{"bad/filepath/", 0, true},
		{`{"initiators":[{"type":"web"}],"tasks":[{"type":"NoOp"}]}`, 1, false},
		{`{"initiators":[{"type":"runAt","time":"2018-01-08T18:12:01.103Z"}],"tasks":[{"type":"NoOp"}]}`, 2, false},
		{"../internal/fixtures/web/end_at_job.json", 3, false},
	}

	for _, tt := range tests {
		test := tt

		set := flag.NewFlagSet("create", 0)
		set.Parse([]string{test.input})
		c := cli.NewContext(nil, set, nil)
		if test.errored {
			assert.Error(t, client.CreateJobSpec(c))
		} else {
			assert.Nil(t, client.CreateJobSpec(c))
		}
		numberOfJobs, _ := app.Store.Jobs()
		assert.Equal(t, test.nJobs, len(numberOfJobs))
	}
}

func TestClient_CreateJobRun(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client, _ := cltest.NewClientAndRenderer(app.Store.Config)

	tests := []struct {
		name    string
		json    string
		jobSpec models.JobSpec
		errored bool
	}{
		{"CreateSuccess", `{"value": 100}`, first(cltest.NewJobWithWebInitiator()), false},
		{"EmptyBody", ``, first(cltest.NewJobWithWebInitiator()), false},
		{"InvalidBody", `{`, first(cltest.NewJobWithWebInitiator()), true},
		{"WithoutWebInitiator", ``, first(cltest.NewJobWithLogInitiator()), true},
		{"NotFound", ``, first(cltest.NewJobWithWebInitiator()), true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			assert.Nil(t, app.Store.SaveJob(&test.jobSpec))

			args := make([]string, 1)
			args[0] = test.jobSpec.ID
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

func TestClient_AddBridge(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client, _ := cltest.NewClientAndRenderer(app.Store.Config)

	tests := []struct {
		name    string
		param   string
		errored bool
	}{
		{"EmptyString", "", true},
		{"ValidString", `{ "name": "TestBridge", "url": "http://localhost:3000/randomNumber" }`, false},
		{"InvalidString", `{ "noname": "", "nourl": "" }`, true},
		{"InvalidChar", `{ "badname": "path/bridge", "nourl": "" }`, true},
		{"ValidPath", "../internal/fixtures/web/create_random_number_bridge_type.json", false},
		{"InvalidPath", "bad/filepath/", true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {

			set := flag.NewFlagSet("bridge", 0)
			set.Parse([]string{test.param})
			c := cli.NewContext(nil, set, nil)
			if test.errored {
				assert.Error(t, client.AddBridge(c))
			} else {
				assert.Nil(t, client.AddBridge(c))
			}
		})
	}
}

func TestClient_GetBridges(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	bt1 := &models.BridgeType{Name: "testingbridges1",
		URL:                  cltest.WebURL("https://testing.com/bridges"),
		DefaultConfirmations: 0}
	app.AddAdapter(bt1)

	bt2 := &models.BridgeType{Name: "testingbridges2",
		URL:                  cltest.WebURL("https://testing.com/bridges"),
		DefaultConfirmations: 0}
	app.AddAdapter(bt2)

	client, r := cltest.NewClientAndRenderer(app.Store.Config)

	assert.Nil(t, client.GetBridges(nil))
	bridges := *r.Renders[0].(*[]models.BridgeType)
	assert.Equal(t, 2, len(bridges))
	assert.Equal(t, bt1.Name, bridges[0].Name)
}

func TestClient_ShowBridge(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	bt := &models.BridgeType{Name: "testingbridges1",
		URL:                  cltest.WebURL("https://testing.com/bridges"),
		DefaultConfirmations: 0}
	app.AddAdapter(bt)

	client, r := cltest.NewClientAndRenderer(app.Store.Config)

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{bt.Name})
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, client.ShowBridge(c))
	assert.Equal(t, 1, len(r.Renders))
	assert.Equal(t, bt.Name, r.Renders[0].(*models.BridgeType).Name)
}

func TestClient_RemoveBridge(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	bt := &models.BridgeType{Name: "testingbridges1",
		URL:                  cltest.WebURL("https://testing.com/bridges"),
		DefaultConfirmations: 0}
	app.AddAdapter(bt)

	client, r := cltest.NewClientAndRenderer(app.Store.Config)

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{bt.Name})
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, client.RemoveBridge(c))
	assert.Equal(t, 1, len(r.Renders))
	assert.Equal(t, bt.Name, r.Renders[0].(*models.BridgeType).Name)
}

func TestClient_BackupDatabase(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client, _ := cltest.NewClientAndRenderer(app.Store.Config)

	job := cltest.NewJob()
	assert.Nil(t, app.Store.SaveJob(&job))

	set := flag.NewFlagSet("backupset", 0)
	path := path.Join(app.Store.Config.RootDir, "backup.bolt")
	set.Parse([]string{path})
	c := cli.NewContext(nil, set, nil)

	err := client.BackupDatabase(c)
	assert.NoError(t, err)

	restored := models.NewORM(path)
	restoredJob, err := restored.FindJob(job.ID)
	assert.NoError(t, err)

	reloaded, err := app.Store.FindJob(job.ID)
	assert.NoError(t, err)
	assert.Equal(t, reloaded, restoredJob)
}

func TestClient_ImportKey(t *testing.T) {
	t.Parallel()

	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client, _ := cltest.NewClientAndRenderer(app.Store.Config)

	os.MkdirAll(app.Store.Config.KeysDir(), os.FileMode(0700))

	set := flag.NewFlagSet("import", 0)
	set.Parse([]string{"../internal/fixtures/keys/3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea.json"})
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, client.ImportKey(c))
	assert.Error(t, client.ImportKey(c))
}

func TestClient_DeleteQuery(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	client, _ := cltest.NewClientAndRenderer(app.Store.Config)

	for i := 0; i < 3; i++ {
		bt := &models.BridgeType{Name: fmt.Sprintf("testbridge%v", i),
			URL:                  cltest.WebURL("http://www.example.com"),
			DefaultConfirmations: 0}
		app.AddAdapter(bt)
	}

	set := flag.NewFlagSet("delete", 0)
	set.Parse([]string{"../internal/fixtures/web/delete_query_bridges.json"})
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, client.DeleteQuery(c))

}

func TestClient_Prune(t *testing.T) {

	entries := []struct {
		name        string
		status      models.RunStatus
		completedAt null.Time
	}{
		{
			"OldCompleted",
			models.RunStatusCompleted,
			cltest.ParseNullableTime("2000-01-01T00:00:00.000Z"),
		},
		{
			"OldPendingConfirmation",
			models.RunStatusPendingConfirmations,
			cltest.ParseNullableTime("2000-01-01T00:00:00.000Z"),
		},
		{
			"OldErrored",
			models.RunStatusErrored,
			cltest.ParseNullableTime("2000-01-01T00:00:00.000Z"),
		},
		{
			"NewCompleted",
			models.RunStatusCompleted,
			cltest.ParseNullableTime("2018-06-01T00:00:00.000Z"),
		},
		{
			"NewErrored",
			models.RunStatusErrored,
			cltest.ParseNullableTime("2018-06-01T00:00:00.000Z"),
		},
		{
			"NewPendingBridge",
			models.RunStatusPendingBridge,
			cltest.ParseNullableTime("2018-06-01T00:00:00.000Z"),
		},
	}

	app, cleanup := cltest.NewApplication()
	defer cleanup()

	for i, entry := range entries {
		job, initr := cltest.NewJobWithWebInitiator()
		run := job.NewRun(initr)
		run.Status = entry.status
		run.CompletedAt = entry.completedAt
		run.CreatedAt = time.Now().AddDate(-1*i, 0, 0)
		assert.NoError(t, app.Store.Save(&run))
	}

	client, r := cltest.NewClientAndRenderer(app.Store.Config)

	set := flag.NewFlagSet("days", 0)
	set.Int("days", 30, "")
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, client.Prune(c))
	runs := *r.Renders[0].(*[]models.JobRun)
	assert.Equal(t, 3, len(runs))

	client, r = cltest.NewClientAndRenderer(app.Store.Config)

	set = flag.NewFlagSet("prune", 0)
	set.Bool("completed", true, "")
	c = cli.NewContext(nil, set, nil)
	assert.Nil(t, client.Prune(c))
	runs = *r.Renders[0].(*[]models.JobRun)
	assert.Equal(t, 1, len(runs))
}

func first(a models.JobSpec, b interface{}) models.JobSpec {
	return a
}
