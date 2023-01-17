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
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	chainlinkutils "github.com/smartcontractkit/chainlink/core/utils"
	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"

	"github.com/rs/zerolog/log"
)

func TestOCR2VRFBasic(t *testing.T) {
	linkEthFeedResponse := big.NewInt(1e18)
	linkFundingAmount := big.NewInt(100)
	beaconPeriodBlocksCount := big.NewInt(3)
	ethFundingAmount := big.NewFloat(1)
	numberOfRandomWordsToRequest := uint16(2)
	confirmationDelay := big.NewInt(1)
	subscriptionID := uint64(1)
	//keyId can be any random value
	keyId := "aee00d81f822f882b6fe28489822f59ebb21ea95c0ae21d9f67c0239461148fc"

	t.Parallel()
	testEnvironment, testNetwork := setupOCR2VRFTest(t)

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
	dkg, err := contractDeployer.DeployDKG()
	require.NoError(t, err)

	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	//4. Deploy VRFCoordinator(beaconPeriodBlocks, linkAddress, linkEthfeedAddress)
	coordinator, err := contractDeployer.DeployOCR2VRFCoordinator(beaconPeriodBlocksCount, linkToken.Address(), mockETHLinkFeed.Address())
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	//5. Deploy VRFBeacon
	vrfBeacon, err := contractDeployer.DeployVRFBeacon(coordinator.Address(), linkToken.Address(), dkg.Address(), keyId)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	//6. Add VRFBeacon as DKG client
	err = dkg.AddClient(keyId, vrfBeacon.Address())
	require.NoError(t, err)

	//7. Adding VRFBeacon as producer in VRFCoordinator
	err = coordinator.SetProducer(vrfBeacon.Address())
	require.NoError(t, err)

	//8. Deploy Consumer Contract
	consumer, err := contractDeployer.DeployVRFBeaconConsumer(coordinator.Address(), beaconPeriodBlocksCount)
	require.NoError(t, err)

	err = chainClient.WaitForEvents()
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
	encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint64"}]`, subscriptionID)
	require.NoError(t, err)

	_, err = linkToken.TransferAndCall(coordinator.Address(), big.NewInt(0).Mul(linkFundingAmount, big.NewInt(1e18)), encodedSubId)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	//10. set Payees for VRFBeacon ((address which gets the reward) for each transmitter)
	nonBootstrapNodeAddresses := nodeAddresses[1:]
	err = vrfBeacon.SetPayees(nonBootstrapNodeAddresses, nonBootstrapNodeAddresses)
	require.NoError(t, err)

	//11. fund OCR Nodes (so that they can transmit)
	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, ethFundingAmount)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	var dkgKeyConfigs []actions.DKGKeyConfig
	var transmitters []string
	var ocrConfigPubKeys []string
	var peerIDs []string
	var ocrOnchainPubKeys []string
	var ocrOffchainPubKeys []string
	var schedule []int

	bootstrapNode := chainlinkNodes[0]
	nonBootstrapNodes := chainlinkNodes[1:]

	for _, node := range nonBootstrapNodes {
		// 11. Create DKG Sign and Encrypt keys for each non-bootstrap node
		dkgSignKey, err := node.MustCreateDkgSignKey()
		require.NoError(t, err)

		dkgEncryptKey, err := node.MustCreateDkgEncryptKey()
		require.NoError(t, err)

		ethKeys, err := node.MustReadETHKeys()
		require.NoError(t, err)
		for _, key := range ethKeys.Data {
			transmitters = append(transmitters, key.Attributes.Address)
		}

		p2pKeys, err := node.MustReadP2PKeys()
		require.NoError(t, err, "Shouldn't fail reading P2P keys from node")

		peerId := p2pKeys.Data[0].Attributes.PeerID
		peerIDs = append(peerIDs, peerId)

		ocr2Keys, err := node.MustReadOCR2Keys()
		require.NoError(t, err, "Shouldn't fail reading OCR2 keys from node")
		var ocr2Config client.OCR2KeyAttributes
		for _, key := range ocr2Keys.Data {
			if key.Attributes.ChainType == string(chaintype.EVM) {
				ocr2Config = key.Attributes
				break
			}
		}

		offchainPubKey := strings.TrimPrefix(ocr2Config.OffChainPublicKey, "ocr2off_evm_")
		ocrOffchainPubKeys = append(ocrOffchainPubKeys, offchainPubKey)

		configPubKey := strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_")
		ocrConfigPubKeys = append(ocrConfigPubKeys, configPubKey)

		onchainPubKey := strings.TrimPrefix(ocr2Config.OnChainPublicKey, "ocr2on_evm_")
		ocrOnchainPubKeys = append(ocrOnchainPubKeys, onchainPubKey)

		schedule = append(schedule, 1)

		dkgKeyConfigs = append(dkgKeyConfigs, actions.DKGKeyConfig{
			DKGEncryptionPublicKey: dkgEncryptKey.Data.Attributes.PublicKey,
			DKGSigningPublicKey:    dkgSignKey.Data.Attributes.PublicKey,
		})
	}

	//12. set Job specs for each node
	ocr2VRFPluginConfig := &actions.OCR2VRFPluginConfig{
		OCR2Config: actions.OCR2Config{
			OnchainPublicKeys:    ocrOnchainPubKeys,
			OffchainPublicKeys:   ocrOffchainPubKeys,
			ConfigPublicKeys:     ocrConfigPubKeys,
			PeerIds:              peerIDs,
			TransmitterAddresses: transmitters,
			Schedule:             schedule,
		},

		DKGConfig: actions.DKGConfig{
			DKGKeyConfigs:      dkgKeyConfigs,
			DKGKeyID:           keyId,
			DKGContractAddress: dkg.Address(),
		},
		VRFBeaconConfig: actions.VRFBeaconConfig{
			VRFBeaconAddress: vrfBeacon.Address(),
			ConfDelays:       []string{"1", "2", "3", "4", "5", "6", "7", "8"},
			CoordinatorConfig: ocr2vrftypes.CoordinatorConfig{
				CacheEvictionWindowSeconds: 60,
				BatchGasLimit:              5_000_000,
				CoordinatorOverhead:        50_000,
				CallbackOverhead:           50_000,
				BlockGasOverhead:           50_000,
				LookbackBlocks:             1_000,
			},
		},
		VRFCoordinatorAddress: coordinator.Address(),
		LinkEthFeedAddress:    mockETHLinkFeed.Address(),
	}
	actions.CreateOCR2VRFJobs(
		t,
		bootstrapNode,
		nonBootstrapNodes,
		ocr2VRFPluginConfig,
		testNetwork.ChainID,
		0,
	)

	ocr2DkgConfig := actions.BuildOCR2DKGConfigVars(t, ocr2VRFPluginConfig)

	//13. set config for DKG OCR
	log.Debug().Interface("OCR2 DKG Config", ocr2DkgConfig).Msg("OCR2 DKG Config prepared")
	err = dkg.SetConfig(
		ocr2DkgConfig.Signers,
		ocr2DkgConfig.Transmitters,
		ocr2DkgConfig.F,
		ocr2DkgConfig.OnchainConfig,
		ocr2DkgConfig.OffchainConfigVersion,
		ocr2DkgConfig.OffchainConfig,
	)
	require.NoError(t, err)

	//14. wait for the event ConfigSet from DKG contract
	dkgConfigSetEvent, err := dkg.WaitForConfigSetEvent()
	log.Info().Interface("Event", dkgConfigSetEvent).Msg("OCR2 DKG Config was set")
	//15. wait for the event Transmitted from DKG contract, meaning that OCR committee has sent out the Public key and Shares
	dkgSharesTransmittedEvent, err := dkg.WaitForTransmittedEvent()
	log.Info().Interface("Event", dkgSharesTransmittedEvent).Msg("DKG Shares were generated and transmitted by OCR Committee")

	ocr2VrfConfig := actions.BuildOCR2VRFConfigVars(t, ocr2VRFPluginConfig)
	log.Debug().Interface("OCR2 VRF Config", ocr2VrfConfig).Msg("OCR2 VRF Config prepared")

	//16. set config for VRFBeacon OCR
	err = vrfBeacon.SetConfig(
		ocr2VrfConfig.Signers,
		ocr2VrfConfig.Transmitters,
		ocr2VrfConfig.F,
		ocr2VrfConfig.OnchainConfig,
		ocr2VrfConfig.OffchainConfigVersion,
		ocr2VrfConfig.OffchainConfig,
	)
	require.NoError(t, err)

	//15. wait for the event ConfigSet from VRFBeacon contract
	vrfConfigSetEvent, err := vrfBeacon.WaitForConfigSetEvent()
	log.Info().Interface("Event", vrfConfigSetEvent).Msg("OCR2 VRF Config was set")

	//16. Trigger Request randomness method from consumer contract
	receipt, err := consumer.RequestRandomness(
		numberOfRandomWordsToRequest,
		subscriptionID,
		confirmationDelay,
	)
	require.NoError(t, err)
	log.Info().Interface("TX Hash", receipt.TxHash).Msg("Randomness requested from Consumer contract")

	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	periodBlocks, err := consumer.IBeaconPeriodBlocks(nil)
	require.NoError(t, err)

	blockNumber := receipt.BlockNumber
	periodOffset := new(big.Int).Mod(blockNumber, periodBlocks)
	nextBeaconOutputHeight := new(big.Int).Sub(new(big.Int).Add(blockNumber, periodBlocks), periodOffset)

	requestID, err := consumer.GetRequestIdsBy(nil, nextBeaconOutputHeight, confirmationDelay)
	require.NoError(t, err)

	//16. Wait for the NewTransmission event - when OCR committee will fulfill and transmit the randomness number
	newTransmissionEvent, err := vrfBeacon.WaitForNewTransmissionEvent()
	log.Info().Interface("NewTransmission event", newTransmissionEvent).Msg("Randomness transmitted by DON")

	//17. Redeem randomness from consumer contract
	err = consumer.RedeemRandomness(requestID)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()

	for i := uint16(0); i < numberOfRandomWordsToRequest; i++ {
		randomness, err := consumer.GetRandomnessByRequestId(nil, requestID, big.NewInt(int64(i)))
		require.NoError(t, err)
		log.Info().Interface("Random Number", randomness).Interface("Randomness Number Index", i).Msg("Randomness retrieved from Consumer contract")
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

	baseTOML := `[Feature]
LogPoller = true

[OCR2]
Enabled = true

[P2P]
[P2P.V2]
Enabled = true
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]`

	networkDetailTOML := `FinalityDepth = 5
[EVM.GasEstimator]
LimitDefault = 1400000
PriceMax = 100000000000
FeeCapDefault = 100000000000`

	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-ocr2vrf-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
	}).
		AddHelm(evmConfig).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"replicas": "6",
			"toml":     client.AddNetworkDetailedConfig(baseTOML, networkDetailTOML, testNetwork),
		}))
	err := testEnvironment.Run()

	require.NoError(t, err, "Error running test environment")

	return testEnvironment, testNetwork
}
