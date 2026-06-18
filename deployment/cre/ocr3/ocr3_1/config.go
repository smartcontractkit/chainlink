package ocr3_1

import (
	"errors"
	"fmt"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1confighelper"

	focr "github.com/smartcontractkit/chainlink-deployments-framework/offchain/ocr"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"

	evm "github.com/smartcontractkit/chainlink-evm/pkg/relay"
	"github.com/smartcontractkit/chainlink/deployment"
)

type V3_1OracleConfig struct {
	DeltaProgressMillis  uint32 `yaml:"deltaProgressMillis" json:"deltaProgressMillis"`
	DeltaRoundMillis     uint32 `yaml:"deltaRoundMillis" json:"deltaRoundMillis"`
	DeltaGraceMillis     uint32 `yaml:"deltaGraceMillis" json:"deltaGraceMillis"`
	DeltaStageMillis     uint32 `yaml:"deltaStageMillis" json:"deltaStageMillis"`
	MaxRoundsPerEpoch    uint64 `yaml:"maxRoundsPerEpoch" json:"maxRoundsPerEpoch"`
	TransmissionSchedule []int  `yaml:"transmissionSchedule" json:"transmissionSchedule"`

	MaxDurationInitializationMillis               uint32 `yaml:"maxDurationInitializationMillis" json:"maxDurationInitializationMillis"`
	MaxDurationShouldAcceptAttestedReportMillis   uint32 `yaml:"maxDurationShouldAcceptAttestedReportMillis" json:"maxDurationShouldAcceptAttestedReportMillis"`
	MaxDurationShouldTransmitAcceptedReportMillis uint32 `yaml:"maxDurationShouldTransmitAcceptedReportMillis" json:"maxDurationShouldTransmitAcceptedReportMillis"`

	WarnDurationQueryMillis               uint32 `yaml:"warnDurationQueryMillis" json:"warnDurationQueryMillis"`
	WarnDurationObservationMillis         uint32 `yaml:"warnDurationObservationMillis" json:"warnDurationObservationMillis"`
	WarnDurationValidateObservationMillis uint32 `yaml:"warnDurationValidateObservationMillis" json:"warnDurationValidateObservationMillis"`
	WarnDurationObservationQuorumMillis   uint32 `yaml:"warnDurationObservationQuorumMillis" json:"warnDurationObservationQuorumMillis"`
	WarnDurationStateTransition           uint32 `yaml:"warnDurationStateTransition" json:"warnDurationStateTransition"`
	WarnDurationCommitted                 uint32 `yaml:"warnDurationCommitted" json:"warnDurationCommitted"`

	MaxFaultyOracles int `yaml:"maxFaultyOracles" json:"maxFaultyOracles"`

	PrevConfigDigest  string `yaml:"prevConfigDigest" json:"prevConfigDigest"`
	PrevSeqNr         uint64 `yaml:"prevSeqNr" json:"prevSeqNr"`
	PrevHistoryDigest string `yaml:"prevHistoryDigest" json:"prevHistoryDigest"`

	DKGOffchainConfig   *DKGOffchainConfig   `yaml:"dkgOffchainConfig,omitempty" json:"dkgOffchainConfig,omitempty"`
	VaultOffchainConfig *VaultOffchainConfig `yaml:"vaultOffchainConfig,omitempty" json:"vaultOffchainConfig,omitempty"`
}

func GenerateOCR3_1ConfigFromNodes(cfg V3_1OracleConfig, nodes []deployment.Node, registryChainSel uint64, secrets focr.OCRSecrets, reportingPluginConfigOverride []byte, extraSignerFamilies []string) (ocr3.OCR2OracleConfig, error) {
	nca := ocr3.MakeNodeKeysSlice(nodes, registryChainSel, extraSignerFamilies)
	return GenerateOCR3_1Config(cfg, nca, secrets, reportingPluginConfigOverride)
}

func GenerateOCR3_1Config(cfg V3_1OracleConfig, nca []ocr3.NodeKeys, secrets focr.OCRSecrets, reportingPluginConfigOverride []byte) (ocr3.OCR2OracleConfig, error) {
	// the transmission schedule is very specific; arguably it should be not be a parameter
	if len(cfg.TransmissionSchedule) != 1 || cfg.TransmissionSchedule[0] != len(nca) {
		return ocr3.OCR2OracleConfig{}, fmt.Errorf("transmission schedule must have exactly one entry, matching the len of the number of nodes want [%d], got %v. Total TransmissionSchedules = %d", len(nca), cfg.TransmissionSchedule, len(cfg.TransmissionSchedule))
	}

	identities, err := ocr3.MakeIdentities(nca)
	if err != nil {
		return ocr3.OCR2OracleConfig{}, fmt.Errorf("failed to make identities: %w", err)
	}

	return genOCR3_1Config(cfg, identities, secrets, reportingPluginConfigOverride)
}

func genOCR3_1Config(cfg V3_1OracleConfig, identities []confighelper.OracleIdentityExtra, secrets focr.OCRSecrets, cfgBytes []byte) (ocr3.OCR2OracleConfig, error) {
	if secrets.IsEmpty() {
		return ocr3.OCR2OracleConfig{}, errors.New("OCRSecrets is required")
	}

	if cfgBytes == nil {
		pc, err := getPluginConfig(cfg)
		if err != nil {
			return ocr3.OCR2OracleConfig{}, err
		}
		if pc != nil {
			cfgBytes, err = pc.Marshal()
			if err != nil {
				return ocr3.OCR2OracleConfig{}, fmt.Errorf("failed to marshal plugin config: %w", err)
			}
		}
	}
	if cfgBytes == nil {
		return ocr3.OCR2OracleConfig{}, errors.New("failed to get offchain config: one of reportingPluginConfigOverride, DKGOffchainConfig, or VaultOffchainConfig is required")
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

	return ocr3.OCR2OracleConfig{
		Signers:               configSigners,
		Transmitters:          transmitterAddresses,
		F:                     f,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}, nil
}
