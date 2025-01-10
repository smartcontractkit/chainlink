package ccipcommit

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	cciporm "github.com/smartcontractkit/chainlink/v2/core/services/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	db "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdb"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/factory"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/observability"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/oraclelib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/promwrapper"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

var defaultNewReportingPluginRetryConfig = ccipdata.RetryConfig{
	InitialDelay: time.Second,
	MaxDelay:     10 * time.Minute,
	// Retry for approximately 4hrs (MaxDelay of 10m = 6 times per hour, times 4 hours, plus 10 because the first
	// 10 retries only take 20 minutes due to an initial retry of 1s and exponential backoff)
	MaxRetries: (6 * 4) + 10,
}

func NewCommitServices(
	ctx context.Context,
	ds sqlutil.DataSource,
	srcProvider commontypes.CCIPCommitProvider,
	dstProvider commontypes.CCIPCommitProvider,
	priceGetter ccip.AllTokensPriceGetter,
	jb job.Job,
	lggr logger.Logger,
	pr pipeline.Runner,
	argsNoPlugin libocr2.OCR2OracleArgs,
	newInstance bool,
	sourceChainID int64,
	destChainID int64,
	logError func(string),
) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec

	var pluginConfig ccipconfig.CommitPluginJobSpecConfig
	err := json.Unmarshal(spec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, err
	}

	commitStoreAddress := common.HexToAddress(spec.ContractID)

	// commit store contract doesn't exist on the source chain, but we have an implementation of it
	// to get access to a gas estimator on the source chain
	srcCommitStore, err := srcProvider.NewCommitStoreReader(ctx, ccipcalc.EvmAddrToGeneric(commitStoreAddress))
	if err != nil {
		return nil, err
	}

	dstCommitStore, err := dstProvider.NewCommitStoreReader(ctx, ccipcalc.EvmAddrToGeneric(commitStoreAddress))
	if err != nil {
		return nil, err
	}

	var commitStoreReader ccipdata.CommitStoreReader
	commitStoreReader = ccip.NewProviderProxyCommitStoreReader(srcCommitStore, dstCommitStore)
	commitLggr := lggr.Named("CCIPCommit").With("sourceChain", sourceChainID, "destChain", destChainID)

	offRampReader, err := dstProvider.NewOffRampReader(ctx, pluginConfig.OffRamp)
	if err != nil {
		return nil, err
	}

	staticConfig, err := commitStoreReader.GetCommitStoreStaticConfig(ctx)
	if err != nil {
		return nil, err
	}
	onRampAddress := staticConfig.OnRamp

	onRampReader, err := srcProvider.NewOnRampReader(ctx, onRampAddress, staticConfig.SourceChainSelector, staticConfig.ChainSelector)
	if err != nil {
		return nil, err
	}

	onRampRouterAddr, err := onRampReader.RouterAddress(ctx)
	if err != nil {
		return nil, err
	}
	sourceNative, err := srcProvider.SourceNativeToken(ctx, onRampRouterAddr)
	if err != nil {
		return nil, err
	}
	// Prom wrappers
	onRampReader = observability.NewObservedOnRampReader(onRampReader, sourceChainID, ccip.CommitPluginLabel)
	commitStoreReader = observability.NewObservedCommitStoreReader(commitStoreReader, destChainID, ccip.CommitPluginLabel)
	offRampReader = observability.NewObservedOffRampReader(offRampReader, destChainID, ccip.CommitPluginLabel)
	metricsCollector := ccip.NewPluginMetricsCollector(ccip.CommitPluginLabel, sourceChainID, destChainID)

	chainHealthCheck := cache.NewObservedChainHealthCheck(
		cache.NewChainHealthcheck(
			// Adding more details to Logger to make healthcheck logs more informative
			// It's safe because healthcheck logs only in case of unhealthy state
			lggr.With(
				"onramp", onRampAddress,
				"commitStore", commitStoreAddress,
				"offramp", pluginConfig.OffRamp,
			),
			onRampReader,
			commitStoreReader,
		),
		ccip.CommitPluginLabel,
		sourceChainID, // assuming this is the chain id?
		destChainID,
		onRampAddress,
	)

	orm, err := cciporm.NewORM(ds, lggr)
	if err != nil {
		return nil, err
	}

	priceService := db.NewPriceService(
		lggr,
		orm,
		jb.ID,
		staticConfig.ChainSelector,
		staticConfig.SourceChainSelector,
		sourceNative,
		priceGetter,
		offRampReader,
	)

	wrappedPluginFactory := NewCommitReportingPluginFactory(CommitPluginStaticConfig{
		lggr:                          lggr,
		newReportingPluginRetryConfig: defaultNewReportingPluginRetryConfig,
		onRampReader:                  onRampReader,
		sourceChainSelector:           staticConfig.SourceChainSelector,
		sourceNative:                  sourceNative,
		offRamp:                       offRampReader,
		commitStore:                   commitStoreReader,
		destChainSelector:             staticConfig.ChainSelector,
		priceRegistryProvider:         ccip.NewChainAgnosticPriceRegistry(dstProvider),
		metricsCollector:              metricsCollector,
		chainHealthcheck:              chainHealthCheck,
		priceService:                  priceService,
	})
	argsNoPlugin.ReportingPluginFactory = promwrapper.NewPromFactory(wrappedPluginFactory, "CCIPCommit", jb.OCR2OracleSpec.Relay, big.NewInt(0).SetInt64(destChainID))
	argsNoPlugin.Logger = commonlogger.NewOCRWrapper(commitLggr, true, logError)
	oracle, err := libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		return nil, err
	}
	// If this is a brand-new job, then we make use of the start blocks. If not then we're rebooting and log poller will pick up where we left off.
	if newInstance {
		return []job.ServiceCtx{
			oraclelib.NewChainAgnosticBackFilledOracle(
				lggr,
				srcProvider,
				dstProvider,
				job.NewServiceAdapter(oracle),
			),
			chainHealthCheck,
			priceService,
		}, nil
	}
	return []job.ServiceCtx{
		job.NewServiceAdapter(oracle),
		chainHealthCheck,
		priceService,
	}, nil
}

func CommitReportToEthTxMeta(typ ccipconfig.ContractType, ver semver.Version) (func(report []byte) (*txmgr.TxMeta, error), error) {
	return factory.CommitReportToEthTxMeta(typ, ver)
}

// UnregisterCommitPluginLpFilters unregisters all the registered filters for both source and dest chains.
// NOTE: The transaction MUST be used here for CLO's monster tx to function as expected
// https://github.com/smartcontractkit/ccip/blob/68e2197472fb017dd4e5630d21e7878d58bc2a44/core/services/feeds/service.go#L716
// TODO once that transaction is broken up, we should be able to simply rely on oracle.Close() to cleanup the filters.
// Until then we have to deterministically reload the readers from the spec (and thus their filters) and close them.
func UnregisterCommitPluginLpFilters(srcProvider commontypes.CCIPCommitProvider, dstProvider commontypes.CCIPCommitProvider) error {
	unregisterFuncs := []func() error{
		func() error {
			return srcProvider.Close()
		},
		func() error {
			return dstProvider.Close()
		},
	}

	var multiErr error
	for _, fn := range unregisterFuncs {
		if err := fn(); err != nil {
			multiErr = multierr.Append(multiErr, err)
		}
	}
	return multiErr
}
