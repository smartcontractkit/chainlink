package ginprom

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrometheusInstrument(t *testing.T) {
	gin.SetMode(gin.TestMode)
	registry := prometheus.NewRegistry()
	p := New(Namespace("service"), Registry(registry))

	engine := gin.New()
	engine.Use(p.Instrument())
	engine.GET("/hello", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/hello", nil)
	require.NoError(t, err)
	engine.ServeHTTP(httptest.NewRecorder(), request)

	families, err := registry.Gather()
	require.NoError(t, err)
	var counter, duration, requestSize, responseSize bool
	for _, family := range families {
		switch family.GetName() {
		case "service_gonic_requests_total":
			counter = true
			require.Len(t, family.GetMetric(), 1)
			metric := family.GetMetric()[0]
			assert.Equal(t, float64(1), metric.GetCounter().GetValue())

			labels := make(map[string]string)
			for _, label := range metric.GetLabel() {
				labels[label.GetName()] = label.GetValue()
			}
			assert.Equal(t, "200", labels["code"])
			assert.Equal(t, http.MethodGet, labels["method"])
			assert.Equal(t, "/hello", labels["path"])
			assert.Contains(t, labels["handler"], "TestPrometheusInstrument")
		case "service_gonic_request_duration":
			duration = true
		case "service_gonic_request_size_bytes":
			requestSize = true
		case "service_gonic_response_size_bytes":
			responseSize = true
		}
	}
	assert.True(t, counter)
	assert.True(t, duration)
	assert.True(t, requestSize)
	assert.True(t, responseSize)

}

func TestPrometheusToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	p := New(Namespace("service"), Registry(prometheus.NewRegistry()), Token("secret"))
	engine := gin.New()
	p.Use(engine)

	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "/metrics", nil)
	require.NoError(t, err)
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, ErrInvalidToken.Error(), recorder.Body.String())

	request, err = http.NewRequestWithContext(t.Context(), http.MethodGet, "/metrics", nil)
	require.NoError(t, err)
	request.Header.Set("Authorization", "Bearer secret")
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
