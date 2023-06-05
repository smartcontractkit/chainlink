package gas

import (
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

// chainSpecificIsUsable allows for additional logic specific to a particular
// Config that determines whether a transaction should be used for gas estimation
func chainSpecificIsUsable(tx evmtypes.Transaction, block evmtypes.Block, cfg Config) bool {
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
	if cfg.ChainType() == config.ChainCelo {
		// Celo specific transaction type that utilizes the feeCurrency field.
		if tx.Type == 0x7c {
			return false
		}
		// Celo has not yet fully migrated to the 0x7c type for special feeCurrency transactions
		// and uses the standard 0x0, 0x2 types instead. We need to discard any invalid transactions
		// and not throw an error since this can happen from time to time and it's an expected behavior
		// until they fully migrate to 0x7c.
		if tx.GasPrice.Cmp(block.BaseFeePerGas) < 0 {
			return false

		}
	}
	return true
}
