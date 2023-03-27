package mercury

import (
	"testing"

	"github.com/stretchr/testify/require"

	mercuryactions "github.com/smartcontractkit/chainlink/integration-tests/actions/mercury"
	"github.com/smartcontractkit/chainlink/integration-tests/smoke/mercury/subtests"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
)

func TestMercuryManyFeedsSingleVerifier(t *testing.T) {
	feedIds := mercuryactions.GenFeedIds(9)

	testEnv, verifierProxyContract, err := mercury.SetupMultiFeedSingleVerifierEnv(t.Name(), "smoke", feedIds, mercury.DefaultResources)
	t.Cleanup(func() {
		testEnv.Cleanup(t)
	})
	require.NoError(t, err)

	exchangerContract, err := testEnv.AddExchangerContract("exchanger1", verifierProxyContract.Address(),
		"", 255)
	require.NoError(t, err)

	for _, feedIdBytes := range feedIds {
		feedIdStr := mercury.Byte32ToString(feedIdBytes)

		subtests.RunTestGetReportByFeedIdStringForRecentBlockNum(t, testEnv, feedIdStr)
		subtests.RunTestReportVerificationWithExchangerContract(t, testEnv, exchangerContract, feedIdStr)
	}
}
