package configtest

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store/config"
	null "gopkg.in/guregu/null.v4"
)

var _ config.EVMConfig = &TestEVMConfig{}

var (
	MinimumContractPayment = assets.NewLink(100)
)

const (
	HeadSamplingIntervalInTest = 0 * time.Millisecond // Head sampling disabled
)

type EVMConfigOverrides struct {
	EvmLogBackfillBatchSize null.Int

	BlockHistoryEstimatorBlockDelay       null.Int
	BlockHistoryEstimatorBlockHistorySize null.Int
	EvmFinalityDepth                      null.Int
	EvmMaxGasPriceWei                     *big.Int
	EvmGasBumpPercent                     null.Int
	EvmGasBumpTxDepth                     null.Int

	EvmGasLimitDefault null.Int

	EvmHeadTrackerHistoryDepth       null.Int
	EvmGasBumpWei                    *big.Int
	EvmGasLimitMultiplier            null.Float
	EvmGasPriceDefault               *big.Int
	EvmHeadTrackerSamplingInterval   *time.Duration
	EvmHeadTrackerMaxBufferSize      null.Int
	EthTxResendAfterThreshold        *time.Duration
	EvmNonceAutoSync                 null.Bool
	EvmRPCDefaultBatchSize           null.Int
	FlagsContractAddress             null.String
	GasEstimatorMode                 null.String
	MinRequiredOutgoingConfirmations null.Int
}

// TestEVMConfig defaults to whatever config.NewEVMConfig()
// gives but allows overriding certain methods
type TestEVMConfig struct {
	config.EVMConfig
	Overrides     EVMConfigOverrides
	GeneralConfig *TestGeneralConfig
	t             testing.TB
}

func NewTestEVMConfig(t testing.TB, generalcfg *TestGeneralConfig) *TestEVMConfig {
	evmcfg := config.NewEVMConfig(generalcfg)
	return &TestEVMConfig{
		evmcfg,
		EVMConfigOverrides{},
		generalcfg,
		t,
	}
}

func (c *TestEVMConfig) GasEstimatorMode() string {
	if c.Overrides.GasEstimatorMode.Valid {
		return c.Overrides.GasEstimatorMode.String
	}
	return "FixedPrice"
}

func (c *TestEVMConfig) EthTxResendAfterThreshold() time.Duration {
	if c.Overrides.EthTxResendAfterThreshold != nil {
		return *c.Overrides.EthTxResendAfterThreshold
	}
	return 0
}

func (c *TestEVMConfig) EvmFinalityDepth() uint {
	if c.Overrides.EvmFinalityDepth.Valid {
		return uint(c.Overrides.EvmFinalityDepth.Int64)
	}
	return 15
}

func (c *TestEVMConfig) EthTxReaperThreshold() time.Duration {
	return 0
}

func (c *TestEVMConfig) EthHeadTrackerSamplingInterval() time.Duration {
	return HeadSamplingIntervalInTest
}

func (c *TestEVMConfig) EvmGasBumpThreshold() uint64 {
	return 3
}

func (c *TestEVMConfig) MinIncomingConfirmations() uint32 {
	return 1
}

func (c *TestEVMConfig) MinRequiredOutgoingConfirmations() uint64 {
	if c.Overrides.MinRequiredOutgoingConfirmations.Valid {
		return uint64(c.Overrides.MinRequiredOutgoingConfirmations.Int64)
	}
	return 1
}

func (c *TestEVMConfig) MinimumContractPayment() *assets.Link {
	return MinimumContractPayment
}

func (c *TestEVMConfig) BalanceMonitorEnabled() bool {
	return false
}

func (c *TestEVMConfig) EvmHeadTrackerMaxBufferSize() uint {
	if c.Overrides.EvmHeadTrackerMaxBufferSize.Valid {
		return uint(c.Overrides.EvmHeadTrackerMaxBufferSize.Int64)
	}
	return c.EVMConfig.EvmHeadTrackerMaxBufferSize()
}

func (c *TestEVMConfig) EvmGasPriceDefault() *big.Int {
	if c.Overrides.EvmGasPriceDefault != nil {
		return c.Overrides.EvmGasPriceDefault
	}
	return c.EVMConfig.EvmGasPriceDefault()
}

func (c *TestEVMConfig) SetEvmGasPriceDefault(p *big.Int) error {
	c.Overrides.EvmGasPriceDefault = p
	return nil
}

