package ocr2vrf_actions

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	"github.com/smartcontractkit/ocr2vrf/dkg"
	"github.com/smartcontractkit/ocr2vrf/ocr2vrf"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

// CreateOCR2VRFJobs bootstraps the first node and to the other nodes sends ocr jobs
func CreateOCR2VRFJobs(
	t *testing.T,
	bootstrapNode *client.Chainlink,
	nonBootstrapNodes []*client.Chainlink,
	OCR2VRFPluginConfig *OCR2VRFPluginConfig,
	chainID int64,
	keyIndex int,
) {
	p2pV2Bootstrapper := createBootstrapJob(t, bootstrapNode, OCR2VRFPluginConfig.DKGConfig.DKGContractAddress, chainID)

	createNonBootstrapJobs(t, nonBootstrapNodes, OCR2VRFPluginConfig, chainID, keyIndex, p2pV2Bootstrapper)
	log.Info().Msg("Done creating OCR automation jobs")
}

func createNonBootstrapJobs(
	t *testing.T,
	nonBootstrapNodes []*client.Chainlink,
	OCR2VRFPluginConfig *OCR2VRFPluginConfig,
	chainID int64,
	keyIndex int,
	P2Pv2Bootstrapper string,
) {
	for index, nonBootstrapNode := range nonBootstrapNodes {
		nodeTransmitterAddress, err := nonBootstrapNode.EthAddresses()
		require.NoError(t, err, "Shouldn't fail getting primary ETH address from OCR node %d", index)
		nodeOCRKeys, err := nonBootstrapNode.MustReadOCR2Keys()
		require.NoError(t, err, "Shouldn't fail getting OCR keys from OCR node %d", index)
		var nodeOCRKeyId []string
		for _, key := range nodeOCRKeys.Data {
			if key.Attributes.ChainType == string(chaintype.EVM) {
				nodeOCRKeyId = append(nodeOCRKeyId, key.ID)
				break
			}
		}

		OCR2VRFJobSpec := client.OCR2TaskJobSpec{
			Name:    "ocr2",
			JobType: "offchainreporting2",
			OCR2OracleSpec: job.OCR2OracleSpec{
				PluginType: "ocr2vrf",
				Relay:      "evm",
				RelayConfig: map[string]interface{}{
					"chainID": int(chainID),
				},
				ContractID:         OCR2VRFPluginConfig.VRFBeaconConfig.VRFBeaconAddress,
				OCRKeyBundleID:     null.StringFrom(nodeOCRKeyId[keyIndex]),
				TransmitterID:      null.StringFrom(nodeTransmitterAddress[keyIndex]),
				P2PV2Bootstrappers: pq.StringArray{P2Pv2Bootstrapper},
				PluginConfig: map[string]interface{}{
					"dkgEncryptionPublicKey": fmt.Sprintf("\"%s\"", OCR2VRFPluginConfig.DKGConfig.DKGKeyConfigs[index].DKGEncryptionPublicKey),
					"dkgSigningPublicKey":    fmt.Sprintf("\"%s\"", OCR2VRFPluginConfig.DKGConfig.DKGKeyConfigs[index].DKGSigningPublicKey),
					"dkgKeyID":               fmt.Sprintf("\"%s\"", OCR2VRFPluginConfig.DKGConfig.DKGKeyID),
					"dkgContractAddress":     fmt.Sprintf("\"%s\"", OCR2VRFPluginConfig.DKGConfig.DKGContractAddress),
					"vrfCoordinatorAddress":  fmt.Sprintf("\"%s\"", OCR2VRFPluginConfig.VRFCoordinatorAddress),
					"linkEthFeedAddress":     fmt.Sprintf("\"%s\"", OCR2VRFPluginConfig.LinkEthFeedAddress),
				},
			},
		}
		_, err = nonBootstrapNode.MustCreateJob(&OCR2VRFJobSpec)
		require.NoError(t, err, "Shouldn't fail creating OCR Task job on OCR node %d", index)
	}
}

func createBootstrapJob(t *testing.T, bootstrapNode *client.Chainlink, dkgAddress string, chainID int64) string {
	bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
	require.NoError(t, err, "Shouldn't fail reading P2P keys from bootstrap node")
	bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID

	bootstrapSpec := &client.OCR2TaskJobSpec{
		Name:    "ocr2 bootstrap node",
		JobType: "bootstrap",
		OCR2OracleSpec: job.OCR2OracleSpec{
			ContractID: dkgAddress,
			Relay:      "evm",
			RelayConfig: map[string]interface{}{
				"chainID": int(chainID),
			},
		},
	}
	_, err = bootstrapNode.MustCreateJob(bootstrapSpec)
	require.NoError(t, err, "Shouldn't fail creating bootstrap job on bootstrap node")
	return fmt.Sprintf("%s@%s:%d", bootstrapP2PId, bootstrapNode.RemoteIP(), 6690)
}

