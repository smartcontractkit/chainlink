package gas

import "math/big"

func BlockHistoryEstimatorFromInterface(bhe Estimator) *BlockHistoryEstimator {
	return bhe.(*BlockHistoryEstimator)
}

func SetRollingBlockHistory(bhe Estimator, blocks []Block) {
	bhe.(*BlockHistoryEstimator).rollingBlockHistory = blocks
}

func GetGasPrice(b *BlockHistoryEstimator) *big.Int {
	b.gasPriceMu.RLock()
	defer b.gasPriceMu.RUnlock()
	return b.gasPrice
}

func SetGasPrice(b *BlockHistoryEstimator, gp *big.Int) {
	b.gasPriceMu.Lock()
	defer b.gasPriceMu.Unlock()
	b.gasPrice = gp
}
