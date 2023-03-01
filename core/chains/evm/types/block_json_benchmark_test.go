package types

import (
	"encoding/binary"
	"encoding/json"
	"math/big"
	"os"
	"testing"
	"time"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ugorji/go/codec"

	"github.com/smartcontractkit/chainlink/core/assets"
)

func makeTestBlock(nTx int) *Block {
	txns := make([]Transaction, nTx)

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
		txns[i] = Transaction{
			GasPrice:             wei,
			GasLimit:             uint32(i),
			MaxFeePerGas:         wei,
			MaxPriorityFeePerGas: wei,
			Type:                 0,
			Hash:                 generateHash(int64(i)),
		}
	}
	return &Block{
		Number:        int64(nTx),
		Hash:          generateHash(int64(1024 * 1024)),
		ParentHash:    generateHash(int64(512 * 1024)),
		BaseFeePerGas: assets.NewWei(big.NewInt(3)),
		Timestamp:     time.Now(),
		Transactions:  txns,
	}
}

var (
	smallBlock  = makeTestBlock(2)
	mediumBlock = makeTestBlock(64)
	largeBlock  = makeTestBlock(512)
)

var expectedSize = unsafe.Sizeof(Block{})

func resetHandler(blocks ...*Block) {
	for _, b := range blocks {
		b.WithJSONCoder(StdLib)
	}
}

func BenchmarkBlock_Small_JSONMarshal(b *testing.B) {
	defer resetHandler(smallBlock)
	smallBlock.WithJSONCoder(StdLib)
	for i := 0; i < b.N; i++ {
		bt, err := json.Marshal(&smallBlock)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
		if bt == nil {
			b.Fatal("nil buf")
		}
	}
}

func BenchmarkBlock_Small_CodecMarshal(b *testing.B) {
	defer resetHandler(smallBlock)
	buf := make([]byte, 0, expectedSize)
	var h codec.Handle = new(codec.JsonHandle)
	enc := codec.NewEncoderBytes(&buf, h)
	smallBlock.WithJSONCoder(GoCodec)
	for i := 0; i < b.N; i++ {
		err := enc.Encode(&smallBlock)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
		buf = buf[:0]
	}
}

func BenchmarkBlock_Small_EasyJSONMarshal(b *testing.B) {
	defer resetHandler(smallBlock)
	smallBlock.WithJSONCoder(EasyJson)
	for i := 0; i < b.N; i++ {
		bt, err := json.Marshal(&smallBlock)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
		if bt == nil {
			b.Fatal("nil buf")
		}
	}
}
func BenchmarkBlock_Medium_JSONMarshal(b *testing.B) {
	defer resetHandler(mediumBlock)
	mediumBlock.WithJSONCoder(StdLib)
	for i := 0; i < b.N; i++ {
		bt, err := json.Marshal(&mediumBlock)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
		if bt == nil {
			b.Fatal("nil buf")
		}
	}
}

func BenchmarkBlock_Medium_CodecMarshal(b *testing.B) {
	defer resetHandler(mediumBlock)
	buf := make([]byte, 0, expectedSize)
	var h codec.Handle = new(codec.JsonHandle)
	enc := codec.NewEncoderBytes(&buf, h)
	mediumBlock.WithJSONCoder(GoCodec)
	for i := 0; i < b.N; i++ {
		err := enc.Encode(&mediumBlock)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
		buf = buf[:0]
	}
}

func BenchmarkBlock_Medium_EasyJSONMarshal(b *testing.B) {
	defer resetHandler(mediumBlock)
	mediumBlock.WithJSONCoder(EasyJson)
	for i := 0; i < b.N; i++ {
		bt, err := json.Marshal(&mediumBlock)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
		if bt == nil {
			b.Fatal("nil buf")
		}
	}
}

func BenchmarkBlock_Large_JSONMarshal(b *testing.B) {
	defer resetHandler(largeBlock)
	largeBlock.WithJSONCoder(StdLib)
	for i := 0; i < b.N; i++ {
		bt, err := json.Marshal(&largeBlock)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
		if bt == nil {
			b.Fatal("nil buf")
		}
	}
}

func BenchmarkBlock_Large_CodecMarshal(b *testing.B) {
	defer resetHandler(largeBlock)
	buf := make([]byte, 0, expectedSize)
	var h codec.Handle = new(codec.JsonHandle)
	enc := codec.NewEncoderBytes(&buf, h)
	largeBlock.WithJSONCoder(GoCodec)
	for i := 0; i < b.N; i++ {
		err := enc.Encode(&largeBlock)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
		buf = buf[:0]
	}
}

func BenchmarkBlock_Large_EasyJSONMarshal(b *testing.B) {
	defer resetHandler(largeBlock)
	largeBlock.WithJSONCoder(EasyJson)
	for i := 0; i < b.N; i++ {
		bt, err := json.Marshal(&largeBlock)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
		if bt == nil {
			b.Fatal("nil buf")
		}
	}
}
func BenchmarkBlock_Small_JSONUnmarshal(b *testing.B) {
	b.StopTimer()
	jsonBytes, err := json.Marshal(&smallBlock)
	if err != nil {
		b.Fatalf("failed to create test json %+v", err)
	}
	b.StartTimer()

	var temp Block
	temp.WithJSONCoder(StdLib)
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(jsonBytes, &temp)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
	}
}

