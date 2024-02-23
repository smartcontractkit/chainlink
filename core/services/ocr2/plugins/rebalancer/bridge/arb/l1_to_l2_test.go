package arb

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclientmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/l2_arbitrum_gateway"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/rebalancer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/mocks/mock_arbitrum_inbox"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
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

func Test_matchingExecutionExists(t *testing.T) {
	type args struct {
		readyCandidate *rebalancer.RebalancerLiquidityTransferred
		receivedLogs   []*rebalancer.RebalancerLiquidityTransferred
	}
	var (
		l2Rebalancer = testutils.NewAddress()
	)
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"matching execution exists",
			args{
				readyCandidate: &rebalancer.RebalancerLiquidityTransferred{
					OcrSeqNum:          1,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l2Rebalancer,
					Amount:             big.NewInt(100),
					BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
					BridgeReturnData:   mustPackReturnData(t, big.NewInt(10)),
				},
				receivedLogs: []*rebalancer.RebalancerLiquidityTransferred{
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2Rebalancer,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackReturnData(t, big.NewInt(9)),
						BridgeReturnData:   []byte{},
					},
					{
						OcrSeqNum:          3,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2Rebalancer,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackReturnData(t, big.NewInt(11)),
						BridgeReturnData:   []byte{},
					},
					{
						OcrSeqNum:          4,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2Rebalancer,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackReturnData(t, big.NewInt(10)),
						BridgeReturnData:   []byte{},
					},
				},
			},
			true,
			false,
		},
		{
			"no matching execution exists",
			args{
				readyCandidate: &rebalancer.RebalancerLiquidityTransferred{
					OcrSeqNum:          1,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l2Rebalancer,
					Amount:             big.NewInt(100),
					BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
					BridgeReturnData:   mustPackReturnData(t, big.NewInt(10)),
				},
				receivedLogs: []*rebalancer.RebalancerLiquidityTransferred{
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2Rebalancer,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackReturnData(t, big.NewInt(9)),
						BridgeReturnData:   []byte{},
					},
					{
						OcrSeqNum:          3,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2Rebalancer,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackReturnData(t, big.NewInt(11)),
						BridgeReturnData:   []byte{},
					},
					{
						OcrSeqNum:          1,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2Rebalancer,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackReturnData(t, big.NewInt(8)),
						BridgeReturnData:   []byte{},
					},
				},
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchingExecutionExists(tt.args.readyCandidate, tt.args.receivedLogs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func mustPackReturnData(t *testing.T, l1ToL2Id *big.Int) []byte {
	packed, err := utils.ABIEncode(`[{"type": "uint256"}]`, l1ToL2Id)
	require.NoError(t, err)
	return packed
}

func Test_filterExecuted(t *testing.T) {
	type args struct {
		readyCandidates []*rebalancer.RebalancerLiquidityTransferred
		receivedLogs    []*rebalancer.RebalancerLiquidityTransferred
	}
	var (
		l2Rebalancer = testutils.NewAddress()
	)
	tests := []struct {
		name      string
		args      args
		wantReady []*rebalancer.RebalancerLiquidityTransferred
		wantErr   bool
	}{
		{
			"empty received list",
			args{
				readyCandidates: []*rebalancer.RebalancerLiquidityTransferred{
					{
						OcrSeqNum:          1,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2Rebalancer,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(10)),
					},
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2Rebalancer,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(11)),
					},
				},
				receivedLogs: []*rebalancer.RebalancerLiquidityTransferred{},
			},
			[]*rebalancer.RebalancerLiquidityTransferred{
				{
					OcrSeqNum:          1,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l2Rebalancer,
					Amount:             big.NewInt(100),
					BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
					BridgeReturnData:   mustPackReturnData(t, big.NewInt(10)),
				},
				{
					OcrSeqNum:          2,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l2Rebalancer,
					Amount:             big.NewInt(100),
					BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
					BridgeReturnData:   mustPackReturnData(t, big.NewInt(11)),
				},
			},
			false,
		},
		{
			"non-empty received list, some executed",
			args{
				readyCandidates: []*rebalancer.RebalancerLiquidityTransferred{
					{
						OcrSeqNum:          1,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2Rebalancer,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(10)),
					},
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2Rebalancer,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(11)),
					},
					{
						OcrSeqNum:          3,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2Rebalancer,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(12)),
					},
				},
				receivedLogs: []*rebalancer.RebalancerLiquidityTransferred{
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2Rebalancer,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackReturnData(t, big.NewInt(10)),
						BridgeReturnData:   []byte{},
					},
				},
			},
			[]*rebalancer.RebalancerLiquidityTransferred{
				{
					OcrSeqNum:          2,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l2Rebalancer,
					Amount:             big.NewInt(100),
					BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
					BridgeReturnData:   mustPackReturnData(t, big.NewInt(11)),
				},
				{
					OcrSeqNum:          3,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l2Rebalancer,
					Amount:             big.NewInt(100),
					BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
					BridgeReturnData:   mustPackReturnData(t, big.NewInt(12)),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReady, err := filterExecuted(tt.args.readyCandidates, tt.args.receivedLogs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantReady, gotReady)
			}
		})
	}
}

