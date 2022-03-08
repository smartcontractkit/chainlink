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

type Relayer interface {
	services.Service
	// Generic for all OCR2 plugins on the given chain.
	NewOCR2Provider(externalJobID uuid.UUID, spec interface{}) (OCR2Provider, error)
	// TODO: Will need some CCIP plugin providers for chain specific implementations
	// of request reading and tracking report status on dest chain.
	// For now, the ocr2/plugins/ccip is EVM specific.
}

// RelayerCtx is replacing Relayer interface
type RelayerCtx interface {
	services.ServiceCtx
	NewOCR2Provider(externalJobID uuid.UUID, spec interface{}) (OCR2ProviderCtx, error)
}

// OCR2Provider contains methods needed for job.OCR2OracleSpec functionality
type OCR2Provider interface {
	services.Service
	ContractTransmitter() types.ContractTransmitter
	ContractConfigTracker() types.ContractConfigTracker
	OffchainConfigDigester() types.OffchainConfigDigester
	OCR2MedianProvider
}

// OCR2ProviderCtx is replacing OCR2Provider interface
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
