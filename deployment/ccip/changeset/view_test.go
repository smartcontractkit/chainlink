package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSmokeView(t *testing.T) {
	t.Parallel()
	tenv := NewMemoryEnvironment(t, WithChains(3))
	_, err := ViewCCIP(tenv.Env)
	require.NoError(t, err)
}
