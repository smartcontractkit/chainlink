package offchainreporting

import (
	"github.com/smartcontractkit/chainlink/core/services/job"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

func NewLocalConfig(cfg ValidationConfig, spec job.OffchainReportingOracleSpec) ocrtypes.LocalConfig {
	concreteSpec := job.LoadEnvConfigVarsLocalOCR(cfg, spec)
	lc := ocrtypes.LocalConfig{
		BlockchainTimeout:                      concreteSpec.BlockchainTimeout.Duration(),
		ContractConfigConfirmations:            concreteSpec.ContractConfigConfirmations,
		SkipContractConfigConfirmations:        cfg.ChainType().IsL2(),
		ContractConfigTrackerPollInterval:      concreteSpec.ContractConfigTrackerPollInterval.Duration(),
		ContractConfigTrackerSubscribeInterval: concreteSpec.ContractConfigTrackerSubscribeInterval.Duration(),
		ContractTransmitterTransmitTimeout:     cfg.OCRContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                        concreteSpec.OCRDatabaseTimeout.Duration(),
		DataSourceTimeout:                      concreteSpec.ObservationTimeout.Duration(),
		DataSourceGracePeriod:                  cfg.OCRObservationGracePeriod(),
	}
	if cfg.Dev() {
		// Skips config validation so we can use any config parameters we want.
		// For example to lower contractConfigTrackerPollInterval to speed up tests.
		lc.DevelopmentMode = ocrtypes.EnableDangerousDevelopmentMode
	}
	return lc
}