func (c *TestEVMConfig) BlockHistoryEstimatorBlockDelay() uint16 {
	if c.Overrides.BlockHistoryEstimatorBlockDelay.Valid {
		return uint16(c.Overrides.BlockHistoryEstimatorBlockDelay.Int64)
	}
	return c.EVMConfig.BlockHistoryEstimatorBlockDelay()
}

func (c *TestEVMConfig) BlockHistoryEstimatorBlockHistorySize() uint16 {
	if c.Overrides.BlockHistoryEstimatorBlockHistorySize.Valid {
		return uint16(c.Overrides.BlockHistoryEstimatorBlockHistorySize.Int64)
	}
	return c.EVMConfig.BlockHistoryEstimatorBlockHistorySize()
}

func (c *TestEVMConfig) EvmGasLimitMultiplier() float32 {
	if c.Overrides.EvmGasLimitMultiplier.Valid {
		return float32(c.Overrides.EvmGasLimitMultiplier.Float64)
	}
	return c.EVMConfig.EvmGasLimitMultiplier()
}

func (c *TestEVMConfig) EvmNonceAutoSync() bool {
	if c.Overrides.EvmNonceAutoSync.Valid {
		return c.Overrides.EvmNonceAutoSync.Bool
	}
	return c.EVMConfig.EvmNonceAutoSync()
}

func (c *TestEVMConfig) EvmGasBumpWei() *big.Int {
	if c.Overrides.EvmGasBumpWei != nil {
		return c.Overrides.EvmGasBumpWei
	}
	return c.EVMConfig.EvmGasBumpWei()
}

func (c *TestEVMConfig) EvmGasBumpPercent() uint16 {
	if c.Overrides.EvmGasBumpPercent.Valid {
		return uint16(c.Overrides.EvmGasBumpPercent.Int64)
	}
	return c.EVMConfig.EvmGasBumpPercent()
}

func (c *TestEVMConfig) EvmRPCDefaultBatchSize() uint32 {
	if c.Overrides.EvmRPCDefaultBatchSize.Valid {
		return uint32(c.Overrides.EvmRPCDefaultBatchSize.Int64)
	}
	return c.EVMConfig.EvmRPCDefaultBatchSize()
}

func (c *TestEVMConfig) EvmMaxGasPriceWei() *big.Int {
	if c.Overrides.EvmMaxGasPriceWei != nil {
		return c.Overrides.EvmMaxGasPriceWei
	}
	return c.EVMConfig.EvmMaxGasPriceWei()
}

func (c *TestEVMConfig) EvmGasBumpTxDepth() uint16 {
	if c.Overrides.EvmGasBumpTxDepth.Valid {
		return uint16(c.Overrides.EvmGasBumpTxDepth.Int64)
	}
	return c.EVMConfig.EvmGasBumpTxDepth()
}

func (c *TestEVMConfig) FlagsContractAddress() string {
	if c.Overrides.FlagsContractAddress.Valid {
		return c.Overrides.FlagsContractAddress.String
	}
	return c.EVMConfig.FlagsContractAddress()
}

func (c *TestEVMConfig) EvmHeadTrackerHistoryDepth() uint {
	if c.Overrides.EvmHeadTrackerHistoryDepth.Valid {
		return uint(c.Overrides.EvmHeadTrackerHistoryDepth.Int64)
	}
	return c.EVMConfig.EvmHeadTrackerHistoryDepth()
}

func (c *TestEVMConfig) EvmHeadTrackerSamplingInterval() time.Duration {
	if c.Overrides.EvmHeadTrackerSamplingInterval != nil {
		return *c.Overrides.EvmHeadTrackerSamplingInterval
	}
	return c.EVMConfig.EvmHeadTrackerSamplingInterval()
}

func (c *TestEVMConfig) EvmLogBackfillBatchSize() uint32 {
	if c.Overrides.EvmLogBackfillBatchSize.Valid {
		return uint32(c.Overrides.EvmLogBackfillBatchSize.Int64)
	}
	return c.EVMConfig.EvmLogBackfillBatchSize()
}

func (c *TestEVMConfig) EvmGasLimitDefault() uint64 {
	if c.Overrides.EvmGasLimitDefault.Valid {
		return uint64(c.Overrides.EvmGasLimitDefault.Int64)
	}
	return c.EVMConfig.EvmGasLimitDefault()
}
