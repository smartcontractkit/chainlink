package relay

import (
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type Type string

var (
	TypeEthereum = Type("ethereum")
	TypeSolana   = Type("solana")
)

type Relayers map[Type]Relayer

type Relayer interface {
	service.Service
	NewOCR2Provider(config interface{}) (OCR2Provider, error)
}

type OCR2Provider interface {
	service.Service
	OffchainKeyring() types.OffchainKeyring
	OnchainKeyring() types.OnchainKeyring
	ContractTransmitter() types.ContractTransmitter
	ContractConfigTracker() types.ContractConfigTracker
	OffchainConfigDigester() types.OffchainConfigDigester
	ReportCodec() median.ReportCodec
	MedianContract() median.MedianContract
}
