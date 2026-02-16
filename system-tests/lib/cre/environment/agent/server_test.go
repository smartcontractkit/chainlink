package agent

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestStartComponentReturnsErrorCodeForUnsupportedSchema(t *testing.T) {
	server := NewServer(zerolog.Nop(), nil)
	handler := server.Handler()

	req := httptest.NewRequest(http.MethodPost, "/v1/components/start", strings.NewReader(`{"schemaVersion":"v0","operation":"StartComponent","payload":{}}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), ErrCodeUnsupportedSchema) {
		t.Fatalf("expected response to include error code %q, got body: %s", ErrCodeUnsupportedSchema, rr.Body.String())
	}
}

func TestStartComponentReturnsErrorCodeForUnsupportedComponent(t *testing.T) {
	server := NewServer(zerolog.Nop(), nil)
	handler := server.Handler()

	body := bytes.NewBufferString(`{"schemaVersion":"v1","operation":"StartComponent","payload":{"componentType":"jd"}}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/components/start", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), ErrCodeUnsupportedComponent) {
		t.Fatalf("expected response to include error code %q, got body: %s", ErrCodeUnsupportedComponent, rr.Body.String())
	}
}
