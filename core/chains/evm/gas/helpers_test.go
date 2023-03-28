package gas

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

func init() {
	// No need to wait 10 seconds in tests
	MaxStartTime = 1 * time.Second
}

func (b *BlockHistoryEstimator) CheckConnectivity(attempts []txmgrtypes.PriorAttempt[EvmFee, common.Hash]) error {
	return b.checkConnectivity(MakeEvmPriorAttempts(attempts))
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

type MockConfig struct {
	BlockHistoryEstimatorBatchSizeF                 uint32
	BlockHistoryEstimatorBlockDelayF                uint16
	BlockHistoryEstimatorBlockHistorySizeF          uint16
	BlockHistoryEstimatorCheckInclusionBlocksF      uint16
	BlockHistoryEstimatorCheckInclusionPercentileF  uint16
	BlockHistoryEstimatorEIP1559FeeCapBufferBlocksF uint16
	BlockHistoryEstimatorTransactionPercentileF     uint16
	ChainTypeF                                      string
	EvmEIP1559DynamicFeesF                          bool
	EvmGasBumpPercentF                              uint16
	EvmGasBumpThresholdF                            uint64
	EvmGasBumpWeiF                                  *assets.Wei
	EvmGasLimitMultiplierF                          float32
	EvmGasTipCapDefaultF                            *assets.Wei
	EvmGasTipCapMinimumF                            *assets.Wei
	EvmMaxGasPriceWeiF                              *assets.Wei
	EvmMinGasPriceWeiF                              *assets.Wei
	EvmGasPriceDefaultF                             *assets.Wei
}

func NewMockConfig() *MockConfig {
	return &MockConfig{}
}

func (m *MockConfig) BlockHistoryEstimatorBatchSize() uint32 {
	return m.BlockHistoryEstimatorBatchSizeF
}

func (m *MockConfig) BlockHistoryEstimatorBlockDelay() uint16 {
	return m.BlockHistoryEstimatorBlockDelayF
}

func (m *MockConfig) BlockHistoryEstimatorBlockHistorySize() uint16 {
	return m.BlockHistoryEstimatorBlockHistorySizeF
}

func (m *MockConfig) BlockHistoryEstimatorCheckInclusionPercentile() uint16 {
	return m.BlockHistoryEstimatorCheckInclusionPercentileF
}

func (m *MockConfig) BlockHistoryEstimatorCheckInclusionBlocks() uint16 {
	return m.BlockHistoryEstimatorCheckInclusionBlocksF
}

func (m *MockConfig) BlockHistoryEstimatorEIP1559FeeCapBufferBlocks() uint16 {
	return m.BlockHistoryEstimatorEIP1559FeeCapBufferBlocksF
}

func (m *MockConfig) BlockHistoryEstimatorTransactionPercentile() uint16 {
	return m.BlockHistoryEstimatorTransactionPercentileF
}

func (m *MockConfig) ChainType() config.ChainType {
	return config.ChainType(m.ChainTypeF)
}

func (m *MockConfig) EvmEIP1559DynamicFees() bool {
	return m.EvmEIP1559DynamicFeesF
}

func (m *MockConfig) EvmFinalityDepth() uint32 {
	panic("not implemented") // TODO: Implement
}

func (m *MockConfig) EvmGasBumpPercent() uint16 {
	return m.EvmGasBumpPercentF
}

func (m *MockConfig) EvmGasBumpThreshold() uint64 {
	return m.EvmGasBumpThresholdF
}

func (m *MockConfig) EvmGasBumpWei() *assets.Wei {
	return m.EvmGasBumpWeiF
}

func (m *MockConfig) EvmGasFeeCapDefault() *assets.Wei {
	panic("not implemented") // TODO: Implement
}

func (m *MockConfig) EvmGasLimitMax() uint32 {
	panic("not implemented") // TODO: Implement
}

func (m *MockConfig) EvmGasLimitMultiplier() float32 {
	return m.EvmGasLimitMultiplierF
}

func (m *MockConfig) EvmGasPriceDefault() *assets.Wei {
	return m.EvmGasPriceDefaultF
}

func (m *MockConfig) EvmGasTipCapDefault() *assets.Wei {
	return m.EvmGasTipCapDefaultF
}

func (m *MockConfig) EvmGasTipCapMinimum() *assets.Wei {
	return m.EvmGasTipCapMinimumF
}

func (m *MockConfig) EvmMaxGasPriceWei() *assets.Wei {
	return m.EvmMaxGasPriceWeiF
}

func (m *MockConfig) EvmMinGasPriceWei() *assets.Wei {
	return m.EvmMinGasPriceWeiF
}

func (m *MockConfig) GasEstimatorMode() string {
	panic("not implemented") // TODO: Implement
}
