package devenv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestPtrVal(t *testing.T) {
	tests.BelongsToCISuite(t, "unit")
	x := "hello"
	xptr := ptr(x)
	got := value(xptr)
	require.Equal(t, x, got)

	var y *string
	got = value(y)
	require.Empty(t, got)
}
