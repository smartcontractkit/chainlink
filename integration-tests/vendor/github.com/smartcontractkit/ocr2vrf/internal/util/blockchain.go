package util

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/downloader"
)

type Blockchain interface {
	downloader.BlockChain
	GetBlockByNumber(number uint64) *types.Block
}
