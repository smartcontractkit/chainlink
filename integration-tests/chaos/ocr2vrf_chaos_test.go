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
	"go.uber.org/zap/zapcore"

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
	t.Parallel()
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

	for testcaseName, tc := range testCases {
		testCase := tc
		t.Run(fmt.Sprintf("OCR2VRF_%s", testcaseName), func(t *testing.T) {
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
			if testEnvironment.WillUseRemoteRunner() {
				return
			}

			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 1, 2, ChaosGroupMinority)
			require.NoError(t, err)
			err = testEnvironment.Client.LabelChaosGroup(testEnvironment.Cfg.Namespace, 3, 5, ChaosGroupMajority)
			require.NoError(t, err)

			chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment)
			require.NoError(t, err, "Error connecting to blockchain")
			contractDeployer, err := contracts.NewContractDeployer(chainClient)
			require.NoError(t, err, "Error building contract deployer")
			chainlinkNodes, err := client.ConnectChainlinkNodes(testEnvironment)
			require.NoError(t, err, "Error connecting to Chainlink nodes")
			nodeAddresses, err := actions.ChainlinkNodeAddresses(chainlinkNodes)
			require.NoError(t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")

			t.Cleanup(func() {
				err := actions.TeardownSuite(t, testEnvironment, utils.ProjectRoot, chainlinkNodes, nil, zapcore.PanicLevel, chainClient)
				require.NoError(t, err, "Error tearing down environment")
			})

			chainClient.ParallelTransactions(true)

			linkToken, err := contractDeployer.DeployLinkTokenContract()
			require.NoError(t, err, "Error deploying LINK token")

			mockETHLinkFeed, err := contractDeployer.DeployMockETHLINKFeed(ocr2vrf_constants.LinkEthFeedResponse)
			require.NoError(t, err, "Error deploying Mock ETH/LINK Feed")

			_, _, vrfBeaconContract, consumerContract, subID := ocr2vrf_actions.SetupOCR2VRFUniverse(
				t,
				linkToken,
				mockETHLinkFeed,
				contractDeployer,
				chainClient,
				nodeAddresses,
				chainlinkNodes,
				testNetwork,
			)

			//Request and Redeem Randomness to verify that process works fine
			requestID := ocr2vrf_actions.RequestAndRedeemRandomness(
				t,
				consumerContract,
				chainClient,
				vrfBeaconContract,
				ocr2vrf_constants.NumberOfRandomWordsToRequest,
				subID,
				ocr2vrf_constants.ConfirmationDelay,
			)

			for i := uint16(0); i < ocr2vrf_constants.NumberOfRandomWordsToRequest; i++ {
				randomness, err := consumerContract.GetRandomnessByRequestId(nil, requestID, big.NewInt(int64(i)))
				require.NoError(t, err)
				log.Info().Interface("Random Number", randomness).Interface("Randomness Number Index", i).Msg("Randomness retrieved from Consumer contract")
				require.NotEqual(t, 0, randomness.Uint64(), "Randomness retrieved from Consumer contract give an answer other than 0")
			}

			id, err := testEnvironment.Chaos.Run(testCase.chaosFunc(testEnvironment.Cfg.Namespace, testCase.chaosProps))
			require.NoError(t, err, "Error running Chaos Experiment")
			log.Info().Msg("Chaos Applied")

			err = testEnvironment.Chaos.WaitForAllRecovered(id)
			require.NoError(t, err, "Error waiting for Chaos Experiment to end")
			log.Info().Msg("Chaos Recovered")

			//Request and Redeem Randomness again to see that after Chaos Experiment whole process is still working
			requestID = ocr2vrf_actions.RequestAndRedeemRandomness(
				t,
				consumerContract,
				chainClient,
				vrfBeaconContract,
				ocr2vrf_constants.NumberOfRandomWordsToRequest,
				subID,
				ocr2vrf_constants.ConfirmationDelay,
			)

			for i := uint16(0); i < ocr2vrf_constants.NumberOfRandomWordsToRequest; i++ {
				randomness, err := consumerContract.GetRandomnessByRequestId(nil, requestID, big.NewInt(int64(i)))
				require.NoError(t, err, "Error getting Randomness result from Consumer Contract")
				log.Info().Interface("Random Number", randomness).Interface("Randomness Number Index", i).Msg("Randomness retrieved from Consumer contract")
				require.NotEqual(t, 0, randomness.Uint64(), "Randomness retrieved from Consumer contract give an answer other than 0")
			}
		})
	}
}