func BuildOCR2DKGConfigVars(
	t *testing.T,
	ocr2VRFPluginConfig *OCR2VRFPluginConfig,
) contracts.OCRConfig {

	var onchainPublicKeys []common.Address
	for _, onchainPublicKey := range ocr2VRFPluginConfig.OCR2Config.OnchainPublicKeys {
		onchainPublicKeys = append(onchainPublicKeys, common.HexToAddress(onchainPublicKey))
	}

	var transmitters []common.Address
	for _, transmitterAddress := range ocr2VRFPluginConfig.OCR2Config.TransmitterAddresses {
		transmitters = append(transmitters, common.HexToAddress(transmitterAddress))
	}
	oracleIdentities, err := toOraclesIdentityList(
		onchainPublicKeys,
		ocr2VRFPluginConfig.OCR2Config.OffchainPublicKeys,
		ocr2VRFPluginConfig.OCR2Config.ConfigPublicKeys,
		ocr2VRFPluginConfig.OCR2Config.PeerIds,
		ocr2VRFPluginConfig.OCR2Config.TransmitterAddresses,
	)
	require.NoError(t, err)

	ed25519Suite := edwards25519.NewBlakeSHA256Ed25519()
	var signingKeys []kyber.Point
	for _, dkgKey := range ocr2VRFPluginConfig.DKGConfig.DKGKeyConfigs {
		signingKeyBytes, err := hex.DecodeString(dkgKey.DKGSigningPublicKey)
		require.NoError(t, err)
		signingKeyPoint := ed25519Suite.Point()
		err = signingKeyPoint.UnmarshalBinary(signingKeyBytes)
		require.NoError(t, err)
		signingKeys = append(signingKeys, signingKeyPoint)
	}

	altbn128Suite := &altbn_128.PairingSuite{}
	var encryptionKeys []kyber.Point
	for _, dkgKey := range ocr2VRFPluginConfig.DKGConfig.DKGKeyConfigs {
		encryptionKeyBytes, err := hex.DecodeString(dkgKey.DKGEncryptionPublicKey)
		require.NoError(t, err)
		encryptionKeyPoint := altbn128Suite.G1().Point()
		err = encryptionKeyPoint.UnmarshalBinary(encryptionKeyBytes)
		require.NoError(t, err)
		encryptionKeys = append(encryptionKeys, encryptionKeyPoint)
	}

	keyIDBytes, err := DecodeHexTo32ByteArray(ocr2VRFPluginConfig.DKGConfig.DKGKeyID)
	require.NoError(t, err, "Shouldn't fail decoding DKG key ID")

	offchainConfig, err := dkg.OffchainConfig(encryptionKeys, signingKeys, &altbn_128.G1{}, &ocr2vrftypes.PairingTranslation{
		Suite: &altbn_128.PairingSuite{},
	})
	require.NoError(t, err)
	onchainConfig, err := dkg.OnchainConfig(dkg.KeyID(keyIDBytes))
	require.NoError(t, err)

	_, _, f, onchainConfig, offchainConfigVersion, offchainConfig, err :=
		confighelper.ContractSetConfigArgsForTests(
			30*time.Second,                          // deltaProgress time.Duration,
			10*time.Second,                          // deltaResend time.Duration,
			10*time.Millisecond,                     // deltaRound time.Duration,
			20*time.Millisecond,                     // deltaGrace time.Duration,
			20*time.Millisecond,                     // deltaStage time.Duration,
			3,                                       // rMax uint8,
			ocr2VRFPluginConfig.OCR2Config.Schedule, // s []int,
			oracleIdentities,                        // oracles []OracleIdentityExtra,
			offchainConfig,
			10*time.Millisecond, // maxDurationQuery time.Duration,
			10*time.Second,      // maxDurationObservation time.Duration,
			10*time.Second,      // maxDurationReport time.Duration,
			10*time.Millisecond, // maxDurationShouldAcceptFinalizedReport time.Duration,
			1*time.Second,       // maxDurationShouldTransmitAcceptedReport time.Duration,
			1,                   // f int,
			onchainConfig,       // onchainConfig []byte,
		)
	require.NoError(t, err, "Shouldn't fail building OCR config")

	log.Info().Msg("Done building DKG OCR config")
	return contracts.OCRConfig{
		Signers:               onchainPublicKeys,
		Transmitters:          transmitters,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}
}

