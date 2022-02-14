// Types are shared with external relay libraries so they can implement
// the interfaces required to run as a core OCR job.
package types

import (
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/services"
)

type Network string

var (
	EVM    Network = "evm"
	Solana Network = "solana"
	Terra  Network = "terra"
)

// Relayer is network specific, but not job specific.
// It implements Checkable since want to block
// node traffic until the network wide services are up (Ready())
// and potentially restart the node if there are unhealthy network
// wide services (Health()), for example if node rpc connections
// are unhealthy.
type Relayer interface {
	services.Service
	// NewOCR2Provider creates an OCR2Provider for a given OCR job spec.
	// It signals on contractReady once contract configuration has been detected
	// and the OCR protocol may begin.
	NewOCR2Provider(externalJobID uuid.UUID, spec interface{}, contractReady chan<- struct{}) (OCR2Provider, error)
}

// OCR2Provider is network and job specific.
// It is a component of an OCR job and so
// only implements job.Service. We do not want to block
// the node traffic or restart the node based on one jobs health.
type OCR2Provider interface {
	Start() error
	Close() error
	ContractTransmitter() types.ContractTransmitter
	ContractConfigTracker() types.ContractConfigTracker
	OffchainConfigDigester() types.OffchainConfigDigester
	ReportCodec() median.ReportCodec
	MedianContract() median.MedianContract
}
