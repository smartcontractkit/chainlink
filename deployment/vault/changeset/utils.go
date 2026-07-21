package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

func recipientAddressesFromNativeTransfers(transfers []types.NativeTransfer) []string {
	addresses := make([]string, len(transfers))
	for i, transfer := range transfers {
		addresses[i] = transfer.To
	}
	return addresses
}

func recipientAddressesFromERC20Transfers(transfers []types.ERC20Transfer) []string {
	addresses := make([]string, len(transfers))
	for i, transfer := range transfers {
		addresses[i] = transfer.Payee
	}
	return addresses
}

// lowestChainSelector returns the smallest chain-selector key in chains, and false when the map
// is empty. Map iteration order is not deterministic in Go, so callers that only need a single
// representative chain (e.g. to build deps / pick a deployer key) use this instead of grabbing an
// arbitrary key, keeping changeset output reproducible for the same input.
func lowestChainSelector[V any](chains map[uint64]V) (uint64, bool) {
	var lowest uint64
	found := false
	for sel := range chains {
		if !found || sel < lowest {
			lowest = sel
			found = true
		}
	}
	return lowest, found
}

func GetContractAddress(ds any, chainSelector uint64, contractType cldf.ContractType) (string, error) {
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

func GetContractAddressWithQualifier(ds any, chainSelector uint64, contractType cldf.ContractType, qualifier string) (string, error) {
	if ds == nil {
		return "", errors.New("datastore is nil")
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
