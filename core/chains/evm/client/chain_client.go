package client

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	commonclient "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type ChainClient struct {
	multiNodeClient commonclient.MultiNodeClient[
		*big.Int,
		*evmtypes.Head,
		RPCClient,
		*types.Transaction,
	]

	chainType config.ChainType
	logger    logger.Logger
}

func (c *ChainClient) NewClient(
	logger logger.Logger,
	selectionMode string,
	noNewHeadsThreshold time.Duration,
	nodes []commonclient.Node[*big.Int, *evmtypes.Head, RPCClient],
	sendonlys []commonclient.SendOnlyNode[*big.Int, RPCClient],
	chainID *big.Int,
	chainType config.ChainType,
) *ChainClient {
	multiNodeClient := commonclient.NewMultiNodeClient[*big.Int, *evmtypes.Head, RPCClient, *types.Transaction](
		logger, selectionMode, noNewHeadsThreshold, nodes, sendonlys, chainID, "EVM",
	)

	lggr := logger.Named("Client").With("chainID", chainID.String())

	client := &ChainClient{
		multiNodeClient: multiNodeClient,
		logger:          lggr,
	}

	return client
}

func (c *ChainClient) Dial(ctx context.Context) error {
	if err := c.multiNodeClient.Dial(ctx); err != nil {
		return errors.Wrap(err, "failed to dial multiNodeClient")
	}
	return nil
}

func (c *ChainClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return c.multiNodeClient.SelectNode().RPCClient().BalanceAt(ctx, account, blockNumber)
}

func (c *ChainClient) BatchCallContext(ctx context.Context, b []any) error {
	return c.multiNodeClient.SelectNode().RPCClient().BatchCallContext(ctx, b)
}

func (c *ChainClient) BatchCallContextAll(ctx context.Context, b []any) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	main := c.multiNodeClient.SelectNode()
	for _, n := range c.multiNodeClient.NodesAsSendOnlys() {
		if n == main {
			// main node is used at the end for the return value
			continue
		}
		// Parallel call made to all other nodes with ignored return value
		wg.Add(1)
		go func(n commonclient.SendOnlyNode[*big.Int, RPCClient]) {
			defer wg.Done()
			err := n.RPCClient().BatchCallContext(ctx, b)
			if err != nil {
				c.logger.Debugw("Secondary node BatchCallContext failed", "err", err)
			} else {
				c.logger.Trace("Secondary node BatchCallContext success")
			}
		}(n)
	}

	return main.RPCClient().BatchCallContext(ctx, b)
}

func (c *ChainClient) ChainType() config.ChainType {
	return c.chainType
}

func (c *ChainClient) Close() {
	c.multiNodeClient.Close()
}

func (c *ChainClient) BlockByHash(ctx context.Context, hash common.Hash) (*evmtypes.Head, error) {
	return c.multiNodeClient.SelectNode().RPCClient().BlockByHash(ctx, hash)
}

func (c *ChainClient) BlockByNumber(ctx context.Context, number *big.Int) (*evmtypes.Head, error) {
	return c.multiNodeClient.SelectNode().RPCClient().BlockByNumber(ctx, number)
}

func (c *ChainClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return c.multiNodeClient.SelectNode().RPCClient().CallContext(ctx, result, method, args...)
}

func (c *ChainClient) CallContract(
	ctx context.Context,
	attempt interface{},
	blockNumber *big.Int,
) (rpcErr []byte, extractErr error) {
	return c.multiNodeClient.SelectNode().RPCClient().CallContract(ctx, attempt, blockNumber)
}

// ChainID makes a direct RPC call. In most cases it should be better to use the configured chain id instead by
// calling ConfiguredChainID.
func (c *ChainClient) ChainID(ctx context.Context) (*big.Int, error) {
	return c.multiNodeClient.SelectNode().RPCClient().ChainID(ctx)
}

func (c *ChainClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return c.multiNodeClient.SelectNode().RPCClient().CodeAt(ctx, account, blockNumber)
}

func (c *ChainClient) ConfiguredChainID() *big.Int {
	return c.multiNodeClient.SelectNode().ConfiguredChainID()
}

func (c *ChainClient) EstimateGas(ctx context.Context, call any) (gas uint64, err error) {
	return c.multiNodeClient.SelectNode().RPCClient().EstimateGas(ctx, call)
}

