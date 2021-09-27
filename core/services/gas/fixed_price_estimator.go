package gas

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

var _ Estimator = &fixedPriceEstimator{}

type fixedPriceEstimator struct {
	config Config
}

func NewFixedPriceEstimator(config Config) Estimator {
	return &fixedPriceEstimator{config}
}

func (f *fixedPriceEstimator) Start() error                                       { return nil }
func (f *fixedPriceEstimator) Close() error                                       { return nil }
func (f *fixedPriceEstimator) OnNewLongestChain(_ context.Context, _ models.Head) {}

func (f *fixedPriceEstimator) EstimateGas(_ []byte, gasLimit uint64, _ ...Opt) (gasPrice *big.Int, chainSpecificGasLimit uint64, err error) {
	gasPrice = f.config.EvmGasPriceDefault()
	chainSpecificGasLimit = applyMultiplier(gasLimit, f.config.EvmGasLimitMultiplier())
	return
}

func (f *fixedPriceEstimator) BumpGas(originalGasPrice *big.Int, originalGasLimit uint64) (gasPrice *big.Int, gasLimit uint64, err error) {
	return BumpGasPriceOnly(f.config, f.config.EvmGasPriceDefault(), originalGasPrice, originalGasLimit)
}
