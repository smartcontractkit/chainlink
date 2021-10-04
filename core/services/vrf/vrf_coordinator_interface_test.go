package vrf_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	keyHash   = secretKey.PublicKey.MustHash()
	jobID     = common.BytesToHash([]byte("1234567890abcdef1234567890abcdef"))
	seed      = big.NewInt(1)
	sender    = common.HexToAddress("0xecfcab0a285d3380e488a39b4bb21e777f8a4eac")
	fee       = assets.NewLinkFromJuels(100)
	requestID = common.HexToHash("0xcafe")
	raw       = vrf.RawRandomnessRequestLog{
		KeyHash:   keyHash,
		Seed:      seed,
		JobID:     jobID,
		Sender:    sender,
		Fee:       (*big.Int)(fee),
		RequestID: requestID,
		Raw: types.Log{
			// A raw, on-the-wire RandomnessRequestLog is the concat of fields as uint256's
			Data: append(append(append(append(
				keyHash.Bytes(),
				common.BigToHash(seed).Bytes()...),
				sender.Hash().Bytes()...),
				fee.ToHash().Bytes()...),
				requestID.Bytes()...),
			Topics: []common.Hash{{}, jobID},
		},
	}
)

func TestVRFParseRandomnessRequestLog(t *testing.T) {
	r := vrf.RawRandomnessRequestLogToRandomnessRequestLog(&raw)
	rawLog, err := r.RawData()
	require.NoError(t, err)
	assert.Equal(t, rawLog, raw.Raw.Data)
	nR, err := vrf.ParseRandomnessRequestLog(types.Log{
		Data:   rawLog,
		Topics: []common.Hash{vrf.VRFRandomnessRequestLogTopic(), jobID},
	})
	require.NoError(t, err)
	require.True(t, r.Equal(*nR),
		"Round-tripping RandomnessRequestLog through serialization and parsing "+
			"resulted in a different log.")
}
