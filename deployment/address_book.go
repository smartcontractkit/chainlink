package deployment

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	chainsel "github.com/smartcontractkit/chain-selectors"
)

var (
	ErrInvalidChainSelector = fmt.Errorf("invalid chain selector")
	ErrInvalidAddress       = fmt.Errorf("invalid address")
	ErrChainNotFound        = fmt.Errorf("chain not found")
)

// ContractType is a simple string type for identifying contract types.
type ContractType string

var (
	Version1_0_0     = *semver.MustParse("1.0.0")
	Version1_1_0     = *semver.MustParse("1.1.0")
	Version1_2_0     = *semver.MustParse("1.2.0")
	Version1_5_0     = *semver.MustParse("1.5.0")
	Version1_6_0_dev = *semver.MustParse("1.6.0-dev")
)

type TypeAndVersion struct {
	Type    ContractType
	Version semver.Version
}

func (tv TypeAndVersion) String() string {
	return fmt.Sprintf("%s %s", tv.Type, tv.Version.String())
}

func (tv TypeAndVersion) Equal(other TypeAndVersion) bool {
	return tv.String() == other.String()
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
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		return TypeAndVersion{}, fmt.Errorf("invalid type and version string: %s", s)
	}
	v, err := semver.NewVersion(parts[1])
	if err != nil {
		return TypeAndVersion{}, err
	}
	return TypeAndVersion{
		Type:    ContractType(parts[0]),
		Version: *v,
	}, nil
}

func NewTypeAndVersion(t ContractType, v semver.Version) TypeAndVersion {
	return TypeAndVersion{
		Type:    t,
		Version: v,
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
}

type AddressBookMap struct {
	AddressesByChain map[uint64]map[string]TypeAndVersion
}

// Save will save an address for a given chain selector. It will error if there is a conflicting existing address.
func (m *AddressBookMap) Save(chainSelector uint64, address string, typeAndVersion TypeAndVersion) error {
	_, exists := chainsel.ChainBySelector(chainSelector)
	if !exists {
		return errors.Wrapf(ErrInvalidChainSelector, "chain selector %d", chainSelector)
	}
	if address == "" || address == common.HexToAddress("0x0").Hex() {
		return errors.Wrap(ErrInvalidAddress, "address cannot be empty")
	}
	if common.IsHexAddress(address) {
		// IMPORTANT: WE ALWAYS STANDARDIZE ETHEREUM ADDRESS STRINGS TO EIP55
		address = common.HexToAddress(address).Hex()
	} else {
		return errors.Wrapf(ErrInvalidAddress, "address %s is not a valid Ethereum address, only Ethereum addresses supported", address)
	}
	if typeAndVersion.Type == "" {
		return fmt.Errorf("type cannot be empty")
	}
	if _, exists := m.AddressesByChain[chainSelector]; !exists {
		// First time chain add, create map
		m.AddressesByChain[chainSelector] = make(map[string]TypeAndVersion)
	}
	if _, exists := m.AddressesByChain[chainSelector][address]; exists {
		return fmt.Errorf("address %s already exists for chain %d", address, chainSelector)
	}
	m.AddressesByChain[chainSelector][address] = typeAndVersion
	return nil
}

func (m *AddressBookMap) Addresses() (map[uint64]map[string]TypeAndVersion, error) {
	return m.AddressesByChain, nil
}

func (m *AddressBookMap) AddressesForChain(chainSelector uint64) (map[string]TypeAndVersion, error) {
	_, exists := chainsel.ChainBySelector(chainSelector)
	if !exists {
		return nil, errors.Wrapf(ErrInvalidChainSelector, "chain selector %d", chainSelector)
	}
	if _, exists := m.AddressesByChain[chainSelector]; !exists {
		return nil, errors.Wrapf(ErrChainNotFound, "chain selector %d", chainSelector)
	}
	return m.AddressesByChain[chainSelector], nil
}

// Merge will merge the addresses from another address book into this one.
// It will error on any existing addresses.
func (m *AddressBookMap) Merge(ab AddressBook) error {
	addresses, err := ab.Addresses()
	if err != nil {
		return err
	}
	for chain, chainAddresses := range addresses {
		for address, typeAndVersions := range chainAddresses {
			if err := m.Save(chain, address, typeAndVersions); err != nil {
				return err
			}
		}
	}
	return nil
}

// TODO: Maybe could add an environment argument
// which would ensure only mainnet/testnet chain selectors are used
// for further safety?
func NewMemoryAddressBook() *AddressBookMap {
	return &AddressBookMap{
		AddressesByChain: make(map[uint64]map[string]TypeAndVersion),
	}
}

func NewMemoryAddressBookFromMap(addressesByChain map[uint64]map[string]TypeAndVersion) *AddressBookMap {
	return &AddressBookMap{
		AddressesByChain: addressesByChain,
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

	return "", fmt.Errorf("not found")
}
