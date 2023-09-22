package ccip

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	relaylogger "github.com/smartcontractkit/chainlink-relay/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/contractutil"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
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
	EXEC_CCIP_SENDS              = "Exec ccip sends"
	EXEC_REPORT_ACCEPTS          = "Exec report accepts"
	EXEC_EXECUTION_STATE_CHANGES = "Exec execution state changes"
	EXEC_TOKEN_POOL_ADDED        = "Token pool added"
	EXEC_TOKEN_POOL_REMOVED      = "Token pool removed"
	FEE_TOKEN_ADDED              = "Fee token added"
	FEE_TOKEN_REMOVED            = "Fee token removed"
)

func NewExecutionServices(lggr logger.Logger, jb job.Job, chainSet evm.LegacyChainContainer, new bool, argsNoPlugin libocr2.OCR2OracleArgs, logError func(string), qopts ...pg.QOpt) ([]job.ServiceCtx, error) {
	spec := jb.OCR2OracleSpec
	var pluginConfig ccipconfig.ExecutionPluginJobSpecConfig
	err := json.Unmarshal(spec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, err
	}

	chainIDInterface, ok := spec.RelayConfig["chainID"]
	if !ok {
		return nil, errors.New("chainID must be provided in relay config")
	}
	destChainID := int64(chainIDInterface.(float64))
	destChain, err := chainSet.Get(strconv.FormatInt(destChainID, 10))
	if err != nil {
		return nil, errors.Wrap(err, "get chainset")
	}
	offRamp, err := contractutil.LoadOffRamp(common.HexToAddress(spec.ContractID), ExecPluginLabel, destChain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "failed loading offRamp")
	}
	offRampConfig, err := offRamp.GetStaticConfig(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	chainId, err := chainselectors.ChainIdFromSelector(offRampConfig.SourceChainSelector)
	if err != nil {
		return nil, err
	}
	sourceChain, err := chainSet.Get(strconv.FormatUint(chainId, 10))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open source chain")
	}
	commitStore, err := contractutil.LoadCommitStore(offRampConfig.CommitStore, ExecPluginLabel, destChain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "failed loading commitStore")
	}
	onRamp, err := contractutil.LoadOnRamp(offRampConfig.OnRamp, ExecPluginLabel, sourceChain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "failed loading onRamp")
	}
	dynamicOnRampConfig, err := contractutil.LoadOnRampDynamicConfig(onRamp, sourceChain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "failed loading onRamp config")
	}
	sourceRouter, err := router.NewRouter(dynamicOnRampConfig.Router, sourceChain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "failed loading source router")
	}
	sourceWrappedNative, err := sourceRouter.GetWrappedNative(&bind.CallOpts{})
	if err != nil {
		return nil, errors.Wrap(err, "could not get source native token")
	}
	sourcePriceRegistry, err := observability.NewObservedPriceRegistry(dynamicOnRampConfig.PriceRegistry, ExecPluginLabel, sourceChain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "could not create source price registry")
	}

	execLggr := lggr.Named("CCIPExecution").With(
		"sourceChain", ChainName(int64(chainId)),
		"destChain", ChainName(destChainID))

	sourceChainEventClient := ccipdata.NewLogPollerReader(sourceChain.LogPoller(), execLggr, sourceChain.Client())

	tokenDataProviders, err := getTokenDataProviders(lggr, pluginConfig, offRampConfig.OnRamp, sourceChainEventClient)
	if err != nil {
		return nil, errors.Wrap(err, "could not get token data providers")
	}

	wrappedPluginFactory := NewExecutionReportingPluginFactory(
		ExecutionPluginConfig{
			lggr:                     execLggr,
			sourceLP:                 sourceChain.LogPoller(),
			destLP:                   destChain.LogPoller(),
			sourceReader:             sourceChainEventClient,
			destReader:               ccipdata.NewLogPollerReader(destChain.LogPoller(), execLggr, destChain.Client()),
			onRamp:                   onRamp,
			offRamp:                  offRamp,
			commitStore:              commitStore,
			sourcePriceRegistry:      sourcePriceRegistry,
			sourceWrappedNativeToken: sourceWrappedNative,
			destClient:               destChain.Client(),
			sourceClient:             sourceChain.Client(),
			destGasEstimator:         destChain.GasEstimator(),
			leafHasher:               hashlib.NewLeafHasher(offRampConfig.SourceChainSelector, offRampConfig.ChainSelector, onRamp.Address(), hashlib.NewKeccakCtx()),
			tokenDataProviders:       tokenDataProviders,
		})

	err = wrappedPluginFactory.UpdateLogPollerFilters(utils.ZeroAddress, qopts...)
	if err != nil {
		return nil, err
	}

	argsNoPlugin.ReportingPluginFactory = promwrapper.NewPromFactory(wrappedPluginFactory, "CCIPExecution", spec.Relay, destChain.ID())
	argsNoPlugin.Logger = relaylogger.NewOCRWrapper(execLggr, true, logError)
	oracle, err := libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		return nil, err
	}
	execLggr.Infow("Initialized exec plugin",
		"pluginConfig", pluginConfig,
		"onRampAddress", onRamp.Address(),
		"sourcePriceRegistry", sourcePriceRegistry.Address(),
		"dynamicOnRampConfig", dynamicOnRampConfig,
		"sourceNative", sourceWrappedNative,
		"sourceRouter", sourceRouter.Address())
	// If this is a brand-new job, then we make use of the start blocks. If not then we're rebooting and log poller will pick up where we left off.
	if new {
		return []job.ServiceCtx{
			oraclelib.NewBackfilledOracle(
				execLggr,
				sourceChain.LogPoller(),
				destChain.LogPoller(),
				pluginConfig.SourceStartBlock,
				pluginConfig.DestStartBlock,
				job.NewServiceAdapter(oracle)),
		}, nil
	}
	return []job.ServiceCtx{job.NewServiceAdapter(oracle)}, nil
}

