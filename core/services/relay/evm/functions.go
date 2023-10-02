package evm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"go.uber.org/multierr"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
	txmgrcommon "github.com/smartcontractkit/chainlink/v2/common/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	txm "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	functionsRelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/functions"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type functionsProvider struct {
	utils.StartStopOnce
	configWatcher       *configWatcher
	contractTransmitter ContractTransmitter
	logPollerWrapper    evmRelayTypes.LogPollerWrapper
}

var _ evmRelayTypes.FunctionsProvider = (*functionsProvider)(nil)

func (p *functionsProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return p.contractTransmitter
}

func (p *functionsProvider) LogPollerWrapper() evmRelayTypes.LogPollerWrapper {
	return p.logPollerWrapper
}

func (p *functionsProvider) FunctionsEvents() relaytypes.FunctionsEvents {
	// TODO (FUN-668): implement
	return nil
}

func (p *functionsProvider) Start(ctx context.Context) error {
	return p.StartOnce("FunctionsProvider", func() error {
		if err := p.configWatcher.Start(ctx); err != nil {
			return err
		}
		return p.logPollerWrapper.Start(ctx)
	})
}

func (p *functionsProvider) Close() error {
	return p.StopOnce("FunctionsProvider", func() (err error) {
		err = multierr.Combine(err, p.logPollerWrapper.Close())
		err = multierr.Combine(err, p.configWatcher.Close())
		return
	})
}

// Forward all calls to the underlying configWatcher
func (p *functionsProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return p.configWatcher.OffchainConfigDigester()
}

func (p *functionsProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return p.configWatcher.ContractConfigTracker()
}

func (p *functionsProvider) HealthReport() map[string]error {
	return p.configWatcher.HealthReport()
}

func (p *functionsProvider) Name() string {
	return p.configWatcher.Name()
}

func NewFunctionsProvider(chain evm.Chain, rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs, lggr logger.Logger, ethKeystore keystore.Eth, pluginType functionsRelay.FunctionsPluginType) (evmRelayTypes.FunctionsProvider, error) {
	relayOpts := evmRelayTypes.NewRelayOpts(rargs)
	relayConfig, err := relayOpts.RelayConfig()
	if err != nil {
		return nil, err
	}
	expectedChainID := relayConfig.ChainID.String()
	if expectedChainID != chain.ID().String() {
		return nil, fmt.Errorf("internal error: chain id in spec does not match this relayer's chain: have %s expected %s", relayConfig.ChainID.String(), chain.ID().String())
	}
	if err != nil {
		return nil, err
	}
	if !common.IsHexAddress(rargs.ContractID) {
		return nil, errors.Errorf("invalid contractID, expected hex address")
	}
	var pluginConfig config.PluginConfig
	if err2 := json.Unmarshal(pargs.PluginConfig, &pluginConfig); err2 != nil {
		return nil, err2
	}
	routerContractAddress := common.HexToAddress(rargs.ContractID)
	logPollerWrapper, err := functionsRelay.NewLogPollerWrapper(routerContractAddress, pluginConfig, chain.Client(), chain.LogPoller(), lggr)
	if err != nil {
		return nil, err
	}
	configWatcher, err := newFunctionsConfigProvider(pluginType, chain, rargs, relayConfig.FromBlock, logPollerWrapper, lggr)
	if err != nil {
		return nil, err
	}
	var contractTransmitter ContractTransmitter
	if relayConfig.SendingKeys != nil {
		contractTransmitter, err = newFunctionsContractTransmitter(pluginConfig.ContractVersion, rargs, pargs.TransmitterID, configWatcher, ethKeystore, logPollerWrapper, lggr)
		if err != nil {
			return nil, err
		}
	} else {
		lggr.Warn("no sending keys configured for functions plugin, not starting contract transmitter")
	}
	return &functionsProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
		logPollerWrapper:    logPollerWrapper,
	}, nil
}

