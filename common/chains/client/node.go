package client

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const queryTimeout = 10 * time.Second

var errInvalidChainID = errors.New("invalid chain id")

var (
	promPoolRPCNodeVerifies = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_verifies",
		Help: "The total number of chain ID verifications for the given RPC node",
	}, []string{"chainID", "nodeName"})
	promPoolRPCNodeVerifiesFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_verifies_failed",
		Help: "The total number of failed chain ID verifications for the given RPC node",
	}, []string{"chainID", "nodeName"})
	promPoolRPCNodeVerifiesSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_verifies_success",
		Help: "The total number of successful chain ID verifications for the given RPC node",
	}, []string{"chainID", "nodeName"})
)

type rawclient struct {
	uri url.URL
}

type Node[
	CHAIN_ID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK_HASH types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	EVENT any,
	EVENT_OPS any, // event filter query options
	TX_RECEIPT any,
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
] interface {
	RPCClient[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]
	// State returns NodeState
	State() NodeState
	// StateAndLatest returns NodeState with the latest received block number & total difficulty.
	StateAndLatest() (NodeState, int64, *utils.Big)
	// Name is a unique identifier for this node.
	Name() string
	String() string

	Order() int32
	Start(context.Context) error
	Close() error
}

type node[
	CHAIN_ID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK_HASH types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	EVENT any,
	EVENT_OPS any, // event filter query options
	TX_RECEIPT any,
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
] struct {
	utils.StartStopOnce
	lfcLog              logger.Logger
	rpcLog              logger.Logger
	name                string
	id                  int32
	chainID             CHAIN_ID
	nodePoolCfg         types.NodePool
	noNewHeadsThreshold time.Duration
	order               int32

	ws   rawclient
	http *rawclient

	rpcClient ChainRPCClient[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]

	stateMu sync.RWMutex // protects state* fields
	state   NodeState
	// Each node is tracking the last received head number and total difficulty
	stateLatestBlockNumber     int64
	stateLatestTotalDifficulty *utils.Big

	// chStopInFlight can be closed to immediately cancel all in-flight requests on
	// this node. Closing and replacing should be serialized through
	// stateMu since it can happen on state transitions as well as node Close.
	chStopInFlight chan struct{}
	// nodeCtx is the node lifetime's context
	nodeCtx context.Context
	// cancelNodeCtx cancels nodeCtx when stopping the node
	cancelNodeCtx context.CancelFunc
	// wg waits for subsidiary goroutines
	wg sync.WaitGroup

	// nLiveNodes is a passed in function that allows this node to:
	//  1. see how many live nodes there are in total, so we can prevent the last alive node in a pool from being
	//  moved to out-of-sync state. It is better to have one out-of-sync node than no nodes at all.
	//  2. compare against the highest head (by number or difficulty) to ensure we don't fall behind too far.
	nLiveNodes func() (count int, blockNumber int64, totalDifficulty *utils.Big)
}

func NewNode[
	CHAIN_ID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK_HASH types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	EVENT any,
	EVENT_OPS any, // event filter query options
	TX_RECEIPT any,
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
](
	nodeCfg types.NodePool,
	noNewHeadsThreshold time.Duration,
	lggr logger.Logger,
	wsuri url.URL,
	httpuri *url.URL,
	name string,
	id int32,
	chainID CHAIN_ID,
	nodeOrder int32,
	rpcClient ChainRPCClient[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB],
) Node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB] {
	n := new(node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB])
	n.name = name
	n.id = id
	n.chainID = chainID
	n.nodePoolCfg = nodeCfg
	n.noNewHeadsThreshold = noNewHeadsThreshold
	n.ws.uri = wsuri
	n.order = nodeOrder
	if httpuri != nil {
		n.http = &rawclient{uri: *httpuri}
	}
	n.chStopInFlight = make(chan struct{})
	n.nodeCtx, n.cancelNodeCtx = context.WithCancel(context.Background())
	lggr = lggr.Named("Node").With(
		"nodeTier", "primary",
		"nodeName", name,
		"node", n.String(),
		"chainID", chainID,
		"nodeOrder", n.order,
	)
	n.lfcLog = lggr.Named("Lifecycle")
	n.rpcLog = lggr.Named("RPC")
	n.stateLatestBlockNumber = -1
	n.rpcClient = rpcClient
	return n
}

