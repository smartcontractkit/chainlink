package helpers

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type CommitStoreReportAccepted struct {
	Report ICommitStoreCommitReport
	Raw    types.Log
}

type ICommitStoreCommitReport struct {
	PriceUpdates InternalPriceUpdates
	Interval     ICommitStoreInterval
	MerkleRoot   [32]byte
}

type InternalPriceUpdates struct {
	TokenPriceUpdates []InternalTokenPriceUpdate
	DestChainId       uint64
	UsdPerUnitGas     *big.Int
}

type InternalTokenPriceUpdate struct {
	SourceToken common.Address
	UsdPerToken *big.Int
}

type ICommitStoreInterval struct {
	Min uint64
	Max uint64
}

type InternalEVM2EVMMessage struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	FeeTokenAmount      *big.Int
	Sender              common.Address
	Nonce               uint64
	GasLimit            *big.Int
	Strict              bool
	Receiver            common.Address
	Data                []byte
	TokenAmounts        []ClientEVMTokenAmount
	FeeToken            common.Address
	MessageId           [32]byte
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

type SendRequestedEvent struct {
	Message InternalEVM2EVMMessage
	Raw     types.Log
}

type InternalExecutionReport struct {
	Messages          []InternalEVM2EVMMessage
	OffchainTokenData [][][]byte
	Proofs            [][32]byte
	ProofFlagBits     *big.Int
}

type EVM2EVMOffRampExecutionStateChanged struct {
	SequenceNumber uint64
	MessageId      [32]byte
	State          uint8
	ReturnData     []byte
	Raw            types.Log
}
