package datastore

import (
	"sort"
	"strings"
)

// LabelSet represents a set of labels on an address book entry.
type LabelSet map[string]struct{}

// NewLabelSet initializes a new LabelSet with any number of labels.
func NewLabelSet(labels ...string) LabelSet {
	set := make(LabelSet)
	for _, lb := range labels {
		set[lb] = struct{}{}
	}
	return set
}

// Add inserts a label into the set.
func (ls LabelSet) Add(label string) {
	ls[label] = struct{}{}
}

// Remove deletes a label from the set, if it exists.
func (ls LabelSet) Remove(label string) {
	delete(ls, label)
}

// Contains checks if the set contains the given label.
func (ls LabelSet) Contains(label string) bool {
	_, ok := ls[label]
	return ok
}

// String returns the labels as a sorted, space-separated string.
// It implements the fmt.Stringer interface.
func (ls LabelSet) String() string {
	labels := ls.List()
	if len(labels) == 0 {
		return ""
	}

	// Concatenate the sorted labels into a single string
	return strings.Join(labels, " ")
}

// List returns the labels as a sorted slice of strings.
func (ls LabelSet) List() []string {
	if len(ls) == 0 {
		return []string{}
	}

	// Collect labels into a slice
	labels := make([]string, 0, len(ls))
	for label := range ls {
		labels = append(labels, label)
	}

	// Sort the labels to ensure consistent ordering
	sort.Strings(labels)

	return labels
}

// Equal checks if two LabelSets are equal.
func (ls LabelSet) Equal(other LabelSet) bool {
	if len(ls) != len(other) {
		return false
	}
	for label := range ls {
		if _, ok := other[label]; !ok {
			return false
		}
	}
	return true
}

// IsEmpty checks if the LabelSet is empty.
func (ls LabelSet) IsEmpty() bool {
	return len(ls) == 0
}
