// TODO: KS-458 copied from https://github.com/smartcontractkit/chainlink/blob/65924811dc53a211613927c814d7f04fd85439a4/core/scripts/keystone/src/88_gen_ocr3_config.go#L1
// to unblock go mod issues when trying to import the scripts package
package keystone

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

type TopLevelConfigSource struct {
	OracleConfig OracleConfigSource
}

type OracleConfigSource struct {
	MaxQueryLengthBytes       uint32
	MaxObservationLengthBytes uint32
	MaxReportLengthBytes      uint32
	MaxRequestBatchSize       uint32
	UniqueReports             bool

	DeltaProgressMillis               uint32
	DeltaResendMillis                 uint32
	DeltaInitialMillis                uint32
	DeltaRoundMillis                  uint32
	DeltaGraceMillis                  uint32
	DeltaCertifiedCommitRequestMillis uint32
	DeltaStageMillis                  uint32
	MaxRoundsPerEpoch                 uint64
	TransmissionSchedule              []int

	MaxDurationQueryMillis       uint32
	MaxDurationObservationMillis uint32
	MaxDurationAcceptMillis      uint32
	MaxDurationTransmitMillis    uint32

	MaxFaultyOracles int
}

type NodeKeys struct {
	EthAddress            string `json:"EthAddress"`
	AptosAccount          string `json:"AptosAccount"`
	AptosBundleID         string `json:"AptosBundleID"`
	AptosOnchainPublicKey string `json:"AptosOnchainPublicKey"`
	P2PPeerID             string `json:"P2PPeerID"`             // p2p_<key>
	OCR2BundleID          string `json:"OCR2BundleID"`          // used only in job spec
	OCR2OnchainPublicKey  string `json:"OCR2OnchainPublicKey"`  // ocr2on_evm_<key>
	OCR2OffchainPublicKey string `json:"OCR2OffchainPublicKey"` // ocr2off_evm_<key>
	OCR2ConfigPublicKey   string `json:"OCR2ConfigPublicKey"`   // ocr2cfg_evm_<key>
	CSAPublicKey          string `json:"CSAPublicKey"`
}

