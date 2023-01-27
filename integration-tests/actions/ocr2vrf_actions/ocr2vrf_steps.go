package ocr2vrf_actions

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/actions/ocr2vrf_actions/ocr2vrf_constants"
	"math/big"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	chainlinkutils "github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func SetAndWaitForVRFBeaconProcessToFinish(t *testing.T, ocr2VRFPluginConfig *OCR2VRFPluginConfig, vrfBeacon contracts.VRFBeacon) {
	ocr2VrfConfig := BuildOCR2VRFConfigVars(t, ocr2VRFPluginConfig)
	log.Debug().Interface("OCR2 VRF Config", ocr2VrfConfig).Msg("OCR2 VRF Config prepared")

	err := vrfBeacon.SetConfig(
		ocr2VrfConfig.Signers,
		ocr2VrfConfig.Transmitters,
		ocr2VrfConfig.F,
		ocr2VrfConfig.OnchainConfig,
		ocr2VrfConfig.OffchainConfigVersion,
		ocr2VrfConfig.OffchainConfig,
	)
	require.NoError(t, err)

	vrfConfigSetEvent, err := vrfBeacon.WaitForConfigSetEvent()
	require.NoError(t, err)
	log.Info().Interface("Event", vrfConfigSetEvent).Msg("OCR2 VRF Config was set")
}

func SetAndWaitForDKGProcessToFinish(t *testing.T, ocr2VRFPluginConfig *OCR2VRFPluginConfig, dkg contracts.DKG) {
	ocr2DkgConfig := BuildOCR2DKGConfigVars(t, ocr2VRFPluginConfig)

	// set config for DKG OCR
	log.Debug().Interface("OCR2 DKG Config", ocr2DkgConfig).Msg("OCR2 DKG Config prepared")
	err := dkg.SetConfig(
		ocr2DkgConfig.Signers,
		ocr2DkgConfig.Transmitters,
		ocr2DkgConfig.F,
		ocr2DkgConfig.OnchainConfig,
		ocr2DkgConfig.OffchainConfigVersion,
		ocr2DkgConfig.OffchainConfig,
	)
	require.NoError(t, err)

	// wait for the event ConfigSet from DKG contract
	dkgConfigSetEvent, err := dkg.WaitForConfigSetEvent()
	require.NoError(t, err)
	log.Info().Interface("Event", dkgConfigSetEvent).Msg("OCR2 DKG Config was set")
	// wait for the event Transmitted from DKG contract, meaning that OCR committee has sent out the Public key and Shares
	dkgSharesTransmittedEvent, err := dkg.WaitForTransmittedEvent()
	require.NoError(t, err)
	log.Info().Interface("Event", dkgSharesTransmittedEvent).Msg("DKG Shares were generated and transmitted by OCR Committee")
}

