package ocr2keeper

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"

	"github.com/jmoiron/sqlx"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	ocr2keepers20 "github.com/smartcontractkit/chainlink-automation/pkg/v2"
	ocr2keepers20coordinator "github.com/smartcontractkit/chainlink-automation/pkg/v2/coordinator"
	ocr2keepers20polling "github.com/smartcontractkit/chainlink-automation/pkg/v2/observer/polling"
	ocr2keepers20runner "github.com/smartcontractkit/chainlink-automation/pkg/v2/runner"
	ocr2keepers21 "github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
	evmregistry20 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v20"
	evmregistry21 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21"
	evmregistry21transmit "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/transmit"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
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

func EVMProvider(db *sqlx.DB, chain legacyevm.Chain, lggr logger.Logger, spec job.Job, ethKeystore keystore.Eth) (evmrelay.OCR2KeeperProvider, error) {
	oSpec := spec.OCR2OracleSpec
	ocr2keeperRelayer := evmrelay.NewOCR2KeeperRelayer(db, chain, lggr.Named("OCR2KeeperRelayer"), ethKeystore)

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
	chain legacyevm.Chain,
	ethKeystore keystore.Eth,
) (evmrelay.OCR2KeeperProvider, *evmregistry20.EvmRegistry, Encoder20, *evmregistry20.LogProvider, error) {
	var err error

	var keeperProvider evmrelay.OCR2KeeperProvider
	var registry *evmregistry20.EvmRegistry

	// the provider will be returned as a dependency
	if keeperProvider, err = EVMProvider(db, chain, lggr, spec, ethKeystore); err != nil {
		return nil, nil, nil, nil, err
	}

	rAddr := ethkey.MustEIP55Address(spec.OCR2OracleSpec.ContractID).Address()
	if registry, err = evmregistry20.NewEVMRegistryService(rAddr, chain, lggr); err != nil {
		return nil, nil, nil, nil, err
	}

	encoder := evmregistry20.EVMAutomationEncoder20{}

	// lookback blocks is hard coded and should provide ample time for logs
	// to be detected in most cases
	var lookbackBlocks int64 = 250
	// TODO: accept a version of the registry contract and use the correct interfaces
	logProvider, err := evmregistry20.NewLogProvider(lggr, chain.LogPoller(), rAddr, chain.Client(), lookbackBlocks)

	return keeperProvider, registry, encoder, logProvider, err
}

func FilterNamesFromSpec20(spec *job.OCR2OracleSpec) (names []string, err error) {
	addr, err := ethkey.NewEIP55Address(spec.ContractID)
	if err != nil {
		return nil, err
	}
	return []string{evmregistry20.LogProviderFilterName(addr.Address()), evmregistry20.UpkeepFilterName(addr.Address())}, err
}

func EVMDependencies21(
	spec job.Job,
	db *sqlx.DB,
	lggr logger.Logger,
	chain legacyevm.Chain,
	mc *models.MercuryCredentials,
	keyring ocrtypes.OnchainKeyring,
	dbCfg pg.QConfig,
) (evmregistry21.AutomationServices, error) {
	var err error

	oSpec := spec.OCR2OracleSpec

	rAddr := ethkey.MustEIP55Address(oSpec.ContractID).Address()
	services, err := evmregistry21.New(rAddr, chain, mc, keyring, lggr, db, dbCfg)
	if err != nil {
		return nil, err
	}

	return services, err
}

func FilterNamesFromSpec21(spec *job.OCR2OracleSpec) (names []string, err error) {
	addr, err := ethkey.NewEIP55Address(spec.ContractID)
	if err != nil {
		return nil, err
	}
	return []string{evmregistry21transmit.EventProviderFilterName(addr.Address()), evmregistry21.RegistryUpkeepFilterName(addr.Address())}, err
}
