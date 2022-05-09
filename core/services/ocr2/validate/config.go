package validate

import (
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/services/job"
)

//go:generate mockery --name Config --output ../mocks/ --case=underscore

// Config contains OCR2 configurations for a job.
type Config interface {
	config.OCR2Config
	Dev() bool
	LogSQL() bool
	JobPipelineResultWriteQueueDepth() uint64
}

// ToLocalConfig creates a OCR2 LocalConfig from the global config and the OCR2 spec.
func ToLocalConfig(config Config, spec job.OCR2OracleSpec) types.LocalConfig {
	var (
		blockchainTimeout     = time.Duration(spec.BlockchainTimeout)
		ccConfirmations       = spec.ContractConfigConfirmations
		ccTrackerPollInterval = time.Duration(spec.ContractConfigTrackerPollInterval)
	)
	if blockchainTimeout == 0 {
		blockchainTimeout = config.OCR2BlockchainTimeout()
	}
	if ccConfirmations == 0 {
		ccConfirmations = config.OCR2ContractConfirmations()
	}
	if ccTrackerPollInterval == 0 {
		ccTrackerPollInterval = config.OCR2ContractPollInterval()
	}
	lc := types.LocalConfig{
		BlockchainTimeout:                  blockchainTimeout,
		ContractConfigConfirmations:        ccConfirmations,
		ContractConfigTrackerPollInterval:  ccTrackerPollInterval,
		ContractTransmitterTransmitTimeout: config.OCR2ContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                    config.OCR2DatabaseTimeout(),
	}
	if config.Dev() {
		// Skips config validation so we can use any config parameters we want.
		// For example to lower contractConfigTrackerPollInterval to speed up tests.
		lc.DevelopmentMode = types.EnableDangerousDevelopmentMode
	}
	return lc
}
