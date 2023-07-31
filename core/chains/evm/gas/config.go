package gas

import (
	"math/big"

	feetypes "github.com/smartcontractkit/chainlink/v2/common/fee/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
)

type bumpConfig interface {
	LimitMultiplier() float32
	PriceMax() *assets.Wei
	BumpPercent() uint16
	BumpMin() *assets.Wei
	TipCapDefault() *assets.Wei
}

type fixedPriceEstimatorConfig interface {
	BumpThreshold() uint64
	FeeCapDefault() *assets.Wei
	LimitMultiplier() float32
	PriceDefault() *assets.Wei
	TipCapDefault() *assets.Wei
	PriceMax() *assets.Wei
	Mode() string
	bumpConfig
}

var _ feetypes.FixedPriceEstimatorConfig = (*wrappedPriceEstimatorConfig)(nil)
var _ feetypes.BumpConfig = (*wrappedPriceEstimatorConfig)(nil)

// wrappedPriceEstimatorConfig is a wrapper that uses *big.Int instead of *assets.Wei
type wrappedPriceEstimatorConfig struct {
	peCfg fixedPriceEstimatorConfig
}

func NewWrappedPriceEstimatorConfig(peCfg fixedPriceEstimatorConfig) wrappedPriceEstimatorConfig {
	return wrappedPriceEstimatorConfig{peCfg: peCfg}
}

func (w wrappedPriceEstimatorConfig) BumpThreshold() uint64 {
	return w.peCfg.BumpThreshold()
}

func (w wrappedPriceEstimatorConfig) FeeCapDefault() *big.Int {
	return w.peCfg.FeeCapDefault().ToInt()
}

func (w wrappedPriceEstimatorConfig) LimitMultiplier() float32 {
	return w.peCfg.LimitMultiplier()
}

func (w wrappedPriceEstimatorConfig) PriceDefault() *big.Int {
	return w.peCfg.PriceDefault().ToInt()
}

func (w wrappedPriceEstimatorConfig) TipCapDefault() *big.Int {
	return w.peCfg.TipCapDefault().ToInt()
}

func (w wrappedPriceEstimatorConfig) PriceMax() *big.Int {
	return w.peCfg.PriceMax().ToInt()
}

func (w wrappedPriceEstimatorConfig) Mode() string {
	return w.peCfg.Mode()
}

func (w wrappedPriceEstimatorConfig) BumpPercent() uint16 {
	return w.peCfg.BumpPercent()
}

func (w wrappedPriceEstimatorConfig) BumpMin() *big.Int {
	return w.peCfg.BumpMin().ToInt()
}
