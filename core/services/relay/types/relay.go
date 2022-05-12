// Types are shared with external relay libraries so they can implement
// the interfaces required to run as a core OCR job.
package types

import (
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/services"
)

type Network string

var (
	EVM    Network = "evm"
	Solana Network = "solana"
	Terra  Network = "terra"
)

// PluginArgs are the args required to create any OCR2 plugin components.
// Its possible that the plugin config might actually be different
// per relay type, so we pass the config directly through.
type PluginArgs struct {
	ConfigWatcherArgs
	TransmitterID string
	PluginConfig  []byte
}

type ConfigWatcherArgs struct {
	ExternalJobID uuid.UUID
	JobID         int32
	ContractID    string
	RelayConfig   []byte
}

type Relayer interface {
	services.ServiceCtx
	NewConfigWatcher(args ConfigWatcherArgs) (ConfigWatcher, error)
	NewMedianProvider(args PluginArgs) (MedianProvider, error)
}

// The bootstrap jobs only watch config.
type ConfigWatcher interface {
	services.ServiceCtx
	OffchainConfigDigester() types.OffchainConfigDigester
	ContractConfigTracker() types.ContractConfigTracker
}

// OCR2Base provides common components for any OCR2 plugin.
// It watches config and is able to transmit.
type OCR2Base interface {
	ConfigWatcher
	ContractTransmitter() types.ContractTransmitter
}

// MedianProvider provides all components needed for a median OCR2 plugin.
type MedianProvider interface {
	OCR2Base
	ReportCodec() median.ReportCodec
	MedianContract() median.MedianContract
}
