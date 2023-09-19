package client

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Client is a generalized client interface that includes methods to interact with different chains.
type Client[
	CHAIN_ID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK_HASH types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	EVENT any,
	EVENT_OPS any,
	TX_RECEIPT types.Receipt[TX_HASH, BLOCK_HASH],
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],
	RPC_CLIENT RPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD],
] interface {
	ClientAPI[
		CHAIN_ID,
		SEQ,
		ADDR,
		BLOCK_HASH,
		TX,
		TX_HASH,
		EVENT,
		EVENT_OPS,
		TX_RECEIPT,
		FEE,
		HEAD,
	]
	BatchCallContextAll(ctx context.Context, b []any) error
	NodeRPC() (rpc RPC_CLIENT, err error)
	ConfiguredChainID() CHAIN_ID
	IsL2() bool
	NodeStates() map[string]string
}

type client[
	CHAIN_ID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK_HASH types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	EVENT any,
	EVENT_OPS any,
	TX_RECEIPT types.Receipt[TX_HASH, BLOCK_HASH],
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],
	RPC_CLIENT RPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD],
] struct {
	multiNode MultiNode[
		CHAIN_ID,
		HEAD,
		RPC_CLIENT,
		TX,
	]

	sendOnlyErrorParser func(err error) (SendTxReturnCode, error)
	chainType           config.ChainType
	logger              logger.Logger
}

func NewClient[
	CHAIN_ID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK_HASH types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	EVENT any,
	EVENT_OPS any,
	TX_RECEIPT types.Receipt[TX_HASH, BLOCK_HASH],
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],
	RPC_CLIENT RPC[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD],
](
	logger logger.Logger,
	selectionMode string,
	noNewHeadsThreshold time.Duration,
	nodes []Node[CHAIN_ID, HEAD, RPC_CLIENT],
	sendonlys []SendOnlyNode[CHAIN_ID, RPC_CLIENT],
	chainID CHAIN_ID,
	chainType config.ChainType,
	sendOnlyErrorParser func(err error) (SendTxReturnCode, error),
) Client[
	CHAIN_ID,
	SEQ,
	ADDR,
	BLOCK_HASH,
	TX,
	TX_HASH,
	EVENT,
	EVENT_OPS,
	TX_RECEIPT,
	FEE,
	HEAD,
	RPC_CLIENT,
] {
	multiNode := NewMultiNode[CHAIN_ID, HEAD, RPC_CLIENT, TX](
		logger, selectionMode, noNewHeadsThreshold, nodes, sendonlys, chainID, "EVM",
	)

	lggr := logger.Named("ChainClient").With("chainID", chainID.String())

	c := &client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]{
		multiNode:           multiNode,
		logger:              lggr,
		sendOnlyErrorParser: sendOnlyErrorParser,
	}

	return c
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) Dial(ctx context.Context) error {
	if err := c.multiNode.Dial(ctx); err != nil {
		return errors.Wrap(err, "failed to dial multiNode")
	}
	return nil
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) NodeRPC() (rpc RPC_CLIENT, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return rpc, err
	}
	return n.RPC(), nil
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BalanceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (*big.Int, error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return nil, err
	}
	return n.RPC().BalanceAt(ctx, account, blockNumber)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BatchCallContext(ctx context.Context, b []any) error {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return err
	}
	return n.RPC().BatchCallContext(ctx, b)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BatchCallContextAll(ctx context.Context, b []any) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	main, selectionErr := c.multiNode.SelectNode()
	for _, n := range c.multiNode.NodesAsSendOnlys() {
		if n == main {
			// main node is used at the end for the return value
			continue
		}
		// Parallel call made to all other nodes with ignored return value
		wg.Add(1)
		go func(n SendOnlyNode[CHAIN_ID, RPC_CLIENT]) {
			defer wg.Done()
			err := n.RPC().BatchCallContext(ctx, b)
			if err != nil {
				c.logger.Debugw("Secondary node BatchCallContext failed", "err", err)
			} else {
				c.logger.Trace("Secondary node BatchCallContext success")
			}
		}(n)
	}

	if selectionErr != nil {
		return selectionErr
	}
	return main.RPC().BatchCallContext(ctx, b)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) ChainType() config.ChainType {
	return c.chainType
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) Close() {
	c.multiNode.Close()
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BlockByHash(ctx context.Context, hash BLOCK_HASH) (h HEAD, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return h, err
	}
	return n.RPC().BlockByHash(ctx, hash)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BlockByNumber(ctx context.Context, number *big.Int) (h HEAD, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return h, err
	}
	return n.RPC().BlockByNumber(ctx, number)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return err
	}
	return n.RPC().CallContext(ctx, result, method, args...)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) CallContract(
	ctx context.Context,
	attempt interface{},
	blockNumber *big.Int,
) (rpcErr []byte, extractErr error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return rpcErr, err
	}
	return n.RPC().CallContract(ctx, attempt, blockNumber)
}

