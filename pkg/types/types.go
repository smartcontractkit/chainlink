package types

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type Service interface {
	Start(context.Context) error
	Close() error
	Ready() error
	Healthy() error
}

// PluginArgs are the args required to create any OCR2 plugin components.
// Its possible that the plugin config might actually be different
// per relay type, so we pass the config directly through.
type PluginArgs struct {
	TransmitterID string
	PluginConfig  []byte
}

type RelayArgs struct {
	ExternalJobID uuid.UUID
	JobID         int32
	ContractID    string
	RelayConfig   []byte
}

type Relayer interface {
	Service
	NewConfigProvider(rargs RelayArgs) (ConfigProvider, error)
	NewMedianProvider(rargs RelayArgs, pargs PluginArgs) (MedianProvider, error)
}

// The bootstrap jobs only watch config.
type ConfigProvider interface {
	Service
	OffchainConfigDigester() types.OffchainConfigDigester
	ContractConfigTracker() types.ContractConfigTracker
}

// Plugin provides common components for any OCR2 plugin.
// It watches config and is able to transmit.
type Plugin interface {
	ConfigProvider
	ContractTransmitter() types.ContractTransmitter
}

// MedianProvider provides all components needed for a median OCR2 plugin.
type MedianProvider interface {
	Plugin
	ReportCodec() median.ReportCodec
	MedianContract() median.MedianContract
}

// OCR2VRFRelayer contains the relayer and instantiating functions for OCR2VRF providers.
type OCR2VRFRelayer interface {
	Relayer
	NewDKGProvider(rargs RelayArgs, transmitterID string) (DKGProvider, error)
	NewOCR2VRFProvider(rargs RelayArgs, transmitterID string) (OCR2VRFProvider, error)
}

// DKGProvider provides all components needed for a DKG plugin.
type DKGProvider interface {
	Plugin
}

// OCR2VRFProvider provides all components needed for a OCR2VRF plugin.
type OCR2VRFProvider interface {
	Plugin
}
