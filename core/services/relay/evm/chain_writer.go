package evm

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonservices "github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	evmtxmgr "github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-framework/chains/txmgr"
	txmgrtypes "github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type ChainWriterService interface {
	services.ServiceCtx
	commontypes.ContractWriter
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

		contracts:       config.Contracts,
		parsedContracts: &codec.ParsedTypes{EncoderDefs: map[string]types.CodecEntry{}, DecoderDefs: map[string]types.CodecEntry{}},
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
		Strategy:       txmgr.NewSendEveryStrategy(),
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
	gasPrice := fee.GasPrice.ToInt()
	if fee.GasFeeCap != nil {
		gasPrice = fee.GasFeeCap.ToInt()
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

// GetEstimateFee returns total cost of TX execution in the underlying chain's currency.
// The value (val) is included in the fee calculation.
func (w *chainWriter) GetEstimateFee(ctx context.Context, contract, method string, args any, toAddress string, meta *commontypes.TxMeta, val *big.Int) (commontypes.EstimateFee, error) {
	calldata, err := w.encoder.Encode(ctx, args, codec.WrapItemType(contract, method, true))
	if err != nil {
		return commontypes.EstimateFee{}, fmt.Errorf("%w: failed to encode args", err)
	}

	to := common.HexToAddress(toAddress)
	var v assets.Eth
	if val != nil {
		v = assets.Eth(*val)
	}

	contractConfig, ok := w.contracts[contract]
	if !ok {
		return commontypes.EstimateFee{}, fmt.Errorf("contract config not found: %v", contract)
	}

	methodConfig, ok := contractConfig.Configs[method]
	if !ok {
		return commontypes.EstimateFee{}, fmt.Errorf("method config not found: %v", method)
	}

	gasLimit := methodConfig.GasLimit
	if meta != nil && meta.GasLimit != nil {
		gasLimit = meta.GasLimit.Uint64()
	}

	from := common.Address{}
	cost, err := w.getMaxCost(ctx, v, calldata, gasLimit, w.maxGasPrice, &from, &to)
	if err != nil {
		return commontypes.EstimateFee{}, err
	}

	return commontypes.EstimateFee{
		Fee:      cost,
		Decimals: 18,
	}, nil
}

func (w *chainWriter) getMaxCost(ctx context.Context, amount assets.Eth, calldata []byte,
	gasLimit uint64, maxGasPrice *assets.Wei, fromAddress, toAddress *common.Address) (*big.Int, error) {
	fee, err := w.GetFeeComponents(ctx)
	var gasPrice *big.Int
	if err != nil {
		w.logger.Warnf("%w: GetFeeComponents failed; use maxFeePrice instead", err)
		gasPrice = fee.ExecutionFee
	} else {
		gasPrice = (*big.Int)(maxGasPrice)
	}

	estimateGas, err := w.client.EstimateGas(ctx, ethereum.CallMsg{To: toAddress, Data: calldata})
	if err != nil {
		w.logger.Warnf("%w: EstimateGas failed; use gasLimit instead", err)
		estimateGas = gasLimit
	}

	totalFee := new(big.Int).Mul(gasPrice, big.NewInt(int64(estimateGas)))
	amountWithFees := new(big.Int).Add(amount.ToInt(), totalFee)

	return amountWithFees, nil
}

func (w *chainWriter) Close() error {
	return w.StopOnce(w.Name(), func() error {
		return nil
	})
}

func (w *chainWriter) HealthReport() map[string]error {
	return map[string]error{w.Name(): w.Healthy()}
}

func (w *chainWriter) Name() string {
	return w.logger.Name()
}

func (w *chainWriter) Start(ctx context.Context) error {
	return w.StartOnce(w.Name(), func() error {
		return nil
	})
}
