package ccip

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	relaylogger "github.com/smartcontractkit/chainlink-relay/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/contractutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/oraclelib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/pricegetter"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/promwrapper"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

const (
	COMMIT_PRICE_UPDATES = "Commit price updates"
)

type BackfillArgs struct {
	sourceLP, destLP                 logpoller.LogPoller
	sourceStartBlock, destStartBlock int64
}

func jobSpecToCommitPluginConfig(lggr logger.Logger, jb job.Job, pr pipeline.Runner, chainSet evm.LegacyChainContainer) (*CommitPluginConfig, *BackfillArgs, error) {
	if jb.OCR2OracleSpec == nil {
		return nil, nil, errors.New("spec is nil")
	}
	spec := jb.OCR2OracleSpec
	var pluginConfig ccipconfig.CommitPluginJobSpecConfig
	err := json.Unmarshal(spec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, nil, err
	}
	chainIDInterface, ok := spec.RelayConfig["chainID"]
	if !ok {
		return nil, nil, errors.New("chainID must be provided in relay config")
	}
	destChainID := int64(chainIDInterface.(float64))
	destChain, err := chainSet.Get(strconv.FormatInt(destChainID, 10))
	if err != nil {
		return nil, nil, errors.Wrap(err, "get chainset")
	}
	commitStore, err := contractutil.LoadCommitStore(common.HexToAddress(spec.ContractID), CommitPluginLabel, destChain.Client())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed loading commitStore")
	}
	staticConfig, err := commitStore.GetStaticConfig(&bind.CallOpts{})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed getting the static config from the commitStore")
	}
	chainId, err := chainselectors.ChainIdFromSelector(staticConfig.SourceChainSelector)
	if err != nil {
		return nil, nil, err
	}
	sourceChain, err := chainSet.Get(strconv.FormatUint(chainId, 10))
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to open source chain")
	}
	offRamp, err := contractutil.LoadOffRamp(common.HexToAddress(pluginConfig.OffRamp), CommitPluginLabel, destChain.Client())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed loading offRamp")
	}
	commitLggr := lggr.Named("CCIPCommit").With(
		"sourceChain", ChainName(int64(chainId)),
		"destChain", ChainName(destChainID))
	onRampReader, err := ccipdata.NewOnRampReader(commitLggr, staticConfig.SourceChainSelector, staticConfig.ChainSelector, staticConfig.OnRamp, sourceChain.LogPoller(), sourceChain.Client(), sourceChain.Config().EVM().FinalityTagEnabled())
	if err != nil {
		return nil, nil, err
	}
	priceGetterObject, err := pricegetter.NewPipelineGetter(pluginConfig.TokenPricesUSDPipeline, pr, jb.ID, jb.ExternalJobID, jb.Name.ValueOrZero(), lggr)
	if err != nil {
		return nil, nil, err
	}
	sourceRouter, err := router.NewRouter(onRampReader.RouterAddress(), sourceChain.Client())
	if err != nil {
		return nil, nil, err
	}
	sourceNative, err := sourceRouter.GetWrappedNative(nil)
	if err != nil {
		return nil, nil, err
	}
	lggr.Infow("NewCommitServices",
		"pluginConfig", pluginConfig,
		"staticConfig", staticConfig,
		// TODO bring back
		//"dynamicOnRampConfig", dynamicOnRampConfig,
		"sourceNative", sourceNative,
		"sourceRouter", sourceRouter.Address())
	return &CommitPluginConfig{
			lggr:                commitLggr,
			destLP:              destChain.LogPoller(),
			onRampReader:        onRampReader,
			destReader:          ccipdata.NewLogPollerReader(destChain.LogPoller(), commitLggr, destChain.Client()),
			offRamp:             offRamp,
			priceGetter:         priceGetterObject,
			sourceNative:        sourceNative,
			sourceFeeEstimator:  sourceChain.GasEstimator(),
			sourceChainSelector: staticConfig.SourceChainSelector,
			destClient:          destChain.Client(),
			commitStore:         commitStore,
		}, &BackfillArgs{
			sourceLP:         sourceChain.LogPoller(),
			destLP:           destChain.LogPoller(),
			sourceStartBlock: pluginConfig.SourceStartBlock,
			destStartBlock:   pluginConfig.DestStartBlock,
		}, nil
}

