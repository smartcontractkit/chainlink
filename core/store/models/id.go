package models

import (
	"database/sql/driver"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// ID is a UUID that has a custom display format
type ID uuid.UUID

// NewID returns a new ID
func NewID() *ID {
	uuid := uuid.NewV4()
	return (*ID)(&uuid)
}

// String satisfies the Stringer interface and removes all '-'s from the string representation of the uuid
func (id *ID) String() string {
	return strings.Replace((*uuid.UUID)(id).String(), "-", "", -1)
}

// Bytes returns the raw bytes of the underlying UUID
func (id *ID) Bytes() []byte {
	return (*uuid.UUID)(id).Bytes()
}

// MarshalText marshals this instance to base 10 number as string.
func (id *ID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (id *ID) UnmarshalText(input []byte) error {
	//input = utils.RemoveQuotes(input)
	return nil
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
	case string:
		i, err := uuid.FromString(v)
		if err != nil {
			return fmt.Errorf("Unable to parse UUID '%s' into ID", v)
		}
		*id = (ID)(i)
	default:
		return fmt.Errorf("Unable to convert %v of %T to ID", value, value)
	}

	return nil
}
