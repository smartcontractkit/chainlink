package client

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	commonconfig "github.com/smartcontractkit/chainlink/v2/common/config"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

const QueryTimeout = 10 * time.Second

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

type NodeConfig interface {
	PollFailureThreshold() uint32
	PollInterval() time.Duration
	SelectionMode() string
	SyncThreshold() uint32
	NodeIsSyncingEnabled() bool
	FinalizedBlockPollInterval() time.Duration
	Errors() config.ClientErrors
}

type ChainConfig interface {
	NodeNoNewHeadsThreshold() time.Duration
	FinalityDepth() uint32
	FinalityTagEnabled() bool
	ChainType() commonconfig.ChainType
}

// ChainInfo - represents RPC’s view of the chain
type ChainInfo struct {
	// BlockNumber - block number of the most recent block observed by the Node
	BlockNumber int64
	// BlockDifficulty - difficulty of the most recent block observed by the Node
	BlockDifficulty *big.Int
	// LatestFinalizedBlock - block number of the most recently finalized block
	LatestFinalizedBlock int64
}

//go:generate mockery --quiet --name Node --structname mockNode --filename "mock_node_test.go" --inpackage --case=underscore
type Node[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC_CLIENT any,
] interface {
	// State returns health state of the underlying RPC
	State() NodeState
	// StateAndLatest returns health state with the latest received block number & total difficulty.
	StateAndLatest() (NodeState, ChainInfo)
	// Name is a unique identifier for this node.
	Name() string
	// String - returns string representation of the node, useful for debugging (name + URLS used to connect to the RPC)
	String() string
	// RPC - returns the underlying RPC_CLIENT
	RPC() RPC_CLIENT
	// UnsubscribeAll - terminates all client subscriptions. Called by MultiNode to trigger clients to resubscribe to
	// new best RPC
	UnsubscribeAll()
	ConfiguredChainID() CHAIN_ID
	// Order - returns priority order configured for the RPC
	Order() int32
	// Start - starts health checks
	Start(context.Context) error
	Close() error
}

type node[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC_CLIENT RPCClient[CHAIN_ID, HEAD],
] struct {
	services.StateMachine
	lfcLog      logger.Logger
	name        string
	id          int32
	chainID     CHAIN_ID
	nodePoolCfg NodeConfig
	chainCfg    ChainConfig
	order       int32
	chainFamily string

	ws   url.URL
	http *url.URL

	rpc RPC_CLIENT

	stateMu sync.RWMutex // protects state* fields
	state   NodeState
	// Each node is tracking the last received head number and total difficulty
	stateLatestBlockNumber          int64
	stateLatestTotalDifficulty      *big.Int
	stateLatestFinalizedBlockNumber int64

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
	nLiveNodes   func() (count int, blockNumber int64, totalDifficulty *big.Int)
	aliveLoopSub types.Subscription
}

func NewNode[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC_CLIENT RPCClient[CHAIN_ID, HEAD],
](
	nodeCfg NodeConfig,
	chainCfg ChainConfig,
	lggr logger.Logger,
	wsuri url.URL,
	httpuri *url.URL,
	name string,
	id int32,
	chainID CHAIN_ID,
	nodeOrder int32,
	rpc RPC_CLIENT,
	chainFamily string,
) Node[CHAIN_ID, HEAD, RPC_CLIENT] {
	n := new(node[CHAIN_ID, HEAD, RPC_CLIENT])
	n.name = name
	n.id = id
	n.chainID = chainID
	n.nodePoolCfg = nodeCfg
	n.chainCfg = chainCfg
	n.ws = wsuri
	n.order = nodeOrder
	if httpuri != nil {
		n.http = httpuri
	}
	n.nodeCtx, n.cancelNodeCtx = context.WithCancel(context.Background())
	lggr = logger.Named(lggr, "Node")
	lggr = logger.With(lggr,
		"nodeTier", Primary.String(),
		"nodeName", name,
		"node", n.String(),
		"chainID", chainID,
		"nodeOrder", n.order,
	)
	n.lfcLog = logger.Named(lggr, "Lifecycle")
	n.stateLatestBlockNumber = -1
	n.rpc = rpc
	n.chainFamily = chainFamily
	n.aliveLoopSub = nil
	return n
}

func (n *node[CHAIN_ID, HEAD, RPC_CLIENT]) String() string {
	s := fmt.Sprintf("(%s)%s:%s", Primary.String(), n.name, n.ws.String())
	if n.http != nil {
		s = s + fmt.Sprintf(":%s", n.http.String())
	}
	return s
}

func (n *node[CHAIN_ID, HEAD, RPC_CLIENT]) ConfiguredChainID() (chainID CHAIN_ID) {
	return n.chainID
}

func (n *node[CHAIN_ID, HEAD, RPC_CLIENT]) Name() string {
	return n.name
}

func (n *node[CHAIN_ID, HEAD, RPC_CLIENT]) RPC() RPC_CLIENT {
	return n.rpc
}

