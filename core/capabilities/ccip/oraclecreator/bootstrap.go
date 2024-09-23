package oraclecreator

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"
	libocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

var _ cctypes.OracleCreator = &bootstrapOracleCreator{}

type bootstrapOracleCreator struct {
	peerWrapper           *ocrcommon.SingletonPeerWrapper
	bootstrapperLocators  []commontypes.BootstrapperLocator
	db                    ocr3types.Database
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	lggr                  logger.Logger
}

func NewBootstrapOracleCreator(
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	bootstrapperLocators []commontypes.BootstrapperLocator,
	db ocr3types.Database,
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator,
	lggr logger.Logger,
) cctypes.OracleCreator {
	return &bootstrapOracleCreator{
		peerWrapper:           peerWrapper,
		bootstrapperLocators:  bootstrapperLocators,
		db:                    db,
		monitoringEndpointGen: monitoringEndpointGen,
		lggr:                  lggr,
	}
}

// Type implements types.OracleCreator.
func (i *bootstrapOracleCreator) Type() cctypes.OracleType {
	return cctypes.OracleTypeBootstrap
}

// Create implements types.OracleCreator.
func (i *bootstrapOracleCreator) Create(_ uint32, config cctypes.OCR3ConfigWithMeta) (cctypes.CCIPOracle, error) {
	// Assuming that the chain selector is referring to an evm chain for now.
	// TODO: add an api that returns chain family.
	// NOTE: this doesn't really matter for the bootstrap node, it doesn't do anything on-chain.
	// Its for the monitoring endpoint generation below.
	chainID, err := chainsel.ChainIdFromSelector(uint64(config.Config.ChainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID from selector: %w", err)
	}

	destChainFamily := chaintype.EVM
	destRelayID := types.NewRelayID(string(destChainFamily), fmt.Sprintf("%d", chainID))

	bootstrapperArgs := libocr3.BootstrapperArgs{
		BootstrapperFactory:   i.peerWrapper.Peer2,
		V2Bootstrappers:       i.bootstrapperLocators,
		ContractConfigTracker: ocrimpls.NewConfigTracker(config),
		Database:              i.db,
		LocalConfig:           defaultLocalConfig(),
		Logger: ocrcommon.NewOCRWrapper(
			i.lggr.
				Named("CCIPBootstrap").
				Named(destRelayID.String()).
				Named(config.Config.ChainSelector.String()).
				Named(hexutil.Encode(config.Config.OfframpAddress)),
			false, /* traceLogging */
			func(ctx context.Context, msg string) {}),
		MonitoringEndpoint: i.monitoringEndpointGen.GenMonitoringEndpoint(
			string(destChainFamily),
			destRelayID.ChainID,
			hexutil.Encode(config.Config.OfframpAddress),
			synchronization.OCR3CCIPBootstrap,
		),
		OffchainConfigDigester: ocrimpls.NewConfigDigester(config.ConfigDigest),
	}
	bootstrapper, err := libocr3.NewBootstrapper(bootstrapperArgs)
	if err != nil {
		return nil, err
	}
	return bootstrapper, nil
}
