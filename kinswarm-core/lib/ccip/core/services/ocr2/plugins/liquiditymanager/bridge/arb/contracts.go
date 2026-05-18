package arb

import (
	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
)

var (
	// Arbitrum Contracts
	// See https://docs.arbitrum.io/for-devs/useful-addresses
	AllContracts map[uint64]map[string]common.Address
)

func init() {
	AllContracts = map[uint64]map[string]common.Address{
		chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector: {
			"L1GatewayRouter": common.HexToAddress("0xcE18836b233C83325Cc8848CA4487e94C6288264"),
			"L1Outbox":        common.HexToAddress("0x65f07C7D521164a4d5DaC6eB8Fac8DA067A3B78F"),
			// labeled "Delayed Inbox" in the arbitrum docs
			"L1Inbox": common.HexToAddress("0xaAe29B0366299461418F5324a79Afc425BE5ae21"),
			"Rollup":  common.HexToAddress("0xd80810638dbDF9081b72C1B33c65375e807281C8"),
			"WETH":    common.HexToAddress("0x7b79995e5f793A07Bc00c21412e50Ecae098E7f9"),
		},
		chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector: {
			"L2GatewayRouter": common.HexToAddress("0x9fDD1C4E4AA24EEc1d913FABea925594a20d43C7"),
			"NodeInterface":   common.HexToAddress("0x00000000000000000000000000000000000000C8"),
			"WETH":            common.HexToAddress("0x980B62Da83eFf3D4576C647993b0c1D7faf17c73"),
		},
	}
}
