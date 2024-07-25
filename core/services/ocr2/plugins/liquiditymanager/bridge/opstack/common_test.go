package opstack

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	bridgetestutils "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/testutils"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_standard_bridge"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	bridgecommon "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/common"
)

func Test_TopicIndexes(t *testing.T) {
	var (
		rebalancerABI     = abihelpers.MustParseABI(liquiditymanager.LiquidityManagerMetaData.ABI)
		standardBridgeABI = abihelpers.MustParseABI(optimism_standard_bridge.OptimismStandardBridgeMetaData.ABI)
	)
	t.Run("liquidity transferred to chain selector idx", func(t *testing.T) {
		ltEvent, ok := rebalancerABI.Events["LiquidityTransferred"]
		require.True(t, ok)

		var toChainSelectorArg abi.Argument
		var topicIndex = 0
		for _, arg := range ltEvent.Inputs {
			if arg.Indexed {
				topicIndex++
			}
			if arg.Name == "toChainSelector" {
				toChainSelectorArg = arg
				break
			}
		}

		require.True(t, toChainSelectorArg.Indexed)
		require.Equal(t, bridgecommon.LiquidityTransferredToChainSelectorTopicIndex, topicIndex)
	})

	t.Run("liquidity transferred from chain selector idx", func(t *testing.T) {
		ltEvent, ok := rebalancerABI.Events["LiquidityTransferred"]
		require.True(t, ok)

		var fromChainSelectorArg abi.Argument
		var topicIndex = 0
		for _, arg := range ltEvent.Inputs {
			if arg.Indexed {
				topicIndex++
			}
			if arg.Name == "fromChainSelector" {
				fromChainSelectorArg = arg
				break
			}
		}

		require.True(t, fromChainSelectorArg.Indexed)
		require.Equal(t, bridgecommon.LiquidityTransferredFromChainSelectorTopicIndex, topicIndex)
	})

	t.Run("ERC20 bridge finalized to address idx", func(t *testing.T) {
		bfEvent, ok := standardBridgeABI.Events["ERC20BridgeFinalized"]
		require.True(t, ok)

		var fromAddressArg abi.Argument
		var topicIndex = 0
		for _, arg := range bfEvent.Inputs {
			if arg.Indexed {
				topicIndex++
			}
			if arg.Name == "from" {
				fromAddressArg = arg
				break
			}
		}

		require.True(t, fromAddressArg.Indexed)
		require.Equal(t, ERC20BridgeFinalizedFromAddressTopicIndex, topicIndex)
	})
}

