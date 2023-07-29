package gas

import (
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"math/big"
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

// Implementing LimitMultiplier
func (bw wrappedBumpConfig) LimitMultiplier() float32 {
	return bw.config.LimitMultiplier()
}

// Implementing PriceMax using *big.Int
func (bw wrappedBumpConfig) PriceMax() *big.Int {
	return bw.config.PriceMax().ToInt()
}

// Implementing BumpPercent
func (bw wrappedBumpConfig) BumpPercent() uint16 {
	return bw.config.BumpPercent()
}

// Implementing BumpMin using *big.Int
func (bw wrappedBumpConfig) BumpMin() *big.Int {
	return bw.config.BumpMin().ToInt()
}

// Implementing TipCapDefault using *big.Int
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

// fixedPriceEstimatorConfigWrapper is a wrapper that uses *big.Int instead of *assets.Wei
type fixedPriceEstimatorConfigWrapper struct {
	config fixedPriceEstimatorConfig
}

// Implementing BumpThreshold
func (fw fixedPriceEstimatorConfigWrapper) BumpThreshold() uint64 {
	return fw.config.BumpThreshold()
}

// Implementing FeeCapDefault using *big.Int
func (fw fixedPriceEstimatorConfigWrapper) FeeCapDefault() *big.Int {
	return fw.config.FeeCapDefault().ToInt()
}

// Implementing LimitMultiplier
func (fw fixedPriceEstimatorConfigWrapper) LimitMultiplier() float32 {
	return fw.config.LimitMultiplier()
}

// Implementing PriceDefault using *big.Int
func (fw fixedPriceEstimatorConfigWrapper) PriceDefault() *big.Int {
	return fw.config.PriceDefault().ToInt()
}

// Implementing TipCapDefault using *big.Int
func (fw fixedPriceEstimatorConfigWrapper) TipCapDefault() *big.Int {
	return fw.config.TipCapDefault().ToInt()
}

// Implementing PriceMax using *big.Int
func (fw fixedPriceEstimatorConfigWrapper) PriceMax() *big.Int {
	return fw.config.PriceMax().ToInt()
}

// Implementing Mode
func (fw fixedPriceEstimatorConfigWrapper) Mode() string {
	return fw.config.Mode()
}
