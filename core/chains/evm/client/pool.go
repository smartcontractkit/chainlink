package client

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"go.uber.org/atomic"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
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

// Dial dials every node in the pool and verifies their chain IDs are consistent.
func (p *Pool) Dial(ctx context.Context) error {
	return p.StartOnce("Pool", func() (merr error) {
		if len(p.nodes) == 0 {
			return errors.Errorf("no available nodes for chain %s", p.chainID.String())
		}
		for _, n := range p.nodes {
			if err := n.Dial(ctx); err != nil {
				p.logger.Errorw("Error dialing node", "node", n, "err", err)
			} else if err := n.Verify(ctx, p.chainID); err != nil {
				p.logger.Errorw("Error verifying node", "node", n, "err", err)
			}
		}
		for _, s := range p.sendonlys {
			// TODO: Deal with sendonly nodes state
			err := s.Dial(ctx)
			if err != nil {
				return err
			}
		}
		p.wg.Add(1)
		go p.runLoop()

		return nil
	})
}

// dialRetryInterval controls how often we try to reconnect a dead node
var dialRetryInterval = 5 * time.Second

func (p *Pool) runLoop() {
	defer p.wg.Done()
	ticker := time.NewTicker(dialRetryInterval)

	for {
		select {
		case <-p.chStop:
			return
		case <-ticker.C:
			// re-dial all dead nodes
			func() {
				ctx, cancel := utils.ContextFromChan(p.chStop)
				defer cancel()
				ctx, cancel = context.WithTimeout(ctx, dialRetryInterval)
				defer cancel()
				// TODO: How does this play with automatic WS reconnects?
				p.redialDeadNodes(ctx)
			}()
		}
	}
}

func (p *Pool) redialDeadNodes(ctx context.Context) {
	for _, n := range p.nodes {
		if n.State() == NodeStateDead {
			if err := n.Dial(ctx); err != nil {
				p.logger.Errorw(fmt.Sprintf("Failed to redial eth node: %v", err), "err", err, "node", n.String())
			}
		}
		if n.State() == NodeStateInvalidChainID || n.State() == NodeStateDialed {
			if err := n.Verify(ctx, p.chainID); err != nil {
				p.logger.Errorw(fmt.Sprintf("Failed to verify eth node: %v", err), "err", err, "node", n.String())
			}
		}
	}
}

func (p *Pool) Close() {
	//nolint:errcheck
	p.StopOnce("Pool", func() error {
		close(p.chStop)
		p.wg.Wait()
		for _, n := range p.nodes {
			n.Close()
		}
		return nil
	})
}

func (p *Pool) ChainID() *big.Int {
	return p.chainID
}

func (p *Pool) getRoundRobin() Node {
	nodes := p.liveNodes()
	nNodes := len(nodes)
	if nNodes == 0 {
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

// Wrapped Geth client methods
func (p *Pool) SendTransaction(ctx context.Context, tx *types.Transaction) error {
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
		// Parallel send to all other nodes with ignored return value
		wg.Add(1)
		go func(n SendOnlyNode) {
			defer wg.Done()
			err := NewSendError(n.SendTransaction(ctx, tx))
			if err == nil || err.IsNonceTooLowError() || err.IsTransactionAlreadyInMempool() {
				// Nonce too low or transaction known errors are expected since
				// the primary SendTransaction may well have succeeded already
				return
			}
			p.logger.Warnw("eth client returned error", "name", n.String(), "err", err, "tx", tx)
		}(n)
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

func (p *Pool) EthSubscribe(ctx context.Context, channel interface{}, args ...interface{}) (ethereum.Subscription, error) {
	return p.getRoundRobin().EthSubscribe(ctx, channel, args...)
}