func SetAndGetOCR2VRFPluginConfig(t *testing.T, nonBootstrapNodes []*client.Chainlink, dkg contracts.DKG, vrfBeacon contracts.VRFBeacon, coordinator contracts.VRFCoordinatorV3, mockETHLinkFeed contracts.MockETHLINKFeed, keyID string, vrfBeaconAllowedConfirmationDelays []string, coordinatorConfig *ocr2vrftypes.CoordinatorConfig) *OCR2VRFPluginConfig {
	var dkgKeyConfigs []DKGKeyConfig
	var transmitters []string
	var ocrConfigPubKeys []string
	var peerIDs []string
	var ocrOnchainPubKeys []string
	var ocrOffchainPubKeys []string
	var schedule []int

	for _, node := range nonBootstrapNodes {
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

func FundVRFCoordinatorSubscription(t *testing.T, linkToken contracts.LinkToken, coordinator contracts.VRFCoordinatorV3, chainClient blockchain.EVMClient, subscriptionID uint64, linkFundingAmount *big.Int) {
	encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint64"}]`, subscriptionID)
	require.NoError(t, err)

	_, err = linkToken.TransferAndCall(coordinator.Address(), big.NewInt(0).Mul(linkFundingAmount, big.NewInt(1e18)), encodedSubId)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)
}

func DeployOCR2VRFContracts(t *testing.T, contractDeployer contracts.ContractDeployer, chainClient blockchain.EVMClient, linkToken contracts.LinkToken, mockETHLinkFeed contracts.MockETHLINKFeed, beaconPeriodBlocksCount *big.Int, keyID string) (contracts.DKG, contracts.VRFCoordinatorV3, contracts.VRFBeacon, contracts.VRFBeaconConsumer) {
	dkg, err := contractDeployer.DeployDKG()
	require.NoError(t, err)

	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	coordinator, err := contractDeployer.DeployOCR2VRFCoordinator(beaconPeriodBlocksCount, linkToken.Address(), mockETHLinkFeed.Address())
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	vrfBeacon, err := contractDeployer.DeployVRFBeacon(coordinator.Address(), linkToken.Address(), dkg.Address(), keyID)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	consumer, err := contractDeployer.DeployVRFBeaconConsumer(coordinator.Address(), beaconPeriodBlocksCount)
	require.NoError(t, err)

	err = chainClient.WaitForEvents()
	require.NoError(t, err)
	return dkg, coordinator, vrfBeacon, consumer
}

func RequestAndRedeemRandomness(t *testing.T, consumer contracts.VRFBeaconConsumer, chainClient blockchain.EVMClient, vrfBeacon contracts.VRFBeacon, numberOfRandomWordsToRequest uint16, subscriptionID uint64, confirmationDelay *big.Int) *big.Int {
	receipt, err := consumer.RequestRandomness(
		numberOfRandomWordsToRequest,
		subscriptionID,
		confirmationDelay,
	)
	require.NoError(t, err)
	log.Info().Interface("TX Hash", receipt.TxHash).Msg("Randomness requested from Consumer contract")

	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	requestID := getRequestId(t, consumer, receipt, confirmationDelay)

	newTransmissionEvent, err := vrfBeacon.WaitForNewTransmissionEvent()
	log.Info().Interface("NewTransmission event", newTransmissionEvent).Msg("Randomness transmitted by DON")

	err = consumer.RedeemRandomness(requestID)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	return requestID
}

func RequestRandomnessFulfillment(
	t *testing.T,
	consumer contracts.VRFBeaconConsumer,
	chainClient blockchain.EVMClient,
	vrfBeacon contracts.VRFBeacon,
	numberOfRandomWordsToRequest uint16,
	subscriptionID uint64,
	confirmationDelay *big.Int,
) *big.Int {
	receipt, err := consumer.RequestRandomnessFulfillment(
		numberOfRandomWordsToRequest,
		subscriptionID,
		confirmationDelay,
		100_000,
		nil,
	)
	require.NoError(t, err)
	log.Info().Interface("TX Hash", receipt.TxHash).Msg("Randomness Fulfillment requested from Consumer contract")

	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	requestID := getRequestId(t, consumer, receipt, confirmationDelay)

	newTransmissionEvent, err := vrfBeacon.WaitForNewTransmissionEvent()
	log.Info().Interface("NewTransmission event", newTransmissionEvent).Msg("Randomness Fulfillment transmitted by DON")

	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	return requestID
}

func getRequestId(t *testing.T, consumer contracts.VRFBeaconConsumer, receipt *types.Receipt, confirmationDelay *big.Int) *big.Int {
	periodBlocks, err := consumer.IBeaconPeriodBlocks(nil)
	require.NoError(t, err)

	blockNumber := receipt.BlockNumber
	periodOffset := new(big.Int).Mod(blockNumber, periodBlocks)
	nextBeaconOutputHeight := new(big.Int).Sub(new(big.Int).Add(blockNumber, periodBlocks), periodOffset)

	requestID, err := consumer.GetRequestIdsBy(nil, nextBeaconOutputHeight, confirmationDelay)
	require.NoError(t, err)
	return requestID
}

func SetupOCR2VRFUniverse(
	t *testing.T,
	linkToken contracts.LinkToken,
	mockETHLinkFeed contracts.MockETHLINKFeed,
	contractDeployer contracts.ContractDeployer,
	chainClient blockchain.EVMClient,
	nodeAddresses []common.Address,
	chainlinkNodes []*client.Chainlink,
	testNetwork blockchain.EVMNetwork,
) (contracts.DKG, contracts.VRFCoordinatorV3, contracts.VRFBeacon, contracts.VRFBeaconConsumer) {

	// Deploy DKG contract
	// Deploy VRFCoordinator(beaconPeriodBlocks, linkAddress, linkEthfeedAddress)
	// Deploy VRFBeacon
	// Deploy Consumer Contract
	dkgContract, coordinatorContract, vrfBeaconContract, consumerContract := DeployOCR2VRFContracts(
		t,
		contractDeployer,
		chainClient,
		linkToken,
		mockETHLinkFeed,
		ocr2vrf_constants.BeaconPeriodBlocksCount,
		ocr2vrf_constants.KeyID,
	)

	// Add VRFBeacon as DKG client
	err := dkgContract.AddClient(ocr2vrf_constants.KeyID, vrfBeaconContract.Address())
	require.NoError(t, err)

	// Adding VRFBeacon as producer in VRFCoordinator
	err = coordinatorContract.SetProducer(vrfBeaconContract.Address())
	require.NoError(t, err)

	// Subscription:
	//1.	Create Subscription
	err = coordinatorContract.CreateSubscription()
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	//2.	Add Consumer to subscription
	err = coordinatorContract.AddConsumer(ocr2vrf_constants.SubscriptionID, consumerContract.Address())
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	//3.	fund subscription with LINK token
	FundVRFCoordinatorSubscription(
		t,
		linkToken,
		coordinatorContract,
		chainClient,
		ocr2vrf_constants.SubscriptionID,
		ocr2vrf_constants.LinkFundingAmount,
	)

	// set Payees for VRFBeacon ((address which gets the reward) for each transmitter)
	nonBootstrapNodeAddresses := nodeAddresses[1:]
	err = vrfBeaconContract.SetPayees(nonBootstrapNodeAddresses, nonBootstrapNodeAddresses)
	require.NoError(t, err)

	// fund OCR Nodes (so that they can transmit)
	err = actions.FundChainlinkNodes(chainlinkNodes, chainClient, ocr2vrf_constants.EthFundingAmount)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

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
	return dkgContract, coordinatorContract, vrfBeaconContract, consumerContract
}
