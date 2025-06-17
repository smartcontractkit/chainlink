// TODO: KS-458 copied from https://github.com/smartcontractkit/chainlink/blob/65924811dc53a211613927c814d7f04fd85439a4/core/scripts/keystone/src/88_gen_ocr3_config.go#L1
// to unblock go mod issues when trying to import the scripts package
package internal

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	capocr3types "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"

	kocr3 "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

type TopLevelConfigSource struct {
	OracleConfig OracleConfig
}

type OracleConfig struct {
	MaxQueryLengthBytes       uint32
	MaxObservationLengthBytes uint32
	MaxReportLengthBytes      uint32
	MaxOutcomeLengthBytes     uint32
	MaxReportCount            uint32
	MaxBatchSize              uint32
	OutcomePruningThreshold   uint64
	UniqueReports             bool
	RequestTimeout            time.Duration

	DeltaProgressMillis               uint32
	DeltaResendMillis                 uint32
	DeltaInitialMillis                uint32
	DeltaRoundMillis                  uint32
	DeltaGraceMillis                  uint32
	DeltaCertifiedCommitRequestMillis uint32
	DeltaStageMillis                  uint32
	MaxRoundsPerEpoch                 uint64
	TransmissionSchedule              []int

	MaxDurationQueryMillis          uint32
	MaxDurationObservationMillis    uint32
	MaxDurationShouldAcceptMillis   uint32
	MaxDurationShouldTransmitMillis uint32

	MaxFaultyOracles int
}

func (oc *OracleConfig) UnmarshalJSON(data []byte) error {
	type aliasT OracleConfig
	temp := &struct {
		RequestTimeout string `json:"RequestTimeout"`
		*aliasT
	}{
		aliasT: (*aliasT)(oc),
	}
	if err := json.Unmarshal(data, temp); err != nil {
		return fmt.Errorf("failed to unmarshal OracleConfig: %w", err)
	}

	if temp.RequestTimeout == "" {
		oc.RequestTimeout = 0
	} else {
		requestTimeout, err := time.ParseDuration(temp.RequestTimeout)
		if err != nil {
			return fmt.Errorf("failed to parse RequestTimeout: %w", err)
		}
		oc.RequestTimeout = requestTimeout
	}

	return nil
}

func (oc OracleConfig) MarshalJSON() ([]byte, error) {
	type aliasT OracleConfig
	return json.Marshal(&struct {
		RequestTimeout string `json:"RequestTimeout"`
		*aliasT
	}{
		RequestTimeout: oc.RequestTimeout.String(),
		aliasT:         (*aliasT)(&oc),
	})
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
	EncryptionPublicKey   string `json:"EncryptionPublicKey"`
}

