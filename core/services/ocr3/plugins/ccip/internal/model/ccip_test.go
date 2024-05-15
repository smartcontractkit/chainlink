package model

import (
	"encoding/json"
	"math/big"
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
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

func TestChainSelector_String(t *testing.T) {
	tests := []struct {
		name     string
		c        ChainSelector
		expected string
	}{
		{"unknown chain", ChainSelector(1), "ChainSelector(1)"},
		{"known chain", ChainSelector(chainsel.ETHEREUM_MAINNET.Selector), "5009297550715157269 (ethereum-mainnet)"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.c.String())
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
		assert.Equal(t, uint64(1000), tp.Price.Uint64())
	})
}

func TestNewGasPriceChain(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		gpc := NewGasPriceChain(big.NewInt(1000), ChainSelector(1))
		assert.Equal(t, uint64(1000), (*big.Int)(gpc.GasPrice).Uint64())
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
