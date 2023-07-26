package client

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// TODO: remove this maybe?
const queryTimeout = 10 * time.Second

var errInvalidChainID = errors.New("invalid chain id")

var (
	promEVMPoolRPCNodeDials = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_dials_total",
		Help: "The total number of dials for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeDialsFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_dials_failed",
		Help: "The total number of failed dials for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeDialsSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_dials_success",
		Help: "The total number of successful dials for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeVerifies = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_verifies",
		Help: "The total number of chain ID verifications for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeVerifiesFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_verifies_failed",
		Help: "The total number of failed chain ID verifications for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeVerifiesSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_verifies_success",
		Help: "The total number of successful chain ID verifications for the given RPC node",
	}, []string{"evmChainID", "nodeName"})

	promEVMPoolRPCNodeCalls = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_calls_total",
		Help: "The approximate total number of RPC calls for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeCallsFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_calls_failed",
		Help: "The approximate total number of failed RPC calls for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCNodeCallsSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "evm_pool_rpc_node_calls_success",
		Help: "The approximate total number of successful RPC calls for the given RPC node",
	}, []string{"evmChainID", "nodeName"})
	promEVMPoolRPCCallTiming = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "evm_pool_rpc_node_rpc_call_time",
		Help: "The duration of an RPC call in nanoseconds",
		Buckets: []float64{
			float64(50 * time.Millisecond),
			float64(100 * time.Millisecond),
			float64(200 * time.Millisecond),
			float64(500 * time.Millisecond),
			float64(1 * time.Second),
			float64(2 * time.Second),
			float64(4 * time.Second),
			float64(8 * time.Second),
		},
	}, []string{"evmChainID", "nodeName", "rpcHost", "isSendOnly", "success", "rpcCallName"})
)

type rawclient struct {
	uri url.URL
}

type Node[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	EVENT any,
	EVENTOPS any, // event filter query options
	TXRECEIPT txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
] interface {
	RPCClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]
	// State returns NodeState
	State() NodeState
	// StateAndLatest returns NodeState with the latest received block number & total difficulty.
	StateAndLatest() (state NodeState, blockNum int64, totalDifficulty *utils.Big)
	// Name is a unique identifier for this node.
	Name() string

	String() string

	Order() int32
}

type node[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	EVENT any,
	EVENTOPS any, // event filter query options
	TXRECEIPT txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
] struct {
	utils.StartStopOnce
	lfcLog              logger.Logger
	rpcLog              logger.Logger
	name                string
	id                  int32
	chainID             CHAINID
	nodePoolCfg         types.NodePool
	noNewHeadsThreshold time.Duration
	order               int32

	ws   rawclient
	http *rawclient

	rpcClient ChainRPCClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]

	stateMu sync.RWMutex // protects state* fields
	state   NodeState
	// Each node is tracking the last received head number and total difficulty
	stateLatestBlockNumber     int64
	stateLatestTotalDifficulty *utils.Big

	// Need to track subscriptions because closing the RPC does not (always?)
	// close the underlying subscription
	subs []types.Subscription

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
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	EVENT any,
	EVENTOPS any,
	TXRECEIPT txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
](
	nodeCfg types.NodePool,
	noNewHeadsThreshold time.Duration,
	lggr logger.Logger,
	wsuri url.URL,
	httpuri *url.URL,
	name string,
	id int32,
	chainID CHAINID,
	nodeOrder int32,
	rpcClient ChainRPCClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE],
) *Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE] {
	n := new(node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE])
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
func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	// ctx, cancel, ws, http, err := n.makeLiveQueryCtx(ctx)
	ctx, cancel, _, _, err := n.makeLiveQueryCtx(ctx)

	if err != nil {
		return err
	}
	defer cancel()
	lggr := n.newRqLggr(switching(n)).With(
		"method", method,
		"args", args,
	)

	lggr.Debug("RPC call: evmclient.Client#CallContext")
	start := time.Now()

	// TODO: Make this differentiate WS and HTTP call context

	// if http != nil {
	// 	err = n.wrapHTTP(http.rpc.CallContext(ctx, result, method, args...))
	// } else {
	// 	err = n.wrapWS(ws.rpc.CallContext(ctx, result, method, args...))
	// }
	n.rpcClient.RPCCallContext(ctx, result, method, args...)
	duration := time.Since(start)

	n.logResult(lggr, err, duration, n.getRPCDomain(), "CallContext")

	return err
}

