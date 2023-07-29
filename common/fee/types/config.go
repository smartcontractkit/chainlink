package types

import "math/big"

type BumpConfig interface {
	LimitMultiplier() float32
	PriceMax() *big.Int
	BumpPercent() uint16
	BumpMin() *big.Int
	TipCapDefault() *big.Int
}

type FixedPriceEstimatorConfig interface {
	BumpThreshold() uint64
	FeeCapDefault() *big.Int
	LimitMultiplier() float32
	PriceDefault() *big.Int
	TipCapDefault() *big.Int
	PriceMax() *big.Int
	Mode() string
}
