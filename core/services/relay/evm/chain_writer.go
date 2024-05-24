package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type ChainWriterService interface {
	services.ServiceCtx
	types.ChainWriter
}

// Compile-time assertion that chainWriter implements the ChainWriterService interface.
var _ ChainWriterService = (*chainWriter)(nil)

func NewChainWriterService(config config.ChainWriter, logger logger.Logger, client evmclient.Client, txm evmtxmgr.TxManager) ChainWriterService {
	return &chainWriter{logger: logger, client: client, config: config, txm: txm}
}

type chainWriter struct {
	commonservices.StateMachine

	logger logger.Logger
	client evmclient.Client
	config config.ChainWriter
	txm    evmtxmgr.TxManager
}

func (w *chainWriter) SubmitSignedTransaction(ctx context.Context, payload []byte, signatures map[string]any, transactionID uuid.UUID, toAddress string, meta *types.TxMeta, value big.Int) error {
	forwarderABI := evmtypes.MustGetABI(w.config.ABI())

	// TODO(nickcorin):
	// Check the required format for the signatures when packing the ABI. The original type was [][]byte, however we'll need strict type assertions
	// translating any -> []byte.
	//
	// Also, double check that the method name is always going to be 'report'?
	calldata, err := forwarderABI.Pack("report", common.HexToAddress(toAddress), payload, signatures)
	if err != nil {
		return fmt.Errorf("pack forwarder abi: %w", err)
	}

	// TODO(nickcorin): Change this to be config driven.
	sendStrategy := txmgr.SendEveryStrategy{}

	var checker evmtxmgr.TransmitCheckerSpec
	if w.config.Checker() != "" {
		checker.CheckerType = txmgrtypes.TransmitCheckerType(w.config.Checker())
	}

	req := evmtxmgr.TxRequest{
		FromAddress:    w.config.FromAddress().Address(),
		ToAddress:      w.config.ForwarderAddress().Address(),
		EncodedPayload: calldata,
		FeeLimit:       w.config.GasLimit(),
		Meta:           &txmgrtypes.TxMeta[common.Address, common.Hash]{WorkflowExecutionID: meta.WorkflowExecutionID},
		Strategy:       sendStrategy,
		Checker:        checker,
	}

	_, err = w.txm.CreateTransaction(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create tx: %w", err)
	}

	return nil
}

func (w *chainWriter) GetTransactionStatus(ctx context.Context, transactionID uuid.UUID) (types.TransactionStatus, error) {
	return types.Unknown, fmt.Errorf("not implemented")
}

func (w *chainWriter) GetFeeComponents(ctx context.Context) (*types.ChainFeeComponents, error) {
	return nil, fmt.Errorf("not implemented")
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
