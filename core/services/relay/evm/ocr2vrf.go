package evm

import (
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
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
	relaytypes.Relayer
	NewDKGProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (DKGProvider, error)
	NewOCR2VRFProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (OCR2VRFProvider, error)
}

// Relayer with added DKG and OCR2VRF provider functions.
type ocr2vrfRelayer struct {
	*Relayer
}

var _ OCR2VRFRelayer = (*ocr2vrfRelayer)(nil)

func NewOCR2VRFRelayer(relayer interface{}) OCR2VRFRelayer {
	return &ocr2vrfRelayer{relayer.(*Relayer)}
}

type dkgProvider struct {
	*configWatcher
	contractTransmitter *ContractTransmitter
}

var _ DKGProvider = (*dkgProvider)(nil)

func (r *ocr2vrfRelayer) NewDKGProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (DKGProvider, error) {
	configWatcher, err := newConfigProvider(r.lggr, r.chainSet, rargs)
	if err != nil {
		return nil, err
	}
	contractTransmitter, err := newContractTransmitter(r.lggr, rargs, pargs.TransmitterID, configWatcher)
	if err != nil {
		return nil, err
	}
	return &dkgProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
	}, nil
}

func (c *dkgProvider) ContractTransmitter() types.ContractTransmitter {
	return c.contractTransmitter
}

type vrfProvider struct {
	*configWatcher
	contractTransmitter *ContractTransmitter
}

var _ OCR2VRFProvider = (*vrfProvider)(nil)

func (r *ocr2vrfRelayer) NewOCR2VRFProvider(rargs relaytypes.RelayArgs, pargs relaytypes.PluginArgs) (OCR2VRFProvider, error) {
	configWatcher, err := newConfigProvider(r.lggr, r.chainSet, rargs)
	if err != nil {
		return nil, err
	}
	contractTransmitter, err := newContractTransmitter(r.lggr, rargs, pargs.TransmitterID, configWatcher)
	if err != nil {
		return nil, err
	}
	return &dkgProvider{
		configWatcher:       configWatcher,
		contractTransmitter: contractTransmitter,
	}, nil
}

func (c *vrfProvider) ContractTransmitter() types.ContractTransmitter {
	return c.contractTransmitter
}
