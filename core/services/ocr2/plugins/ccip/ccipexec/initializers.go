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
	chainselectors "github.com/smartcontractkit/chain-selectors"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"go.uber.org/multierr"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const numTokenDataWorkers = 5

func NewExecutionServices(ctx context.Context, lggr logger.Logger, jb job.Job, chainSet legacyevm.LegacyChainContainer, new bool, argsNoPlugin libocr2.OCR2OracleArgs, logError func(string), qopts ...pg.QOpt) ([]job.ServiceCtx, error) {
	execPluginConfig, backfillArgs, err := jobSpecToExecPluginConfig(ctx, lggr, jb, chainSet, qopts...)
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
				job.NewServiceAdapter(oracle)),
		}, nil
	}
	return []job.ServiceCtx{job.NewServiceAdapter(oracle)}, nil
}

// UnregisterExecPluginLpFilters unregisters all the registered filters for both source and dest chains.
// See comment in UnregisterCommitPluginLpFilters
// It MUST mirror the filters registered in NewExecutionServices.
func UnregisterExecPluginLpFilters(ctx context.Context, lggr logger.Logger, jb job.Job, chainSet legacyevm.LegacyChainContainer, qopts ...pg.QOpt) error {
	params, err := extractJobSpecParams(lggr, jb, chainSet, false, qopts...)
	if err != nil {
		return err
	}

	versionFinder := factory.NewEvmVersionFinder()
	unregisterFuncs := []func() error{
		func() error {
			return factory.CloseCommitStoreReader(lggr, versionFinder, params.offRampConfig.CommitStore, params.destChain.Client(), params.destChain.LogPoller(), params.sourceChain.GasEstimator(), qopts...)
		},
		func() error {
			return factory.CloseOnRampReader(lggr, versionFinder, params.offRampConfig.SourceChainSelector, params.offRampConfig.ChainSelector, params.offRampConfig.OnRamp, params.sourceChain.LogPoller(), params.sourceChain.Client(), qopts...)
		},
		func() error {
			return factory.CloseOffRampReader(lggr, versionFinder, params.offRampReader.Address(), params.destChain.Client(), params.destChain.LogPoller(), params.destChain.GasEstimator(), qopts...)
		},
		func() error { // usdc token data reader
			if usdcDisabled := params.pluginConfig.USDCConfig.AttestationAPI == ""; usdcDisabled {
				return nil
			}
			return ccipdata.CloseUSDCReader(lggr, jobIDToString(jb.ID), params.pluginConfig.USDCConfig.SourceMessageTransmitterAddress, params.sourceChain.LogPoller(), qopts...)
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
func ExecReportToEthTxMeta(typ ccipconfig.ContractType, ver semver.Version) (func(report []byte) (*txmgr.TxMeta, error), error) {
	return factory.ExecReportToEthTxMeta(typ, ver)
}

func initTokenDataProviders(lggr logger.Logger, jobID string, pluginConfig ccipconfig.ExecutionPluginJobSpecConfig, sourceLP logpoller.LogPoller, qopts ...pg.QOpt) (map[cciptypes.Address]tokendata.Reader, error) {
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
				pluginConfig.USDCConfig.AttestationAPITimeoutSeconds,
			)
	}

	return tokenDataProviders, nil
}

func jobSpecToExecPluginConfig(ctx context.Context, lggr logger.Logger, jb job.Job, chainSet legacyevm.LegacyChainContainer, qopts ...pg.QOpt) (*ExecutionPluginStaticConfig, *ccipcommon.BackfillArgs, error) {
	params, err := extractJobSpecParams(lggr, jb, chainSet, true, qopts...)
	if err != nil {
		return nil, nil, err
	}

	sourceChainID := params.sourceChain.ID().Int64()
	destChainID := params.destChain.ID().Int64()
	versionFinder := factory.NewEvmVersionFinder()

	sourceChainName, destChainName, err := ccipconfig.ResolveChainNames(sourceChainID, destChainID)
	if err != nil {
		return nil, nil, err
	}
	execLggr := lggr.Named("CCIPExecution").With("sourceChain", sourceChainName, "destChain", destChainName)
	onRampReader, err := factory.NewOnRampReader(execLggr, versionFinder, params.offRampConfig.SourceChainSelector, params.offRampConfig.ChainSelector, params.offRampConfig.OnRamp, params.sourceChain.LogPoller(), params.sourceChain.Client(), qopts...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "create onramp reader")
	}
	dynamicOnRampConfig, err := onRampReader.GetDynamicConfig()
	if err != nil {
		return nil, nil, errors.Wrap(err, "get onramp dynamic config")
	}

	routerAddr, err := ccipcalc.GenericAddrToEvm(dynamicOnRampConfig.Router)
	if err != nil {
		return nil, nil, err
	}
	sourceRouter, err := router.NewRouter(routerAddr, params.sourceChain.Client())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed loading source router")
	}
	sourceWrappedNative, err := sourceRouter.GetWrappedNative(&bind.CallOpts{})
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get source native token")
	}

	// TODO: we don't support onramp source registry changes without a reboot yet?
	sourcePriceRegistry, err := factory.NewPriceRegistryReader(lggr, versionFinder, dynamicOnRampConfig.PriceRegistry, params.sourceChain.LogPoller(), params.sourceChain.Client())
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not load source registry")
	}

	commitStoreReader, err := factory.NewCommitStoreReader(lggr, versionFinder, params.offRampConfig.CommitStore, params.destChain.Client(), params.destChain.LogPoller(), params.sourceChain.GasEstimator(), qopts...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not load commitStoreReader reader")
	}

	tokenDataProviders, err := initTokenDataProviders(lggr, jobIDToString(jb.ID), params.pluginConfig, params.sourceChain.LogPoller(), qopts...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get token data providers")
	}

	// Prom wrappers
	onRampReader = observability.NewObservedOnRampReader(onRampReader, sourceChainID, ccip.ExecPluginLabel)
	sourcePriceRegistry = observability.NewPriceRegistryReader(sourcePriceRegistry, sourceChainID, ccip.ExecPluginLabel)
	commitStoreReader = observability.NewObservedCommitStoreReader(commitStoreReader, destChainID, ccip.ExecPluginLabel)
	offRampReader := observability.NewObservedOffRampReader(params.offRampReader, destChainID, ccip.ExecPluginLabel)
	metricsCollector := ccip.NewPluginMetricsCollector(ccip.ExecPluginLabel, sourceChainID, destChainID)

	destChainSelector, err := chainselectors.SelectorFromChainId(uint64(destChainID))
	if err != nil {
		return nil, nil, fmt.Errorf("get chain %d selector: %w", destChainID, err)
	}
	sourceChainSelector, err := chainselectors.SelectorFromChainId(uint64(sourceChainID))
	if err != nil {
		return nil, nil, fmt.Errorf("get chain %d selector: %w", sourceChainID, err)
	}

	execLggr.Infow("Initialized exec plugin",
		"pluginConfig", params.pluginConfig,
		"onRampAddress", params.offRampConfig.OnRamp,
		"sourcePriceRegistry", sourcePriceRegistry.Address(),
		"dynamicOnRampConfig", dynamicOnRampConfig,
		"sourceNative", sourceWrappedNative,
		"sourceRouter", sourceRouter.Address())

	batchCaller := rpclib.NewDynamicLimitedBatchCaller(lggr, params.destChain.Client(), rpclib.DefaultRpcBatchSizeLimit, rpclib.DefaultRpcBatchBackOffMultiplier)

	tokenPoolBatchedReader, err := batchreader.NewEVMTokenPoolBatchedReader(execLggr, sourceChainSelector, offRampReader.Address(), batchCaller)
	if err != nil {
		return nil, nil, fmt.Errorf("new token pool batched reader: %w", err)
	}

	return &ExecutionPluginStaticConfig{
			lggr:                     execLggr,
			onRampReader:             onRampReader,
			commitStoreReader:        commitStoreReader,
			offRampReader:            offRampReader,
			sourcePriceRegistry:      sourcePriceRegistry,
			sourceWrappedNativeToken: cciptypes.Address(sourceWrappedNative.String()),
			destChainSelector:        destChainSelector,
			priceRegistryProvider:    ccipdataprovider.NewEvmPriceRegistry(params.destChain.LogPoller(), params.destChain.Client(), execLggr, ccip.ExecPluginLabel),
			tokenPoolBatchedReader:   tokenPoolBatchedReader,
			tokenDataWorker: tokendata.NewBackgroundWorker(
				ctx,
				tokenDataProviders,
				numTokenDataWorkers,
				5*time.Second,
				offRampReader.OnchainConfig().PermissionLessExecutionThresholdSeconds,
			),
			metricsCollector: metricsCollector,
		}, &ccipcommon.BackfillArgs{
			SourceLP:         params.sourceChain.LogPoller(),
			DestLP:           params.destChain.LogPoller(),
			SourceStartBlock: params.pluginConfig.SourceStartBlock,
			DestStartBlock:   params.pluginConfig.DestStartBlock,
		}, nil
}

type jobSpecParams struct {
	pluginConfig  ccipconfig.ExecutionPluginJobSpecConfig
	offRampConfig cciptypes.OffRampStaticConfig
	offRampReader ccipdata.OffRampReader
	sourceChain   legacyevm.Chain
	destChain     legacyevm.Chain
}

func extractJobSpecParams(lggr logger.Logger, jb job.Job, chainSet legacyevm.LegacyChainContainer, registerFilters bool, qopts ...pg.QOpt) (*jobSpecParams, error) {
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
	offRampReader, err := factory.NewOffRampReader(lggr, versionFinder, offRampAddress, destChain.Client(), destChain.LogPoller(), destChain.GasEstimator(), registerFilters, qopts...)
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
