package feestest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFail(t *testing.T) {
	require.Equal(t, 0, 1)
}
