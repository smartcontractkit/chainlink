package ccipexec

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/batchreader"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/ccipdataprovider"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/factory"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/observability"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/oraclelib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata/usdc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/promwrapper"
)

const numTokenDataWorkers = 5

func NewExecutionServices(ctx context.Context, lggr logger.Logger, jb job.Job, chainSet legacyevm.LegacyChainContainer, new bool, argsNoPlugin libocr2.OCR2OracleArgs, logError func(string)) ([]job.ServiceCtx, error) {
	execPluginConfig, backfillArgs, chainHealthcheck, tokenWorker, err := jobSpecToExecPluginConfig(ctx, lggr, jb, chainSet)
	if err != nil {
		return nil, err
	}
	wrappedPluginFactory := NewExecutionReportingPluginFactory(*execPluginConfig)
	destChainID, err := chainselectors.ChainIdFromSelector(execPluginConfig.destChainSelector)
	if err != nil {
		return nil, err
	}
	argsNoPlugin.ReportingPluginFactory = promwrapper.NewPromFactory(wrappedPluginFactory, "CCIPExecution", jb.OCR2OracleSpec.Relay, big.NewInt(0).SetUint64(destChainID))
	argsNoPlugin.Logger = commonlogger.NewOCRWrapper(execPluginConfig.lggr, true, logError)
	oracle, err := libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		return nil, err
	}
	// If this is a brand-new job, then we make use of the start blocks. If not then we're rebooting and log poller will pick up where we left off.
	if new {
		return []job.ServiceCtx{
			oraclelib.NewBackfilledOracle(
				execPluginConfig.lggr,
				backfillArgs.SourceLP,
				backfillArgs.DestLP,
				backfillArgs.SourceStartBlock,
				backfillArgs.DestStartBlock,
				job.NewServiceAdapter(oracle),
			),
			chainHealthcheck,
			tokenWorker,
		}, nil
	}
	return []job.ServiceCtx{
		job.NewServiceAdapter(oracle),
		chainHealthcheck,
		tokenWorker,
	}, nil
}

// UnregisterExecPluginLpFilters unregisters all the registered filters for both source and dest chains.
// See comment in UnregisterCommitPluginLpFilters
// It MUST mirror the filters registered in NewExecutionServices.
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
			return factory.CloseCommitStoreReader(lggr, versionFinder, params.offRampConfig.CommitStore, params.destChain.Client(), params.destChain.LogPoller(), params.sourceChain.GasEstimator(), params.sourceChain.Config().EVM().GasEstimator().PriceMax().ToInt())
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

func initTokenDataProviders(lggr logger.Logger, jobID string, pluginConfig ccipconfig.ExecutionPluginJobSpecConfig, sourceLP logpoller.LogPoller) (map[cciptypes.Address]tokendata.Reader, error) {
	tokenDataProviders := make(map[cciptypes.Address]tokendata.Reader)

	// init usdc token data provider
	if pluginConfig.USDCConfig.AttestationAPI != "" {
		lggr.Infof("USDC token data provider enabled")
		err := pluginConfig.USDCConfig.ValidateUSDCConfig()
		if err != nil {
			return nil, err
		}

		attestationURI, err := url.ParseRequestURI(pluginConfig.USDCConfig.AttestationAPI)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse USDC attestation API")
		}

		usdcReader, err := ccipdata.NewUSDCReader(lggr, jobID, pluginConfig.USDCConfig.SourceMessageTransmitterAddress, sourceLP, true)
		if err != nil {
			return nil, errors.Wrap(err, "new usdc reader")
		}

		tokenDataProviders[cciptypes.Address(pluginConfig.USDCConfig.SourceTokenAddress.String())] =
			usdc.NewUSDCTokenDataReader(
				lggr,
				usdcReader,
				attestationURI,
				int(pluginConfig.USDCConfig.AttestationAPITimeoutSeconds),
				pluginConfig.USDCConfig.SourceTokenAddress,
				time.Duration(pluginConfig.USDCConfig.AttestationAPIIntervalMilliseconds)*time.Millisecond,
			)
	}

	return tokenDataProviders, nil
}

