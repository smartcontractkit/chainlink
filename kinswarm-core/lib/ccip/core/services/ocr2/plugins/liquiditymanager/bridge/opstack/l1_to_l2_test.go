package opstack

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_standard_bridge"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/abiutils"
	bridgetestutils "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/testutils"
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
			"non-matching nonces/payloads",
			args{
				[][]byte{
					bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
					bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000003"),
				},
				1,
			},
			nil,
			true,
		},
		{
			"happy path",
			args{
				[][]byte{
					bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
					bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
					bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
				},
				1,
			},
			bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
			false,
		},
		{
			"happy path, fewer payloads",
			args{
				[][]byte{
					bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
					bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
				},
				1,
			},
			bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
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
				unpackedGot, err := abiutils.UnpackUint256(got)
				require.NoError(t, err)
				unpackedExpected, err := abiutils.UnpackUint256(tt.want)
				require.NoError(t, err)
				require.Equal(t, unpackedExpected, unpackedGot)
			}
		})
	}
}

func Test_partitionTransfers(t *testing.T) {
	l1LocalToken := models.Address(common.HexToAddress("0x123"))
	l2LocalToken := models.Address(common.HexToAddress("0x456"))
	l1BridgeAdapterAddress := common.HexToAddress("0x789")
	l2LiquidityManagerAddress := common.HexToAddress("0xabc")
	l1ChainSelector := chainsel.ETHEREUM_MAINNET.Selector
	l2ChainSelector := chainsel.ETHEREUM_MAINNET_OPTIMISM_1.Selector

	type args struct {
		localToken                models.Address
		l1BridgeAdapterAddress    common.Address
		l2LiquidityManagerAddress common.Address
		sentLogs                  []*liquiditymanager.LiquidityManagerLiquidityTransferred
		erc20BridgeFinalizedLogs  []*optimism_standard_bridge.OptimismStandardBridgeERC20BridgeFinalized
		receivedLogs              []*liquiditymanager.LiquidityManagerLiquidityTransferred
	}
	tests := []struct {
		name            string
		args            args
		wantNotReady    []*liquiditymanager.LiquidityManagerLiquidityTransferred
		wantReady       []*liquiditymanager.LiquidityManagerLiquidityTransferred
		wantMissingSent []*optimism_standard_bridge.OptimismStandardBridgeERC20BridgeFinalized
		wantErr         bool
	}{
		{
			name: "empty",
			args: args{
				localToken:                models.Address{},
				l1BridgeAdapterAddress:    common.Address{},
				l2LiquidityManagerAddress: common.Address{},
				sentLogs:                  nil,
				erc20BridgeFinalizedLogs:  nil,
				receivedLogs:              nil,
			},
			wantNotReady:    nil,
			wantReady:       nil,
			wantMissingSent: nil,
			wantErr:         false,
		},
		{
			name: "happy path - one ready, one not ready, one missing, one done",
			args: args{
				localToken:                l1LocalToken,
				l1BridgeAdapterAddress:    l1BridgeAdapterAddress,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
				sentLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					// This one is not ready (only present in sentLogs)
					{
						FromChainSelector:  l1ChainSelector,
						ToChainSelector:    l2ChainSelector,
						To:                 l2LiquidityManagerAddress,
						Amount:             big.NewInt(1),
						BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
						BridgeSpecificData: []byte{},
					},
					// This one is ready (present in sentLogs and finalized logs)
					{
						FromChainSelector:  l1ChainSelector,
						ToChainSelector:    l2ChainSelector,
						To:                 l2LiquidityManagerAddress,
						Amount:             big.NewInt(1),
						BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
						BridgeSpecificData: []byte{},
					},
					// This one is done/already received (present in sentLogs, finalized logs, and receivedLogs), should not be included in any output slices
					{
						FromChainSelector:  l1ChainSelector,
						ToChainSelector:    l2ChainSelector,
						To:                 l2LiquidityManagerAddress,
						Amount:             big.NewInt(1),
						BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000003"),
						BridgeSpecificData: []byte{},
					},
				},
				erc20BridgeFinalizedLogs: []*optimism_standard_bridge.OptimismStandardBridgeERC20BridgeFinalized{
					// This one is ready (present in sentLogs and finalized logs)
					{
						LocalToken:  common.Address(l2LocalToken),
						RemoteToken: common.Address(l1LocalToken),
						From:        l1BridgeAdapterAddress,
						To:          l2LiquidityManagerAddress,
						Amount:      big.NewInt(1),
						ExtraData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
					},
					// This one is already done (present in sentLogs, finalized logs, and receivedLogs)
					{
						LocalToken:  common.Address(l2LocalToken),
						RemoteToken: common.Address(l1LocalToken),
						From:        l1BridgeAdapterAddress,
						To:          l2LiquidityManagerAddress,
						Amount:      big.NewInt(1),
						ExtraData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000003"),
					},
					// This one is missing (not present in sentLogs)
					{
						LocalToken:  common.Address(l2LocalToken),
						RemoteToken: common.Address(l1LocalToken),
						From:        l1BridgeAdapterAddress,
						To:          l2LiquidityManagerAddress,
						Amount:      big.NewInt(1),
						ExtraData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000065"),
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						FromChainSelector:  l1ChainSelector,
						ToChainSelector:    l2ChainSelector,
						To:                 l2LiquidityManagerAddress,
						Amount:             big.NewInt(1),
						BridgeReturnData:   []byte{},
						BridgeSpecificData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000003"),
					},
				},
			},
			wantNotReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
				// This one is not ready (only present in sentLogs)
				{
					FromChainSelector:  l1ChainSelector,
					ToChainSelector:    l2ChainSelector,
					To:                 l2LiquidityManagerAddress,
					Amount:             big.NewInt(1),
					BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					BridgeSpecificData: []byte{},
				},
			},
			wantReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
				// This one is ready (present in sentLogs and finalized logs)
				{
					FromChainSelector:  l1ChainSelector,
					ToChainSelector:    l2ChainSelector,
					To:                 l2LiquidityManagerAddress,
					Amount:             big.NewInt(1),
					BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
					BridgeSpecificData: []byte{},
				},
			},
			wantMissingSent: []*optimism_standard_bridge.OptimismStandardBridgeERC20BridgeFinalized{
				// This one is missing (not present in sentLogs)
				{
					LocalToken:  common.Address(l2LocalToken),
					RemoteToken: common.Address(l1LocalToken),
					From:        l1BridgeAdapterAddress,
					To:          l2LiquidityManagerAddress,
					Amount:      big.NewInt(1),
					ExtraData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000065"),
				},
			},
			wantErr: false,
		},
		{
			name: "L2 standard bridge finalized event - mismatched from, to, remote token fields",
			args: args{
				localToken:                l1LocalToken,
				l1BridgeAdapterAddress:    l1BridgeAdapterAddress,
				l2LiquidityManagerAddress: l2LiquidityManagerAddress,
				sentLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					// Mismatched finalized event 'from' field
					{
						FromChainSelector:  l1ChainSelector,
						ToChainSelector:    l2ChainSelector,
						To:                 l2LiquidityManagerAddress,
						Amount:             big.NewInt(1),
						BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
						BridgeSpecificData: []byte{},
					},
					// Mismatched finalized event 'to' field
					{
						FromChainSelector:  l1ChainSelector,
						ToChainSelector:    l2ChainSelector,
						To:                 l2LiquidityManagerAddress,
						Amount:             big.NewInt(1),
						BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
						BridgeSpecificData: []byte{},
					},
					// Mismatched finalized event 'remote_token' field
					{
						FromChainSelector:  l1ChainSelector,
						ToChainSelector:    l2ChainSelector,
						To:                 l2LiquidityManagerAddress,
						Amount:             big.NewInt(1),
						BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000003"),
						BridgeSpecificData: []byte{},
					},
				},
				erc20BridgeFinalizedLogs: []*optimism_standard_bridge.OptimismStandardBridgeERC20BridgeFinalized{
					// Mismatched finalized event 'from' field
					{
						LocalToken:  common.Address(l2LocalToken),
						RemoteToken: common.Address(l1LocalToken),
						From:        common.HexToAddress("0x123"),
						To:          l2LiquidityManagerAddress,
						Amount:      big.NewInt(1),
						ExtraData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					},
					// Mismatched finalized event 'to' field
					{
						LocalToken:  common.Address(l2LocalToken),
						RemoteToken: common.Address(l1LocalToken),
						From:        l1BridgeAdapterAddress,
						To:          common.HexToAddress("0x456"),
						Amount:      big.NewInt(1),
						ExtraData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
					},
					// Mismatched finalized event 'remote_token' field
					{
						LocalToken:  common.Address(l2LocalToken),
						RemoteToken: common.HexToAddress("0xabcd"),
						From:        l1BridgeAdapterAddress,
						To:          l2LiquidityManagerAddress,
						Amount:      big.NewInt(1),
						ExtraData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000003"),
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			},
			wantNotReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
				// Mismatched finalized event 'from' field
				{
					FromChainSelector:  l1ChainSelector,
					ToChainSelector:    l2ChainSelector,
					To:                 l2LiquidityManagerAddress,
					Amount:             big.NewInt(1),
					BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					BridgeSpecificData: []byte{},
				},
				// Mismatched finalized event 'to' field
				{
					FromChainSelector:  l1ChainSelector,
					ToChainSelector:    l2ChainSelector,
					To:                 l2LiquidityManagerAddress,
					Amount:             big.NewInt(1),
					BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
					BridgeSpecificData: []byte{},
				},
				// Mismatched finalized event 'remote_token' field
				{
					FromChainSelector:  l1ChainSelector,
					ToChainSelector:    l2ChainSelector,
					To:                 l2LiquidityManagerAddress,
					Amount:             big.NewInt(1),
					BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000003"),
					BridgeSpecificData: []byte{},
				}},
			wantReady:       nil,
			wantMissingSent: nil,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNotReady, gotReady, gotMissingSent, err := partitionTransfers(tt.args.localToken, tt.args.l1BridgeAdapterAddress, tt.args.l2LiquidityManagerAddress, tt.args.sentLogs, tt.args.erc20BridgeFinalizedLogs, tt.args.receivedLogs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				bridgetestutils.AssertLiquidityTransferredEventSlicesEqual(t, tt.wantNotReady, gotNotReady, bridgetestutils.SortByBridgeReturnData)
				bridgetestutils.AssertLiquidityTransferredEventSlicesEqual(t, tt.wantReady, gotReady, bridgetestutils.SortByBridgeReturnData)
				assert.Equal(t, tt.wantMissingSent, gotMissingSent)
			}
		})
	}
}

func Test_l1ToL2Bridge_Close(t *testing.T) {
	type fields struct {
		l1LogPoller  *lpmocks.LogPoller
		l2LogPoller  *lpmocks.LogPoller
		l1FilterName string
		l2FilterName string
	}
	tests := []struct {
		name       string
		fields     fields
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

			err := l.Close(testutils.Context(t))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
