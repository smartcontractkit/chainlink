package ccipexec

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
)

func Test_validateSendRequests(t *testing.T) {
	testCases := []struct {
		name             string
		seqNums          []uint64
		providedInterval cciptypes.CommitStoreInterval
		expErr           bool
	}{
		{
			name:             "zero interval no seq nums",
			seqNums:          nil,
			providedInterval: cciptypes.CommitStoreInterval{Min: 0, Max: 0},
			expErr:           true,
		},
		{
			name:             "exp 1 seq num got none",
			seqNums:          nil,
			providedInterval: cciptypes.CommitStoreInterval{Min: 1, Max: 1},
			expErr:           true,
		},
		{
			name:             "exp 10 seq num got none",
			seqNums:          nil,
			providedInterval: cciptypes.CommitStoreInterval{Min: 1, Max: 10},
			expErr:           true,
		},
		{
			name:             "got 1 seq num as expected",
			seqNums:          []uint64{1},
			providedInterval: cciptypes.CommitStoreInterval{Min: 1, Max: 1},
			expErr:           false,
		},
		{
			name:             "got 5 seq num as expected",
			seqNums:          []uint64{11, 12, 13, 14, 15},
			providedInterval: cciptypes.CommitStoreInterval{Min: 11, Max: 15},
			expErr:           false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sendReqs := make([]cciptypes.EVM2EVMMessageWithTxMeta, 0, len(tc.seqNums))
			for _, seqNum := range tc.seqNums {
				sendReqs = append(sendReqs, cciptypes.EVM2EVMMessageWithTxMeta{
					EVM2EVMMessage: cciptypes.EVM2EVMMessage{SequenceNumber: seqNum},
				})
			}
			err := validateSendRequests(sendReqs, tc.providedInterval)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
