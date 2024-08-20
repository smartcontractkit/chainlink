package deployment

import (
	"errors"
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
	offRamp100 := NewTypeAndVersion("OffRamp", Version1_0_0)
	addr1 := common.HexToAddress("0x1").String()
	addr2 := common.HexToAddress("0x2").String()
	addr3 := common.HexToAddress("0x3").String()

	err := ab.Save(chainsel.TEST_90000001.Selector, addr1, onRamp100)
	require.NoError(t, err)

	// Check input validation
	err = ab.Save(chainsel.TEST_90000001.Selector, "asdlfkj", onRamp100)
	require.Error(t, err)
	assert.Equal(t, errors.Is(err, ErrInvalidAddress), true, "err %s", err)
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

	// Test merge
	ab2 := NewMemoryAddressBook()
	require.NoError(t, ab2.Save(chainsel.TEST_90000003.Selector, addr3, onRamp100))
	require.NoError(t, ab.Merge(ab2))
	// Other address book should remain unchanged.
	addresses, err = ab2.Addresses()
	require.NoError(t, err)
	assert.DeepEqual(t, addresses, map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000003.Selector: {
			addr3: onRamp100,
		},
	})
	// Existing addressbook should contain the new elements.
	addresses, err = ab.Addresses()
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
		chainsel.TEST_90000003.Selector: {
			addr3: onRamp100,
		},
	})

	// Merge to an existing chain.
	require.NoError(t, ab2.Save(chainsel.TEST_90000002.Selector, addr3, offRamp100))
}
