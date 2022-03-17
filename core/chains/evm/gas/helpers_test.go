package gas

import (
	"math/big"
	"time"
)

func init() {
	// No need to wait 10 seconds in tests
	MaxStartTime = 1 * time.Second
}

func BlockHistoryEstimatorFromInterface(bhe Estimator) *BlockHistoryEstimator {
	return bhe.(*BlockHistoryEstimator)
}

func SetRollingBlockHistory(bhe Estimator, blocks []Block) {
	bhe.(*BlockHistoryEstimator).rollingBlockHistory = blocks
}

func SetGasPrice(b *BlockHistoryEstimator, gp *big.Int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.gasPrice = gp
}

func SetTipCap(b *BlockHistoryEstimator, gp *big.Int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tipCap = gp
}

func GetGasPrice(b *BlockHistoryEstimator) *big.Int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.gasPrice
}

func GetTipCap(b *BlockHistoryEstimator) *big.Int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.tipCap
}

func GetLatestBaseFee(b *BlockHistoryEstimator) *big.Int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.latestBaseFee
}

func SimulateStart(b *BlockHistoryEstimator) {
	b.StartOnce("BlockHistoryEstimatorSimulatedStart", func() error { return nil })
}
