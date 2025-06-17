package read

import (
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type readValues struct {
	address  string
	contract string
	readName string
}

// lookup provides basic utilities for mapping a complete readIdentifier to
// finite contract read information
type lookup struct {
	mu sync.RWMutex
	// contractReadNames maps a contract name to all available readNames (method, log, event, etc.)
	contractReadNames map[string][]string
	// readIdentifiers maps from a complete readIdentifier string to finite read data
	// a readIdentifier is a combination of address, contract, and readName as a concatenated string
	readIdentifiers map[string]readValues
}

func newLookup() *lookup {
	return &lookup{
		contractReadNames: make(map[string][]string),
		readIdentifiers:   make(map[string]readValues),
	}
}

func (l *lookup) addReadNameForContract(contract, readName string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.contractReadNames[contract] = append(l.contractReadNames[contract], readName)
}

func (l *lookup) bindAddressForContract(contract, address string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, readName := range l.contractReadNames[contract] {
		readIdentifier := types.BoundContract{
			Address: address,
			Name:    contract,
		}.ReadIdentifier(readName)

		l.readIdentifiers[readIdentifier] = readValues{
			address:  address,
			contract: contract,
			readName: readName,
		}
	}
}

func (l *lookup) unbindAddressForContract(contract, address string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, readName := range l.contractReadNames[contract] {
		readIdentifier := types.BoundContract{
			Address: address,
			Name:    contract,
		}.ReadIdentifier(readName)

		delete(l.readIdentifiers, readIdentifier)
	}
}

func (l *lookup) getContractForReadName(readName string) (readValues, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	contract, ok := l.readIdentifiers[readName]

	return contract, ok
}
