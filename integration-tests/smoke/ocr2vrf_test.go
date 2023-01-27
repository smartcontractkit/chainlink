package smoke

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	eth "github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/ocr2vrf_actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/ocr2vrf_actions/ocr2vrf_constants"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	"github.com/rs/zerolog/log"
)

func TestOCR2VRFBasic(t *testing.T) {
	t.Parallel()
	testEnvironment, testNetwork := setupOCR2VRFTest(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
	require.NoError(t, err)
	contractDeployer, err := contracts.NewContractDeployer(chainClient)
	require.NoError(t, err)
	chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
	require.NoError(t, err)
	nodeAddresses, err := actions.ChainlinkNodeAddresses(chainlinkNodes)
	require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")

	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})

	chainClient.ParallelTransactions(true)

	//1. DEPLOY LINK TOKEN
	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err)

	//2. DEPLOY ETHLINK FEED
	mockETHLinkFeed, err := contractDeployer.DeployMockETHLINKFeed(ocr2vrf_constants.LinkEthFeedResponse)
	require.NoError(t, err)

	//3. Deploy OCR2VRF Contracts (VRFRouter, VRFCoordinator, VRFBeacon, Consumer Contract)
	// Deploy Consumer Contract
	dkg, router, coordinator, vrfBeacon, consumer := ocr2vrf_actions.DeployOCR2VRFContracts(
		t,
		contractDeployer,
		chainClient,
		linkToken,
		mockETHLinkFeed,
		ocr2vrf_constants.BeaconPeriodBlocksCount,
		ocr2vrf_constants.KeyID,
	)

	//4. Register coordinator to router
	err = router.RegisterCoordinator(coordinator.Address())
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	//5. Add VRFBeacon as DKG client
	err = dkg.AddClient(ocr2vrf_constants.KeyID, vrfBeacon.Address())
	require.NoError(t, err)

	//6. Adding VRFBeacon as producer in VRFCoordinator
	err = coordinator.SetProducer(vrfBeacon.Address())
	require.NoError(t, err)

	//7. Subscription:
	//7.1	Create Subscription
	err = coordinator.CreateSubscription()
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)
	subID, err := coordinator.FindSubscriptionID()
	require.NoError(t, err)

	//7.2	Add Consumer to subscription
	err = coordinator.AddConsumer(subID, consumer.Address())
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	//7.3	fund subscription with LINK token
	ocr2vrf_actions.FundVRFCoordinatorSubscription(
		t,
		linkToken,
		coordinator,
		chainClient,
		subID,
		ocr2vrf_constants.LinkFundingAmount,
	)

	//10. set Payees for VRFBeacon ((address which gets the reward) for each transmitter)
	nonBootstrapNodeAddresses := nodeAddresses[1:]
	err = vrfBeacon.SetPayees(nonBootstrapNodeAddresses, nonBootstrapNodeAddresses)
	require.NoError(t, err)

	//11. fund OCR Nodes (so that they can transmit)
	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, ocr2vrf_constants.EthFundingAmount)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	bootstrapNode := chainlinkNodes[0]
	nonBootstrapNodes := chainlinkNodes[1:]

	//11. Create DKG Sign and Encrypt keys for each non-bootstrap node
	// set Job specs for each node
	ocr2VRFPluginConfig := ocr2vrf_actions.SetAndGetOCR2VRFPluginConfig(
		t,
		nonBootstrapNodes,
		dkg,
		vrfBeacon,
		coordinator,
		mockETHLinkFeed,
		ocr2vrf_constants.KeyID,
		ocr2vrf_constants.VRFBeaconAllowedConfirmationDelays,
		ocr2vrf_constants.CoordinatorConfig,
	)
	//12. Create Jobs for Bootstrap and non-boostrap nodes
	ocr2vrf_actions.CreateOCR2VRFJobs(
		t,
		bootstrapNode,
		nonBootstrapNodes,
		ocr2VRFPluginConfig,
		testNetwork.ChainID,
		0,
	)

	//13. set config for DKG OCR,
	//14. wait for the event ConfigSet from DKG contract
	//15. wait for the event Transmitted from DKG contract, meaning that OCR committee has sent out the Public key and Shares
	ocr2vrf_actions.SetAndWaitForDKGProcessToFinish(t, ocr2VRFPluginConfig, dkg)

	//16. set config for VRFBeacon OCR,
	//17. wait for the event ConfigSet from VRFBeacon contract
	ocr2vrf_actions.SetAndWaitForVRFBeaconProcessToFinish(t, ocr2VRFPluginConfig, vrfBeacon)

	//Request and Redeem Randomness
	requestID := ocr2vrf_actions.RequestAndRedeemRandomness(
		t,
		consumer,
		chainClient,
		vrfBeacon,
		ocr2vrf_constants.NumberOfRandomWordsToRequest,
		subID,
		ocr2vrf_constants.ConfirmationDelay,
	)

	for i := uint16(0); i < ocr2vrf_constants.NumberOfRandomWordsToRequest; i++ {
		randomness, err := consumer.GetRandomnessByRequestId(nil, requestID, big.NewInt(int64(i)))
		require.NoError(t, err)
		log.Info().Interface("Random Number", randomness).Interface("Randomness Number Index", i).Msg("Randomness retrieved from Consumer contract")
		require.NotEqual(t, 0, randomness.Uint64(), "Randomness retrieved from Consumer contract give an answer other than 0")
	}
}

func setupOCR2VRFTest(t *testing.T) (testEnvironment *environment.Environment, testNetwork blockchain.EVMNetwork) {
	testNetwork = networks.SelectedNetwork
	evmConfig := eth.New(nil)
	if !testNetwork.Simulated {
		evmConfig = eth.New(&eth.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}

	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-ocr2vrf-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(evmConfig).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": "6",
			"toml": client.AddNetworkDetailedConfig(
				config.BaseOCR2VRFTomlConfig,
				config.DefaultOCR2VRFNetworkDetailTomlConfig, testNetwork),
		}))
	err := testEnvironment.Run()

	require.NoError(t, err, "Error running test environment")

	return testEnvironment, testNetwork
}
