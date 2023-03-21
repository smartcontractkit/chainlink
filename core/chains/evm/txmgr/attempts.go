package txmgr

import (
	"bytes"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	txmgrtypes "github.com/smartcontractkit/chainlink/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
)

// AttemptBuilder takes the base unsigned transaction + optional parameters (tx type, gas parameters)
// and returns a signed TxAttempt
// it is able to estimate fees and sign transactions
type AttemptBuilder interface {
	// interfaces for running the underlying estimator
	services.ServiceCtx
	txmgrtypes.HeadTrackable[*evmtypes.Head]

	// NewAttempt builds a transaction using the configured transaction type and fee estimator (new estimation)
	NewAttempt(ctx context.Context, etx EthTx, lggr logger.Logger, opts ...txmgrtypes.Opt) (attempt EthTxAttempt, fee gas.EvmFee, feeLimit uint32, retryable bool, err error)

	// NewAttemptWithType builds a transaction using the configured fee estimator (new estimation) + passed in tx type
	NewAttemptWithType(ctx context.Context, etx EthTx, lggr logger.Logger, txType int, opts ...txmgrtypes.Opt) (attempt EthTxAttempt, fee gas.EvmFee, feeLimit uint32, retryable bool, err error)

	// NewBumpAttempt builds a transaction using the configured fee estimator (bumping) + passed in tx type
	// this should only be used after an initial attempt has been broadcast and the underlying gas estimator only needs to bump the fee
	NewBumpAttempt(ctx context.Context, etx EthTx, previousAttempt EthTxAttempt, txType int, priorAttempts []txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash], lggr logger.Logger) (attempt EthTxAttempt, bumpedFee gas.EvmFee, bumpedFeeLimit uint32, retryable bool, err error)

	// NewCustomAttempt builds a transaction using the passed in fee + tx type
	NewCustomAttempt(etx EthTx, fee gas.EvmFee, gasLimit uint32, txType int, lggr logger.Logger) (attempt EthTxAttempt, retryable bool, err error)

	// FeeEstimator returns the underlying gas estimator
	FeeEstimator() txmgrtypes.FeeEstimator[*evmtypes.Head, gas.EvmFee, *assets.Wei, common.Hash]
}

var _ AttemptBuilder = (*evmAttemptBuilder)(nil)

type evmAttemptBuilder struct {
	chainID  big.Int
	config   Config
	keystore KeyStore
	gas.EvmFeeEstimator
}

func NewEvmAttemptBuilder(chainID big.Int, config Config, keystore KeyStore, estimator gas.EvmFeeEstimator) *evmAttemptBuilder {
	return &evmAttemptBuilder{chainID, config, keystore, estimator}
}

func (c *evmAttemptBuilder) FeeEstimator() txmgrtypes.FeeEstimator[*evmtypes.Head, gas.EvmFee, *assets.Wei, common.Hash] {
	return c.EvmFeeEstimator
}

func (c *evmAttemptBuilder) SignTx(address common.Address, tx *gethTypes.Transaction) (common.Hash, []byte, error) {
	signedTx, err := c.keystore.SignTx(address, tx, &c.chainID)
	if err != nil {
		return common.Hash{}, nil, errors.Wrap(err, "SignTx failed")
	}
	rlp := new(bytes.Buffer)
	if err = signedTx.EncodeRLP(rlp); err != nil {
		return common.Hash{}, nil, errors.Wrap(err, "SignTx failed")
	}
	return signedTx.Hash(), rlp.Bytes(), nil
}

func (c *evmAttemptBuilder) NewAttempt(ctx context.Context, etx EthTx, lggr logger.Logger, opts ...txmgrtypes.Opt) (attempt EthTxAttempt, fee gas.EvmFee, feeLimit uint32, retryable bool, err error) {
	txType := 0x0
	if c.config.EvmEIP1559DynamicFees() {
		txType = 0x2
	}
	return c.NewAttemptWithType(ctx, etx, lggr, txType, opts...)
}

func (c *evmAttemptBuilder) NewAttemptWithType(ctx context.Context, etx EthTx, lggr logger.Logger, txType int, opts ...txmgrtypes.Opt) (attempt EthTxAttempt, fee gas.EvmFee, feeLimit uint32, retryable bool, err error) {
	keySpecificMaxGasPriceWei := c.config.KeySpecificMaxGasPriceWei(etx.FromAddress)
	fee, feeLimit, err = c.EvmFeeEstimator.GetFee(ctx, etx.EncodedPayload, etx.GasLimit, keySpecificMaxGasPriceWei, opts...)
	if err != nil {
		return attempt, fee, feeLimit, true, errors.Wrap(err, "failed to get fee") // estimator errors are retryable
	}

	attempt, retryable, err = c.NewCustomAttempt(etx, fee, feeLimit, txType, lggr)
	return attempt, fee, feeLimit, retryable, err
}

