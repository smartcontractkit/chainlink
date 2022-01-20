// Types are shared with external relay libraries so they can implement
// the interfaces required to run as a core OCR job.
package types

import (
	"fmt"

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

type Relayer interface {
	services.Service
	NewOCR2Provider(externalJobID uuid.UUID, spec interface{}) (OCR2Provider, error)
}

type OCR2Provider interface {
	services.Service
	ContractTransmitter() types.ContractTransmitter
	ContractConfigTracker() types.ContractConfigTracker
	OffchainConfigDigester() types.OffchainConfigDigester
	ReportCodec() median.ReportCodec
	MedianContract() median.MedianContract
}

// DisabledRelayer returns a Relayer for n which returns disabled errors instead of OCR2Provider.
func DisabledRelayer(n Network) Relayer {
	return &disabledRelayer{n}
}

type disabledRelayer struct {
	n Network
}

func (d disabledRelayer) Start() error {
	return nil
}

func (d disabledRelayer) Close() error {
	return nil
}

func (d disabledRelayer) Ready() error {
	return nil
}

func (d disabledRelayer) Healthy() error {
	return nil
}

func (d disabledRelayer) NewOCR2Provider(_ uuid.UUID, _ interface{}) (OCR2Provider, error) {
	return nil, fmt.Errorf("relayer disabled for %s network", d.n)
}
