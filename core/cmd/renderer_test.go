package cmd_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"regexp"
	"testing"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/smartcontractkit/chainlink/core/web"
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
	job := cltest.NewJobWithWebInitiator()
	p := presenters.JobSpec{JobSpec: job}
	assert.Nil(t, r.Render(&p))
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
			run := job.NewRun(job.Initiators[0])
			assert.Nil(t, test.renderer.Render(&presenters.JobRun{run}))
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
	if bytes.Index(actual, []byte(w.expected)) != -1 {
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

			assert.Nil(t, r.Render(bridge))
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

			assert.Nil(t, r.Render(bridge))
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

			assert.Nil(t, r.Render(&[]models.BridgeType{*bridge}))
			assert.Equal(t, test.wantFound, tw.found)
		})
	}
}

func TestRendererTable_Render_TxAttempts(t *testing.T) {
	t.Parallel()

	attempts := []models.TxAttempt{
		models.TxAttempt{
			Hash:      cltest.NewHash(),
			TxID:      1,
			GasPrice:  models.NewBig(big.NewInt(1)),
			Confirmed: false,
			SentAt:    1,
		},
	}

	buffer := bytes.NewBufferString("")
	r := cmd.RendererTable{Writer: buffer}
	assert.NoError(t, r.Render(&attempts))
	output := buffer.String()

	assert.Contains(t, output, fmt.Sprint(attempts[0].TxID))
	assert.Contains(t, output, attempts[0].Hash.Hex())
	assert.Contains(t, output, fmt.Sprint(attempts[0].GasPrice))
	assert.Contains(t, output, fmt.Sprint(attempts[0].SentAt))
	assert.Contains(t, output, fmt.Sprint(attempts[0].Confirmed))
}

func TestRendererTable_Render_Tx(t *testing.T) {
	t.Parallel()

	from := cltest.NewAddress()
	to := cltest.NewAddress()
	tx := presenters.Tx{
		Hash:      cltest.NewHash(),
		Nonce:     "1",
		From:      &from,
		To:        &to,
		GasPrice:  "2",
		Confirmed: false,
		SentAt:    "3",
	}

	buffer := bytes.NewBufferString("")
	r := cmd.RendererTable{Writer: buffer}
	assert.NoError(t, r.Render(&tx))
	output := buffer.String()

	assert.Contains(t, output, tx.Hash.Hex())
	assert.Contains(t, output, tx.Nonce)
	assert.Contains(t, output, from.Hex())
	assert.Contains(t, output, to.Hex())
	assert.Contains(t, output, fmt.Sprint(tx.Confirmed))
}

func TestRendererTable_Render_Txs(t *testing.T) {
	t.Parallel()

	a := cltest.NewAddress()
	txs := []presenters.Tx{
		{
			Hash:      cltest.NewHash(),
			Nonce:     "1",
			From:      &a,
			GasPrice:  "2",
			Confirmed: false,
			SentAt:    "3",
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
	assert.Contains(t, output, fmt.Sprint(txs[0].Confirmed))
}

func TestRendererTable_ServiceAgreementShow(t *testing.T) {
	t.Parallel()

	sa, err := cltest.ServiceAgreementFromString(string(cltest.MustReadFile(t, "testdata/hello_world_agreement.json")))
	assert.NoError(t, err)
	psa := presenters.ServiceAgreement{ServiceAgreement: sa}

	buffer := bytes.NewBufferString("")
	r := cmd.RendererTable{Writer: buffer}

	assert.NoError(t, r.Render(&psa))
	output := buffer.String()
	assert.Regexp(t, regexp.MustCompile("0x[0-9a-zA-Z]{64}"), output)
	assert.Regexp(t, regexp.MustCompile("1.000000000000000000 LINK"), output)
	assert.Regexp(t, regexp.MustCompile("300 seconds"), output)
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
