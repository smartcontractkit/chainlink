package cmd_test

import (
	"io/ioutil"
	"testing"

	"github.com/smartcontractkit/chainlink/cmd"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/presenters"
	"github.com/stretchr/testify/assert"
)

func TestRendererJSONRenderJobs(t *testing.T) {
	r := cmd.RendererJSON{Writer: ioutil.Discard}
	job := cltest.NewJob()
	jobs := []models.JobSpec{job}
	assert.Nil(t, r.Render(&jobs))
}

func TestRendererTableRenderJobs(t *testing.T) {
	r := cmd.RendererTable{Writer: ioutil.Discard}
	job := cltest.NewJob()
	jobs := []models.JobSpec{job}
	assert.Nil(t, r.Render(&jobs))
}

func TestRendererTableRenderShowJob(t *testing.T) {
	r := cmd.RendererTable{Writer: ioutil.Discard}
	job, initr := cltest.NewJobWithWebInitiator()
	run := job.NewRun(initr)
	p := presenters.JobSpec{JobSpec: job, Runs: []models.JobRun{run}}
	assert.Nil(t, r.Render(&p))
}

func TestRendererTableRenderBridge(t *testing.T) {
	r := cmd.RendererTable{Writer: ioutil.Discard}
	bridge := models.BridgeType{}
	assert.Nil(t, r.Render(&bridge))
}

func TestRendererTableRenderUnknown(t *testing.T) {
	r := cmd.RendererTable{Writer: ioutil.Discard}
	anon := struct{ Name string }{"Romeo"}
	assert.NotNil(t, r.Render(&anon))
}
