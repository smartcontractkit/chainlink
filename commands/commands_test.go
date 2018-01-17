package commands_test

import (
	"flag"
	"testing"

	"github.com/h2non/gock"
	"github.com/smartcontractkit/chainlink/commands"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCommandShowJob(t *testing.T) {
	defer cltest.CloseGock(t)
	job := cltest.NewJob()
	gock.New("http://localhost:8080").
		Get("/v2/jobs/" + job.ID).
		Reply(200).
		JSON(job)

	client := commands.Client{commands.RendererNoOp{}}

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{job.ID})
	c := cli.NewContext(nil, set, nil)
	assert.Nil(t, client.ShowJob(c))
}

func TestCommandShowJobNotFound(t *testing.T) {
	defer cltest.CloseGock(t)
	gock.New("http://localhost:8080").
		Get("/v2/jobs/bogus-ID").
		Reply(404)

	client := commands.Client{commands.RendererNoOp{}}

	set := flag.NewFlagSet("test", 0)
	set.Parse([]string{"bogus-ID"})
	c := cli.NewContext(nil, set, nil)
	assert.NotNil(t, client.ShowJob(c))
}
