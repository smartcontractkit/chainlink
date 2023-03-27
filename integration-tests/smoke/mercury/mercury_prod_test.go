package mercury

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/smoke/mercury/subtests"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
)

// To run this test, provide path to env config file in MERCURY_ENV_CONFIG_PATH
// Example:
//
// {
//     "id": "TestSmokeMercuryProd",
//     "chainId": 420,
//     "feedId": "feed-1",
//     "contracts": {
//         "verifierProxy": "0x42973a598f94Dd6A14a1F2E9CB336Fe88672Fa79"
//     },
//     "mercuryServer": {
//         "remoteUrl": "http://10.14.90.115:3000",
//         "userId": "02185d5a-f1ee-40d1-a52a-bf39871b614c",
//         "userKey": "admintestkey"
//     }
// }

func TestSmokeMercuryProd(t *testing.T) {
	testEnv, err := mercury.NewEnv(t.Name(), "smoke", mercury.DefaultResources)
	require.NoError(t, err)
	if testEnv.C == nil {
		t.Skip("Test is skipped because env config file was not provided")
	}
	feedId := testEnv.C.FeedId
	require.NotEmpty(t, feedId, "'feedId' needs to be provided in the env config file")

	t.Cleanup(func() {
		testEnv.Cleanup(t)
	})

	err = testEnv.AddEvmNetwork()
	require.NoError(t, err)

	msClient := client.NewMercuryServerClient(
		testEnv.MSInfo.RemoteUrl, testEnv.MSInfo.UserId, testEnv.MSInfo.UserKey)
	testEnv.MSClient = msClient

	// verifierProxyContract, err := testEnv.AddVerifierProxyContract("verifierProxy")
	// require.NoError(t, err)
	// exchangerContract, err := testEnv.AddExchangerContract("exchanger", verifierProxyContract.Address(),
	// 	"", 255)
	// _ = exchangerContract
	// require.NoError(t, err)

	subtests.RunTestGetReportByFeedIdForRecentBlockNum(t, &testEnv, feedId, client.StringFeedId)

	// subtests.RunTestGetReportByFeedIdForRecentBlockNum(t, &testEnv, feedId, client.HexFeedId)

	subtests.RunTestGetReportNotFound(t, &testEnv, feedId)

	for i := 0; i < 10000; i++ {
		subtests.RunTestGetReportByFeedIdStrFromWS(t, &testEnv, feedId)
	}

	// subtests.RunTestReportVerificationWithVerifierContract(t, &testEnv, verifierProxyContract, feedId)

	// subtests.RunTestReportVerificationWithExchangerContract(t, &testEnv, exchangerContract, feedId)
}