func NewCommitServices(lggr logger.Logger, jb job.Job, chainSet evm.LegacyChainContainer, new bool, pr pipeline.Runner, argsNoPlugin libocr2.OCR2OracleArgs, logError func(string), qopts ...pg.QOpt) ([]job.ServiceCtx, error) {
	pluginConfig, backfillArgs, err := jobSpecToCommitPluginConfig(lggr, jb, pr, chainSet)
	if err != nil {
		return nil, err
	}
	wrappedPluginFactory := NewCommitReportingPluginFactory(*pluginConfig)
	err = wrappedPluginFactory.UpdateLogPollerFilters(utils.ZeroAddress, qopts...)
	if err != nil {
		return nil, err
	}

	argsNoPlugin.ReportingPluginFactory = promwrapper.NewPromFactory(wrappedPluginFactory, "CCIPCommit", jb.OCR2OracleSpec.Relay, pluginConfig.destChainEVMID)
	argsNoPlugin.Logger = relaylogger.NewOCRWrapper(pluginConfig.lggr, true, logError)
	oracle, err := libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		return nil, err
	}
	// If this is a brand-new job, then we make use of the start blocks. If not then we're rebooting and log poller will pick up where we left off.
	if new {
		return []job.ServiceCtx{oraclelib.NewBackfilledOracle(
			pluginConfig.lggr,
			backfillArgs.sourceLP,
			backfillArgs.destLP,
			backfillArgs.sourceStartBlock,
			backfillArgs.destStartBlock,
			job.NewServiceAdapter(oracle)),
		}, nil
	}
	return []job.ServiceCtx{job.NewServiceAdapter(oracle)}, nil
}

// CommitReportToEthTxMeta generates a txmgr.EthTxMeta from the given commit report.
// sequence numbers of the committed messages will be added to tx metadata
func CommitReportToEthTxMeta(report []byte) (*txmgr.TxMeta, error) {
	commitReport, err := abihelpers.DecodeCommitReport(report)
	if err != nil {
		return nil, err
	}
	n := int(commitReport.Interval.Max-commitReport.Interval.Min) + 1
	seqRange := make([]uint64, n)
	for i := 0; i < n; i++ {
		seqRange[i] = uint64(i) + commitReport.Interval.Min
	}
	return &txmgr.TxMeta{
		SeqNumbers: seqRange,
	}, nil
}

func getCommitPluginDestLpFilters(priceRegistry common.Address, offRamp common.Address) []logpoller.Filter {
	return []logpoller.Filter{
		{
			Name:      logpoller.FilterName(COMMIT_PRICE_UPDATES, priceRegistry.String()),
			EventSigs: []common.Hash{abihelpers.EventSignatures.UsdPerUnitGasUpdated, abihelpers.EventSignatures.UsdPerTokenUpdated},
			Addresses: []common.Address{priceRegistry},
		},
		{
			Name:      logpoller.FilterName(FEE_TOKEN_ADDED, priceRegistry),
			EventSigs: []common.Hash{abihelpers.EventSignatures.FeeTokenAdded},
			Addresses: []common.Address{priceRegistry},
		},
		{
			Name:      logpoller.FilterName(FEE_TOKEN_REMOVED, priceRegistry),
			EventSigs: []common.Hash{abihelpers.EventSignatures.FeeTokenRemoved},
			Addresses: []common.Address{priceRegistry},
		},
		{
			Name:      logpoller.FilterName(EXEC_TOKEN_POOL_ADDED, offRamp),
			EventSigs: []common.Hash{abihelpers.EventSignatures.PoolAdded},
			Addresses: []common.Address{offRamp},
		},
		{
			Name:      logpoller.FilterName(EXEC_TOKEN_POOL_REMOVED, offRamp),
			EventSigs: []common.Hash{abihelpers.EventSignatures.PoolRemoved},
			Addresses: []common.Address{offRamp},
		},
	}
}

// UnregisterCommitPluginLpFilters unregisters all the registered filters for both source and dest chains.
func UnregisterCommitPluginLpFilters(ctx context.Context, lggr logger.Logger, jb job.Job, pr pipeline.Runner, chainSet evm.LegacyChainContainer, qopts ...pg.QOpt) error {
	commitPluginConfig, _, err := jobSpecToCommitPluginConfig(lggr, jb, pr, chainSet)
	if err != nil {
		return errors.New("spec is nil")
	}
	if err := commitPluginConfig.onRampReader.Close(); err != nil {
		return err
	}

	// TODO: once offramp/commit are abstracted, we can call Close on the offramp/commit readers to unregister filters.
	return unregisterCommitPluginFilters(ctx, commitPluginConfig.destLP, commitPluginConfig.commitStore,
		commitPluginConfig.offRamp.Address(), qopts...)
}

func unregisterCommitPluginFilters(ctx context.Context, destLP logpoller.LogPoller, destCommitStore commit_store.CommitStoreInterface, offRamp common.Address, qopts ...pg.QOpt) error {
	dynamicCfg, err := destCommitStore.GetDynamicConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return err
	}
	return logpollerutil.UnregisterLpFilters(
		destLP,
		getCommitPluginDestLpFilters(dynamicCfg.PriceRegistry, offRamp),
		qopts...,
	)
}
