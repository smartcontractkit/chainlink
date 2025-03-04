package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLabelSet(t *testing.T) {
	t.Run("no labels", func(t *testing.T) {
		ms := NewLabelSet()
		assert.Empty(t, ms, "expected empty set")
	})

	t.Run("some labels", func(t *testing.T) {
		ms := NewLabelSet("foo", "bar")
		assert.Len(t, ms, 2)
		assert.True(t, ms.Contains("foo"))
		assert.True(t, ms.Contains("bar"))
		assert.False(t, ms.Contains("baz"))
	})
}

func TestLabelSet_Add(t *testing.T) {
	ms := NewLabelSet("initial")
	ms.Add("new")

	assert.True(t, ms.Contains("initial"), "expected 'initial' in set")
	assert.True(t, ms.Contains("new"), "expected 'new' in set")
	assert.Len(t, ms, 2, "expected 2 distinct labels in set")

	// Add duplicate "new" again; size should remain 2
	ms.Add("new")
	assert.Len(t, ms, 2, "expected size to remain 2 after adding a duplicate")
}

func TestLabelSet_Remove(t *testing.T) {
	ms := NewLabelSet("remove_me", "keep")
	ms.Remove("remove_me")

	assert.False(t, ms.Contains("remove_me"), "expected 'remove_me' to be removed")
	assert.True(t, ms.Contains("keep"), "expected 'keep' to remain")
	assert.Len(t, ms, 1, "expected set size to be 1 after removal")

	// Removing a non-existent item shouldn't change the size
	ms.Remove("non_existent")
	assert.Len(t, ms, 1, "expected size to remain 1 after removing a non-existent item")
}

func TestLabelSet_Contains(t *testing.T) {
	ms := NewLabelSet("foo", "bar")

	assert.True(t, ms.Contains("foo"))
	assert.True(t, ms.Contains("bar"))
	assert.False(t, ms.Contains("baz"))
}

// TestLabelSet_String tests the String() method of the LabelSet type.
func TestLabelSet_String(t *testing.T) {
	tests := []struct {
		name     string
		labels   LabelSet
		expected string
	}{
		{
			name:     "Empty LabelSet",
			labels:   NewLabelSet(),
			expected: "",
		},
		{
			name:     "Single label",
			labels:   NewLabelSet("alpha"),
			expected: "alpha",
		},
		{
			name:     "Multiple labels in random order",
			labels:   NewLabelSet("beta", "gamma", "alpha"),
			expected: "alpha beta gamma",
		},
		{
			name:     "Labels with special characters",
			labels:   NewLabelSet("beta", "gamma!", "@alpha"),
			expected: "@alpha beta gamma!",
		},
		{
			name:     "Labels with spaces",
			labels:   NewLabelSet("beta", "gamma delta", "alpha"),
			expected: "alpha beta gamma delta",
		},
		{
			name:     "Labels added in different orders",
			labels:   NewLabelSet("delta", "beta", "alpha"),
			expected: "alpha beta delta",
		},
		{
			name:     "Labels with duplicate additions",
			labels:   NewLabelSet("alpha", "beta", "alpha", "gamma", "beta"),
			expected: "alpha beta gamma",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := tt.labels.String()
			assert.Equal(t, tt.expected, result, "LabelSet.String() should return the expected sorted string")
		})
	}
}

func TestLabelSet_Equal(t *testing.T) {
	tests := []struct {
		name     string
		set1     LabelSet
		set2     LabelSet
		expected bool
	}{
		{
			name:     "Both sets empty",
			set1:     NewLabelSet(),
			set2:     NewLabelSet(),
			expected: true,
		},
		{
			name:     "First set empty, second set non-empty",
			set1:     NewLabelSet(),
			set2:     NewLabelSet("label1"),
			expected: false,
		},
		{
			name:     "First set non-empty, second set empty",
			set1:     NewLabelSet("label1"),
			set2:     NewLabelSet(),
			expected: false,
		},
		{
			name:     "Identical sets with single label",
			set1:     NewLabelSet("label1"),
			set2:     NewLabelSet("label1"),
			expected: true,
		},
		{
			name:     "Identical sets with multiple labels",
			set1:     NewLabelSet("label1", "label2", "label3"),
			set2:     NewLabelSet("label3", "label2", "label1"), // Different order
			expected: true,
		},
		{
			name:     "Different sets, same size",
			set1:     NewLabelSet("label1", "label2", "label3"),
			set2:     NewLabelSet("label1", "label2", "label4"),
			expected: false,
		},
		{
			name:     "Different sets, different sizes",
			set1:     NewLabelSet("label1", "label2"),
			set2:     NewLabelSet("label1", "label2", "label3"),
			expected: false,
		},
		{
			name:     "Subset sets",
			set1:     NewLabelSet("label1", "label2"),
			set2:     NewLabelSet("label1", "label2", "label3"),
			expected: false,
		},
		{
			name:     "Disjoint sets",
			set1:     NewLabelSet("label1", "label2"),
			set2:     NewLabelSet("label3", "label4"),
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.set1.Equal(tt.set2)
			assert.Equal(t, tt.expected, result, "Equal(%v, %v) should be %v", tt.set1, tt.set2, tt.expected)
		})
	}
}
