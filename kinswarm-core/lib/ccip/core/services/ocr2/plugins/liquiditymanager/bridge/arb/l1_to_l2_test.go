package arb

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	bridgetestutils "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/testutils"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclientmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/l2_arbitrum_gateway"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/mocks/mock_arbitrum_inbox"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	bridgecommon "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
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
		ctx context.Context //nolint:containedctx
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
				f.l1LogPoller.On("UnregisterFilter", mock.Anything, f.l1FilterName).Return(nil)
				f.l2LogPoller.On("UnregisterFilter", mock.Anything, f.l2FilterName).Return(nil)
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
				f.l1LogPoller.On("UnregisterFilter", mock.Anything, f.l1FilterName).Return(errors.New("unregister error"))
				f.l2LogPoller.On("UnregisterFilter", mock.Anything, f.l2FilterName).Return(nil)
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
				f.l1LogPoller.On("UnregisterFilter", mock.Anything, f.l1FilterName).Return(nil)
				f.l2LogPoller.On("UnregisterFilter", mock.Anything, f.l2FilterName).Return(errors.New("unregister error"))
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
		ctx context.Context //nolint:containedctx
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
		ctx context.Context //nolint:containedctx
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
		ctx        context.Context //nolint:containedctx
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
		readyCandidate *liquiditymanager.LiquidityManagerLiquidityTransferred
		receivedLogs   []*liquiditymanager.LiquidityManagerLiquidityTransferred
	}
	var (
		l2LiquidityManager = testutils.NewAddress()
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
				readyCandidate: &liquiditymanager.LiquidityManagerLiquidityTransferred{
					OcrSeqNum:          1,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l2LiquidityManager,
					Amount:             big.NewInt(100),
					BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
					BridgeReturnData:   mustPackReturnData(t, big.NewInt(10)),
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackReturnData(t, big.NewInt(9)),
						BridgeReturnData:   []byte{},
					},
					{
						OcrSeqNum:          3,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackReturnData(t, big.NewInt(11)),
						BridgeReturnData:   []byte{},
					},
					{
						OcrSeqNum:          4,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2LiquidityManager,
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
				readyCandidate: &liquiditymanager.LiquidityManagerLiquidityTransferred{
					OcrSeqNum:          1,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l2LiquidityManager,
					Amount:             big.NewInt(100),
					BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
					BridgeReturnData:   mustPackReturnData(t, big.NewInt(10)),
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackReturnData(t, big.NewInt(9)),
						BridgeReturnData:   []byte{},
					},
					{
						OcrSeqNum:          3,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackReturnData(t, big.NewInt(11)),
						BridgeReturnData:   []byte{},
					},
					{
						OcrSeqNum:          1,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackReturnData(t, big.NewInt(8)),
						BridgeReturnData:   []byte{},
					},
				},
			},
			false,
			false,
		},
		{
			"bad bridge return data",
			args{
				readyCandidate: &liquiditymanager.LiquidityManagerLiquidityTransferred{
					BridgeReturnData: []byte{1, 2, 3},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			},
			false,
			true,
		},
		{
			"bad bridge specific data",
			args{
				readyCandidate: &liquiditymanager.LiquidityManagerLiquidityTransferred{
					BridgeReturnData: mustPackReturnData(t, big.NewInt(10)),
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						BridgeSpecificData: []byte{1, 2, 3},
					},
				},
			},
			false,
			true,
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
		readyCandidates []*liquiditymanager.LiquidityManagerLiquidityTransferred
		receivedLogs    []*liquiditymanager.LiquidityManagerLiquidityTransferred
	}
	var (
		l2LiquidityManager = testutils.NewAddress()
	)
	tests := []struct {
		name      string
		args      args
		wantReady []*liquiditymanager.LiquidityManagerLiquidityTransferred
		wantErr   bool
	}{
		{
			"empty received list",
			args{
				readyCandidates: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						OcrSeqNum:          1,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(10)),
					},
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(11)),
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			},
			[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					OcrSeqNum:          1,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l2LiquidityManager,
					Amount:             big.NewInt(100),
					BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
					BridgeReturnData:   mustPackReturnData(t, big.NewInt(10)),
				},
				{
					OcrSeqNum:          2,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l2LiquidityManager,
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
				readyCandidates: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						OcrSeqNum:          1,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(10)),
					},
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(11)),
					},
					{
						OcrSeqNum:          3,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(12)),
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l2LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackReturnData(t, big.NewInt(10)),
						BridgeReturnData:   []byte{},
					},
				},
			},
			[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					OcrSeqNum:          2,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l2LiquidityManager,
					Amount:             big.NewInt(100),
					BridgeSpecificData: mustPackSendPayload(t, big.NewInt(100_000), big.NewInt(250_000), assets.GWei(3).ToInt()),
					BridgeReturnData:   mustPackReturnData(t, big.NewInt(11)),
				},
				{
					OcrSeqNum:          3,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l2LiquidityManager,
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
	var (
		localToken                = testutils.NewAddress()
		l1BridgeAdapterAddress    = testutils.NewAddress()
		l2LiquidityManagerAddress = testutils.NewAddress()
	)
	type args struct {
		localToken                models.Address
		l1BridgeAdapterAddress    common.Address
		l2LiquidityManagerAddress common.Address
		sentLogs                  []*liquiditymanager.LiquidityManagerLiquidityTransferred
		depositFinalizedLogs      []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized
		receivedLogs              []*liquiditymanager.LiquidityManagerLiquidityTransferred
	}
	tests := []struct {
		name          string
		args          args
		wantNotReady  []*liquiditymanager.LiquidityManagerLiquidityTransferred
		wantReady     []*liquiditymanager.LiquidityManagerLiquidityTransferred
		wantReadyData [][]byte
		wantErr       bool
	}{
		{
			name: "happy path - one ready, one not ready, one already received",
			args: args{
				localToken:                models.Address(localToken),
				l1BridgeAdapterAddress:    l1BridgeAdapterAddress,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
				sentLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					// Amount = 100, ready
					{
						To:               l2LiquidityManagerAddress,
						Amount:           big.NewInt(100),
						BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					},
					// Amount = 200, not ready
					{
						To:               l2LiquidityManagerAddress,
						Amount:           big.NewInt(200),
						BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
					},
					// Amount = 300, already received
					{
						To:               l2LiquidityManagerAddress,
						Amount:           big.NewInt(300),
						BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000003"),
					},
				},
				depositFinalizedLogs: []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
					// Amount = 100, ready
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      l2LiquidityManagerAddress,
						Amount:  big.NewInt(100),
					},
					// Amount = 300, already received
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      l2LiquidityManagerAddress,
						Amount:  big.NewInt(300),
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					// Amount = 300, already received
					{
						To:                 l2LiquidityManagerAddress,
						Amount:             big.NewInt(300),
						BridgeSpecificData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000003"),
					},
				},
			},
			wantNotReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
				// Amount = 200, not ready
				{
					To:               l2LiquidityManagerAddress,
					Amount:           big.NewInt(200),
					BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
				},
			},
			wantReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
				// Amount = 100, ready
				{
					To:               l2LiquidityManagerAddress,
					Amount:           big.NewInt(100),
					BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
				}},
			wantReadyData: [][]byte{bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001")},
			wantErr:       false,
		},
		{
			name: "mismatched token address",
			args: args{
				localToken:                models.Address(localToken),
				l1BridgeAdapterAddress:    l1BridgeAdapterAddress,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
				sentLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						To:               l2LiquidityManagerAddress,
						Amount:           big.NewInt(100),
						BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					},
				},
				depositFinalizedLogs: []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
					{
						L1Token: testutils.NewAddress(), // Mismatched address
						From:    l1BridgeAdapterAddress,
						To:      l2LiquidityManagerAddress,
						Amount:  big.NewInt(100),
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			},
			wantNotReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					To:               l2LiquidityManagerAddress,
					Amount:           big.NewInt(100),
					BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
				},
			},
			wantReady:     []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			wantReadyData: nil,
			wantErr:       false,
		},
		{
			name: "mismatched deposit finalized From address",
			args: args{
				localToken:                models.Address(localToken),
				l1BridgeAdapterAddress:    l1BridgeAdapterAddress,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
				sentLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						To:               l2LiquidityManagerAddress,
						Amount:           big.NewInt(100),
						BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					},
				},
				depositFinalizedLogs: []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
					{
						L1Token: localToken,
						From:    testutils.NewAddress(), // Mismatched address
						To:      l2LiquidityManagerAddress,
						Amount:  big.NewInt(100),
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			},
			wantNotReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					To:               l2LiquidityManagerAddress,
					Amount:           big.NewInt(100),
					BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
				},
			},
			wantReady:     []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			wantReadyData: nil,
			wantErr:       false,
		},
		{
			name: "mismatched deposit finalized To address",
			args: args{
				localToken:                models.Address(localToken),
				l1BridgeAdapterAddress:    l1BridgeAdapterAddress,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
				sentLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						To:               l2LiquidityManagerAddress,
						Amount:           big.NewInt(100),
						BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					},
				},
				depositFinalizedLogs: []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      testutils.NewAddress(), // Mismatched address
						Amount:  big.NewInt(100),
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			},
			wantNotReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					To:               l2LiquidityManagerAddress,
					Amount:           big.NewInt(100),
					BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
				},
			},
			wantReady:     []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			wantReadyData: nil,
			wantErr:       false,
		},
		{
			name: "mismatched deposit finalized amount",
			args: args{
				localToken:                models.Address(localToken),
				l1BridgeAdapterAddress:    l1BridgeAdapterAddress,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
				sentLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						To:               l2LiquidityManagerAddress,
						Amount:           big.NewInt(100),
						BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					},
				},
				depositFinalizedLogs: []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      l2LiquidityManagerAddress,
						Amount:  big.NewInt(200), // Mismatched amount
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			},
			wantNotReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					To:               l2LiquidityManagerAddress,
					Amount:           big.NewInt(100),
					BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
				},
			},
			wantReady:     []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			wantReadyData: nil,
			wantErr:       false,
		},
		{
			name: "amount matching for dep finalized event but mismatched for received log, should never happen",
			args: args{
				localToken:                models.Address(localToken),
				l1BridgeAdapterAddress:    l1BridgeAdapterAddress,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
				sentLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						To:               l2LiquidityManagerAddress,
						Amount:           big.NewInt(100),
						BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					},
				},
				depositFinalizedLogs: []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      l2LiquidityManagerAddress,
						Amount:  big.NewInt(100),
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						To:                 l2LiquidityManagerAddress,
						Amount:             big.NewInt(200), // Mismatched amount
						BridgeSpecificData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					},
				},
			},
			wantNotReady:  nil,
			wantReady:     nil,
			wantReadyData: nil,
			wantErr:       true,
		},
		{
			name: "mismatched bridge specific data for received event",
			args: args{
				localToken:                models.Address(localToken),
				l1BridgeAdapterAddress:    l1BridgeAdapterAddress,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
				sentLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						To:               l2LiquidityManagerAddress,
						Amount:           big.NewInt(100),
						BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					},
				},
				depositFinalizedLogs: []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      l2LiquidityManagerAddress,
						Amount:  big.NewInt(100),
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						To:     l2LiquidityManagerAddress,
						Amount: big.NewInt(100),
						// Mismatched bridge specific data
						BridgeSpecificData: bridgetestutils.MustPackBridgeData(t, "0x1111000000000000000000000000000000000000000000000000000000000001"),
					},
				},
			},
			wantNotReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			wantReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					To:               l2LiquidityManagerAddress,
					Amount:           big.NewInt(100),
					BridgeReturnData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
				},
			},
			wantReadyData: [][]byte{bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001")},
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNotReady, gotReady, gotReadyData, err := partitionTransfers(tt.args.localToken, tt.args.l1BridgeAdapterAddress, tt.args.l2LiquidityManagerAddress, tt.args.sentLogs, tt.args.depositFinalizedLogs, tt.args.receivedLogs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				bridgetestutils.AssertLiquidityTransferredEventSlicesEqual(t, tt.wantNotReady, gotNotReady, bridgetestutils.SortByBridgeReturnData)
				bridgetestutils.AssertLiquidityTransferredEventSlicesEqual(t, tt.wantReady, gotReady, bridgetestutils.SortByBridgeReturnData)
				require.Equal(t, tt.wantReadyData, gotReadyData)
			}
		})
	}
}

