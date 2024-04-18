package ocr

import (
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// Config contains OCR configurations for a job.
type Config interface {
	Insecure() config.Insecure
	JobPipeline() config.JobPipeline
	OCR() config.OCR
}

func toLocalConfig(cfg ValidationConfig, evmOcrConfig evmconfig.OCR, insecureCfg insecureConfig, spec job.OCROracleSpec, ocrConfig job.OCRConfig) ocrtypes.LocalConfig {
	concreteSpec := job.LoadConfigVarsLocalOCR(evmOcrConfig, spec, ocrConfig)
	lc := ocrtypes.LocalConfig{
		BlockchainTimeout:                      concreteSpec.BlockchainTimeout.Duration(),
		ContractConfigConfirmations:            concreteSpec.ContractConfigConfirmations,
		SkipContractConfigConfirmations:        cfg.ChainType().IsL2(),
		ContractConfigTrackerPollInterval:      concreteSpec.ContractConfigTrackerPollInterval.Duration(),
		ContractConfigTrackerSubscribeInterval: concreteSpec.ContractConfigTrackerSubscribeInterval.Duration(),
		ContractTransmitterTransmitTimeout:     concreteSpec.ContractTransmitterTransmitTimeout.Duration(),
		DatabaseTimeout:                        concreteSpec.DatabaseTimeout.Duration(),
		DataSourceTimeout:                      concreteSpec.ObservationTimeout.Duration(),
		DataSourceGracePeriod:                  concreteSpec.ObservationGracePeriod.Duration(),
	}
	if insecureCfg.OCRDevelopmentMode() {
		// Skips config validation so we can use any config parameters we want.
		// For example to lower contractConfigTrackerPollInterval to speed up tests.
		lc.DevelopmentMode = ocrtypes.EnableDangerousDevelopmentMode
	}
	return lc
}
