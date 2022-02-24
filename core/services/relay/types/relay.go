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
	NewOCR2Provider(externalJobID uuid.UUID, spec interface{}) (OCR2Provider, error)
}

// OCR2Provider contains methods needed for generic job.OCR2OracleSpec functionality
type OCR2Provider interface {
	services.Service
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
