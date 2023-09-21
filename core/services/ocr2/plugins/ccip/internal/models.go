package internal

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
)

// EVM2EVMOnRampCCIPSendRequestedWithMeta helper struct to hold the send request and some metadata
type EVM2EVMOnRampCCIPSendRequestedWithMeta struct {
	evm_2_evm_offramp.InternalEVM2EVMMessage
	BlockTimestamp time.Time
	Executed       bool
	Finalized      bool
	LogIndex       uint
	TxHash         common.Hash
}
