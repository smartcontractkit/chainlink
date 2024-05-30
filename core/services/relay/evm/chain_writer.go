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
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

type ChainWriterService interface {
	services.ServiceCtx
	commontypes.ChainWriter
}

// Compile-time assertion that chainWriter implements the ChainWriterService interface.
var _ ChainWriterService = (*chainWriter)(nil)

func NewChainWriterService(logger logger.Logger, client evmclient.Client, txm evmtxmgr.TxManager, config types.ChainWriterConfig) ChainWriterService {
	return &chainWriter{logger: logger, client: client, config: config, txm: txm}
}

type chainWriter struct {
	commonservices.StateMachine

	logger logger.Logger
	client evmclient.Client
	config types.ChainWriterConfig
	txm    evmtxmgr.TxManager
}

func (w *chainWriter) SubmitTransaction(ctx context.Context, contract, method string, args []any, transactionID uuid.UUID, toAddress string, meta *commontypes.TxMeta, value big.Int) error {
	if !common.IsHexAddress(toAddress) {
		return fmt.Errorf("toAddress is not a valid ethereum address: %v", toAddress)
	}

	contractConfig, ok := w.config.Contracts[contract]
	if !ok {
		return fmt.Errorf("contract config not found: %v", contract)
	}
	methodConfig, ok := contractConfig.Configs[method]
	if !ok {
		return fmt.Errorf("method config not found: %v", method)
	}

	forwarderABI := evmtypes.MustGetABI(contractConfig.ContractABI)

	calldata, err := forwarderABI.Pack(methodConfig.ChainSpecificName, args...)
	if err != nil {
		return fmt.Errorf("pack forwarder abi: %w", err)
	}

	var checker evmtxmgr.TransmitCheckerSpec
	if methodConfig.Checker != "" {
		checker.CheckerType = txmgrtypes.TransmitCheckerType(methodConfig.Checker)
	}

	var sendStrategy txmgrtypes.TxStrategy = txmgr.SendEveryStrategy{}
	if w.config.SendStrategy != nil {
		sendStrategy = w.config.SendStrategy
	}

	req := evmtxmgr.TxRequest{
		FromAddress:    methodConfig.FromAddress,
		ToAddress:      common.HexToAddress(toAddress),
		EncodedPayload: calldata,
		FeeLimit:       methodConfig.GasLimit,
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

func (w *chainWriter) GetTransactionStatus(ctx context.Context, transactionID uuid.UUID) (commontypes.TransactionStatus, error) {
	return commontypes.Unknown, fmt.Errorf("not implemented")
}

func (w *chainWriter) GetFeeComponents(ctx context.Context) (*commontypes.ChainFeeComponents, error) {
	return nil, fmt.Errorf("not implemented")
}

func (w *chainWriter) Close() error {
	return w.StopOnce(w.Name(), func() error {
		_, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

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
	return nil
}

func (w *chainWriter) Start(ctx context.Context) error {
	return w.StartOnce(w.Name(), func() error {
		return nil
	})
}
