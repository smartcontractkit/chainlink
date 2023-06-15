package validate

import (
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// OCR2Config contains OCR2 configurations for a job.
type OCR2Config interface {
	BlockchainTimeout() time.Duration
	ContractConfirmations() uint16
	ContractPollInterval() time.Duration
	ContractTransmitterTransmitTimeout() time.Duration
	DatabaseTimeout() time.Duration
	TraceLogging() bool
}

type InsecureConfig interface {
	OCRDevelopmentMode() bool
}

// ToLocalConfig creates a OCR2 LocalConfig from the global config and the OCR2 spec.
func ToLocalConfig(ocr2Config OCR2Config, insConf InsecureConfig, spec job.OCR2OracleSpec) types.LocalConfig {
	var (
		blockchainTimeout     = time.Duration(spec.BlockchainTimeout)
		ccConfirmations       = spec.ContractConfigConfirmations
		ccTrackerPollInterval = time.Duration(spec.ContractConfigTrackerPollInterval)
	)
	if blockchainTimeout == 0 {
		blockchainTimeout = ocr2Config.BlockchainTimeout()
	}
	if ccConfirmations == 0 {
		ccConfirmations = ocr2Config.ContractConfirmations()
	}
	if ccTrackerPollInterval == 0 {
		ccTrackerPollInterval = ocr2Config.ContractPollInterval()
	}
	lc := types.LocalConfig{
		BlockchainTimeout:                  blockchainTimeout,
		ContractConfigConfirmations:        ccConfirmations,
		ContractConfigTrackerPollInterval:  ccTrackerPollInterval,
		ContractTransmitterTransmitTimeout: ocr2Config.ContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                    ocr2Config.DatabaseTimeout(),
	}
	if insConf.OCRDevelopmentMode() {
		// Skips config validation so we can use any config parameters we want.
		// For example to lower contractConfigTrackerPollInterval to speed up tests.
		lc.DevelopmentMode = types.EnableDangerousDevelopmentMode
	}
	return lc
}
