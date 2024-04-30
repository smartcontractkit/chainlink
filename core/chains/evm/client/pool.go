package client

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	pkgerrors "github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink/v2/common/config"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

var (
	// PromEVMPoolRPCNodeStates reports current RPC node state
	PromEVMPoolRPCNodeStates = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "evm_pool_rpc_node_states",
		Help: "The number of RPC nodes currently in the given state for the given chain",
	}, []string{"evmChainID", "state"})
)

const (
	NodeSelectionMode_HighestHead     = "HighestHead"
	NodeSelectionMode_RoundRobin      = "RoundRobin"
	NodeSelectionMode_TotalDifficulty = "TotalDifficulty"
	NodeSelectionMode_PriorityLevel   = "PriorityLevel"
)

// NodeSelector represents a strategy to select the next node from the pool.
//
// Deprecated: use [pkg/github.com/smartcontractkit/chainlink/v2/common/client.NodeSelector]
type NodeSelector interface {
	// Select returns a Node, or nil if none can be selected.
	// Implementation must be thread-safe.
	Select() Node
	// Name returns the strategy name, e.g. "HighestHead" or "RoundRobin"
	Name() string
}

// PoolConfig represents settings for the Pool
//
// Deprecated: to be removed
type PoolConfig interface {
	NodeSelectionMode() string
	NodeNoNewHeadsThreshold() time.Duration
	LeaseDuration() time.Duration
}

// Pool represents an abstraction over one or more primary nodes
// It is responsible for liveness checking and balancing queries across live nodes
//
// Deprecated: use [pkg/github.com/smartcontractkit/chainlink/v2/common/client.MultiNode]
type Pool struct {
	services.StateMachine
	nodes               []Node
	sendonlys           []SendOnlyNode
	chainID             *big.Int
	chainType           config.ChainType
	logger              logger.SugaredLogger
	selectionMode       string
	noNewHeadsThreshold time.Duration
	nodeSelector        NodeSelector
	leaseDuration       time.Duration
	leaseTicker         *time.Ticker

	activeMu   sync.RWMutex
	activeNode Node

	chStop services.StopChan
	wg     sync.WaitGroup
}

// NewPool - creates new instance of [Pool]
//
// Deprecated: use [pkg/github.com/smartcontractkit/chainlink/v2/common/client.NewMultiNode]
func NewPool(lggr logger.Logger, selectionMode string, leaseDuration time.Duration, noNewHeadsTreshold time.Duration, nodes []Node, sendonlys []SendOnlyNode, chainID *big.Int, chainType config.ChainType) *Pool {
	if chainID == nil {
		panic("chainID is required")
	}

	nodeSelector := func() NodeSelector {
		switch selectionMode {
		case NodeSelectionMode_HighestHead:
			return NewHighestHeadNodeSelector(nodes)
		case NodeSelectionMode_RoundRobin:
			return NewRoundRobinSelector(nodes)
		case NodeSelectionMode_TotalDifficulty:
			return NewTotalDifficultyNodeSelector(nodes)
		case NodeSelectionMode_PriorityLevel:
			return NewPriorityLevelNodeSelector(nodes)
		default:
			panic(fmt.Sprintf("unsupported NodeSelectionMode: %s", selectionMode))
		}
	}()

	lggr = logger.Named(lggr, "Pool")
	lggr = logger.With(lggr, "evmChainID", chainID.String())

	p := &Pool{
		nodes:               nodes,
		sendonlys:           sendonlys,
		chainID:             chainID,
		chainType:           chainType,
		logger:              logger.Sugared(lggr),
		selectionMode:       selectionMode,
		noNewHeadsThreshold: noNewHeadsTreshold,
		nodeSelector:        nodeSelector,
		chStop:              make(chan struct{}),
		leaseDuration:       leaseDuration,
	}

	p.logger.Debugf("The pool is configured to use NodeSelectionMode: %s", selectionMode)

	return p
}

// Dial starts every node in the pool
//
// Nodes handle their own redialing and runloops, so this function does not
// return any error if the nodes aren't available
func (p *Pool) Dial(ctx context.Context) error {
	return p.StartOnce("Pool", func() (merr error) {
		if len(p.nodes) == 0 {
			return pkgerrors.Errorf("no available nodes for chain %s", p.chainID.String())
		}
		var ms services.MultiStart
		for _, n := range p.nodes {
			if n.ChainID().Cmp(p.chainID) != 0 {
				return ms.CloseBecause(pkgerrors.Errorf("node %s has chain ID %s which does not match pool chain ID of %s", n.String(), n.ChainID().String(), p.chainID.String()))
			}
			rawNode, ok := n.(*node)
			if ok {
				// This is a bit hacky but it allows the node to be aware of
				// pool state and prevent certain state transitions that might
				// otherwise leave no nodes available. It is better to have one
				// node in a degraded state than no nodes at all.
				rawNode.nLiveNodes = p.nLiveNodes
			}
			// node will handle its own redialing and automatic recovery
			if err := ms.Start(ctx, n); err != nil {
				return err
			}
		}
		for _, s := range p.sendonlys {
			if s.ChainID().Cmp(p.chainID) != 0 {
				return ms.CloseBecause(pkgerrors.Errorf("sendonly node %s has chain ID %s which does not match pool chain ID of %s", s.String(), s.ChainID().String(), p.chainID.String()))
			}
			if err := ms.Start(ctx, s); err != nil {
				return err
			}
		}
		p.wg.Add(1)
		go p.runLoop()

		if p.leaseDuration.Seconds() > 0 && p.selectionMode != NodeSelectionMode_RoundRobin {
			p.logger.Infof("The pool will switch to best node every %s", p.leaseDuration.String())
			p.wg.Add(1)
			go p.checkLeaseLoop()
		} else {
			p.logger.Info("Best node switching is disabled")
		}

		return nil
	})
}

