package ocrkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOCRKeys_OffChainPrivateKey(t *testing.T) {
	t.Parallel()

	k, err := NewV2()
	require.NoError(t, err)

	sig, err := k.OffChainSigning.Sign([]byte("hello world"))

	assert.NoError(t, err)
	assert.NotEmpty(t, sig)
}
