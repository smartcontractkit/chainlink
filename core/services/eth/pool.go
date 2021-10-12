package eth

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"go.uber.org/multierr"

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

	nodesMu sync.RWMutex
}

func NewPool(logger logger.Logger, nodes []Node, sendonlys []SendOnlyNode, chainID *big.Int) *Pool {
	return &Pool{utils.StartStopOnce{}, nodes, sendonlys, chainID, atomic.Uint32{}, logger, sync.RWMutex{}}
}

func (p *Pool) AddNode(ctx context.Context, n Node) (err error) {
	if n.ChainID() == nil || (n.ChainID().Cmp(p.chainID) != 0) {
		return errors.Errorf("cannot add node with chain ID %s to pool with chain ID %s", n.ChainID().String(), p.chainID.String())
	}
	ok := p.IfStarted(func() {
		if err = n.Dial(ctx); err != nil {
			err = errors.Wrap(err, "AddNode: failed to dial node")
			return
		}

		p.nodesMu.Lock()
		defer p.nodesMu.Unlock()
		if p.hasNodeWithName(n.Name()) {
			n.Close()
			err = errors.Errorf("node already exists with name %s", n.Name())
			return
		}
		p.nodes = append(p.nodes, n)
	})
	if !ok {
		return errors.New("cannot add node; pool is not started")
	}
	return err
}

func (p *Pool) AddSendOnlyNode(ctx context.Context, n SendOnlyNode) (err error) {
	if n.ChainID() == nil || (n.ChainID().Cmp(p.chainID) != 0) {
		return errors.Errorf("cannot add send only node with chain ID %s to pool with chain ID %s", n.ChainID().String(), p.chainID.String())
	}
	ok := p.IfStarted(func() {
		if err = n.Dial(ctx); err != nil {
			err = errors.Wrap(err, "AddNode: failed to dial node")
			return
		}

		p.nodesMu.Lock()
		defer p.nodesMu.Unlock()
		if p.hasNodeWithName(n.Name()) {
			err = errors.Errorf("node already exists with name %s", n.Name())
			return
		}
		p.sendonlys = append(p.sendonlys, n)
	})
	if !ok {
		return errors.New("cannot add send only node; pool is not started")
	}
	return err
}

func (p *Pool) hasNodeWithName(s string) bool {
	for _, n := range p.nodes {
		if s == n.Name() {
			return true
		}
	}
	for _, n := range p.sendonlys {
		if s == n.Name() {
			return true
		}
	}
	return false
}

func (p *Pool) Dial(ctx context.Context) (err error) {
	return p.StartOnce("Pool", func() (merr error) {
		p.nodesMu.Lock()
		defer p.nodesMu.Unlock()

		for _, n := range p.nodes {
			err = multierr.Combine(err, n.Dial(ctx))
		}
		for _, s := range p.sendonlys {
			err = multierr.Combine(err, s.Dial(ctx))
		}
		return err
	})
}

func (p *Pool) Close() error {
	return p.StopOnce("Pool", func() (merr error) {
		p.nodesMu.Lock()
		defer p.nodesMu.Unlock()

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
	p.nodesMu.RLock()
	defer p.nodesMu.RUnlock()

	nNodes := len(p.nodes)
	if nNodes == 0 {
		return &erroringNode{errMsg: fmt.Sprintf("no nodes available for chain %s", p.chainID.String())}
	}

	// NOTE: Inc returns the number after addition, so we must -1 to get the "current" counter
	count := p.roundRobinCount.Inc() - 1
	idx := int(count % uint32(nNodes))

	return p.nodes[idx]
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

	p.nodesMu.RLock()
	defer p.nodesMu.RUnlock()

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