// nLiveNodes returns the number of currently alive nodes, as well as the highest block number and greatest total difficulty.
// totalDifficulty will be 0 if all nodes return nil.
func (p *Pool) nLiveNodes() (nLiveNodes int, blockNumber int64, totalDifficulty *big.Int) {
	totalDifficulty = big.NewInt(0)
	for _, n := range p.nodes {
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

func (p *Pool) checkLease() {
	bestNode := p.nodeSelector.Select()
	for _, n := range p.nodes {
		// Terminate client subscriptions. Services are responsible for reconnecting, which will be routed to the new
		// best node. Only terminate connections with more than 1 subscription to account for the aliveLoop subscription
		if n.State() == NodeStateAlive && n != bestNode && n.SubscribersCount() > 1 {
			p.logger.Infof("Switching to best node from %q to %q", n.String(), bestNode.String())
			n.UnsubscribeAllExceptAliveLoop()
		}
	}

	if bestNode != p.activeNode {
		p.activeMu.Lock()
		p.activeNode = bestNode
		p.activeMu.Unlock()
	}
}

func (p *Pool) checkLeaseLoop() {
	defer p.wg.Done()
	p.leaseTicker = time.NewTicker(p.leaseDuration)
	defer p.leaseTicker.Stop()

	for {
		select {
		case <-p.leaseTicker.C:
			p.checkLease()
		case <-p.chStop:
			return
		}
	}
}

func (p *Pool) runLoop() {
	defer p.wg.Done()

	p.report()

	// Prometheus' default interval is 15s, set this to under 7.5s to avoid
	// aliasing (see: https://en.wikipedia.org/wiki/Nyquist_frequency)
	reportInterval := 6500 * time.Millisecond
	monitor := time.NewTicker(utils.WithJitter(reportInterval))
	defer monitor.Stop()

	for {
		select {
		case <-monitor.C:
			p.report()
		case <-p.chStop:
			return
		}
	}
}

func (p *Pool) report() {
	type nodeWithState struct {
		Node  string
		State string
	}

	var total, dead int
	counts := make(map[NodeState]int)
	nodeStates := make([]nodeWithState, len(p.nodes))
	for i, n := range p.nodes {
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
		PromEVMPoolRPCNodeStates.WithLabelValues(p.chainID.String(), state.String()).Set(float64(count))
	}

	live := total - dead
	p.logger.Tracew(fmt.Sprintf("Pool state: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	if total == dead {
		rerr := fmt.Errorf("no EVM primary nodes available: 0/%d nodes are alive", total)
		p.logger.Criticalw(rerr.Error(), "nodeStates", nodeStates)
		p.SvcErrBuffer.Append(rerr)
	} else if dead > 0 {
		p.logger.Errorw(fmt.Sprintf("At least one EVM primary node is dead: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	}
}

// Close tears down the pool and closes all nodes
func (p *Pool) Close() error {
	return p.StopOnce("Pool", func() error {
		close(p.chStop)
		p.wg.Wait()

		return services.CloseAll(services.MultiCloser(p.nodes), services.MultiCloser(p.sendonlys))
	})
}

func (p *Pool) ChainID() *big.Int {
	return p.selectNode().ChainID()
}

func (p *Pool) ChainType() config.ChainType {
	return p.chainType
}

// selectNode returns the active Node, if it is still NodeStateAlive, otherwise it selects a new one from the NodeSelector.
func (p *Pool) selectNode() (node Node) {
	p.activeMu.RLock()
	node = p.activeNode
	p.activeMu.RUnlock()
	if node != nil && node.State() == NodeStateAlive {
		return // still alive
	}

	// select a new one
	p.activeMu.Lock()
	defer p.activeMu.Unlock()
	node = p.activeNode
	if node != nil && node.State() == NodeStateAlive {
		return // another goroutine beat us here
	}

	p.activeNode = p.nodeSelector.Select()

	if p.activeNode == nil {
		p.logger.Criticalw("No live RPC nodes available", "NodeSelectionMode", p.nodeSelector.Name())
		errmsg := fmt.Errorf("no live nodes available for chain %s", p.chainID.String())
		p.SvcErrBuffer.Append(errmsg)
		return &erroringNode{errMsg: errmsg.Error()}
	}

	if p.leaseTicker != nil {
		p.leaseTicker.Reset(p.leaseDuration)
	}
	return p.activeNode
}

func (p *Pool) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return p.selectNode().CallContext(ctx, result, method, args...)
}

func (p *Pool) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return p.selectNode().BatchCallContext(ctx, b)
}

// BatchCallContextAll calls BatchCallContext for every single node including
// sendonlys.
// CAUTION: This should only be used for mass re-transmitting transactions, it
// might have unexpected effects to use it for anything else.
func (p *Pool) BatchCallContextAll(ctx context.Context, b []rpc.BatchElem) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	main := p.selectNode()
	var all []SendOnlyNode
	for _, n := range p.nodes {
		all = append(all, n)
	}
	all = append(all, p.sendonlys...)
	for _, n := range all {
		if n == main {
			// main node is used at the end for the return value
			continue
		}
		// Parallel call made to all other nodes with ignored return value
		wg.Add(1)
		go func(n SendOnlyNode) {
			defer wg.Done()
			err := n.BatchCallContext(ctx, b)
			if err != nil {
				p.logger.Debugw("Secondary node BatchCallContext failed", "err", err)
			} else {
				p.logger.Trace("Secondary node BatchCallContext success")
			}
		}(n)
	}

	return main.BatchCallContext(ctx, b)
}

// SendTransaction wrapped Geth client methods
func (p *Pool) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	main := p.selectNode()
	var all []SendOnlyNode
	for _, n := range p.nodes {
		all = append(all, n)
	}
	all = append(all, p.sendonlys...)
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
		ok := p.IfNotStopped(func() {
			// Must wrap inside IfNotStopped to avoid waitgroup racing with Close
			p.wg.Add(1)
			go func(n SendOnlyNode) {
				defer p.wg.Done()

				sendCtx, cancel := p.chStop.CtxCancel(ContextWithDefaultTimeout())
				defer cancel()

				err := NewSendError(n.SendTransaction(sendCtx, tx))
				p.logger.Debugw("Sendonly node sent transaction", "name", n.String(), "tx", tx, "err", err)
				if err == nil || err.IsNonceTooLowError(nil) || err.IsTransactionAlreadyMined(nil) || err.IsTransactionAlreadyInMempool(nil) {
					// Nonce too low or transaction known errors are expected since
					// the primary SendTransaction may well have succeeded already
					return
				}

				p.logger.Warnw("Eth client returned error", "name", n.String(), "err", err, "tx", tx)
			}(n)
		})
		if !ok {
			p.logger.Debug("Cannot send transaction on sendonly node; pool is stopped", "node", n.String())
		}
	}

	return main.SendTransaction(ctx, tx)
}

func (p *Pool) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return p.selectNode().PendingCodeAt(ctx, account)
}

