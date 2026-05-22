package deployment

import (
	"slices"

	"github.com/ethereum/go-ethereum/common"
)

func IsAddressListUnique(addresses []common.Address) bool {
	addressSet := make(map[common.Address]struct{})
	for _, address := range addresses {
		if _, exists := addressSet[address]; exists {
			return false
		}
		addressSet[address] = struct{}{}
	}
	return true
}

func AddressListContainsEmptyAddress(addresses []common.Address) bool {
	return slices.Contains(addresses, (common.Address{}))
}
