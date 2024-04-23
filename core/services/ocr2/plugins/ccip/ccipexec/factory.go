package ccipexec

import (
	"context"
	"fmt"
	"sync"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

type ExecutionReportingPluginFactory struct {
	// Config derived from job specs and does not change between instances.
	config ExecutionPluginStaticConfig

	destPriceRegReader ccipdata.PriceRegistryReader
	destPriceRegAddr   cciptypes.Address
	readersMu          *sync.Mutex
}

func NewExecutionReportingPluginFactory(config ExecutionPluginStaticConfig) *ExecutionReportingPluginFactory {
	return &ExecutionReportingPluginFactory{
		config:    config,
		readersMu: &sync.Mutex{},

		// the fields below are initially empty and populated on demand
		destPriceRegReader: nil,
		destPriceRegAddr:   "",
	}
}

func (rf *ExecutionReportingPluginFactory) UpdateDynamicReaders(ctx context.Context, newPriceRegAddr cciptypes.Address) error {
	rf.readersMu.Lock()
	defer rf.readersMu.Unlock()
	// TODO: Investigate use of Close() to cleanup.
	// TODO: a true price registry upgrade on an existing lane may want some kind of start block in its config? Right now we
	// essentially assume that plugins don't care about historical price reg logs.
	if rf.destPriceRegAddr == newPriceRegAddr {
		// No-op
		return nil
	}
	// Close old reader (if present) and open new reader if address changed.
	if rf.destPriceRegReader != nil {
		if err := rf.destPriceRegReader.Close(); err != nil {
			return err
		}
	}

	destPriceRegistryReader, err := rf.config.priceRegistryProvider.NewPriceRegistryReader(context.Background(), newPriceRegAddr)
	if err != nil {
		return err
	}
	rf.destPriceRegReader = destPriceRegistryReader
	rf.destPriceRegAddr = newPriceRegAddr
	return nil
}

func (rf *ExecutionReportingPluginFactory) NewReportingPlugin(config types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	ctx := context.Background() // todo: consider setting a timeout

	destPriceRegistry, destWrappedNative, err := rf.config.offRampReader.ChangeConfig(ctx, config.OnchainConfig, config.OffchainConfig)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}
	// Open dynamic readers
	err = rf.UpdateDynamicReaders(ctx, destPriceRegistry)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}

	offchainConfig, err := rf.config.offRampReader.OffchainConfig(ctx)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, fmt.Errorf("get offchain config from offramp: %w", err)
	}

	gasPriceEstimator, err := rf.config.offRampReader.GasPriceEstimator(ctx)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, fmt.Errorf("get gas price estimator from offramp: %w", err)
	}

	onchainConfig, err := rf.config.offRampReader.OnchainConfig(ctx)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, fmt.Errorf("get onchain config from offramp: %w", err)
	}

	lggr := rf.config.lggr.Named("ExecutionReportingPlugin")
	return &ExecutionReportingPlugin{
			F:                           config.F,
			lggr:                        lggr,
			offchainConfig:              offchainConfig,
			tokenDataWorker:             rf.config.tokenDataWorker,
			gasPriceEstimator:           gasPriceEstimator,
			sourcePriceRegistryProvider: rf.config.sourcePriceRegistryProvider,
			sourcePriceRegistryLock:     sync.RWMutex{},
			sourceWrappedNativeToken:    rf.config.sourceWrappedNativeToken,
			onRampReader:                rf.config.onRampReader,
			commitStoreReader:           rf.config.commitStoreReader,
			destPriceRegistry:           rf.destPriceRegReader,
			destWrappedNative:           destWrappedNative,
			onchainConfig:               onchainConfig,
			offRampReader:               rf.config.offRampReader,
			tokenPoolBatchedReader:      rf.config.tokenPoolBatchedReader,
			inflightReports:             newInflightExecReportsContainer(offchainConfig.InflightCacheExpiry.Duration()),
			commitRootsCache:            cache.NewCommitRootsCache(lggr, onchainConfig.PermissionLessExecutionThresholdSeconds, offchainConfig.RootSnoozeTime.Duration()),
			metricsCollector:            rf.config.metricsCollector,
			chainHealthcheck:            rf.config.chainHealthcheck,
		}, types.ReportingPluginInfo{
			Name: "CCIPExecution",
			// Setting this to false saves on calldata since OffRamp doesn't require agreement between NOPs
			// (OffRamp is only able to execute committed messages).
			UniqueReports: false,
			Limits: types.ReportingPluginLimits{
				MaxObservationLength: ccip.MaxObservationLength,
				MaxReportLength:      MaxExecutionReportLength,
			},
		}, nil
}
