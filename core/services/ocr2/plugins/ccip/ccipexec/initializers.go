package ccipexec

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/Masterminds/semver/v3"
	"go.uber.org/multierr"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
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
	tokenDataWorkerNumWorkers = 5
)

var defaultNewReportingPluginRetryConfig = ccipdata.RetryConfig{InitialDelay: time.Second, MaxDelay: 5 * time.Minute}

func NewExecServices(ctx context.Context, lggr logger.Logger, jb job.Job, srcProvider types.CCIPExecProvider, dstProvider types.CCIPExecProvider, srcChain legacyevm.Chain, dstChain legacyevm.Chain, srcChainID int64, dstChainID int64, chainSet legacyevm.LegacyChainContainer, new bool, argsNoPlugin libocr2.OCR2OracleArgs, logError func(string)) ([]job.ServiceCtx, error) {
	if jb.OCR2OracleSpec == nil {
		return nil, fmt.Errorf("spec is nil")
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

	versionFinder := ccip.NewEvmVersionFinder()
	commitStoreReader, err := factory.NewCommitStoreReader(lggr, versionFinder, offRampConfig.CommitStore, dstChain.Client(), dstChain.LogPoller())
	if err != nil {
		return nil, fmt.Errorf("could not load commitStoreReader reader: %w", err)
	}

	err = commitStoreReader.SetGasEstimator(ctx, srcChain.GasEstimator())
	if err != nil {
		return nil, fmt.Errorf("could not set gas estimator: %w", err)
	}

	err = commitStoreReader.SetSourceMaxGasPrice(ctx, srcChain.Config().EVM().GasEstimator().PriceMax().ToInt())
	if err != nil {
		return nil, fmt.Errorf("could not set source max gas price: %w", err)
	}

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
		2*tokenDataWorkerTimeout,
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
func UnregisterExecPluginLpFilters(ctx context.Context, lggr logger.Logger, jb job.Job, chainSet legacyevm.LegacyChainContainer) error {
	params, err := extractJobSpecParams(lggr, jb, chainSet, false)
	if err != nil {
		return err
	}

	offRampAddress, err := params.offRampReader.Address(ctx)
	if err != nil {
		return fmt.Errorf("get offramp reader address: %w", err)
	}

	versionFinder := factory.NewEvmVersionFinder()
	unregisterFuncs := []func() error{
		func() error {
			return factory.CloseCommitStoreReader(lggr, versionFinder, params.offRampConfig.CommitStore, params.destChain.Client(), params.destChain.LogPoller())
		},
		func() error {
			return factory.CloseOnRampReader(lggr, versionFinder, params.offRampConfig.SourceChainSelector, params.offRampConfig.ChainSelector, params.offRampConfig.OnRamp, params.sourceChain.LogPoller(), params.sourceChain.Client())
		},
		func() error {
			return factory.CloseOffRampReader(lggr, versionFinder, offRampAddress, params.destChain.Client(), params.destChain.LogPoller(), params.destChain.GasEstimator(), params.destChain.Config().EVM().GasEstimator().PriceMax().ToInt())
		},
		func() error { // usdc token data reader
			if usdcDisabled := params.pluginConfig.USDCConfig.AttestationAPI == ""; usdcDisabled {
				return nil
			}
			return ccipdata.CloseUSDCReader(lggr, jobIDToString(jb.ID), params.pluginConfig.USDCConfig.SourceMessageTransmitterAddress, params.sourceChain.LogPoller())
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

type jobSpecParams struct {
	pluginConfig  ccipconfig.ExecPluginJobSpecConfig
	offRampConfig cciptypes.OffRampStaticConfig
	offRampReader ccipdata.OffRampReader
	sourceChain   legacyevm.Chain
	destChain     legacyevm.Chain
}

func extractJobSpecParams(lggr logger.Logger, jb job.Job, chainSet legacyevm.LegacyChainContainer, registerFilters bool) (*jobSpecParams, error) {
	if jb.OCR2OracleSpec == nil {
		return nil, fmt.Errorf("spec is nil")
	}
	spec := jb.OCR2OracleSpec
	var pluginConfig ccipconfig.ExecPluginJobSpecConfig
	err := json.Unmarshal(spec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, err
	}

	destChain, _, err := ccipconfig.GetChainFromSpec(spec, chainSet)
	if err != nil {
		return nil, err
	}

	versionFinder := factory.NewEvmVersionFinder()
	offRampAddress := ccipcalc.HexToAddress(spec.ContractID)
	offRampReader, err := factory.NewOffRampReader(lggr, versionFinder, offRampAddress, destChain.Client(), destChain.LogPoller(), destChain.GasEstimator(), destChain.Config().EVM().GasEstimator().PriceMax().ToInt(), registerFilters)
	if err != nil {
		return nil, fmt.Errorf("create offRampReader: %w", err)
	}

	offRampConfig, err := offRampReader.GetStaticConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get offRamp static config: %w", err)
	}

	chainID, err := chainselectors.ChainIdFromSelector(offRampConfig.SourceChainSelector)
	if err != nil {
		return nil, err
	}

	sourceChain, err := chainSet.Get(strconv.FormatUint(chainID, 10))
	if err != nil {
		return nil, fmt.Errorf("open source chain: %w", err)
	}

	return &jobSpecParams{
		pluginConfig:  pluginConfig,
		offRampConfig: offRampConfig,
		offRampReader: offRampReader,
		sourceChain:   sourceChain,
		destChain:     destChain,
	}, nil
}

func jobIDToString(id int32) string {
	return fmt.Sprintf("job_%d", id)
}
