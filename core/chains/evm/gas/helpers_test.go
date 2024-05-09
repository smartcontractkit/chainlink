package gas

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/common/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
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

func SetRollingBlockHistory(bhe *BlockHistoryEstimator, blocks []evmtypes.Block) {
	bhe.blocksMu.Lock()
	defer bhe.blocksMu.Unlock()
	bhe.blocks = blocks
}

func GetRollingBlockHistory(bhe *BlockHistoryEstimator) []evmtypes.Block {
	return bhe.getBlocks()
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

func BlockHistoryEstimatorV2FromInterface(bhe EvmEstimator) *BlockHistoryEstimatorV2 {
	return bhe.(*BlockHistoryEstimatorV2)
}

func (b *BlockHistoryEstimatorV2) CheckConnectivityV2(attempts []EvmPriorAttempt) error {
	return b.checkConnectivity(attempts)
}

func GetRollingBlockHistoryV2(bhe *BlockHistoryEstimatorV2) []evmtypes.Block {
	return bhe.getBlocks()
}

func SetRollingBlockHistoryV2(bhe *BlockHistoryEstimatorV2, blocks []evmtypes.Block) {
	bhe.blocksMu.Lock()
	defer bhe.blocksMu.Unlock()
	bhe.blocks = blocks
}

func SetGasCost(b *BlockHistoryEstimatorV2, gc *assets.Wei) {
	b.costMu.Lock()
	defer b.costMu.Unlock()
	b.gasCost = gc
}

func GetGasCost(b *BlockHistoryEstimatorV2) *assets.Wei {
	b.costMu.RLock()
	defer b.costMu.RUnlock()
	return b.gasCost
}

func GetLatestBaseFeeV2(b *BlockHistoryEstimatorV2) *assets.Wei {
	b.latestMu.RLock()
	defer b.latestMu.RUnlock()
	if b.latest == nil {
		return nil
	}
	return b.latest.BaseFeePerGas
}

func SimulateStartV2(t *testing.T, b *BlockHistoryEstimatorV2) {
	require.NoError(t, b.StartOnce("BlockHistoryEstimatorV2SimulatedStart", func() error { return nil }))
}

type MockBlockHistoryConfig struct {
	BatchSizeF                 uint32
	BlockDelayF                uint16
	BlockHistorySizeF          uint16
	CheckInclusionBlocksF      uint16
	CheckInclusionPercentileF  uint16
	EIP1559FeeCapBufferBlocksF uint16
	TransactionPercentileF     uint16
	FinalityTagEnabledF        bool
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
	ChainTypeF          string
	FinalityTagEnabledF bool
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

func (m *MockConfig) FinalityTagEnabled() bool {
	return m.FinalityTagEnabledF
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
	LimitMaxF           uint64
	ModeF               string
}

func NewMockGasConfig() *MockGasEstimatorConfig {
	return &MockGasEstimatorConfig{}
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

func (m *MockGasEstimatorConfig) LimitMax() uint64 {
	return m.LimitMaxF
}

func (m *MockGasEstimatorConfig) Mode() string {
	return m.ModeF
}
