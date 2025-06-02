package ccipexec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/Masterminds/semver/v3"
	"go.uber.org/multierr"

	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/statuschecker"

	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/factory"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/observability"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/oraclelib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/promwrapper"
)

var (
	// tokenDataWorkerTimeout defines 1) The timeout while waiting for a bg call to the token data 3P provider.
	// 2) When a client requests token data and does not specify a timeout this value is used as a default.
	// 5 seconds is a reasonable value for a timeout.
	// At this moment, minimum OCR Delta Round is set to 30s and deltaGrace to 5s. Based on this configuration
	// 5s for token data worker timeout is a reasonable default.
	tokenDataWorkerTimeout = 5 * time.Second
	// tokenDataWorkerNumWorkers is the number of workers that will be processing token data in parallel.
	tokenDataWorkerNumWorkers = 10
	// expirationDur is the duration for which the token data will be cached.
	expirationDurTokenData = 10 * time.Minute
)

var defaultNewReportingPluginRetryConfig = ccipdata.RetryConfig{
	InitialDelay: time.Second,
	MaxDelay:     10 * time.Minute,
	// Retry for approximately 4hrs (MaxDelay of 10m = 6 times per hour, times 4 hours, plus 10 because the first
	// 10 retries only take 20 minutes due to an initial retry of 1s and exponential backoff)
	MaxRetries: (6 * 4) + 10,
}

