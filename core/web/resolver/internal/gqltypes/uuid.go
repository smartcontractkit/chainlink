package gqltypes

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type UUID struct {
	uuid.UUID
}

func (UUID) ImplementsGraphQLType(name string) bool {
	return name == "uuid"
}

func (id *UUID) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {
	case uuid.UUID:
		id.UUID = input

		return nil
	default:
		return fmt.Errorf("wrong type for UUID: %T", input)
	}
}

// MarshalJSON is a custom marshaler for UUID
//
// This function will be called whenever you
// query for fields that use the Time type
func (id UUID) MarshalJSON() ([]byte, error) {
	return id.Bytes(), nil
}
