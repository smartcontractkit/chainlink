// Types are shared with external relay libraries so they can implement
// the interfaces required to run as a core OCR job.
package types

import (
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services"
)

type Network string

var (
	EVM    Network = "evm"
	Solana Network = "solana"
	Terra  Network = "terra"
)

// OCR2Args are the args required to create any OCR2 plugin provider.
// Its possible that the plugin config might actually be different
// per relay type, so we pass the config directly through.
type OCR2Args struct {
	ExternalJobID uuid.UUID
	JobID         int32
	ContractID    string
	TransmitterID null.String
	RelayConfig   map[string]interface{}
	PluginConfig  map[string]interface{}
	IsBootstrap   bool
}

// RelayerCtx represents a relayer
type RelayerCtx interface {
	services.ServiceCtx
	NewMedianProvider(args OCR2Args) (MedianProvider, error)
	// TODO: Will need some CCIP plugin providers for chain specific implementations
	// of request reading and tracking report status on dest chain.
	// For now, the ocr2/plugins/ccip is EVM specific.
}

// OCR2Provider provides common components for any OCR2 plugin.
type OCR2Provider interface {
	services.ServiceCtx
	ContractTransmitter() types.ContractTransmitter
	ContractConfigTracker() types.ContractConfigTracker
	OffchainConfigDigester() types.OffchainConfigDigester
}

// MedianProvider provides all components needed for a median OCR2 plugin.
type MedianProvider interface {
	OCR2Provider
	ReportCodec() median.ReportCodec
	MedianContract() median.MedianContract
}
