package gas

import (
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
)

// chainSpecificIsUsable allows for additional logic specific to a particular
// Config that determines whether a transaction should be used for gas estimation
func chainSpecificIsUsable(tx evmtypes.Transaction, cfg Config) bool {
	if cfg.ChainType() == config.ChainXDai {
		// GasPrice 0 on most chains is great since it indicates cheap/free transactions.
		// However, xDai reserves a special type of "bridge" transaction with 0 gas
		// price that is always processed at top priority. Ordinary transactions
		// must be priced at least 1GWei, so we have to discard anything priced
		// below that (unless the contract is whitelisted).
		if tx.GasPrice != nil && tx.GasPrice.Cmp(cfg.EvmMinGasPriceWei()) < 0 {
			return false
		}
	}
	if cfg.ChainType() == config.ChainOptimismBedrock {
		// This is a special deposit transaction type introduced in Bedrock upgrade.
		// This is a system transaction that it will occur at least one time per block.
		// We should discard this type before even processing it to avoid flooding the
		// logs with warnings.
		// https://github.com/ethereum-optimism/optimism/blob/develop/specs/deposits.md
		if tx.Type == 0x7e {
			return false
		}
	}
	return true
}
