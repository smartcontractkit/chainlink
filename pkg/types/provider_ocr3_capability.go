package types

import "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

type OCR3CapabilityProvider interface {
	PluginProvider
	OCR3ContractTransmitter() ocr3types.ContractTransmitter[[]byte]
}
