package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetadataSet(t *testing.T) {
	t.Run("no metadata", func(t *testing.T) {
		ms := NewMetadataSet()
		assert.Empty(t, ms, "expected empty set")
	})

	t.Run("some metadata", func(t *testing.T) {
		ms := NewMetadataSet("foo", "bar")
		assert.Len(t, ms, 2)
		assert.True(t, ms.Contains("foo"))
		assert.True(t, ms.Contains("bar"))
		assert.False(t, ms.Contains("baz"))
	})
}

func TestMetadataSet_Add(t *testing.T) {
	ms := NewMetadataSet("initial")
	ms.Add("new")

	assert.True(t, ms.Contains("initial"), "expected 'initial' in set")
	assert.True(t, ms.Contains("new"), "expected 'new' in set")
	assert.Len(t, ms, 2, "expected 2 distinct metadata in set")

	// Add duplicate "new" again; size should remain 2
	ms.Add("new")
	assert.Len(t, ms, 2, "expected size to remain 2 after adding a duplicate")
}

func TestMetadataSet_Remove(t *testing.T) {
	ms := NewMetadataSet("remove_me", "keep")
	ms.Remove("remove_me")

	assert.False(t, ms.Contains("remove_me"), "expected 'remove_me' to be removed")
	assert.True(t, ms.Contains("keep"), "expected 'keep' to remain")
	assert.Len(t, ms, 1, "expected set size to be 1 after removal")

	// Removing a non-existent item shouldn't change the size
	ms.Remove("non_existent")
	assert.Len(t, ms, 1, "expected size to remain 1 after removing a non-existent item")
}

func TestMetadataSet_Contains(t *testing.T) {
	ms := NewMetadataSet("foo", "bar")

	assert.True(t, ms.Contains("foo"))
	assert.True(t, ms.Contains("bar"))
	assert.False(t, ms.Contains("baz"))
}

func TestMetadataSet_AsSlice(t *testing.T) {
	ms := NewMetadataSet("foo", "bar")
	slice := ms.AsSlice()

	// We can't rely on order in a map-based set, so we only check membership and length
	assert.Len(t, slice, 2, "expected 2 distinct metadata in slice")

	// Convert slice to a map for quick membership checks
	found := make(map[string]bool)
	for _, item := range slice {
		found[item] = true
	}
	assert.True(t, found["foo"], "expected 'foo' in slice")
	assert.True(t, found["bar"], "expected 'bar' in slice")
}
