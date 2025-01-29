package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
)

func TestSmokeState(t *testing.T) {
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := changeset.LoadOnchainState(tenv.Env)
	require.NoError(t, err)
	_, err = state.View(tenv.Env.AllChainSelectors())
	require.NoError(t, err)
}

// TODO: add solana state test
