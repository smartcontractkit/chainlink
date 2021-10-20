package cmd_test

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/web"
	webpresenters "github.com/smartcontractkit/chainlink/core/web/presenters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRendererJSON_RenderVRFKeys(t *testing.T) {
	t.Parallel()

	now := time.Now()
	r := cmd.RendererJSON{Writer: ioutil.Discard}
	keys := []cmd.VRFKeyPresenter{
		{
			VRFKeyResource: webpresenters.VRFKeyResource{
				Compressed:   "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4d01",
				Uncompressed: "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4db44652a69526181101d4aa9a58ecf43b1be972330de99ea5e540f56f4e0a672f",
				Hash:         "0x9926c5f19ec3b3ce005e1c183612f05cfc042966fcdd82ec6e78bf128d91695a",
				CreatedAt:    now,
				UpdatedAt:    now,
				DeletedAt:    nil,
			},
		},
	}
	assert.NoError(t, r.Render(&keys))
}

func TestRendererJSON_RenderJobs(t *testing.T) {
	t.Parallel()
	r := cmd.RendererJSON{Writer: ioutil.Discard}
	job := cltest.NewJob()
	jobs := []models.JobSpec{job}
	assert.NoError(t, r.Render(&jobs))
}

func TestRendererTable_RenderJobs(t *testing.T) {
	t.Parallel()

	buffer := bytes.NewBufferString("")
	r := cmd.RendererTable{Writer: buffer}
	job := cltest.NewJob()
	jobs := []models.JobSpec{job}
	assert.NoError(t, r.Render(&jobs))

	output := buffer.String()
	assert.Contains(t, output, "noop")
}

func TestRendererTable_RenderConfiguration(t *testing.T) {
	t.Parallel()

	ethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		ethClient,
	)
	defer cleanup()
	require.NoError(t, app.Start())
	client := app.NewHTTPClient()

	resp, cleanup := client.Get("/v2/config")
	defer cleanup()
	cp := presenters.ConfigPrinter{}
	require.NoError(t, cltest.ParseJSONAPIResponse(t, resp, &cp))

	r := cmd.RendererTable{Writer: ioutil.Discard}
	assert.NoError(t, r.Render(&cp))
}

func TestRendererTable_RenderShowJob(t *testing.T) {
	t.Parallel()
	r := cmd.RendererTable{Writer: ioutil.Discard}
	job := cltest.NewJobWithWebInitiator()
	p := presenters.JobSpec{JobSpec: job}
	assert.NoError(t, r.Render(&p))
}

func TestRenderer_RenderJobRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		renderer cmd.Renderer
	}{
		{"json", cmd.RendererJSON{Writer: ioutil.Discard}},
		{"table", cmd.RendererTable{Writer: ioutil.Discard}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			job := cltest.NewJobWithWebInitiator()
			run := cltest.NewJobRun(job)
			assert.NoError(t, test.renderer.Render(&presenters.JobRun{JobRun: run}))
		})
	}
}

func TestRendererTable_RenderJobRun(t *testing.T) {
	t.Parallel()
	r := cmd.RendererTable{Writer: ioutil.Discard}
	job := cltest.NewJob()
	jobs := []models.JobSpec{job}
	assert.NoError(t, r.Render(&jobs))
}

type testWriter struct {
	expected string
	t        testing.TB
	found    bool
}

func (w *testWriter) Write(actual []byte) (int, error) {
	if bytes.Contains(actual, []byte(w.expected)) {
		w.found = true
	}
	return len(actual), nil
}

func TestRendererTable_RenderExternalInitiatorAuthentication(t *testing.T) {
	t.Parallel()

	eia := webpresenters.ExternalInitiatorAuthentication{
		Name:           "bitcoin",
		URL:            cltest.WebURL(t, "http://localhost:8888"),
		AccessKey:      "accesskey",
		Secret:         "secret",
		OutgoingToken:  "outgoingToken",
		OutgoingSecret: "outgoingSecret",
	}
	tests := []struct {
		name, content string
	}{
		{"Name", eia.Name},
		{"URL", eia.URL.String()},
		{"AccessKey", eia.AccessKey},
		{"Secret", eia.Secret},
		{"OutgoingToken", eia.OutgoingToken},
		{"OutgoingSecret", eia.OutgoingSecret},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tw := &testWriter{test.content, t, false}
			r := cmd.RendererTable{Writer: tw}

			assert.NoError(t, r.Render(&eia))
			assert.True(t, tw.found)
		})
	}
}

func checkPresence(t *testing.T, s, output string) { assert.Regexp(t, regexp.MustCompile(s), output) }

func TestRendererTable_ServiceAgreementShow(t *testing.T) {
	t.Parallel()

	sa, err := cltest.ServiceAgreementFromString(string(cltest.MustReadFile(t, "../testdata/jsonspecs/hello_world_agreement.json")))
	assert.NoError(t, err)
	psa := presenters.ServiceAgreement{ServiceAgreement: sa}

	buffer := bytes.NewBufferString("")
	r := cmd.RendererTable{Writer: buffer}

	require.NoError(t, r.Render(&psa))
	output := buffer.String()
	checkPresence(t, "0x[0-9a-zA-Z]{64}", output)
	checkPresence(t, "1.000000000000000000 LINK", output)
	checkPresence(t, "300 seconds", output)
	checkPresence(t, "0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF", output) // Aggregator address
	checkPresence(t, "0xd0771e55", output)                                 // AggInitiateJobSelector
	checkPresence(t, "0xbadc0de5", output)                                 // AggFulfillSelector
}

func TestRendererTable_PatchResponse(t *testing.T) {
	t.Parallel()

	buffer := bytes.NewBufferString("")
	r := cmd.RendererTable{Writer: buffer}

	patchResponse := web.ConfigPatchResponse{
		EthGasPriceDefault: web.Change{
			From: "98721",
			To:   "53276",
		},
	}

	assert.NoError(t, r.Render(&patchResponse))
	output := buffer.String()
	assert.Regexp(t, regexp.MustCompile("98721"), output)
	assert.Regexp(t, regexp.MustCompile("53276"), output)
}

func TestRendererTable_RenderUnknown(t *testing.T) {
	t.Parallel()
	r := cmd.RendererTable{Writer: ioutil.Discard}
	anon := struct{ Name string }{"Romeo"}
	assert.Error(t, r.Render(&anon))
}
