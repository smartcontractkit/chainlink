package txmgr

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var _ EvmTxmClient = (*evmTxmClient)(nil)

type evmTxmClient struct {
	client evmclient.Client
}

func NewEvmTxmClient(c evmclient.Client) *evmTxmClient {
	return &evmTxmClient{client: c}
}

func (c *evmTxmClient) ConfiguredChainID() *big.Int {
	return c.client.ConfiguredChainID()
}

func (c *evmTxmClient) BatchSendTransactions(ctx context.Context, txStore EvmTxStore, attempts []EvmTxAttempt, batchSize int, lggr logger.Logger) (codes []clienttypes.SendTxReturnCode, txErrs []error, err error) {
	// preallocate
	codes = make([]clienttypes.SendTxReturnCode, len(attempts))
	txErrs = make([]error, len(attempts))

	reqs, batchErr := batchSendTransactions(ctx, txStore, attempts, batchSize, lggr, c.client)
	err = errors.Join(err, batchErr) // this error does not block processing

	// safety check - exits before processing
	if len(reqs) != len(attempts) {
		lenErr := fmt.Errorf("Returned request data length (%d) != number of tx attempts (%d)", len(reqs), len(attempts))
		err = errors.Join(err, lenErr)
		lggr.Criticalw("Mismatched length", "error", err)
		return
	}

	// for each batched tx convert response to standard error code
	for i := range reqs {
		// convert to tx for logging purposes - exits early if error occurs
		tx, signedErr := GetGethSignedTx(attempts[i].SignedRawTx)
		if signedErr != nil {
			err = errors.Join(err, signedErr)
			return
		}
		codes[i], txErrs[i] = evmclient.NewSendErrorReturnCode(reqs[i].Error, lggr, tx, attempts[i].Tx.FromAddress, c.client.IsL2())
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

func (c *evmTxmClient) PendingNonceAt(ctx context.Context, fromAddress common.Address) (n evmtypes.Nonce, err error) {
	nextNonce, err := c.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return n, err
	}

	if nextNonce > math.MaxInt64 {
		return n, fmt.Errorf("nonce overflow, got: %v", nextNonce)
	}
	return evmtypes.Nonce(nextNonce), nil
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
		return nil, nil, fmt.Errorf("EthConfirmer#batchFetchReceipts error fetching receipts with BatchCallContext: %w", err)
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
	newTxAttempt func(seq evmtypes.Nonce, feeLimit uint32, fee gas.EvmFee, fromAddress common.Address) (attempt EvmTxAttempt, err error),
	seq evmtypes.Nonce,
	gasLimit uint32,
	fee gas.EvmFee,
	fromAddress common.Address,
) (txhash string, err error) {
	defer utils.WrapIfError(&err, "sendEmptyTransaction failed")

	attempt, err := newTxAttempt(seq, gasLimit, fee, fromAddress)
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
