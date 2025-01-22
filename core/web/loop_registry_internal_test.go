package web

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type responseWriter struct {
	statusCode int
	header     http.Header
}

func (r *responseWriter) Write(b []byte) (int, error) {
	return 0, errors.New("could not write to response")
}

func (r *responseWriter) Header() http.Header {
	return http.Header{}
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

func newResponseWriter() *responseWriter {
	return &responseWriter{header: make(http.Header)}
}

func TestLoopRegistryServer_CantWriteToResponse(t *testing.T) {
	l, o := logger.TestLoggerObserved(t, zap.ErrorLevel)
	s := &LoopRegistryServer{
		exposedPromPort: 1,
		registry:        plugins.NewTestLoopRegistry(l),
		logger:          l.(logger.SugaredLogger),
		jsonMarshalFn:   json.Marshal,
	}

	rw := newResponseWriter()
	s.discoveryHandler(rw, &http.Request{})
	assert.Equal(t, rw.statusCode, http.StatusInternalServerError)
	assert.Equal(t, 1, o.FilterMessageSnippet("could not write to response").Len())
}

func TestLoopRegistryServer_CantMarshal(t *testing.T) {
	l, o := logger.TestLoggerObserved(t, zap.ErrorLevel)
	s := &LoopRegistryServer{
		exposedPromPort: 1,
		registry:        plugins.NewTestLoopRegistry(l),
		logger:          l.(logger.SugaredLogger),
		jsonMarshalFn: func(any) ([]byte, error) {
			return []byte(""), errors.New("can't unmarshal")
		},
	}

	rw := newResponseWriter()
	s.discoveryHandler(rw, &http.Request{})
	assert.Equal(t, rw.statusCode, http.StatusInternalServerError)
	assert.Equal(t, 1, o.FilterMessageSnippet("could not write to response").Len())
}
