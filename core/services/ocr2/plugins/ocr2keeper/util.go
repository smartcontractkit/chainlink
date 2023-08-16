package ocr2keeper

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocr2keepers20 "github.com/smartcontractkit/ocr2keepers/pkg/v2"
	ocr2keepers20coordinator "github.com/smartcontractkit/ocr2keepers/pkg/v2/coordinator"
	ocr2keepers20polling "github.com/smartcontractkit/ocr2keepers/pkg/v2/observer/polling"
	ocr2keepers20runner "github.com/smartcontractkit/ocr2keepers/pkg/v2/runner"
	ocr2keepers21plugin "github.com/smartcontractkit/ocr2keepers/pkg/v3/plugin"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	ocr2keepers21 "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
	kevm20 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm20"
	kevm21 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/upkeepstate"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

type Encoder20 interface {
	ocr2keepers20.Encoder
	ocr2keepers20coordinator.Encoder
	ocr2keepers20polling.Encoder
	ocr2keepers20runner.Encoder
	ocr2keepers20coordinator.Encoder
}

type Encoder21 interface {
	ocr2keepers21.Encoder
}

var (
	ErrNoChainFromSpec = fmt.Errorf("could not create chain from spec")
)

func EVMProvider(db *sqlx.DB, chain evm.Chain, lggr logger.Logger, spec job.Job, pr pipeline.Runner) (evmrelay.OCR2KeeperProvider, error) {
	oSpec := spec.OCR2OracleSpec
	ocr2keeperRelayer := evmrelay.NewOCR2KeeperRelayer(db, chain, pr, spec, lggr.Named("OCR2KeeperRelayer"))

	keeperProvider, err := ocr2keeperRelayer.NewOCR2KeeperProvider(
		types.RelayArgs{
			ExternalJobID: spec.ExternalJobID,
			JobID:         oSpec.ID,
			ContractID:    oSpec.ContractID,
			RelayConfig:   oSpec.RelayConfig.Bytes(),
		},
		types.PluginArgs{
			TransmitterID: oSpec.TransmitterID.String,
			PluginConfig:  oSpec.PluginConfig.Bytes(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create new ocr2keeper provider", err)
	}

	return keeperProvider, nil
}

func EVMDependencies20(
	spec job.Job,
	db *sqlx.DB,
	lggr logger.Logger,
	set evm.ChainSet,
	pr pipeline.Runner,
) (evmrelay.OCR2KeeperProvider, *kevm20.EvmRegistry, Encoder20, *kevm20.LogProvider, error) {
	var err error
	var chain evm.Chain
	var keeperProvider evmrelay.OCR2KeeperProvider
	var registry *kevm20.EvmRegistry

	oSpec := spec.OCR2OracleSpec

	// get the chain from the config
	chainID, err2 := spec.OCR2OracleSpec.RelayConfig.EVMChainID()
	if err2 != nil {
		return nil, nil, nil, nil, err2
	}
	chain, err2 = set.Get(big.NewInt(chainID))
	if err2 != nil {
		return nil, nil, nil, nil, fmt.Errorf("%w: %s", ErrNoChainFromSpec, err2)
	}

	// the provider will be returned as a dependency
	if keeperProvider, err = EVMProvider(db, chain, lggr, spec, pr); err != nil {
		return nil, nil, nil, nil, err
	}

	rAddr := ethkey.MustEIP55Address(oSpec.ContractID).Address()
	if registry, err = kevm20.NewEVMRegistryService(rAddr, chain, lggr); err != nil {
		return nil, nil, nil, nil, err
	}

	encoder := kevm20.EVMAutomationEncoder20{}

	// lookback blocks is hard coded and should provide ample time for logs
	// to be detected in most cases
	var lookbackBlocks int64 = 250
	// TODO: accept a version of the registry contract and use the correct interfaces
	logProvider, err := kevm20.NewLogProvider(lggr, chain.LogPoller(), rAddr, chain.Client(), lookbackBlocks)

	return keeperProvider, registry, encoder, logProvider, err
}

func FilterNamesFromSpec20(spec *job.OCR2OracleSpec) (names []string, err error) {
	addr, err := ethkey.NewEIP55Address(spec.ContractID)
	if err != nil {
		return nil, err
	}
	return []string{kevm20.LogProviderFilterName(addr.Address()), kevm20.UpkeepFilterName(addr.Address())}, err
}

func EVMDependencies21(
	spec job.Job,
	db *sqlx.DB,
	lggr logger.Logger,
	set evm.ChainSet,
	pr pipeline.Runner,
	mc *models.MercuryCredentials,
	keyring ocrtypes.OnchainKeyring,
) (evmrelay.OCR2KeeperProvider, *kevm21.EvmRegistry, Encoder21, *kevm21.TransmitEventProvider, ocr2keepers21.LogEventProvider, ocr3types.OnchainKeyring[ocr2keepers21plugin.AutomationReportInfo], *kevm21.BlockSubscriber, ocr2keepers21.PayloadBuilder, ocr2keepers21.UpkeepStateUpdater, ocr2keepers21.ConditionalUpkeepProvider, error) {
	var err error
	var chain evm.Chain
	var keeperProvider evmrelay.OCR2KeeperProvider
	var registry *kevm21.EvmRegistry
	var encoder *kevm21.EVMAutomationEncoder21

	oSpec := spec.OCR2OracleSpec

	// get the chain from the config
	chainID, err2 := spec.OCR2OracleSpec.RelayConfig.EVMChainID()
	if err2 != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err2
	}
	chain, err2 = set.Get(big.NewInt(chainID))
	if err2 != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("%w: %s", ErrNoChainFromSpec, err2)
	}

	// the provider will be returned as a dependency
	if keeperProvider, err = EVMProvider(db, chain, lggr, spec, pr); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	bs := kevm21.NewBlockSubscriber(chain.HeadBroadcaster(), chain.LogPoller(), lggr)
	rAddr := ethkey.MustEIP55Address(oSpec.ContractID).Address()
	if registry, encoder, err = kevm21.NewEVMRegistryService(rAddr, chain, mc, bs, lggr); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	pb := kevm21.NewPayloadBuilder(lggr)
	scanner := upkeepstate.NewPerformedEventsScanner(
		lggr,
		chain.LogPoller(),
		rAddr,
	)
	us := upkeepstate.NewUpkeepStateStore(lggr, scanner)
	up := kevm21.NewUpkeepProvider(registry, chain.LogPoller())

	// lookback blocks is hard coded and should provide ample time for logs
	// to be detected in most cases
	var lookbackBlocks int64 = 250
	// TODO: accept a version of the registry contract and use the correct interfaces
	logTransmitter, err := kevm21.NewTransmitEventProvider(lggr, chain.LogPoller(), rAddr, chain.Client(), lookbackBlocks)

	return keeperProvider, registry, encoder, logTransmitter, registry.LogEventProvider(), kevm21.NewOnchainKeyringV3Wrapper(keyring), bs, pb, us, up, err
}

func FilterNamesFromSpec21(spec *job.OCR2OracleSpec) (names []string, err error) {
	addr, err := ethkey.NewEIP55Address(spec.ContractID)
	if err != nil {
		return nil, err
	}
	return []string{kevm21.TransmitEventProviderFilterName(addr.Address()), kevm21.UpkeepFilterName(addr.Address())}, err
}
