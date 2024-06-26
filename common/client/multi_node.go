package client

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"slices"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

var (
	// PromMultiNodeRPCNodeStates reports current RPC node state
	PromMultiNodeRPCNodeStates = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "multi_node_states",
		Help: "The number of RPC nodes currently in the given state for the given chain",
	}, []string{"network", "chainId", "state"})
	// PromMultiNodeInvariantViolations reports violation of our assumptions
	PromMultiNodeInvariantViolations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "multi_node_invariant_violations",
		Help: "The number of invariant violations",
	}, []string{"network", "chainId", "invariant"})
	ErroringNodeError = fmt.Errorf("no live nodes available")
)

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
	RPC_CLIENT RPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, BATCH_ELEM],
	BATCH_ELEM any,
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
		BATCH_ELEM,
	]
	Close() error
	NodeStates() map[string]string
	SelectNodeRPC() (RPC_CLIENT, error)

	BatchCallContextAll(ctx context.Context, b []BATCH_ELEM) error
	ConfiguredChainID() CHAIN_ID
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
	RPC_CLIENT RPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, BATCH_ELEM],
	BATCH_ELEM any,
] struct {
	services.StateMachine
	nodes               []Node[CHAIN_ID, HEAD, RPC_CLIENT]
	sendonlys           []SendOnlyNode[CHAIN_ID, RPC_CLIENT]
	chainID             CHAIN_ID
	lggr                logger.SugaredLogger
	selectionMode       string
	noNewHeadsThreshold time.Duration
	nodeSelector        NodeSelector[CHAIN_ID, HEAD, RPC_CLIENT]
	leaseDuration       time.Duration
	leaseTicker         *time.Ticker
	chainFamily         string
	reportInterval      time.Duration
	sendTxSoftTimeout   time.Duration // defines max waiting time from first response til responses evaluation

	activeMu   sync.RWMutex
	activeNode Node[CHAIN_ID, HEAD, RPC_CLIENT]

	chStop services.StopChan
	wg     sync.WaitGroup

	classifySendTxError func(tx TX, err error) SendTxReturnCode
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
	RPC_CLIENT RPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, BATCH_ELEM],
	BATCH_ELEM any,
](
	lggr logger.Logger,
	selectionMode string,
	leaseDuration time.Duration,
	noNewHeadsThreshold time.Duration,
	nodes []Node[CHAIN_ID, HEAD, RPC_CLIENT],
	sendonlys []SendOnlyNode[CHAIN_ID, RPC_CLIENT],
	chainID CHAIN_ID,
	chainFamily string,
	classifySendTxError func(tx TX, err error) SendTxReturnCode,
	sendTxSoftTimeout time.Duration,
) MultiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM] {
	nodeSelector := newNodeSelector(selectionMode, nodes)
	// Prometheus' default interval is 15s, set this to under 7.5s to avoid
	// aliasing (see: https://en.wikipedia.org/wiki/Nyquist_frequency)
	const reportInterval = 6500 * time.Millisecond
	if sendTxSoftTimeout == 0 {
		sendTxSoftTimeout = QueryTimeout / 2
	}
	c := &multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]{
		nodes:               nodes,
		sendonlys:           sendonlys,
		chainID:             chainID,
		lggr:                logger.Sugared(lggr).Named("MultiNode").With("chainID", chainID.String()),
		selectionMode:       selectionMode,
		noNewHeadsThreshold: noNewHeadsThreshold,
		nodeSelector:        nodeSelector,
		chStop:              make(services.StopChan),
		leaseDuration:       leaseDuration,
		chainFamily:         chainFamily,
		classifySendTxError: classifySendTxError,
		reportInterval:      reportInterval,
		sendTxSoftTimeout:   sendTxSoftTimeout,
	}

	c.lggr.Debugf("The MultiNode is configured to use NodeSelectionMode: %s", selectionMode)

	return c
}

