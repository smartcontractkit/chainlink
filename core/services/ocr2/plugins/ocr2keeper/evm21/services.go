package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/ocr2keepers/pkg/v3/plugin"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/logprovider"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/transmit"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/upkeepstate"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type AutomationServices interface {
	Registry() *EvmRegistry
	Encoder() ocr2keepers.Encoder
	TransmitEventProvider() *transmit.EventProvider
	BlockSubscriber() *BlockSubscriber
	PayloadBuilder() ocr2keepers.PayloadBuilder
	UpkeepStateStore() upkeepstate.UpkeepStateStore
	LogEventProvider() logprovider.LogEventProvider
	LogRecoverer() logprovider.LogRecoverer
	UpkeepProvider() ocr2keepers.ConditionalUpkeepProvider
	Keyring() ocr3types.OnchainKeyring[plugin.AutomationReportInfo]
}

func New(addr common.Address, client evm.Chain, mc *models.MercuryCredentials, keyring ocrtypes.OnchainKeyring, lggr logger.Logger, db *sqlx.DB, dbCfg pg.QConfig) (AutomationServices, error) {
	registryContract, err := iregistry21.NewIKeeperRegistryMaster(addr, client.Client())
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create caller for address and backend", ErrInitializationFailure)
	}
	// lookback blocks for transmit event is hard coded and should provide ample time for logs
	// to be detected in most cases
	var transmitLookbackBlocks int64 = 250
	transmitEventProvider, err := transmit.NewTransmitEventProvider(lggr, client.LogPoller(), addr, client.Client(), transmitLookbackBlocks)
	if err != nil {
		return nil, err
	}

	services := new(automationServices)
	services.transmitEventProvider = transmitEventProvider

	packer := encoding.NewAbiPacker()
	services.encoder = encoding.NewReportEncoder(packer)

	finalityDepth := client.Config().EVM().FinalityDepth()

	orm := upkeepstate.NewORM(client.ID(), db, lggr, dbCfg)
	scanner := upkeepstate.NewPerformedEventsScanner(lggr, client.LogPoller(), addr, finalityDepth)
	services.upkeepState = upkeepstate.NewUpkeepStateStore(orm, lggr, scanner)

	logProvider, logRecoverer := logprovider.New(lggr, client.LogPoller(), client.Client(), services.upkeepState, finalityDepth)
	services.logProvider = logProvider
	services.logRecoverer = logRecoverer
	services.blockSub = NewBlockSubscriber(client.HeadBroadcaster(), client.LogPoller(), finalityDepth, lggr)

	services.keyring = NewOnchainKeyringV3Wrapper(keyring)

	al := NewActiveUpkeepList()
	services.payloadBuilder = NewPayloadBuilder(al, logRecoverer, lggr)

	services.reg = NewEvmRegistry(lggr, addr, client,
		registryContract, mc, al, services.logProvider,
		packer, services.blockSub, finalityDepth)

	services.upkeepProvider = NewUpkeepProvider(al, services.blockSub, client.LogPoller())

	return services, nil
}

type automationServices struct {
	reg                   *EvmRegistry
	encoder               ocr2keepers.Encoder
	transmitEventProvider *transmit.EventProvider
	blockSub              *BlockSubscriber
	payloadBuilder        ocr2keepers.PayloadBuilder
	upkeepState           upkeepstate.UpkeepStateStore
	logProvider           logprovider.LogEventProvider
	logRecoverer          logprovider.LogRecoverer
	upkeepProvider        *upkeepProvider
	keyring               *onchainKeyringV3Wrapper
}

var _ AutomationServices = &automationServices{}

func (f *automationServices) Registry() *EvmRegistry {
	return f.reg
}

func (f *automationServices) Encoder() ocr2keepers.Encoder {
	return f.encoder
}

func (f *automationServices) TransmitEventProvider() *transmit.EventProvider {
	return f.transmitEventProvider
}

func (f *automationServices) BlockSubscriber() *BlockSubscriber {
	return f.blockSub
}

func (f *automationServices) PayloadBuilder() ocr2keepers.PayloadBuilder {
	return f.payloadBuilder
}

func (f *automationServices) UpkeepStateStore() upkeepstate.UpkeepStateStore {
	return f.upkeepState
}

func (f *automationServices) LogEventProvider() logprovider.LogEventProvider {
	return f.logProvider
}

func (f *automationServices) LogRecoverer() logprovider.LogRecoverer {
	return f.logRecoverer
}

func (f *automationServices) UpkeepProvider() ocr2keepers.ConditionalUpkeepProvider {
	return f.upkeepProvider
}

func (f *automationServices) Keyring() ocr3types.OnchainKeyring[plugin.AutomationReportInfo] {
	return f.keyring
}
