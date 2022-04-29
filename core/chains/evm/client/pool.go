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
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/atomic"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	// PromEVMPoolRPCNodeStates reports current RPC node state
	PromEVMPoolRPCNodeStates = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "evm_pool_rpc_node_states",
		Help: "The number of RPC nodes currently in the given state for the given chain",
	}, []string{"evmChainID", "state"})
)

// Pool represents an abstraction over one or more primary nodes
// It is responsible for liveness checking and balancing queries across live nodes
type Pool struct {
	utils.StartStopOnce
	nodes           []Node
	sendonlys       []SendOnlyNode
	chainID         *big.Int
	roundRobinCount atomic.Uint32
	logger          logger.Logger

	chStop chan struct{}
	wg     sync.WaitGroup
}

func NewPool(logger logger.Logger, nodes []Node, sendonlys []SendOnlyNode, chainID *big.Int) *Pool {
	if chainID == nil {
		panic("chainID is required")
	}
	p := &Pool{
		utils.StartStopOnce{},
		nodes,
		sendonlys,
		chainID,
		atomic.Uint32{},
		logger.Named("Pool").With("evmChainID", chainID.String()),
		make(chan struct{}),
		sync.WaitGroup{},
	}
	return p
}

// Dial starts every node in the pool
func (p *Pool) Dial(ctx context.Context) error {
	return p.StartOnce("Pool", func() (merr error) {
		if len(p.nodes) == 0 {
			return errors.Errorf("no available nodes for chain %s", p.chainID.String())
		}
		for _, n := range p.nodes {
			if n.ChainID().Cmp(p.chainID) != 0 {
				return errors.Errorf("node %s has chain ID %s which does not match pool chain ID of %s", n.String(), n.ChainID().String(), p.chainID.String())
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
			if err := n.Start(ctx); err != nil {
				return err
			}
		}
		for _, s := range p.sendonlys {
			if s.ChainID().Cmp(p.chainID) != 0 {
				return errors.Errorf("sendonly node %s has chain ID %s which does not match pool chain ID of %s", s.String(), s.ChainID().String(), p.chainID.String())
			}
			err := s.Start(ctx)
			if err != nil {
				return err
			}
		}
		p.wg.Add(1)
		go p.runLoop()

		return nil
	})
}

// nLiveNodes returns the number of currently alive nodes
func (p *Pool) nLiveNodes() (nLiveNodes int) {
	for _, n := range p.nodes {
		if n.State() == NodeStateAlive {
			nLiveNodes++
		}
	}
	return
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
		p.logger.Criticalw(fmt.Sprintf("No EVM primary nodes available: 0/%d nodes are alive", total), "nodeStates", nodeStates)
	} else if dead > 0 {
		p.logger.Errorw(fmt.Sprintf("At least one EVM primary node is dead: %d/%d nodes are alive", live, total), "nodeStates", nodeStates)
	}
}

