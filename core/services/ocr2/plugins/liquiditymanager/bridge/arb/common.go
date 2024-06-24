package arb

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arb_node_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_rollup_core"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbsys"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/l2_arbitrum_gateway"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/l2_arbitrum_messenger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
)

const (
	// DepositFinalizedToAddressTopicIndex is the index of the topic in the DepositFinalized event
	// that contains the "to" address.
	DepositFinalizedToAddressTopicIndex = 3
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

	nodeInterfaceABI = abihelpers.MustParseABI(arb_node_interface.NodeInterfaceMetaData.ABI)
	l1AdapterABI     = abihelpers.MustParseABI(arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterMetaData.ABI)
)