func Test_getEffectiveEvents(t *testing.T) {
	type args struct {
		localToken                models.Address
		l1BridgeAdapterAddress    common.Address
		l2LiquidityManagerAddress common.Address
		depositFinalizedLogs      []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized
	}
	var (
		localToken                = testutils.NewAddress()
		l1BridgeAdapterAddress    = testutils.NewAddress()
		l2LiquidityManagerAddress = testutils.NewAddress()
	)
	tests := []struct {
		name                          string
		args                          args
		wantEffectiveDepositFinalized []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized
	}{
		{
			"empty deposit finalized list",
			args{
				localToken:                models.Address(testutils.NewAddress()),
				l1BridgeAdapterAddress:    testutils.NewAddress(),
				l2LiquidityManagerAddress: testutils.NewAddress(),
				depositFinalizedLogs:      []*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{},
			},
			nil,
		},
		{
			"none applicable",
			args{
				localToken:                models.Address(testutils.NewAddress()),
				l1BridgeAdapterAddress:    testutils.NewAddress(),
				l2LiquidityManagerAddress: testutils.NewAddress(),
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
				localToken:                models.Address(localToken),
				l1BridgeAdapterAddress:    l1BridgeAdapterAddress,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
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
						To:      l2LiquidityManagerAddress,
						Amount:  big.NewInt(100),
					},
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      l2LiquidityManagerAddress,
						Amount:  big.NewInt(200),
					},
				},
			},
			[]*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
				{
					L1Token: localToken,
					From:    l1BridgeAdapterAddress,
					To:      l2LiquidityManagerAddress,
					Amount:  big.NewInt(100),
				},
				{
					L1Token: localToken,
					From:    l1BridgeAdapterAddress,
					To:      l2LiquidityManagerAddress,
					Amount:  big.NewInt(200),
				},
			},
		},
		{
			"some partially applicable but still not included",
			args{
				localToken:                models.Address(localToken),
				l1BridgeAdapterAddress:    l1BridgeAdapterAddress,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
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
						To:      l2LiquidityManagerAddress,
						Amount:  big.NewInt(100),
					},
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      testutils.NewAddress(), // not to liquidityManager
						Amount:  big.NewInt(200),
					},
				},
			},
			nil,
		},
		{
			"some fully applicable and some partially applicable but still not included",
			args{
				localToken:                models.Address(localToken),
				l1BridgeAdapterAddress:    l1BridgeAdapterAddress,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
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
						To:      l2LiquidityManagerAddress,
						Amount:  big.NewInt(100),
					},
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      testutils.NewAddress(), // not to liquidityManager
						Amount:  big.NewInt(200),
					},
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      l2LiquidityManagerAddress,
						Amount:  big.NewInt(100),
					},
					{
						L1Token: localToken,
						From:    l1BridgeAdapterAddress,
						To:      l2LiquidityManagerAddress,
						Amount:  big.NewInt(200),
					},
				},
			},
			[]*l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{
				{
					L1Token: localToken,
					From:    l1BridgeAdapterAddress,
					To:      l2LiquidityManagerAddress,
					Amount:  big.NewInt(100),
				},
				{
					L1Token: localToken,
					From:    l1BridgeAdapterAddress,
					To:      l2LiquidityManagerAddress,
					Amount:  big.NewInt(200),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEffectiveDepositFinalized := getEffectiveEvents(tt.args.localToken, tt.args.l1BridgeAdapterAddress, tt.args.l2LiquidityManagerAddress, tt.args.depositFinalizedLogs)
			require.Equal(t, tt.wantEffectiveDepositFinalized, gotEffectiveDepositFinalized)
		})
	}
}

