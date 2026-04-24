package devenv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPtrVal(t *testing.T) {
	x := "hello"
	xptr := ptr(x)
	got := value(xptr)
	require.Equal(t, x, got)

	var y *string
	got = value(y)
	require.Empty(t, got)
}
