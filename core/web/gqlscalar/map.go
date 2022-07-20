package gqlscalar

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Map to contain configuration
type Map map[string]interface{}

// ImplementsGraphQLType implements GraphQL type for Map
func (Map) ImplementsGraphQLType(name string) bool { return name == "Map" }

// UnmarshalGraphQL sets the Map
func (m *Map) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {
	case Map:
		*m = input
		return nil
	default:
		return errors.New("wrong type")
	}
}

// MarshalJSON returns json
func (m Map) MarshalJSON() ([]byte, error) {
	// Cast this so we don't have infinite recursion
	// (don't want json.Marshal calling the MarshalJSON method on m)
	return json.Marshal(map[string]interface{}(m))
}
