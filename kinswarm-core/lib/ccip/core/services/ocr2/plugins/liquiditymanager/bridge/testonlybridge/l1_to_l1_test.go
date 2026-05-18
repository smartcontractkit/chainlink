package testonlybridge

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmclient_mocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lp_mocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/mock_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func Test_testBridge_toPendingTransfers(t *testing.T) {
	var (
		sourceSelector              = models.NetworkSelector(1)
		destSelector                = models.NetworkSelector(2)
		localToken                  = models.Address(testutils.NewAddress())
		remoteToken                 = models.Address(testutils.NewAddress())
		sourceAdapterAddress        = models.Address(testutils.NewAddress())
		destLiquidityManagerAddress = models.Address(testutils.NewAddress())
	)
	sourceAdapter, err := mock_l1_bridge_adapter.NewMockL1BridgeAdapter(common.Address(sourceAdapterAddress), nil)
	require.NoError(t, err)
	destLiquidityManager, err := liquiditymanager.NewLiquidityManager(common.Address(destLiquidityManagerAddress), nil)
	require.NoError(t, err)
	type fields struct {
		sourceSelector       models.NetworkSelector
		destSelector         models.NetworkSelector
		destLiquidityManager liquiditymanager.LiquidityManagerInterface
		sourceAdapter        mock_l1_bridge_adapter.MockL1BridgeAdapterInterface
		lggr                 logger.Logger
	}
	type args struct {
		localToken      models.Address
		remoteToken     models.Address
		readyToProve    []*liquiditymanager.LiquidityManagerLiquidityTransferred
		readyToFinalize []*liquiditymanager.LiquidityManagerLiquidityTransferred
		parsedToLP      map[logKey]logpoller.Log
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []models.PendingTransfer
	}{
		{
			"no pending transfers",
			fields{
				sourceSelector:       sourceSelector,
				destSelector:         destSelector,
				destLiquidityManager: destLiquidityManager,
				sourceAdapter:        sourceAdapter,
				lggr:                 logger.TestLogger(t),
			},
			args{
				localToken:      localToken,
				remoteToken:     remoteToken,
				readyToProve:    []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
				readyToFinalize: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
				parsedToLP:      map[logKey]logpoller.Log{},
			},
			nil,
		},
		{
			"some pending transfers, all ready to prove",
			fields{
				sourceSelector:       sourceSelector,
				destSelector:         destSelector,
				destLiquidityManager: destLiquidityManager,
				sourceAdapter:        sourceAdapter,
				lggr:                 logger.TestLogger(t),
			},
			args{
				localToken:  localToken,
				remoteToken: remoteToken,
				readyToProve: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						Amount:           big.NewInt(1),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(1)),
						Raw: types.Log{
							TxHash: common.HexToHash("0x1"),
							Index:  1,
						},
					},
					{
						Amount:           big.NewInt(2),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(2)),
						Raw: types.Log{
							TxHash: common.HexToHash("0x2"),
							Index:  2,
						},
					},
				},
				readyToFinalize: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
				parsedToLP: map[logKey]logpoller.Log{
					{txHash: common.HexToHash("0x1"), logIndex: 1}: {
						BlockTimestamp: mustParseTime(t, "2021-01-01T00:00:00Z"),
					},
					{txHash: common.HexToHash("0x2"), logIndex: 2}: {
						BlockTimestamp: mustParseTime(t, "2021-01-02T00:00:00Z"),
					},
				},
			},
			[]models.PendingTransfer{
				{
					Transfer: models.Transfer{
						From:               sourceSelector,
						To:                 destSelector,
						Sender:             sourceAdapterAddress,
						Receiver:           destLiquidityManagerAddress,
						LocalTokenAddress:  localToken,
						RemoteTokenAddress: remoteToken,
						Date:               mustParseTime(t, "2021-01-01T00:00:00Z"),
						Amount:             ubig.NewI(1),
						BridgeData:         mustPackProvePayload(t, big.NewInt(1)),
						Stage:              1,
					},
					Status: models.TransferStatusReady,
					ID:     fmt.Sprintf("%s-%d-prove", common.HexToHash("0x1"), 1),
				},
				{
					Transfer: models.Transfer{
						From:               sourceSelector,
						To:                 destSelector,
						Sender:             sourceAdapterAddress,
						Receiver:           destLiquidityManagerAddress,
						LocalTokenAddress:  localToken,
						RemoteTokenAddress: remoteToken,
						Date:               mustParseTime(t, "2021-01-02T00:00:00Z"),
						Amount:             ubig.NewI(2),
						BridgeData:         mustPackProvePayload(t, big.NewInt(2)),
						Stage:              1,
					},
					Status: models.TransferStatusReady,
					ID:     fmt.Sprintf("%s-%d-prove", common.HexToHash("0x2"), 2),
				},
			},
		},
		{
			"some pending transfers, all ready to finalize",
			fields{
				sourceSelector:       sourceSelector,
				destSelector:         destSelector,
				destLiquidityManager: destLiquidityManager,
				sourceAdapter:        sourceAdapter,
				lggr:                 logger.TestLogger(t),
			},
			args{
				localToken:   localToken,
				remoteToken:  remoteToken,
				readyToProve: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
				readyToFinalize: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						Amount:           big.NewInt(1),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(1)),
						Raw: types.Log{
							TxHash: common.HexToHash("0x1"),
							Index:  1,
						},
					},
					{
						Amount:           big.NewInt(2),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(2)),
						Raw: types.Log{
							TxHash: common.HexToHash("0x2"),
							Index:  2,
						},
					},
				},
				parsedToLP: map[logKey]logpoller.Log{
					{txHash: common.HexToHash("0x1"), logIndex: 1}: {
						BlockTimestamp: mustParseTime(t, "2021-01-01T00:00:00Z"),
					},
					{txHash: common.HexToHash("0x2"), logIndex: 2}: {
						BlockTimestamp: mustParseTime(t, "2021-01-02T00:00:00Z"),
					},
				},
			},
			[]models.PendingTransfer{
				{
					Transfer: models.Transfer{
						From:               sourceSelector,
						To:                 destSelector,
						Sender:             sourceAdapterAddress,
						Receiver:           destLiquidityManagerAddress,
						LocalTokenAddress:  localToken,
						RemoteTokenAddress: remoteToken,
						Date:               mustParseTime(t, "2021-01-01T00:00:00Z"),
						Amount:             ubig.NewI(1),
						BridgeData:         mustPackFinalizePayload(t, big.NewInt(1), big.NewInt(1)),
						Stage:              2,
					},
					Status: models.TransferStatusReady,
					ID:     fmt.Sprintf("%s-%d-finalize", common.HexToHash("0x1"), 1),
				},
				{
					Transfer: models.Transfer{
						From:               sourceSelector,
						To:                 destSelector,
						Sender:             sourceAdapterAddress,
						Receiver:           destLiquidityManagerAddress,
						LocalTokenAddress:  localToken,
						RemoteTokenAddress: remoteToken,
						Date:               mustParseTime(t, "2021-01-02T00:00:00Z"),
						Amount:             ubig.NewI(2),
						BridgeData:         mustPackFinalizePayload(t, big.NewInt(2), big.NewInt(2)),
						Stage:              2,
					},
					Status: models.TransferStatusReady,
					ID:     fmt.Sprintf("%s-%d-finalize", common.HexToHash("0x2"), 2),
				},
			},
		},
		{
			"some pending transfers, someready to finalize, some ready to prove",
			fields{
				sourceSelector:       sourceSelector,
				destSelector:         destSelector,
				destLiquidityManager: destLiquidityManager,
				sourceAdapter:        sourceAdapter,
				lggr:                 logger.TestLogger(t),
			},
			args{
				localToken:  localToken,
				remoteToken: remoteToken,
				readyToProve: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						Amount:           big.NewInt(3),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(3)),
						Raw: types.Log{
							TxHash: common.HexToHash("0x3"),
							Index:  3,
						},
					},
					{
						Amount:           big.NewInt(4),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(4)),
						Raw: types.Log{
							TxHash: common.HexToHash("0x4"),
							Index:  4,
						},
					},
				},
				readyToFinalize: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						Amount:           big.NewInt(1),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(1)),
						Raw: types.Log{
							TxHash: common.HexToHash("0x1"),
							Index:  1,
						},
					},
					{
						Amount:           big.NewInt(2),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(2)),
						Raw: types.Log{
							TxHash: common.HexToHash("0x2"),
							Index:  2,
						},
					},
				},
				parsedToLP: map[logKey]logpoller.Log{
					{txHash: common.HexToHash("0x1"), logIndex: 1}: {
						BlockTimestamp: mustParseTime(t, "2021-01-01T00:00:00Z"),
					},
					{txHash: common.HexToHash("0x2"), logIndex: 2}: {
						BlockTimestamp: mustParseTime(t, "2021-01-02T00:00:00Z"),
					},
					{txHash: common.HexToHash("0x3"), logIndex: 3}: {
						BlockTimestamp: mustParseTime(t, "2021-01-03T00:00:00Z"),
					},
					{txHash: common.HexToHash("0x4"), logIndex: 4}: {
						BlockTimestamp: mustParseTime(t, "2021-01-04T00:00:00Z"),
					},
				},
			},
			[]models.PendingTransfer{
				{
					Transfer: models.Transfer{
						From:               sourceSelector,
						To:                 destSelector,
						Sender:             sourceAdapterAddress,
						Receiver:           destLiquidityManagerAddress,
						LocalTokenAddress:  localToken,
						RemoteTokenAddress: remoteToken,
						Date:               mustParseTime(t, "2021-01-03T00:00:00Z"),
						Amount:             ubig.NewI(3),
						BridgeData:         mustPackProvePayload(t, big.NewInt(3)),
						Stage:              1,
					},
					Status: models.TransferStatusReady,
					ID:     fmt.Sprintf("%s-%d-prove", common.HexToHash("0x3"), 3),
				},
				{
					Transfer: models.Transfer{
						From:               sourceSelector,
						To:                 destSelector,
						Sender:             sourceAdapterAddress,
						Receiver:           destLiquidityManagerAddress,
						LocalTokenAddress:  localToken,
						RemoteTokenAddress: remoteToken,
						Date:               mustParseTime(t, "2021-01-04T00:00:00Z"),
						Amount:             ubig.NewI(4),
						BridgeData:         mustPackProvePayload(t, big.NewInt(4)),
						Stage:              1,
					},
					Status: models.TransferStatusReady,
					ID:     fmt.Sprintf("%s-%d-prove", common.HexToHash("0x4"), 4),
				},
				{
					Transfer: models.Transfer{
						From:               sourceSelector,
						To:                 destSelector,
						Sender:             sourceAdapterAddress,
						Receiver:           destLiquidityManagerAddress,
						LocalTokenAddress:  localToken,
						RemoteTokenAddress: remoteToken,
						Date:               mustParseTime(t, "2021-01-01T00:00:00Z"),
						Amount:             ubig.NewI(1),
						BridgeData:         mustPackFinalizePayload(t, big.NewInt(1), big.NewInt(1)),
						Stage:              2,
					},
					Status: models.TransferStatusReady,
					ID:     fmt.Sprintf("%s-%d-finalize", common.HexToHash("0x1"), 1),
				},
				{
					Transfer: models.Transfer{
						From:               sourceSelector,
						To:                 destSelector,
						Sender:             sourceAdapterAddress,
						Receiver:           destLiquidityManagerAddress,
						LocalTokenAddress:  localToken,
						RemoteTokenAddress: remoteToken,
						Date:               mustParseTime(t, "2021-01-02T00:00:00Z"),
						Amount:             ubig.NewI(2),
						BridgeData:         mustPackFinalizePayload(t, big.NewInt(2), big.NewInt(2)),
						Stage:              2,
					},
					Status: models.TransferStatusReady,
					ID:     fmt.Sprintf("%s-%d-finalize", common.HexToHash("0x2"), 2),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &testBridge{
				sourceSelector:       tt.fields.sourceSelector,
				destSelector:         tt.fields.destSelector,
				destLiquidityManager: tt.fields.destLiquidityManager,
				sourceAdapter:        tt.fields.sourceAdapter,
				lggr:                 tt.fields.lggr,
			}
			got, err := tr.toPendingTransfers(tt.args.localToken, tt.args.remoteToken, tt.args.readyToProve, tt.args.readyToFinalize, tt.args.parsedToLP)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func mustParseTime(t *testing.T, s string) time.Time {
	tm, err := time.Parse(time.RFC3339, s)
	require.NoError(t, err)
	return tm
}

func Test_filterFinalized(t *testing.T) {
	type args struct {
		sends     []*liquiditymanager.LiquidityManagerLiquidityTransferred
		finalizes []*liquiditymanager.LiquidityManagerLiquidityTransferred
	}
	tests := []struct {
		name    string
		args    args
		want    []*liquiditymanager.LiquidityManagerLiquidityTransferred
		wantErr bool
	}{
		{
			"no finalizes",
			args{
				[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						Amount:           big.NewInt(1),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(1)),
					},
					{
						Amount:           big.NewInt(2),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(2)),
					},
					{
						Amount:           big.NewInt(3),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(3)),
					},
				},
				[]*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			},
			[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					Amount:           big.NewInt(1),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(1)),
				},
				{
					Amount:           big.NewInt(2),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(2)),
				},
				{
					Amount:           big.NewInt(3),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(3)),
				},
			},
			false,
		},
		{
			"some finalizes",
			args{
				[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						Amount:           big.NewInt(1),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(1)),
					},
					{
						Amount:           big.NewInt(2),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(2)),
					},
					{
						Amount:           big.NewInt(3),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(3)),
					},
				},
				[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						Amount:             big.NewInt(2),
						BridgeSpecificData: mustPackFinalizePayload(t, big.NewInt(2), big.NewInt(2)),
					},
				},
			},
			[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					Amount:           big.NewInt(1),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(1)),
				},
				{
					Amount:           big.NewInt(3),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(3)),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterFinalized(tt.args.sends, tt.args.finalizes)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_groupByStage(t *testing.T) {
	type args struct {
		unfinalized    []*liquiditymanager.LiquidityManagerLiquidityTransferred
		stepsCompleted []*liquiditymanager.LiquidityManagerFinalizationStepCompleted
	}
	tests := []struct {
		name                string
		args                args
		wantReadyToProve    []*liquiditymanager.LiquidityManagerLiquidityTransferred
		wantReadyToFinalize []*liquiditymanager.LiquidityManagerLiquidityTransferred
		wantErr             bool
	}{
		{
			"all ready to prove",
			args{
				[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						Amount:           big.NewInt(1),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(1)),
					},
					{
						Amount:           big.NewInt(2),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(2)),
					},
					{
						Amount:           big.NewInt(3),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(3)),
					},
				},
				[]*liquiditymanager.LiquidityManagerFinalizationStepCompleted{}, // none proven
			},
			[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					Amount:           big.NewInt(1),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(1)),
				},
				{
					Amount:           big.NewInt(2),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(2)),
				},
				{
					Amount:           big.NewInt(3),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(3)),
				},
			},
			nil,
			false,
		},
		{
			"all ready to finalize",
			args{
				[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						Amount:           big.NewInt(1),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(1)),
					},
					{
						Amount:           big.NewInt(2),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(2)),
					},
					{
						Amount:           big.NewInt(3),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(3)),
					},
				},
				[]*liquiditymanager.LiquidityManagerFinalizationStepCompleted{ // all proven
					{
						BridgeSpecificData: mustPackProvePayload(t, big.NewInt(1)),
					},
					{
						BridgeSpecificData: mustPackProvePayload(t, big.NewInt(2)),
					},
					{
						BridgeSpecificData: mustPackProvePayload(t, big.NewInt(3)),
					},
				},
			},
			nil,
			[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					Amount:           big.NewInt(1),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(1)),
				},
				{
					Amount:           big.NewInt(2),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(2)),
				},
				{
					Amount:           big.NewInt(3),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(3)),
				},
			},
			false,
		},
		{
			"mix of ready to prove and ready to finalize",
			args{
				[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						Amount:           big.NewInt(1),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(1)),
					},
					{
						Amount:           big.NewInt(2),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(2)),
					},
					{
						Amount:           big.NewInt(3),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(3)),
					},
					{
						Amount:           big.NewInt(4),
						BridgeReturnData: mustPackSendReturnData(t, big.NewInt(4)),
					},
				},
				[]*liquiditymanager.LiquidityManagerFinalizationStepCompleted{ // 1 and 3 already proven, ready to finalize
					{
						BridgeSpecificData: mustPackProvePayload(t, big.NewInt(1)),
					},
					{
						BridgeSpecificData: mustPackProvePayload(t, big.NewInt(3)),
					},
				},
			},
			[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					Amount:           big.NewInt(2),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(2)),
				},
				{
					Amount:           big.NewInt(4),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(4)),
				},
			},
			[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					Amount:           big.NewInt(1),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(1)),
				},
				{
					Amount:           big.NewInt(3),
					BridgeReturnData: mustPackSendReturnData(t, big.NewInt(3)),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReadyToProve, gotReadyToFinalize, err := groupByStage(tt.args.unfinalized, tt.args.stepsCompleted)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantReadyToProve, gotReadyToProve)
				require.Equal(t, tt.wantReadyToFinalize, gotReadyToFinalize)
			}
		})
	}
}

