package client

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-relay/pkg/services"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	// PromMultiNodeRPCNodeStates reports current RPC node state
	PromMultiNodeRPCNodeStates = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "multi_node_states",
		Help: "The number of RPC nodes currently in the given state for the given chain",
	}, []string{"network", "chainId", "state"})
	ErroringNodeError = fmt.Errorf("no live nodes available")
)

const (
	NodeSelectionModeHighestHead     = "HighestHead"
	NodeSelectionModeRoundRobin      = "RoundRobin"
	NodeSelectionModeTotalDifficulty = "TotalDifficulty"
	NodeSelectionModePriorityLevel   = "PriorityLevel"
)

type NodeSelector[
	CHAIN_ID types.ID,
	HEAD Head,
	RPC NodeClient[CHAIN_ID, HEAD],
] interface {
	// Select returns a Node, or nil if none can be selected.
	// Implementation must be thread-safe.
	Select() Node[CHAIN_ID, HEAD, RPC]
	// Name returns the strategy name, e.g. "HighestHead" or "RoundRobin"
	Name() string
}

// MultiNode is a generalized multi node client interface that includes methods to interact with different chains.
// It also handles multiple node RPC connections simultaneously.
type MultiNode[
	CHAIN_ID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK_HASH types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	EVENT any,
	EVENT_OPS any,
	TX_RECEIPT types.Receipt[TX_HASH, BLOCK_HASH],
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],
	RPC_CLIENT RPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD],
] interface {
	clientAPI[
		CHAIN_ID,
		SEQ,
		ADDR,
		BLOCK_HASH,
		TX,
		TX_HASH,
		EVENT,
		EVENT_OPS,
		TX_RECEIPT,
		FEE,
		HEAD,
	]
	Close() error
	NodeStates() map[string]string
	SelectNodeRPC() (RPC_CLIENT, error)

	BatchCallContextAll(ctx context.Context, b []any) error
	ConfiguredChainID() CHAIN_ID
	IsL2() bool
}

type multiNode[
	CHAIN_ID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK_HASH types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	EVENT any,
	EVENT_OPS any,
	TX_RECEIPT types.Receipt[TX_HASH, BLOCK_HASH],
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],
	RPC_CLIENT RPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD],
] struct {
	services.StateMachine
	nodes               []Node[CHAIN_ID, HEAD, RPC_CLIENT]
	sendonlys           []SendOnlyNode[CHAIN_ID, RPC_CLIENT]
	chainID             CHAIN_ID
	chainType           config.ChainType
	logger              logger.Logger
	selectionMode       string
	noNewHeadsThreshold time.Duration
	nodeSelector        NodeSelector[CHAIN_ID, HEAD, RPC_CLIENT]
	leaseDuration       time.Duration
	leaseTicker         *time.Ticker
	chainFamily         string

	activeMu   sync.RWMutex
	activeNode Node[CHAIN_ID, HEAD, RPC_CLIENT]

	chStop utils.StopChan
	wg     sync.WaitGroup

	sendOnlyErrorParser func(err error) SendTxReturnCode
}

