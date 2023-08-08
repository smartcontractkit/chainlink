package client

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type chainClient struct {
	utils.StartStopOnce
	nodes     []clienttypes.Node[*big.Int, common.Hash, *evmtypes.Head, ethereum.Subscription, RPCClient]
	sendonlys []clienttypes.Node[*big.Int, common.Hash, *evmtypes.Head, ethereum.Subscription, RPCClient]

	multiNodeClient clienttypes.MultiNodeClient[
		*big.Int,
		common.Hash,
		*evmtypes.Head,
		ethereum.Subscription,
		RPCClient,
	]

	chainType config.ChainType
	logger    logger.Logger
	chStop    utils.StopChan
	wg        sync.WaitGroup
}

func (c *chainClient) NewClient(
	logger logger.Logger,
	selectionMode string,
	noNewHeadsThreshold time.Duration,
	nodes []clienttypes.Node[*big.Int, common.Hash, *evmtypes.Head, ethereum.Subscription, RPCClient],
	sendonlys []clienttypes.Node[*big.Int, common.Hash, *evmtypes.Head, ethereum.Subscription, RPCClient],
	chainID *big.Int,
	chainType config.ChainType,
) *chainClient {
	multiNodeClient := clienttypes.NewMultiNodeClient[*big.Int, common.Hash, *evmtypes.Head, ethereum.Subscription, RPCClient](
		logger, selectionMode, noNewHeadsThreshold, nodes, sendonlys, chainID, chainType,
	)

	lggr := logger.Named("Client").With("chainID", chainID.String())

	client := &chainClient{
		nodes:           nodes,
		sendonlys:       sendonlys,
		multiNodeClient: multiNodeClient,
		logger:          lggr,
		chStop:          make(chan struct{}),
	}

	return client
}

func (c *chainClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return c.multiNodeClient.SelectNode().RPCClient().BalanceAt(ctx, account, blockNumber)
}

func (c *chainClient) BatchCallContext(ctx context.Context, b []any) error {
	return c.multiNodeClient.SelectNode().RPCClient().BatchCallContext(ctx, b)
}

