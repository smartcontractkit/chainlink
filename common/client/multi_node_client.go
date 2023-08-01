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
	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/config"
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
	}, []string{"chainId", "state"})
)

const (
	NodeSelectionMode_HighestHead     = "HighestHead"
	NodeSelectionMode_RoundRobin      = "RoundRobin"
	NodeSelectionMode_TotalDifficulty = "TotalDifficulty"
	NodeSelectionMode_PriorityLevel   = "PriorityLevel"
)

type NodeSelector[
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
	HEAD *types.Head[BLOCKHASH],
] interface {
	// Select returns a Node, or nil if none can be selected.
	// Implementation must be thread-safe.
	Select() Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]
	// Name returns the strategy name, e.g. "HighestHead" or "RoundRobin"
	Name() string
}

type MultiNodeClient[
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
	HEAD *types.Head[BLOCKHASH],
] interface {
	RPCClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]
	Dial(context.Context) error
	Close() error
	NodeStates() map[string]string
	BatchCallContextAll(ctx context.Context, b []any) error
	runLoop()
	nLiveNodes() (int, int64, *utils.Big)
	report()
}

func ContextWithDefaultTimeout() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), queryTimeout)
}

type multiNodeClient[
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
	HEAD *types.Head[BLOCKHASH],
] struct {
	utils.StartStopOnce
	nodes               []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]
	sendonlys           []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]
	chainID             CHAINID
	chainType           config.ChainType
	logger              logger.Logger
	selectionMode       string
	noNewHeadsThreshold time.Duration
	nodeSelector        NodeSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]

	activeMu   sync.RWMutex
	activeNode Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]

	chStop utils.StopChan
	wg     sync.WaitGroup
}

func NewMultiNodeClient[
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
	HEAD *types.Head[BLOCKHASH],
](logger logger.Logger, selectionMode string, noNewHeadsTreshold time.Duration, nodes []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD], sendonlys []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD], chainID CHAINID, chainType config.ChainType,
) MultiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD] {
	if &chainID == nil {
		panic("chainID is required")
	}

	nodeSelector := func() NodeSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD] {
		switch selectionMode {
		case NodeSelectionMode_HighestHead:
			return NewHighestHeadNodeSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD](nodes)
		case NodeSelectionMode_RoundRobin:
			return NewRoundRobinSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD](nodes)
		case NodeSelectionMode_TotalDifficulty:
			return NewTotalDifficultyNodeSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD](nodes)
		case NodeSelectionMode_PriorityLevel:
			return NewPriorityLevelNodeSelector[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD](nodes)
		default:
			panic(fmt.Sprintf("unsupported NodeSelectionMode: %s", selectionMode))
		}
	}()

	lggr := logger.Named("MultiNodeClient").With("evmChainID", chainID.String())

	c := &multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]{
		nodes:               nodes,
		sendonlys:           sendonlys,
		chainID:             chainID,
		chainType:           chainType,
		logger:              lggr,
		selectionMode:       selectionMode,
		noNewHeadsThreshold: noNewHeadsTreshold,
		nodeSelector:        nodeSelector,
		chStop:              make(chan struct{}),
	}

	c.logger.Debugf("The MultiNodeClient is configured to use NodeSelectionMode: %s", selectionMode)

	return c
}

// selectNode returns the active Node, if it is still NodeStateAlive, otherwise it selects a new one from the NodeSelector.
func (c *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) selectNode() (node Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) {
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
		return &erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]{errMsg: errmsg.Error()}
	}

	return c.activeNode
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) Dial(ctx context.Context) error {
	return client.StartOnce("Client", func() (merr error) {
		if len(client.nodes) == 0 {
			return errors.Errorf("no available nodes for chain %s", client.chainID.String())
		}
		var ms services.MultiStart
		for _, n := range client.nodes {
			chainID, err := n.ChainID()
			if err != nil {
				return errors.Errorf("Invalid ChainID")
			}
			if chainID.String() != client.chainID.String() {
				return ms.CloseBecause(errors.Errorf("node %s has chain ID %s which does not match client chain ID of %s", n.String(), chainID.String(), client.chainID.String()))
			}
			rawNode, ok := n.(*node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD])
			if ok {
				// This is a bit hacky but it allows the node to be aware of
				// client / pool state and prevent certain state transitions that might
				// otherwise leave no nodes available. It is better to have one
				// node in a degraded state than no nodes at all.
				rawNode.nLiveNodes = rawNode.nLiveNodes
			}
			// node will handle its own redialing and automatic recovery
			if err := ms.Start(ctx, n); err != nil {
				return err
			}
		}
		for _, s := range client.sendonlys {
			chainID, err := s.ChainID()
			if err != nil {
				return errors.Errorf("Invalid ChainID")
			}
			if chainID.String() != client.chainID.String() {
				return ms.CloseBecause(errors.Errorf("sendonly node %s has chain ID %s which does not match client chain ID of %s", s.String(), chainID.String(), client.chainID.String()))
			}
			if err := ms.Start(ctx, s); err != nil {
				return err
			}
		}
		client.wg.Add(1)
		go client.runLoop()

		return nil
	})
}

