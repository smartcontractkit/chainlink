package chaos

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-env/chaos"
	"github.com/smartcontractkit/chainlink-env/environment"
	a "github.com/smartcontractkit/chainlink-env/pkg/alias"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
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

var (
	defaultOCR2VRFSettings = map[string]interface{}{
		"replicas": "6",
		"toml": client.AddNetworkDetailedConfig(
			config.BaseOCR2VRFTomlConfig,
			config.DefaultOCR2VRFNetworkDetailTomlConfig,
			networks.SelectedNetwork),
	}

	defaultOCR2VRFEthereumSettings = &ethereum.Props{
		NetworkName: networks.SelectedNetwork.Name,
		Simulated:   networks.SelectedNetwork.Simulated,
		WsURLs:      networks.SelectedNetwork.URLs,
	}
)

func TestOCR2VRFChaos(t *testing.T) {

	testCases := map[string]struct {
		networkChart environment.ConnectedChart
		clChart      environment.ConnectedChart
		chaosFunc    chaos.ManifestFunc
		chaosProps   *chaos.Props
	}{
		// network-* and pods-* are split intentionally into 2 parallel groups
		// we can't use chaos.NewNetworkPartition and chaos.NewFailPods in parallel
		// because of jsii runtime bug, see Makefile

		// PodChaosFailMinorityNodes Test description:
		//1. DKG and VRF beacon processes are set and VRF request gets fulfilled
		//2. Apply chaos experiment - take down 2 nodes out of 5 non-bootstrap
		//3. Bring back all nodes to normal
		//4. verify VRF request gets fulfilled
		PodChaosFailMinorityNodes: {
			ethereum.New(defaultOCR2VRFEthereumSettings),
			chainlink.New(0, defaultOCR2VRFSettings),
			chaos.NewFailPods,
			&chaos.Props{
				LabelsSelector: &map[string]*string{ChaosGroupMinority: a.Str("1")},
				DurationStr:    "1m",
			},
		},
		//todo - currently failing, need to investigate deeper
		//PodChaosFailMajorityNodes: {
		//	ethereum.New(defaultOCR2VRFEthereumSettings),
		//	chainlink.New(0, defaultOCR2VRFSettings),
		//	chaos.NewFailPods,
		//	&chaos.Props{
		//		LabelsSelector: &map[string]*string{ChaosGroupMajority: a.Str("1")},
		//		DurationStr:    "1m",
		//	},
		//},
		//todo - do we need these chaos tests?
		//PodChaosFailMajorityDB: {
		//	ethereum.New(defaultOCR2VRFEthereumSettings),
		//	chainlink.New(0, defaultOCR2VRFSettings),
		//	chaos.NewFailPods,
		//	&chaos.Props{
		//		LabelsSelector: &map[string]*string{ChaosGroupMajority: a.Str("1")},
		//		DurationStr:    "1m",
		//		ContainerNames: &[]*string{a.Str("chainlink-db")},
		//	},
		//},
		//NetworkChaosFailMajorityNetwork: {
		//	ethereum.New(defaultOCR2VRFEthereumSettings),
		//	chainlink.New(0, defaultOCR2VRFSettings),
		//	chaos.NewNetworkPartition,
		//	&chaos.Props{
		//		FromLabels:  &map[string]*string{ChaosGroupMajority: a.Str("1")},
		//		ToLabels:    &map[string]*string{ChaosGroupMinority: a.Str("1")},
		//		DurationStr: "1m",
		//	},
		//},
		//NetworkChaosFailBlockchainNode: {
		//	ethereum.New(defaultOCR2VRFEthereumSettings),
		//	chainlink.New(0, defaultOCR2VRFSettings),
		//	chaos.NewNetworkPartition,
		//	&chaos.Props{
		//		FromLabels:  &map[string]*string{"app": a.Str("geth")},
		//		ToLabels:    &map[string]*string{ChaosGroupMajority: a.Str("1")},
		//		DurationStr: "1m",
		//	},
		//},
	}

	for testcaseName, testCase := range testCases {
		t.Run(testcaseName, func(t *testing.T) {
			t.Parallel()
			testNetwork := networks.SelectedNetwork
			testEnvironment := environment.
				New(&environment.Config{
					NamespacePrefix: fmt.Sprintf(
						"chaos-ocr2vrf-%s",
						strings.ReplaceAll(strings.ToLower(testNetwork.Name),
							" ",
							"-")),
					Test: t,
				}).
				AddHelm(testCase.networkChart).
				AddHelm(testCase.clChart)
			err := testEnvironment.Run()
			require.NoError(t, err, "Error running test environment")

			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 1, 2, ChaosGroupMinority)
			require.NoError(t, err)
			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 3, 5, ChaosGroupMajority)
			require.NoError(t, err)

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

			//3. Deploy DKG contract
			//4. Deploy VRFRouter
			//5. Deploy VRFCoordinator(beaconPeriodBlocks, linkAddress, linkEthfeedAddress)
			//6. Deploy VRFBeacon
			//7. Deploy Consumer Contract
			dkg, router, coordinator, vrfBeacon, consumer := ocr2vrf_actions.DeployOCR2VRFContracts(
				t,
				contractDeployer,
				chainClient,
				linkToken,
				mockETHLinkFeed,
				ocr2vrf_constants.BeaconPeriodBlocksCount,
				ocr2vrf_constants.KeyID,
			)

			//8. Register coordnator
			err = router.RegisterCoordinator(coordinator.Address())
			require.NoError(t, err)

			//9. Add VRFBeacon as DKG client
			err = dkg.AddClient(ocr2vrf_constants.KeyID, vrfBeacon.Address())
			require.NoError(t, err)

			//10. Adding VRFBeacon as producer in VRFCoordinator
			err = coordinator.SetProducer(vrfBeacon.Address())
			require.NoError(t, err)

			//11. Subscription:
			//11.1	Create Subscription
			err = coordinator.CreateSubscription()
			require.NoError(t, err)
			err = chainClient.WaitForEvents()
			require.NoError(t, err)
			subID, err := coordinator.FindSubscriptionID()
			require.NoError(t, err)

			//11.2	Add Consumer to subscription
			err = coordinator.AddConsumer(subID, consumer.Address())
			require.NoError(t, err)
			err = chainClient.WaitForEvents()
			require.NoError(t, err)

			//11.3	fund subscription with LINK token
			ocr2vrf_actions.FundVRFCoordinatorSubscription(
				t,
				linkToken,
				coordinator,
				chainClient,
				subID,
				ocr2vrf_constants.LinkFundingAmount,
			)

			//12. set Payees for VRFBeacon ((address which gets the reward) for each transmitter)
			nonBootstrapNodeAddresses := nodeAddresses[1:]
			err = vrfBeacon.SetPayees(nonBootstrapNodeAddresses, nonBootstrapNodeAddresses)
			require.NoError(t, err)

			//13. fund OCR Nodes (so that they can transmit)
			err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, ocr2vrf_constants.EthFundingAmount)
			require.NoError(t, err)
			err = chainClient.WaitForEvents()
			require.NoError(t, err)

			bootstrapNode := chainlinkNodes[0]
			nonBootstrapNodes := chainlinkNodes[1:]

			//14. Create DKG Sign and Encrypt keys for each non-bootstrap node
			//15. set Job specs for each node
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
			//16. Create Jobs for Bootstrap and non-boostrap nodes
			ocr2vrf_actions.CreateOCR2VRFJobs(
				t,
				bootstrapNode,
				nonBootstrapNodes,
				ocr2VRFPluginConfig,
				testNetwork.ChainID,
				0,
			)

			//17. set config for DKG OCR,
			//18. wait for the event ConfigSet from DKG contract
			//19. wait for the event Transmitted from DKG contract, meaning that OCR committee has sent out the Public key and Shares
			ocr2vrf_actions.SetAndWaitForDKGProcessToFinish(t, ocr2VRFPluginConfig, dkg)

			//20. set config for VRFBeacon OCR,
			//21. wait for the event ConfigSet from VRFBeacon contract
			ocr2vrf_actions.SetAndWaitForVRFBeaconProcessToFinish(t, ocr2VRFPluginConfig, vrfBeacon)

			//Request and Redeem Randomness to verify that process works fine
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

			id, err := testEnvironment.Chaos.Run(testCase.chaosFunc(testEnvironment.Cfg.Namespace, testCase.chaosProps))
			require.NoError(t, err)
			log.Info().Msg("Chaos Applied")

			err = testEnvironment.Chaos.WaitForAllRecovered(id)
			require.NoError(t, err)
			log.Info().Msg("Chaos Recovered")

			//Request and Redeem Randomness again to see that after Chaos Experiment whole process is still working
			requestID = ocr2vrf_actions.RequestAndRedeemRandomness(
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
		})
	}
}