func NewMultiNode[
	CHAIN_ID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK_HASH types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	EVENT any,
	EVENT_OPS any,
	TX_RECEIPT types.Receipt[TX_HASH, BLOCK_HASH],
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],
	RPC_CLIENT RPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD],
](
	logger logger.Logger,
	selectionMode string,
	leaseDuration time.Duration,
	noNewHeadsThreshold time.Duration,
	nodes []Node[CHAIN_ID, HEAD, RPC_CLIENT],
	sendonlys []SendOnlyNode[CHAIN_ID, RPC_CLIENT],
	chainID CHAIN_ID,
	chainType config.ChainType,
	chainFamily string,
	sendOnlyErrorParser func(err error) SendTxReturnCode,
) MultiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT] {
	nodeSelector := func() NodeSelector[CHAIN_ID, HEAD, RPC_CLIENT] {
		switch selectionMode {
		case NodeSelectionModeHighestHead:
			return NewHighestHeadNodeSelector[CHAIN_ID, HEAD, RPC_CLIENT](nodes)
		case NodeSelectionModeRoundRobin:
			return NewRoundRobinSelector[CHAIN_ID, HEAD, RPC_CLIENT](nodes)
		case NodeSelectionModeTotalDifficulty:
			return NewTotalDifficultyNodeSelector[CHAIN_ID, HEAD, RPC_CLIENT](nodes)
		case NodeSelectionModePriorityLevel:
			return NewPriorityLevelNodeSelector[CHAIN_ID, HEAD, RPC_CLIENT](nodes)
		default:
			panic(fmt.Sprintf("unsupported NodeSelectionMode: %s", selectionMode))
		}
	}()

	lggr := logger.Named("MultiNode").With("chainID", chainID.String())

	c := &multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]{
		nodes:               nodes,
		sendonlys:           sendonlys,
		chainID:             chainID,
		chainType:           chainType,
		logger:              lggr,
		selectionMode:       selectionMode,
		noNewHeadsThreshold: noNewHeadsThreshold,
		nodeSelector:        nodeSelector,
		chStop:              make(chan struct{}),
		leaseDuration:       leaseDuration,
		chainFamily:         chainFamily,
		sendOnlyErrorParser: sendOnlyErrorParser,
	}

	c.logger.Debugf("The MultiNode is configured to use NodeSelectionMode: %s", selectionMode)

	return c
}

// Dial starts every node in the pool
//
// Nodes handle their own redialing and runloops, so this function does not
// return any error if the nodes aren't available
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) Dial(ctx context.Context) error {
	return c.StartOnce("MultiNode", func() (merr error) {
		if len(c.nodes) == 0 {
			return errors.Errorf("no available nodes for chain %s", c.chainID.String())
		}
		var ms services.MultiStart
		for _, n := range c.nodes {
			if n.ConfiguredChainID().String() != c.chainID.String() {
				return ms.CloseBecause(errors.Errorf("node %s has configured chain ID %s which does not match multinode configured chain ID of %s", n.String(), n.ConfiguredChainID().String(), c.chainID.String()))
			}
			rawNode, ok := n.(*node[CHAIN_ID, HEAD, RPC_CLIENT])
			if ok {
				// This is a bit hacky but it allows the node to be aware of
				// pool state and prevent certain state transitions that might
				// otherwise leave no nodes available. It is better to have one
				// node in a degraded state than no nodes at all.
				rawNode.nLiveNodes = c.nLiveNodes
			}
			// node will handle its own redialing and automatic recovery
			if err := ms.Start(ctx, n); err != nil {
				return err
			}
		}
		for _, s := range c.sendonlys {
			if s.ConfiguredChainID().String() != c.chainID.String() {
				return ms.CloseBecause(errors.Errorf("sendonly node %s has configured chain ID %s which does not match multinode configured chain ID of %s", s.String(), s.ConfiguredChainID().String(), c.chainID.String()))
			}
			if err := ms.Start(ctx, s); err != nil {
				return err
			}
		}
		c.wg.Add(1)
		go c.runLoop()

		if c.leaseDuration.Seconds() > 0 && c.selectionMode != NodeSelectionModeRoundRobin {
			c.logger.Infof("The MultiNode will switch to best node every %s", c.leaseDuration.String())
			c.wg.Add(1)
			go c.checkLeaseLoop()
		} else {
			c.logger.Info("Best node switching is disabled")
		}

		return nil
	})
}

// Close tears down the MultiNode and closes all nodes
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) Close() error {
	return c.StopOnce("MultiNode", func() error {
		close(c.chStop)
		c.wg.Wait()

		return services.CloseAll(services.MultiCloser(c.nodes), services.MultiCloser(c.sendonlys))
	})
}

// SelectNodeRPC returns an RPC of an active node. If there are no active nodes it returns an error.
// Call this method from your chain-specific client implementation to access any chain-specific rpc calls.
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) SelectNodeRPC() (rpc RPC_CLIENT, err error) {
	n, err := c.selectNode()
	if err != nil {
		return rpc, err
	}
	return n.RPC(), nil

}

