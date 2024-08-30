package deployment

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestAddressBook(t *testing.T) {
	ab := NewMemoryAddressBook()
	err := ab.Save(1, "0x1", "OnRamp 1.0.0")
	require.NoError(t, err)
	// Duplicate address will error
	err = ab.Save(1, "0x1", "OnRamp 1.0.0")
	require.Error(t, err)
	// Distinct address same TV will not
	err = ab.Save(1, "0x2", "OnRamp 1.0.0")
	require.NoError(t, err)
	// Same address different chain will not error
	err = ab.Save(2, "0x1", "OnRamp 1.0.0")
	require.NoError(t, err)
	// We can save different versions of the same contract
	err = ab.Save(2, "0x2", "OnRamp 1.2.0")
	require.NoError(t, err)

	addresses, err := ab.Addresses()
	require.NoError(t, err)
	assert.DeepEqual(t, addresses, map[uint64]map[string]string{
		1: {
			"0x1": "OnRamp 1.0.0",
			"0x2": "OnRamp 1.0.0",
		},
		2: {
			"0x1": "OnRamp 1.0.0",
			"0x2": "OnRamp 1.2.0",
		},
	})

	// Test merge
	ab2 := NewMemoryAddressBook()
	require.NoError(t, ab2.Save(3, "0x3", "OnRamp 1.0.0"))
	require.NoError(t, ab.Merge(ab2))
	// Other address book should remain unchanged.
	addresses, err = ab2.Addresses()
	require.NoError(t, err)
	assert.DeepEqual(t, addresses, map[uint64]map[string]string{
		3: {
			"0x3": "OnRamp 1.0.0",
		},
	})
	// Existing addressbook should contain the new elements.
	addresses, err = ab.Addresses()
	require.NoError(t, err)
	assert.DeepEqual(t, addresses, map[uint64]map[string]string{
		1: {
			"0x1": "OnRamp 1.0.0",
			"0x2": "OnRamp 1.0.0",
		},
		2: {
			"0x1": "OnRamp 1.0.0",
			"0x2": "OnRamp 1.2.0",
		},
		3: {
			"0x3": "OnRamp 1.0.0",
		},
	})

	// Merge to an existing chain.
	require.NoError(t, ab2.Save(2, "0x3", "OffRamp 1.0.0"))
}
