package contracts

import (
	"time"

	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/ocr3_1"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

// values supplied by Alexandr Yepishev as the expected values for OCR3 config
func DefaultOCR3Config() *ocr3.OracleConfig {
	// values supplied by Alexandr Yepishev as the expected values for OCR3 config
	oracleConfig := &keystone_changeset.OracleConfig{
		DeltaProgressMillis:               5000,
		DeltaResendMillis:                 5000,
		DeltaInitialMillis:                5000,
		DeltaRoundMillis:                  2000,
		DeltaGraceMillis:                  500,
		DeltaCertifiedCommitRequestMillis: 1000,
		DeltaStageMillis:                  30000,
		MaxRoundsPerEpoch:                 10,
		MaxDurationQueryMillis:            1000,
		MaxDurationObservationMillis:      1000,
		MaxDurationShouldAcceptMillis:     1000,
		MaxDurationShouldTransmitMillis:   1000,
		MaxFaultyOracles:                  1,
		ConsensusCapOffchainConfig: &ocr3.ConsensusCapOffchainConfig{
			MaxQueryLengthBytes:       1000000,
			MaxObservationLengthBytes: 1000000,
			MaxOutcomeLengthBytes:     1000000,
			MaxReportLengthBytes:      1000000,
			MaxBatchSize:              1000,
			RequestTimeout:            30 * time.Second,
		},
		UniqueReports: true,
	}

	return oracleConfig
}

func DefaultOCR3_1Config(numWorkers int) *ocr3_1.V3_1OracleConfig {
	return &ocr3_1.V3_1OracleConfig{
		DeltaProgressMillis:  5000, // DKG 10-15 seconds; Vault 5 sec // check bandwidth from nops
		DeltaRoundMillis:     200,
		DeltaGraceMillis:     0,
		DeltaStageMillis:     0,
		MaxRoundsPerEpoch:    10,
		TransmissionSchedule: []int{numWorkers},

		MaxDurationInitializationMillis:               10000,
		MaxDurationShouldAcceptAttestedReportMillis:   1000,
		MaxDurationShouldTransmitAcceptedReportMillis: 1000,

		WarnDurationQueryMillis:               1000,
		WarnDurationObservationMillis:         1000,
		WarnDurationValidateObservationMillis: 1000,
		WarnDurationObservationQuorumMillis:   1000,
		WarnDurationStateTransition:           1000,
		WarnDurationCommitted:                 1000,

		MaxFaultyOracles: 1,

		PrevConfigDigest:  "",
		PrevSeqNr:         0,
		PrevHistoryDigest: "",
	}
}

func DefaultChainCapabilityOCR3Config() *ocr3.OracleConfig {
	cfg := DefaultOCR3Config()

	cfg.DeltaRoundMillis = 1000
	const kib = 1024
	const mib = 1024 * kib
	cfg.ConsensusCapOffchainConfig = nil
	cfg.ChainCapOffchainConfig = &ocr3.ChainCapOffchainConfig{
		MaxQueryLengthBytes:       mib,
		MaxObservationLengthBytes: 97 * kib,
		MaxReportLengthBytes:      mib,
		MaxOutcomeLengthBytes:     mib,
		MaxReportCount:            1000,
		MaxBatchSize:              200,
	}
	return cfg
}
