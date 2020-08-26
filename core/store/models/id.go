package models

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink/core/utils"

	uuid "github.com/satori/go.uuid"
)

// ID is a UUID that has a custom display format
type ID uuid.UUID

// UUID converts it back into a uuid.UUID
func (id ID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

// NewID returns a new ID
func NewID() *ID {
	uuid := uuid.NewV4()
	return (*ID)(&uuid)
}

// NewIDFromString is a convenience function to return an id from an input string
func NewIDFromString(input string) (*ID, error) {
	id := new(ID)
	return id, id.UnmarshalString(input)
}

// String satisfies the Stringer interface and removes all '-'s from the string representation of the uuid
func (id *ID) String() string {
	return strings.Replace((*uuid.UUID)(id).String(), "-", "", -1)
}

// Bytes returns the raw bytes of the underlying UUID
func (id *ID) Bytes() []byte {
	return (*uuid.UUID)(id).Bytes()
}

// MarshalText implements encoding.TextMarshaler, using String()
func (id *ID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (id *ID) UnmarshalText(input []byte) error {
	input = utils.RemoveQuotes(input)
	return (*uuid.UUID)(id).UnmarshalText(input)
}

// UnmarshalString is a wrapper for UnmarshalText which takes a string
func (id *ID) UnmarshalString(input string) error {
	return id.UnmarshalText([]byte(input))
}

// Value returns this instance serialized for database storage.
func (id *ID) Value() (driver.Value, error) {
	return id.String(), nil
}

// Scan reads the database value and returns an instance.
func (id *ID) Scan(value interface{}) error {
	switch v := value.(type) {
	case []uint8:
		return id.UnmarshalText(v)
	case string:
		return id.UnmarshalString(v)
	default:
		return fmt.Errorf("unable to convert %v of %T to ID", value, value)
	}
}
