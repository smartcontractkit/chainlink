package txmgr

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	types2 "github.com/smartcontractkit/chainlink/common/types"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/gas"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
)

func (c *ChainKeyStore[ADDR, TX_HASH]) NewDynamicFeeAttempt(etx EthTx[ADDR, TX_HASH], fee gas.DynamicFee, gasLimit uint32) (attempt EthTxAttempt[ADDR, TX_HASH], err error) {
	if err = validateDynamicFeeGas(c.config, fee, gasLimit, etx); err != nil {
		return attempt, errors.Wrap(err, "error validating gas")
	}

	var al types.AccessList
	if etx.AccessList.Valid {
		al = etx.AccessList.AccessList
	}
	nativeToAddress := *etx.ToAddress.(*evmtypes.Address).NativeAddress()
	d := newDynamicFeeTransaction(
		uint64(*etx.Nonce),
		nativeToAddress,
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
func validateDynamicFeeGas[ADDR types2.Hashable, TX_HASH types2.Hashable](cfg Config, fee gas.DynamicFee, gasLimit uint32, etx EthTx[ADDR, TX_HASH]) error {
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
	nativeFromAddress := *etx.FromAddress.(*evmtypes.Address).NativeAddress()
	max := cfg.KeySpecificMaxGasPriceWei(nativeFromAddress)
	if gasFeeCap.Cmp(max) > 0 {
		return errors.Errorf("cannot create tx attempt: specified gas fee cap of %s would exceed max configured gas price of %s for key %s", gasFeeCap.String(), max.String(), nativeFromAddress.Hex())
	}
	// Tip must be above minimum
	minTip := cfg.EvmGasTipCapMinimum()
	if gasTipCap.Cmp(minTip) < 0 {
		return errors.Errorf("cannot create tx attempt: specified gas tip cap of %s is below min configured gas tip of %s for key %s", gasTipCap.String(), minTip.String(), nativeFromAddress.Hex())
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

func (c *ChainKeyStore[ADDR, TX_HASH]) NewLegacyAttempt(etx EthTx[ADDR, TX_HASH], gasPrice *assets.Wei, gasLimit uint32) (attempt EthTxAttempt[ADDR, TX_HASH], err error) {
	if err = validateLegacyGas(c.config, gasPrice, gasLimit, *etx.FromAddress.(*evmtypes.Address).NativeAddress()); err != nil {
		return attempt, errors.Wrap(err, "error validating gas")
	}

	tx := newLegacyTransaction(
		uint64(*etx.Nonce),
		*etx.ToAddress.(*evmtypes.Address).NativeAddress(),
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
	attempt.Hash = *hash
	attempt.TxType = 0
	attempt.ChainSpecificGasLimit = gasLimit
	attempt.EthTx = etx

	return attempt, nil
}

// validateLegacyGas is a sanity check - we have other checks elsewhere, but this
// makes sure we _never_ create an invalid attempt
func validateLegacyGas(cfg Config, gasPrice *assets.Wei, gasLimit uint32, fromAddress common.Address) error {
	if gasPrice == nil {
		panic("gas price missing")
	}
	max := cfg.KeySpecificMaxGasPriceWei(fromAddress)
	if gasPrice.Cmp(max) > 0 {
		return errors.Errorf("cannot create tx attempt: specified gas price of %s would exceed max configured gas price of %s for key %s", gasPrice.String(), max.String(), fromAddress.String())
	}
	min := cfg.EvmMinGasPriceWei()
	if gasPrice.Cmp(min) < 0 {
		return errors.Errorf("cannot create tx attempt: specified gas price of %s is below min configured gas price of %s for key %s", gasPrice.String(), min.String(), fromAddress.String())
	}
	return nil
}

func (c *ChainKeyStore[ADDR, TX_HASH]) newSignedAttempt(etx EthTx[ADDR, TX_HASH], tx *types.Transaction) (attempt EthTxAttempt[ADDR, TX_HASH], err error) {
	hash, signedTxBytes, err := c.signTx(etx.FromAddress, tx)
	if err != nil {
		return attempt, errors.Wrapf(err, "error using account %s to sign transaction %v", etx.FromAddress.String(), etx.ID)
	}

	attempt.State = EthTxAttemptInProgress
	attempt.SignedRawTx = signedTxBytes
	attempt.EthTxID = etx.ID
	attempt.EthTx = etx
	attempt.Hash = *hash

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

func (c *ChainKeyStore[ADDR, TX_HASH]) signTx(address ADDR, tx *types.Transaction) (*TX_HASH, []byte, error) {
	// Native EVM types used here will be removed in later PRs.
	signedTx, txHash, err := c.keystore.SignTx(address, tx, &c.chainID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "signTx failed")
	}
	rlp := new(bytes.Buffer)
	if err := signedTx.EncodeRLP(rlp); err != nil {
		return nil, nil, errors.Wrap(err, "signTx failed")
	}
	return txHash, rlp.Bytes(), err
}

func getEvmAddress(h types2.Hashable) *evmtypes.Address {
	// The runtime cast here will be removed when the Config interface start implementing Hashable.
	return h.(*evmtypes.Address)
}
