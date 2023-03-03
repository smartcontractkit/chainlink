package keepers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	ktypes "github.com/smartcontractkit/ocr2keepers/pkg/types"
)

const maxObservationLength = 1_000

var _ types.ReportingPluginFactory = (*keepersReportingFactory)(nil)

type ReportingFactoryConfig struct {
	CacheExpiration       time.Duration
	CacheEvictionInterval time.Duration
	MaxServiceWorkers     int
	ServiceQueueLength    int
}

type keepersReportingFactory struct {
	headSubscriber ktypes.HeadSubscriber
	registry       ktypes.Registry
	encoder        ktypes.ReportEncoder
	perfLogs       ktypes.PerformLogProvider
	logger         *log.Logger
	config         ReportingFactoryConfig
	upkeepService  *onDemandUpkeepService
}

// NewReportingPluginFactory returns an OCR ReportingPluginFactory. When the plugin
// starts, a separate service is started as a separate go-routine automatically. There
// is no start or stop function for this service so stopping this service relies on
// releasing references to the plugin such that the Go garbage collector cleans up
// hanging routines automatically.
func NewReportingPluginFactory(
	headSubscriber ktypes.HeadSubscriber,
	registry ktypes.Registry,
	perfLogs ktypes.PerformLogProvider,
	encoder ktypes.ReportEncoder,
	logger *log.Logger,
	config ReportingFactoryConfig,
) types.ReportingPluginFactory {
	return &keepersReportingFactory{
		headSubscriber: headSubscriber,
		registry:       registry,
		perfLogs:       perfLogs,
		encoder:        encoder,
		logger:         logger,
		config:         config,
	}
}

// NewReportingPlugin implements the libocr/offchainreporting2/types ReportingPluginFactory interface
func (d *keepersReportingFactory) NewReportingPlugin(c types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	offChainCfg, err := ktypes.DecodeOffchainConfig(c.OffchainConfig)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, fmt.Errorf("%w: failed to decode off chain config", err)
	}

	d.logger.Printf("successfully decoded offchain config when creating plugin: %+v", offChainCfg)

	info := types.ReportingPluginInfo{
		Name: fmt.Sprintf("Oracle %d: Keepers Plugin Instance w/ Digest '%s'", c.OracleID, c.ConfigDigest),
		Limits: types.ReportingPluginLimits{
			// queries should be empty anyway with the current implementation
			MaxQueryLength: 0,
			// an upkeep key is composed of a block number and upkeep id (~40 bytes)
			// an observation is multiple upkeeps to be performed
			// 100 upkeeps to be performed would be a very high upper limit
			// 100 * 10 = 1_000 bytes
			MaxObservationLength: maxObservationLength,
			// a report is composed of 1 or more abi encoded perform calls
			// with performData of arbitrary length
			MaxReportLength: 10_000, // TODO (config): pick sane limit based on expected performData size. maybe set this to block size limit or 2/3 block size limit?
		},
		// UniqueReports increases the threshold of signatures needed for quorum to (n+f)/2 so that it's guaranteed a unique report is generated per round.
		// Fixed to false for ocr2keepers, as we always expect f+1 signatures on a report on contract and do not support uniqueReports quorum
		UniqueReports: false,
	}

	var p float64
	p, err = strconv.ParseFloat(offChainCfg.TargetProbability, 32)
	if err != nil {
		return nil, info, fmt.Errorf("%w: failed to parse configured probability", err)
	}

	sample, err := sampleFromProbability(offChainCfg.TargetInRounds, c.N-c.F, float32(p))
	if err != nil {
		return nil, info, fmt.Errorf("%w: failed to create plugin", err)
	}

	if d.upkeepService != nil {
		d.upkeepService.stop()
	}
	d.upkeepService = newOnDemandUpkeepService(
		sample,
		d.headSubscriber,
		d.registry,
		d.logger,
		time.Duration(offChainCfg.SamplingJobDuration)*time.Millisecond,
		d.config.CacheExpiration,
		d.config.CacheEvictionInterval,
		d.config.MaxServiceWorkers,
		d.config.ServiceQueueLength,
	)

	return &keepers{
		id:      c.OracleID,
		service: d.upkeepService,
		encoder: d.encoder,
		logger:  d.logger,
		filter: newReportCoordinator(
			d.registry,
			time.Duration(offChainCfg.PerformLockoutWindow)*time.Millisecond,
			d.config.CacheEvictionInterval,
			d.perfLogs,
			offChainCfg.MinConfirmations,
			d.logger,
		),
		reportGasLimit:     offChainCfg.GasLimitPerReport,
		upkeepGasOverhead:  offChainCfg.GasOverheadPerUpkeep,
		maxUpkeepBatchSize: offChainCfg.MaxUpkeepBatchSize,
		reportBlockLag:     offChainCfg.ReportBlockLag,
	}, info, nil
}
