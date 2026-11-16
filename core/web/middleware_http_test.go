package web

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestSizeLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(requestSizeLimiter(5))
	engine.POST("/", func(c *gin.Context) {
		_, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return
		}
		c.String(http.StatusOK, "ok")
	})

	request, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader("hello!"))
	require.NoError(t, err)
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusRequestEntityTooLarge, recorder.Code)

	request, err = http.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader("hi"))
	require.NoError(t, err)
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestSecureMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(secureMiddleware(true, "secure.example.com", false))
	engine.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/path?q=1", nil)
	require.NoError(t, err)
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusMovedPermanently, recorder.Code)
	assert.Equal(t, "https://secure.example.com/path?q=1", recorder.Header().Get("Location"))

	engine = gin.New()
	engine.Use(secureMiddleware(true, "", true))
	engine.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	request, err = http.NewRequestWithContext(t.Context(), http.MethodGet, "http://example.com/", nil)
	require.NoError(t, err)
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "DENY", recorder.Header().Get("X-Frame-Options"))
}
