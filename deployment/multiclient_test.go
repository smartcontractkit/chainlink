package deployment

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestMultiClient(t *testing.T) {
	lggr := logger.TestLogger(t)
	// Expect an error if no RPCs supplied.
	s := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		b, err := ioutil.ReadAll(request.Body)
		require.NoError(t, err)
		// TODO: Helper struct somewhere for this?
		if string(b) == "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_chainId\"}" {
			writer.WriteHeader(http.StatusOK)
			// Respond with 1337
			_, err = writer.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x539"}`))
			require.NoError(t, err)
			return
		} else {
			// Dial
			writer.WriteHeader(http.StatusOK)
			_, err = writer.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":true}`))
			require.NoError(t, err)
		}
	}))
	defer s.Close()
	// Expect defaults to be set if not provided.
	mc, err := NewMultiClient(lggr, []RPC{{WSURL: s.URL}})
	require.NoError(t, err)
	require.NotNil(t, mc)
	assert.Equal(t, mc.RetryConfig.Attempts, uint(RPC_DEFAULT_RETRY_ATTEMPTS))
	assert.Equal(t, RPC_DEFAULT_RETRY_DELAY, mc.RetryConfig.Delay)

	_, err = NewMultiClient(lggr, []RPC{})
	require.Error(t, err)

	// Expect second client to be set as backup.
	mc, err = NewMultiClient(lggr, []RPC{
		{WSURL: s.URL},
		{WSURL: s.URL},
	})
	require.NoError(t, err)
	require.Len(t, mc.Backups, 1)
}
