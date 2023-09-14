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
	NodeRPC() RPC_CLIENT
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

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) NodeRPC() RPC_CLIENT {
	return c.multiNode.SelectNode().RPC()
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BalanceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (*big.Int, error) {
	return c.multiNode.SelectNode().RPC().BalanceAt(ctx, account, blockNumber)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BatchCallContext(ctx context.Context, b []any) error {
	return c.multiNode.SelectNode().RPC().BatchCallContext(ctx, b)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BatchCallContextAll(ctx context.Context, b []any) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	main := c.multiNode.SelectNode()
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

	return main.RPC().BatchCallContext(ctx, b)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) ChainType() config.ChainType {
	return c.chainType
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) Close() {
	c.multiNode.Close()
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BlockByHash(ctx context.Context, hash BLOCK_HASH) (HEAD, error) {
	return c.multiNode.SelectNode().RPC().BlockByHash(ctx, hash)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) BlockByNumber(ctx context.Context, number *big.Int) (HEAD, error) {
	return c.multiNode.SelectNode().RPC().BlockByNumber(ctx, number)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return c.multiNode.SelectNode().RPC().CallContext(ctx, result, method, args...)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) CallContract(
	ctx context.Context,
	attempt interface{},
	blockNumber *big.Int,
) (rpcErr []byte, extractErr error) {
	return c.multiNode.SelectNode().RPC().CallContract(ctx, attempt, blockNumber)
}

// ChainID makes a direct RPC call. In most cases it should be better to use the configured chain id instead by
// calling ConfiguredChainID.
func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) ChainID(ctx context.Context) (CHAIN_ID, error) {
	return c.multiNode.SelectNode().RPC().ChainID(ctx)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) CodeAt(ctx context.Context, account ADDR, blockNumber *big.Int) ([]byte, error) {
	return c.multiNode.SelectNode().RPC().CodeAt(ctx, account, blockNumber)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) ConfiguredChainID() CHAIN_ID {
	return c.multiNode.SelectNode().ConfiguredChainID()
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) EstimateGas(ctx context.Context, call any) (gas uint64, err error) {
	return c.multiNode.SelectNode().RPC().EstimateGas(ctx, call)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) FilterEvents(ctx context.Context, query EVENT_OPS) ([]EVENT, error) {
	return c.multiNode.SelectNode().RPC().FilterEvents(ctx, query)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) IsL2() bool {
	return c.ChainType().IsL2()
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
	return c.multiNode.SelectNode().RPC().LatestBlockHeight(ctx)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) LINKBalance(ctx context.Context, accountAddress ADDR, linkAddress ADDR) (*assets.Link, error) {
	return c.multiNode.SelectNode().RPC().LINKBalance(ctx, accountAddress, linkAddress)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) NodeStates() map[string]string {
	return c.multiNode.NodeStates()
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) PendingSequenceAt(ctx context.Context, addr ADDR) (SEQ, error) {
	return c.multiNode.SelectNode().RPC().PendingSequenceAt(ctx, addr)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) SendEmptyTransaction(
	ctx context.Context,
	newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt any, err error),
	seq SEQ,
	gasLimit uint32,
	fee FEE,
	fromAddress ADDR,
) (txhash string, err error) {
	return c.multiNode.SelectNode().RPC().SendEmptyTransaction(ctx, newTxAttempt, seq, gasLimit, fee, fromAddress)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) SendTransaction(ctx context.Context, tx TX) error {
	main := c.multiNode.SelectNode()
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

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) SequenceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (SEQ, error) {
	return c.multiNode.SelectNode().RPC().SequenceAt(ctx, account, blockNumber)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) SimulateTransaction(ctx context.Context, tx TX) error {
	return c.multiNode.SelectNode().RPC().SimulateTransaction(ctx, tx)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (types.Subscription, error) {
	return c.multiNode.SelectNode().RPC().Subscribe(ctx, channel, args)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) TokenBalance(ctx context.Context, account ADDR, tokenAddr ADDR) (*big.Int, error) {
	return c.multiNode.SelectNode().RPC().TokenBalance(ctx, account, tokenAddr)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) TransactionByHash(ctx context.Context, txHash TX_HASH) (TX, error) {
	return c.multiNode.SelectNode().RPC().TransactionByHash(ctx, txHash)
}

func (c *client[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, RPC_CLIENT]) TransactionReceipt(ctx context.Context, txHash TX_HASH) (TX_RECEIPT, error) {
	return c.multiNode.SelectNode().RPC().TransactionReceipt(ctx, txHash)
}
