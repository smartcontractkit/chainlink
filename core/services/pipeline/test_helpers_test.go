package pipeline_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"

	"github.com/smartcontractkit/sqlx"
)

func fakeExternalAdapter(t *testing.T, expectedRequest, response interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Helper()

		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		expectedBody := &bytes.Buffer{}
		err = json.NewEncoder(expectedBody).Encode(expectedRequest)
		require.NoError(t, err)
		require.Equal(t, string(bytes.TrimSpace(expectedBody.Bytes())), string(body))

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		require.NoError(t, err)
	})
}

func makeBridge(t *testing.T, db *sqlx.DB, expectedRequest, response interface{}, cfg pg.QConfig) (*httptest.Server, bridges.BridgeType) {
	t.Helper()

	server := httptest.NewServer(fakeExternalAdapter(t, expectedRequest, response))

	bridgeFeedURL, err := url.ParseRequestURI(server.URL)
	require.NoError(t, err)

	_, bt := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: bridgeFeedURL.String()}, cfg)

	return server, *bt
}

func mustNewObjectParam(t *testing.T, val interface{}) *pipeline.ObjectParam {
	var value pipeline.ObjectParam
	if err := value.UnmarshalPipelineParam(val); err != nil {
		t.Fatalf("failed to init ObjectParam from %v, err: %v", val, err)
	}
	return &value
}
