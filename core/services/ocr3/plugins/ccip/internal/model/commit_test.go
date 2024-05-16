package model

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitPluginObservation_EncodeAndDecode(t *testing.T) {
	obs := NewCommitPluginObservation(
		[]CCIPMsgBaseDetails{
			{ID: [32]byte{1}, SourceChain: math.MaxUint64, SeqNum: 123},
			{ID: [32]byte{2}, SourceChain: 321, SeqNum: math.MaxUint64},
		},
		[]GasPriceChain{},
		[]TokenPrice{},
		[]SeqNumChain{},
	)

	b, err := obs.Encode()
	assert.NoError(t, err)
	assert.Equal(t, `{"newMsgs":[{"id":"0x0100000000000000000000000000000000000000000000000000000000000000","sourceChain":"18446744073709551615","seqNum":"123"},{"id":"0x0200000000000000000000000000000000000000000000000000000000000000","sourceChain":"321","seqNum":"18446744073709551615"}],"gasPrices":[],"tokenPrices":[],"maxSeqNums":[]}`, string(b))

	obs2, err := DecodeCommitPluginObservation(b)
	assert.NoError(t, err)
	assert.Equal(t, obs, obs2)
}

func TestCommitPluginOutcome_EncodeAndDecode(t *testing.T) {
	o := NewCommitPluginOutcome(
		[]SeqNumChain{
			NewSeqNumChain(ChainSelector(1), SeqNum(20)),
			NewSeqNumChain(ChainSelector(2), SeqNum(25)),
		},
		[]MerkleRootChain{
			NewMerkleRootChain(ChainSelector(1), NewSeqNumRange(21, 22), [32]byte{1}),
			NewMerkleRootChain(ChainSelector(2), NewSeqNumRange(25, 35), [32]byte{2}),
		},
		[]TokenPrice{
			NewTokenPrice("0x123", big.NewInt(1234)),
			NewTokenPrice("0x125", big.NewInt(0).Mul(big.NewInt(999999999999), big.NewInt(999999999999))),
		},
	)

	b, err := o.Encode()
	assert.NoError(t, err)
	assert.Equal(t, `{"maxSeqNums":[{"chainSel":1,"seqNum":20},{"chainSel":2,"seqNum":25}],"merkleRoots":[{"chain":1,"seqNumsRange":[21,22],"merkleRoot":"0x0100000000000000000000000000000000000000000000000000000000000000"},{"chain":2,"seqNumsRange":[25,35],"merkleRoot":"0x0200000000000000000000000000000000000000000000000000000000000000"}],"tokenPrices":[{"tokenID":"0x123","price":"1234"},{"tokenID":"0x125","price":"999999999998000000000001"}]}`, string(b))

	o2, err := DecodeCommitPluginOutcome(b)
	assert.NoError(t, err)
	assert.Equal(t, o, o2)

	assert.Equal(t, `{MaxSeqNums: [{ChainSelector(1) 20} {ChainSelector(2) 25}], MerkleRoots: [{ChainSelector(1) [21 -> 22] 0x0100000000000000000000000000000000000000000000000000000000000000} {ChainSelector(2) [25 -> 35] 0x0200000000000000000000000000000000000000000000000000000000000000}]}`, o.String())
}

func TestCommitPluginReport(t *testing.T) {
	t.Run("is empty", func(t *testing.T) {
		r := NewCommitPluginReport(nil, nil)
		assert.True(t, r.IsEmpty())
	})

	t.Run("is not empty", func(t *testing.T) {
		r := NewCommitPluginReport(make([]MerkleRootChain, 1), nil)
		assert.False(t, r.IsEmpty())

		r = NewCommitPluginReport(nil, make([]TokenPrice, 1))
		assert.False(t, r.IsEmpty())

		r = NewCommitPluginReport(make([]MerkleRootChain, 1), make([]TokenPrice, 1))
		assert.False(t, r.IsEmpty())
	})
}
