package arb

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/mock"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/arbitrum_rollup_core"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/mocks/mock_arbitrum_rollup_core"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

func Test_L2ToL1Bridge_GetBridgePayloadAndFee(t *testing.T) {
	bridge := &l2ToL1Bridge{}
	payload, fee, err := bridge.GetBridgePayloadAndFee(testutils.Context(t), models.Transfer{})
	require.NoError(t, err)
	require.Empty(t, payload)
	require.Equal(t, big.NewInt(0), fee)
}

func Test_l2ToL1Bridge_getLatestNodeConfirmed(t *testing.T) {
	type fields struct {
		l1LogPoller *lpmocks.LogPoller
		rollupCore  *mock_arbitrum_rollup_core.ArbRollupCoreInterface
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       *arbitrum_rollup_core.ArbRollupCoreNodeConfirmed
		wantErr    bool
		before     func(*testing.T, fields, *arbitrum_rollup_core.ArbRollupCoreNodeConfirmed)
		assertions func(*testing.T, fields)
	}{
		{
			"log found",
			fields{
				l1LogPoller: lpmocks.NewLogPoller(t),
				rollupCore:  mock_arbitrum_rollup_core.NewArbRollupCoreInterface(t),
			},
			args{
				ctx: testutils.Context(t),
			},
			&arbitrum_rollup_core.ArbRollupCoreNodeConfirmed{
				NodeNum:   1,
				BlockHash: testutils.Random32Byte(),
				SendRoot:  testutils.Random32Byte(),
			},
			false,
			func(t *testing.T, f fields, want *arbitrum_rollup_core.ArbRollupCoreNodeConfirmed) {
				topic1 := common.HexToHash(hexutil.EncodeUint64(want.NodeNum))
				data, err := utils.ABIEncode(`[{"type": "bytes32"}, {"type": "bytes32"}]`, want.BlockHash, want.SendRoot)
				require.NoError(t, err)
				rollupAddress := testutils.NewAddress()
				f.l1LogPoller.On("LatestLogByEventSigWithConfs", NodeConfirmedTopic, rollupAddress, logpoller.Finalized, mock.Anything).
					Return(&logpoller.Log{
						Topics: [][]byte{
							NodeConfirmedTopic[:],
							topic1[:],
						},
						Data: data,
					}, nil)
				f.rollupCore.On("Address").Return(rollupAddress)
				f.rollupCore.On("ParseNodeConfirmed", mock.Anything).Return(want, nil)
			},
			func(t *testing.T, f fields) {
				f.l1LogPoller.AssertExpectations(t)
				f.rollupCore.AssertExpectations(t)
			},
		},
		{
			"log not found",
			fields{
				l1LogPoller: lpmocks.NewLogPoller(t),
				rollupCore:  mock_arbitrum_rollup_core.NewArbRollupCoreInterface(t),
			},
			args{
				ctx: testutils.Context(t),
			},
			nil,
			true,
			func(t *testing.T, f fields, want *arbitrum_rollup_core.ArbRollupCoreNodeConfirmed) {
				rollupAddress := testutils.NewAddress()
				f.l1LogPoller.On("LatestLogByEventSigWithConfs", NodeConfirmedTopic, rollupAddress, logpoller.Finalized, mock.Anything).
					Return(nil, errors.New("not found"))
				f.rollupCore.On("Address").Return(rollupAddress)
			},
			func(t *testing.T, f fields) {
				f.l1LogPoller.AssertExpectations(t)
				f.rollupCore.AssertExpectations(t)
			},
		},
		{
			"parse error",
			fields{
				l1LogPoller: lpmocks.NewLogPoller(t),
				rollupCore:  mock_arbitrum_rollup_core.NewArbRollupCoreInterface(t),
			},
			args{
				ctx: testutils.Context(t),
			},
			&arbitrum_rollup_core.ArbRollupCoreNodeConfirmed{
				NodeNum:   1,
				BlockHash: testutils.Random32Byte(),
				SendRoot:  testutils.Random32Byte(),
			},
			true,
			func(t *testing.T, f fields, want *arbitrum_rollup_core.ArbRollupCoreNodeConfirmed) {
				topic1 := common.HexToHash(hexutil.EncodeUint64(want.NodeNum))
				data, err := utils.ABIEncode(`[{"type": "bytes32"}, {"type": "bytes32"}]`, want.BlockHash, want.SendRoot)
				require.NoError(t, err)
				rollupAddress := testutils.NewAddress()
				f.l1LogPoller.On("LatestLogByEventSigWithConfs", NodeConfirmedTopic, rollupAddress, logpoller.Finalized, mock.Anything).
					Return(&logpoller.Log{
						Topics: [][]byte{
							NodeConfirmedTopic[:],
							topic1[:],
						},
						Data: data,
					}, nil)
				f.rollupCore.On("Address").Return(rollupAddress)
				f.rollupCore.On("ParseNodeConfirmed", mock.Anything).Return(nil, errors.New("parse error"))
			},
			func(t *testing.T, f fields) {
				f.l1LogPoller.AssertExpectations(t)
				f.rollupCore.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &l2ToL1Bridge{
				l1LogPoller: tt.fields.l1LogPoller,
				rollupCore:  tt.fields.rollupCore,
			}
			if tt.before != nil {
				tt.before(t, tt.fields, tt.want)
				defer tt.assertions(t, tt.fields)
			}
			got, err := l.getLatestNodeConfirmed(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
