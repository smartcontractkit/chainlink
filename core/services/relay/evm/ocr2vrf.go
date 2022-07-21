package evm

import (
	"encoding/json"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/sqlx"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/dkg/config"
)

// DKGProvider provides all components needed for a DKG plugin.
type DKGProvider interface {
	relaytypes.Plugin
}

// OCR2VRFProvider provides all components needed for a OCR2VRF plugin.
type OCR2VRFProvider interface {
	relaytypes.Plugin
}

// OCR2VRFRelayer contains the relayer and instantiating functions for OCR2VRF providers.
type OCR2VRFRelayer interface {
	NewDKGProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (DKGProvider, error)
	NewOCR2VRFProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (OCR2VRFProvider, error)
}

var (
	_ OCR2VRFRelayer  = (*ocr2vrfRelayer)(nil)
	_ DKGProvider     = (*dkgProvider)(nil)
	_ OCR2VRFProvider = (*ocr2vrfProvider)(nil)
)

// Relayer with added DKG and OCR2VRF provider functions.
type ocr2vrfRelayer struct {
	db       *sqlx.DB
	chainSet evm.ChainSet
	lggr     logger.Logger
}

func NewOCR2VRFRelayer(db *sqlx.DB, chainSet evm.ChainSet, lggr logger.Logger) OCR2VRFRelayer {
	return &ocr2vrfRelayer{
		db:       db,
		chainSet: chainSet,
		lggr:     lggr,
	}
}

func (r *ocr2vrfRelayer) NewDKGProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (DKGProvider, error) {
	configWatcher, err := newConfigProvider(r.lggr, r.chainSet, rargs)
	if err != nil {
		return nil, err
	}
	contractTransmitter, err := newContractTransmitter(r.lggr, rargs, pargs.TransmitterID, configWatcher)
	if err != nil {
		return nil, err
	}

	var pluginConfig config.PluginConfig
	err = json.Unmarshal(pargs.PluginConfig, &pluginConfig)
	if err != nil {
		return nil, err
	}

	return &dkgProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
		pluginConfig:        pluginConfig,
	}, nil
}

func (r *ocr2vrfRelayer) NewOCR2VRFProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (OCR2VRFProvider, error) {
	configWatcher, err := newConfigProvider(r.lggr, r.chainSet, rargs)
	if err != nil {
		return nil, err
	}
	contractTransmitter, err := newContractTransmitter(r.lggr, rargs, pargs.TransmitterID, configWatcher)
	if err != nil {
		return nil, err
	}
	return &ocr2vrfProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
	}, nil
}

type dkgProvider struct {
	*configWatcher
	contractTransmitter *ContractTransmitter
	pluginConfig        config.PluginConfig
}

func (c *dkgProvider) ContractTransmitter() types.ContractTransmitter {
	return c.contractTransmitter
}

type ocr2vrfProvider struct {
	*configWatcher
	contractTransmitter *ContractTransmitter
}

func (c *ocr2vrfProvider) ContractTransmitter() types.ContractTransmitter {
	return c.contractTransmitter
}
