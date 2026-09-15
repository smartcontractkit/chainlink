package shared

import (
	"context"
	"fmt"
	"strings"

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

// QualifierFromParts joins the parts that identify one instance into a datastore qualifier.
//
// Each part is quoted rather than joined on a bare separator: parts are often caller-supplied
// free-form text, so a plain join would let ("A", "B/C") and ("A/B", "C") produce the same
// qualifier and collide. Quoting escapes any separator inside a part, so distinct inputs always
// yield distinct qualifiers.
func QualifierFromParts(parts ...string) string {
	quoted := make([]string, 0, len(parts))
	for _, p := range parts {
		quoted = append(quoted, fmt.Sprintf("%q", p))
	}

	return strings.Join(quoted, "/")
}

// TokenPoolLookupTableQualifier returns the datastore qualifier for a Solana token-pool lookup
// table, which is uniquely identified by (token mint, pool type, metadata).
//
// Each component is quoted rather than joined on a bare separator. metadata is caller-supplied
// free-form text, so a plain "a/b/c" join would let ("A", "B/C") and ("A/B", "C") produce the same
// qualifier and collide in the datastore. Quoting escapes any separator inside a component, so
// distinct inputs always yield distinct qualifiers.
func TokenPoolLookupTableQualifier(tokenPubKey, poolType, metadata string) string {
	return QualifierFromParts(tokenPubKey, poolType, metadata)
}
