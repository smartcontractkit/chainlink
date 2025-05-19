package resolver

import (
	"context"
	"errors"
	"fmt"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/cosmoskey"
)

func TestResolver_CosmosKeys(t *testing.T) {
	t.Parallel()

	query := `
		query GetCosmosKeys {
			cosmosKeys {
				results {
					id
				}
			}
		}`
	k := cosmoskey.MustNewInsecure(keystest.NewRandReaderFromSeed(1))
	result := fmt.Sprintf(`
	{
		"cosmosKeys": {
			"results": [
				{
					"id": "%s"
				}
			]
		}
	}`, k.PublicKeyStr())
	gError := errors.New("error")

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "cosmosKeys"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.cosmos.On("GetAll").Return([]cosmoskey.Key{k}, nil)
				f.Mocks.keystore.On("Cosmos").Return(f.Mocks.cosmos)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: result,
		},
		{
			name:          "no keys returned by GetAll",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.Mocks.cosmos.On("GetAll").Return([]cosmoskey.Key{}, gError)
				f.Mocks.keystore.On("Cosmos").Return(f.Mocks.cosmos)
				f.App.On("GetKeyStore").Return(f.Mocks.keystore)
			},
			query:  query,
			result: `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: gError,
					Path:          []interface{}{"cosmosKeys"},
					Message:       gError.Error(),
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
