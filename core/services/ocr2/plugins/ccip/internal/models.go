package internal

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// EVM2EVMOnRampCCIPSendRequestedWithMeta helper struct to hold the send request and some metadata
type EVM2EVMOnRampCCIPSendRequestedWithMeta struct {
	EVM2EVMMessage
	BlockTimestamp time.Time
	Executed       bool
	Finalized      bool
	LogIndex       uint
	TxHash         common.Hash
}

type Hash [32]byte

func (h Hash) String() string {
	return hexutil.Encode(h[:])
}

type TokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

// EVM2EVMMessage is the interface for a message sent from the offramp to the onramp
// Plugin can operate against any lane version which has a message satisfying this interface.
type EVM2EVMMessage struct {
	SequenceNumber      uint64
	GasLimit            *big.Int
	Nonce               uint64
	MessageId           Hash
	SourceChainSelector uint64
	Sender              common.Address
	Receiver            common.Address
	Strict              bool
	FeeToken            common.Address
	FeeTokenAmount      *big.Int
	Data                []byte
	TokenAmounts        []TokenAmount
	SourceTokenData     [][]byte

	// Computed
	Hash Hash
}
