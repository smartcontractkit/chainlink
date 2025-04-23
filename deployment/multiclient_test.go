package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// TODO(giogam): This test is incomplete, it should be completed with support for websockets URLS
func TestMultiClient(t *testing.T) {
	var (
		lggr                 = logger.TestLogger(t)
		chainSelector uint64 = 16015286601757825753 // "ethereum-testnet-sepolia"
		wsURL                = "ws://example.com"
		httpURL              = "http://example.com"
	)

	// Expect defaults to be set if not provided.
	mc, err := NewMultiClient(lggr, RPCConfig{ChainSelector: chainSelector, RPCs: []RPC{
		{Name: "test-rpc", WSURL: wsURL, HTTPURL: httpURL, PreferredURLScheme: URLSchemePreferenceHTTP},
	}})

	require.NoError(t, err)
	require.NotNil(t, mc)

	assert.Equal(t, "ethereum-testnet-sepolia", mc.chainName)
	assert.Equal(t, mc.RetryConfig.Attempts, uint(RPCDefaultRetryAttempts))
	assert.Equal(t, RPCDefaultRetryDelay, mc.RetryConfig.Delay)

	// Expect error if no RPCs provided.
	_, err = NewMultiClient(lggr, RPCConfig{ChainSelector: chainSelector, RPCs: []RPC{}})
	require.Error(t, err)

	// Expect second client to be set as backup.
	mc, err = NewMultiClient(lggr, RPCConfig{ChainSelector: chainSelector, RPCs: []RPC{
		{Name: "test-rpc", WSURL: wsURL, HTTPURL: httpURL, PreferredURLScheme: URLSchemePreferenceHTTP},
		{Name: "test-rpc", WSURL: wsURL, HTTPURL: httpURL, PreferredURLScheme: URLSchemePreferenceHTTP},
	}})
	require.NoError(t, err)
	require.Len(t, mc.Backups, 1)
}
