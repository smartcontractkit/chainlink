package ocr2vrf_actions

import (
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/seth"

	ocr2vrftypes "github.com/smartcontractkit/chainlink-vrf/types"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"

	actions_seth "github.com/smartcontractkit/chainlink/integration-tests/actions/seth"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"

	"github.com/smartcontractkit/chainlink/integration-tests/actions/ocr2vrf_actions/ocr2vrf_constants"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func SetAndWaitForVRFBeaconProcessToFinish(t *testing.T, ocr2VRFPluginConfig *OCR2VRFPluginConfig, vrfBeacon contracts.VRFBeacon) {
	l := logging.GetTestLogger(t)
	ocr2VrfConfig := BuildOCR2VRFConfigVars(t, ocr2VRFPluginConfig)
	l.Debug().Interface("OCR2 VRF Config", ocr2VrfConfig).Msg("OCR2 VRF Config prepared")

	err := vrfBeacon.SetConfig(
		ocr2VrfConfig.Signers,
		ocr2VrfConfig.Transmitters,
		ocr2VrfConfig.F,
		ocr2VrfConfig.OnchainConfig,
		ocr2VrfConfig.OffchainConfigVersion,
		ocr2VrfConfig.OffchainConfig,
	)
	require.NoError(t, err, "Error setting OCR config for VRFBeacon contract")

	vrfConfigSetEvent, err := vrfBeacon.WaitForConfigSetEvent(time.Minute)
	require.NoError(t, err, "Error waiting for ConfigSet Event for VRFBeacon contract")
	l.Info().Interface("Event", vrfConfigSetEvent).Msg("OCR2 VRF Config was set")
}

func SetAndWaitForDKGProcessToFinish(t *testing.T, ocr2VRFPluginConfig *OCR2VRFPluginConfig, dkg contracts.DKG) {
	l := logging.GetTestLogger(t)
	ocr2DkgConfig := BuildOCR2DKGConfigVars(t, ocr2VRFPluginConfig)

	// set config for DKG OCR
	l.Debug().Interface("OCR2 DKG Config", ocr2DkgConfig).Msg("OCR2 DKG Config prepared")
	err := dkg.SetConfig(
		ocr2DkgConfig.Signers,
		ocr2DkgConfig.Transmitters,
		ocr2DkgConfig.F,
		ocr2DkgConfig.OnchainConfig,
		ocr2DkgConfig.OffchainConfigVersion,
		ocr2DkgConfig.OffchainConfig,
	)
	require.NoError(t, err, "Error setting OCR config for DKG contract")

	// wait for the event ConfigSet from DKG contract
	dkgConfigSetEvent, err := dkg.WaitForConfigSetEvent(time.Minute)
	require.NoError(t, err, "Error waiting for ConfigSet Event for DKG contract")
	l.Info().Interface("Event", dkgConfigSetEvent).Msg("OCR2 DKG Config Set")
	// wait for the event Transmitted from DKG contract, meaning that OCR committee has sent out the Public key and Shares
	dkgSharesTransmittedEvent, err := dkg.WaitForTransmittedEvent(time.Minute * 5)
	require.NoError(t, err)
	l.Info().Interface("Event", dkgSharesTransmittedEvent).Msg("DKG Shares were generated and transmitted by OCR Committee")
}

func SetAndGetOCR2VRFPluginConfig(
	t *testing.T,
	nonBootstrapNodes []*client.ChainlinkK8sClient,
	dkg contracts.DKG,
	vrfBeacon contracts.VRFBeacon,
	coordinator contracts.VRFCoordinatorV3,
	mockETHLinkFeed contracts.MockETHLINKFeed,
	keyID string,
	vrfBeaconAllowedConfirmationDelays []string,
	coordinatorConfig *ocr2vrftypes.CoordinatorConfig,
) *OCR2VRFPluginConfig {
	var (
		dkgKeyConfigs      []DKGKeyConfig
		transmitters       []string
		ocrConfigPubKeys   []string
		peerIDs            []string
		ocrOnchainPubKeys  []string
		ocrOffchainPubKeys []string
		schedule           []int
	)

	for _, node := range nonBootstrapNodes {
		dkgSignKey, err := node.MustCreateDkgSignKey()
		require.NoError(t, err, "Error creating DKG Sign Keys")

		dkgEncryptKey, err := node.MustCreateDkgEncryptKey()
		require.NoError(t, err, "Error creating DKG Encrypt Keys")

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

		dkgKeyConfigs = append(dkgKeyConfigs, DKGKeyConfig{
			DKGEncryptionPublicKey: dkgEncryptKey.Data.Attributes.PublicKey,
			DKGSigningPublicKey:    dkgSignKey.Data.Attributes.PublicKey,
		})
	}

	ocr2VRFPluginConfig := &OCR2VRFPluginConfig{
		OCR2Config: OCR2Config{
			OnchainPublicKeys:    ocrOnchainPubKeys,
			OffchainPublicKeys:   ocrOffchainPubKeys,
			ConfigPublicKeys:     ocrConfigPubKeys,
			PeerIds:              peerIDs,
			TransmitterAddresses: transmitters,
			Schedule:             schedule,
		},

		DKGConfig: DKGConfig{
			DKGKeyConfigs:      dkgKeyConfigs,
			DKGKeyID:           keyID,
			DKGContractAddress: dkg.Address(),
		},
		VRFBeaconConfig: VRFBeaconConfig{
			VRFBeaconAddress:  vrfBeacon.Address(),
			ConfDelays:        vrfBeaconAllowedConfirmationDelays,
			CoordinatorConfig: coordinatorConfig,
		},
		VRFCoordinatorAddress: coordinator.Address(),
		LinkEthFeedAddress:    mockETHLinkFeed.Address(),
	}
	return ocr2VRFPluginConfig
}

func FundVRFCoordinatorV3Subscription(t *testing.T, linkToken contracts.LinkToken, coordinator contracts.VRFCoordinatorV3, subscriptionID, linkFundingAmount *big.Int) {
	encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint256"}]`, subscriptionID)
	require.NoError(t, err, "Error Abi encoding subscriptionID")
	_, err = linkToken.TransferAndCall(coordinator.Address(), big.NewInt(0).Mul(linkFundingAmount, big.NewInt(1e18)), encodedSubId)
	require.NoError(t, err, "Error sending Link token")
}

func DeployOCR2VRFContracts(t *testing.T, chainClient *seth.Client, linkToken contracts.LinkToken, beaconPeriodBlocksCount *big.Int, keyID string) (contracts.DKG, contracts.VRFCoordinatorV3, contracts.VRFBeacon, contracts.VRFBeaconConsumer) {
	dkg, err := contracts.DeployDKG(chainClient)
	require.NoError(t, err, "Error deploying DKG Contract")

	coordinator, err := contracts.DeployOCR2VRFCoordinator(chainClient, beaconPeriodBlocksCount, linkToken.Address())
	require.NoError(t, err, "Error deploying OCR2VRFCoordinator Contract")

	vrfBeacon, err := contracts.DeployVRFBeacon(chainClient, coordinator.Address(), linkToken.Address(), dkg.Address(), keyID)
	require.NoError(t, err, "Error deploying VRFBeacon Contract")

	consumer, err := contracts.DeployVRFBeaconConsumer(chainClient, coordinator.Address(), beaconPeriodBlocksCount)
	require.NoError(t, err, "Error deploying VRFBeaconConsumer Contract")

	return dkg, coordinator, vrfBeacon, consumer
}

func RequestAndRedeemRandomness(
	t *testing.T,
	consumer contracts.VRFBeaconConsumer,
	vrfBeacon contracts.VRFBeacon,
	numberOfRandomWordsToRequest uint16,
	subscriptionID,
	confirmationDelay *big.Int,
	randomnessTransmissionEventTimeout time.Duration,
) *big.Int {
	l := logging.GetTestLogger(t)
	receipt, err := consumer.RequestRandomness(
		numberOfRandomWordsToRequest,
		subscriptionID,
		confirmationDelay,
	)
	require.NoError(t, err, "Error requesting randomness from Consumer Contract")
	l.Info().Interface("TX Hash", receipt.TxHash).Msg("Randomness requested from Consumer contract")

	requestID := getRequestId(t, consumer, receipt, confirmationDelay)

	newTransmissionEvent, err := vrfBeacon.WaitForNewTransmissionEvent(randomnessTransmissionEventTimeout)
	require.NoError(t, err, "Error waiting for NewTransmission event from VRF Beacon Contract")
	l.Info().Interface("NewTransmission event", newTransmissionEvent).Msg("Randomness transmitted by DON")

	err = consumer.RedeemRandomness(subscriptionID, requestID)
	require.NoError(t, err, "Error redeeming randomness from Consumer Contract")

	return requestID
}

func RequestRandomnessFulfillmentAndWaitForFulfilment(
	t *testing.T,
	consumer contracts.VRFBeaconConsumer,
	vrfBeacon contracts.VRFBeacon,
	numberOfRandomWordsToRequest uint16,
	subscriptionID *big.Int,
	confirmationDelay *big.Int,
	randomnessTransmissionEventTimeout time.Duration,
) *big.Int {
	l := logging.GetTestLogger(t)
	receipt, err := consumer.RequestRandomnessFulfillment(
		numberOfRandomWordsToRequest,
		subscriptionID,
		confirmationDelay,
		200_000,
		100_000,
		nil,
	)
	require.NoError(t, err, "Error requesting Randomness Fulfillment")
	l.Info().Interface("TX Hash", receipt.TxHash).Msg("Randomness Fulfillment requested from Consumer contract")

	requestID := getRequestId(t, consumer, receipt, confirmationDelay)

	newTransmissionEvent, err := vrfBeacon.WaitForNewTransmissionEvent(randomnessTransmissionEventTimeout)
	require.NoError(t, err, "Error waiting for NewTransmission event from VRF Beacon Contract")
	l.Info().Interface("NewTransmission event", newTransmissionEvent).Msg("Randomness Fulfillment transmitted by DON")

	return requestID
}

func getRequestId(t *testing.T, consumer contracts.VRFBeaconConsumer, receipt *types.Receipt, confirmationDelay *big.Int) *big.Int {
	periodBlocks, err := consumer.IBeaconPeriodBlocks(testcontext.Get(t))
	require.NoError(t, err, "Error getting Beacon Period block count")

	blockNumber := receipt.BlockNumber
	periodOffset := new(big.Int).Mod(blockNumber, periodBlocks)
	nextBeaconOutputHeight := new(big.Int).Sub(new(big.Int).Add(blockNumber, periodBlocks), periodOffset)

	requestID, err := consumer.GetRequestIdsBy(testcontext.Get(t), nextBeaconOutputHeight, confirmationDelay)
	require.NoError(t, err, "Error getting requestID from consumer contract")

	return requestID
}

func SetupOCR2VRFUniverse(
	t *testing.T,
	linkToken contracts.LinkToken,
	mockETHLinkFeed contracts.MockETHLINKFeed,
	chainClient *seth.Client,
	nodeAddresses []common.Address,
	chainlinkNodes []*client.ChainlinkK8sClient,
	testNetwork blockchain.EVMNetwork,
) (contracts.DKG, contracts.VRFCoordinatorV3, contracts.VRFBeacon, contracts.VRFBeaconConsumer, *big.Int) {
	l := logging.GetTestLogger(t)

	// Deploy DKG contract
	// Deploy VRFCoordinator(beaconPeriodBlocks, linkAddress, linkEthfeedAddress)
	// Deploy VRFBeacon
	// Deploy Consumer Contract
	dkgContract, coordinatorContract, vrfBeaconContract, consumerContract := DeployOCR2VRFContracts(
		t,
		chainClient,
		linkToken,
		ocr2vrf_constants.BeaconPeriodBlocksCount,
		ocr2vrf_constants.KeyID,
	)

	// Add VRFBeacon as DKG client
	err := dkgContract.AddClient(ocr2vrf_constants.KeyID, vrfBeaconContract.Address())
	require.NoError(t, err, "Error adding client to DKG Contract")
	// Adding VRFBeacon as producer in VRFCoordinator
	err = coordinatorContract.SetProducer(vrfBeaconContract.Address())
	require.NoError(t, err, "Error setting Producer for VRFCoordinator contract")
	err = coordinatorContract.SetConfig(2.5e6, 160 /* 5 EVM words */)
	require.NoError(t, err, "Error setting config for VRFCoordinator contract")

	// Subscription:
	//1.	Create Subscription
	err = coordinatorContract.CreateSubscription()
	require.NoError(t, err, "Error creating subscription in VRFCoordinator contract")
	subID, err := coordinatorContract.FindSubscriptionID()
	require.NoError(t, err)

	//2.	Add Consumer to subscription
	err = coordinatorContract.AddConsumer(subID, consumerContract.Address())
	require.NoError(t, err, "Error adding a consumer to a subscription in VRFCoordinator contract")

	//3.	fund subscription with LINK token
	FundVRFCoordinatorV3Subscription(
		t,
		linkToken,
		coordinatorContract,
		subID,
		ocr2vrf_constants.LinkFundingAmount,
	)

	// set Payees for VRFBeacon ((address which gets the reward) for each transmitter)
	nonBootstrapNodeAddresses := nodeAddresses[1:]
	err = vrfBeaconContract.SetPayees(nonBootstrapNodeAddresses, nonBootstrapNodeAddresses)
	require.NoError(t, err, "Error setting Payees in VRFBeacon Contract")

	// fund OCR Nodes (so that they can transmit)
	nodes := contracts.ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(chainlinkNodes)
	err = actions_seth.FundChainlinkNodesFromRootAddress(l, chainClient, nodes, ocr2vrf_constants.EthFundingAmount)
	require.NoError(t, err, "Error funding Nodes")

	bootstrapNode := chainlinkNodes[0]
	nonBootstrapNodes := chainlinkNodes[1:]

	// Create DKG Sign and Encrypt keys for each non-bootstrap node
	// set Job specs for each node
	ocr2VRFPluginConfig := SetAndGetOCR2VRFPluginConfig(
		t,
		nonBootstrapNodes,
		dkgContract,
		vrfBeaconContract,
		coordinatorContract,
		mockETHLinkFeed,
		ocr2vrf_constants.KeyID,
		ocr2vrf_constants.VRFBeaconAllowedConfirmationDelays,
		ocr2vrf_constants.CoordinatorConfig,
	)
	// Create Jobs for Bootstrap and non-boostrap nodes
	CreateOCR2VRFJobs(
		t,
		bootstrapNode,
		nonBootstrapNodes,
		ocr2VRFPluginConfig,
		testNetwork.ChainID,
		0,
	)

	// set config for DKG OCR,
	// wait for the event ConfigSet from DKG contract
	// wait for the event Transmitted from DKG contract, meaning that OCR committee has sent out the Public key and Shares
	SetAndWaitForDKGProcessToFinish(t, ocr2VRFPluginConfig, dkgContract)

	// set config for VRFBeacon OCR,
	// wait for the event ConfigSet from VRFBeacon contract
	SetAndWaitForVRFBeaconProcessToFinish(t, ocr2VRFPluginConfig, vrfBeaconContract)
	return dkgContract, coordinatorContract, vrfBeaconContract, consumerContract, subID
}
