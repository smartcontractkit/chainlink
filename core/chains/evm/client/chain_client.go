package client

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/client"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

var _ Client = (*chainClient)(nil)

// TODO-1663: rename this to client, once the client.go file is deprecated.
type chainClient struct {
	multiNode commonclient.MultiNode[
		*big.Int,
		evmtypes.Nonce,
		common.Address,
		common.Hash,
		*types.Transaction,
		common.Hash,
		types.Log,
		ethereum.FilterQuery,
		*evmtypes.Receipt,
		*assets.Wei,
		*evmtypes.Head,
		RPCCLient,
	]
	logger logger.Logger
}

func NewChainClient(
	logger logger.Logger,
	selectionMode string,
	leaseDuration time.Duration,
	noNewHeadsThreshold time.Duration,
	nodes []commonclient.Node[*big.Int, *evmtypes.Head, RPCCLient],
	sendonlys []commonclient.SendOnlyNode[*big.Int, RPCCLient],
	chainID *big.Int,
	chainType config.ChainType,
) Client {
	multiNode := commonclient.NewMultiNode[
		*big.Int,
		evmtypes.Nonce,
		common.Address,
		common.Hash,
		*types.Transaction,
		common.Hash,
		types.Log,
		ethereum.FilterQuery,
		*evmtypes.Receipt,
		*assets.Wei,
		*evmtypes.Head,
		RPCCLient,
	](
		logger,
		selectionMode,
		leaseDuration,
		noNewHeadsThreshold,
		nodes,
		sendonlys,
		chainID,
		chainType,
		"EVM",
		ClassifySendOnlyError,
	)
	return &chainClient{
		multiNode: multiNode,
		logger:    logger,
	}
}

func (c *chainClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return c.multiNode.BalanceAt(ctx, account, blockNumber)
}

func (c *chainClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	batch := make([]any, len(b))
	for i, arg := range b {
		batch[i] = any(arg)
	}
	return c.multiNode.BatchCallContext(ctx, batch)
}

func (c *chainClient) BatchCallContextAll(ctx context.Context, b []rpc.BatchElem) error {
	batch := make([]any, len(b))
	for i, arg := range b {
		batch[i] = any(arg)
	}
	return c.multiNode.BatchCallContextAll(ctx, batch)
}

// TODO-1663: return custom Block type instead of geth's once client.go is deprecated.
func (c *chainClient) BlockByHash(ctx context.Context, hash common.Hash) (b *types.Block, err error) {
	rpc, err := c.multiNode.SelectNodeRPC()
	if err != nil {
		return b, err
	}
	return rpc.BlockByHashGeth(ctx, hash)
}

// TODO-1663: return custom Block type instead of geth's once client.go is deprecated.
func (c *chainClient) BlockByNumber(ctx context.Context, number *big.Int) (b *types.Block, err error) {
	rpc, err := c.multiNode.SelectNodeRPC()
	if err != nil {
		return b, err
	}
	return rpc.BlockByNumberGeth(ctx, number)
}

func (c *chainClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return c.multiNode.CallContext(ctx, result, method)
}

func (c *chainClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return c.multiNode.CallContract(ctx, msg, blockNumber)
}

// TODO-1663: change this to actual ChainID() call once client.go is deprecated.
func (c *chainClient) ChainID() (*big.Int, error) {
	//return c.multiNode.ChainID(ctx), nil
	return c.multiNode.ConfiguredChainID(), nil
}

func (c *chainClient) Close() {
	c.multiNode.Close()
}

func (c *chainClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return c.multiNode.CodeAt(ctx, account, blockNumber)
}

func (c *chainClient) ConfiguredChainID() *big.Int {
	return c.multiNode.ConfiguredChainID()
}

func (c *chainClient) Dial(ctx context.Context) error {
	return c.multiNode.Dial(ctx)
}

func (c *chainClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	return c.multiNode.EstimateGas(ctx, call)
}
func (c *chainClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return c.multiNode.FilterEvents(ctx, q)
}

func (c *chainClient) HeaderByHash(ctx context.Context, h common.Hash) (head *types.Header, err error) {
	rpc, err := c.multiNode.SelectNodeRPC()
	if err != nil {
		return head, err
	}
	return rpc.HeaderByHash(ctx, h)
}

