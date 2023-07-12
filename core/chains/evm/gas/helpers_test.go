package gas

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

func init() {
	// No need to wait 10 seconds in tests
	MaxStartTime = 1 * time.Second
}

func (b *BlockHistoryEstimator) CheckConnectivity(attempts []EvmPriorAttempt) error {
	return b.checkConnectivity(attempts)
}

func BlockHistoryEstimatorFromInterface(bhe EvmEstimator) *BlockHistoryEstimator {
	return bhe.(*BlockHistoryEstimator)
}

func SetRollingBlockHistory(bhe EvmEstimator, blocks []evmtypes.Block) {
	bhe.(*BlockHistoryEstimator).blocksMu.Lock()
	defer bhe.(*BlockHistoryEstimator).blocksMu.Unlock()
	bhe.(*BlockHistoryEstimator).blocks = blocks
}

func GetRollingBlockHistory(bhe EvmEstimator) []evmtypes.Block {
	return bhe.(*BlockHistoryEstimator).getBlocks()
}

func SetGasPrice(b *BlockHistoryEstimator, gp *assets.Wei) {
	b.priceMu.Lock()
	defer b.priceMu.Unlock()
	b.gasPrice = gp
}

func SetTipCap(b *BlockHistoryEstimator, gp *assets.Wei) {
	b.priceMu.Lock()
	defer b.priceMu.Unlock()
	b.tipCap = gp
}

func GetGasPrice(b *BlockHistoryEstimator) *assets.Wei {
	b.priceMu.RLock()
	defer b.priceMu.RUnlock()
	return b.gasPrice
}

func GetTipCap(b *BlockHistoryEstimator) *assets.Wei {
	b.priceMu.RLock()
	defer b.priceMu.RUnlock()
	return b.tipCap
}

func GetLatestBaseFee(b *BlockHistoryEstimator) *assets.Wei {
	b.latestMu.RLock()
	defer b.latestMu.RUnlock()
	if b.latest == nil {
		return nil
	}
	return b.latest.BaseFeePerGas
}

func SimulateStart(t *testing.T, b *BlockHistoryEstimator) {
	require.NoError(t, b.StartOnce("BlockHistoryEstimatorSimulatedStart", func() error { return nil }))
}

type MockBlockHistoryConfig struct {
	BatchSizeF                 uint32
	BlockDelayF                uint16
	BlockHistorySizeF          uint16
	CheckInclusionBlocksF      uint16
	CheckInclusionPercentileF  uint16
	EIP1559FeeCapBufferBlocksF uint16
	TransactionPercentileF     uint16
}

func (m *MockBlockHistoryConfig) BatchSize() uint32 {
	return m.BatchSizeF
}

func (m *MockBlockHistoryConfig) BlockDelay() uint16 {
	return m.BlockDelayF
}

func (m *MockBlockHistoryConfig) BlockHistorySize() uint16 {
	return m.BlockHistorySizeF
}

func (m *MockBlockHistoryConfig) CheckInclusionPercentile() uint16 {
	return m.CheckInclusionPercentileF
}

func (m *MockBlockHistoryConfig) CheckInclusionBlocks() uint16 {
	return m.CheckInclusionBlocksF
}

func (m *MockBlockHistoryConfig) EIP1559FeeCapBufferBlocks() uint16 {
	return m.EIP1559FeeCapBufferBlocksF
}

func (m *MockBlockHistoryConfig) TransactionPercentile() uint16 {
	return m.TransactionPercentileF
}

type MockConfig struct {
	ChainTypeF string
}

func NewMockConfig() *MockConfig {
	return &MockConfig{}
}

func (m *MockConfig) ChainType() config.ChainType {
	return config.ChainType(m.ChainTypeF)
}

func (m *MockConfig) FinalityDepth() uint32 {
	panic("not implemented") // TODO: Implement
}

type MockGasEstimatorConfig struct {
	EIP1559DynamicFeesF bool
	BumpPercentF        uint16
	BumpThresholdF      uint64
	BumpMinF            *assets.Wei
	LimitMultiplierF    float32
	TipCapDefaultF      *assets.Wei
	TipCapMinF          *assets.Wei
	PriceMaxF           *assets.Wei
	PriceMinF           *assets.Wei
	PriceDefaultF       *assets.Wei
	FeeCapDefaultF      *assets.Wei
	LimitMaxF           uint32
	ModeF               string
}

func (m *MockGasEstimatorConfig) BumpPercent() uint16 {
	return m.BumpPercentF
}

func (m *MockGasEstimatorConfig) BumpThreshold() uint64 {
	return m.BumpThresholdF
}

func (m *MockGasEstimatorConfig) BumpMin() *assets.Wei {
	return m.BumpMinF
}

func (m *MockGasEstimatorConfig) EIP1559DynamicFees() bool {
	return m.EIP1559DynamicFeesF
}

func (m *MockGasEstimatorConfig) LimitMultiplier() float32 {
	return m.LimitMultiplierF
}

func (m *MockGasEstimatorConfig) PriceDefault() *assets.Wei {
	return m.PriceDefaultF
}

func (m *MockGasEstimatorConfig) TipCapDefault() *assets.Wei {
	return m.TipCapDefaultF
}

func (m *MockGasEstimatorConfig) TipCapMin() *assets.Wei {
	return m.TipCapMinF
}

func (m *MockGasEstimatorConfig) PriceMax() *assets.Wei {
	return m.PriceMaxF
}

func (m *MockGasEstimatorConfig) PriceMin() *assets.Wei {
	return m.PriceMinF
}

func (m *MockGasEstimatorConfig) FeeCapDefault() *assets.Wei {
	return m.FeeCapDefaultF
}

func (m *MockGasEstimatorConfig) LimitMax() uint32 {
	return m.LimitMaxF
}

func (m *MockGasEstimatorConfig) Mode() string {
	return m.ModeF
}
