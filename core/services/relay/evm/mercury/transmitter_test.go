package mercury

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func Test_MercuryTransmitter_MercuryReport_Marshal_Unmarshal(t *testing.T) {
	mr := MercuryReport{samplePayload, sampleFromAccount}

	b, err := json.Marshal(mr)
	require.NoError(t, err)

	assert.JSONEq(t, sampleMercuryReport, string(b))

	nmr := MercuryReport{}
	err = json.Unmarshal(b, &nmr)
	require.NoError(t, err)

	assert.Equal(t, mr, nmr)
}

type MockHTTPClient struct {
	do func(req *http.Request) (*http.Response, error)
}

func (m MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.do(req)
}

func Test_MercuryTransmitter_Transmit(t *testing.T) {
	lggr := logger.TestLogger(t)
	reportURL := "http://report.test/foo"
	username := "my username"
	password := "my password"

	t.Run("successful HTTP POST", func(t *testing.T) {
		httpClient := MockHTTPClient{
			do: func(req *http.Request) (resp *http.Response, err error) {
				assert.Equal(t, "POST", req.Method)
				buf, err := req.GetBody()
				require.NoError(t, err)
				d := json.NewDecoder(buf)
				mr := MercuryReport{}
				err = d.Decode(&mr)
				require.NoError(t, err, "expected JSON to unmarshal into MercuryReport{}")
				assert.Equal(t, d.InputOffset(), req.ContentLength)
				assert.Equal(t, samplePayloadHex, mr.Payload.String())
				assert.Equal(t, sampleFromAccount, mr.FromAccount)
				assert.Equal(t, "report.test", req.Host)
				assert.Equal(t, "/foo", req.URL.Path)
				resp = new(http.Response)
				resp.Body = io.NopCloser(bytes.NewBuffer([]byte{}))
				resp.Status = "200 OK"
				resp.StatusCode = 200
				return resp, nil
			},
		}
		mt := NewTransmitter(lggr, httpClient, sampleFromAccount, reportURL, username, password)
		err := mt.Transmit(testutils.Context(t), sampleReportContext, sampleReport, sampleSigs)

		require.NoError(t, err)
	})

	t.Run("failing HTTP POST", func(t *testing.T) {
		httpClient := MockHTTPClient{
			do: func(req *http.Request) (resp *http.Response, err error) {
				return nil, errors.New("foo error")
			},
		}
		mt := NewTransmitter(lggr, httpClient, sampleFromAccount, reportURL, username, password)
		err := mt.Transmit(testutils.Context(t), sampleReportContext, sampleReport, sampleSigs)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "foo error")
	})
}
