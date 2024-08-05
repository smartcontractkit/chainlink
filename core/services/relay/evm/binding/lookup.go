package binding

import (
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type readValues struct {
	address  string
	contract string
	method   string
}

type lookup struct {
	mu              sync.RWMutex
	contractMethods map[string][]string
	readNameValues  map[string]readValues
}

func newLookup() *lookup {
	return &lookup{
		contractMethods: make(map[string][]string),
		readNameValues:  make(map[string]readValues),
	}
}

func (l *lookup) addMethodForContract(contract, method string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	methods, exists := l.contractMethods[contract]
	if !exists {
		methods = []string{}
	}

	l.contractMethods[contract] = append(methods, method)
}

func (l *lookup) bindAddressForContract(contract, address string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, method := range l.contractMethods[contract] {
		readName := types.BoundContract{
			Address: address,
			Name:    contract,
		}.ReadIdentifier(method)

		l.readNameValues[readName] = readValues{
			address:  address,
			contract: contract,
			method:   method,
		}
	}
}

func (l *lookup) unbindAddressForContract(contract, address string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, method := range l.contractMethods[contract] {
		readName := types.BoundContract{
			Address: address,
			Name:    contract,
		}.ReadIdentifier(method)

		delete(l.readNameValues, readName)
	}
}

func (l *lookup) getContractForReadName(readName string) (readValues, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	contract, ok := l.readNameValues[readName]

	return contract, ok
}
