package models

import (
	"database/sql/driver"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// JobID is a UUID that has a custom display format
type JobID uuid.UUID

// NilJobID is special form of JobID that is specified to have all
// 128 bits set to zero.
var NilJobID = JobID{}

// UUID converts it back into a uuid.UUID
func (id JobID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

// NewJobID returns a new JobID
func NewJobID() JobID {
	return (JobID)(uuid.NewV4())
}

// NewIDFromString is a convenience function to return an id from an input string
func NewIDFromString(input string) (JobID, error) {
	id := new(JobID)
	return *id, id.UnmarshalString(input)
}

// String satisfies the Stringer interface and removes all '-'s from the string representation of the uuid
func (id JobID) String() string {
	return strings.Replace(id.UUID().String(), "-", "", -1)
}

// MarshalText implements encoding.TextMarshaler, using String()
func (id JobID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (id *JobID) UnmarshalText(input []byte) error {
	input = utils.RemoveQuotes(input)
	return (*uuid.UUID)(id).UnmarshalText(input)
}

// UnmarshalString is a wrapper for UnmarshalText which takes a string
func (id *JobID) UnmarshalString(input string) error {
	return id.UnmarshalText([]byte(input))
}

// IsZero returns true if the JobID is the zero ID
func (id JobID) IsZero() bool {
	return id.UUID() == uuid.Nil
}

// Value hands off to the uuid lib
func (id JobID) Value() (driver.Value, error) {
	return id.UUID().Value()
}

// Scan hands off to the uuid lib
func (id *JobID) Scan(src interface{}) error {
	return (*uuid.UUID)(id).Scan(src)
}
