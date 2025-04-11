package gasprice

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink/v2/core/cbor"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
)

const (
	// feeLimit is a placeholder when getting current price from gas estimator. it does not impact gas price calculation
	feeLimit = uint64(1_000_000)
	// maxFeePrice is a placeholder when getting current price from gas estimator. it caps the returned gas price from
	// the estimator. it's set to a very high value because the gas price will be compared with user-defined gas price
	// later.
	maxFeePrice = 1_000_000_000_000_000
)

type UpkeepOffchainConfig struct {
	MaxGasPrice *big.Int `json:"maxGasPrice" cbor:"maxGasPrice"`
}

// CheckGasPrice retrieves the current gas price and compare against the max gas price configured in upkeep's offchain config
// any errors in offchain config decoding will result in max gas price check disabled
func CheckGasPrice(ctx context.Context, upkeepId *big.Int, offchainConfigBytes []byte, ge gas.EvmFeeEstimator, lggr logger.Logger) encoding.UpkeepFailureReason {
	// check for empty offchain config
	if len(offchainConfigBytes) == 0 {
		return encoding.UpkeepFailureReasonNone
	}

	var offchainConfig UpkeepOffchainConfig
	if err := cbor.ParseDietCBORToStruct(offchainConfigBytes, &offchainConfig); err != nil {
		lggr.Warnw("failed to parse upkeep offchain config, gas price check is disabled", "offchainconfig", hexutil.Encode(offchainConfigBytes), "upkeepId", upkeepId.String(), "err", err)
		return encoding.UpkeepFailureReasonNone
	}
	if offchainConfig.MaxGasPrice == nil || offchainConfig.MaxGasPrice.Int64() <= 0 {
		lggr.Debugw("maxGasPrice is not configured or incorrectly configured in upkeep offchain config, gas price check is disabled", "offchainconfig", hexutil.Encode(offchainConfigBytes), "upkeepId", upkeepId.String())
		return encoding.UpkeepFailureReasonNone
	}
	lggr.Debugf("successfully decode offchain config for %s, max gas price is %s", upkeepId.String(), offchainConfig.MaxGasPrice.String())

	fee, _, err := ge.GetFee(ctx, []byte{}, feeLimit, assets.NewWei(big.NewInt(maxFeePrice)), nil, nil)
	if err != nil {
		lggr.Errorw("failed to get fee, gas price check is disabled", "upkeepId", upkeepId.String(), "err", err)
		return encoding.UpkeepFailureReasonNone
	}

	if fee.ValidDynamic() {
		lggr.Debugf("current gas price EIP-1559 is fee cap %s, tip cap %s", fee.GasFeeCap.String(), fee.GasTipCap.String())
		if fee.GasFeeCap.Cmp(assets.NewWei(offchainConfig.MaxGasPrice)) > 0 {
			// current gas price is higher than max gas price
			lggr.Warnf("maxGasPrice %s for %s is LOWER than current gas price %d", offchainConfig.MaxGasPrice.String(), upkeepId.String(), fee.GasFeeCap.Int64())
			return encoding.UpkeepFailureReasonGasPriceTooHigh
		}
		lggr.Debugf("maxGasPrice %s for %s is HIGHER than current gas price %d", offchainConfig.MaxGasPrice.String(), upkeepId.String(), fee.GasFeeCap.Int64())
	} else {
		lggr.Debugf("current gas price legacy is %s", fee.GasPrice.String())
		if fee.GasPrice.Cmp(assets.NewWei(offchainConfig.MaxGasPrice)) > 0 {
			// current gas price is higher than max gas price
			lggr.Warnf("maxGasPrice %s for %s is LOWER than current gas price %d", offchainConfig.MaxGasPrice.String(), upkeepId.String(), fee.GasPrice.Int64())
			return encoding.UpkeepFailureReasonGasPriceTooHigh
		}
		lggr.Debugf("maxGasPrice %s for %s is HIGHER than current gas price %d", offchainConfig.MaxGasPrice.String(), upkeepId.String(), fee.GasPrice.Int64())
	}

	return encoding.UpkeepFailureReasonNone
}
