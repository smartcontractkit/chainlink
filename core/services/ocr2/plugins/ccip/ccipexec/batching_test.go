package ccipexec

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

func Test_validateSendRequests(t *testing.T) {
	testCases := []struct {
		name             string
		seqNums          []uint64
		providedInterval ccipdata.CommitStoreInterval
		expErr           bool
	}{
		{
			name:             "zero interval no seq nums",
			seqNums:          nil,
			providedInterval: ccipdata.CommitStoreInterval{Min: 0, Max: 0},
			expErr:           true,
		},
		{
			name:             "exp 1 seq num got none",
			seqNums:          nil,
			providedInterval: ccipdata.CommitStoreInterval{Min: 1, Max: 1},
			expErr:           true,
		},
		{
			name:             "exp 10 seq num got none",
			seqNums:          nil,
			providedInterval: ccipdata.CommitStoreInterval{Min: 1, Max: 10},
			expErr:           true,
		},
		{
			name:             "got 1 seq num as expected",
			seqNums:          []uint64{1},
			providedInterval: ccipdata.CommitStoreInterval{Min: 1, Max: 1},
			expErr:           false,
		},
		{
			name:             "got 5 seq num as expected",
			seqNums:          []uint64{11, 12, 13, 14, 15},
			providedInterval: ccipdata.CommitStoreInterval{Min: 11, Max: 15},
			expErr:           false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sendReqs := make([]ccipdata.Event[internal.EVM2EVMMessage], 0, len(tc.seqNums))
			for _, seqNum := range tc.seqNums {
				sendReqs = append(sendReqs, ccipdata.Event[internal.EVM2EVMMessage]{
					Data: internal.EVM2EVMMessage{SequenceNumber: seqNum},
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
