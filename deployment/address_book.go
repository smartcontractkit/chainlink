package deployment

import (
	"fmt"
	"strings"
	"sync"

	"golang.org/x/exp/maps"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	chainsel "github.com/smartcontractkit/chain-selectors"
)

var (
	ErrInvalidChainSelector = errors.New("invalid chain selector")
	ErrInvalidAddress       = errors.New("invalid address")
	ErrChainNotFound        = errors.New("chain not found")
)

// ContractType is a simple string type for identifying contract types.
type ContractType string

var (
	Version1_0_0 = *semver.MustParse("1.0.0")
	Version1_1_0 = *semver.MustParse("1.1.0")
	Version1_2_0 = *semver.MustParse("1.2.0")
	Version1_5_0 = *semver.MustParse("1.5.0")
	Version1_5_1 = *semver.MustParse("1.5.1")
	Version1_6_0 = *semver.MustParse("1.6.0")
)

type TypeAndVersion struct {
	Type    ContractType   `json:"Type"`
	Version semver.Version `json:"Version"`
	Labels  LabelSet       `json:"Labels,omitempty"`
}

func (tv TypeAndVersion) String() string {
	if len(tv.Labels) == 0 {
		return fmt.Sprintf("%s %s", tv.Type, tv.Version.String())
	}

	// Use the LabelSet's String method for sorted labels
	sortedLabels := tv.Labels.String()
	return fmt.Sprintf("%s %s %s",
		tv.Type,
		tv.Version.String(),
		sortedLabels,
	)
}

func (tv TypeAndVersion) Equal(other TypeAndVersion) bool {
	// Compare Type
	if tv.Type != other.Type {
		return false
	}
	// Compare Versions
	if !tv.Version.Equal(&other.Version) {
		return false
	}
	// Compare Labels
	return tv.Labels.Equal(other.Labels)
}

func MustTypeAndVersionFromString(s string) TypeAndVersion {
	tv, err := TypeAndVersionFromString(s)
	if err != nil {
		panic(err)
	}
	return tv
}

// Note this will become useful for validation. When we want
// to assert an onchain call to typeAndVersion yields whats expected.
func TypeAndVersionFromString(s string) (TypeAndVersion, error) {
	parts := strings.Fields(s) // Ignores consecutive spaces
	if len(parts) < 2 {
		return TypeAndVersion{}, fmt.Errorf("invalid type and version string: %s", s)
	}
	v, err := semver.NewVersion(parts[1])
	if err != nil {
		return TypeAndVersion{}, err
	}
	labels := make(LabelSet)
	if len(parts) > 2 {
		labels = NewLabelSet(parts[2:]...)
	}
	return TypeAndVersion{
		Type:    ContractType(parts[0]),
		Version: *v,
		Labels:  labels,
	}, nil
}

func NewTypeAndVersion(t ContractType, v semver.Version) TypeAndVersion {
	return TypeAndVersion{
		Type:    t,
		Version: v,
		Labels:  make(LabelSet), // empty set,
	}
}

// AddressBook is a simple interface for storing and retrieving contract addresses across
// chains. It is family agnostic as the keys are chain selectors.
// We store rather than derive typeAndVersion as some contracts do not support it.
// For ethereum addresses are always stored in EIP55 format.
type AddressBook interface {
	Save(chainSelector uint64, address string, tv TypeAndVersion) error
	Addresses() (map[uint64]map[string]TypeAndVersion, error)
	AddressesForChain(chain uint64) (map[string]TypeAndVersion, error)
	// Allows for merging address books (e.g. new deployments with existing ones)
	Merge(other AddressBook) error
	Remove(ab AddressBook) error
}

type AddressesByChain map[uint64]map[string]TypeAndVersion

type AddressBookMap struct {
	addressesByChain AddressesByChain
	mtx              sync.RWMutex
}

// Save will save an address for a given chain selector. It will error if there is a conflicting existing address.
func (m *AddressBookMap) save(chainSelector uint64, address string, typeAndVersion TypeAndVersion) error {
	family, err := chainsel.GetSelectorFamily(chainSelector)
	if err != nil {
		return errors.Wrapf(ErrInvalidChainSelector, "chain selector %d", chainSelector)
	}
	if family == chainsel.FamilyEVM {
		if address == "" || address == common.HexToAddress("0x0").Hex() {
			return errors.Wrap(ErrInvalidAddress, "address cannot be empty")
		}
		if common.IsHexAddress(address) {
			// IMPORTANT: WE ALWAYS STANDARDIZE ETHEREUM ADDRESS STRINGS TO EIP55
			address = common.HexToAddress(address).Hex()
		} else {
			return errors.Wrapf(ErrInvalidAddress, "address %s is not a valid Ethereum address, only Ethereum addresses supported for EVM chains", address)
		}
	}

	// TODO NONEVM-960: Add validation for non-EVM chain addresses

	if typeAndVersion.Type == "" {
		return errors.New("type cannot be empty")
	}

	if _, exists := m.addressesByChain[chainSelector]; !exists {
		// First time chain add, create map
		m.addressesByChain[chainSelector] = make(map[string]TypeAndVersion)
	}
	if _, exists := m.addressesByChain[chainSelector][address]; exists {
		return fmt.Errorf("address %s already exists for chain %d", address, chainSelector)
	}
	m.addressesByChain[chainSelector][address] = typeAndVersion
	return nil
}

// Save will save an address for a given chain selector. It will error if there is a conflicting existing address.
// thread safety version of the save method
func (m *AddressBookMap) Save(chainSelector uint64, address string, typeAndVersion TypeAndVersion) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return m.save(chainSelector, address, typeAndVersion)
}

