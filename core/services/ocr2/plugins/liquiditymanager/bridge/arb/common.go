package arb

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arb_node_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_rollup_core"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbsys"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/l2_arbitrum_gateway"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/l2_arbitrum_messenger"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/liquiditymanager"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

const (
	DurationMonth = 720 * time.Hour

	// LiquidityTransferredToChainSelectorTopicIndex is the index of the topic in the LiquidityTransferred event
	// that contains the "to" chain selector.
	LiquidityTransferredToChainSelectorTopicIndex = 3
	// LiquidityTransferredFromChainSelectorTopicIndex is the index of the topic in the LiquidityTransferred event
	// that contains the "from" chain selector.
	LiquidityTransferredFromChainSelectorTopicIndex = 2
	// DepositFinalizedToAddressTopicIndex is the index of the topic in the DepositFinalized event
	// that contains the "to" address.
	DepositFinalizedToAddressTopicIndex = 3

	// Arbitrum stages
	// StageRebalanceConfirmed is set as the transfer stage when the rebalanceLiquidity tx is confirmed onchain.
	StageRebalanceConfirmed = 1
	// StageFinalizeReady is set as the transfer stage when the finalization is ready to execute onchain.
	StageFinalizeReady = 2
	// StageFinalizeConfirmed is set as the transfer stage when the finalization is confirmed onchain.
	// This is a terminal stage.
	StageFinalizeConfirmed = 3
)

var (
	// Arbitrum events emitted on L1
	NodeConfirmedTopic = arbitrum_rollup_core.ArbRollupCoreNodeConfirmed{}.Topic()

	// Arbitrum events emitted on L2
	TxToL1Topic              = l2_arbitrum_messenger.L2ArbitrumMessengerTxToL1{}.Topic()
	WithdrawalInitiatedTopic = l2_arbitrum_gateway.L2ArbitrumGatewayWithdrawalInitiated{}.Topic()
	L2ToL1TxTopic            = arbsys.ArbSysL2ToL1Tx{}.Topic()
	DepositFinalizedTopic    = l2_arbitrum_gateway.L2ArbitrumGatewayDepositFinalized{}.Topic()

	// Important addresses on L2
	// These are precompiles so their addresses will never change
	NodeInterfaceAddress = common.HexToAddress("0x00000000000000000000000000000000000000c8")
	ArbSysAddress        = common.HexToAddress("0x0000000000000000000000000000000000000064")

	// Multipliers to ensure our L1 -> L2 tx goes through
	// These values match the arbitrum SDK
	// TODO: should these be configurable?
	l2BaseFeeMultiplier     = big.NewInt(3)
	submissionFeeMultiplier = big.NewInt(4)

	// liquidityManager event - emitted on both L1 and L2
	LiquidityTransferredTopic = liquiditymanager.LiquidityManagerLiquidityTransferred{}.Topic()

	nodeInterfaceABI = abihelpers.MustParseABI(arb_node_interface.NodeInterfaceMetaData.ABI)
	l1AdapterABI     = abihelpers.MustParseABI(arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterMetaData.ABI)
)

type logKey struct {
	txHash   common.Hash
	logIndex int64
}

func parseLiquidityTransferred(parseFunc func(gethtypes.Log) (*liquiditymanager.LiquidityManagerLiquidityTransferred, error), lgs []logpoller.Log) ([]*liquiditymanager.LiquidityManagerLiquidityTransferred, map[logKey]logpoller.Log, error) {
	transferred := make([]*liquiditymanager.LiquidityManagerLiquidityTransferred, len(lgs))
	toLP := make(map[logKey]logpoller.Log)
	for i, lg := range lgs {
		parsed, err := parseFunc(lg.ToGethLog())
		if err != nil {
			// should never happen
			return nil, nil, fmt.Errorf("parse LiquidityTransferred log: %w", err)
		}
		transferred[i] = parsed
		toLP[logKey{
			txHash:   lg.TxHash,
			logIndex: lg.LogIndex,
		}] = lg
	}
	return transferred, toLP, nil
}

func toHash(selector models.NetworkSelector) common.Hash {
	encoded := hexutil.EncodeUint64(uint64(selector))
	return common.HexToHash(encoded)
}
