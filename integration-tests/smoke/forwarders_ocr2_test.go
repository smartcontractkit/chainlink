package smoke

import (
	"fmt"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestForwarderOCR2Basic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.ForwarderOcr2)
	if err != nil {
		t.Fatal(err)
	}

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
		WithMockAdapter().
		WithCLNodeConfig(node.NewConfig(node.NewBaseConfig(),
			node.WithOCR2(),
			node.WithP2Pv2(),
		)).
		WithForwarders().
		WithCLNodes(6).
		WithFunding(big.NewFloat(.1)).
		WithStandardCleanup().
		WithSeth().
		Build()
	require.NoError(t, err)

	nodeClients := env.ClCluster.NodeAPIs()
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	workerNodeAddresses, err := actions.ChainlinkNodeAddressesLocal(workerNodes)
	require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")

	selectedNetwork := networks.MustGetSelectedNetworkConfig(config.Network)[0]
	sethClient, err := env.GetSethClient(selectedNetwork.ChainID)
	require.NoError(t, err, "Error getting seth client")

	lt, err := contracts.DeployLinkTokenContract(l, sethClient)
	require.NoError(t, err, "Deploying Link Token Contract shouldn't fail")

	fundingAmount := big.NewFloat(.05)
	l.Info().Str("ETH amount per node", fundingAmount.String()).Msg("Funding Chainlink nodes")
	err = actions_seth.FundChainlinkNodesFromRootAddress(l, sethClient, contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(workerNodes), fundingAmount)
	require.NoError(t, err, "Error funding Chainlink nodes")

	operators, authorizedForwarders, _ := actions_seth.DeployForwarderContracts(
		t, sethClient, common.HexToAddress(lt.Address()), len(workerNodes),
	)

	require.Equal(t, len(workerNodes), len(operators), "Number of operators should match number of worker nodes")

	for i := range workerNodes {
		actions_seth.AcceptAuthorizedReceiversOperator(
			t, l, sethClient, operators[i], authorizedForwarders[i], []common.Address{workerNodeAddresses[i]},
		)
		require.NoError(t, err, "Accepting Authorize Receivers on Operator shouldn't fail")
		actions_seth.TrackForwarder(t, sethClient, authorizedForwarders[i], workerNodes[i])
	}

	// Gather transmitters
	var transmitters []string
	for _, forwarderCommonAddress := range authorizedForwarders {
		transmitters = append(transmitters, forwarderCommonAddress.Hex())
	}

	ocrOffchainOptions := contracts.DefaultOffChainAggregatorOptions()
	ocrInstances, err := actions_seth.DeployOCRv2Contracts(l, sethClient, 1, common.HexToAddress(lt.Address()), transmitters, ocrOffchainOptions)
	require.NoError(t, err, "Error deploying OCRv2 contracts with forwarders")

	err = actions.CreateOCRv2JobsLocal(ocrInstances, bootstrapNode, workerNodes, env.MockAdapter, "ocr2", 5, uint64(sethClient.ChainID), true, false)
	require.NoError(t, err, "Error creating OCRv2 jobs with forwarders")

	ocrv2Config, err := actions.BuildMedianOCR2ConfigLocal(workerNodes, ocrOffchainOptions)
	require.NoError(t, err, "Error building OCRv2 config")
	ocrv2Config.Transmitters = authorizedForwarders

	err = actions_seth.ConfigureOCRv2AggregatorContracts(ocrv2Config, ocrInstances)
	require.NoError(t, err, "Error configuring OCRv2 aggregator contracts")

	err = actions_seth.WatchNewOCRRound(l, sethClient, 1, contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(10*time.Minute))
	require.NoError(t, err, "error watching for new OCRv2 round")

	answer, err := ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Getting latest answer from OCRv2 contract shouldn't fail")
	require.Equal(t, int64(5), answer.Int64(), "Expected latest answer from OCRw contract to be 5 but got %d", answer.Int64())

	for i := 2; i <= 3; i++ {
		ocrRoundVal := (5 + i) % 10
		err = env.MockAdapter.SetAdapterBasedIntValuePath("ocr2", []string{http.MethodGet, http.MethodPost}, ocrRoundVal)
		require.NoError(t, err)
		err = actions_seth.WatchNewOCRRound(l, sethClient, int64(i), contracts.V2OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(10*time.Minute))
		require.NoError(t, err, "error watching for new OCRv2 round")
		answer, err = ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
		require.NoError(t, err, "Error getting latest OCRv2 answer")
		require.Equal(t, int64(ocrRoundVal), answer.Int64(), fmt.Sprintf("Expected latest answer from OCRv2 contract to be %d but got %d", ocrRoundVal, answer.Int64()))
	}
}
