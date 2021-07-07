package gas

import (
	"context"
	"math/big"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/eth"
)

var _ Estimator = &fixedPriceEstimator{}

type fixedPriceEstimator struct {
	config Config
}

func NewFixedPriceEstimator(config Config) Estimator {
	return &fixedPriceEstimator{config}
}

func (f *fixedPriceEstimator) Start() error                                    { return nil }
func (f *fixedPriceEstimator) Close() error                                    { return nil }
func (f *fixedPriceEstimator) OnNewLongestChain(_ context.Context, _ eth.Head) {}

func (f *fixedPriceEstimator) GetLegacyGas(_ []byte, gasLimit uint64, _ ...Opt) (gasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
	gasPrice = f.config.EvmGasPriceDefault()
	chainSpecificGasLimit = applyMultiplier(gasLimit, f.config.EvmGasLimitMultiplier())
	return
}

func (f *fixedPriceEstimator) BumpLegacyGas(originalGasPrice *big.Int, originalGasLimit uint64) (gasPrice *big.Int, gasLimit uint64, err error) {
	return BumpLegacyGasPriceOnly(f.config, originalGasPrice, originalGasLimit)
}

func (f *fixedPriceEstimator) GetDynamicFee(originalGasLimit uint64) (d DynamicFee, chainSpecificGasLimit uint64, err error) {
	gasTipCap := f.config.EvmGasTipCapDefault()
	if gasTipCap == nil {
		return d, 0, errors.New("cannot calculate dynamic fee: EthGasTipCapDefault was not set")
	}
	chainSpecificGasLimit = applyMultiplier(originalGasLimit, f.config.EvmGasLimitMultiplier())
	return DynamicFee{
		FeeCap: f.config.EvmGasFeeCap(),
		TipCap: gasTipCap,
	}, chainSpecificGasLimit, nil
}

func (f *fixedPriceEstimator) BumpDynamicFee(originalFee DynamicFee, originalGasLimit uint64) (bumped DynamicFee, chainSpecificGasLimit uint64, err error) {
	return BumpDynamicFeeOnly(f.config, originalFee, originalGasLimit)
}
