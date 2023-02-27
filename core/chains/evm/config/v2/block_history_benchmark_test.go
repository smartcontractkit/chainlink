package v2

import (
	"encoding/json"
	"testing"
	"unsafe"

	"github.com/ugorji/go/codec"
)

var static = BlockHistoryEstimator{
	BatchSize:                 ptr(uint32(4096)),
	BlockHistorySize:          ptr(uint16(1024)),
	CheckInclusionBlocks:      ptr(uint16(8)),
	CheckInclusionPercentile:  ptr(uint16(50)),
	EIP1559FeeCapBufferBlocks: ptr(uint16(10)),
	TransactionPercentile:     ptr(uint16(75)),
}

var staticWithNil = BlockHistoryEstimator{
	BatchSize:                 ptr(uint32(4096)),
	BlockHistorySize:          ptr(uint16(1024)),
	CheckInclusionBlocks:      ptr(uint16(8)),
	CheckInclusionPercentile:  ptr(uint16(50)),
	EIP1559FeeCapBufferBlocks: nil,
	TransactionPercentile:     ptr(uint16(75)),
}
var expectedSize = unsafe.Sizeof(BlockHistoryEstimator{})

var staticBytes []byte

func init() {

	var err error
	staticBytes, err = json.Marshal(&static)
	if err != nil {
		panic(err)
	}
}
func BenchmarkBlockHistoryEstimator_Marshal_JSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bt, err := json.Marshal(&static)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
		if bt == nil {
			b.Fatal("nil buf")
		}
	}
}

func BenchmarkBlockHistoryEstimator_Unarshal_JSON(b *testing.B) {
	var temp BlockHistoryEstimator
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(staticBytes, &temp)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
	}
}

func BenchmarkBlockHistoryEstimator_withNil_Marshal_Marshal_JSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bt, err := json.Marshal(&staticWithNil)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
		if bt == nil {
			b.Fatal("nil buf")
		}
	}
}

func BenchmarkBlockHistoryEstimator_Marshal_Codec(b *testing.B) {
	buf := make([]byte, 0, expectedSize)
	var h codec.Handle = new(codec.JsonHandle)
	enc := codec.NewEncoderBytes(&buf, h)
	for i := 0; i < b.N; i++ {
		err := enc.Encode(&static)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
		buf = buf[:0]
	}
}

func BenchmarkBlockHistoryEstimator_Unmarshal_Codec(b *testing.B) {
	var h codec.Handle = new(codec.JsonHandle)
	var temp BlockHistoryEstimator
	for i := 0; i < b.N; i++ {
		dec := codec.NewDecoderBytes(staticBytes, h)
		err := dec.Decode(&temp)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
	}
}

func BenchmarkBlockHistoryEstimator_withNil_Marshal_Codec(b *testing.B) {
	buf := make([]byte, 0, expectedSize)
	var h codec.Handle = new(codec.JsonHandle)
	enc := codec.NewEncoderBytes(&buf, h)
	for i := 0; i < b.N; i++ {
		err := enc.Encode(&staticWithNil)
		if err != nil {
			b.Fatalf("err %+v", err)
		}
		buf = buf[:0]
	}
}

func ptr[T any](t T) *T {
	return &t
}
