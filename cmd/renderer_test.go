package cmd_test

import (
	"bytes"
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

type testWriter struct {
	expected string
	t        testing.TB
	found    bool
}

func (w *testWriter) Write(actual []byte) (int, error) {
	if bytes.Index(actual, []byte(w.expected)) != -1 {
		w.found = true
	}
	return len(actual), nil
}

func TestRendererTableRenderBridge(t *testing.T) {
	bridge := models.BridgeType{Name: "hapax",
		URL:                  cltest.WebURL("http://hap.ax"),
		DefaultConfirmations: 0}
	tw := &testWriter{bridge.Name, t, false}
	r := cmd.RendererTable{Writer: tw}
	assert.Nil(t, r.Render(&bridge))
	assert.Equal(t, tw.found, true)
}

func TestRendererTableRenderUnknown(t *testing.T) {
	r := cmd.RendererTable{Writer: ioutil.Discard}
	anon := struct{ Name string }{"Romeo"}
	assert.NotNil(t, r.Render(&anon))
}
