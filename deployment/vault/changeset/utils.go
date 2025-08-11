package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

func GetContractAddress(ds interface{}, chainSelector uint64, contractType cldf.ContractType) (string, error) {
	if ds == nil {
		return "", errors.New("datastore is nil")
	}

	var addresses []datastore.AddressRef

	switch v := ds.(type) {
	case datastore.DataStore:
		addresses = v.Addresses().Filter(
			datastore.AddressRefByChainSelector(chainSelector),
			datastore.AddressRefByType(datastore.ContractType(contractType)),
		)
	case datastore.MutableDataStore:
		addresses = v.Addresses().Filter(
			datastore.AddressRefByChainSelector(chainSelector),
			datastore.AddressRefByType(datastore.ContractType(contractType)),
		)
	default:
		return "", fmt.Errorf("unsupported datastore type: %T", ds)
	}

	// Return the first match since we expect only one contract of each type per chain (for now)
	if len(addresses) > 0 {
		return addresses[0].Address, nil
	}

	return "", fmt.Errorf("contract of type %s not found for chain %d", contractType, chainSelector)
}
