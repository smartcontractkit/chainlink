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
	r := cmd.RendererJSON{ioutil.Discard}
	job := cltest.NewJob()
	jobs := []models.JobSpec{job}
	assert.Nil(t, r.Render(&jobs))
}

func TestRendererTableRenderJobs(t *testing.T) {
	r := cmd.RendererTable{ioutil.Discard}
	job := cltest.NewJob()
	jobs := []models.JobSpec{job}
	assert.Nil(t, r.Render(&jobs))
}

func TestRendererTableRenderShowJob(t *testing.T) {
	r := cmd.RendererTable{ioutil.Discard}
	job := cltest.NewJobWithWebInitiator()
	run := job.NewRun()
	p := presenters.JobSpec{job, []models.JobRun{run}}
	assert.Nil(t, r.Render(&p))
}

func TestRendererTableRenderUnknown(t *testing.T) {
	r := cmd.RendererTable{ioutil.Discard}
	anon := struct{ Name string }{"Romeo"}
	assert.NotNil(t, r.Render(&anon))
}
