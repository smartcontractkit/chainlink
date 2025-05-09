package utils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

func MaybeFindEthAddress(ab cldf.AddressBook, chain uint64, typ cldf.ContractType) (common.Address, error) {
	addressHex, err := cldf.SearchAddressBook(ab, chain, typ)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to find contract %s address: %w", typ, err)
	}
	address := common.HexToAddress(addressHex)
	return address, nil
}