func getTokenDataProviders(lggr logger.Logger, pluginConfig ccipconfig.ExecutionPluginJobSpecConfig, onRampAddress common.Address, sourceChainEventClient *ccipdata.LogPollerReader) (map[common.Address]tokendata.Reader, error) {
	tokenDataProviders := make(map[common.Address]tokendata.Reader)

	if pluginConfig.USDCConfig.AttestationAPI != "" {
		lggr.Infof("USDC token data provider enabled")
		err := pluginConfig.USDCConfig.ValidateUSDCConfig()
		if err != nil {
			return nil, err
		}

		attestationURI, err2 := url.ParseRequestURI(pluginConfig.USDCConfig.AttestationAPI)
		if err2 != nil {
			return nil, errors.Wrap(err2, "failed to parse USDC attestation API")
		}

		tokenDataProviders[pluginConfig.USDCConfig.SourceTokenAddress] = tokendata.NewCachedReader(
			usdc.NewUSDCTokenDataReader(
				lggr,
				sourceChainEventClient,
				pluginConfig.USDCConfig.SourceTokenAddress,
				pluginConfig.USDCConfig.SourceMessageTransmitterAddress,
				onRampAddress,
				attestationURI,
			),
		)
	}

	return tokenDataProviders, nil
}

func getExecutionPluginSourceLpChainFilters(onRamp, priceRegistry common.Address, tokenDataProviders map[common.Address]tokendata.Reader) []logpoller.Filter {
	var filters []logpoller.Filter
	for _, provider := range tokenDataProviders {
		filters = append(filters, provider.GetSourceLogPollerFilters()...)
	}

	return append(filters, []logpoller.Filter{
		{
			Name:      logpoller.FilterName(EXEC_CCIP_SENDS, onRamp.String()),
			EventSigs: []common.Hash{abihelpers.EventSignatures.SendRequested},
			Addresses: []common.Address{onRamp},
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
	}...)
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
func UnregisterExecPluginLpFilters(ctx context.Context, lggr logger.Logger, spec *job.OCR2OracleSpec, chainSet evm.LegacyChainContainer, qopts ...pg.QOpt) error {
	if spec == nil {
		return errors.New("spec is nil")
	}

	var pluginConfig ccipconfig.ExecutionPluginJobSpecConfig
	err := json.Unmarshal(spec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return err
	}

	destChainIDInterface, ok := spec.RelayConfig["chainID"]
	if !ok {
		return errors.New("chainID must be provided in relay config")
	}
	destChainIDf64, is := destChainIDInterface.(float64)
	if !is {
		return fmt.Errorf("chain id '%v' is not float64", destChainIDInterface)
	}
	destChain, err := chainSet.Get(strconv.FormatInt(int64(destChainIDf64), 10))
	if err != nil {
		return err
	}

	offRampAddress := common.HexToAddress(spec.ContractID)
	offRamp, err := contractutil.LoadOffRamp(offRampAddress, ExecPluginLabel, destChain.Client())
	if err != nil {
		return err
	}

	offRampConfig, err := offRamp.GetStaticConfig(&bind.CallOpts{})
	if err != nil {
		return err
	}
	chainId, err := chainselectors.ChainIdFromSelector(offRampConfig.SourceChainSelector)
	if err != nil {
		return err
	}
	sourceChain, err := chainSet.Get(strconv.FormatUint(chainId, 10))
	if err != nil {
		return errors.Wrap(err, "unable to open source chain")
	}
	sourceOnRamp, err := contractutil.LoadOnRamp(offRampConfig.OnRamp, ExecPluginLabel, sourceChain.Client())
	if err != nil {
		return errors.Wrap(err, "failed loading onRamp")
	}

	return unregisterExecutionPluginLpFilters(ctx, lggr, sourceChain.LogPoller(), destChain.LogPoller(), offRamp, offRampConfig, sourceOnRamp, sourceChain.Client(), pluginConfig, qopts...)
}

func unregisterExecutionPluginLpFilters(
	ctx context.Context,
	lggr logger.Logger,
	sourceLP logpoller.LogPoller,
	destLP logpoller.LogPoller,
	destOffRamp evm_2_evm_offramp.EVM2EVMOffRampInterface,
	destOffRampConfig evm_2_evm_offramp.EVM2EVMOffRampStaticConfig,
	sourceOnRamp evm_2_evm_onramp.EVM2EVMOnRampInterface,
	sourceChainClient client.Client,
	pluginConfig ccipconfig.ExecutionPluginJobSpecConfig,
	qopts ...pg.QOpt) error {
	destOffRampDynCfg, err := destOffRamp.GetDynamicConfig(&bind.CallOpts{Context: ctx})
	if err != nil {
		return err
	}

	onRampDynCfg, err := contractutil.LoadOnRampDynamicConfig(sourceOnRamp, sourceChainClient)
	if err != nil {
		return err
	}

	// SourceChainEventClient can be nil because it is not used in unregisterExecutionPluginLpFilters
	tokenDataProviders, err := getTokenDataProviders(lggr, pluginConfig, destOffRampConfig.OnRamp, nil)
	if err != nil {
		return err
	}

	if err = logpollerutil.UnregisterLpFilters(
		sourceLP,
		getExecutionPluginSourceLpChainFilters(destOffRampConfig.OnRamp, onRampDynCfg.PriceRegistry, tokenDataProviders),
		qopts...,
	); err != nil {
		return err
	}

	return logpollerutil.UnregisterLpFilters(
		destLP,
		getExecutionPluginDestLpChainFilters(destOffRampConfig.CommitStore, destOffRamp.Address(), destOffRampDynCfg.PriceRegistry),
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
