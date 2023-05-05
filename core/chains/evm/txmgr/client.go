package txmgr

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type TxmClient[
	CHAIN_ID txmgrtypes.ID,
	HEAD txmgrtypes.Head,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH, BLOCK_HASH],
	SEQ txmgrtypes.Sequence,
	FEE txmgrtypes.Fee,
	ADD any,
] interface {
	BatchSendTransactions(
		ctx context.Context,
		store txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD],
		attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD],
		bathSize int,
		lggr logger.Logger,
	) ([]clienttypes.SendTxReturnCode, []error, error)
	SendTransactionReturnCode(
		ctx context.Context,
		tx txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD],
		attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD],
		lggr logger.Logger,
	) (clienttypes.SendTxReturnCode, error)
	PendingNonceAt(ctx context.Context, addr ADDR) (int64, error)
	SequenceAt(ctx context.Context, addr ADDR, blockNum *big.Int) (SEQ, error)
	BatchGetReceipts(
		ctx context.Context,
		attempts []txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD],
	) (txReceipt []R, txErr []error, err error)
	SendEmptyTransaction(
		ctx context.Context,
		txAttemptBuilder txmgrtypes.TxAttemptBuilder[HEAD, FEE, ADDR, TX_HASH, txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD], SEQ],
		seq SEQ,
		gasLimit uint32,
		fee FEE,
		fromAddress ADDR,
	) (txhash string, err error)
	CallContract(
		ctx context.Context,
		attempt txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD],
		blockNumber *big.Int,
	) (rpcErr fmt.Stringer, extractErr error)
}

type EvmTxmClient = TxmClient[*big.Int, *evmtypes.Head, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, EvmAccessList]

var _ EvmTxmClient = (*evmTxmClient)(nil)

type evmTxmClient struct {
	client evmclient.Client
}

func NewEvmTxmClient(c evmclient.Client) *evmTxmClient {
	return &evmTxmClient{client: c}
}

func (c *evmTxmClient) BatchSendTransactions(ctx context.Context, txStore EvmTxStore, attempts []EvmTxAttempt, batchSize int, lggr logger.Logger) (codes []clienttypes.SendTxReturnCode, txErr []error, err error) {
	reqs, err := batchSendTransactions(ctx, txStore, attempts, batchSize, lggr, c.client)
	if err != nil {
		return nil, nil, err
	}

	// for each batched tx convert response to standard error code
	codes = make([]clienttypes.SendTxReturnCode, len(reqs))
	txErr = make([]error, len(reqs))
	for i := range reqs {
		// convert to tx for logging purposes
		tx, err := GetGethSignedTx(attempts[i].SignedRawTx)
		if err != nil {
			return nil, nil, err
		}
		codes[i], txErr[i] = c.client.NewSendErrorReturnCode(tx, attempts[i].Tx.FromAddress, reqs[i].Error)
	}
	return
}

func (c *evmTxmClient) SendTransactionReturnCode(ctx context.Context, etx EvmTx, attempt EvmTxAttempt, lggr logger.Logger) (clienttypes.SendTxReturnCode, error) {
	signedTx, err := GetGethSignedTx(attempt.SignedRawTx)
	if err != nil {
		lggr.Criticalw("Fatal error signing transaction", "err", err, "etx", etx)
		return clienttypes.Fatal, err
	}
	return c.client.SendTransactionReturnCode(ctx, signedTx, etx.FromAddress)
}

func (c *evmTxmClient) PendingNonceAt(ctx context.Context, fromAddress common.Address) (n int64, err error) {
	nextNonce, err := c.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return n, err
	}

	if nextNonce > math.MaxInt64 {
		return n, fmt.Errorf("nonce overflow, got: %v", nextNonce)
	}
	return int64(nextNonce), nil
}

func (c *evmTxmClient) SequenceAt(ctx context.Context, addr common.Address, blockNum *big.Int) (evmtypes.Nonce, error) {
	return c.client.SequenceAt(ctx, addr, blockNum)
}

func (c *evmTxmClient) BatchGetReceipts(ctx context.Context, attempts []EvmTxAttempt) (txReceipt []*evmtypes.Receipt, txErr []error, funcErr error) {
	var reqs []rpc.BatchElem
	for _, attempt := range attempts {
		req := rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{attempt.Hash},
			Result: &evmtypes.Receipt{},
		}
		reqs = append(reqs, req)
	}

	if err := c.client.BatchCallContext(ctx, reqs); err != nil {
		return nil, nil, errors.Wrap(err, "EthConfirmer#batchFetchReceipts error fetching receipts with BatchCallContext")
	}

	for _, req := range reqs {
		result, err := req.Result, req.Error

		receipt, ok := result.(*evmtypes.Receipt)
		if !ok {
			return nil, nil, fmt.Errorf("expected result to be a %T, got %T", (*evmtypes.Receipt)(nil), result)
		}

		txReceipt = append(txReceipt, receipt)
		txErr = append(txErr, err)
	}
	return txReceipt, txErr, nil
}

// sendEmptyTransaction sends a transaction with 0 Eth and an empty payload to the burn address
// May be useful for clearing stuck nonces
func (c *evmTxmClient) SendEmptyTransaction(
	ctx context.Context,
	txAttemptBuilder EvmTxAttemptBuilder,
	seq evmtypes.Nonce,
	gasLimit uint32,
	fee gas.EvmFee,
	fromAddress common.Address,
) (txhash string, err error) {
	defer utils.WrapIfError(&err, "sendEmptyTransaction failed")

	attempt, err := txAttemptBuilder.NewEmptyTxAttempt(seq, gasLimit, fee, fromAddress)
	if err != nil {
		return txhash, err
	}

	signedTx, err := GetGethSignedTx(attempt.SignedRawTx)
	if err != nil {
		return txhash, err
	}

	err = c.client.SendTransaction(ctx, signedTx)
	return signedTx.Hash().String(), err
}

func (c *evmTxmClient) CallContract(ctx context.Context, a EvmTxAttempt, blockNumber *big.Int) (rpcErr fmt.Stringer, extractErr error) {
	_, errCall := c.client.CallContract(ctx, ethereum.CallMsg{
		From:       a.Tx.FromAddress,
		To:         &a.Tx.ToAddress,
		Gas:        uint64(a.Tx.FeeLimit),
		GasPrice:   a.TxFee.Legacy.ToInt(),
		GasFeeCap:  a.TxFee.DynamicFeeCap.ToInt(),
		GasTipCap:  a.TxFee.DynamicTipCap.ToInt(),
		Value:      nil,
		Data:       a.Tx.EncodedPayload,
		AccessList: nil,
	}, blockNumber)
	return evmclient.ExtractRPCError(errCall)
}
