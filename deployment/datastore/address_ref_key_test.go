package datastore

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
)

func TestAddressRefKey_Equals(t *testing.T) {
	tests := []struct {
		name     string
		key1     AddressRefKey
		key2     AddressRefKey
		expected bool
	}{
		{
			name:     "Identical keys",
			key1:     NewAddressRefKey(1, ContractType("typeA"), semver.MustParse("1.0.0"), "qualifier1"),
			key2:     NewAddressRefKey(1, ContractType("typeA"), semver.MustParse("1.0.0"), "qualifier1"),
			expected: true,
		},
		{
			name:     "Different chainSelector",
			key1:     NewAddressRefKey(1, ContractType("typeA"), semver.MustParse("1.0.0"), "qualifier1"),
			key2:     NewAddressRefKey(2, ContractType("typeA"), semver.MustParse("1.0.0"), "qualifier1"),
			expected: false,
		},
		{
			name:     "Different contractType",
			key1:     NewAddressRefKey(1, ContractType("typeA"), semver.MustParse("1.0.0"), "qualifier1"),
			key2:     NewAddressRefKey(1, ContractType("typeB"), semver.MustParse("1.0.0"), "qualifier1"),
			expected: false,
		},
		{
			name:     "Different version",
			key1:     NewAddressRefKey(1, ContractType("typeA"), semver.MustParse("1.0.0"), "qualifier1"),
			key2:     NewAddressRefKey(1, ContractType("typeA"), semver.MustParse("2.0.0"), "qualifier1"),
			expected: false,
		},
		{
			name:     "Different qualifier",
			key1:     NewAddressRefKey(1, ContractType("typeA"), semver.MustParse("1.0.0"), "qualifier1"),
			key2:     NewAddressRefKey(1, ContractType("typeA"), semver.MustParse("1.0.0"), "qualifier2"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.key1.Equals(tt.key2))
		})
	}
}

func TestNewAddressRefKey(t *testing.T) {
	version := semver.MustParse("1.0.0")
	key := NewAddressRefKey(1, ContractType("typeA"), version, "qualifier1")

	assert.Equal(t, uint64(1), key.ChainSelector(), "ChainSelector should match")
	assert.Equal(t, ContractType("typeA"), key.Type(), "ContractType should match")
	assert.Equal(t, version, key.Version(), "Version should match")
	assert.Equal(t, "qualifier1", key.Qualifier(), "Qualifier should match")
}
