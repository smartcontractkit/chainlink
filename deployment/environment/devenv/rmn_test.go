package devenv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractPeerID(t *testing.T) {
	data := []byte("{\"Salt\":\"yAUVSNGPnY78hJwJwKDBmA==\",\"Nonce\":\"8p0vNKpnQH1P+7/cg2dM3vNlc60tzPl0\",\"Ciphertext\":\"VboRchQzDx/+zIVDWyTXIEwo4Ej0s7O6kQqmgCqW+A+HEXl6bI02W1Y2T88XuY0m44eC9DvSpBPhs/VLtgymj6nBcl+nzfJHLMp2pBMPjKHxgNsEk34mDgnmqYKTtIkgzv7hEyy7j0CAcR/RxkjfQNKpfWOlNtAFryFO+w==\",\"AdditionalData\":\"{\\\"PeerID\\\":\\\"12D3KooWNqugYSJw9thwwu1PC3aEpPsBxMZM2EdzpJXGesRG4E8n\\\"}\"}")
	peerID, err := extractPeerID(data)
	require.NoError(t, err)
	assert.Equal(t, "12D3KooWNqugYSJw9thwwu1PC3aEpPsBxMZM2EdzpJXGesRG4E8n", peerID.String())
}
