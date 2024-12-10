package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSmokeView(t *testing.T) {
	tenv := NewMemoryEnvironment(t, WithChains(3))
	_, err := ViewCCIP(tenv.Env)
	require.NoError(t, err)
}