// makeLiveQueryCtx wraps makeQueryCtx but returns error if node is not NodeStateAlive.
func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) makeLiveQueryCtx(parentCtx context.Context) (ctx context.Context, cancel context.CancelFunc, ws rawclient, http *rawclient, err error) {
	// Need to wrap in mutex because state transition can cancel and replace the
	// context
	n.stateMu.RLock()
	if n.state != NodeStateAlive {
		err = errors.Errorf("cannot execute RPC call on node with state: %s", n.state)
		n.stateMu.RUnlock()
		return
	}
	cancelCh := n.chStopInFlight
	ws = n.ws
	if n.http != nil {
		cp := *n.http
		http = &cp
	}
	n.stateMu.RUnlock()
	ctx, cancel = makeQueryCtx(parentCtx, cancelCh)
	return
}

func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) wrapHTTP(err error) error {
	err = wrap(err, fmt.Sprintf("primary http (%s)", n.http.uri.Redacted()))
	if err != nil {
		n.rpcLog.Debugw("Call failed", "err", err)
	} else {
		n.rpcLog.Trace("Call succeeded")
	}
	return err
}

func wrap(err error, tp string) error {
	if err == nil {
		return nil
	}
	if errors.Cause(err).Error() == "context deadline exceeded" {
		err = errors.Wrap(err, "remote eth node timed out")
	}
	return errors.Wrapf(err, "%s call failed", tp)
}

func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) String() string {
	s := fmt.Sprintf("(primary)%s:%s", n.name, n.ws.uri)
	if n.http != nil {
		s = s + fmt.Sprintf(":%s", n.http.uri)
	}
	return s
}

func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) ChainID() (chainID CHAINID, err error) {
	return n.chainID, nil
}

func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) Name() string {
	return n.name
}

func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) makeQueryCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	return makeQueryCtx(ctx, n.getChStopInflight())
}

// getChStopInflight provides a convenience helper that mutex wraps a
// read to the chStopInFlight
func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) getChStopInflight() chan struct{} {
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

// disconnectAll disconnects all clients connected to the node
// WARNING: NOT THREAD-SAFE
// This must be called from within the n.stateMu lock
func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) disconnectAll() {
	if n.rpcClient != nil {
		n.rpcClient.Close()
	}
	n.cancelInflightRequests()
	n.unsubscribeAll()
}

// cancelInflightRequests closes and replaces the chStopInFlight
// WARNING: NOT THREAD-SAFE
// This must be called from within the n.stateMu lock
func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) cancelInflightRequests() {
	close(n.chStopInFlight)
	n.chStopInFlight = make(chan struct{})
}

// unsubscribeAll unsubscribes all subscriptions
// WARNING: NOT THREAD-SAFE
// This must be called from within the n.stateMu lock
func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) unsubscribeAll() {
	for _, sub := range n.subs {
		sub.Unsubscribe()
	}
	n.subs = nil
}

func switching[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	EVENT any,
	EVENTOPS any,
	TXRECEIPT txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
](n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) string {
	if n.http != nil {
		return "http"
	}
	return "websocket"
}

func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) wrapWS(err error) error {
	err = wrap(err, fmt.Sprintf("primary websocket (%s)", n.ws.uri.Redacted()))
	return err
}

// newRqLggr generates a new logger with a unique request ID
func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) newRqLggr(mode string) logger.Logger {
	return n.rpcLog.With(
		"requestID", uuid.New(),
		"mode", mode,
	)
}

func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) logResult(
	lggr logger.Logger,
	err error,
	callDuration time.Duration,
	rpcDomain,
	callName string,
	results ...interface{},
) {
	lggr = lggr.With("duration", callDuration, "rpcDomain", rpcDomain, "callName", callName)
	promEVMPoolRPCNodeCalls.WithLabelValues(n.chainID.String(), n.name).Inc()
	if err == nil {
		promEVMPoolRPCNodeCallsSuccess.WithLabelValues(n.chainID.String(), n.name).Inc()
		lggr.Tracew(
			fmt.Sprintf("evmclient.Client#%s RPC call success", callName),
			results...,
		)
	} else {
		promEVMPoolRPCNodeCallsFailed.WithLabelValues(n.chainID.String(), n.name).Inc()
		lggr.Debugw(
			fmt.Sprintf("evmclient.Client#%s RPC call failure", callName),
			append(results, "err", err)...,
		)
	}
	promEVMPoolRPCCallTiming.
		WithLabelValues(
			n.chainID.String(),             // chain id
			n.name,                         // node name
			rpcDomain,                      // rpc domain
			"false",                        // is send only
			strconv.FormatBool(err == nil), // is successful
			callName,                       // rpc call name
		).
		Observe(float64(callDuration))
}

func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) getRPCDomain() string {
	if n.http != nil {
		return n.http.uri.Host
	}
	return n.ws.uri.Host
}

func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) Order() int32 {
	return n.order
}

func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BalanceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (*big.Int, error) {
	return n.rpcClient.BalanceAt(ctx, account, blockNumber)
}

func (n *node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE]) BatchCallContext(ctx context.Context, b []any) error {
	return n.rpcClient.BatchCallContext(ctx, b)
}
