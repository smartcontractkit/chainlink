package bulletprooftxmanager

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/gas"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func NewDynamicFeeAttempt(cfg Config, ks KeyStore, chainID *big.Int, etx EthTx, fee gas.DynamicFee, gasLimit uint64) (attempt EthTxAttempt, err error) {
	if err = validateDynamicFeeGas(cfg, fee, gasLimit, etx); err != nil {
		return attempt, errors.Wrap(err, "error validating gas")
	}

	var al types.AccessList
	if etx.AccessList.Valid {
		al = etx.AccessList.AccessList
	}
	d := newDynamicFeeTransaction(
		uint64(*etx.Nonce),
		etx.ToAddress,
		etx.Value.ToInt(),
		gasLimit,
		chainID,
		fee.TipCap,
		fee.FeeCap,
		etx.EncodedPayload,
		al,
	)
	tx := types.NewTx(&d)
	attempt, err = newSignedAttempt(ks, etx, chainID, tx)
	if err != nil {
		return attempt, err
	}
	attempt.GasTipCap = utils.NewBig(fee.TipCap)
	attempt.GasFeeCap = utils.NewBig(fee.FeeCap)
	attempt.ChainSpecificGasLimit = gasLimit
	attempt.TxType = 2
	return attempt, nil
}

var Max256BitUInt = big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil)

// validateDynamicFeeGas is a sanity check - we have other checks elsewhere, but this
// makes sure we _never_ create an invalid attempt
func validateDynamicFeeGas(cfg Config, fee gas.DynamicFee, gasLimit uint64, etx EthTx) error {
	gasTipCap, gasFeeCap := fee.TipCap, fee.FeeCap

	if gasTipCap == nil {
		panic("gas tip cap missing")
	}
	if gasFeeCap == nil {
		panic("gas fee cap missing")
	}
	// Assertions from:	https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1559.md
	// Prevent impossibly large numbers
	if gasFeeCap.Cmp(Max256BitUInt) > 0 {
		return errors.New("impossibly large fee cap")
	}
	if gasTipCap.Cmp(Max256BitUInt) > 0 {
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

func newDynamicFeeTransaction(nonce uint64, to common.Address, value *big.Int, gasLimit uint64, chainID, gasTipCap, gasFeeCap *big.Int, data []byte, accessList types.AccessList) types.DynamicFeeTx {
	return types.DynamicFeeTx{
		ChainID:    chainID,
		Nonce:      nonce,
		GasTipCap:  gasTipCap,
		GasFeeCap:  gasFeeCap,
		Gas:        gasLimit,
		To:         &to,
		Value:      value,
		Data:       data,
		AccessList: accessList,
	}
}

func NewLegacyAttempt(cfg Config, ks KeyStore, chainID *big.Int, etx EthTx, gasPrice *big.Int, gasLimit uint64) (attempt EthTxAttempt, err error) {
	if err = validateLegacyGas(cfg, gasPrice, gasLimit, etx); err != nil {
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
	hash, signedTxBytes, err := SignTx(ks, etx.FromAddress, transaction, chainID)
	if err != nil {
		return attempt, errors.Wrapf(err, "error using account %s to sign transaction %v", etx.FromAddress.String(), etx.ID)
	}

	attempt.State = EthTxAttemptInProgress
	attempt.SignedRawTx = signedTxBytes
	attempt.EthTxID = etx.ID
	attempt.GasPrice = utils.NewBig(gasPrice)
	attempt.Hash = hash
	attempt.TxType = 0
	attempt.ChainSpecificGasLimit = gasLimit

	return attempt, nil
}

// validateLegacyGas is a sanity check - we have other checks elsewhere, but this
// makes sure we _never_ create an invalid attempt
func validateLegacyGas(cfg Config, gasPrice *big.Int, gasLimit uint64, etx EthTx) error {
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

func newSignedAttempt(ks KeyStore, etx EthTx, chainID *big.Int, tx *types.Transaction) (attempt EthTxAttempt, err error) {
	hash, signedTxBytes, err := signTx(ks, etx.FromAddress, tx, chainID)
	if err != nil {
		return attempt, errors.Wrapf(err, "error using account %s to sign transaction %v", etx.FromAddress.String(), etx.ID)
	}

	attempt.State = EthTxAttemptInProgress
	attempt.SignedRawTx = signedTxBytes
	attempt.EthTxID = etx.ID
	attempt.Hash = hash

	return attempt, nil
}

func newLegacyTransaction(nonce uint64, to common.Address, value *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) types.LegacyTx {
	return types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	}
}

func signTx(keyStore KeyStore, address common.Address, tx *types.Transaction, chainID *big.Int) (common.Hash, []byte, error) {
	signedTx, err := keyStore.SignTx(address, tx, chainID)
	if err != nil {
		return common.Hash{}, nil, errors.Wrap(err, "signTx failed")
	}
	rlp := new(bytes.Buffer)
	if err := signedTx.EncodeRLP(rlp); err != nil {
		return common.Hash{}, nil, errors.Wrap(err, "signTx failed")
	}
	return signedTx.Hash(), rlp.Bytes(), nil

}
