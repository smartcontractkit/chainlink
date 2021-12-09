// Types are shared with external relay libraries so they can implement
// the interfaces required to run as a core OCR job.
package types

import (
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/service"
)

type Network string

var (
	EVM    Network = "evm"
	Solana Network = "solana"
)

type Relayer interface {
	service.Service
	NewOCR2Provider(externalJobID uuid.UUID, spec interface{}) (OCR2Provider, error)
}

type OCR2Provider interface {
	service.Service
	ContractTransmitter() types.ContractTransmitter
	ContractConfigTracker() types.ContractConfigTracker
	OffchainConfigDigester() types.OffchainConfigDigester
	ReportCodec() median.ReportCodec
	MedianContract() median.MedianContract
}
