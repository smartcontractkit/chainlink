package gas

import "math/big"

func BlockHistoryEstimatorFromInterface(bhe Estimator) *BlockHistoryEstimator {
	return bhe.(*BlockHistoryEstimator)
}

func SetRollingBlockHistory(bhe Estimator, blocks []Block) {
	bhe.(*BlockHistoryEstimator).rollingBlockHistory = blocks
}

<<<<<<< HEAD
func GetGasPrice(b *BlockHistoryEstimator) *big.Int {
	b.gasPriceMu.Lock()
	defer b.gasPriceMu.Unlock()
	return b.gasPrice
}
=======
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
>>>>>>> 588b3892f (Fix regression where gas estimator would not use current gas price)
