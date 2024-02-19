package ccipcommit

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"go.uber.org/multierr"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/ccipdataprovider"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/rpclib"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/factory"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/observability"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/oraclelib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/promwrapper"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

func NewCommitServices(ctx context.Context, lggr logger.Logger, jb job.Job, chainSet legacyevm.LegacyChainContainer, new bool, pr pipeline.Runner, argsNoPlugin libocr2.OCR2OracleArgs, logError func(string), qopts ...pg.QOpt) ([]job.ServiceCtx, error) {
	pluginConfig, backfillArgs, err := jobSpecToCommitPluginConfig(lggr, jb, pr, chainSet, qopts...)
	if err != nil {
		return nil, err
	}
	wrappedPluginFactory := NewCommitReportingPluginFactory(*pluginConfig)
	destChainID, err := chainselectors.ChainIdFromSelector(pluginConfig.destChainSelector)
	if err != nil {
		return nil, err
	}
	argsNoPlugin.ReportingPluginFactory = promwrapper.NewPromFactory(wrappedPluginFactory, "CCIPCommit", jb.OCR2OracleSpec.Relay, big.NewInt(0).SetUint64(destChainID))
	argsNoPlugin.Logger = commonlogger.NewOCRWrapper(pluginConfig.lggr, true, logError)
	oracle, err := libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		return nil, err
	}
	// If this is a brand-new job, then we make use of the start blocks. If not then we're rebooting and log poller will pick up where we left off.
	if new {
		return []job.ServiceCtx{oraclelib.NewBackfilledOracle(
			pluginConfig.lggr,
			backfillArgs.SourceLP,
			backfillArgs.DestLP,
			backfillArgs.SourceStartBlock,
			backfillArgs.DestStartBlock,
			job.NewServiceAdapter(oracle)),
		}, nil
	}
	return []job.ServiceCtx{job.NewServiceAdapter(oracle)}, nil
}

func CommitReportToEthTxMeta(typ ccipconfig.ContractType, ver semver.Version) (func(report []byte) (*txmgr.TxMeta, error), error) {
	return factory.CommitReportToEthTxMeta(typ, ver)
}

