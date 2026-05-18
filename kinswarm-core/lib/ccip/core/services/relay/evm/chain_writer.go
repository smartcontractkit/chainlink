package evm

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/common/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	evmtxmgr "github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type ChainWriterService interface {
	services.ServiceCtx
	commontypes.ChainWriter
}

// Compile-time assertion that chainWriter implements the ChainWriterService interface.
var _ ChainWriterService = (*chainWriter)(nil)

func NewChainWriterService(logger logger.Logger, client evmclient.Client, txm evmtxmgr.TxManager, estimator gas.EvmFeeEstimator, config types.ChainWriterConfig) (ChainWriterService, error) {
	if config.MaxGasPrice == nil {
		return nil, fmt.Errorf("max gas price is required")
	}

	w := chainWriter{
		logger:      logger,
		client:      client,
		txm:         txm,
		ge:          estimator,
		maxGasPrice: config.MaxGasPrice,

		sendStrategy:    txmgr.NewSendEveryStrategy(),
		contracts:       config.Contracts,
		parsedContracts: &codec.ParsedTypes{EncoderDefs: map[string]types.CodecEntry{}, DecoderDefs: map[string]types.CodecEntry{}},
	}

	if config.SendStrategy != nil {
		w.sendStrategy = config.SendStrategy
	}

	if err := w.parseContracts(); err != nil {
		return nil, fmt.Errorf("%w: failed to parse contracts", err)
	}

	var err error
	if w.encoder, err = w.parsedContracts.ToCodec(); err != nil {
		return nil, fmt.Errorf("%w: failed to create codec", err)
	}

	return &w, nil
}

type chainWriter struct {
	commonservices.StateMachine

	logger      logger.Logger
	client      evmclient.Client
	txm         evmtxmgr.TxManager
	ge          gas.EvmFeeEstimator
	maxGasPrice *assets.Wei

	sendStrategy    txmgrtypes.TxStrategy
	contracts       map[string]*types.ContractConfig
	parsedContracts *codec.ParsedTypes

	encoder commontypes.Encoder
}

// SubmitTransaction ...
//
// Note: The codec that ChainWriter uses to encode the parameters for the contract ABI cannot handle
// `nil` values, including for slices. Until the bug is fixed we need to ensure that there are no
// `nil` values passed in the request.
func (w *chainWriter) SubmitTransaction(ctx context.Context, contract, method string, args any, transactionID string, toAddress string, meta *commontypes.TxMeta, value *big.Int) error {
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

	calldata, err := w.encoder.Encode(ctx, args, codec.WrapItemType(contract, method, true))
	if err != nil {
		return fmt.Errorf("%w: failed to encode args", err)
	}

	var checker evmtxmgr.TransmitCheckerSpec
	if methodConfig.Checker != "" {
		checker.CheckerType = txmgrtypes.TransmitCheckerType(methodConfig.Checker)
	}

	v := big.NewInt(0)
	if value != nil {
		v = value
	}

	var txMeta *txmgrtypes.TxMeta[common.Address, common.Hash]
	if meta != nil && meta.WorkflowExecutionID != nil {
		txMeta = &txmgrtypes.TxMeta[common.Address, common.Hash]{
			WorkflowExecutionID: meta.WorkflowExecutionID,
		}
	}

	gasLimit := methodConfig.GasLimit
	if meta != nil && meta.GasLimit != nil {
		gasLimit = meta.GasLimit.Uint64()
	}

	req := evmtxmgr.TxRequest{
		FromAddress:    methodConfig.FromAddress,
		ToAddress:      common.HexToAddress(toAddress),
		EncodedPayload: calldata,
		FeeLimit:       gasLimit,
		Meta:           txMeta,
		IdempotencyKey: &transactionID,
		Strategy:       w.sendStrategy,
		Checker:        checker,
		Value:          *v,
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
			inputMod, err := methodConfig.InputModifications.ToModifier(codec.DecoderHooks...)
			if err != nil {
				return fmt.Errorf("%w: failed to create input mods", err)
			}

			input := types.NewCodecEntry(abiMethod.Inputs, abiMethod.ID, inputMod)

			if err = input.Init(); err != nil {
				return fmt.Errorf("%w: failed to init codec entry for method %s", err, method)
			}

			w.parsedContracts.EncoderDefs[codec.WrapItemType(contract, method, true)] = input
		}
	}

	return nil
}

func (w *chainWriter) GetTransactionStatus(ctx context.Context, transactionID string) (commontypes.TransactionStatus, error) {
	return w.txm.GetTransactionStatus(ctx, transactionID)
}

// GetFeeComponents the execution and data availability (L1Oracle) fees for the chain.
// Dynamic fees (introduced in EIP-1559) include a fee cap and a tip cap. If the dyanmic fee is not available,
// (if the chain doesn't support dynamic TXs) the legacy GasPrice is used.
func (w *chainWriter) GetFeeComponents(ctx context.Context) (*commontypes.ChainFeeComponents, error) {
	if w.ge == nil {
		return nil, fmt.Errorf("gas estimator not available")
	}

	fee, _, err := w.ge.GetFee(ctx, nil, 0, w.maxGasPrice, nil, nil)
	if err != nil {
		return nil, err
	}
	// Use legacy if no dynamic is available.
	gasPrice := fee.Legacy.ToInt()
	if fee.DynamicFeeCap != nil {
		gasPrice = fee.DynamicFeeCap.ToInt()
	}
	if gasPrice == nil {
		return nil, fmt.Errorf("dynamic fee and legacy gas price missing %+v", fee)
	}
	l1Oracle := w.ge.L1Oracle()
	if l1Oracle == nil {
		return &commontypes.ChainFeeComponents{
			ExecutionFee:        gasPrice,
			DataAvailabilityFee: big.NewInt(0),
		}, nil
	}
	l1OracleFee, err := l1Oracle.GasPrice(ctx)
	if err != nil {
		return nil, err
	}

	return &commontypes.ChainFeeComponents{
		ExecutionFee:        gasPrice,
		DataAvailabilityFee: big.NewInt(l1OracleFee.Int64()),
	}, nil
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
