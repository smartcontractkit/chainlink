package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

type ChainWriter interface {
	// SubmitSignedTransaction packs and broadcasts a transaction to the underlying chain.
	//
	// The `transactionID` will be used by the underlying TXM as an idempotency key, and unique reference to track transaction attempts.
	SubmitSignedTransaction(ctx context.Context, payload []byte, signature map[string]any, transactionID uuid.UUID, toAddress string, meta *TxMeta, value big.Int) (int64, error)

	// GetTransactionStatus returns the current status of a transaction in the underlying chain's TXM.
	GetTransactionStatus(ctx context.Context, transactionID uuid.UUID) (TransactionStatus, error)

	// GetFeeComponents retrieves the associated gas costs for executing a transaction.
	GetFeeComponents(ctx context.Context) (ChainFeeComponents, error)
}

// TxMeta contains metadata fields for a transaction.
//
// Eventually this will replace, or be replaced by (via a move), the `TxMeta` in core:
// https://github.com/smartcontractkit/chainlink/blob/dfc399da715f16af1fcf6441ea5fc47b71800fa1/common/txmgr/types/tx.go#L121
type TxMeta = map[string]string

// TransactionStatus are the status we expect every TXM to support and that can be returned by StatusForUUID.
type TransactionStatus int

const (
	Unknown TransactionStatus = iota
	Unconfirmed
	Finalized
	Failed
	Fatal
)

// ChainFeeComponents contains the different cost components of executing a transaction.
type ChainFeeComponents struct {
	// The cost of executing transaction in the chain's EVM (or the L2 environment).
	ExecutionFee big.Int

	// The cost associated with an L2 posting a transaction's data to the L1.
	DataAvailabilityFee big.Int
}

type ChainWriterService interface {
	services.ServiceCtx
	ChainWriter
}

// Compile-time assertion that chainWriter implements the ChainWriterService interface.
var _ ChainWriterService = (*chainWriter)(nil)

func NewChainWriterService(config config.ChainWriter, logger logger.Logger, client evmclient.Client) ChainWriterService {
	return &chainWriter{logger: logger, client: client, config: config}
}

type chainWriter struct {
	commonservices.StateMachine

	logger logger.Logger
	client evmclient.Client
	config config.ChainWriter
}

func (w *chainWriter) SubmitSignedTransaction(ctx context.Context, payload []byte, signature map[string]any, transactionID uuid.UUID, toAddress string, meta *TxMeta, value big.Int) (int64, error) {
	return 0, fmt.Errorf("not implemented")
}

func (w *chainWriter) GetTransactionStatus(ctx context.Context, transactionID uuid.UUID) (TransactionStatus, error) {
	return Unknown, fmt.Errorf("not implemented")
}

func (w *chainWriter) GetFeeComponents(ctx context.Context) (ChainFeeComponents, error) {
	return ChainFeeComponents{}, fmt.Errorf("not implemented")
}

func (w *chainWriter) Close() error {
	return w.StopOnce(w.Name(), func() error {
		_, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// TODO(nickcorin): Add shutdown steps here.
		return nil
	})
}

func (w *chainWriter) HealthReport() map[string]error {
	return map[string]error{
		w.Name(): nil,
	}
}

func (w *chainWriter) Name() string {
	return "chain-writer"
}

func (w *chainWriter) Ready() error {
	// TODO(nickcorin): Return nil here once the implementation is done.
	return fmt.Errorf("not fully implemented")
}

func (w *chainWriter) Start(ctx context.Context) error {
	return w.StartOnce(w.Name(), func() error {
		// TODO(nickcorin): Add startup steps here.
		return nil
	})
}
