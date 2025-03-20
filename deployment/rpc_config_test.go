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

func TestURLSchemePreferenceFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected URLSchemePreference
		wantErr  string
	}{
		{input: "none", expected: URLSchemePreferenceNone},
		{input: "ws", expected: URLSchemePreferenceWS},
		{input: "http", expected: URLSchemePreferenceHTTP},
		{input: "invalid", wantErr: "invalid URLSchemePreference: invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := URLSchemePreferenceFromString(tt.input)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestURLSchemePreference_UnmarshalText(t *testing.T) {
	tests := []struct {
		input    []byte
		expected URLSchemePreference
		wantErr  string
	}{
		{input: []byte("none"), expected: URLSchemePreferenceNone},
		{input: []byte("ws"), expected: URLSchemePreferenceWS},
		{input: []byte("http"), expected: URLSchemePreferenceHTTP},
		{input: []byte("invalid"), wantErr: "invalid URLSchemePreference: invalid"},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			var got URLSchemePreference
			err := got.UnmarshalText(tt.input)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, got)
			}
		})
	}
}
