package cmd_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/eth"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRendererJSON_RenderVRFKeys(t *testing.T) {
	t.Parallel()

	now := time.Now()
	r := cmd.RendererJSON{Writer: ioutil.Discard}
	keys := []cmd.VRFKeyPresenter{
		{
			Compressed:   "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4d01",
			Uncompressed: "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4db44652a69526181101d4aa9a58ecf43b1be972330de99ea5e540f56f4e0a672f",
			Hash:         "0x9926c5f19ec3b3ce005e1c183612f05cfc042966fcdd82ec6e78bf128d91695a",
			CreatedAt:    &now,
			UpdatedAt:    &now,
			DeletedAt:    nil,
		},
	}
	assert.NoError(t, r.Render(&keys))
}

func TestRendererTable_RenderVRKKeys(t *testing.T) {
	t.Parallel()

	now := time.Now()
	buffer := bytes.NewBufferString("")
	r := cmd.RendererTable{Writer: buffer}
	keys := []cmd.VRFKeyPresenter{
		{
			Compressed:   "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4d01",
			Uncompressed: "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4db44652a69526181101d4aa9a58ecf43b1be972330de99ea5e540f56f4e0a672f",
			Hash:         "0x9926c5f19ec3b3ce005e1c183612f05cfc042966fcdd82ec6e78bf128d91695a",
			CreatedAt:    &now,
			UpdatedAt:    &now,
			DeletedAt:    nil,
		},
	}
	assert.NoError(t, r.Render(&keys))
	output := buffer.String()
	assert.Contains(t, output, "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4d01")
	assert.Contains(t, output, "0xe2c659dd73ded1663c0caf02304aac5ccd247047b3993d273a8920bba0402f4db44652a69526181101d4aa9a58ecf43b1be972330de99ea5e540f56f4e0a672f")
	assert.Contains(t, output, "0x9926c5f19ec3b3ce005e1c183612f05cfc042966fcdd82ec6e78bf128d91695a")

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

func TestRendererTable_RenderJobsV2(t *testing.T) {
	t.Parallel()
	now := time.Now()

	buffer := bytes.NewBufferString("")
	r := cmd.RendererTable{Writer: buffer}
	jobs := []cmd.Job{
		{
			JAID: cmd.JAID{ID: "1"},
			Name: "Test Job",
			Type: cmd.DirectRequestJob,
			DirectRequestSpec: &cmd.DirectRequestSpec{
				CreatedAt: now,
			},
			PipelineSpec: cmd.PipelineSpec{
				DotDAGSource: "    ds1          [type=http method=GET url=\"example.com\" allowunrestrictednetworkaccess=\"true\"];\n    ds1_parse    [type=jsonparse path=\"USD\"];\n    ds1_multiply [type=multiply times=100];\n    ds1 -\u003e ds1_parse -\u003e ds1_multiply;\n",
			},
		},
	}
	assert.NoError(t, r.Render(&jobs))

	output := buffer.String()
	assert.Contains(t, output, "1")
	assert.Contains(t, output, "Test Job")
	assert.Contains(t, output, now.Format(time.RFC3339))
	assert.Contains(t, output, "directrequest")
	assert.Contains(t, output, "ds1 http")
	assert.Contains(t, output, "ds1_parse jsonparse")
	assert.Contains(t, output, "ds1_multiply multiply")
}

func TestRendererTable_RenderConfiguration(t *testing.T) {
	t.Parallel()

	rpcClient, gethClient, _, assertMocksCalled := cltest.NewEthMocksWithStartupAssertions(t)
	defer assertMocksCalled()
	app, cleanup := cltest.NewApplicationWithKey(t,
		eth.NewClientWith(rpcClient, gethClient),
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

func TestRendererTable_RenderBridgeShow(t *testing.T) {
	t.Parallel()
	_, bridge := cltest.NewBridgeType(t, "hapax", "http://hap.ax")
	bridge.Confirmations = 0

	tests := []struct {
		name, content string
	}{
		{"name", bridge.Name.String()},
		{"outgoing token", bridge.OutgoingToken},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tw := &testWriter{test.content, t, false}
			r := cmd.RendererTable{Writer: tw}

			assert.NoError(t, r.Render(bridge))
			assert.True(t, tw.found)
		})
	}
}

func TestRendererTable_RenderBridgeAdd(t *testing.T) {
	t.Parallel()
	bridge, _ := cltest.NewBridgeType(t, "hapax", "http://hap.ax")
	bridge.Confirmations = 0

	tests := []struct {
		name, content string
	}{
		{"name", bridge.Name.String()},
		{"outgoing token", bridge.OutgoingToken},
		{"incoming token", bridge.IncomingToken},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tw := &testWriter{test.content, t, false}
			r := cmd.RendererTable{Writer: tw}

			assert.NoError(t, r.Render(bridge))
			assert.True(t, tw.found)
		})
	}
}