func (m *AddressBookMap) Addresses() (map[uint64]map[string]TypeAndVersion, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	// maps are mutable and pass via a pointer
	// creating a copy of the map to prevent concurrency
	// read and changes outside object-bound
	return m.cloneAddresses(m.addressesByChain), nil
}

func (m *AddressBookMap) AddressesForChain(chainSelector uint64) (map[string]TypeAndVersion, error) {
	_, err := chainsel.GetChainIDFromSelector(chainSelector)
	if err != nil {
		return nil, errors.Wrapf(ErrInvalidChainSelector, "chain selector %d", chainSelector)
	}

	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if _, exists := m.addressesByChain[chainSelector]; !exists {
		return nil, errors.Wrapf(ErrChainNotFound, "chain selector %d", chainSelector)
	}

	// maps are mutable and pass via a pointer
	// creating a copy of the map to prevent concurrency
	// read and changes outside object-bound
	return maps.Clone(m.addressesByChain[chainSelector]), nil
}

// Merge will merge the addresses from another address book into this one.
// It will error on any existing addresses.
func (m *AddressBookMap) Merge(ab AddressBook) error {
	addresses, err := ab.Addresses()
	if err != nil {
		return err
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()

	for chainSelector, chainAddresses := range addresses {
		for address, typeAndVersion := range chainAddresses {
			if err := m.save(chainSelector, address, typeAndVersion); err != nil {
				return err
			}
		}
	}
	return nil
}

// Remove removes the address book addresses specified via the argument from the AddressBookMap.
// Errors if all the addresses in the given address book are not contained in the AddressBookMap.
func (m *AddressBookMap) Remove(ab AddressBook) error {
	addresses, err := ab.Addresses()
	if err != nil {
		return err
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()

	// State of m.addressesByChain storage must not be changed in case of an error
	// need to do double iteration over the address book. First validation, second actual deletion
	for chainSelector, chainAddresses := range addresses {
		for address := range chainAddresses {
			if _, exists := m.addressesByChain[chainSelector][address]; !exists {
				return errors.New("AddressBookMap does not contain address from the given address book")
			}
		}
	}

	for chainSelector, chainAddresses := range addresses {
		for address := range chainAddresses {
			delete(m.addressesByChain[chainSelector], address)
		}
	}

	return nil
}

// cloneAddresses creates a deep copy of map[uint64]map[string]TypeAndVersion object
func (m *AddressBookMap) cloneAddresses(input map[uint64]map[string]TypeAndVersion) map[uint64]map[string]TypeAndVersion {
	result := make(map[uint64]map[string]TypeAndVersion)
	for chainSelector, chainAddresses := range input {
		result[chainSelector] = maps.Clone(chainAddresses)
	}
	return result
}

// TODO: Maybe could add an environment argument
// which would ensure only mainnet/testnet chain selectors are used
// for further safety?
func NewMemoryAddressBook() *AddressBookMap {
	return &AddressBookMap{
		addressesByChain: make(map[uint64]map[string]TypeAndVersion),
	}
}

func NewMemoryAddressBookFromMap(addressesByChain map[uint64]map[string]TypeAndVersion) *AddressBookMap {
	return &AddressBookMap{
		addressesByChain: addressesByChain,
	}
}

// SearchAddressBook search an address book for a given chain and contract type and return the first matching address.
func SearchAddressBook(ab AddressBook, chain uint64, typ ContractType) (string, error) {
	addrs, err := ab.AddressesForChain(chain)
	if err != nil {
		return "", err
	}

	for addr, tv := range addrs {
		if tv.Type == typ {
			return addr, nil
		}
	}

	return "", errors.New("not found")
}

func AddressBookContains(ab AddressBook, chain uint64, addrToFind string) (bool, error) {
	addrs, err := ab.AddressesForChain(chain)
	if err != nil {
		return false, err
	}

	for addr := range addrs {
		if addr == addrToFind {
			return true, nil
		}
	}

	return false, nil
}

type typeVersionKey struct {
	Type    ContractType
	Version string
	Labels  string // store labels in a canonical form (comma-joined sorted list)
}

func tvKey(tv TypeAndVersion) typeVersionKey {
	sortedLabels := tv.Labels.String()
	return typeVersionKey{
		Type:    tv.Type,
		Version: tv.Version.String(),
		Labels:  sortedLabels,
	}
}

// AddressesContainBundle checks if the addresses
// contains a single instance of all the addresses in the bundle.
// It returns an error if there are more than one instance of a contract.
func AddressesContainBundle(addrs map[string]TypeAndVersion, wantTypes []TypeAndVersion) (bool, error) {
	// Count how many times each wanted TypeAndVersion is found
	counts := make(map[typeVersionKey]int)
	for _, wantTV := range wantTypes {
		wantKey := tvKey(wantTV)
		for _, haveTV := range addrs {
			if wantTV.Equal(haveTV) {
				// They match exactly (Type, Version, Labels)
				counts[wantKey]++
				if counts[wantKey] > 1 {
					return false, fmt.Errorf("found more than one instance of contract %s %s (labels=%s)",
						wantTV.Type, wantTV.Version.String(), wantTV.Labels.String())
				}
			}
		}
	}

	// Ensure we found *all* wantTypes exactly once
	return len(counts) == len(wantTypes), nil
}

// AddLabel adds a string to the LabelSet in the TypeAndVersion.
func (tv *TypeAndVersion) AddLabel(label string) {
	if tv.Labels == nil {
		tv.Labels = make(LabelSet)
	}
	tv.Labels.Add(label)
}
