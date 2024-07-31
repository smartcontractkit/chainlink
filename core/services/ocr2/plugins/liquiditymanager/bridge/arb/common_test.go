package arb_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/l2_arbitrum_gateway"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/arb"
	bridgecommon "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/common"
)

func Test_TopicIndexes(t *testing.T) {
	var (
		rebalancerABI = abihelpers.MustParseABI(liquiditymanager.LiquidityManagerMetaData.ABI)
		l2GatewayABI  = abihelpers.MustParseABI(l2_arbitrum_gateway.L2ArbitrumGatewayMetaData.ABI)
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

	t.Run("deposit finalized to address idx", func(t *testing.T) {
		dfEvent, ok := l2GatewayABI.Events["DepositFinalized"]
		require.True(t, ok)

		var toAddressArg abi.Argument
		var topicIndex = 0
		for _, arg := range dfEvent.Inputs {
			if arg.Indexed {
				topicIndex++
			}
			if arg.Name == "_to" {
				toAddressArg = arg
				break
			}
		}

		require.True(t, toAddressArg.Indexed)
		require.Equal(t, arb.DepositFinalizedToAddressTopicIndex, topicIndex)
	})
}
