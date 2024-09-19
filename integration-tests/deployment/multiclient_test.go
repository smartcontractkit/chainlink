package deployment

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiClient(t *testing.T) {
	// Expect an error if no RPCs supplied.
	s := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, err := writer.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":true}`))
		require.NoError(t, err)
	}))
	defer s.Close()
	_, err := NewMultiClient([]RPC{})
	require.Error(t, err)

	// Expect defaults to be set if not provided.
	mc, err := NewMultiClient([]RPC{{WSURL: s.URL}})
	require.NoError(t, err)
	assert.Equal(t, mc.RetryConfig.Attempts, uint(RPC_DEFAULT_RETRY_ATTEMPTS))
	assert.Equal(t, mc.RetryConfig.Delay, RPC_DEFAULT_RETRY_DELAY)

	// Expect second client to be set as backup.
	mc, err = NewMultiClient([]RPC{
		{WSURL: s.URL},
		{WSURL: s.URL},
	})
	require.NoError(t, err)
	require.Equal(t, len(mc.Backups), 1)
}