type Orc2drOracleConfig struct {
	Signers               [][]byte
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

func (c Orc2drOracleConfig) MarshalJSON() ([]byte, error) {
	alias := struct {
		Signers               []string
		Transmitters          []string
		F                     uint8
		OnchainConfig         string
		OffchainConfigVersion uint64
		OffchainConfig        string
	}{
		Signers:               make([]string, len(c.Signers)),
		Transmitters:          make([]string, len(c.Transmitters)),
		F:                     c.F,
		OnchainConfig:         "0x" + hex.EncodeToString(c.OnchainConfig),
		OffchainConfigVersion: c.OffchainConfigVersion,
		OffchainConfig:        "0x" + hex.EncodeToString(c.OffchainConfig),
	}

	for i, signer := range c.Signers {
		alias.Signers[i] = hex.EncodeToString(signer)
	}

	for i, transmitter := range c.Transmitters {
		alias.Transmitters[i] = transmitter.Hex()
	}

	return json.Marshal(alias)
}

func GenerateOCR3Config(cfg OracleConfigSource, nca []NodeKeys) (Orc2drOracleConfig, error) {
	onchainPubKeys := [][]byte{}
	allPubKeys := map[string]any{}
	for _, n := range nca {
		// evm keys always required
		if n.OCR2OnchainPublicKey == "" {
			return Orc2drOracleConfig{}, errors.New("OCR2OnchainPublicKey is required")
		}
		ethPubKey := common.HexToAddress(n.OCR2OnchainPublicKey)
		pubKeys := map[string]types.OnchainPublicKey{
			string(chaintype.EVM): ethPubKey.Bytes(),
		}
		// add aptos key if present
		if n.AptosOnchainPublicKey != "" {
			aptosPubKey, err := hex.DecodeString(n.AptosOnchainPublicKey)
			if err != nil {
				return Orc2drOracleConfig{}, fmt.Errorf("failed to decode AptosOnchainPublicKey: %w", err)
			}
			pubKeys[string(chaintype.Aptos)] = aptosPubKey
		}
		// validate uniqueness of each individual key
		for _, key := range pubKeys {
			raw := hex.EncodeToString(key)
			_, exists := allPubKeys[raw]
			if exists {
				return Orc2drOracleConfig{}, fmt.Errorf("Duplicate onchain public key: '%s'", raw)
			}
			allPubKeys[raw] = struct{}{}
		}
		pubKey, err := ocrcommon.MarshalMultichainPublicKey(pubKeys)
		if err != nil {
			return Orc2drOracleConfig{}, fmt.Errorf("failed to marshal multichain public key: %w", err)
		}
		onchainPubKeys = append(onchainPubKeys, pubKey)
	}

	offchainPubKeysBytes := []types.OffchainPublicKey{}
	for _, n := range nca {
		pkBytes, err := hex.DecodeString(n.OCR2OffchainPublicKey)
		if err != nil {
			return Orc2drOracleConfig{}, fmt.Errorf("failed to decode OCR2OffchainPublicKey: %w", err)
		}

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		nCopied := copy(pkBytesFixed[:], pkBytes)
		if nCopied != ed25519.PublicKeySize {
			return Orc2drOracleConfig{}, fmt.Errorf("wrong num elements copied from ocr2 offchain public key. expected %d but got %d", ed25519.PublicKeySize, nCopied)
		}

		offchainPubKeysBytes = append(offchainPubKeysBytes, types.OffchainPublicKey(pkBytesFixed))
	}

	configPubKeysBytes := []types.ConfigEncryptionPublicKey{}
	for _, n := range nca {
		pkBytes, err := hex.DecodeString(n.OCR2ConfigPublicKey)
		if err != nil {
			return Orc2drOracleConfig{}, fmt.Errorf("failed to decode OCR2ConfigPublicKey: %w", err)
		}

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		n := copy(pkBytesFixed[:], pkBytes)
		if n != ed25519.PublicKeySize {
			return Orc2drOracleConfig{}, fmt.Errorf("wrong num elements copied from ocr2 config public key. expected %d but got %d", ed25519.PublicKeySize, n)
		}

		configPubKeysBytes = append(configPubKeysBytes, types.ConfigEncryptionPublicKey(pkBytesFixed))
	}

	identities := []confighelper.OracleIdentityExtra{}
	for index := range nca {
		identities = append(identities, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  onchainPubKeys[index][:],
				OffchainPublicKey: offchainPubKeysBytes[index],
				PeerID:            nca[index].P2PPeerID,
				TransmitAccount:   types.Account(nca[index].EthAddress),
			},
			ConfigEncryptionPublicKey: configPubKeysBytes[index],
		})
	}

	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		time.Duration(cfg.DeltaProgressMillis)*time.Millisecond,
		time.Duration(cfg.DeltaResendMillis)*time.Millisecond,
		time.Duration(cfg.DeltaInitialMillis)*time.Millisecond,
		time.Duration(cfg.DeltaRoundMillis)*time.Millisecond,
		time.Duration(cfg.DeltaGraceMillis)*time.Millisecond,
		time.Duration(cfg.DeltaCertifiedCommitRequestMillis)*time.Millisecond,
		time.Duration(cfg.DeltaStageMillis)*time.Millisecond,
		cfg.MaxRoundsPerEpoch,
		cfg.TransmissionSchedule,
		identities,
		nil, // reportingPluginConfig
		nil, // maxDurationInitialization
		time.Duration(cfg.MaxDurationQueryMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationObservationMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationAcceptMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationTransmitMillis)*time.Millisecond,
		cfg.MaxFaultyOracles,
		nil, // empty onChain config
	)
	if err != nil {
		return Orc2drOracleConfig{}, fmt.Errorf("failed to generate contract config args: %w", err)
	}

	var configSigners [][]byte
	for _, signer := range signers {
		configSigners = append(configSigners, signer)
	}

	transmitterAddresses, err := evm.AccountToAddress(transmitters)
	if err != nil {
		return Orc2drOracleConfig{}, fmt.Errorf("failed to convert transmitters to addresses: %w", err)
	}

	config := Orc2drOracleConfig{
		Signers:               configSigners,
		Transmitters:          transmitterAddresses,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}

	return config, nil
}
