package client

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const queryTimeout = 10 * time.Second

var errInvalidChainID = errors.New("invalid chain id")

var (
	promPoolRPCNodeVerifies = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_verifies",
		Help: "The total number of chain ID verifications for the given RPC node",
	}, []string{"network", "chainID", "nodeName"})
	promPoolRPCNodeVerifiesFailed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_verifies_failed",
		Help: "The total number of failed chain ID verifications for the given RPC node",
	}, []string{"network", "chainID", "nodeName"})
	promPoolRPCNodeVerifiesSuccess = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pool_rpc_node_verifies_success",
		Help: "The total number of successful chain ID verifications for the given RPC node",
	}, []string{"network", "chainID", "nodeName"})
)

type rawclient struct {
	uri url.URL
}

type Node[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
	RPC_CLIENT NodeClientAPI[CHAIN_ID, BLOCK_HASH, HEAD, SUB],
] interface {
	// State returns NodeState
	State() NodeState
	// StateAndLatest returns NodeState with the latest received block number & total difficulty.
	StateAndLatest() (NodeState, int64, *utils.Big)
	// Name is a unique identifier for this node.
	Name() string
	String() string
	RPCClient() RPC_CLIENT
	ChainID() (CHAIN_ID, error)
	Order() int32
	Start(context.Context) error
	Close() error
}

type node[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
	RPC_CLIENT NodeClientAPI[CHAIN_ID, BLOCK_HASH, HEAD, SUB],
] struct {
	utils.StartStopOnce
	lfcLog              logger.Logger
	name                string
	id                  int32
	chainID             CHAIN_ID
	nodePoolCfg         types.NodePool
	noNewHeadsThreshold time.Duration
	order               int32
	chainFamily         string

	ws   rawclient
	http *rawclient

	rpcClient RPC_CLIENT

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
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
	RPC_CLIENT NodeClientAPI[CHAIN_ID, BLOCK_HASH, HEAD, SUB],
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
	rpcClient RPC_CLIENT,
	chainFamily string,
) Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT] {
	n := new(node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT])
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
	n.stateLatestBlockNumber = -1
	n.rpcClient = rpcClient
	n.chainFamily = chainFamily
	return n
}

func (n *node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) String() string {
	s := fmt.Sprintf("(primary)%s:%s", n.name, n.ws.uri.String())
	if n.http != nil {
		s = s + fmt.Sprintf(":%s", n.http.uri.String())
	}
	return s
}

func (n *node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) ChainID() (chainID CHAIN_ID, err error) {
	return n.chainID, nil
}

func (n *node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Name() string {
	return n.name
}

func (n *node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) RPCClient() RPC_CLIENT {
	return n.rpcClient
}

func (n *node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) makeQueryCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	return makeQueryCtx(ctx, n.getChStopInflight())
}

// getChStopInflight provides a convenience helper that mutex wraps a
// read to the chStopInFlight
func (n *node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) getChStopInflight() chan struct{} {
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
func (n *node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Start(startCtx context.Context) error {
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
func (n *node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) start(startCtx context.Context) {
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
func (n *node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) verify(callerCtx context.Context) (err error) {
	ctx, cancel := n.makeQueryCtx(callerCtx)
	defer cancel()

	promPoolRPCNodeVerifies.WithLabelValues(n.chainFamily, n.chainID.String(), n.name).Inc()
	promFailed := func() {
		promPoolRPCNodeVerifiesFailed.WithLabelValues(n.chainFamily, n.chainID.String(), n.name).Inc()
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

	promPoolRPCNodeVerifiesSuccess.WithLabelValues(n.chainFamily, n.chainID.String(), n.name).Inc()

	return nil
}

func (n *node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Close() error {
	return n.rpcClient.Close()
}

// disconnectAll disconnects all clients connected to the node
// WARNING: NOT THREAD-SAFE
// This must be called from within the n.stateMu lock
func (n *node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) disconnectAll() {
	n.rpcClient.DisconnectAll()
}

func (n *node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Order() int32 {
	return n.order
}

func (n *node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) SetState(state NodeState) {
	n.state = state
	n.rpcClient.SetState(state)
}

func (n *node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) ChainFamily() string {
	return n.chainFamily
}
