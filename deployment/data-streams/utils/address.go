package utils

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

func EnvironmentAddresses(e deployment.Environment) (addresses map[string]deployment.TypeAndVersion, err error) {
	addresses = make(map[string]deployment.TypeAndVersion)
	records, err := e.DataStore.Addresses().Fetch()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch addresses from datastore: %w", err)
	}
	for _, record := range records {
		addresses[record.Address] = deployment.TypeAndVersion{
			Type:    deployment.ContractType(record.Type),
			Version: *record.Version,
		}
	}
	return addresses, nil
}

// GetContractAddress returns the address for a specific contract type. Used when expecting only one contract
func GetContractAddress(addresses datastore.AddressRefStore, contractType deployment.ContractType) (string, error) {
	records := addresses.Filter(datastore.AddressRefByType(datastore.ContractType(contractType)))
	if len(records) != 1 {
		return "", fmt.Errorf("expected 1 %s address, found %d", contractType, len(records))
	}
	return records[0].Address, nil
}
