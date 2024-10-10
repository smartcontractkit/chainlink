package config

import (
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
)

type gasEstimatorConfig struct {
	c                       toml.GasEstimator
	k                       toml.KeySpecificConfig
	blockDelay              *uint16
	transactionsMaxInFlight *uint32
}

func (g *gasEstimatorConfig) PriceMaxKey(addr gethcommon.Address) *assets.Wei {
	var keySpecific *assets.Wei
	for i := range g.k {
		ks := g.k[i]
		if ks.Key.Address() == addr {
			keySpecific = ks.GasEstimator.PriceMax
			break
		}
	}

	chainSpecific := g.c.PriceMax
	if keySpecific != nil && keySpecific.Cmp(chainSpecific) < 0 {
		return keySpecific
	}

	return g.c.PriceMax
}

func (g *gasEstimatorConfig) BlockHistory() BlockHistory {
	return &blockHistoryConfig{c: g.c.BlockHistory, blockDelay: g.blockDelay, bumpThreshold: g.c.BumpThreshold}
}

func (g *gasEstimatorConfig) FeeHistory() FeeHistory {
	return &feeHistoryConfig{c: g.c.FeeHistory}
}

func (g *gasEstimatorConfig) DAOracle() DAOracle {
	return &daOracleConfig{c: g.c.DAOracle}
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

func (g *gasEstimatorConfig) LimitDefault() uint64 {
	return *g.c.LimitDefault
}

func (g *gasEstimatorConfig) LimitMax() uint64 {
	return *g.c.LimitMax
}

func (g *gasEstimatorConfig) LimitMultiplier() float32 {
	f, _ := g.c.LimitMultiplier.BigFloat().Float32()
	return f
}

func (g *gasEstimatorConfig) LimitTransfer() uint64 {
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

func (g *gasEstimatorConfig) LimitJobType() LimitJobType {
	return &limitJobTypeConfig{c: g.c.LimitJobType}
}

func (g *gasEstimatorConfig) EstimateLimit() bool {
	return *g.c.EstimateLimit
}

// ====== DAOracle ======
type daOracleConfig struct {
	c toml.DAOracle
}

func (d *daOracleConfig) OracleType() toml.OracleType {
	return d.c.OracleType
}

// OracleAddress returns the address of the oracle contract and is only supported on the OP stack for now.
func (d *daOracleConfig) OracleAddress() *types.EIP55Address {
	if d.c.OracleType != toml.OP {
		return nil
	}
	return d.c.OracleAddress
}

// CustomGasPriceCalldata returns the calldata for a custom gas price API and is only supported on the OP stack
// for now.
func (d *daOracleConfig) CustomGasPriceCalldata() string {
	if d.c.OracleType != toml.OP {
		return ""
	}
	return d.c.CustomGasPriceCalldata
}

// ====== LimitJobType ======
type limitJobTypeConfig struct {
	c toml.GasLimitJobType
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

// ====== BlockHistory ======
type blockHistoryConfig struct {
	c             toml.BlockHistoryEstimator
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

// ====== FeeHistory ======
type feeHistoryConfig struct {
	c toml.FeeHistoryEstimator
}

func (u *feeHistoryConfig) CacheTimeout() time.Duration {
	return u.c.CacheTimeout.Duration()
}
