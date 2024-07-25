package arb

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_rollup_core"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/mocks/mock_arbitrum_rollup_core"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
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
		ctx context.Context //nolint:containedctx
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
				f.l1LogPoller.On("LatestLogByEventSigWithConfs", mock.Anything, NodeConfirmedTopic, rollupAddress, evmtypes.Finalized).
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
				f.l1LogPoller.On("LatestLogByEventSigWithConfs", mock.Anything, NodeConfirmedTopic, rollupAddress, evmtypes.Finalized, mock.Anything).
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
				f.l1LogPoller.On("LatestLogByEventSigWithConfs", mock.Anything, NodeConfirmedTopic, rollupAddress, evmtypes.Finalized, mock.Anything).
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

func Test_unpackFinalizationPayload(t *testing.T) {
	type args struct {
		calldata []byte
	}
	tests := []struct {
		name    string
		args    args
		want    arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterArbitrumFinalizationPayload
		wantErr bool
	}{
		{
			// from an actual tx: https://sepolia.etherscan.io/tx/0xdbb29a504d33d575ea8ed0342efccc051f409509b96399e00999c4aa98e174a8#eventlog
			"success",
			args{
				calldata: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000001200000000000000000000000000000000000000000000000000000000000005DD6000000000000000000000000CFB1F08A4852699A979909E22C30263CA249556D000000000000000000000000A8AD8D7E13CBF556EE75CB0324C13535D8100E1E0000000000000000000000000000000000000000000000000000000000F83B630000000000000000000000000000000000000000000000000000000000516A3E0000000000000000000000000000000000000000000000000000000065D6369D000000000000000000000000000000000000000000000000000000003B9ACA000000000000000000000000000000000000000000000000000000000000000320000000000000000000000000000000000000000000000000000000000000000F716CB7F674AF62C558D059D92581103870CCC17902F0B8D50730B651DF418997260A4DCC8EAD2B50B03DFF96F31F63AB26FA9F71AAB423D3E9DF04AA185131D4DF4FDEF41AAD6F03C68D859C4E53D92005AB6B3485AA1F048ED9FDD66DD5ED71AB99AECFDD7475CE024858197D70D5F02CBAB94593423544C5F655185864AAE1F02C4844E5510583C027EFED35E7CB7B205E102842AC1CD81D222C36AB138904000000000000000000000000000000000000000000000000000000000000000097046C56B8CBA34F73B51D0A200A7DC0D17A2FC07DC86D1E0FC3F6F6E263234BE9B1964C5F4428F87CB49C7EF72656490927F4786370A1B655CAE8F2CAEFED08A82CEDA610907DF10408E3602DAB2098FAE7AAE0FA9DB2DDFD88F903DC82F91A0000000000000000000000000000000000000000000000000000000000000000989E42AF081500D142FE5584375A58B95BE061BB5CDF38A3579F5043D90052B284AE48F37A0C35683DF949536B6576CE97CAA855C176CCD7D769494015B9BA3E4BDE3FBFA26B03BFCBBC80884A5C229D987992CE2366BEF5882742CAC7CC7EFE000000000000000000000000000000000000000000000000000000000000000037EF20C8140B820C15303931A9E839810EC66F4D9523C3FA8187434076A25B1300000000000000000000000000000000000000000000000000000000000001242E567B360000000000000000000000007B79995E5F793A07BC00C21412E50ECAE098E7F9000000000000000000000000E97467A3BDA1FAC5051D2CFDAC8B5F28FAA65788000000000000000000000000F9BB721BD68F5EED40CB7C9DDCC0F0BA9D0B1B96000000000000000000000000000000000000000000000000000000003B9ACA0000000000000000000000000000000000000000000000000000000000000000A0000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000110000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
			},
			arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterArbitrumFinalizationPayload{
				L2Sender:    common.HexToAddress("0xCFB1f08A4852699a979909e22c30263ca249556D"),
				To:          common.HexToAddress("0xA8aD8d7e13cbf556eE75CB0324c13535d8100e1E"),
				L2Block:     big.NewInt(16268131),
				L1Block:     big.NewInt(5335614),
				L2Timestamp: big.NewInt(1708537501),
				Value:       big.NewInt(1000000000),
				Proof: [][32]uint8{
					{0x71, 0x6c, 0xb7, 0xf6, 0x74, 0xaf, 0x62, 0xc5, 0x58, 0xd0, 0x59, 0xd9, 0x25, 0x81, 0x10, 0x38, 0x70, 0xcc, 0xc1, 0x79, 0x2, 0xf0, 0xb8, 0xd5, 0x7, 0x30, 0xb6, 0x51, 0xdf, 0x41, 0x89, 0x97},
					{0x26, 0xa, 0x4d, 0xcc, 0x8e, 0xad, 0x2b, 0x50, 0xb0, 0x3d, 0xff, 0x96, 0xf3, 0x1f, 0x63, 0xab, 0x26, 0xfa, 0x9f, 0x71, 0xaa, 0xb4, 0x23, 0xd3, 0xe9, 0xdf, 0x4, 0xaa, 0x18, 0x51, 0x31, 0xd4},
					{0xdf, 0x4f, 0xde, 0xf4, 0x1a, 0xad, 0x6f, 0x3, 0xc6, 0x8d, 0x85, 0x9c, 0x4e, 0x53, 0xd9, 0x20, 0x5, 0xab, 0x6b, 0x34, 0x85, 0xaa, 0x1f, 0x4, 0x8e, 0xd9, 0xfd, 0xd6, 0x6d, 0xd5, 0xed, 0x71},
					{0xab, 0x99, 0xae, 0xcf, 0xdd, 0x74, 0x75, 0xce, 0x2, 0x48, 0x58, 0x19, 0x7d, 0x70, 0xd5, 0xf0, 0x2c, 0xba, 0xb9, 0x45, 0x93, 0x42, 0x35, 0x44, 0xc5, 0xf6, 0x55, 0x18, 0x58, 0x64, 0xaa, 0xe1},
					{0xf0, 0x2c, 0x48, 0x44, 0xe5, 0x51, 0x5, 0x83, 0xc0, 0x27, 0xef, 0xed, 0x35, 0xe7, 0xcb, 0x7b, 0x20, 0x5e, 0x10, 0x28, 0x42, 0xac, 0x1c, 0xd8, 0x1d, 0x22, 0x2c, 0x36, 0xab, 0x13, 0x89, 0x4},
					{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
					{0x97, 0x4, 0x6c, 0x56, 0xb8, 0xcb, 0xa3, 0x4f, 0x73, 0xb5, 0x1d, 0xa, 0x20, 0xa, 0x7d, 0xc0, 0xd1, 0x7a, 0x2f, 0xc0, 0x7d, 0xc8, 0x6d, 0x1e, 0xf, 0xc3, 0xf6, 0xf6, 0xe2, 0x63, 0x23, 0x4b},
					{0xe9, 0xb1, 0x96, 0x4c, 0x5f, 0x44, 0x28, 0xf8, 0x7c, 0xb4, 0x9c, 0x7e, 0xf7, 0x26, 0x56, 0x49, 0x9, 0x27, 0xf4, 0x78, 0x63, 0x70, 0xa1, 0xb6, 0x55, 0xca, 0xe8, 0xf2, 0xca, 0xef, 0xed, 0x8},
					{0xa8, 0x2c, 0xed, 0xa6, 0x10, 0x90, 0x7d, 0xf1, 0x4, 0x8, 0xe3, 0x60, 0x2d, 0xab, 0x20, 0x98, 0xfa, 0xe7, 0xaa, 0xe0, 0xfa, 0x9d, 0xb2, 0xdd, 0xfd, 0x88, 0xf9, 0x3, 0xdc, 0x82, 0xf9, 0x1a},
					{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
					{0x98, 0x9e, 0x42, 0xaf, 0x8, 0x15, 0x0, 0xd1, 0x42, 0xfe, 0x55, 0x84, 0x37, 0x5a, 0x58, 0xb9, 0x5b, 0xe0, 0x61, 0xbb, 0x5c, 0xdf, 0x38, 0xa3, 0x57, 0x9f, 0x50, 0x43, 0xd9, 0x0, 0x52, 0xb2},
					{0x84, 0xae, 0x48, 0xf3, 0x7a, 0xc, 0x35, 0x68, 0x3d, 0xf9, 0x49, 0x53, 0x6b, 0x65, 0x76, 0xce, 0x97, 0xca, 0xa8, 0x55, 0xc1, 0x76, 0xcc, 0xd7, 0xd7, 0x69, 0x49, 0x40, 0x15, 0xb9, 0xba, 0x3e},
					{0x4b, 0xde, 0x3f, 0xbf, 0xa2, 0x6b, 0x3, 0xbf, 0xcb, 0xbc, 0x80, 0x88, 0x4a, 0x5c, 0x22, 0x9d, 0x98, 0x79, 0x92, 0xce, 0x23, 0x66, 0xbe, 0xf5, 0x88, 0x27, 0x42, 0xca, 0xc7, 0xcc, 0x7e, 0xfe},
					{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
					{0x37, 0xef, 0x20, 0xc8, 0x14, 0xb, 0x82, 0xc, 0x15, 0x30, 0x39, 0x31, 0xa9, 0xe8, 0x39, 0x81, 0xe, 0xc6, 0x6f, 0x4d, 0x95, 0x23, 0xc3, 0xfa, 0x81, 0x87, 0x43, 0x40, 0x76, 0xa2, 0x5b, 0x13},
				},
				Index: big.NewInt(24022),
				Data: []uint8{0x2e, 0x56, 0x7b, 0x36, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b, 0x79, 0x99, 0x5e, 0x5f, 0x79, 0x3a, 0x7, 0xbc, 0x0, 0xc2, 0x14, 0x12, 0xe5, 0xe, 0xca, 0xe0,
					0x98, 0xe7, 0xf9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xe9, 0x74, 0x67, 0xa3, 0xbd, 0xa1, 0xfa, 0xc5, 0x5, 0x1d, 0x2c, 0xfd, 0xac, 0x8b, 0x5f, 0x28, 0xfa, 0xa6, 0x57,
					0x88, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf9, 0xbb, 0x72, 0x1b, 0xd6, 0x8f, 0x5e, 0xed, 0x40, 0xcb, 0x7c, 0x9d, 0xdc, 0xc0, 0xf0, 0xba, 0x9d, 0xb, 0x1b, 0x96, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3b, 0x9a, 0xca, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x60, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x11, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40,
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			},
			false,
		},
		{
			"invalid calldata",
			args{
				calldata: []byte{0x01, 0x02, 0x03},
			},
			arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterArbitrumFinalizationPayload{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unpackFinalizationPayload(tt.args.calldata)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_filterUnfinalizedTransfers(t *testing.T) {
	type args struct {
		sentLogs     []*liquiditymanager.LiquidityManagerLiquidityTransferred
		receivedLogs []*liquiditymanager.LiquidityManagerLiquidityTransferred
	}
	var (
		l1LiquidityManager = testutils.NewAddress()
	)
	tests := []struct {
		name    string
		args    args
		want    []*liquiditymanager.LiquidityManagerLiquidityTransferred
		wantErr bool
	}{
		{
			"no sent or received",
			args{
				sentLogs:     nil,
				receivedLogs: nil,
			},
			nil,
			false,
		},
		{
			"some sent no received",
			args{
				sentLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						OcrSeqNum:          1,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l1LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: []byte{},
						BridgeReturnData:   []byte{},
					},
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l1LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: []byte{},
						BridgeReturnData:   []byte{},
					},
				},
				receivedLogs: nil,
			},
			[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					OcrSeqNum:          1,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l1LiquidityManager,
					Amount:             big.NewInt(100),
					BridgeSpecificData: []byte{},
					BridgeReturnData:   []byte{},
				},
				{
					OcrSeqNum:          2,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l1LiquidityManager,
					Amount:             big.NewInt(100),
					BridgeSpecificData: []byte{},
					BridgeReturnData:   []byte{},
				},
			},
			false,
		},
		{
			"some sent some received don't match",
			args{
				sentLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						OcrSeqNum:          1,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l1LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: []byte{},
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(3)),
					},
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l1LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: []byte{},
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(4)),
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						OcrSeqNum:          1,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 testutils.NewAddress(),
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackFinalizationPayload(t, big.NewInt(1)),
						BridgeReturnData:   []byte{},
					},
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 testutils.NewAddress(),
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackFinalizationPayload(t, big.NewInt(2)),
						BridgeReturnData:   []byte{},
					},
				},
			},
			[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					OcrSeqNum:          1,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l1LiquidityManager,
					Amount:             big.NewInt(100),
					BridgeSpecificData: []byte{},
					BridgeReturnData:   mustPackReturnData(t, big.NewInt(3)),
				},
				{
					OcrSeqNum:          2,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l1LiquidityManager,
					Amount:             big.NewInt(100),
					BridgeSpecificData: []byte{},
					BridgeReturnData:   mustPackReturnData(t, big.NewInt(4)),
				},
			},
			false,
		},
		{
			"some sent some received match",
			args{
				sentLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						OcrSeqNum:          1,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l1LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: []byte{},
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(3)),
					},
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l1LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: []byte{},
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(4)),
					},
					{
						OcrSeqNum:          3,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 l1LiquidityManager,
						Amount:             big.NewInt(100),
						BridgeSpecificData: []byte{},
						BridgeReturnData:   mustPackReturnData(t, big.NewInt(5)),
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						OcrSeqNum:          1,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 testutils.NewAddress(),
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackFinalizationPayload(t, big.NewInt(1)),
						BridgeReturnData:   []byte{},
					},
					{
						OcrSeqNum:          2,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 testutils.NewAddress(),
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackFinalizationPayload(t, big.NewInt(2)),
						BridgeReturnData:   []byte{},
					},
					{
						OcrSeqNum:          3,
						FromChainSelector:  10,
						ToChainSelector:    20,
						To:                 testutils.NewAddress(),
						Amount:             big.NewInt(100),
						BridgeSpecificData: mustPackFinalizationPayload(t, big.NewInt(3)),
						BridgeReturnData:   []byte{},
					},
				},
			},
			[]*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					OcrSeqNum:          2,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l1LiquidityManager,
					Amount:             big.NewInt(100),
					BridgeSpecificData: []byte{},
					BridgeReturnData:   mustPackReturnData(t, big.NewInt(4)),
				},
				{
					OcrSeqNum:          3,
					FromChainSelector:  10,
					ToChainSelector:    20,
					To:                 l1LiquidityManager,
					Amount:             big.NewInt(100),
					BridgeSpecificData: []byte{},
					BridgeReturnData:   mustPackReturnData(t, big.NewInt(5)),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterUnfinalizedTransfers(tt.args.sentLogs, tt.args.receivedLogs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func mustPackFinalizationPayload(t *testing.T, finalizationIndex *big.Int) []byte {
	payload := arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterArbitrumFinalizationPayload{
		Proof:       make([][32]byte, 1),
		L2Sender:    testutils.NewAddress(),
		To:          testutils.NewAddress(),
		L2Block:     big.NewInt(100),
		L1Block:     big.NewInt(100),
		L2Timestamp: big.NewInt(2000),
		Value:       big.NewInt(0),
		Data:        []byte{},
		// only thing that matters for detecting finalization
		Index: finalizationIndex,
	}
	packed, err := l1AdapterABI.Methods["exposeArbitrumFinalizationPayload"].Inputs.Pack(payload)
	require.NoError(t, err)
	return packed
}
