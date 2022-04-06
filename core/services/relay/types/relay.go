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

	// Custom Endpoint just posts price feeds to a custom API endpoint
	// which the endpoint controls. There's no Chain for such cases.
	// The core OCR2 components for customendpoint are written to ensure that
	// the protocol works fine, and able to report accurate telemetry data.
	CustomEndpoint Network = "customendpoint"
)

// RelayerCtx represents a relayer
type RelayerCtx interface {
	services.ServiceCtx
	// NewOCR2Provider is generic for all OCR2 plugins on the given chain.
	NewOCR2Provider(externalJobID uuid.UUID, spec interface{}) (OCR2ProviderCtx, error)
	// TODO: Will need some CCIP plugin providers for chain specific implementations
	// of request reading and tracking report status on dest chain.
	// For now, the ocr2/plugins/ccip is EVM specific.
}

// OCR2ProviderCtx contains methods needed for job.OCR2OracleSpec functionality
type OCR2ProviderCtx interface {
	services.ServiceCtx
	ContractTransmitter() types.ContractTransmitter
	ContractConfigTracker() types.ContractConfigTracker
	OffchainConfigDigester() types.OffchainConfigDigester
	OCR2MedianProvider
}

// OCR2MedianProvider contains methods needed for the median.Median plugin
type OCR2MedianProvider interface {
	ReportCodec() median.ReportCodec
	MedianContract() median.MedianContract
}
