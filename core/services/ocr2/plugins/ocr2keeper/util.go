package ocr2keeper

import (
	"fmt"
	"math/big"

	kchain "github.com/smartcontractkit/ocr2keepers/pkg/chain"
	ktypes "github.com/smartcontractkit/ocr2keepers/pkg/types"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	kevm "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2keeper/evm"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	evmrelay "github.com/smartcontractkit/chainlink/core/services/relay/evm"
)

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

func EVMDependencies(spec job.Job, db *sqlx.DB, lggr logger.Logger, set evm.ChainSet, pr pipeline.Runner) (evmrelay.OCR2KeeperProvider, *kevm.EvmRegistry, ktypes.ReportEncoder, *LogProvider, error) {
	var err error
	var chain evm.Chain
	var keeperProvider evmrelay.OCR2KeeperProvider
	var registry *kevm.EvmRegistry
	var encoder ktypes.ReportEncoder

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
	if registry, err = kevm.NewEVMRegistryServiceV2_0(rAddr, chain, lggr); err != nil {
		return nil, nil, nil, nil, err
	}

	encoder = kchain.NewEVMReportEncoder()

	// lookback blocks is hard coded and should provide ample time for logs
	// to be detected in most cases
	var lookbackBlocks int64 = 250
	logProvider, err := NewLogProvider(lggr, chain.LogPoller(), rAddr, chain.Client(), lookbackBlocks)

	return keeperProvider, registry, encoder, logProvider, err
}
