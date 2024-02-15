package ccip

import (
	"context"
	"math/big"
)

type OnRampDynamicConfig struct {
	Router                            Address
	MaxNumberOfTokensPerMsg           uint16
	DestGasOverhead                   uint32
	DestGasPerPayloadByte             uint16
	DestDataAvailabilityOverheadGas   uint32
	DestGasPerDataAvailabilityByte    uint16
	DestDataAvailabilityMultiplierBps uint16
	PriceRegistry                     Address
	MaxDataBytes                      uint32
	MaxPerMsgGasLimit                 uint32
}

// EVM2EVMMessage is the interface for a message sent from the off-ramp to the on-ramp
// Plugin can operate against any lane version which has a message satisfying this interface.
type EVM2EVMMessage struct {
	SequenceNumber      uint64
	GasLimit            *big.Int
	Nonce               uint64
	MessageID           Hash
	SourceChainSelector uint64
	Sender              Address
	Receiver            Address
	Strict              bool
	FeeToken            Address
	FeeTokenAmount      *big.Int
	Data                []byte
	TokenAmounts        []TokenAmount
	SourceTokenData     [][]byte

	// Computed
	Hash Hash
}

type EVM2EVMMessageWithTxMeta struct {
	TxMeta
	EVM2EVMMessage
}

type TokenAmount struct {
	Token  Address
	Amount *big.Int
}

type OnRampReader interface {
	// GetSendRequestsBetweenSeqNums returns all the finalized message send requests in the provided sequence numbers range (inclusive).
	GetSendRequestsBetweenSeqNums(ctx context.Context, seqNumMin, seqNumMax uint64, finalized bool) ([]EVM2EVMMessageWithTxMeta, error)

	// RouterAddress returns the router address that is configured on the onRamp
	RouterAddress() (Address, error)

	Address() (Address, error)

	GetDynamicConfig() (OnRampDynamicConfig, error)
}
