package evm

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
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

func NewChainWriterService(logger logger.Logger, client evmclient.Client, txm evmtxmgr.TxManager, config types.ChainWriterConfig) (ChainWriterService, error) {
	w := chainWriter{
		logger: logger,
		client: client,
		txm:    txm,

		sendStrategy:    txmgr.NewSendEveryStrategy(),
		contracts:       config.Contracts,
		parsedContracts: &parsedTypes{encoderDefs: map[string]types.CodecEntry{}, decoderDefs: map[string]types.CodecEntry{}},
	}

	if config.SendStrategy != nil {
		w.sendStrategy = config.SendStrategy
	}

	if err := w.parseContracts(); err != nil {
		return nil, fmt.Errorf("%w: failed to parse contracts", err)
	}

	var err error
	if w.encoder, err = w.parsedContracts.toCodec(); err != nil {
		return nil, fmt.Errorf("%w: failed to create codec", err)
	}

	return &w, nil
}

type chainWriter struct {
	commonservices.StateMachine

	logger logger.Logger
	client evmclient.Client
	txm    evmtxmgr.TxManager

	sendStrategy    txmgrtypes.TxStrategy
	contracts       map[string]*types.ContractConfig
	parsedContracts *parsedTypes

	encoder commontypes.Encoder
}

// SubmitTransaction ...
//
// Note: The codec that ChainWriter uses to encode the parameters for the contract ABI cannot handle
// `nil` values, including for slices. Until the bug is fixed we need to ensure that there are no
// `nil` values passed in the request.
func (w *chainWriter) SubmitTransaction(ctx context.Context, contract, method string, args any, transactionID uuid.UUID, toAddress string, meta *commontypes.TxMeta, value big.Int) error {
	if !common.IsHexAddress(toAddress) {
		return fmt.Errorf("toAddress is not a valid ethereum address: %v", toAddress)
	}

	contractConfig, ok := w.contracts[contract]
	if !ok {
		return fmt.Errorf("contract config not found: %v", contract)
	}

	methodConfig, ok := contractConfig.Configs[method]
	if !ok {
		return fmt.Errorf("method config not found: %v", method)
	}

	calldata, err := w.encoder.Encode(ctx, args, wrapItemType(contract, method, true))
	if err != nil {
		return fmt.Errorf("%w: failed to encode args", err)
	}

	var checker evmtxmgr.TransmitCheckerSpec
	if methodConfig.Checker != "" {
		checker.CheckerType = txmgrtypes.TransmitCheckerType(methodConfig.Checker)
	}

	req := evmtxmgr.TxRequest{
		FromAddress:    methodConfig.FromAddress,
		ToAddress:      common.HexToAddress(toAddress),
		EncodedPayload: calldata,
		FeeLimit:       methodConfig.GasLimit,
		Meta:           &txmgrtypes.TxMeta[common.Address, common.Hash]{WorkflowExecutionID: meta.WorkflowExecutionID},
		Strategy:       w.sendStrategy,
		Checker:        checker,
	}

	_, err = w.txm.CreateTransaction(ctx, req)
	if err != nil {
		return fmt.Errorf("%w; failed to create tx", err)
	}

	return nil
}

func (w *chainWriter) parseContracts() error {
	for contract, contractConfig := range w.contracts {
		abi, err := abi.JSON(strings.NewReader(contractConfig.ContractABI))
		if err != nil {
			return fmt.Errorf("%w: failed to parse contract abi", err)
		}

		for method, methodConfig := range contractConfig.Configs {
			abiMethod, ok := abi.Methods[methodConfig.ChainSpecificName]
			if !ok {
				return fmt.Errorf("%w: method %s doesn't exist", commontypes.ErrInvalidConfig, methodConfig.ChainSpecificName)
			}

			// ABI.Pack prepends the method.ID to the encodings, we'll need the encoder to do the same.
			inputMod, err := methodConfig.InputModifications.ToModifier(evmDecoderHooks...)
			if err != nil {
				return fmt.Errorf("%w: failed to create input mods", err)
			}

			input := types.NewCodecEntry(abiMethod.Inputs, abiMethod.ID, inputMod)

			if err = input.Init(); err != nil {
				return fmt.Errorf("%w: failed to init codec entry for method %s", err, method)
			}

			w.parsedContracts.encoderDefs[wrapItemType(contract, method, true)] = input
		}
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
