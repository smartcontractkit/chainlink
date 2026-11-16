package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func corsEngine(allowOrigins string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(uiCorsHandler(allowOrigins))
	engine.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	engine.POST("/", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	return engine
}

func serveCORS(t *testing.T, engine *gin.Engine, method, origin string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	request, err := http.NewRequestWithContext(t.Context(), method, "http://example.com/", nil)
	require.NoError(t, err)
	if origin != "" {
		request.Header.Set("Origin", origin)
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	return recorder
}

func TestLocalCORS_NoOrigin(t *testing.T) {
	recorder := serveCORS(t, corsEngine("http://localhost:3000"), http.MethodGet, "", nil)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Empty(t, recorder.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, recorder.Header().Get("Vary"))
}

func TestLocalCORS_AllowedOrigin(t *testing.T) {
	recorder := serveCORS(t, corsEngine("http://localhost:3000, https://example.com"), http.MethodGet, "http://localhost:3000", nil)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "http://localhost:3000", recorder.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", recorder.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "Content-Length", recorder.Header().Get("Access-Control-Expose-Headers"))
	assert.Equal(t, "Origin", recorder.Header().Get("Vary"))
}

func TestLocalCORS_DisallowedOrigin(t *testing.T) {
	recorder := serveCORS(t, corsEngine("http://localhost:3000"), http.MethodGet, "http://localhost:1234", nil)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
	assert.Empty(t, recorder.Header().Get("Access-Control-Allow-Origin"))
}

func TestLocalCORS_AllowAllOrigins(t *testing.T) {
	recorder := serveCORS(t, corsEngine("*"), http.MethodGet, "http://localhost:1234", nil)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "*", recorder.Header().Get("Access-Control-Allow-Origin"))
	assert.Empty(t, recorder.Header().Get("Vary"))
}

func TestLocalCORS_SameOrigin(t *testing.T) {
	recorder := serveCORS(t, corsEngine("http://localhost:3000"), http.MethodGet, "http://example.com", nil)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Empty(t, recorder.Header().Get("Access-Control-Allow-Origin"))
}

func TestLocalCORS_Preflight(t *testing.T) {
	recorder := serveCORS(t, corsEngine("http://localhost:3000"), http.MethodOptions, "http://localhost:3000", map[string]string{
		"Access-Control-Request-Method":  "POST",
		"Access-Control-Request-Headers": "Content-Type",
	})
	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.Equal(t, "http://localhost:3000", recorder.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", recorder.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "GET,POST,PATCH,DELETE", recorder.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Origin,Content-Type,Accept", recorder.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "2147483647", recorder.Header().Get("Access-Control-Max-Age"))
	assert.ElementsMatch(t, []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"}, recorder.Header().Values("Vary"))
}

func TestLocalCORS_PreflightDisallowedOrigin(t *testing.T) {
	recorder := serveCORS(t, corsEngine("http://localhost:3000"), http.MethodOptions, "http://localhost:1234", nil)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}
