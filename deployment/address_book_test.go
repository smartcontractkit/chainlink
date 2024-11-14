package deployment

import (
	"errors"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestAddressBook_Save(t *testing.T) {
	ab := NewMemoryAddressBook()
	onRamp100 := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp110 := NewTypeAndVersion("OnRamp", Version1_1_0)
	addr1 := common.HexToAddress("0x1").String()
	addr2 := common.HexToAddress("0x2").String()

	err := ab.Save(chainsel.TEST_90000001.Selector, addr1, onRamp100)
	require.NoError(t, err)

	// Invalid address
	err = ab.Save(chainsel.TEST_90000001.Selector, "asdlfkj", onRamp100)
	require.Error(t, err)
	assert.Equal(t, errors.Is(err, ErrInvalidAddress), true, "err %s", err)

	// Valid chain but not present.
	_, err = ab.AddressesForChain(chainsel.TEST_90000002.Selector)
	assert.Equal(t, errors.Is(err, ErrChainNotFound), true, "err %s", err)

	// Invalid selector
	err = ab.Save(0, addr1, onRamp100)
	require.Error(t, err)
	assert.Equal(t, errors.Is(err, ErrInvalidChainSelector), true)

	// Duplicate
	err = ab.Save(chainsel.TEST_90000001.Selector, addr1, onRamp100)
	require.Error(t, err)

	// Zero address
	err = ab.Save(chainsel.TEST_90000001.Selector, common.HexToAddress("0x0").Hex(), onRamp100)
	require.Error(t, err)

	// Distinct address same TV will not
	err = ab.Save(chainsel.TEST_90000001.Selector, addr2, onRamp100)
	require.NoError(t, err)
	// Same address different chain will not error
	err = ab.Save(chainsel.TEST_90000002.Selector, addr1, onRamp100)
	require.NoError(t, err)
	// We can save different versions of the same contract
	err = ab.Save(chainsel.TEST_90000002.Selector, addr2, onRamp110)
	require.NoError(t, err)

	addresses, err := ab.Addresses()
	require.NoError(t, err)
	assert.DeepEqual(t, addresses, map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp100,
			addr2: onRamp110,
		},
	})
}

func TestAddressBook_Merge(t *testing.T) {
	onRamp100 := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp110 := NewTypeAndVersion("OnRamp", Version1_1_0)
	addr1 := common.HexToAddress("0x1").String()
	addr2 := common.HexToAddress("0x2").String()
	a1 := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
		},
	})
	a2 := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110,
		},
	})
	require.NoError(t, a1.Merge(a2))

	addresses, err := a1.Addresses()
	require.NoError(t, err)
	assert.DeepEqual(t, addresses, map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110,
		},
	})

	// Merge with conflicting addresses should error
	a3 := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
		},
	})
	require.Error(t, a1.Merge(a3))
	// a1 should not have changed
	addresses, err = a1.Addresses()
	require.NoError(t, err)
	assert.DeepEqual(t, addresses, map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110,
		},
	})
}

func TestAddressBook_Remove(t *testing.T) {
	onRamp100 := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp110 := NewTypeAndVersion("OnRamp", Version1_1_0)
	addr1 := common.HexToAddress("0x1").String()
	addr2 := common.HexToAddress("0x2").String()
	addr3 := common.HexToAddress("0x3").String()

	baseAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110,
			addr3: onRamp110,
		},
	})

	copyOfBaseAB := NewMemoryAddressBookFromMap(baseAB.cloneAddresses(baseAB.addressesByChain))

	// this address book shouldn't be removed (state of baseAB not changed, error thrown)
	failAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr3: onRamp100, // doesn't exist in TEST_90000001.Selector
		},
	})
	require.Error(t, baseAB.Remove(failAB))
	require.EqualValues(t, baseAB, copyOfBaseAB)

	// this Address book should be removed without error
	successAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000002.Selector: {
			addr3: onRamp100,
		},
		chainsel.TEST_90000001.Selector: {
			addr2: onRamp100,
		},
	})

	expectingAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110},
	})

	require.NoError(t, baseAB.Remove(successAB))
	require.EqualValues(t, baseAB, expectingAB)
}

func TestAddressBook_ConcurrencyAndDeadlock(t *testing.T) {
	onRamp100 := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp110 := NewTypeAndVersion("OnRamp", Version1_1_0)

	baseAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			common.BigToAddress(big.NewInt(1)).String(): onRamp100,
		},
	})

	// concurrent writes
	var i int64
	wg := sync.WaitGroup{}
	for i = 2; i < 1000; i++ {
		wg.Add(1)
		go func(input int64) {
			require.NoError(t, baseAB.Save(
				chainsel.TEST_90000001.Selector,
				common.BigToAddress(big.NewInt(input)).String(),
				onRamp100,
			))
			wg.Done()
		}(i)
	}

	// concurrent reads
	for i = 0; i < 100; i++ {
		wg.Add(1)
		go func(input int64) {
			addresses, err := baseAB.Addresses()
			require.NoError(t, err)
			for chainSelector, chainAddresses := range addresses {
				// concurrent read chainAddresses from Addresses() method
				for address, _ := range chainAddresses {
					addresses[chainSelector][address] = onRamp110
				}

				// concurrent read chainAddresses from AddressesForChain() method
				chainAddresses, err = baseAB.AddressesForChain(chainSelector)
				require.NoError(t, err)
				for address, _ := range chainAddresses {
					_ = addresses[chainSelector][address]
				}
			}
			require.NoError(t, err)
			wg.Done()
		}(i)
	}

	// concurrent merges, starts from 1001 to avoid address conflicts
	for i = 1001; i < 1100; i++ {
		wg.Add(1)
		go func(input int64) {
			// concurrent merge
			additionalAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
				chainsel.TEST_90000002.Selector: {
					common.BigToAddress(big.NewInt(input)).String(): onRamp100,
				},
			})
			require.NoError(t, baseAB.Merge(additionalAB))
			wg.Done()
		}(i)
	}

	wg.Wait()
}