func TestRendererTable_RenderBridgeList(t *testing.T) {
	t.Parallel()
	_, bridge := cltest.NewBridgeType(t, "hapax", "http://hap.ax")
	bridge.Confirmations = 0

	tests := []struct {
		name, content string
		wantFound     bool
	}{
		{"name", bridge.Name.String(), true},
		{"outgoing token", bridge.OutgoingToken, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tw := &testWriter{test.content, t, false}
			r := cmd.RendererTable{Writer: tw}

			assert.NoError(t, r.Render(&[]models.BridgeType{*bridge}))
			assert.Equal(t, test.wantFound, tw.found)
		})
	}
}

func TestRendererTable_RenderExternalInitiatorAuthentication(t *testing.T) {
	t.Parallel()

	eia := presenters.ExternalInitiatorAuthentication{
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

func TestRendererTable_Render_Tx(t *testing.T) {
	t.Parallel()

	from := cltest.NewAddress()
	to := cltest.NewAddress()
	tx := presenters.EthTx{
		Hash:     cltest.NewHash(),
		Nonce:    "1",
		From:     &from,
		To:       &to,
		GasPrice: "2",
		State:    "confirmed",
		SentAt:   "3",
	}

	buffer := bytes.NewBufferString("")
	r := cmd.RendererTable{Writer: buffer}
	assert.NoError(t, r.Render(&tx))
	output := buffer.String()

	assert.NotContains(t, output, tx.Hash.Hex())
	assert.Contains(t, output, tx.Nonce)
	assert.Contains(t, output, from.Hex())
	assert.Contains(t, output, to.Hex())
	assert.Contains(t, output, fmt.Sprint(tx.State))
}

func TestRendererTable_Render_Txs(t *testing.T) {
	t.Parallel()

	a := cltest.NewAddress()
	txs := []presenters.EthTx{
		{
			Hash:     cltest.NewHash(),
			Nonce:    "1",
			From:     &a,
			GasPrice: "2",
			State:    "confirmed",
			SentAt:   "3",
		},
	}

	buffer := bytes.NewBufferString("")
	r := cmd.RendererTable{Writer: buffer}
	assert.NoError(t, r.Render(&txs))
	output := buffer.String()

	assert.Contains(t, output, txs[0].Nonce)
	assert.Contains(t, output, txs[0].Hash.Hex())
	assert.Contains(t, output, txs[0].GasPrice)
	assert.Contains(t, output, txs[0].SentAt)
	assert.Contains(t, output, a.Hex())
	assert.Contains(t, output, fmt.Sprint(txs[0].State))
}

func checkPresence(t *testing.T, s, output string) { assert.Regexp(t, regexp.MustCompile(s), output) }

func TestRendererTable_ServiceAgreementShow(t *testing.T) {
	t.Parallel()

	sa, err := cltest.ServiceAgreementFromString(string(cltest.MustReadFile(t, "testdata/hello_world_agreement.json")))
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
