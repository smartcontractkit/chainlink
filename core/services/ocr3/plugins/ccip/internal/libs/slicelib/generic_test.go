package slicelib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupBy(t *testing.T) {
	type person struct {
		id   string
		name string
		age  int
	}

	testCases := []struct {
		name          string
		items         []person
		expGroupNames []string
		expGroups     map[string][]person
	}{
		{
			name:          "empty slice",
			items:         []person{},
			expGroupNames: []string{},
			expGroups:     map[string][]person{},
		},
		{
			name: "no duplicate",
			items: []person{
				{id: "2", name: "Bob", age: 25},
				{id: "1", name: "Alice", age: 23},
				{id: "3", name: "Charlie", age: 22},
				{id: "4", name: "Dim", age: 13},
			},
			expGroupNames: []string{"2", "1", "3", "4"}, // should be deterministic
			expGroups: map[string][]person{
				"1": {{id: "1", name: "Alice", age: 23}},
				"2": {{id: "2", name: "Bob", age: 25}},
				"3": {{id: "3", name: "Charlie", age: 22}},
				"4": {{id: "4", name: "Dim", age: 13}},
			},
		},
		{
			name: "with duplicate",
			items: []person{
				{id: "1", name: "Alice", age: 23},
				{id: "1", name: "Bob", age: 25},
				{id: "3", name: "Charlie", age: 22},
			},
			expGroupNames: []string{"1", "3"},
			expGroups: map[string][]person{
				"1": {{id: "1", name: "Alice", age: 23}, {id: "1", name: "Bob", age: 25}},
				"3": {{id: "3", name: "Charlie", age: 22}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			keys, groups := GroupBy(tc.items, func(p person) string { return p.id })
			assert.Equal(t, tc.expGroupNames, keys)
			assert.Equal(t, len(tc.expGroups), len(groups))
			for _, k := range keys {
				assert.Equal(t, tc.expGroups[k], groups[k])
			}
			return
		})
	}
}

func TestCountUnique(t *testing.T) {
	testCases := []struct {
		name     string
		items    []string
		expCount int
	}{
		{
			name:     "empty slice",
			items:    []string{},
			expCount: 0,
		},
		{
			name:     "no duplicate",
			items:    []string{"a", "b", "c"},
			expCount: 3,
		},
		{
			name:     "with duplicate",
			items:    []string{"a", "a", "b", "c", "b"},
			expCount: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expCount, CountUnique(tc.items))
		})
	}
}

func TestFlatten(t *testing.T) {
	testCases := []struct {
		name       string
		slices     [][]int
		expFlatten []int
	}{
		{
			name:       "empty slice",
			slices:     [][]int{},
			expFlatten: []int{},
		},
		{
			name:       "no duplicate",
			slices:     [][]int{{1, 2}, {3, 4}, {5, 6}},
			expFlatten: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:       "with duplicate",
			slices:     [][]int{{1, 2}, {1, 2}, {3, 4}, {5, 6}},
			expFlatten: []int{1, 2, 1, 2, 3, 4, 5, 6},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expFlatten, Flatten(tc.slices))
		})
	}
}

func TestFilter(t *testing.T) {
	type person struct {
		id   string
		name string
		age  int
	}

	testCases := []struct {
		name       string
		items      []person
		valid      func(person) bool
		expResults []person
	}{
		{
			name:       "empty slice",
			items:      []person{},
			valid:      func(p person) bool { return p.age > 20 },
			expResults: []person{},
		},
		{
			name: "no valid item",
			items: []person{
				{id: "1", name: "Alice", age: 18},
				{id: "2", name: "Bob", age: 20},
				{id: "3", name: "Charlie", age: 19},
			},
			valid:      func(p person) bool { return p.age > 20 },
			expResults: []person{},
		},
		{
			name: "with valid item",
			items: []person{
				{id: "1", name: "Alice", age: 18},
				{id: "2", name: "Bob", age: 25},
				{id: "3", name: "Charlie", age: 19},
			},
			valid: func(p person) bool { return p.age > 20 },
			expResults: []person{
				{id: "2", name: "Bob", age: 25},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expResults, Filter(tc.items, tc.valid))
		})
	}
}
