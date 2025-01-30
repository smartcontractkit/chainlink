package deployment

// MetadataSet represents a set of metadata on an address book entry.
type MetadataSet map[string]struct{}

// NewMetadataSet initializes a new MetadataSet with any number of metadata.
func NewMetadataSet(metadata ...string) MetadataSet {
	set := make(MetadataSet)
	for _, md := range metadata {
		set[md] = struct{}{}
	}
	return set
}

// Add inserts a metadata into the set.
func (ms MetadataSet) Add(metadata string) {
	ms[metadata] = struct{}{}
}

// Remove deletes a metadata from the set, if it exists.
func (ms MetadataSet) Remove(metadata string) {
	delete(ms, metadata)
}

// Contains checks if the set contains the given metadata.
func (ms MetadataSet) Contains(metadata string) bool {
	_, ok := ms[metadata]
	return ok
}

// AsSlice returns the labels in a slice. Useful for printing or serialization.
func (ms MetadataSet) AsSlice() []string {
	out := make([]string, 0, len(ms))
	for metadata := range ms {
		out = append(out, metadata)
	}
	return out
}
