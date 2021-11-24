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
)

// Pool represents an abstraction over one or more primary nodes
// It is responsible for liveness checking and balancing queries across live nodes
type Pool struct {
	nodes           []Node
	sendonlys       []SendOnlyNode
	chainID         *big.Int
	roundRobinCount atomic.Uint32
	logger          logger.Logger
}

func NewPool(logger logger.Logger, nodes []Node, sendonlys []SendOnlyNode, chainID *big.Int) *Pool {
	if len(nodes) == 0 {
		panic("must provide at least one node")
	}
	p := &Pool{
		nodes:     nodes,
		sendonlys: sendonlys,
		logger:    logger.Named("Pool"),
	}
	if chainID != nil {
		p.initChainID(chainID)
	}
	return p
}

func (p *Pool) initChainID(chainID *big.Int) {
	p.chainID = chainID
	p.logger = p.logger.With("evmChainID", chainID.String())
}

// Dial dials every node in the pool and verifies their chain IDs are consistent.
func (p *Pool) Dial(ctx context.Context) (err error) {
	for _, n := range p.nodes {
		err = multierr.Combine(err, n.Dial(ctx))
	}
	for _, s := range p.sendonlys {
		err = multierr.Combine(err, s.Dial(ctx))
	}
	if err != nil {
		return err
	}
	return p.verifyChainIDs(ctx)
}

// verifyChainIDs checks that every node's chain ID is consistent, initializing from the first node if nil.
func (p *Pool) verifyChainIDs(ctx context.Context) (err error) {
	if p.chainID == nil {
		chainID, err2 := p.nodes[0].ChainID(ctx)
		if err2 != nil {
			return errors.Wrap(err, "failed to get chain ID from first node")
		}
		p.initChainID(chainID)
	}
	for _, n := range p.nodes {
		err = multierr.Combine(err, n.Verify(ctx, p.chainID))
	}
	for _, s := range p.sendonlys {
		err = multierr.Combine(err, s.Verify(ctx, p.chainID))
	}
	return err
}

func (p *Pool) Close() {
	for _, n := range p.nodes {
		n.Close()
	}
}

func (p *Pool) ChainID() *big.Int {
	return p.chainID
}

func (p *Pool) getRoundRobin() Node {
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
