package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

var (
	// PromMultiNodeClientRPCNodeStates reports current RPC node state
	PromMultiNodeClientRPCNodeStates = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "client_rpc_node_states",
		Help: "The number of RPC nodes currently in the given state for the given chain",
	}, []string{"network", "chainId", "state"})
)

const (
	NodeSelectionMode_HighestHead     = "HighestHead"
	NodeSelectionMode_RoundRobin      = "RoundRobin"
	NodeSelectionMode_TotalDifficulty = "TotalDifficulty"
	NodeSelectionMode_PriorityLevel   = "PriorityLevel"
)

type NodeSelector[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
	RPC_CLIENT NodeClientAPI[CHAIN_ID, BLOCK_HASH, HEAD, SUB],
] interface {
	// Select returns a Node, or nil if none can be selected.
	// Implementation must be thread-safe.
	Select() Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]
	// Name returns the strategy name, e.g. "HighestHead" or "RoundRobin"
	Name() string
}

type MultiNodeClient[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
	RPC_CLIENT NodeClientAPI[CHAIN_ID, BLOCK_HASH, HEAD, SUB],
] interface {
	Dial(context.Context) error
	Close() error
	NodeStates() map[string]string
	SelectNode() Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]

	runLoop()
	nLiveNodes() (int, int64, *utils.Big)
	report()
}

func ContextWithDefaultTimeout() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), queryTimeout)
}

type multiNodeClient[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
	RPC_CLIENT NodeClientAPI[CHAIN_ID, BLOCK_HASH, HEAD, SUB],
] struct {
	utils.StartStopOnce
	nodes               []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]
	sendonlys           []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]
	chainID             CHAIN_ID
	logger              logger.Logger
	selectionMode       string
	noNewHeadsThreshold time.Duration
	nodeSelector        NodeSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]
	chainFamily         string

	activeMu   sync.RWMutex
	activeNode Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]

	chStop utils.StopChan
	wg     sync.WaitGroup
}

func NewMultiNodeClient[
	CHAIN_ID types.ID,
	BLOCK_HASH types.Hashable,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
	RPC_CLIENT NodeClientAPI[CHAIN_ID, BLOCK_HASH, HEAD, SUB],
](
	logger logger.Logger,
	selectionMode string,
	noNewHeadsThreshold time.Duration,
	nodes []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT],
	sendonlys []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT],
	chainID CHAIN_ID,
	chainFamily string,
) MultiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT] {
	nodeSelector := func() NodeSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT] {
		switch selectionMode {
		case NodeSelectionMode_HighestHead:
			return NewHighestHeadNodeSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT](nodes)
		case NodeSelectionMode_RoundRobin:
			return NewRoundRobinSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT](nodes)
		case NodeSelectionMode_TotalDifficulty:
			return NewTotalDifficultyNodeSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT](nodes)
		case NodeSelectionMode_PriorityLevel:
			return NewPriorityLevelNodeSelector[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT](nodes)
		default:
			panic(fmt.Sprintf("unsupported NodeSelectionMode: %s", selectionMode))
		}
	}()

	lggr := logger.Named("MultiNodeClient").With("chainID", chainID.String())

	c := &multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]{
		nodes:               nodes,
		sendonlys:           sendonlys,
		chainID:             chainID,
		logger:              lggr,
		selectionMode:       selectionMode,
		noNewHeadsThreshold: noNewHeadsThreshold,
		nodeSelector:        nodeSelector,
		chStop:              make(chan struct{}),
		chainFamily:         chainFamily,
	}

	c.logger.Debugf("The MultiNodeClient is configured to use NodeSelectionMode: %s", selectionMode)

	return c
}

// selectNode returns the active Node, if it is still NodeStateAlive, otherwise it selects a new one from the NodeSelector.
func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) SelectNode() (node Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) {
	c.activeMu.RLock()
	node = c.activeNode
	c.activeMu.RUnlock()
	if node != nil && node.State() == NodeStateAlive {
		return // still alive
	}

	// select a new one
	c.activeMu.Lock()
	defer c.activeMu.Unlock()
	node = c.activeNode
	if node != nil && node.State() == NodeStateAlive {
		return // another goroutine beat us here
	}

	c.activeNode = c.nodeSelector.Select()

	if c.activeNode == nil {
		c.logger.Criticalw("No live RPC nodes available", "NodeSelectionMode", c.nodeSelector.Name())
		errmsg := fmt.Errorf("no live nodes available for chain %s", c.chainID.String())
		c.SvcErrBuffer.Append(errmsg)
		return &erroringNode[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]{errMsg: errmsg.Error()}
	}

	return c.activeNode
}