// Close tears down the pool and closes all nodes
func (p *Pool) Close() {
	err := p.StopOnce("Pool", func() error {
		close(p.chStop)
		p.wg.Wait()

		var closeWg sync.WaitGroup
		closeWg.Add(len(p.nodes))
		for _, n := range p.nodes {
			go func(node Node) {
				defer closeWg.Done()
				node.Close()
			}(n)
		}
		closeWg.Add(len(p.sendonlys))
		for _, s := range p.sendonlys {
			go func(sNode SendOnlyNode) {
				defer closeWg.Done()
				sNode.Close()
			}(s)
		}
		closeWg.Wait()
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (p *Pool) ChainID() *big.Int {
	return p.chainID
}

func (p *Pool) getRoundRobin() Node {
	nodes := p.liveNodes()
	nNodes := len(nodes)
	if nNodes == 0 {
		p.logger.Critical("No live RPC nodes available")
		return &erroringNode{errMsg: fmt.Sprintf("no live nodes available for chain %s", p.chainID.String())}
	}

	// NOTE: Inc returns the number after addition, so we must -1 to get the "current" counter
	count := p.roundRobinCount.Inc() - 1
	idx := int(count % uint32(nNodes))

	return nodes[idx]
}

func (p *Pool) liveNodes() (liveNodes []Node) {
	for _, n := range p.nodes {
		if n.State() == NodeStateAlive {
			liveNodes = append(liveNodes, n)
		}
	}
	return
}

func (p *Pool) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return p.getRoundRobin().CallContext(ctx, result, method, args...)
}

func (p *Pool) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return p.getRoundRobin().BatchCallContext(ctx, b)
}

// BatchCallContextAll calls BatchCallContext for every single node including
// sendonlys.
// CAUTION: This should only be used for mass re-transmitting transactions, it
// might have unexpected effects to use it for anything else.
func (p *Pool) BatchCallContextAll(ctx context.Context, b []rpc.BatchElem) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	main := p.getRoundRobin()
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

// Wrapped Geth client methods
func (p *Pool) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	main := p.getRoundRobin()
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
			go func(n SendOnlyNode, txCp types.Transaction) {
				defer p.wg.Done()
				timeoutCtx, cancel := DefaultQueryCtx()
				defer cancel()
				sendCtx, cancel2 := utils.WithCloseChan(timeoutCtx, p.chStop)
				defer cancel2()
				err := NewSendError(n.SendTransaction(sendCtx, &txCp))
				p.logger.Debugw("Sendonly node sent transaction", "name", n.String(), "tx", tx, "err", err)
				if err == nil || err.IsNonceTooLowError() || err.IsTransactionAlreadyMined() || err.IsTransactionAlreadyInMempool() {
					// Nonce too low or transaction known errors are expected since
					// the primary SendTransaction may well have succeeded already
					return
				}

				p.logger.Warnw("Eth client returned error", "name", n.String(), "err", err, "tx", tx)
			}(n, *tx) // copy tx here in case it is mutated after the function returns
		})
		if !ok {
			p.logger.Debug("Cannot send transaction on sendonly node; pool is stopped", "node", n.String())
		}
	}

	return main.SendTransaction(ctx, tx)
}

func (p *Pool) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return p.getRoundRobin().PendingCodeAt(ctx, account)
}

func (p *Pool) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return p.getRoundRobin().PendingNonceAt(ctx, account)
}

func (p *Pool) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return p.getRoundRobin().NonceAt(ctx, account, blockNumber)
}

func (p *Pool) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return p.getRoundRobin().TransactionReceipt(ctx, txHash)
}

func (p *Pool) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return p.getRoundRobin().BlockByNumber(ctx, number)
}

func (p *Pool) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return p.getRoundRobin().BlockByHash(ctx, hash)
}

func (p *Pool) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return p.getRoundRobin().BalanceAt(ctx, account, blockNumber)
}

func (p *Pool) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return p.getRoundRobin().FilterLogs(ctx, q)
}

func (p *Pool) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return p.getRoundRobin().SubscribeFilterLogs(ctx, q, ch)
}

func (p *Pool) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	return p.getRoundRobin().EstimateGas(ctx, call)
}

func (p *Pool) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return p.getRoundRobin().SuggestGasPrice(ctx)
}

func (p *Pool) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return p.getRoundRobin().CallContract(ctx, msg, blockNumber)
}

func (p *Pool) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return p.getRoundRobin().CodeAt(ctx, account, blockNumber)
}

// bind.ContractBackend methods
func (p *Pool) HeaderByNumber(ctx context.Context, n *big.Int) (*types.Header, error) {
	return p.getRoundRobin().HeaderByNumber(ctx, n)
}

func (p *Pool) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return p.getRoundRobin().SuggestGasTipCap(ctx)
}

// EthSubscribe implements evmclient.Client
func (p *Pool) EthSubscribe(ctx context.Context, channel chan<- *evmtypes.Head, args ...interface{}) (ethereum.Subscription, error) {
	return p.getRoundRobin().EthSubscribe(ctx, channel, args...)
}
