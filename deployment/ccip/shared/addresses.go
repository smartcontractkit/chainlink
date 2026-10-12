package shared

import (
	"strings"

	chainsel "github.com/smartcontractkit/chain-selectors"
)

// NormalizeAddress lower-cases an address only on families where case carries no meaning.
//
// EVM addresses are hex and are commonly written in EIP-55 checksum form, so the same contract
// can appear in two casings and must compare equal. Solana and Sui addresses are case-sensitive,
// where folding case would map two different contracts onto one value. An unknown selector is
// left alone: preserving a distinction that does not exist is harmless, whereas erasing one that
// does silently merges two contracts.
func NormalizeAddress(chainSelector uint64, address string) string {
	family, err := chainsel.GetSelectorFamily(chainSelector)
	if err != nil || family != chainsel.FamilyEVM {
		return address
	}

	return strings.ToLower(address)
}

// AddressesEqual reports whether two addresses on the same chain refer to the same contract,
// using the case rules of that chain's family. Use it instead of == whenever an address that came
// from one source is compared against an address that came from another, since only one of them
// may be checksummed.
func AddressesEqual(chainSelector uint64, a, b string) bool {
	return NormalizeAddress(chainSelector, a) == NormalizeAddress(chainSelector, b)
}
