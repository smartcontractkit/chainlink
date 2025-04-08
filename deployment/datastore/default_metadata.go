package datastore

// DefaultMetadata is a struct that can be used as a default metadata type.
type DefaultMetadata struct {
	Data string `json:"data"`
}

// DefaultMetadata implements the Cloneable interface
func (d DefaultMetadata) Clone() DefaultMetadata { return d }
