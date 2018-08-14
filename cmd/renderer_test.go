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

func TestRendererJSON_RenderJobs(t *testing.T) {
	t.Parallel()
	r := cmd.RendererJSON{Writer: ioutil.Discard}
	job := cltest.NewJob()
	jobs := []models.JobSpec{job}
	assert.Nil(t, r.Render(&jobs))
}

func TestRendererTable_RenderJobs(t *testing.T) {
	t.Parallel()
	r := cmd.RendererTable{Writer: ioutil.Discard}
	job := cltest.NewJob()
	jobs := []models.JobSpec{job}
	assert.NoError(t, r.Render(&jobs))
}

func TestRendererTable_RenderShowJob(t *testing.T) {
	t.Parallel()
	r := cmd.RendererTable{Writer: ioutil.Discard}
	job, initr := cltest.NewJobWithWebInitiator()
	run := job.NewRun(initr)
	p := presenters.JobSpec{JobSpec: job, Runs: []presenters.JobRun{presenters.JobRun{run}}}
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

func TestRendererTable_RenderBridgeShow(t *testing.T) {
	t.Parallel()
	bridge := cltest.NewBridgeType("hapax", "http://hap.ax")
	bridge.DefaultConfirmations = 0

	tests := []struct {
		name, content string
	}{
		{"name", bridge.Name.String()},
		{"incoming token", bridge.IncomingToken},
		{"outgoing token", bridge.OutgoingToken},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tw := &testWriter{test.content, t, false}
			r := cmd.RendererTable{Writer: tw}

			assert.Nil(t, r.Render(&bridge))
			assert.True(t, tw.found)
		})
	}
}

func TestRendererTable_RenderBridgeList(t *testing.T) {
	t.Parallel()
	bridge := cltest.NewBridgeType("hapax", "http://hap.ax")
	bridge.DefaultConfirmations = 0

	tests := []struct {
		name, content string
		wantFound     bool
	}{
		{"name", bridge.Name.String(), true},
		{"incoming token", bridge.IncomingToken, false},
		{"outgoing token", bridge.OutgoingToken, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tw := &testWriter{test.content, t, false}
			r := cmd.RendererTable{Writer: tw}

			assert.Nil(t, r.Render(&[]models.BridgeType{bridge}))
			assert.Equal(t, test.wantFound, tw.found)
		})
	}
}

func TestRendererTable_ServiceAgreementShow(t *testing.T) {
	t.Parallel()

	sa := cltest.ServiceAgreementFromString(string(cltest.LoadJSON("../internal/fixtures/web/hello_world_agreement.json")))
	psa := presenters.ServiceAgreement{ServiceAgreement: sa}

	tests := []struct {
		name, content string
	}{
		{"ID", "0x8de92e4d6a2b527f1e3c39022ee0f1c177a6d874224fe54f5c4c0d5bcaa57d50"},
		{"payment amount", "1.000000000000000000 LINK"},
		{"expiration", "300 seconds"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tw := &testWriter{test.content, t, false}
			r := cmd.RendererTable{Writer: tw}

			assert.NoError(t, r.Render(&psa))
			assert.True(t, tw.found)
		})
	}
}

func TestRendererTable_RenderUnknown(t *testing.T) {
	t.Parallel()
	r := cmd.RendererTable{Writer: ioutil.Discard}
	anon := struct{ Name string }{"Romeo"}
	assert.Error(t, r.Render(&anon))
}
