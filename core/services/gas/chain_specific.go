package gas

import (
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
)

// chainSpecificIsUsableTx allows for additional logic specific to a particular
// chain that determines whether a transaction should be used for gas estimation
func chainSpecificIsUsableTx(tx Transaction, cfg Config) bool {
	if cfg.ChainType() == evmtypes.XDai {
		// GasPrice 0 on most chains is great since it indicates cheap/free transactions.
		// However, xDai reserves a special type of "bridge" transaction with 0 gas
		// price that is always processed at top priority. Ordinary transactions
		// must be priced at least 1GWei, so we have to discard anything priced
		// below that (unless the contract is whitelisted).
		if tx.GasPrice != nil && tx.GasPrice.Cmp(cfg.EvmMinGasPriceWei()) < 0 {
			return false
		}
	}
	return true
}
