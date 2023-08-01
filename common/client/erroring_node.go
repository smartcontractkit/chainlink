package client

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/pkg/errors"
	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
)

type erroringNode[
	CHAINID types.ID,
	SEQ types.Sequence,
	ADDR types.Hashable,
	BLOCK any,
	BLOCKHASH types.Hashable,
	TX any,
	TXHASH types.Hashable,
	EVENT any,
	EVENTOPS any, // event filter query options
	TXRECEIPT txmgrtypes.ChainReceipt[TXHASH, BLOCKHASH],
	FEE feetypes.Fee,
	HEAD *types.Head[BLOCKHASH],
] struct {
	errMsg string
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) ChainID() (chainID CHAINID, err error) {
	return chainID, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) ConfiguredChainID() (chainID CHAINID) {
	return chainID
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) FilterEvents(ctx context.Context, query EVENTOPS) ([]EVENT, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) Start(ctx context.Context) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) Close() error {
	return nil
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) Verify(ctx context.Context, expectedChainID *big.Int) (err error) {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) BatchCallContext(ctx context.Context, b []any) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) BatchCallContextAll(ctx context.Context, b []any) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) BatchGetReceipts(
	ctx context.Context,
	attempts []txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
) (txReceipt []TXRECEIPT, txErr []error, err error) {
	return nil, nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) BatchSendTransactions(
	ctx context.Context,
	updateBroadcastTime func(now time.Time, txIDs []int64) error,
	attempts []txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
	bathSize int,
	lggr logger.Logger,
) ([]clienttypes.SendTxReturnCode, []error, error) {
	return nil, nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) BlockByHash(ctx context.Context, hash BLOCKHASH) (*BLOCK, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) SendTransaction(ctx context.Context, tx *TX) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) LatestBlockHeight(ctx context.Context) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) PendingSequenceAt(ctx context.Context, addr ADDR) (sequence SEQ, err error) {
	return sequence, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) SendEmptyTransaction(
	ctx context.Context,
	newTxAttempt func(seq SEQ, feeLimit uint32, fee FEE, fromAddress ADDR) (attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE], err error),
	seq SEQ,
	gasLimit uint32,
	fee FEE,
	fromAddress ADDR,
) (txhash string, err error) {
	return "", errors.New(e.errMsg)
}

// func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) SendTransactionReturnCode(
// 	ctx context.Context,
// 	tx txmgrtypes.Tx[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
// 	attempt txmgrtypes.TxAttempt[CHAINID, ADDR, TXHASH, BLOCKHASH, SEQ, FEE],
// 	lggr logger.Logger,
// ) (returnCode clienttypes.SendTxReturnCode, err error) {
// 	return returnCode, errors.New(e.errMsg)
// }

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) SendTransactionReturnCode(
	ctx context.Context,
	tx any,
) (returnCode clienttypes.SendTxReturnCode, err error) {
	return returnCode, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) SequenceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (sequence SEQ, err error) {
	return sequence, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) SimulateTransaction(ctx context.Context, tx *TX) error {
	return errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) TokenBalance(ctx context.Context, account ADDR, tokenAddr ADDR) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) TransactionByHash(ctx context.Context, txHash TXHASH) (*TX, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) TransactionReceipt(ctx context.Context, txHash TXHASH) (*TXRECEIPT, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) PendingCodeAt(ctx context.Context, account ADDR) ([]byte, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) PendingNonceAt(ctx context.Context, account ADDR) (uint64, error) {
	return 0, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) NonceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (uint64, error) {
	return 0, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) BlockByNumber(ctx context.Context, number *big.Int) (*BLOCK, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) BlockNumber(ctx context.Context) (uint64, error) {
	return 0, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) BalanceAt(ctx context.Context, account ADDR, blockNumber *big.Int) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) EstimateGas(ctx context.Context, call any) (uint64, error) {
	return 0, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) HeadByNumber(ctx context.Context, number *big.Int) (head HEAD, err error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) HeadByHash(ctx context.Context, hash BLOCKHASH) (head HEAD, err error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) IsL2() bool {
	return false
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) LINKBalance(ctx context.Context, accountAddress ADDR, linkAddress ADDR) (*assets.Link, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) CallContract(
	ctx context.Context,
	attempt interface{},
	blockNumber *big.Int,
) (rpcErr []byte, extractErr error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) CodeAt(ctx context.Context, account ADDR, blockNumber *big.Int) ([]byte, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) Subscribe(ctx context.Context, channel chan<- types.Head[BLOCKHASH], args ...interface{}) (types.Subscription, error) {
	return nil, errors.New(e.errMsg)
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) String() string {
	return "<erroring node>"
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, TXRECEIPT, EVENT, EVENTOPS, FEE, HEAD]) State() NodeState {
	return NodeStateUnreachable
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) StateAndLatest() (NodeState, int64, *utils.Big) {
	return NodeStateUnreachable, -1, nil
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) Order() int32 {
	return 100
}

func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) DeclareOutOfSync() {
}
func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) DeclareInSync() {
}
func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) DeclareUnreachable() {
}
func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) Name() string {
	return ""
}
func (e *erroringNode[CHAINID, SEQ, ADDR, BLOCK, BLOCKHASH, TX, TXHASH, EVENT, EVENTOPS, TXRECEIPT, FEE, HEAD]) NodeStates() map[int32]string {
	return nil
}
