package blocks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//go:generate codecgen -o internal_types_codecgen.go  -j true -d 1709 transactions.go block.go

// BlockInternal is JSON-serialization optimized intermediate representation between EVM blocks
// and our public representation
type BlockInternal struct {
	Number        string                `json:"number"`
	Hash          common.Hash           `json:"hash"`
	ParentHash    common.Hash           `json:"parentHash"`
	BaseFeePerGas *hexutil.Big          `json:"baseFeePerGas"`
	Timestamp     hexutil.Uint64        `json:"timestamp"`
	Transactions  []TransactionInternal `json:"transactions"`
}

func (bi BlockInternal) Empty() bool {
	var dflt BlockInternal

	return len(bi.Transactions) == 0 &&
		bi.Hash == dflt.Hash &&
		bi.ParentHash == dflt.ParentHash &&
		bi.BaseFeePerGas == dflt.BaseFeePerGas &&
		bi.Timestamp == dflt.Timestamp
}
