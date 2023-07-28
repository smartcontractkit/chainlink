package evm

import (
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

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
)

type functionsProvider struct {
	*configWatcher
	contractTransmitter ContractTransmitter
}

var (
	_ relaytypes.Plugin = (*functionsProvider)(nil)
)

func (p *functionsProvider) ContractTransmitter() types.ContractTransmitter {
	return p.contractTransmitter
}

func NewFunctionsProvider(chainSet evm.ChainSet, rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs, lggr logger.Logger, ethKeystore keystore.Eth, pluginType functionsRelay.FunctionsPluginType) (relaytypes.Plugin, error) {
	configWatcher, err := newFunctionsConfigProvider(pluginType, chainSet, rargs, lggr)
	if err != nil {
		return nil, err
	}
	var pluginConfig config.PluginConfig
	if err2 := json.Unmarshal(pargs.PluginConfig, &pluginConfig); err2 != nil {
		return nil, err2
	}
	contractTransmitter, err := newFunctionsContractTransmitter(pluginConfig.ContractVersion, rargs, pargs.TransmitterID, configWatcher, ethKeystore, lggr)
	if err != nil {
		return nil, err
	}
	return &functionsProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
	}, nil
}

func newFunctionsConfigProvider(pluginType functionsRelay.FunctionsPluginType, chainSet evm.ChainSet, args relaytypes.RelayArgs, lggr logger.Logger) (*configWatcher, error) {
	var relayConfig evmRelayTypes.RelayConfig
	err := json.Unmarshal(args.RelayConfig, &relayConfig)
	if err != nil {
		return nil, err
	}
	chain, err := chainSet.Get(relayConfig.ChainID.ToInt())
	if err != nil {
		return nil, err
	}
	if !common.IsHexAddress(args.ContractID) {
		return nil, errors.Errorf("invalid contractID, expected hex address")
	}

	contractAddress := common.HexToAddress(args.ContractID)
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorMetaData.ABI))
	if err != nil {
		return nil, errors.Wrap(err, "could not get contract ABI JSON")
	}

	cp, err := functionsRelay.NewFunctionsConfigPoller(pluginType, chain.LogPoller(), contractAddress, lggr)
	if err != nil {
		return nil, err
	}

	offchainConfigDigester := functionsRelay.FunctionsOffchainConfigDigester{
		PluginType: pluginType,
		BaseDigester: evmutil.EVMOffchainConfigDigester{
			ChainID:         chain.ID().Uint64(),
			ContractAddress: contractAddress,
		},
	}

	return newConfigWatcher(lggr, contractAddress, contractABI, offchainConfigDigester, cp, chain, relayConfig.FromBlock, args.New), nil
}

func newFunctionsContractTransmitter(contractVersion uint32, rargs relaytypes.RelayArgs, transmitterID string, configWatcher *configWatcher, ethKeystore keystore.Eth, lggr logger.Logger) (ContractTransmitter, error) {
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

	return functionsRelay.NewFunctionsContractTransmitter(
		configWatcher.contractAddress,
		configWatcher.chain.Client(),
		configWatcher.contractABI,
		transmitter,
		configWatcher.chain.LogPoller(),
		lggr,
		nil,
		contractVersion,
	)
}
