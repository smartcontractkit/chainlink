package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func GetContractAddress(ds any, chainSelector uint64, contractType cldf.ContractType) (string, error) {
	return GetContractAddressWithQualifier(ds, chainSelector, contractType, "")
}

// GetContractAddressWithQualifier returns the contract address for the given chain and type, optionally filtered by qualifier.
// When qualifier is empty, it is treated as commonchangeset.DefaultTimelockQualifier ("default") so the first timelock per chain is used (matches deploy_timelock and migrated refs).
func GetContractAddressWithQualifier(ds any, chainSelector uint64, contractType cldf.ContractType, qualifier string) (string, error) {
	if ds == nil {
		return "", errors.New("datastore is nil")
	}
	if qualifier == "" {
		qualifier = commonchangeset.DefaultTimelockQualifier
	}

	var addresses []datastore.AddressRef

	switch v := ds.(type) {
	case datastore.DataStore:
		filters := []datastore.FilterFunc[datastore.AddressRefKey, datastore.AddressRef]{
			datastore.AddressRefByChainSelector(chainSelector),
			datastore.AddressRefByType(datastore.ContractType(contractType)),
			datastore.AddressRefByQualifier(qualifier),
		}
		addresses = v.Addresses().Filter(filters...)
	case datastore.MutableDataStore:
		filters := []datastore.FilterFunc[datastore.AddressRefKey, datastore.AddressRef]{
			datastore.AddressRefByChainSelector(chainSelector),
			datastore.AddressRefByType(datastore.ContractType(contractType)),
			datastore.AddressRefByQualifier(qualifier),
		}
		addresses = v.Addresses().Filter(filters...)
	default:
		return "", fmt.Errorf("unsupported datastore type: %T", ds)
	}

	if len(addresses) > 0 {
		return addresses[0].Address, nil
	}
	return "", fmt.Errorf("contract of type %s not found for chain %d with qualifier %q", contractType, chainSelector, qualifier)
}
