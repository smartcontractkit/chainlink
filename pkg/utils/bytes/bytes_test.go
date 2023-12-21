package bytes_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/bytes"
	"github.com/stretchr/testify/require"
)

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	require.True(t, bytes.IsEmpty([]byte{0, 0, 0}))
	require.True(t, bytes.IsEmpty([]byte{}))
	require.False(t, bytes.IsEmpty([]byte{1, 2, 3, 5}))
}