func (c *ChainClient) FilterEvents(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return c.multiNodeClient.SelectNode().RPCClient().FilterEvents(ctx, query)
}

func (c *ChainClient) IsL2() bool {
	return c.ChainType().IsL2()
}

func (c *ChainClient) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
	return c.multiNodeClient.SelectNode().RPCClient().LatestBlockHeight(ctx)
}

func (c *ChainClient) LINKBalance(ctx context.Context, accountAddress common.Address, linkAddress common.Address) (*assets.Link, error) {
	return c.multiNodeClient.SelectNode().RPCClient().LINKBalance(ctx, accountAddress, linkAddress)
}

func (c *ChainClient) PendingSequenceAt(ctx context.Context, addr common.Address) (evmtypes.Nonce, error) {
	return c.multiNodeClient.SelectNode().RPCClient().PendingSequenceAt(ctx, addr)
}

func (c *ChainClient) SendEmptyTransaction(
	ctx context.Context,
	newTxAttempt func(seq evmtypes.Nonce, feeLimit uint32, fee *assets.Wei, fromAddress common.Address) (attempt any, err error),
	seq evmtypes.Nonce,
	gasLimit uint32,
	fee *assets.Wei,
	fromAddress common.Address,
) (txhash string, err error) {
	return c.multiNodeClient.SelectNode().RPCClient().SendEmptyTransaction(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
}

func (c *ChainClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	main := c.multiNodeClient.SelectNode()
	for _, n := range c.multiNodeClient.NodesAsSendOnlys() {
		if n == main {
			// main node is used at the end for the return value
			continue
		}
		// Parallel send to all other nodes with ignored return value
		// Async - we do not want to block the main thread with secondary nodes
		// in case they are unreliable/slow.
		// It is purely a "best effort" send.
		// Resource is not unbounded because the default context has a timeout.
		c.multiNodeClient.WrapSendOnlyTransaction(ctx, c.logger, tx, n, sendOnlyTransaction)
	}

	return main.RPCClient().SendTransaction(ctx, tx)
}

func sendOnlyTransaction(ctx context.Context, lggr logger.Logger, tx *types.Transaction, n commonclient.SendOnlyNode[*big.Int, RPCClient]) {
	err := n.RPCClient().SendTransaction(ctx, tx)
	lggr.Debugw("Sendonly node sent transaction", "name", n.String(), "tx", tx, "err", err)
	sendOnlyError, err := NewSendOnlyErrorReturnCode(err)
	if sendOnlyError != commonclient.Successful {
		lggr.Warnw("Eth client returned error", "name", n.String(), "err", err, "tx", tx)
	}
}

func (c *ChainClient) SendTransactionReturnCode(ctx context.Context, tx *types.Transaction, fromAddress common.Address) (commonclient.SendTxReturnCode, error) {
	err := c.SendTransaction(ctx, tx)
	return NewSendErrorReturnCode(err, c.logger, tx, fromAddress, c.ChainType().IsL2())
}

func (c *ChainClient) SequenceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (evmtypes.Nonce, error) {
	return c.multiNodeClient.SelectNode().RPCClient().SequenceAt(ctx, account, blockNumber)
}

func (c *ChainClient) SimulateTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.multiNodeClient.SelectNode().RPCClient().SimulateTransaction(ctx, tx)
}

func (c *ChainClient) Subscribe(ctx context.Context, channel chan<- *evmtypes.Head, args ...interface{}) (ethereum.Subscription, error) {
	return c.multiNodeClient.SelectNode().RPCClient().Subscribe(ctx, channel, args)
}

func (c *ChainClient) TokenBalance(ctx context.Context, account common.Address, tokenAddr common.Address) (*big.Int, error) {
	return c.multiNodeClient.SelectNode().RPCClient().TokenBalance(ctx, account, tokenAddr)
}

func (c *ChainClient) TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, error) {
	return c.multiNodeClient.SelectNode().RPCClient().TransactionByHash(ctx, txHash)
}

func (c *ChainClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*evmtypes.Receipt, error) {
	return c.multiNodeClient.SelectNode().RPCClient().TransactionReceipt(ctx, txHash)
}
