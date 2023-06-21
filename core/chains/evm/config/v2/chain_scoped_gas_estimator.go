package v2

import (
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
)

type gasEstimatorConfig struct {
	c                       GasEstimator
	blockDelay              *uint16
	transactionsMaxInFlight *uint32
}

func (g *gasEstimatorConfig) BlockHistory() config.BlockHistory {
	return &blockHistoryConfig{c: g.c.BlockHistory, blockDelay: g.blockDelay, bumpThreshold: g.c.BumpThreshold}
}

func (g *gasEstimatorConfig) EIP1559DynamicFees() bool {
	return *g.c.EIP1559DynamicFees
}

func (g *gasEstimatorConfig) BumpPercent() uint16 {
	return *g.c.BumpPercent
}

func (g *gasEstimatorConfig) BumpThreshold() uint64 {
	return uint64(*g.c.BumpThreshold)
}

func (g *gasEstimatorConfig) BumpTxDepth() uint32 {
	if g.c.BumpTxDepth != nil {
		return *g.c.BumpTxDepth
	}
	return *g.transactionsMaxInFlight
}

func (g *gasEstimatorConfig) BumpMin() *assets.Wei {
	return g.c.BumpMin
}

func (g *gasEstimatorConfig) FeeCapDefault() *assets.Wei {
	return g.c.FeeCapDefault
}

func (g *gasEstimatorConfig) LimitDefault() uint32 {
	return *g.c.LimitDefault
}

func (g *gasEstimatorConfig) LimitMax() uint32 {
	return *g.c.LimitMax
}

func (g *gasEstimatorConfig) LimitMultiplier() float32 {
	f, _ := g.c.LimitMultiplier.BigFloat().Float32()
	return f
}

func (g *gasEstimatorConfig) LimitTransfer() uint32 {
	return *g.c.LimitTransfer
}

func (g *gasEstimatorConfig) PriceDefault() *assets.Wei {
	return g.c.PriceDefault
}

func (g *gasEstimatorConfig) PriceMin() *assets.Wei {
	return g.c.PriceMin
}

func (g *gasEstimatorConfig) PriceMax() *assets.Wei {
	return g.c.PriceMax
}

func (g *gasEstimatorConfig) TipCapDefault() *assets.Wei {
	return g.c.TipCapDefault
}

func (g *gasEstimatorConfig) TipCapMin() *assets.Wei {
	return g.c.TipCapMin
}

func (g *gasEstimatorConfig) Mode() string {
	return *g.c.Mode
}

func (g *gasEstimatorConfig) LimitJobType() config.LimitJobType {
	return &limitJobTypeConfig{c: g.c.LimitJobType}
}

type limitJobTypeConfig struct {
	c GasLimitJobType
}

func (l *limitJobTypeConfig) OCR() *uint32 {
	return l.c.OCR
}

func (l *limitJobTypeConfig) OCR2() *uint32 {
	return l.c.OCR2
}

func (l *limitJobTypeConfig) DR() *uint32 {
	return l.c.DR
}

func (l *limitJobTypeConfig) FM() *uint32 {
	return l.c.FM
}

func (l *limitJobTypeConfig) Keeper() *uint32 {
	return l.c.Keeper
}

func (l *limitJobTypeConfig) VRF() *uint32 {
	return l.c.VRF
}

type blockHistoryConfig struct {
	c             BlockHistoryEstimator
	blockDelay    *uint16
	bumpThreshold *uint32
}

func (b *blockHistoryConfig) BatchSize() uint32 {
	return *b.c.BatchSize
}

func (b *blockHistoryConfig) BlockHistorySize() uint16 {
	return *b.c.BlockHistorySize
}

func (b *blockHistoryConfig) CheckInclusionBlocks() uint16 {
	return *b.c.CheckInclusionBlocks
}

func (b *blockHistoryConfig) CheckInclusionPercentile() uint16 {
	return *b.c.CheckInclusionPercentile
}

func (b *blockHistoryConfig) EIP1559FeeCapBufferBlocks() uint16 {
	if b.c.EIP1559FeeCapBufferBlocks == nil {
		return uint16(*b.bumpThreshold) + 1
	}
	return *b.c.EIP1559FeeCapBufferBlocks
}

func (b *blockHistoryConfig) TransactionPercentile() uint16 {
	return *b.c.TransactionPercentile
}

func (b *blockHistoryConfig) BlockDelay() uint16 {
	return *b.blockDelay
}