func (p *Pool) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return p.selectNode().PendingNonceAt(ctx, account)
}

func (p *Pool) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return p.selectNode().NonceAt(ctx, account, blockNumber)
}

func (p *Pool) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return p.selectNode().TransactionReceipt(ctx, txHash)
}

func (p *Pool) TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, error) {
	return p.selectNode().TransactionByHash(ctx, txHash)
}

func (p *Pool) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return p.selectNode().BlockByNumber(ctx, number)
}

func (p *Pool) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return p.selectNode().BlockByHash(ctx, hash)
}

func (p *Pool) BlockNumber(ctx context.Context) (uint64, error) {
	return p.selectNode().BlockNumber(ctx)
}

func (p *Pool) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return p.selectNode().BalanceAt(ctx, account, blockNumber)
}

func (p *Pool) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return p.selectNode().FilterLogs(ctx, q)
}

func (p *Pool) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return p.selectNode().SubscribeFilterLogs(ctx, q, ch)
}

func (p *Pool) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	return p.selectNode().EstimateGas(ctx, call)
}

func (p *Pool) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return p.selectNode().SuggestGasPrice(ctx)
}

func (p *Pool) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return p.selectNode().CallContract(ctx, msg, blockNumber)
}

func (p *Pool) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	return p.selectNode().PendingCallContract(ctx, msg)
}

func (p *Pool) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return p.selectNode().CodeAt(ctx, account, blockNumber)
}

// bind.ContractBackend methods
func (p *Pool) HeaderByNumber(ctx context.Context, n *big.Int) (*types.Header, error) {
	return p.selectNode().HeaderByNumber(ctx, n)
}
func (p *Pool) HeaderByHash(ctx context.Context, h common.Hash) (*types.Header, error) {
	return p.selectNode().HeaderByHash(ctx, h)
}

func (p *Pool) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return p.selectNode().SuggestGasTipCap(ctx)
}

// EthSubscribe implements evmclient.Client
func (p *Pool) EthSubscribe(ctx context.Context, channel chan<- *evmtypes.Head, args ...interface{}) (ethereum.Subscription, error) {
	return p.selectNode().EthSubscribe(ctx, channel, args...)
}