func NewExecServices(ctx context.Context, lggr logger.Logger, jb job.Job, srcProvider types.CCIPExecProvider, dstProvider types.CCIPExecProvider, srcChainID int64, dstChainID int64, new bool, argsNoPlugin libocr2.OCR2OracleArgs, logError func(string)) ([]job.ServiceCtx, error) {
	if jb.OCR2OracleSpec == nil {
		return nil, errors.New("spec is nil")
	}
	spec := jb.OCR2OracleSpec
	var pluginConfig ccipconfig.ExecPluginJobSpecConfig
	err := json.Unmarshal(spec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, err
	}

	offRampAddress := ccipcalc.HexToAddress(spec.ContractID)
	offRampReader, err := dstProvider.NewOffRampReader(ctx, offRampAddress)
	if err != nil {
		return nil, fmt.Errorf("create offRampReader: %w", err)
	}

	offRampConfig, err := offRampReader.GetStaticConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("get offRamp static config: %w", err)
	}

	srcChainSelector := offRampConfig.SourceChainSelector
	dstChainSelector := offRampConfig.ChainSelector
	onRampReader, err := srcProvider.NewOnRampReader(ctx, offRampConfig.OnRamp, srcChainSelector, dstChainSelector)
	if err != nil {
		return nil, fmt.Errorf("create onRampReader: %w", err)
	}

	dynamicOnRampConfig, err := onRampReader.GetDynamicConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("get onramp dynamic config: %w", err)
	}

	sourceWrappedNative, err := srcProvider.SourceNativeToken(ctx, dynamicOnRampConfig.Router)
	if err != nil {
		return nil, fmt.Errorf("get source wrapped native token: %w", err)
	}

	srcCommitStore, err := srcProvider.NewCommitStoreReader(ctx, offRampConfig.CommitStore)
	if err != nil {
		return nil, fmt.Errorf("could not create src commitStoreReader reader: %w", err)
	}

	dstCommitStore, err := dstProvider.NewCommitStoreReader(ctx, offRampConfig.CommitStore)
	if err != nil {
		return nil, fmt.Errorf("could not create dst commitStoreReader reader: %w", err)
	}

	var commitStoreReader ccipdata.CommitStoreReader
	commitStoreReader = ccip.NewProviderProxyCommitStoreReader(srcCommitStore, dstCommitStore)

	tokenDataProviders := make(map[cciptypes.Address]tokendata.Reader)
	// init usdc token data provider
	if pluginConfig.USDCConfig.AttestationAPI != "" {
		lggr.Infof("USDC token data provider enabled")
		err2 := pluginConfig.USDCConfig.ValidateUSDCConfig()
		if err2 != nil {
			return nil, err2
		}

		usdcReader, err2 := srcProvider.NewTokenDataReader(ctx, ccip.EvmAddrToGeneric(pluginConfig.USDCConfig.SourceTokenAddress))
		if err2 != nil {
			return nil, fmt.Errorf("new usdc reader: %w", err2)
		}
		tokenDataProviders[cciptypes.Address(pluginConfig.USDCConfig.SourceTokenAddress.String())] = usdcReader
	}
	// init lbtc token data provider
	if pluginConfig.LBTCConfig.AttestationAPI != "" {
		lggr.Infof("LBTC token data provider enabled")
		err2 := pluginConfig.LBTCConfig.ValidateLBTCConfig()
		if err2 != nil {
			return nil, err2
		}

		lbtcReader, err2 := srcProvider.NewTokenDataReader(ctx, ccip.EvmAddrToGeneric(pluginConfig.LBTCConfig.SourceTokenAddress))
		if err2 != nil {
			return nil, fmt.Errorf("new lbtc reader: %w", err2)
		}
		tokenDataProviders[cciptypes.Address(pluginConfig.LBTCConfig.SourceTokenAddress.String())] = lbtcReader
	}

	// Prom wrappers
	onRampReader = observability.NewObservedOnRampReader(onRampReader, srcChainID, ccip.ExecPluginLabel)
	commitStoreReader = observability.NewObservedCommitStoreReader(commitStoreReader, dstChainID, ccip.ExecPluginLabel)
	offRampReader = observability.NewObservedOffRampReader(offRampReader, dstChainID, ccip.ExecPluginLabel)
	metricsCollector := ccip.NewPluginMetricsCollector(ccip.ExecPluginLabel, srcChainID, dstChainID)

	tokenPoolBatchedReader, err := dstProvider.NewTokenPoolBatchedReader(ctx, offRampAddress, srcChainSelector)
	if err != nil {
		return nil, fmt.Errorf("new token pool batched reader: %w", err)
	}

	chainHealthcheck := cache.NewObservedChainHealthCheck(
		cache.NewChainHealthcheck(
			// Adding more details to Logger to make healthcheck logs more informative
			// It's safe because healthcheck logs only in case of unhealthy state
			lggr.With(
				"onramp", offRampConfig.OnRamp,
				"commitStore", offRampConfig.CommitStore,
				"offramp", offRampAddress,
			),
			onRampReader,
			commitStoreReader,
		),
		ccip.ExecPluginLabel,
		srcChainID,
		dstChainID,
		offRampConfig.OnRamp,
	)

	tokenBackgroundWorker := tokendata.NewBackgroundWorker(
		tokenDataProviders,
		tokenDataWorkerNumWorkers,
		tokenDataWorkerTimeout,
		expirationDurTokenData,
	)

	wrappedPluginFactory := NewExecutionReportingPluginFactory(ExecutionPluginStaticConfig{
		lggr:                          lggr,
		onRampReader:                  onRampReader,
		commitStoreReader:             commitStoreReader,
		offRampReader:                 offRampReader,
		sourcePriceRegistryProvider:   ccip.NewChainAgnosticPriceRegistry(srcProvider),
		sourceWrappedNativeToken:      sourceWrappedNative,
		destChainSelector:             dstChainSelector,
		priceRegistryProvider:         ccip.NewChainAgnosticPriceRegistry(dstProvider),
		tokenPoolBatchedReader:        tokenPoolBatchedReader,
		tokenDataWorker:               tokenBackgroundWorker,
		metricsCollector:              metricsCollector,
		chainHealthcheck:              chainHealthcheck,
		newReportingPluginRetryConfig: defaultNewReportingPluginRetryConfig,
		txmStatusChecker:              statuschecker.NewTxmStatusChecker(dstProvider.GetTransactionStatus),
	})

	argsNoPlugin.ReportingPluginFactory = promwrapper.NewPromFactory(wrappedPluginFactory, "CCIPExecution", jb.OCR2OracleSpec.Relay, big.NewInt(0).SetInt64(dstChainID))
	argsNoPlugin.Logger = commonlogger.NewOCRWrapper(lggr, true, logError)
	oracle, err := libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		return nil, err
	}
	// If this is a brand-new job, then we make use of the start blocks. If not then we're rebooting and log poller will pick up where we left off.
	if new {
		return []job.ServiceCtx{
			oraclelib.NewChainAgnosticBackFilledOracle(
				lggr,
				srcProvider,
				dstProvider,
				job.NewServiceAdapter(oracle),
			),
			chainHealthcheck,
			tokenBackgroundWorker,
		}, nil
	}
	return []job.ServiceCtx{
		job.NewServiceAdapter(oracle),
		chainHealthcheck,
		tokenBackgroundWorker,
	}, nil
}

// UnregisterExecPluginLpFilters unregisters all the registered filters for both source and dest chains.
// See comment in UnregisterCommitPluginLpFilters
// It MUST mirror the filters registered in NewExecServices.
// This currently works because the filters registered by the created custom providers when the job is first added
// are stored in the db. Those same filters are unregistered (i.e. deleted from the db) by the newly created providers
// that are passed in from cleanupEVM, as while the providers have no knowledge of each other, they are created
// on the same source and dest relayer.
func UnregisterExecPluginLpFilters(srcProvider types.CCIPExecProvider, dstProvider types.CCIPExecProvider) error {
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

// ExecReportToEthTxMeta generates a txmgr.EthTxMeta from the given report.
// Only MessageIDs will be populated in the TxMeta.
func ExecReportToEthTxMeta(ctx context.Context, typ ccipconfig.ContractType, ver semver.Version) (func(report []byte) (*txmgr.TxMeta, error), error) {
	return factory.ExecReportToEthTxMeta(ctx, typ, ver)
}
