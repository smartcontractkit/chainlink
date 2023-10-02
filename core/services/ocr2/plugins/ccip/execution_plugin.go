package ccip

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	relaylogger "github.com/smartcontractkit/chainlink-relay/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/contractutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/logpollerutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/oraclelib"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/observability"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata/usdc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/promwrapper"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const (
	EXEC_REPORT_ACCEPTS          = "Exec report accepts"
	EXEC_EXECUTION_STATE_CHANGES = "Exec execution state changes"
	EXEC_TOKEN_POOL_ADDED        = "Token pool added"
	EXEC_TOKEN_POOL_REMOVED      = "Token pool removed"
	FEE_TOKEN_ADDED              = "Fee token added"
	FEE_TOKEN_REMOVED            = "Fee token removed"
)

func jobSpecToExecPluginConfig(lggr logger.Logger, jb job.Job, chainSet evm.LegacyChainContainer) (*ExecutionPluginConfig, *BackfillArgs, error) {
	if jb.OCR2OracleSpec == nil {
		return nil, nil, errors.New("spec is nil")
	}
	spec := jb.OCR2OracleSpec
	var pluginConfig ccipconfig.ExecutionPluginJobSpecConfig
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
	offRamp, _, err := contractutil.LoadOffRamp(common.HexToAddress(spec.ContractID), ExecPluginLabel, destChain.Client())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed loading offRamp")
	}
	offRampConfig, err := offRamp.GetStaticConfig(&bind.CallOpts{})
	if err != nil {
		return nil, nil, err
	}
	chainId, err := chainselectors.ChainIdFromSelector(offRampConfig.SourceChainSelector)
	if err != nil {
		return nil, nil, err
	}
	sourceChain, err := chainSet.Get(strconv.FormatUint(chainId, 10))
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to open source chain")
	}
	commitStore, commitStoreVersion, err := contractutil.LoadCommitStore(offRampConfig.CommitStore, ExecPluginLabel, destChain.Client())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed loading commitStore")
	}
	onRamp, onRampVersion, err := contractutil.LoadOnRamp(offRampConfig.OnRamp, ExecPluginLabel, sourceChain.Client())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed loading onRamp")
	}
	dynamicOnRampConfig, err := contractutil.LoadOnRampDynamicConfig(onRamp, onRampVersion, sourceChain.Client())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed loading onRamp config")
	}
	sourceRouter, err := router.NewRouter(dynamicOnRampConfig.Router, sourceChain.Client())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed loading source router")
	}
	sourceWrappedNative, err := sourceRouter.GetWrappedNative(&bind.CallOpts{})
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get source native token")
	}
	sourcePriceRegistry, err := observability.NewObservedPriceRegistry(dynamicOnRampConfig.PriceRegistry, ExecPluginLabel, sourceChain.Client())
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create source price registry")
	}

	execLggr := lggr.Named("CCIPExecution").With(
		"sourceChain", ChainName(int64(chainId)),
		"destChain", ChainName(destChainID))
	onRampReader, err := ccipdata.NewOnRampReader(execLggr, offRampConfig.SourceChainSelector,
		offRampConfig.ChainSelector, offRampConfig.OnRamp, sourceChain.LogPoller(), sourceChain.Client(), sourceChain.Config().EVM().FinalityTagEnabled())
	if err != nil {
		return nil, nil, err
	}
	tokenDataProviders, err := getTokenDataProviders(lggr, pluginConfig, sourceChain.LogPoller())
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not get token data providers")
	}
	execLggr.Infow("Initialized exec plugin",
		"pluginConfig", pluginConfig,
		"onRampAddress", onRamp.Address(),
		"sourcePriceRegistry", sourcePriceRegistry.Address(),
		"dynamicOnRampConfig", dynamicOnRampConfig,
		"sourceNative", sourceWrappedNative,
		"sourceRouter", sourceRouter.Address())
	return &ExecutionPluginConfig{
			lggr:                     execLggr,
			sourceLP:                 sourceChain.LogPoller(),
			destLP:                   destChain.LogPoller(),
			onRampReader:             onRampReader,
			destReader:               ccipdata.NewLogPollerReader(destChain.LogPoller(), execLggr, destChain.Client()),
			onRamp:                   onRamp,
			onRampVersion:            onRampVersion,
			offRamp:                  offRamp,
			commitStore:              commitStore,
			commitStoreVersion:       commitStoreVersion,
			sourcePriceRegistry:      sourcePriceRegistry,
			sourceWrappedNativeToken: sourceWrappedNative,
			destClient:               destChain.Client(),
			sourceClient:             sourceChain.Client(),
			destGasEstimator:         destChain.GasEstimator(),
			destChainEVMID:           destChain.ID(),
			tokenDataProviders:       tokenDataProviders,
		}, &BackfillArgs{
			sourceLP:         sourceChain.LogPoller(),
			destLP:           destChain.LogPoller(),
			sourceStartBlock: pluginConfig.SourceStartBlock,
			destStartBlock:   pluginConfig.DestStartBlock,
		}, nil
}

