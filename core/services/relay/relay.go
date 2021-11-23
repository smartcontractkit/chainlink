package relay

import (
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type Relays struct {
	// TODO: support Ethereum as a relay
	// Ethereum Relay
	Solana Relay
}

type Relay interface {
	service.Service
	NewOCR2Service(config interface{}) OCR2Service
}

type OCR2Service interface {
	service.Service
	ContractTransmitter() types.ContractTransmitter
	ContractConfigTracker() types.ContractConfigTracker
	OffchainConfigDigester() types.OffchainConfigDigester
	ReportCodec() median.ReportCodec
	MedianContract() median.MedianContract
}
