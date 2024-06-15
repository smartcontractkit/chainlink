package smoke

import (
	"fmt"
	"math/big"
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
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestForwarderOCRBasic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig("Smoke", tc.ForwarderOcr)
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
		WithCLNodes(6).
		WithFunding(big.NewFloat(*config.Common.ChainlinkNodeFunding)).
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
	ocrInstances, err := actions_seth.DeployOCRContractsForwarderFlow(
		l,
		sethClient,
		1,
		common.HexToAddress(lt.Address()),
		contracts.ChainlinkClientToChainlinkNodeWithKeysAndAddress(workerNodes),
		authorizedForwarders,
	)
	require.NoError(t, err, "Error deploying OCR contracts")

	err = actions.CreateOCRJobsWithForwarderLocal(ocrInstances, bootstrapNode, workerNodes, 5, env.MockAdapter, fmt.Sprint(sethClient.ChainID))
	require.NoError(t, err, "failed to setup forwarder jobs")
	err = actions_seth.WatchNewOCRRound(l, sethClient, 1, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(10*time.Minute))
	require.NoError(t, err, "error watching for new OCR round")

	answer, err := ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(5), answer.Int64(), "Expected latest answer from OCR contract to be 5 but got %d", answer.Int64())

	err = actions.SetAllAdapterResponsesToTheSameValueLocal(10, ocrInstances, workerNodes, env.MockAdapter)
	require.NoError(t, err)
	err = actions_seth.WatchNewOCRRound(l, sethClient, 2, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(10*time.Minute))
	require.NoError(t, err, "error watching for new OCR round")

	answer, err = ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
}
