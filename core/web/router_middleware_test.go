package web

import (
	"expvar"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelmetDefault(t *testing.T) {
	engine := gin.New()
	engine.Use(helmetDefault())
	engine.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	recorder := httptest.NewRecorder()
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	require.NoError(t, err)
	engine.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "ok", recorder.Body.String())
	for key, expected := range map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-DNS-Prefetch-Control":    "off",
		"X-Frame-Options":           "DENY",
		"Strict-Transport-Security": "max-age=5184000; includeSubDomains",
		"X-Download-Options":        "noopen",
		"X-XSS-Protection":          "1; mode=block",
	} {
		assert.Equal(t, expected, recorder.Header().Get(key), key)
	}
}

func TestExpvarHandlerFormatting(t *testing.T) {
	number := new(expvar.Int)
	number.Set(42)
	quoted := new(expvar.String)
	quoted.Set("a\"b")

	engine := gin.New()
	engine.GET("/debug/vars", expvarHandler(func(yield func(expvar.KeyValue)) {
		yield(expvar.KeyValue{Key: "number", Value: number})
		yield(expvar.KeyValue{Key: "quoted", Value: quoted})
	}))

	recorder := httptest.NewRecorder()
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/debug/vars", nil)
	require.NoError(t, err)
	engine.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("Content-Type"))
	assert.Equal(t, "{\n\"number\": 42,\n\"quoted\": \"a\\\"b\"\n}\n", recorder.Body.String())
}

func TestExpvarHandlerEmpty(t *testing.T) {
	engine := gin.New()
	engine.GET("/debug/vars", expvarHandler(func(func(expvar.KeyValue)) {}))

	recorder := httptest.NewRecorder()
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/debug/vars", nil)
	require.NoError(t, err)
	engine.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "{\n\n}\n", recorder.Body.String())
}
