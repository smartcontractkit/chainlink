package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGinTestEngine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "creates gin engine without data race",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			engine := gin.New()
			engine.GET("/ping", func(c *gin.Context) {
				c.String(http.StatusOK, "pong")
			})

			w := httptest.NewRecorder()
			req, err := http.NewRequestWithContext(t.Context(), "GET", "/ping", nil)
			require.NoError(t, err)

			engine.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "pong", w.Body.String())
		})
	}
}
