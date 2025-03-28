package datastore

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/deployment"
)

func TestAddressRefKey_Equals(t *testing.T) {
	version1 := semver.MustParse("1.0.0")
	version2 := semver.MustParse("2.0.0")

	key1 := NewAddressRefKey(1, deployment.ContractType("typeA"), version1, "qualifier1")
	key2 := NewAddressRefKey(1, deployment.ContractType("typeA"), version1, "qualifier1")
	key3 := NewAddressRefKey(2, deployment.ContractType("typeA"), version1, "qualifier1")
	key4 := NewAddressRefKey(1, deployment.ContractType("typeB"), version1, "qualifier1")
	key5 := NewAddressRefKey(1, deployment.ContractType("typeA"), version2, "qualifier1")
	key6 := NewAddressRefKey(1, deployment.ContractType("typeA"), version1, "qualifier2")

	assert.True(t, key1.Equals(key2), "Keys with identical fields should be equal")
	assert.False(t, key1.Equals(key3), "Keys with different chainSelector should not be equal")
	assert.False(t, key1.Equals(key4), "Keys with different contractType should not be equal")
	assert.False(t, key1.Equals(key5), "Keys with different version should not be equal")
	assert.False(t, key1.Equals(key6), "Keys with different qualifier should not be equal")
}

func TestNewAddressRefKey(t *testing.T) {
	version := semver.MustParse("1.0.0")
	key := NewAddressRefKey(1, deployment.ContractType("typeA"), version, "qualifier1")

	assert.Equal(t, uint64(1), key.ChainSelector(), "ChainSelector should match")
	assert.Equal(t, deployment.ContractType("typeA"), key.Type(), "ContractType should match")
	assert.Equal(t, version, key.Version(), "Version should match")
	assert.Equal(t, "qualifier1", key.Qualifier(), "Qualifier should match")
}
