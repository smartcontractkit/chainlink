package arb

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/liquiditymanager/multienv"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	utilsbig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arb_node_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbitrum_rollup_core"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/arbsys"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/l2_arbitrum_gateway"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/l2_arbitrum_messenger"
)

var (
	// Events emitted on L2
	TxToL1Topic              = l2_arbitrum_messenger.L2ArbitrumMessengerTxToL1{}.Topic()
	WithdrawalInitiatedTopic = l2_arbitrum_gateway.L2ArbitrumGatewayWithdrawalInitiated{}.Topic()
	L2ToL1TxTopic            = arbsys.ArbSysL2ToL1Tx{}.Topic()

	// Important addresses on L2
	NodeInterfaceAddress = common.HexToAddress("0x00000000000000000000000000000000000000c8")
	ArbSysAddress        = common.HexToAddress("0x0000000000000000000000000000000000000064")

	// Events emitted on L1
	NodeConfirmedTopic = arbitrum_rollup_core.ArbRollupCoreNodeConfirmed{}.Topic()
)

// nolint
// function executeTransaction(
//
//	  bytes32[] calldata proof,
//	  uint256 index,
//	  address l2Sender,
//	  address to,
//	  uint256 l2Block,
//	  uint256 l1Block,
//	  uint256 l2Timestamp,
//	  uint256 value,
//	  bytes calldata data
//	) external;
//
// Arg 0: proof. This takes multiple steps:
// 1. Get the latest NodeConfirmed event on L1, which indicates the latest node that was confirmed by the rollup.
// 2. Call eth_getBlockByHash on L2 specifying the L2 block hash in the NodeConfirmed event.
// 3. Get the `sendCount` field from the response.
// 4. Get the `l2ToL1Id` field from the `WithdrawalInitiated` log from the L2 withdrawal tx.
// 5. Call `constructOutboxProof` on the L2 node interface contract with the `sendCount` as the first argument and `l2ToL1Id` as the second argument.
// Arg 1: index. Fetch the index from the TxToL1 log in the L2 tx.
// Arg 2: l2Sender. Fetch the source of the WithdrawalInitiated log in the L2 tx.
// Arg 3: to. Fetch the `to` field of the WithdrawalInitiated log in the L2 tx.
// Arg 4: l1Block. Fetch the `l1BlockNumber` field of the JSON-RPC response to eth_getTransactionReceipt
// passing in the L2 tx hash as the param.
// Arg 5: l2Block. This is the l2 block number in which the withdrawal tx was included.
// Arg 6: l2Timestamp. Get the `timestamp` field from the L2ToL1Tx event emitted by ArbSys (0x64).
// Arg 7: value. Fetch the `value` field from the WithdrawalInitiated log in the L2 tx.
// Arg 8: data. Fetch the `data` field from the TxToL1 log in the L2 tx.
func FinalizeL1(
	env multienv.Env,
	l1ChainID uint64,
	l2ChainID uint64,
	l1BridgeAdapterAddress common.Address,
	l2TxHash common.Hash,
) {
	// get the logs we care about from the L2 tx:
	// 1. L2ToL1Tx
	// 2. WithdrawalInitiated
	// 3. TxToL1
	l2Client := env.Clients[l2ChainID]
	receipt, err := l2Client.TransactionReceipt(context.Background(), l2TxHash)
	helpers.PanicErr(err)
	var (
		l2ToL1TxLog, withdrawalInitiatedLog, txToL1Log *types.Log
	)
	for _, lg := range receipt.Logs {
		if lg.Topics[0] == L2ToL1TxTopic {
			l2ToL1TxLog = lg
		} else if lg.Topics[0] == WithdrawalInitiatedTopic {
			withdrawalInitiatedLog = lg
		} else if lg.Topics[0] == TxToL1Topic {
			txToL1Log = lg
		}
	}
	if l2ToL1TxLog == nil || withdrawalInitiatedLog == nil || txToL1Log == nil {
		helpers.PanicErr(fmt.Errorf("missing logs in L2 tx %s", l2TxHash.String()))
		return
	}
	arbSys, err := arbsys.NewArbSys(ArbSysAddress, env.Clients[l2ChainID])
	helpers.PanicErr(err)
	// parse logs
	l2ToL1Tx, err := arbSys.ParseL2ToL1Tx(*l2ToL1TxLog)
	helpers.PanicErr(err)
	withdrawalInitiated := parseWithdrawalInitiated(env, l2ChainID, withdrawalInitiatedLog)
	txToL1 := parseTxToL1(env, l2ChainID, txToL1Log)
	// get the proof
	arg0Proof := getProof(env, l1ChainID, l2ChainID, withdrawalInitiated.L2ToL1Id)
	// argument 1: index
	arg1Index := withdrawalInitiated.L2ToL1Id
	// argument 2: l2Sender
	arg2L2Sender := withdrawalInitiatedLog.Address
	// argument 3: to
	arg3To := txToL1.To
	// argument 4: l1Block
	arg4L1Block := getL1BlockFromRPC(env, l2ChainID, l2TxHash)
	// argument 5: l2Block
	arg5L2Block := receipt.BlockNumber
	// argument 6: l2Timestamp
	arg6L2Timestamp := l2ToL1Tx.Timestamp
	// argument 7: value
	arg7Value := withdrawalInitiated.Amount
	// argument 8: data
	arg8Data := txToL1.Data

	// print the arguments for the executeTransaction call
	fmt.Println("proof:", encodeProofToHex(arg0Proof), "\n",
		"index:", arg1Index, "\n",
		"l2Sender:", arg2L2Sender, "\n",
		"to:", arg3To, "\n",
		"l1Block:", arg4L1Block, "\n",
		"l2Block:", arg5L2Block, "\n",
		"l2Timestamp:", arg6L2Timestamp, "\n",
		"value:", arg7Value, "\n",
		"data:", hexutil.Encode(arg8Data))

	// execute the transaction
	fmt.Println("executing transaction on the bridge adapter with the above data")

	adapter, err := arbitrum_l1_bridge_adapter.NewArbitrumL1BridgeAdapter(l1BridgeAdapterAddress, env.Clients[l1ChainID])
	helpers.PanicErr(err)

	adapterABI, err := arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterMetaData.GetAbi()
	helpers.PanicErr(err)

	finalizationPayload, err := adapterABI.Pack("exposeArbitrumFinalizationPayload", arbitrum_l1_bridge_adapter.ArbitrumL1BridgeAdapterArbitrumFinalizationPayload{
		Proof:       arg0Proof,
		Index:       arg1Index,
		L2Sender:    arg2L2Sender,
		To:          arg3To,
		L1Block:     arg4L1Block,
		L2Block:     arg5L2Block,
		L2Timestamp: arg6L2Timestamp,
		Value:       arg7Value,
		Data:        arg8Data,
	})
	helpers.PanicErr(err)
	// trim first four bytes (function signature)
	finalizationPayload = finalizationPayload[4:]

	tx, err := adapter.FinalizeWithdrawERC20(env.Transactors[l1ChainID], common.HexToAddress("0x0"), common.HexToAddress("0x0"), finalizationPayload)
	helpers.PanicErr(err)
	receipt = helpers.ConfirmTXMined(context.Background(), env.Clients[l1ChainID], tx, int64(l1ChainID))
	fmt.Println("transaction mined:", receipt.TxHash.String(), "status:", receiptStatusToString(receipt.Status))
}

