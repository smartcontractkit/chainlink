package ocr3_1

import (
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

const offchainPublicKeyType byte = 0x8

func oCR3CapabilityCompatibleOnchainPublicKey(offchainPublicKey types.OffchainPublicKey) types.OnchainPublicKey {
	result := make([]byte, 0, 1+2+len(offchainPublicKey))
	result = append(result, offchainPublicKeyType)
	result = binary.LittleEndian.AppendUint16(result, uint16(len(offchainPublicKey)))
	result = append(result, offchainPublicKey[:]...)

	return result
}

func dkgIdentityKeys(nca []ocr3.NodeKeys) ([]types.OnchainPublicKey, []types.Account, error) {
	onchainKeys := make([]types.OnchainPublicKey, 0, len(nca))
	transmitAccounts := make([]types.Account, 0, len(nca))
	for _, n := range nca {
		pkBytes, err := hex.DecodeString(n.OCR2OffchainPublicKey)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decode OCR2OffchainPublicKey: %w", err)
		}

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		nCopied := copy(pkBytesFixed[:], pkBytes)
		if nCopied != ed25519.PublicKeySize {
			return nil, nil, fmt.Errorf("wrong num elements copied from ocr2 offchain public key. expected %d but got %d", ed25519.PublicKeySize, nCopied)
		}

		onchainKeys = append(onchainKeys, oCR3CapabilityCompatibleOnchainPublicKey(pkBytesFixed))
		transmitAccounts = append(transmitAccounts, types.Account(common.HexToAddress(fmt.Sprintf("0xc1c1c1c1%x", pkBytesFixed[:16])).Hex()))
	}
	return onchainKeys, transmitAccounts, nil
}

func GenerateDKGConfigFromNodes(cfg V3_1OracleConfig, nodes []deployment.Node, registryChainSel uint64, secrets ocr.OCRSecrets, extraSignerFamilies []string) (ocr3.OCR2OracleConfig, error) {
	nca := ocr3.MakeNodeKeysSlice(nodes, registryChainSel, extraSignerFamilies)
	if cfg.DKGOffchainConfig == nil {
		return ocr3.OCR2OracleConfig{}, errors.New("cannot provide empty DKGOffchainConfig")
	}
	cfgBytes, err := cfg.DKGOffchainConfig.Marshal()
	if err != nil {
		return ocr3.OCR2OracleConfig{}, fmt.Errorf("failed to marshal ReportingPluginConfig: %w", err)
	}
	onchainKeys, transmitAccounts, err := dkgIdentityKeys(nca)
	if err != nil {
		return ocr3.OCR2OracleConfig{}, err
	}
	return GenerateDKGConfig(cfg, nca, secrets, cfgBytes, onchainKeys, transmitAccounts)
}

func GenerateDKGConfig(cfg V3_1OracleConfig, nca []ocr3.NodeKeys, secrets ocr.OCRSecrets, cfgBytes []byte, onchainKeys []types.OnchainPublicKey, transmitAccounts []types.Account) (ocr3.OCR2OracleConfig, error) {
	// the transmission schedule is very specific; arguably it should be not be a parameter
	if len(cfg.TransmissionSchedule) != 1 || cfg.TransmissionSchedule[0] != len(nca) {
		return ocr3.OCR2OracleConfig{}, fmt.Errorf("transmission schedule must have exactly one entry, matching the len of the number of nodes want [%d], got %v. Total TransmissionSchedules = %d", len(nca), cfg.TransmissionSchedule, len(cfg.TransmissionSchedule))
	}

	offchainPubKeysBytes := []types.OffchainPublicKey{}
	for _, n := range nca {
		pkBytes, err := hex.DecodeString(n.OCR2OffchainPublicKey)
		if err != nil {
			return ocr3.OCR2OracleConfig{}, fmt.Errorf("failed to decode OCR2OffchainPublicKey: %w", err)
		}

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		nCopied := copy(pkBytesFixed[:], pkBytes)
		if nCopied != ed25519.PublicKeySize {
			return ocr3.OCR2OracleConfig{}, fmt.Errorf("wrong num elements copied from ocr2 offchain public key. expected %d but got %d", ed25519.PublicKeySize, nCopied)
		}

		offchainPubKeysBytes = append(offchainPubKeysBytes, pkBytesFixed)
	}

	configPubKeysBytes := []types.ConfigEncryptionPublicKey{}
	for _, n := range nca {
		pkBytes, err := hex.DecodeString(n.OCR2ConfigPublicKey)
		if err != nil {
			return ocr3.OCR2OracleConfig{}, fmt.Errorf("failed to decode OCR2ConfigPublicKey: %w", err)
		}

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		n := copy(pkBytesFixed[:], pkBytes)
		if n != ed25519.PublicKeySize {
			return ocr3.OCR2OracleConfig{}, fmt.Errorf("wrong num elements copied from ocr2 config public key. expected %d but got %d", ed25519.PublicKeySize, n)
		}

		configPubKeysBytes = append(configPubKeysBytes, pkBytesFixed)
	}

	identities := make([]confighelper.OracleIdentityExtra, 0, len(nca))
	for i := range nca {
		identities = append(identities, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  onchainKeys[i],
				OffchainPublicKey: offchainPubKeysBytes[i],
				PeerID:            nca[i].P2PPeerID,
				TransmitAccount:   transmitAccounts[i],
			},
			ConfigEncryptionPublicKey: configPubKeysBytes[i],
		})
	}

	return genOCR3_1Config(cfg, identities, secrets, cfgBytes)
}