// Dial starts every node in the pool
//
// Nodes handle their own redialing and runloops, so this function does not
// return any error if the nodes aren't available
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) Dial(ctx context.Context) error {
	return c.StartOnce("MultiNode", func() (merr error) {
		if len(c.nodes) == 0 {
			return fmt.Errorf("no available nodes for chain %s", c.chainID.String())
		}
		var ms services.MultiStart
		for _, n := range c.nodes {
			if n.ConfiguredChainID().String() != c.chainID.String() {
				return ms.CloseBecause(fmt.Errorf("node %s has configured chain ID %s which does not match multinode configured chain ID of %s", n.String(), n.ConfiguredChainID().String(), c.chainID.String()))
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
				return ms.CloseBecause(fmt.Errorf("sendonly node %s has configured chain ID %s which does not match multinode configured chain ID of %s", s.String(), s.ConfiguredChainID().String(), c.chainID.String()))
			}
			if err := ms.Start(ctx, s); err != nil {
				return err
			}
		}
		c.wg.Add(1)
		go c.runLoop()

		if c.leaseDuration.Seconds() > 0 && c.selectionMode != NodeSelectionModeRoundRobin {
			c.lggr.Infof("The MultiNode will switch to best node every %s", c.leaseDuration.String())
			c.wg.Add(1)
			go c.checkLeaseLoop()
		} else {
			c.lggr.Info("Best node switching is disabled")
		}

		return nil
	})
}

// Close tears down the MultiNode and closes all nodes
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) Close() error {
	return c.StopOnce("MultiNode", func() error {
		close(c.chStop)
		c.wg.Wait()

		return services.CloseAll(services.MultiCloser(c.nodes), services.MultiCloser(c.sendonlys))
	})
}

// SelectNodeRPC returns an RPC of an active node. If there are no active nodes it returns an error.
// Call this method from your chain-specific client implementation to access any chain-specific rpc calls.
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) SelectNodeRPC() (rpc RPC_CLIENT, err error) {
	n, err := c.selectNode()
	if err != nil {
		return rpc, err
	}
	return n.RPC(), nil
}

// selectNode returns the active Node, if it is still nodeStateAlive, otherwise it selects a new one from the NodeSelector.
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) selectNode() (node Node[CHAIN_ID, HEAD, RPC_CLIENT], err error) {
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
		c.lggr.Criticalw("No live RPC nodes available", "NodeSelectionMode", c.nodeSelector.Name())
		errmsg := fmt.Errorf("no live nodes available for chain %s", c.chainID.String())
		c.SvcErrBuffer.Append(errmsg)
		err = ErroringNodeError
	}

	return c.activeNode, err
}

