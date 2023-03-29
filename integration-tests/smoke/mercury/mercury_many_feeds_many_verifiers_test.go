package mercury

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	mercuryactions "github.com/smartcontractkit/chainlink/integration-tests/actions/mercury"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/smoke/mercury/subtests"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
)

func TestMercuryManyFeedsManyVerifiers(t *testing.T) {
	feedIds := mercuryactions.GenFeedIds(9)

	testEnv, err := mercury.NewEnv(t.Name(), "smoke", mercury.DefaultResources)
	t.Cleanup(func() {
		testEnv.Cleanup(t)
	})
	require.NoError(t, err)

	testEnv.AddEvmNetwork()

	err = testEnv.AddDON(mercury.GetMockserverResources(len(feedIds)))
	require.NoError(t, err)

	ocrConfig, err := testEnv.BuildOCRConfig()
	require.NoError(t, err)

	_, _, err = testEnv.AddMercuryServer(nil)
	require.NoError(t, err)

	verifierProxyContract, err := testEnv.AddVerifierProxyContract("verifierProxy1")
	require.NoError(t, err)
	exchangerContract, err := testEnv.AddExchangerContract("exchanger1", verifierProxyContract.Address(),
		"", 255)
	require.NoError(t, err)

	// Use separate verifier contracts for each feed
	for i, feedId := range feedIds {
		verifierContractId := fmt.Sprintf("verifier_%d", i)
		verifierContract, err := testEnv.AddVerifierContract(verifierContractId, verifierProxyContract.Address())
		require.NoError(t, err)

		blockNumber, err := testEnv.SetConfigAndInitializeVerifierContract(
			fmt.Sprintf("setAndInitializeVerifier%d", i),
			verifierContractId,
			"verifierProxy1",
			feedId,
			*ocrConfig,
		)
		require.NoError(t, err)

		err = testEnv.AddBootstrapJob(fmt.Sprintf("createBoostrap%d", i), verifierContract.Address(), uint64(blockNumber), feedId)
		require.NoError(t, err)

		err = testEnv.AddOCRJobs(fmt.Sprintf("createOcrJobs%d", i), verifierContract.Address(), uint64(blockNumber), feedId)
		require.NoError(t, err)
	}

	err = testEnv.WaitForReportsInMercuryDb(feedIds)
	require.NoError(t, err)

	for _, feedId := range feedIds {
		feedIdStr := mercury.Byte32ToString(feedId)

		subtests.RunTestGetReportByFeedIdForRecentBlockNum(t, &testEnv, feedIdStr, client.StringFeedId)
		subtests.RunTestGetReportByFeedIdForRecentBlockNum(t, &testEnv, feedIdStr, client.HexFeedId)
		subtests.RunTestReportVerificationWithExchangerContract(t, &testEnv, exchangerContract, feedIdStr)
	}
}
