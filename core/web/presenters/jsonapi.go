package presenters

import (
	"fmt"
	"strconv"
)

// JAID represents a JSON API ID.
// It implements the api2go MarshalIdentifier and UnmarshalIdentitier interface.
type JAID struct {
	ID string `json:"-"`
}

func NewJAID(id string) JAID {
	return JAID{id}
}

// NewPrefixedJAID prefixes JAID with chain id in %s/%s format.
func NewPrefixedJAID(id string, prefix string) JAID {
	return JAID{ID: fmt.Sprintf("%s/%s", prefix, id)}
}

// NewJAIDInt32 converts an int32 into a JAID
func NewJAIDInt32(id int32) JAID {
	return JAID{strconv.Itoa(int(id))}
}

// NewJAIDInt64 converts an int64 into a JAID
func NewJAIDInt64(id int64) JAID {
	return JAID{strconv.Itoa(int(id))}
}

// GetID implements the api2go MarshalIdentifier interface.
func (jaid JAID) GetID() string {
	return jaid.ID
}

// SetID implements the api2go UnmarshalIdentitier interface.
func (jaid *JAID) SetID(value string) error {
	jaid.ID = value

	return nil
}
