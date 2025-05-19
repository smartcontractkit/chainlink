package resolver

import (
	"context"
	"errors"
	"fmt"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/tronkey"
)

func TestResolver_TronKeys(t *testing.T) {
	t.Parallel()

	query := `
		query GetTronKeys {
			tronKeys {
				results {
					id
				}
			}
		}`
	k := tronkey.MustNewInsecure(keystest.NewRandReaderFromSeed(1))
	result := fmt.Sprintf(`
	{
		"tronKeys": {
			"results": [
				{
					"id": "%s"
				}
			]
		}
	}`, k.ID())
	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "tronKeys"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.tron.On("GetAll").Return([]tronkey.Key{k}, nil)
				f.Mocks.keystore.On("Tron").Return(f.Mocks.tron)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: result,
		},
		{
			name:          "no keys returned by GetAll",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.tron.On("GetAll").Return([]tronkey.Key{}, gError)
				f.Mocks.keystore.On("Tron").Return(f.Mocks.tron)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"tronKeys"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