func newFunctionsConfigProvider(pluginType functionsRelay.FunctionsPluginType, chain evm.Chain, args relaytypes.RelayArgs, fromBlock uint64, logPollerWrapper evmRelayTypes.LogPollerWrapper, lggr logger.Logger) (*configWatcher, error) {
	if !common.IsHexAddress(args.ContractID) {
		return nil, errors.Errorf("invalid contractID, expected hex address")
	}

	routerContractAddress := common.HexToAddress(args.ContractID)
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorMetaData.ABI))
	if err != nil {
		return nil, errors.Wrap(err, "could not get contract ABI JSON")
	}

	cp, err := functionsRelay.NewFunctionsConfigPoller(pluginType, chain.LogPoller(), lggr)
	if err != nil {
		return nil, err
	}
	logPollerWrapper.SubscribeToUpdates("FunctionsConfigPoller", cp)

	offchainConfigDigester := functionsRelay.NewFunctionsOffchainConfigDigester(pluginType, chain.ID().Uint64())
	logPollerWrapper.SubscribeToUpdates("FunctionsOffchainConfigDigester", offchainConfigDigester)

	return newConfigWatcher(lggr, routerContractAddress, contractABI, offchainConfigDigester, cp, chain, fromBlock, args.New), nil
}

func newFunctionsContractTransmitter(contractVersion uint32, rargs relaytypes.RelayArgs, transmitterID string, configWatcher *configWatcher, ethKeystore keystore.Eth, logPollerWrapper evmRelayTypes.LogPollerWrapper, lggr logger.Logger) (ContractTransmitter, error) {
	var relayConfig evmRelayTypes.RelayConfig
	if err := json.Unmarshal(rargs.RelayConfig, &relayConfig); err != nil {
		return nil, err
	}
	var fromAddresses []common.Address
	sendingKeys := relayConfig.SendingKeys
	if !relayConfig.EffectiveTransmitterID.Valid {
		return nil, errors.New("EffectiveTransmitterID must be specified")
	}
	effectiveTransmitterAddress := common.HexToAddress(relayConfig.EffectiveTransmitterID.String)

	sendingKeysLength := len(sendingKeys)
	if sendingKeysLength == 0 {
		return nil, errors.New("no sending keys provided")
	}

	// If we are using multiple sending keys, then a forwarder is needed to rotate transmissions.
	// Ensure that this forwarder is not set to a local sending key, and ensure our sending keys are enabled.
	for _, s := range sendingKeys {
		if sendingKeysLength > 1 && s == effectiveTransmitterAddress.String() {
			return nil, errors.New("the transmitter is a local sending key with transaction forwarding enabled")
		}
		if err := ethKeystore.CheckEnabled(common.HexToAddress(s), configWatcher.chain.Config().EVM().ChainID()); err != nil {
			return nil, errors.Wrap(err, "one of the sending keys given is not enabled")
		}
		fromAddresses = append(fromAddresses, common.HexToAddress(s))
	}

	scoped := configWatcher.chain.Config()
	strategy := txmgrcommon.NewQueueingTxStrategy(rargs.ExternalJobID, scoped.OCR2().DefaultTransactionQueueDepth(), scoped.Database().DefaultQueryTimeout())

	var checker txm.TransmitCheckerSpec
	if configWatcher.chain.Config().OCR2().SimulateTransactions() {
		checker.CheckerType = txm.TransmitCheckerTypeSimulate
	}

	gasLimit := configWatcher.chain.Config().EVM().GasEstimator().LimitDefault()
	ocr2Limit := configWatcher.chain.Config().EVM().GasEstimator().LimitJobType().OCR2()
	if ocr2Limit != nil {
		gasLimit = *ocr2Limit
	}

	transmitter, err := ocrcommon.NewTransmitter(
		configWatcher.chain.TxManager(),
		fromAddresses,
		gasLimit,
		effectiveTransmitterAddress,
		strategy,
		checker,
		configWatcher.chain.ID(),
		ethKeystore,
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create transmitter")
	}

	functionsTransmitter, err := functionsRelay.NewFunctionsContractTransmitter(
		configWatcher.chain.Client(),
		configWatcher.contractABI,
		transmitter,
		configWatcher.chain.LogPoller(),
		lggr,
		nil,
		contractVersion,
	)
	if err != nil {
		return nil, err
	}
	logPollerWrapper.SubscribeToUpdates("FunctionsConfigTransmitter", functionsTransmitter)
	return functionsTransmitter, err
}