func Test_l1ToL2Bridge_toPendingTransfers(t *testing.T) {
	type fields struct {
		localSelector             models.NetworkSelector
		remoteSelector            models.NetworkSelector
		l1LiquidityManager        liquiditymanager.LiquidityManagerInterface
		l2LiquidityManagerAddress common.Address
	}
	var (
		localSelector             = models.NetworkSelector(1)
		remoteSelector            = models.NetworkSelector(2)
		l2LiquidityManager        = testutils.NewAddress()
		localToken                = models.Address(testutils.NewAddress())
		remoteToken               = models.Address(testutils.NewAddress())
		l1LiquidityManagerAddress = testutils.NewAddress()
		l1LiquidityManager, err   = liquiditymanager.NewLiquidityManager(l1LiquidityManagerAddress, nil)
	)
	require.NoError(t, err)
	type args struct {
		localToken  models.Address
		remoteToken models.Address
		notReady    []*liquiditymanager.LiquidityManagerLiquidityTransferred
		ready       []*liquiditymanager.LiquidityManagerLiquidityTransferred
		readyData   [][]byte
		parsedToLP  map[bridgecommon.LogKey]logpoller.Log
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.PendingTransfer
		wantErr bool
	}{
		{
			"len(ready) not equal to len(readyData)",
			fields{
				localSelector:             localSelector,
				remoteSelector:            remoteSelector,
				l1LiquidityManager:        l1LiquidityManager,
				l2LiquidityManagerAddress: l2LiquidityManager,
			},
			args{
				localToken:  localToken,
				remoteToken: remoteToken,
				notReady:    nil,
				ready: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{},
					{},
				},
				readyData: [][]byte{
					{},
				},
				parsedToLP: nil,
			},
			nil,
			true,
		},
		{
			"not ready and ready",
			fields{
				localSelector:             localSelector,
				remoteSelector:            remoteSelector,
				l1LiquidityManager:        l1LiquidityManager,
				l2LiquidityManagerAddress: l2LiquidityManager,
			},
			args{
				localToken:  localToken,
				remoteToken: remoteToken,
				notReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						Amount: big.NewInt(100),
						Raw: types.Log{
							TxHash: common.HexToHash("0x1"),
							Index:  1,
						},
					},
				},
				ready: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						Amount: big.NewInt(300),
						Raw: types.Log{
							TxHash: common.HexToHash("0x2"),
							Index:  2,
						},
					},
				},
				readyData: [][]byte{
					{1, 2, 3},
				},
				parsedToLP: make(map[bridgecommon.LogKey]logpoller.Log),
			},
			[]models.PendingTransfer{
				{
					Transfer: models.Transfer{
						From:               localSelector,
						To:                 remoteSelector,
						Sender:             models.Address(l1LiquidityManagerAddress),
						Receiver:           models.Address(l2LiquidityManager),
						LocalTokenAddress:  localToken,
						RemoteTokenAddress: remoteToken,
						Amount:             ubig.NewI(100),
						Date:               time.Time{},
						BridgeData:         []byte{},
						Stage:              1,
						NativeBridgeFee:    ubig.NewI(0),
					},
					Status: models.TransferStatusNotReady,
					ID:     fmt.Sprintf("%s-%d", common.HexToHash("0x1"), 1),
				},
				{
					Transfer: models.Transfer{
						From:               localSelector,
						To:                 remoteSelector,
						Sender:             models.Address(l1LiquidityManagerAddress),
						Receiver:           models.Address(l2LiquidityManager),
						LocalTokenAddress:  localToken,
						RemoteTokenAddress: remoteToken,
						Amount:             ubig.NewI(300),
						Date:               time.Time{},
						BridgeData:         []byte{1, 2, 3},
						Stage:              2,
						NativeBridgeFee:    ubig.NewI(0),
					},
					Status: models.TransferStatusReady,
					ID:     fmt.Sprintf("%s-%d", common.HexToHash("0x2"), 2),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &l1ToL2Bridge{
				localSelector:             tt.fields.localSelector,
				remoteSelector:            tt.fields.remoteSelector,
				l1LiquidityManager:        tt.fields.l1LiquidityManager,
				l2LiquidityManagerAddress: tt.fields.l2LiquidityManagerAddress,
			}
			got, err := l.toPendingTransfers(tt.args.localToken, tt.args.remoteToken, tt.args.notReady, tt.args.ready, tt.args.readyData, tt.args.parsedToLP)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_l1ToL2Bridge_getLogs(t *testing.T) {
	type fields struct {
		localSelector             models.NetworkSelector
		remoteSelector            models.NetworkSelector
		l1LiquidityManager        liquiditymanager.LiquidityManagerInterface
		l2LiquidityManagerAddress common.Address
		l2Gateway                 l2_arbitrum_gateway.L2ArbitrumGatewayInterface
		l1LogPoller               *lpmocks.LogPoller
		l2LogPoller               *lpmocks.LogPoller
	}
	type args struct {
		ctx    context.Context //nolint:containedctx
		fromTs time.Time
	}
	var (
		localSelector             = models.NetworkSelector(1)
		remoteSelector            = models.NetworkSelector(2)
		l2LiquidityManagerAddress = testutils.NewAddress()
		l2GatewayAddress          = testutils.NewAddress()
		l1LiquidityManagerAddress = testutils.NewAddress()
	)
	l1LiquidityManager, err := liquiditymanager.NewLiquidityManager(l1LiquidityManagerAddress, nil)
	require.NoError(t, err)
	l2Gateway, err := l2_arbitrum_gateway.NewL2ArbitrumGateway(l2GatewayAddress, nil)
	require.NoError(t, err)
	tests := []struct {
		name                     string
		fields                   fields
		args                     args
		before                   func(t *testing.T, f fields, a args)
		assertions               func(t *testing.T, f fields)
		wantSendLogs             []logpoller.Log
		wantDepositFinalizedLogs []logpoller.Log
		wantReceiveLogs          []logpoller.Log
		wantErr                  bool
	}{
		{
			"error getting l1 liquidity transferred events",
			fields{
				localSelector:             localSelector,
				remoteSelector:            remoteSelector,
				l1LiquidityManager:        l1LiquidityManager,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
				l2Gateway:                 l2Gateway,
				l1LogPoller:               lpmocks.NewLogPoller(t),
				l2LogPoller:               lpmocks.NewLogPoller(t),
			},
			args{
				ctx:    testutils.Context(t),
				fromTs: time.Now(),
			},
			func(t *testing.T, f fields, a args) {
				f.l1LogPoller.On("IndexedLogsCreatedAfter",
					mock.Anything,
					bridgecommon.LiquidityTransferredTopic,
					l1LiquidityManager.Address(),
					bridgecommon.LiquidityTransferredToChainSelectorTopicIndex,
					[]common.Hash{bridgecommon.NetworkSelectorToHash(remoteSelector)},
					a.fromTs,
					evmtypes.Confirmations(1),
				).Return(nil, errors.New("error"))
			},
			func(t *testing.T, f fields) {
				f.l1LogPoller.AssertExpectations(t)
			},
			nil,
			nil,
			nil,
			true,
		},
		{
			"error getting l2 deposit finalized events",
			fields{
				localSelector:             localSelector,
				remoteSelector:            remoteSelector,
				l1LiquidityManager:        l1LiquidityManager,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
				l2Gateway:                 l2Gateway,
				l1LogPoller:               lpmocks.NewLogPoller(t),
				l2LogPoller:               lpmocks.NewLogPoller(t),
			},
			args{
				ctx:    testutils.Context(t),
				fromTs: time.Now(),
			},
			func(t *testing.T, f fields, a args) {
				f.l1LogPoller.On("IndexedLogsCreatedAfter",
					mock.Anything,
					bridgecommon.LiquidityTransferredTopic,
					l1LiquidityManager.Address(),
					bridgecommon.LiquidityTransferredToChainSelectorTopicIndex,
					[]common.Hash{bridgecommon.NetworkSelectorToHash(remoteSelector)},
					a.fromTs,
					evmtypes.Confirmations(1),
				).Return([]logpoller.Log{{}, {}}, nil)
				f.l2LogPoller.On("IndexedLogsCreatedAfter",
					mock.Anything,
					DepositFinalizedTopic,
					l2Gateway.Address(),
					DepositFinalizedToAddressTopicIndex,
					[]common.Hash{common.HexToHash(l2LiquidityManagerAddress.Hex())},
					a.fromTs,
					evmtypes.Finalized,
				).Return(nil, errors.New("error"))
			},
			func(t *testing.T, f fields) {
				f.l1LogPoller.AssertExpectations(t)
				f.l2LogPoller.AssertExpectations(t)
			},
			nil,
			nil,
			nil,
			true,
		},
		{
			"error getting l2 liquidity transferred events",
			fields{
				localSelector:             localSelector,
				remoteSelector:            remoteSelector,
				l1LiquidityManager:        l1LiquidityManager,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
				l2Gateway:                 l2Gateway,
				l1LogPoller:               lpmocks.NewLogPoller(t),
				l2LogPoller:               lpmocks.NewLogPoller(t),
			},
			args{
				ctx:    testutils.Context(t),
				fromTs: time.Now(),
			},
			func(t *testing.T, f fields, a args) {
				f.l1LogPoller.On("IndexedLogsCreatedAfter",
					mock.Anything,
					bridgecommon.LiquidityTransferredTopic,
					l1LiquidityManager.Address(),
					bridgecommon.LiquidityTransferredToChainSelectorTopicIndex,
					[]common.Hash{bridgecommon.NetworkSelectorToHash(remoteSelector)},
					a.fromTs,
					evmtypes.Confirmations(1),
				).Return([]logpoller.Log{{}, {}}, nil)
				f.l2LogPoller.On("IndexedLogsCreatedAfter",
					mock.Anything,
					DepositFinalizedTopic,
					l2Gateway.Address(),
					DepositFinalizedToAddressTopicIndex,
					[]common.Hash{common.HexToHash(l2LiquidityManagerAddress.Hex())},
					a.fromTs,
					evmtypes.Finalized,
				).Return([]logpoller.Log{{}, {}}, nil)
				f.l2LogPoller.On("IndexedLogsCreatedAfter",
					mock.Anything,
					bridgecommon.LiquidityTransferredTopic,
					l2LiquidityManagerAddress,
					bridgecommon.LiquidityTransferredFromChainSelectorTopicIndex,
					[]common.Hash{bridgecommon.NetworkSelectorToHash(localSelector)},
					a.fromTs,
					evmtypes.Confirmations(1),
				).Return(nil, errors.New("error"))
			},
			func(t *testing.T, f fields) {
				f.l1LogPoller.AssertExpectations(t)
				f.l2LogPoller.AssertExpectations(t)
			},
			nil,
			nil,
			nil,
			true,
		},
		{
			"happy path",
			fields{
				localSelector:             localSelector,
				remoteSelector:            remoteSelector,
				l1LiquidityManager:        l1LiquidityManager,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
				l2Gateway:                 l2Gateway,
				l1LogPoller:               lpmocks.NewLogPoller(t),
				l2LogPoller:               lpmocks.NewLogPoller(t),
			},
			args{
				ctx:    testutils.Context(t),
				fromTs: time.Now(),
			},
			func(t *testing.T, f fields, a args) {
				f.l1LogPoller.On("IndexedLogsCreatedAfter",
					mock.Anything,
					bridgecommon.LiquidityTransferredTopic,
					l1LiquidityManager.Address(),
					bridgecommon.LiquidityTransferredToChainSelectorTopicIndex,
					[]common.Hash{bridgecommon.NetworkSelectorToHash(remoteSelector)},
					a.fromTs,
					evmtypes.Confirmations(1),
				).Return([]logpoller.Log{
					{EventSig: bridgecommon.LiquidityTransferredTopic, TxHash: common.HexToHash("0x1")},
					{EventSig: bridgecommon.LiquidityTransferredTopic, TxHash: common.HexToHash("0x2")},
				}, nil)
				f.l2LogPoller.On("IndexedLogsCreatedAfter",
					mock.Anything,
					DepositFinalizedTopic,
					l2Gateway.Address(),
					DepositFinalizedToAddressTopicIndex,
					[]common.Hash{common.HexToHash(l2LiquidityManagerAddress.Hex())},
					a.fromTs,
					evmtypes.Finalized,
				).Return([]logpoller.Log{
					{EventSig: DepositFinalizedTopic, TxHash: common.HexToHash("0x3")},
					{EventSig: DepositFinalizedTopic, TxHash: common.HexToHash("0x4")},
				}, nil)
				f.l2LogPoller.On("IndexedLogsCreatedAfter",
					mock.Anything,
					bridgecommon.LiquidityTransferredTopic,
					l2LiquidityManagerAddress,
					bridgecommon.LiquidityTransferredFromChainSelectorTopicIndex,
					[]common.Hash{bridgecommon.NetworkSelectorToHash(localSelector)},
					a.fromTs,
					evmtypes.Confirmations(1),
				).Return([]logpoller.Log{
					{EventSig: bridgecommon.LiquidityTransferredTopic, TxHash: common.HexToHash("0x5")},
					{EventSig: bridgecommon.LiquidityTransferredTopic, TxHash: common.HexToHash("0x6")},
				}, nil)
			},
			func(t *testing.T, f fields) {
				f.l1LogPoller.AssertExpectations(t)
				f.l2LogPoller.AssertExpectations(t)
			},
			[]logpoller.Log{
				{EventSig: bridgecommon.LiquidityTransferredTopic, TxHash: common.HexToHash("0x1")},
				{EventSig: bridgecommon.LiquidityTransferredTopic, TxHash: common.HexToHash("0x2")},
			},
			[]logpoller.Log{
				{EventSig: DepositFinalizedTopic, TxHash: common.HexToHash("0x3")},
				{EventSig: DepositFinalizedTopic, TxHash: common.HexToHash("0x4")},
			},
			[]logpoller.Log{
				{EventSig: bridgecommon.LiquidityTransferredTopic, TxHash: common.HexToHash("0x5")},
				{EventSig: bridgecommon.LiquidityTransferredTopic, TxHash: common.HexToHash("0x6")},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &l1ToL2Bridge{
				localSelector:             tt.fields.localSelector,
				remoteSelector:            tt.fields.remoteSelector,
				l1LiquidityManager:        tt.fields.l1LiquidityManager,
				l2LiquidityManagerAddress: tt.fields.l2LiquidityManagerAddress,
				l2Gateway:                 tt.fields.l2Gateway,
				l1LogPoller:               tt.fields.l1LogPoller,
				l2LogPoller:               tt.fields.l2LogPoller,
			}
			tt.before(t, tt.fields, tt.args)
			defer tt.assertions(t, tt.fields)
			gotSendLogs, gotDepositFinalizedLogs, gotReceiveLogs, err := l.getLogs(tt.args.ctx, tt.args.fromTs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantSendLogs, gotSendLogs)
				require.Equal(t, tt.wantDepositFinalizedLogs, gotDepositFinalizedLogs)
				require.Equal(t, tt.wantReceiveLogs, gotReceiveLogs)
			}
		})
	}
}
