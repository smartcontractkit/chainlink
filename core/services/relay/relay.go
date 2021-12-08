package relay

import (
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type Network string

var (
	EVM               Network = "evm"
	Solana            Network = "solana"
	SupportedRelayers         = map[Network]struct{}{
		EVM:    {},
		Solana: {},
	}
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
