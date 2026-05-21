package evm

import (
	"slices"

	"github.com/ethereum/go-ethereum/common"
)

// isAddressListUnique checks if the list of addresses is unique
func isAddressListUnique(addresses []common.Address) bool {
	addressSet := make(map[common.Address]struct{})
	for _, address := range addresses {
		if _, exists := addressSet[address]; exists {
			return false
		}
		addressSet[address] = struct{}{}
	}
	return true
}

// addressListContainsEmptyAddress checks if the list of addresses contains an empty address
func addressListContainsEmptyAddress(addresses []common.Address) bool {
	return slices.Contains(addresses, (common.Address{}))
}