// ChainID makes a direct RPC call. In most cases it should be better to use the configured chain id instead by
// calling ConfiguredChainID.
func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) ChainID(ctx context.Context) (id CHAIN_ID, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return id, err
	}
	return n.RPC().ChainID(ctx)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) CodeAt(ctx context.Context, account ADDR, blockNumber *big.Int) (code []byte, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return code, err
	}
	return n.RPC().CodeAt(ctx, account, blockNumber)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) ConfiguredChainID() (id CHAIN_ID) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return id
	}
	return n.ConfiguredChainID()
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) EstimateGas(ctx context.Context, call any) (gas uint64, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return gas, err
	}
	return n.RPC().EstimateGas(ctx, call)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) FilterEvents(ctx context.Context, query EVENT_OPS) (e []EVENT, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return e, err
	}
	return n.RPC().FilterEvents(ctx, query)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) IsL2() bool {
	return c.ChainType().IsL2()
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) LatestBlockHeight(ctx context.Context) (h *big.Int, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return h, err
	}
	return n.RPC().LatestBlockHeight(ctx)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) LINKBalance(ctx context.Context, accountAddress ADDR, linkAddress ADDR) (b *assets.Link, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return b, err
	}
	return n.RPC().LINKBalance(ctx, accountAddress, linkAddress)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) NodeStates() map[string]string {
	return c.multiNode.NodeStates()
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) PendingSequenceAt(ctx context.Context, addr ADDR) (s SEQ, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return s, err
	}
	return n.RPC().PendingSequenceAt(ctx, addr)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) SendEmptyTransaction(
	ctx context.Context,
	newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt any, err error),
	seq SEQ,
	gasLimit uint32,
	fee FEE,
	fromAddress ADDR,
) (txhash string, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return txhash, err
	}
	return n.RPC().SendEmptyTransaction(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) SendTransaction(ctx context.Context, tx TX) error {
	main, err := c.multiNode.SelectNode()
	for _, n := range c.multiNode.NodesAsSendOnlys() {
		if n == main {
			// main node is used at the end for the return value
			continue
		}
		// Parallel send to all other nodes with ignored return value
		// Async - we do not want to block the main thread with secondary nodes
		// in case they are unreliable/slow.
		// It is purely a "best effort" send.
		// Resource is not unbounded because the default context has a timeout.
		c.multiNode.WrapSendOnlyTransaction(ctx, c.logger, tx, n, c.sendOnlyTransaction)
	}

	if err != nil {
		return err
	}
	return main.RPC().SendTransaction(ctx, tx)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) sendOnlyTransaction(ctx context.Context, lggr logger.Logger, tx TX, n SendOnlyNode[CHAIN_ID, RPC_CLIENT]) {
	err := n.RPC().SendTransaction(ctx, tx)
	lggr.Debugw("Sendonly node sent transaction", "name", n.String(), "tx", tx, "err", err)
	sendOnlyError, err := c.sendOnlyErrorParser(err)
	if sendOnlyError != Successful {
		lggr.Warnw("Eth client returned error", "name", n.String(), "err", err, "tx", tx)
	}
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) SequenceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (s SEQ, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return s, err
	}
	return n.RPC().SequenceAt(ctx, account, blockNumber)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) SimulateTransaction(ctx context.Context, tx TX) error {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return err
	}
	return n.RPC().SimulateTransaction(ctx, tx)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (s types.Subscription, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return s, err
	}
	return n.RPC().Subscribe(ctx, channel, args...)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) TokenBalance(ctx context.Context, account ADDR, tokenAddr ADDR) (b *big.Int, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return b, err
	}
	return n.RPC().TokenBalance(ctx, account, tokenAddr)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) TransactionByHash(ctx context.Context, txHash TX_HASH) (tx TX, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return tx, err
	}
	return n.RPC().TransactionByHash(ctx, txHash)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) TransactionReceipt(ctx context.Context, txHash TX_HASH) (txr TX_RECEIPT, err error) {
	n, err := c.multiNode.SelectNode()
	if err != nil {
		return txr, err
	}
	return n.RPC().TransactionReceipt(ctx, txHash)
}
