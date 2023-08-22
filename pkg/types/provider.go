package types

import ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

// The bootstrap jobs only watch config.
type ConfigProvider interface {
	Service
	OffchainConfigDigester() ocrtypes.OffchainConfigDigester
	ContractConfigTracker() ocrtypes.ContractConfigTracker
}

// Plugin is an alias for PluginProvider, for compatibility.
// Deprecated
type Plugin = PluginProvider

// PluginProvider provides common components for any OCR2 plugin.
// It watches config and is able to transmit.
type PluginProvider interface {
	ConfigProvider
	ContractTransmitter() ocrtypes.ContractTransmitter
}