func (c *chainClient) BatchCallContextAll(ctx context.Context, b []any) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	main := c.multiNodeClient.SelectNode()
	var all []clienttypes.Node[*big.Int, common.Hash, *evmtypes.Head, ethereum.Subscription, RPCClient]
	all = append(all, c.nodes...)
	all = append(all, c.sendonlys...)
	for _, n := range all {
		if n == main {
			// main node is used at the end for the return value
			continue
		}
		// Parallel call made to all other nodes with ignored return value
		wg.Add(1)
		go func(n clienttypes.Node[*big.Int, common.Hash, *evmtypes.Head, ethereum.Subscription, RPCClient]) {
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

func (c *chainClient) ChainType() config.ChainType {
	return c.chainType
}

func (c *chainClient) BlockByHash(ctx context.Context, hash common.Hash) (*evmtypes.Head, error) {
	return c.multiNodeClient.SelectNode().RPCClient().BlockByHash(ctx, hash)
}

func (c *chainClient) BlockByNumber(ctx context.Context, number *big.Int) (*evmtypes.Head, error) {
	return c.multiNodeClient.SelectNode().RPCClient().BlockByNumber(ctx, number)
}

func (c *chainClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return c.multiNodeClient.SelectNode().RPCClient().CallContext(ctx, result, method, args...)
}

func (c *chainClient) CallContract(
	ctx context.Context,
	attempt interface{},
	blockNumber *big.Int,
) (rpcErr []byte, extractErr error) {
	return c.multiNodeClient.SelectNode().RPCClient().CallContract(ctx, attempt, blockNumber)
}

func (c *chainClient) ChainID() (*big.Int, error) {
	return c.multiNodeClient.SelectNode().RPCClient().ChainID()
}

func (c *chainClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return c.multiNodeClient.SelectNode().RPCClient().CodeAt(ctx, account, blockNumber)
}

func (c *chainClient) ConfiguredChainID() *big.Int {
	return c.multiNodeClient.SelectNode().RPCClient().ConfiguredChainID()
}

func (c *chainClient) EstimateGas(ctx context.Context, call any) (gas uint64, err error) {
	return c.multiNodeClient.SelectNode().RPCClient().EstimateGas(ctx, call)
}

func (c *chainClient) FilterEvents(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return c.multiNodeClient.SelectNode().RPCClient().FilterEvents(ctx, query)
}

func (c *chainClient) IsL2() bool {
	return c.ChainType().IsL2()
}

func (c *chainClient) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
	return c.multiNodeClient.SelectNode().RPCClient().LatestBlockHeight(ctx)
}

func (c *chainClient) LINKBalance(ctx context.Context, accountAddress common.Address, linkAddress common.Address) (*assets.Link, error) {
	return c.multiNodeClient.SelectNode().RPCClient().LINKBalance(ctx, accountAddress, linkAddress)
}

func (c *chainClient) PendingSequenceAt(ctx context.Context, addr common.Address) (evmtypes.Nonce, error) {
	return c.multiNodeClient.SelectNode().RPCClient().PendingSequenceAt(ctx, addr)
}

func (c *chainClient) SendEmptyTransaction(
	ctx context.Context,
	newTxAttempt func(seq evmtypes.Nonce, feeLimit uint32, fee *assets.Wei, fromAddress common.Address) (attempt any, err error),
	seq evmtypes.Nonce,
	gasLimit uint32,
	fee *assets.Wei,
	fromAddress common.Address,
) (txhash string, err error) {
	return c.multiNodeClient.SelectNode().RPCClient().SendEmptyTransaction(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
}

func (c *chainClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	main := c.multiNodeClient.SelectNode()
	var all []clienttypes.Node[*big.Int, common.Hash, *evmtypes.Head, ethereum.Subscription, RPCClient]
	all = append(all, c.nodes...)
	all = append(all, c.sendonlys...)
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
		ok := c.IfNotStopped(func() {
			// Must wrap inside IfNotStopped to avoid waitgroup racing with Close
			c.wg.Add(1)
			go func(n clienttypes.Node[*big.Int, common.Hash, *evmtypes.Head, ethereum.Subscription, RPCClient]) {
				defer c.wg.Done()

				sendCtx, cancel := c.chStop.CtxCancel(ContextWithDefaultTimeout())
				defer cancel()
				err, _ := n.RPCClient().SendTransactionReturnCode(sendCtx, tx)
				c.logger.Debugw("Sendonly node sent transaction", "name", n.String(), "tx", tx, "err", err)
				if err == clienttypes.TransactionAlreadyKnown || err == clienttypes.Successful {
					// Nonce too low or transaction known errors are expected since
					// the primary SendTransaction may well have succeeded already
					return
				}
				c.logger.Warnw("Eth client returned error", "name", n.String(), "err", err, "tx", tx)
			}(n)
		})
		if !ok {
			c.logger.Debug("Cannot send transaction on sendonly node; pool is stopped", "node", n.String())
		}
	}

	return main.RPCClient().SendTransaction(ctx, tx)
}

func (c *chainClient) SendTransactionReturnCode(
	ctx context.Context,
	tx *types.Transaction,
) (clienttypes.SendTxReturnCode, error) {
	return c.multiNodeClient.SelectNode().RPCClient().SendTransactionReturnCode(ctx, tx)
}

func (c *chainClient) SequenceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (evmtypes.Nonce, error) {
	return c.multiNodeClient.SelectNode().RPCClient().SequenceAt(ctx, account, blockNumber)
}

func (c *chainClient) SimulateTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.multiNodeClient.SelectNode().RPCClient().SimulateTransaction(ctx, tx)
}

func (c *chainClient) Subscribe(ctx context.Context, channel chan<- *evmtypes.Head, args ...interface{}) (ethereum.Subscription, error) {
	return c.multiNodeClient.SelectNode().RPCClient().Subscribe(ctx, channel, args)
}

func (c *chainClient) TokenBalance(ctx context.Context, account common.Address, tokenAddr common.Address) (*big.Int, error) {
	return c.multiNodeClient.SelectNode().RPCClient().TokenBalance(ctx, account, tokenAddr)
}

func (c *chainClient) TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, error) {
	return c.multiNodeClient.SelectNode().RPCClient().TransactionByHash(ctx, txHash)
}

func (c *chainClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*evmtypes.Receipt, error) {
	return c.multiNodeClient.SelectNode().RPCClient().TransactionReceipt(ctx, txHash)
}
