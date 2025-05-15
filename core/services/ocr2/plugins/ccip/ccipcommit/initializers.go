package ccipcommit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"go.uber.org/multierr"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	cciporm "github.com/smartcontractkit/chainlink/v2/core/services/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/cache"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/factory"
	db "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdb"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/observability"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/oraclelib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/promwrapper"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
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
	jb job.Job,
	lggr logger.Logger,
	pr pipeline.Runner,
	argsNoPlugin libocr2.OCR2OracleArgs,
	newInstance bool,
	sourceChainID int64,
	destChainID int64,
	logError func(string),
	pluginJobSpecConfig ccipconfig.CommitPluginJobSpecConfig,
	relayGetter RelayGetter,
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

	if sourceChainID < 0 || destChainID < 0 {
		return nil, errors.New("source and dest chain IDs must be positive")
	}
	srcChain, ok := chainselectors.ChainByEvmChainID(uint64(sourceChainID))
	if !ok {
		return nil, fmt.Errorf("failed to get source chain by evm ID %d", destChainID)
	}
	dstChain, ok2 := chainselectors.ChainByEvmChainID(uint64(destChainID))
	if !ok2 {
		return nil, fmt.Errorf("failed to get dest chain by evm ID %d", destChainID)
	}

	priceGetter, err := initCommitPriceGetter(ctx, lggr, pluginJobSpecConfig, jb, sourceNative,
		pr, relayGetter, srcChain.Selector, dstChain.Selector)
	if err != nil {
		return nil, fmt.Errorf("failed to create price getter: %w", err)
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

	// --------------------------------------------------------------------------------
	// Backwards compatibility for old job spec price getter dynamic config.
	// Should be removed after all jobSpecs migrate to the new format.
	dynamicPriceGetter, is := priceGetter.(*pricegetter.DynamicPriceGetter)
	if is {
		sourceNativeEvmAddr, err2 := ccipcalc.GenericAddrToEvm(sourceNative)
		if err2 != nil {
			return nil, fmt.Errorf("convert source native token address %s to evm address: %w", sourceNative, err2)
		}
		err = dynamicPriceGetter.MoveDeprecatedFields(srcChain.Selector, dstChain.Selector, sourceNativeEvmAddr)
		if err != nil {
			return nil, fmt.Errorf("move deprecated fields: %w", err)
		}
	}
	// --------------------------------------------------------------------------------

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

func initCommitPriceGetter(
	ctx context.Context,
	lggr logger.Logger,
	pluginJobSpecConfig ccipconfig.CommitPluginJobSpecConfig,
	jb job.Job,
	sourceNativeTokenAddr cciptypes.Address,
	pipelineRunner pipeline.Runner,
	relayGetter RelayGetter,
	sourceChainSelector uint64,
	destChainSelector uint64,
) (priceGetter ccip.AllTokensPriceGetter, err error) {
	spec := jb.OCR2OracleSpec
	withPipeline := strings.Trim(pluginJobSpecConfig.TokenPricesUSDPipeline, "\n\t ") != ""
	if withPipeline {
		priceGetter, err = ccip.NewPipelineGetter(
			pluginJobSpecConfig.TokenPricesUSDPipeline,
			pipelineRunner,
			jb.ID,
			jb.ExternalJobID,
			jb.Name.ValueOrZero(),
			lggr,
			sourceNativeTokenAddr,
			sourceChainSelector,
			destChainSelector,
		)
		if err != nil {
			return nil, fmt.Errorf("creating pipeline price getter: %w", err)
		}
	} else {
		// Use dynamic price getter.
		if pluginJobSpecConfig.PriceGetterConfig == nil {
			return nil, errors.New("priceGetterConfig is nil")
		}

		// Configure contract readers for all chains specified in the aggregator configurations.
		// Some lanes (e.g. Wemix/Kroma) requires other clients than source and destination, since they use feeds from other chains.
		aggregatorChainsToContracts := make(map[uint64][]common.Address)
		for _, aggCfg := range pluginJobSpecConfig.PriceGetterConfig.AggregatorPrices {
			if _, ok := aggregatorChainsToContracts[aggCfg.ChainID]; !ok {
				aggregatorChainsToContracts[aggCfg.ChainID] = make([]common.Address, 0)
			}

			aggregatorChainsToContracts[aggCfg.ChainID] = append(aggregatorChainsToContracts[aggCfg.ChainID], aggCfg.AggregatorContractAddress)
		}

		for _, priceCfg := range pluginJobSpecConfig.PriceGetterConfig.TokenPrices {
			if priceCfg.AggregatorConfig == nil {
				continue
			}
			aggCfg := *priceCfg.AggregatorConfig
			contractAddrs, ok := aggregatorChainsToContracts[aggCfg.ChainID]
			if !ok {
				aggregatorChainsToContracts[aggCfg.ChainID] = make([]common.Address, 0)
			}
			if !slices.Contains(contractAddrs, aggCfg.AggregatorContractAddress) {
				aggregatorChainsToContracts[aggCfg.ChainID] = append(aggregatorChainsToContracts[aggCfg.ChainID],
					aggCfg.AggregatorContractAddress)
			}
		}

		contractReaders := map[uint64]commontypes.ContractReader{}

		for chainID, aggregatorContracts := range aggregatorChainsToContracts {
			relayID := commontypes.RelayID{Network: spec.Relay, ChainID: strconv.FormatUint(chainID, 10)}
			relay, rerr := relayGetter.Get(relayID)
			if rerr != nil {
				return nil, fmt.Errorf("get relay by id=%v: %w", relayID, err)
			}

			contractsConfig := make(map[string]evmrelaytypes.ChainContractReader, len(aggregatorContracts))
			for i := range aggregatorContracts {
				contractsConfig[fmt.Sprintf("%v_%v", ccip.OffchainAggregator, i)] = evmrelaytypes.ChainContractReader{
					ContractABI: ccip.OffChainAggregatorABI,
					Configs: map[string]*evmrelaytypes.ChainReaderDefinition{
						"decimals": { // CR consumers choose an alias
							ChainSpecificName: "decimals",
						},
						"latestRoundData": {
							ChainSpecificName: "latestRoundData",
						},
					},
				}
			}
			contractReaderConfig := evmrelaytypes.ChainReaderConfig{
				Contracts: contractsConfig,
			}

			contractReaderConfigJSONBytes, jerr := json.Marshal(contractReaderConfig)
			if jerr != nil {
				return nil, fmt.Errorf("marshal contract reader config: %w", jerr)
			}

			contractReader, cerr := relay.NewContractReader(ctx, contractReaderConfigJSONBytes)
			if cerr != nil {
				return nil, fmt.Errorf("new ccip commit contract reader %w", cerr)
			}

			contractReaders[chainID] = contractReader
		}

		priceGetter, err = ccip.NewDynamicPriceGetter(*pluginJobSpecConfig.PriceGetterConfig, contractReaders)
		if err != nil {
			return nil, fmt.Errorf("creating dynamic price getter: %w", err)
		}
	}
	return priceGetter, nil
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

type RelayGetter interface {
	Get(id commontypes.RelayID) (loop.Relayer, error)
	GetIDToRelayerMap() map[commontypes.RelayID]loop.Relayer
}