func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Dial(ctx context.Context) error {
	return c.StartOnce("Client", func() (merr error) {
		if len(c.nodes) == 0 {
			return errors.Errorf("no available nodes for chain %s", c.chainID.String())
		}
		var ms services.MultiStart
		for _, n := range c.nodes {
			chainID, err := n.ChainID()
			if err != nil {
				return errors.Errorf("Invalid ChainID")
			}
			if chainID.String() != c.chainID.String() {
				return ms.CloseBecause(errors.Errorf("node %s has chain ID %s which does not match client chain ID of %s", n.String(), chainID.String(), c.chainID.String()))
			}
			rawNode, ok := n.(*node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT])
			if ok {
				// This is a bit hacky but it allows the node to be aware of
				// client / pool state and prevent certain state transitions that might
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
			chainID, err := s.ChainID()
			if err != nil {
				return errors.Errorf("Invalid ChainID")
			}
			if chainID.String() != c.chainID.String() {
				return ms.CloseBecause(errors.Errorf("sendonly node %s has chain ID %s which does not match client chain ID of %s", s.String(), chainID.String(), c.chainID.String()))
			}
			if err := ms.Start(ctx, s); err != nil {
				return err
			}
		}
		c.wg.Add(1)
		go c.runLoop()

		return nil
	})
}

// Close tears down the pool and closes all nodes
func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Close() error {
	return c.StopOnce("Client", func() error {
		close(c.chStop)
		c.wg.Wait()

		return services.CloseAll(services.MultiCloser(c.nodes), services.MultiCloser(c.sendonlys))
	})
}

func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) NodeStates() (states map[string]string) {
	states = make(map[string]string)
	for _, n := range c.nodes {
		states[n.Name()] = n.State().String()
	}
	for _, s := range c.sendonlys {
		states[s.Name()] = s.State().String()
	}
	return
}

