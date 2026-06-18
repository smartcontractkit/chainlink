package modifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_aptosPublicKeyHexToAccountAddress(t *testing.T) {
	tests := []struct {
		name      string
		hexPubKey string
		wantAddr  string
		wantErr   string
	}{
		{
			name:      "valid ed25519 public key",
			hexPubKey: "0da83f37b1491159df6064accee571393825affeb2e58ea70914cf1e45939d2d",
			wantAddr:  "6ac0430dfd584bfa13d6b2090ea17c872e2e8d9c567046881ac61b21ca67e368",
		},
		{
			name:      "valid with 0x prefix",
			hexPubKey: "0x0da83f37b1491159df6064accee571393825affeb2e58ea70914cf1e45939d2d",
			wantAddr:  "6ac0430dfd584bfa13d6b2090ea17c872e2e8d9c567046881ac61b21ca67e368",
		},
		{
			name:      "invalid hex",
			hexPubKey: "zzzz",
			wantErr:   "parse ed25519 public key",
		},
		{
			name:      "wrong length",
			hexPubKey: "aabbccdd",
			wantErr:   "parse ed25519 public key",
		},
		{
			name:      "empty string",
			hexPubKey: "",
			wantErr:   "parse ed25519 public key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := aptosPublicKeyHexToAccountAddress(tt.hexPubKey)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantAddr, got)
		})
	}
}
