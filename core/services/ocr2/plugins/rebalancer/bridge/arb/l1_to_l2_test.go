package arb

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/stretchr/testify/mock"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclientmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/mocks/mock_arbitrum_inbox"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func Test_l1ToL2Bridge_QuorumizedBridgePayload(t *testing.T) {
	type args struct {
		payloads [][]byte
		f        int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"not enough payloads",
			args{
				[][]byte{},
				1,
			},
			nil,
			true,
		},
		{
			"happy path",
			args{
				[][]byte{
					mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
					mustPackSendPayload(t, big.NewInt(105_000), big.NewInt(225_000), assets.GWei(2).ToInt()),
					mustPackSendPayload(t, big.NewInt(103_300), big.NewInt(255_000), assets.GWei(5).ToInt()),
				},
				1,
			},
			mustPackSendPayload(t, big.NewInt(103_300), big.NewInt(250_000), assets.GWei(3).ToInt()), // second highest in each category
			false,
		},
		{
			"happy path, less payloads",
			args{
				[][]byte{
					mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
					mustPackSendPayload(t, big.NewInt(103_300), big.NewInt(255_000), assets.GWei(5).ToInt()),
				},
				1,
			},
			mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()), // second highest in each category
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &l1ToL2Bridge{}
			got, err := l.QuorumizedBridgePayload(tt.args.payloads, tt.args.f)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				unpackedGot, err := UnpackL1ToL2SendBridgePayload(got)
				require.NoError(t, err)
				unpackedExpected, err := UnpackL1ToL2SendBridgePayload(tt.want)
				require.NoError(t, err)
				require.Equal(t, unpackedExpected.GasLimit, unpackedGot.GasLimit)
				require.Equal(t, unpackedExpected.MaxSubmissionCost, unpackedGot.MaxSubmissionCost)
				require.Equal(t, unpackedExpected.MaxFeePerGas, unpackedGot.MaxFeePerGas)
				// require.Equal(t, tt.want, got)
			}
		})
	}
}

func mustPackSendPayload(t *testing.T, gasLimit, maxSubmissionCost, maxFeePerGas *big.Int) []byte {
	packed, err := PackL1ToL2SendBridgePayload(gasLimit, maxSubmissionCost, maxFeePerGas)
	require.NoError(t, err)
	return packed
}

