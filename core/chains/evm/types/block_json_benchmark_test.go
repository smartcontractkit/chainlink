package types_test

import (
	"encoding/binary"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

func makeTestBlock(nTx int) *evmtypes.Block {
	txns := make([]evmtypes.Transaction, nTx)

	generateHash := func(x int64) common.Hash {
		out := make([]byte, 0, 32)

		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(x))

		for i := 0; i < 4; i++ {
			out = append(out, b...)
		}
		return common.BytesToHash(out)
	}
	for i := 0; i < nTx; i++ {
		wei := assets.NewWei(big.NewInt(int64(i)))
		txns[i] = evmtypes.Transaction{
			GasPrice:             wei,
			GasLimit:             uint32(i),
			MaxFeePerGas:         wei,
			MaxPriorityFeePerGas: wei,
			Type:                 evmtypes.TxType(i),
			Hash:                 generateHash(int64(i)),
		}
	}
	return &evmtypes.Block{
		Number:        int64(nTx),
		Hash:          generateHash(int64(1024 * 1024)),
		ParentHash:    generateHash(int64(512 * 1024)),
		BaseFeePerGas: assets.NewWei(big.NewInt(3)),
		Timestamp:     time.Unix(0, 0),
		Transactions:  txns,
	}
}

var (
	smallBlock  = makeTestBlock(2)
	mediumBlock = makeTestBlock(64)
	largeBlock  = makeTestBlock(512)
	xlBlock     = makeTestBlock(4 * 1024)
)

func unmarshal_block(b *testing.B, block *evmtypes.Block) {
	jsonBytes, err := json.Marshal(&block)
	if err != nil {
		b.Fatalf("failed to create test json %+v", err)
	}
	b.ResetTimer()

	var temp evmtypes.Block
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(jsonBytes, &temp)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
	}
}

func BenchmarkBlock_Small_JSONUnmarshal(b *testing.B) {
	unmarshal_block(b, smallBlock)

}

func BenchmarkBlock_Medium_JSONUnmarshal(b *testing.B) {
	unmarshal_block(b, mediumBlock)
}

func BenchmarkBlock_Large_JSONUnmarshal(b *testing.B) {
	unmarshal_block(b, largeBlock)
}

func BenchmarkBlock_XL_JSONUnmarshal(b *testing.B) {
	unmarshal_block(b, xlBlock)
}