func (n *node[CHAIN_ID, HEAD, RPC_CLIENT]) UnsubscribeAll() {
	n.rpc.UnsubscribeAllExcept(n.aliveLoopSub)
}

func (n *node[CHAIN_ID, HEAD, RPC_CLIENT]) Close() error {
	return n.StopOnce(n.name, n.close)
}

func (n *node[CHAIN_ID, HEAD, RPC_CLIENT]) close() error {
	defer func() {
		n.wg.Wait()
		n.rpc.Close()
	}()

	n.stateMu.Lock()
	defer n.stateMu.Unlock()

	n.cancelNodeCtx()
	n.state = nodeStateClosed
	return nil
}

// Start dials and verifies the node
// Should only be called once in a node's lifecycle
// Return value is necessary to conform to interface but this will never
// actually return an error.
func (n *node[CHAIN_ID, HEAD, RPC_CLIENT]) Start(startCtx context.Context) error {
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
func (n *node[CHAIN_ID, HEAD, RPC_CLIENT]) start(startCtx context.Context) {
	if n.state != nodeStateUndialed {
		panic(fmt.Sprintf("cannot dial node with state %v", n.state))
	}

	if err := n.rpc.Dial(startCtx); err != nil {
		n.lfcLog.Errorw("Dial failed: Node is unreachable", "err", err)
		n.declareUnreachable()
		return
	}
	n.setState(nodeStateDialed)

	state := n.verifyConn(startCtx, n.lfcLog)
	n.declareState(state)
}

// verifyChainID checks that connection to the node matches the given chain ID
// Not thread-safe
// Pure verifyChainID: does not mutate node "state" field.
func (n *node[CHAIN_ID, HEAD, RPC_CLIENT]) verifyChainID(callerCtx context.Context, lggr logger.Logger) NodeState {
	promPoolRPCNodeVerifies.WithLabelValues(n.chainFamily, n.chainID.String(), n.name).Inc()
	promFailed := func() {
		promPoolRPCNodeVerifiesFailed.WithLabelValues(n.chainFamily, n.chainID.String(), n.name).Inc()
	}

	st := n.State()
	switch st {
	case nodeStateClosed:
		// The node is already closed, and any subsequent transition is invalid.
		// To make spotting such transitions a bit easier, return the invalid node state.
		return nodeStateLen
	case nodeStateDialed, nodeStateOutOfSync, nodeStateInvalidChainID, nodeStateSyncing:
	default:
		panic(fmt.Sprintf("cannot verify node in state %v", st))
	}

	var chainID CHAIN_ID
	var err error
	if chainID, err = n.rpc.ChainID(callerCtx); err != nil {
		promFailed()
		lggr.Errorw("Failed to verify chain ID for node", "err", err, "NodeState", n.State())
		return nodeStateUnreachable
	} else if chainID.String() != n.chainID.String() {
		promFailed()
		err = fmt.Errorf(
			"rpc ChainID doesn't match local chain ID: RPC ID=%s, local ID=%s, node name=%s: %w",
			chainID.String(),
			n.chainID.String(),
			n.name,
			errInvalidChainID,
		)
		lggr.Errorw("Failed to verify RPC node; remote endpoint returned the wrong chain ID", "err", err, "NodeState", n.State())
		return nodeStateInvalidChainID
	}

	promPoolRPCNodeVerifiesSuccess.WithLabelValues(n.chainFamily, n.chainID.String(), n.name).Inc()

	return nodeStateAlive
}

// createVerifiedConn - establishes new connection with the RPC and verifies that it's valid: chainID matches, and it's not syncing.
// Returns desired state if one of the verifications fails. Otherwise, returns nodeStateAlive.
func (n *node[CHAIN_ID, HEAD, RPC_CLIENT]) createVerifiedConn(ctx context.Context, lggr logger.Logger) NodeState {
	if err := n.rpc.Dial(ctx); err != nil {
		n.lfcLog.Errorw("Dial failed: Node is unreachable", "err", err, "NodeState", n.State())
		return nodeStateUnreachable
	}

	return n.verifyConn(ctx, lggr)
}

// verifyConn - verifies that current connection is valid: chainID matches, and it's not syncing.
// Returns desired state if one of the verifications fails. Otherwise, returns nodeStateAlive.
func (n *node[CHAIN_ID, HEAD, RPC_CLIENT]) verifyConn(ctx context.Context, lggr logger.Logger) NodeState {
	state := n.verifyChainID(ctx, lggr)
	if state != nodeStateAlive {
		return state
	}

	if n.nodePoolCfg.NodeIsSyncingEnabled() {
		isSyncing, err := n.rpc.IsSyncing(ctx)
		if err != nil {
			lggr.Errorw("Unexpected error while verifying RPC node synchronization status", "err", err, "NodeState", n.State())
			return nodeStateUnreachable
		}

		if isSyncing {
			lggr.Errorw("Verification failed: Node is syncing", "NodeState", n.State())
			return nodeStateSyncing
		}
	}

	return nodeStateAlive
}

func (n *node[CHAIN_ID, HEAD, RPC_CLIENT]) Order() int32 {
	return n.order
}