func (c *chainClient) HeaderByNumber(ctx context.Context, n *big.Int) (head *types.Header, err error) {
	rpc, err := c.multiNode.SelectNodeRPC()
	if err != nil {
		return head, err
	}
	return rpc.HeaderByNumber(ctx, n)
}

func (c *chainClient) HeadByHash(ctx context.Context, h common.Hash) (*evmtypes.Head, error) {
	return c.multiNode.BlockByHash(ctx, h)
}

func (c *chainClient) HeadByNumber(ctx context.Context, n *big.Int) (*evmtypes.Head, error) {
	return c.multiNode.BlockByNumber(ctx, n)
}

func (c *chainClient) IsL2() bool {
	return c.multiNode.IsL2()
}

func (c *chainClient) LINKBalance(ctx context.Context, address common.Address, linkAddress common.Address) (*assets.Link, error) {
	return c.multiNode.LINKBalance(ctx, address, linkAddress)
}

func (c *chainClient) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
	return c.multiNode.LatestBlockHeight(ctx)
}

func (c *chainClient) NodeStates() map[string]string {
	return c.multiNode.NodeStates()
}

func (c *chainClient) PendingCodeAt(ctx context.Context, account common.Address) (b []byte, err error) {
	rpc, err := c.multiNode.SelectNodeRPC()
	if err != nil {
		return b, err
	}
	return rpc.PendingCodeAt(ctx, account)
}

// TODO-1663: change this to evmtypes.Nonce(int64) once client.go is deprecated.
func (c *chainClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	n, err := c.multiNode.PendingSequenceAt(ctx, account)
	return uint64(n), err
}

func (c *chainClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.multiNode.SendTransaction(ctx, tx)
}

func (c *chainClient) SendTransactionReturnCode(ctx context.Context, tx *types.Transaction, fromAddress common.Address) (commonclient.SendTxReturnCode, error) {
	err := c.SendTransaction(ctx, tx)
	return ClassifySendError(err, c.logger, tx, fromAddress, c.IsL2())
}

func (c *chainClient) SequenceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (evmtypes.Nonce, error) {
	return c.multiNode.SequenceAt(ctx, account, blockNumber)
}

func (c *chainClient) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (s ethereum.Subscription, err error) {
	rpc, err := c.multiNode.SelectNodeRPC()
	if err != nil {
		return s, err
	}
	return rpc.SubscribeFilterLogs(ctx, q, ch)
}

func (c *chainClient) SubscribeNewHead(ctx context.Context, ch chan<- *evmtypes.Head) (ethereum.Subscription, error) {
	csf := newChainIDSubForwarder(c.ConfiguredChainID(), ch)
	err := csf.start(c.multiNode.Subscribe(ctx, csf.srcCh, "newHeads"))
	if err != nil {
		return nil, err
	}
	return csf, nil
}

func (c *chainClient) SuggestGasPrice(ctx context.Context) (p *big.Int, err error) {
	rpc, err := c.multiNode.SelectNodeRPC()
	if err != nil {
		return p, err
	}
	return rpc.SuggestGasPrice(ctx)
}

func (c *chainClient) SuggestGasTipCap(ctx context.Context) (t *big.Int, err error) {
	rpc, err := c.multiNode.SelectNodeRPC()
	if err != nil {
		return t, err
	}
	return rpc.SuggestGasTipCap(ctx)
}

func (c *chainClient) TokenBalance(ctx context.Context, address common.Address, contractAddress common.Address) (*big.Int, error) {
	return c.multiNode.TokenBalance(ctx, address, contractAddress)
}

func (c *chainClient) TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, error) {
	return c.multiNode.TransactionByHash(ctx, txHash)
}

// TODO-1663: return custom Receipt type instead of geth's once client.go is deprecated.
func (c *chainClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (r *types.Receipt, err error) {
	rpc, err := c.multiNode.SelectNodeRPC()
	if err != nil {
		return r, err
	}
	//return rpc.TransactionReceipt(ctx, txHash)
	return rpc.TransactionReceiptGeth(ctx, txHash)
}