func (c *evmAttemptBuilder) NewBumpAttempt(ctx context.Context, etx EthTx, previousAttempt EthTxAttempt, txType int, priorAttempts []txmgrtypes.PriorAttempt[gas.EvmFee, common.Hash], lggr logger.Logger) (attempt EthTxAttempt, bumpedFee gas.EvmFee, bumpedFeeLimit uint32, retryable bool, err error) {
	keySpecificMaxGasPriceWei := c.config.KeySpecificMaxGasPriceWei(etx.FromAddress)
	bumpedFee, bumpedFeeLimit, err = c.EvmFeeEstimator.BumpFee(ctx, previousAttempt.Fee(), etx.GasLimit, keySpecificMaxGasPriceWei, priorAttempts)
	if err != nil {
		return attempt, bumpedFee, bumpedFeeLimit, true, errors.Wrap(err, "failed to bump fee") // estimator errors are retryable
	}

	attempt, retryable, err = c.NewCustomAttempt(etx, bumpedFee, bumpedFeeLimit, txType, lggr)
	return attempt, bumpedFee, bumpedFeeLimit, retryable, err
}

func (c *evmAttemptBuilder) NewCustomAttempt(etx EthTx, fee gas.EvmFee, gasLimit uint32, txType int, lggr logger.Logger) (attempt EthTxAttempt, retryable bool, err error) {
	switch txType {
	case 0x0: // legacy
		if fee.Legacy == nil {
			err = errors.Errorf("Attempt %v is a type 0 transaction but estimator did not return legacy fee bump", attempt.ID)
			logger.Sugared(lggr).AssumptionViolation(err.Error())
			return attempt, false, err // not retryable
		}
		attempt, err = c.newLegacyAttempt(etx, fee.Legacy, gasLimit)
		return attempt, true, err
	case 0x2: // dynamic, EIP1559
		if fee.Dynamic == nil {
			err = errors.Errorf("Attempt %v is a type 2 transaction but estimator did not return dynamic fee bump", attempt.ID)
			logger.Sugared(lggr).AssumptionViolation(err.Error())
			return attempt, false, err // not retryable
		}
		attempt, err = c.newDynamicFeeAttempt(etx, *fee.Dynamic, gasLimit)
		return attempt, true, err
	default:
		err = errors.Errorf("invariant violation: Attempt %v had unrecognised transaction type %v"+
			"This is a bug! Please report to https://github.com/smartcontractkit/chainlink/issues", attempt.ID, attempt.TxType)
		logger.Sugared(lggr).AssumptionViolation(err.Error())
		return attempt, false, err // not retryable
	}
}

func (c *evmAttemptBuilder) newDynamicFeeAttempt(etx EthTx, fee gas.DynamicFee, gasLimit uint32) (attempt EthTxAttempt, err error) {
	if err = validateDynamicFeeGas(c.config, fee, gasLimit, etx); err != nil {
		return attempt, errors.Wrap(err, "error validating gas")
	}

	var al types.AccessList
	if etx.AccessList.Valid {
		al = etx.AccessList.AccessList
	}
	d := newDynamicFeeTransaction(
		uint64(*etx.Nonce),
		etx.ToAddress,
		&etx.Value,
		gasLimit,
		&c.chainID,
		fee.TipCap,
		fee.FeeCap,
		etx.EncodedPayload,
		al,
	)
	tx := types.NewTx(&d)
	attempt, err = c.newSignedAttempt(etx, tx)
	if err != nil {
		return attempt, err
	}
	attempt.GasTipCap = fee.TipCap
	attempt.GasFeeCap = fee.FeeCap
	attempt.ChainSpecificGasLimit = gasLimit
	attempt.TxType = 2
	return attempt, nil
}

var Max256BitUInt = big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil)

// validateDynamicFeeGas is a sanity check - we have other checks elsewhere, but this
// makes sure we _never_ create an invalid attempt
func validateDynamicFeeGas(cfg Config, fee gas.DynamicFee, gasLimit uint32, etx EthTx) error {
	gasTipCap, gasFeeCap := fee.TipCap, fee.FeeCap

	if gasTipCap == nil {
		panic("gas tip cap missing")
	}
	if gasFeeCap == nil {
		panic("gas fee cap missing")
	}
	// Assertions from:	https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1559.md
	// Prevent impossibly large numbers
	if gasFeeCap.ToInt().Cmp(Max256BitUInt) > 0 {
		return errors.New("impossibly large fee cap")
	}
	if gasTipCap.ToInt().Cmp(Max256BitUInt) > 0 {
		return errors.New("impossibly large tip cap")
	}
	// The total must be at least as large as the tip
	if gasFeeCap.Cmp(gasTipCap) < 0 {
		return errors.Errorf("gas fee cap must be greater than or equal to gas tip cap (fee cap: %s, tip cap: %s)", gasFeeCap.String(), gasTipCap.String())
	}

	// Configuration sanity-check
	max := cfg.KeySpecificMaxGasPriceWei(etx.FromAddress)
	if gasFeeCap.Cmp(max) > 0 {
		return errors.Errorf("cannot create tx attempt: specified gas fee cap of %s would exceed max configured gas price of %s for key %s", gasFeeCap.String(), max.String(), etx.FromAddress.Hex())
	}
	// Tip must be above minimum
	minTip := cfg.EvmGasTipCapMinimum()
	if gasTipCap.Cmp(minTip) < 0 {
		return errors.Errorf("cannot create tx attempt: specified gas tip cap of %s is below min configured gas tip of %s for key %s", gasTipCap.String(), minTip.String(), etx.FromAddress.Hex())
	}
	return nil
}

