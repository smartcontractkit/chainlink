package ocr2keeper

import (
	"fmt"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocr2keepers20 "github.com/smartcontractkit/ocr2keepers/pkg/v2"
	ocr2keepers20coordinator "github.com/smartcontractkit/ocr2keepers/pkg/v2/coordinator"
	ocr2keepers20polling "github.com/smartcontractkit/ocr2keepers/pkg/v2/observer/polling"
	ocr2keepers20runner "github.com/smartcontractkit/ocr2keepers/pkg/v2/runner"
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
	kevm21transmit "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/transmit"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
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
	chain evm.Chain,
	pr pipeline.Runner,
) (evmrelay.OCR2KeeperProvider, *kevm20.EvmRegistry, Encoder20, *kevm20.LogProvider, error) {
	var err error

	var keeperProvider evmrelay.OCR2KeeperProvider
	var registry *kevm20.EvmRegistry

	// the provider will be returned as a dependency
	if keeperProvider, err = EVMProvider(db, chain, lggr, spec, pr); err != nil {
		return nil, nil, nil, nil, err
	}

	rAddr := ethkey.MustEIP55Address(spec.OCR2OracleSpec.ContractID).Address()
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
	chain evm.Chain,
	pr pipeline.Runner,
	mc *models.MercuryCredentials,
	keyring ocrtypes.OnchainKeyring,
	dbCfg pg.QConfig,
) (evmrelay.OCR2KeeperProvider, kevm21.AutomationServices, error) {
	var err error
	var keeperProvider evmrelay.OCR2KeeperProvider

	oSpec := spec.OCR2OracleSpec
	// the provider will be returned as a dependency
	if keeperProvider, err = EVMProvider(db, chain, lggr, spec, pr); err != nil {
		return nil, nil, err
	}

	rAddr := ethkey.MustEIP55Address(oSpec.ContractID).Address()
	services, err := kevm21.New(rAddr, chain, mc, keyring, lggr, db, dbCfg)
	if err != nil {
		return nil, nil, err
	}

	return keeperProvider, services, err
}

func FilterNamesFromSpec21(spec *job.OCR2OracleSpec) (names []string, err error) {
	addr, err := ethkey.NewEIP55Address(spec.ContractID)
	if err != nil {
		return nil, err
	}
	return []string{kevm21transmit.EventProviderFilterName(addr.Address()), kevm21.RegistryUpkeepFilterName(addr.Address())}, err
}
