package web_test

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	appmocks "github.com/smartcontractkit/chainlink/v2/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web"
)

func TestLPSkipController_LPSkipToBlock(t *testing.T) {
	t.Parallel()
	cfg := configtest.NewGeneralConfig(t, func(config *chainlink.Config, secrets *chainlink.Secrets) {
		config.Feature.LogPoller = new(true)
	})
	ec := setupEthClientForControllerTests(t)
	app := cltest.NewApplicationWithConfigAndKey(t, cfg, cltest.DefaultP2PKey, ec)
	require.NoError(t, app.Start(t.Context()))
	client := app.NewHTTPClient(nil)

	postSkip := func(t *testing.T, request web.LPSkipToBlockRequest) *http.Response {
		t.Helper()
		body, err := json.Marshal(request)
		require.NoError(t, err)
		resp, cleanup := client.Post("/v2/lp_skip_to_block", bytes.NewReader(body))
		t.Cleanup(cleanup)
		return resp
	}

	t.Run("missing chain family", func(t *testing.T) {
		t.Parallel()
		resp := postSkip(t, web.LPSkipToBlockRequest{
			BlockNumber: 100,
			ChainID:     "1",
		})
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(b), "chain family was not provided")
	})

	t.Run("unsupported chain family", func(t *testing.T) {
		t.Parallel()
		resp := postSkip(t, web.LPSkipToBlockRequest{
			BlockNumber: 100,
			Family:      "solana",
			ChainID:     "1",
		})
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(b), "only evm is supported")
	})

	t.Run("missing chain id", func(t *testing.T) {
		t.Parallel()
		resp := postSkip(t, web.LPSkipToBlockRequest{
			BlockNumber: 100,
			Family:      "evm",
		})
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(b), "chain-id was not provided")
	})

	t.Run("invalid block number", func(t *testing.T) {
		t.Parallel()
		resp := postSkip(t, web.LPSkipToBlockRequest{
			BlockNumber: 1,
			Family:      "evm",
			ChainID:     "1",
		})
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(b), "block number must be")
	})

	t.Run("unknown relayer", func(t *testing.T) {
		t.Parallel()
		resp := postSkip(t, web.LPSkipToBlockRequest{
			BlockNumber: 100,
			Family:      "evm",
			ChainID:     "99999",
		})
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(b), "relayer does not exist")
	})
}

func TestLPSkipController_LPSkipToBlock_HappyPath(t *testing.T) {
	t.Parallel()

	mockApp := appmocks.NewApplication(t)
	mockApp.EXPECT().LPSkipToBlock(mock.Anything, "evm", "1", int64(100)).Return(nil)

	controller := web.LPSkipController{App: mockApp}

	w := httptest.NewRecorder()
	engine := gin.New()
	engine.POST("/v2/lp_skip_to_block", controller.LPSkipToBlock)

	body, err := json.Marshal(web.LPSkipToBlockRequest{
		BlockNumber: 100,
		Family:      "evm",
		ChainID:     "1",
	})
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(t.Context(), "POST", "/v2/lp_skip_to_block", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response web.LPSkipToBlockResponse
	err = web.ParseJSONAPIResponse(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Log poller will start processing from the new block on next tick", response.Message)
	assert.Equal(t, "1", response.ChainID)
	assert.Equal(t, int64(100), response.BlockNumber)
}
