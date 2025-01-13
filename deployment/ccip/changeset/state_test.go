package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSmokeState(t *testing.T) {
	tenv, _ := NewMemoryEnvironment(t, WithChains(3))
	state, err := LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	_, err = state.View(tenv.Env.AllChainSelectors())
	require.NoError(t, err)
}
