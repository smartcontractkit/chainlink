package txm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

type AttemptBuilderKeystore interface {
	SignTx(ctx context.Context, fromAddress common.Address, tx *evmtypes.Transaction, chainID *big.Int) (*evmtypes.Transaction, error)
}

type attemptBuilder struct {
	gas.EvmFeeEstimator
	chainID     *big.Int
	priceMaxKey func(common.Address) *assets.Wei
	keystore    AttemptBuilderKeystore
}

func NewAttemptBuilder(chainID *big.Int, priceMaxKey func(common.Address) *assets.Wei, estimator gas.EvmFeeEstimator, keystore AttemptBuilderKeystore) *attemptBuilder {
	return &attemptBuilder{
		chainID:         chainID,
		priceMaxKey:     priceMaxKey,
		EvmFeeEstimator: estimator,
		keystore:        keystore,
	}
}

func (a *attemptBuilder) NewAttempt(ctx context.Context, lggr logger.Logger, tx *types.Transaction, dynamic bool) (*types.Attempt, error) {
	fee, estimatedGasLimit, err := a.EvmFeeEstimator.GetFee(ctx, tx.Data, tx.SpecifiedGasLimit, a.priceMaxKey(tx.FromAddress), &tx.FromAddress, &tx.ToAddress)
	if err != nil {
		return nil, err
	}
	txType := evmtypes.LegacyTxType
	if dynamic {
		txType = evmtypes.DynamicFeeTxType
	}
	return a.newCustomAttempt(ctx, tx, fee, estimatedGasLimit, byte(txType), lggr)
}

func (a *attemptBuilder) NewBumpAttempt(ctx context.Context, lggr logger.Logger, tx *types.Transaction, previousAttempt types.Attempt) (*types.Attempt, error) {
	bumpedFee, bumpedFeeLimit, err := a.EvmFeeEstimator.BumpFee(ctx, previousAttempt.Fee, tx.SpecifiedGasLimit, a.priceMaxKey(tx.FromAddress), nil)
	if err != nil {
		return nil, err
	}
	return a.newCustomAttempt(ctx, tx, bumpedFee, bumpedFeeLimit, previousAttempt.Type, lggr)
}

func (a *attemptBuilder) newCustomAttempt(
	ctx context.Context,
	tx *types.Transaction,
	fee gas.EvmFee,
	estimatedGasLimit uint64,
	txType byte,
	lggr logger.Logger,
) (attempt *types.Attempt, err error) {
	switch txType {
	case 0x0:
		if fee.GasPrice == nil {
			err = fmt.Errorf("tried to create attempt of type %v for txID: %v but estimator did not return legacy fee", txType, tx.ID)
			logger.Sugared(lggr).AssumptionViolation(err.Error())
			return
		}
		return a.newLegacyAttempt(ctx, tx, fee.GasPrice, estimatedGasLimit)
	case 0x2:
		if !fee.ValidDynamic() {
			err = fmt.Errorf("tried to create attempt of type %v for txID: %v but estimator did not return dynamic fee", txType, tx.ID)
			logger.Sugared(lggr).AssumptionViolation(err.Error())
			return
		}
		return a.newDynamicFeeAttempt(ctx, tx, fee.DynamicFee, estimatedGasLimit)
	default:
		return nil, fmt.Errorf("cannot build attempt, unrecognized transaction type: %v", txType)
	}
}

func (a *attemptBuilder) newLegacyAttempt(ctx context.Context, tx *types.Transaction, gasPrice *assets.Wei, estimatedGasLimit uint64) (*types.Attempt, error) {
	var data []byte
	var toAddress common.Address
	value := big.NewInt(0)
	if !tx.IsPurgeable {
		data = tx.Data
		toAddress = tx.ToAddress
		value = tx.Value
	}
	if tx.Nonce == nil {
		return nil, fmt.Errorf("failed to create attempt for txID: %v: nonce empty", tx.ID)
	}
	legacyTx := evmtypes.LegacyTx{
		Nonce:    *tx.Nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      estimatedGasLimit,
		GasPrice: gasPrice.ToInt(),
		Data:     data,
	}

	signedTx, err := a.keystore.SignTx(ctx, tx.FromAddress, evmtypes.NewTx(&legacyTx), a.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to sign attempt for txID: %v, err: %w", tx.ID, err)
	}

	attempt := &types.Attempt{
		TxID:              tx.ID,
		Fee:               gas.EvmFee{GasPrice: gasPrice},
		Hash:              signedTx.Hash(),
		GasLimit:          estimatedGasLimit,
		Type:              evmtypes.LegacyTxType,
		SignedTransaction: signedTx,
	}

	return attempt, nil
}

func (a *attemptBuilder) newDynamicFeeAttempt(ctx context.Context, tx *types.Transaction, dynamicFee gas.DynamicFee, estimatedGasLimit uint64) (*types.Attempt, error) {
	var data []byte
	var toAddress common.Address
	value := big.NewInt(0)
	if !tx.IsPurgeable {
		data = tx.Data
		toAddress = tx.ToAddress
		value = tx.Value
	}
	if tx.Nonce == nil {
		return nil, fmt.Errorf("failed to create attempt for txID: %v: nonce empty", tx.ID)
	}
	dynamicTx := evmtypes.DynamicFeeTx{
		Nonce:     *tx.Nonce,
		To:        &toAddress,
		Value:     value,
		Gas:       estimatedGasLimit,
		GasFeeCap: dynamicFee.GasFeeCap.ToInt(),
		GasTipCap: dynamicFee.GasTipCap.ToInt(),
		Data:      data,
	}

	signedTx, err := a.keystore.SignTx(ctx, tx.FromAddress, evmtypes.NewTx(&dynamicTx), a.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to sign attempt for txID: %v, err: %w", tx.ID, err)
	}

	attempt := &types.Attempt{
		TxID:              tx.ID,
		Fee:               gas.EvmFee{DynamicFee: gas.DynamicFee{GasFeeCap: dynamicFee.GasFeeCap, GasTipCap: dynamicFee.GasTipCap}},
		Hash:              signedTx.Hash(),
		GasLimit:          estimatedGasLimit,
		Type:              evmtypes.DynamicFeeTxType,
		SignedTransaction: signedTx,
	}

	return attempt, nil
}