func BenchmarkBlock_Small_CodecUnmarshal(b *testing.B) {
	b.StopTimer()
	smallBlock.WithJSONCoder(GoCodec)
	defer resetHandler(smallBlock)
	jsonBytes, err := json.Marshal(&smallBlock)
	if err != nil {
		b.Fatalf("failed to create test json %+v", err)
	}
	b.StartTimer()
	var h codec.Handle = new(codec.JsonHandle)
	var temp Block
	temp.WithJSONCoder(GoCodec)
	os.Setenv(EnvHack, "codec")
	defer os.Unsetenv(EnvHack)
	for i := 0; i < b.N; i++ {
		dec := codec.NewDecoderBytes(jsonBytes, h)
		err := dec.Decode(&temp)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
	}
}

func BenchmarkBlock_Small_EasyJSONUnmarshal(b *testing.B) {
	b.StopTimer()
	jsonBytes, err := json.Marshal(&smallBlock)
	smallBlock.WithJSONCoder(EasyJson)
	defer resetHandler(smallBlock)
	if err != nil {
		b.Fatalf("failed to create test json %+v", err)
	}
	b.StartTimer()

	var temp Block
	temp.WithJSONCoder(EasyJson)
	os.Setenv(EnvHack, "easy")
	defer os.Unsetenv(EnvHack)
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(jsonBytes, &temp)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
	}
}

func BenchmarkBlock_Medium_JSONUnmarshal(b *testing.B) {
	b.StopTimer()
	jsonBytes, err := json.Marshal(&mediumBlock)
	if err != nil {
		b.Fatalf("failed to create test json %+v", err)
	}
	b.StartTimer()

	var temp Block
	temp.WithJSONCoder(StdLib)
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(jsonBytes, &temp)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
	}
}

func BenchmarkBlock_Medium_CodecUnmarshal(b *testing.B) {
	b.StopTimer()
	mediumBlock.WithJSONCoder(GoCodec)
	defer resetHandler(mediumBlock)
	jsonBytes, err := json.Marshal(&mediumBlock)
	if err != nil {
		b.Fatalf("failed to create test json %+v", err)
	}
	b.StartTimer()
	var h codec.Handle = new(codec.JsonHandle)
	var temp Block
	temp.WithJSONCoder(GoCodec)
	os.Setenv(EnvHack, "codec")
	defer os.Unsetenv(EnvHack)
	for i := 0; i < b.N; i++ {
		dec := codec.NewDecoderBytes(jsonBytes, h)
		err := dec.Decode(&temp)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
	}
}

func BenchmarkBlock_Medium_EasyJSONUnmarshal(b *testing.B) {
	b.StopTimer()
	mediumBlock.WithJSONCoder(EasyJson)
	defer resetHandler(mediumBlock)
	jsonBytes, err := json.Marshal(&mediumBlock)
	if err != nil {
		b.Fatalf("failed to create test json %+v", err)
	}
	b.StartTimer()

	temp := new(Block)
	temp.WithJSONCoder(EasyJson)
	os.Setenv(EnvHack, "easy")
	defer os.Unsetenv(EnvHack)
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(jsonBytes, temp)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
	}
}

func BenchmarkBlock_Large_JSONUnmarshal(b *testing.B) {

	b.StopTimer()
	jsonBytes, err := json.Marshal(&largeBlock)
	if err != nil {
		b.Fatalf("failed to create test json %+v", err)
	}
	b.StartTimer()

	var temp Block
	temp.WithJSONCoder(StdLib)
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(jsonBytes, &temp)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
	}
}

func BenchmarkBlock_Large_CodecUnmarshal(b *testing.B) {
	b.StopTimer()
	largeBlock.WithJSONCoder(GoCodec)
	defer resetHandler(largeBlock)
	jsonBytes, err := json.Marshal(&largeBlock)
	if err != nil {
		b.Fatalf("failed to create test json %+v", err)
	}
	b.StartTimer()
	var h codec.Handle = new(codec.JsonHandle)
	var temp Block
	temp.WithJSONCoder(GoCodec)
	os.Setenv(EnvHack, "codec")
	defer os.Unsetenv(EnvHack)
	for i := 0; i < b.N; i++ {
		dec := codec.NewDecoderBytes(jsonBytes, h)
		err := dec.Decode(&temp)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
	}
}

func BenchmarkBlock_Large_EasyJSONUnmarshal(b *testing.B) {

	b.StopTimer()
	mediumBlock.WithJSONCoder(EasyJson)
	defer resetHandler(mediumBlock)
	jsonBytes, err := json.Marshal(&largeBlock)
	if err != nil {
		b.Fatalf("failed to create test json %+v", err)
	}
	b.StartTimer()

	var temp Block
	temp.WithJSONCoder(EasyJson)
	os.Setenv(EnvHack, "easy")
	defer os.Unsetenv(EnvHack)
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(jsonBytes, &temp)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
	}
}
func ptr[T any](t T) *T {
	return &t
}
