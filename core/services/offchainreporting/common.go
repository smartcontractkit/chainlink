package offchainreporting

import (
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
)

func NewLocalConfig(cfg ValidationConfig, spec job.OffchainReportingOracleSpec) types.LocalConfig {
	spec = *job.LoadDynamicConfigVars(cfg, spec)
	lc := types.LocalConfig{
		BlockchainTimeout:                      spec.BlockchainTimeout.Duration(),
		ContractConfigConfirmations:            spec.ContractConfigConfirmations,
		SkipContractConfigConfirmations:        cfg.Chain().IsL2(),
		ContractConfigTrackerPollInterval:      spec.ContractConfigTrackerPollInterval.Duration(),
		ContractConfigTrackerSubscribeInterval: spec.ContractConfigTrackerSubscribeInterval.Duration(),
		ContractTransmitterTransmitTimeout:     cfg.OCRContractTransmitterTransmitTimeout(),
		DatabaseTimeout:                        cfg.OCRDatabaseTimeout(),
		DataSourceTimeout:                      spec.ObservationTimeout.Duration(),
		DataSourceGracePeriod:                  cfg.OCRObservationGracePeriod(),
	}
	if cfg.Dev() {
		// Skips config validation so we can use any config parameters we want.
		// For example to lower contractConfigTrackerPollInterval to speed up tests.
		lc.DevelopmentMode = types.EnableDangerousDevelopmentMode
	}
	return lc
}
