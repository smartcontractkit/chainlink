package validate

import (
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
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
func ToLocalConfig(ocr2Config OCR2Config, insConf InsecureConfig, spec job.OCR2OracleSpec) (types.LocalConfig, error) {
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
	if spec.Relay == relay.Solana && env.MedianPluginCmd.Get() != "" {
		// Work around for Solana Feeds configured with zero values to support LOOP Plugins.
		minOCR2MaxDurationQuery, err := getMinOCR2MaxDurationQuery()
		if err != nil {
			return types.LocalConfig{}, err
		}
		lc.MinOCR2MaxDurationQuery = minOCR2MaxDurationQuery
	}
	if insConf.OCRDevelopmentMode() {
		// Skips config validation so we can use any config parameters we want.
		// For example to lower contractConfigTrackerPollInterval to speed up tests.
		lc.DevelopmentMode = types.EnableDangerousDevelopmentMode
	}
	return lc, nil
}

var (
	minOCR2MaxDurationQuery     = 20 * time.Millisecond
	minOCR2MaxDurationQueryErr  error
	minOCR2MaxDurationQueryOnce sync.Once
)

func getMinOCR2MaxDurationQuery() (time.Duration, error) {
	minOCR2MaxDurationQueryOnce.Do(func() {
		if v := env.MinOCR2MaxDurationQuery.Get(); v != "" {
			minOCR2MaxDurationQuery, minOCR2MaxDurationQueryErr = time.ParseDuration(v)
			if minOCR2MaxDurationQueryErr != nil {
				minOCR2MaxDurationQueryErr = fmt.Errorf("failed to parse %s: %w", env.MinOCR2MaxDurationQuery, minOCR2MaxDurationQueryErr)
			}
		}
	})
	return minOCR2MaxDurationQuery, minOCR2MaxDurationQueryErr
}
