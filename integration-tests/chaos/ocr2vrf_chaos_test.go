package chaos

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-env/chaos"
	"github.com/smartcontractkit/chainlink-env/environment"
	a "github.com/smartcontractkit/chainlink-env/pkg/alias"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/ocr2vrf_actions"

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

	linkEthFeedResponse          = big.NewInt(1e18)
	linkFundingAmount            = big.NewInt(100)
	beaconPeriodBlocksCount      = big.NewInt(3)
	ethFundingAmount             = big.NewFloat(1)
	numberOfRandomWordsToRequest = uint16(1)
	confirmationDelay            = big.NewInt(1)
	subscriptionID               = uint64(1)
	//keyID can be any random value
	keyID = "aee00d81f822f882b6fe28489822f59ebb21ea95c0ae21d9f67c0239461148fc"

	coordinatorConfig = &ocr2vrftypes.CoordinatorConfig{
		CacheEvictionWindowSeconds: 60,
		BatchGasLimit:              5_000_000,
		CoordinatorOverhead:        50_000,
		CallbackOverhead:           50_000,
		BlockGasOverhead:           50_000,
		LookbackBlocks:             1_000,
	}
	vrfBeaconAllowedConfirmationDelays = []string{"1", "2", "3", "4", "5", "6", "7", "8"}
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

			//todo - 3 nodes instead of 2 are down when specifying 1-2
			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 1, 1, ChaosGroupMinority)
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
			mockETHLinkFeed, err := contractDeployer.DeployMockETHLINKFeed(linkEthFeedResponse)
			require.NoError(t, err)

			//3. Deploy DKG contract
			//4. Deploy VRFCoordinator(beaconPeriodBlocks, linkAddress, linkEthfeedAddress)
			//5. Deploy VRFBeacon
			//8. Deploy Consumer Contract
			dkg, coordinator, vrfBeacon, consumer := ocr2vrf_actions.DeployOCR2VRFContracts(t, contractDeployer, chainClient, linkToken, mockETHLinkFeed, beaconPeriodBlocksCount, keyID)

			//6. Add VRFBeacon as DKG client
			err = dkg.AddClient(keyID, vrfBeacon.Address())
			require.NoError(t, err)

			//7. Adding VRFBeacon as producer in VRFCoordinator
			err = coordinator.SetProducer(vrfBeacon.Address())
			require.NoError(t, err)

			//9. Subscription:
			//9.1	Create Subscription
			err = coordinator.CreateSubscription()
			require.NoError(t, err)
			err = chainClient.WaitForEvents()
			require.NoError(t, err)

			//9.2	Add Consumer to subscription
			err = coordinator.AddConsumer(subscriptionID, consumer.Address())
			require.NoError(t, err)
			err = chainClient.WaitForEvents()
			require.NoError(t, err)

			//9.3	fund subscription with LINK token
			ocr2vrf_actions.FundVRFCoordinatorSubscription(t, linkToken, coordinator, chainClient, subscriptionID, linkFundingAmount)

			//10. set Payees for VRFBeacon ((address which gets the reward) for each transmitter)
			nonBootstrapNodeAddresses := nodeAddresses[1:]
			err = vrfBeacon.SetPayees(nonBootstrapNodeAddresses, nonBootstrapNodeAddresses)
			require.NoError(t, err)

			//11. fund OCR Nodes (so that they can transmit)
			err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, ethFundingAmount)
			require.NoError(t, err)
			err = chainClient.WaitForEvents()
			require.NoError(t, err)

			bootstrapNode := chainlinkNodes[0]
			nonBootstrapNodes := chainlinkNodes[1:]

			//11. Create DKG Sign and Encrypt keys for each non-bootstrap node
			//12. set Job specs for each node
			ocr2VRFPluginConfig := ocr2vrf_actions.SetAndGetOCR2VRFPluginConfig(t, nonBootstrapNodes, dkg, vrfBeacon, coordinator, mockETHLinkFeed, keyID, vrfBeaconAllowedConfirmationDelays, coordinatorConfig)
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
			ocr2vrf_actions.SetAndWaitForDKGProcessToFinish(t, ocr2VRFPluginConfig, err, dkg)

			//16. set config for VRFBeacon OCR,
			//17. wait for the event ConfigSet from VRFBeacon contract
			ocr2vrf_actions.SetAndWaitForVRFBeaconProcessToFinish(t, ocr2VRFPluginConfig, err, vrfBeacon)

			//todo - do we really need to perform requestAndRedeemRandomness() before Chaos experiment is applied?
			//Request and Redeem Randomness
			requestID := ocr2vrf_actions.RequestAndRedeemRandomness(t, consumer, chainClient, vrfBeacon, numberOfRandomWordsToRequest, subscriptionID, confirmationDelay)

			g := gomega.NewGomegaWithT(t)
			for i := uint16(0); i < numberOfRandomWordsToRequest; i++ {
				randomness, err := consumer.GetRandomnessByRequestId(nil, requestID, big.NewInt(int64(i)))
				require.NoError(t, err)
				log.Info().Interface("Random Number", randomness).Interface("Randomness Number Index", i).Msg("Randomness retrieved from Consumer contract")
				g.Expect(randomness.Uint64()).ShouldNot(gomega.BeNumerically("==", 0), "Randomness retrieved from Consumer contract give an answer other than 0")
			}

			id, err := testEnvironment.Chaos.Run(testCase.chaosFunc(testEnvironment.Cfg.Namespace, testCase.chaosProps))
			require.NoError(t, err)
			log.Info().Msg("Chaos Applied")

			err = testEnvironment.Chaos.WaitForAllRecovered(id)
			require.NoError(t, err)
			log.Info().Msg("Chaos Recovered")

			//Request and Redeem Randomness again to see that after Chaos Experiment whole process is still working
			requestID = ocr2vrf_actions.RequestAndRedeemRandomness(t, consumer, chainClient, vrfBeacon, numberOfRandomWordsToRequest, subscriptionID, confirmationDelay)

			for i := uint16(0); i < numberOfRandomWordsToRequest; i++ {
				randomness, err := consumer.GetRandomnessByRequestId(nil, requestID, big.NewInt(int64(i)))
				require.NoError(t, err)
				log.Info().Interface("Random Number", randomness).Interface("Randomness Number Index", i).Msg("Randomness retrieved from Consumer contract")
				g.Expect(randomness.Uint64()).ShouldNot(gomega.BeNumerically("==", 0), "Randomness retrieved from Consumer contract give an answer other than 0")
			}

		})
	}
}
