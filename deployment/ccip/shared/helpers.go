package shared

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	fqv2ops "github.com/smartcontractkit/chainlink-ccip/chains/evm/deployment/v2_0_0/operations/fee_quoter"
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

// AdminSlot is the specific storage location defined by EIP-1967 to store the Proxy Admin address.
//
// Background:
// Proxies must store their admin address in a storage slot that does not collide
// with the storage layout of the Logic (Implementation) contract.
//
// Formula:
// bytes32(uint256(keccak256('eip1967.proxy.admin')) - 1)
//
// Result:
// 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103
//
// This guarantees that reading this slot returns the Admin address without
// interfering with the ERC20 token state (Balances, Supply, etc).
var AdminSlot = common.HexToHash("0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103")

// TUPImplementationSlot is the specific storage location defined by EIP-1967 to store the
// address of the Logic (Implementation) contract.
//
// Background:
// Proxies must store the address of the code they delegate to in a slot that does
// not collide with the storage layout of the Logic contract itself.
//
// Formula:
// bytes32(uint256(keccak256('eip1967.proxy.implementation')) - 1)
//
// Result:
// 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc
//
// Reading this slot returns the address of the BurnMintERC20Transparent contract.
var TUPImplementationSlot = common.HexToHash("0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc")

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

// ResolveFeeQuoterAddressAndVersion returns the FeeQuoter with the highest semver for a chain.
func ResolveFeeQuoterAddressAndVersion(
	addresses []datastore.AddressRef,
	chainSel uint64,
) (common.Address, semver.Version, error) {
	var bestRef datastore.AddressRef
	var bestVersion *semver.Version

	for _, ref := range addresses {
		if ref.ChainSelector != chainSel {
			continue
		}
		if ref.Type != datastore.ContractType(fqv2ops.ContractType) {
			continue
		}
		if ref.Version == nil {
			continue
		}
		if bestVersion == nil || ref.Version.GreaterThan(bestVersion) {
			bestVersion = ref.Version
			bestRef = ref
		}
	}

	if bestVersion == nil {
		return common.Address{}, semver.Version{}, fmt.Errorf("no fee quoter address found for chain %d", chainSel)
	}

	if !common.IsHexAddress(bestRef.Address) {
		return common.Address{}, semver.Version{}, fmt.Errorf("invalid fee quoter address %q for chain %d", bestRef.Address, chainSel)
	}

	return common.HexToAddress(bestRef.Address), *bestVersion, nil
}
