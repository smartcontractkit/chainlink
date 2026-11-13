package web_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	appmocks "github.com/smartcontractkit/chainlink/v2/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/jobsync"
	"github.com/smartcontractkit/chainlink/v2/core/web"
)

func TestReloadController_Reload(t *testing.T) {
	t.Parallel()
	raw := []byte(`{"revision":2,"jobs":[]}`)
	mockApp := appmocks.NewApplication(t)
	mockApp.EXPECT().ReloadDynamicConfig(mock.Anything).Return(jobsync.Applied{
		Result: jobsync.Result{Revision: 2},
		Raw:    raw,
	}, nil)

	controller := web.ReloadController{App: mockApp}
	w := httptest.NewRecorder()
	engine := gin.New()
	engine.POST("/reload", controller.Reload)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/reload", nil)
	require.NoError(t, err)
	req.RemoteAddr = "127.0.0.1:54321"

	engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var got web.ReloadResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, uint64(2), got.Revision)
	// sha256 over the canonical form of the config, which is what the sidecar
	// compares against the hash it verified.
	assert.Equal(t, "sha256:f45d6516ba178b894883313f2a184a5ab458a1b7a2be9e957bd169c97cdfab72", got.ConfigHash)
}

func TestReloadController_RejectsNonLoopback(t *testing.T) {
	t.Parallel()
	mockApp := appmocks.NewApplication(t)

	controller := web.ReloadController{App: mockApp}
	w := httptest.NewRecorder()
	engine := gin.New()
	engine.POST("/reload", controller.Reload)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/reload", nil)
	require.NoError(t, err)
	req.RemoteAddr = "10.0.0.1:54321"

	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
	mockApp.AssertNotCalled(t, "ReloadDynamicConfig", mock.Anything)
}

func TestReloadController_ReportsFailure(t *testing.T) {
	t.Parallel()
	mockApp := appmocks.NewApplication(t)
	mockApp.EXPECT().ReloadDynamicConfig(mock.Anything).Return(jobsync.Applied{}, jobsync.ErrInvalidConfig)

	controller := web.ReloadController{App: mockApp}
	w := httptest.NewRecorder()
	engine := gin.New()
	engine.POST("/reload", controller.Reload)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/reload", nil)
	require.NoError(t, err)
	req.RemoteAddr = "127.0.0.1:54321"

	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
