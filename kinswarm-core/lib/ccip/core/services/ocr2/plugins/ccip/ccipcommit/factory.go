package ccipcommit

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

type CommitReportingPluginFactory struct {
	// Configuration derived from the job spec which does not change
	// between plugin instances (ie between SetConfigs onchain)
	config CommitPluginStaticConfig

	// Dynamic readers
	readersMu          *sync.Mutex
	destPriceRegReader ccipdata.PriceRegistryReader
	destPriceRegAddr   common.Address
}

// NewCommitReportingPluginFactory return a new CommitReportingPluginFactory.
func NewCommitReportingPluginFactory(config CommitPluginStaticConfig) *CommitReportingPluginFactory {
	return &CommitReportingPluginFactory{
		config:    config,
		readersMu: &sync.Mutex{},

		// the fields below are initially empty and populated on demand
		destPriceRegReader: nil,
		destPriceRegAddr:   common.Address{},
	}
}

func (rf *CommitReportingPluginFactory) UpdateDynamicReaders(ctx context.Context, newPriceRegAddr common.Address) error {
	rf.readersMu.Lock()
	defer rf.readersMu.Unlock()
	// TODO: Investigate use of Close() to cleanup.
	// TODO: a true price registry upgrade on an existing lane may want some kind of start block in its config? Right now we
	// essentially assume that plugins don't care about historical price reg logs.
	if rf.destPriceRegAddr == newPriceRegAddr {
		// No-op
		return nil
	}
	// Close old reader if present and open new reader if address changed
	if rf.destPriceRegReader != nil {
		if err := rf.destPriceRegReader.Close(); err != nil {
			return err
		}
	}

	destPriceRegistryReader, err := rf.config.priceRegistryProvider.NewPriceRegistryReader(context.Background(), cciptypes.Address(newPriceRegAddr.String()))
	if err != nil {
		return fmt.Errorf("init dynamic price registry: %w", err)
	}
	rf.destPriceRegReader = destPriceRegistryReader
	rf.destPriceRegAddr = newPriceRegAddr
	return nil
}

type reportingPluginAndInfo struct {
	plugin     types.ReportingPlugin
	pluginInfo types.ReportingPluginInfo
}

// NewReportingPlugin registers a new ReportingPlugin
func (rf *CommitReportingPluginFactory) NewReportingPlugin(config types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	initialRetryDelay := rf.config.newReportingPluginRetryConfig.InitialDelay
	maxDelay := rf.config.newReportingPluginRetryConfig.MaxDelay
	maxRetries := rf.config.newReportingPluginRetryConfig.MaxRetries

	pluginAndInfo, err := ccipcommon.RetryUntilSuccess(
		rf.NewReportingPluginFn(config), initialRetryDelay, maxDelay, maxRetries,
	)
	if err != nil {
		return nil, types.ReportingPluginInfo{}, err
	}
	return pluginAndInfo.plugin, pluginAndInfo.pluginInfo, err
}

// NewReportingPluginFn implements the NewReportingPlugin logic. It is defined as a function so that it can easily be
// retried via RetryUntilSuccess. NewReportingPlugin must return successfully in order for the Commit plugin to
// function, hence why we can only keep retrying it until it succeeds.
func (rf *CommitReportingPluginFactory) NewReportingPluginFn(config types.ReportingPluginConfig) func() (reportingPluginAndInfo, error) {
	newReportingPluginFn := func() (reportingPluginAndInfo, error) {
		ctx := context.Background() // todo: consider adding some timeout

		destPriceReg, err := rf.config.commitStore.ChangeConfig(ctx, config.OnchainConfig, config.OffchainConfig)
		if err != nil {
			return reportingPluginAndInfo{}, fmt.Errorf("commitStore.ChangeConfig error: %w", err)
		}

		priceRegEvmAddr, err := ccipcalc.GenericAddrToEvm(destPriceReg)
		if err != nil {
			return reportingPluginAndInfo{}, fmt.Errorf("GenericAddrToEvm error: %w", err)
		}
		if err = rf.UpdateDynamicReaders(ctx, priceRegEvmAddr); err != nil {
			return reportingPluginAndInfo{}, fmt.Errorf("UpdateDynamicReaders error: %w", err)
		}

		pluginOffChainConfig, err := rf.config.commitStore.OffchainConfig(ctx)
		if err != nil {
			return reportingPluginAndInfo{}, fmt.Errorf("commitStore.OffchainConfig error: %w", err)
		}

		gasPriceEstimator, err := rf.config.commitStore.GasPriceEstimator(ctx)
		if err != nil {
			return reportingPluginAndInfo{}, fmt.Errorf("commitStore.GasPriceEstimator error: %w", err)
		}

		err = rf.config.priceService.UpdateDynamicConfig(ctx, gasPriceEstimator, rf.destPriceRegReader)
		if err != nil {
			return reportingPluginAndInfo{}, fmt.Errorf("priceService.UpdateDynamicConfig error: %w", err)
		}

		lggr := rf.config.lggr.Named("CommitReportingPlugin")
		plugin := &CommitReportingPlugin{
			sourceChainSelector:     rf.config.sourceChainSelector,
			sourceNative:            rf.config.sourceNative,
			onRampReader:            rf.config.onRampReader,
			destChainSelector:       rf.config.destChainSelector,
			commitStoreReader:       rf.config.commitStore,
			F:                       config.F,
			lggr:                    lggr,
			destPriceRegistryReader: rf.destPriceRegReader,
			offRampReader:           rf.config.offRamp,
			gasPriceEstimator:       gasPriceEstimator,
			offchainConfig:          pluginOffChainConfig,
			metricsCollector:        rf.config.metricsCollector,
			chainHealthcheck:        rf.config.chainHealthcheck,
			priceService:            rf.config.priceService,
		}

		pluginInfo := types.ReportingPluginInfo{
			Name:          "CCIPCommit",
			UniqueReports: false, // See comment in CommitStore constructor.
			Limits: types.ReportingPluginLimits{
				MaxQueryLength:       ccip.MaxQueryLength,
				MaxObservationLength: ccip.MaxObservationLength,
				MaxReportLength:      MaxCommitReportLength,
			},
		}

		return reportingPluginAndInfo{plugin, pluginInfo}, nil
	}

	return func() (reportingPluginAndInfo, error) {
		result, err := newReportingPluginFn()
		if err != nil {
			rf.config.lggr.Errorw("NewReportingPlugin failed", "err", err)
			rf.config.metricsCollector.NewReportingPluginError()
		}

		return result, err
	}
}
