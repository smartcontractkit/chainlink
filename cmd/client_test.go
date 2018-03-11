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
	set := flag.NewFlagSet("create", 0)
	set.Parse([]string{"badinput"})
	c := cli.NewContext(nil, set, nil)
	assert.NotNil(t, client.CreateJobSpec(c))
}