func Test_l1ToL2Bridge_Close(t *testing.T) {
	type fields struct {
		l1LogPoller  *lpmocks.LogPoller
		l2LogPoller  *lpmocks.LogPoller
		l1FilterName string
		l2FilterName string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		before     func(*testing.T, fields)
		assertions func(*testing.T, fields)
	}{
		{
			"happy path",
			fields{
				l1LogPoller:  lpmocks.NewLogPoller(t),
				l2LogPoller:  lpmocks.NewLogPoller(t),
				l1FilterName: "l1FilterName",
				l2FilterName: "l2FilterName",
			},
			args{
				testutils.Context(t),
			},
			false,
			func(t *testing.T, f fields) {
				f.l1LogPoller.On("UnregisterFilter", f.l1FilterName).Return(nil)
				f.l2LogPoller.On("UnregisterFilter", f.l2FilterName).Return(nil)
			},
			func(t *testing.T, f fields) {
				f.l1LogPoller.AssertExpectations(t)
				f.l2LogPoller.AssertExpectations(t)
			},
		},
		{
			"l1 unregister error",
			fields{
				l1LogPoller:  lpmocks.NewLogPoller(t),
				l2LogPoller:  lpmocks.NewLogPoller(t),
				l1FilterName: "l1FilterName",
				l2FilterName: "l2FilterName",
			},
			args{
				testutils.Context(t),
			},
			true,
			func(t *testing.T, f fields) {
				f.l1LogPoller.On("UnregisterFilter", f.l1FilterName).Return(errors.New("unregister error"))
				f.l2LogPoller.On("UnregisterFilter", f.l2FilterName).Return(nil)
			},
			func(t *testing.T, f fields) {
				f.l1LogPoller.AssertExpectations(t)
				f.l2LogPoller.AssertExpectations(t)
			},
		},
		{
			"l2 unregister error",
			fields{
				l1LogPoller:  lpmocks.NewLogPoller(t),
				l2LogPoller:  lpmocks.NewLogPoller(t),
				l1FilterName: "l1FilterName",
				l2FilterName: "l2FilterName",
			},
			args{
				testutils.Context(t),
			},
			true,
			func(t *testing.T, f fields) {
				f.l1LogPoller.On("UnregisterFilter", f.l1FilterName).Return(nil)
				f.l2LogPoller.On("UnregisterFilter", f.l2FilterName).Return(errors.New("unregister error"))
			},
			func(t *testing.T, f fields) {
				f.l1LogPoller.AssertExpectations(t)
				f.l2LogPoller.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &l1ToL2Bridge{
				l1LogPoller:  tt.fields.l1LogPoller,
				l2LogPoller:  tt.fields.l2LogPoller,
				l1FilterName: tt.fields.l1FilterName,
				l2FilterName: tt.fields.l2FilterName,
			}
			if tt.before != nil {
				tt.before(t, tt.fields)
				defer tt.assertions(t, tt.fields)
			}

			err := l.Close(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_l1ToL2Bridge_estimateMaxFeePerGasOnL2(t *testing.T) {
	type fields struct {
		l2Client *evmclientmocks.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       *big.Int
		wantErr    bool
		before     func(*testing.T, fields)
		assertions func(*testing.T, fields)
	}{
		{
			"happy path",
			fields{
				l2Client: evmclientmocks.NewClient(t),
			},
			args{
				testutils.Context(t),
			},
			assets.GWei(3).ToInt(),
			false,
			func(t *testing.T, f fields) {
				f.l2Client.On("SuggestGasPrice", mock.Anything).Return(assets.GWei(1).ToInt(), nil)
			},
			func(t *testing.T, f fields) {
				f.l2Client.AssertExpectations(t)
			},
		},
		{
			"suggest gas price error",
			fields{
				l2Client: evmclientmocks.NewClient(t),
			},
			args{
				testutils.Context(t),
			},
			nil,
			true,
			func(t *testing.T, f fields) {
				f.l2Client.On("SuggestGasPrice", mock.Anything).Return(nil, errors.New("error"))
			},
			func(t *testing.T, f fields) {
				f.l2Client.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &l1ToL2Bridge{
				l2Client: tt.fields.l2Client,
			}
			if tt.before != nil {
				tt.before(t, tt.fields)
				defer tt.assertions(t, tt.fields)
			}

			got, err := l.estimateMaxFeePerGasOnL2(tt.args.ctx)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_l1ToL2Bridge_estimateRetryableGasLimit(t *testing.T) {
	type fields struct {
		l2Client *evmclientmocks.Client
	}
	type args struct {
		ctx context.Context
		rd  RetryableData
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       *big.Int
		wantErr    bool
		before     func(*testing.T, fields)
		assertions func(*testing.T, fields)
	}{
		{
			"happy path",
			fields{
				l2Client: evmclientmocks.NewClient(t),
			},
			args{
				testutils.Context(t),
				RetryableData{
					From:                testutils.NewAddress(),
					To:                  testutils.NewAddress(),
					L2CallValue:         big.NewInt(2),
					ExcessFeeRefundAddr: testutils.NewAddress(),
					CallValueRefundAddr: testutils.NewAddress(),
					Data:                []byte{1, 2, 3},
				},
			},
			big.NewInt(125_000),
			false,
			func(t *testing.T, f fields) {
				f.l2Client.On("EstimateGas", mock.Anything, mock.MatchedBy(func(e ethereum.CallMsg) bool {
					return e.To != nil && *e.To == NodeInterfaceAddress
				})).Return(uint64(125_000), nil)
			},
			func(t *testing.T, f fields) {
				f.l2Client.AssertExpectations(t)
			},
		},
		{
			"estimate gas error",
			fields{
				l2Client: evmclientmocks.NewClient(t),
			},
			args{
				testutils.Context(t),
				RetryableData{
					From:                testutils.NewAddress(),
					To:                  testutils.NewAddress(),
					L2CallValue:         big.NewInt(2),
					ExcessFeeRefundAddr: testutils.NewAddress(),
					CallValueRefundAddr: testutils.NewAddress(),
					Data:                []byte{1, 2, 3},
				},
			},
			nil,
			true,
			func(t *testing.T, f fields) {
				f.l2Client.On("EstimateGas", mock.Anything, mock.MatchedBy(func(e ethereum.CallMsg) bool {
					return e.To != nil && *e.To == NodeInterfaceAddress
				})).Return(uint64(0), errors.New("error"))
			},
			func(t *testing.T, f fields) {
				f.l2Client.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &l1ToL2Bridge{
				l2Client: tt.fields.l2Client,
			}
			if tt.before != nil {
				tt.before(t, tt.fields)
				defer tt.assertions(t, tt.fields)
			}

			got, err := l.estimateRetryableGasLimit(tt.args.ctx, tt.args.rd)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_l1ToL2Bridge_estimateMaxSubmissionFee(t *testing.T) {
	type fields struct {
		l1Inbox *mock_arbitrum_inbox.ArbitrumInboxInterface
	}
	type args struct {
		ctx        context.Context
		l1BaseFee  *big.Int
		dataLength int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       *big.Int
		wantErr    bool
		before     func(*testing.T, fields, args)
		assertions func(*testing.T, fields)
	}{
		{
			"happy path",
			fields{
				l1Inbox: mock_arbitrum_inbox.NewArbitrumInboxInterface(t),
			},
			args{
				testutils.Context(t),
				assets.GWei(1).ToInt(),
				100,
			},
			big.NewInt(400_000),
			false,
			func(t *testing.T, f fields, a args) {
				f.l1Inbox.On(
					"CalculateRetryableSubmissionFee",
					mock.Anything,
					big.NewInt(int64(a.dataLength)),
					a.l1BaseFee,
				).
					Return(big.NewInt(100_000), nil)
			},
			func(t *testing.T, f fields) {
				f.l1Inbox.AssertExpectations(t)
			},
		},
		{
			"calculate retryable submission fee error",
			fields{
				l1Inbox: mock_arbitrum_inbox.NewArbitrumInboxInterface(t),
			},
			args{
				testutils.Context(t),
				assets.GWei(1).ToInt(),
				100,
			},
			nil,
			true,
			func(t *testing.T, f fields, a args) {
				f.l1Inbox.On(
					"CalculateRetryableSubmissionFee",
					mock.Anything,
					big.NewInt(int64(a.dataLength)),
					a.l1BaseFee,
				).
					Return(nil, errors.New("error"))
			},
			func(t *testing.T, f fields) {
				f.l1Inbox.AssertExpectations(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &l1ToL2Bridge{
				l1Inbox: tt.fields.l1Inbox,
			}
			if tt.before != nil {
				tt.before(t, tt.fields, tt.args)
				defer tt.assertions(t, tt.fields)
			}

			got, err := l.estimateMaxSubmissionFee(tt.args.ctx, tt.args.l1BaseFee, tt.args.dataLength)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
