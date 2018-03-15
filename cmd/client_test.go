package cmd_test

import (
	"flag"
	"testing"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestClient_RunNode(t *testing.T) {
	app, _ := cltest.NewApplicationWithKeyStore() // cleanup invoked in client.RunNode
	r := &cltest.RendererMock{}
	var called bool
	auth := cltest.CallbackAuthenticator{func(*store.Store, string) { called = true }}
	client := cmd.Client{
		r,
		app.Store.Config,
		cltest.InstanceAppFactory{App: app},
		auth,
		cltest.EmptyRunner{}}

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{""})
	c := cli.NewContext(nil, set, nil)

	client.RunNode(c)
	assert.True(t, called)
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

func TestClient_ShowJobSpec(t *testing.T) {
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
	assert.NotNil(t, client.ShowJobSpec(c))
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
		{"{bad json}", 0, true},
		{"bad/filepath/", 0, true},
		{"{\"initiators\":[{\"type\":\"web\"}], \"tasks\":[{\"type\": \"NoOp\"}]}", 1, false},
		{"{\"initiators\":[{\"type\":\"ethLog\", \"address\": \"0x3cCad4715152693fE3BC4460591e3D3Fbd071b42\"}],\"tasks\": [ { \"type\": \"NoOp\" } ]}", 2, false},
		{"../internal/fixtures/web/eth_log_job.json", 3, false},
		{"~/go/src/github.com/smartcontractkit/chainlink/internal/fixtures/web/hello_world_job.json", 4, false},
		{"~/go/src/github.com/smartcontractkit/chainlink/internal/fixtures/web/invalid_cron.json", 4, true},
	}

	for _, tt := range tests {
		test := tt

		set := flag.NewFlagSet("create", 0)
		set.Parse([]string{test.input})
		c := cli.NewContext(nil, set, nil)
		if test.errored {
			assert.NotNil(t, client.CreateJobSpec(c))
		} else {
			assert.Nil(t, client.CreateJobSpec(c))
		}
		numberOfJobs, _ := app.Store.Jobs()
		assert.Equal(t, len(numberOfJobs), test.nJobs)
	}
}