func Test_partitionTransfers(t *testing.T) {
	type args struct {
		localToken             models.Address
		l1BridgeAdapterAddress common.Address
		l2RebalancerAddress    common.Address
		sentLogs               []*rebalancer.RebalancerLiquidityTransferred
		depositFinalizedLogs   []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized
		receivedLogs           []*rebalancer.RebalancerLiquidityTransferred
	}
	tests := []struct {
		name          string
		args          args
		wantNotReady  []*rebalancer.RebalancerLiquidityTransferred
		wantReady     []*rebalancer.RebalancerLiquidityTransferred
		wantReadyData [][]byte
		wantErr       bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNotReady, gotReady, gotReadyData, err := partitionTransfers(tt.args.localToken, tt.args.l1BridgeAdapterAddress, tt.args.l2RebalancerAddress, tt.args.sentLogs, tt.args.depositFinalizedLogs, tt.args.receivedLogs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantNotReady, gotNotReady)
				require.Equal(t, tt.wantReady, gotReady)
				require.Equal(t, tt.wantReadyData, gotReadyData)
			}
		})
	}
}

func Test_getEffectiveEvents(t *testing.T) {
	type args struct {
		localToken             models.Address
		l1BridgeAdapterAddress common.Address
		l2RebalancerAddress    common.Address
		depositFinalizedLogs   []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized
	}
	var (
		localToken             = testutils.NewAddress()
		l1BridgeAdapterAddress = testutils.NewAddress()
		l2RebalancerAddress    = testutils.NewAddress()
	)
	tests := []struct {
		name                          string
		args                          args
		wantEffectiveDepositFinalized []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized
	}{
		{
			"empty deposit finalized list",
			args{
				localToken:             models.Address(testutils.NewAddress()),
				l1BridgeAdapterAddress: testutils.NewAddress(),
				l2RebalancerAddress:    testutils.NewAddress(),
				depositFinalizedLogs:   []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{},
			},
			nil,
		},
		{
			"none applicable",
			args{
				localToken:             models.Address(testutils.NewAddress()),
				l1BridgeAdapterAddress: testutils.NewAddress(),
				l2RebalancerAddress:    testutils.NewAddress(),
				depositFinalizedLogs: []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
					{
						L1Token: testutils.NewAddress(),
						From:    testutils.NewAddress(),
						To:      testutils.NewAddress(),
						Amount:  big.NewInt(100),
					},
					{
						L1Token: testutils.NewAddress(),
						From:    testutils.NewAddress(),
						To:      testutils.NewAddress(),
						Amount:  big.NewInt(100),
					},
				},
			},
			nil,
		},
		{
			"some exactly applicable",
			args{
				localToken:             models.Address(localToken),
				l1BridgeAdapterAddress: l1BridgeAdapterAddress,
				l2RebalancerAddress:    l2RebalancerAddress,
				depositFinalizedLogs: []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
					{
						L1Token: testutils.NewAddress(),
						From:    testutils.NewAddress(),
						To:      testutils.NewAddress(),
						Amount:  big.NewInt(100),
					},
					{
						L1Token: testutils.NewAddress(),
						From:    testutils.NewAddress(),
						To:      testutils.NewAddress(),
						Amount:  big.NewInt(100),
					},
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      l2RebalancerAddress,
						Amount:  big.NewInt(100),
					},
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      l2RebalancerAddress,
						Amount:  big.NewInt(200),
					},
				},
			},
			[]*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
				{
					L1Token: localToken,
					From:    l1BridgeAdapterAddress,
					To:      l2RebalancerAddress,
					Amount:  big.NewInt(100),
				},
				{
					L1Token: localToken,
					From:    l1BridgeAdapterAddress,
					To:      l2RebalancerAddress,
					Amount:  big.NewInt(200),
				},
			},
		},
		{
			"some partially applicable but still not included",
			args{
				localToken:             models.Address(localToken),
				l1BridgeAdapterAddress: l1BridgeAdapterAddress,
				l2RebalancerAddress:    l2RebalancerAddress,
				depositFinalizedLogs: []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
					{
						L1Token: testutils.NewAddress(),
						From:    testutils.NewAddress(),
						To:      testutils.NewAddress(),
						Amount:  big.NewInt(100),
					},
					{
						L1Token: testutils.NewAddress(),
						From:    testutils.NewAddress(),
						To:      testutils.NewAddress(),
						Amount:  big.NewInt(100),
					},
					{
						L1Token: localToken,
						From:    testutils.NewAddress(), // not from bridge adapter
						To:      l2RebalancerAddress,
						Amount:  big.NewInt(100),
					},
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      testutils.NewAddress(), // not to rebalancer
						Amount:  big.NewInt(200),
					},
				},
			},
			nil,
		},
		{
			"some fully applicable and some partially applicable but still not included",
			args{
				localToken:             models.Address(localToken),
				l1BridgeAdapterAddress: l1BridgeAdapterAddress,
				l2RebalancerAddress:    l2RebalancerAddress,
				depositFinalizedLogs: []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
					{
						L1Token: testutils.NewAddress(),
						From:    testutils.NewAddress(),
						To:      testutils.NewAddress(),
						Amount:  big.NewInt(100),
					},
					{
						L1Token: testutils.NewAddress(),
						From:    testutils.NewAddress(),
						To:      testutils.NewAddress(),
						Amount:  big.NewInt(100),
					},
					{
						L1Token: localToken,
						From:    testutils.NewAddress(), // not from bridge adapter
						To:      l2RebalancerAddress,
						Amount:  big.NewInt(100),
					},
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      testutils.NewAddress(), // not to rebalancer
						Amount:  big.NewInt(200),
					},
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      l2RebalancerAddress,
						Amount:  big.NewInt(100),
					},
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      l2RebalancerAddress,
						Amount:  big.NewInt(200),
					},
				},
			},
			[]*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
				{
					L1Token: localToken,
					From:    l1BridgeAdapterAddress,
					To:      l2RebalancerAddress,
					Amount:  big.NewInt(100),
				},
				{
					L1Token: localToken,
					From:    l1BridgeAdapterAddress,
					To:      l2RebalancerAddress,
					Amount:  big.NewInt(200),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEffectiveDepositFinalized := getEffectiveEvents(tt.args.localToken, tt.args.l1BridgeAdapterAddress, tt.args.l2RebalancerAddress, tt.args.depositFinalizedLogs)
			require.Equal(t, tt.wantEffectiveDepositFinalized, gotEffectiveDepositFinalized)
		})
	}
}
