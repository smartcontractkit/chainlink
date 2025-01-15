package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSmokeView(t *testing.T) {
	t.Parallel()
	tenv, _ := NewMemoryEnvironment(t, WithChains(3))
	_, err := ViewCCIP(tenv.Env)
	require.NoError(t, err)
}
