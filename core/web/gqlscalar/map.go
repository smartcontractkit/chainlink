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