// OCR2OracleConfig is the input configuration for an OCR2/3 contract.
type OCR2OracleConfig struct {
	Signers               [][]byte
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

func (c OCR2OracleConfig) MarshalJSON() ([]byte, error) {
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

func (c *OCR2OracleConfig) UnmarshalJSON(data []byte) error {
	type aliasT struct {
		Signers               []string
		Transmitters          []string
		F                     uint8
		OnchainConfig         string
		OffchainConfigVersion uint64
		OffchainConfig        string
	}
	var alias aliasT
	err := json.Unmarshal(data, &alias)
	if err != nil {
		return fmt.Errorf("failed to unmarshal OCR2OracleConfig alias: %w", err)
	}
	c.F = alias.F
	c.OffchainConfigVersion = alias.OffchainConfigVersion
	c.Signers = make([][]byte, len(alias.Signers))
	for i, signer := range alias.Signers {
		c.Signers[i], err = hex.DecodeString(strings.TrimPrefix(signer, "0x"))
		if err != nil {
			return fmt.Errorf("failed to decode signer: %w", err)
		}
	}
	c.Transmitters = make([]common.Address, len(alias.Transmitters))
	for i, transmitter := range alias.Transmitters {
		c.Transmitters[i] = common.HexToAddress(transmitter)
	}
	c.OnchainConfig, err = hex.DecodeString(strings.TrimPrefix(alias.OnchainConfig, "0x"))
	if err != nil {
		return fmt.Errorf("failed to decode onchain config: %w", err)
	}
	c.OffchainConfig, err = hex.DecodeString(strings.TrimPrefix(alias.OffchainConfig, "0x"))
	if err != nil {
		return fmt.Errorf("failed to decode offchain config: %w", err)
	}
	return nil
}

func GenerateOCR3Config(cfg OracleConfig, nca []NodeKeys, secrets cldf.OCRSecrets) (OCR2OracleConfig, error) {
	// the transmission schedule is very specific; arguably it should be not be a parameter
	if len(cfg.TransmissionSchedule) != 1 || cfg.TransmissionSchedule[0] != len(nca) {
		return OCR2OracleConfig{}, fmt.Errorf("transmission schedule must have exactly one entry, matching the len of the number of nodes want [%d], got %v", len(nca), cfg.TransmissionSchedule)
	}
	onchainPubKeys := [][]byte{}
	allPubKeys := map[string]any{}
	if secrets.IsEmpty() {
		return OCR2OracleConfig{}, errors.New("OCRSecrets is required")
	}
	for _, n := range nca {
		// evm keys always required
		if n.OCR2OnchainPublicKey == "" {
			return OCR2OracleConfig{}, errors.New("OCR2OnchainPublicKey is required")
		}
		ethPubKey := common.HexToAddress(n.OCR2OnchainPublicKey)
		pubKeys := map[string]types.OnchainPublicKey{
			string(chaintype.EVM): ethPubKey.Bytes(),
		}
		// add aptos key if present
		if n.AptosOnchainPublicKey != "" {
			aptosPubKey, err := hex.DecodeString(n.AptosOnchainPublicKey)
			if err != nil {
				return OCR2OracleConfig{}, fmt.Errorf("failed to decode AptosOnchainPublicKey: %w", err)
			}
			pubKeys[string(chaintype.Aptos)] = aptosPubKey
		}
		// validate uniqueness of each individual key
		for _, key := range pubKeys {
			raw := hex.EncodeToString(key)
			_, exists := allPubKeys[raw]
			if exists {
				return OCR2OracleConfig{}, fmt.Errorf("Duplicate onchain public key: '%s'", raw)
			}
			allPubKeys[raw] = struct{}{}
		}
		pubKey, err := ocrcommon.MarshalMultichainPublicKey(pubKeys)
		if err != nil {
			return OCR2OracleConfig{}, fmt.Errorf("failed to marshal multichain public key: %w", err)
		}
		onchainPubKeys = append(onchainPubKeys, pubKey)
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
				OnchainPublicKey:  onchainPubKeys[index],
				OffchainPublicKey: offchainPubKeysBytes[index],
				PeerID:            nca[index].P2PPeerID,
				TransmitAccount:   types.Account(nca[index].EthAddress),
			},
			ConfigEncryptionPublicKey: configPubKeysBytes[index],
		})
	}

	// let's keep reqTimeout as nil if it's 0, so we can use the default value within `chainlink-common`.
	// See: https://github.com/smartcontractkit/chainlink-common/blob/main/pkg/capabilities/consensus/ocr3/factory.go#L73
	var reqTimeout *durationpb.Duration
	if cfg.RequestTimeout > 0 {
		reqTimeout = durationpb.New(cfg.RequestTimeout)
	}
	cfgBytes, err := proto.Marshal(&capocr3types.ReportingPluginConfig{
		MaxQueryLengthBytes:       cfg.MaxQueryLengthBytes,
		MaxObservationLengthBytes: cfg.MaxObservationLengthBytes,
		MaxReportLengthBytes:      cfg.MaxReportLengthBytes,
		MaxOutcomeLengthBytes:     cfg.MaxOutcomeLengthBytes,
		MaxReportCount:            cfg.MaxReportCount,
		MaxBatchSize:              cfg.MaxBatchSize,
		OutcomePruningThreshold:   cfg.OutcomePruningThreshold,
		RequestTimeout:            reqTimeout,
	})
	if err != nil {
		return OCR2OracleConfig{}, fmt.Errorf("failed to marshal ReportingPluginConfig: %w", err)
	}

	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsDeterministic(
		secrets.EphemeralSk,
		secrets.SharedSecret,
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
		cfgBytes, // reportingPluginConfig
		nil,      // maxDurationInitialization
		time.Duration(cfg.MaxDurationQueryMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationObservationMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationShouldAcceptMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationShouldTransmitMillis)*time.Millisecond,
		cfg.MaxFaultyOracles,
		nil, // empty onChain config
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

type configureOCR3Request struct {
	cfg        *OracleConfig
	chain      cldf_evm.Chain
	contract   *kocr3.OCR3Capability
	nodes      []deployment.Node
	dryRun     bool
	ocrSecrets cldf.OCRSecrets

	useMCMS bool
}

func (r configureOCR3Request) generateOCR3Config() (OCR2OracleConfig, error) {
	nks := makeNodeKeysSlice(r.nodes, r.chain.Selector)
	return GenerateOCR3Config(*r.cfg, nks, r.ocrSecrets)
}

type configureOCR3Response struct {
	ocrConfig OCR2OracleConfig
	ops       *mcmstypes.BatchOperation
}

func configureOCR3contract(req configureOCR3Request) (*configureOCR3Response, error) {
	if req.contract == nil {
		return nil, errors.New("OCR3 contract is nil")
	}
	ocrConfig, err := req.generateOCR3Config()
	if err != nil {
		return nil, fmt.Errorf("failed to generate OCR3 config: %w", err)
	}
	if req.dryRun {
		return &configureOCR3Response{ocrConfig, nil}, nil
	}

	txOpt := req.chain.DeployerKey
	if req.useMCMS {
		txOpt = cldf.SimTransactOpts()
	}

	tx, err := req.contract.SetConfig(txOpt,
		ocrConfig.Signers,
		ocrConfig.Transmitters,
		ocrConfig.F,
		ocrConfig.OnchainConfig,
		ocrConfig.OffchainConfigVersion,
		ocrConfig.OffchainConfig,
	)
	if err != nil {
		err = cldf.DecodeErr(kocr3.OCR3CapabilityABI, err)
		return nil, fmt.Errorf("failed to call SetConfig for OCR3 contract %s using mcms: %T: %w", req.contract.Address().String(), req.useMCMS, err)
	}

	var ops mcmstypes.BatchOperation
	if !req.useMCMS {
		_, err = req.chain.Confirm(tx)
		if err != nil {
			err = cldf.DecodeErr(kocr3.OCR3CapabilityABI, err)
			return nil, fmt.Errorf("failed to confirm SetConfig for OCR3 contract %s: %w", req.contract.Address().String(), err)
		}
	} else {
		ops, err = proposalutils.BatchOperationForChain(req.chain.Selector, req.contract.Address().Hex(), tx.Data(), big.NewInt(0), string(OCR3Capability), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create batch operation: %w", err)
		}
	}

	return &configureOCR3Response{ocrConfig, &ops}, nil
}
