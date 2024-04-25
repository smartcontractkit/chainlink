package gasprice

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/cbor"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
)

const (
	// feeLimit is a placeholder when getting current price from gas estimator. it does not impact gas price calculation
	feeLimit              = uint64(1_000_000)
	// maxFeePrice is a placeholder when getting current price from gas estimator. it caps the returned gas price from
	// the estimator. it's set to a very high value because the gas price will be compared with user-defined gas price
	// later.
	maxFeePrice           = 1_000_000_000_000_000
)

type UpkeepOffchainConfig struct {
	MaxGasPrice *big.Int `json:"maxGasPrice" cbor:"maxGasPrice"`
}

func CheckGasPrice(ctx context.Context, upkeepId *big.Int, oc []byte, ge gas.EvmFeeEstimator, lggr logger.Logger) encoding.UpkeepFailureReason {
	if len(oc) == 0 {
		return encoding.UpkeepFailureReasonNone
	}

	var offchainConfig UpkeepOffchainConfig
	if err := cbor.ParseDietCBORToStruct(oc, &offchainConfig); err != nil {
		lggr.Errorw("failed to parse upkeep offchain config", "upkeepId", upkeepId.String(), "err", err)
		return encoding.UpkeepFailureReasonFailToParseOffchainConfig
	}
	if offchainConfig.MaxGasPrice == nil {
		lggr.Infow("maxGasPrice is not configured in upkeep offchain config", "upkeepId", upkeepId.String())
		return encoding.UpkeepFailureReasonNone
	}
	lggr.Infof("successfully decode offchain config for %s", upkeepId.String())
	lggr.Infof("max gas price for %s is %s", upkeepId.String(), offchainConfig.MaxGasPrice.String())

	fee, _, err := ge.GetFee(ctx, []byte{}, feeLimit, assets.NewWei(big.NewInt(maxFeePrice)))
	if err != nil {
		lggr.Errorw("failed to get fee", "upkeepId", upkeepId.String(), "err", err)
		return encoding.UpkeepFailureReasonFailToRetrieveGasPrice
	}

	if fee.ValidDynamic() {
		lggr.Infof("current gas price EIP-1559 is fee cap %s, tip cap %s", fee.DynamicFeeCap.String(), fee.DynamicTipCap.String())
		if fee.DynamicFeeCap.Cmp(assets.NewWei(offchainConfig.MaxGasPrice)) > 0 {
			// current gas price is higher than max gas price
			lggr.Warnf("max gas price %s for %s is LOWER than current gas price %s", offchainConfig.MaxGasPrice.String(), upkeepId.String(), fee.DynamicFeeCap.String())
			return encoding.UpkeepFailureReasonGasPriceTooHigh
		} else {
			lggr.Infof("max gas price %s for %s is HIGHER than current gas price %s", offchainConfig.MaxGasPrice.String(), upkeepId.String(), fee.DynamicFeeCap.String())
		}
	} else {
		lggr.Infof("current gas price legacy is %s", fee.Legacy.String())
		if fee.Legacy.Cmp(assets.NewWei(offchainConfig.MaxGasPrice)) > 0 {
			// current gas price is higher than max gas price
			lggr.Infof("max gas price %s for %s is LOWER than current gas price %s", offchainConfig.MaxGasPrice.String(), upkeepId.String(), fee.Legacy.String())
			return encoding.UpkeepFailureReasonGasPriceTooHigh
		} else {
			lggr.Infof("max gas price %s for %s is HIGHER than current gas price %s", offchainConfig.MaxGasPrice.String(), upkeepId.String(), fee.Legacy.String())
		}
	}

	return encoding.UpkeepFailureReasonNone
}