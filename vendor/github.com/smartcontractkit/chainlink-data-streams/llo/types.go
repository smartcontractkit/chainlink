package llo

import (
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

type ChannelDefinitionWithID struct {
	llotypes.ChannelDefinition
	ChannelID llotypes.ChannelID
}

type ChannelHash [32]byte

type OnchainConfigCodec interface {
	Encode(OnchainConfig) ([]byte, error)
	Decode([]byte) (OnchainConfig, error)
}

type Transmitter interface {
	// NOTE: Mercury doesn't actually transmit on-chain, so there is no
	// "contract" involved with the transmitter.
	// - Transmit should be implemented and send to Mercury server
	// - FromAccount() should return CSA public key
	ocr3types.ContractTransmitter[llotypes.ReportInfo]
}
