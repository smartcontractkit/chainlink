package ocr2vrf_actions

import (
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

func FundVRFCoordinatorSubscription(t *testing.T, linkToken contracts.LinkToken, coordinator contracts.VRFCoordinatorV3, chainClient blockchain.EVMClient, subscriptionID, linkFundingAmount *big.Int) {
	encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint256"}]`, subscriptionID)
	require.NoError(t, err)

	_, err = linkToken.TransferAndCall(coordinator.Address(), big.NewInt(0).Mul(linkFundingAmount, big.NewInt(1e18)), encodedSubId)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)
}

func DeployOCR2VRFContracts(t *testing.T, contractDeployer contracts.ContractDeployer, chainClient blockchain.EVMClient, linkToken contracts.LinkToken, mockETHLinkFeed contracts.MockETHLINKFeed, beaconPeriodBlocksCount *big.Int, keyID string) (contracts.DKG, contracts.VRFRouter, contracts.VRFCoordinatorV3, contracts.VRFBeacon, contracts.VRFBeaconConsumer) {
	dkg, err := contractDeployer.DeployDKG()
	require.NoError(t, err)

	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	router, err := contractDeployer.DeployVRFRouter()
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	coordinator, err := contractDeployer.DeployOCR2VRFCoordinator(beaconPeriodBlocksCount, linkToken.Address(), mockETHLinkFeed.Address(), router.Address())
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
	return dkg, router, coordinator, vrfBeacon, consumer
}

func RequestAndRedeemRandomness(t *testing.T, consumer contracts.VRFBeaconConsumer, chainClient blockchain.EVMClient, vrfBeacon contracts.VRFBeacon, numberOfRandomWordsToRequest uint16, subscriptionID, confirmationDelay *big.Int) *big.Int {
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

	newTransmissionEvent, err := vrfBeacon.WaitForNewTransmissionEvent()
	log.Info().Interface("NewTransmission event", newTransmissionEvent).Msg("Randomness transmitted by DON")

	err = consumer.RedeemRandomness(subscriptionID, requestID)
	require.NoError(t, err)
	err = chainClient.WaitForEvents()
	require.NoError(t, err)

	return requestID
}