func mustPackSendReturnData(t *testing.T, nonce *big.Int) []byte {
	packed, err := utils.ABIEncode(`[{"type": "uint256"}]`, nonce)
	require.NoError(t, err)
	return packed
}

func mustPackFinalizePayload(t *testing.T, nonce, amount *big.Int) []byte {
	packed, err := PackFinalizeBridgePayload(amount, nonce)
	require.NoError(t, err)
	return packed
}

func mustPackProvePayload(t *testing.T, nonce *big.Int) []byte {
	packed, err := PackProveBridgePayload(nonce)
	require.NoError(t, err)
	return packed
}

func TestNew(t *testing.T) {
	var (
		sourceSelector                = models.NetworkSelector(1)
		destSelector                  = models.NetworkSelector(2)
		sourceLiquidityManagerAddress = models.Address(testutils.NewAddress())
		destLiquidityManagerAddress   = models.Address(testutils.NewAddress())
		sourceAdapterAddress          = models.Address(testutils.NewAddress())
		destAdapterAddress            = models.Address(testutils.NewAddress())
	)
	type args struct {
		sourceSelector                models.NetworkSelector
		destSelector                  models.NetworkSelector
		sourceLiquidityManagerAddress models.Address
		destLiquidityManagerAddress   models.Address
		sourceAdapter                 models.Address
		destAdapter                   models.Address
		sourceLogPoller               *lp_mocks.LogPoller
		destLogPoller                 *lp_mocks.LogPoller
		sourceClient                  *evmclient_mocks.Client
		destClient                    *evmclient_mocks.Client
		lggr                          logger.Logger
	}
	tests := []struct {
		name    string
		args    args
		expect  func(t *testing.T, args args)
		assert  func(t *testing.T, args args)
		wantErr bool
	}{
		{
			"happy path",
			args{
				sourceSelector:                sourceSelector,
				destSelector:                  destSelector,
				sourceLiquidityManagerAddress: sourceLiquidityManagerAddress,
				destLiquidityManagerAddress:   destLiquidityManagerAddress,
				sourceAdapter:                 sourceAdapterAddress,
				destAdapter:                   destAdapterAddress,
				sourceLogPoller:               lp_mocks.NewLogPoller(t),
				destLogPoller:                 lp_mocks.NewLogPoller(t),
				sourceClient:                  evmclient_mocks.NewClient(t),
				destClient:                    evmclient_mocks.NewClient(t),
				lggr:                          logger.TestLogger(t),
			},
			func(t *testing.T, args args) {
				args.sourceLogPoller.On("RegisterFilter", mock.Anything, mock.MatchedBy(func(f logpoller.Filter) bool {
					ok := len(f.Addresses) == 1
					ok = ok && f.Addresses[0] == common.Address(args.sourceLiquidityManagerAddress)
					ok = ok && len(f.EventSigs) == 2
					ok = ok && f.EventSigs[0] == LiquidityTransferredTopic
					ok = ok && f.EventSigs[1] == FinalizationStepCompletedTopic
					ok = ok && strings.HasPrefix(f.Name, "Local-LiquidityTransferred-FinalizationCompleted")
					return ok
				})).Return(nil)
				args.destLogPoller.On("RegisterFilter", mock.Anything, mock.MatchedBy(func(f logpoller.Filter) bool {
					ok := len(f.Addresses) == 1
					ok = ok && f.Addresses[0] == common.Address(args.destLiquidityManagerAddress)
					ok = ok && len(f.EventSigs) == 2
					ok = ok && f.EventSigs[0] == LiquidityTransferredTopic
					ok = ok && f.EventSigs[1] == FinalizationStepCompletedTopic
					ok = ok && strings.HasPrefix(f.Name, "Remote-LiquidityTransferred-FinalizationCompleted")
					return ok
				})).Return(nil)
			},
			func(t *testing.T, args args) {
				args.sourceLogPoller.AssertExpectations(t)
				args.destLogPoller.AssertExpectations(t)
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expect(t, tt.args)
			defer tt.assert(t, tt.args)
			got, err := New(testutils.Context(t), tt.args.sourceSelector, tt.args.destSelector, tt.args.sourceLiquidityManagerAddress, tt.args.destLiquidityManagerAddress, tt.args.sourceAdapter, tt.args.destAdapter, tt.args.sourceLogPoller, tt.args.destLogPoller, tt.args.sourceClient, tt.args.destClient, tt.args.lggr)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
				tt.expect(t, tt.args)
			}
		})
	}
}
