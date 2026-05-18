package opstack

import (
	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
)

var (
	// Optimism Contracts
	// See https://docs.optimism.io/chain/addresses
	OptimismContractsByChainID map[uint64]map[string]common.Address
)

func init() {
	OptimismContractsByChainID = map[uint64]map[string]common.Address{
		chainsel.ETHEREUM_TESTNET_SEPOLIA.EvmChainID: {
			"L1StandardBridge":       common.HexToAddress("0xFBb0621E0B23b5478B630BD55a5f21f67730B0F1"),
			"L1CrossDomainMessenger": common.HexToAddress("0x58Cc85b8D04EA49cC6DBd3CbFFd00B4B8D6cb3ef"),
			"WETH":                   common.HexToAddress("0x7b79995e5f793a07bc00c21412e50ecae098e7f9"),
			"FaucetTestingToken":     common.HexToAddress("0x5589BB8228C07c4e15558875fAf2B859f678d129"),
			"OptimismPortalProxy":    common.HexToAddress("0x16Fc5058F25648194471939df75CF27A2fdC48BC"),
			"L2OutputOracle":         common.HexToAddress("0x90E9c4f8a994a250F6aEfd61CAFb4F2e895D458F"), // Removed after FPAC upgrade
		},
		chainsel.ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1.EvmChainID: {
			"WETH":                common.HexToAddress("0x4200000000000000000000000000000000000006"),
			"FaucetTestingToken":  common.HexToAddress("0xD08a2917653d4E460893203471f0000826fb4034"),
			"L2ToL1MessagePasser": common.HexToAddress("0x4200000000000000000000000000000000000016"),
		},
	}
}
