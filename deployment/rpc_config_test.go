package deployment

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRPC_ToEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		rpc      RPC
		expected string
		wantErr  string
	}{
		{
			name: "No URL scheme specified defaults to URLSchemePreferenceNone, returns WS URL",
			rpc: RPC{
				Name:    "TestRPC",
				WSURL:   "ws://example.com",
				HTTPURL: "http://example.com",
			},
			expected: "ws://example.com",
		},
		{
			name: "Preferred URL scheme None, returns WS URL",
			rpc: RPC{
				Name:               "TestRPC",
				WSURL:              "ws://example.com",
				HTTPURL:            "http://example.com",
				PreferredURLScheme: URLSchemePreferenceNone,
			},
			expected: "ws://example.com",
		},
		{
			name: "Preferred URL scheme WS, returns WS URL",
			rpc: RPC{
				Name:               "TestRPC",
				WSURL:              "ws://example.com",
				HTTPURL:            "http://example.com",
				PreferredURLScheme: URLSchemePreferenceWS,
			},
			expected: "ws://example.com",
		},
		{
			name: "Preferred URL scheme HTTP, returns HTTP URL",
			rpc: RPC{
				Name:               "TestRPC",
				WSURL:              "ws://example.com",
				HTTPURL:            "http://example.com",
				PreferredURLScheme: URLSchemePreferenceHTTP,
			},
			expected: "http://example.com",
		},
		{
			name: "Unknown URL scheme, returns error",
			rpc: RPC{
				Name:               "TestRPC",
				WSURL:              "ws://example.com",
				HTTPURL:            "http://example.com",
				PreferredURLScheme: URLSchemePreference(999),
			},
			wantErr: "Unknown URLSchemePreference",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.rpc.ToEndpoint()

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, got)
			}
		})
	}
}
