package evm

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/ocr2keepers/pkg/v3/plugin"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/feed_lookup_compatible_interface"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/upkeepstate"
)

type AutomationFacade interface {
	Registry() *EvmRegistry
	Encoder() ocr2keepers.Encoder
	TransmitEventProvider() *TransmitEventProvider
	BlockSubscriber() *BlockSubscriber
	PayloadBuilder() ocr2keepers.PayloadBuilder
	UpkeepStateStore() upkeepstate.UpkeepStateStore
	LogEventProvider() logprovider.LogEventProvider
	LogRecoverer() logprovider.LogRecoverer
	UpkeepProvider() ocr2keepers.ConditionalUpkeepProvider
	Keyring() ocr3types.OnchainKeyring[plugin.AutomationReportInfo]

	Start(context.Context) error
	io.Closer
}

func New(addr common.Address, client evm.Chain, mc *models.MercuryCredentials, keyring ocrtypes.OnchainKeyring, lggr logger.Logger) (AutomationFacade, error) {
	feedLookupCompatibleABI, err := abi.JSON(strings.NewReader(feed_lookup_compatible_interface.FeedLookupCompatibleInterfaceABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}
	keeperRegistryABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}
	utilsABI, err := abi.JSON(strings.NewReader(automation_utils_2_1.AutomationUtilsABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}
	registryContract, err := iregistry21.NewIKeeperRegistryMaster(addr, client.Client())
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create caller for address and backend", ErrInitializationFailure)
	}
	// lookback blocks for transmit event is hard coded and should provide ample time for logs
	// to be detected in most cases
	var transmitLookbackBlocks int64 = 250
	transmitter, err := NewTransmitEventProvider(lggr, client.LogPoller(), addr, client.Client(), transmitLookbackBlocks)
	if err != nil {
		return nil, err
	}

	f := new(automationFacade)
	f.transmitter = transmitter

	packer := encoding.NewAbiPacker(keeperRegistryABI, utilsABI)
	f.encoder = encoding.NewReportEncoder(packer)

	f.logProvider, f.logRecoverer = logprovider.New(lggr, client.LogPoller(), utilsABI)

	f.blockSub = NewBlockSubscriber(client.HeadBroadcaster(), client.LogPoller(), lggr)

	f.payloadBuilder = NewPayloadBuilder(lggr)

	scanner := upkeepstate.NewPerformedEventsScanner(lggr, client.LogPoller(), addr)
	f.upkeepState = upkeepstate.NewUpkeepStateStore(lggr, scanner)

	f.keyring = NewOnchainKeyringV3Wrapper(keyring)

	al := NewActiveUpkeepList()

	f.reg = NewEvmRegistry(lggr, addr, client, feedLookupCompatibleABI,
		keeperRegistryABI, registryContract, mc, al, f.logProvider, f.encoder, packer)

	f.upkeepProvider = NewUpkeepProvider(al, f.reg, client.LogPoller())

	return f, nil
}

type automationFacade struct {
	reg            *EvmRegistry
	encoder        ocr2keepers.Encoder
	transmitter    *TransmitEventProvider
	blockSub       *BlockSubscriber
	payloadBuilder ocr2keepers.PayloadBuilder
	upkeepState    upkeepstate.UpkeepStateStore
	logProvider    logprovider.LogEventProvider
	logRecoverer   logprovider.LogRecoverer
	upkeepProvider *upkeepProvider
	keyring        *onchainKeyringV3Wrapper
}

var _ AutomationFacade = &automationFacade{}

func (f *automationFacade) Start(_ context.Context) error {
	// TODO: implement
	return nil
}

func (f *automationFacade) Close() error {
	// TODO: implement
	return nil
}

func (f *automationFacade) Registry() *EvmRegistry {
	return f.reg
}

func (f *automationFacade) Encoder() ocr2keepers.Encoder {
	return f.encoder
}

func (f *automationFacade) TransmitEventProvider() *TransmitEventProvider {
	return f.transmitter
}

func (f *automationFacade) BlockSubscriber() *BlockSubscriber {
	return f.blockSub
}

func (f *automationFacade) PayloadBuilder() ocr2keepers.PayloadBuilder {
	return f.payloadBuilder
}

func (f *automationFacade) UpkeepStateStore() upkeepstate.UpkeepStateStore {
	return f.upkeepState
}

func (f *automationFacade) LogEventProvider() logprovider.LogEventProvider {
	return f.logProvider
}

func (f *automationFacade) LogRecoverer() logprovider.LogRecoverer {
	return f.logRecoverer
}

func (f *automationFacade) UpkeepProvider() ocr2keepers.ConditionalUpkeepProvider {
	return f.upkeepProvider
}

func (f *automationFacade) Keyring() ocr3types.OnchainKeyring[plugin.AutomationReportInfo] {
	return f.keyring
}
