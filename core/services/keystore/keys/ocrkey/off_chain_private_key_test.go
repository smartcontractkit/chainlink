package ocrkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestOCRKeys_OffChainPrivateKey(t *testing.T) {
	t.Parallel()
	tests.BelongsToCISuite(t, "unit")

	k, err := NewV2()
	require.NoError(t, err)

	sig, err := k.OffChainSigning.Sign([]byte("hello world"))

	assert.NoError(t, err)
	assert.NotEmpty(t, sig)
}