// CallContext implementation
func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return n.rpcClient.CallContext(ctx, result, method, args)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) String() string {
	s := fmt.Sprintf("(primary)%s:%s", n.name, n.ws.uri.String())
	if n.http != nil {
		s = s + fmt.Sprintf(":%s", n.http.uri.String())
	}
	return s
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) ChainID() (chainID CHAIN_ID, err error) {
	return n.chainID, nil
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Name() string {
	return n.name
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) makeQueryCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	return makeQueryCtx(ctx, n.getChStopInflight())
}

// getChStopInflight provides a convenience helper that mutex wraps a
// read to the chStopInFlight
func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) getChStopInflight() chan struct{} {
	n.stateMu.RLock()
	defer n.stateMu.RUnlock()
	return n.chStopInFlight
}

// makeQueryCtx returns a context that cancels if:
// 1. Passed in ctx cancels
// 2. Passed in channel is closed
// 3. Default timeout is reached (queryTimeout)
func makeQueryCtx(ctx context.Context, ch utils.StopChan) (context.Context, context.CancelFunc) {
	var chCancel, timeoutCancel context.CancelFunc
	ctx, chCancel = ch.Ctx(ctx)
	ctx, timeoutCancel = context.WithTimeout(ctx, queryTimeout)
	cancel := func() {
		chCancel()
		timeoutCancel()
	}
	return ctx, cancel
}

// Start dials and verifies the node
// Should only be called once in a node's lifecycle
// Return value is necessary to conform to interface but this will never
// actually return an error.
func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Start(startCtx context.Context) error {
	return n.StartOnce(n.name, func() error {
		n.start(startCtx)
		return nil
	})
}

// start initially dials the node and verifies chain ID
// This spins off lifecycle goroutines.
// Not thread-safe.
// Node lifecycle is synchronous: only one goroutine should be running at a
// time.
func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) start(startCtx context.Context) {
	if n.state != NodeStateUndialed {
		panic(fmt.Sprintf("cannot dial node with state %v", n.state))
	}

	dialCtx, dialCancel := n.makeQueryCtx(startCtx)
	defer dialCancel()
	if err := n.rpcClient.Dial(dialCtx); err != nil {
		n.lfcLog.Errorw("Dial failed: Node is unreachable", "err", err)
		n.declareUnreachable()
		return
	}
	n.setState(NodeStateDialed)

	verifyCtx, verifyCancel := n.makeQueryCtx(startCtx)
	defer verifyCancel()
	if err := n.verify(verifyCtx); errors.Is(err, errInvalidChainID) {
		n.lfcLog.Errorw("Verify failed: Node has the wrong chain ID", "err", err)
		n.declareInvalidChainID()
		return
	} else if err != nil {
		n.lfcLog.Errorw(fmt.Sprintf("Verify failed: %v", err), "err", err)
		n.declareUnreachable()
		return
	}

	n.declareAlive()
}

