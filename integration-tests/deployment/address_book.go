package deployment

import "fmt"

// AddressBook is a simple interface for storing and retrieving contract addresses across
// chains. It is family agnostic.
type AddressBook interface {
	Save(chainSelector uint64, address string, typeAndVersion string) error
	Addresses() (map[uint64]map[string]string, error)
	AddressesForChain(chain uint64) (map[string]string, error)
	// Allows for merging address books (e.g. new deployments with existing ones)
	Merge(other AddressBook) error
}

type AddressBookMap struct {
	AddressesByChain map[uint64]map[string]string
}

func (m *AddressBookMap) Save(chainSelector uint64, address string, typeAndVersion string) error {
	if _, exists := m.AddressesByChain[chainSelector]; !exists {
		// First time chain add, create map
		m.AddressesByChain[chainSelector] = make(map[string]string)
	}
	if _, exists := m.AddressesByChain[chainSelector][address]; exists {
		return fmt.Errorf("address %s already exists for chain %d", address, chainSelector)
	}
	m.AddressesByChain[chainSelector][address] = typeAndVersion
	return nil
}

func (m *AddressBookMap) Addresses() (map[uint64]map[string]string, error) {
	return m.AddressesByChain, nil
}

func (m *AddressBookMap) AddressesForChain(chain uint64) (map[string]string, error) {
	if _, exists := m.AddressesByChain[chain]; !exists {
		return nil, fmt.Errorf("chain %d not found", chain)
	}
	return m.AddressesByChain[chain], nil
}

// Attention this will mutate existing book
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

func NewMemoryAddressBook() *AddressBookMap {
	return &AddressBookMap{
		AddressesByChain: make(map[uint64]map[string]string),
	}
}