func encodeProofToHex(proof [][32]byte) []string {
	proofHex := make([]string, len(proof))
	for i, step := range proof {
		proofHex[i] = hexutil.Encode(step[:])
	}
	return proofHex
}

// Arg 0: proof. This takes multiple steps:
// 1. Get the latest NodeConfirmed event on L1, which indicates the latest node that was confirmed by the rollup.
// 2. Call eth_getBlockByHash on L2 specifying the L2 block hash in the NodeConfirmed event.
// 3. Get the `sendCount` field from the response.
// 4. Get the `l2ToL1Id` field from the `WithdrawalInitiated` log from the L2 withdrawal tx.
// 5. Call `constructOutboxProof` on the L2 node interface contract with the `sendCount` as the first argument and `l2ToL1Id` as the second argument.
// Note that this may fail, specifically, ConstructOutboxProof. If it does, that means that the L2 batch that has the bridge tx
// has not yet been committed to L1, and that we should wait and try again.
func getProof(env multienv.Env, l1ChainID, l2ChainID uint64, l2ToL1Id *big.Int) [][32]byte {
	l1Client := env.Clients[l1ChainID]
	latestHeader, err := l1Client.HeaderByNumber(context.Background(), nil)
	helpers.PanicErr(err)
	// start four hours back in terms of blocks
	// 12 seconds per block => 5 * 240 = 1200 blocks
	startBlock := big.NewInt(0).Sub(latestHeader.Number, big.NewInt(1200))
	lgs, err := l1Client.FilterLogs(context.Background(), ethereum.FilterQuery{
		Addresses: []common.Address{ArbitrumContracts[l1ChainID]["Rollup"]},
		Topics: [][]common.Hash{{
			NodeConfirmedTopic,
		}},
		FromBlock: startBlock,
	})
	helpers.PanicErr(err)
	var latestNodeConfirmed *types.Log
	for _, lg := range lgs {
		lg := lg // exportloopref
		if latestNodeConfirmed == nil || lg.BlockNumber > latestNodeConfirmed.BlockNumber {
			latestNodeConfirmed = &lg
		}
	}
	if latestNodeConfirmed == nil {
		helpers.PanicErr(fmt.Errorf("no node confirmed event found"))
	}
	// parse latest nodeconfirmed event
	nodeConfirmed := parseNodeConfirmed(env, l1ChainID, latestNodeConfirmed)
	fmt.Println("latest node confirmed:", nodeConfirmedToString(nodeConfirmed))
	type Response struct {
		SendCount *utilsbig.Big `json:"sendCount"`
	}
	response := new(Response)
	l2Rpc := env.JRPCs[l2ChainID]
	err = l2Rpc.Call(response, "eth_getBlockByHash", hexutil.Encode(nodeConfirmed.BlockHash[:]), false)
	helpers.PanicErr(err)
	nodeInterface, err := arb_node_interface.NewNodeInterface(NodeInterfaceAddress, env.Clients[l2ChainID])
	helpers.PanicErr(err)
	fmt.Println("send count:", response.SendCount, "l2 to l1 id:", l2ToL1Id)
	outboxProof, err := nodeInterface.ConstructOutboxProof(nil, response.SendCount.ToInt().Uint64(), l2ToL1Id.Uint64())
	helpers.PanicErr(err)
	return outboxProof.Proof
}

