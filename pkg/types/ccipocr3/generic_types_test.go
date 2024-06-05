package ccipocr3

import (
	"math/big"

	"github.com/stretchr/testify/assert"
)

import (
	"encoding/json"
	"testing"
)

func TestSeqNumRange(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		rng := NewSeqNumRange(1, 2)
		assert.Equal(t, SeqNum(1), rng.Start())
		assert.Equal(t, SeqNum(2), rng.End())
	})

	t.Run("empty", func(t *testing.T) {
		rng := SeqNumRange{}
		assert.Equal(t, SeqNum(0), rng.Start())
		assert.Equal(t, SeqNum(0), rng.End())
	})

	t.Run("override start and end", func(t *testing.T) {
		rng := NewSeqNumRange(1, 2)
		rng.SetStart(10)
		rng.SetEnd(20)
		assert.Equal(t, SeqNum(10), rng.Start())
		assert.Equal(t, SeqNum(20), rng.End())
	})

	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "[1 -> 2]", NewSeqNumRange(1, 2).String())
		assert.Equal(t, "[0 -> 0]", SeqNumRange{}.String())
	})
}

func TestSeqNumRange_Overlap(t *testing.T) {
	testCases := []struct {
		name string
		r1   SeqNumRange
		r2   SeqNumRange
		exp  bool
	}{
		{"OverlapMiddle", SeqNumRange{5, 10}, SeqNumRange{8, 12}, true},
		{"OverlapStart", SeqNumRange{5, 10}, SeqNumRange{10, 15}, true},
		{"OverlapEnd", SeqNumRange{5, 10}, SeqNumRange{0, 5}, true},
		{"NoOverlapBefore", SeqNumRange{5, 10}, SeqNumRange{0, 4}, false},
		{"NoOverlapAfter", SeqNumRange{5, 10}, SeqNumRange{11, 15}, false},
		{"SameRange", SeqNumRange{5, 10}, SeqNumRange{5, 10}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.exp, tc.r1.Overlaps(tc.r2))
		})
	}
}

func TestSeqNumRange_Contains(t *testing.T) {
	tests := []struct {
		name     string
		r        SeqNumRange
		seq      SeqNum
		expected bool
	}{
		{"ContainsMiddle", SeqNumRange{5, 10}, SeqNum(7), true},
		{"ContainsStart", SeqNumRange{5, 10}, SeqNum(5), true},
		{"ContainsEnd", SeqNumRange{5, 10}, SeqNum(10), true},
		{"BeforeRange", SeqNumRange{5, 10}, SeqNum(4), false},
		{"AfterRange", SeqNumRange{5, 10}, SeqNum(11), false},
		{"EmptyRange", SeqNumRange{5, 5}, SeqNum(5), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.r.Contains(tt.seq))
		})
	}
}

func TestCCIPMsg_String(t *testing.T) {
	tests := []struct {
		name     string
		c        CCIPMsg
		expected string
	}{
		{
			"base",
			CCIPMsg{CCIPMsgBaseDetails{ID: [32]byte{123}, SourceChain: ChainSelector(1), SeqNum: 2}},
			`{"id":"0x7b00000000000000000000000000000000000000000000000000000000000000","sourceChain":"1","seqNum":"2"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.c.String())
		})
	}
}

func TestNewTokenPrice(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		tp := NewTokenPrice("link", big.NewInt(1000))
		assert.Equal(t, "link", string(tp.TokenID))
		assert.Equal(t, uint64(1000), tp.Price.Int.Uint64())
	})
}

func TestNewGasPriceChain(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		gpc := NewGasPriceChain(big.NewInt(1000), ChainSelector(1))
		assert.Equal(t, uint64(1000), (gpc.GasPrice).Uint64())
		assert.Equal(t, ChainSelector(1), gpc.ChainSel)
	})
}

func TestMerkleRoot(t *testing.T) {
	t.Run("str", func(t *testing.T) {
		mr := Bytes32([32]byte{1})
		assert.Equal(t, "0x0100000000000000000000000000000000000000000000000000000000000000", mr.String())
	})

	t.Run("json", func(t *testing.T) {
		mr := Bytes32([32]byte{1})
		b, err := json.Marshal(mr)
		assert.NoError(t, err)
		assert.Equal(t, `"0x0100000000000000000000000000000000000000000000000000000000000000"`, string(b))

		mr2 := Bytes32{}
		err = json.Unmarshal(b, &mr2)
		assert.NoError(t, err)
		assert.Equal(t, mr, mr2)

		mr3 := Bytes32{}
		err = json.Unmarshal([]byte(`"123"`), &mr3)
		assert.Error(t, err)

		err = json.Unmarshal([]byte(`""`), &mr3)
		assert.Error(t, err)
	})
}
