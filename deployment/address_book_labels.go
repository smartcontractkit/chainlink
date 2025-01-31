package deployment

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

// Add inserts a labels into the set.
func (ls LabelSet) Add(labels string) {
	ls[labels] = struct{}{}
}

// Remove deletes a labels from the set, if it exists.
func (ls LabelSet) Remove(labels string) {
	delete(ls, labels)
}

// Contains checks if the set contains the given labels.
func (ls LabelSet) Contains(labels string) bool {
	_, ok := ls[labels]
	return ok
}

// AsSlice returns the labels in a slice. Useful for printing or serialization.
func (ls LabelSet) AsSlice() []string {
	out := make([]string, 0, len(ls))
	for labels := range ls {
		out = append(out, labels)
	}
	return out
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