// selectNode returns the active Node, if it is still nodeStateAlive, otherwise it selects a new one from the NodeSelector.
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) selectNode() (node Node[CHAIN_ID, HEAD, RPC_CLIENT], err error) {
	c.activeMu.RLock()
	node = c.activeNode
	c.activeMu.RUnlock()
	if node != nil && node.State() == nodeStateAlive {
		return // still alive
	}

	// select a new one
	c.activeMu.Lock()
	defer c.activeMu.Unlock()
	node = c.activeNode
	if node != nil && node.State() == nodeStateAlive {
		return // another goroutine beat us here
	}

	c.activeNode = c.nodeSelector.Select()

	if c.activeNode == nil {
		c.logger.Criticalw("No live RPC nodes available", "NodeSelectionMode", c.nodeSelector.Name())
		errmsg := fmt.Errorf("no live nodes available for chain %s", c.chainID.String())
		c.SvcErrBuffer.Append(errmsg)
		err = ErroringNodeError
	}

	return c.activeNode, err
}

// nLiveNodes returns the number of currently alive nodes, as well as the highest block number and greatest total difficulty.
// totalDifficulty will be 0 if all nodes return nil.
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) nLiveNodes() (nLiveNodes int, blockNumber int64, totalDifficulty *utils.Big) {
	totalDifficulty = utils.NewBigI(0)
	for _, n := range c.nodes {
		if s, num, td := n.StateAndLatest(); s == nodeStateAlive {
			nLiveNodes++
			if num > blockNumber {
				blockNumber = num
			}
			if td != nil && td.Cmp(totalDifficulty) > 0 {
				totalDifficulty = td
			}
		}
	}
	return
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) checkLease() {
	bestNode := c.nodeSelector.Select()
	for _, n := range c.nodes {
		// Terminate client subscriptions. Services are responsible for reconnecting, which will be routed to the new
		// best node. Only terminate connections with more than 1 subscription to account for the aliveLoop subscription
		if n.State() == nodeStateAlive && n != bestNode && n.SubscribersCount() > 1 {
			c.logger.Infof("Switching to best node from %q to %q", n.String(), bestNode.String())
			n.UnsubscribeAllExceptAliveLoop()
		}
	}

	c.activeMu.Lock()
	if bestNode != c.activeNode {
		c.activeNode = bestNode
	}
	c.activeMu.Unlock()
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) checkLeaseLoop() {
	defer c.wg.Done()
	c.leaseTicker = time.NewTicker(c.leaseDuration)
	defer c.leaseTicker.Stop()

	for {
		select {
		case <-c.leaseTicker.C:
			c.checkLease()
		case <-c.chStop:
			return
		}
	}
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) runLoop() {
	defer c.wg.Done()

	c.report()

	// Prometheus' default interval is 15s, set this to under 7.5s to avoid
	// aliasing (see: https://en.wikipedia.org/wiki/Nyquist_frequency)
	reportInterval := 6500 * time.Millisecond
	monitor := time.NewTicker(utils.WithJitter(reportInterval))
	defer monitor.Stop()

	for {
		select {
		case <-monitor.C:
			c.report()
		case <-c.chStop:
			return
		}
	}
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) report() {
	type nodeWithState struct {
		Node  string
		State string
	}

	var total, dead int
	counts := make(map[nodeState]int)
	nodeStates := make([]nodeWithState, len(c.nodes))
	for i, n := range c.nodes {
		state := n.State()
		nodeStates[i] = nodeWithState{n.String(), state.String()}
		total++
		if state != nodeStateAlive {
			dead++
		}
		counts[state]++
	}
	for _, state := range allNodeStates {
		count := counts[state]
		PromMultiNodeRPCNodeStates.WithLabelValues(c.chainFamily, c.chainID.String(), state.String()).Set(float64(count))
	}

	live := total - dead
	c.logger.Tracew(fmt.Sprintf("MultiNode state: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	if total == dead {
		rerr := fmt.Errorf("no primary nodes available: 0/%d nodes are alive", total)
		c.logger.Criticalw(rerr.Error(), "nodeStates", nodeStates)
		c.SvcErrBuffer.Append(rerr)
	} else if dead > 0 {
		c.logger.Errorw(fmt.Sprintf("At least one primary node is dead: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	}
}

// ClientAPI methods
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BalanceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (*big.Int, error) {
	n, err := c.selectNode()
	if err != nil {
		return nil, err
	}
	return n.RPC().BalanceAt(ctx, account, blockNumber)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BatchCallContext(ctx context.Context, b []any) error {
	n, err := c.selectNode()
	if err != nil {
		return err
	}
	return n.RPC().BatchCallContext(ctx, b)
}

// BatchCallContextAll calls BatchCallContext for every single node including
// sendonlys.
// CAUTION: This should only be used for mass re-transmitting transactions, it
// might have unexpected effects to use it for anything else.
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BatchCallContextAll(ctx context.Context, b []any) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	main, selectionErr := c.selectNode()
	var all []SendOnlyNode[CHAIN_ID, RPC_CLIENT]
	for _, n := range c.nodes {
		all = append(all, n)
	}
	all = append(all, c.sendonlys...)
	for _, n := range all {
		if n == main {
			// main node is used at the end for the return value
			continue
		}
		// Parallel call made to all other nodes with ignored return value
		wg.Add(1)
		go func(n SendOnlyNode[CHAIN_ID, RPC_CLIENT]) {
			defer wg.Done()
			err := n.RPC().BatchCallContext(ctx, b)
			if err != nil {
				c.logger.Debugw("Secondary node BatchCallContext failed", "err", err)
			} else {
				c.logger.Trace("Secondary node BatchCallContext success")
			}
		}(n)
	}

	if selectionErr != nil {
		return selectionErr
	}
	return main.RPC().BatchCallContext(ctx, b)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BlockByHash(ctx context.Context, hash BLOCK_HASH) (h HEAD, err error) {
	n, err := c.selectNode()
	if err != nil {
		return h, err
	}
	return n.RPC().BlockByHash(ctx, hash)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BlockByNumber(ctx context.Context, number *big.Int) (h HEAD, err error) {
	n, err := c.selectNode()
	if err != nil {
		return h, err
	}
	return n.RPC().BlockByNumber(ctx, number)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	n, err := c.selectNode()
	if err != nil {
		return err
	}
	return n.RPC().CallContext(ctx, result, method, args...)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) CallContract(
	ctx context.Context,
	attempt interface{},
	blockNumber *big.Int,
) (rpcErr []byte, extractErr error) {
	n, err := c.selectNode()
	if err != nil {
		return rpcErr, err
	}
	return n.RPC().CallContract(ctx, attempt, blockNumber)
}

// ChainID makes a direct RPC call. In most cases it should be better to use the configured chain id instead by
// calling ConfiguredChainID.
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) ChainID(ctx context.Context) (id CHAIN_ID, err error) {
	n, err := c.selectNode()
	if err != nil {
		return id, err
	}
	return n.RPC().ChainID(ctx)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) ChainType() config.ChainType {
	return c.chainType
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) CodeAt(ctx context.Context, account ADDR, blockNumber *big.Int) (code []byte, err error) {
	n, err := c.selectNode()
	if err != nil {
		return code, err
	}
	return n.RPC().CodeAt(ctx, account, blockNumber)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) ConfiguredChainID() CHAIN_ID {
	return c.chainID
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) EstimateGas(ctx context.Context, call any) (gas uint64, err error) {
	n, err := c.selectNode()
	if err != nil {
		return gas, err
	}
	return n.RPC().EstimateGas(ctx, call)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) FilterEvents(ctx context.Context, query EVENT_OPS) (e []EVENT, err error) {
	n, err := c.selectNode()
	if err != nil {
		return e, err
	}
	return n.RPC().FilterEvents(ctx, query)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) IsL2() bool {
	return c.ChainType().IsL2()
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) LatestBlockHeight(ctx context.Context) (h *big.Int, err error) {
	n, err := c.selectNode()
	if err != nil {
		return h, err
	}
	return n.RPC().LatestBlockHeight(ctx)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) LINKBalance(ctx context.Context, accountAddress ADDR, linkAddress ADDR) (b *assets.Link, err error) {
	n, err := c.selectNode()
	if err != nil {
		return b, err
	}
	return n.RPC().LINKBalance(ctx, accountAddress, linkAddress)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) NodeStates() (states map[string]string) {
	states = make(map[string]string)
	for _, n := range c.nodes {
		states[n.Name()] = n.State().String()
	}
	for _, s := range c.sendonlys {
		states[s.Name()] = s.State().String()
	}
	return
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) PendingSequenceAt(ctx context.Context, addr ADDR) (s SEQ, err error) {
	n, err := c.selectNode()
	if err != nil {
		return s, err
	}
	return n.RPC().PendingSequenceAt(ctx, addr)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) SendEmptyTransaction(
	ctx context.Context,
	newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt any, err error),
	seq SEQ,
	gasLimit uint32,
	fee FEE,
	fromAddress ADDR,
) (txhash string, err error) {
	n, err := c.selectNode()
	if err != nil {
		return txhash, err
	}
	return n.RPC().SendEmptyTransaction(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) SendTransaction(ctx context.Context, tx TX) error {
	main, nodeError := c.selectNode()
	var all []SendOnlyNode[CHAIN_ID, RPC_CLIENT]
	for _, n := range c.nodes {
		all = append(all, n)
	}
	all = append(all, c.sendonlys...)
	for _, n := range all {
		if n == main {
			// main node is used at the end for the return value
			continue
		}
		// Parallel send to all other nodes with ignored return value
		// Async - we do not want to block the main thread with secondary nodes
		// in case they are unreliable/slow.
		// It is purely a "best effort" send.
		// Resource is not unbounded because the default context has a timeout.
		ok := c.IfNotStopped(func() {
			// Must wrap inside IfNotStopped to avoid waitgroup racing with Close
			c.wg.Add(1)
			go func(n SendOnlyNode[CHAIN_ID, RPC_CLIENT]) {
				defer c.wg.Done()

				txErr := n.RPC().SendTransaction(ctx, tx)
				c.logger.Debugw("Sendonly node sent transaction", "name", n.String(), "tx", tx, "err", txErr)
				sendOnlyError := c.sendOnlyErrorParser(txErr)
				if sendOnlyError != Successful {
					c.logger.Warnw("RPC returned error", "name", n.String(), "tx", tx, "err", txErr)
				}
			}(n)
		})
		if !ok {
			c.logger.Debug("Cannot send transaction on sendonly node; MultiNode is stopped", "node", n.String())
		}
	}
	if nodeError != nil {
		return nodeError
	}
	return main.RPC().SendTransaction(ctx, tx)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) SequenceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (s SEQ, err error) {
	n, err := c.selectNode()
	if err != nil {
		return s, err
	}
	return n.RPC().SequenceAt(ctx, account, blockNumber)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) SimulateTransaction(ctx context.Context, tx TX) error {
	n, err := c.selectNode()
	if err != nil {
		return err
	}
	return n.RPC().SimulateTransaction(ctx, tx)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (s types.Subscription, err error) {
	n, err := c.selectNode()
	if err != nil {
		return s, err
	}
	return n.RPC().Subscribe(ctx, channel, args...)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) TokenBalance(ctx context.Context, account ADDR, tokenAddr ADDR) (b *big.Int, err error) {
	n, err := c.selectNode()
	if err != nil {
		return b, err
	}
	return n.RPC().TokenBalance(ctx, account, tokenAddr)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) TransactionByHash(ctx context.Context, txHash TX_HASH) (tx TX, err error) {
	n, err := c.selectNode()
	if err != nil {
		return tx, err
	}
	return n.RPC().TransactionByHash(ctx, txHash)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) TransactionReceipt(ctx context.Context, txHash TX_HASH) (txr TX_RECEIPT, err error) {
	n, err := c.selectNode()
	if err != nil {
		return txr, err
	}
	return n.RPC().TransactionReceipt(ctx, txHash)
}