func Test_filterExecuted(t *testing.T) {
	l2LiquidityManagerAddress := common.HexToAddress("0xabc")
	l1ChainSelector := chainsel.ETHEREUM_MAINNET.Selector
	l2ChainSelector := chainsel.ETHEREUM_MAINNET_OPTIMISM_1.Selector
	type args struct {
		readyCandidates []*liquiditymanager.LiquidityManagerLiquidityTransferred
		receivedLogs    []*liquiditymanager.LiquidityManagerLiquidityTransferred
	}
	tests := []struct {
		name      string
		args      args
		wantReady []*liquiditymanager.LiquidityManagerLiquidityTransferred
		wantErr   bool
	}{
		{
			name: "no logs",
			args: args{
				readyCandidates: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
				receivedLogs:    []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			},
			wantReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			wantErr:   false,
		},
		{
			name: "no received logs",
			args: args{
				readyCandidates: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						FromChainSelector:  l1ChainSelector,
						ToChainSelector:    l2ChainSelector,
						To:                 l2LiquidityManagerAddress,
						Amount:             big.NewInt(1),
						BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
						BridgeSpecificData: []byte{},
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			},
			wantReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					FromChainSelector:  l1ChainSelector,
					ToChainSelector:    l2ChainSelector,
					To:                 l2LiquidityManagerAddress,
					Amount:             big.NewInt(1),
					BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					BridgeSpecificData: []byte{},
				},
			},
			wantErr: false,
		},
		{
			name: "mismatched nonces",
			args: args{
				readyCandidates: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						FromChainSelector: l1ChainSelector,
						ToChainSelector:   l2ChainSelector,
						To:                l2LiquidityManagerAddress,
						Amount:            big.NewInt(1),
						// nonce = 1
						BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
						BridgeSpecificData: []byte{},
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						FromChainSelector: l1ChainSelector,
						ToChainSelector:   l2ChainSelector,
						To:                l2LiquidityManagerAddress,
						Amount:            big.NewInt(1),
						BridgeReturnData:  []byte{},
						// nonce = 2
						BridgeSpecificData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
					}},
			},
			wantReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
				{
					FromChainSelector: l1ChainSelector,
					ToChainSelector:   l2ChainSelector,
					To:                l2LiquidityManagerAddress,
					Amount:            big.NewInt(1),
					// nonce = 1
					BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					BridgeSpecificData: []byte{},
				},
			},
			wantErr: false,
		},
		{
			name: "matching nonces",
			args: args{
				readyCandidates: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						FromChainSelector: l1ChainSelector,
						ToChainSelector:   l2ChainSelector,
						To:                l2LiquidityManagerAddress,
						Amount:            big.NewInt(1),
						// nonce = 1
						BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
						BridgeSpecificData: []byte{},
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						FromChainSelector: l1ChainSelector,
						ToChainSelector:   l2ChainSelector,
						To:                l2LiquidityManagerAddress,
						Amount:            big.NewInt(1),
						BridgeReturnData:  []byte{},
						// nonce = 1
						BridgeSpecificData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					}},
			},
			wantReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{},
			wantErr:   false,
		},
		{
			name: "multiple logs, some matching, some not",
			args: args{
				readyCandidates: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						FromChainSelector: l1ChainSelector,
						ToChainSelector:   l2ChainSelector,
						To:                l2LiquidityManagerAddress,
						Amount:            big.NewInt(1),
						// nonce = 1, is executed
						BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
						BridgeSpecificData: []byte{},
					},
					{
						FromChainSelector: l1ChainSelector,
						ToChainSelector:   l2ChainSelector,
						To:                l2LiquidityManagerAddress,
						Amount:            big.NewInt(1),
						// nonce = 2, not executed
						BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
						BridgeSpecificData: []byte{},
					},
					{
						FromChainSelector: l1ChainSelector,
						ToChainSelector:   l2ChainSelector,
						To:                l2LiquidityManagerAddress,
						Amount:            big.NewInt(1),
						// nonce = 3, is executed
						BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000003"),
						BridgeSpecificData: []byte{},
					},
				},
				receivedLogs: []*liquiditymanager.LiquidityManagerLiquidityTransferred{
					{
						FromChainSelector: l1ChainSelector,
						ToChainSelector:   l2ChainSelector,
						To:                l2LiquidityManagerAddress,
						Amount:            big.NewInt(1),
						BridgeReturnData:  []byte{},
						// nonce = 1
						BridgeSpecificData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000001"),
					},
					{
						FromChainSelector: l1ChainSelector,
						ToChainSelector:   l2ChainSelector,
						To:                l2LiquidityManagerAddress,
						Amount:            big.NewInt(1),
						BridgeReturnData:  []byte{},
						// nonce = 3
						BridgeSpecificData: bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000003"),
					},
				},
			},
			wantReady: []*liquiditymanager.LiquidityManagerLiquidityTransferred{

				{
					FromChainSelector: l1ChainSelector,
					ToChainSelector:   l2ChainSelector,
					To:                l2LiquidityManagerAddress,
					Amount:            big.NewInt(1),
					// nonce = 2
					BridgeReturnData:   bridgetestutils.MustPackBridgeData(t, "0x0000000000000000000000000000000000000000000000000000000000000002"),
					BridgeSpecificData: []byte{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReady, err := filterExecuted(tt.args.readyCandidates, tt.args.receivedLogs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				bridgetestutils.AssertLiquidityTransferredEventSlicesEqual(t, tt.wantReady, gotReady, bridgetestutils.SortByBridgeReturnData)
			}
		})
	}
}