// nLiveNodes returns the number of currently alive nodes, as well as the highest block number and greatest total difficulty.
// totalDifficulty will be 0 if all nodes return nil.
func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) nLiveNodes() (nLiveNodes int, blockNumber int64, totalDifficulty *utils.Big) {
	totalDifficulty = utils.NewBigI(0)
	for _, n := range c.nodes {
		if s, num, td := n.StateAndLatest(); s == NodeStateAlive {
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

func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) runLoop() {
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

func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) report() {
	type nodeWithState struct {
		Node  string
		State string
	}

	var total, dead int
	counts := make(map[NodeState]int)
	nodeStates := make([]nodeWithState, len(c.nodes))
	for i, n := range c.nodes {
		state := n.State()
		nodeStates[i] = nodeWithState{n.String(), state.String()}
		total++
		if state != NodeStateAlive {
			dead++
		}
		counts[state]++
	}
	for _, state := range allNodeStates {
		count := counts[state]
		PromMultiNodeClientRPCNodeStates.WithLabelValues(c.chainFamily, c.chainID.String(), state.String()).Set(float64(count))
	}

	live := total - dead
	c.logger.Tracew(fmt.Sprintf("Client state: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	if total == dead {
		rerr := fmt.Errorf("no primary nodes available: 0/%d nodes are alive", total)
		c.logger.Criticalw(rerr.Error(), "nodeStates", nodeStates)
		c.SvcErrBuffer.Append(rerr)
	} else if dead > 0 {
		c.logger.Errorw(fmt.Sprintf("At least one primary node is dead: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	}
}

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) BalanceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (*big.Int, error) {
// 	return c.selectNode().RPCClient().BalanceAt(ctx, account, blockNumber)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) BatchCallContext(ctx context.Context, b []any) error {
// 	return c.selectNode().RPCClient().BatchCallContext(ctx, b)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) BatchCallContextAll(ctx context.Context, b []any) error {
// 	var wg sync.WaitGroup
// 	defer wg.Wait()

// 	main := c.selectNode()
// 	var all []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]
// 	all = append(all, c.nodes...)
// 	all = append(all, c.sendonlys...)
// 	for _, n := range all {
// 		if n == main {
// 			// main node is used at the end for the return value
// 			continue
// 		}
// 		// Parallel call made to all other nodes with ignored return value
// 		wg.Add(1)
// 		go func(n Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) {
// 			defer wg.Done()
// 			err := n.RPCClient().BatchCallContext(ctx, b)
// 			if err != nil {
// 				c.logger.Debugw("Secondary node BatchCallContext failed", "err", err)
// 			} else {
// 				c.logger.Trace("Secondary node BatchCallContext success")
// 			}
// 		}(n)
// 	}

// 	return main.RPCClient().BatchCallContext(ctx, b)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) BlockByHash(ctx context.Context, hash BLOCK_HASH) (HEAD, error) {
// 	return c.selectNode().RPCClient().BlockByHash(ctx, hash)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) BlockByNumber(ctx context.Context, number *big.Int) (HEAD, error) {
// 	return c.selectNode().RPCClient().BlockByNumber(ctx, number)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
// 	return c.selectNode().RPCClient().CallContext(ctx, result, method, args...)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) CallContract(
// 	ctx context.Context,
// 	attempt interface{},
// 	blockNumber *big.Int,
// ) (rpcErr []byte, extractErr error) {
// 	return c.selectNode().RPCClient().CallContract(ctx, attempt, blockNumber)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) ChainID() (CHAIN_ID, error) {
// 	return c.selectNode().RPCClient().ChainID()
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) CodeAt(ctx context.Context, account ADDR, blockNumber *big.Int) ([]byte, error) {
// 	return c.selectNode().RPCClient().CodeAt(ctx, account, blockNumber)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) ConfiguredChainID() CHAIN_ID {
// 	return c.selectNode().RPCClient().ConfiguredChainID()
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) EstimateGas(ctx context.Context, call any) (gas uint64, err error) {
// 	return c.selectNode().RPCClient().EstimateGas(ctx, call)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) FilterEvents(ctx context.Context, query EVENT_OPS) ([]EVENT, error) {
// 	return c.selectNode().RPCClient().FilterEvents(ctx, query)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) IsL2() bool {
// 	return c.ChainType().IsL2()
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
// 	return c.selectNode().RPCClient().LatestBlockHeight(ctx)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) LINKBalance(ctx context.Context, accountAddress ADDR, linkAddress ADDR) (*assets.Link, error) {
// 	return c.selectNode().RPCClient().LINKBalance(ctx, accountAddress, linkAddress)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) PendingSequenceAt(ctx context.Context, addr ADDR) (SEQ, error) {
// 	return c.selectNode().RPCClient().PendingSequenceAt(ctx, addr)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) SendEmptyTransaction(
// 	ctx context.Context,
// 	newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt any, err error),
// 	seq SEQ,
// 	gasLimit uint32,
// 	fee FEE,
// 	fromAddress ADDR,
// ) (txhash string, err error) {
// 	return c.selectNode().RPCClient().SendEmptyTransaction(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) SendTransaction(ctx context.Context, tx *TX) error {
// 	main := c.selectNode()
// 	var all []Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]
// 	all = append(all, c.nodes...)
// 	all = append(all, c.sendonlys...)
// 	for _, n := range all {
// 		if n == main {
// 			// main node is used at the end for the return value
// 			continue
// 		}
// 		// Parallel send to all other nodes with ignored return value
// 		// Async - we do not want to block the main thread with secondary nodes
// 		// in case they are unreliable/slow.
// 		// It is purely a "best effort" send.
// 		// Resource is not unbounded because the default context has a timeout.
// 		ok := c.IfNotStopped(func() {
// 			// Must wrap inside IfNotStopped to avoid waitgroup racing with Close
// 			c.wg.Add(1)
// 			go func(n Node[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) {
// 				defer c.wg.Done()

// 				sendCtx, cancel := c.chStop.CtxCancel(ContextWithDefaultTimeout())
// 				defer cancel()
// 				err, _ := n.RPCClient().SendTransactionReturnCode(sendCtx, tx)
// 				c.logger.Debugw("Sendonly node sent transaction", "name", n.String(), "tx", tx, "err", err)
// 				if err == TransactionAlreadyKnown || err == Successful {
// 					// Nonce too low or transaction known errors are expected since
// 					// the primary SendTransaction may well have succeeded already
// 					return
// 				}
// 				c.logger.Warnw("Eth client returned error", "name", n.String(), "err", err, "tx", tx)
// 			}(n)
// 		})
// 		if !ok {
// 			c.logger.Debug("Cannot send transaction on sendonly node; pool is stopped", "node", n.String())
// 		}
// 	}

// 	return main.RPCClient().SendTransaction(ctx, tx)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) SendTransactionReturnCode(
// 	ctx context.Context,
// 	tx *TX,
// ) (SendTxReturnCode, error) {
// 	return c.selectNode().RPCClient().SendTransactionReturnCode(ctx, tx)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) SequenceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (SEQ, error) {
// 	return c.selectNode().RPCClient().SequenceAt(ctx, account, blockNumber)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) SimulateTransaction(ctx context.Context, tx *TX) error {
// 	return c.selectNode().RPCClient().SimulateTransaction(ctx, tx)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (SUB, error) {
// 	return c.selectNode().RPCClient().Subscribe(ctx, channel, args)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) TokenBalance(ctx context.Context, account ADDR, tokenAddr ADDR) (*big.Int, error) {
// 	return c.selectNode().RPCClient().TokenBalance(ctx, account, tokenAddr)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) TransactionByHash(ctx context.Context, txHash TX_HASH) (*TX, error) {
// 	return c.selectNode().RPCClient().TransactionByHash(ctx, txHash)
// }

// func (c *multiNodeClient[CHAIN_ID, BLOCK_HASH, HEAD, SUB, RPC_CLIENT]) TransactionReceipt(ctx context.Context, txHash TX_HASH) (TX_RECEIPT, error) {
// 	return c.selectNode().RPCClient().TransactionReceipt(ctx, txHash)
// }
