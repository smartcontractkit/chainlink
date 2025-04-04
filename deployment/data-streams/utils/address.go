package utils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment"
)

func MaybeFindEthAddress(ab deployment.AddressBook, chain uint64, typ deployment.ContractType) (common.Address, error) {
	addressHex, err := deployment.SearchAddressBook(ab, chain, typ)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to find contract %s address: %w", typ, err)
	}
	address := common.HexToAddress(addressHex)
	return address, nil
}
