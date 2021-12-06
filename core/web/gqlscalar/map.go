package gqlscalar

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type Map models.JSON

func (Map) ImplementsGraphQLType(name string) bool { return name == "Map" }

func (m *Map) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {
	case Map:
		*m = input
		return nil
	default:
		return errors.New("wrong type")
	}
}

func (m Map) MarshalJSON() ([]byte, error) {
	return []byte(m.Raw), nil
}

// MarshalJSON returns the JSON data if it already exists, returns
// an empty JSON object as bytes if not.
// func (j JSON) MarshalJSON() ([]byte, error) {
// 	if j.Exists() {
// 		return j.Bytes(), nil
// 	}
// 	return []byte("{}"), nil
// }
