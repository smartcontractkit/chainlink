package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	t.Parallel()

	t.Run("multiple slices", func(t *testing.T) {
		t.Parallel()
		result := Flatten([][]int{{1, 2}, {3, 4}, {5}})
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("empty input", func(t *testing.T) {
		t.Parallel()
		result := Flatten[int](nil)
		assert.Empty(t, result)
	})

	t.Run("slices with empty sub-slices", func(t *testing.T) {
		t.Parallel()
		result := Flatten([][]int{{1}, {}, {2, 3}, {}})
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

func TestUniqueValues(t *testing.T) {
	t.Parallel()

	t.Run("with duplicates", func(t *testing.T) {
		t.Parallel()
		m := map[string]int{"a": 1, "b": 2, "c": 1, "d": 3}
		result := UniqueValues(m)
		assert.Len(t, result, 3)
		assert.ElementsMatch(t, []int{1, 2, 3}, result)
	})

	t.Run("empty map", func(t *testing.T) {
		t.Parallel()
		result := UniqueValues(map[string]int{})
		assert.Empty(t, result)
	})
}

func TestMergeMaps(t *testing.T) {
	t.Parallel()

	t.Run("overlapping keys", func(t *testing.T) {
		t.Parallel()
		a := map[string]int{"x": 1, "y": 2}
		b := map[string]int{"y": 3, "z": 4}
		result := MergeMaps(a, b)
		assert.Equal(t, map[string]int{"x": 1, "y": 3, "z": 4}, result)
	})

	t.Run("no maps", func(t *testing.T) {
		t.Parallel()
		result := MergeMaps[string, int]()
		assert.Empty(t, result)
	})

	t.Run("single map", func(t *testing.T) {
		t.Parallel()
		m := map[string]int{"a": 1}
		result := MergeMaps(m)
		assert.Equal(t, map[string]int{"a": 1}, result)
	})
}

func TestFilterSlice(t *testing.T) {
	t.Parallel()

	t.Run("filter evens", func(t *testing.T) {
		t.Parallel()
		result := FilterSlice([]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n%2 == 0 })
		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()
		result := FilterSlice([]int{}, func(n int) bool { return true })
		assert.Empty(t, result)
	})

	t.Run("none match", func(t *testing.T) {
		t.Parallel()
		result := FilterSlice([]int{1, 3, 5}, func(n int) bool { return n%2 == 0 })
		assert.Empty(t, result)
	})
}

func TestMapSlice(t *testing.T) {
	t.Parallel()

	t.Run("double values", func(t *testing.T) {
		t.Parallel()
		result := MapSlice([]int{1, 2, 3}, func(n int) int { return n * 2 })
		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("int to string", func(t *testing.T) {
		t.Parallel()
		result := MapSlice([]int{1, 2}, func(n int) string {
			return string(rune('a' + n - 1))
		})
		assert.Equal(t, []string{"a", "b"}, result)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		result := MapSlice([]int{}, func(n int) int { return n })
		assert.Empty(t, result)
	})
}

func TestContainsDuplicate(t *testing.T) {
	t.Parallel()

	assert.True(t, ContainsDuplicate([]int{1, 2, 3, 2}))
	assert.False(t, ContainsDuplicate([]int{1, 2, 3}))
	assert.False(t, ContainsDuplicate([]int{}))
	assert.False(t, ContainsDuplicate([]int{1}))
	assert.True(t, ContainsDuplicate([]string{"a", "b", "a"}))
}

func TestBatchSplit(t *testing.T) {
	list := []int{}
	for i := range 100 {
		list = append(list, i)
	}

	runs := []struct {
		name      string
		input     []int
		max       int // max per batch
		num       int // expected number of batches
		lastLen   int // expected number in last batch
		expectErr bool
	}{
		{"max=1", list, 1, len(list), 1, false},
		{"max=25", list, 25, 4, 25, false},
		{"max=33", list, 33, 4, 1, false},
		{"max=87", list, 87, 2, 13, false},
		{"max=len", list, len(list), 1, 100, false},
		{"max=len+1", list, len(list) + 1, 1, len(list), false}, // max exceeds len of list
		{"zero-list", []int{}, 1, 1, 0, false},                  // zero length list
		{"zero-max", list, 0, 0, 0, true},                       // zero as max input
		{"negative-max", list, -1, 0, 0, true},                  // negative max input
	}

	for _, r := range runs {
		t.Run(r.name, func(t *testing.T) {
			batch, err := BatchSplit(r.input, r.max)
			if r.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, batch, r.num) // check number of batches

			temp := []int{}
			for i := range batch {
				expectedLen := r.max
				if i == len(batch)-1 {
					expectedLen = r.lastLen // expect last batch to be less than max
				}
				assert.Len(t, batch[i], expectedLen) // check length of batch

				temp = append(temp, batch[i]...)
			}
			// assert order has not changed when list is reconstructed
			assert.Equal(t, r.input, temp)
		})
	}
}
