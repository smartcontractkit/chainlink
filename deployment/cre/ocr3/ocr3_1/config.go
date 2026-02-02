package ocr3_1

import (
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1confighelper"

	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"

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

	PrevConfigDigest  string
	PrevSeqNr         uint64
	PrevHistoryDigest string
}

func GenerateOCR3_1ConfigFromNodes(cfg V3_1OracleConfig, nodes []deployment.Node, registryChainSel uint64, secrets focr.OCRSecrets, reportingPluginConfigOverride []byte) (ocr3.OCR2OracleConfig, error) {
	nca := ocr3.MakeNodeKeysSlice(nodes, registryChainSel)
	return GenerateOCR3_1Config(cfg, nca, secrets, reportingPluginConfigOverride)
}

func GenerateOCR3_1Config(cfg V3_1OracleConfig, nca []ocr3.NodeKeys, secrets focr.OCRSecrets, reportingPluginConfigOverride []byte) (ocr3.OCR2OracleConfig, error) {
	// the transmission schedule is very specific; arguably it should be not be a parameter
	if len(cfg.TransmissionSchedule) != 1 || cfg.TransmissionSchedule[0] != len(nca) {
		return ocr3.OCR2OracleConfig{}, fmt.Errorf("transmission schedule must have exactly one entry, matching the len of the number of nodes want [%d], got %v. Total TransmissionSchedules = %d", len(nca), cfg.TransmissionSchedule, len(cfg.TransmissionSchedule))
	}

	if secrets.IsEmpty() {
		return ocr3.OCR2OracleConfig{}, errors.New("OCRSecrets is required")
	}

	identities, err := ocr3.MakeIdentities(nca)
	if err != nil {
		return ocr3.OCR2OracleConfig{}, fmt.Errorf("failed to make identities: %w", err)
	}

	cfgBytes := reportingPluginConfigOverride
	if cfgBytes == nil {
		return ocr3.OCR2OracleConfig{}, errors.New("failed to get offchain config: reportingPluginConfigOverride is required for OCR3.1")
	}
	prevConfigDigest, prevHistoryDigest, err := VerifyAndExtractOCR3_1Fields(cfg.PrevConfigDigest, cfg.PrevSeqNr, cfg.PrevHistoryDigest)
	if err != nil {
		return ocr3.OCR2OracleConfig{}, errors.New("VerifyAndExtractOCR3_1Fields failed to verify and extract OCR3.1 fields: " + err.Error())
	}
	var prevSeqNr *uint64
	if cfg.PrevSeqNr != 0 {
		prevSeqNr = &cfg.PrevSeqNr
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
		ocr3_1confighelper.ContractSetConfigArgsOptionalConfig{
			PrevConfigDigest:  prevConfigDigest,
			PrevSeqNr:         prevSeqNr,
			PrevHistoryDigest: prevHistoryDigest,
		},
	)
	if err != nil {
		return ocr3.OCR2OracleConfig{}, fmt.Errorf("failed to generate contract config args: %w", err)
	}

	var configSigners [][]byte
	for _, signer := range signers {
		configSigners = append(configSigners, signer)
	}

	transmitterAddresses, err := evm.AccountToAddress(transmitters)
	if err != nil {
		return ocr3.OCR2OracleConfig{}, fmt.Errorf("failed to convert transmitters to addresses: %w", err)
	}

	config := ocr3.OCR2OracleConfig{
		Signers:               configSigners,
		Transmitters:          transmitterAddresses,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}

	return config, nil
}
