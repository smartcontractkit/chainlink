package shared

import (
	"context"
	"errors"
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

// AddressKey identifies an address in an address book. The chain selector is part of the key
// because the same address can legitimately be used on multiple chains.
type AddressKey struct {
	ChainSelector uint64
	Address       string
}

// QualifiersForAddressBook returns explicit qualifiers for every address in an address book.
// Use it only when the address book is scoped to one logical instance, such as one token pool.
func QualifiersForAddressBook(addressBook deployment.AddressBook, qualifier string) (map[AddressKey]string, error) {
	addrs, err := addressBook.Addresses()
	if err != nil {
		return nil, err
	}
	qualifiers := make(map[AddressKey]string)
	for chainSelector, chainAddresses := range addrs {
		for address := range chainAddresses {
			qualifiers[AddressKey{ChainSelector: chainSelector, Address: address}] = qualifier
		}
	}
	return qualifiers, nil
}

// QualifiersByAddress scopes an address-only qualifier map to the chains present in an address
// book. It is useful for legacy callers that already collect qualifiers by address while
// constructing a single-chain address book.
func QualifiersByAddress(addressBook deployment.AddressBook, byAddress map[string]string) (map[AddressKey]string, error) {
	addrs, err := addressBook.Addresses()
	if err != nil {
		return nil, err
	}
	qualifiers := make(map[AddressKey]string, len(byAddress))
	for chainSelector, chainAddresses := range addrs {
		for address := range chainAddresses {
			if qualifier, ok := byAddress[address]; ok {
				qualifiers[AddressKey{ChainSelector: chainSelector, Address: address}] = qualifier
			}
		}
	}
	return qualifiers, nil
}

// PopulateDataStore converts a changeset's OWN freshly-deployed address book into a datastore.
// Callers must provide a qualifier for every multi-instance address they deployed; all other refs
// use the empty qualifier. It FAILS if two refs map to the same datastore key
// (chain, type, version, qualifier) — the address is not part of the key — because silently
// keeping only one would lose data nondeterministically. It must never be called on a merged or
// pre-existing address book that cannot carry qualifiers; resolve existing state from a ref slice
// via CollectAddressRefs instead.
func PopulateDataStore(addressBook deployment.AddressBook, qualifiers map[AddressKey]string) (*datastore.MemoryDataStore, error) {
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
				Qualifier:     qualifiers[AddressKey{ChainSelector: chainselector, Address: addr}],
			}

			// If the address book has labels, we need to add them to the addressRef
			if !typever.Labels.IsEmpty() {
				ref.Labels = datastore.NewLabelSet(typever.Labels.List()...)
			}

			if err = ds.Addresses().Add(ref); err != nil {
				if errors.Is(err, datastore.ErrAddressRefExists) {
					return nil, fmt.Errorf(
						"address book has multiple refs that map to the same datastore key (chain=%d, type=%s, version=%s, qualifier=%q); provide distinct qualifiers for each multi-instance contract",
						chainselector, typever.Type, typever.Version, ref.Qualifier)
				}
				return nil, err
			}
		}
	}

	return ds, nil
}

// CollectAddressRefs returns a plain (unkeyed) slice of address refs drawn from both the
// environment's datastore and its existing address book, including labels. It does not key or
// dedup, so duplicate/multi-instance refs are all retained. Use it to resolve chain-singleton
// contracts (e.g. FeeQuoter) over existing state without going through a keyed datastore. It
// returns an error if either source cannot be read, so callers do not silently resolve from an
// incomplete set.
func CollectAddressRefs(e deployment.Environment) ([]datastore.AddressRef, error) {
	var refs []datastore.AddressRef
	if e.DataStore != nil {
		r, err := e.DataStore.Addresses().Fetch()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch address refs from environment datastore: %w", err)
		}
		refs = append(refs, r...)
	}
	if e.ExistingAddresses != nil {
		addrs, err := e.ExistingAddresses.Addresses()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch addresses from environment address book: %w", err)
		}
		for sel, chainAddresses := range addrs {
			for addr, tv := range chainAddresses {
				ref := datastore.AddressRef{
					ChainSelector: sel,
					Address:       addr,
					Type:          datastore.ContractType(tv.Type),
					Version:       &tv.Version,
				}
				if !tv.Labels.IsEmpty() {
					ref.Labels = datastore.NewLabelSet(tv.Labels.List()...)
				}
				refs = append(refs, ref)
			}
		}
	}
	return refs, nil
}

// TokenPoolLookupTableQualifier returns the datastore qualifier for a Solana token-pool lookup
// table, which is uniquely identified by (token mint, pool type, metadata).
func TokenPoolLookupTableQualifier(tokenPubKey, poolType, metadata string) string {
	return fmt.Sprintf("%s/%s/%s", tokenPubKey, poolType, metadata)
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
