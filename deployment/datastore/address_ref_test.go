package datastore

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
)

func TestAddressRef_Clone(t *testing.T) {
	original := AddressRef{
		Address:       "0x123",
		ChainSelector: 1,
		Labels:        NewLabelSet("label1", "label2"),
		Qualifier:     "qualifier",
		Type:          "contractType",
		Version:       semver.MustParse("1.0.0"),
	}

	clone := original.Clone()

	assert.Equal(t, original, clone, "Clone should produce an identical copy")
	assert.NotSame(t, &original.Labels, &clone.Labels, "Labels should be deeply cloned")
}

func TestAddressRef_Key(t *testing.T) {
	ref := AddressRef{
		Address:       "0x123",
		ChainSelector: 1,
		Labels:        NewLabelSet("label1", "label2"),
		Qualifier:     "qualifier",
		Type:          "contractType",
		Version:       semver.MustParse("1.0.0"),
	}

	key := ref.Key()

	expectedKey := NewAddressRefKey(ref.ChainSelector, ref.Type, ref.Version, ref.Qualifier)
	assert.Equal(t, expectedKey, key, "Key should match the expected AddressRefKey")
}
