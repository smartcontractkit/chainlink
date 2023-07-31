package gas

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
)

type bumpConfig interface {
	LimitMultiplier() float32
	PriceMax() *assets.Wei
	BumpPercent() uint16
	BumpMin() *assets.Wei
	TipCapDefault() *assets.Wei
}

// wrappedBumpConfig is a wrapper that uses *big.Int instead of *assets.Wei
type wrappedBumpConfig struct {
	config bumpConfig
}

func NewWrappedBumpConfig(config bumpConfig) wrappedBumpConfig {
	return wrappedBumpConfig{config: config}
}

func (bw wrappedBumpConfig) LimitMultiplier() float32 {
	return bw.config.LimitMultiplier()
}

func (bw wrappedBumpConfig) PriceMax() *big.Int {
	return bw.config.PriceMax().ToInt()
}

func (bw wrappedBumpConfig) BumpPercent() uint16 {
	return bw.config.BumpPercent()
}

func (bw wrappedBumpConfig) BumpMin() *big.Int {
	return bw.config.BumpMin().ToInt()
}

func (bw wrappedBumpConfig) TipCapDefault() *big.Int {
	return bw.config.TipCapDefault().ToInt()
}

type fixedPriceEstimatorConfig interface {
	BumpThreshold() uint64
	FeeCapDefault() *assets.Wei
	LimitMultiplier() float32
	PriceDefault() *assets.Wei
	TipCapDefault() *assets.Wei
	PriceMax() *assets.Wei
	Mode() string
}

// wrappedPriceEstimatorConfig is a wrapper that uses *big.Int instead of *assets.Wei
type wrappedPriceEstimatorConfig struct {
	config fixedPriceEstimatorConfig
}

func NewWrappedPriceEstimatorConfig(config fixedPriceEstimatorConfig) wrappedPriceEstimatorConfig {
	return wrappedPriceEstimatorConfig{config: config}
}

func (fw wrappedPriceEstimatorConfig) BumpThreshold() uint64 {
	return fw.config.BumpThreshold()
}

func (fw wrappedPriceEstimatorConfig) FeeCapDefault() *big.Int {
	return fw.config.FeeCapDefault().ToInt()
}

func (fw wrappedPriceEstimatorConfig) LimitMultiplier() float32 {
	return fw.config.LimitMultiplier()
}

func (fw wrappedPriceEstimatorConfig) PriceDefault() *big.Int {
	return fw.config.PriceDefault().ToInt()
}

func (fw wrappedPriceEstimatorConfig) TipCapDefault() *big.Int {
	return fw.config.TipCapDefault().ToInt()
}

func (fw wrappedPriceEstimatorConfig) PriceMax() *big.Int {
	return fw.config.PriceMax().ToInt()
}

func (fw wrappedPriceEstimatorConfig) Mode() string {
	return fw.config.Mode()
}
