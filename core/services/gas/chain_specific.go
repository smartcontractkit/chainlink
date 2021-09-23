package gas

import "math/big"

// chainSpecificIsUsableTx allows for additional logic specific to a
// particular chain that determines whether a transction should be used for gas
// estimation
func chainSpecificIsUsableTx(tx Transaction, minGasPriceWei, chainID *big.Int) bool {
	if isXDai(chainID) {
		// GasPrice 0 on most chains is great since it indicates cheap/free transctions.
		// However, xDai reserves a special type of "bridge" transaction with 0 gas
		// price that is always processed at top priority. Ordinary transactions
		// must be priced at least 1GWei, so we have to discard anything priced
		// below that (unless the contract is whitelisted).
		if tx.GasPrice != nil && tx.GasPrice.Cmp(minGasPriceWei) < 0 {
			return false
		}
	}
	return true
}

func isXDai(chainID *big.Int) bool {
	return chainID.Cmp(big.NewInt(100)) == 0
}
