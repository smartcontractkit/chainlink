package smoke

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/integration-tests/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestForwarderOCR2Basic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig([]string{"Smoke"}, tc.ForwarderOcr2)
	require.NoError(t, err, "Error getting config")

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
		WithMockAdapter().
		WithCLNodes(6).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	nodeClients := env.ClCluster.NodeAPIs()
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	workerNodeAddresses, err := actions.ChainlinkNodeAddressesLocal(workerNodes)
	require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")

	evmNetwork, err := env.GetFirstEvmNetwork()
	require.NoError(t, err, "Error getting first evm network")

	sethClient, err := utils.TestAwareSethClient(t, config, evmNetwork)
	require.NoError(t, err, "Error getting seth client")

	err = actions.FundChainlinkNodesFromRootAddress(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(env.ClCluster.NodeAPIs()), big.NewFloat(*config.Common.ChainlinkNodeFunding))
	require.NoError(t, err, "Failed to fund the nodes")

	t.Cleanup(func() {
		// ignore error, we will see failures in the logs anyway
		_ = actions.ReturnFundsFromNodes(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(env.ClCluster.NodeAPIs()))
	})

	linkContract, err := actions.LinkTokenContract(l, sethClient, config.OCR2)
	require.NoError(t, err, "Error loading/deploying link token contract")

	fundingAmount := big.NewFloat(.05)
	l.Info().Str("ETH amount per node", fundingAmount.String()).Msg("Funding Chainlink nodes")
	err = actions.FundChainlinkNodesFromRootAddress(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(workerNodes), fundingAmount)
	require.NoError(t, err, "Error funding Chainlink nodes")

	operators, authorizedForwarders, _ := actions.DeployForwarderContracts(
		t, sethClient, common.HexToAddress(linkContract.Address()), len(workerNodes),
	)

	require.Equal(t, len(workerNodes), len(operators), "Number of operators should match number of worker nodes")

	for i := range workerNodes {
		actions.AcceptAuthorizedReceiversOperator(
			t, l, sethClient, operators[i], authorizedForwarders[i], []common.Address{workerNodeAddresses[i]},
		)
		require.NoError(t, err, "Accepting Authorize Receivers on Operator shouldn't fail")
		actions.TrackForwarder(t, sethClient, authorizedForwarders[i], workerNodes[i])
	}

	// Gather transmitters
	var transmitters []string
	for _, forwarderCommonAddress := range authorizedForwarders {
		transmitters = append(transmitters, forwarderCommonAddress.Hex())
	}

	ocrOffchainOptions := contracts.DefaultOffChainAggregatorOptions()
	ocrInstances, err := actions.SetupOCRv2Contracts(l, sethClient, config.OCR2, common.HexToAddress(linkContract.Address()), transmitters, ocrOffchainOptions)
	require.NoError(t, err, "Error deploying OCRv2 contracts with forwarders")

	ocrv2Config, err := actions.BuildMedianOCR2ConfigLocal(workerNodes, ocrOffchainOptions)
	require.NoError(t, err, "Error building OCRv2 config")
	ocrv2Config.Transmitters = authorizedForwarders

	err = actions.ConfigureOCRv2AggregatorContracts(ocrv2Config, ocrInstances)
	require.NoError(t, err, "Error configuring OCRv2 aggregator contracts")

	err = actions.CreateOCRv2JobsLocal(ocrInstances, bootstrapNode, workerNodes, env.MockAdapter, "ocr2", 5, uint64(sethClient.ChainID), true, false)
	require.NoError(t, err, "Error creating OCRv2 jobs with forwarders")

	err = actions.WatchNewOCRRound(l, sethClient, 1, contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(10*time.Minute))
	require.NoError(t, err, "error watching for new OCRv2 round")

	answer, err := ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Getting latest answer from OCRv2 contract shouldn't fail")
	require.Equal(t, int64(5), answer.Int64(), "Expected latest answer from OCRw contract to be 5 but got %d", answer.Int64())

	for i := 2; i <= 3; i++ {
		ocrRoundVal := (5 + i) % 10
		err = env.MockAdapter.SetAdapterBasedIntValuePath("ocr2", []string{http.MethodGet, http.MethodPost}, ocrRoundVal)
		require.NoError(t, err)
		err = actions.WatchNewOCRRound(l, sethClient, int64(i), contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(10*time.Minute))
		require.NoError(t, err, "error watching for new OCRv2 round")
		answer, err = ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
		require.NoError(t, err, "Error getting latest OCRv2 answer")
		require.Equal(t, int64(ocrRoundVal), answer.Int64(), fmt.Sprintf("Expected latest answer from OCRv2 contract to be %d but got %d", ocrRoundVal, answer.Int64()))
	}
}