// UnregisterCommitPluginLpFilters unregisters all the registered filters for both source and dest chains.
// NOTE: The transaction MUST be used here for CLO's monster tx to function as expected
// https://github.com/smartcontractkit/ccip/blob/68e2197472fb017dd4e5630d21e7878d58bc2a44/core/services/feeds/service.go#L716
// TODO once that transaction is broken up, we should be able to simply rely on oracle.Close() to cleanup the filters.
// Until then we have to deterministically reload the readers from the spec (and thus their filters) and close them.
func UnregisterCommitPluginLpFilters(ctx context.Context, lggr logger.Logger, jb job.Job, chainSet legacyevm.LegacyChainContainer, qopts ...pg.QOpt) error {
	params, err := extractJobSpecParams(jb, chainSet)
	if err != nil {
		return err
	}
	versionFinder := factory.NewEvmVersionFinder()
	unregisterFuncs := []func() error{
		func() error {
			return factory.CloseCommitStoreReader(lggr, versionFinder, params.commitStoreAddress, params.destChain.Client(), params.destChain.LogPoller(), params.sourceChain.GasEstimator(), qopts...)
		},
		func() error {
			return factory.CloseOnRampReader(lggr, versionFinder, params.commitStoreStaticCfg.SourceChainSelector, params.commitStoreStaticCfg.ChainSelector, cciptypes.Address(params.commitStoreStaticCfg.OnRamp.String()), params.sourceChain.LogPoller(), params.sourceChain.Client(), qopts...)
		},
		func() error {
			return factory.CloseOffRampReader(lggr, versionFinder, params.pluginConfig.OffRamp, params.destChain.Client(), params.destChain.LogPoller(), params.destChain.GasEstimator(), qopts...)
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

func jobSpecToCommitPluginConfig(lggr logger.Logger, jb job.Job, pr pipeline.Runner, chainSet legacyevm.LegacyChainContainer, qopts ...pg.QOpt) (*CommitPluginStaticConfig, *ccipcommon.BackfillArgs, error) {
	params, err := extractJobSpecParams(jb, chainSet)
	if err != nil {
		return nil, nil, err
	}

	versionFinder := factory.NewEvmVersionFinder()
	commitStoreReader, err := factory.NewCommitStoreReader(lggr, versionFinder, params.commitStoreAddress, params.destChain.Client(), params.destChain.LogPoller(), params.sourceChain.GasEstimator(), qopts...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create commitStore reader")
	}
	sourceChainName, destChainName, err := ccipconfig.ResolveChainNames(params.sourceChain.ID().Int64(), params.destChain.ID().Int64())
	if err != nil {
		return nil, nil, err
	}
	commitLggr := lggr.Named("CCIPCommit").With("sourceChain", sourceChainName, "destChain", destChainName)

	var priceGetter pricegetter.PriceGetter
	withPipeline := strings.Trim(params.pluginConfig.TokenPricesUSDPipeline, "\n\t ") != ""
	if withPipeline {
		priceGetter, err = pricegetter.NewPipelineGetter(params.pluginConfig.TokenPricesUSDPipeline, pr, jb.ID, jb.ExternalJobID, jb.Name.ValueOrZero(), lggr)
		if err != nil {
			return nil, nil, fmt.Errorf("creating pipeline price getter: %w", err)
		}
	} else {
		// Use dynamic price getter.
		if params.pluginConfig.PriceGetterConfig == nil {
			return nil, nil, fmt.Errorf("priceGetterConfig is nil")
		}

		// Build price getter clients for all chains specified in the aggregator configurations.
		// Some lanes (e.g. Wemix/Kroma) requires other clients than source and destination, since they use feeds from other chains.
		priceGetterClients := map[uint64]pricegetter.DynamicPriceGetterClient{}
		for _, aggCfg := range params.pluginConfig.PriceGetterConfig.AggregatorPrices {
			chainID := aggCfg.ChainID
			// Retrieve the chain.
			chain, _, err2 := ccipconfig.GetChainByChainID(chainSet, chainID)
			if err2 != nil {
				return nil, nil, fmt.Errorf("retrieving chain for chainID %d: %w", chainID, err2)
			}
			caller := rpclib.NewDynamicLimitedBatchCaller(
				lggr,
				chain.Client(),
				rpclib.DefaultRpcBatchSizeLimit,
				rpclib.DefaultRpcBatchBackOffMultiplier,
			)
			priceGetterClients[chainID] = pricegetter.NewDynamicPriceGetterClient(caller)
		}

		priceGetter, err = pricegetter.NewDynamicPriceGetter(*params.pluginConfig.PriceGetterConfig, priceGetterClients)
		if err != nil {
			return nil, nil, fmt.Errorf("creating dynamic price getter: %w", err)
		}
	}

	// Load all the readers relevant for this plugin.
	onRampReader, err := factory.NewOnRampReader(commitLggr, versionFinder, params.commitStoreStaticCfg.SourceChainSelector, params.commitStoreStaticCfg.ChainSelector, cciptypes.Address(params.commitStoreStaticCfg.OnRamp.String()), params.sourceChain.LogPoller(), params.sourceChain.Client(), qopts...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed onramp reader")
	}
	offRampReader, err := factory.NewOffRampReader(commitLggr, versionFinder, params.pluginConfig.OffRamp, params.destChain.Client(), params.destChain.LogPoller(), params.destChain.GasEstimator(), true, qopts...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed offramp reader")
	}
	onRampRouterAddr, err := onRampReader.RouterAddress()
	if err != nil {
		return nil, nil, err
	}
	routerAddr, err := ccipcalc.GenericAddrToEvm(onRampRouterAddr)
	if err != nil {
		return nil, nil, err
	}
	sourceRouter, err := router.NewRouter(routerAddr, params.sourceChain.Client())
	if err != nil {
		return nil, nil, err
	}
	sourceNative, err := sourceRouter.GetWrappedNative(nil)
	if err != nil {
		return nil, nil, err
	}

	// Prom wrappers
	onRampReader = observability.NewObservedOnRampReader(onRampReader, params.sourceChain.ID().Int64(), ccip.CommitPluginLabel)
	offRampReader = observability.NewObservedOffRampReader(offRampReader, params.destChain.ID().Int64(), ccip.CommitPluginLabel)
	commitStoreReader = observability.NewObservedCommitStoreReader(commitStoreReader, params.destChain.ID().Int64(), ccip.CommitPluginLabel)
	metricsCollector := ccip.NewPluginMetricsCollector(ccip.CommitPluginLabel, params.sourceChain.ID().Int64(), params.destChain.ID().Int64())

	lggr.Infow("NewCommitServices",
		"pluginConfig", params.pluginConfig,
		"staticConfig", params.commitStoreStaticCfg,
		// TODO bring back
		//"dynamicOnRampConfig", dynamicOnRampConfig,
		"sourceNative", sourceNative,
		"sourceRouter", sourceRouter.Address())
	return &CommitPluginStaticConfig{
			lggr:                  commitLggr,
			onRampReader:          onRampReader,
			offRamp:               offRampReader,
			sourceNative:          ccipcalc.EvmAddrToGeneric(sourceNative),
			priceGetter:           priceGetter,
			sourceChainSelector:   params.commitStoreStaticCfg.SourceChainSelector,
			destChainSelector:     params.commitStoreStaticCfg.ChainSelector,
			commitStore:           commitStoreReader,
			priceRegistryProvider: ccipdataprovider.NewEvmPriceRegistry(params.destChain.LogPoller(), params.destChain.Client(), commitLggr, ccip.CommitPluginLabel),
			metricsCollector:      metricsCollector,
		}, &ccipcommon.BackfillArgs{
			SourceLP:         params.sourceChain.LogPoller(),
			DestLP:           params.destChain.LogPoller(),
			SourceStartBlock: params.pluginConfig.SourceStartBlock,
			DestStartBlock:   params.pluginConfig.DestStartBlock,
		}, nil
}

type jobSpecParams struct {
	pluginConfig         ccipconfig.CommitPluginJobSpecConfig
	commitStoreAddress   cciptypes.Address
	commitStoreStaticCfg commit_store.CommitStoreStaticConfig
	sourceChain          legacyevm.Chain
	destChain            legacyevm.Chain
}

func extractJobSpecParams(jb job.Job, chainSet legacyevm.LegacyChainContainer) (*jobSpecParams, error) {
	if jb.OCR2OracleSpec == nil {
		return nil, errors.New("spec is nil")
	}
	spec := jb.OCR2OracleSpec

	var pluginConfig ccipconfig.CommitPluginJobSpecConfig
	err := json.Unmarshal(spec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, err
	}

	destChain, _, err := ccipconfig.GetChainFromSpec(spec, chainSet)
	if err != nil {
		return nil, err
	}

	commitStoreAddress := common.HexToAddress(spec.ContractID)
	staticConfig, err := ccipdata.FetchCommitStoreStaticConfig(commitStoreAddress, destChain.Client())
	if err != nil {
		return nil, fmt.Errorf("get commit store static config: %w", err)
	}

	sourceChain, _, err := ccipconfig.GetChainByChainSelector(chainSet, staticConfig.SourceChainSelector)
	if err != nil {
		return nil, err
	}

	return &jobSpecParams{
		pluginConfig:         pluginConfig,
		commitStoreAddress:   ccipcalc.EvmAddrToGeneric(commitStoreAddress),
		commitStoreStaticCfg: staticConfig,
		sourceChain:          sourceChain,
		destChain:            destChain,
	}, nil
}