// nLiveNodes returns the number of currently alive nodes, as well as the highest block number and greatest total difficulty.
// totalDifficulty will be 0 if all nodes return nil.
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) nLiveNodes() (nLiveNodes int, blockNumber int64, totalDifficulty *big.Int) {
	totalDifficulty = big.NewInt(0)
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

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) checkLease() {
	bestNode := c.nodeSelector.Select()
	for _, n := range c.nodes {
		// Terminate client subscriptions. Services are responsible for reconnecting, which will be routed to the new
		// best node. Only terminate connections with more than 1 subscription to account for the aliveLoop subscription
		if n.State() == nodeStateAlive && n != bestNode && n.SubscribersCount() > 1 {
			c.lggr.Infof("Switching to best node from %q to %q", n.String(), bestNode.String())
			n.UnsubscribeAllExceptAliveLoop()
		}
	}

	c.activeMu.Lock()
	if bestNode != c.activeNode {
		c.activeNode = bestNode
	}
	c.activeMu.Unlock()
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) checkLeaseLoop() {
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

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) runLoop() {
	defer c.wg.Done()

	c.report()

	monitor := time.NewTicker(utils.WithJitter(c.reportInterval))
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

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) report() {
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
	c.lggr.Tracew(fmt.Sprintf("MultiNode state: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	if total == dead {
		rerr := fmt.Errorf("no primary nodes available: 0/%d nodes are alive", total)
		c.lggr.Criticalw(rerr.Error(), "nodeStates", nodeStates)
		c.SvcErrBuffer.Append(rerr)
	} else if dead > 0 {
		c.lggr.Errorw(fmt.Sprintf("At least one primary node is dead: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	}
}

// ClientAPI methods
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) BalanceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (*big.Int, error) {
	n, err := c.selectNode()
	if err != nil {
		return nil, err
	}
	return n.RPC().BalanceAt(ctx, account, blockNumber)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) BatchCallContext(ctx context.Context, b []BATCH_ELEM) error {
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
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) BatchCallContextAll(ctx context.Context, b []BATCH_ELEM) error {
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

		if n.State() != nodeStateAlive {
			continue
		}
		// Parallel call made to all other nodes with ignored return value
		wg.Add(1)
		go func(n SendOnlyNode[CHAIN_ID, RPC_CLIENT]) {
			defer wg.Done()
			err := n.RPC().BatchCallContext(ctx, b)
			if err != nil {
				c.lggr.Debugw("Secondary node BatchCallContext failed", "err", err)
			} else {
				c.lggr.Trace("Secondary node BatchCallContext success")
			}
		}(n)
	}

	if selectionErr != nil {
		return selectionErr
	}
	return main.RPC().BatchCallContext(ctx, b)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) BlockByHash(ctx context.Context, hash BLOCK_HASH) (h HEAD, err error) {
	n, err := c.selectNode()
	if err != nil {
		return h, err
	}
	return n.RPC().BlockByHash(ctx, hash)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) BlockByNumber(ctx context.Context, number *big.Int) (h HEAD, err error) {
	n, err := c.selectNode()
	if err != nil {
		return h, err
	}
	return n.RPC().BlockByNumber(ctx, number)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	n, err := c.selectNode()
	if err != nil {
		return err
	}
	return n.RPC().CallContext(ctx, result, method, args...)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) CallContract(
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

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) PendingCallContract(
	ctx context.Context,
	attempt interface{},
) (rpcErr []byte, extractErr error) {
	n, err := c.selectNode()
	if err != nil {
		return rpcErr, err
	}
	return n.RPC().PendingCallContract(ctx, attempt)
}

// ChainID makes a direct RPC call. In most cases it should be better to use the configured chain id instead by
// calling ConfiguredChainID.
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) ChainID(ctx context.Context) (id CHAIN_ID, err error) {
	n, err := c.selectNode()
	if err != nil {
		return id, err
	}
	return n.RPC().ChainID(ctx)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) CodeAt(ctx context.Context, account ADDR, blockNumber *big.Int) (code []byte, err error) {
	n, err := c.selectNode()
	if err != nil {
		return code, err
	}
	return n.RPC().CodeAt(ctx, account, blockNumber)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) ConfiguredChainID() CHAIN_ID {
	return c.chainID
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) EstimateGas(ctx context.Context, call any) (gas uint64, err error) {
	n, err := c.selectNode()
	if err != nil {
		return gas, err
	}
	return n.RPC().EstimateGas(ctx, call)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) FilterEvents(ctx context.Context, query EVENT_OPS) (e []EVENT, err error) {
	n, err := c.selectNode()
	if err != nil {
		return e, err
	}
	return n.RPC().FilterEvents(ctx, query)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) LatestBlockHeight(ctx context.Context) (h *big.Int, err error) {
	n, err := c.selectNode()
	if err != nil {
		return h, err
	}
	return n.RPC().LatestBlockHeight(ctx)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) LINKBalance(ctx context.Context, accountAddress ADDR, linkAddress ADDR) (b *assets.Link, err error) {
	n, err := c.selectNode()
	if err != nil {
		return b, err
	}
	return n.RPC().LINKBalance(ctx, accountAddress, linkAddress)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) NodeStates() (states map[string]string) {
	states = make(map[string]string)
	for _, n := range c.nodes {
		states[n.Name()] = n.State().String()
	}
	for _, s := range c.sendonlys {
		states[s.Name()] = s.State().String()
	}
	return
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) PendingSequenceAt(ctx context.Context, addr ADDR) (s SEQ, err error) {
	n, err := c.selectNode()
	if err != nil {
		return s, err
	}
	return n.RPC().PendingSequenceAt(ctx, addr)
}

type sendTxErrors map[SendTxReturnCode][]error

// String - returns string representation of the errors map. Required by logger to properly represent the value
func (errs sendTxErrors) String() string {
	return fmt.Sprint(map[SendTxReturnCode][]error(errs))
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) SendEmptyTransaction(
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

type sendTxResult struct {
	Err        error
	ResultCode SendTxReturnCode
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) broadcastTxAsync(ctx context.Context,
	n SendOnlyNode[CHAIN_ID, RPC_CLIENT], tx TX) sendTxResult {
	txErr := n.RPC().SendTransaction(ctx, tx)
	c.lggr.Debugw("Node sent transaction", "name", n.String(), "tx", tx, "err", txErr)
	resultCode := c.classifySendTxError(tx, txErr)
	if !slices.Contains(sendTxSuccessfulCodes, resultCode) {
		c.lggr.Warnw("RPC returned error", "name", n.String(), "tx", tx, "err", txErr)
	}

	return sendTxResult{Err: txErr, ResultCode: resultCode}
}

// collectTxResults - refer to SendTransaction comment for implementation details,
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) collectTxResults(ctx context.Context, tx TX, healthyNodesNum int, txResults <-chan sendTxResult) error {
	if healthyNodesNum == 0 {
		return ErroringNodeError
	}
	// combine context and stop channel to ensure we stop, when signal received
	ctx, cancel := c.chStop.Ctx(ctx)
	defer cancel()
	requiredResults := int(math.Ceil(float64(healthyNodesNum) * sendTxQuorum))
	errorsByCode := sendTxErrors{}
	var softTimeoutChan <-chan time.Time
	var resultsCount int
loop:
	for {
		select {
		case <-ctx.Done():
			c.lggr.Debugw("Failed to collect of the results before context was done", "tx", tx, "errorsByCode", errorsByCode)
			return ctx.Err()
		case result := <-txResults:
			errorsByCode[result.ResultCode] = append(errorsByCode[result.ResultCode], result.Err)
			resultsCount++
			if slices.Contains(sendTxSuccessfulCodes, result.ResultCode) || resultsCount >= requiredResults {
				break loop
			}
		case <-softTimeoutChan:
			c.lggr.Debugw("Send Tx soft timeout expired - returning responses we've collected so far", "tx", tx, "resultsCount", resultsCount, "requiredResults", requiredResults)
			break loop
		}

		if softTimeoutChan == nil {
			tm := time.NewTimer(c.sendTxSoftTimeout)
			softTimeoutChan = tm.C
			// we are fine with stopping timer at the end of function
			//nolint
			defer tm.Stop()
		}
	}

	// ignore critical error as it's reported in reportSendTxAnomalies
	result, _ := aggregateTxResults(errorsByCode)
	return result
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) reportSendTxAnomalies(tx TX, txResults <-chan sendTxResult) {
	defer c.wg.Done()
	resultsByCode := sendTxErrors{}
	// txResults eventually will be closed
	for txResult := range txResults {
		resultsByCode[txResult.ResultCode] = append(resultsByCode[txResult.ResultCode], txResult.Err)
	}

	_, criticalErr := aggregateTxResults(resultsByCode)
	if criticalErr != nil {
		c.lggr.Criticalw("observed invariant violation on SendTransaction", "tx", tx, "resultsByCode", resultsByCode, "err", criticalErr)
		c.SvcErrBuffer.Append(criticalErr)
		PromMultiNodeInvariantViolations.WithLabelValues(c.chainFamily, c.chainID.String(), criticalErr.Error()).Inc()
	}
}

func aggregateTxResults(resultsByCode sendTxErrors) (txResult error, err error) {
	severeErrors, hasSevereErrors := findFirstIn(resultsByCode, sendTxSevereErrors)
	successResults, hasSuccess := findFirstIn(resultsByCode, sendTxSuccessfulCodes)
	if hasSuccess {
		// We assume that primary node would never report false positive txResult for a transaction.
		// Thus, if such case occurs it's probably due to misconfiguration or a bug and requires manual intervention.
		if hasSevereErrors {
			const errMsg = "found contradictions in nodes replies on SendTransaction: got success and severe error"
			// return success, since at least 1 node has accepted our broadcasted Tx, and thus it can now be included onchain
			return successResults[0], fmt.Errorf(errMsg)
		}

		// other errors are temporary - we are safe to return success
		return successResults[0], nil
	}

	if hasSevereErrors {
		return severeErrors[0], nil
	}

	// return temporary error
	for _, result := range resultsByCode {
		return result[0], nil
	}

	err = fmt.Errorf("expected at least one response on SendTransaction")
	return err, err
}

const sendTxQuorum = 0.7

// SendTransaction - broadcasts transaction to all the send-only and primary nodes regardless of their health.
// A returned nil or error does not guarantee that the transaction will or won't be included. Additional checks must be
// performed to determine the final state.
//
// Send-only nodes' results are ignored as they tend to return false-positive responses. Broadcast to them is necessary
// to speed up the propagation of TX in the network.
//
// Handling of primary nodes' results consists of collection and aggregation.
// In the collection step, we gather as many results as possible while minimizing waiting time. This operation succeeds
// on one of the following conditions:
// * Received at least one success
// * Received at least one result and `sendTxSoftTimeout` expired
// * Received results from the sufficient number of nodes defined by sendTxQuorum.
// The aggregation is based on the following conditions:
// * If there is at least one success - returns success
// * If there is at least one terminal error - returns terminal error
// * If there is both success and terminal error - returns success and reports invariant violation
// * Otherwise, returns any (effectively random) of the errors.
func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) SendTransaction(ctx context.Context, tx TX) error {
	if len(c.nodes) == 0 {
		return ErroringNodeError
	}

	healthyNodesNum := 0
	txResults := make(chan sendTxResult, len(c.nodes))
	// Must wrap inside IfNotStopped to avoid waitgroup racing with Close
	ok := c.IfNotStopped(func() {
		// fire-n-forget, as sendOnlyNodes can not be trusted with result reporting
		for _, n := range c.sendonlys {
			if n.State() != nodeStateAlive {
				continue
			}
			c.wg.Add(1)
			go func(n SendOnlyNode[CHAIN_ID, RPC_CLIENT]) {
				defer c.wg.Done()
				c.broadcastTxAsync(ctx, n, tx)
			}(n)
		}

		var primaryBroadcastWg sync.WaitGroup
		txResultsToReport := make(chan sendTxResult, len(c.nodes))
		for _, n := range c.nodes {
			if n.State() != nodeStateAlive {
				continue
			}

			healthyNodesNum++
			primaryBroadcastWg.Add(1)
			go func(n SendOnlyNode[CHAIN_ID, RPC_CLIENT]) {
				defer primaryBroadcastWg.Done()
				result := c.broadcastTxAsync(ctx, n, tx)
				// both channels are sufficiently buffered, so we won't be locked
				txResultsToReport <- result
				txResults <- result
			}(n)
		}

		c.wg.Add(1)
		go func() {
			// wait for primary nodes to finish the broadcast before closing the channel
			primaryBroadcastWg.Wait()
			close(txResultsToReport)
			close(txResults)
			c.wg.Done()
		}()

		c.wg.Add(1)
		go c.reportSendTxAnomalies(tx, txResultsToReport)
	})
	if !ok {
		return fmt.Errorf("aborted while broadcasting tx - multiNode is stopped: %w", context.Canceled)
	}

	return c.collectTxResults(ctx, tx, healthyNodesNum, txResults)
}

