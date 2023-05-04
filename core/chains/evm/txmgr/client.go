package txmgr

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	clienttypes "github.com/smartcontractkit/chainlink/v2/common/chains/client"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type TxmClient[
	CHAIN_ID txmgrtypes.ID,
	ADDR types.Hashable,
	TX_HASH types.Hashable,
	BLOCK_HASH types.Hashable,
	R txmgrtypes.ChainReceipt[TX_HASH],
	SEQ txmgrtypes.Sequence,
	FEE txmgrtypes.Fee,
	ADD any,
] interface {
	BatchSendTransactions(
		context.Context,
		txmgrtypes.TxStore[ADDR, CHAIN_ID, TX_HASH, BLOCK_HASH, R, SEQ, FEE, ADD],
		[]txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD],
		int,
		logger.Logger,
	) ([]clienttypes.SendTxReturnCode, []error, error)
	SendTransactionReturnCode(
		context.Context,
		txmgrtypes.Tx[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD],
		txmgrtypes.TxAttempt[CHAIN_ID, ADDR, TX_HASH, BLOCK_HASH, R, FEE, ADD],
		logger.Logger,
	) (clienttypes.SendTxReturnCode, error)
	PendingNonceAt(context.Context, ADDR) (int64, error)
}

type EvmTxmClient = TxmClient[*big.Int, common.Address, common.Hash, common.Hash, *evmtypes.Receipt, evmtypes.Nonce, gas.EvmFee, EvmAccessList]

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