package memory

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
)

// Backend is a wrapper struct which implements
// OnchainClient but also exposes backend methods.
type Backend struct {
	mu  sync.Mutex
	Sim *simulated.Backend
}

func (b *Backend) Commit() common.Hash {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Sim.Commit()
}

func (b *Backend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return b.Sim.Client().CodeAt(ctx, contract, blockNumber)
}

func (b *Backend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return b.Sim.Client().CallContract(ctx, call, blockNumber)
}

func (b *Backend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	return b.Sim.Client().EstimateGas(ctx, call)
}

func (b *Backend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return b.Sim.Client().SuggestGasPrice(ctx)
}

func (b *Backend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return b.Sim.Client().SuggestGasTipCap(ctx)
}

func (b *Backend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return b.Sim.Client().SendTransaction(ctx, tx)
}

func (b *Backend) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return b.Sim.Client().HeaderByNumber(ctx, number)
}

func (b *Backend) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return b.Sim.Client().PendingCodeAt(ctx, account)
}

func (b *Backend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return b.Sim.Client().PendingNonceAt(ctx, account)
}

func (b *Backend) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return b.Sim.Client().FilterLogs(ctx, q)
}

func (b *Backend) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return b.Sim.Client().SubscribeFilterLogs(ctx, q, ch)
}

func (b *Backend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return b.Sim.Client().TransactionReceipt(ctx, txHash)
}

func (b *Backend) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return b.Sim.Client().BalanceAt(ctx, account, blockNumber)
}

func (b *Backend) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return b.Sim.Client().NonceAt(ctx, account, blockNumber)
}

func NewBackend(sim *simulated.Backend) *Backend {
	if sim == nil {
		panic("simulated backend is nil")
	}
	return &Backend{
		Sim: sim,
	}
}
