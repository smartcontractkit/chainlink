package client

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

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
}

//go:generate mockery --quiet --name Node --structname mockNode --filename "mock_node_test.go" --inpackage --case=underscore
type Node[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC NodeClient[CHAIN_ID, HEAD],
] interface {
	// State returns nodeState
	State() nodeState
	// StateAndLatest returns nodeState with the latest received block number & total difficulty.
	StateAndLatest() (nodeState, int64, *big.Int)
	// Name is a unique identifier for this node.
	Name() string
	String() string
	RPC() RPC
	SubscribersCount() int32
	UnsubscribeAllExceptAliveLoop()
	ConfiguredChainID() CHAIN_ID
	Order() int32
	Start(context.Context) error
	Close() error
}

type node[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC NodeClient[CHAIN_ID, HEAD],
] struct {
	services.StateMachine
	lfcLog              logger.Logger
	name                string
	id                  int32
	chainID             CHAIN_ID
	nodePoolCfg         NodeConfig
	noNewHeadsThreshold time.Duration
	order               int32
	chainFamily         string

	ws   url.URL
	http *url.URL

	rpc RPC

	stateMu sync.RWMutex // protects state* fields
	state   nodeState
	// Each node is tracking the last received head number and total difficulty
	stateLatestBlockNumber     int64
	stateLatestTotalDifficulty *big.Int

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
	nLiveNodes func() (count int, blockNumber int64, totalDifficulty *big.Int)
}

func NewNode[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC NodeClient[CHAIN_ID, HEAD],
](
	nodeCfg NodeConfig,
	noNewHeadsThreshold time.Duration,
	lggr logger.Logger,
	wsuri url.URL,
	httpuri *url.URL,
	name string,
	id int32,
	chainID CHAIN_ID,
	nodeOrder int32,
	rpc RPC,
	chainFamily string,
) Node[CHAIN_ID, HEAD, RPC] {
	n := new(node[CHAIN_ID, HEAD, RPC])
	n.name = name
	n.id = id
	n.chainID = chainID
	n.nodePoolCfg = nodeCfg
	n.noNewHeadsThreshold = noNewHeadsThreshold
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
	return n
}

func (n *node[CHAIN_ID, HEAD, RPC]) String() string {
	s := fmt.Sprintf("(%s)%s:%s", Primary.String(), n.name, n.ws.String())
	if n.http != nil {
		s = s + fmt.Sprintf(":%s", n.http.String())
	}
	return s
}

func (n *node[CHAIN_ID, HEAD, RPC]) ConfiguredChainID() (chainID CHAIN_ID) {
	return n.chainID
}

func (n *node[CHAIN_ID, HEAD, RPC]) Name() string {
	return n.name
}

func (n *node[CHAIN_ID, HEAD, RPC]) RPC() RPC {
	return n.rpc
}

func (n *node[CHAIN_ID, HEAD, RPC]) SubscribersCount() int32 {
	return n.rpc.SubscribersCount()
}

func (n *node[CHAIN_ID, HEAD, RPC]) UnsubscribeAllExceptAliveLoop() {
	n.rpc.UnsubscribeAllExceptAliveLoop()
}

func (n *node[CHAIN_ID, HEAD, RPC]) Close() error {
	return n.StopOnce(n.name, n.close)
}

func (n *node[CHAIN_ID, HEAD, RPC]) close() error {
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
func (n *node[CHAIN_ID, HEAD, RPC]) Start(startCtx context.Context) error {
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
func (n *node[CHAIN_ID, HEAD, RPC]) start(startCtx context.Context) {
	if n.state != nodeStateUndialed {
		panic(fmt.Sprintf("cannot dial node with state %v", n.state))
	}

	if err := n.rpc.Dial(startCtx); err != nil {
		n.lfcLog.Errorw("Dial failed: Node is unreachable", "err", err)
		n.declareUnreachable()
		return
	}
	n.setState(nodeStateDialed)

	if err := n.verify(startCtx); errors.Is(err, errInvalidChainID) {
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
func (n *node[CHAIN_ID, HEAD, RPC]) verify(callerCtx context.Context) (err error) {
	promPoolRPCNodeVerifies.WithLabelValues(n.chainFamily, n.chainID.String(), n.name).Inc()
	promFailed := func() {
		promPoolRPCNodeVerifiesFailed.WithLabelValues(n.chainFamily, n.chainID.String(), n.name).Inc()
	}

	st := n.State()
	switch st {
	case nodeStateDialed, nodeStateOutOfSync, nodeStateInvalidChainID:
	default:
		panic(fmt.Sprintf("cannot verify node in state %v", st))
	}

	var chainID CHAIN_ID
	if chainID, err = n.rpc.ChainID(callerCtx); err != nil {
		promFailed()
		return fmt.Errorf("failed to verify chain ID for node %s: %w", n.name, err)
	} else if chainID.String() != n.chainID.String() {
		promFailed()
		return fmt.Errorf(
			"rpc ChainID doesn't match local chain ID: RPC ID=%s, local ID=%s, node name=%s: %w",
			chainID.String(),
			n.chainID.String(),
			n.name,
			errInvalidChainID,
		)
	}

	promPoolRPCNodeVerifiesSuccess.WithLabelValues(n.chainFamily, n.chainID.String(), n.name).Inc()

	return nil
}

// disconnectAll disconnects all clients connected to the node
// WARNING: NOT THREAD-SAFE
// This must be called from within the n.stateMu lock
func (n *node[CHAIN_ID, HEAD, RPC]) disconnectAll() {
	n.rpc.DisconnectAll()
}

func (n *node[CHAIN_ID, HEAD, RPC]) Order() int32 {
	return n.order
}
