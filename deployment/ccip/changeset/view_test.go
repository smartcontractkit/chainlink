package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
)

func TestSmokeView(t *testing.T) {
	t.Parallel()
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	jsonData, err := changeset.ViewCCIP(tenv.Env)
	require.NoError(t, err)
	// to ensure the view is valid
	_, err = jsonData.MarshalJSON()
	require.NoError(t, err)
}