func newDynamicFeeTransaction(nonce uint64, to common.Address, value *assets.Eth, gasLimit uint32, chainID *big.Int, gasTipCap, gasFeeCap *assets.Wei, data []byte, accessList types.AccessList) types.DynamicFeeTx {
	return types.DynamicFeeTx{
		ChainID:    chainID,
		Nonce:      nonce,
		GasTipCap:  gasTipCap.ToInt(),
		GasFeeCap:  gasFeeCap.ToInt(),
		Gas:        uint64(gasLimit),
		To:         &to,
		Value:      value.ToInt(),
		Data:       data,
		AccessList: accessList,
	}
}

func (c *evmAttemptBuilder) newLegacyAttempt(etx EthTx, gasPrice *assets.Wei, gasLimit uint32) (attempt EthTxAttempt, err error) {
	if err = validateLegacyGas(c.config, gasPrice, gasLimit, etx); err != nil {
		return attempt, errors.Wrap(err, "error validating gas")
	}

	tx := newLegacyTransaction(
		uint64(*etx.Nonce),
		etx.ToAddress,
		etx.Value.ToInt(),
		gasLimit,
		gasPrice,
		etx.EncodedPayload,
	)

	transaction := types.NewTx(&tx)
	hash, signedTxBytes, err := c.SignTx(etx.FromAddress, transaction)
	if err != nil {
		return attempt, errors.Wrapf(err, "error using account %s to sign transaction %v", etx.FromAddress.String(), etx.ID)
	}

	attempt.State = EthTxAttemptInProgress
	attempt.SignedRawTx = signedTxBytes
	attempt.EthTxID = etx.ID
	attempt.GasPrice = gasPrice
	attempt.Hash = hash
	attempt.TxType = 0
	attempt.ChainSpecificGasLimit = gasLimit
	attempt.EthTx = etx

	return attempt, nil
}

// validateLegacyGas is a sanity check - we have other checks elsewhere, but this
// makes sure we _never_ create an invalid attempt
func validateLegacyGas(cfg Config, gasPrice *assets.Wei, gasLimit uint32, etx EthTx) error {
	if gasPrice == nil {
		panic("gas price missing")
	}
	max := cfg.KeySpecificMaxGasPriceWei(etx.FromAddress)
	if gasPrice.Cmp(max) > 0 {
		return errors.Errorf("cannot create tx attempt: specified gas price of %s would exceed max configured gas price of %s for key %s", gasPrice.String(), max.String(), etx.FromAddress.Hex())
	}
	min := cfg.EvmMinGasPriceWei()
	if gasPrice.Cmp(min) < 0 {
		return errors.Errorf("cannot create tx attempt: specified gas price of %s is below min configured gas price of %s for key %s", gasPrice.String(), min.String(), etx.FromAddress.Hex())
	}
	return nil
}

func (c *evmAttemptBuilder) newSignedAttempt(etx EthTx, tx *types.Transaction) (attempt EthTxAttempt, err error) {
	hash, signedTxBytes, err := c.signTx(etx.FromAddress, tx)
	if err != nil {
		return attempt, errors.Wrapf(err, "error using account %s to sign transaction %v", etx.FromAddress.String(), etx.ID)
	}

	attempt.State = EthTxAttemptInProgress
	attempt.SignedRawTx = signedTxBytes
	attempt.EthTxID = etx.ID
	attempt.EthTx = etx
	attempt.Hash = hash

	return attempt, nil
}

func newLegacyTransaction(nonce uint64, to common.Address, value *big.Int, gasLimit uint32, gasPrice *assets.Wei, data []byte) types.LegacyTx {
	return types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      uint64(gasLimit),
		GasPrice: gasPrice.ToInt(),
		Data:     data,
	}
}

func (c *evmAttemptBuilder) signTx(address common.Address, tx *types.Transaction) (common.Hash, []byte, error) {
	signedTx, err := c.keystore.SignTx(address, tx, &c.chainID)
	if err != nil {
		return common.Hash{}, nil, errors.Wrap(err, "signTx failed")
	}
	rlp := new(bytes.Buffer)
	if err := signedTx.EncodeRLP(rlp); err != nil {
		return common.Hash{}, nil, errors.Wrap(err, "signTx failed")
	}
	return signedTx.Hash(), rlp.Bytes(), nil
}