func NewExecutionServices(lggr logger.Logger, jb job.Job, chainSet evm.LegacyChainContainer, new bool, argsNoPlugin libocr2.OCR2OracleArgs, logError func(string), qopts ...pg.QOpt) ([]job.ServiceCtx, error) {
	execPluginConfig, backfillArgs, err := jobSpecToExecPluginConfig(lggr, jb, chainSet)
	if err != nil {
		return nil, err
	}
	wrappedPluginFactory := NewExecutionReportingPluginFactory(*execPluginConfig)
	err = wrappedPluginFactory.UpdateLogPollerFilters(utils.ZeroAddress, qopts...)
	if err != nil {
		return nil, err
	}

	argsNoPlugin.ReportingPluginFactory = promwrapper.NewPromFactory(wrappedPluginFactory, "CCIPExecution", jb.OCR2OracleSpec.Relay, execPluginConfig.destChainEVMID)
	argsNoPlugin.Logger = relaylogger.NewOCRWrapper(execPluginConfig.lggr, true, logError)
	oracle, err := libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		return nil, err
	}
	// If this is a brand-new job, then we make use of the start blocks. If not then we're rebooting and log poller will pick up where we left off.
	if new {
		return []job.ServiceCtx{
			oraclelib.NewBackfilledOracle(
				execPluginConfig.lggr,
				backfillArgs.sourceLP,
				backfillArgs.destLP,
				backfillArgs.sourceStartBlock,
				backfillArgs.destStartBlock,
				job.NewServiceAdapter(oracle)),
		}, nil
	}
	return []job.ServiceCtx{job.NewServiceAdapter(oracle)}, nil
}

func getTokenDataProviders(lggr logger.Logger, pluginConfig ccipconfig.ExecutionPluginJobSpecConfig, sourceLP logpoller.LogPoller) (map[common.Address]tokendata.Reader, error) {
	tokenDataProviders := make(map[common.Address]tokendata.Reader)

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

		usdcReader, err := ccipdata.NewUSDCReader(lggr, pluginConfig.USDCConfig.SourceMessageTransmitterAddress, sourceLP)
		if err != nil {
			return nil, err
		}
		tokenDataProviders[pluginConfig.USDCConfig.SourceTokenAddress] = tokendata.NewCachedReader(
			usdc.NewUSDCTokenDataReader(
				lggr,
				usdcReader,
				attestationURI,
			),
		)
	}

	return tokenDataProviders, nil
}

func getExecutionPluginSourceLpChainFilters(priceRegistry common.Address) []logpoller.Filter {
	return []logpoller.Filter{
		{
			Name:      logpoller.FilterName(FEE_TOKEN_ADDED, priceRegistry.String()),
			EventSigs: []common.Hash{abihelpers.EventSignatures.FeeTokenAdded},
			Addresses: []common.Address{priceRegistry},
		},
		{
			Name:      logpoller.FilterName(FEE_TOKEN_REMOVED, priceRegistry.String()),
			EventSigs: []common.Hash{abihelpers.EventSignatures.FeeTokenRemoved},
			Addresses: []common.Address{priceRegistry},
		},
	}
}