// Close tears down the pool and closes all nodes
func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) Close() error {
	return client.StopOnce("Client", func() error {
		close(client.chStop)
		client.wg.Wait()

		return services.CloseAll(services.MultiCloser(client.nodes), services.MultiCloser(client.sendonlys))
	})
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) NodeStates() (states map[string]string) {
	states = make(map[string]string)
	for _, n := range client.nodes {
		states[n.Name()] = n.State().String()
	}
	for _, s := range client.sendonlys {
		states[s.Name()] = s.State().String()
	}
	return
}

// nLiveNodes returns the number of currently alive nodes, as well as the highest block number and greatest total difficulty.
// totalDifficulty will be 0 if all nodes return nil.
func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) nLiveNodes() (nLiveNodes int, blockNumber int64, totalDifficulty *utils.Big) {
	totalDifficulty = utils.NewBigI(0)
	for _, n := range client.nodes {
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

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) runLoop() {
	defer client.wg.Done()

	client.report()

	// Prometheus' default interval is 15s, set this to under 7.5s to avoid
	// aliasing (see: https://en.wikipedia.org/wiki/Nyquist_frequency)
	reportInterval := 6500 * time.Millisecond
	monitor := time.NewTicker(utils.WithJitter(reportInterval))
	defer monitor.Stop()

	for {
		select {
		case <-monitor.C:
			client.report()
		case <-client.chStop:
			return
		}
	}
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) ChainType() config.ChainType {
	return client.chainType
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) report() {
	type nodeWithState struct {
		Node  string
		State string
	}

	var total, dead int
	counts := make(map[NodeState]int)
	nodeStates := make([]nodeWithState, len(client.nodes))
	for i, n := range client.nodes {
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
		PromMultiNodeClientRPCNodeStates.WithLabelValues(client.chainID.String(), state.String()).Set(float64(count))
	}

	live := total - dead
	client.logger.Tracew(fmt.Sprintf("Client state: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	if total == dead {
		rerr := fmt.Errorf("no EVM primary nodes available: 0/%d nodes are alive", total)
		client.logger.Criticalw(rerr.Error(), "nodeStates", nodeStates)
		client.SvcErrBuffer.Append(rerr)
	} else if dead > 0 {
		client.logger.Errorw(fmt.Sprintf("At least one EVM primary node is dead: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	}
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) BalanceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (*big.Int, error) {
	return client.selectNode().BalanceAt(ctx, account, blockNumber)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) BatchCallContext(ctx context.Context, b []any) error {
	return client.selectNode().BatchCallContext(ctx, b)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) BatchCallContextAll(ctx context.Context, b []any) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	main := client.selectNode()
	var all []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]
	for _, n := range client.nodes {
		all = append(all, n)
	}
	all = append(all, client.sendonlys...)
	for _, n := range all {
		if n == main {
			// main node is used at the end for the return value
			continue
		}
		// Parallel call made to all other nodes with ignored return value
		wg.Add(1)
		go func(n Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) {
			defer wg.Done()
			err := n.BatchCallContext(ctx, b)
			if err != nil {
				client.logger.Debugw("Secondary node BatchCallContext failed", "err", err)
			} else {
				client.logger.Trace("Secondary node BatchCallContext success")
			}
		}(n)
	}

	return main.BatchCallContext(ctx, b)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) BlockByHash(ctx context.Context, hash BLOCKHASH) (*BLOCK, error) {
	return client.selectNode().BlockByHash(ctx, hash)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) BlockByNumber(ctx context.Context, number *big.Int) (*BLOCK, error) {
	return client.selectNode().BlockByNumber(ctx, number)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return client.selectNode().CallContext(ctx, result, method, args...)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) CallContract(
	ctx context.Context,
	attempt interface{},
	blockNumber *big.Int,
) (rpcErr []byte, extractErr error) {
	return client.selectNode().CallContract(ctx, attempt, blockNumber)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) ChainID() (CHAINID, error) {
	return client.selectNode().ChainID()
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) CodeAt(ctx context.Context, account ADDR, blockNumber *big.Int) ([]byte, error) {
	return client.selectNode().CodeAt(ctx, account, blockNumber)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) ConfiguredChainID() CHAINID {
	return client.selectNode().ConfiguredChainID()
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) EstimateGas(ctx context.Context, call any) (gas uint64, err error) {
	return client.selectNode().EstimateGas(ctx, call)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) FilterEvents(ctx context.Context, query EVENTOPS) ([]EVENT, error) {
	return client.selectNode().FilterEvents(ctx, query)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) HeadByNumber(ctx context.Context, number *big.Int) (head HEAD, err error) {
	return client.selectNode().HeadByNumber(ctx, number)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) HeadByHash(ctx context.Context, hash BLOCKHASH) (head HEAD, err error) {
	return client.selectNode().HeadByHash(ctx, hash)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) IsL2() bool {
	return client.selectNode().IsL2()
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
	return client.selectNode().LatestBlockHeight(ctx)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) LINKBalance(ctx context.Context, accountAddress ADDR, linkAddress ADDR) (*assets.Link, error) {
	return client.selectNode().LINKBalance(ctx, accountAddress, linkAddress)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) PendingSequenceAt(ctx context.Context, addr ADDR) (SEQ, error) {
	return client.selectNode().PendingSequenceAt(ctx, addr)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) SendEmptyTransaction(
	ctx context.Context,
	newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE], err error),
	seq SEQ,
	gasLimit uint32,
	fee FEE,
	fromAddress ADDR,
) (txhash string, err error) {
	return client.selectNode().SendEmptyTransaction(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) SendTransaction(ctx context.Context, tx *TX) error {
	main := client.selectNode()
	var all []Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]
	for _, n := range client.nodes {
		all = append(all, n)
	}
	all = append(all, client.sendonlys...)
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
		ok := client.IfNotStopped(func() {
			// Must wrap inside IfNotStopped to avoid waitgroup racing with Close
			client.wg.Add(1)
			go func(n Node[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) {
				defer client.wg.Done()

				sendCtx, cancel := client.chStop.CtxCancel(ContextWithDefaultTimeout())
				defer cancel()
				err, _ := n.SendTransactionReturnCode(sendCtx, tx)
				client.logger.Debugw("Sendonly node sent transaction", "name", n.String(), "tx", tx, "err", err)
				if err == clienttypes.TransactionAlreadyKnown || err == clienttypes.Successful {
					// Nonce too low or transaction known errors are expected since
					// the primary SendTransaction may well have succeeded already
					return
				}
				client.logger.Warnw("Eth client returned error", "name", n.String(), "err", err, "tx", tx)
			}(n)
		})
		if !ok {
			client.logger.Debug("Cannot send transaction on sendonly node; pool is stopped", "node", n.String())
		}
	}

	return main.SendTransaction(ctx, tx)
}

// func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) SendTransactionReturnCode(
// 	ctx context.Context,
// 	TX any,
// 	attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
// 	lggr logger.Logger,
// ) (clienttypes.SendTxReturnCode, error) {
// 	return client.selectNode().SendTransactionReturnCode(ctx, tx, attempt, lggr)
// }

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) SendTransactionReturnCode(
	ctx context.Context,
	tx any,
) (clienttypes.SendTxReturnCode, error) {
	return client.selectNode().SendTransactionReturnCode(ctx, tx)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) SequenceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (SEQ, error) {
	return client.selectNode().SequenceAt(ctx, account, blockNumber)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) SimulateTransaction(ctx context.Context, tx *TX) error {
	return client.selectNode().SimulateTransaction(ctx, tx)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) Subscribe(ctx context.Context, channel chan<- types.Head[BLOCKHASH], args ...interface{}) (types.Subscription, error) {
	return client.selectNode().Subscribe(ctx, channel, args)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) TokenBalance(ctx context.Context, account ADDR, tokenAddr ADDR) (*big.Int, error) {
	return client.selectNode().TokenBalance(ctx, account, tokenAddr)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) TransactionByHash(ctx context.Context, txHash TXHASH) (*TX, error) {
	return client.selectNode().TransactionByHash(ctx, txHash)
}

func (client *multiNodeClient[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) TransactionReceipt(ctx context.Context, txHash TXHASH) (*TXRECEIPT, error) {
	return client.selectNode().TransactionReceipt(ctx, txHash)
}
