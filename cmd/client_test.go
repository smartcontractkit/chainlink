package cmd_test

import (
	"flag"
	"testing"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestClientShowJob(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	job := cltest.NewJob()
	app.Store.SaveJob(job)

	client := cmd.Client{cmd.RendererNoOp{}, app.Store.Config}

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{job.ID})
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, client.ShowJob(c))
}

func TestClientShowJobNotFound(t *testing.T) {
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := cmd.Client{cmd.RendererNoOp{}, app.Store.Config}

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"bogus-ID"})
	c := cli.NewContext(nil, set, nil)
	assert.NotNil(t, client.ShowJob(c))
}
