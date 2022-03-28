package logpoller

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/utils"
)

type LogPollerBlock struct {
	EvmChainId *utils.Big
	BlockHash common.Hash
	// TODO: Following existing data types for block numbers, though not sure why we don't use big.Ints elsewhere?
	// Seems conceivable we could run out on some new EVM chain.
	BlockNumber int64
	CreatedAt time.Time
}

type Log struct {
	EvmChainId *utils.Big
	LogIndex int64
	BlockHash common.Hash
	BlockNumber int64
	EventSignature string
	Address common.Address
	TxHash common.Hash
	Data []byte
	CreatedAt time.Time
}