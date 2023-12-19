package evm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/plugin"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var (
	_ OCR2KeeperRelayer  = (*ocr2keeperRelayer)(nil)
	_ OCR2KeeperProvider = (*ocr2keeperProvider)(nil)
)

// OCR2KeeperProviderOpts is the custom options to create a keeper provider
type OCR2KeeperProviderOpts struct {
	RArgs      commontypes.RelayArgs
	PArgs      commontypes.PluginArgs
	InstanceID int
}

// OCR2KeeperProvider provides all components needed for a OCR2Keeper plugin.
type OCR2KeeperProvider interface {
	commontypes.Plugin
}

// OCR2KeeperRelayer contains the relayer and instantiating functions for OCR2Keeper providers.
type OCR2KeeperRelayer interface {
	NewOCR2KeeperProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (OCR2KeeperProvider, error)
}

// ocr2keeperRelayer is the relayer with added DKG and OCR2Keeper provider functions.
type ocr2keeperRelayer struct {
	db          *sqlx.DB
	chain       legacyevm.Chain
	lggr        logger.Logger
	ethKeystore keystore.Eth
}

// NewOCR2KeeperRelayer is the constructor of ocr2keeperRelayer
func NewOCR2KeeperRelayer(db *sqlx.DB, chain legacyevm.Chain, lggr logger.Logger, ethKeystore keystore.Eth) OCR2KeeperRelayer {
	return &ocr2keeperRelayer{
		db:          db,
		chain:       chain,
		lggr:        lggr,
		ethKeystore: ethKeystore,
	}
}

func (r *ocr2keeperRelayer) NewOCR2KeeperProvider(rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (OCR2KeeperProvider, error) {
	cfgWatcher, err := newOCR2KeeperConfigProvider(r.lggr, r.chain, rargs)
	if err != nil {
		return nil, err
	}

	gasLimit := cfgWatcher.chain.Config().EVM().OCR2().Automation().GasLimit()
	contractTransmitter, err := newContractTransmitter(r.lggr, rargs, pargs.TransmitterID, r.ethKeystore, cfgWatcher, configTransmitterOpts{pluginGasLimit: &gasLimit}, nil)
	if err != nil {
		return nil, err
	}

	return &ocr2keeperProvider{
		configWatcher:       cfgWatcher,
		contractTransmitter: contractTransmitter,
	}, nil
}

type ocr3keeperProviderContractTransmitter struct {
	contractTransmitter ocrtypes.ContractTransmitter
}

var _ ocr3types.ContractTransmitter[plugin.AutomationReportInfo] = &ocr3keeperProviderContractTransmitter{}

func NewKeepersOCR3ContractTransmitter(ocr2ContractTransmitter ocrtypes.ContractTransmitter) *ocr3keeperProviderContractTransmitter {
	return &ocr3keeperProviderContractTransmitter{ocr2ContractTransmitter}
}

func (t *ocr3keeperProviderContractTransmitter) Transmit(
	ctx context.Context,
	digest ocrtypes.ConfigDigest,
	seqNr uint64,
	reportWithInfo ocr3types.ReportWithInfo[plugin.AutomationReportInfo],
	aoss []ocrtypes.AttributedOnchainSignature,
) error {
	return t.contractTransmitter.Transmit(
		ctx,
		ocrtypes.ReportContext{
			ReportTimestamp: ocrtypes.ReportTimestamp{
				ConfigDigest: digest,
				Epoch:        uint32(seqNr),
			},
		},
		reportWithInfo.Report,
		aoss,
	)
}

func (t *ocr3keeperProviderContractTransmitter) FromAccount() (ocrtypes.Account, error) {
	return t.contractTransmitter.FromAccount()
}

type ocr2keeperProvider struct {
	*configWatcher
	contractTransmitter ContractTransmitter
}

func (c *ocr2keeperProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return c.contractTransmitter
}

func (c *ocr2keeperProvider) ChainReader() commontypes.ChainReader {
	return nil
}

func newOCR2KeeperConfigProvider(lggr logger.Logger, chain legacyevm.Chain, rargs commontypes.RelayArgs) (*configWatcher, error) {
	var relayConfig types.RelayConfig
	err := json.Unmarshal(rargs.RelayConfig, &relayConfig)
	if err != nil {
		return nil, err
	}
	if !common.IsHexAddress(rargs.ContractID) {
		return nil, fmt.Errorf("invalid contract address '%s'", rargs.ContractID)
	}

	contractAddress := common.HexToAddress(rargs.ContractID)
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorMetaData.ABI))
	if err != nil {
		return nil, errors.Wrap(err, "could not get OCR2Aggregator ABI JSON")
	}

	configPoller, err := NewConfigPoller(
		lggr.With("contractID", rargs.ContractID),
		chain.Client(),
		chain.LogPoller(),
		contractAddress,
		// TODO: Does ocr2keeper need to support config contract? DF-19182
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create config poller")
	}

	offchainConfigDigester := evmutil.EVMOffchainConfigDigester{
		ChainID:         chain.Config().EVM().ChainID().Uint64(),
		ContractAddress: contractAddress,
	}

	return newConfigWatcher(
		lggr,
		contractAddress,
		contractABI,
		offchainConfigDigester,
		configPoller,
		chain,
		relayConfig.FromBlock,
		rargs.New,
	), nil
}
