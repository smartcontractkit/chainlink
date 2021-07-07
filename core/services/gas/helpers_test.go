package gas

import "math/big"

func BlockHistoryEstimatorFromInterface(bhe Estimator) *BlockHistoryEstimator {
	return bhe.(*BlockHistoryEstimator)
}

func SetRollingBlockHistory(bhe Estimator, blocks []Block) {
	bhe.(*BlockHistoryEstimator).rollingBlockHistory = blocks
}

func GetGasPrice(b *BlockHistoryEstimator) *big.Int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.gasPrice
}

func GetTipCap(b *BlockHistoryEstimator) *big.Int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.tipCap
}
