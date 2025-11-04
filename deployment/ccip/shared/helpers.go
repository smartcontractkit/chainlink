package shared

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
)

const (
	CapabilityLabelledName = "ccip"
	CapabilityVersion      = "v1.0.0"
)

var CCIPCapabilityID = utils.Keccak256Fixed(MustABIEncode(`[{"type": "string"}, {"type": "string"}]`, CapabilityLabelledName, CapabilityVersion))

func GetCCIPDonsFromCapRegistry(ctx context.Context, capRegistry *capabilities_registry.CapabilitiesRegistry) ([]capabilities_registry.CapabilitiesRegistryDONInfo, error) {
	if capRegistry == nil {
		return nil, nil
	}
	// Get the all Dons from the capabilities registry
	allDons, err := capRegistry.GetDONs(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("failed to get all Dons from capabilities registry: %w", err)
	}
	ccipDons := make([]capabilities_registry.CapabilitiesRegistryDONInfo, 0, len(allDons))
	for _, don := range allDons {
		for _, capConfig := range don.CapabilityConfigurations {
			if capConfig.CapabilityId == CCIPCapabilityID {
				ccipDons = append(ccipDons, don)
				break
			}
		}
	}

	return ccipDons, nil
}

func MustABIEncode(abiString string, args ...any) []byte {
	encoded, err := utils.ABIEncode(abiString, args...)
	if err != nil {
		panic(err)
	}
	return encoded
}

func PopulateDataStore(addressBook deployment.AddressBook) (*datastore.MemoryDataStore, error) {
	addrs, err := addressBook.Addresses()
	if err != nil {
		return nil, err
	}

	ds := datastore.NewMemoryDataStore()
	for chainselector, chainAddresses := range addrs {
		for addr, typever := range chainAddresses {
			ref := datastore.AddressRef{
				ChainSelector: chainselector,
				Address:       addr,
				Type:          datastore.ContractType(typever.Type),
				Version:       &typever.Version,
				// Since the address book does not have a qualifier, we use the address and type as a
				// unique identifier for the addressRef. Otherwise, we would have some clashes in the
				// between address refs.
				Qualifier: fmt.Sprintf("%s-%s", addr, typever.Type),
			}

			// If the address book has labels, we need to add them to the addressRef
			if !typever.Labels.IsEmpty() {
				ref.Labels = datastore.NewLabelSet(typever.Labels.List()...)
			}

			if err = ds.Addresses().Add(ref); err != nil {
				return nil, err
			}
		}
	}

	return ds, nil
}
