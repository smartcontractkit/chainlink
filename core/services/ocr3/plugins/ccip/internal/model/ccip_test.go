package model

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
)

func TestSeqNumRange_String(t *testing.T) {
	tests := []struct {
		name     string
		s        SeqNumRange
		expected string
	}{
		{"base", SeqNumRange{1, 2}, "[1 -> 2]"},
		{"empty", SeqNumRange{}, "[0 -> 0]"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.s.String())
		})
	}
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
			CCIPMsg{CCIPMsgBaseDetails{SourceChain: ChainSelector(1), SeqNum: 2}},
			`{"sourceChain":"1","seqNum":"2"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.c.String())
		})
	}
}

func TestSeqNumRange_StartEnd(t *testing.T) {
	s := SeqNumRange{1, 2}
	assert.Equal(t, SeqNum(1), s.Start())
	assert.Equal(t, SeqNum(2), s.End())

	s = SeqNumRange{}
	assert.Equal(t, SeqNum(0), s.Start())
	assert.Equal(t, SeqNum(0), s.End())
}