// findFirstIn - returns first existing value for the slice of keys
func findFirstIn[K comparable, V any](set map[K]V, keys []K) (V, bool) {
	for _, k := range keys {
		if v, ok := set[k]; ok {
			return v, true
		}
	}
	var v V
	return v, false
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) SequenceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (s SEQ, err error) {
	n, err := c.selectNode()
	if err != nil {
		return s, err
	}
	return n.RPC().SequenceAt(ctx, account, blockNumber)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) SimulateTransaction(ctx context.Context, tx TX) error {
	n, err := c.selectNode()
	if err != nil {
		return err
	}
	return n.RPC().SimulateTransaction(ctx, tx)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (s types.Subscription, err error) {
	n, err := c.selectNode()
	if err != nil {
		return s, err
	}
	return n.RPC().Subscribe(ctx, channel, args...)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) TokenBalance(ctx context.Context, account ADDR, tokenAddr ADDR) (b *big.Int, err error) {
	n, err := c.selectNode()
	if err != nil {
		return b, err
	}
	return n.RPC().TokenBalance(ctx, account, tokenAddr)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) TransactionByHash(ctx context.Context, txHash TX_HASH) (tx TX, err error) {
	n, err := c.selectNode()
	if err != nil {
		return tx, err
	}
	return n.RPC().TransactionByHash(ctx, txHash)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) TransactionReceipt(ctx context.Context, txHash TX_HASH) (txr TX_RECEIPT, err error) {
	n, err := c.selectNode()
	if err != nil {
		return txr, err
	}
	return n.RPC().TransactionReceipt(ctx, txHash)
}

func (c *multiNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT, BATCH_ELEM]) LatestFinalizedBlock(ctx context.Context) (head HEAD, err error) {
	n, err := c.selectNode()
	if err != nil {
		return head, err
	}

	return n.RPC().LatestFinalizedBlock(ctx)
}
