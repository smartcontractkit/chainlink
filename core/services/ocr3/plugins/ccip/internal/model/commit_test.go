package model

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitPluginObservation_EncodeAndDecode(t *testing.T) {
	obs := NewCommitPluginObservation(
		"nodeID",
		[]CCIPMsgBaseDetails{
			{SourceChain: math.MaxUint64, SeqNum: 123},
			{SourceChain: 321, SeqNum: math.MaxUint64},
		},
	)

	b, err := obs.Encode()
	assert.NoError(t, err)
	assert.Equal(t, `{"nodeID":"nodeID","newMsgs":[{"sourceChain":"18446744073709551615","seqNum":"123"},{"sourceChain":"321","seqNum":"18446744073709551615"}]}`, string(b))

	obs2, err := DecodeCommitPluginObservation(b)
	assert.NoError(t, err)
	assert.Equal(t, obs, obs2)
}