func nodeConfirmedToString(nodeConfirmed *arbitrum_rollup_core.ArbRollupCoreNodeConfirmed) string {
	return fmt.Sprintf("NodeConfirmed{BlockHash: %s, SendRoot: %s, NodeNum: %d}",
		hexutil.Encode(nodeConfirmed.BlockHash[:]),
		hexutil.Encode(nodeConfirmed.SendRoot[:]),
		nodeConfirmed.NodeNum,
	)
}

func getL1BlockFromRPC(env multienv.Env, l2ChainID uint64, l2TxHash common.Hash) *big.Int {
	l1Rpc := env.JRPCs[l2ChainID]
	type Response struct {
		L1BlockNumber hexutil.Big `json:"l1BlockNumber"`
	}
	response := new(Response)
	err := l1Rpc.Call(response, "eth_getTransactionReceipt", l2TxHash.String())
	helpers.PanicErr(err)
	return response.L1BlockNumber.ToInt()
}

func parseWithdrawalInitiated(env multienv.Env, l2ChainID uint64, lg *types.Log) *l2_arbitrum_gateway.L2ArbitrumGatewayWithdrawalInitiated {
	// Address provided doesn't matter, we're just going to parse the log
	l2ArbGateway, err := l2_arbitrum_gateway.NewL2ArbitrumGateway(common.HexToAddress("0x0"), env.Clients[l2ChainID])
	helpers.PanicErr(err)
	parsed, err := l2ArbGateway.ParseWithdrawalInitiated(*lg)
	helpers.PanicErr(err)
	return parsed
}

func parseTxToL1(env multienv.Env, l2ChainID uint64, lg *types.Log) *l2_arbitrum_messenger.L2ArbitrumMessengerTxToL1 {
	// Address provided doesn't matter, we're just going to parse the log
	l2ArbMessenger, err := l2_arbitrum_messenger.NewL2ArbitrumMessenger(common.HexToAddress("0x0"), env.Clients[l2ChainID])
	helpers.PanicErr(err)
	parsed, err := l2ArbMessenger.ParseTxToL1(*lg)
	helpers.PanicErr(err)
	return parsed
}

func parseNodeConfirmed(env multienv.Env, l1ChainID uint64, lg *types.Log) *arbitrum_rollup_core.ArbRollupCoreNodeConfirmed {
	// Address provided doesn't matter, we're just going to parse the log
	rollupCore, err := arbitrum_rollup_core.NewArbRollupCore(common.HexToAddress("0x0"), env.Clients[l1ChainID])
	helpers.PanicErr(err)
	parsed, err := rollupCore.ParseNodeConfirmed(*lg)
	helpers.PanicErr(err)
	return parsed
}

func receiptStatusToString(status uint64) string {
	if status == 0 {
		return "failed"
	}
	return "successful"
}