func jobSpecToExecPluginConfig(ctx context.Context, lggr logger.Logger, jb job.Job, chainSet legacyevm.LegacyChainContainer) (*ExecutionPluginStaticConfig, *ccipcommon.BackfillArgs, *cache.ObservedChainHealthcheck, *tokendata.BackgroundWorker, error) {
	params, err := extractJobSpecParams(lggr, jb, chainSet, true)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	lggr.Infow("Initializing exec plugin",
		"CommitStore", params.offRampConfig.CommitStore,
		"OnRamp", params.offRampConfig.OnRamp,
		"ArmProxy", params.offRampConfig.ArmProxy,
		"SourceChainSelector", params.offRampConfig.SourceChainSelector,
		"DestChainSelector", params.offRampConfig.ChainSelector)

	sourceChainID := params.sourceChain.ID().Int64()
	destChainID := params.destChain.ID().Int64()
	versionFinder := factory.NewEvmVersionFinder()

	sourceChainName, destChainName, err := ccipconfig.ResolveChainNames(sourceChainID, destChainID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	execLggr := lggr.Named("CCIPExecution").With("sourceChain", sourceChainName, "destChain", destChainName)
	onRampReader, err := factory.NewOnRampReader(execLggr, versionFinder, params.offRampConfig.SourceChainSelector, params.offRampConfig.ChainSelector, params.offRampConfig.OnRamp, params.sourceChain.LogPoller(), params.sourceChain.Client())
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "create onramp reader")
	}
	dynamicOnRampConfig, err := onRampReader.GetDynamicConfig(ctx)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "get onramp dynamic config")
	}

	routerAddr, err := ccipcalc.GenericAddrToEvm(dynamicOnRampConfig.Router)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	sourceRouter, err := router.NewRouter(routerAddr, params.sourceChain.Client())
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "failed loading source router")
	}
	sourceWrappedNative, err := sourceRouter.GetWrappedNative(&bind.CallOpts{})
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "could not get source native token")
	}

	commitStoreReader, err := factory.NewCommitStoreReader(lggr, versionFinder, params.offRampConfig.CommitStore, params.destChain.Client(), params.destChain.LogPoller(), params.sourceChain.GasEstimator(), params.sourceChain.Config().EVM().GasEstimator().PriceMax().ToInt())
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "could not load commitStoreReader reader")
	}

	tokenDataProviders, err := initTokenDataProviders(lggr, jobIDToString(jb.ID), params.pluginConfig, params.sourceChain.LogPoller())
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "could not get token data providers")
	}

	// Prom wrappers
	onRampReader = observability.NewObservedOnRampReader(onRampReader, sourceChainID, ccip.ExecPluginLabel)
	commitStoreReader = observability.NewObservedCommitStoreReader(commitStoreReader, destChainID, ccip.ExecPluginLabel)
	offRampReader := observability.NewObservedOffRampReader(params.offRampReader, destChainID, ccip.ExecPluginLabel)
	metricsCollector := ccip.NewPluginMetricsCollector(ccip.ExecPluginLabel, sourceChainID, destChainID)

	destChainSelector, err := chainselectors.SelectorFromChainId(uint64(destChainID))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get chain %d selector: %w", destChainID, err)
	}
	sourceChainSelector, err := chainselectors.SelectorFromChainId(uint64(sourceChainID))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get chain %d selector: %w", sourceChainID, err)
	}

	execLggr.Infow("Initialized exec plugin",
		"pluginConfig", params.pluginConfig,
		"onRampAddress", params.offRampConfig.OnRamp,
		"dynamicOnRampConfig", dynamicOnRampConfig,
		"sourceNative", sourceWrappedNative,
		"sourceRouter", sourceRouter.Address())

	batchCaller := rpclib.NewDynamicLimitedBatchCaller(
		lggr,
		params.destChain.Client(),
		rpclib.DefaultRpcBatchSizeLimit,
		rpclib.DefaultRpcBatchBackOffMultiplier,
		rpclib.DefaultMaxParallelRpcCalls,
	)

	offrampAddress, err := offRampReader.Address(ctx)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get offramp reader address: %w", err)
	}

	tokenPoolBatchedReader, err := batchreader.NewEVMTokenPoolBatchedReader(execLggr, sourceChainSelector, offrampAddress, batchCaller)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("new token pool batched reader: %w", err)
	}

	chainHealthcheck := cache.NewObservedChainHealthCheck(
		cache.NewChainHealthcheck(
			// Adding more details to Logger to make healthcheck logs more informative
			// It's safe because healthcheck logs only in case of unhealthy state
			lggr.With(
				"onramp", params.offRampConfig.OnRamp,
				"commitStore", params.offRampConfig.CommitStore,
				"offramp", offrampAddress,
			),
			onRampReader,
			commitStoreReader,
		),
		ccip.ExecPluginLabel,
		sourceChainID,
		destChainID,
		params.offRampConfig.OnRamp,
	)

	onchainConfig, err := offRampReader.OnchainConfig(ctx)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get onchain config from offramp reader: %w", err)
	}

	tokenBackgroundWorker := tokendata.NewBackgroundWorker(
		tokenDataProviders,
		numTokenDataWorkers,
		5*time.Second,
		onchainConfig.PermissionLessExecutionThresholdSeconds,
	)
	return &ExecutionPluginStaticConfig{
			lggr:                        execLggr,
			onRampReader:                onRampReader,
			commitStoreReader:           commitStoreReader,
			offRampReader:               offRampReader,
			sourcePriceRegistryProvider: ccipdataprovider.NewEvmPriceRegistry(params.sourceChain.LogPoller(), params.sourceChain.Client(), execLggr, ccip.ExecPluginLabel),
			sourceWrappedNativeToken:    cciptypes.Address(sourceWrappedNative.String()),
			destChainSelector:           destChainSelector,
			priceRegistryProvider:       ccipdataprovider.NewEvmPriceRegistry(params.destChain.LogPoller(), params.destChain.Client(), execLggr, ccip.ExecPluginLabel),
			tokenPoolBatchedReader:      tokenPoolBatchedReader,
			tokenDataWorker:             tokenBackgroundWorker,
			metricsCollector:            metricsCollector,
			chainHealthcheck:            chainHealthcheck,
		}, &ccipcommon.BackfillArgs{
			SourceLP:         params.sourceChain.LogPoller(),
			DestLP:           params.destChain.LogPoller(),
			SourceStartBlock: params.pluginConfig.SourceStartBlock,
			DestStartBlock:   params.pluginConfig.DestStartBlock,
		},
		chainHealthcheck,
		tokenBackgroundWorker,
		nil
}

type jobSpecParams struct {
	pluginConfig  ccipconfig.ExecutionPluginJobSpecConfig
	offRampConfig cciptypes.OffRampStaticConfig
	offRampReader ccipdata.OffRampReader
	sourceChain   legacyevm.Chain
	destChain     legacyevm.Chain
}

func extractJobSpecParams(lggr logger.Logger, jb job.Job, chainSet legacyevm.LegacyChainContainer, registerFilters bool) (*jobSpecParams, error) {
	if jb.OCR2OracleSpec == nil {
		return nil, errors.New("spec is nil")
	}
	spec := jb.OCR2OracleSpec
	var pluginConfig ccipconfig.ExecutionPluginJobSpecConfig
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
		return nil, errors.Wrap(err, "create offRampReader")
	}

	offRampConfig, err := offRampReader.GetStaticConfig(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "get offRamp static config")
	}

	chainID, err := chainselectors.ChainIdFromSelector(offRampConfig.SourceChainSelector)
	if err != nil {
		return nil, err
	}

	sourceChain, err := chainSet.Get(strconv.FormatUint(chainID, 10))
	if err != nil {
		return nil, errors.Wrap(err, "open source chain")
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