// verify checks that all connections to eth nodes match the given chain ID
// Not thread-safe
// Pure verify: does not mutate node "state" field.
func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) verify(callerCtx context.Context) (err error) {
	ctx, cancel := n.makeQueryCtx(callerCtx)
	defer cancel()

	promPoolRPCNodeVerifies.WithLabelValues(n.chainID.String(), n.name).Inc()
	promFailed := func() {
		promPoolRPCNodeVerifiesFailed.WithLabelValues(n.chainID.String(), n.name).Inc()
	}

	st := n.State()
	switch st {
	case NodeStateDialed, NodeStateOutOfSync, NodeStateInvalidChainID:
	default:
		panic(fmt.Sprintf("cannot verify node in state %v", st))
	}

	var chainID CHAIN_ID
	if chainID, err = n.rpcClient.ClientChainID(ctx); err != nil {
		promFailed()
		return errors.Wrapf(err, "failed to verify chain ID for node %s", n.name)
	} else if chainID.String() != n.chainID.String() {
		promFailed()
		return errors.Wrapf(
			errInvalidChainID,
			"websocket rpc ChainID doesn't match local chain ID: RPC ID=%s, local ID=%s, node name=%s",
			chainID.String(),
			n.chainID.String(),
			n.name,
		)
	}
	if n.http != nil {
		if chainID, err = n.rpcClient.ClientChainID(ctx); err != nil {
			promFailed()
			return errors.Wrapf(err, "failed to verify chain ID for node %s", n.name)
		} else if chainID.String() != n.chainID.String() {
			promFailed()
			return errors.Wrapf(
				errInvalidChainID,
				"http rpc ChainID doesn't match local chain ID: RPC ID=%s, local ID=%s, node name=%s",
				chainID.String(),
				n.chainID.String(),
				n.name,
			)
		}
	}

	promPoolRPCNodeVerifiesSuccess.WithLabelValues(n.chainID.String(), n.name).Inc()

	return nil
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Close() error {
	return n.rpcClient.Close()
}

// disconnectAll disconnects all clients connected to the node
// WARNING: NOT THREAD-SAFE
// This must be called from within the n.stateMu lock
func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) disconnectAll() {
	n.rpcClient.DisconnectAll()
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Order() int32 {
	return n.order
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) BalanceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (*big.Int, error) {
	return n.rpcClient.BalanceAt(ctx, account, blockNumber)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) BatchCallContext(ctx context.Context, b []any) error {
	return n.rpcClient.BatchCallContext(ctx, b)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) BlockByHash(ctx context.Context, hash BLOCK_HASH) (HEAD, error) {
	return n.rpcClient.BlockByHash(ctx, hash)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) BlockByNumber(ctx context.Context, number *big.Int) (HEAD, error) {
	return n.rpcClient.BlockByNumber(ctx, number)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) CallContract(
	ctx context.Context,
	attempt interface{},
	blockNumber *big.Int,
) (rpcErr []byte, extractErr error) {
	return n.rpcClient.CallContract(ctx, attempt, blockNumber)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) CodeAt(ctx context.Context, account ADDR, blockNumber *big.Int) ([]byte, error) {
	return n.rpcClient.CodeAt(ctx, account, blockNumber)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) ConfiguredChainID() CHAIN_ID {
	return n.rpcClient.ConfiguredChainID()
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) EstimateGas(ctx context.Context, call any) (gas uint64, err error) {
	return n.rpcClient.EstimateGas(ctx, call)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) FilterEvents(ctx context.Context, query EVENT_OPS) ([]EVENT, error) {
	return n.rpcClient.FilterEvents(ctx, query)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
	return n.rpcClient.LatestBlockHeight(ctx)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) LINKBalance(ctx context.Context, accountAddress ADDR, linkAddress ADDR) (*assets.Link, error) {
	return n.rpcClient.LINKBalance(ctx, accountAddress, linkAddress)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) PendingSequenceAt(ctx context.Context, addr ADDR) (SEQ, error) {
	return n.rpcClient.PendingSequenceAt(ctx, addr)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) SendEmptyTransaction(
	ctx context.Context,
	newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt any, err error),
	seq SEQ,
	gasLimit uint32,
	fee FEE,
	fromAddress ADDR,
) (txhash string, err error) {
	return n.rpcClient.SendEmptyTransaction(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) SendTransaction(ctx context.Context, tx *TX) error {
	return n.rpcClient.SendTransaction(ctx, tx)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) SendTransactionReturnCode(
	ctx context.Context,
	tx *TX,
) (SendTxReturnCode, error) {
	return n.rpcClient.SendTransactionReturnCode(ctx, tx)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) SetState(state NodeState) {
	n.state = state
	n.rpcClient.SetState(state)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) SequenceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (SEQ, error) {
	return n.rpcClient.SequenceAt(ctx, account, blockNumber)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) SimulateTransaction(ctx context.Context, tx *TX) error {
	return n.rpcClient.SimulateTransaction(ctx, tx)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (SUB, error) {
	return n.rpcClient.Subscribe(ctx, channel, args)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) TokenBalance(ctx context.Context, account ADDR, tokenAddr ADDR) (*big.Int, error) {
	return n.rpcClient.TokenBalance(ctx, account, tokenAddr)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) TransactionByHash(ctx context.Context, txHash TX_HASH) (*TX, error) {
	return n.rpcClient.TransactionByHash(ctx, txHash)
}

func (n *node[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) TransactionReceipt(ctx context.Context, txHash TX_HASH) (*TX_RECEIPT, error) {
	return n.rpcClient.TransactionReceipt(ctx, txHash)
}
