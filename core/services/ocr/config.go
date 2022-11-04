package ocr

import (
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// Config contains OCR configurations for a job.
type Config interface {
	pg.QConfig
}

func toLocalConfig(cfg ValidationConfig, spec job.OCROracleSpec) ocrtypes.LocalConfig {
	concreteSpec := job.LoadEnvConfigVarsLocalOCR(cfg, spec)
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
	if cfg.Dev() {
		// Skips config validation so we can use any config parameters we want.
		// For example to lower contractConfigTrackerPollInterval to speed up tests.
		lc.DevelopmentMode = ocrtypes.EnableDangerousDevelopmentMode
	}
	return lc
}
