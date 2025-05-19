package v1_2_0

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/common/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

func TestLogPollerClient_GetSendRequestsBetweenSeqNumsV1_2_0(t *testing.T) {
	onRampAddr := utils.RandomAddress()
	seqNum := uint64(100)
	limit := uint64(10)
	lggr := logger.TestLogger(t)

	tests := []struct {
		name          string
		finalized     bool
		confirmations evmtypes.Confirmations
	}{
		{"finalized", true, evmtypes.Finalized},
		{"unfinalized", false, evmtypes.Confirmations(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lp := mocks.NewLogPoller(t)
			onRampV2, err := NewOnRamp(lggr, 1, 1, onRampAddr, lp, nil)
			require.NoError(t, err)

			lp.On("LogsDataWordRange",
				mock.Anything,
				onRampV2.sendRequestedEventSig,
				onRampAddr,
				onRampV2.sendRequestedSeqNumberWord,
				abihelpers.EvmWord(seqNum),
				abihelpers.EvmWord(seqNum+limit),
				tt.confirmations,
			).Once().Return([]logpoller.Log{}, nil)

			events, err1 := onRampV2.GetSendRequestsBetweenSeqNums(context.Background(), seqNum, seqNum+limit, tt.finalized)
			assert.NoError(t, err1)
			assert.Empty(t, events)

			lp.AssertExpectations(t)
		})
	}
}
