package evm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/plugin"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/automation"
	ac "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_automation_v21_plus_common"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"

	evm "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/logprovider"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/transmit"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/upkeepstate"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

var (
	_                        OCR2KeeperRelayer  = (*ocr2keeperRelayer)(nil)
	_                        OCR2KeeperProvider = (*ocr2keeperProvider)(nil)
	ErrInitializationFailure                    = fmt.Errorf("failed to initialize registry")
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
	Registry() automation.Registry
	Encoder() automation.Encoder
	TransmitEventProvider() automation.EventProvider
	BlockSubscriber() automation.BlockSubscriber
	PayloadBuilder() automation.PayloadBuilder
	UpkeepStateStore() automation.UpkeepStateStore
	LogEventProvider() automation.LogEventProvider
	LogRecoverer() automation.LogRecoverer
	UpkeepProvider() automation.ConditionalUpkeepProvider
}

// OCR2KeeperRelayer contains the relayer and instantiating functions for OCR2Keeper providers.
type OCR2KeeperRelayer interface {
	NewOCR2KeeperProvider(ctx context.Context, rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (OCR2KeeperProvider, error)
}

// ocr2keeperRelayer is the relayer with added DKG and OCR2Keeper provider functions.
type ocr2keeperRelayer struct {
	ds          sqlutil.DataSource
	chain       legacyevm.Chain
	lggr        logger.Logger
	ethKeystore keys.Store
}

// NewOCR2KeeperRelayer is the constructor of ocr2keeperRelayer
func NewOCR2KeeperRelayer(ds sqlutil.DataSource, chain legacyevm.Chain, lggr logger.Logger, ethKeystore keys.Store) OCR2KeeperRelayer {
	return &ocr2keeperRelayer{
		ds:          ds,
		chain:       chain,
		lggr:        lggr,
		ethKeystore: ethKeystore,
	}
}

func (r *ocr2keeperRelayer) NewOCR2KeeperProvider(ctx context.Context, rargs commontypes.RelayArgs, pargs commontypes.PluginArgs) (OCR2KeeperProvider, error) {
	cfgWatcher, err := newOCR2KeeperConfigProvider(ctx, r.lggr, r.chain, rargs)
	if err != nil {
		return nil, err
	}

	gasLimit := cfgWatcher.chain.Config().EVM().OCR2().Automation().GasLimit()
	contractTransmitter, err := newOnChainContractTransmitter(ctx, r.lggr, rargs, r.ethKeystore, cfgWatcher, configTransmitterOpts{pluginGasLimit: &gasLimit}, OCR2AggregatorTransmissionContractABI)
	if err != nil {
		return nil, err
	}

	client := r.chain

	services := new(ocr2keeperProvider)
	services.configWatcher = cfgWatcher
	services.contractTransmitter = contractTransmitter

	addr := evmtypes.MustEIP55Address(rargs.ContractID).Address()

	registryContract, err := ac.NewIAutomationV21PlusCommon(addr, client.Client())
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create caller for address and backend", ErrInitializationFailure)
	}
	// lookback blocks for transmit event is hard coded and should provide ample time for logs
	// to be detected in most cases
	var transmitLookbackBlocks int64 = 250
	transmitEventProvider, err := transmit.NewTransmitEventProvider(ctx, r.lggr, client.LogPoller(), addr, client.Client(), transmitLookbackBlocks)
	if err != nil {
		return nil, err
	}

	services.transmitEventProvider = transmitEventProvider

	packer := encoding.NewAbiPacker()
	services.encoder = encoding.NewReportEncoder(packer)

	finalityDepth := client.Config().EVM().FinalityDepth()

	orm := upkeepstate.NewORM(client.ID(), r.ds)
	scanner := upkeepstate.NewPerformedEventsScanner(r.lggr, client.LogPoller(), addr, finalityDepth)
	services.upkeepStateStore = upkeepstate.NewUpkeepStateStore(orm, r.lggr, scanner)

	logProvider, logRecoverer := logprovider.New(r.lggr, client.LogPoller(), client.Client(), services.upkeepStateStore, finalityDepth, client.ID())
	services.logEventProvider = logProvider
	services.logRecoverer = logRecoverer
	blockSubscriber := evm.NewBlockSubscriber(client.HeadBroadcaster(), client.LogPoller(), finalityDepth, r.lggr)
	services.blockSubscriber = blockSubscriber

	al := evm.NewActiveUpkeepList()
	services.payloadBuilder = evm.NewPayloadBuilder(al, logRecoverer, r.lggr)

	services.registry = evm.NewEvmRegistry(r.lggr, addr, client,
		registryContract, rargs.MercuryCredentials, al, logProvider,
		packer, blockSubscriber, finalityDepth)

	services.conditionalUpkeepProvider = evm.NewUpkeepProvider(al, blockSubscriber, client.LogPoller())

	return services, nil
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

func (t *ocr3keeperProviderContractTransmitter) FromAccount(ctx context.Context) (ocrtypes.Account, error) {
	return t.contractTransmitter.FromAccount(ctx)
}

type ocr2keeperProvider struct {
	*configWatcher
	contractTransmitter       ContractTransmitter
	registry                  automation.Registry
	encoder                   automation.Encoder
	transmitEventProvider     automation.EventProvider
	blockSubscriber           automation.BlockSubscriber
	payloadBuilder            automation.PayloadBuilder
	upkeepStateStore          automation.UpkeepStateStore
	logEventProvider          automation.LogEventProvider
	logRecoverer              automation.LogRecoverer
	conditionalUpkeepProvider automation.ConditionalUpkeepProvider
}

func (c *ocr2keeperProvider) ContractTransmitter() ocrtypes.ContractTransmitter {
	return c.contractTransmitter
}

func (c *ocr2keeperProvider) ContractReader() commontypes.ContractReader {
	return nil
}

func (c *ocr2keeperProvider) Codec() commontypes.Codec {
	return nil
}

func newOCR2KeeperConfigProvider(ctx context.Context, lggr logger.Logger, chain legacyevm.Chain, rargs commontypes.RelayArgs) (*configWatcher, error) {
	var relayConfig types.RelayConfig
	err := json.Unmarshal(rargs.RelayConfig, &relayConfig)
	if err != nil {
		return nil, err
	}
	if !common.IsHexAddress(rargs.ContractID) {
		return nil, fmt.Errorf("invalid contract address '%s'", rargs.ContractID)
	}

	contractAddress := common.HexToAddress(rargs.ContractID)

	configPoller, err := NewConfigPoller(
		ctx,
		logger.With(lggr, "contractID", rargs.ContractID),
		CPConfig{
			chain.Client(),
			chain.LogPoller(),
			contractAddress,
			// TODO: Does ocr2keeper need to support config contract? DF-19182
			nil,
			OCR2AggregatorLogDecoder,
		},
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
		offchainConfigDigester,
		configPoller,
		chain,
		relayConfig.FromBlock,
		rargs.New,
	), nil
}

func (c *ocr2keeperProvider) Registry() automation.Registry {
	return c.registry
}

func (c *ocr2keeperProvider) Encoder() automation.Encoder {
	return c.encoder
}

func (c *ocr2keeperProvider) TransmitEventProvider() automation.EventProvider {
	return c.transmitEventProvider
}

func (c *ocr2keeperProvider) BlockSubscriber() automation.BlockSubscriber {
	return c.blockSubscriber
}

func (c *ocr2keeperProvider) PayloadBuilder() automation.PayloadBuilder {
	return c.payloadBuilder
}

func (c *ocr2keeperProvider) UpkeepStateStore() automation.UpkeepStateStore {
	return c.upkeepStateStore
}

func (c *ocr2keeperProvider) LogEventProvider() automation.LogEventProvider {
	return c.logEventProvider
}

func (c *ocr2keeperProvider) LogRecoverer() automation.LogRecoverer {
	return c.logRecoverer
}

func (c *ocr2keeperProvider) UpkeepProvider() automation.ConditionalUpkeepProvider {
	return c.conditionalUpkeepProvider
}
