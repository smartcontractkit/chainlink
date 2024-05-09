package evm

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type ocr3CapabilityProvider struct {
	types.PluginProvider
	transmitter ocr3types.ContractTransmitter[[]byte]
}

func (o *ocr3CapabilityProvider) OCR3ContractTransmitter() ocr3types.ContractTransmitter[[]byte] {
	return o.transmitter
}