func toOraclesIdentityList(onchainPubKeys []common.Address, offchainPubKeys, configPubKeys, peerIDs, transmitters []string) ([]confighelper.OracleIdentityExtra, error) {
	offchainPubKeysBytes := []types.OffchainPublicKey{}
	for _, pkHex := range offchainPubKeys {
		pkBytes, err := hex.DecodeString(pkHex)
		if err != nil {
			return nil, err
		}
		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		n := copy(pkBytesFixed[:], pkBytes)
		if n != ed25519.PublicKeySize {
			return nil, errors.New("wrong num elements copied")
		}

		offchainPubKeysBytes = append(offchainPubKeysBytes, types.OffchainPublicKey(pkBytesFixed))
	}

	configPubKeysBytes := []types.ConfigEncryptionPublicKey{}
	for _, pkHex := range configPubKeys {
		pkBytes, err := hex.DecodeString(pkHex)
		if err != nil {
			return nil, err
		}
		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		n := copy(pkBytesFixed[:], pkBytes)
		if n != ed25519.PublicKeySize {
			return nil, errors.New("wrong num elements copied")
		}

		configPubKeysBytes = append(configPubKeysBytes, types.ConfigEncryptionPublicKey(pkBytesFixed))
	}

	o := []confighelper.OracleIdentityExtra{}
	for index := range configPubKeys {
		o = append(o, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  onchainPubKeys[index][:],
				OffchainPublicKey: offchainPubKeysBytes[index],
				PeerID:            peerIDs[index],
				TransmitAccount:   types.Account(transmitters[index]),
			},
			ConfigEncryptionPublicKey: configPubKeysBytes[index],
		})
	}
	return o, nil
}

func DecodeHexTo32ByteArray(val string) ([32]byte, error) {
	var byteArray [32]byte
	decoded, err := hex.DecodeString(val)
	if err != nil {
		return [32]byte{}, err
	}
	if len(decoded) != 32 {
		return [32]byte{}, fmt.Errorf("expected value to be 32 bytes but received %d bytes", len(decoded))
	}
	copy(byteArray[:], decoded)
	return byteArray, err
}

func BuildOCR2VRFConfigVars(
	t *testing.T,
	ocr2VRFPluginConfig *OCR2VRFPluginConfig,
) contracts.OCRConfig {

	var onchainPublicKeys []common.Address
	for _, onchainPublicKey := range ocr2VRFPluginConfig.OCR2Config.OnchainPublicKeys {
		onchainPublicKeys = append(onchainPublicKeys, common.HexToAddress(onchainPublicKey))
	}

	var transmitters []common.Address
	for _, transmitterAddress := range ocr2VRFPluginConfig.OCR2Config.TransmitterAddresses {
		transmitters = append(transmitters, common.HexToAddress(transmitterAddress))
	}

	oracleIdentities, err := toOraclesIdentityList(
		onchainPublicKeys,
		ocr2VRFPluginConfig.OCR2Config.OffchainPublicKeys,
		ocr2VRFPluginConfig.OCR2Config.ConfigPublicKeys,
		ocr2VRFPluginConfig.OCR2Config.PeerIds,
		ocr2VRFPluginConfig.OCR2Config.TransmitterAddresses,
	)
	require.NoError(t, err)

	confDelays := make(map[uint32]struct{})

	for _, c := range ocr2VRFPluginConfig.VRFBeaconConfig.ConfDelays {
		confDelay, err := strconv.ParseUint(c, 0, 32)
		require.NoError(t, err)
		confDelays[uint32(confDelay)] = struct{}{}
	}

	onchainConfig := ocr2vrf.OnchainConfig(confDelays)

	_, _, f, onchainConfig, offchainConfigVersion, offchainConfig, err :=
		confighelper.ContractSetConfigArgsForTests(
			30*time.Second,                          // deltaProgress time.Duration,
			10*time.Second,                          // deltaResend time.Duration,
			10*time.Second,                          // deltaRound time.Duration,
			20*time.Second,                          // deltaGrace time.Duration,
			20*time.Second,                          // deltaStage time.Duration,
			3,                                       // rMax uint8,
			ocr2VRFPluginConfig.OCR2Config.Schedule, // s []int,
			oracleIdentities,                        // oracles []OracleIdentityExtra,
			ocr2vrf.OffchainConfig(ocr2VRFPluginConfig.VRFBeaconConfig.CoordinatorConfig),
			10*time.Millisecond, // maxDurationQuery time.Duration,
			10*time.Second,      // maxDurationObservation time.Duration,
			10*time.Second,      // maxDurationReport time.Duration,
			5*time.Second,       // maxDurationShouldAcceptFinalizedReport time.Duration,
			1*time.Second,       // maxDurationShouldTransmitAcceptedReport time.Duration,
			1,                   // f int,
			onchainConfig,       // onchainConfig []byte,
		)
	require.NoError(t, err)

	log.Info().Msg("Done building VRF OCR config")
	return contracts.OCRConfig{
		Signers:               onchainPublicKeys,
		Transmitters:          transmitters,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}

}
