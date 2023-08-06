package client

import (
	"context"
	"math/big"
	"time"

	"github.com/pkg/errors"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type erroringNode[
	CHAIN_ID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK_HASH types.Hashable,
	TX any,
	TX_HASH types.Hashable,
	EVENT any,
	EVENT_OPS any, // event filter query options
	TX_RECEIPT any,
	FEE feetypes.Fee,
	HEAD types.Head[BLOCK_HASH],
	SUB types.Subscription,
] struct {
	errMsg string
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) ChainID() (chainID CHAIN_ID, err error) {
	return chainID, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) ConfiguredChainID() (chainID CHAIN_ID) {
	return chainID
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) FilterEvents(ctx context.Context, query EVENT_OPS) ([]EVENT, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Start(ctx context.Context) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Close() error {
	return nil
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Verify(ctx context.Context, expectedChainID *big.Int) (err error) {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) BatchCallContext(ctx context.Context, b []any) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) BatchCallContextAll(ctx context.Context, b []any) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) BatchSendTransactions(
	ctx context.Context,
	updateBroadcastTime func(now time.Time, txIDs []int64) error,
	attempts []any,
	bathSize int,
	lggr logger.Logger,
) ([]SendTxReturnCode, []error, error) {
	return nil, nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) BlockByHash(ctx context.Context, hash BLOCK_HASH) (head HEAD, err error) {
	return head, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) SendTransaction(ctx context.Context, tx *TX) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) SendEmptyTransaction(
	ctx context.Context,
	newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt any, err error),
	seq SEQ,
	gasLimit uint32,
	fee FEE,
	fromAddress ADDR,
) (txhash string, err error) {
	return "", errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) PendingSequenceAt(ctx context.Context, addr ADDR) (sequence SEQ, err error) {
	return sequence, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) SendTransactionReturnCode(
	ctx context.Context,
	tx *TX,
) (returnCode SendTxReturnCode, err error) {
	return returnCode, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) SequenceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (sequence SEQ, err error) {
	return sequence, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) SimulateTransaction(ctx context.Context, tx *TX) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) TokenBalance(ctx context.Context, account ADDR, tokenAddr ADDR) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) TransactionByHash(ctx context.Context, txHash TX_HASH) (*TX, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) TransactionReceipt(ctx context.Context, txHash TX_HASH) (*TX_RECEIPT, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) PendingCodeAt(ctx context.Context, account ADDR) ([]byte, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) PendingNonceAt(ctx context.Context, account ADDR) (uint64, error) {
	return 0, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) NonceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (uint64, error) {
	return 0, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) BlockByNumber(ctx context.Context, number *big.Int) (head HEAD, err error) {
	return head, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) BlockNumber(ctx context.Context) (uint64, error) {
	return 0, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) BalanceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) EstimateGas(ctx context.Context, call any) (uint64, error) {
	return 0, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) IsL2() bool {
	return false
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) LINKBalance(ctx context.Context, accountAddress ADDR, linkAddress ADDR) (*assets.Link, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) CallContract(
	ctx context.Context,
	attempt interface{},
	blockNumber *big.Int,
) (rpcErr []byte, extractErr error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) CodeAt(ctx context.Context, account ADDR, blockNumber *big.Int) ([]byte, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Subscribe(ctx context.Context, channel chan<- HEAD, args ...interface{}) (subscription SUB, err error) {
	return subscription, errors.New(e.errMsg)
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) String() string {
	return "<erroring node>"
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) State() NodeState {
	return NodeStateUnreachable
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) StateAndLatest() (NodeState, int64, *utils.Big) {
	return NodeStateUnreachable, -1, nil
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Order() int32 {
	return 100
}

func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) DeclareOutOfSync() {
	// no return value
}
func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) DeclareInSync() {
	// no return value
}
func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) DeclareUnreachable() {
	// no return value
}
func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) Name() string {
	return ""
}
func (e *erroringNode[CHAIN_ID, SEQ, ADDR, BLOCK_HASH, TX, TX_HASH, EVENT, EVENT_OPS, TX_RECEIPT, FEE, HEAD, SUB]) NodeStates() map[int32]string {
	return nil
}
