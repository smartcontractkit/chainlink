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

func TestRunNode(t *testing.T) {
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

func TestClientGetJobs(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j1 := cltest.NewJob()
	app.Store.SaveJob(&j1)
	j2 := cltest.NewJob()
	app.Store.SaveJob(&j2)

	client, r := cltest.NewClientAndRenderer(app.Store.Config)

	assert.Nil(t, client.GetJobs(nil))
	jobs := *r.Renders[0].(*[]models.Job)
	assert.Equal(t, 2, len(jobs))
	assert.Equal(t, j1.ID, jobs[0].ID)
}

func TestClientShowJob(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	job := cltest.NewJob()
	app.Store.SaveJob(&job)

	client, r := cltest.NewClientAndRenderer(app.Store.Config)

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{job.ID})
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, client.ShowJob(c))
	assert.Equal(t, 1, len(r.Renders))
	assert.Equal(t, job.ID, r.Renders[0].(*presenters.Job).ID)
}

func TestClientShowJobNotFound(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	client, r := cltest.NewClientAndRenderer(app.Store.Config)

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"bogus-ID"})
	c := cli.NewContext(nil, set, nil)
	assert.NotNil(t, client.ShowJob(c))
	assert.Empty(t, r.Renders)
}
