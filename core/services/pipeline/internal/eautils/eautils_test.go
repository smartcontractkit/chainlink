package eautils

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBestEffortExtractEAStatus(t *testing.T) {
	tests := []struct {
		name       string
		arg        []byte
		expectCode int
		expectOk   bool
	}{
		{
			name:       "invalid object",
			arg:        []byte(`{"error": "invalid json object" `),
			expectCode: 0,
			expectOk:   false,
		},
		{
			name:       "no status code in object",
			arg:        []byte(`{}`),
			expectCode: 0,
			expectOk:   false,
		},
		{
			name:       "invalid status code",
			arg:        []byte(`{"statusCode":400}`),
			expectCode: http.StatusBadRequest,
			expectOk:   true,
		},
		{
			name:       "invalid provider status code",
			arg:        []byte(`{"statusCode":200, "providerStatusCode":500}`),
			expectCode: http.StatusInternalServerError,
			expectOk:   true,
		},
		{
			name:       "valid statuses with error message",
			arg:        []byte(`{"statusCode":200, "providerStatusCode":200, "error": "unexpected error"}`),
			expectCode: http.StatusInternalServerError,
			expectOk:   true,
		},
		{
			name:       "valid status code",
			arg:        []byte(`{"statusCode":200}`),
			expectCode: http.StatusOK,
			expectOk:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, ok := BestEffortExtractEAStatus(tt.arg)
			assert.Equal(t, tt.expectCode, code)
			assert.Equal(t, tt.expectOk, ok)
		})
	}
}