func getExecutionPluginDestLpChainFilters(commitStore, offRamp, priceRegistry common.Address) []logpoller.Filter {
	return []logpoller.Filter{
		{
			Name:      logpoller.FilterName(EXEC_REPORT_ACCEPTS, commitStore.String()),
			EventSigs: []common.Hash{abihelpers.EventSignatures.ReportAccepted},
			Addresses: []common.Address{commitStore},
		},
		{
			Name:      logpoller.FilterName(EXEC_EXECUTION_STATE_CHANGES, offRamp.String()),
			EventSigs: []common.Hash{abihelpers.EventSignatures.ExecutionStateChanged},
			Addresses: []common.Address{offRamp},
		},
		{
			Name:      logpoller.FilterName(EXEC_TOKEN_POOL_ADDED, offRamp.String()),
			EventSigs: []common.Hash{abihelpers.EventSignatures.PoolAdded},
			Addresses: []common.Address{offRamp},
		},
		{
			Name:      logpoller.FilterName(EXEC_TOKEN_POOL_REMOVED, offRamp.String()),
			EventSigs: []common.Hash{abihelpers.EventSignatures.PoolRemoved},
			Addresses: []common.Address{offRamp},
		},
		{
			Name:      logpoller.FilterName(FEE_TOKEN_ADDED, priceRegistry.String()),
			EventSigs: []common.Hash{abihelpers.EventSignatures.FeeTokenAdded},
			Addresses: []common.Address{priceRegistry},
		},
		{
			Name:      logpoller.FilterName(FEE_TOKEN_REMOVED, priceRegistry.String()),
			EventSigs: []common.Hash{abihelpers.EventSignatures.FeeTokenRemoved},
			Addresses: []common.Address{priceRegistry},
		},
	}
}

// UnregisterExecPluginLpFilters unregisters all the registered filters for both source and dest chains.
func UnregisterExecPluginLpFilters(ctx context.Context, lggr logger.Logger, jb job.Job, chainSet evm.LegacyChainContainer, qopts ...pg.QOpt) error {
	execPluginConfig, _, err := jobSpecToExecPluginConfig(lggr, jb, chainSet)
	if err != nil {
		return err
	}
	if err := execPluginConfig.onRampReader.Close(); err != nil {
		return err
	}
	for _, tokenReader := range execPluginConfig.tokenDataProviders {
		if err := tokenReader.Close(); err != nil {
			return err
		}
	}
	// TODO: once offramp/commit/pricereg are abstracted, we can call Close on the offramp/commit readers to unregister filters.
	return unregisterExecutionPluginLpFilters(ctx, execPluginConfig.sourceLP, execPluginConfig.destLP, execPluginConfig.offRamp,
		execPluginConfig.commitStore.Address(), execPluginConfig.onRamp, execPluginConfig.sourceClient, qopts...)
}

func unregisterExecutionPluginLpFilters(
	ctx context.Context,
	sourceLP logpoller.LogPoller,
	destLP logpoller.LogPoller,
	destOffRamp evm_2_evm_offramp.EVM2EVMOffRampInterface,
	commitStore common.Address,
	sourceOnRamp evm_2_evm_onramp.EVM2EVMOnRampInterface,
	sourceChainClient client.Client,
	qopts ...pg.QOpt) error {
	destOffRampDynCfg, err := destOffRamp.GetDynamicConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return err
	}

	// TODO stopgap solution before compatibility phase-2
	tvStr, err := sourceOnRamp.TypeAndVersion(&bind.CallOpts{Context: ctx})
	if err != nil {
		return err
	}
	_, versionStr, err := ccipconfig.ParseTypeAndVersion(tvStr)
	if err != nil {
		return err
	}
	version, err := semver.NewVersion(versionStr)
	if err != nil {
		return err
	}

	onRampDynCfg, err := contractutil.LoadOnRampDynamicConfig(sourceOnRamp, *version, sourceChainClient)
	if err != nil {
		return err
	}

	if err = logpollerutil.UnregisterLpFilters(
		sourceLP,
		getExecutionPluginSourceLpChainFilters(onRampDynCfg.PriceRegistry),
		qopts...,
	); err != nil {
		return err
	}

	return logpollerutil.UnregisterLpFilters(
		destLP,
		getExecutionPluginDestLpChainFilters(commitStore, destOffRamp.Address(), destOffRampDynCfg.PriceRegistry),
		qopts...,
	)
}

// ExecutionReportToEthTxMeta generates a txmgr.EthTxMeta from the given report.
// Only MessageIDs will be populated in the TxMeta.
func ExecutionReportToEthTxMeta(report []byte) (*txmgr.TxMeta, error) {
	execReport, err := abihelpers.DecodeExecutionReport(report)
	if err != nil {
		return nil, err
	}

	msgIDs := make([]string, len(execReport.Messages))
	for i, msg := range execReport.Messages {
		msgIDs[i] = hexutil.Encode(msg.MessageId[:])
	}

	return &txmgr.TxMeta{
		MessageIDs: msgIDs,
	}, nil
}
