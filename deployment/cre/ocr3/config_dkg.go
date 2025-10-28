package ocr3

import (
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"

	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

type V3_1OracleConfig struct {
	DeltaProgressMillis  uint32
	DeltaRoundMillis     uint32
	DeltaGraceMillis     uint32
	DeltaStageMillis     uint32
	MaxRoundsPerEpoch    uint64
	TransmissionSchedule []int

	MaxDurationInitializationMillis               uint32
	MaxDurationShouldAcceptAttestedReportMillis   uint32
	MaxDurationShouldTransmitAcceptedReportMillis uint32

	WarnDurationQueryMillis               uint32
	WarnDurationObservationMillis         uint32
	WarnDurationValidateObservationMillis uint32
	WarnDurationObservationQuorumMillis   uint32
	WarnDurationStateTransition           uint32
	WarnDurationCommitted                 uint32

	MaxFaultyOracles int
}

const offchainPublicKeyType byte = 0x8

func oCR3CapabilityCompatibleOnchainPublicKey(offchainPublicKey types.OffchainPublicKey) types.OnchainPublicKey {
	result := make([]byte, 0, 1+2+len(offchainPublicKey))
	result = append(result, offchainPublicKeyType)
	result = binary.LittleEndian.AppendUint16(result, uint16(len(offchainPublicKey)))
	result = append(result, offchainPublicKey[:]...)

	return result
}

func GenerateDKGConfigFromNodes(cfg V3_1OracleConfig, nodes []deployment.Node, registryChainSel uint64, secrets ocr.OCRSecrets, dkgCfg dkgocrtypes.ReportingPluginConfig) (OCR2OracleConfig, error) {
	nca := makeNodeKeysSlice(nodes, registryChainSel)
	return GenerateDKGConfig(cfg, nca, secrets, dkgCfg)
}

func GenerateDKGConfig(cfg V3_1OracleConfig, nca []NodeKeys, secrets ocr.OCRSecrets, dkgCfg dkgocrtypes.ReportingPluginConfig) (OCR2OracleConfig, error) {
	// the transmission schedule is very specific; arguably it should be not be a parameter
	if len(cfg.TransmissionSchedule) != 1 || cfg.TransmissionSchedule[0] != len(nca) {
		return OCR2OracleConfig{}, fmt.Errorf("transmission schedule must have exactly one entry, matching the len of the number of nodes want [%d], got %v. Total TransmissionSchedules = %d", len(nca), cfg.TransmissionSchedule, len(cfg.TransmissionSchedule))
	}

	offchainPubKeysBytes := []types.OffchainPublicKey{}
	for _, n := range nca {
		pkBytes, err := hex.DecodeString(n.OCR2OffchainPublicKey)
		if err != nil {
			return OCR2OracleConfig{}, fmt.Errorf("failed to decode OCR2OffchainPublicKey: %w", err)
		}

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		nCopied := copy(pkBytesFixed[:], pkBytes)
		if nCopied != ed25519.PublicKeySize {
			return OCR2OracleConfig{}, fmt.Errorf("wrong num elements copied from ocr2 offchain public key. expected %d but got %d", ed25519.PublicKeySize, nCopied)
		}

		offchainPubKeysBytes = append(offchainPubKeysBytes, pkBytesFixed)
	}

	onChainPublicKeys := make([]types.OnchainPublicKey, 0, len(offchainPubKeysBytes))
	for _, pk := range offchainPubKeysBytes {
		onChainPublicKeys = append(onChainPublicKeys, oCR3CapabilityCompatibleOnchainPublicKey(pk))
	}

	configPubKeysBytes := []types.ConfigEncryptionPublicKey{}
	for _, n := range nca {
		pkBytes, err := hex.DecodeString(n.OCR2ConfigPublicKey)
		if err != nil {
			return OCR2OracleConfig{}, fmt.Errorf("failed to decode OCR2ConfigPublicKey: %w", err)
		}

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		n := copy(pkBytesFixed[:], pkBytes)
		if n != ed25519.PublicKeySize {
			return OCR2OracleConfig{}, fmt.Errorf("wrong num elements copied from ocr2 config public key. expected %d but got %d", ed25519.PublicKeySize, n)
		}

		configPubKeysBytes = append(configPubKeysBytes, pkBytesFixed)
	}

	identities := []confighelper.OracleIdentityExtra{}
	for index := range nca {
		identities = append(identities, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  onChainPublicKeys[index],
				OffchainPublicKey: offchainPubKeysBytes[index],
				PeerID:            nca[index].P2PPeerID,
				TransmitAccount:   types.Account(common.HexToAddress(fmt.Sprintf("0xc1c1c1c1%x", offchainPubKeysBytes[index][:16])).Hex()),
			},
			ConfigEncryptionPublicKey: configPubKeysBytes[index],
		})
	}

	cfgBytes, err := dkgCfg.MarshalBinary()
	if err != nil {
		return OCR2OracleConfig{}, fmt.Errorf("failed to marshal ReportingPluginConfig: %w", err)
	}

	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3_1confighelper.ContractSetConfigArgsDeterministic(
		ocr3_1confighelper.CheckPublicConfigLevelDefault,
		secrets.EphemeralSk,
		secrets.SharedSecret,
		identities,
		cfg.MaxFaultyOracles,
		time.Duration(cfg.DeltaProgressMillis)*time.Millisecond,
		time.Duration(cfg.DeltaRoundMillis)*time.Millisecond,
		time.Duration(cfg.DeltaGraceMillis)*time.Millisecond,
		cfg.MaxRoundsPerEpoch,
		time.Duration(cfg.DeltaStageMillis)*time.Millisecond,
		cfg.TransmissionSchedule,
		cfgBytes,
		nil, // onchainConfig
		time.Duration(cfg.MaxDurationInitializationMillis)*time.Millisecond,
		time.Duration(cfg.WarnDurationQueryMillis)*time.Millisecond,
		time.Duration(cfg.WarnDurationObservationMillis)*time.Millisecond,
		time.Duration(cfg.WarnDurationValidateObservationMillis)*time.Millisecond,
		time.Duration(cfg.WarnDurationObservationQuorumMillis)*time.Millisecond,
		time.Duration(cfg.WarnDurationStateTransition)*time.Millisecond,
		time.Duration(cfg.WarnDurationCommitted)*time.Millisecond,
		time.Duration(cfg.MaxDurationShouldAcceptAttestedReportMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationShouldTransmitAcceptedReportMillis)*time.Millisecond,
		ocr3_1confighelper.ContractSetConfigArgsOptionalConfig{},
	)
	if err != nil {
		return OCR2OracleConfig{}, fmt.Errorf("failed to generate contract config args: %w", err)
	}

	var configSigners [][]byte
	for _, signer := range signers {
		configSigners = append(configSigners, signer)
	}

	transmitterAddresses, err := evm.AccountToAddress(transmitters)
	if err != nil {
		return OCR2OracleConfig{}, fmt.Errorf("failed to convert transmitters to addresses: %w", err)
	}

	config := OCR2OracleConfig{
		Signers:               configSigners,
		Transmitters:          transmitterAddresses,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}

	return config, nil
}
