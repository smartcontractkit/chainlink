package txmgr

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink-framework/multinode"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
)

var _ TxmClient = (*evmTxmClient)(nil)

type evmTxmClient struct {
	client       client.Client
	clientErrors config.ClientErrors
}

func NewEvmTxmClient(c client.Client, clientErrors config.ClientErrors) *evmTxmClient {
	return &evmTxmClient{client: c, clientErrors: clientErrors}
}

func (c *evmTxmClient) PendingSequenceAt(ctx context.Context, addr common.Address) (types.Nonce, error) {
	return c.PendingNonceAt(ctx, addr)
}

func (c *evmTxmClient) ConfiguredChainID() *big.Int {
	return c.client.ConfiguredChainID()
}

func (c *evmTxmClient) BatchSendTransactions(
	ctx context.Context,
	attempts []TxAttempt,
	batchSize int,
	lggr logger.SugaredLogger,
) (
	codes []multinode.SendTxReturnCode,
	txErrs []error,
	broadcastTime time.Time,
	successfulTxIDs []int64,
	err error,
) {
	// preallocate
	codes = make([]multinode.SendTxReturnCode, len(attempts))
	txErrs = make([]error, len(attempts))

	reqs, broadcastTime, successfulTxIDs, batchErr := batchSendTransactions(ctx, attempts, batchSize, lggr, c.client)
	err = errors.Join(err, batchErr) // this error does not block processing

	// safety check - exits before processing
	if len(reqs) != len(attempts) {
		lenErr := fmt.Errorf("Returned request data length (%d) != number of tx attempts (%d)", len(reqs), len(attempts))
		err = errors.Join(err, lenErr)
		lggr.Criticalw("Mismatched length", "err", err)
		return
	}

	// for each batched tx convert response to standard error code
	var wg sync.WaitGroup
	wg.Add(len(reqs))
	processingErr := make([]error, len(attempts))
	for index := range reqs {
		go func(i int) {
			defer wg.Done()

			// convert to tx for logging purposes - exits early if error occurs
			tx, signedErr := GetGethSignedTx(attempts[i].SignedRawTx)
			if signedErr != nil {
				signedErrMsg := fmt.Sprintf("failed to process tx (index %d)", i)
				lggr.Errorw(signedErrMsg, "err", signedErr)
				processingErr[i] = fmt.Errorf("%s: %w", signedErrMsg, signedErr)
				return
			}
			sendErr := reqs[i].Error
			codes[i] = client.ClassifySendError(sendErr, c.clientErrors, lggr, tx, attempts[i].Tx.FromAddress, c.client.IsL2())
			txErrs[i] = sendErr
		}(index)
	}
	wg.Wait()
	err = errors.Join(err, errors.Join(processingErr...)) // merge errors together
	return
}

func (c *evmTxmClient) SendTransactionReturnCode(ctx context.Context, etx Tx, attempt TxAttempt, lggr logger.SugaredLogger) (multinode.SendTxReturnCode, error) {
	signedTx, err := GetGethSignedTx(attempt.SignedRawTx)
	if err != nil {
		lggr.Criticalw("Fatal error signing transaction", "err", err, "etx", etx)
		return multinode.Fatal, err
	}
	return c.client.SendTransactionReturnCode(ctx, signedTx, etx.FromAddress)
}

func (c *evmTxmClient) PendingNonceAt(ctx context.Context, fromAddress common.Address) (n types.Nonce, err error) {
	nextNonce, err := c.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return n, err
	}

	if nextNonce > math.MaxInt64 {
		return n, fmt.Errorf("nonce overflow, got: %v", nextNonce)
	}
	return types.Nonce(nextNonce), nil
}

func (c *evmTxmClient) SequenceAt(ctx context.Context, addr common.Address, blockNum *big.Int) (types.Nonce, error) {
	nonce, err := c.client.NonceAt(ctx, addr, blockNum)
	if nonce > math.MaxInt64 {
		return 0, fmt.Errorf("overflow for nonce: %d", nonce)
	}

	return types.Nonce(nonce), err
}

func (c *evmTxmClient) BatchGetReceipts(ctx context.Context, attempts []TxAttempt) (txReceipt []*types.Receipt, txErr []error, funcErr error) {
	var reqs []rpc.BatchElem
	for _, attempt := range attempts {
		res := &types.Receipt{}
		req := rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{attempt.Hash},
			Result: res,
		}
		txReceipt = append(txReceipt, res)
		reqs = append(reqs, req)
	}

	if err := c.client.BatchCallContext(ctx, reqs); err != nil {
		return nil, nil, fmt.Errorf("error fetching receipts with BatchCallContext: %w", err)
	}

	for _, req := range reqs {
		txErr = append(txErr, req.Error)
	}
	return txReceipt, txErr, nil
}

// sendEmptyTransaction sends a transaction with 0 Eth and an empty payload to the burn address
// May be useful for clearing stuck nonces
func (c *evmTxmClient) SendEmptyTransaction(
	ctx context.Context,
	newTxAttempt func(ctx context.Context, seq types.Nonce, feeLimit uint64, fee gas.EvmFee, fromAddress common.Address) (attempt TxAttempt, err error),
	seq types.Nonce,
	gasLimit uint64,
	fee gas.EvmFee,
	fromAddress common.Address,
) (txhash string, err error) {
	defer utils.WrapIfError(&err, "sendEmptyTransaction failed")

	attempt, err := newTxAttempt(ctx, seq, gasLimit, fee, fromAddress)
	if err != nil {
		return txhash, err
	}

	signedTx, err := GetGethSignedTx(attempt.SignedRawTx)
	if err != nil {
		return txhash, err
	}

	_, err = c.client.SendTransactionReturnCode(ctx, signedTx, fromAddress)
	return signedTx.Hash().String(), err
}

func (c *evmTxmClient) CallContract(ctx context.Context, a TxAttempt, blockNumber *big.Int) (rpcErr fmt.Stringer, extractErr error) {
	_, errCall := c.client.CallContract(ctx, ethereum.CallMsg{
		From:       a.Tx.FromAddress,
		To:         &a.Tx.ToAddress,
		Gas:        a.Tx.FeeLimit,
		GasPrice:   a.TxFee.GasPrice.ToInt(),
		GasFeeCap:  a.TxFee.GasFeeCap.ToInt(),
		GasTipCap:  a.TxFee.GasTipCap.ToInt(),
		Value:      nil,
		Data:       a.Tx.EncodedPayload,
		AccessList: nil,
	}, blockNumber)
	return client.ExtractRPCError(errCall)
}

func (c *evmTxmClient) HeadByHash(ctx context.Context, hash common.Hash) (*types.Head, error) {
	return c.client.HeadByHash(ctx, hash)
}

func (c *evmTxmClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	return c.client.BatchCallContext(ctx, b)
}
